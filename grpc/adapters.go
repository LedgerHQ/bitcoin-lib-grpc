package grpc

import (
	"encoding/hex"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	pb "github.com/ledgerhq/bitcoin-lib-grpc/pb/bitcoin"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/bitcoin"
	"github.com/pkg/errors"
)

func BitcoinNetworkParams(network pb.BitcoinNetwork) (*chaincfg.Params, error) {
	switch network {
	case pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET:
		return &chaincfg.MainNetParams, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3:
		return &chaincfg.TestNet3Params, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST:
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, errors.Wrapf(ErrUnknownNetwork,
			"failed to decode network params from network %s", network.String())
	}
}

func BitcoinChainParams(chainParams *pb.ChainParams) (bitcoin.ChainParams, error) {
	switch network := chainParams.GetBitcoinNetwork(); network {
	case pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET:
		return bitcoin.Mainnet, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3:
		return bitcoin.Testnet3, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST:
		return bitcoin.Regtest, nil
	default:
		return nil, errors.Wrapf(ErrUnknownNetwork,
			"failed to decode chain params from network %s", network.String())
	}
}

func BitcoinAddressEncoding(encoding pb.AddressEncoding) (bitcoin.AddressEncoding, error) {
	switch encoding {
	case pb.AddressEncoding_ADDRESS_ENCODING_P2PKH:
		return bitcoin.Legacy, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH:
		return bitcoin.WrappedSegwit, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH:
		return bitcoin.NativeSegwit, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_UNSPECIFIED:
		return -1, errors.Wrapf(bitcoin.ErrUnknownAddressType,
			"invalid address encoding %s", encoding)
	default:
		return -1, errors.Wrapf(bitcoin.ErrUnknownAddressType,
			"invalid address encoding %s", encoding)
	}
}

// Tx is an adapter function to build a *bitcoin.Tx object from a gRPC message.
// It also converts raw gRPC values to a format that is acceptable to btcd.
func Tx(txProto *pb.CreateTransactionRequest) (*bitcoin.Tx, error) {
	var inputs []bitcoin.Input
	for _, inputProto := range txProto.Inputs {
		inputs = append(inputs, bitcoin.Input{
			OutputHash:  inputProto.OutputHash,
			OutputIndex: uint32(inputProto.OutputIndex),
		})
	}

	var outputs []bitcoin.Output
	for _, outputProto := range txProto.Outputs {
		value, err := strconv.ParseInt(outputProto.Value, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err,
				"invalid output value: %s", outputProto.Value)
		}

		outputs = append(outputs, bitcoin.Output{
			Address: outputProto.Address,
			Value:   value,
		})
	}

	return &bitcoin.Tx{
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: txProto.LockTime,
	}, nil
}

// RawTx is an adapter function to build a *bitcoin.RawTx object from a gRPC message.
func RawTx(rawTxProto *pb.RawTransactionResponse) *bitcoin.RawTx {
	return &bitcoin.RawTx{
		Hex:         rawTxProto.Hex,
		Hash:        rawTxProto.Hash,
		WitnessHash: rawTxProto.WitnessHash,
	}
}

// Utxo is an adapter function to build a *bitcoin.Utxo object from a gRPC message.
func Utxo(proto *pb.Utxo) (*bitcoin.Utxo, error) {
	value, err := strconv.ParseInt(proto.Value, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid utxo value: %s", proto.Value)
	}

	return &bitcoin.Utxo{
		Script:     proto.Script,
		Value:      value,
		Derivation: proto.Derivation,
	}, nil
}

// SignatureMetadata is an adapter function to build a *bitcoin.SignatureMetadata object from a gRPC message.
func SignatureMetadata(proto *pb.SignatureMetadata, chainParams bitcoin.ChainParams) (*bitcoin.SignatureMetadata, error) {
	addrEncoding, err := BitcoinAddressEncoding(proto.AddrEncoding)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid output value: %s", proto.AddrEncoding)
	}

	serializedPubKey, err := hex.DecodeString(proto.PublicKey)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to parse serialized pub key from %s", proto.PublicKey)
	}

	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, chainParams)
	if err != nil {
		return nil, errors.Wrap(err,
			"failed to parse pub key from signature")
	}

	return &bitcoin.SignatureMetadata{
		DerSig:       proto.DerSignature,
		PubKey:       addressPubKey.PubKey(),
		AddrEncoding: addrEncoding,
	}, nil
}

package grpc

import (
	"encoding/hex"
	"strconv"

	"github.com/btcsuite/btcutil"
	pb "github.com/ledgerhq/bitcoin-lib-grpc/pb/bitcoin"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/chaincfg"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/core"
	"github.com/pkg/errors"
)

func ChainParams(chainParams *pb.ChainParams) (chaincfg.ChainParams, error) {
	switch network := chainParams.GetBitcoinNetwork(); network {
	case pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET:
		return chaincfg.BitcoinMainNetParams, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3:
		return chaincfg.BitcoinTestNet3Params, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST:
		return chaincfg.BitcoinRegressionNetParams, nil
	}

	switch network := chainParams.GetLitecoinNetwork(); network {
	case pb.LitecoinNetwork_LITECOIN_NETWORK_MAINNET:
		return chaincfg.LitecoinMainNetParams, nil
	}

	switch network := chainParams.GetBitcoinCashNetwork(); network {
	case pb.BitcoinCashNetwork_BITCOIN_CASH_NETWORK_MAINNET:
		return chaincfg.BitcoinCashMainNetParams, nil
	default:
		return nil, errors.Wrapf(ErrUnknownNetwork,
			"failed to decode chain params from network %s", network.String())
	}
}

func BitcoinAddressEncoding(encoding pb.AddressEncoding) (core.AddressEncoding, error) {
	switch encoding {
	case pb.AddressEncoding_ADDRESS_ENCODING_P2PKH:
		return core.Legacy, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH:
		return core.WrappedSegwit, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH:
		return core.NativeSegwit, nil
	case pb.AddressEncoding_ADDRESS_ENCODING_UNSPECIFIED:
		return -1, errors.Wrapf(core.ErrUnknownAddressType,
			"invalid address encoding %s", encoding)
	default:
		return -1, errors.Wrapf(core.ErrUnknownAddressType,
			"invalid address encoding %s", encoding)
	}
}

// Tx is an adapter function to build a *core.Tx object from a gRPC message.
// It also converts raw gRPC values to a format that is acceptable to btcd.
func Tx(txProto *pb.CreateTransactionRequest) (*core.Tx, error) {
	var inputs []core.Input
	for _, inputProto := range txProto.Inputs {
		inputs = append(inputs, core.Input{
			OutputHash:  inputProto.OutputHash,
			OutputIndex: uint32(inputProto.OutputIndex),
			Script:      inputProto.Script,
			Value:       inputProto.Value,
		})
	}

	var outputs []core.Output
	for _, outputProto := range txProto.Outputs {
		value, err := strconv.ParseInt(outputProto.Value, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err,
				"invalid output value: %s", outputProto.Value)
		}

		outputs = append(outputs, core.Output{
			Address: outputProto.Address,
			Value:   value,
		})
	}

	return &core.Tx{
		Inputs:        inputs,
		Outputs:       outputs,
		ChangeAddress: txProto.ChangeAddress,
		FeeSatPerKb:   txProto.FeeSatPerKb,
		LockTime:      txProto.LockTime,
	}, nil
}

// RawTx is an adapter function to build a *core.RawTx object from a gRPC message.
func RawTx(rawTxProto *pb.RawTransactionResponse) *core.RawTx {
	return &core.RawTx{
		Hex:         rawTxProto.Hex,
		Hash:        rawTxProto.Hash,
		WitnessHash: rawTxProto.WitnessHash,
	}
}

// Utxo is an adapter function to build a *bitcoin.Utxo object from a gRPC message.
func Utxo(proto *pb.Utxo) (*core.Utxo, error) {
	value, err := strconv.ParseInt(proto.Value, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid utxo value: %s", proto.Value)
	}

	script, err := hex.DecodeString(proto.ScriptHex)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid utxo script hex: %s", proto.ScriptHex)
	}

	return &core.Utxo{
		Script:     script,
		Value:      value,
		Derivation: proto.Derivation,
	}, nil
}

// SignatureMetadata is an adapter function to build a *bitcoin.SignatureMetadata object from a gRPC message.
func SignatureMetadata(proto *pb.SignatureMetadata, chainParams chaincfg.ChainParams) (*core.SignatureMetadata, error) {
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

	return &core.SignatureMetadata{
		DerSig:       proto.DerSignature,
		PubKey:       addressPubKey.PubKey(),
		AddrEncoding: addrEncoding,
	}, nil
}

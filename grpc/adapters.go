package grpc

import (
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
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
			"failed to decode network %s", network.String())
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
			"failed to decode network %s", network.String())
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
	var inputs []bitcoin.UnsignedInput
	for _, inputProto := range txProto.Inputs {
		inputs = append(inputs, bitcoin.UnsignedInput{
			OutputHash:  inputProto.OutputHash,
			OutputIndex: uint32(inputProto.OutputIndex),
			Script:      inputProto.Script,
			Sequence:    inputProto.Sequence,
		})
	}

	var recipients []bitcoin.Recipient
	for _, recipientProto := range txProto.Recipients {
		value, err := strconv.ParseInt(recipientProto.Value, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err,
				"invalid recipient value: %s", recipientProto.Value)
		}

		recipients = append(recipients, bitcoin.Recipient{
			Address: recipientProto.Address,
			Value:   value,
		})
	}

	return &bitcoin.Tx{
		Inputs:     inputs,
		Recipients: recipients,
		LockTime:   txProto.LockTime,
	}, nil
}

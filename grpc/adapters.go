package grpc

import (
	"github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"github.com/ledgerhq/lama-bitcoin-svc/pkg/bitcoin"
	"github.com/pkg/errors"
)

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

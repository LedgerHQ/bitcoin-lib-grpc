package bitcoin

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct{}

func (s *Service) ValidateAddress(
	ctx context.Context, request *pb.ValidateAddressRequest,
) (*pb.ValidateAddressResponse, error) {
	params, err := getChainParams(request.ChainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	addr, err := btcutil.DecodeAddress(request.Address, params)
	if err != nil {
		return &pb.ValidateAddressResponse{
			Address:       request.Address,
			IsValid:       false,
			InvalidReason: err.Error(),
		}, nil
	}

	return &pb.ValidateAddressResponse{
		Address: addr.EncodeAddress(), // Normalize the original address
		IsValid: true,
	}, nil
}

func getChainParams(params *pb.ChainParams) (*chaincfg.Params, error) {
	switch network := params.GetBitcoinNetwork(); network {
	case pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET:
		return &chaincfg.MainNetParams, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3:
		return &chaincfg.TestNet3Params, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST:
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, ErrUnknownNetwork(network.String())
	}
}

package grpc

import (
	"context"

	pb "github.com/ledgerhq/bitcoin-lib-grpc/pb/bitcoin"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/bitcoin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controller struct {
	svc bitcoin.Service
}

func NewBitcoinController() *controller {
	return &controller{
		svc: bitcoin.Service{},
	}
}

func (c *controller) ValidateAddress(
	ctx context.Context, request *pb.ValidateAddressRequest,
) (*pb.ValidateAddressResponse, error) {
	chainParams, err := BitcoinChainParams(request.ChainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	addr, err := c.svc.ValidateAddress(request.Address, chainParams)
	if err != nil {
		return &pb.ValidateAddressResponse{
			Address:       request.Address,
			IsValid:       false,
			InvalidReason: err.Error(),
		}, nil
	}

	return &pb.ValidateAddressResponse{
		Address: addr,
		IsValid: true,
	}, nil
}

func (c *controller) EncodeAddress(
	ctx context.Context, request *pb.EncodeAddressRequest,
) (*pb.EncodeAddressResponse, error) {
	chainParams, err := BitcoinChainParams(request.ChainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	encoding, err := BitcoinAddressEncoding(request.Encoding)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	address, err := c.svc.EncodeAddress(request.PublicKey, encoding, chainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.EncodeAddressResponse{
		Address: address,
	}, nil
}

func (c *controller) DeriveExtendedKey(
	ctx context.Context, request *pb.DeriveExtendedKeyRequest,
) (*pb.DeriveExtendedKeyResponse, error) {
	response, err := c.svc.DeriveExtendedKey(request.ExtendedKey, request.Derivation)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.DeriveExtendedKeyResponse{
		ExtendedKey: response.ExtendedKey,
		PublicKey:   response.PublicKey,
		ChainCode:   response.ChainCode,
	}, nil
}

func (c *controller) CreateTransaction(
	ctx context.Context, txRequest *pb.CreateTransactionRequest,
) (*pb.CreateTransactionResponse, error) {

	chainParams, err := BitcoinNetworkParams(txRequest.Network)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	tx, err := Tx(txRequest)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	rawTx, err := c.svc.CreateTransaction(tx, chainParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := pb.CreateTransactionResponse{
		Hex:         rawTx.Hex,
		Hash:        rawTx.Hash,
		WitnessHash: rawTx.WitnessHash,
	}

	return &response, nil
}

func (c *controller) GetKeypair(
	ctx context.Context, request *pb.GetKeypairRequest,
) (*pb.GetKeypairResponse, error) {

	chainParams, err := BitcoinNetworkParams(request.Network)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	keypair, err := c.svc.GetKeypair(request.Seed, chainParams, request.Derivation)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := pb.GetKeypairResponse{ExtendedPublicKey: keypair.ExtendedPublicKey, PrivateKey: keypair.PrivateKey}

	return &response, nil
}

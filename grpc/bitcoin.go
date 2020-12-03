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
) (*pb.RawTransactionResponse, error) {

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

	var notEnoughtUtxo *pb.NotEnoughUtxo
	if rawTx.NotEnoughUtxo != nil {
		notEnoughtUtxo = &pb.NotEnoughUtxo{MissingAmount: rawTx.NotEnoughUtxo.MissingAmount}
	}

	response := pb.RawTransactionResponse{
		Hex:           rawTx.Hex,
		Hash:          rawTx.Hash,
		WitnessHash:   rawTx.WitnessHash,
		NotEnoughUtxo: notEnoughtUtxo,
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

func (c *controller) GenerateDerSignatures(
	ctx context.Context, request *pb.GenerateDerSignaturesRequest,
) (*pb.GenerateDerSignaturesResponse, error) {

	rawTx := RawTx(request.RawTx)

	utxos := make([]bitcoin.Utxo, len(request.Utxos))
	for idx, utxoProto := range request.Utxos {
		utxo, err := Utxo(utxoProto)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		utxos[idx] = *utxo
	}

	msgTx, err := c.svc.DeserializeMsgTx(rawTx)
	if err != nil {
		return nil, err
	}

	derSignatures, err := c.svc.GenerateDerSignatures(msgTx, utxos, request.PrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.GenerateDerSignaturesResponse{DerSignatures: derSignatures}, nil
}

func (c *controller) SignTransaction(
	ctx context.Context, request *pb.SignTransactionRequest,
) (*pb.RawTransactionResponse, error) {

	chainParams, err := BitcoinNetworkParams(request.Network)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rawTx := RawTx(request.RawTx)

	signatures := make([]bitcoin.SignatureMetadata, len(request.Signatures))

	for idx, signature := range request.Signatures {
		sigMetadata, err := SignatureMetadata(signature, chainParams)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		signatures[idx] = *sigMetadata
	}

	msgTx, err := c.svc.DeserializeMsgTx(rawTx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	signedRawTx, err := c.svc.SignTransaction(msgTx, chainParams, signatures)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := pb.RawTransactionResponse{
		Hex:         signedRawTx.Hex,
		Hash:        signedRawTx.Hash,
		WitnessHash: signedRawTx.WitnessHash,
	}

	return &response, nil
}

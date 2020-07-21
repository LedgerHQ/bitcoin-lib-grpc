package bitcoin

import (
	"context"

	"github.com/btcsuite/btcutil/hdkeychain"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

// DeriveExtendedKey provides an API to derive hierarchical deterministic
// extended keys.
//
// There are no restrictions on the extended keys that can be derived, as
// long as BIP0032 rules are followed. However, it is intended to be used
// for deriving child keys from public extended keys at the account level
// (HD depth 3).
//
// The derivation is agnostic of chain parameters. Derived extended keys
// are associated to the same network as the parent extended key.
//
// The method's response includes the following fields:
//     ExtendedKey: extended key as a human-readable base58-encoded string.
//     PublicKey:   33-byte compressed public key of the derived extended key.
//     ChainCode:   32-byte chain code of the derived extended key.
func (s *Service) DeriveExtendedKey(
	ctx context.Context, request *pb.DeriveExtendedKeyRequest,
) (*pb.DeriveExtendedKeyResponse, error) {
	extendedKey, err := hdkeychain.NewKeyFromString(request.ExtendedKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Derive len(request.Derivation) HD levels, starting from extendedKey
	// as the parent node.
	for _, childIndex := range request.Derivation {
		extendedKey, err = extendedKey.Child(childIndex)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	pubKey, err := extendedKey.ECPubKey()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.DeriveExtendedKeyResponse{
		ExtendedKey: extendedKey.String(),
		PublicKey:   pubKey.SerializeCompressed(),
		ChainCode:   extendedKey.ChainCode(),
	}, nil
}

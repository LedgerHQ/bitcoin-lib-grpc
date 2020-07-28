package bitcoin

import (
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

// PublicKeyMaterial contains an extended public key, and the corresponding
// public key - chain code pair.
type PublicKeyMaterial struct {
	ExtendedKey string
	PublicKey   []byte
	ChainCode   []byte
}

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
	extendedKey string, derivation []uint32,
) (PublicKeyMaterial, error) {
	response := PublicKeyMaterial{}

	xKey, err := hdkeychain.NewKeyFromString(extendedKey)
	if err != nil {
		return response, errors.Wrapf(err, "failed to decode xkey %s",
			extendedKey)
	}

	// Derive len(request.Derivation) HD levels, starting from extendedKey
	// as the parent node.
	for _, childIndex := range derivation {
		xKey, err = xKey.Child(childIndex)
		if err != nil {
			return response, errors.Wrapf(err, "failed to derive xkey %s at index %d",
				extendedKey, childIndex)
		}
	}

	pubKey, err := xKey.ECPubKey()
	if err != nil {
		return response, errors.Wrapf(err, "failed to get public key from xkey %s",
			extendedKey)
	}

	response.ExtendedKey = xKey.String()
	response.PublicKey = pubKey.SerializeCompressed()
	response.ChainCode = xKey.ChainCode()
	return response, nil
}

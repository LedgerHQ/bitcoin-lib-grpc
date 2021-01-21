package bitcoin

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
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

// Keypair contains en extended public key and the corresponding private key
type Keypair struct {
	ExtendedPublicKey string
	PrivateKey        string
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
		xKey, err = xKey.Derive(childIndex)
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

// GetAccountExtendedKey returns the serialized extended key from public key
// material, and various parameters. This is typically provided by the HSM.
//
// Certain assumptions have been made for the parent fingerprint. Please read
// the corresponding note in the code.
//
// accountIndex must NOT add the BIP32 harden bit. The account MUST have
// been derived using the following scheme:
//   m / purpose' / coin_type' / account'
//
// It also implies that accountIndex is at BIP32 level 3.
func (s *Service) GetAccountExtendedKey(
	publicKey []byte,
	chainCode []byte,
	accountIndex uint32,
	chainParams ChainParams,
) (string, error) {
	// Load the serialized public key to a btcec.PublicKey type, in order to
	// ensure that the:
	//   * public point is on the secp256k1 elliptic curve.
	//   * public point coordinates belong to the finite field of secp256k1.
	//   * public key is well formed (valid magic, length, etc).
	//   * public key used for serializing the extended key is compressed.
	//
	// Both compressed and uncompressed public keys are accepted.
	loadedPublicKey, err := btcec.ParsePubKey(publicKey, btcec.S256())
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse public key %s",
			hex.EncodeToString(publicKey))
	}

	serializedPublicKey := loadedPublicKey.SerializeCompressed()
	const depth = 3
	childNum := accountIndex + hdkeychain.HardenedKeyStart

	// The fingerprint of the parent for the derived child is the first 4
	// bytes of the RIPEMD160(SHA256(parentPubKey)).
	//
	// Caution: The HSM does NOT provide the parent fingerprint, so we use
	// the fingerprint of the child (BIP32 depth 3). While this is incorrect,
	// the fingerprint has no impact on the derived addresses.
	parentFP := btcutil.Hash160(serializedPublicKey)[:4]

	key := hdkeychain.NewExtendedKey(
		chainParams.HDPublicKeyID[:],
		serializedPublicKey,
		chainCode,
		parentFP,
		depth,
		childNum,
		false,
	)

	return key.String(), nil
}

// Useful service to get keypair (xpub + privKey) from a seed for testing.
// Random seed is generated if no seed is provided.
func (s *Service) GetKeypair(seed string, chainParams ChainParams, derivation []uint32) (Keypair, error) {
	var (
		seedBytes []byte
		response  Keypair
	)

	if seed == "" {
		// Generate a random seed at the recommended length.
		generatedSeed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			return response, err
		}

		seedBytes = generatedSeed

	} else {
		seedBytes = []byte(seed)
	}

	// Generate a new master node using the seed.
	extendedKey, err := hdkeychain.NewMaster(seedBytes, chainParams)
	if err != nil {
		return response, err
	}

	// Derive the extended key for given derivation path
	for _, childIndex := range derivation {
		extendedKey, err = extendedKey.Derive(childIndex)
		if err != nil {
			return response, errors.Wrapf(err, "failed to derive extendedKey %s at index %d",
				extendedKey, childIndex)
		}
	}

	// Get the human readable extended public key
	accountExtendedPublicKey, err := extendedKey.Neuter()
	if err != nil {
		return response, err
	}

	response.ExtendedPublicKey = accountExtendedPublicKey.String()
	response.PrivateKey = extendedKey.String()

	return response, nil
}

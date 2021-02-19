package bitcoin

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"

	"github.com/btcsuite/btcutil"
)

// AddressEncoding is an enum type for the various address encoding
// schemes supported for Bitcoin.
type AddressEncoding int

const (
	// Legacy indicates the P2PKH address encoding scheme.
	Legacy AddressEncoding = iota

	// WrappedSegwit indicates the P2WPKH-in-P2SH address encoding
	// scheme.
	WrappedSegwit

	// NativeSegwit indicates the P2WPKH address encoding scheme.
	NativeSegwit
)

// ValidateAddress returns an error if the given address is malformed.
// It returns the normalized address otherwise.
func (s *Service) ValidateAddress(address string, chainParams ChainParams) (string, error) {
	addr, err := btcutil.DecodeAddress(address, chainParams)
	if err != nil {
		return "", errors.Wrapf(err, "failed to decode address %s", address)
	}

	// Normalize the original address
	return addr.EncodeAddress(), nil
}

// EncodeAddress serializes a public key into a string, based on the
// encoding and the chain parameters.
//
// References:
//   [Learn me a Bitcoin]: P2PKH - Pay To Pubkey Hash
//   https://learnmeabitcoin.com/technical/p2pkh
//
//   [BIP13]: BIP0013 - Address Format for pay-to-script-hash
//   https://github.com/bitcoin/bips/blob/master/bip-0013.mediawiki
//
//   [BIP173]: BIP0173 - Base32 address format for native v0-16 witness outputs
//   https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
func (s *Service) EncodeAddress(
	publicKey []byte, encoding AddressEncoding, chainParams ChainParams,
) (string, error) {
	// Load the serialized public key to a btcec.PublicKey type, in order to
	// ensure that the:
	//   * public point is on the secp256k1 elliptic curve.
	//   * public point coordinates belong to the finite field of secp256k1.
	//   * public key is well formed (valid magic, length, etc).
	//
	// Both compressed and uncompressed public keys are accepted.
	//
	// Using addresses encoded from incorrect public keys may lead to
	// irrevocable fund loss.
	loadedPublicKey, err := btcec.ParsePubKey(publicKey, btcec.S256())
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse public key %s",
			hex.EncodeToString(publicKey))
	}

	// Calculate the RIPEMD160 of the SHA256 of a public key, aka HASH160.
	//
	// As per Bitcoin protocol, the serialized public key MUST be compressed
	// for P2SH-P2WPKH and P2WPKH addresses, whereas P2PKH addresses could be
	// either compressed or uncompressed. Regarding P2PKH addresses, the
	// convention at Ledger and in Bitcoin Wiki examples is to use compressed
	// public keys.
	publicKeyHash := btcutil.Hash160(loadedPublicKey.SerializeCompressed())

	address, err := func() (btcutil.Address, error) {
		switch encoding {
		case Legacy:
			// Ref: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			return btcutil.NewAddressPubKeyHash(publicKeyHash, chainParams)
		case WrappedSegwit:
			// Ref: https://bitcoincore.org/en/segwit_wallet_dev/#creation-of-p2sh-p2wpkh-address

			// Create a P2WPKH native-segwit address
			p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, chainParams)
			if err != nil {
				return nil, err
			}

			// Create a P2SH redeemScript that pays to the P2WPKH address.
			//
			// The redeemScript is 22 bytes long, and starts with a OP_0,
			// followed by a canonical push of the keyhash. The keyhash
			// is HASH160 of the 33-byte compressed public key.
			//
			// scriptSig: OP_0 <hash160(compressed public key)>
			redeemScript, err := txscript.PayToAddrScript(p2wpkhAddress)
			if err != nil {
				return nil, err
			}

			return btcutil.NewAddressScriptHash(redeemScript, chainParams)
		case NativeSegwit:
			// Ref: https://bitcoincore.org/en/segwit_wallet_dev/#native-pay-to-witness-public-key-hash-p2wpkh
			return btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, chainParams)
		default:
			return nil, ErrUnknownAddressType
		}
	}()

	if err != nil {
		return "", errors.Wrapf(err, "unable to encode public key %s to address",
			hex.EncodeToString(publicKey))
	}

	return address.EncodeAddress(), nil
}

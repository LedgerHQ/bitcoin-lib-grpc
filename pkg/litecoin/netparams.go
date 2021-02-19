package litecoin

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// MainNetParams defines the network parameters for the main Litecoin network.
// For reference, see: https://github.com/ltcsuite/ltcd/blob/master/chaincfg/params.go#L229
var MainNetParams *chaincfg.Params

func init() {
	// Copy of Btc main net params to construct LTC MainNetParams
	fromBtcParams := chaincfg.MainNetParams

	MainNetParams = &fromBtcParams

	// Magic number
	MainNetParams.Net = 0xdbb6c0fb

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	MainNetParams.Bech32HRPSegwit = "ltc" // always ltc for main net

	// Address encoding magics
	MainNetParams.PubKeyHashAddrID = 0x30        // starts with L
	MainNetParams.ScriptHashAddrID = 0x32        // starts with M
	MainNetParams.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	MainNetParams.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	MainNetParams.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	MainNetParams.HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4} // starts with Ltpv
	MainNetParams.HDPublicKeyID = [4]byte{0x04, 0x88, 0xb2, 0x1e}  // starts with Ltub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	MainNetParams.HDCoinType = 2

	// Register litecoin network params to the changcfg
	chaincfg.Register(MainNetParams)
}

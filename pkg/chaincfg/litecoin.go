package chaincfg

import "github.com/btcsuite/btcd/chaincfg"

// LitecoinMainNetParams defines the network parameters for the main Litecoin network.
// For reference, see: https://github.com/ltcsuite/ltcd/blob/master/chaincfg/params.go#L229
var LitecoinMainNetParams *chaincfg.Params

func init() {
	// Copy of Btc main net params to construct LTC LitecoinMainNetParams
	fromBtcParams := chaincfg.MainNetParams

	LitecoinMainNetParams = &fromBtcParams

	// Magic number
	LitecoinMainNetParams.Net = 0xdbb6c0fb

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	LitecoinMainNetParams.Bech32HRPSegwit = "ltc" // always ltc for main net

	// Address encoding magics
	LitecoinMainNetParams.PubKeyHashAddrID = 0x30        // starts with L
	LitecoinMainNetParams.ScriptHashAddrID = 0x32        // starts with M
	LitecoinMainNetParams.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	LitecoinMainNetParams.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	LitecoinMainNetParams.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	LitecoinMainNetParams.HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4} // starts with Ltpv
	LitecoinMainNetParams.HDPublicKeyID = [4]byte{0x04, 0x88, 0xb2, 0x1e}  // starts with Ltub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	LitecoinMainNetParams.HDCoinType = 2

	// Register litecoin network params to the changcfg
	if err := chaincfg.Register(LitecoinMainNetParams); err != nil {
		panic(err)
	}
}

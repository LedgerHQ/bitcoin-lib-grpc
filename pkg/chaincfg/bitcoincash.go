package chaincfg

import "github.com/btcsuite/btcd/chaincfg"

// BitcoinCashMainNetParams defines the network parameters for the main Bitcoin cash network.
// For reference, see: https://github.com/gcash/bchd/blob/master/chaincfg/params.go#L261
var BitcoinCashMainNetParams *chaincfg.Params

func init() {
	// Copy of Btc main net params to construct BCH BitcoinCashMainNetParams
	fromBtcParams := chaincfg.MainNetParams

	BitcoinCashMainNetParams = &fromBtcParams

	// Magic number
	BitcoinCashMainNetParams.Net = 0xe8f3e1e3

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	BitcoinCashMainNetParams.Bech32HRPSegwit = "bitcoincash" // always bch for main net

	// Address encoding magics
	BitcoinCashMainNetParams.PubKeyHashAddrID = 0x00 // starts with 1
	BitcoinCashMainNetParams.ScriptHashAddrID = 0x05 // starts with 3
	BitcoinCashMainNetParams.PrivateKeyID = 0x80     // starts with 5 (uncompressed) or K (compressed)

	// BIP32 hierarchical deterministic extended key magics
	BitcoinCashMainNetParams.HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4} // starts with xprv
	BitcoinCashMainNetParams.HDPublicKeyID = [4]byte{0x04, 0x88, 0xb2, 0x1e}  // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	BitcoinCashMainNetParams.HDCoinType = 145

	// Register bitcoin cash network params to the changcfg
	if err := chaincfg.Register(BitcoinCashMainNetParams); err != nil {
		panic(err)
	}
}

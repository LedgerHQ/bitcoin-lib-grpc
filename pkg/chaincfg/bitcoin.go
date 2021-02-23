package chaincfg

import "github.com/btcsuite/btcd/chaincfg"

// Bitcoin network params
var (
	// BitcoinMainNetParams defines the network parameters for the main Bitcoin network.
	BitcoinMainNetParams = &chaincfg.MainNetParams

	// BitcoinTestNet3Params defines the network parameters for the test Bitcoin network
	// (version 3).
	BitcoinTestNet3Params = &chaincfg.TestNet3Params

	// BitcoinRegressionNetParams defines the network parameters for the regression test
	// Bitcoin network.
	BitcoinRegressionNetParams = &chaincfg.RegressionNetParams
)

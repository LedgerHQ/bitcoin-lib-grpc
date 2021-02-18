package bitcoin

import "github.com/btcsuite/btcd/chaincfg"

// ChainParams is a type alias for chaincfg.Params, to allow external
// packages to refer to the chain parameters without importing btcd.
type ChainParams = *chaincfg.Params

// Bitcoin network params
var (
	// MainNetParams defines the network parameters for the main Bitcoin network.
	MainNetParams = &chaincfg.MainNetParams

	// TestNet3Params defines the network parameters for the test Bitcoin network
	// (version 3).
	TestNet3Params = &chaincfg.TestNet3Params

	// RegTestParams defines the network parameters for the regression test
	// Bitcoin network.
	RegTestParams = &chaincfg.RegressionNetParams
)

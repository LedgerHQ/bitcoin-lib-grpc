package chaincfg

import "github.com/btcsuite/btcd/chaincfg"

// ChainParams is a type alias for chaincfg.Params, to allow external
// packages to refer to the chain parameters without importing btcd.
type ChainParams = *chaincfg.Params

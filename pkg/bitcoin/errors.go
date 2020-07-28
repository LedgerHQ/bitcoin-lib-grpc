package bitcoin

import "github.com/btcsuite/btcutil"

// ErrUnknownAddressType is a type alias to allow reference in external
// packages without importing btcutil.
var ErrUnknownAddressType = btcutil.ErrUnknownAddressType

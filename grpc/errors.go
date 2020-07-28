package grpc

import "errors"

// ErrUnknownNetwork is returned when a string representing an unknown network
// is found.
var ErrUnknownNetwork = errors.New("invalid network")

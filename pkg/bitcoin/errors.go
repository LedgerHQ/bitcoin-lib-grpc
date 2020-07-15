package bitcoin

import "fmt"

// ErrUnknownNetwork is returned when a string representing an unknown network
// is found.
type ErrUnknownNetwork string

func (e ErrUnknownNetwork) Error() string {
	return fmt.Sprintf("invalid network in chain params: %s", string(e))
}

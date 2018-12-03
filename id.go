package errorx

import (
	"sync"
)

var internalID uint64 = 0x123456789abcdef0
var idMtx sync.Mutex

// nextInternalID creates next unique id for errorx entities.
// All equality comparison should take id into account, lest there be some false positive matches.
func nextInternalID() uint64 {
	idMtx.Lock()
	// xorshift
	x := internalID
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	id := x + internalID
	internalID = x
	idMtx.Unlock()
	return id
}

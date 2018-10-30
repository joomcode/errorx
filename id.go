package errorx

import "sync/atomic"

var internalID int64

// nextInternalID creates next unique id for errorx entities.
// All equality comparison should take id into account, lest there be some false positive matches.
func nextInternalID() int64 {
	return atomic.AddInt64(&internalID, 1)
}

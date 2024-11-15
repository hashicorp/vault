package types

import (
	"math"
	"time"
)

const (
	// CITRUSLEAF_EPOCH defines the citrusleaf epoc: Jan 01 2010 00:00:00 GMT
	CITRUSLEAF_EPOCH = 1262304000
)

// TTL converts an Expiration time from citrusleaf epoc to TTL in seconds.
func TTL(secsFromCitrusLeafEpoc uint32) uint32 {
	switch secsFromCitrusLeafEpoc {
	// don't convert magic values
	case 0: // when set to don't expire, this value is returned
		return math.MaxUint32
	default:
		// Record may not have expired on server, but delay or clock differences may
		// cause it to look expired on client. Floor at 1, not 0, to avoid old
		// "never expires" interpretation.
		now := time.Now().Unix()
		expiration := int64(CITRUSLEAF_EPOCH + secsFromCitrusLeafEpoc)
		if (expiration < 0 && now >= 0) || expiration > now {
			return uint32(expiration - now)
		}
		return 1
	}
}

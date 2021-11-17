package memd

import (
	"math"
	"time"
)

// EncodeSrvDura16 takes a standard go time duration and encodes it into
// the appropriate format for the server.
func EncodeSrvDura16(dura time.Duration) uint16 {
	serverDurationUs := dura / time.Microsecond
	serverDurationEnc := int(math.Pow(float64(serverDurationUs)*2, 1.0/1.74))
	if serverDurationEnc > 65535 {
		serverDurationEnc = 65535
	}
	return uint16(serverDurationEnc)
}

// DecodeSrvDura16 takes an encoded operation duration from the server
// and converts it to a standard Go time duration.
func DecodeSrvDura16(enc uint16) time.Duration {
	serverDurationUs := math.Round(math.Pow(float64(enc), 1.74) / 2)
	return time.Duration(serverDurationUs) * time.Microsecond
}

package canoe

import (
	"encoding/binary"
	"github.com/satori/go.uuid"
)

// Uint64UUID returns a UUID encoded to uint64
func Uint64UUID() uint64 {
	return binary.LittleEndian.Uint64(uuid.NewV4().Bytes())
}

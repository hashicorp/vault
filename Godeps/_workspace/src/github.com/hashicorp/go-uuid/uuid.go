package uuid

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUID is used to generate a random UUID
func GenerateUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}

	return FormatUUID(buf)
}

func FormatUUID(buf []byte) (string, error) {
	if len(buf) != 16 {
		return "", fmt.Errorf("wrong length byte slice (%d)", len(buf))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16]), nil
}

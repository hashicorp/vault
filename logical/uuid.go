package logical

import (
	"crypto/rand"
	"fmt"
	"time"
)

// UUID returns a UUID.
func UUID() (string, error) {
	unix := uint32(time.Now().UTC().Unix())

	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%04x%08x",
			unix, b[0:2], b[2:4], b[4:6], b[6:8], b[8:]),
		nil
}

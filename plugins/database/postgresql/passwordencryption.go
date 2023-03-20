package postgresql

import "fmt"

type passwordEncryption string

var (
	passwordEncryptionSCRAMSHA256 passwordEncryption = "scram-sha-256"
	passwordEncryptionNone        passwordEncryption = "none"
)

var passwordEncryptions = map[passwordEncryption]struct{}{
	passwordEncryptionSCRAMSHA256: {},
	passwordEncryptionNone:        {},
}

func parsePasswordEncryption(s string) (passwordEncryption, error) {
	if _, ok := passwordEncryptions[passwordEncryption(s)]; !ok {
		return "", fmt.Errorf("'%s' is not a valid password encryption type", s)
	}

	return passwordEncryption(s), nil
}

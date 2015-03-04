// password is a package for reading a password securely from a terminal.
// The code in this package disables echo in the terminal so that the
// password is not echoed back in plaintext to the user.
package password

import (
	"io"
	"os"
)

// Read reads the password from the given os.File. The password
// will not be echoed back to the user.
func Read(f *os.File) (string, error) {
	return read(f)
}

func readline(f *os.File) (string, error) {
	var buf [1]byte
	resultBuf := make([]byte, 0, 64)
	for {
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		resultBuf = append(resultBuf, buf[0])
	}

	return string(resultBuf), nil
}

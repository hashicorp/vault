// password is a package for reading a password securely from a terminal.
// The code in this package disables echo in the terminal so that the
// password is not echoed back in plaintext to the user.
package password

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"strings"
)

var ErrInterrupted = errors.New("interrupted")

// Read reads the password from the given os.File. The password
// will not be echoed back to the user. Ctrl-C will automatically return
// from this function with a blank string and an ErrInterrupted.
func Read(f *os.File) (string, error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)

	// Run the actual read in a go-routine so that we can still detect signals
	var result string
	var resultErr error
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		result, resultErr = read(f)
	}()

	// Wait on either the read to finish or the signal to come through
	select {
	case <-ch:
		return "", ErrInterrupted
	case <-doneCh:
		return removeiTermDelete(result), resultErr
	}
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

		// ASCII code 3 is what is sent for a Ctrl-C while reading raw.
		// If we see that, then get the interrupt. We have to do this here
		// because terminals in raw mode won't catch it at the shell level.
		if buf[0] == 3 {
			return "", ErrInterrupted
		}

		resultBuf = append(resultBuf, buf[0])
	}

	return string(resultBuf), nil
}

func removeiTermDelete(input string) string {
	return strings.TrimPrefix(input, "\x20\x7f")
}

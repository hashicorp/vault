package diff

import (
	"bufio"
	"bytes"
	"io"
)

// readLine is a helper that mimics the functionality of calling bufio.Scanner.Scan() and
// bufio.Scanner.Bytes(), but without the token size limitation. It will read and return
// the next line in the Reader with the trailing newline stripped. It will return an
// io.EOF error when there is nothing left to read (at the start of the function call). It
// will return any other errors it receives from the underlying call to ReadBytes.
func readLine(r *bufio.Reader) ([]byte, error) {
	line_, err := r.ReadBytes('\n')
	if err == io.EOF {
		if len(line_) == 0 {
			return nil, io.EOF
		}

		// ReadBytes returned io.EOF, because it didn't find another newline, but there is
		// still the remainder of the file to return as a line.
		line := line_
		return line, nil
	} else if err != nil {
		return nil, err
	}
	line := bytes.TrimSuffix(line_, []byte{'\n'})
	return line, nil
}

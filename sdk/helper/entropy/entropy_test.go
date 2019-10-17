package entropy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/errwrap"
	"testing"
)

type mockSourcer struct{}
type mockSourcerWithError struct{}
type mockSourcerFailureWithError struct{}
type mockSourcerFailureWithoutError struct{}

// simulates a successful sourcer
func (m *mockSourcer) GetRandom(bytes int) ([]byte, error) {
	return make([]byte, bytes), nil
}

// simulates a sourcer that reads in the requested number of bytes but encounters an error.
// Read should drop any error if the number of bytes specified were successfully read.
func (m *mockSourcerWithError) GetRandom(bytes int) ([]byte, error) {
	return make([]byte, bytes), errors.New("boom but you shouldn't care")
}

func (m *mockSourcerFailureWithError) GetRandom(bytes int) ([]byte, error) {
	numRetBytes := bytes - 1
	return make([]byte, numRetBytes), fmt.Errorf("requested %d bytes of entropy but only filled %d", bytes, numRetBytes)
}

func (m *mockSourcerFailureWithoutError) GetRandom(bytes int) ([]byte, error) {
	numRetBytes := bytes - 1
	return make([]byte, numRetBytes), nil
}

func TestRead(t *testing.T) {
	var tests = []struct {
		sourcer      Sourcer
		preReadBuff  []byte
		postReadBuff []byte
		outErr       error
	}{
		{
			new(mockSourcer),
			[]byte{1, 2, 3, 4},
			[]byte{0, 0, 0, 0},
			nil,
		},
		{
			new(mockSourcerWithError),
			[]byte{1, 2, 3, 4},
			[]byte{0, 0, 0, 0},
			nil,
		},
		{
			new(mockSourcerFailureWithError),
			[]byte{1, 2, 3, 4},
			nil,
			errwrap.Wrapf("unable to fill provided buffer with entropy: {{err}}", fmt.Errorf("requested %d bytes of entropy but only filled %d", 4, 3)),
		},
		{
			new(mockSourcerFailureWithoutError),
			[]byte{1, 2, 3, 4},
			nil,
			fmt.Errorf("unable to fill provided buffer with entropy"),
		},
	}

	for _, test := range tests {
		mockReader := NewReader(test.sourcer)
		buff := make([]byte, len(test.preReadBuff))
		copy(buff, test.preReadBuff)
		_, err := mockReader.Read(buff)
		// validate the error, both should be nil or have the same Error()
		switch {
		case err != nil && test.outErr != nil:
			if err.Error() != test.outErr.Error() {
				t.Fatalf("error mismatch: expected %#v got %#v", err, test.outErr)
			}
		case err != test.outErr:
			t.Fatalf("error mismatch: expected %#v got %#v", err, test.outErr)
		case err == nil && !bytes.Equal(buff, test.postReadBuff):
			t.Fatalf("after read expected buff to be: %#v but got: %#v", test.postReadBuff, buff)
		}
	}
}

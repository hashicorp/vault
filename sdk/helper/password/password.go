// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package password

import (
	"os"

	extpassword "github.com/hashicorp/go-secure-stdlib/password"
)

var ErrInterrupted = extpassword.ErrInterrupted

func Read(f *os.File) (string, error) {
	return extpassword.Read(f)
}

// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package base62

import (
	"io"

	extbase62 "github.com/hashicorp/go-secure-stdlib/base62"
)

func Random(length int) (string, error) {
	return extbase62.Random(length)
}

func RandomWithReader(length int, reader io.Reader) (string, error) {
	return extbase62.RandomWithReader(length, reader)
}

// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package mlock

import (
	extmlock "github.com/hashicorp/go-secure-stdlib/mlock"
)

func Supported() bool {
	return extmlock.Supported()
}

func LockMemory() error {
	return extmlock.LockMemory()
}

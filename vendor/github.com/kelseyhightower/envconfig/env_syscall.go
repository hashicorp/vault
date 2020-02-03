// +build !appengine

package envconfig

import "syscall"

var lookupEnv = syscall.Getenv

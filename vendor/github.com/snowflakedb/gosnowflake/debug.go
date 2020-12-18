// +build sfdebug

package gosnowflake

import "log"

func debugPanicf(fmt string, args ...interface{}) {
	log.Panicf(fmt, args...)
}

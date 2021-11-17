//go:build !sfdebug
// +build !sfdebug

package gosnowflake

func debugPanicf(fmt string, args ...interface{}) {}

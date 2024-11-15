// Copyright (c) 2018-2022 Snowflake Computing Inc. All rights reserved.

//go:build sfdebug
// +build sfdebug

package gosnowflake

import "log"

func debugPanicf(fmt string, args ...interface{}) {
	log.Panicf(fmt, args...)
}

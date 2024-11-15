// Copyright (c) 2018-2022 Snowflake Computing Inc. All rights reserved.

//go:build !sfdebug
// +build !sfdebug

package gosnowflake

func debugPanicf(fmt string, args ...interface{}) {}

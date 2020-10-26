// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build !cse

package mongocrypt

// Error represents an error from an operation on a MongoCrypt instance.
type Error struct {
	Code    int32
	Message string
}

// Error implements the error interface
func (Error) Error() string {
	panic(cseNotSupportedMsg)
}

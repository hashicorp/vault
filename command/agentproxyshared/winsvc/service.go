// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package winsvc

var chanGraceExit = make(chan int)

// ShutdownChannel returns a channel that sends a message that a shutdown
// signal has been received for the service.
func ShutdownChannel() <-chan int {
	return chanGraceExit
}

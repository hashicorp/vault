// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
)

// ExecuteExhaust reads a response from the provided StreamerConnection. This will error if the connection's
// CurrentlyStreaming function returns false.
func (op Operation) ExecuteExhaust(ctx context.Context, conn StreamerConnection) error {
	if !conn.CurrentlyStreaming() {
		return errors.New("exhaust read must be done with a connection that is currently streaming")
	}

	res, err := op.readWireMessage(ctx, conn)
	if err != nil {
		return err
	}
	if op.ProcessResponseFn != nil {
		// Server, ConnectionDescription, and CurrentIndex are unused in this mode.
		info := ResponseInfo{
			ServerResponse: res,
			Connection:     conn,
		}
		if err = op.ProcessResponseFn(info); err != nil {
			return err
		}
	}

	return nil
}

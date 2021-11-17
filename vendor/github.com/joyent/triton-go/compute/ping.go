//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joyent/triton-go/client"
	pkgerrors "github.com/pkg/errors"
)

const pingEndpoint = "/--ping"

type CloudAPI struct {
	Versions []string `json:"versions"`
}

type PingOutput struct {
	Ping     string   `json:"ping"`
	CloudAPI CloudAPI `json:"cloudapi"`
}

// Ping sends a request to the '/--ping' endpoint and returns a `pong` as well
// as a list of API version numbers your instance of CloudAPI is presenting.
func (c *ComputeClient) Ping(ctx context.Context) (*PingOutput, error) {
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   pingEndpoint,
	}
	response, err := c.Client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to ping")
	}
	if response == nil {
		return nil, pkgerrors.Wrap(err, "unable to ping")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	var result *PingOutput
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		if err != nil {
			return nil, pkgerrors.Wrap(err, "unable to decode ping response")
		}
	}

	return result, nil
}

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
	"path"
	"time"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type SnapshotsClient struct {
	client *client.Client
}

type Snapshot struct {
	Name    string
	State   string
	Created time.Time
	Updated time.Time
}

type ListSnapshotsInput struct {
	MachineID string
}

func (c *SnapshotsClient) List(ctx context.Context, input *ListSnapshotsInput) ([]*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list snapshots")
	}

	var result []*Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list snapshots response")
	}

	return result, nil
}

type GetSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Get(ctx context.Context, input *GetSnapshotInput) (*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get snapshot")
	}

	var result *Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get snapshot response")
	}

	return result, nil
}

type DeleteSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Delete(ctx context.Context, input *DeleteSnapshotInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to delete snapshot")
	}

	return nil
}

type StartMachineFromSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) StartMachine(ctx context.Context, input *StartMachineFromSnapshotInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to start machine")
	}

	return nil
}

type CreateSnapshotInput struct {
	MachineID string
	Name      string
}

func (c *SnapshotsClient) Create(ctx context.Context, input *CreateSnapshotInput) (*Snapshot, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.MachineID, "snapshots")

	data := make(map[string]interface{})
	data["name"] = input.Name

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   data,
	}

	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to create snapshot")
	}

	var result *Snapshot
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode create snapshot response")
	}

	return result, nil
}

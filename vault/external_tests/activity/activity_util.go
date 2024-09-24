// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package activity

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

type QueryResponse struct {
	StartTime   string                     `json:"start_time" mapstructure:"start_time"`
	EndTime     string                     `json:"end_time" mapstructure:"end_time"`
	ByNamespace []*vault.ResponseNamespace `json:"by_namespace"`
	Total       *vault.ResponseCounts      `json:"total"`
}

func expectEndTime(t *testing.T, expected time.Time, resp *api.Secret) {
	t.Helper()

	qr := QueryResponse{}
	mapstructure.Decode(resp.Data, &qr)
	parsedTime, err := time.Parse(time.RFC3339, qr.EndTime)
	if err != nil {
		t.Fatal(err)
	}
	if !expected.Equal(parsedTime) {
		t.Errorf("wrong end time, expected %v actual %v", expected, parsedTime)
	}
}

func expectStartTime(t *testing.T, expected time.Time, resp *api.Secret) {
	t.Helper()

	var qr QueryResponse
	mapstructure.Decode(resp.Data, &qr)
	parsedTime, err := time.Parse(time.RFC3339, qr.StartTime)
	if err != nil {
		t.Fatal(err)
	}
	if !expected.Equal(parsedTime) {
		t.Errorf("wrong start time, expected %v actual %v", expected, parsedTime)
	}
}

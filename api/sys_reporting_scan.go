// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) ReportingScan() (*ReportingScanOutput, error) {
	return c.ReportingScanWithContext(context.Background())
}

func (c *Sys) ReportingScanWithContext(ctx context.Context) (*ReportingScanOutput, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPost, "/v1/sys/reporting/scan")

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result ReportingScanOutput
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

type ReportingScanOutput struct {
	Timestamp         string `json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	FullDirectoryPath string `json:"full_directory_path" structs:"full_directory_path" mapstructure:"full_directory_path"`
}

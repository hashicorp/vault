// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package puller

import (
	"context"
	"errors"
	"io"
)

// type PullRequest struct {
// 	Type    consts.PluginType
// 	Name    string
// 	Version string
// 	SHA256  string
// }

// type PullResponse struct {
// 	SHA256 string
// }

// type Puller interface {
// 	Pull(*PullRequest) (*PullResponse, error)
// }

var errNotFound = errors.New("not found")

type build struct {
	Arch string `json:"arch"`
	OS   string `json:"os"`
	URL  string `json:"url"`
}

type metadata struct {
	Builds               []build  `json:"builds"`
	URLSHASums           string   `json:"url_shasums"`
	URLSHASumsSignatures []string `json:"url_shasums_signatures"`
	Version              string   `json:"version"`
}

type pluginSource interface {
	listMetadata(ctx context.Context, plugin string) ([]metadata, error)
	getMetadata(ctx context.Context, plugin, version string) (metadata, error)
	getContentReader(ctx context.Context, url string) (io.ReadCloser, error)
}

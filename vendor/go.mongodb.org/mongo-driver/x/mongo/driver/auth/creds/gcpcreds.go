// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package creds

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// GCPCredentialProvider provides GCP credentials.
type GCPCredentialProvider struct {
	httpClient *http.Client
}

// NewGCPCredentialProvider generates new GCPCredentialProvider
func NewGCPCredentialProvider(httpClient *http.Client) GCPCredentialProvider {
	return GCPCredentialProvider{httpClient}
}

// GetCredentialsDoc generates GCP credentials.
func (p GCPCredentialProvider) GetCredentialsDoc(ctx context.Context) (bsoncore.Document, error) {
	metadataHost := "metadata.google.internal"
	if envhost := os.Getenv("GCE_METADATA_HOST"); envhost != "" {
		metadataHost = envhost
	}
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", metadataHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: %w", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := p.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unable to retrieve GCP credentials: expected StatusCode 200, got StatusCode: %v. Response body: %s",
			resp.StatusCode,
			body)
	}
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	// Attempt to read body as JSON
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve GCP credentials: error reading body JSON: %w (response body: %s)",
			err,
			body)
	}
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: got unexpected empty accessToken from GCP Metadata Server. Response body: %s", body)
	}

	builder := bsoncore.NewDocumentBuilder().AppendString("accessToken", tokenResponse.AccessToken)
	return builder.Build(), nil
}

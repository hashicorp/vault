// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package credproviders

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)

const (
	// AzureProviderName provides a name of Azure provider
	AzureProviderName = "AzureProvider"

	azureURI = "http://169.254.169.254/metadata/identity/oauth2/token"
)

// An AzureProvider retrieves credentials from Azure IMDS.
type AzureProvider struct {
	httpClient   *http.Client
	expiration   time.Time
	expiryWindow time.Duration
}

// NewAzureProvider returns a pointer to an Azure credential provider.
func NewAzureProvider(httpClient *http.Client, expiryWindow time.Duration) *AzureProvider {
	return &AzureProvider{
		httpClient:   httpClient,
		expiration:   time.Time{},
		expiryWindow: expiryWindow,
	}
}

// RetrieveWithContext retrieves the keys from the Azure service.
func (a *AzureProvider) RetrieveWithContext(ctx context.Context) (credentials.Value, error) {
	v := credentials.Value{ProviderName: AzureProviderName}
	req, err := http.NewRequest(http.MethodGet, azureURI, nil)
	if err != nil {
		return v, fmt.Errorf("unable to retrieve Azure credentials: %w", err)
	}
	q := make(url.Values)
	q.Set("api-version", "2018-02-01")
	q.Set("resource", "https://vault.azure.net")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata", "true")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return v, fmt.Errorf("unable to retrieve Azure credentials: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, fmt.Errorf("unable to retrieve Azure credentials: error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("unable to retrieve Azure credentials: expected StatusCode 200, got StatusCode: %v. Response body: %s", resp.StatusCode, body)
	}
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}
	// Attempt to read body as JSON
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return v, fmt.Errorf("unable to retrieve Azure credentials: error reading body JSON: %w (response body: %s)", err, body)
	}
	if tokenResponse.AccessToken == "" {
		return v, fmt.Errorf("unable to retrieve Azure credentials: got unexpected empty accessToken from Azure Metadata Server. Response body: %s", body)
	}
	v.SessionToken = tokenResponse.AccessToken

	expiresIn, err := time.ParseDuration(tokenResponse.ExpiresIn + "s")
	if err != nil {
		return v, err
	}
	if expiration := expiresIn - a.expiryWindow; expiration > 0 {
		a.expiration = time.Now().Add(expiration)
	}

	return v, err
}

// Retrieve retrieves the keys from the Azure service.
func (a *AzureProvider) Retrieve() (credentials.Value, error) {
	return a.RetrieveWithContext(context.Background())
}

// IsExpired returns if the credentials have been retrieved.
func (a *AzureProvider) IsExpired() bool {
	return a.expiration.Before(time.Now())
}

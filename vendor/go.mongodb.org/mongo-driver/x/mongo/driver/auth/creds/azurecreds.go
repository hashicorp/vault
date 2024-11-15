// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package creds

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
	"go.mongodb.org/mongo-driver/internal/credproviders"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// AzureCredentialProvider provides Azure credentials.
type AzureCredentialProvider struct {
	cred *credentials.Credentials
}

// NewAzureCredentialProvider generates new AzureCredentialProvider
func NewAzureCredentialProvider(httpClient *http.Client) AzureCredentialProvider {
	return AzureCredentialProvider{
		credentials.NewCredentials(credproviders.NewAzureProvider(httpClient, 1*time.Minute)),
	}
}

// GetCredentialsDoc generates Azure credentials.
func (p AzureCredentialProvider) GetCredentialsDoc(ctx context.Context) (bsoncore.Document, error) {
	creds, err := p.cred.GetWithContext(ctx)
	if err != nil {
		return nil, err
	}
	builder := bsoncore.NewDocumentBuilder().
		AppendString("accessToken", creds.SessionToken)
	return builder.Build(), nil
}

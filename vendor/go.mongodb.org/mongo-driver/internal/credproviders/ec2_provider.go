// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package credproviders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)

const (
	// ec2ProviderName provides a name of EC2 provider
	ec2ProviderName = "EC2Provider"

	awsEC2URI       = "http://169.254.169.254/"
	awsEC2RolePath  = "latest/meta-data/iam/security-credentials/"
	awsEC2TokenPath = "latest/api/token"

	defaultHTTPTimeout = 10 * time.Second
)

// An EC2Provider retrieves credentials from EC2 metadata.
type EC2Provider struct {
	httpClient *http.Client
	expiration time.Time

	// expiryWindow will allow the credentials to trigger refreshing prior to the credentials actually expiring.
	// This is beneficial so expiring credentials do not cause request to fail unexpectedly due to exceptions.
	//
	// So a ExpiryWindow of 10s would cause calls to IsExpired() to return true
	// 10 seconds before the credentials are actually expired.
	expiryWindow time.Duration
}

// NewEC2Provider returns a pointer to an EC2 credential provider.
func NewEC2Provider(httpClient *http.Client, expiryWindow time.Duration) *EC2Provider {
	return &EC2Provider{
		httpClient:   httpClient,
		expiryWindow: expiryWindow,
	}
}

func (e *EC2Provider) getToken(ctx context.Context) (string, error) {
	req, err := http.NewRequest(http.MethodPut, awsEC2URI+awsEC2TokenPath, nil)
	if err != nil {
		return "", err
	}
	const defaultEC2TTLSeconds = "30"
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", defaultEC2TTLSeconds)

	ctx, cancel := context.WithTimeout(ctx, defaultHTTPTimeout)
	defer cancel()
	resp, err := e.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s %s failed: %s", req.Method, req.URL.String(), resp.Status)
	}

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(token) == 0 {
		return "", errors.New("unable to retrieve token from EC2 metadata")
	}
	return string(token), nil
}

func (e *EC2Provider) getRoleName(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, awsEC2URI+awsEC2RolePath, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-aws-ec2-metadata-token", token)

	ctx, cancel := context.WithTimeout(ctx, defaultHTTPTimeout)
	defer cancel()
	resp, err := e.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s %s failed: %s", req.Method, req.URL.String(), resp.Status)
	}

	role, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(role) == 0 {
		return "", errors.New("unable to retrieve role_name from EC2 metadata")
	}
	return string(role), nil
}

func (e *EC2Provider) getCredentials(ctx context.Context, token string, role string) (credentials.Value, time.Time, error) {
	v := credentials.Value{ProviderName: ec2ProviderName}

	pathWithRole := awsEC2URI + awsEC2RolePath + role
	req, err := http.NewRequest(http.MethodGet, pathWithRole, nil)
	if err != nil {
		return v, time.Time{}, err
	}
	req.Header.Set("X-aws-ec2-metadata-token", token)
	ctx, cancel := context.WithTimeout(ctx, defaultHTTPTimeout)
	defer cancel()
	resp, err := e.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return v, time.Time{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return v, time.Time{}, fmt.Errorf("%s %s failed: %s", req.Method, req.URL.String(), resp.Status)
	}

	var ec2Resp struct {
		AccessKeyID     string    `json:"AccessKeyId"`
		SecretAccessKey string    `json:"SecretAccessKey"`
		Token           string    `json:"Token"`
		Expiration      time.Time `json:"Expiration"`
	}

	err = json.NewDecoder(resp.Body).Decode(&ec2Resp)
	if err != nil {
		return v, time.Time{}, err
	}

	v.AccessKeyID = ec2Resp.AccessKeyID
	v.SecretAccessKey = ec2Resp.SecretAccessKey
	v.SessionToken = ec2Resp.Token

	return v, ec2Resp.Expiration, nil
}

// RetrieveWithContext retrieves the keys from the AWS service.
func (e *EC2Provider) RetrieveWithContext(ctx context.Context) (credentials.Value, error) {
	v := credentials.Value{ProviderName: ec2ProviderName}

	token, err := e.getToken(ctx)
	if err != nil {
		return v, err
	}

	role, err := e.getRoleName(ctx, token)
	if err != nil {
		return v, err
	}

	v, exp, err := e.getCredentials(ctx, token, role)
	if err != nil {
		return v, err
	}
	if !v.HasKeys() {
		return v, errors.New("failed to retrieve EC2 keys")
	}
	e.expiration = exp.Add(-e.expiryWindow)

	return v, nil
}

// Retrieve retrieves the keys from the AWS service.
func (e *EC2Provider) Retrieve() (credentials.Value, error) {
	return e.RetrieveWithContext(context.Background())
}

// IsExpired returns true if the credentials are expired.
func (e *EC2Provider) IsExpired() bool {
	return e.expiration.Before(time.Now())
}

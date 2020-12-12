// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type clientState int

const (
	clientStarting clientState = iota
	clientFirst
	clientFinal
	clientDone
)

type awsConversation struct {
	state    clientState
	valid    bool
	nonce    []byte
	username string
	password string
	token    string
}

type serverMessage struct {
	Nonce primitive.Binary `bson:"s"`
	Host  string           `bson:"h"`
}

type ecsResponse struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Token           string `json:"Token"`
}

const (
	amzDateFormat       = "20060102T150405Z"
	awsRelativeURI      = "http://169.254.170.2/"
	awsEC2URI           = "http://169.254.169.254/"
	awsEC2RolePath      = "latest/meta-data/iam/security-credentials/"
	awsEC2TokenPath     = "latest/api/token"
	defaultRegion       = "us-east-1"
	maxHostLength       = 255
	defaultHTTPTimeout  = 10 * time.Second
	responceNonceLength = 64
)

// Step takes a string provided from a server (or just an empty string for the
// very first conversation step) and attempts to move the authentication
// conversation forward.  It returns a string to be sent to the server or an
// error if the server message is invalid.  Calling Step after a conversation
// completes is also an error.
func (ac *awsConversation) Step(challenge []byte) (response []byte, err error) {
	switch ac.state {
	case clientStarting:
		ac.state = clientFirst
		response, err = ac.firstMsg()
	case clientFirst:
		ac.state = clientFinal
		response, err = ac.finalMsg(challenge)
	case clientFinal:
		ac.state = clientDone
		ac.valid = true
	default:
		response, err = nil, errors.New("Conversation already completed")
	}
	return
}

// Done returns true if the conversation is completed or has errored.
func (ac *awsConversation) Done() bool {
	return ac.state == clientDone
}

// Valid returns true if the conversation successfully authenticated with the
// server, including counter-validation that the server actually has the
// user's stored credentials.
func (ac *awsConversation) Valid() bool {
	return ac.valid
}

func getRegion(host string) (string, error) {
	region := defaultRegion

	if len(host) == 0 {
		return "", errors.New("invalid STS host: empty")
	}
	if len(host) > maxHostLength {
		return "", errors.New("invalid STS host: too large")
	}
	// The implicit region for sts.amazonaws.com is us-east-1
	if host == "sts.amazonaws.com" {
		return region, nil
	}
	if strings.HasPrefix(host, ".") || strings.HasSuffix(host, ".") || strings.Contains(host, "..") {
		return "", errors.New("invalid STS host: empty part")
	}

	// If the host has multiple parts, the second part is the region
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		region = parts[1]
	}

	return region, nil
}

func (ac *awsConversation) validateAndMakeCredentials() (*credentials.Credentials, error) {
	if ac.username != "" && ac.password == "" {
		return nil, errors.New("ACCESS_KEY_ID is set, but SECRET_ACCESS_KEY is missing")
	}
	if ac.username == "" && ac.password != "" {
		return nil, errors.New("SECRET_ACCESS_KEY is set, but ACCESS_KEY_ID is missing")
	}
	if ac.username == "" && ac.password == "" && ac.token != "" {
		return nil, errors.New("AWS_SESSION_TOKEN is set, but ACCESS_KEY_ID and SECRET_ACCESS_KEY are missing")
	}
	if ac.username != "" || ac.password != "" || ac.token != "" {
		return credentials.NewStaticCredentials(ac.username, ac.password, ac.token), nil
	}
	return nil, nil
}

func executeAWSHTTPRequest(req *http.Request) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return ioutil.ReadAll(resp.Body)
}

func (ac *awsConversation) getEC2Credentials() (*credentials.Credentials, error) {
	// get token
	req, err := http.NewRequest("PUT", awsEC2URI+awsEC2TokenPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "30")

	token, err := executeAWSHTTPRequest(req)
	if err != nil {
		return nil, err
	}
	if len(token) == 0 {
		return nil, errors.New("unable to retrieve token from EC2 metadata")
	}
	tokenStr := string(token)

	// get role name
	req, err = http.NewRequest("GET", awsEC2URI+awsEC2RolePath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-aws-ec2-metadata-token", tokenStr)

	role, err := executeAWSHTTPRequest(req)
	if err != nil {
		return nil, err
	}
	if len(role) == 0 {
		return nil, errors.New("unable to retrieve role_name from EC2 metadata")
	}

	// get credentials
	pathWithRole := awsEC2URI + awsEC2RolePath + string(role)
	req, err = http.NewRequest("GET", pathWithRole, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-aws-ec2-metadata-token", tokenStr)
	creds, err := executeAWSHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	var es2Resp ecsResponse
	err = json.Unmarshal(creds, &es2Resp)
	if err != nil {
		return nil, err
	}
	ac.username = es2Resp.AccessKeyID
	ac.password = es2Resp.SecretAccessKey
	ac.token = es2Resp.Token

	return ac.validateAndMakeCredentials()
}

func (ac *awsConversation) getCredentials() (*credentials.Credentials, error) {
	// Credentials passed through URI
	creds, err := ac.validateAndMakeCredentials()
	if creds != nil || err != nil {
		return creds, err
	}

	// Credentials from environment variables
	ac.username = os.Getenv("AWS_ACCESS_KEY_ID")
	ac.password = os.Getenv("AWS_SECRET_ACCESS_KEY")
	ac.token = os.Getenv("AWS_SESSION_TOKEN")

	creds, err = ac.validateAndMakeCredentials()
	if creds != nil || err != nil {
		return creds, err
	}

	// Credentials from ECS metadata
	relativeEcsURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI")
	if len(relativeEcsURI) > 0 {
		fullURI := awsRelativeURI + relativeEcsURI

		req, err := http.NewRequest("GET", fullURI, nil)
		if err != nil {
			return nil, err
		}

		body, err := executeAWSHTTPRequest(req)
		if err != nil {
			return nil, err
		}

		var espResp ecsResponse
		err = json.Unmarshal(body, &espResp)
		if err != nil {
			return nil, err
		}
		ac.username = espResp.AccessKeyID
		ac.password = espResp.SecretAccessKey
		ac.token = espResp.Token

		creds, err = ac.validateAndMakeCredentials()
		if creds != nil || err != nil {
			return creds, err
		}
	}

	// Credentials from EC2 metadata
	creds, err = ac.getEC2Credentials()
	if creds == nil && err == nil {
		return nil, errors.New("unable to get credentials")
	}
	return creds, err
}

func (ac *awsConversation) firstMsg() ([]byte, error) {
	// Values are cached for use in final message parameters
	ac.nonce = make([]byte, 32)
	_, _ = rand.Read(ac.nonce)

	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendInt32Element(msg, "p", 110)
	msg = bsoncore.AppendBinaryElement(msg, "r", 0x00, ac.nonce)
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)
	return msg, nil
}

func (ac *awsConversation) finalMsg(s1 []byte) ([]byte, error) {
	var sm serverMessage
	err := bson.Unmarshal(s1, &sm)
	if err != nil {
		return nil, err
	}

	// Check nonce prefix
	if sm.Nonce.Subtype != 0x00 {
		return nil, errors.New("server reply contained unexpected binary subtype")
	}
	if len(sm.Nonce.Data) != responceNonceLength {
		return nil, fmt.Errorf("server reply nonce was not %v bytes", responceNonceLength)
	}
	if !bytes.HasPrefix(sm.Nonce.Data, ac.nonce) {
		return nil, errors.New("server nonce did not extend client nonce")
	}

	region, err := getRegion(sm.Host)
	if err != nil {
		return nil, err
	}

	creds, err := ac.getCredentials()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UTC()
	body := "Action=GetCallerIdentity&Version=2011-06-15"

	// Create http.Request
	req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "43")
	req.Host = sm.Host
	req.Header.Set("X-Amz-Date", currentTime.Format(amzDateFormat))
	if len(ac.token) > 0 {
		req.Header.Set("X-Amz-Security-Token", ac.token)
	}
	req.Header.Set("X-MongoDB-Server-Nonce", base64.StdEncoding.EncodeToString(sm.Nonce.Data))
	req.Header.Set("X-MongoDB-GS2-CB-Flag", "n")

	// Create signer with credentials
	signer := v4.Signer{
		Credentials: creds,
	}

	// Get signed header
	_, err = signer.Sign(req, strings.NewReader(body), "sts", region, currentTime)
	if err != nil {
		return nil, err
	}

	// create message
	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendStringElement(msg, "a", req.Header.Get("Authorization"))
	msg = bsoncore.AppendStringElement(msg, "d", req.Header.Get("X-Amz-Date"))
	if len(ac.token) > 0 {
		msg = bsoncore.AppendStringElement(msg, "t", ac.token)
	}
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)

	return msg, nil
}

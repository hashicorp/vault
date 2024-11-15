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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/aws/credentials"
	v4signer "go.mongodb.org/mongo-driver/internal/aws/signer/v4"
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
	state       clientState
	valid       bool
	nonce       []byte
	credentials *credentials.Credentials
}

type serverMessage struct {
	Nonce primitive.Binary `bson:"s"`
	Host  string           `bson:"h"`
}

const (
	amzDateFormat       = "20060102T150405Z"
	defaultRegion       = "us-east-1"
	maxHostLength       = 255
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
		response = ac.firstMsg()
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

func (ac *awsConversation) firstMsg() []byte {
	// Values are cached for use in final message parameters
	ac.nonce = make([]byte, 32)
	_, _ = rand.Read(ac.nonce)

	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendInt32Element(msg, "p", 110)
	msg = bsoncore.AppendBinaryElement(msg, "r", 0x00, ac.nonce)
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)
	return msg
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

	creds, err := ac.credentials.GetWithContext(context.Background())
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
	if len(creds.SessionToken) > 0 {
		req.Header.Set("X-Amz-Security-Token", creds.SessionToken)
	}
	req.Header.Set("X-MongoDB-Server-Nonce", base64.StdEncoding.EncodeToString(sm.Nonce.Data))
	req.Header.Set("X-MongoDB-GS2-CB-Flag", "n")

	// Create signer with credentials
	signer := v4signer.NewSigner(ac.credentials)

	// Get signed header
	_, err = signer.Sign(req, strings.NewReader(body), "sts", region, currentTime)
	if err != nil {
		return nil, err
	}

	// create message
	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendStringElement(msg, "a", req.Header.Get("Authorization"))
	msg = bsoncore.AppendStringElement(msg, "d", req.Header.Get("X-Amz-Date"))
	if len(creds.SessionToken) > 0 {
		msg = bsoncore.AppendStringElement(msg, "t", creds.SessionToken)
	}
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)

	return msg, nil
}

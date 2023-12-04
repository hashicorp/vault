// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sqs

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hashicorp/vault/sdk/event"
	"github.com/hashicorp/vault/version"
)

var _ event.EventSubscriptionPlugin = (*SQS)(nil)

const pluginType = "sqs"

func New() *SQS {
	return &SQS{
		clients:   map[string]*sqs.Client{},
		queueUrls: map[string]string{},
	}
}

type SQS struct {
	clients    map[string]*sqs.Client
	queueUrls  map[string]string
	clientLock sync.RWMutex
}

func (s *SQS) Initialize(_ context.Context) error {
	return nil
}

func (s *SQS) Subscribe(ctx context.Context, request *event.SubscribeRequest) error {
	var options []func(*config.LoadOptions) error
	// TODO: support creating the queue
	accessKey := fmt.Sprintf("%v", request.Config["access_key_id"])
	secretKey := fmt.Sprintf("%v", request.Config["secret_access_key"])
	sessionToken := fmt.Sprintf("%v", request.Config["session_token"])
	region := fmt.Sprintf("%v", request.Config["region"])
	queueName := fmt.Sprintf("%v", request.Config["queue_name"])
	if accessKey != "" && secretKey != "" {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)))
	}
	if region != "" {
		options = append(options, config.WithRegion(region))
	}
	awsConfig, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return err
	}
	client := sqs.NewFromConfig(awsConfig)
	resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return err
	}
	if resp == nil || resp.QueueUrl == nil {
		return fmt.Errorf("invalid response from AWS: missing queue URL")
	}

	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.queueUrls[request.SubscriptionID] = *resp.QueueUrl
	s.clients[request.SubscriptionID] = client
	return nil
}

func (s *SQS) getClientAndURL(subscriptionID string) (*sqs.Client, string, error) {
	s.clientLock.RLock()
	defer s.clientLock.RUnlock()
	queueUrl := s.queueUrls[subscriptionID]
	client := s.clients[subscriptionID]
	if queueUrl == "" || client == nil {
		return nil, "", fmt.Errorf("invalid subscription_id")
	}
	return client, queueUrl, nil
}

func (s *SQS) SendSubscriptionEvent(subscriptionID string, eventJson string) error {
	client, queueURL, err := s.getClientAndURL(subscriptionID)
	if err != nil {
		return err
	}
	// TODO: if an error happens, we should kill the subscription
	_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody: &eventJson,
		QueueUrl:    &queueURL,
	})
	return err
}

func (s *SQS) Unsubscribe(_ context.Context, subscriptionID string) error {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	delete(s.queueUrls, subscriptionID)
	delete(s.clients, subscriptionID)
	return nil
}

func (s *SQS) Type() (string, string) {
	return pluginType, version.Version
}

func (s *SQS) Close(_ context.Context) error {
	go func() {
		s.clientLock.Lock()
		defer s.clientLock.Unlock()
		clear(s.clients)
		clear(s.queueUrls)
	}()
	return nil
}

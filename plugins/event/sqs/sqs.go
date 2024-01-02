// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sqs

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hashicorp/vault/sdk/event"
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/mapstructure"
)

var _ event.EventSubscriptionPlugin = (*sqsBackend)(nil)

const pluginType = "sqs"

// ErrQueueRequired is returned if the required queue parameters are not present.
var ErrQueueRequired = errors.New("queue_name or queue_url must be specified")

// New returns a new instance of the SQS plugin backend.
func New() *sqsBackend {
	return &sqsBackend{
		connections: map[string]*sqsConnection{},
	}
}

type sqsBackend struct {
	connections map[string]*sqsConnection
	clientLock  sync.RWMutex
}

type sqsConnection struct {
	client        *sqs.Client
	config        *sqsConfig
	queueURL      string
	ctx           context.Context
	ctxCancelFunc func() // used when we are destroying the connection or when there is an error
}

type sqsConfig struct {
	event.SubscribeConfigDefaults
	CreateQueue     bool   `mapstructure:"create_queue"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	SessionToken    string `mapstructure:"session_token"`
	Region          string `mapstructure:"region"`
	QueueName       string `mapstructure:"queue_name"`
	QueueURL        string `mapstructure:"queue_url"`
}

func (s *sqsBackend) Initialize(_ context.Context) error {
	return nil
}

func (s *sqsBackend) Subscribe(ctx context.Context, request *event.SubscribeRequest) error {
	var options []func(*config.LoadOptions) error
	var sconfig sqsConfig
	err := mapstructure.Decode(request.Config, &sconfig)
	if err != nil {
		return err
	}
	if sconfig.QueueName == "" && sconfig.QueueURL == "" {
		return ErrQueueRequired
	}
	if sconfig.AccessKeyID != "" && sconfig.SecretAccessKey != "" {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sconfig.AccessKeyID, sconfig.SecretAccessKey, sconfig.SessionToken)))
	}
	if sconfig.Region != "" {
		options = append(options, config.WithRegion(sconfig.Region))
	}
	awsConfig, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return err
	}
	client := sqs.NewFromConfig(awsConfig)

	var queueURL string
	if sconfig.CreateQueue && sconfig.QueueName != "" {
		resp, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: &sconfig.QueueName,
		})
		if err != nil {
			return err
		}
		if resp == nil || resp.QueueUrl == nil {
			return fmt.Errorf("invalid response from AWS: missing queue URL")
		}
		queueURL = *resp.QueueUrl
	} else if sconfig.QueueURL != "" {
		queueURL = sconfig.QueueURL
	} else {
		resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: &sconfig.QueueName,
		})
		if err != nil {
			return err
		}
		if resp == nil || resp.QueueUrl == nil {
			return fmt.Errorf("invalid response from AWS: missing queue URL")
		}
		queueURL = *resp.QueueUrl
	}

	connCtx, connCancel := context.WithCancel(context.Background())
	conn := &sqsConnection{
		client:        client,
		config:        &sconfig,
		queueURL:      queueURL,
		ctx:           connCtx,
		ctxCancelFunc: connCancel,
	}
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	if _, ok := s.connections[request.SubscriptionID]; ok {
		s.killConnectionWithLock(request.SubscriptionID)
	}
	s.connections[request.SubscriptionID] = conn
	return nil
}

func (s *sqsBackend) killConnection(subscriptionID string) {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	s.killConnectionWithLock(subscriptionID)
}

func (s *sqsBackend) killConnectionWithLock(subscriptionID string) {
	conn := s.connections[subscriptionID]
	conn.ctxCancelFunc()
	delete(s.connections, subscriptionID)
}

func (s *sqsBackend) getConn(subscriptionID string) (*sqsConnection, error) {
	s.clientLock.RLock()
	defer s.clientLock.RUnlock()
	conn, ok := s.connections[subscriptionID]
	if !ok {
		return nil, fmt.Errorf("invalid subscription_id")
	}
	return conn, nil
}

func (s *sqsBackend) SendSubscriptionEvent(subscriptionID string, eventJson string) error {
	conn, err := s.getConn(subscriptionID)
	if err != nil {
		return err
	}
	backoff := conn.config.NewRetryBackoff()
	for {
		_, err = conn.client.SendMessage(context.Background(), &sqs.SendMessageInput{
			MessageBody: &eventJson,
			QueueUrl:    &conn.queueURL,
		})
		if err == nil {
			return nil
		} else {
			err2 := backoff.NextSleep()
			if err2 != nil {
				err = errors.Join(err2, err)
				break
			}
		}
	}
	if err != nil {
		s.killConnection(subscriptionID)
		return err
	}
	return nil
}

func (s *sqsBackend) Unsubscribe(_ context.Context, subscriptionID string) error {
	s.killConnection(subscriptionID)
	return nil
}

func (s *sqsBackend) Type() (string, string) {
	return pluginType, version.Version
}

func (s *sqsBackend) Close(_ context.Context) error {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	var subscriptions []string
	for k := range s.connections {
		subscriptions = append(subscriptions, k)
	}
	for _, subscription := range subscriptions {
		s.killConnectionWithLock(subscription)
	}
	return nil
}

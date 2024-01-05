// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sqs

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/event"
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/mapstructure"
)

var _ event.EventSubscriptionPlugin = (*sqsBackend)(nil)

const pluginType = "sqs"

// ErrQueueRequired is returned if the required queue parameters are not present.
var ErrQueueRequired = errors.New("queue_name or queue_url must be specified")

// New returns a new instance of the SQS plugin backend.
func New() event.EventSubscriptionPlugin {
	return &sqsBackend{
		connections: map[string]*sqsConnection{},
	}
}

type sqsBackend struct {
	connections map[string]*sqsConnection
	clientLock  sync.RWMutex
}

type sqsConnection struct {
	client        *sqs.SQS
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
	Region          string `mapstructure:"region"`
	QueueName       string `mapstructure:"queue_name"`
	QueueURL        string `mapstructure:"queue_url"`
}

func newClient(sconfig *sqsConfig) (*sqs.SQS, error) {
	var options []awsutil.Option
	if sconfig.AccessKeyID != "" && sconfig.SecretAccessKey != "" {
		options = append(options, awsutil.WithAccessKey(sconfig.AccessKeyID))
		options = append(options, awsutil.WithSecretKey(sconfig.SecretAccessKey))
	}
	if sconfig.Region != "" {
		options = append(options, awsutil.WithRegion(sconfig.Region))
	}
	options = append(options, awsutil.WithEnvironmentCredentials(true))
	credConfig, err := awsutil.NewCredentialsConfig(options...)
	if err != nil {
		return nil, err
	}
	session, err := credConfig.GetSession()
	if err != nil {
		return nil, err
	}
	return sqs.New(session), nil
}

func (s *sqsBackend) Initialize(_ context.Context) error {
	return nil
}

func (s *sqsBackend) Subscribe(_ context.Context, request *event.SubscribeRequest) error {
	var sconfig sqsConfig
	err := mapstructure.Decode(request.Config, &sconfig)
	if err != nil {
		return err
	}
	if sconfig.QueueName == "" && sconfig.QueueURL == "" {
		return ErrQueueRequired
	}
	client, err := newClient(&sconfig)
	if err != nil {
		return err
	}
	var queueURL string
	if sconfig.CreateQueue && sconfig.QueueName != "" {
		resp, err := client.CreateQueue(&sqs.CreateQueueInput{
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
		resp, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
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
	err = backoff.Retry(func() error {
		_, err = conn.client.SendMessage(&sqs.SendMessageInput{
			MessageBody: &eventJson,
			QueueUrl:    &conn.queueURL,
		})
		return err
	})
	if err != nil {
		// refresh client and try again
		s.killConnection(subscriptionID)
		return err
	}
	return nil
}

func (s *sqsBackend) refreshClient(subscriptionID string) error {
	conn, err := s.getConn(subscriptionID)
	if err != nil {
		return err
	}
	client, err := newClient(conn.config)
	if err != nil {
		return err
	}
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	conn.client = client
	// probably not necessary, but just in case
	s.connections[subscriptionID] = conn
	return nil
}

func (s *sqsBackend) sendSubscriptionEventInternal(subscriptionID string, eventJson string, isRetry bool) error {
	conn, err := s.getConn(subscriptionID)
	if err != nil {
		return err
	}
	backoff := conn.config.NewRetryBackoff()
	err = backoff.Retry(func() error {
		_, err = conn.client.SendMessage(&sqs.SendMessageInput{
			MessageBody: &eventJson,
			QueueUrl:    &conn.queueURL,
		})
		return err
	})
	if err != nil && !isRetry {
		// refresh client and try again, once
		err2 := s.refreshClient(subscriptionID)
		if err2 != nil {
			return errors.Join(err, err2)
		}
		return s.sendSubscriptionEventInternal(subscriptionID, eventJson, true)
	} else if err != nil && isRetry {
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

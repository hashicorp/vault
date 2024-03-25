// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/plugins/event"
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	_ event.Factory            = New
	_ event.SubscriptionPlugin = (*kafkaPlugin)(nil)
)

const pluginName = "kafka"

// ErrAddressesRequired is returned if the addresses field is empty.
var ErrAddressesRequired = errors.New("must specify at least one address")

// New returns a new instance of the Kafka plugin backend.
func New(_ context.Context) (event.SubscriptionPlugin, error) {
	return &kafkaPlugin{
		clients: map[string]*kafkaClient{},
	}, nil
}

type kafkaPlugin struct {
	clientLock sync.RWMutex
	clients    map[string]*kafkaClient
}

type kafkaClient struct {
	w      *kafka.Writer
	config *kafkaConfig
}

type kafkaConfig struct {
	event.SubscribeConfigDefaults

	Topic string `mapstructure:"topic"` // Kafka topic to write to.

	Addresses []string `mapstructure:"addresses"` // Kafka broker addresses.
	ClientID  string   `mapstructure:"client_id"` // Client identifier, only used as metadata.

	Balancer        string `mapstructure:"balancer"`          // Client side load balancing algorithm. Defaults to "round_robin".
	RequiredAcks    string `mapstructure:"required_acks"`     // Required acks from the broker. Defaults to "none".
	Async           bool   `mapstructure:"async"`             // Fire-and-forget on the client side; incompatible with requiring acks.
	AutoCreateTopic bool   `mapstructure:"auto_create_topic"` // Create the event topic on the broker if it doesn't exist.

	BatchSize    int           `mapstructure:"batch_size"`    // Maximum number of messages to batch.
	BatchBytes   int64         `mapstructure:"batch_bytes"`   // Maximum size of a batch.
	BatchTimeout time.Duration `mapstructure:"batch_timeout"` // Maximum time to wait for a batch to fill.

	CAPem         string `mapstructure:"ca_pem"`          // PEM encoded CA certificate.
	TLSServerName string `mapstructure:"tls_server_name"` // ServerName to verify for TLS connections.
	TLSSkipVerify bool   `mapstructure:"tls_skip_verify"` // Set to skip server cert verification.
	TLSDisabled   bool   `mapstructure:"tls_disabled"`    // Set for PLAINTEXT listeners. Default assumes SSL listeners.

	SASLType string `mapstructure:"sasl_type"` // SASL (auth) mechanism type.
	Username string `mapstructure:"username"`  // SASL username.
	Password string `mapstructure:"password"`  // SASL password.
}

func newClient(kconfig *kafkaConfig) (*kafkaClient, error) {
	var balancer kafka.Balancer
	switch kconfig.Balancer {
	case "round_robin":
		balancer = &kafka.RoundRobin{}
	case "hash":
		balancer = &kafka.Hash{}
	case "least_bytes":
		balancer = &kafka.LeastBytes{}
	}

	var requiredAcks kafka.RequiredAcks
	switch kconfig.RequiredAcks {
	case "none":
		requiredAcks = kafka.RequireNone
	case "one":
		requiredAcks = kafka.RequireOne
	case "all":
		requiredAcks = kafka.RequireAll
	}

	var sasl sasl.Mechanism
	var err error
	// TODO: support more SASL types.
	switch kconfig.SASLType {
	case "plain":
		sasl = &plain.Mechanism{Username: kconfig.Username, Password: kconfig.Password}
	case "scram_sha256":
		sasl, err = scram.Mechanism(scram.SHA256, kconfig.Username, kconfig.Password)
	case "scram_sha512":
		sasl, err = scram.Mechanism(scram.SHA512, kconfig.Username, kconfig.Password)
	}
	if err != nil {
		return nil, err
	}

	clientID := "vault.hashicorp.com"
	if kconfig.ClientID != "" {
		clientID = kconfig.ClientID
	}

	var tlsConfig *tls.Config
	if !kconfig.TLSDisabled {
		var certPool *x509.CertPool
		if kconfig.CAPem != "" {
			certPool = x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(kconfig.CAPem))
			if !ok {
				return nil, errors.New("failed to parse any certificates from ca_pem contents")
			}
		}

		tlsConfig = &tls.Config{
			// TODO: support client cert auth.
			RootCAs:            certPool,
			ServerName:         kconfig.TLSServerName,
			InsecureSkipVerify: kconfig.TLSSkipVerify,
		}
	}

	writer := &kafka.Writer{
		Addr: kafka.TCP(kconfig.Addresses...),

		// TODO: writers could be shared by multiple subscriptions if they have
		// the same config except for the topic, and we instead set the topic
		// on individual messages. Sharing clients more broadly would improve
		// performance for the client side load balancing algorithms if users
		// set up multiple subscriptions to the same kafka cluster.
		Topic: kconfig.Topic,

		Balancer:               balancer,
		RequiredAcks:           requiredAcks,
		Async:                  kconfig.Async,
		AllowAutoTopicCreation: kconfig.AutoCreateTopic,

		BatchSize:    kconfig.BatchSize,
		BatchBytes:   kconfig.BatchBytes,
		BatchTimeout: kconfig.BatchTimeout,

		MaxAttempts:     kconfig.GetRetries() + 1,
		WriteBackoffMin: kconfig.GetRetryMinBackoff(),
		WriteBackoffMax: kconfig.GetRetryMaxBackoff(),

		Transport: &kafka.Transport{
			ClientID: clientID,
			SASL:     sasl,
			TLS:      tlsConfig,
		},
	}

	return &kafkaClient{
		w:      writer,
		config: kconfig,
	}, nil
}

func (k *kafkaPlugin) Subscribe(_ context.Context, request *event.SubscribeRequest) error {
	var cfg kafkaConfig
	err := mapstructure.Decode(request.Config, &cfg)
	if err != nil {
		return err
	}
	if len(cfg.Addresses) == 0 {
		return ErrAddressesRequired
	}

	k.clientLock.Lock()
	defer k.clientLock.Unlock()

	if _, ok := k.clients[request.SubscriptionID]; ok {
		k.killClientWithLock(request.SubscriptionID)
	}

	client, err := newClient(&cfg)
	if err != nil {
		return err
	}

	// TODO: Handle ValidateConnection param.

	k.clients[request.SubscriptionID] = client
	return nil
}

func (k *kafkaPlugin) killClient(subscriptionID string) error {
	k.clientLock.Lock()
	defer k.clientLock.Unlock()
	return k.killClientWithLock(subscriptionID)
}

func (k *kafkaPlugin) killClientWithLock(subscriptionID string) error {
	client, ok := k.clients[subscriptionID]
	if !ok {
		return nil
	}

	delete(k.clients, subscriptionID)

	return client.w.Close()
}

func (k *kafkaPlugin) getClient(subscriptionID string) (*kafkaClient, error) {
	k.clientLock.RLock()
	defer k.clientLock.RUnlock()
	client, ok := k.clients[subscriptionID]
	if !ok {
		return nil, fmt.Errorf("invalid subscription_id")
	}

	return client, nil
}

func (k *kafkaPlugin) Send(ctx context.Context, send *event.SendRequest) error {
	client, err := k.getClient(send.SubscriptionID)
	if err != nil {
		return err
	}

	backoff := client.config.NewRetryBackoff()
	return backoff.Retry(func() error {
		return client.w.WriteMessages(ctx, kafka.Message{
			// TODO: Validate that this key is a sensible choice. It puts all
			// messages for a subscription in the same partition.
			Key:   []byte(send.SubscriptionID),
			Value: []byte(send.EventJSON),
		})
	})
}

func (k *kafkaPlugin) Unsubscribe(_ context.Context, request *event.UnsubscribeRequest) error {
	return k.killClient(request.SubscriptionID)
}

func (k *kafkaPlugin) PluginMetadata() *event.PluginMetadata {
	return &event.PluginMetadata{
		Name:    pluginName,
		Version: version.GetVersion().Version,
	}
}

func (k *kafkaPlugin) Close(_ context.Context) error {
	k.clientLock.Lock()
	defer k.clientLock.Unlock()
	var subscriptions []string
	for k := range k.clients {
		subscriptions = append(subscriptions, k)
	}

	var err error
	for _, subscription := range subscriptions {
		err = errors.Join(err, k.killClientWithLock(subscription))
	}

	return err
}

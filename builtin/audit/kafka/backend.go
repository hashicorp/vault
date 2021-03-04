package kafka

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

type MessageType string

const (
	Request  MessageType = "request"
	Response MessageType = "response"
)

func Factory(ctx context.Context, conf *audit.BackendConfig) (audit.Backend, error) {
	address, ok := conf.Config["address"]
	if !ok {
		return nil, fmt.Errorf("address is required")
	}

	topic, ok := conf.Config["topic"]
	if !ok {
		return nil, fmt.Errorf("topic is required")
	}

	producerCfg, err := producerFromConfig(conf.Config)
	if err != nil {
		return nil, err
	}
	err = producerCfg.Validate()
	if err != nil {
		return nil, err
	}

	syncProducer, err := sarama.NewSyncProducer([]string{address}, producerCfg)
	if err != nil {
		return nil, err
	}

	// Check if hashing of accessor is disabled
	hmacAccessor := true
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		value, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		hmacAccessor = value
	}

	// Check if raw logging is enabled
	logRaw := false
	if raw, ok := conf.Config["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	b := &Backend{
		topic:       topic,
		address:     address,
		producer:    syncProducer,
		kafkaConfig: producerCfg,
		formatConfig: audit.FormatterConfig{
			Raw:          logRaw,
			HMACAccessor: hmacAccessor,
		},
	}

	return b, nil
}

type Backend struct {
	address      string
	topic        string
	producer     sarama.SyncProducer
	kafkaConfig  *sarama.Config
	formatConfig audit.FormatterConfig

	saltMutex  sync.RWMutex
	salt       *salt.Salt
	saltConfig *salt.Config
	saltView   logical.Storage
}

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	return b.sendToKafka(ctx, Request, in)
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	return b.sendToKafka(ctx, Response, in)
}

func (b *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", err
	}
	return audit.HashString(salt, data), nil
}

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	b.saltMutex.RLock()
	if b.salt != nil {
		defer b.saltMutex.RUnlock()
		return b.salt, nil
	}
	b.saltMutex.RUnlock()
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	salt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

func (b *Backend) Reload(context.Context) error {
	b.producer.Close()

	syncProducer, err := sarama.NewSyncProducer([]string{b.address}, b.kafkaConfig)
	b.producer = syncProducer
	return err
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt = nil
}

func (b *Backend) sendToKafka(ctx context.Context, messageType MessageType, in *logical.LogInput) error {
	key := in.Request.ID
	value, err := logToMessage(messageType, ctx, b, in)

	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: b.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(*value),
	}

	_, _, err = b.producer.SendMessage(msg)

	return err
}

func logToMessage(messageType MessageType, ctx context.Context, b *Backend, in *logical.LogInput) (*string, error) {
	var buf bytes.Buffer
	formatter := &audit.AuditFormatter{
		AuditFormatWriter: &audit.JSONFormatWriter{
			SaltFunc: b.Salt,
		},
	}

	if messageType == Response {
		if err := formatter.FormatResponse(ctx, &buf, b.formatConfig, in); err != nil {
			return nil, err
		}
	}
	if messageType == Request {
		if err := formatter.FormatRequest(ctx, &buf, b.formatConfig, in); err != nil {
			return nil, err
		}
	}
	message := buf.String()

	return &message, nil
}

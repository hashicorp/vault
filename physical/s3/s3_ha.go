package s3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	LockRenewInterval      = 5 * time.Second
	LockRetryInterval      = 5 * time.Second
	LockTTL                = 15 * time.Second
	LockWatchRetryInterval = 5 * time.Second
	LockWatchRetryMax      = 5
)

var (
	_ physical.HABackend = (*S3Backend)(nil)
	_ physical.Lock      = (*S3Lock)(nil)
)

var (
	// metricLockUnlock is the metric to register for a lock delete.
	metricLockUnlock = []string{"s3", "lock", "unlock"}

	// metricLockLock is the metric to register for a lock get.
	metricLockLock = []string{"s3", "lock", "lock"}

	// metricLockValue is the metric to register for a lock create/update.
	metricLockValue = []string{"s3", "lock", "value"}
)

type S3Lock struct {
	backend  *S3Backend
	key      string
	value    string
	held     bool
	identity string

	stopCh  chan struct{}
	stopped bool

	renewInterval      time.Duration
	retryInterval      time.Duration
	ttl                time.Duration
	watchRetryInterval time.Duration
	watchRetryMax      int
}

func (l *S3Lock) Unlock() error {
	defer metrics.MeasureSince(metricLockUnlock, time.Now())

	// First verify the lock exists and we own it
	exists, _, err := l.Value()
	if err != nil {
		return err
	}
	if !exists {
		// Lock is already gone
		return nil
	}

	// Delete the lock file
	_, err = l.backend.client.DeleteObject(l.backend.context, &s3.DeleteObjectInput{
		Bucket: aws.String(l.backend.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", l.backend.path, l.key)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete lock: %w", err)
	}

	// We are no longer holding the lock
	l.held = false

	return nil
}

type LockRecord struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Identity  string    `json:"identity"`
	Timestamp time.Time `json:"timestamp"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *S3Backend) HAEnabled() bool {
	return s.haEnabled
}

func (s *S3Backend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate identity: %w", err)
	}

	return &S3Lock{
		backend:  s,
		key:      key,
		value:    value,
		identity: identity,

		renewInterval:      LockRenewInterval,
		retryInterval:      LockRetryInterval,
		ttl:                LockTTL,
		watchRetryInterval: LockWatchRetryInterval,
		watchRetryMax:      LockWatchRetryMax,
	}, nil
}

func (l *S3Lock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	defer metrics.MeasureSince(metricLockLock, time.Now())

	// Attempt to lock - this function blocks until a lock is acquired or an error
	// occurs.

	acquired, err := l.attemptLock(stopCh)
	if err != nil {
		return nil, fmt.Errorf("lock: %w", err)
	}
	if !acquired {
		return nil, nil
	}

	// We have the lock now
	l.held = true

	// Build the locks
	l.stopCh = make(chan struct{})

	// Periodically renew and watch the lock
	go l.watchLock()

	return l.stopCh, nil
}

func (l *S3Lock) Value() (bool, string, error) {
	defer metrics.MeasureSince(metricLockValue, time.Now())
	lockPath := fmt.Sprintf("%s/%s", l.backend.path, l.key)

	input := &s3.GetObjectInput{
		Bucket: aws.String(l.backend.bucket),
		Key:    aws.String(lockPath),
	}

	result, err := l.backend.client.GetObject(l.backend.context, input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to get lock: %w", err)
	}
	defer result.Body.Close()

	var lockInfo LockRecord
	if err := json.NewDecoder(result.Body).Decode(&lockInfo); err != nil {
		return false, "", fmt.Errorf("failed to decode lock info: %w", err)
	}

	return true, lockInfo.Value, nil
}

func (l *S3Lock) attemptLock(stopCh <-chan struct{}) (bool, error) {
	ticker := time.NewTicker(l.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			acquired, err := l.writeLock()
			if err != nil {
				return false, fmt.Errorf("attempt lock: %w", err)
			}
			if !acquired {
				continue
			}

			return true, nil
		case <-stopCh:
			return false, nil
		}
	}
}

func (l *S3Lock) writeLock() (bool, error) {
	lockPath := fmt.Sprintf("%s/%s", l.backend.path, l.key)

	lockInfo := LockRecord{
		Key:       lockPath,
		Value:     l.value,
		Identity:  l.identity,
		Timestamp: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(l.ttl),
	}

	lockData, err := json.Marshal(lockInfo)
	if err != nil {
		return false, fmt.Errorf("failed to marshal lock info: %w", err)
	}

	// Check if the object already exists
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(l.backend.bucket),
		Key:    aws.String(lockPath),
	}
	_, err = l.backend.client.HeadObject(l.backend.context, headInput)
	if err == nil {
		return false, nil // Lock already exists
	}

	var nsk *types.NoSuchKey
	if errors.As(err, &nsk) {
		return false, fmt.Errorf("failed to verify lock existence: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(l.backend.bucket),
		Key:         aws.String(lockPath),
		Body:        bytes.NewReader(lockData),
		ContentType: aws.String("application/json"),
		IfNoneMatch: aws.String("*"),
		Expires:     aws.Time(time.Now().Add(l.ttl)),
	}

	if l.backend.kmsKeyId != "" {
		input.ServerSideEncryption = types.ServerSideEncryptionAwsKms
		input.SSEKMSKeyId = aws.String(l.backend.kmsKeyId)
	}

	_, err = l.backend.client.PutObject(l.backend.context, input)
	if err != nil {
		return false, fmt.Errorf("failed to create lock: %w", err)
	}

	return true, nil
}

func (l *S3Lock) watchLock() {
	ticker := time.NewTicker(l.watchRetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lockPath := fmt.Sprintf("%s/%s", l.backend.path, l.key)
			input := &s3.GetObjectInput{
				Bucket: aws.String(l.backend.bucket),
				Key:    aws.String(lockPath),
			}

			result, err := l.backend.client.GetObject(l.backend.context, input)
			if err != nil {
				var nsk *types.NoSuchKey
				if errors.As(err, &nsk) {
					return
				}
				continue
			}

			var lockInfo LockRecord
			if err := json.NewDecoder(result.Body).Decode(&lockInfo); err != nil {
				continue
			}
			result.Body.Close()

			if lockInfo.Identity != l.identity {
				return
			}
		case <-l.stopCh:
			return
		}
	}
}

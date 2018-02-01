package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
)

type credentialRotationManager struct {
	l sync.Mutex

	initialized  bool
	nextRotation map[string]time.Time
}

// initialized must be read with the lock
func (m *credentialRotationManager) initialize(ctx context.Context, s logical.Storage) error {
	if m.initialized {
		return nil
	}

	m.nextRotation = make(map[string]time.Time)

	dbs, err := s.List(ctx, "config/")
	if err != nil {
		return err
	}

	for _, dbName := range dbs {
		configEntry, err := s.Get(ctx, fmt.Sprintf("config/%s", dbName))
		if err != nil {
			return err
		}

		var config DatabaseConfig
		if err := configEntry.DecodeJSON(&config); err != nil {
			return err
		}
		m.add(s, dbName, config.RootCredentialsRotateInterval)
	}

	m.initialized = true
	return nil
}

func (m *credentialRotationManager) Add(ctx context.Context, s logical.Storage, name string, nextRotation time.Duration) error {
	m.l.Lock()
	defer m.l.Unlock()

	if err := m.initialize(ctx, s); err != nil {
		return err
	}

	m.add(s, name, nextRotation)
	return nil
}

func (m *credentialRotationManager) add(s logical.Storage, name string, updateInterval time.Duration) {
	if updateInterval == 0 {
		m.nextRotation[name] = time.Unix(1<<63-1, 999999999) // Max date
		return
	}
	m.nextRotation[name] = time.Now().Add(updateInterval)
}

func (m *credentialRotationManager) Remove(name string) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.nextRotation, name)
}

func (m *credentialRotationManager) NeedRotate(ctx context.Context, s logical.Storage) ([]string, error) {
	m.l.Lock()
	defer m.l.Unlock()

	if err := m.initialize(ctx, s); err != nil {
		return nil, err
	}

	var toRotate []string
	now := time.Now()
	for db, nextRotation := range m.nextRotation {
		if nextRotation.Before(now) {
			toRotate = append(toRotate, db)
		}
	}

	return toRotate, nil
}

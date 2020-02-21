package credsutil

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

const (
	DefaultExpirationFormat = "2006-01-02 15:04:05-0700"
)

// UsernamePasswordProducer is a combined wrapper around UsernameProducer and PasswordProducer that adheres
// to the CredentialsProducer interface
type UsernamePasswordProducer struct {
	UsernameProducer
	PasswordProducer

	ExpirationFormat string
}

type UsernamePasswordOpt func(*UsernamePasswordProducer) error

// UsernameOpts is a collection of options passed into the UsernameProducer.
// This is a passthrough to the UsernameProducer. Multiple calls to this will not build upon each other.
func UsernameOpts(opts ...UsernameOpt) UsernamePasswordOpt {
	return func(up *UsernamePasswordProducer) error {
		user, err := NewUsernameProducer(opts...)
		if err != nil {
			return fmt.Errorf("unable to create username producer: %w", err)
		}

		up.UsernameProducer = user
		return nil
	}
}

// PasswordOpts is a collection of options passed into the PasswordProducer.
// This is a passthrough to the PasswordProducer. Multiple calls to this will not build upon each other.
func PasswordOpts(opts ...PasswordOpt) UsernamePasswordOpt {
	return func(up *UsernamePasswordProducer) error {
		pass, err := NewPasswordProducer(opts...)
		if err != nil {
			return fmt.Errorf("unable to create username producer: %w", err)
		}

		up.PasswordProducer = pass
		return nil
	}
}

// ExpirationFormat for use when generating an expiration string.
func ExpirationFormat(format string) UsernamePasswordOpt {
	return func(up *UsernamePasswordProducer) error {
		up.ExpirationFormat = format
		return nil
	}
}

// NewUsernamePasswordProducer creates a UsernamePasswordProducer that can be used when the CredentialsProducer
// interface is embedded in a struct. This adheres to the entire CredentialsProducer interface.
func NewUsernamePasswordProducer(opts ...UsernamePasswordOpt) (up UsernamePasswordProducer, err error) {
	merr := &multierror.Error{}
	for _, opt := range opts {
		merr = multierror.Append(merr, opt(&up))
	}

	return up, merr.ErrorOrNil()
}

func (up UsernamePasswordProducer) GenerateUsername(config dbplugin.UsernameConfig) (string, error) {
	return up.UsernameProducer.GenerateUsername(config)
}

func (up UsernamePasswordProducer) GeneratePassword() (string, error) {
	return up.PasswordProducer.GeneratePassword()
}

func (up UsernamePasswordProducer) GenerateCredentials(ctx context.Context) (string, error) {
	return up.GeneratePassword()
}

func (up UsernamePasswordProducer) GenerateExpiration(t time.Time) (string, error) {
	if up.ExpirationFormat == "" {
		return t.Format(DefaultExpirationFormat), nil
	}
	return t.Format(up.ExpirationFormat), nil
}

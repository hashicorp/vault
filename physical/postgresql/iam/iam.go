package iam

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type DBConfig struct {
	URL string

	UseAWSIAMAuth bool
	AWSDBRegion   string
	Logger        hclog.Logger
}

type authToken struct {
	token       string
	valid       bool
	lock        sync.Mutex
	createdTime time.Time
}

func newAuthToken(dbConfig DBConfig, pgConfig pgx.ConnConfig) (*authToken, error) {
	// Will have a switch case here for different clouds
	token, err := fetchAuthToken(dbConfig, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("fetching aws token: %v", err)
	}

	return &authToken{
		token:       token,
		valid:       false,
		lock:        sync.Mutex{},
		createdTime: time.Now(),
	}, nil
}

func (t *authToken) getTokenString(dbConfig DBConfig, pgConfig pgx.ConnConfig) (string, error) {
	if time.Since(t.createdTime) <= 10*time.Minute {
		return t.token, nil
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	if time.Since(t.createdTime) <= 2*time.Minute {
		return t.token, nil
	}

	token, err := fetchAuthToken(dbConfig, pgConfig)
	if err != nil {
		return "", err
	}

	t.token = token
	t.createdTime = time.Now()

	return t.token, nil
}

func DBHandler(dbConfig DBConfig) (*sql.DB, error) {
	// We would need some validations here

	connConfig, err := pgx.ParseConfig(dbConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %v", err)
	}

	// create new token since this is the entry point
	token, err := newAuthToken(dbConfig, *connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token: %v", err)
	}

	beforeConnect, err := BeforeConnectFn(token, dbConfig, *connConfig)
	if err != nil {
		return nil, fmt.Errorf("generating before connect function: %v", err)
	}

	db := stdlib.OpenDB(*connConfig, stdlib.OptionBeforeConnect(beforeConnect))
	return db, nil
}

func BeforeConnectFn(token *authToken, dbConfig DBConfig, pgConfig pgx.ConnConfig) (func(context.Context, *pgx.ConnConfig) error, error) {
	var beforeConnect func(context.Context, *pgx.ConnConfig) error

	if dbConfig.UseAWSIAMAuth {
		beforeConnect = func(ctx context.Context, config *pgx.ConnConfig) error {
			tokenVal, err := token.getTokenString(dbConfig, *config)
			if err != nil {
				return fmt.Errorf("fetching aws token value: %v", err)
			}

			dbConfig.Logger.Info("setting password for AWS IAM auth", "host", config.Host, "token", tokenVal)
			config.Password = tokenVal
			return nil
		}
	}

	return beforeConnect, nil
}

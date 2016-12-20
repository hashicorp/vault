package dbs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/strutil"
)

type Cassandra struct {
	// Session is goroutine safe, however, since we reinitialize
	// it when connection info changes, we want to make sure we
	// can close it and use a new connection; hence the lock
	ConnectionProducer
	CredentialsProducer
	sync.RWMutex
}

func (c *Cassandra) Type() string {
	return cassandraTypeName
}

func (c *Cassandra) getConnection() (*gocql.Session, error) {
	session, err := c.Connection()
	if err != nil {
		return nil, err
	}

	return session.(*gocql.Session), nil
}

func (c *Cassandra) CreateUser(createStmt, rollbackStmt, username, password, expiration string) error {
	// Get the connection
	session, err := c.getConnection()
	if err != nil {
		return err
	}

	// TODO: This is racey
	// Grab a read lock
	c.RLock()
	defer c.RUnlock()

	// Set consistency
	/*	if .Consistency != "" {
		consistencyValue, err := gocql.ParseConsistencyWrapper(role.Consistency)
		if err != nil {
			return err
		}

		session.SetConsistency(consistencyValue)
	}*/

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(createStmt, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		err = session.Query(queryHelper(query, map[string]string{
			"username": username,
			"password": password,
		})).Exec()
		if err != nil {
			for _, query := range strutil.ParseArbitraryStringSlice(rollbackStmt, ";") {
				query = strings.TrimSpace(query)
				if len(query) == 0 {
					continue
				}

				session.Query(queryHelper(query, map[string]string{
					"username": username,
					"password": password,
				})).Exec()
			}
			return err
		}
	}

	return nil
}

func (c *Cassandra) RenewUser(username, expiration string) error {
	// NOOP
	return nil
}

func (c *Cassandra) RevokeUser(username, revocationSQL string) error {
	session, err := c.getConnection()
	if err != nil {
		return err
	}
	// TODO: this is Racey
	c.RLock()
	defer c.RUnlock()

	err = session.Query(fmt.Sprintf("DROP USER '%s'", username)).Exec()
	if err != nil {
		return fmt.Errorf("error removing user %s", username)
	}

	return nil
}

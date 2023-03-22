// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package couchdb

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestCouchDBBackend(t *testing.T) {
	cleanup, config := prepareCouchdbDBTestContainer(t)
	defer cleanup()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewCouchDBBackend(map[string]string{
		"endpoint": config.URL().String(),
		"username": config.username,
		"password": config.password,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestTransactionalCouchDBBackend(t *testing.T) {
	cleanup, config := prepareCouchdbDBTestContainer(t)
	defer cleanup()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewTransactionalCouchDBBackend(map[string]string{
		"endpoint": config.URL().String(),
		"username": config.username,
		"password": config.password,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

type couchDB struct {
	baseURL  url.URL
	dbname   string
	username string
	password string
}

func (c couchDB) Address() string {
	return c.baseURL.Host
}

func (c couchDB) URL() *url.URL {
	u := c.baseURL
	u.Path = c.dbname
	return &u
}

var _ docker.ServiceConfig = &couchDB{}

func prepareCouchdbDBTestContainer(t *testing.T) (func(), *couchDB) {
	// If environment variable is set, assume caller wants to target a real
	// DynamoDB.
	if os.Getenv("COUCHDB_ENDPOINT") != "" {
		return func() {}, &couchDB{
			baseURL:  url.URL{Host: os.Getenv("COUCHDB_ENDPOINT")},
			username: os.Getenv("COUCHDB_USERNAME"),
			password: os.Getenv("COUCHDB_PASSWORD"),
		}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName:   "couchdb",
		ImageRepo:       "docker.mirror.hashicorp.services/library/couchdb",
		ImageTag:        "1.6",
		Ports:           []string{"5984/tcp"},
		DoNotAutoRemove: true,
	})
	if err != nil {
		t.Fatalf("Could not start local CouchDB: %s", err)
	}

	svc, err := runner.StartService(context.Background(), setupCouchDB)
	if err != nil {
		t.Fatalf("Could not start local CouchDB: %s", err)
	}

	return svc.Cleanup, svc.Config.(*couchDB)
}

func setupCouchDB(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	c := &couchDB{
		baseURL:  url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)},
		dbname:   fmt.Sprintf("vault-test-%d", time.Now().Unix()),
		username: "admin",
		password: "admin",
	}

	{
		resp, err := http.Get(c.baseURL.String())
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("expected couchdb to return status code 200, got (%s) instead", resp.Status)
		}
	}

	{
		req, err := http.NewRequest("PUT", c.URL().String(), nil)
		if err != nil {
			return nil, fmt.Errorf("could not create create database request: %q", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not create database: %q", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			bs, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to create database: %s %s\n", resp.Status, string(bs))
		}
	}

	{
		u := c.baseURL
		u.Path = fmt.Sprintf("_config/admins/%s", c.username)
		req, err := http.NewRequest("PUT", u.String(), strings.NewReader(fmt.Sprintf(`"%s"`, c.password)))
		if err != nil {
			return nil, fmt.Errorf("Could not create admin user request: %q", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Could not create admin user: %q", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bs, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("Failed to create admin user: %s %s\n", resp.Status, string(bs))
		}
	}

	return c, nil
}

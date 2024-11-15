// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package couchbase

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/mitchellh/mapstructure"
)

type couchbaseDBConnectionProducer struct {
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	ProjectID   string `json:"project_id"`
	Hosts       string `json:"hosts"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	TLS         bool   `json:"tls"`
	InsecureTLS bool   `json:"insecure_tls"`
	Base64Pem   string `json:"base64pem"`
	BucketName  string `json:"bucket_name"`

	Initialized bool
	rawConfig   map[string]interface{}
	Type        string
	cluster     *gocb.Cluster
	sync.RWMutex
}

func (c *couchbaseDBConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password: "[password]",
		c.Username: "[username]",
	}
}

func (c *couchbaseDBConnectionProducer) Init(ctx context.Context, initConfig map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error) {
	// Don't let anyone read or write the config while we're using it
	c.Lock()
	defer c.Unlock()

	c.rawConfig = initConfig

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           c,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(initConfig)
	if err != nil {
		return nil, err
	}

	switch {
	case len(c.Hosts) == 0:
		return nil, fmt.Errorf("hosts cannot be empty")
	case len(c.Username) == 0:
		return nil, fmt.Errorf("username cannot be empty")
	case len(c.Password) == 0:
		return nil, fmt.Errorf("password cannot be empty")
	}

	if c.TLS {
		if len(c.Base64Pem) == 0 {
			return nil, fmt.Errorf("base64pem cannot be empty")
		}

		if !strings.HasPrefix(c.Hosts, "couchbases://") {
			return nil, fmt.Errorf("hosts list must start with couchbases:// for TLS connection")
		}
	}

	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			c.close()
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return initConfig, nil
}

func (c *couchbaseDBConnectionProducer) Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, config, verifyConnection)
	return err
}

func (c *couchbaseDBConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	// This is intentionally not grabbing the lock since the calling functions
	// (e.g. CreateUser) are claiming it.

	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	if c.cluster != nil {
		return c.cluster, nil
	}
	var err error
	var sec gocb.SecurityConfig
	var pem []byte

	if c.TLS {
		pem, err = base64.StdEncoding.DecodeString(c.Base64Pem)
		if err != nil {
			return nil, errwrap.Wrapf("error decoding Base64Pem: {{err}}", err)
		}
		rootCAs := x509.NewCertPool()
		ok := rootCAs.AppendCertsFromPEM([]byte(pem))
		if !ok {
			return nil, fmt.Errorf("failed to parse root certificate")
		}
		sec = gocb.SecurityConfig{
			TLSRootCAs:    rootCAs,
			TLSSkipVerify: c.InsecureTLS,
		}
	}

	c.cluster, err = gocb.Connect(
		c.Hosts,
		gocb.ClusterOptions{
			Username:       c.Username,
			Password:       c.Password,
			SecurityConfig: sec,
		})
	if err != nil {
		return nil, errwrap.Wrapf("error in Connection: {{err}}", err)
	}

	// For databases 6.0 and earlier, we will need to open a `Bucket instance before connecting to any other
	// HTTP services such as UserManager.

	if c.BucketName != "" {
		bucket := c.cluster.Bucket(c.BucketName)
		// We wait until the bucket is definitely connected and setup.
		err = bucket.WaitUntilReady(computeTimeout(ctx), nil)
		if err != nil {
			return nil, errwrap.Wrapf("error in Connection waiting for bucket: {{err}}", err)
		}
	} else {
		err = c.cluster.WaitUntilReady(computeTimeout(ctx), nil)

		if err != nil {
			return nil, errwrap.Wrapf("error in Connection waiting for cluster: {{err}}", err)
		}
	}

	return c.cluster, nil
}

// close terminates the database connection without locking
func (c *couchbaseDBConnectionProducer) close() error {
	if c.cluster != nil {
		if err := c.cluster.Close(&gocb.ClusterCloseOptions{}); err != nil {
			return err
		}
	}

	c.cluster = nil
	return nil
}

// Close terminates the database connection with locking
func (c *couchbaseDBConnectionProducer) Close() error {
	// Don't let anyone read or write the config while we're using it
	c.Lock()
	defer c.Unlock()

	return c.close()
}

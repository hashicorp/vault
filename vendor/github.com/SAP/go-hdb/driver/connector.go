/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"sync"
)

/*
A Connector represents a hdb driver in a fixed configuration.
A Connector can be passed to sql.OpenDB (starting from go 1.10) allowing users to bypass a string based data source name.
*/
type Connector struct {
	mu                             sync.RWMutex
	host, username, password       string
	locale                         string
	bufferSize, fetchSize, timeout int
	tlsConfig                      *tls.Config
}

func newConnector() *Connector {
	return &Connector{
		fetchSize: DefaultFetchSize,
		timeout:   DefaultTimeout,
	}
}

// NewBasicAuthConnector creates a connector for basic authentication.
func NewBasicAuthConnector(host, username, password string) *Connector {
	c := newConnector()
	c.host = host
	c.username = username
	c.password = password
	return c
}

// NewDSNConnector creates a connector from a data source name.
func NewDSNConnector(dsn string) (*Connector, error) {
	c := newConnector()

	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	c.host = url.Host

	if url.User != nil {
		c.username = url.User.Username()
		c.password, _ = url.User.Password()
	}

	var certPool *x509.CertPool

	for k, v := range url.Query() {
		switch k {

		default:
			return nil, fmt.Errorf("URL parameter %s is not supported", k)

		case DSNFetchSize:
			if len(v) == 0 {
				continue
			}
			fetchSize, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("failed to parse fetchSize: %s", v[0])
			}
			if fetchSize < minFetchSize {
				c.fetchSize = minFetchSize
			} else {
				c.fetchSize = fetchSize
			}

		case DSNTimeout:
			if len(v) == 0 {
				continue
			}
			timeout, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("failed to parse timeout: %s", v[0])
			}
			if timeout < minTimeout {
				c.timeout = minTimeout
			} else {
				c.timeout = timeout
			}

		case DSNLocale:
			if len(v) == 0 {
				continue
			}
			c.locale = v[0]

		case DSNTLSServerName:
			if len(v) == 0 {
				continue
			}
			if c.tlsConfig == nil {
				c.tlsConfig = &tls.Config{}
			}
			c.tlsConfig.ServerName = v[0]

		case DSNTLSInsecureSkipVerify:
			if len(v) == 0 {
				continue
			}
			var err error
			b := true
			if v[0] != "" {
				b, err = strconv.ParseBool(v[0])
				if err != nil {
					return nil, fmt.Errorf("failed to parse InsecureSkipVerify (bool): %s", v[0])
				}
			}
			if c.tlsConfig == nil {
				c.tlsConfig = &tls.Config{}
			}
			c.tlsConfig.InsecureSkipVerify = b

		case DSNTLSRootCAFile:
			for _, fn := range v {
				rootPEM, err := ioutil.ReadFile(fn)
				if err != nil {
					return nil, err
				}
				if certPool == nil {
					certPool = x509.NewCertPool()
				}
				if ok := certPool.AppendCertsFromPEM(rootPEM); !ok {
					return nil, fmt.Errorf("failed to parse root certificate - filename: %s", fn)
				}
			}
			if certPool != nil {
				if c.tlsConfig == nil {
					c.tlsConfig = &tls.Config{}
				}
				c.tlsConfig.RootCAs = certPool
			}
		}
	}
	return c, nil
}

// Host returns the host of the connector.
func (c *Connector) Host() string {
	return c.host
}

// Username returns the username of the connector.
func (c *Connector) Username() string {
	return c.username
}

// Password returns the password of the connector.
func (c *Connector) Password() string {
	return c.password
}

// Locale returns the locale of the connector.
func (c *Connector) Locale() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.locale
}

/*
SetLocale sets the locale of the connector.

For more information please see DSNLocale.
*/
func (c *Connector) SetLocale(locale string) {
	c.mu.Lock()
	c.locale = locale
	c.mu.Unlock()
}

// FetchSize returns the fetchSize of the connector.
func (c *Connector) FetchSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.fetchSize
}

/*
SetFetchSize sets the fetchSize of the connector.

For more information please see DSNFetchSize.
*/
func (c *Connector) SetFetchSize(fetchSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if fetchSize < minFetchSize {
		fetchSize = minFetchSize
	}
	c.fetchSize = fetchSize
	return nil
}

// Timeout returns the timeout of the connector.
func (c *Connector) Timeout() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeout
}

/*
SetTimeout sets the timeout of the connector.

For more information please see DSNTimeout.
*/
func (c *Connector) SetTimeout(timeout int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout < minTimeout {
		timeout = minTimeout
	}
	c.timeout = timeout
	return nil
}

// TLSConfig returns the TLS configuration of the connector.
func (c *Connector) TLSConfig() *tls.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tlsConfig
}

// SetTLSConfig sets the TLS configuration of the connector.
func (c *Connector) SetTLSConfig(tlsConfig *tls.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tlsConfig = tlsConfig
	return nil
}

// BasicAuthDSN return the connector DSN for basic authentication.
func (c *Connector) BasicAuthDSN() string {
	values := url.Values{}
	if c.locale != "" {
		values.Set(DSNLocale, c.locale)
	}
	if c.fetchSize != 0 {
		values.Set(DSNFetchSize, fmt.Sprintf("%d", c.fetchSize))
	}
	if c.timeout != 0 {
		values.Set(DSNTimeout, fmt.Sprintf("%d", c.timeout))
	}
	return (&url.URL{
		Scheme:   DriverName,
		User:     url.UserPassword(c.username, c.password),
		Host:     c.host,
		RawQuery: values.Encode(),
	}).String()
}

// Connect implements the database/sql/driver/Connector interface.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return newConn(ctx, c)
}

// Driver implements the database/sql/driver/Connector interface.
func (c *Connector) Driver() driver.Driver {
	return drv
}

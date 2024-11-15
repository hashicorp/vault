package driver

import (
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/protocol/auth"
)

type certKeyFiles struct {
	certFile, keyFile string
}

func newCertKeyFiles(certFile, keyFile string) *certKeyFiles {
	return &certKeyFiles{certFile: path.Clean(certFile), keyFile: path.Clean(keyFile)}
}

func (f *certKeyFiles) read() ([]byte, []byte, error) {
	cert, err := os.ReadFile(f.certFile)
	if err != nil {
		return nil, nil, err
	}
	key, err := os.ReadFile(f.keyFile)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

// authAttrs is holding authentication relevant attributes.
type authAttrs struct {
	hasCookie            atomic.Bool
	version              atomic.Uint64 // auth attributes version
	mu                   sync.RWMutex
	_username, _password string // basic authentication
	_certKeyFiles        *certKeyFiles
	_certKey             *auth.CertKey // X509
	_token               string        // JWT
	_logonname           string        // session cookie login does need logon name provided by JWT authentication.
	_sessionCookie       []byte        // authentication via session cookie (HDB currently does support only SAML and JWT - go-hdb JWT)
	_refreshPassword     func() (password string, ok bool)
	_refreshClientCert   func() (clientCert, clientKey []byte, ok bool)
	_refreshToken        func() (token string, ok bool)
	cbmu                 sync.Mutex // prevents refresh callbacks from being called in parallel
}

func isJWTToken(token string) bool { return strings.HasPrefix(token, "ey") }

/*
	keep c as the instance name, so that the generated help does have
	the same instance variable name when included in connector
*/

func (c *authAttrs) clone() *authAttrs {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &authAttrs{
		_username:          c._username,
		_password:          c._password,
		_certKey:           c._certKey,
		_token:             c._token,
		_refreshPassword:   c._refreshPassword,
		_refreshClientCert: c._refreshClientCert,
		_refreshToken:      c._refreshToken,
	}
}

func (c *authAttrs) cookieAuth() *p.AuthHnd {
	if !c.hasCookie.Load() { // fastpath without lock
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	auth := p.NewAuthHnd(c._logonname)                              // important: for session cookie auth we do need the logonname from JWT auth,
	auth.AddSessionCookie(c._sessionCookie, c._logonname, clientID) // and for HANA onPrem the final session cookie req needs the logonname as well.
	return auth
}

func (c *authAttrs) authHnd() *p.AuthHnd {
	c.mu.RLock()
	defer c.mu.RUnlock()

	authHnd := p.NewAuthHnd(c._username) // use username as logonname
	if c._certKey != nil {
		authHnd.AddX509(c._certKey)
	}
	if c._token != "" {
		authHnd.AddJWT(c._token)
	}
	// mimic standard drivers and use password as token if user is empty
	if c._token == "" && c._username == "" && isJWTToken(c._password) {
		authHnd.AddJWT(c._password)
	}
	if c._password != "" {
		authHnd.AddBasic(c._username, c._password)
	}
	return authHnd
}

func (c *authAttrs) callRefreshPasswordWithLock(refreshPassword func() (string, bool)) (string, bool) {
	defer c.mu.Lock() // finally lock attr again
	c.mu.Unlock()     // unlock attr, so that callback can call attr methods
	return refreshPassword()
}

func (c *authAttrs) callRefreshTokenWithLock(refreshToken func() (token string, ok bool)) (string, bool) {
	defer c.mu.Lock() // finally lock attr again
	c.mu.Unlock()     // unlock attr, so that callback can call attr methods
	return refreshToken()
}

func (c *authAttrs) callRefreshClientCertWithLock(refreshClientCert func() (clientCert, clientKey []byte, ok bool)) ([]byte, []byte, bool) {
	defer c.mu.Lock() // finally lock attr again
	c.mu.Unlock()     // unlock attr, so that callback can call attr methods
	return refreshClientCert()
}

func (c *authAttrs) refresh() error {
	c.cbmu.Lock() // synchronize refresh calls
	defer c.cbmu.Unlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c._refreshPassword != nil {
		if password, ok := c.callRefreshPasswordWithLock(c._refreshPassword); ok {
			if password != c._password {
				c._password = password
				c.version.Add(1)
			}
		}
	}
	if c._refreshToken != nil {
		if token, ok := c.callRefreshTokenWithLock(c._refreshToken); ok {
			if token != c._token {
				c._token = token
				c.version.Add(1)
			}
		}
	}
	switch {
	case c._certKeyFiles != nil && c._refreshClientCert == nil:
		if clientCert, clientKey, err := c._certKeyFiles.read(); err != nil {
			if c._certKey == nil || !c._certKey.Equal(clientCert, clientKey) {
				certKey, err := auth.NewCertKey(clientCert, clientKey)
				if err != nil {
					return err
				}
				c._certKey = certKey
				c.version.Add(1)
			}
		}
	case c._refreshClientCert != nil:
		if clientCert, clientKey, ok := c.callRefreshClientCertWithLock(c._refreshClientCert); ok {
			if c._certKey == nil || !c._certKey.Equal(clientCert, clientKey) {
				certKey, err := auth.NewCertKey(clientCert, clientKey)
				if err != nil {
					return err
				}
				c._certKey = certKey
				c.version.Add(1)
			}
		}
	}
	return nil
}

func (c *authAttrs) invalidateCookie() { c.hasCookie.Store(false) }

func (c *authAttrs) setCookie(logonname string, sessionCookie []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hasCookie.Store(true)
	c._logonname = logonname
	c._sessionCookie = sessionCookie
}

// Username returns the username of the connector.
func (c *authAttrs) Username() string { c.mu.RLock(); defer c.mu.RUnlock(); return c._username }

// Password returns the basic authentication password of the connector.
func (c *authAttrs) Password() string { c.mu.RLock(); defer c.mu.RUnlock(); return c._password }

// SetPassword sets the basic authentication password of the connector.
func (c *authAttrs) SetPassword(password string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._password = password
}

// RefreshPassword returns the callback function for basic authentication password refresh.
func (c *authAttrs) RefreshPassword() func() (password string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._refreshPassword
}

// SetRefreshPassword sets the callback function for basic authentication password refresh.
// The callback function might be called simultaneously from multiple goroutines only if registered
// for more than one Connector.
func (c *authAttrs) SetRefreshPassword(refreshPassword func() (password string, ok bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._refreshPassword = refreshPassword
}

// ClientCert returns the X509 authentication client certificate and key of the connector.
func (c *authAttrs) ClientCert() (clientCert, clientKey []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._certKey.Cert(), c._certKey.Key()
}

// RefreshClientCert returns the callback function for X509 authentication client certificate and key refresh.
func (c *authAttrs) RefreshClientCert() func() (clientCert, clientKey []byte, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._refreshClientCert
}

// SetRefreshClientCert sets the callback function for X509 authentication client certificate and key refresh.
// The callback function might be called simultaneously from multiple goroutines only if registered
// for more than one Connector.
func (c *authAttrs) SetRefreshClientCert(refreshClientCert func() (clientCert, clientKey []byte, ok bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._refreshClientCert = refreshClientCert
}

// Token returns the JWT authentication token of the connector.
func (c *authAttrs) Token() string { c.mu.RLock(); defer c.mu.RUnlock(); return c._token }

// RefreshToken returns the callback function for JWT authentication token refresh.
func (c *authAttrs) RefreshToken() func() (token string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._refreshToken
}

// SetRefreshToken sets the callback function for JWT authentication token refresh.
// The callback function might be called simultaneously from multiple goroutines only if registered
// for more than one Connector.
func (c *authAttrs) SetRefreshToken(refreshToken func() (token string, ok bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._refreshToken = refreshToken
}

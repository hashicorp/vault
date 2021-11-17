package gocbcore

import "crypto/tls"

// UserPassPair represents a username and password pair.
type UserPassPair struct {
	Username string
	Password string
}

// AuthCredsRequest represents an authentication details request from the agent.
type AuthCredsRequest struct {
	Service  ServiceType
	Endpoint string
}

// AuthCertRequest represents a certificate details request from the agent.
type AuthCertRequest struct {
	Service  ServiceType
	Endpoint string
}

// AuthProvider is an interface to allow the agent to fetch authentication
// credentials on-demand from the application.
type AuthProvider interface {
	SupportsTLS() bool
	SupportsNonTLS() bool
	Certificate(req AuthCertRequest) (*tls.Certificate, error)
	Credentials(req AuthCredsRequest) ([]UserPassPair, error)
}

func getSingleAuthCreds(auth AuthProvider, req AuthCredsRequest) (UserPassPair, error) {
	creds, err := auth.Credentials(req)
	if err != nil {
		return UserPassPair{}, err
	}

	if len(creds) != 1 {
		return UserPassPair{}, errInvalidCredentials
	}

	return creds[0], nil
}

func getKvAuthCreds(auth AuthProvider, endpoint string) (UserPassPair, error) {
	return getSingleAuthCreds(auth, AuthCredsRequest{
		Service:  MemdService,
		Endpoint: endpoint,
	})
}

// PasswordAuthProvider provides a standard AuthProvider implementation
// for use with a standard username/password pair (for example, RBAC).
type PasswordAuthProvider struct {
	Username string
	Password string
}

// SupportsNonTLS specifies whether this authenticator supports non-TLS connections.
func (auth PasswordAuthProvider) SupportsNonTLS() bool {
	return true
}

// SupportsTLS specifies whether this authenticator supports TLS connections.
func (auth PasswordAuthProvider) SupportsTLS() bool {
	return true
}

// Certificate directly returns a certificate chain to present for the connection.
func (auth PasswordAuthProvider) Certificate(req AuthCertRequest) (*tls.Certificate, error) {
	return nil, nil
}

// Credentials directly returns the username/password from the provider.
func (auth PasswordAuthProvider) Credentials(req AuthCredsRequest) ([]UserPassPair, error) {
	return []UserPassPair{{
		Username: auth.Username,
		Password: auth.Password,
	}}, nil
}

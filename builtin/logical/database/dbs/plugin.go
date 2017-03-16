package dbs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/rpc"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-plugin"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type DatabasePlugin struct {
	impl DatabaseType
}

func (d DatabasePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &databasePluginRPCServer{impl: d.impl}, nil
}

func (DatabasePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &databasePluginRPCClient{client: c}, nil
}

type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	*databasePluginRPCClient
}

func (dc *DatabasePluginClient) Close() error {
	err := dc.databasePluginRPCClient.Close()
	dc.client.Kill()

	return err
}

func generateX509Cert() ([]byte, *x509.Certificate, *ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		//	c.logger.Error("core: failed to generate replicated cluster signing key", "error", err)
		return nil, nil, nil, err
	}

	//c.logger.Trace("core: generating replicated cluster certificate")

	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, nil, err
	}
	host = "localhost"
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		// 30 years of single-active uptime ought to be enough for anybody
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		//	c.logger.Error("core: error generating self-signed cert for replication", "error", err)
		return nil, nil, nil, fmt.Errorf("unable to generate replicated cluster certificate: %v", err)
	}

	caCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		//	c.logger.Error("core: error parsing replicated self-signed cert", "error", err)
		return nil, nil, nil, fmt.Errorf("error parsing generated replication certificate: %v", err)
	}

	return certBytes, caCert, key, nil
}

func generateClientCert(CACert *x509.Certificate, CAKey *ecdsa.PrivateKey) ([]byte, *x509.Certificate, []byte, error) {
	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, nil, err
	}
	host = "localhost"
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	clientKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, nil, errwrap.Wrapf("error generating client key: {{err}}", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, CACert, clientKey.Public(), CAKey)
	if err != nil {
		return nil, nil, nil, errwrap.Wrapf("unable to generate client certificate: {{err}}", err)
	}

	clientCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		//	c.logger.Error("core: error parsing replicated self-signed cert", "error", err)
		return nil, nil, nil, fmt.Errorf("error parsing generated replication certificate: %v", err)
	}

	keyBytes, err := x509.MarshalECPrivateKey(clientKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return certBytes, clientCert, keyBytes, nil
}

func newPluginClient(sys logical.SystemView, command, checksum string) (DatabaseType, error) {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": new(DatabasePlugin),
	}

	CACertBytes, CACert, CAKey, err := generateX509Cert()
	if err != nil {
		return nil, err
	}

	clientCertBytes, clientCert, clientKey, err := generateClientCert(CACert, CAKey)
	if err != nil {
		return nil, err
	}

	/*	serverCert, serverKey, err := generateClientCert(CACert, CAKey)
		if err != nil {
			return nil, err
		}*/
	serverKey, err := x509.MarshalECPrivateKey(CAKey)
	if err != nil {
		return nil, err
	}
	cert := tls.Certificate{
		Certificate: [][]byte{clientCertBytes, CACertBytes},
		PrivateKey:  clientKey,
		Leaf:        clientCert,
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AddCert(CACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
		ClientCAs:    clientCertPool,
		ServerName:   CACert.Subject.CommonName,
		MinVersion:   tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	wrapToken, err := sys.ResponseWrapData(map[string]interface{}{
		"CACert":     CACertBytes,
		"ServerCert": CACertBytes,
		"ServerKey":  serverKey,
	}, time.Second*10, true)

	cmd := exec.Command(command)
	cmd.Env = append(cmd.Env, fmt.Sprintf("VAULT_WRAP_TOKEN=%s", wrapToken))

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		TLSConfig:       tlsConfig,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	databaseRPC := raw.(*databasePluginRPCClient)

	return &DatabasePluginClient{
		client:                  client,
		databasePluginRPCClient: databaseRPC,
	}, nil
}

func NewPluginServer(db DatabaseType) {
	dbPlugin := &DatabasePlugin{
		impl: db,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": dbPlugin,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     VaultPluginTLSProvider,
	})
}

func VaultPluginTLSProvider() (*tls.Config, error) {
	unwrapToken := os.Getenv("VAULT_WRAP_TOKEN")
	if strings.Count(unwrapToken, ".") != 2 {
		return nil, errors.New("Could not parse unwraptoken")
	}

	wt, err := jws.ParseJWT([]byte(unwrapToken))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error decoding token: %s", err))
	}
	if wt == nil {
		return nil, errors.New("nil decoded token")
	}

	addrRaw := wt.Claims().Get("addr")
	if addrRaw == nil {
		return nil, errors.New("decoded token does not contain primary cluster address")
	}
	vaultAddr, ok := addrRaw.(string)
	if !ok {
		return nil, errors.New("decoded token's address not valid")
	}
	if vaultAddr == "" {
		return nil, errors.New(`no address for the vault found`)
	}

	// Sanity check the value
	if _, err := url.Parse(vaultAddr); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the vault address: %s", err))
	}

	clientConf := api.DefaultConfig()
	clientConf.Address = vaultAddr
	client, err := api.NewClient(clientConf)
	if err != nil {
		return nil, errwrap.Wrapf("error during token unwrap request: {{err}}", err)
	}

	secret, err := client.Logical().Unwrap(unwrapToken)
	if err != nil {
		return nil, errwrap.Wrapf("error during token unwrap request: {{err}}", err)
	}

	CABytesRaw, ok := secret.Data["CACert"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	CABytes, err := base64.StdEncoding.DecodeString(CABytesRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	CACert, err := x509.ParseCertificate(CABytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	serverCertBytesRaw, ok := secret.Data["ServerCert"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	serverCertBytes, err := base64.StdEncoding.DecodeString(serverCertBytesRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	serverCert, err := x509.ParseCertificate(serverCertBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	serverKeyRaw, ok := secret.Data["ServerKey"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	serverKey, err := base64.StdEncoding.DecodeString(serverKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(CACert)

	cert := tls.Certificate{
		Certificate: [][]byte{serverCertBytes},
		PrivateKey:  serverKey,
		Leaf:        serverCert,
	}

	// Setup TLS config
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		RootCAs:    caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		// TLS 1.2 minimum
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// ---- RPC client domain ----

type databasePluginRPCClient struct {
	client *rpc.Client
}

func (dr *databasePluginRPCClient) Type() string {
	return "plugin"
}

func (dr *databasePluginRPCClient) CreateUser(statements Statements, username, password, expiration string) error {
	req := CreateUserRequest{
		Statements: statements,
		Username:   username,
		Password:   password,
		Expiration: expiration,
	}

	err := dr.client.Call("Plugin.CreateUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) RenewUser(statements Statements, username, expiration string) error {
	req := RenewUserRequest{
		Statements: statements,
		Username:   username,
		Expiration: expiration,
	}

	err := dr.client.Call("Plugin.RenewUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) RevokeUser(statements Statements, username string) error {
	req := RevokeUserRequest{
		Statements: statements,
		Username:   username,
	}

	err := dr.client.Call("Plugin.RevokeUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Initialize(conf map[string]interface{}) error {
	err := dr.client.Call("Plugin.Initialize", conf, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Close() error {
	err := dr.client.Call("Plugin.Close", struct{}{}, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) GenerateUsername(displayName string) (string, error) {
	var username string
	err := dr.client.Call("Plugin.GenerateUsername", displayName, &username)

	return username, err
}

func (dr *databasePluginRPCClient) GeneratePassword() (string, error) {
	var password string
	err := dr.client.Call("Plugin.GeneratePassword", struct{}{}, &password)

	return password, err
}

func (dr *databasePluginRPCClient) GenerateExpiration(duration time.Duration) (string, error) {
	var expiration string
	err := dr.client.Call("Plugin.GenerateExpiration", duration, &expiration)

	return expiration, err
}

// ---- RPC server domain ----
type databasePluginRPCServer struct {
	impl DatabaseType
}

func (ds *databasePluginRPCServer) Type(_ struct{}, resp *string) error {
	*resp = ds.impl.Type()
	return nil
}

func (ds *databasePluginRPCServer) CreateUser(args *CreateUserRequest, _ *struct{}) error {
	err := ds.impl.CreateUser(args.Statements, args.Username, args.Password, args.Expiration)

	return err
}

func (ds *databasePluginRPCServer) RenewUser(args *RenewUserRequest, _ *struct{}) error {
	err := ds.impl.RenewUser(args.Statements, args.Username, args.Expiration)

	return err
}

func (ds *databasePluginRPCServer) RevokeUser(args *RevokeUserRequest, _ *struct{}) error {
	err := ds.impl.RevokeUser(args.Statements, args.Username)

	return err
}

func (ds *databasePluginRPCServer) Initialize(args map[string]interface{}, _ *struct{}) error {
	err := ds.impl.Initialize(args)

	return err
}

func (ds *databasePluginRPCServer) Close(_ interface{}, _ *struct{}) error {
	ds.impl.Close()
	return nil
}

func (ds *databasePluginRPCServer) GenerateUsername(args string, resp *string) error {
	var err error
	*resp, err = ds.impl.GenerateUsername(args)

	return err
}

func (ds *databasePluginRPCServer) GeneratePassword(_ struct{}, resp *string) error {
	var err error
	*resp, err = ds.impl.GeneratePassword()

	return err
}

func (ds *databasePluginRPCServer) GenerateExpiration(args time.Duration, resp *string) error {
	var err error
	*resp, err = ds.impl.GenerateExpiration(args)

	return err
}

// ---- Request Args domain ----

type CreateUserRequest struct {
	Statements Statements
	Username   string
	Password   string
	Expiration string
}

type RenewUserRequest struct {
	Statements Statements
	Username   string
	Expiration string
}

type RevokeUserRequest struct {
	Statements Statements
	Username   string
}

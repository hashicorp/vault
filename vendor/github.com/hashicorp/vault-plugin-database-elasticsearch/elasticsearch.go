// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
)

const (
	defaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 15) (.RoleName | truncate 15) (random 20) (unix_time) | truncate 100 }}`
)

var _ dbplugin.Database = (*Elasticsearch)(nil)

// New returns a new Elasticsearch instance
func New() (interface{}, error) {
	db := &Elasticsearch{}
	return dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues), nil
}

// Elasticsearch implements dbplugin's Database interface.
type Elasticsearch struct {
	// This protects the config from races while also allowing multiple threads
	// to read the config simultaneously when it's not changing.
	mux sync.RWMutex

	// The root credential config.
	config map[string]interface{}

	usernameProducer template.StringTemplate
}

// Type returns the TypeName for this backend
func (es *Elasticsearch) Type() (string, error) {
	return "elasticsearch", nil
}

// SecretValues is used by some error-sanitizing middleware in Vault that basically
// replaces the keys in the map with the values given so they're not leaked via
// error messages.
func (es *Elasticsearch) SecretValues() map[string]string {
	es.mux.RLock()
	defer es.mux.RUnlock()

	replacements := make(map[string]string)
	for _, secretKey := range []string{"password", "client_key"} {
		vIfc, found := es.config[secretKey]
		if !found {
			continue
		}
		secretVal, ok := vIfc.(string)
		if !ok {
			continue
		}
		// So, supposing a password of "0pen5e5ame",
		// this will cause that string to get replaced with "[password]".
		replacements[secretVal] = "[" + secretKey + "]"
	}
	return replacements
}

// Initialize is called on `$ vault write database/config/:db-name`,
// or when you do a creds call after Vault's been restarted.
func (es *Elasticsearch) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	usernameTemplate, err := strutil.GetString(req.Config, "username_template")
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to retrieve username_template: %w", err)
	}
	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}

	up, err := template.NewTemplate(template.Template(usernameTemplate))
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to initialize username template: %w", err)
	}

	_, err = up.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	// Validate the config to provide immediate feedback to the user.
	// Ensure required string fields are provided in the expected format.
	for _, requiredField := range []string{"username", "password", "url"} {
		raw, ok := req.Config[requiredField]
		if !ok {
			return dbplugin.InitializeResponse{}, fmt.Errorf(`%q must be provided`, requiredField)
		}
		if _, ok := raw.(string); !ok {
			return dbplugin.InitializeResponse{}, fmt.Errorf(`%q must be a string`, requiredField)
		}
	}

	// Ensure optional string fields are provided in the expected format.
	for _, optionalField := range []string{"ca_cert", "ca_path", "client_cert", "client_key", "tls_server_name"} {
		raw, ok := req.Config[optionalField]
		if !ok {
			continue
		}
		if _, ok = raw.(string); !ok {
			return dbplugin.InitializeResponse{}, fmt.Errorf(`%q must be a string`, optionalField)
		}
	}

	// Ensure optional bool fields are provided in the expected format.
	for _, optionalBool := range []string{"insecure", "use_old_xpack"} {
		raw, ok := req.Config[optionalBool]
		if !ok {
			continue
		}

		boolVal, err := parseBool(raw)
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf(`%q must be a bool`, optionalBool)
		}
		req.Config[optionalBool] = boolVal
	}

	// Test the given config to see if we can make a client.
	client, err := buildClient(req.Config)
	if err != nil {
		return dbplugin.InitializeResponse{}, errwrap.Wrapf("couldn't make client with inbound config: {{err}}", err)
	}

	// Optionally, test the given config to see if we can make a successful call.
	if req.VerifyConnection {
		// Whether this role is found or unfound, if we're configured correctly there will
		// be no err from the call. However, if something is misconfigured, this will yield
		// an error response, which will be described in the returned error.
		if _, err := client.GetRole(ctx, "vault-test"); err != nil {
			return dbplugin.InitializeResponse{}, errwrap.Wrapf("client test of getting a role failed: {{err}}", err)
		}
	}

	// Everything's working, write the new config to memory and storage.
	es.mux.Lock()
	defer es.mux.Unlock()
	es.usernameProducer = up
	es.config = req.Config
	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

// NewUser is called on `$ vault read database/creds/:role-name`
// and it's the first time anything is touched from `$ vault write database/roles/:role-name`.
// This is likely to be the highest-throughput method for this plugin.
func (es *Elasticsearch) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	username, err := es.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	stmt, err := newCreationStatement(req.Statements)
	if err != nil {
		return dbplugin.NewUserResponse{}, errwrap.Wrapf("unable to read creation_statements: {{err}}", err)
	}

	user := &User{
		Password: req.Password,
		Roles:    stmt.PreexistingRoles,
	}

	// Don't let anyone write the config while we're using it for our current client.
	es.mux.RLock()
	defer es.mux.RUnlock()

	client, err := buildClient(es.config)
	if err != nil {
		return dbplugin.NewUserResponse{}, errwrap.Wrapf("unable to get client: {{err}}", err)
	}

	// If the RoleToCreate map has been populated with any data, we have one role to create.
	// There can either be one RoleToCreate and no PreexistingRoles, or >= 1 PreexistingRoles
	// and no RoleToCreate. They're mutually exclusive.
	if len(stmt.RoleToCreate) > 0 {
		// We'll simply name the role the same thing as the username, making it easy to tie back to this user.
		if err := client.CreateRole(ctx, username, stmt.RoleToCreate); err != nil {
			return dbplugin.NewUserResponse{}, errwrap.Wrapf(fmt.Sprintf("unable to create role name %s, role definition %q: {{err}}", username, stmt.RoleToCreate), err)
		}
		user.Roles = []string{username}
	}
	if err := client.CreateUser(ctx, username, user); err != nil {
		return dbplugin.NewUserResponse{}, errwrap.Wrapf(fmt.Sprintf("unable to create user name %s, user %q: {{err}}", username, user), err)
	}
	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

// DeleteUser is used to delete users from elasticsearch
func (es *Elasticsearch) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	// Don't let anyone write the config while we're using it for our current client.
	es.mux.RLock()
	defer es.mux.RUnlock()

	client, err := buildClient(es.config)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, errwrap.Wrapf("unable to get client: {{err}}", err)
	}

	var errs *multierror.Error
	// If the role already doesn't exist, either it wasn't created for this
	// user, or it was successfully deleted on a previous attempt to run this
	// code, there will be no error, so it's harmless to try.
	if err := client.DeleteRole(ctx, req.Username); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("unable to delete role name %s: %w", req.Username, err))
	}

	// Same with the user. If it was already deleted on a previous attempt, there won't be an
	// error.
	if err := client.DeleteUser(ctx, req.Username); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("unable to delete user name %s: %w", req.Username, err))
	}
	return dbplugin.DeleteUserResponse{}, errs.ErrorOrNil()
}

// UpdateUser doesn't require any statements from the user because it's not configurable in any
// way. We simply generate a new password and hit a pre-defined Elasticsearch REST API to rotate them.
func (es *Elasticsearch) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	// Don't let anyone read or write the config while we're in the process of rotating the password.
	es.mux.Lock()
	defer es.mux.Unlock()

	client, err := buildClient(es.config)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("unable to get client: %w", err)
	}

	if req.Password != nil {
		if err := client.ChangePassword(ctx, req.Username, req.Password.NewPassword); err != nil {
			return dbplugin.UpdateUserResponse{}, fmt.Errorf("unable to change password: %w", err)
		}
		// Note: changing the expiration of a user is a no-op for Elasticsearch,
		// and therefore ignored here
	}

	return dbplugin.UpdateUserResponse{}, nil
}

// Close for Elasticsearch is a NOOP, nothing to close
func (es *Elasticsearch) Close() error {
	return nil
}

func newCreationStatement(statements dbplugin.Statements) (*creationStatement, error) {
	if len(statements.Commands) == 0 {
		return nil, dbutil.ErrEmptyCreationStatement
	}
	if len(statements.Commands) > 1 {
		return nil, fmt.Errorf("only 1 creation statement supported for creation")
	}
	stmt := &creationStatement{}
	if err := json.Unmarshal([]byte(statements.Commands[0]), stmt); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %s: %w", []byte(statements.Commands[0]), err)
	}
	if len(stmt.PreexistingRoles) > 0 && len(stmt.RoleToCreate) > 0 {
		return nil, errors.New(`"elasticsearch_roles" and "elasticsearch_role_definition" are mutually exclusive`)
	}
	return stmt, nil
}

type creationStatement struct {
	PreexistingRoles []string               `json:"elasticsearch_roles"`
	RoleToCreate     map[string]interface{} `json:"elasticsearch_role_definition"`
}

// buildClient is a helper method for building a client from the present config,
// which is done often.
func buildClient(config map[string]interface{}) (*Client, error) {
	// We can presume these required fields are provided by strings
	// because they're validated in Init.
	clientConfig := &ClientConfig{
		Username: config["username"].(string),
		Password: config["password"].(string),
		BaseURL:  config["url"].(string),
	}

	hasTLSConf := false
	tlsConf := &TLSConfig{}

	// We can presume that if these are provided, they're in the expected format
	// because they're also validated in Init.
	if raw, ok := config["ca_cert"]; ok {
		tlsConf.CACert = raw.(string)
		hasTLSConf = true
	}
	if raw, ok := config["ca_path"]; ok {
		tlsConf.CAPath = raw.(string)
		hasTLSConf = true
	}
	if raw, ok := config["client_cert"]; ok {
		tlsConf.ClientCert = raw.(string)
		hasTLSConf = true
	}
	if raw, ok := config["client_key"]; ok {
		tlsConf.ClientKey = raw.(string)
		hasTLSConf = true
	}
	if raw, ok := config["tls_server_name"]; ok {
		tlsConf.TLSServerName = raw.(string)
		hasTLSConf = true
	}
	if raw, ok := config["insecure"]; ok {
		tlsConf.Insecure = raw.(bool)
		hasTLSConf = true
	}
	if raw, ok := config["use_old_xpack"]; ok {
		clientConfig.UseOldSecurityPath = raw.(bool)
	}

	// We should only fulfill the clientConfig's TLSConfig pointer if we actually
	// want the client to use TLS.
	if hasTLSConf {
		clientConfig.TLSConfig = tlsConf
	}

	client, err := NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func parseBool(rawVal interface{}) (bool, error) {
	switch val := rawVal.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	}

	return false, fmt.Errorf("could not parse type")
}

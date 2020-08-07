package couchbase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
)

const (
	couchbaseTypeName        = "couchbase"
	defaultCouchbaseUserRole = `{"Roles": [{"role":"ro_admin"}]}`
)

var (
	_ dbplugin.Database = &CouchbaseDB{}
)

// Type that combines the custom plugins Couchbase database connection configuration options and the Vault CredentialsProducer
// used for generating user information for the Couchbase database.
type CouchbaseDB struct {
	*couchbaseDBConnectionProducer
	credsutil.CredentialsProducer
}

// Type that combines the Couchbase Roles and Groups representing specific account permissions. Used to pass roles and or
// groups between the Vault server and the custom plugin in the dbplugin.Statements
type RolesAndGroups struct {
	Roles  []gocb.Role `json:"roles"`
	Groups []string    `json:"groups"`
}

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *CouchbaseDB {
	connProducer := &couchbaseDBConnectionProducer{}
	connProducer.Type = couchbaseTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 50,
		RoleNameLen:    50,
		UsernameLen:    50,
		Separator:      "-",
	}

	db := &CouchbaseDB{
		couchbaseDBConnectionProducer: connProducer,
		CredentialsProducer:           credsProducer,
	}

	return db
}

// Run instantiates a CouchbaseDB object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

func (c *CouchbaseDB) Type() (string, error) {
	return couchbaseTypeName, nil
}

func computeTimeout(ctx context.Context) (timeout time.Duration) {
	deadline, ok := ctx.Deadline()
	if ok {
		return time.Until(deadline)
	}
	return 5 * time.Second
}
func (c *CouchbaseDB) getConnection(ctx context.Context) (*gocb.Cluster, error) {
	db, err := c.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return db.(*gocb.Cluster), nil
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (c *CouchbaseDB) SetCredentials(ctx context.Context, _ dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	username = staticUser.Username
	password = staticUser.Password
	if username == "" || password == "" {
		return "", "", errors.New("must provide both username and password")
	}

	// Grab the lock
	c.Lock()
	defer c.Unlock()

	// Get the connection
	db, err := c.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	// Close the database connection to ensure no new connections come in
	defer func() {
		if err := c.close(); err != nil {
			logger := hclog.New(&hclog.LoggerOptions{})
			logger.Error("defer close failed", "error", err)
		}
	}()

	// Get the UserManager

	mgr := db.Users()

	// Get the User and error out if it does not exist.

	userOpts, err := mgr.GetUser(username, nil)
	if err != nil {
		return "", "", err
	}

	user := gocb.User{
		Username:    username,
		Password:    password,
		Roles:       userOpts.Roles,
		Groups:      userOpts.Groups,
		DisplayName: userOpts.DisplayName,
	}

	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout:    computeTimeout(ctx),
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func (c *CouchbaseDB) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, _ time.Time) (username string, password string, err error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Creation) == 0 {
		statements.Creation = append(statements.Creation, defaultCouchbaseUserRole)
	}

	jsonRoleAndGroupData := []byte(statements.Creation[0])

	var rag RolesAndGroups

	err = json.Unmarshal(jsonRoleAndGroupData, &rag)
	if err != nil {
		return "", "", errwrap.Wrapf("error unmarshaling roles and groups creation statement JSON: {{err}}", err)
	}

	username, err = c.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	password, err = c.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	// Get the connection
	db, err := c.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	// Close the database connection to ensure no new connections come in
	defer func() {
		if err := c.close(); err != nil {
			logger := hclog.New(&hclog.LoggerOptions{})
			logger.Error("defer close failed", "error", err)
		}
	}()

	// Get the UserManager

	mgr := db.Users()

	user := gocb.User{
		Username:    username,
		DisplayName: usernameConfig.DisplayName,
		Password:    password,
		Roles:       rag.Roles,
		Groups:      rag.Groups,
	}

	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout:    computeTimeout(ctx),
			DomainName: "local",
		})
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

// RenewUser is not supported by Couchbase, so this is a no-op.
func (p *CouchbaseDB) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP

	return nil
}

func (c *CouchbaseDB) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	db, err := c.getConnection(ctx)
	if err != nil {
		return err
	}

	// Close the database connection to ensure no new connections come in
	defer func() {
		if err := c.close(); err != nil {
			logger := hclog.New(&hclog.LoggerOptions{})
			logger.Error("defer close failed", "error", err)
		}
	}()

	// Get the UserManager
	mgr := db.Users()

	err = mgr.DropUser(username, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *CouchbaseDB) RotateRootCredentials(ctx context.Context, _ []string) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if len(c.Username) == 0 || len(c.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	password, err := c.GeneratePassword()
	if err != nil {
		return nil, err
	}

	db, err := c.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	// Close the database connection to ensure no new connections come in
	defer func() {
		if err := c.close(); err != nil {
			logger := hclog.New(&hclog.LoggerOptions{})
			logger.Error("defer close failed", "error", err)
		}
	}()

	// Get the UserManager

	mgr := db.Users()

	// Get the User

	userOpts, err := mgr.GetUser(c.Username, nil)
	if err != nil {
		return nil, err
	}

	user := gocb.User{
		Username:    c.Username,
		Password:    password,
		Roles:       userOpts.Roles,
		Groups:      userOpts.Groups,
		DisplayName: userOpts.DisplayName,
	}

	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout:    computeTimeout(ctx),
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return nil, err
	}

	c.rawConfig["password"] = password

	return c.rawConfig, nil
}

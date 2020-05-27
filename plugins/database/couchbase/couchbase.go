package couchbase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
)

var _ = fmt.Printf

const (
	couchbaseTypeName        = "couchbase"
	defaultCouchbaseUserRole = `[{"name":"ro_admin"}]`
)

var (
	_ dbplugin.Database = &CouchbaseDB{}
)

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

type CouchbaseDB struct {
	*couchbaseDBConnectionProducer
	credsutil.CredentialsProducer
}

func (c *CouchbaseDB) Type() (string, error) {
	return couchbaseTypeName, nil
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
func (c *CouchbaseDB) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	if len(statements.Rotation) == 0 {
		statements.Rotation = []string{}
	}

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
			Timeout:    1 * time.Second,
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return "", "", err
	}

	// Close the database connection to ensure no new connections come in
	if err := c.close(); err != nil {
		return "", "", err
	}

	return username, password, nil
}

func (c *CouchbaseDB) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Creation) == 0 {
		statements.Creation[0] = defaultCouchbaseUserRole
		//return "", "", dbutil.ErrEmptyCreationStatement
	}

	jsonRoleData := []byte(statements.Creation[0])
	var v []interface{}
	var roles []gocb.Role
	err = json.Unmarshal(jsonRoleData, &v)
	if err != nil {
		return "", "", errwrap.Wrapf("error unmarshaling JSON: {{err}}", err)
	}

	err = mapstructure.Decode(v, &roles)
	if err != nil {
		return "", "", errwrap.Wrapf("error mapping roles: {{err}}", err)
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

	// Get the UserManager

	mgr := db.Users()

	user := gocb.User{
		Username:    username,
		DisplayName: usernameConfig.DisplayName,
		Password:    password,
		Roles:       roles,
		Groups:      []string{""},
	}

	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout:    1 * time.Second,
			DomainName: "local",
		})
	if err != nil {
		return "", "", err
	}

	c.close()

	return username, password, nil
}

// RenewUser is not supported by Couchbase, so this is a no-op.
func (p *CouchbaseDB) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP

	return nil
}

func (p *CouchbaseDB) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	p.Lock()
	defer p.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Revocation) == 0 {
		return p.defaultRevokeUser(ctx, username)
	}

	return p.customRevokeUser(ctx, username, statements.Revocation)
}

func (c *CouchbaseDB) customRevokeUser(ctx context.Context, username string, revocationStmts []string) error {
	db, err := c.getConnection(ctx)
	if err != nil {
		return err
	}

	// Get the UserManager
	mgr := db.Users()

	err = mgr.DropUser(username, nil)

	if err != nil {
		panic(err)
	}

	c.close()

	return nil
}

func (c *CouchbaseDB) defaultRevokeUser(ctx context.Context, username string) error {
	_, err := c.getConnection(ctx)
	if err != nil {
		return err
	}
	c.close()

	return nil
}

func (c *CouchbaseDB) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if len(c.Username) == 0 || len(c.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatements := statements
	if len(rotateStatements) == 0 {
		rotateStatements = []string{""}
	}

	password, err := c.GeneratePassword()
	if err != nil {
		return nil, err
	}

	db, err := c.getConnection(ctx)
	if err != nil {
		return nil, err
	}

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
			Timeout:    1 * time.Second,
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return nil, err
	}

	// Close the database connection to ensure no new connections come in
	if err := c.close(); err != nil {
		return nil, err
	}

	c.rawConfig["password"] = password

	return c.rawConfig, nil
}

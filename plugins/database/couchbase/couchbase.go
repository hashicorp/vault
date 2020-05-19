package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"encoding/json"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
	"github.com/couchbase/gocb/v2"
)

var _ = fmt.Printf

const (
	couchbaseTypeName      = "couchbase"
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
		CredentialsProducer:   credsProducer,
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
		Username: username,
		Password: password,
		Roles: userOpts.Roles,
		Groups: userOpts.Groups,
		DisplayName: userOpts.DisplayName,
	}
		
	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout: 1 * time.Second,
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return "", "", err
	}

	// Close the database connection to ensure no new connections come in
	if err := c.Close(); err != nil {
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
		Username: username,
		DisplayName: usernameConfig.DisplayName,
		Password: password,
		Roles: roles,
		Groups: []string{""},
	}
		
	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout: 1 * time.Second,
			DomainName: "local",
		})
	if err != nil {
		return "", "", err
	}
	
	//	c.Close()

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

func (p *CouchbaseDB) customRevokeUser(ctx context.Context, username string, revocationStmts []string) error {
	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	// Get the UserManager
	mgr := db.Users()

	err = mgr.DropUser(username, nil)

	if err != nil {
		panic(err)
	}

	db.Close(&gocb.ClusterCloseOptions{})
	p.cluster = nil
	
	return nil

/*	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	for _, stmt := range revocationStmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":     username,
				"username": username,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return tx.Commit()*/ return err
}

func (p *CouchbaseDB) defaultRevokeUser(ctx context.Context, username string) error {
	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}
	db.Close(&gocb.ClusterCloseOptions{})

	// Check if the role exists
/*	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !exists {
		return nil
	}

	// Query for permissions; we need to revoke permissions before we can drop
	// the role
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err := db.PrepareContext(ctx, "SELECT DISTINCT table_schema FROM information_schema.role_column_grants WHERE grantee=$1;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return err
	}
	defer rows.Close()

	const initialNumRevocations = 16
	revocationStmts := make([]string, 0, initialNumRevocations)
	for rows.Next() {
		var schema string
		err = rows.Scan(&schema)
		if err != nil {
			// keep going; remove as many permissions as possible right now
			continue
		}
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA %s FROM %s;`,
			schema,
			username))

		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE USAGE ON SCHEMA %s FROM %s;`,
			schema,
			username))
	}

	// for good measure, revoke all privileges and usage on schema public
	revocationStmts = append(revocationStmts, fmt.Sprintf(
		`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM %s;`,
		username))

	revocationStmts = append(revocationStmts, fmt.Sprintf(
		"REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM %s;",
		username))

	revocationStmts = append(revocationStmts, fmt.Sprintf(
		"REVOKE USAGE ON SCHEMA public FROM %s;",
		username))

	// get the current database name so we can issue a REVOKE CONNECT for
	// this username
	var dbname sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbname); err != nil {
		return err
	}

	if dbname.Valid {
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE CONNECT ON DATABASE %s FROM %s;`,
			dbname.String,
			username))
	}

	// again, here, we do not stop on error, as we want to remove as
	// many permissions as possible right now
	var lastStmtError error
	for _, query := range revocationStmts {
		if err := dbtxn.ExecuteDBQuery(ctx, db, nil, query); err != nil {
			lastStmtError = err
		}
	}

	// can't drop if not all privileges are revoked
	if rows.Err() != nil {
		return errwrap.Wrapf("could not generate revocation statements for all rows: {{err}}", rows.Err())
	}
	if lastStmtError != nil {
		return errwrap.Wrapf("could not perform all revocation statements: {{err}}", lastStmtError)
	}

	// Drop this user
	stmt, err = db.PrepareContext(ctx, fmt.Sprintf(
		`DROP ROLE IF EXISTS %s;`, username))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}
*/
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
		Username: c.Username,
		Password: password,
		Roles: userOpts.Roles,
		Groups: userOpts.Groups,
		DisplayName: userOpts.DisplayName,
	}
		
	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout: 1 * time.Second,
			DomainName: string(userOpts.Domain),
		})

	if err != nil {
		return nil, err
	}

	// Close the database connection to ensure no new connections come in
	if err := c.Close(); err != nil {
		return nil, err
	}

	c.rawConfig["password"] = password
	
	return c.rawConfig, nil
}

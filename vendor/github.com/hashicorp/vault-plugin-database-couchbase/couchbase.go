package couchbase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/template"
)

const (
	couchbaseTypeName        = "couchbase"
	defaultCouchbaseUserRole = `{"Roles": [{"role":"ro_admin"}]}`
	defaultTimeout           = 20000 * time.Millisecond

	defaultUserNameTemplate = `V_{{.DisplayName | uppercase | truncate 64}}_{{.RoleName | uppercase | truncate 64}}_{{random 20 | uppercase}}_{{unix_time}}`
)

var (
	_ dbplugin.Database = &CouchbaseDB{}
)

// Type that combines the custom plugins Couchbase database connection configuration options and the Vault CredentialsProducer
// used for generating user information for the Couchbase database.
type CouchbaseDB struct {
	*couchbaseDBConnectionProducer
	credsutil.CredentialsProducer

	usernameProducer template.StringTemplate
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

	db := &CouchbaseDB{
		couchbaseDBConnectionProducer: connProducer,
	}

	return db
}

func (c *CouchbaseDB) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
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
	c.usernameProducer = up

	err = c.couchbaseDBConnectionProducer.Initialize(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}
	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (c *CouchbaseDB) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	username, err := c.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to generate username: %w", err)
	}
	username = strings.ToUpper(username)

	db, err := c.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to get connection: %w", err)
	}

	err = newUser(ctx, db, username, req)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}

	return resp, nil
}

func (c *CouchbaseDB) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password != nil {
		err := c.changeUserPassword(ctx, req.Username, req.Password.NewPassword)
		return dbplugin.UpdateUserResponse{}, err
	}
	return dbplugin.UpdateUserResponse{}, nil
}

func (c *CouchbaseDB) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	c.Lock()
	defer c.Unlock()

	db, err := c.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("failed to make connection: %w", err)
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

	err = mgr.DropUser(req.Username, nil)

	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	return dbplugin.DeleteUserResponse{}, nil
}

func newUser(ctx context.Context, db *gocb.Cluster, username string, req dbplugin.NewUserRequest) error {
	statements := removeEmpty(req.Statements.Commands)
	if len(statements) == 0 {
		statements = append(statements, defaultCouchbaseUserRole)
	}

	jsonRoleAndGroupData := []byte(statements[0])

	var rag RolesAndGroups

	err := json.Unmarshal(jsonRoleAndGroupData, &rag)
	if err != nil {
		return errwrap.Wrapf("error unmarshalling roles and groups creation statement JSON: {{err}}", err)
	}

	// Get the UserManager

	mgr := db.Users()

	user := gocb.User{
		Username:    username,
		DisplayName: req.UsernameConfig.DisplayName,
		Password:    req.Password,
		Roles:       rag.Roles,
		Groups:      rag.Groups,
	}

	err = mgr.UpsertUser(user,
		&gocb.UpsertUserOptions{
			Timeout:    computeTimeout(ctx),
			DomainName: "local",
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *CouchbaseDB) changeUserPassword(ctx context.Context, username, password string) error {
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
	user, err := mgr.GetUser(username, nil)

	if err != nil {
		return fmt.Errorf("unable to retrieve user %s: %w", username, err)
	}
	user.User.Password = password

	err = mgr.UpsertUser(user.User,
		&gocb.UpsertUserOptions{
			Timeout:    computeTimeout(ctx),
			DomainName: "local",
		})

	if err != nil {
		return err
	}

	return nil
}

func removeEmpty(strs []string) []string {
	var newStrs []string
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		newStrs = append(newStrs, str)
	}

	return newStrs
}

func computeTimeout(ctx context.Context) (timeout time.Duration) {
	deadline, ok := ctx.Deadline()
	if ok {
		return time.Until(deadline)
	}
	return defaultTimeout
}

func (c *CouchbaseDB) getConnection(ctx context.Context) (*gocb.Cluster, error) {
	db, err := c.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return db.(*gocb.Cluster), nil
}

func (c *CouchbaseDB) Type() (string, error) {
	return couchbaseTypeName, nil
}

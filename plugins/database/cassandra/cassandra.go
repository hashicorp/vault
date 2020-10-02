package cassandra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

const (
	defaultUserCreationCQL   = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultUserDeletionCQL   = `DROP USER '{{username}}';`
	defaultChangePasswordCQL = `ALTER USER {{username}} WITH PASSWORD '{{password}}';`
	cassandraTypeName        = "cassandra"
)

var _ newdbplugin.Database = &Cassandra{}

// Cassandra is an implementation of Database interface
type Cassandra struct {
	*cassandraConnectionProducer
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	db := new()
	dbType := newdbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *Cassandra {
	connProducer := &cassandraConnectionProducer{}
	connProducer.Type = cassandraTypeName

	return &Cassandra{
		cassandraConnectionProducer: connProducer,
	}
}

// Run instantiates a Cassandra object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	newdbplugin.Serve(dbType.(newdbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

// Type returns the TypeName for this backend
func (c *Cassandra) Type() (string, error) {
	return cassandraTypeName, nil
}

func (c *Cassandra) getConnection(ctx context.Context) (*gocql.Session, error) {
	session, err := c.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return session.(*gocql.Session), nil
}

// NewUser generates the username/password on the underlying Cassandra secret backend as instructed by
// the statements provided.
func (c *Cassandra) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (newdbplugin.NewUserResponse, error) {
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection(ctx)
	if err != nil {
		return newdbplugin.NewUserResponse{}, err
	}

	creationCQL := req.Statements.Commands
	if len(creationCQL) == 0 {
		creationCQL = []string{defaultUserCreationCQL}
	}

	rollbackCQL := req.RollbackStatements.Commands
	if len(rollbackCQL) == 0 {
		rollbackCQL = []string{defaultUserDeletionCQL}
	}

	username, err := credsutil.GenerateUsername(
		credsutil.DisplayName(req.UsernameConfig.DisplayName, 15),
		credsutil.RoleName(req.UsernameConfig.RoleName, 15),
		credsutil.Separator("_"),
		credsutil.MaxLength(100),
		credsutil.ToLower(),
	)
	if err != nil {
		return newdbplugin.NewUserResponse{}, err
	}
	username = strings.ReplaceAll(username, "-", "_")

	for _, stmt := range creationCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
				"password": req.Password,
			}
			err = session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			if err != nil {
				rollbackErr := rollbackUser(ctx, session, username, rollbackCQL)
				if rollbackErr != nil {
					err = multierror.Append(err, rollbackErr)
				}
				return newdbplugin.NewUserResponse{}, err
			}
		}
	}

	resp := newdbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func (c *Cassandra) newUsername(config newdbplugin.UsernameMetadata) (string, error) {
	displayName := trunc(config.DisplayName, 15)
	roleName := trunc(config.RoleName, 15)

	userUUID, err := credsutil.RandomAlphaNumeric(20, false)
	if err != nil {
		return "", err
	}

	now := fmt.Sprint(time.Now().Unix())

	parts := []string{
		"v",
		displayName,
		roleName,
		userUUID,
		now,
	}
	username := joinNonEmpty("_", parts...)
	username = trunc(username, 100)
	username = strings.ReplaceAll(username, "-", "_")
	username = strings.ToLower(username)

	return username, nil
}

func trunc(str string, l int) string {
	if len(str) < l {
		return str
	}
	return str[:l]
}

func joinNonEmpty(sep string, vals ...string) string {
	if sep == "" {
		return strings.Join(vals, sep)
	}
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	}
	builder := &strings.Builder{}
	for _, val := range vals {
		if val == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(val)
	}
	return builder.String()
}

func rollbackUser(ctx context.Context, session *gocql.Session, username string, rollbackCQL []string) error {
	for _, stmt := range rollbackCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			if err != nil {
				return fmt.Errorf("failed to roll back user %s: %w", username, err)
			}
		}
	}
	return nil
}

func (c *Cassandra) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest) (newdbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return newdbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	if req.Password != nil {
		err := c.changeUserPassword(ctx, req.Username, req.Password)
		return newdbplugin.UpdateUserResponse{}, err
	}
	// Expiration is no-op
	return newdbplugin.UpdateUserResponse{}, nil
}

func (c *Cassandra) changeUserPassword(ctx context.Context, username string, changePass *newdbplugin.ChangePassword) error {
	session, err := c.getConnection(ctx)
	if err != nil {
		return err
	}

	rotateCQL := changePass.Statements.Commands
	if len(rotateCQL) == 0 {
		rotateCQL = []string{defaultChangePasswordCQL}
	}

	var result *multierror.Error
	for _, stmt := range rotateCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
				"password": changePass.NewPassword,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// DeleteUser attempts to drop the specified user.
func (c *Cassandra) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection(ctx)
	if err != nil {
		return newdbplugin.DeleteUserResponse{}, err
	}

	revocationCQL := req.Statements.Commands
	if len(revocationCQL) == 0 {
		revocationCQL = []string{defaultUserDeletionCQL}
	}

	var result *multierror.Error
	for _, stmt := range revocationCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": req.Username,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()

			result = multierror.Append(result, err)
		}
	}

	return newdbplugin.DeleteUserResponse{}, result.ErrorOrNil()
}

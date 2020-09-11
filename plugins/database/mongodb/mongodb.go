package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const mongoDBTypeName = "mongodb"

// MongoDB is an implementation of Database interface
type MongoDB struct {
	*mongoDBConnectionProducer
}

var _ newdbplugin.Database = &MongoDB{}

// New returns a new MongoDB instance
func New() (interface{}, error) {
	db := new()
	dbType := newdbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *MongoDB {
	connProducer := &mongoDBConnectionProducer{
		Type: mongoDBTypeName,
	}

	return &MongoDB{
		mongoDBConnectionProducer: connProducer,
	}
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

// Type returns the TypeName for this backend
func (m *MongoDB) Type() (string, error) {
	return mongoDBTypeName, nil
}

func (m *MongoDB) getConnection(ctx context.Context) (*mongo.Client, error) {
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return client.(*mongo.Client), nil
}

func (m *MongoDB) Initialize(ctx context.Context, req newdbplugin.InitializeRequest) (newdbplugin.InitializeResponse, error) {
	m.Lock()
	defer m.Unlock()

	m.RawConfig = req.Config

	err := mapstructure.WeakDecode(req.Config, m.mongoDBConnectionProducer)
	if err != nil {
		return newdbplugin.InitializeResponse{}, err
	}

	if len(m.ConnectionURL) == 0 {
		return newdbplugin.InitializeResponse{}, fmt.Errorf("connection_url cannot be empty-mongo fail")
	}

	writeOpts, err := m.getWriteConcern()
	if err != nil {
		return newdbplugin.InitializeResponse{}, err
	}

	authOpts, err := m.getTLSAuth()
	if err != nil {
		return newdbplugin.InitializeResponse{}, err
	}

	m.clientOptions = options.MergeClientOptions(writeOpts, authOpts)

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	m.Initialized = true

	if req.VerifyConnection {
		_, err := m.Connection(ctx)
		if err != nil {
			return newdbplugin.InitializeResponse{}, fmt.Errorf("failed to verify connection: %w", err)
		}

		err = m.client.Ping(ctx, readpref.Primary())
		if err != nil {
			return newdbplugin.InitializeResponse{}, fmt.Errorf("failed to verify connection: %w", err)
		}
	}

	resp := newdbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (m *MongoDB) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (newdbplugin.NewUserResponse, error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	if len(req.Statements.Commands) == 0 {
		return newdbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	username, err := newUsername(req.UsernameConfig)
	if err != nil {
		return newdbplugin.NewUserResponse{}, err
	}

	// Unmarshal statements.CreationStatements into mongodbRoles
	var mongoCS mongoDBStatement
	err = json.Unmarshal([]byte(req.Statements.Commands[0]), &mongoCS)
	if err != nil {
		return newdbplugin.NewUserResponse{}, err
	}

	// Default to "admin" if no db provided
	if mongoCS.DB == "" {
		mongoCS.DB = "admin"
	}

	if len(mongoCS.Roles) == 0 {
		return newdbplugin.NewUserResponse{}, fmt.Errorf("roles array is required in creation statement")
	}

	createUserCmd := createUserCommand{
		Username: username,
		Password: req.Password,
		Roles:    mongoCS.Roles.toStandardRolesArray(),
	}

	if err := m.runCommandWithRetry(ctx, mongoCS.DB, createUserCmd); err != nil {
		return newdbplugin.NewUserResponse{}, err
	}

	resp := newdbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func newUsername(config newdbplugin.UsernameMetadata) (string, error) {
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
	username := joinNonEmpty("-", parts...)
	username = trunc(username, 100)

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

func removeEmpty(input []string) []string {
	output := []string{}
	for _, val := range input {
		if val != "" {
			output = append(output, val)
		}
	}
	return output
}

func (m *MongoDB) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest) (newdbplugin.UpdateUserResponse, error) {
	if req.Password != nil {
		err := m.changeUserPassword(ctx, req.Username, req.Password.NewPassword)
		return newdbplugin.UpdateUserResponse{}, err
	}
	return newdbplugin.UpdateUserResponse{}, nil
}

func (m *MongoDB) changeUserPassword(ctx context.Context, username, password string) error {
	m.Lock()
	defer m.Unlock()

	connURL := m.getConnectionURL()
	cs, err := connstring.Parse(connURL)
	if err != nil {
		return err
	}

	// Currently doesn't support custom statements for changing the user's password
	changeUserCmd := &updateUserCommand{
		Username: username,
		Password: password,
	}

	database := cs.Database
	if username == m.Username || database == "" {
		database = "admin"
	}

	err = m.runCommandWithRetry(ctx, database, changeUserCmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	m.Lock()
	defer m.Unlock()

	// If no revocation statements provided, pass in empty JSON
	var revocationStatement string
	switch len(req.Statements.Commands) {
	case 0:
		revocationStatement = `{}`
	case 1:
		revocationStatement = req.Statements.Commands[0]
	default:
		return newdbplugin.DeleteUserResponse{}, fmt.Errorf("expected 0 or 1 revocation statements, got %d", len(req.Statements.Commands))
	}

	// Unmarshal revocation statements into mongodbRoles
	var mongoCS mongoDBStatement
	err := json.Unmarshal([]byte(revocationStatement), &mongoCS)
	if err != nil {
		return newdbplugin.DeleteUserResponse{}, err
	}

	db := mongoCS.DB
	// If db is not specified, use the default authenticationDatabase "admin"
	if db == "" {
		db = "admin"
	}

	dropUserCmd := &dropUserCommand{
		Username:     req.Username,
		WriteConcern: writeconcern.New(writeconcern.WMajority()),
	}

	err = m.runCommandWithRetry(ctx, db, dropUserCmd)
	return newdbplugin.DeleteUserResponse{}, err
}

// runCommandWithRetry runs a command and retries once more if there's a failure
// on the first attempt. This should be called with the lock held
func (m *MongoDB) runCommandWithRetry(ctx context.Context, db string, cmd interface{}) error {
	// Get the client
	client, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	// Run command
	result := client.Database(db).RunCommand(ctx, cmd, nil)

	// Error check on the first attempt
	err = result.Err()
	switch {
	case err == nil:
		return nil
	case err == io.EOF, strings.Contains(err.Error(), "EOF"):
		// Call getConnection to reset and retry query if we get an EOF error on first attempt.
		client, err = m.getConnection(ctx)
		if err != nil {
			return err
		}
		result = client.Database(db).RunCommand(ctx, cmd, nil)
		if err := result.Err(); err != nil {
			return err
		}
	default:
		return err
	}

	return nil
}

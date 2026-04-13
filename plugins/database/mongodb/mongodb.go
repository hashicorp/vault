// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	mongoDBTypeName = "mongodb"

	defaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 15) (.RoleName | truncate 15) (random 20) (unix_time) | replace "." "-" | truncate 100 }}`
)

// MongoDB is an implementation of Database interface
type MongoDB struct {
	*mongoDBConnectionProducer

	usernameProducer template.StringTemplate
}

var _ dbplugin.Database = &MongoDB{}

// New returns a new MongoDB instance
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
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

// Type returns the TypeName for this backend
func (m *MongoDB) Type() (string, error) {
	return mongoDBTypeName, nil
}

func (m *MongoDB) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	m.Lock()
	defer m.Unlock()

	m.RawConfig = req.Config

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
	m.usernameProducer = up

	_, err = m.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	err = m.mongoDBConnectionProducer.loadConfig(req.Config)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	m.Initialized = true

	if req.VerifyConnection {
		client, err := m.mongoDBConnectionProducer.createClient(ctx)
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("failed to verify connection: %w", err)
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			_ = client.Disconnect(ctx) // Try to prevent any sort of resource leak
			return dbplugin.InitializeResponse{}, fmt.Errorf("failed to verify connection: %w", err)
		}
		m.mongoDBConnectionProducer.client = client
	}

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (m *MongoDB) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	username, err := m.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// Unmarshal statements.CreationStatements into mongodbRoles
	var mongoCS mongoDBStatement
	err = json.Unmarshal([]byte(req.Statements.Commands[0]), &mongoCS)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// Default to "admin" if no db provided
	if mongoCS.DB == "" {
		mongoCS.DB = "admin"
	}

	if len(mongoCS.Roles) == 0 {
		return dbplugin.NewUserResponse{}, fmt.Errorf("roles array is required in creation statement")
	}

	createUserCmd := createUserCommand{
		Username: username,
		Password: req.Password,
		Roles:    mongoCS.Roles.toStandardRolesArray(),
	}

	if err := m.runCommandWithRetry(ctx, mongoCS.DB, createUserCmd); err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func (m *MongoDB) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password != nil {
		err := m.changeUserPassword(ctx, req.Username, req.Password.NewPassword, req.Password.Statements.Commands)
		return dbplugin.UpdateUserResponse{}, err
	}
	return dbplugin.UpdateUserResponse{}, nil
}

func (m *MongoDB) changeUserPassword(ctx context.Context, username, password string, rotateStatements []string) error {
	connURL := m.getConnectionURL()
	cs, err := connstring.Parse(connURL)
	if err != nil {
		return err
	}

	database := cs.Database
	if database == "" {
		database = "admin"
	}

	if len(rotateStatements) == 0 {
		changeUserCmd := &updateUserCommand{
			Username: username,
			Password: password,
		}

		if err := m.runCommandWithRetry(ctx, database, changeUserCmd); err != nil {
			return err
		}
		return nil
	}

	queryMap := map[string]string{
		"name":     username,
		"username": username,
		"password": password,
	}

	for i, stmt := range rotateStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		stmt = dbutil.QueryHelper(stmt, queryMap)

		var rotationStmt mongoDBRotationStatement
		if err := json.Unmarshal([]byte(stmt), &rotationStmt); err != nil {
			return fmt.Errorf("failed to parse rotation statement %d as JSON: %w", i, err)
		}

		if err := m.executeRotationStatement(ctx, i, database, rotationStmt); err != nil {
			return err
		}
	}

	return nil
}

func (m *MongoDB) executeRotationStatement(ctx context.Context, idx int, defaultDB string, rotationStmt mongoDBRotationStatement) error {
	db := rotationStmt.DB
	if db == "" {
		db = defaultDB
	}
	if db == "" {
		db = "admin"
	}

	if len(rotationStmt.Command) == 0 {
		return fmt.Errorf("rotation statement %d must contain a non-empty command object", idx)
	}

	var cmd bson.D
	if err := bson.UnmarshalExtJSON(rotationStmt.Command, true, &cmd); err != nil {
		return fmt.Errorf("rotation statement %d has invalid command JSON: %w", idx, err)
	}
	if len(cmd) == 0 {
		return fmt.Errorf("rotation statement %d must contain a non-empty command object", idx)
	}

	if err := m.runCommandWithRetry(ctx, db, cmd); err != nil {
		if shouldIgnoreRotationError(err, rotationStmt.IgnoreErrors) {
			return nil
		}
		return fmt.Errorf("rotation statement %d failed: %w", idx, err)
	}
	return nil
}

func shouldIgnoreRotationError(err error, ignoreErrors []string) bool {
	if err == nil || len(ignoreErrors) == 0 {
		return false
	}

	// Allow ergonomic ignore tokens that map to well-known codes.
	// This avoids forcing users to configure numeric codes like "51003".
	wellKnownErrorCodes := map[string][]int32{
		// createUser: "User \"...\" already exists"
		"UserAlreadyExists": {51003},
		// Common duplicate key codes
		"DuplicateKey": {11000, 11001},
	}

	cErr, ok := err.(mongo.CommandError)
	if ok {
		errText := strings.ToLower(err.Error())
		msgText := strings.ToLower(cErr.Message)

		for _, token := range ignoreErrors {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			// Match by command error name, when available.
			if cErr.Name != "" && cErr.Name == token {
				return true
			}
			// Also allow matching by numeric error code, e.g. "51003".
			if code, convErr := strconv.Atoi(token); convErr == nil && int32(code) == cErr.Code {
				return true
			}
			// Allow matching by well-known symbolic token -> code mapping.
			if codes, ok := wellKnownErrorCodes[token]; ok {
				for _, c := range codes {
					if cErr.Code == c {
						return true
					}
				}
			}
			// As a last resort, allow ignoring by message semantics for common cases.
			// Some server errors expose a code but no stable "Name".
			switch token {
			case "UserAlreadyExists":
				if strings.Contains(msgText, "already exists") || strings.Contains(errText, "already exists") {
					return true
				}
			}
		}
		return false
	}

	// Some drivers may return a WriteException; allow matching by write error codes too.
	if wErr, ok := err.(mongo.WriteException); ok {
		for _, token := range ignoreErrors {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			code, convErr := strconv.Atoi(token)
			if convErr != nil {
				// Also allow well-known symbolic token -> code mapping.
				if codes, ok := wellKnownErrorCodes[token]; ok {
					for _, c := range codes {
						for _, we := range wErr.WriteErrors {
							if int32(we.Code) == c {
								return true
							}
						}
					}
				}
				continue
			}
			for _, we := range wErr.WriteErrors {
				if int32(we.Code) == int32(code) {
					return true
				}
			}
		}
	}

	return false
}

func (m *MongoDB) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	// If no revocation statements provided, pass in empty JSON
	var revocationStatement string
	switch len(req.Statements.Commands) {
	case 0:
		revocationStatement = `{}`
	case 1:
		revocationStatement = req.Statements.Commands[0]
	default:
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("expected 0 or 1 revocation statements, got %d", len(req.Statements.Commands))
	}

	// Unmarshal revocation statements into mongodbRoles
	var mongoCS mongoDBStatement
	err := json.Unmarshal([]byte(revocationStatement), &mongoCS)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	db := mongoCS.DB
	// If db is not specified, use the default authenticationDatabase "admin"
	if db == "" {
		db = "admin"
	}

	// Set the write concern. The default is majority.
	writeConcern := writeconcern.New(writeconcern.WMajority())
	opts, err := m.getWriteConcern()
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	if opts != nil {
		writeConcern = opts.WriteConcern
	}

	dropUserCmd := &dropUserCommand{
		Username:     req.Username,
		WriteConcern: writeConcern,
	}

	err = m.runCommandWithRetry(ctx, db, dropUserCmd)
	cErr, ok := err.(mongo.CommandError)
	if ok && cErr.Name == "UserNotFound" { // User already removed, don't retry needlessly
		log.Default().Warn("MongoDB user was deleted prior to lease revocation", "user", req.Username)
		return dbplugin.DeleteUserResponse{}, nil
	}

	return dbplugin.DeleteUserResponse{}, err
}

// runCommandWithRetry runs a command and retries once more if there's a failure
// on the first attempt. This should be called with the lock held
func (m *MongoDB) runCommandWithRetry(ctx context.Context, db string, cmd interface{}) error {
	// Get the client
	client, err := m.Connection(ctx)
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
		client, err = m.Connection(ctx)
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

package database

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type databaseVersionWrapper struct {
	database       newdbplugin.Database
	legacyDatabase dbplugin.Database
}

func makeDatabase(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (dbw databaseVersionWrapper, err error) {
	newDB, err := newdbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw.database = newDB
		return dbw, nil
	}

	legacyDB, err := dbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw.legacyDatabase = legacyDB
		return dbw, nil
	}

	return dbw, fmt.Errorf("invalid database version")
}

func (db databaseVersionWrapper) Type() (string, error) {
	if db.database != nil {
		return db.database.Type()
	}
	return db.legacyDatabase.Type()
}

func (db databaseVersionWrapper) Close() error {
	if db.database != nil {
		return db.database.Close()
	}
	return db.legacyDatabase.Close()
}

// /////////////////////////////////////////////////////////////////////////////////
// Initialization
// /////////////////////////////////////////////////////////////////////////////////

func initDatabase(ctx context.Context, dbw databaseVersionWrapper, connDetails map[string]interface{}, verifyConnection bool) (newConfig map[string]interface{}, err error) {
	if dbw.database != nil {
		return initNewDatabase(ctx, dbw, connDetails, verifyConnection)
	}
	return initLegacyDatabase(ctx, dbw, connDetails, verifyConnection)
}

func initNewDatabase(ctx context.Context, dbw databaseVersionWrapper, connDetails map[string]interface{}, verifyConnection bool) (newConfig map[string]interface{}, err error) {
	req := newdbplugin.InitializeRequest{
		Config:           connDetails,
		VerifyConnection: verifyConnection,
	}
	resp, err := dbw.database.Initialize(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

func initLegacyDatabase(ctx context.Context, dbw databaseVersionWrapper, connDetails map[string]interface{}, verifyConnection bool) (newConfig map[string]interface{}, err error) {
	return dbw.legacyDatabase.Init(ctx, connDetails, verifyConnection)
}

// /////////////////////////////////////////////////////////////////////////////////
// Password generation
// /////////////////////////////////////////////////////////////////////////////////

type passwordGenerator interface {
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

var (
	defaultPasswordGenerator = random.DefaultStringGenerator
)

func generatePassword(ctx context.Context, dbw databaseVersionWrapper, generator passwordGenerator, passwordPolicy string) (password string, err error) {
	// If using the legacy database, use GenerateCredentials instead of password policies
	// This will keep the existing behavior even though passwords can be generated with a policy
	if dbw.legacyDatabase != nil {
		password, err := dbw.legacyDatabase.GenerateCredentials(ctx)
		if err != nil {
			return "", err
		}
		return password, nil
	}

	if passwordPolicy == "" {
		return defaultPasswordGenerator.Generate(ctx, rand.Reader)
	}
	return generator.GeneratePasswordFromPolicy(ctx, passwordPolicy)
}

// /////////////////////////////////////////////////////////////////////////////////
// Change new user
// /////////////////////////////////////////////////////////////////////////////////

func createUser(ctx context.Context, dbw databaseVersionWrapper, pg passwordGenerator, statements dbplugin.Statements, displayName, roleName string, expiration time.Time, passwordPolicy string) (username, password string, err error) {
	if dbw.database != nil {
		return createNewUser(ctx, dbw, pg, displayName, roleName, expiration, passwordPolicy, statements)
	}
	return createLegacyUser(ctx, dbw, statements, displayName, roleName, expiration)
}

// createNewUser creates a user with the v5 Database interface
func createNewUser(ctx context.Context,
	dbw databaseVersionWrapper,
	pg passwordGenerator,
	displayName, roleName string,
	expiration time.Time,
	passwordPolicy string,
	statements dbplugin.Statements) (username, password string, err error) {

	pass, err := generatePassword(ctx, dbw, pg, passwordPolicy)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate password: %w", err)
	}

	req := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: displayName,
			RoleName:    roleName,
		},
		Statements: newdbplugin.Statements{
			Commands: statements.Creation,
		},
		RollbackStatements: newdbplugin.Statements{
			Commands: statements.Rollback,
		},
		Password:   pass,
		Expiration: expiration,
	}

	resp, err := dbw.database.NewUser(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %w", err)
	}
	return resp.Username, pass, nil
}

// createLegacyUser creates a user with the v4 Database interface
func createLegacyUser(
	ctx context.Context,
	dbw databaseVersionWrapper,
	statements dbplugin.Statements,
	displayName, roleName string,
	expiration time.Time) (username, password string, err error) {

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: displayName,
		RoleName:    roleName,
	}

	return dbw.legacyDatabase.CreateUser(ctx, statements, usernameConfig, expiration)
}

// /////////////////////////////////////////////////////////////////////////////////
// Change user password
// /////////////////////////////////////////////////////////////////////////////////

func changeUserPassword(ctx context.Context, dbw databaseVersionWrapper, username, newpassword string, rotationCmds []string) error {
	if dbw.database != nil {
		return changeUserPasswordNew(ctx, dbw, username, newpassword, rotationCmds)
	}
	return changeUserPasswordLegacy(ctx, dbw, username, newpassword, rotationCmds)
}

func changeUserPasswordNew(ctx context.Context, dbw databaseVersionWrapper, username, newPassword string, rotationCmds []string) error {
	req := newdbplugin.UpdateUserRequest{
		Username: username,
		Password: &newdbplugin.ChangePassword{
			NewPassword: newPassword,
			Statements: newdbplugin.Statements{
				Commands: rotationCmds,
			},
		},
	}
	_, err := dbw.database.UpdateUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func changeUserPasswordLegacy(
	ctx context.Context,
	dbw databaseVersionWrapper,
	username, newPassword string,
	rotationCmds []string) (err error) {

	// Attempt to use SetCredentials for the root credential rotation
	statements := dbplugin.Statements{Rotation: rotationCmds}
	userConfig := dbplugin.StaticUserConfig{
		Username: username,
		Password: newPassword,
	}

	_, _, err = dbw.legacyDatabase.SetCredentials(ctx, statements, userConfig)
	return err
}

// /////////////////////////////////////////////////////////////////////////////////
// Renewing expiration
// /////////////////////////////////////////////////////////////////////////////////

func renewUser(ctx context.Context, dbw databaseVersionWrapper, username string, newExpiration time.Time, commands []string) error {
	if dbw.database != nil {
		return renewUserNew(ctx, dbw, username, newExpiration, commands)
	}
	return renewUserLegacy(ctx, dbw, username, newExpiration, commands)
}

func renewUserNew(ctx context.Context, dbw databaseVersionWrapper, username string, newExpiration time.Time, commands []string) error {
	req := newdbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &newdbplugin.ChangeExpiration{
			NewExpiration: newExpiration,
			Statements: newdbplugin.Statements{
				Commands: commands,
			},
		},
	}
	_, err := dbw.database.UpdateUser(ctx, req)
	return err
}

func renewUserLegacy(ctx context.Context, dbw databaseVersionWrapper, username string, newExpiration time.Time, commands []string) error {
	statements := dbplugin.Statements{
		Renewal: commands,
	}
	return dbw.legacyDatabase.RenewUser(ctx, statements, username, newExpiration)
}

// /////////////////////////////////////////////////////////////////////////////////
// Deleting user
// /////////////////////////////////////////////////////////////////////////////////

func deleteUser(ctx context.Context, dbw databaseVersionWrapper, username string, commands []string) error {
	if dbw.database != nil {
		return deleteUserNew(ctx, dbw, username, commands)
	}
	return deleteUserLegacy(ctx, dbw, username, commands)
}

func deleteUserNew(ctx context.Context, dbw databaseVersionWrapper, username string, commands []string) error {
	req := newdbplugin.DeleteUserRequest{
		Username: username,
		Statements: newdbplugin.Statements{
			Commands: commands,
		},
	}
	_, err := dbw.database.DeleteUser(ctx, req)
	return err
}

func deleteUserLegacy(ctx context.Context, dbw databaseVersionWrapper, username string, commands []string) error {
	statements := dbplugin.Statements{
		Revocation: commands,
	}
	return dbw.legacyDatabase.RevokeUser(ctx, statements, username)
}

// /////////////////////////////////////////////////////////////////////////////////
// Storage helpers
// /////////////////////////////////////////////////////////////////////////////////

func storeConfig(ctx context.Context, storage logical.Storage, name string, config *DatabaseConfig) error {
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
	if err != nil {
		return fmt.Errorf("unable to marshal object to JSON: %w", err)
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to save object: %w", err)
	}
	return nil
}

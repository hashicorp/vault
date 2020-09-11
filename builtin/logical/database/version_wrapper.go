package database

import (
	"context"
	"crypto/rand"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
)

type databaseVersionWrapper struct {
	v4 dbplugin.Database
	v5 newdbplugin.Database
}

// newDatabaseWrapper figures out which version of the database the pluginName is referring to and returns a wrapper object
// that can be used to make operations on the underlying database plugin.
func newDatabaseWrapper(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (dbw databaseVersionWrapper, err error) {
	newDB, err := newdbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw = databaseVersionWrapper{
			v5: newDB,
		}
		return dbw, nil
	}

	legacyDB, err := dbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw = databaseVersionWrapper{
			v4: legacyDB,
		}
		return dbw, nil
	}

	return dbw, fmt.Errorf("invalid database version")
}

// Initialize the underlying database. This is analogous to a constructor on the database plugin object.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Initialize(ctx context.Context, req newdbplugin.InitializeRequest) (newdbplugin.InitializeResponse, error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return newdbplugin.InitializeResponse{}, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isNewDB() {
		return d.v5.Initialize(ctx, req)
	}

	// v4 Database
	saveConfig, err := d.v4.Init(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return newdbplugin.InitializeResponse{}, err
	}
	resp := newdbplugin.InitializeResponse{
		Config: saveConfig,
	}
	return resp, nil
}

// NewUser in the database. This is different from the v5 Database in that it returns a password as well.
// This is done because the v4 Database is expected to generate a password and return it. The NewUserReponse
// does not have a way of returning the password so this function signature needs to be different.
// The password returned here should be considered the source of truth, not the provided password.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (resp newdbplugin.NewUserResponse, password string, err error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return newdbplugin.NewUserResponse{}, "", fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isNewDB() {
		resp, err = d.v5.NewUser(ctx, req)
		return resp, req.Password, err
	}

	// v4 Database
	stmts := dbplugin.Statements{
		Creation: req.Statements.Commands,
		Rollback: req.RollbackStatements.Commands,
	}
	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: req.UsernameConfig.DisplayName,
		RoleName:    req.UsernameConfig.RoleName,
	}
	username, password, err := d.v4.CreateUser(ctx, stmts, usernameConfig, req.Expiration)
	if err != nil {
		return resp, "", err
	}

	resp = newdbplugin.NewUserResponse{
		Username: username,
	}
	return resp, password, nil
}

// UpdateUser in the underlying database. This is used to update any information currently supported
// in the UpdateUserRequest such as password credentials or user TTL.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest, isRootUser bool) (saveConfig map[string]interface{}, err error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return nil, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isNewDB() {
		_, err := d.v5.UpdateUser(ctx, req)
		return nil, err
	}

	// /////////////////////////////////////////////
	// v4 Database

	if req.Password == nil && req.Expiration == nil {
		return nil, fmt.Errorf("missing change to be sent to the database")
	}
	if req.Password != nil && req.Expiration != nil {
		// We could support this, but it would require handling partial
		// errors which I'm punting on since we don't need it for now
		return nil, fmt.Errorf("cannot specify both password and expiration change at the same time")
	}

	// Change password
	if req.Password != nil {
		return d.changePasswordLegacy(ctx, req.Username, req.Password, isRootUser)
	}

	// Change expiration date
	if req.Expiration != nil {
		stmts := dbplugin.Statements{
			Renewal: req.Expiration.Statements.Commands,
		}
		err := d.v4.RenewUser(ctx, stmts, req.Username, req.Expiration.NewExpiration)
		return nil, err
	}
	return nil, nil
}

// changePasswordLegacy attempts to use SetCredentials to change the password for the user with the password provided
// in ChangePassword. If that user is the root user and SetCredentials is unimplemented, it will fall back to using
// RotateRootCredentials. If not the root user, this will not use RotateRootCredentials.
func (d databaseVersionWrapper) changePasswordLegacy(ctx context.Context, username string, passwordChange *newdbplugin.ChangePassword, isRootUser bool) (saveConfig map[string]interface{}, err error) {
	err = d.changeUserPasswordLegacy(ctx, username, passwordChange)

	// If changing the root user's password but SetCredentials is unimplemented, fall back to RotateRootCredentials
	if isRootUser && status.Code(err) == codes.Unimplemented {
		saveConfig, err = d.changeRootUserPasswordLegacy(ctx, passwordChange)
		if err != nil {
			return nil, err
		}
		return saveConfig, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (d databaseVersionWrapper) changeUserPasswordLegacy(ctx context.Context, username string, passwordChange *newdbplugin.ChangePassword) (err error) {
	stmts := dbplugin.Statements{
		Rotation: passwordChange.Statements.Commands,
	}
	staticConfig := dbplugin.StaticUserConfig{
		Username: username,
		Password: passwordChange.NewPassword,
	}
	_, _, err = d.v4.SetCredentials(ctx, stmts, staticConfig)
	return err
}

func (d databaseVersionWrapper) changeRootUserPasswordLegacy(ctx context.Context, passwordChange *newdbplugin.ChangePassword) (saveConfig map[string]interface{}, err error) {
	return d.v4.RotateRootCredentials(ctx, passwordChange.Statements.Commands)
}

// DeleteUser in the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return newdbplugin.DeleteUserResponse{}, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isNewDB() {
		return d.v5.DeleteUser(ctx, req)
	}

	// v4 Database
	stmts := dbplugin.Statements{
		Revocation: req.Statements.Commands,
	}
	err := d.v4.RevokeUser(ctx, stmts, req.Username)
	return newdbplugin.DeleteUserResponse{}, err
}

// Type of the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Type() (string, error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return "", fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isNewDB() {
		return d.v5.Type()
	}

	// v4 Database
	return d.v4.Type()
}

// Close the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Close() error {
	if !d.isNewDB() && !d.isLegacyDB() {
		return fmt.Errorf("no underlying database specified")
	}
	// v5 Database
	if d.isNewDB() {
		return d.v5.Close()
	}

	// v4 Database
	return d.v4.Close()
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

// GeneratePassword either from the v4 database or by using the provided password policy. If using a v5 database
// and no password policy is specified, this will have a reasonable default password generator.
func (d databaseVersionWrapper) GeneratePassword(ctx context.Context, generator passwordGenerator, passwordPolicy string) (password string, err error) {
	if !d.isNewDB() && !d.isLegacyDB() {
		return "", fmt.Errorf("no underlying database specified")
	}

	// If using the legacy database, use GenerateCredentials instead of password policies
	// This will keep the existing behavior even though passwords can be generated with a policy
	if d.isLegacyDB() {
		password, err := d.v4.GenerateCredentials(ctx)
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

func (d databaseVersionWrapper) isNewDB() bool {
	return d.v5 != nil
}

func (d databaseVersionWrapper) isLegacyDB() bool {
	return d.v4 != nil
}

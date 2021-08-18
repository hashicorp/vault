package database

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/random"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type databaseVersionWrapper struct {
	passwordGenerator passwordGenerator

	v4 v4.Database
	v5 v5.Database
}

type system interface {
	pluginutil.LookRunnerUtil
	passwordGenerator
}

type passwordGenerator interface {
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

// newDatabaseWrapper figures out which version of the database the pluginName is referring to and returns a wrapper object
// that can be used to make operations on the underlying database plugin.
func newDatabaseWrapper(ctx context.Context, pluginName string, sys system, logger log.Logger) (dbw databaseVersionWrapper, err error) {
	newDB, err := v5.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw = databaseVersionWrapper{
			passwordGenerator: sys,
			v5:                newDB,
		}
		return dbw, nil
	}

	merr := &multierror.Error{}
	merr = multierror.Append(merr, err)

	legacyDB, err := v4.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw = databaseVersionWrapper{
			// passwordGenerator isn't needed for v4
			v4: legacyDB,
		}
		return dbw, nil
	}
	merr = multierror.Append(merr, err)

	return dbw, fmt.Errorf("invalid database version: %s", merr)
}

// Initialize the underlying database. This is analogous to a constructor on the database plugin object.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Initialize(ctx context.Context, req v5.InitializeRequest) (v5.InitializeResponse, error) {
	if !d.isV5() && !d.isV4() {
		return v5.InitializeResponse{}, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isV5() {
		return d.v5.Initialize(ctx, req)
	}

	// v4 Database
	saveConfig, err := d.v4.Init(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return v5.InitializeResponse{}, err
	}
	resp := v5.InitializeResponse{
		Config: saveConfig,
	}
	return resp, nil
}

// NewUser in the database. This is different from the v5 Database in that it returns a password as well.
// This is done because the v4 Database is expected to generate a password and return it. The NewUserResponse
// does not have a way of returning the password so this function signature needs to be different.
// The password returned here should be considered the source of truth, not the provided password.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) NewUser(ctx context.Context, dbConfig *DatabaseConfig, role *roleEntry, roleName string, dispName string, expiration time.Time) (respData map[string]interface{}, internalData map[string]interface{}, err error) {
	if !d.isV5() && !d.isV4() {
		return nil, nil, fmt.Errorf("no underlying database specified")
	}

	if d.isV5() {
		return d.newUserV5(ctx, dbConfig, role, roleName, dispName, expiration)
	}

	return d.newUserV4(ctx, role, roleName, dispName, expiration)
}

func (d databaseVersionWrapper) newUserV5(ctx context.Context, dbConfig *DatabaseConfig, role *roleEntry, roleName string, dispName string, expiration time.Time) (respData map[string]interface{}, internalData map[string]interface{}, err error) {
	password, err := d.GeneratePassword(ctx, d.passwordGenerator, dbConfig.PasswordPolicy)

	newUserReq := v5.NewUserRequest{
		UsernameConfig: v5.UsernameMetadata{
			DisplayName: dispName,
			RoleName:    roleName,
		},
		Statements: v5.Statements{
			Commands: role.Statements.Creation,
		},
		RollbackStatements: v5.Statements{
			Commands: role.Statements.Rollback,
		},
		Expiration: expiration,
		Password:   password,
	}

	newUserResp, err := d.v5.NewUser(ctx, newUserReq)
	if err != nil {
		return nil, nil, err
	}

	respData = getNewUserResponseData(newUserReq, newUserResp)
	internalData = map[string]interface{}{
		"username":              newUserResp.Username,
		"role":                  roleName,
		"db_name":               role.DBName,
		"revocation_statements": role.Statements.Revocation,
	}

	return respData, internalData, nil
}

func getNewUserResponseData(req v5.NewUserRequest, resp v5.NewUserResponse) map[string]interface{} {
	respData := map[string]interface{}{
		"username": resp.Username,
	}

	if req.Password != "" {
		respData["password"] = req.Password
	}

	return respData
}

func (d databaseVersionWrapper) newUserV4(ctx context.Context, role *roleEntry, roleName string, dispName string, expiration time.Time) (respData map[string]interface{}, internalData map[string]interface{}, err error) {
	usernameConfig := v4.UsernameConfig{
		DisplayName: dispName,
		RoleName:    roleName,
	}
	username, password, err := d.v4.CreateUser(ctx, role.Statements, usernameConfig, expiration)
	if err != nil {
		return nil, nil, err
	}

	respData = map[string]interface{}{
		"username": username,
		"password": password,
	}
	internalData = map[string]interface{}{
		"username":              username,
		"role":                  roleName,
		"db_name":               role.DBName,
		"revocation_statements": role.Statements.Revocation,
	}
	return respData, internalData, nil
}

// UpdateUser in the underlying database. This is used to update any information currently supported
// in the UpdateUserRequest such as password credentials or user TTL.
// Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) UpdateUser(ctx context.Context, req v5.UpdateUserRequest, isRootUser bool) (saveConfig map[string]interface{}, err error) {
	if !d.isV5() && !d.isV4() {
		return nil, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isV5() {
		_, err := d.v5.UpdateUser(ctx, req)
		return nil, err
	}

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
		stmts := v4.Statements{
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
func (d databaseVersionWrapper) changePasswordLegacy(ctx context.Context, username string, passwordChange *v5.ChangePassword, isRootUser bool) (saveConfig map[string]interface{}, err error) {
	err = d.changeUserPasswordLegacy(ctx, username, passwordChange)

	// If changing the root user's password but SetCredentials is unimplemented, fall back to RotateRootCredentials
	if isRootUser && (err == v4.ErrPluginStaticUnsupported || status.Code(err) == codes.Unimplemented) {
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

func (d databaseVersionWrapper) changeUserPasswordLegacy(ctx context.Context, username string, passwordChange *v5.ChangePassword) (err error) {
	stmts := v4.Statements{
		Rotation: passwordChange.Statements.Commands,
	}
	staticConfig := v4.StaticUserConfig{
		Username: username,
		Password: passwordChange.NewPassword,
	}
	_, _, err = d.v4.SetCredentials(ctx, stmts, staticConfig)
	return err
}

func (d databaseVersionWrapper) changeRootUserPasswordLegacy(ctx context.Context, passwordChange *v5.ChangePassword) (saveConfig map[string]interface{}, err error) {
	return d.v4.RotateRootCredentials(ctx, passwordChange.Statements.Commands)
}

// DeleteUser in the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) DeleteUser(ctx context.Context, req v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	if !d.isV5() && !d.isV4() {
		return v5.DeleteUserResponse{}, fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isV5() {
		return d.v5.DeleteUser(ctx, req)
	}

	// v4 Database
	stmts := v4.Statements{
		Revocation: req.Statements.Commands,
	}
	err := d.v4.RevokeUser(ctx, stmts, req.Username)
	return v5.DeleteUserResponse{}, err
}

// Type of the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Type() (string, error) {
	if !d.isV5() && !d.isV4() {
		return "", fmt.Errorf("no underlying database specified")
	}

	// v5 Database
	if d.isV5() {
		return d.v5.Type()
	}

	// v4 Database
	return d.v4.Type()
}

// Close the underlying database. Errors if the wrapper does not contain an underlying database.
func (d databaseVersionWrapper) Close() error {
	if !d.isV5() && !d.isV4() {
		return fmt.Errorf("no underlying database specified")
	}
	// v5 Database
	if d.isV5() {
		return d.v5.Close()
	}

	// v4 Database
	return d.v4.Close()
}

// /////////////////////////////////////////////////////////////////////////////////
// Password generation
// /////////////////////////////////////////////////////////////////////////////////

var defaultPasswordGenerator = random.DefaultStringGenerator

// GeneratePassword either from the v4 database or by using the provided password policy. If using a v5 database
// and no password policy is specified, this will have a reasonable default password generator.
func (d databaseVersionWrapper) GeneratePassword(ctx context.Context, generator passwordGenerator, passwordPolicy string) (password string, err error) {
	if !d.isV5() && !d.isV4() {
		return "", fmt.Errorf("no underlying database specified")
	}

	// If using the legacy database, use GenerateCredentials instead of password policies
	// This will keep the existing behavior even though passwords can be generated with a policy
	if d.isV4() {
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

func (d databaseVersionWrapper) isV5() bool {
	return d.v5 != nil
}

func (d databaseVersionWrapper) isV4() bool {
	return d.v4 != nil
}

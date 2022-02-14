package influxdbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/influxdata/influxdb-client-go/v2"
)

const (
	influxdbTypeName = "influxdbv2"

	defaultUserNameTemplate = `{{ printf "v_%s_%s_%s_%s" (.DisplayName | truncate 15) (.RoleName | truncate 15) (random 20) (unix_time) | truncate 100 | replace "-" "_" | lowercase }}`
)

var _ dbplugin.Database = &InfluxdbV2{}

// InfluxdbV2 is an implementation of Database interface
type InfluxdbV2 struct {
	*influxdbConnectionProducer

	usernameProducer template.StringTemplate
}

// New returns a new InfluxDBv2 instance
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *InfluxdbV2 {
	connProducer := &influxdbConnectionProducer{}
	connProducer.Type = influxdbTypeName

	return &InfluxdbV2{
		influxdbConnectionProducer: connProducer,
	}
}

// Type returns the TypeName for this backend
func (i *InfluxdbV2) Type() (string, error) {
	return influxdbTypeName, nil
}

func (i *InfluxdbV2) getConnection(ctx context.Context) (influxdb2.Client, error) {
	cli, err := i.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return cli.(influxdb2.Client), nil
}

func (i *InfluxdbV2) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (resp dbplugin.InitializeResponse, err error) {
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
	i.usernameProducer = up

	_, err = i.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	return i.influxdbConnectionProducer.Initialize(ctx, req)
}

// NewUser generates the username/password on the underlying Influxdb secret backend
func (i *InfluxdbV2) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (resp dbplugin.NewUserResponse, err error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	username, err := i.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	user, err := cli.UsersAPI().CreateUserWithName(ctx, username)
	if err != nil {
		// Attempt rollback only when the response has an error
		err2 := cli.UsersAPI().DeleteUser(ctx, user)
		if err2 != nil {
			return dbplugin.NewUserResponse{}, fmt.Errorf("failed to rollback query in InfluxDB: %w : %s", err, err2)
		}
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", err)
	}
	err = cli.UsersAPI().UpdateUserPassword(ctx, user, req.Password)
	if err != nil {
		// Attempt rollback only when the response has an error
		err2 := cli.UsersAPI().DeleteUser(ctx, user)
		if err2 != nil {
			return dbplugin.NewUserResponse{}, fmt.Errorf("failed to rollback query in InfluxDB: %w : %s", err, err2)
		}
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", err)
	}
	organization, err := cli.OrganizationsAPI().FindOrganizationByName(ctx, i.Organization)
	if err != nil {
		// Attempt rollback only when the response has an error
		err2 := cli.UsersAPI().DeleteUser(ctx, user)
		if err2 != nil {
			return dbplugin.NewUserResponse{}, fmt.Errorf("failed to rollback query in InfluxDB: %w : %s", err, err2)
		}
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", err)
	}
	_, err = cli.OrganizationsAPI().AddMember(ctx, organization, user)
	if err != nil {
		// Attempt rollback only when the response has an error
		err2 := cli.UsersAPI().DeleteUser(ctx, user)
		if err2 != nil {
			return dbplugin.NewUserResponse{}, fmt.Errorf("failed to rollback query in InfluxDB: %w : %s", err, err2)
		}
		return dbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", err)
	}
	resp = dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func deleteUser(ctx context.Context, cli influxdb2.Client, username string) error {
	user, err := cli.UsersAPI().FindUserByName(ctx, username)
	if err != nil {
		return err
	}
	err = cli.UsersAPI().DeleteUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (i *InfluxdbV2) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	err = deleteUser(ctx, cli, req.Username)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("failed to delete user cleanly: %w", err)
	}
	return dbplugin.DeleteUserResponse{}, nil
}

func (i *InfluxdbV2) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	i.Lock()
	defer i.Unlock()

	if req.Password != nil {
		err := i.changeUserPassword(ctx, req.Username, req.Password)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, fmt.Errorf("failed to change %q password: %w", req.Username, err)
		}
	}
	// Expiration is a no-op
	return dbplugin.UpdateUserResponse{}, nil
}

func (i *InfluxdbV2) changeUserPassword(ctx context.Context, username string, changePassword *dbplugin.ChangePassword) error {
	cli, err := i.getConnection(ctx)
	if err != nil {
		return fmt.Errorf("unable to get connection: %w", err)
	}
	user, err := cli.UsersAPI().FindUserByName(ctx, username)
	if err != nil {
		return err
	}
	err = cli.UsersAPI().UpdateUserPassword(ctx, user, changePassword.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

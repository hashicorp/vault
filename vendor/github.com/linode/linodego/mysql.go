package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type MySQLDatabaseTarget string

type MySQLDatabaseMaintenanceWindow = DatabaseMaintenanceWindow

const (
	MySQLDatabaseTargetPrimary   MySQLDatabaseTarget = "primary"
	MySQLDatabaseTargetSecondary MySQLDatabaseTarget = "secondary"
)

// A MySQLDatabase is an instance of Linode MySQL Managed Databases
type MySQLDatabase struct {
	ID              int                       `json:"id"`
	Status          DatabaseStatus            `json:"status"`
	Label           string                    `json:"label"`
	Hosts           DatabaseHost              `json:"hosts"`
	Region          string                    `json:"region"`
	Type            string                    `json:"type"`
	Engine          string                    `json:"engine"`
	Version         string                    `json:"version"`
	ClusterSize     int                       `json:"cluster_size"`
	ReplicationType string                    `json:"replication_type"`
	SSLConnection   bool                      `json:"ssl_connection"`
	Encrypted       bool                      `json:"encrypted"`
	AllowList       []string                  `json:"allow_list"`
	InstanceURI     string                    `json:"instance_uri"`
	Created         *time.Time                `json:"-"`
	Updated         *time.Time                `json:"-"`
	Updates         DatabaseMaintenanceWindow `json:"updates"`
}

func (d *MySQLDatabase) UnmarshalJSON(b []byte) error {
	type Mask MySQLDatabase

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	d.Updated = (*time.Time)(p.Updated)
	return nil
}

// MySQLCreateOptions fields are used when creating a new MySQL Database
type MySQLCreateOptions struct {
	Label           string   `json:"label"`
	Region          string   `json:"region"`
	Type            string   `json:"type"`
	Engine          string   `json:"engine"`
	AllowList       []string `json:"allow_list,omitempty"`
	ReplicationType string   `json:"replication_type,omitempty"`
	ClusterSize     int      `json:"cluster_size,omitempty"`
	Encrypted       bool     `json:"encrypted,omitempty"`
	SSLConnection   bool     `json:"ssl_connection,omitempty"`
}

// MySQLUpdateOptions fields are used when altering the existing MySQL Database
type MySQLUpdateOptions struct {
	Label     string                     `json:"label,omitempty"`
	AllowList *[]string                  `json:"allow_list,omitempty"`
	Updates   *DatabaseMaintenanceWindow `json:"updates,omitempty"`
}

// MySQLDatabaseBackup is information for interacting with a backup for the existing MySQL Database
type MySQLDatabaseBackup struct {
	ID      int        `json:"id"`
	Label   string     `json:"label"`
	Type    string     `json:"type"`
	Created *time.Time `json:"-"`
}

// MySQLBackupCreateOptions are options used for CreateMySQLDatabaseBackup(...)
type MySQLBackupCreateOptions struct {
	Label  string              `json:"label"`
	Target MySQLDatabaseTarget `json:"target"`
}

func (d *MySQLDatabaseBackup) UnmarshalJSON(b []byte) error {
	type Mask MySQLDatabaseBackup

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(d),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	d.Created = (*time.Time)(p.Created)
	return nil
}

// MySQLDatabaseCredential is the Root Credentials to access the Linode Managed Database
type MySQLDatabaseCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MySQLDatabaseSSL is the SSL Certificate to access the Linode Managed MySQL Database
type MySQLDatabaseSSL struct {
	CACertificate []byte `json:"ca_certificate"`
}

// ListMySQLDatabases lists all MySQL Databases associated with the account
func (c *Client) ListMySQLDatabases(ctx context.Context, opts *ListOptions) ([]MySQLDatabase, error) {
	response, err := getPaginatedResults[MySQLDatabase](ctx, c, "databases/mysql/instances", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListMySQLDatabaseBackups lists all MySQL Database Backups associated with the given MySQL Database
func (c *Client) ListMySQLDatabaseBackups(ctx context.Context, databaseID int, opts *ListOptions) ([]MySQLDatabaseBackup, error) {
	response, err := getPaginatedResults[MySQLDatabaseBackup](ctx, c, formatAPIPath("databases/mysql/instances/%d/backups", databaseID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetMySQLDatabase returns a single MySQL Database matching the id
func (c *Client) GetMySQLDatabase(ctx context.Context, databaseID int) (*MySQLDatabase, error) {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	response, err := doGETRequest[MySQLDatabase](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateMySQLDatabase creates a new MySQL Database using the createOpts as configuration, returns the new MySQL Database
func (c *Client) CreateMySQLDatabase(ctx context.Context, opts MySQLCreateOptions) (*MySQLDatabase, error) {
	e := "databases/mysql/instances"
	response, err := doPOSTRequest[MySQLDatabase](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteMySQLDatabase deletes an existing MySQL Database with the given id
func (c *Client) DeleteMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	err := doDELETERequest(ctx, c, e)
	return err
}

// UpdateMySQLDatabase updates the given MySQL Database with the provided opts, returns the MySQLDatabase with the new settings
func (c *Client) UpdateMySQLDatabase(ctx context.Context, databaseID int, opts MySQLUpdateOptions) (*MySQLDatabase, error) {
	e := formatAPIPath("databases/mysql/instances/%d", databaseID)
	response, err := doPUTRequest[MySQLDatabase](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetMySQLDatabaseSSL returns the SSL Certificate for the given MySQL Database
func (c *Client) GetMySQLDatabaseSSL(ctx context.Context, databaseID int) (*MySQLDatabaseSSL, error) {
	e := formatAPIPath("databases/mysql/instances/%d/ssl", databaseID)
	response, err := doGETRequest[MySQLDatabaseSSL](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetMySQLDatabaseCredentials returns the Root Credentials for the given MySQL Database
func (c *Client) GetMySQLDatabaseCredentials(ctx context.Context, databaseID int) (*MySQLDatabaseCredential, error) {
	e := formatAPIPath("databases/mysql/instances/%d/credentials", databaseID)
	response, err := doGETRequest[MySQLDatabaseCredential](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ResetMySQLDatabaseCredentials returns the Root Credentials for the given MySQL Database (may take a few seconds to work)
func (c *Client) ResetMySQLDatabaseCredentials(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/credentials/reset", databaseID)
	_, err := doPOSTRequest[MySQLDatabaseCredential, any](ctx, c, e)
	return err
}

// GetMySQLDatabaseBackup returns a specific MySQL Database Backup with the given ids
func (c *Client) GetMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) (*MySQLDatabaseBackup, error) {
	e := formatAPIPath("databases/mysql/instances/%d/backups/%d", databaseID, backupID)
	response, err := doGETRequest[MySQLDatabaseBackup](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RestoreMySQLDatabaseBackup returns the given MySQL Database with the given Backup
func (c *Client) RestoreMySQLDatabaseBackup(ctx context.Context, databaseID int, backupID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/backups/%d/restore", databaseID, backupID)
	_, err := doPOSTRequest[MySQLDatabaseBackup, any](ctx, c, e)
	return err
}

// CreateMySQLDatabaseBackup creates a snapshot for the given MySQL database
func (c *Client) CreateMySQLDatabaseBackup(ctx context.Context, databaseID int, opts MySQLBackupCreateOptions) error {
	e := formatAPIPath("databases/mysql/instances/%d/backups", databaseID)
	_, err := doPOSTRequest[MySQLDatabaseBackup](ctx, c, e, opts)
	return err
}

// PatchMySQLDatabase applies security patches and updates to the underlying operating system of the Managed MySQL Database
func (c *Client) PatchMySQLDatabase(ctx context.Context, databaseID int) error {
	e := formatAPIPath("databases/mysql/instances/%d/patch", databaseID)
	_, err := doPOSTRequest[MySQLDatabase, any](ctx, c, e)
	return err
}

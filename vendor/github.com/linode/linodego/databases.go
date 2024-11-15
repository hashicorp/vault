package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type (
	DatabaseEngineType           string
	DatabaseDayOfWeek            int
	DatabaseMaintenanceFrequency string
	DatabaseStatus               string
)

const (
	DatabaseMaintenanceDayMonday DatabaseDayOfWeek = iota + 1
	DatabaseMaintenanceDayTuesday
	DatabaseMaintenanceDayWednesday
	DatabaseMaintenanceDayThursday
	DatabaseMaintenanceDayFriday
	DatabaseMaintenanceDaySaturday
	DatabaseMaintenanceDaySunday
)

const (
	DatabaseMaintenanceFrequencyWeekly  DatabaseMaintenanceFrequency = "weekly"
	DatabaseMaintenanceFrequencyMonthly DatabaseMaintenanceFrequency = "monthly"
)

const (
	DatabaseEngineTypeMySQL    DatabaseEngineType = "mysql"
	DatabaseEngineTypePostgres DatabaseEngineType = "postgresql"
)

const (
	DatabaseStatusProvisioning DatabaseStatus = "provisioning"
	DatabaseStatusActive       DatabaseStatus = "active"
	DatabaseStatusDeleting     DatabaseStatus = "deleting"
	DatabaseStatusDeleted      DatabaseStatus = "deleted"
	DatabaseStatusSuspending   DatabaseStatus = "suspending"
	DatabaseStatusSuspended    DatabaseStatus = "suspended"
	DatabaseStatusResuming     DatabaseStatus = "resuming"
	DatabaseStatusRestoring    DatabaseStatus = "restoring"
	DatabaseStatusFailed       DatabaseStatus = "failed"
	DatabaseStatusDegraded     DatabaseStatus = "degraded"
	DatabaseStatusUpdating     DatabaseStatus = "updating"
	DatabaseStatusBackingUp    DatabaseStatus = "backing_up"
)

// A Database is a instance of Linode Managed Databases
type Database struct {
	ID              int            `json:"id"`
	Status          DatabaseStatus `json:"status"`
	Label           string         `json:"label"`
	Hosts           DatabaseHost   `json:"hosts"`
	Region          string         `json:"region"`
	Type            string         `json:"type"`
	Engine          string         `json:"engine"`
	Version         string         `json:"version"`
	ClusterSize     int            `json:"cluster_size"`
	ReplicationType string         `json:"replication_type"`
	SSLConnection   bool           `json:"ssl_connection"`
	Encrypted       bool           `json:"encrypted"`
	AllowList       []string       `json:"allow_list"`
	InstanceURI     string         `json:"instance_uri"`
	Created         *time.Time     `json:"-"`
	Updated         *time.Time     `json:"-"`
}

// DatabaseHost for Primary/Secondary of Database
type DatabaseHost struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary,omitempty"`
}

// DatabaseEngine is information about Engines supported by Linode Managed Databases
type DatabaseEngine struct {
	ID      string `json:"id"`
	Engine  string `json:"engine"`
	Version string `json:"version"`
}

// DatabaseMaintenanceWindow stores information about a MySQL cluster's maintenance window
type DatabaseMaintenanceWindow struct {
	DayOfWeek   DatabaseDayOfWeek            `json:"day_of_week"`
	Duration    int                          `json:"duration"`
	Frequency   DatabaseMaintenanceFrequency `json:"frequency"`
	HourOfDay   int                          `json:"hour_of_day"`
	WeekOfMonth *int                         `json:"week_of_month"`
}

// DatabaseType is information about the supported Database Types by Linode Managed Databases
type DatabaseType struct {
	ID          string                `json:"id"`
	Label       string                `json:"label"`
	Class       string                `json:"class"`
	VirtualCPUs int                   `json:"vcpus"`
	Disk        int                   `json:"disk"`
	Memory      int                   `json:"memory"`
	Engines     DatabaseTypeEngineMap `json:"engines"`
}

// DatabaseTypeEngineMap stores a list of Database Engine types by engine
type DatabaseTypeEngineMap struct {
	MySQL []DatabaseTypeEngine `json:"mysql"`
}

// DatabaseTypeEngine Sizes and Prices
type DatabaseTypeEngine struct {
	Quantity int          `json:"quantity"`
	Price    ClusterPrice `json:"price"`
}

// ClusterPrice for Hourly and Monthly price models
type ClusterPrice struct {
	Hourly  float32 `json:"hourly"`
	Monthly float32 `json:"monthly"`
}

func (d *Database) UnmarshalJSON(b []byte) error {
	type Mask Database

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

// ListDatabases lists all Database instances in Linode Managed Databases for the account
func (c *Client) ListDatabases(ctx context.Context, opts *ListOptions) ([]Database, error) {
	response, err := getPaginatedResults[Database](ctx, c, "databases/instances", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListDatabaseEngines lists all Database Engines. This endpoint is cached by default.
func (c *Client) ListDatabaseEngines(ctx context.Context, opts *ListOptions) ([]DatabaseEngine, error) {
	response, err := getPaginatedResults[DatabaseEngine](ctx, c, "databases/engines", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDatabaseEngine returns a specific Database Engine. This endpoint is cached by default.
func (c *Client) GetDatabaseEngine(ctx context.Context, _ *ListOptions, engineID string) (*DatabaseEngine, error) {
	e := formatAPIPath("databases/engines/%s", engineID)
	response, err := doGETRequest[DatabaseEngine](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListDatabaseTypes lists all Types of Database provided in Linode Managed Databases. This endpoint is cached by default.
func (c *Client) ListDatabaseTypes(ctx context.Context, opts *ListOptions) ([]DatabaseType, error) {
	response, err := getPaginatedResults[DatabaseType](ctx, c, "databases/types", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDatabaseType returns a specific Database Type. This endpoint is cached by default.
func (c *Client) GetDatabaseType(ctx context.Context, _ *ListOptions, typeID string) (*DatabaseType, error) {
	e := formatAPIPath("databases/types/%s", typeID)
	response, err := doGETRequest[DatabaseType](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

package azuresql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/management"
	azsql "github.com/Azure/azure-sdk-for-go/management/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathConfigLease(&b),
			pathConfigSubscription(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	azureClient *azsql.SQLDatabaseClient
	db          *sql.DB
	defaultDb   string
	server      string
	database    string
	lock        sync.Mutex
}

func (b *backend) AzureClient(s logical.Storage) (*azsql.SQLDatabaseClient, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If we already have a client, just return it!
	if b.azureClient != nil {
		return b.azureClient, nil
	}

	// Otherwise, attempt to make connection
	entry, err := s.Get("config/subscription")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("configure the Azure client connection with config/subscription first")
	}

	var azureConfig subscriptionConfig
	if err := entry.DecodeJSON(&azureConfig); err != nil {
		return nil, err
	}
	// Use the Azure Go SDK
	var client management.Client
	if len(azureConfig.PublishSettings) > 0 {
		client, err = management.ClientFromPublishSettingsFile(azureConfig.PublishSettings, azureConfig.SubscriptionID)
		if err != nil {
			return nil, err
		}
	} else {
		cert, err := ioutil.ReadFile(azureConfig.ManagementCert)
		if err != nil {
			return nil, err
		}
		client, err = management.NewClient(azureConfig.SubscriptionID, cert)
		if err != nil {
			return nil, err
		}
	}

	//Test out the client
	azClient := azsql.NewClient(client)
	_, err = azClient.GetDatabase(azureConfig.Server, azureConfig.Database)
	if err != nil {
		return nil, err
	}

	b.azureClient = &azClient
	b.server = azureConfig.Server
	b.database = azureConfig.Database
	return b.azureClient, nil
}

// DB returns the default database connection.
func (b *backend) DB(s logical.Storage) (*sql.DB, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If we already have a DB, we got it!
	if b.db != nil {
		return b.db, nil
	}

	// Otherwise, attempt to make connection
	entry, err := s.Get("config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("configure the DB connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}
	connString := connConfig.ConnectionString

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	db.SetMaxOpenConns(connConfig.MaxOpenConnections)

	//Testing out the connection
	stmt, err := db.Prepare("SELECT db_name();")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&b.defaultDb)
	if err != nil {
		return nil, err
	}

	b.db = db
	return b.db, nil
}

// ResetDB forces a connection next time DB() is called.
func (b *backend) ResetDB() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.db != nil {
		b.db.Close()
	}

	b.db = nil
}

// LeaseConfig returns the lease configuration
func (b *backend) LeaseConfig(s logical.Storage) (*configLease, error) {
	entry, err := s.Get("config/lease")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const backendHelp = `
The Azure SQL backend dynamically generates database users and 
optionally create firewall rules for the users.

After mounting this backend, configure it using the endpoints within
the "config/" path.

This backend only support Azure SQL Database V12 (contained user mode, 
which is the default for new Azure SQL Database setups).
`

package mssqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

const mssqlPassword = "yourStrong(!)Password"

// This constant is used in retrying the mssql container restart, since
// intermittently the container starts but mssql within the container
// is unreachable.
const numRetries = 3

func PrepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	var err error
	for i := 0; i < numRetries; i++ {
		var svc *docker.Service
		runner, err := docker.NewServiceRunner(docker.RunOptions{
			ContainerName: "sqlserver",
			ImageRepo:     "mcr.microsoft.com/mssql/server",
			ImageTag:      "2017-latest-ubuntu",
			Env:           []string{"ACCEPT_EULA=Y", "SA_PASSWORD=" + mssqlPassword},
			Ports:         []string{"1433/tcp"},
			LogConsumer: func(s string) {
				if t.Failed() {
					t.Logf("container logs: %s", s)
				}
			},
		})
		if err != nil {
			t.Fatalf("Could not start docker MSSQL: %s", err)
		}

		svc, err = runner.StartService(context.Background(), connectMSSQL)
		if err == nil {
			return svc.Cleanup, svc.Config.URL().String()
		}
	}

	t.Fatalf("Could not start docker MSSQL: %s", err)
	return nil, ""
}

func connectMSSQL(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("sa", mssqlPassword),
		Host:   fmt.Sprintf("%s:%d", host, port),
	}
	// Attempt to address connection flakiness within tests such as "Failed to initialize: error verifying connection ..."
	u.Query().Add("Connection Timeout", "30")

	db, err := sql.Open("mssql", u.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return docker.NewServiceURL(u), nil
}

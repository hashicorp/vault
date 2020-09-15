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

func PrepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "sqlserver",
		ImageRepo:     "mcr.microsoft.com/mssql/server",
		ImageTag:      "2017-latest-ubuntu",
		Env:           []string{"ACCEPT_EULA=Y", "SA_PASSWORD=" + mssqlPassword},
		Ports:         []string{"1433/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker MSSQL: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectMSSQL)
	if err != nil {
		t.Fatalf("Could not start docker MSSQL: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectMSSQL(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("sa", mssqlPassword),
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

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

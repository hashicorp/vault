// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mssqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

const mssqlPassword = "yourStrong(!)Password"

// This constant is used in retrying the mssql container restart, since
// intermittently the container starts but mssql within the container
// is unreachable.
const numRetries = 3

func PrepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	logger := corehelpers.NewTestLogger(t)

	var err error
	for i := 0; i < numRetries; i++ {
		var svc *docker.Service
		var runner *docker.Runner
		runner, err = docker.NewServiceRunner(docker.RunOptions{
			ContainerName: "sqlserver",
			ImageRepo:     "mcr.microsoft.com/mssql/server",
			ImageTag:      "2022-latest",
			Env:           []string{"ACCEPT_EULA=Y", "SA_PASSWORD=" + mssqlPassword},
			Ports:         []string{"1433/tcp"},
			LogConsumer: func(s string) {
				logger.Info(s)
			},
		})
		if err != nil {
			logger.Error("failed creating new service runner", "error", err.Error())
			continue
		}

		svc, err = runner.StartService(context.Background(), connectMSSQL)
		if err == nil {
			return svc.Cleanup, svc.Config.URL().String()
		}

		logger.Error("failed starting service", "error", err.Error())
	}

	t.Fatalf("Could not start docker MSSQL last error: %v", err)
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

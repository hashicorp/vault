// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mssqlhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

const mssqlPassword = "yourStrong(!)Password"

// This constant is used in retrying the mssql container restart, since
// intermittently the container starts but mssql within the container
// is unreachable.
const numRetries = 3

// SQL Server accepts sa logins before startup recovery, which includes system
// database upgrades, has finished; statements issued in that window can stall.
// SQL Server writes this line to its log once that startup work is done.
const recoveryCompleteLog = "Recovery is complete."

// PrepareMSSQLTestContainer returns a connection URL for a SQL Server that tests
// can use, and a cleanup function. If MSSQL_URL is set, that server is used as is.
// Otherwise it starts a fresh SQL Server container and returns only after the
// server has finished starting up, not merely after it first accepts logins.
func PrepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	logger := corehelpers.NewTestLogger(t)

	// Workaround for https://github.com/microsoft/mssql-docker/issues/895 and us temporary seeing
	// tls: failed to parse certificate from server: x509: negative serial number in test case failures.
	containerfile := `
FROM mcr.microsoft.com/mssql/server:2022-latest
USER root
ENV MSDIR=/var/opt/mssql
RUN mkdir -p $MSDIR \
    && openssl req -x509 -nodes -newkey rsa:2048 -subj '/CN=mssql' -addext "subjectAltName = DNS:mssql" -keyout $MSDIR/mssql.key -out $MSDIR/mssql.pem -days 1 \
	&& chmod 400 $MSDIR/mssql.key \
	&& chmod 400 $MSDIR/mssql.pem \
    && chown -R mssql $MSDIR

RUN echo "[network]" > $MSDIR/mssql.conf \
	&& echo "tlscert = $MSDIR/mssql.pem" >> $MSDIR/mssql.conf \
	&& echo "tlskey = $MSDIR/mssql.key" >> $MSDIR/mssql.conf \ 
	&& echo "tlsprotocols = 1.2" >> $MSDIR/mssql.conf \ 
	&& echo "forceencryption = 1" >> $MSDIR/mssql.conf 

USER mssql
`
	bCtx := docker.NewBuildContext()
	imageName := "mssql-workaround-895"
	imageTag := "latest"

	// runOptions describes the container to start. Its LogConsumer receives every
	// line the container logs and sets recovered once SQL Server reports that
	// startup recovery is complete.
	runOptions := func(recovered *atomic.Bool) docker.RunOptions {
		return docker.RunOptions{
			ContainerName: "sqlserver",
			ImageRepo:     imageName,
			ImageTag:      imageTag,
			Env:           []string{"ACCEPT_EULA=Y", "SA_PASSWORD=" + mssqlPassword},
			Ports:         []string{"1433/tcp"},
			LogConsumer: func(s string) {
				logger.Info(s)
				if strings.Contains(s, recoveryCompleteLog) {
					recovered.Store(true)
				}
			},
		}
	}

	// This runner only builds the image, so its recovered flag is never read.
	runner, err := docker.NewServiceRunner(runOptions(new(atomic.Bool)))
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	// Sometimes we see timeouts and issues pulling the zlint code from GitHub
	testhelpers.RetryUntil(t, 30*time.Second, func() error {
		output, err := runner.BuildImage(context.Background(), containerfile, bCtx,
			docker.BuildRemove(true),
			docker.BuildForceRemove(true),
			docker.BuildPullParent(true),
			docker.BuildTags([]string{imageName + ":" + imageTag}))
		if err != nil {
			return fmt.Errorf("could not build new mssql image: %w", err)
		}

		t.Logf("Image build output: %v", string(output))
		return nil
	})

	for i := 0; i < numRetries; i++ {
		var svc *docker.Service

		// A runner per attempt gives each container its own log consumer, so
		// output from a replaced container cannot mark the next one as recovered.
		recovered := new(atomic.Bool)
		runner, err = docker.NewServiceRunner(runOptions(recovered))
		if err != nil {
			t.Fatalf("Could not provision docker service runner: %s", err)
		}
		// StartService starts the container, then keeps calling this readiness
		// check with backoff (up to about two minutes) until it returns no error.
		// The server counts as ready only after it has logged that recovery is
		// complete and an sa login succeeds.
		svc, err = runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
			if !recovered.Load() {
				return nil, errors.New("SQL Server has not completed startup recovery")
			}
			return connectMSSQL(ctx, host, port)
		})
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

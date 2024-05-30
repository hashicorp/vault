// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mysqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/docker"
)

type Config struct {
	docker.ServiceHostPort
	ConnString string
}

var _ docker.ServiceConfig = &Config{}

func PrepareTestContainer(t *testing.T, legacy bool, pw string) (func(), string) {
	if os.Getenv("MYSQL_URL") != "" {
		return func() {}, os.Getenv("MYSQL_URL")
	}

	// ARM64 is only supported on MySQL 8.0 and above. If we update
	// our image and support to 8.0, we can unskip these tests.
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as MySQL 5.7 is not supported on ARM architectures")
	}

	imageVersion := "5.7"
	if legacy {
		imageVersion = "5.6"
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "mysql",
		ImageRepo:     "docker.mirror.hashicorp.services/library/mysql",
		ImageTag:      imageVersion,
		Ports:         []string{"3306/tcp"},
		Env:           []string{"MYSQL_ROOT_PASSWORD=" + pw},
	})
	if err != nil {
		t.Fatalf("could not start docker mysql: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		hostIP := docker.NewServiceHostPort(host, port)
		connString := fmt.Sprintf("root:%s@tcp(%s)/mysql?parseTime=true", pw, hostIP.Address())
		db, err := sql.Open("mysql", connString)
		if err != nil {
			return nil, err
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			return nil, err
		}
		return &Config{ServiceHostPort: *hostIP, ConnString: connString}, nil
	})
	if err != nil {
		t.Fatalf("could not start docker mysql: %s", err)
	}
	return svc.Cleanup, svc.Config.(*Config).ConnString
}

func TestCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	connURL = strings.Replace(connURL, "root:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

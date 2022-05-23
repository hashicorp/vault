package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	return prepareTestContainer(t, version, "secret", "database")
}

func PrepareTestContainerWithPassword(t *testing.T, version, password string) (func(), string) {
	return prepareTestContainer(t, version, password, "database")
}

func prepareTestContainer(t *testing.T, version, password, db string) (func(), string) {
	if os.Getenv("PG_URL") != "" {
		return func() {}, os.Getenv("PG_URL")
	}

	if version == "" {
		version = "11"
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "postgres",
		ImageTag:  version,
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		Ports: []string{"5432/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectPostgres(password))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectPostgres(password string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}

		db, err := sql.Open("pgx", u.String())
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
}

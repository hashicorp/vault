package neo4j

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/url"
	"os"
	"testing"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	if os.Getenv("NEO4J_URL") != "" {
		return func() {}, os.Getenv("NEO4J_URL")
	}

	if version == "" {
		version = "enterprise"
	}
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "neo4j",
		ImageTag:  version,
		Env:       []string{"NEO4J_ACCEPT_LICENSE_AGREEMENT=yes", "NEO4J_AUTH=neo4j/secret"},
		Ports:     []string{"7687/tcp", "7474/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker Neo4j: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectNeo4j)
	if err != nil {
		t.Fatalf("Could not start docker Neo4j: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectNeo4j(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "neo4j",
		Host:   fmt.Sprintf("%s:%d", host, port),
		User:   url.UserPassword("neo4j", "secret"),
	}
	pass, _ := u.User.Password()
	driver, err := neo4j.NewDriverWithContext(u.String(), neo4j.BasicAuth(u.User.Username(), pass, ""))
	if err != nil {
		return nil, err
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not verify neo4j connection: %w", err)
	}
	return docker.NewServiceURL(u), nil
}

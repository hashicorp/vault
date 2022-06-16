package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/url"
	"testing"
)

func TestContainer(t *testing.T) {
	cleanup, urlStr := PrepareTestContainer(t, "enterprise")
	defer cleanup()
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err != nil {
		t.Fatal(err)
	}
	pass, _ := u.User.Password()
	driver, err := neo4j.NewDriverWithContext(u.String(), neo4j.BasicAuth(u.User.Username(), pass, ""))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

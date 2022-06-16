package neo4j

import (
	"context"
	"testing"
)

func TestInit_Simple(t *testing.T) {
	tags := []string{"enterprise", "latest"}
	ctx := context.Background()
	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			db, cleanup := getNeo4j(t, nil, tag)
			defer cleanup()
			driver, err := db.getConnection(ctx)
			if err != nil {
				t.Fatal(err)
			}

			err = driver.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

package physical

import (
	"fmt"
	"os"
	"testing"

	_ "gopkg.in/mgo.v2"
)

func TestMongoBackend(t *testing.T) {
	address := os.Getenv("MONGO_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MONGO_DB")
	if database == "" {
		database = "vaultTest"
	}

	collection := os.Getenv("MONGO_COLLECTION")
	if collection == "" {
		collection = "testVaultCollection"
	}

	username := os.Getenv("MONGO_USERNAME")
	if username == "" {
		fmt.Println("env variable $MONGO_USERNAME not set")
	}
	password := os.Getenv("MONGO_PASSWORD")
	if password == "" {
		fmt.Println("env variable $MONGO_PASSWORD not set")
	}

	// Run vault tests
	b, err := NewBackend("mongodb", map[string]string{
		"address":  	address,
		"database": 	database,
		"collection": collection,
		"username": 	username,
		"password": 	password,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	mongodb := b.(*MongoBackend)

	defer func() {
		err := mongodb.session.DB(mongodb.database).DropDatabase()
		if err != nil {
			t.Fatalf("Failed to drop database: %v", err)
		}
	}()

	testBackend(t, mongodb)
	testBackend_ListPrefix(t, mongodb)

}

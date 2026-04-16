// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultMongoImage   = "docker.mirror.hashicorp.services/mongo"
	defaultMongoVersion = "7.0"
	defaultMongoUser    = "admin"
	defaultMongoPass    = "secret"
)

// defaultRunOpts returns default Docker run options for MongoDB container
// Uses test name to ensure unique container names for parallel execution
func defaultRunOpts(t *testing.T) docker.RunOptions {
	return docker.RunOptions{
		ContainerName: fmt.Sprintf("mongodb-%s", sanitize(t.Name())),
		ImageRepo:     defaultMongoImage,
		ImageTag:      defaultMongoVersion,
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=" + defaultMongoUser,
			"MONGO_INITDB_ROOT_PASSWORD=" + defaultMongoPass,
			"MONGO_INITDB_DATABASE=admin",
		},
		Ports:             []string{"27017/tcp"},
		DoNotAutoRemove:   false,
		OmitLogTimestamps: true,
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	}
}

// requireVaultEnv skips the test if required Vault environment variables are not set
func requireVaultEnv(t *testing.T) {
	t.Helper()

	if os.Getenv("VAULT_ADDR") == "" || os.Getenv("VAULT_TOKEN") == "" {
		t.Skip("skipping blackbox test: VAULT_ADDR and VAULT_TOKEN are required")
	}
}

// sanitize converts test name to a valid container name
// Removes special characters and converts to lowercase
func sanitize(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	// Remove any remaining special characters
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// PrepareTestContainer starts a MongoDB container for testing
// Returns cleanup function and connection URL
// If MONGO_URL environment variable is set, uses that instead of starting a container
func PrepareTestContainer(t *testing.T) (func(), string) {
	_, cleanup, connURL, _ := prepareTestContainer(t, defaultRunOpts(t), defaultMongoPass, true, false)
	return cleanup, connURL
}

// prepareTestContainer is the internal function that handles container setup
// Supports both Docker container creation and external MongoDB via environment variable
func prepareTestContainer(
	t *testing.T,
	runOpts docker.RunOptions,
	password string,
	addSuffix bool,
	forceLocalAddr bool,
) (*docker.Runner, func(), string, string) {
	requireVaultEnv(t)

	// Check for external MongoDB URL
	if os.Getenv("MONGO_URL") != "" {
		envMongoURL := os.Getenv("MONGO_URL")

		// Create unique database for this test
		dbName := fmt.Sprintf("test_%s_%d", sanitize(t.Name()), time.Now().Unix())
		testURL := replaceDatabase(envMongoURL, dbName)

		// Create the database
		if err := createDatabase(t, envMongoURL, dbName); err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}

		cleanup := func() {
			dropDatabase(t, envMongoURL, dbName)
		}

		return nil, cleanup, testURL, ""
	}

	// Start Docker container
	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "docker") &&
			(strings.Contains(errStr, "daemon") || strings.Contains(errStr, "connect")) {
			t.Skipf("skipping blackbox test: docker not available: %v", err)
		}
		t.Fatalf("Could not start docker MongoDB: %s", err)
	}

	svc, containerID, err := runner.StartNewService(
		context.Background(),
		addSuffix,
		forceLocalAddr,
		connectMongoDB(password),
	)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "docker") &&
			(strings.Contains(errStr, "daemon") || strings.Contains(errStr, "connect")) {
			t.Skipf("skipping blackbox test: docker not available: %v", err)
		}
		t.Fatalf("Could not start docker MongoDB: %s", err)
	}

	connURL := svc.Config.URL().String()

	// Create unique database for this test
	dbName := fmt.Sprintf("test_%s_%d", sanitize(t.Name()), time.Now().Unix())
	testURL := replaceDatabase(connURL, dbName)

	// Create the database
	if err := createDatabase(t, connURL, dbName); err != nil {
		svc.Cleanup()
		t.Fatalf("Failed to create test database: %v", err)
	}

	cleanup := func() {
		dropDatabase(t, connURL, dbName)
		svc.Cleanup()
	}

	return runner, cleanup, testURL, containerID
}

// connectMongoDB returns a ServiceAdapter that connects to MongoDB
// Includes retry logic with 30-second timeout for container startup
func connectMongoDB(password string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme: "mongodb",
			User:   url.UserPassword(defaultMongoUser, password),
			Host:   fmt.Sprintf("%s:%d", host, port),
			Path:   "/admin",
		}

		// Retry connection with timeout
		deadline := time.Now().Add(30 * time.Second)
		var lastErr error

		for time.Now().Before(deadline) {
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(u.String()))
			if err != nil {
				lastErr = err
				time.Sleep(1 * time.Second)
				continue
			}

			// Ping to verify connection
			if err = client.Ping(ctx, nil); err != nil {
				client.Disconnect(ctx)
				lastErr = err
				time.Sleep(1 * time.Second)
				continue
			}

			// Connection successful
			client.Disconnect(ctx)
			return docker.NewServiceURL(u), nil
		}

		return nil, fmt.Errorf("mongodb not ready after 30s: %w", lastErr)
	}
}

// replaceDatabase replaces the database name in a MongoDB connection URL
func replaceDatabase(connURL, dbName string) string {
	u, err := url.Parse(connURL)
	if err != nil {
		return connURL
	}

	u.Path = "/" + dbName
	return u.String()
}

// createDatabase creates a new database in MongoDB
// MongoDB creates databases lazily, so we create a collection to ensure it exists
func createDatabase(t *testing.T, connURL, dbName string) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	// Create a collection to ensure database exists
	db := client.Database(dbName)
	if err := db.CreateCollection(ctx, "_init"); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	t.Logf("Created test database: %s", dbName)
	return nil
}

// dropDatabase drops a database from MongoDB
func dropDatabase(t *testing.T, connURL, dbName string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		t.Logf("Warning: failed to connect for cleanup: %v", err)
		return
	}
	defer client.Disconnect(ctx)

	if err := client.Database(dbName).Drop(ctx); err != nil {
		t.Logf("Warning: failed to drop database %s: %v", dbName, err)
		return
	}

	t.Logf("Dropped test database: %s", dbName)
}

// mongoConnectionConfigPayload returns a standard MongoDB connection configuration payload
func mongoConnectionConfigPayload(connURL, allowedRoles string, verifyConnection bool) map[string]any {
	return map[string]any{
		"plugin_name":       "mongodb-database-plugin",
		"connection_url":    connURL,
		"allowed_roles":     allowedRoles,
		"verify_connection": verifyConnection,
	}
}

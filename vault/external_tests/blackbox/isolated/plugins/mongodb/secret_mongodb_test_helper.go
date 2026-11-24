// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultMongoImage   = "docker.mirror.hashicorp.services/mongo"
	defaultMongoVersion = "7.0"
	defaultMongoUser    = "admin"
	defaultMongoPass    = "secret"

	testConnectionName  = "my-mongodb-db"
	testInitialPassword = "initialpass"
	testRotationPeriod  = 86400 // 24 hours in seconds
)

// defaultRunOpts returns default Docker run options for MongoDB container
// Uses test name to ensure unique container names for parallel execution
func defaultRunOpts(t *testing.T) docker.RunOptions {
	return docker.RunOptions{
		ContainerName: fmt.Sprintf("mongo-%s", sanitize(t.Name())),
		ImageRepo:     defaultMongoImage,
		ImageTag:      defaultMongoVersion,
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=" + defaultMongoUser,
			"MONGO_INITDB_ROOT_PASSWORD=" + defaultMongoPass,
			"MONGO_INITDB_DATABASE=admin",
		},
		Ports:             []string{"27017/tcp"},
		DoNotAutoRemove:   false,
		PreDelete:         true,
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

var sanitizeRegex = regexp.MustCompile(`[^a-z0-9]+`)

// sanitize converts test name to a valid identifier with smart truncation
// Replaces non-alphanumeric characters with dashes and truncates long names
// with a hash suffix for uniqueness
func sanitize(name string) string {
	lower := strings.ToLower(name)
	out := sanitizeRegex.ReplaceAllString(lower, "-")
	out = strings.Trim(out, "-")

	if out == "" {
		return "test"
	}

	// Truncate long names with hash suffix for uniqueness
	if len(out) > 54 {
		const hashLen = 8
		sum := sha256.Sum256([]byte(out))
		hash := hex.EncodeToString(sum[:])[:hashLen]
		prefixLen := 54 - 1 - hashLen
		out = out[:prefixLen] + "-" + hash
	}

	return out
}

// PrepareTestContainer starts a MongoDB container for testing.
// Returns cleanup function, Vault connection URL, test runner connection URL,
// and the generated per-test database name.
// If MONGO_URL environment variable is set, uses that instead of starting a container.
// In CI: vaultURL is private (same VPC), testRunnerURL is public (different VPC)
func PrepareTestContainer(t *testing.T) (func(), string, string, string) {
	_, cleanup, vaultURL, testRunnerURL, dbName := prepareTestContainer(t, defaultRunOpts(t), defaultMongoPass, false, true)
	return cleanup, vaultURL, testRunnerURL, dbName
}

// prepareTestContainer is the internal function that handles container setup
// Supports both Docker container creation and external MongoDB via environment variable
func prepareTestContainer(
	t *testing.T,
	runOpts docker.RunOptions,
	password string,
	addSuffix bool,
	forceLocalAddr bool,
) (*docker.Runner, func(), string, string, string) {
	requireVaultEnv(t)

	// Check for external MongoDB URL
	if os.Getenv("MONGO_URL") != "" {
		envMongoURL := os.Getenv("MONGO_URL")

		// Use private URL for Vault, fall back to public
		vaultMongoURL := os.Getenv("MONGO_URL_PRIVATE")
		if vaultMongoURL == "" {
			vaultMongoURL = envMongoURL
		}

		// Create unique database for this test (max 63 chars for MongoDB)
		sanitized := sanitize(t.Name())
		timestamp := time.Now().Unix()
		// Format: test_{name}_{ts} - ensure total length <= 63
		// Reserve 5 for "test_", 10 for timestamp, 1 for underscore = 16 chars overhead
		maxNameLen := 63 - 16
		if len(sanitized) > maxNameLen {
			sanitized = sanitized[:maxNameLen]
		}
		dbName := fmt.Sprintf("test_%s_%d", sanitized, timestamp)

		// Vault uses private URL (same VPC), test runner uses public URL (different VPC)
		vaultURL := replaceDatabase(vaultMongoURL, dbName)
		testRunnerURL := replaceDatabase(envMongoURL, dbName)

		// Create the database (test runner uses public URL)
		if err := createDatabase(t, envMongoURL, dbName); err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}

		cleanup := func() {
			dropDatabase(t, envMongoURL, dbName)
		}

		return nil, cleanup, vaultURL, testRunnerURL, dbName
	}

	// Start Docker container
	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			t.Fatalf("skipping blackbox test: docker daemon not available: %v", err)
		}
		t.Fatalf("Could not start docker MongoDB: %s", err)
	}

	// Retry StartNewService with small delays to handle port mapping timing
	var svc *docker.Service
	for attempt := 0; attempt < 5; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}

		svc, _, err = runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectMongoDB(password))
		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "no port mapping found") {
			break
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			t.Fatalf("skipping blackbox test: docker daemon not available: %v", err)
		}
		t.Fatalf("Could not start docker MongoDB: %s", err)
	}

	connURL := svc.Config.URL().String()

	// Create unique database for this test (max 63 chars for MongoDB)
	sanitized := sanitize(t.Name())
	timestamp := time.Now().Unix()
	// Format: test_{name}_{ts} - ensure total length <= 63
	// Reserve 5 for "test_", 10 for timestamp, 1 for underscore = 16 chars overhead
	maxNameLen := 63 - 16
	if len(sanitized) > maxNameLen {
		sanitized = sanitized[:maxNameLen]
	}
	dbName := fmt.Sprintf("test_%s_%d", sanitized, timestamp)
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

	// For Docker, both Vault and test runner use the same URL (localhost)
	return runner, cleanup, testURL, testURL, dbName
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

	u.Path = dbName

	// Ensure authSource=admin so authentication works against admin database
	// even when connecting to a different database
	q := u.Query()
	if q.Get("authSource") == "" {
		q.Set("authSource", "admin")
		u.RawQuery = q.Encode()
	}

	return u.String()
}

// createDatabase creates a new database in MongoDB
// MongoDB creates databases lazily, so we create a collection to ensure it exists
// Uses 60-second timeout with retry logic for CI environments with network latency.
func createDatabase(t *testing.T, connURL, dbName string) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Retry connection with backoff for CI network latency
	var client *mongo.Client
	var err error
	deadline := time.Now().Add(60 * time.Second)

	for time.Now().Before(deadline) {
		// Set shorter server selection timeout (10s) to allow multiple retries within 60s window
		clientOpts := options.Client().
			ApplyURI(connURL).
			SetServerSelectionTimeout(10 * time.Second).
			SetConnectTimeout(10 * time.Second)

		client, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		// Verify connection with ping
		if err = client.Ping(ctx, nil); err != nil {
			client.Disconnect(ctx)
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB after retries: %w", err)
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
// Uses 60-second timeout with retry logic for CI environments with network latency.
func dropDatabase(t *testing.T, connURL, dbName string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Retry connection with backoff for CI network latency
	var client *mongo.Client
	var err error
	deadline := time.Now().Add(60 * time.Second)

	for time.Now().Before(deadline) {
		// Set shorter server selection timeout (10s) to allow multiple retries within 60s window
		clientOpts := options.Client().
			ApplyURI(connURL).
			SetServerSelectionTimeout(10 * time.Second).
			SetConnectTimeout(10 * time.Second)

		client, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		// Verify connection with ping
		if err = client.Ping(ctx, nil); err != nil {
			client.Disconnect(ctx)
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if err != nil {
		t.Logf("Warning: failed to connect for cleanup after retries: %v", err)
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

// setupMongoDBTest performs common test setup: creates container, enables mount, configures connection.
// Returns mount path, connection URL for verifying credentials, the generated
// per-test database name, and a reusable MongoDB client for test operations.
// In CI: Vault uses private URL (same VPC), test runner uses public URL (different VPC)
func setupMongoDBTest(t *testing.T, v *blackbox.Session) (string, string, string, *mongo.Client) {
	t.Helper()

	requireVaultEnv(t)
	cleanup, vaultURL, testRunnerURL, dbName := PrepareTestContainer(t)
	t.Cleanup(cleanup)

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	// Vault uses private URL (same VPC as MongoDB in CI)
	v.MustWrite(
		mount+"/config/"+testConnectionName,
		mongoConnectionConfigPayload(vaultURL, "*", false),
	)

	// Test runner uses public URL (different VPC in CI, needs public access)
	client := getMongoClient(t, testRunnerURL)
	t.Cleanup(func() {
		if err := client.Disconnect(context.Background()); err != nil {
			t.Logf("Warning: failed to disconnect MongoDB client: %v", err)
		}
	})

	// Return testRunnerURL for credential verification (test runner needs public access)
	return mount, testRunnerURL, dbName, client
}

// getMongoClient creates a MongoDB client with optimized settings for CI environments.
// Uses connection pooling and shorter timeouts for faster failure detection.
func getMongoClient(t *testing.T, connURL string) *mongo.Client {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(connURL).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(10).
		SetMinPoolSize(2)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("failed to create MongoDB client: %v", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		t.Fatalf("failed to ping MongoDB: %v", err)
	}

	return client
}

// createMongoDBUser creates a MongoDB user for testing static roles.
// Uses the provided client connection to avoid connection overhead.
// The user is created in the database identified by dbName, while auth
// continues to use authSource=admin.
func createMongoDBUser(t *testing.T, client *mongo.Client, dbName, username, password string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(dbName)

	// First, try to drop the user if it exists (for test cleanup/retry scenarios).
	_ = db.RunCommand(ctx, bson.D{{Key: "dropUser", Value: username}}).Err()

	err := db.RunCommand(ctx, bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password},
		{Key: "roles", Value: bson.A{
			bson.D{
				{Key: "role", Value: "readWrite"},
				{Key: "db", Value: dbName},
			},
		}},
	}).Err()
	if err != nil {
		t.Fatalf("failed to create MongoDB user: %v", err)
	}

	t.Logf("Created MongoDB user: %s in database %s", username, dbName)
}

// verifyMongoDBCredentials verifies that the given credentials work for MongoDB.
// Creates a new client with the provided credentials to test authentication.
func verifyMongoDBCredentials(t *testing.T, connURL, username, password string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Replace credentials in connection URL
	u, err := url.Parse(connURL)
	if err != nil {
		t.Fatalf("failed to parse connection URL: %v", err)
	}
	u.User = url.UserPassword(username, password)

	// Remove the auth source so we verify with the default database
	values := u.Query()
	if values.Get("authSource") != "" {
		values.Del("authSource")
		u.RawQuery = values.Encode()
	}

	// Create client with new credentials
	clientOpts := options.Client().
		ApplyURI(u.String()).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v: url %s", err, u.String())
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("failed to ping mongodb: %v: url %s", err, u.String())
	}

	t.Logf("Verified MongoDB credentials for user: %s", username)
}

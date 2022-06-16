package neo4j

import (
	"context"
	neo4jtest "github.com/hashicorp/vault/helper/testhelpers/neo4j"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"net/url"
	"testing"
	"time"
)

func getUsernamePasswordFromUrl(t *testing.T, url2 string) (string, string) {
	u, err := url.Parse(url2)
	if err != nil {
		t.Fatal(err)
	}
	pass, _ := u.User.Password()
	return u.User.Username(), pass
}

func getNeo4j(t *testing.T, options map[string]interface{}, tag string) (*Neo4j, func()) {
	cleanup, connURL := neo4jtest.PrepareTestContainer(t, tag)

	username, password := getUsernamePasswordFromUrl(t, connURL)

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"username":       username,
		"password":       password,
	}
	for k, v := range options {
		connectionDetails[k] = v
	}

	neo := newDB()

	_, err := neo.Initialize(context.Background(), dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !neo.Initialized {
		t.Fatal("Database should be initialized")
	}
	return neo, cleanup
}

func cloneConfig(c *ConnectionProducer) func(c *neo4j.Config) {
	config := &neo4j.Config{
		TlsConfig:                    c.tlsConf,
		MaxTransactionRetryTime:      c.maxTransactionRetryTime,
		MaxConnectionPoolSize:        c.MaxConnectionPoolSize,
		MaxConnectionLifetime:        c.maxConnectionLifetime,
		ConnectionAcquisitionTimeout: c.connectionAcquisitionTimeout,
		SocketConnectTimeout:         c.socketConnectTimeout,
		SocketKeepalive:              true,
		UserAgent:                    neo4j.UserAgent,
		FetchSize:                    neo4j.FetchDefault,
	}

	return func(c *neo4j.Config) {
		*c = *config
	}
}

func TestUpdateUser(t *testing.T) {
	tags := []string{"enterprise", "latest"}
	ctx := context.Background()
	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			neo, cleanup := getNeo4j(t, nil, tag)
			defer cleanup()
			driver, err := neo.getConnection(ctx)
			if err != nil {
				t.Fatal(err)
			}
			session := driver.NewSession(ctx, neo4j.SessionConfig{})
			defer session.Close(ctx)

			_, err = session.ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (interface{}, error) {
					return tx.Run(ctx, "CREATE OR REPLACE USER testuser SET PLAINTEXT PASSWORD '123'", nil)
				})
			if err != nil {
				t.Fatal(err)
			}

			username := "testuser"
			password := "456"
			req := dbplugin.UpdateUserRequest{
				Username:       username,
				CredentialType: dbplugin.CredentialTypePassword,
				Password: &dbplugin.ChangePassword{
					NewPassword: password,
					Statements:  dbplugin.Statements{},
				},
			}
			_, err = neo.UpdateUser(ctx, req)

			if err != nil {
				t.Fatal(err)
			}

			testConnOld, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth("testuser", "123", ""), cloneConfig(neo.ConnectionProducer))

			// Auth doesn't actually get hit until you VerifyConnectivity()
			if err != nil {
				t.Fatal(err)
			}

			defer testConnOld.Close(ctx)

			err = testConnOld.VerifyConnectivity(ctx)

			if err == nil {
				t.Fatalf("expected old password to fail, %s:%s", username, password)
			} else {
				dbError, isDbError := err.(*db.Neo4jError)
				if !isDbError {
					t.Fatalf("expected *db.Neo4jError, got %T", err)
				}
				if !dbError.IsAuthenticationFailed() {
					t.Fatalf("expected error to be Neo.ClientError.Security.Unauthorized, but got %s", err)
				}
			}

			testConnNew, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(username, password, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConnNew.Close(ctx)

			err = testConnNew.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestNewUser(t *testing.T) {
	tags := []string{"enterprise", "latest"}
	ctx := context.Background()

	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			neo, cleanup := getNeo4j(t, nil, tag)
			defer cleanup()

			password := "12345"

			req := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "arse", RoleName: "arse2",
				},
				Statements: dbplugin.Statements{
					Commands: []string{"CREATE OR REPLACE USER $username SET PLAINTEXT PASSWORD $password"},
				},
				RollbackStatements: dbplugin.Statements{},
				CredentialType:     dbplugin.CredentialTypePassword,
				Password:           password,
				Expiration:         time.Now().Add(time.Minute),
			}

			resp, err := neo.NewUser(ctx, req)
			if err != nil {
				t.Fatal(err)
			}

			testConn, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(resp.Username, password, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConn.Close(ctx)

			err = testConn.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatal(err)
			}

			// Check that statements which add DB data are successful too (executed non-transactionally)
			req2 := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "arse", RoleName: "arse2",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE OR REPLACE USER $username SET PLAINTEXT PASSWORD $password`,
						`CREATE (n:NEO4J_USER {username: $username})`,
					},
				},
				RollbackStatements: dbplugin.Statements{},
				CredentialType:     dbplugin.CredentialTypePassword,
				Password:           password,
				Expiration:         time.Now().Add(time.Minute),
			}

			resp, err = neo.NewUser(ctx, req2)
			if err != nil {
				t.Fatal(err)
			}

			testConn2, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(resp.Username, password, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConn2.Close(ctx)

			err = testConn2.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatalf("fook (%s:%s): %s", resp.Username, password, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tags := []string{"enterprise", "latest"}
	ctx := context.Background()

	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			neo, cleanup := getNeo4j(t, nil, tag)
			defer cleanup()

			password := "12345"
			req := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "arse", RoleName: "arse2",
				},
				Statements: dbplugin.Statements{
					Commands: []string{"CREATE OR REPLACE USER $username SET PLAINTEXT PASSWORD $password"},
				},
				RollbackStatements: dbplugin.Statements{},
				CredentialType:     dbplugin.CredentialTypePassword,
				Password:           password,
				Expiration:         time.Now().Add(time.Minute),
			}

			resp, err := neo.NewUser(ctx, req)

			if err != nil {
				t.Fatal(err)
			}

			testConn, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(resp.Username, password, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConn.Close(ctx)

			err = testConn.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatal(err)
			}

			delReq := dbplugin.DeleteUserRequest{
				Username: resp.Username,
				Statements: dbplugin.Statements{
					Commands: []string{
						`CALL dbms.listConnections() YIELD connectionId,  username WHERE username = $username WITH collect(connectionId) AS conns
					CALL dbms.killConnections(conns) YIELD connectionId, message RETURN *`,
						"DROP USER $username",
					},
				},
			}
			_, err = neo.DeleteUser(ctx, delReq)

			if err != nil {
				t.Fatal(err)
			}

			err = testConn.VerifyConnectivity(ctx)

			if err == nil {
				t.Fatalf("expected connection to fail, %s:%s", resp.Username, password)
			} else {
				_, isDbError := err.(*neo4j.ConnectivityError)
				if !isDbError {
					t.Fatalf("expected *neo4j.ConnectivityError, got %T", err)
				}
			}
			testConn2, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(resp.Username, password, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConn2.Close(ctx)

			err = testConn.VerifyConnectivity(ctx)

			if err == nil {
				t.Fatalf("expected old password to fail, %s:%s", resp.Username, password)
			} else {
				dbError, isDbError := err.(*db.Neo4jError)
				if !isDbError {
					t.Fatalf("expected *db.Neo4jError, got %T", err)
				}
				if !dbError.IsAuthenticationFailed() {
					t.Fatalf("expected error to be Neo.ClientError.Security.Unauthorized, but got %s", err)
				}
			}
		})
	}
}

func TestRotateRootCredentials(t *testing.T) {
	tags := []string{"enterprise", "latest"}
	ctx := context.Background()

	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			neo, cleanup := getNeo4j(t, nil, tag)
			defer cleanup()

			oldUsername := neo.Username
			oldPassword := neo.Password
			req := dbplugin.UpdateUserRequest{
				Username:       "neo4j",
				CredentialType: dbplugin.CredentialTypePassword,
				Password: &dbplugin.ChangePassword{
					NewPassword: "12345",
					Statements:  dbplugin.Statements{},
				},
			}
			_, err := neo.UpdateUser(ctx, req)
			if err != nil {
				t.Fatal(err)
			}

			testConnOld, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(oldUsername, oldPassword, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConnOld.Close(ctx)

			err = testConnOld.VerifyConnectivity(ctx)

			if err == nil {
				t.Fatalf("expected old password to fail, %s:%s", oldUsername, oldPassword)
			} else {
				dbError, isDbError := err.(*db.Neo4jError)
				if !isDbError {
					t.Fatalf("expected *db.Neo4jError, got %T", err)
				}
				if !dbError.IsAuthenticationFailed() {
					t.Fatalf("expected error to be Neo.ClientError.Security.Unauthorized, but got %s", err)
				}
			}

			newPassword := neo.Password
			testConnNew, err := neo4j.NewDriverWithContext(neo.ConnectionURL, neo4j.BasicAuth(oldUsername, newPassword, ""), cloneConfig(neo.ConnectionProducer))

			if err != nil {
				t.Fatal(err)
			}
			defer testConnNew.Close(ctx)

			err = testConnNew.VerifyConnectivity(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

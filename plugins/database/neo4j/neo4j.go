package neo4j

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"strings"
)

const (
	neo4jTypeName = "neo4j"

	defaultNeo4jRotateCredentialsCypher = `
ALTER USER $username SET PASSWORD $password CHANGE NOT REQUIRED
`
	defaultRevokeUserCypher = `
DROP USER $username;
`
	expirationFormat = "2006-01-02 15:04:05-0700"

	defaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 4) (.RoleName | truncate 40) (random 4) (unix_time) | truncate 63 }}`
)

var (
	_ dbplugin.Database = (*Neo4j)(nil)
)

type Neo4j struct {
	*ConnectionProducer
	usernameProducer template.StringTemplate
}

func New() (interface{}, error) {
	neo := newDB()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(neo, neo.SecretValues)
	return dbType, nil
}

func newDB() *Neo4j {
	connProducer := &ConnectionProducer{}
	neo := &Neo4j{
		ConnectionProducer: connProducer,
	}
	return neo
}

func (n *Neo4j) Type() (string, error) {
	return neo4jTypeName, nil
}

func (n *Neo4j) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	newConf, err := n.ConnectionProducer.Init(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}

	usernameTemplate, err := strutil.GetString(req.Config, "username_template")
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to retrieve username_template: %w", err)
	}
	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}

	up, err := template.NewTemplate(template.Template(usernameTemplate))
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to initialize username template: %w", err)
	}
	n.usernameProducer = up

	_, err = n.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	resp := dbplugin.InitializeResponse{
		Config: newConf,
	}
	return resp, nil
}

func (n *Neo4j) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	n.Lock()
	defer n.Unlock()
	username, err := n.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	expirationStr := req.Expiration.Format(expirationFormat)
	// Get the connection
	neo, err := n.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// Start a transaction
	session := neo.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: n.DatabaseName,
	})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	rollbackCyper := req.RollbackStatements.Commands
	if len(rollbackCyper) == 0 {
		rollbackCyper = []string{defaultRevokeUserCypher}
	}
	// Execute each query
	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]interface{}{
				"name":       username,
				"username":   username,
				"password":   req.Password,
				"expiration": expirationStr,
			}

			// In case someone uses {{templated}} variables instead of $params (support both)
			query = dbutil.QueryHelper(query, mapStrUnsafe(m))

			if _, err := tx.Run(ctx, stmt, m); err != nil {
				if n.shouldExecNonTransactional(err) {
					// Run each statement in a separate transaction (original transaction will be rolled back)
					err = n.execNonTransactional(req.Statements.Commands, username, req.Password, session, ctx)
					if err != nil {
						for _, rollbackStmt := range rollbackCyper {
							rollbackStmt = strings.TrimSpace(rollbackStmt)
							if len(rollbackStmt) == 0 {
								continue
							}
							_, errRollback := tx.Run(ctx, rollbackStmt, m)
							if errRollback != nil {
								if n.shouldExecNonTransactional(errRollback) {
									errRollback = n.execNonTransactional(req.Statements.Commands, username, req.Password, session, ctx)
									if errRollback != nil {
										return dbplugin.NewUserResponse{}, fmt.Errorf("error during user creation (%s) and rollback statements failed (%w)", err, errRollback)
									} else {
										return dbplugin.NewUserResponse{
											Username: username,
										}, nil
									}
								}
								return dbplugin.NewUserResponse{}, err
							}
						}
						return dbplugin.NewUserResponse{}, err
					} else {
						return dbplugin.NewUserResponse{
							Username: username,
						}, nil
					}
				} else {
					return dbplugin.NewUserResponse{}, err
				}

			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func (n *Neo4j) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Username == "" {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("missing username")
	}
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	merr := &multierror.Error{}
	if req.Password != nil {
		err := n.changeUserPassword(ctx, req.Username, req.Password)
		merr = multierror.Append(merr, err)

		if req.Username == n.Username {
			n.Password = req.Password.NewPassword
			n.RawConfig["password"] = req.Password.NewPassword
		}
	}
	if req.Expiration != nil {
		err := n.changeUserExpiration(ctx, req.Username, req.Expiration)
		merr = multierror.Append(merr, err)
	}
	return dbplugin.UpdateUserResponse{}, merr.ErrorOrNil()
}

func (n *Neo4j) changeUserExpiration(ctx context.Context, username string, changeExp *dbplugin.ChangeExpiration) error {
	n.Lock()
	defer n.Unlock()

	renewStmts := changeExp.Statements.Commands

	// Although user expiry/renewal isn't a concept in neo4j, a non-empty renewal statement
	// allows users to store information about expiry times in neo4j itself
	if len(renewStmts) == 0 {
		return nil
	}

	neo, err := n.getConnection(ctx)
	if err != nil {
		return err
	}

	// Start a transaction
	session := neo.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: n.DatabaseName,
	})
	defer session.Close(ctx)
	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	expirationStr := changeExp.NewExpiration.Format(expirationFormat)
	if err != nil {
		return err
	}

	for _, stmt := range renewStmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]interface{}{
				"name":       username,
				"username":   username,
				"expiration": expirationStr,
			}

			// In case someone uses {{templated}} variables instead of $params (support both)
			query = dbutil.QueryHelper(query, mapStrUnsafe(m))

			if _, err := tx.Run(ctx, stmt, m); err != nil {
				// These custom statements shouldn't fail due to ForbiddenDueToTransactionType
				// (since there's no admin operation relating to user expiry)
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (n *Neo4j) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	n.Lock()
	defer n.Unlock()

	if len(req.Statements.Commands) == 0 {
		// Dropping the user drops all permissions, but doesn't disconnect connections
		// If you wanted, you could instead configure the revocation statements to be:
		/*
			DROP USER $username;
			CALL dbms.listConnections() YIELD connectionId,  username WHERE username = $username WITH collect(connectionId) AS conns
			CALL dbms.killConnections(conns) YIELD connectionId, message RETURN *;
		*/
		req.Statements.Commands = []string{defaultRevokeUserCypher}
	}

	neo, err := n.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	// Start a transaction
	session := neo.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: n.DatabaseName,
	})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]interface{}{
				"name":     req.Username,
				"username": req.Username,
			}

			// In case someone uses {{templated}} variables instead of $params (support both)
			query = dbutil.QueryHelper(query, mapStrUnsafe(m))

			if _, err := tx.Run(ctx, stmt, m); err != nil {
				if n.shouldExecNonTransactional(err) {
					err = n.execNonTransactional(req.Statements.Commands, req.Username, "", session, ctx)
				}
				return dbplugin.DeleteUserResponse{}, err

			}
		}
	}

	return dbplugin.DeleteUserResponse{}, tx.Commit(ctx)
}

func (n *Neo4j) shouldExecNonTransactional(err error) bool {
	var dbError *db.Neo4jError
	if errors.As(err, &dbError) {
		if dbError.Code == "Neo.ClientError.Transaction.ForbiddenDueToTransactionType" {
			return true
		}
	}
	return false
}

func (n *Neo4j) execNonTransactional(stmts []string, username, password string, session neo4j.SessionWithContext, ctx context.Context) error {
	var errs []error
	for _, stmt := range stmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]interface{}{
				"name":     username,
				"username": username,
			}

			if password != "" {
				m["password"] = password
			}

			// In case someone uses {{templated}} variables instead of $params (support both)
			query = dbutil.QueryHelper(query, mapStrUnsafe(m))

			_, err := session.ExecuteWrite(ctx, func(tx2 neo4j.ManagedTransaction) (interface{}, error) {
				return tx2.Run(ctx, stmt, m)
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}
	return nil
}

func (n *Neo4j) changeUserPassword(ctx context.Context, username string, changePass *dbplugin.ChangePassword) (err error) {
	stmts := changePass.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{defaultNeo4jRotateCredentialsCypher}
	}

	password := changePass.NewPassword
	if password == "" {
		return fmt.Errorf("missing password")
	}

	n.Lock()
	defer n.Unlock()

	// Get the connection
	neo, err := n.getConnection(ctx)
	if err != nil {
		return fmt.Errorf("unable to get connection: %w", err)
	}

	session := neo.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: n.DatabaseName,
	})
	defer session.Close(ctx)

	userResult, err := session.Run(ctx, "SHOW USERS YIELD user, suspended WHERE user = $username", map[string]interface{}{"username": username})

	// Check if the role exists
	_, err = userResult.Single(ctx)

	if err != nil {
		// Zero or  > 1 result
		return fmt.Errorf("expected username %s to exist but got error: %w", username, err)
	}

	// Vault requires the database user already exist, and that the credentials
	// used to execute the rotation statements has sufficient privileges.

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Execute each query
	for _, stmt := range stmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}
			m := map[string]interface{}{
				"name":     username,
				"username": username,
				"password": password,
			}

			// In case someone uses {{templated}} variables instead of $params (support both)
			query = dbutil.QueryHelper(query, mapStrUnsafe(m))

			if _, err := tx.Run(ctx, query, m); err != nil {
				if n.shouldExecNonTransactional(err) {
					// Run each statement in a separate transaction (original transaction will be rolled back)
					err = n.execNonTransactional(stmts, username, password, session, ctx)
				}
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func mapStrUnsafe(input map[string]interface{}) map[string]string {
	output := map[string]string{}
	for k, v := range input {
		output[k] = v.(string)
	}
	return output
}

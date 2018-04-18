package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/dbtxn"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/lib/pq"
)

const SecretCredsType = "creds"

func secretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password",
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("usernameRaw is not a string")
	}
	// Get our connection
	db, err := b.DB(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Get the lease information
	lease, err := b.Lease(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	// Make sure we increase the VALID UNTIL endpoint for this user.
	ttl, _, err := framework.CalculateTTL(b.System(), req.Secret.Increment, lease.Lease, 0, lease.LeaseMax, 0, req.Secret.IssueTime)
	if err != nil {
		return nil, err
	}
	if ttl > 0 {
		expireTime := time.Now().Add(ttl)
		// Adding a small buffer since the TTL will be calculated again afeter this call
		// to ensure the database credential does not expire before the lease
		expireTime = expireTime.Add(5 * time.Second)
		expiration := expireTime.Format("2006-01-02 15:04:05-0700")

		query := fmt.Sprintf(
			"ALTER ROLE %s VALID UNTIL '%s';",
			pq.QuoteIdentifier(username),
			expiration)
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.Lease
	resp.Secret.MaxTTL = lease.LeaseMax
	return resp, nil
}

func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("usernameRaw is not a string")
	}
	var revocationSQL string
	var resp *logical.Response

	roleNameRaw, ok := req.Secret.InternalData["role"]
	if ok {
		role, err := b.Role(ctx, req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			if resp == nil {
				resp = &logical.Response{}
			}
			resp.AddWarning(fmt.Sprintf("Role %q cannot be found. Using default revocation SQL.", roleNameRaw.(string)))
		} else {
			revocationSQL = role.RevocationSQL
		}
	}

	// Get our connection
	db, err := b.DB(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	switch revocationSQL {

	// This is the default revocation logic. If revocation SQL is provided it
	// is simply executed as-is.
	case "":
		// Check if the role exists
		var exists bool
		err = db.QueryRow("SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if exists == false {
			return resp, nil
		}

		// Query for permissions; we need to revoke permissions before we can drop
		// the role
		// This isn't done in a transaction because even if we fail along the way,
		// we want to remove as much access as possible
		stmt, err := db.Prepare("SELECT DISTINCT table_schema FROM information_schema.role_column_grants WHERE grantee=$1;")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		rows, err := stmt.Query(username)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		const initialNumRevocations = 16
		revocationStmts := make([]string, 0, initialNumRevocations)
		for rows.Next() {
			var schema string
			err = rows.Scan(&schema)
			if err != nil {
				// keep going; remove as many permissions as possible right now
				continue
			}
			revocationStmts = append(revocationStmts, fmt.Sprintf(
				`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA %s FROM %s;`,
				pq.QuoteIdentifier(schema),
				pq.QuoteIdentifier(username)))

			revocationStmts = append(revocationStmts, fmt.Sprintf(
				`REVOKE USAGE ON SCHEMA %s FROM %s;`,
				pq.QuoteIdentifier(schema),
				pq.QuoteIdentifier(username)))
		}

		// for good measure, revoke all privileges and usage on schema public
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM %s;`,
			pq.QuoteIdentifier(username)))

		revocationStmts = append(revocationStmts, fmt.Sprintf(
			"REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM %s;",
			pq.QuoteIdentifier(username)))

		revocationStmts = append(revocationStmts, fmt.Sprintf(
			"REVOKE USAGE ON SCHEMA public FROM %s;",
			pq.QuoteIdentifier(username)))

		// get the current database name so we can issue a REVOKE CONNECT for
		// this username
		var dbname sql.NullString
		if err := db.QueryRow("SELECT current_database();").Scan(&dbname); err != nil {
			return nil, err
		}

		if dbname.Valid {
			revocationStmts = append(revocationStmts, fmt.Sprintf(
				`REVOKE CONNECT ON DATABASE %s FROM %s;`,
				pq.QuoteIdentifier(dbname.String),
				pq.QuoteIdentifier(username)))
		}

		// again, here, we do not stop on error, as we want to remove as
		// many permissions as possible right now
		var lastStmtError error
		for _, query := range revocationStmts {
			if err := dbtxn.ExecuteDBQuery(ctx, db, nil, query); err != nil {
				lastStmtError = err
			}
		}

		// can't drop if not all privileges are revoked
		if rows.Err() != nil {
			return nil, errwrap.Wrapf("could not generate revocation statements for all rows: {{err}}", rows.Err())
		}
		if lastStmtError != nil {
			return nil, errwrap.Wrapf("could not perform all revocation statements: {{err}}", lastStmtError)
		}

		// Drop this user
		stmt, err = db.Prepare(fmt.Sprintf(
			`DROP ROLE IF EXISTS %s;`, pq.QuoteIdentifier(username)))
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return nil, err
		}

	// We have revocation SQL, execute directly, within a transaction
	default:
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer func() {
			tx.Rollback()
		}()

		for _, query := range strutil.ParseArbitraryStringSlice(revocationSQL, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name": username,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return nil, err
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

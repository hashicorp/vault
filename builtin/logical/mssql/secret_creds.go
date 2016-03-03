package mssql

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

func (b *backend) secretCredsRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, b.System())
	return f(req, d)
}

func (b *backend) secretCredsRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// First disable server login
	stmt, err := db.Prepare(fmt.Sprintf("ALTER LOGIN [%s] DISABLE;", username))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	// Query for sessions for the login so that we can kill any outstanding
	// sessions.  There cannot be any active sessions before we drop the logins
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err = db.Prepare(fmt.Sprintf(
		"SELECT session_id FROM sys.dm_exec_sessions WHERE login_name = '%s';", username))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revokeStmts []string
	for rows.Next() {
		var sessionID int
		err = rows.Scan(&sessionID)
		if err != nil {
			continue
		}
		revokeStmts = append(revokeStmts, fmt.Sprintf("KILL %d;", sessionID))
	}

	// Query for database users using undocumented stored procedure for now since
	// it is the easiest way to get this information;
	// we need to drop the database users before we can drop the login and the role
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err = db.Prepare(fmt.Sprintf("EXEC sp_msloginmappings '%s';", username))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var loginName, dbName, qUsername, aliasName string
		err = rows.Scan(&loginName, &dbName, &qUsername, &aliasName)
		if err != nil {
			// keep going; remove as many permissions as possible right now
			continue
		}
		revokeStmts = append(revokeStmts, fmt.Sprintf(
			"USE [%s]; DROP USER IF EXISTS [%s];", dbName, username))
	}

	// we do not stop on error, as we want to remove as
	// many permissions as possible right now
	var lastStmtError error
	for _, query := range revokeStmts {
		stmt, err := db.Prepare(query)
		if err != nil {
			lastStmtError = err
			continue
		}
		_, err = stmt.Exec()
		if err != nil {
			lastStmtError = err
		}
	}

	// can't drop if not all database users are dropped
	if rows.Err() != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"could not generate sql statements for all rows: %v", rows.Err())), nil
	}
	if lastStmtError != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"could not perform all sql statements: %v", lastStmtError)), nil
	}

	// Drop this login
	stmt, err = db.Prepare(fmt.Sprintf(dropLoginSQL, username, username))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	return nil, nil
}

const dropLoginSQL = `
IF EXISTS
  (SELECT name
   FROM master.sys.server_principals
   WHERE name = '%s')
BEGIN
    DROP LOGIN [%s]
END
`

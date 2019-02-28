package database

import (
        "context"
        "fmt"
        "time"

        "github.com/hashicorp/vault/helper/consts"
        "github.com/hashicorp/vault/logical"
        "github.com/hashicorp/vault/logical/framework"
        "github.com/mitchellh/mapstructure"
)

// walSetCredentials is used to store information in a WAL that can re-try a
// credential setting or rotation in the event of partial failure.
type walSetCredentials struct {
        Username          string    `json:"username"`
        Password          string    `json:"password"`
        RoleName          string    `json:"role_name"`
        Statements        []string  `json:"statements"`
        LastVaultRotation time.Time `json:"last_vault_rotation"`
}

func (b *databaseBackend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
        walRollbackMap := map[string]framework.WALRollbackFunc{
                "setcreds": b.pathCredentialRollback,
        }

        if !b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
                return nil
        }

        f, ok := walRollbackMap[kind]
        if !ok {
                return fmt.Errorf("unknown type to rollback")
        }

        return f(ctx, req, kind, data)
}

func (b *databaseBackend) pathCredentialRollback(ctx context.Context, req *logical.Request, _kind string, data interface{}) error {
        var entry walSetCredentials
        if err := mapstructure.Decode(data, &entry); err != nil {
                return err
        }
        // check the LastVaultRotation times. If the role has had a password change
        // since the wal's LastVaultRotation, we can assume things are fine here

        // attempt to rollback the password to a known value

        return nil
}

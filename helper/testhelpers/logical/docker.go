package testing

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/acctest"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

func dockerTest(tt TestT, c TestCase) {
	// Create an in-memory Vault core
	logger := logging.NewVaultLogger(log.Trace)
	// TODO: function or method to make this safe
	if acctest.TestHelper == nil {
		tt.Fatal("expected test helper")
	}
	client := acctest.TestHelper.Client

	// TODO custom path to avoid conflict
	err := client.Sys().Mount("plug-transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		tt.Fatal(err)
	}
	defer func() {
		err := client.Sys().Unmount("plug-transit")
		if err != nil {
			tt.Fatal(err)
		}
	}()

	tokenInfo, err := client.Auth().Token().LookupSelf()
	if err != nil {
		tt.Fatal("error looking up token: ", err)
		return
	}
	var tokenPolicies []string
	if tokenPoliciesRaw, ok := tokenInfo.Data["policies"]; ok {
		if tokenPoliciesSliceRaw, ok := tokenPoliciesRaw.([]interface{}); ok {
			for _, p := range tokenPoliciesSliceRaw {
				tokenPolicies = append(tokenPolicies, p.(string))
			}
		}
	}

	// go through steps
	// Make requests
	var revoke []*logical.Request
	for i, s := range c.Steps {
		if logger.IsWarn() {
			logger.Warn("Executing test step", "step_number", i+1)
		}

		// Create the request
		// TODO translate into client.Logical.Write/et. al
		// req := &logical.Request{
		// 	Operation: s.Operation,
		// 	Path:      s.Path,
		// 	Data:      s.Data,
		// }

		// TODO hard coded path here:
		path := fmt.Sprintf("plug-transit/%s", s.Path)
		var err error
		var resp *api.Secret
		// TODO should check expect none here?
		var lr *logical.Response
		switch s.Operation {
		case logical.CreateOperation, logical.UpdateOperation:
			// resp, err = client.Logical().Write(s.Path, s.Data)
			resp, err = client.Logical().Write(path, s.Data)
		case logical.ReadOperation:
			resp, err = client.Logical().Read(path)
		case logical.ListOperation:
			resp, err = client.Logical().List(path)
			// TODO why though
			lr = &logical.Response{}
		case logical.DeleteOperation:
			resp, err = client.Logical().Delete(path)
		default:
			panic("bad operation")
		}

		// TODO: verify this check
		// error at this point is a problem with the request?
		if err != nil && !s.ErrorOk {
			tt.Fatal(err)
		}
		// TODO: unauth, preflight
		// if !s.Unauthenticated {
		// 	req.ClientToken = client.Token()
		// 	req.SetTokenEntry(&logical.TokenEntry{
		// 		ID:          req.ClientToken,
		// 		NamespaceID: namespace.RootNamespaceID,
		// 		Policies:    tokenPolicies,
		// 		DisplayName: tokenInfo.Data["display_name"].(string),
		// 	})
		// }
		// req.Connection = &logical.Connection{RemoteAddr: s.RemoteAddr}
		// if s.ConnState != nil {
		// 	req.Connection.ConnState = s.ConnState
		// }

		// if s.PreFlight != nil {
		// 	ct := req.ClientToken
		// 	req.ClientToken = ""
		// 	if err := s.PreFlight(req); err != nil {
		// 		tt.Error(fmt.Sprintf("Failed preflight for step %d: %s", i+1, err))
		// 		break
		// 	}
		// 	req.ClientToken = ct
		// }

		// Make sure to prefix the path with where we mounted the thing
		// TODO setup needs to know prefix/mount
		// prefix := "transit"
		// req.Path = fmt.Sprintf("%s/%s", prefix, req.Path)

		// if isAuthBackend {
		// 	// Prepend the path with "auth"
		// 	req.Path = "auth/" + req.Path
		// }

		// Make the request
		// resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		// if resp != nil && resp.Secret != nil {
		if resp != nil {
			// Revoke this secret later
			revoke = append(revoke, &logical.Request{
				Operation: logical.UpdateOperation,
				// Path:      "sys/revoke/" + resp.Secret.LeaseID,
				Path: "sys/revoke/" + resp.LeaseID,
			})
		}

		// Test step returned an error.
		if err != nil {
			// But if an error is expected, do not fail the test step,
			// regardless of whether the error is a 'logical.ErrorResponse'
			// or not. Set the err to nil. If the error is a logical.ErrorResponse,
			// it will be handled later.
			if s.ErrorOk {
				err = nil
			} else {
				// If the error is not expected, fail right away.
				tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
				break
			}
		}

		// If the error is a 'logical.ErrorResponse' and if error was not expected,
		// set the error so that this can be caught below.
		// TODO logical error here
		// if resp.IsError() && !s.ErrorOk {
		// 	err = fmt.Errorf("erroneous response:\n\n%#v", resp)
		// }

		// Either the 'err' was nil or if an error was expected, it was set to nil.
		// Call the 'Check' function if there is one.
		//
		// TODO: This works perfectly for now, but it would be better if 'Check'
		// function takes in both the response object and the error, and decide on
		// the action on its own.
		if err == nil && s.Check != nil {
			// Call the test method
			// TODO faking logical response right here

			if resp != nil && lr == nil {
				lr = &logical.Response{}
			}

			if resp != nil {
				lr.Secret = &logical.Secret{
					LeaseID: resp.LeaseID,
				}
				lr.Data = resp.Data
			}

			err = s.Check(lr)
		}

		// TODO check err ok here
		if err != nil && !s.ErrorOk {
			tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
			break
		}
	}
	// TODO revoke secrets
	// Revoke any secrets we might have.
	// var failedRevokes []*logical.Secret
	// for _, req := range revoke {
	// 	if logger.IsWarn() {
	// 		logger.Warn("Revoking secret", "secret", fmt.Sprintf("%#v", req))
	// 	}
	// 	req.ClientToken = client.Token()
	// 	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	// 	if err == nil && resp.IsError() {
	// 		err = fmt.Errorf("erroneous response:\n\n%#v", resp)
	// 	}
	// 	if err != nil {
	// 		failedRevokes = append(failedRevokes, req.Secret)
	// 		tt.Error(fmt.Sprintf("Revoke error: %s", err))
	// 	}
	// }

	// Perform any rollbacks. This should no-op if there aren't any.
	// We set the "immediate" flag here that any backend can pick up on
	// to do all rollbacks immediately even if the WAL entries are new.
	// logger.Warn("Requesting RollbackOperation")
	// rollbackPath := prefix + "/"
	// if c.CredentialFactory != nil || c.CredentialBackend != nil {
	// 	rollbackPath = "auth/" + rollbackPath
	// }
	// req := logical.RollbackRequest(rollbackPath)
	// req.Data["immediate"] = true
	// req.ClientToken = client.Token()
	// resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	// if err == nil && resp.IsError() {
	// 	err = fmt.Errorf("erroneous response:\n\n%#v", resp)
	// }
	// if err != nil {
	// 	if !errwrap.Contains(err, logical.ErrUnsupportedOperation.Error()) {
	// 		tt.Error(fmt.Sprintf("[ERR] Rollback error: %s", err))
	// 	}
	// }

	// // If we have any failed revokes, log it.
	// if len(failedRevokes) > 0 {
	// 	for _, s := range failedRevokes {
	// 		tt.Error(fmt.Sprintf(
	// 			"WARNING: Revoking the following secret failed. It may\n"+
	// 				"still exist. Please verify:\n\n%#v",
	// 			s))
	// 	}
	// }
}

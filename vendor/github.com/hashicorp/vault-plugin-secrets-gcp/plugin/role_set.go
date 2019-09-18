package gcpsecrets

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"regexp"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iam/v1"
)

const (
	serviceAccountMaxLen          = 30
	serviceAccountDisplayNameTmpl = "Service account for Vault secrets backend role set %s"
)

type RoleSet struct {
	Name       string
	SecretType string

	RawBindings string
	Bindings    ResourceBindings

	AccountId *gcputil.ServiceAccountId
	TokenGen  *TokenGenerator
}

func (rs *RoleSet) validate() error {
	var err *multierror.Error
	if rs.Name == "" {
		err = multierror.Append(err, errors.New("role set name is empty"))
	}

	if rs.SecretType == "" {
		err = multierror.Append(err, errors.New("role set secret type is empty"))
	}

	if rs.AccountId == nil {
		err = multierror.Append(err, fmt.Errorf("role set should have account associated"))
	}

	if len(rs.Bindings) == 0 {
		err = multierror.Append(err, fmt.Errorf("role set bindings cannot be empty"))
	}

	if len(rs.RawBindings) == 0 {
		err = multierror.Append(err, fmt.Errorf("role set raw bindings cannot be empty string"))
	}

	switch rs.SecretType {
	case SecretTypeAccessToken:
		if rs.TokenGen == nil {
			err = multierror.Append(err, fmt.Errorf("access token role set should have initialized token generator"))
		} else if len(rs.TokenGen.Scopes) == 0 {
			err = multierror.Append(err, fmt.Errorf("access token role set should have defined scopes"))
		}
	case SecretTypeKey:
		break
	default:
		err = multierror.Append(err, fmt.Errorf("unknown secret type: %s", rs.SecretType))
	}
	return err.ErrorOrNil()
}

func (rs *RoleSet) save(ctx context.Context, s logical.Storage) error {
	if err := rs.validate(); err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", rolesetStoragePrefix, rs.Name), rs)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func (rs *RoleSet) bindingHash() string {
	return getStringHash(rs.RawBindings)
}

func (rs *RoleSet) getServiceAccount(iamAdmin *iam.Service) (*iam.ServiceAccount, error) {
	if rs.AccountId == nil {
		return nil, fmt.Errorf("role set '%s' is invalid, has no associated service account", rs.Name)
	}

	account, err := iamAdmin.Projects.ServiceAccounts.Get(rs.AccountId.ResourceName()).Do()
	if err != nil {
		return nil, fmt.Errorf("could not find service account: %v. If account was deleted, role set must be updated (write to roleset/%s/rotate) before generating new secrets", err, rs.Name)
	} else if account == nil {
		return nil, fmt.Errorf("roleset service account was removed - role set must be updated (path roleset/%s/rotate) before generating new secrets", rs.Name)
	}

	return account, nil
}

type ResourceBindings map[string]util.StringSet

func (rb ResourceBindings) asOutput() map[string][]string {
	out := make(map[string][]string)
	for k, v := range rb {
		out[k] = v.ToSlice()
	}
	return out
}

type TokenGenerator struct {
	KeyName    string
	B64KeyJSON string

	Scopes []string
}

func (b *backend) saveRoleSetWithNewAccount(ctx context.Context, s logical.Storage, rs *RoleSet, project string, newBinds ResourceBindings, scopes []string) (warning []string, err error) {
	b.rolesetLock.Lock()
	defer b.rolesetLock.Unlock()

	httpC, err := b.HTTPClient(s)
	if err != nil {
		return nil, err
	}

	iamAdmin, err := b.IAMClient(s)
	if err != nil {
		return nil, err
	}

	iamHandle := iamutil.GetIamHandle(httpC, useragent.String())

	oldAccount := rs.AccountId
	oldBindings := rs.Bindings
	oldTokenKey := rs.TokenGen

	oldWals, err := rs.addWALsForCurrentAccount(ctx, s)
	if err != nil {
		tryDeleteWALs(ctx, s, oldWals...)
		return nil, errwrap.Wrapf("failed to create WAL for cleaning up old account: {{err}}", err)
	}

	newWals := make([]string, 0, len(newBinds)+2)
	walId, err := rs.newServiceAccount(ctx, s, iamAdmin, project)
	if err != nil {
		tryDeleteWALs(ctx, s, oldWals...)
		return nil, err
	}
	newWals = append(newWals, walId)

	binds := rs.Bindings
	if newBinds != nil {
		binds = newBinds
		rs.Bindings = newBinds
	}
	walIds, err := rs.updateIamPolicies(ctx, s, b.iamResources, iamHandle, binds)
	if err != nil {
		tryDeleteWALs(ctx, s, oldWals...)
		return nil, err
	}
	newWals = append(newWals, walIds...)

	if rs.SecretType == SecretTypeAccessToken {
		walId, err := rs.newKeyForTokenGen(ctx, s, iamAdmin, scopes)
		if err != nil {
			tryDeleteWALs(ctx, s, oldWals...)
			return nil, err
		}
		newWals = append(newWals, walId)
	}

	if err := rs.save(ctx, s); err != nil {
		tryDeleteWALs(ctx, s, oldWals...)
		return nil, err
	}

	// Delete WALs for cleaning up new resources now that they have been saved.
	tryDeleteWALs(ctx, s, newWals...)

	// Try deleting old resources (WALs exist so we can ignore failures)
	if oldAccount == nil || oldAccount.EmailOrId == "" {
		// nothing to clean up
		return nil, nil
	}

	// Return any errors as warnings so user knows immediate cleanup failed
	warnings := make([]string, 0)
	if errs := b.removeBindings(ctx, iamHandle, oldAccount.EmailOrId, oldBindings); errs != nil {
		warnings = make([]string, len(errs.Errors), len(errs.Errors)+2)
		for idx, err := range errs.Errors {
			warnings[idx] = fmt.Sprintf("unable to immediately delete old binding (WAL cleanup entry has been added): %v", err)
		}
	}
	if err := b.deleteServiceAccount(ctx, iamAdmin, oldAccount); err != nil {
		warnings = append(warnings, fmt.Sprintf("unable to immediately delete old account (WAL cleanup entry has been added): %v", err))
	}
	if err := b.deleteTokenGenKey(ctx, iamAdmin, oldTokenKey); err != nil {
		warnings = append(warnings, fmt.Sprintf("unable to immediately delete old key (WAL cleanup entry has been added): %v", err))
	}
	return warnings, nil
}

func (b *backend) saveRoleSetWithNewTokenKey(ctx context.Context, s logical.Storage, rs *RoleSet, scopes []string) (warning string, err error) {
	b.rolesetLock.Lock()
	defer b.rolesetLock.Unlock()

	if rs.SecretType != SecretTypeAccessToken {
		return "", fmt.Errorf("a key is not saved or used for non-access-token role set '%s'", rs.Name)
	}

	iamAdmin, err := b.IAMClient(s)
	if err != nil {
		return "", err
	}

	oldKeyWalId := ""
	if rs.TokenGen != nil {
		if oldKeyWalId, err = framework.PutWAL(ctx, s, walTypeAccountKey, &walAccountKey{
			RoleSet:            rs.Name,
			KeyName:            rs.TokenGen.KeyName,
			ServiceAccountName: rs.AccountId.ResourceName(),
		}); err != nil {
			return "", errwrap.Wrapf("unable to create WAL for deleting old key: {{err}}", err)
		}
	}
	oldKeyGen := rs.TokenGen

	newKeyWalId, err := rs.newKeyForTokenGen(ctx, s, iamAdmin, scopes)
	if err != nil {
		tryDeleteWALs(ctx, s, oldKeyWalId)
		return "", err
	}

	if err := rs.save(ctx, s); err != nil {
		tryDeleteWALs(ctx, s, oldKeyWalId)
		return "", err
	}

	// Delete WALs for cleaning up new key now that it's been saved.
	tryDeleteWALs(ctx, s, newKeyWalId)
	if err := b.deleteTokenGenKey(ctx, iamAdmin, oldKeyGen); err != nil {
		return errwrap.Wrapf("unable to delete old key (delayed cleaned up WAL entry added): {{err}}", err).Error(), nil
	}

	return "", nil
}

func (rs *RoleSet) addWALsForCurrentAccount(ctx context.Context, s logical.Storage) ([]string, error) {
	if rs.AccountId == nil {
		return nil, nil
	}
	wals := make([]string, 0, len(rs.Bindings)+2)
	walId, err := framework.PutWAL(ctx, s, walTypeAccount, &walAccount{
		RoleSet: rs.Name,
		Id: gcputil.ServiceAccountId{
			Project:   rs.AccountId.Project,
			EmailOrId: rs.AccountId.EmailOrId,
		},
	})
	if err != nil {
		return nil, err
	}
	wals = append(wals, walId)
	for resource, roles := range rs.Bindings {
		var walId string
		walId, err = framework.PutWAL(ctx, s, walTypeIamPolicy, &walIamPolicy{
			RoleSet: rs.Name,
			AccountId: gcputil.ServiceAccountId{
				Project:   rs.AccountId.Project,
				EmailOrId: rs.AccountId.EmailOrId,
			},
			Resource: resource,
			Roles:    roles.ToSlice(),
		})
		if err != nil {
			return nil, err
		}
		wals = append(wals, walId)
	}

	if rs.SecretType == SecretTypeAccessToken && rs.TokenGen != nil {
		walId, err := framework.PutWAL(ctx, s, walTypeAccountKey, &walAccountKey{
			RoleSet:            rs.Name,
			KeyName:            rs.TokenGen.KeyName,
			ServiceAccountName: rs.AccountId.ResourceName(),
		})
		if err != nil {
			return nil, err
		}
		wals = append(wals, walId)
	}
	return wals, nil
}

func (rs *RoleSet) newServiceAccount(ctx context.Context, s logical.Storage, iamAdmin *iam.Service, project string) (string, error) {
	saEmailPrefix := roleSetServiceAccountName(rs.Name)
	projectName := fmt.Sprintf("projects/%s", project)
	displayName := fmt.Sprintf(serviceAccountDisplayNameTmpl, rs.Name)

	walId, err := framework.PutWAL(ctx, s, walTypeAccount, &walAccount{
		RoleSet: rs.Name,
		Id: gcputil.ServiceAccountId{
			Project:   project,
			EmailOrId: fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saEmailPrefix, project),
		},
	})
	if err != nil {
		return "", errwrap.Wrapf("unable to create WAL entry for generating new service account: {{err}}", err)
	}

	sa, err := iamAdmin.Projects.ServiceAccounts.Create(
		projectName, &iam.CreateServiceAccountRequest{
			AccountId:      saEmailPrefix,
			ServiceAccount: &iam.ServiceAccount{DisplayName: displayName},
		}).Do()
	if err != nil {
		return walId, errwrap.Wrapf(fmt.Sprintf("unable to create new service account under project '%s': {{err}}", projectName), err)
	}
	rs.AccountId = &gcputil.ServiceAccountId{
		Project:   project,
		EmailOrId: sa.Email,
	}
	return walId, nil
}

func (rs *RoleSet) newKeyForTokenGen(ctx context.Context, s logical.Storage, iamAdmin *iam.Service, scopes []string) (string, error) {
	walId, err := framework.PutWAL(ctx, s, walTypeAccountKey, &walAccountKey{
		RoleSet:            rs.Name,
		KeyName:            "",
		ServiceAccountName: rs.AccountId.ResourceName(),
	})
	if err != nil {
		return "", err
	}

	key, err := iamAdmin.Projects.ServiceAccounts.Keys.Create(rs.AccountId.ResourceName(),
		&iam.CreateServiceAccountKeyRequest{
			PrivateKeyType: privateKeyTypeJson,
		}).Do()
	if err != nil {
		framework.DeleteWAL(ctx, s, walId)
		return "", err
	}
	rs.TokenGen = &TokenGenerator{
		KeyName:    key.Name,
		B64KeyJSON: key.PrivateKeyData,
		Scopes:     scopes,
	}
	return walId, nil
}

func (rs *RoleSet) updateIamPolicies(ctx context.Context, s logical.Storage, enabledIamResources iamutil.IamResourceParser, iamHandle *iamutil.IamHandle, rb ResourceBindings) ([]string, error) {
	wals := make([]string, 0, len(rb))
	for rName, roles := range rb {
		walId, err := framework.PutWAL(ctx, s, walTypeIamPolicy, &walIamPolicy{
			RoleSet: rs.Name,
			AccountId: gcputil.ServiceAccountId{
				Project:   rs.AccountId.Project,
				EmailOrId: rs.AccountId.EmailOrId,
			},
			Resource: rName,
			Roles:    roles.ToSlice(),
		})
		if err != nil {
			return wals, err
		}

		resource, err := enabledIamResources.Parse(rName)
		if err != nil {
			return wals, err
		}

		p, err := iamHandle.GetIamPolicy(ctx, resource)
		if err != nil {
			return wals, err
		}

		changed, newP := p.AddBindings(&iamutil.PolicyDelta{
			Roles: roles,
			Email: rs.AccountId.EmailOrId,
		})
		if !changed || newP == nil {
			continue
		}

		if _, err := iamHandle.SetIamPolicy(ctx, resource, newP); err != nil {
			return wals, err
		}
		wals = append(wals, walId)
	}
	return wals, nil
}

func roleSetServiceAccountName(rsName string) (name string) {
	// Sanitize role name
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	rsName = reg.ReplaceAllString(rsName, "-")

	intSuffix := fmt.Sprintf("%d", time.Now().Unix())
	fullName := fmt.Sprintf("vault%s-%s", rsName, intSuffix)
	name = fullName
	if len(fullName) > serviceAccountMaxLen {
		toTrunc := len(fullName) - serviceAccountMaxLen
		name = fmt.Sprintf("vault%s-%s", rsName[:len(rsName)-toTrunc], intSuffix)
	}
	return name
}

func getStringHash(bindingsRaw string) string {
	ssum := sha256.Sum256([]byte(bindingsRaw)[:])
	return base64.StdEncoding.EncodeToString(ssum[:])
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
	"github.com/hashicorp/vault/sdk/framework"
)

type inputParams struct {
	name       string
	secretType string

	hasBindings bool
	rawBindings string
	bindings    ResourceBindings

	project             string
	serviceAccountEmail string

	scopes []string
}

func (input *inputParams) parseOkInputSecretType(d *framework.FieldData) (warnings []string, err error) {
	secretType := d.Get("secret_type").(string)
	if secretType == "" && input.secretType == "" {
		return nil, fmt.Errorf("secret_type is required")
	}
	if input.secretType != "" && secretType != "" && input.secretType != secretType {
		return nil, fmt.Errorf("cannot update secret_type")
	}

	switch secretType {
	case SecretTypeKey, SecretTypeAccessToken:
		input.secretType = secretType
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid secret_type %q", secretType)
	}
}

// parseOkInputServiceAccountEmail checks that when creating a static acocunt, a service account
// email is provided. A service account email can be provide while updating the static account
// but it must be the same as the one in the static account and cannot be updated.
func (input *inputParams) parseOkInputServiceAccountEmail(d *framework.FieldData) (warnings []string, err error) {
	email := d.Get("service_account_email").(string)
	if email == "" && input.serviceAccountEmail == "" {
		return nil, fmt.Errorf("email is required")
	}
	if input.serviceAccountEmail != "" && email != "" && input.serviceAccountEmail != email {
		return nil, fmt.Errorf("cannot update email")
	}

	input.serviceAccountEmail = email
	return nil, nil
}

func (input *inputParams) parseOkInputTokenScopes(d *framework.FieldData) (warnings []string, err error) {
	// Parse secretType if not yet parsed
	if input.secretType == "" {
		warnings, err = input.parseOkInputSecretType(d)
		if err != nil {
			return nil, err
		}
	}

	v, ok := d.GetOk("token_scopes")
	if ok {
		scopes, castOk := v.([]string)
		if !castOk {
			return nil, fmt.Errorf("scopes unexpected type %T, expected []string", v)
		}
		input.scopes = scopes
	}

	if input.secretType == SecretTypeAccessToken && len(input.scopes) == 0 {
		return nil, fmt.Errorf("non-empty token_scopes must be provided for generating access_token secrets")
	}

	if input.secretType != SecretTypeAccessToken && ok && len(input.scopes) > 0 {
		warnings = append(warnings, "ignoring non-empty token_scopes, secret type not access_token")
	}
	return
}

func (input *inputParams) parseOkInputBindings(d *framework.FieldData) (warnings []string, err error) {
	bRaw, ok := d.GetOk("bindings")
	if !ok {
		input.hasBindings = false
		return nil, nil
	}

	rawBindings, castok := bRaw.(string)
	if !castok {
		return nil, fmt.Errorf("bindings are not a string")
	}

	bindings, err := util.ParseBindings(bRaw.(string))
	if err != nil {
		return nil, errwrap.Wrapf("unable to parse bindings: {{err}}", err)
	}

	input.hasBindings = true
	input.rawBindings = rawBindings
	input.bindings = bindings
	return nil, nil
}

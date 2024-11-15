// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/authmetadata"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

// gcpConfig contains all config required for the GCP backend.
type gcpConfig struct {
	Credentials     *gcputil.GcpCredentials `json:"credentials"`
	IAMAliasType    string                  `json:"iam_alias"`
	IAMAuthMetadata *authmetadata.Handler   `json:"iam_auth_metadata_handler"`
	GCEAliasType    string                  `json:"gce_alias"`
	GCEAuthMetadata *authmetadata.Handler   `json:"gce_auth_metadata_handler"`

	// APICustomEndpoint overrides the service endpoint for www.googleapis.com
	APICustomEndpoint string `json:"api_custom_endpoint"`
	// IAMCustomEndpoint overrides the service endpoint for api.googleapis.com
	IAMCustomEndpoint string `json:"iam_custom_endpoint"`
	// CRMCustomEndpoint overrides the service endpoint for cloudresourcemananger.googleapis.com
	CRMCustomEndpoint string `json:"crm_custom_endpoint"`
	// ComputeCustomEndpoint overrides the service endpoint for compute.googleapis.com
	ComputeCustomEndpoint string `json:"compute_custom_endpoint"`

	pluginidentityutil.PluginIdentityTokenParams
	ServiceAccountEmail string `json:"service_account_email"`
}

// standardizedCreds wraps gcputil.GcpCredentials with a type to allow
// parsing through Google libraries, since the google libraries struct is not
// exposed.
type standardizedCreds struct {
	*gcputil.GcpCredentials
	CredType string `json:"type"`
}

const serviceAccountCredsType = "service_account"

// formatAsCredentialJSON converts and marshals the config credentials
// into a parsable format by Google libraries.
func (c *gcpConfig) formatAndMarshalCredentials() ([]byte, error) {
	if c == nil || c.Credentials == nil {
		return []byte{}, nil
	}

	return json.Marshal(standardizedCreds{
		GcpCredentials: c.Credentials,
		CredType:       serviceAccountCredsType,
	})
}

// Update sets gcpConfig values parsed from the FieldData.
func (c *gcpConfig) Update(d *framework.FieldData) error {
	if d == nil {
		return nil
	}

	if v, ok := d.GetOk("credentials"); ok {
		credentials := v.(string)

		// If the given credentials are empty, reset them so that application default
		// credentials are used. Otherwise, parse and validate the given credentials.
		if credentials == "" {
			c.Credentials = nil
		} else {
			creds, err := gcputil.Credentials(credentials)
			if err != nil {
				return fmt.Errorf("failed to read credentials: %w", err)
			}

			if len(creds.PrivateKeyId) == 0 {
				return errors.New("missing private key in credentials")
			}

			c.Credentials = creds
		}
	}

	rawIamAlias, exists := d.GetOk("iam_alias")
	if exists {
		c.IAMAliasType = rawIamAlias.(string)
	}

	if err := c.IAMAuthMetadata.ParseAuthMetadata(d); err != nil {
		return fmt.Errorf("failed to parse iam metadata: %w", err)
	}

	rawGceAlias, exists := d.GetOk("gce_alias")
	if exists {
		c.GCEAliasType = rawGceAlias.(string)
	}

	if err := c.GCEAuthMetadata.ParseAuthMetadata(d); err != nil {
		return fmt.Errorf("failed to parse gce metadata: %w", err)
	}

	rawEndpoint, exists := d.GetOk("custom_endpoint")
	if exists {
		for k, v := range rawEndpoint.(map[string]string) {
			switch k {
			case "api":
				c.APICustomEndpoint = v
			case "iam":
				c.IAMCustomEndpoint = v
			case "crm":
				c.CRMCustomEndpoint = v
			case "compute":
				c.ComputeCustomEndpoint = v
			default:
				return fmt.Errorf("invalid custom endpoint type %q. Available types are: 'api', 'iam', 'crm', 'compute'", k)
			}
		}
	}

	// set plugin identity token fields
	if err := c.ParsePluginIdentityTokenFields(d); err != nil {
		return err
	}

	// set Service Account email
	saEmail, ok := d.GetOk("service_account_email")
	if ok {
		c.ServiceAccountEmail = saEmail.(string)
	}

	if c.IdentityTokenAudience != "" && c.Credentials != nil {
		return fmt.Errorf("only one of 'credentials' or 'identity_token_audience' can be set")
	}

	if c.IdentityTokenAudience != "" && c.ServiceAccountEmail == "" {
		return fmt.Errorf("missing required 'service_account_email' when 'identity_token_audience' is set")
	}

	return nil
}

func (c *gcpConfig) getIAMAlias(role *gcpRole, svcAccount *iam.ServiceAccount) (alias string, err error) {
	aliaser, exists := allowedIAMAliases[c.IAMAliasType]
	if !exists {
		return "", fmt.Errorf("invalid IAM alias type: must be one of: %s", strings.Join(allowedIAMAliasesSlice, ", "))
	}
	return aliaser(role, svcAccount), nil
}

func (c *gcpConfig) getGCEAlias(role *gcpRole, instance *compute.Instance) (alias string, err error) {
	aliaser, exists := allowedGCEAliases[c.GCEAliasType]
	if !exists {
		return "", fmt.Errorf("invalid GCE alias type: must be one of: %s", strings.Join(allowedGCEAliasesSlice, ", "))
	}
	return aliaser(role, instance), nil
}

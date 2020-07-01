package gcpauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/authmetadata"
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
		creds, err := gcputil.Credentials(v.(string))
		if err != nil {
			return errwrap.Wrapf("failed to read credentials: {{err}}", err)
		}

		if len(creds.PrivateKeyId) == 0 {
			return errors.New("missing private key in credentials")
		}

		c.Credentials = creds
	}

	rawIamAlias, exists := d.GetOk("iam_alias")
	if exists {
		iamAlias := rawIamAlias.(string)
		if iamAlias != c.IAMAliasType {
			c.IAMAliasType = iamAlias
		}
	}
	if err := c.IAMAuthMetadata.ParseAuthMetadata(d); err != nil {
		return errwrap.Wrapf("failed to parse iam metadata: {{err}}", err)
	}

	rawGceAlias, exists := d.GetOk("gce_alias")
	if exists {
		gceAlias := rawGceAlias.(string)
		if gceAlias != c.GCEAliasType {
			c.GCEAliasType = gceAlias
		}
	}
	if err := c.GCEAuthMetadata.ParseAuthMetadata(d); err != nil {
		return errwrap.Wrapf("failed to parse gce metadata: {{err}}", err)
	}
	return nil
}

func (c *gcpConfig) getIAMAlias(role *gcpRole, svcAccount *iam.ServiceAccount) (alias string, err error) {
	aliasType := c.IAMAliasType
	if aliasType == "" {
		aliasType = defaultIAMAlias
	}

	aliaser, exists := allowedIAMAliases[aliasType]
	if !exists {
		return "", fmt.Errorf("invalid IAM alias type: must be one of: %s", strings.Join(allowedIAMAliasesSlice, ", "))
	}
	return aliaser(role, svcAccount), nil
}

func (c *gcpConfig) getGCEAlias(role *gcpRole, instance *compute.Instance) (alias string, err error) {
	aliasType := c.GCEAliasType
	if aliasType == "" {
		aliasType = defaultGCEAlias
	}

	aliaser, exists := allowedGCEAliases[aliasType]
	if !exists {
		return "", fmt.Errorf("invalid GCE alias type: must be one of: %s", strings.Join(allowedGCEAliasesSlice, ", "))
	}
	return aliaser(role, instance), nil
}

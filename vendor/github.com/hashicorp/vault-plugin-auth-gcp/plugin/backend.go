package gcpauth

import (
	"context"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
	"net/http"
)

const defaultCloudScope = "https://www.googleapis.com/auth/cloud-platform"

type GcpAuthBackend struct {
	*framework.Backend

	// OAuth scopes for generating HTTP and GCP service clients.
	oauthScopes []string
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *GcpAuthBackend {
	b := &GcpAuthBackend{
		oauthScopes: []string{defaultCloudScope},
	}

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(b),
				pathLogin(b),
			},
			pathsRole(b),
		),
	}
	return b
}

func (b *GcpAuthBackend) httpClient(ctx context.Context, s logical.Storage) (*http.Client, error) {
	config, err := b.config(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf(
			"could not check to see if GCP credentials were configured, error"+
				"reading config: {{err}}", err)
	}

	credsBytes, err := config.formatAndMarshalCredentials()
	if err != nil {
		return nil, errwrap.Wrapf(
			"unable to marshal given GCP credential JSON: {{err}}", err)
	}

	var creds *google.Credentials
	if config != nil && config.Credentials != nil {
		creds, err = google.CredentialsFromJSON(ctx, credsBytes, b.oauthScopes...)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse credentials: {{err}}", err)
		}
	} else {
		creds, err = google.FindDefaultCredentials(ctx, b.oauthScopes...)
		if err != nil {
			return nil, errwrap.Wrapf(
				"credentials were not configured and Vault could not find "+
					"Application Default Credentials (ADC). Either set ADC or "+
					"configure this auth backend at auth/$MOUNT/config "+
					"(default auth/gcp/config). Error: {{err}}", err)
		}
	}

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())
	client := oauth2.NewClient(cleanCtx, creds.TokenSource)
	return client, nil
}

func (b *GcpAuthBackend) newGcpClients(ctx context.Context, s logical.Storage) (*clientHandles, error) {
	httpC, err := b.httpClient(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf("could not obtain HTTP client: {{err}}", err)
	}

	iamClient, err := iam.New(httpC)
	if err != nil {
		return nil, fmt.Errorf(clientErrorTemplate, "IAM", err)
	}
	iamClient.UserAgent = useragent.String()

	gceClient, err := compute.New(httpC)
	if err != nil {
		return nil, fmt.Errorf(clientErrorTemplate, "Compute", err)
	}
	iamClient.UserAgent = useragent.String()

	crmClient, err := cloudresourcemanager.New(httpC)
	if err != nil {
		return nil, fmt.Errorf(clientErrorTemplate, "Cloud Resource Manager", err)
	}
	crmClient.UserAgent = useragent.String()

	return &clientHandles{
		iam:             iamClient,
		gce:             gceClient,
		resourceManager: crmClient,
	}, nil
}

type clientHandles struct {
	iam             *iam.Service
	gce             *compute.Service
	resourceManager *cloudresourcemanager.Service
}

const backendHelp = `
The GCP backend plugin allows authentication for Google Cloud Platform entities.
Currently, it supports authentication for:

* IAM Service Accounts:
	IAM service accounts provide a signed JSON Web Token (JWT), signed by
	calling GCP APIs directly or via the Vault CL helper.

* GCE VM Instances:
	GCE provide a signed instance metadata JSON Web Token (JWT), obtained from the
	GCE instance metadata server  (http://metadata.google.internal/computeMetadata/v1/instance).
	Using the /service-accounts/<service-account-name>/identity	endpoint, the instance
	can obtain this JWT and pass it to Vault on login.
`

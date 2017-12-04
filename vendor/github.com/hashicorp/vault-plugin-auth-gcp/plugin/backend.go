package gcpauth

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/version"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

type GcpAuthBackend struct {
	*framework.Backend

	// OAuth scopes for generating HTTP and GCP service clients.
	oauthScopes []string

	// Locks for guarding service clients
	clientMutex sync.RWMutex

	// -- GCP service clients --
	iamClient *iam.Service
	gceClient *compute.Service
}

// Factory returns a new backend as logical.Backend.
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *GcpAuthBackend {
	b := &GcpAuthBackend{
		oauthScopes: []string{
			iam.CloudPlatformScope,
			compute.ComputeReadonlyScope,
		},
	}

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Invalidate:  b.invalidate,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
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

func (b *GcpAuthBackend) invalidate(key string) {
	switch key {
	case "config":
		b.Close()
	}
}

// Close deletes created GCP clients in backend.
func (b *GcpAuthBackend) Close() {
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	b.iamClient = nil
	b.gceClient = nil
}

func (b *GcpAuthBackend) IAM(s logical.Storage) (*iam.Service, error) {
	b.clientMutex.RLock()
	if b.iamClient != nil {
		defer b.clientMutex.RUnlock()
		return b.iamClient, nil
	}

	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// Check if client was created during lock switch.
	if b.iamClient == nil {
		err := b.initClients(s)
		if err != nil {
			return nil, err
		}
	}

	return b.iamClient, nil
}

func (b *GcpAuthBackend) GCE(s logical.Storage) (*compute.Service, error) {
	b.clientMutex.RLock()
	if b.gceClient != nil {
		defer b.clientMutex.RUnlock()
		return b.gceClient, nil
	}

	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// Check if client was created during lock switch.
	if b.gceClient == nil {
		err := b.initClients(s)
		if err != nil {
			return nil, err
		}
	}

	return b.gceClient, nil
}

// Initialize attempts to create GCP clients from stored config.
// It does not attempt to claim the client lock.
func (b *GcpAuthBackend) initClients(s logical.Storage) (err error) {
	config, err := b.config(s)
	if err != nil {
		return err
	}

	var httpClient *http.Client
	if config == nil || config.Credentials == nil {
		// Use Application Default Credentials
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())

		httpClient, err = google.DefaultClient(ctx, b.oauthScopes...)
		if err != nil {
			return fmt.Errorf("credentials were not configured and fallback to application default credentials failed: %v", err)
		}
	} else {
		httpClient, err = util.GetHttpClient(config.Credentials, b.oauthScopes...)
		if err != nil {
			return err
		}
	}

	userAgentStr := fmt.Sprintf("(%s %s) Vault/%s", runtime.GOOS, runtime.GOARCH, version.GetVersion().FullVersionNumber(true))

	b.iamClient, err = iam.New(httpClient)
	if err != nil {
		b.Close()
		return err
	}
	b.iamClient.UserAgent = userAgentStr

	b.gceClient, err = compute.New(httpClient)
	if err != nil {
		b.Close()
		return err
	}
	b.gceClient.UserAgent = userAgentStr

	return nil
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

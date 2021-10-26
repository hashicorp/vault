package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/compute/metadata"
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"github.com/hashicorp/vault/api"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

type GCPAuth struct {
	roleName            string
	mountPath           string
	authType            string
	serviceAccountEmail string
}

var _ api.AuthMethod = (*GCPAuth)(nil)

type LoginOption func(a *GCPAuth) error

const (
	iamType          = "iam"
	gceType          = "gce"
	defaultMountPath = "gcp"
	defaultAuthType  = gceType
)

// NewGCPAuth initializes a new GCP auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath, WithIAMAuth
func NewGCPAuth(roleName string, opts ...LoginOption) (*GCPAuth, error) {
	if roleName == "" {
		return nil, fmt.Errorf("no role name provided for login")
	}

	a := &GCPAuth{
		mountPath: defaultMountPath,
		authType:  defaultAuthType,
		roleName:  roleName,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *GCPAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

// Login sets up the required request body for the GCP auth method's /login
// endpoint, and performs a write to it. This method defaults to the "gce"
// auth type unless NewGCPAuth is called with WithIAMAuth().
func (a *GCPAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	loginData := map[string]interface{}{
		"role": a.roleName,
	}
	switch a.authType {
	case gceType:
		if !metadata.OnGCE() {
			return nil, fmt.Errorf("GCE metadata service not available")
		}
		// loginData["jwt"] =
	case iamType:
		jwtResp, err := a.signJWT()
		if err != nil {
			return nil, fmt.Errorf("unable to sign JWT for authenticating to GCP: %w", err)
		}
		loginData["jwt"] = jwtResp.SignedJwt
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with GCP auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *GCPAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

func WithIAMAuth(serviceAccountEmail string) LoginOption {
	return func(a *GCPAuth) error {
		a.serviceAccountEmail = serviceAccountEmail
		a.authType = iamType
		return nil
	}
}

// generate signed JWT token from GCP IAM.
func (a *GCPAuth) signJWT() (*credentialspb.SignJwtResponse, error) {
	ctx := context.Background()
	iamClient, err := credentials.NewIamCredentialsClient(ctx) // can pass option.WithCredentialsFile("path/to/creds.json") as second param if GOOGLE_APPLICATION_CREDENTIALS env var not set
	if err != nil {
		return nil, fmt.Errorf("unable to initialize IAM credentials client: %w", err)
	}
	defer iamClient.Close()

	resourceName := fmt.Sprintf("projects/-/serviceAccounts/%s", a.serviceAccountEmail)
	jwtPayload := map[string]interface{}{
		"aud": fmt.Sprintf("vault/%s", a.roleName),
		"sub": a.serviceAccountEmail,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal jwt payload to json: %w", err)
	}

	signJWTReq := &credentialspb.SignJwtRequest{
		Name:    resourceName,
		Payload: string(payloadBytes),
	}

	jwtResp, err := iamClient.SignJwt(ctx, signJWTReq)
	if err != nil {
		return nil, fmt.Errorf("unable to sign JWT: %w", err)
	}

	return jwtResp, nil
}

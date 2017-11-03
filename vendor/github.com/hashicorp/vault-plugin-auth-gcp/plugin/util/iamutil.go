package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iam/v1"
	"time"

	// TODO(emilymye): Currently square's JOSE library doesn't allow for obtaining claims without providing a key.
	// Replace SermoDigital in this file once it is possible to verify the JWT using claims.
	"github.com/SermoDigital/jose/jws"
)

const (
	serviceAccountTemplate    string = "projects/%s/serviceAccounts/%s"
	serviceAccountKeyTemplate string = "projects/%s/serviceAccounts/%s/keys/%s"
	serviceAccountKeyFileType string = "TYPE_X509_PEM_FILE"
)

func ServiceAccountLoginJwt(
	iamClient *iam.Service, exp time.Time, aud, project, serviceAccount string) (*iam.SignJwtResponse, error) {
	accountResource := fmt.Sprintf(serviceAccountTemplate, project, serviceAccount)

	payload, err := json.Marshal(map[string]interface{}{
		"sub": serviceAccount,
		"aud": aud,
		"exp": exp.Unix(),
	})
	if err != nil {
		return nil, err
	}
	jwtReq := &iam.SignJwtRequest{
		Payload: string(payload),
	}
	return iamClient.Projects.ServiceAccounts.SignJwt(accountResource, jwtReq).Do()
}

// serviceAccount wraps a call to the GCP IAM API to get a service account.
func ServiceAccount(iamClient *iam.Service, accountId, projectName string) (*iam.ServiceAccount, error) {
	accountResource := fmt.Sprintf(serviceAccountTemplate, projectName, accountId)
	account, err := iamClient.Projects.ServiceAccounts.Get(accountResource).Do()
	if err != nil {
		return nil, fmt.Errorf("service account '%s' does not exist", accountResource)
	}

	return account, nil
}

// ServiceAccountKey wraps a call to the GCP IAM API to get a service account key.
func ServiceAccountKey(iamClient *iam.Service, keyId, accountId, projectName string) (*iam.ServiceAccountKey, error) {
	keyResource := fmt.Sprintf(serviceAccountKeyTemplate, projectName, accountId, keyId)
	key, err := iamClient.Projects.ServiceAccounts.Keys.Get(keyResource).PublicKeyType(serviceAccountKeyFileType).Do()
	if err != nil {
		return nil, fmt.Errorf("service account key '%s' does not exist: %v", keyResource, err)
	}
	return key, nil
}

// ParseServiceAccountFromIAMJWT parses the service account from the 'sub' claim given a serialized signed JWT.
func ParseServiceAccountFromIAMJWT(signedJwt string) (string, error) {
	jwtVal, err := jws.ParseJWT([]byte(signedJwt))
	if err != nil {
		return "", fmt.Errorf("could not parse service account from JWT 'sub' claim: %v", err)
	}
	accountId, ok := jwtVal.Claims().Subject()
	if !ok {
		return "", errors.New("expected 'sub' claim with service account ID or name")
	}
	return accountId, nil
}

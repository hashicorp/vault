package gcputil

import (
	"fmt"

	"google.golang.org/api/iam/v1"
)

const (
	// ServiceAccountTemplate is used with Google IAM v1.
	//
	// Deprecated: Use ServiceAccountCredentialsTemplate with Service Account Credentials API v1
	// instead. See https://cloud.google.com/iam/docs/migrating-to-credentials-api
	// ServiceAccountTemplate is used with
	// https://pkg.go.dev/google.golang.org/api@v0.3.0/iam/v1
	ServiceAccountTemplate = "projects/%s/serviceAccounts/%s"

	// ServiceAccountCredentialsTemplate is used with
	// https://pkg.go.dev/google.golang.org/api@v0.3.0/iamcredentials/v1
	ServiceAccountCredentialsTemplate = "projects/-/serviceAccounts/%s"
	ServiceAccountKeyTemplate         = "projects/%s/serviceAccounts/%s/keys/%s"
	ServiceAccountKeyFileType         = "TYPE_X509_PEM_FILE"
)

type ServiceAccountId struct {
	Project   string
	EmailOrId string
}

func (id *ServiceAccountId) ResourceName() string {
	return fmt.Sprintf(ServiceAccountTemplate, id.Project, id.EmailOrId)
}

type ServiceAccountKeyId struct {
	Project   string
	EmailOrId string
	Key       string
}

func (id *ServiceAccountKeyId) ResourceName() string {
	return fmt.Sprintf(ServiceAccountKeyTemplate, id.Project, id.EmailOrId, id.Key)
}

// ServiceAccount wraps a call to the GCP IAM API to get a service account.
func ServiceAccount(iamClient *iam.Service, accountId *ServiceAccountId) (*iam.ServiceAccount, error) {
	account, err := iamClient.Projects.ServiceAccounts.Get(accountId.ResourceName()).Do()
	if err != nil {
		return nil, fmt.Errorf("could not find service account '%s': %v", accountId.ResourceName(), err)
	}

	return account, nil
}

// ServiceAccountKey wraps a call to the GCP IAM API to get a service account key.
func ServiceAccountKey(iamClient *iam.Service, keyId *ServiceAccountKeyId) (*iam.ServiceAccountKey, error) {
	keyResource := keyId.ResourceName()
	key, err := iamClient.Projects.ServiceAccounts.Keys.Get(keyId.ResourceName()).PublicKeyType(ServiceAccountKeyFileType).Do()
	if err != nil {
		return nil, fmt.Errorf("could not find service account key '%s': %v", keyResource, err)
	}
	return key, nil
}

package provider

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
)

// Environmental virables that may be used by the provider
const (
	// Deprecated: don't use it outside of this project
	ENVAccessKeyID = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	// Deprecated: don't use it outside of this project
	ENVAccessKeySecret = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
	// Deprecated: don't use it outside of this project
	ENVCredentialFile = "ALIBABA_CLOUD_CREDENTIALS_FILE"
	// Deprecated: don't use it outside of this project
	ENVEcsMetadata = "ALIBABA_CLOUD_ECS_METADATA"
)

// When you want to customize the provider, you only need to implement the method of the interface.
type Provider interface {
	Resolve() (auth.Credential, error)
}

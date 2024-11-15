package credentials

import (
	"fmt"
	"os"
)

type EnvironmentVariableCredentialsProvider struct {
}

func NewEnvironmentVariableCredentialsProvider() (provider *EnvironmentVariableCredentialsProvider) {
	return &EnvironmentVariableCredentialsProvider{}
}

func (provider *EnvironmentVariableCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")

	if accessKeyId == "" {
		err = fmt.Errorf("unable to get credentials from enviroment variables, Access key ID must be specified via environment variable (ALIBABA_CLOUD_ACCESS_KEY_ID)")
		return
	}

	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

	if accessKeySecret == "" {
		err = fmt.Errorf("unable to get credentials from enviroment variables, Access key secret must be specified via environment variable (ALIBABA_CLOUD_ACCESS_KEY_SECRET)")
		return
	}

	securityToken := os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN")

	cc = &Credentials{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		SecurityToken:   securityToken,
		ProviderName:    provider.GetProviderName(),
	}

	return
}

func (provider *EnvironmentVariableCredentialsProvider) GetProviderName() string {
	return "env"
}

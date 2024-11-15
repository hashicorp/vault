// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"os"
	"path"
)

const (
	//ResourcePrincipalVersion2_2 supported version for resource principals
	ResourcePrincipalVersion2_2 = "2.2"
	//ResourcePrincipalVersionEnvVar environment var name for version
	ResourcePrincipalVersionEnvVar = "OCI_RESOURCE_PRINCIPAL_VERSION"
	//ResourcePrincipalRPSTEnvVar environment var name holding the token or a path to the token
	ResourcePrincipalRPSTEnvVar = "OCI_RESOURCE_PRINCIPAL_RPST"
	//ResourcePrincipalPrivatePEMEnvVar environment var holding a rsa private key in pem format or a path to one
	ResourcePrincipalPrivatePEMEnvVar = "OCI_RESOURCE_PRINCIPAL_PRIVATE_PEM"
	//ResourcePrincipalPrivatePEMPassphraseEnvVar environment var holding the passphrase to a key or a path to one
	ResourcePrincipalPrivatePEMPassphraseEnvVar = "OCI_RESOURCE_PRINCIPAL_PRIVATE_PEM_PASSPHRASE"
	//ResourcePrincipalRegionEnvVar environment variable holding a region
	ResourcePrincipalRegionEnvVar = "OCI_RESOURCE_PRINCIPAL_REGION"

	// TenancyOCIDClaimKey is the key used to look up the resource tenancy in an RPST
	TenancyOCIDClaimKey = "res_tenant"
	// CompartmentOCIDClaimKey is the key used to look up the resource compartment in an RPST
	CompartmentOCIDClaimKey = "res_compartment"
)

// ConfigurationProviderWithClaimAccess mixes in a method to access the claims held on the underlying security token
type ConfigurationProviderWithClaimAccess interface {
	common.ConfigurationProvider
	ClaimHolder
}

// ResourcePrincipalConfigurationProvider returns a resource principal configuration provider using well known
// environment variables to look up token information. The environment variables can either paths or contain the material value
// of the keys. However in the case of the keys and tokens paths and values can not be mixed
func ResourcePrincipalConfigurationProvider() (ConfigurationProviderWithClaimAccess, error) {
	var version string
	var ok bool
	if version, ok = os.LookupEnv(ResourcePrincipalVersionEnvVar); !ok {
		return nil, fmt.Errorf("can not create resource principal, environment variable: %s, not present", ResourcePrincipalVersionEnvVar)
	}

	switch version {
	case ResourcePrincipalVersion2_2:
		rpst := requireEnv(ResourcePrincipalRPSTEnvVar)
		if rpst == nil {
			return nil, fmt.Errorf("can not create resource principal, environment variable: %s, not present", ResourcePrincipalRPSTEnvVar)
		}
		private := requireEnv(ResourcePrincipalPrivatePEMEnvVar)
		if private == nil {
			return nil, fmt.Errorf("can not create resource principal, environment variable: %s, not present", ResourcePrincipalPrivatePEMEnvVar)
		}
		passphrase := requireEnv(ResourcePrincipalPrivatePEMPassphraseEnvVar)
		region := requireEnv(ResourcePrincipalRegionEnvVar)
		if region == nil {
			return nil, fmt.Errorf("can not create resource principal, environment variable: %s, not present", ResourcePrincipalRegionEnvVar)
		}
		return newResourcePrincipalKeyProvider22(
			*rpst, *private, passphrase, *region)
	default:
		return nil, fmt.Errorf("can not create resource principal, environment variable: %s, must be valid", ResourcePrincipalVersionEnvVar)
	}
}

func requireEnv(key string) *string {
	if val, ok := os.LookupEnv(key); ok {
		return &val
	}
	return nil
}

// resourcePrincipalKeyProvider22 is key provider that reads from specified the specified environment variables
// the environment variables can host the material keys/passphrases or they can be paths to files that need to be read
type resourcePrincipalKeyProvider struct {
	FederationClient  federationClient
	KeyProviderRegion common.Region
}

func newResourcePrincipalKeyProvider22(sessionTokenLocation, privatePemLocation string,
	passphraseLocation *string, region string) (*resourcePrincipalKeyProvider, error) {

	//Check both the the passphrase and the key are paths
	if passphraseLocation != nil && (!isPath(privatePemLocation) && isPath(*passphraseLocation) ||
		isPath(privatePemLocation) && !isPath(*passphraseLocation)) {
		return nil, fmt.Errorf("cant not create resource principal: both key and passphrase need to be path or none needs to be path")
	}

	var supplier sessionKeySupplier
	var err error

	//File based case
	if isPath(privatePemLocation) {
		supplier, err = newFileBasedKeySessionSupplier(privatePemLocation, passphraseLocation)
		if err != nil {
			return nil, fmt.Errorf("can not create resource principal, due to: %s ", err.Error())
		}
	} else {
		//else the content is in the env vars
		var passphrase []byte
		if passphraseLocation != nil {
			passphrase = []byte(*passphraseLocation)
		}
		supplier, err = newStaticKeySessionSupplier([]byte(privatePemLocation), passphrase)
		if err != nil {
			return nil, fmt.Errorf("can not create resource principal, due to: %s ", err.Error())
		}
	}

	var fd federationClient
	if isPath(sessionTokenLocation) {
		fd, _ = newFileBasedFederationClient(sessionTokenLocation, supplier)
	} else {
		fd, err = newStaticFederationClient(sessionTokenLocation, supplier)
		if err != nil {
			return nil, fmt.Errorf("can not create resource principal, due to: %s ", err.Error())
		}
	}

	rs := resourcePrincipalKeyProvider{
		FederationClient:  fd,
		KeyProviderRegion: common.StringToRegion(region),
	}
	return &rs, nil
}

func (p *resourcePrincipalKeyProvider) PrivateRSAKey() (privateKey *rsa.PrivateKey, err error) {
	if privateKey, err = p.FederationClient.PrivateKey(); err != nil {
		err = fmt.Errorf("failed to get private key: %s", err.Error())
		return nil, err
	}
	return privateKey, nil
}

func (p *resourcePrincipalKeyProvider) KeyID() (string, error) {
	var securityToken string
	var err error
	if securityToken, err = p.FederationClient.SecurityToken(); err != nil {
		return "", fmt.Errorf("failed to get security token: %s", err.Error())
	}
	return fmt.Sprintf("ST$%s", securityToken), nil
}

func (p *resourcePrincipalKeyProvider) Region() (string, error) {
	return string(p.KeyProviderRegion), nil
}

var (
	// ErrNonStringClaim is returned if the token has a claim for a key, but it's not a string value
	ErrNonStringClaim = errors.New("claim does not have a string value")
)

func (p *resourcePrincipalKeyProvider) TenancyOCID() (string, error) {
	if claim, err := p.GetClaim(TenancyOCIDClaimKey); err != nil {
		return "", err
	} else if tenancy, ok := claim.(string); ok {
		return tenancy, nil
	} else {
		return "", ErrNonStringClaim
	}
}

func (p *resourcePrincipalKeyProvider) GetClaim(claim string) (interface{}, error) {
	return p.FederationClient.GetClaim(claim)
}

func (p *resourcePrincipalKeyProvider) KeyFingerprint() (string, error) {
	return "", nil
}

func (p *resourcePrincipalKeyProvider) UserOCID() (string, error) {
	return "", nil
}

// By contract for the the content of a resource principal to be considered path, it needs to be
// an absolute path.
func isPath(str string) bool {
	return path.IsAbs(str)
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/transit/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ConfigureWrapper             = configureWrapper
	CreateSecureRandomReaderFunc = createSecureRandomReader

	AliCloudKMSEnvVars = []string{"ALICLOUD_REGION", "ALICLOUD_DOMAIN", "ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "VAULT_ALICLOUDKMS_SEAL_KEY_ID"}
	AWSKMSEnvVars      = []string{"AWS_REGION", "AWS_DEFAULT_REGION", "AWS_ACCESS_KEY_ID", "AWS_SESSION_TOKEN", "AWS_SECRET_ACCESS_KEY", "VAULT_AWSKMS_SEAL_KEY_ID", "AWS_KMS_ENDPOINT"}
	AzureEnvVars       = []string{"AZURE_TENANT_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET", "AZURE_ENVIRONMENT", "VAULT_AZUREKEYVAULT_VAULT_NAME", "VAULT_AZUREKEYVAULT_KEY_NAME", "AZURE_AD_RESOURCE"}
	GCPCKMSEnvVars     = []string{"GOOGLE_CREDENTIALS", "GOOGLE_APPLICATION_CREDENTIALS", "GOOGLE_PROJECT", "GOOGLE_REGION", "VAULT_GCPCKMS_SEAL_KEY_RING", "VAULT_GCPCKMS_SEAL_CRYPTO_KEY"}
	OCIKMSEnvVars      = []string{"VAULT_OCIKMS_SEAL_KEY_ID", "VAULT_OCIKMS_CRYPTO_ENDPOINT", "VAULT_OCIKMS_MANAGEMENT_ENDPOINT"}
	TransitEnvVars     = []string{"VAULT_ADDR", "VAULT_TOKEN", "VAULT_TRANSIT_SEAL_KEY_NAME", "VAULT_TRANSIT_SEAL_MOUNT_PATH", "VAULT_NAMESPACE", "VAULT_TRANSIT_SEAL_DISABLE_RENEWAL", "VAULT_CACERT", "VAULT_CLIENT_CERT", "VAULT_CLIENT_KEY", "VAULT_TLS_SERVER_NAME", "VAULT_SKIP_VERIFY"}
)

// Entropy contains Entropy configuration for the server
type EntropyMode int

const (
	EntropyUnknown EntropyMode = iota
	EntropyAugmentation
)

type Entropy struct {
	Mode EntropyMode
}

// KMS contains KMS configuration for the server
type KMS struct {
	UnusedKeys []string `hcl:",unusedKeys"`
	Type       string
	// Purpose can be used to allow a string-based specification of what this
	// KMS is designated for, in situations where we want to allow more than
	// one KMS to be specified
	Purpose []string `hcl:"-"`

	Disabled bool
	Config   map[string]string

	Priority int
	Name     string
}

func (k *KMS) GoString() string {
	return fmt.Sprintf("*%#v", *k)
}

func parseKMS(result *[]*KMS, list *ast.ObjectList, blockName string, maxKMS int) error {
	if len(list.Items) > maxKMS {
		return fmt.Errorf("only two or less %q blocks are permitted", blockName)
	}

	seals := make([]*KMS, 0, len(list.Items))
	for _, item := range list.Items {
		key := blockName
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		// We first decode into a map[string]interface{} because purpose isn't
		// necessarily a string. Then we migrate everything else over to
		// map[string]string and error if it doesn't work.
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
		}

		var purpose []string
		var err error
		if v, ok := m["purpose"]; ok {
			if purpose, err = parseutil.ParseCommaStringSlice(v); err != nil {
				return multierror.Prefix(fmt.Errorf("unable to parse 'purpose' in kms type %q: %w", key, err), fmt.Sprintf("%s.%s:", blockName, key))
			}
			for i, p := range purpose {
				purpose[i] = strings.ToLower(p)
			}
			delete(m, "purpose")
		}

		var disabled bool
		if v, ok := m["disabled"]; ok {
			disabled, err = parseutil.ParseBool(v)
			if err != nil {
				return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
			}
			delete(m, "disabled")
		}

		var priority int
		if v, ok := m["priority"]; ok {
			priority, err = parseutil.SafeParseInt(v)
			if err != nil {
				return multierror.Prefix(fmt.Errorf("unable to parse 'priority' in kms type %q: %w", key, err), fmt.Sprintf("%s.%s", blockName, key))
			}
			delete(m, "priority")
		}

		name := strings.ToLower(key)
		if v, ok := m["name"]; ok {
			name, ok = v.(string)
			if !ok {
				return multierror.Prefix(fmt.Errorf("unable to parse 'name' in kms type %q: unexpected type %T", key, v), fmt.Sprintf("%s.%s", blockName, key))
			}
			delete(m, "name")
		}

		strMap := make(map[string]string, len(m))
		for k, v := range m {
			s, err := parseutil.ParseString(v)
			if err != nil {
				return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
			}
			strMap[k] = s
		}

		seal := &KMS{
			Type:     strings.ToLower(key),
			Purpose:  purpose,
			Disabled: disabled,
			Priority: priority,
			Name:     name,
		}
		if len(strMap) > 0 {
			seal.Config = strMap
		}
		seals = append(seals, seal)
	}

	*result = append(*result, seals...)

	return nil
}

func ParseKMSes(d string) ([]*KMS, error) {
	// Parse!
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, err
	}

	// Start building the result
	var result struct {
		Seals []*KMS `hcl:"-"`
	}

	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	if o := list.Filter("seal"); len(o.Items) > 0 {
		if err := parseKMS(&result.Seals, o, "seal", 3); err != nil {
			return nil, fmt.Errorf("error parsing 'seal': %w", err)
		}
	}

	if o := list.Filter("kms"); len(o.Items) > 0 {
		if err := parseKMS(&result.Seals, o, "kms", 3); err != nil {
			return nil, fmt.Errorf("error parsing 'kms': %w", err)
		}
	}

	return result.Seals, nil
}

func configureWrapper(configKMS *KMS, infoKeys *[]string, info *map[string]string, logger hclog.Logger, opts ...wrapping.Option) (wrapping.Wrapper, error) {
	var wrapper wrapping.Wrapper
	var kmsInfo map[string]string
	var err error

	switch wrapping.WrapperType(configKMS.Type) {
	case wrapping.WrapperTypeShamir:
		return nil, nil

	case wrapping.WrapperTypeAead:
		wrapper, kmsInfo, err = GetAEADKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypeAliCloudKms:
		wrapper, kmsInfo, err = GetAliCloudKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypeAwsKms:
		wrapper, kmsInfo, err = GetAWSKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypeAzureKeyVault:
		wrapper, kmsInfo, err = GetAzureKeyVaultKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypeGcpCkms:
		wrapper, kmsInfo, err = GetGCPCKMSKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypeOciKms:
		if keyId, ok := configKMS.Config["key_id"]; ok {
			opts = append(opts, wrapping.WithKeyId(keyId))
		}
		wrapper, kmsInfo, err = GetOCIKMSKMSFunc(configKMS, opts...)
	case wrapping.WrapperTypeTransit:
		wrapper, kmsInfo, err = GetTransitKMSFunc(configKMS, opts...)

	case wrapping.WrapperTypePkcs11:
		return nil, fmt.Errorf("KMS type 'pkcs11' requires the Vault Enterprise HSM binary")

	default:
		return nil, fmt.Errorf("Unknown KMS type %q", configKMS.Type)
	}

	if err != nil {
		return nil, err
	}

	if infoKeys != nil && info != nil {
		for k, v := range kmsInfo {
			*infoKeys = append(*infoKeys, k)
			(*info)[k] = v
		}
	}

	return wrapper, nil
}

func GetAEADKMSFunc(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	wrapper := aeadwrapper.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		str := "AEAD Type"
		if len(kms.Purpose) > 0 {
			str = fmt.Sprintf("%v %s", kms.Purpose, str)
		}
		info[str] = wrapperInfo.Metadata["aead_type"]
	}
	return wrapper, info, nil
}

func GetAliCloudKMSFunc(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getAliCloudEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := alicloudkms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["AliCloud KMS Region"] = wrapperInfo.Metadata["region"]
		info["AliCloud KMS KeyID"] = wrapperInfo.Metadata["kms_key_id"]
		if domain, ok := wrapperInfo.Metadata["domain"]; ok {
			info["AliCloud KMS Domain"] = domain
		}
	}
	return wrapper, info, nil
}

var GetAWSKMSFunc = func(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getAWSKMSEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := awskms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["AWS KMS Region"] = wrapperInfo.Metadata["region"]
		info["AWS KMS KeyID"] = wrapperInfo.Metadata["kms_key_id"]
		if endpoint, ok := wrapperInfo.Metadata["endpoint"]; ok {
			info["AWS KMS Endpoint"] = endpoint
		}
	}
	return wrapper, info, nil
}

func GetAzureKeyVaultKMSFunc(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getAzureEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := azurekeyvault.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["Azure Environment"] = wrapperInfo.Metadata["environment"]
		info["Azure Vault Name"] = wrapperInfo.Metadata["vault_name"]
		info["Azure Key Name"] = wrapperInfo.Metadata["key_name"]
	}
	return wrapper, info, nil
}

func GetGCPCKMSKMSFunc(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getGCPCKMSEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := gcpckms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["GCP KMS Project"] = wrapperInfo.Metadata["project"]
		info["GCP KMS Region"] = wrapperInfo.Metadata["region"]
		info["GCP KMS Key Ring"] = wrapperInfo.Metadata["key_ring"]
		info["GCP KMS Crypto Key"] = wrapperInfo.Metadata["crypto_key"]
	}
	return wrapper, info, nil
}

func GetOCIKMSKMSFunc(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getOCIKMSEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := ocikms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["OCI KMS KeyID"] = wrapperInfo.Metadata[ocikms.KmsConfigKeyId]
		info["OCI KMS Crypto Endpoint"] = wrapperInfo.Metadata[ocikms.KmsConfigCryptoEndpoint]
		info["OCI KMS Management Endpoint"] = wrapperInfo.Metadata[ocikms.KmsConfigManagementEndpoint]
		info["OCI KMS Principal Type"] = wrapperInfo.Metadata["principal_type"]
	}
	return wrapper, info, nil
}

var GetTransitKMSFunc = func(kms *KMS, opts ...wrapping.Option) (wrapping.Wrapper, map[string]string, error) {
	envVarSuffix := ""
	if kms.Priority > 1 {
		envVarSuffix = kms.Name
	}

	envConfig := getTransitEnvConfig(envVarSuffix)
	for name, val := range envConfig {
		kms.Config[name] = val
	}

	wrapper := transit.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithConfigMap(kms.Config))...)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["Transit Address"] = wrapperInfo.Metadata["address"]
		info["Transit Mount Path"] = wrapperInfo.Metadata["mount_path"]
		info["Transit Key Name"] = wrapperInfo.Metadata["key_name"]
		if namespace, ok := wrapperInfo.Metadata["namespace"]; ok {
			info["Transit Namespace"] = namespace
		}
	}
	return wrapper, info, nil
}

func createSecureRandomReader(conf *SharedConfig, wrapper wrapping.Wrapper) (io.Reader, error) {
	return rand.Reader, nil
}

func getAzureEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range AzureEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "AZURE_TENANT_ID":
				envValues["tenant_id"] = val
			case "AZURE_CLIENT_ID":
				envValues["client_id"] = val
			case "AZURE_CLIENT_SECRET":
				envValues["client_secret"] = val
			case "AZURE_ENVIRONMENT":
				envValues["environment"] = val
			case "VAULT_AZUREKEYVAULT_VAULT_NAME":
				envValues["vault_name"] = val
			case "VAULT_AZUREKEYVAULT_KEY_NAME":
				envValues["key_name"] = val
			case "AZURE_AD_RESOURCE":
				envValues["resource"] = val
			}
		}
	}

	return envValues
}

func getAliCloudEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range AliCloudKMSEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "ALICLOUD_REGION":
				envValues["region"] = val
			case "ALICLOUD_DOMAIN":
				envValues["domain"] = val
			case "ALICLOUD_ACCESS_KEY":
				envValues["access_key"] = val
			case "ALICLOUD_SECRET_KEY":
				envValues["secret_key"] = val
			case "VAULT_ALICLOUDKMS_SEAL_KEY_ID":
				envValues["kms_key_id"] = val
			}
		}
	}

	return envValues
}

func getAWSKMSEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range AWSKMSEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "AWS_REGION", "AWS_DEFAULT_REGION":
				envValues["region"] = val
			case "AWS_ACCESS_KEY_ID":
				envValues["access_key"] = val
			case "AWS_SESSION_TOKEN":
				envValues["session_token"] = val
			case "AWS_SECRET_ACCESS_KEY":
				envValues["secret_key"] = val
			case "AWS_AWSKMS_SEAL_KEY_ID":
				envValues["kms_key_id"] = val
			case "AWS_KMS_ENDPOINT":
				envValues["endpoint"] = val
			}
		}
	}

	return envValues
}

func getGCPCKMSEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range GCPCKMSEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "GOOGLE_CREDENTIALS", "GOOGLE_APPLICATION_CREDENTIALS":
				envValues["credentials"] = val
			case "GOOGLE_PROJECT":
				envValues["project"] = val
			case "GOOGLE_REGION":
				envValues["region"] = val
			case "VAULT_GCPCKMS_SEAL_KEY_RING":
				envValues["key_ring"] = val
			case "VAULT_GCPCKMS_SEAL_CRYPTO_KEY":
				envValues["crypto_key"] = val
			}
		}
	}

	return envValues
}

func getOCIKMSEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range OCIKMSEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "VAULT_OCIKMS_SEAL_KEY_ID":
				envValues["key_id"] = val
			case "VAULT_OCIKMS_CRYPTO_ENDPOINT":
				envValues["crypto_endpoint"] = val
			case "VAULT_OCIKMS_MANAGEMENT_ENDPOINT":
				envValues["management_endpoint"] = val
			}
		}
	}
	return envValues
}

func getTransitEnvConfig(suffix string) map[string]string {
	envValues := make(map[string]string)

	for _, envVar := range TransitEnvVars {
		val := os.Getenv(fmt.Sprintf("%s_%s", envVar, suffix))
		if val != "" {
			switch envVar {
			case "VAULT_ADDR":
				envValues["address"] = val
			case "VAULT_TOKEN":
				envValues["token"] = val
			case "VAULT_TRANSIT_SEAL_KEY_NAME":
				envValues["key_name"] = val
			case "VAULT_TRANSIT_SEAL_MOUNT_PATH":
				envValues["mount_paht"] = val
			case "VAULT_NAMESPACE":
				envValues["namespace"] = val
			case "VAULT_TRANSIT_SEAL_DISABLE_RENEWAL":
				envValues["disable_renewal"] = val
			case "VAULT_CACERT":
				envValues["tls_ca_cert"] = val
			case "VAULT_CLIENT_CERT":
				envValues["tls_client_cert"] = val
			case "VAULT_CLIENT_KEY":
				envValues["tls_client_key"] = val
			case "VAULT_TLS_SERVER_NAME":
				envValues["tls_server_name"] = val
			case "VAULT_SKIP_VERIFY":
				envValues["tls_skip_verify"] = val
			}
		}
	}
	return envValues
}

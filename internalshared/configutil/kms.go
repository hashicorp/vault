// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-kms-wrapping/entropy/v2"
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
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ConfigureWrapper             = configureWrapper
	CreateSecureRandomReaderFunc = createSecureRandomReader
	GetEnvConfigFunc             = getEnvConfig
)

// Entropy contains Entropy configuration for the server
type EntropyMode int

const (
	EntropyUnknown EntropyMode = iota
	EntropyAugmentation
)

type Entropy struct {
	Mode     EntropyMode
	SealName string
}

type EntropySourcerInfo struct {
	Sourcer entropy.Sourcer
	Name    string
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

	Priority int    `hcl:"priority"`
	Name     string `hcl:"name"`
}

func (k *KMS) GoString() string {
	return fmt.Sprintf("*%#v", *k)
}

func parseKMS(result *[]*KMS, list *ast.ObjectList, blockName string, maxKMS int) error {
	if len(list.Items) > maxKMS {
		return fmt.Errorf("only %d or less %q blocks are permitted", maxKMS, blockName)
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

			if priority < 1 {
				return multierror.Prefix(fmt.Errorf("invalid priority in kms type %q: %d", key, priority), fmt.Sprintf("%s.%s", blockName, key))
			}
		}

		name := strings.ToLower(key)
		// ensure that seals of the same type will have unique names for seal migration
		if disabled {
			name += "-disabled"
		}
		if v, ok := m["name"]; ok {
			name, ok = v.(string)
			if !ok {
				return multierror.Prefix(fmt.Errorf("unable to parse 'name' in kms type %q: unexpected type %T", key, v), fmt.Sprintf("%s.%s", blockName, key))
			}
			delete(m, "name")

			if !regexp.MustCompile("^[a-zA-Z0-9-_]+$").MatchString(name) {
				return multierror.Prefix(errors.New("'name' field can only include alphanumeric characters, hyphens, and underscores"), fmt.Sprintf("%s.%s", blockName, key))
			}
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

	envConfig := GetEnvConfigFunc(configKMS)
	// transit is a special case, because some config values take precedence over env vars
	if configKMS.Type == wrapping.WrapperTypeTransit.String() {
		mergeTransitConfig(configKMS.Config, envConfig)
	} else {
		for name, val := range envConfig {
			configKMS.Config[name] = val
		}
	}

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
	wrapper := alicloudkms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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
	wrapper := awskms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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
	wrapper := azurekeyvault.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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
	wrapper := gcpckms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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
	wrapper := ocikms.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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
	wrapper := transit.NewWrapper()
	wrapperInfo, err := wrapper.SetConfig(context.Background(), append(opts, wrapping.WithDisallowEnvVars(true), wrapping.WithConfigMap(kms.Config))...)
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

func createSecureRandomReader(_ *SharedConfig, _ []*EntropySourcerInfo, _ hclog.Logger) (io.Reader, error) {
	return rand.Reader, nil
}

func getEnvConfig(kms *KMS) map[string]string {
	envValues := make(map[string]string)

	var wrapperEnvVars map[string]string
	switch wrapping.WrapperType(kms.Type) {
	case wrapping.WrapperTypeAliCloudKms:
		wrapperEnvVars = AliCloudKMSEnvVars
	case wrapping.WrapperTypeAwsKms:
		wrapperEnvVars = AWSKMSEnvVars
	case wrapping.WrapperTypeAzureKeyVault:
		wrapperEnvVars = AzureEnvVars
	case wrapping.WrapperTypeGcpCkms:
		wrapperEnvVars = GCPCKMSEnvVars
	case wrapping.WrapperTypeOciKms:
		wrapperEnvVars = OCIKMSEnvVars
	case wrapping.WrapperTypeTransit:
		wrapperEnvVars = TransitEnvVars
	default:
		return nil
	}

	for envVar, configName := range wrapperEnvVars {
		val := os.Getenv(envVar)
		if val != "" {
			envValues[configName] = val
		}
	}

	return envValues
}

func mergeTransitConfig(config map[string]string, envConfig map[string]string) {
	useFileTlsConfig := false
	for _, varName := range TransitTLSConfigVars {
		if _, ok := config[varName]; ok {
			useFileTlsConfig = true
			break
		}
	}

	if useFileTlsConfig {
		for _, varName := range TransitTLSConfigVars {
			delete(envConfig, varName)
		}
	}

	for varName, val := range envConfig {
		// for some values, file config takes precedence
		if strutil.StrListContains(TransitPrioritizeConfigValues, varName) && config[varName] != "" {
			continue
		}

		config[varName] = val
	}
}

func (k *KMS) Clone() *KMS {
	ret := &KMS{
		UnusedKeys: k.UnusedKeys,
		Type:       k.Type,
		Purpose:    k.Purpose,
		Config:     k.Config,
		Name:       k.Name,
		Disabled:   k.Disabled,
		Priority:   k.Priority,
	}
	return ret
}

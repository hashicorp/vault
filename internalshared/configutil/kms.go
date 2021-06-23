package configutil

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault"
	"github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/ocikms"
	"github.com/hashicorp/go-kms-wrapping/wrappers/transit"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ConfigureWrapper             = configureWrapper
	CreateSecureRandomReaderFunc = createSecureRandomReader
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

func configureWrapper(configKMS *KMS, infoKeys *[]string, info *map[string]string, logger hclog.Logger) (wrapping.Wrapper, error) {
	var wrapper wrapping.Wrapper
	var kmsInfo map[string]string
	var err error

	opts := &wrapping.WrapperOptions{
		Logger: logger,
	}

	switch configKMS.Type {
	case wrapping.Shamir:
		return nil, nil

	case wrapping.AEAD:
		wrapper, kmsInfo, err = GetAEADKMSFunc(opts, configKMS)

	case wrapping.AliCloudKMS:
		wrapper, kmsInfo, err = GetAliCloudKMSFunc(opts, configKMS)

	case wrapping.AWSKMS:
		wrapper, kmsInfo, err = GetAWSKMSFunc(opts, configKMS)

	case wrapping.AzureKeyVault:
		wrapper, kmsInfo, err = GetAzureKeyVaultKMSFunc(opts, configKMS)

	case wrapping.GCPCKMS:
		wrapper, kmsInfo, err = GetGCPCKMSKMSFunc(opts, configKMS)

	case wrapping.OCIKMS:
		wrapper, kmsInfo, err = GetOCIKMSKMSFunc(opts, configKMS)

	case wrapping.Transit:
		wrapper, kmsInfo, err = GetTransitKMSFunc(opts, configKMS)

	case wrapping.PKCS11:
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

func GetAEADKMSFunc(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := aeadwrapper.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		str := "AEAD Type"
		if len(kms.Purpose) > 0 {
			str = fmt.Sprintf("%v %s", kms.Purpose, str)
		}
		info[str] = wrapperInfo["aead_type"]
	}
	return wrapper, info, nil
}

func GetAliCloudKMSFunc(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := alicloudkms.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["AliCloud KMS Region"] = wrapperInfo["region"]
		info["AliCloud KMS KeyID"] = wrapperInfo["kms_key_id"]
		if domain, ok := wrapperInfo["domain"]; ok {
			info["AliCloud KMS Domain"] = domain
		}
	}
	return wrapper, info, nil
}

var GetAWSKMSFunc = func(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := awskms.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["AWS KMS Region"] = wrapperInfo["region"]
		info["AWS KMS KeyID"] = wrapperInfo["kms_key_id"]
		if endpoint, ok := wrapperInfo["endpoint"]; ok {
			info["AWS KMS Endpoint"] = endpoint
		}
	}
	return wrapper, info, nil
}

func GetAzureKeyVaultKMSFunc(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := azurekeyvault.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["Azure Environment"] = wrapperInfo["environment"]
		info["Azure Vault Name"] = wrapperInfo["vault_name"]
		info["Azure Key Name"] = wrapperInfo["key_name"]
	}
	return wrapper, info, nil
}

func GetGCPCKMSKMSFunc(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := gcpckms.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["GCP KMS Project"] = wrapperInfo["project"]
		info["GCP KMS Region"] = wrapperInfo["region"]
		info["GCP KMS Key Ring"] = wrapperInfo["key_ring"]
		info["GCP KMS Crypto Key"] = wrapperInfo["crypto_key"]
	}
	return wrapper, info, nil
}

func GetOCIKMSKMSFunc(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := ocikms.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["OCI KMS KeyID"] = wrapperInfo[ocikms.KMSConfigKeyID]
		info["OCI KMS Crypto Endpoint"] = wrapperInfo[ocikms.KMSConfigCryptoEndpoint]
		info["OCI KMS Management Endpoint"] = wrapperInfo[ocikms.KMSConfigManagementEndpoint]
		info["OCI KMS Principal Type"] = wrapperInfo["principal_type"]
	}
	return wrapper, info, nil
}

var GetTransitKMSFunc = func(opts *wrapping.WrapperOptions, kms *KMS) (wrapping.Wrapper, map[string]string, error) {
	wrapper := transit.NewWrapper(opts)
	wrapperInfo, err := wrapper.SetConfig(kms.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	info := make(map[string]string)
	if wrapperInfo != nil {
		info["Transit Address"] = wrapperInfo["address"]
		info["Transit Mount Path"] = wrapperInfo["mount_path"]
		info["Transit Key Name"] = wrapperInfo["key_name"]
		if namespace, ok := wrapperInfo["namespace"]; ok {
			info["Transit Namespace"] = namespace
		}
	}
	return wrapper, info, nil
}

func createSecureRandomReader(conf *SharedConfig, wrapper wrapping.Wrapper) (io.Reader, error) {
	return rand.Reader, nil
}

package mfa

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pquerna/otp"
)

func methodsListPaths(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "methods/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.mfaBackendMethodList,
		},

		HelpSynopsis:    mfaListMethodsHelp,
		HelpDescription: mfaListMethodsHelp,
	}
}

func methodListPaths(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "method/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.mfaBackendMethodList,
		},

		HelpSynopsis:    mfaListMethodsHelp,
		HelpDescription: mfaListMethodsHelp,
	}
}

func methodPaths(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "method/" + framework.GenericNameRegex("method_name") + "$",
		Fields: map[string]*framework.FieldSchema{
			"method_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: mfaMethodNameHelp,
			},

			"type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: mfaTypesHelp,
			},

			"totp_hash_algorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "sha1",
				Description: mfaTOTPHashAlgorithmHelp,
			},

			"duo_host": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "",
			},

			"duo_ikey": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "",
			},

			"duo_skey": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "",
			},

			"duo_user_agent": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.mfaBackendMethodRead,
			logical.CreateOperation: b.mfaBackendMethodCreateUpdate,
			logical.UpdateOperation: b.mfaBackendMethodCreateUpdate,
			logical.DeleteOperation: b.mfaBackendMethodDelete,
		},

		ExistenceCheck: b.mfaBackendMethodExistenceCheck,

		HelpSynopsis:    mfaPathMethodsHelp,
		HelpDescription: mfaPathMethodsHelp,
	}
}

func (b *backend) mfaBackendMethod(methodName string) (*mfaMethodEntry, error) {
	b.RLock()
	defer b.RUnlock()

	return b.mfaBackendMethodInternal(methodName)
}

func (b *backend) mfaBackendMethodInternal(methodName string) (*mfaMethodEntry, error) {
	entry, err := b.storage.Get(fmt.Sprintf("method/%s/config", strings.ToLower(methodName)))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result mfaMethodEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) mfaBackendMethodList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.RLock()
	defer b.RUnlock()

	entries, err := b.storage.List("method/")
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(entries))
	for i, entry := range entries {
		ret[i] = strings.TrimPrefix(entry, "method/")
	}

	return logical.ListResponse(ret), nil
}

func (b *backend) mfaBackendMethodDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	b.Lock()
	defer b.Unlock()

	err := b.storage.Delete(fmt.Sprintf("method/%s/config", strings.ToLower(methodName)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) mfaBackendMethodRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	method, err := b.mfaBackendMethod(methodName)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(method).Map(),
	}

	// Make sure not to return the secret key
	delete(resp.Data, "duo_skey")

	return resp, nil
}

func (b *backend) mfaBackendMethodExistenceCheck(
	req *logical.Request, data *framework.FieldData) (bool, error) {
	name := data.Get("method_name").(string)
	if name == "" {
		return false, fmt.Errorf("method name cannot be empty")
	}

	method, err := b.mfaBackendMethod(name)
	if err != nil {
		return false, err
	}

	return method != nil, nil
}

func (b *backend) mfaBackendMethodCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	methodType := data.Get("type").(string)
	if methodType == "" {
		return logical.ErrorResponse("type cannot be empty"), nil
	}

	b.Lock()
	defer b.Unlock()

	entry, err := b.mfaBackendMethodInternal(methodName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		entry = &mfaMethodEntry{
			Name: methodName,
		}
	}

	typeInt, ok := data.GetOk("type")
	if ok {
		entry.Type = typeInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.Type = data.Get("type").(string)
	}
	switch entry.Type {
	case "duo", "totp":
	case "":
		return logical.ErrorResponse("type cannot be empty"), nil
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown MFA type %s", entry.Type)), nil
	}

	totpHashAlgorithmInt, ok := data.GetOk("totp_hash_algorithm")
	if ok {
		entry.TOTPHashAlgorithm = totpHashAlgorithmInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.TOTPHashAlgorithm = data.Get("totp_hash_algorithm").(string)
	}
	if entry.Type == "totp" {
		if _, err := entry.totpAlgorithm(); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	duoHostInt, ok := data.GetOk("duo_host")
	if ok {
		entry.DuoHost = duoHostInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.DuoHost = data.Get("duo_host").(string)
	}

	duoIKeyInt, ok := data.GetOk("duo_ikey")
	if ok {
		entry.DuoIKey = duoIKeyInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.DuoIKey = data.Get("duo_ikey").(string)
	}

	duoSKeyInt, ok := data.GetOk("duo_skey")
	if ok {
		entry.DuoSKey = duoSKeyInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.DuoSKey = data.Get("duo_skey").(string)
	}

	duoUserAgentInt, ok := data.GetOk("duo_user_agent")
	if ok {
		entry.DuoUserAgent = duoUserAgentInt.(string)
	} else if req.Operation == logical.CreateOperation {
		entry.DuoUserAgent = data.Get("duo_user_agent").(string)
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON(fmt.Sprintf("method/%s/config", strings.ToLower(methodName)), entry)
	if err != nil {
		return nil, err
	}
	if err := b.storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

type mfaMethodEntry struct {
	// Name, available here for use in other parts of the code
	Name string `json:"name" mapstructure:"name" structs:"name"`

	// The type, such as "duo" or "totp"
	Type string `json:"type" mapstructure:"type" structs:"type"`

	// The hash type, such as "sha1"
	TOTPHashAlgorithm string `json:"totp_hash_algorithm" mapstructure:"totp_hash_algorithm" structs:"totp_hash_algorithm"`

	// The host to use for Duo authentication
	DuoHost string `json:"duo_host" mapstructure:"duo_host" structs:"duo_host"`

	// The integration key for Duo authentication
	DuoIKey string `json:"duo_ikey" mapstructure:"duo_ikey" structs:"duo_ikey"`

	// The secret key for Duo authentication
	DuoSKey string `json:"duo_skey" mapstructure:"duo_skey" structs:"duo_skey"`

	// The user agent for Duo authentication
	DuoUserAgent string `json:"duo_user_agent" mapstructure:"duo_user_agent" structs:"duo_user_agent"`
}

func (me *mfaMethodEntry) totpAlgorithm() (otp.Algorithm, error) {
	switch me.TOTPHashAlgorithm {
	case "sha1":
		return otp.AlgorithmSHA1, nil
	case "sha256":
		return otp.AlgorithmSHA256, nil
	case "sha512":
		return otp.AlgorithmSHA512, nil
	}

	return 0, fmt.Errorf("unknown method hash algorithm %s", me.TOTPHashAlgorithm)
}

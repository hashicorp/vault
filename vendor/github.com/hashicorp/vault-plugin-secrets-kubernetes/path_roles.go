// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultRoleType     = "Role"
	rolesPath           = "roles/"
	defaultNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 8) (.RoleName | truncate 8) (unix_time) (random 24) | truncate 62 | lowercase }}`
)

type roleEntry struct {
	Name                  string            `json:"name" mapstructure:"name"`
	K8sNamespaces         []string          `json:"allowed_kubernetes_namespaces" mapstructure:"allowed_kubernetes_namespaces"`
	K8sNamespaceSelector  string            `json:"allowed_kubernetes_namespace_selector" mapstructure:"allowed_kubernetes_namespace_selector"`
	TokenMaxTTL           time.Duration     `json:"token_max_ttl" mapstructure:"token_max_ttl"`
	TokenDefaultTTL       time.Duration     `json:"token_default_ttl" mapstructure:"token_default_ttl"`
	TokenDefaultAudiences []string          `json:"token_default_audiences" mapstructure:"token_default_audiences"`
	ServiceAccountName    string            `json:"service_account_name" mapstructure:"service_account_name"`
	K8sRoleName           string            `json:"kubernetes_role_name" mapstructure:"kubernetes_role_name"`
	K8sRoleType           string            `json:"kubernetes_role_type" mapstructure:"kubernetes_role_type"`
	RoleRules             string            `json:"generated_role_rules" mapstructure:"generated_role_rules"`
	NameTemplate          string            `json:"name_template" mapstructure:"name_template"`
	ExtraLabels           map[string]string `json:"extra_labels" mapstructure:"extra_labels"`
	ExtraAnnotations      map[string]string `json:"extra_annotations" mapstructure:"extra_annotations"`
}

// HasSingleK8sNamespace returns true if the role has a single namespace specified
// and the label selector for Kubernetes namespaces is empty
func (r *roleEntry) HasSingleK8sNamespace() bool {
	return r.K8sNamespaceSelector == "" &&
		len(r.K8sNamespaces) == 1 && r.K8sNamespaces[0] != "" && r.K8sNamespaces[0] != "*"
}

func (r *roleEntry) toResponseData() (map[string]interface{}, error) {
	respData := map[string]interface{}{}
	if err := mapstructure.Decode(r, &respData); err != nil {
		return nil, err
	}
	// Format the TTLs as seconds
	respData["token_default_ttl"] = r.TokenDefaultTTL.Seconds()
	respData["token_max_ttl"] = r.TokenMaxTTL.Seconds()

	return respData, nil
}

func (b *backend) pathRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: rolesPath + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKubernetes,
				OperationSuffix: "role",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role",
					Required:    true,
				},
				"allowed_kubernetes_namespaces": {
					Type:        framework.TypeCommaStringSlice,
					Description: `A list of the Kubernetes namespaces in which credentials can be generated. If set to "*" all namespaces are allowed.`,
					Required:    false,
				},
				"allowed_kubernetes_namespace_selector": {
					Type:        framework.TypeString,
					Description: `A label selector for Kubernetes namespaces in which credentials can be generated. Accepts either a JSON or YAML object. If set with allowed_kubernetes_namespaces, the conditions are conjuncted.`,
					Required:    false,
				},
				"token_max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The maximum ttl for generated Kubernetes service account tokens. If not set or set to 0, will use system default.",
					Required:    false,
				},
				"token_default_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The default ttl for generated Kubernetes service account tokens. If not set or set to 0, will use system default.",
					Required:    false,
				},
				"token_default_audiences": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The default audiences for generated Kubernetes service account tokens. If not set or set to \"\", will use k8s cluster default.",
					Required:    false,
				},
				"service_account_name": {
					Type:        framework.TypeString,
					Description: "The pre-existing service account to generate tokens for. Mutually exclusive with all role parameters. If set, only a Kubernetes service account token will be created.",
					Required:    false,
				},
				"kubernetes_role_name": {
					Type:        framework.TypeString,
					Description: "The pre-existing Role or ClusterRole to bind a generated service account to. If set, Kubernetes token, service account, and role binding objects will be created.",
					Required:    false,
				},
				"kubernetes_role_type": {
					Type:        framework.TypeString,
					Description: "Specifies whether the Kubernetes role is a Role or ClusterRole.",
					Required:    false,
					Default:     defaultRoleType,
				},
				"generated_role_rules": {
					Type:        framework.TypeString,
					Description: "The Role or ClusterRole rules to use when generating a role. Accepts either a JSON or YAML object. If set, the entire chain of Kubernetes objects will be generated.",
					Required:    false,
				},
				"name_template": {
					Type:        framework.TypeString,
					Description: "The name template to use when generating service accounts, roles and role bindings. If unset, a default template is used.",
					Required:    false,
				},
				"extra_labels": {
					Type:        framework.TypeKVPairs,
					Description: "Additional labels to apply to all generated Kubernetes objects.",
					Required:    false,
				},
				"extra_annotations": {
					Type:        framework.TypeKVPairs,
					Description: "Additional annotations to apply to all generated Kubernetes objects.",
					Required:    false,
				},
			},
			ExistenceCheck: b.pathRoleExistenceCheck("name"),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRolesRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRolesDelete,
				},
			},
			HelpSynopsis:    rolesHelpSynopsis,
			HelpDescription: rolesHelpDescription,
		},
		{
			Pattern: rolesPath + "?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKubernetes,
				OperationSuffix: "roles",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRolesList,
				},
			},
			HelpSynopsis:    pathRolesListHelpSynopsis,
			HelpDescription: pathRolesListHelpDescription,
		},
	}
}

func (b *backend) pathRoleExistenceCheck(roleFieldName string) framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
		rName := d.Get(roleFieldName).(string)
		r, err := getRole(ctx, req.Storage, rName)
		if err != nil {
			return false, err
		}
		return r != nil, nil
	}
}

func (b *backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	respData, err := entry.toResponseData()
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: respData,
	}, nil
}

func (b *backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("role name must be specified"), nil
	}

	entry, err := getRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		entry = &roleEntry{
			Name: name,
		}
	}

	if k8sNamespaces, ok := d.GetOk("allowed_kubernetes_namespaces"); ok {
		// K8s namespaces need to be lowercase
		entry.K8sNamespaces = strutil.RemoveDuplicates(k8sNamespaces.([]string), true)
	}
	if k8sNamespaceSelector, ok := d.GetOk("allowed_kubernetes_namespace_selector"); ok {
		entry.K8sNamespaceSelector = k8sNamespaceSelector.(string)
	}
	if tokenMaxTTLRaw, ok := d.GetOk("token_max_ttl"); ok {
		entry.TokenMaxTTL = time.Duration(tokenMaxTTLRaw.(int)) * time.Second
	}
	if tokenTTLRaw, ok := d.GetOk("token_default_ttl"); ok {
		entry.TokenDefaultTTL = time.Duration(tokenTTLRaw.(int)) * time.Second
	}
	if tokenAudiencesRaw, ok := d.GetOk("token_default_audiences"); ok {
		entry.TokenDefaultAudiences = strutil.RemoveDuplicates(tokenAudiencesRaw.([]string), false)
	}
	if svcAccount, ok := d.GetOk("service_account_name"); ok {
		entry.ServiceAccountName = svcAccount.(string)
	}
	if k8sRoleName, ok := d.GetOk("kubernetes_role_name"); ok {
		entry.K8sRoleName = k8sRoleName.(string)
	}

	if k8sRoleType, ok := d.GetOk("kubernetes_role_type"); ok {
		entry.K8sRoleType = k8sRoleType.(string)
	}
	if entry.K8sRoleType == "" {
		entry.K8sRoleType = defaultRoleType
	}

	if roleRules, ok := d.GetOk("generated_role_rules"); ok {
		entry.RoleRules = roleRules.(string)
	}
	if nameTemplate, ok := d.GetOk("name_template"); ok {
		entry.NameTemplate = nameTemplate.(string)
	}
	if extraLabels, ok := d.GetOk("extra_labels"); ok {
		entry.ExtraLabels = extraLabels.(map[string]string)
	}
	if extraAnnotations, ok := d.GetOk("extra_annotations"); ok {
		entry.ExtraAnnotations = extraAnnotations.(map[string]string)
	}

	// Validate the entry
	if len(entry.K8sNamespaces) == 0 && entry.K8sNamespaceSelector == "" {
		return logical.ErrorResponse("one (at least) of allowed_kubernetes_namespaces or allowed_kubernetes_namespace_selector must be set"), nil
	}
	if !onlyOneSet(entry.ServiceAccountName, entry.K8sRoleName, entry.RoleRules) {
		return logical.ErrorResponse("one (and only one) of service_account_name, kubernetes_role_name or generated_role_rules must be set"), nil
	}
	if entry.TokenMaxTTL > 0 && entry.TokenDefaultTTL > entry.TokenMaxTTL {
		return logical.ErrorResponse("token_default_ttl %s cannot be greater than token_max_ttl %s", entry.TokenDefaultTTL, entry.TokenMaxTTL), nil
	}

	casedRoleType := makeRoleType(entry.K8sRoleType)
	if casedRoleType != "Role" && casedRoleType != "ClusterRole" {
		return logical.ErrorResponse("kubernetes_role_type must be either 'Role' or 'ClusterRole'"), nil
	}
	entry.K8sRoleType = casedRoleType

	// Try parsing the label selector as json or yaml
	if entry.K8sNamespaceSelector != "" {
		if _, err := makeLabelSelector(entry.K8sNamespaceSelector); err != nil {
			return logical.ErrorResponse("failed to parse 'allowed_kubernetes_namespace_selector' as k8s.io/api/meta/v1/LabelSelector object"), nil
		}
	}

	// Try parsing the role rules as json or yaml
	if entry.RoleRules != "" {
		if _, err := makeRules(entry.RoleRules); err != nil {
			return logical.ErrorResponse("failed to parse 'generated_role_rules' as k8s.io/api/rbac/v1/Policy object"), nil
		}
	}

	// verify the template is valid
	nameTemplate := entry.NameTemplate
	if nameTemplate == "" {
		nameTemplate = defaultNameTemplate
	}
	_, err = template.NewTemplate(template.Template(nameTemplate))
	if err != nil {
		return logical.ErrorResponse("unable to initialize name template: %s", err), nil
	}

	if err := setRole(ctx, req.Storage, name, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (resp *logical.Response, err error) {
	rName := d.Get("name").(string)
	if err := req.Storage.Delete(ctx, rolesPath+rName); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathRolesList(ctx context.Context, req *logical.Request, d *framework.FieldData) (resp *logical.Response, err error) {
	roles, err := req.Storage.List(ctx, rolesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return logical.ListResponse(roles), nil
}

func onlyOneSet(vars ...string) bool {
	count := 0
	for _, v := range vars {
		if v != "" {
			count++
		}
	}
	return count == 1
}

func getRole(ctx context.Context, s logical.Storage, name string) (*roleEntry, error) {
	if name == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, rolesPath+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var role roleEntry

	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func setRole(ctx context.Context, s logical.Storage, name string, entry *roleEntry) error {
	jsonEntry, err := logical.StorageEntryJSON(rolesPath+name, entry)
	if err != nil {
		return err
	}

	if jsonEntry == nil {
		return fmt.Errorf("failed to create storage entry for role %q", name)
	}

	return s.Put(ctx, jsonEntry)
}

const (
	rolesHelpSynopsis            = `Manage the roles that can be created with this secrets engine.`
	rolesHelpDescription         = `This path lets you manage the roles that can be created with this secrets engine.`
	pathRolesListHelpSynopsis    = `List the existing roles in this secrets engine.`
	pathRolesListHelpDescription = `A list of existing role names will be returned.`
)

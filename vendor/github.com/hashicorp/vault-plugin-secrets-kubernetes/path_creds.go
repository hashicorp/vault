// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	josejwt "github.com/go-jose/go-jose/v4/jwt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	pathCreds     = "creds/"
	kubeTokenType = "kube_token"

	pathCredsHelpSyn  = `Request Kubernetes service account credentials for a given Vault role.`
	pathCredsHelpDesc = `
This path creates dynamic Kubernetes service account credentials.
The associated Vault role can be configured to generate tokens for an
existing service account, create a new service account bound to an
existing Role/ClusterRole, or create a new service account and role
bindings. The service account token and any other objects created in
Kubernetes will be automatically deleted when the lease has expired.
`
)

// AllowedSigningAlgs contains all signing algorithms supported by k8s OIDC.
// ref: https://github.com/kubernetes/kubernetes/blob/b4935d910dcf256288694391ef675acfbdb8e7a3/staging/src/k8s.io/apiserver/plugin/pkg/authenticator/token/oidc/oidc.go#L222-L233
var AllowedSigningAlgs = []jose.SignatureAlgorithm{
	jose.RS256,
	jose.RS384,
	jose.RS512,
	jose.ES256,
	jose.ES384,
	jose.ES512,
	jose.PS256,
	jose.PS384,
	jose.PS512,
}

type credsRequest struct {
	Namespace          string        `json:"kubernetes_namespace"`
	ClusterRoleBinding bool          `json:"cluster_role_binding"`
	TTL                time.Duration `json:"ttl"`
	RoleName           string        `json:"role_name"`
	Audiences          []string      `json:"audiences"`
}

// The fields in nameMetadata are used for templated name generation
type nameMetadata struct {
	DisplayName string
	RoleName    string
}

func (b *backend) pathCredentials() *framework.Path {
	forwardOperation := &framework.PathOperation{
		Callback:                    b.pathCredentialsRead,
		ForwardPerformanceSecondary: true,
		ForwardPerformanceStandby:   true,
	}
	return &framework.Path{
		Pattern: pathCreds + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKubernetes,
			OperationVerb:   "generate",
			OperationSuffix: "credentials",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the Vault role",
				Required:    true,
			},
			"kubernetes_namespace": {
				Type:        framework.TypeString,
				Description: "The name of the Kubernetes namespace in which to generate the credentials",
				Required:    true,
			},
			"cluster_role_binding": {
				Type:        framework.TypeBool,
				Description: "If true, generate a ClusterRoleBinding to grant permissions across the whole cluster instead of within a namespace. Requires the Vault role to have kubernetes_role_type set to ClusterRole.",
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "The TTL of the generated credentials",
			},
			"audiences": {
				Type:        framework.TypeCommaStringSlice,
				Description: "The intended audiences of the generated credentials",
			},
		},

		HelpSynopsis:    pathCredsHelpSyn,
		HelpDescription: pathCredsHelpDesc,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: forwardOperation,
		},
	}
}

func (b *backend) pathCredentialsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	roleEntry, err := getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' does not exist", roleName)), nil
	}

	request := &credsRequest{
		RoleName: roleName,
	}
	requestNamespace, ok := d.GetOk("kubernetes_namespace")
	if ok {
		request.Namespace = requestNamespace.(string)
	}

	request.ClusterRoleBinding = d.Get("cluster_role_binding").(bool)

	ttlRaw, ok := d.GetOk("ttl")
	if ok {
		request.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}

	audiences, ok := d.Get("audiences").([]string)
	if ok {
		request.Audiences = audiences
	}

	// Validate the request
	isValidNs, err := b.isValidKubernetesNamespace(ctx, req, request, roleEntry)
	if err != nil {
		return nil, fmt.Errorf("error verifying namespace: %w", err)
	}
	if !isValidNs {
		return logical.ErrorResponse(fmt.Sprintf("kubernetes_namespace '%s' is not present in role's allowed_kubernetes_namespaces or does not match role's label selector allowed_kubernetes_namespace_selector", request.Namespace)), nil
	}
	if request.ClusterRoleBinding && roleEntry.K8sRoleType == "Role" {
		return logical.ErrorResponse("a ClusterRoleBinding cannot ref a Role"), nil
	}

	return b.createCreds(ctx, req, roleEntry, request)
}

func (b *backend) isValidKubernetesNamespace(ctx context.Context, req *logical.Request, request *credsRequest, role *roleEntry) (bool, error) {
	if request.Namespace == "" {
		if role.HasSingleK8sNamespace() {
			// Assign the single namespace to the creds request namespace
			request.Namespace = role.K8sNamespaces[0]
			return true, nil
		}

		return false, fmt.Errorf("'kubernetes_namespace' is required unless the Vault role has a single namespace specified")
	}

	if strutil.StrListContains(role.K8sNamespaces, "*") || strutil.StrListContains(role.K8sNamespaces, request.Namespace) {
		return true, nil
	}

	if role.K8sNamespaceSelector == "" {
		return false, nil
	}
	selector, err := makeLabelSelector(role.K8sNamespaceSelector)
	if err != nil {
		return false, err
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	nsLabels, err := client.getNamespaceLabelSet(ctx, request.Namespace)
	if err != nil {
		return false, err
	}
	labelSelector, err := metav1.LabelSelectorAsSelector(&selector)
	if err != nil {
		return false, err
	}
	return labelSelector.Matches(labels.Set(nsLabels)), nil
}

func (b *backend) createCreds(ctx context.Context, req *logical.Request, role *roleEntry, reqPayload *credsRequest) (*logical.Response, error) {
	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	nameTemplate := role.NameTemplate
	if nameTemplate == "" {
		nameTemplate = defaultNameTemplate
	}

	up, err := template.NewTemplate(template.Template(nameTemplate))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize name template: %w", err)
	}
	um := nameMetadata{
		DisplayName: req.DisplayName,
		RoleName:    role.Name,
	}
	genName, err := up.Generate(um)
	if err != nil {
		return nil, fmt.Errorf("failed to generate name: %w", err)
	}

	// Determine the TTL here, since it might come from the mount if nothing on
	// the vault role or creds payload is specified, and we need to know it
	// before creating K8s Token
	theTTL := time.Duration(0)
	switch {
	case reqPayload.TTL > 0:
		theTTL = reqPayload.TTL
	case role.TokenDefaultTTL > 0:
		theTTL = role.TokenDefaultTTL
	default:
		theTTL = b.System().DefaultLeaseTTL()
	}

	var respWarning []string
	// If the calculated TTL is greater than the role's max ttl, it'll be capped
	// by the framework when returned. Catch it here so that the k8s token has
	// the same capped TTL.
	if role.TokenMaxTTL > 0 && theTTL > role.TokenMaxTTL {
		respWarning = append(respWarning, fmt.Sprintf("ttl of %s is greater than the role's token_max_ttl of %s; capping accordingly", theTTL.String(), role.TokenMaxTTL.String()))
		theTTL = role.TokenMaxTTL
	}
	// Similarly, if the calculated TTL is greater than the system's max lease
	// ttl, cap accordingly here.
	if theTTL > b.System().MaxLeaseTTL() {
		respWarning = append(respWarning, fmt.Sprintf("ttl of %s is greater than Vault's max lease ttl %s; capping accordingly", theTTL.String(), b.System().MaxLeaseTTL().String()))
		theTTL = b.System().MaxLeaseTTL()
	}

	theAudiences := role.TokenDefaultAudiences
	if len(reqPayload.Audiences) != 0 {
		theAudiences = reqPayload.Audiences
	}

	// These are created items to save internally and/or return to the caller
	token := ""
	serviceAccountName := ""
	createdServiceAccountName := ""
	createdK8sRoleBinding := ""
	createdK8sRole := ""

	var walID string

	switch {
	case role.ServiceAccountName != "":
		// Create token for existing service account
		status, err := client.createToken(ctx, reqPayload.Namespace, role.ServiceAccountName, theTTL, theAudiences)
		if err != nil {
			return nil, fmt.Errorf("failed to create a service account token for %s/%s: %s", reqPayload.Namespace, role.ServiceAccountName, err)
		}
		serviceAccountName = role.ServiceAccountName
		token = status.Token
	case role.K8sRoleName != "":
		// Create rolebinding for existing role
		// Create service account for existing role
		// then token
		// RoleBinding/ClusterRoleBinding will be the owning object
		ownerRef := metav1.OwnerReference{}
		walID, ownerRef, err = createRoleBindingWithWAL(ctx, client, req.Storage, reqPayload.Namespace, genName, role.K8sRoleName, reqPayload.ClusterRoleBinding, role)
		if err != nil {
			return nil, err
		}

		err = createServiceAccount(ctx, client, reqPayload.Namespace, genName, role, ownerRef)
		if err != nil {
			return nil, err
		}

		status, err := client.createToken(ctx, reqPayload.Namespace, genName, theTTL, theAudiences)
		if err != nil {
			return nil, fmt.Errorf("failed to create a service account token for %s/%s: %s", reqPayload.Namespace, genName, err)
		}
		token = status.Token
		serviceAccountName = genName
		createdServiceAccountName = genName
		createdK8sRoleBinding = genName
	case role.RoleRules != "":
		// Create role, rolebinding, service account, token
		// Role/ClusterRole will be the owning object
		ownerRef := metav1.OwnerReference{}
		walID, ownerRef, err = createRoleWithWAL(ctx, client, req.Storage, reqPayload.Namespace, genName, role)
		if err != nil {
			return nil, err
		}

		err = createRoleBinding(ctx, client, reqPayload.Namespace, genName, genName, reqPayload.ClusterRoleBinding, role, ownerRef)
		if err != nil {
			return nil, err
		}

		err = createServiceAccount(ctx, client, reqPayload.Namespace, genName, role, ownerRef)
		if err != nil {
			return nil, err
		}

		status, err := client.createToken(ctx, reqPayload.Namespace, genName, theTTL, theAudiences)
		if err != nil {
			return nil, fmt.Errorf("failed to create a service account token for %s/%s: %s", reqPayload.Namespace, genName, err)
		}
		token = status.Token
		createdK8sRole = genName
		serviceAccountName = genName
		createdServiceAccountName = genName
		createdK8sRoleBinding = genName

	default:
		return nil, fmt.Errorf("one of service_account_name, kubernetes_role_name, or generated_role_rules must be set")
	}

	resp := b.Secret(kubeTokenType).Response(map[string]interface{}{
		"service_account_namespace": reqPayload.Namespace,
		"service_account_name":      serviceAccountName,
		"service_account_token":     token,
	}, map[string]interface{}{
		// the internal data is whatever we need to cleanup on revoke
		// (service_account_name, role, role_binding).
		"role":                      reqPayload.RoleName,
		"service_account_namespace": reqPayload.Namespace,
		"cluster_role_binding":      reqPayload.ClusterRoleBinding,
		"created_service_account":   createdServiceAccountName,
		"created_role_binding":      createdK8sRoleBinding,
		"created_role":              createdK8sRole,
		"created_role_type":         role.K8sRoleType,
	})

	resp.Secret.TTL = theTTL
	if role.TokenMaxTTL > 0 {
		resp.Secret.MaxTTL = role.TokenMaxTTL
	}

	createdTokenTTL, err := getTokenTTL(token)
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to read TTL of created Kubernetes token for %s/%s: %s", reqPayload.Namespace, genName, err)
	case createdTokenTTL > theTTL:
		respWarning = append(respWarning, fmt.Sprintf("the created Kubernetes service accout token TTL %v is greater than the Vault lease TTL %v", createdTokenTTL, theTTL))
	case createdTokenTTL < theTTL:
		respWarning = append(respWarning, fmt.Sprintf("the created Kubernetes service accout token TTL %v is less than the Vault lease TTL %v; capping the lease TTL accordingly", createdTokenTTL, theTTL))
		resp.Secret.TTL = createdTokenTTL
	}

	if len(respWarning) > 0 {
		resp.Warnings = respWarning
	}

	// Delete the WAL entry that was created, since all the k8s objects were
	// created successfully (no need to rollback anymore)
	if walID != "" {
		if err := framework.DeleteWAL(ctx, req.Storage, walID); err != nil {
			return nil, fmt.Errorf("error deleting WAL: %w", err)
		}
	}

	return resp, nil
}

func (b *backend) getClient(ctx context.Context, s logical.Storage) (*client, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	client := b.client
	if client != nil {
		return client, nil
	}

	config, err := b.configWithDynamicValues(ctx, s)
	if err != nil {
		return nil, err
	}

	if b.client == nil && config == nil {
		config = new(kubeConfig)
	}

	b.client, err = newClient(config)
	if err != nil {
		return nil, err
	}

	return b.client, nil
}

// create service account
func createServiceAccount(ctx context.Context, client *client, namespace, name string, vaultRole *roleEntry, ownerRef metav1.OwnerReference) error {
	_, err := client.createServiceAccount(ctx, namespace, name, vaultRole, ownerRef)
	if err != nil {
		return fmt.Errorf("failed to create service account '%s/%s': %s", namespace, name, err)
	}

	return nil
}

// create role binding and put a WAL entry
func createRoleBindingWithWAL(ctx context.Context, client *client, s logical.Storage, namespace, name, k8sRoleName string, isClusterRoleBinding bool, vaultRole *roleEntry) (string, metav1.OwnerReference, error) {
	// Write a WAL entry in case the role binding create doesn't complete
	walId, err := framework.PutWAL(ctx, s, walBindingKind, &walRoleBinding{
		Namespace:  namespace,
		Name:       name,
		IsCluster:  isClusterRoleBinding,
		Expiration: time.Now().Add(maxWALAge),
	})
	if err != nil {
		return "", metav1.OwnerReference{}, fmt.Errorf("error writing role binding WAL: %w", err)
	}

	ownerRef, err := client.createRoleBinding(ctx, namespace, name, k8sRoleName, isClusterRoleBinding, vaultRole, nil)
	if err != nil {
		return "", ownerRef, fmt.Errorf("failed to create RoleBinding/ClusterRoleBinding '%s' for %s: %s", name, k8sRoleName, err)
	}

	return walId, ownerRef, nil
}

func createRoleBinding(ctx context.Context, client *client, namespace, name, k8sRoleName string, isClusterRoleBinding bool, vaultRole *roleEntry, ownerRef metav1.OwnerReference) error {
	_, err := client.createRoleBinding(ctx, namespace, name, k8sRoleName, isClusterRoleBinding, vaultRole, &ownerRef)
	if err != nil {
		return fmt.Errorf("failed to create RoleBinding/ClusterRoleBinding '%s' for %s: %s", name, k8sRoleName, err)
	}
	return nil
}

// create a role and put a WAL entry
func createRoleWithWAL(ctx context.Context, client *client, s logical.Storage, namespace, name string, vaultRole *roleEntry) (string, metav1.OwnerReference, error) {
	// Write a WAL entry in case subsequent parts don't complete
	walId, err := framework.PutWAL(ctx, s, walRoleKind, &walRole{
		Namespace:  namespace,
		Name:       name,
		RoleType:   vaultRole.K8sRoleType,
		Expiration: time.Now().Add(maxWALAge),
	})
	if err != nil {
		return "", metav1.OwnerReference{}, fmt.Errorf("error writing service account WAL: %w", err)
	}

	ownerRef, err := client.createRole(ctx, namespace, name, vaultRole)
	if err != nil {
		return "", ownerRef, fmt.Errorf("failed to create Role/ClusterRole '%s/%s: %s", namespace, name, err)
	}

	return walId, ownerRef, nil
}

func getTokenTTL(token string) (time.Duration, error) {
	parsed, err := josejwt.ParseSigned(token, AllowedSigningAlgs)
	if err != nil {
		return 0, err
	}
	claims := map[string]interface{}{}
	err = parsed.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return 0, err
	}
	sa := struct {
		Expiration int64 `mapstructure:"exp"`
		IssuedAt   int64 `mapstructure:"iat"`
	}{}
	err = mapstructure.Decode(claims, &sa)
	if err != nil {
		return 0, err
	}
	return time.Duration(sa.Expiration-sa.IssuedAt) * time.Second, nil
}

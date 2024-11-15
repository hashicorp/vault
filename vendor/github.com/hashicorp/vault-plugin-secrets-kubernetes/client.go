// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var standardLabels = map[string]string{
	"app.kubernetes.io/managed-by": "HashiCorp-Vault",
	"app.kubernetes.io/created-by": "vault-plugin-secrets-kubernetes",
}

type client struct {
	k8s kubernetes.Interface
}

func newClient(config *kubeConfig) (*client, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	clientConfig := rest.Config{
		Host:        config.Host,
		BearerToken: config.ServiceAccountJwt,
	}
	if config.CACert != "" {
		clientConfig.TLSClientConfig.CAData = []byte(config.CACert)
	}
	k8sClient, err := kubernetes.NewForConfig(&clientConfig)
	if err != nil {
		return nil, err
	}
	return &client{k8sClient}, nil
}

func (c *client) createToken(ctx context.Context, namespace, name string, ttl time.Duration, audiences []string) (*authenticationv1.TokenRequestStatus, error) {
	intTTL := int64(ttl.Seconds())
	resp, err := c.k8s.CoreV1().ServiceAccounts(namespace).CreateToken(ctx, name, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &intTTL,
			Audiences:         audiences,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	c.k8s.CoreV1().ServiceAccounts(namespace)
	return &resp.Status, nil
}

func (c *client) createServiceAccount(ctx context.Context, namespace, name string, vaultRole *roleEntry, ownerRef metav1.OwnerReference) (*v1.ServiceAccount, error) {
	// Set standardLabels last so that users can't override them
	labels := combineMaps(vaultRole.ExtraLabels, standardLabels)
	serviceAccountConfig := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			Annotations:     vaultRole.ExtraAnnotations,
			OwnerReferences: []metav1.OwnerReference{ownerRef},
		},
	}
	return c.k8s.CoreV1().ServiceAccounts(namespace).Create(ctx, serviceAccountConfig, metav1.CreateOptions{})
}

func (c *client) deleteServiceAccount(ctx context.Context, namespace, name string) error {
	err := c.k8s.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !k8s_errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *client) createRole(ctx context.Context, namespace, name string, vaultRole *roleEntry) (metav1.OwnerReference, error) {
	thisOwnerRef := metav1.OwnerReference{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Name:       name,
	}
	roleRules, err := makeRules(vaultRole.RoleRules)
	if err != nil {
		return thisOwnerRef, err
	}
	// Set standardLabels last so that users can't override them
	labels := combineMaps(vaultRole.ExtraLabels, standardLabels)
	objectMeta := metav1.ObjectMeta{
		Name:        name,
		Labels:      labels,
		Annotations: vaultRole.ExtraAnnotations,
	}

	switch vaultRole.K8sRoleType {
	case "Role":
		objectMeta.Namespace = namespace
		roleConfig := &rbacv1.Role{
			ObjectMeta: objectMeta,
			Rules:      roleRules,
		}
		resp, err := c.k8s.RbacV1().Roles(namespace).Create(ctx, roleConfig, metav1.CreateOptions{})
		if resp != nil {
			thisOwnerRef.Kind = "Role"
			thisOwnerRef.UID = resp.UID
		}
		return thisOwnerRef, err

	case "ClusterRole":
		roleConfig := &rbacv1.ClusterRole{
			ObjectMeta: objectMeta,
			Rules:      roleRules,
		}
		resp, err := c.k8s.RbacV1().ClusterRoles().Create(ctx, roleConfig, metav1.CreateOptions{})
		if resp != nil {
			thisOwnerRef.Kind = "ClusterRole"
			thisOwnerRef.UID = resp.UID
		}
		return thisOwnerRef, err

	default:
		return thisOwnerRef, fmt.Errorf("unknown role type '%s'", vaultRole.K8sRoleType)
	}
}

func (c *client) deleteRole(ctx context.Context, namespace, name, roleType string) error {
	var err error
	switch roleType {
	case "Role":
		err = c.k8s.RbacV1().Roles(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	case "ClusterRole":
		err = c.k8s.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	default:
		return fmt.Errorf("unsupported role type '%s'", roleType)
	}
	if err != nil && !k8s_errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *client) createRoleBinding(ctx context.Context, namespace, name, k8sRoleName string, isClusterRoleBinding bool, vaultRole *roleEntry, ownerRef *metav1.OwnerReference) (metav1.OwnerReference, error) {
	thisOwnerRef := metav1.OwnerReference{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Name:       name,
	}
	// Set standardLabels last so that users can't override them
	labels := combineMaps(vaultRole.ExtraLabels, standardLabels)
	objectMeta := metav1.ObjectMeta{
		Name:        name,
		Labels:      labels,
		Annotations: vaultRole.ExtraAnnotations,
	}
	if ownerRef != nil {
		objectMeta.OwnerReferences = []metav1.OwnerReference{*ownerRef}
	}
	subjects := []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      name,
			Namespace: namespace,
		},
	}
	roleRef := rbacv1.RoleRef{
		Kind: vaultRole.K8sRoleType,
		Name: k8sRoleName,
	}

	if isClusterRoleBinding {
		roleConfig := &rbacv1.ClusterRoleBinding{
			ObjectMeta: objectMeta,
			Subjects:   subjects,
			RoleRef:    roleRef,
		}
		resp, err := c.k8s.RbacV1().ClusterRoleBindings().Create(ctx, roleConfig, metav1.CreateOptions{})
		if resp != nil {
			thisOwnerRef.Kind = "ClusterRoleBinding"
			thisOwnerRef.UID = resp.UID
		}
		return thisOwnerRef, err
	}

	objectMeta.Namespace = namespace
	roleConfig := &rbacv1.RoleBinding{
		ObjectMeta: objectMeta,
		Subjects:   subjects,
		RoleRef:    roleRef,
	}
	resp, err := c.k8s.RbacV1().RoleBindings(namespace).Create(ctx, roleConfig, metav1.CreateOptions{})
	if resp != nil {
		thisOwnerRef.Kind = "RoleBinding"
		thisOwnerRef.UID = resp.UID
	}
	return thisOwnerRef, err
}

func (c *client) deleteRoleBinding(ctx context.Context, namespace, name string, isClusterRoleBinding bool) error {
	var err error
	if isClusterRoleBinding {
		err = c.k8s.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
	} else {
		err = c.k8s.RbacV1().RoleBindings(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	}
	if err != nil && !k8s_errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *client) getNamespaceLabelSet(ctx context.Context, namespace string) (map[string]string, error) {
	ns, err := c.k8s.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return map[string]string{}, err
	}
	return ns.Labels, nil
}

func makeRules(rules string) ([]rbacv1.PolicyRule, error) {
	policyRules := struct {
		Rules []rbacv1.PolicyRule `json:"rules"`
	}{}
	decoder := k8s_yaml.NewYAMLOrJSONDecoder(strings.NewReader(rules), len(rules))
	err := decoder.Decode(&policyRules)
	if err != nil {
		return nil, err
	}
	return policyRules.Rules, nil
}

func makeLabelSelector(selector string) (metav1.LabelSelector, error) {
	labelSelector := metav1.LabelSelector{}
	decoder := k8s_yaml.NewYAMLOrJSONDecoder(strings.NewReader(selector), len(selector))
	err := decoder.Decode(&labelSelector)
	if err != nil {
		return labelSelector, err
	}
	return labelSelector, nil
}

func makeRoleType(roleType string) string {
	switch strings.ToLower(roleType) {
	case "role":
		return "Role"
	case "clusterrole":
		return "ClusterRole"
	default:
		return roleType
	}
}

func combineMaps(maps ...map[string]string) map[string]string {
	newMap := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			newMap[k] = v
		}
	}
	return newMap
}

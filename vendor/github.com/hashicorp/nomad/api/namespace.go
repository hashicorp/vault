// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"sort"
)

// Namespaces is used to query the namespace endpoints.
type Namespaces struct {
	client *Client
}

// Namespaces returns a new handle on the namespaces.
func (c *Client) Namespaces() *Namespaces {
	return &Namespaces{client: c}
}

// List is used to dump all of the namespaces.
func (n *Namespaces) List(q *QueryOptions) ([]*Namespace, *QueryMeta, error) {
	var resp []*Namespace
	qm, err := n.client.query("/v1/namespaces", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(NamespaceIndexSort(resp))
	return resp, qm, nil
}

// PrefixList is used to do a PrefixList search over namespaces
func (n *Namespaces) PrefixList(prefix string, q *QueryOptions) ([]*Namespace, *QueryMeta, error) {
	if q == nil {
		q = &QueryOptions{Prefix: prefix}
	} else {
		q.Prefix = prefix
	}

	return n.List(q)
}

// Info is used to query a single namespace by its name.
func (n *Namespaces) Info(name string, q *QueryOptions) (*Namespace, *QueryMeta, error) {
	var resp Namespace
	qm, err := n.client.query("/v1/namespace/"+name, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Register is used to register a namespace.
func (n *Namespaces) Register(namespace *Namespace, q *WriteOptions) (*WriteMeta, error) {
	wm, err := n.client.put("/v1/namespace", namespace, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Delete is used to delete a namespace
func (n *Namespaces) Delete(namespace string, q *WriteOptions) (*WriteMeta, error) {
	wm, err := n.client.delete(fmt.Sprintf("/v1/namespace/%s", namespace), nil, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Namespace is used to serialize a namespace.
type Namespace struct {
	Name                  string
	Description           string
	Quota                 string
	Capabilities          *NamespaceCapabilities          `hcl:"capabilities,block"`
	NodePoolConfiguration *NamespaceNodePoolConfiguration `hcl:"node_pool_config,block"`
	VaultConfiguration    *NamespaceVaultConfiguration    `hcl:"vault,block"`
	ConsulConfiguration   *NamespaceConsulConfiguration   `hcl:"consul,block"`
	Meta                  map[string]string
	CreateIndex           uint64
	ModifyIndex           uint64
}

// NamespaceCapabilities represents a set of capabilities allowed for this
// namespace, to be checked at job submission time.
type NamespaceCapabilities struct {
	EnabledTaskDrivers   []string `hcl:"enabled_task_drivers"`
	DisabledTaskDrivers  []string `hcl:"disabled_task_drivers"`
	EnabledNetworkModes  []string `hcl:"enabled_network_modes"`
	DisabledNetworkModes []string `hcl:"disabled_network_modes"`
}

// NamespaceNodePoolConfiguration stores configuration about node pools for a
// namespace.
type NamespaceNodePoolConfiguration struct {
	Default string
	Allowed []string
	Denied  []string
}

// NamespaceVaultConfiguration stores configuration about permissions to Vault
// clusters for a namespace, for use with Nomad Enterprise.
type NamespaceVaultConfiguration struct {
	// Default is the Vault cluster used by jobs in this namespace that don't
	// specify a cluster of their own.
	Default string

	// Allowed specifies the Vault clusters that are allowed to be used by jobs
	// in this namespace. By default, all clusters are allowed. If an empty list
	// is provided only the namespace's default cluster is allowed. This field
	// supports wildcard globbing through the use of `*` for multi-character
	// matching. This field cannot be used with Denied.
	Allowed []string

	// Denied specifies the Vault clusters that are not allowed to be used by
	// jobs in this namespace. This field supports wildcard globbing through the
	// use of `*` for multi-character matching. If specified, any cluster is
	// allowed to be used, except for those that match any of these patterns.
	// This field cannot be used with Allowed.
	Denied []string
}

// NamespaceConsulConfiguration stores configuration about permissions to Consul
// clusters for a namespace, for use with Nomad Enterprise.
type NamespaceConsulConfiguration struct {
	// Default is the Consul cluster used by jobs in this namespace that don't
	// specify a cluster of their own.
	Default string

	// Allowed specifies the Consul clusters that are allowed to be used by jobs
	// in this namespace. By default, all clusters are allowed. If an empty list
	// is provided only the namespace's default cluster is allowed. This field
	// supports wildcard globbing through the use of `*` for multi-character
	// matching. This field cannot be used with Denied.
	Allowed []string

	// Denied specifies the Consul clusters that are not allowed to be used by
	// jobs in this namespace. This field supports wildcard globbing through the
	// use of `*` for multi-character matching. If specified, any cluster is
	// allowed to be used, except for those that match any of these patterns.
	// This field cannot be used with Allowed.
	Denied []string
}

// NamespaceIndexSort is a wrapper to sort Namespaces by CreateIndex. We
// reverse the test so that we get the highest index first.
type NamespaceIndexSort []*Namespace

func (n NamespaceIndexSort) Len() int {
	return len(n)
}

func (n NamespaceIndexSort) Less(i, j int) bool {
	return n[i].CreateIndex > n[j].CreateIndex
}

func (n NamespaceIndexSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// NamespacedID is used for things that are unique only per-namespace,
// such as jobs.
type NamespacedID struct {
	// Namespace is the Name of the Namespace
	Namespace string
	// ID is the ID of the namespaced object (e.g. Job ID)
	ID string
}

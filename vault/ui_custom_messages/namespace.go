// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import "github.com/hashicorp/vault/helper/namespace"

// NamespaceManager is the interface needed of a NamespaceManager by this
// package. This interface allows setting a dummy NamespaceManager in the
// community edition that can be replaced with the real
// namespace.NamespaceManager in the enterprise edition.
type NamespaceManager interface {
	GetParentNamespace(string) *namespace.Namespace
}

// CommunityEditionNamespaceManager is a struct that implements the
// NamespaceManager interface. This struct is used as a placeholder in the
// community edition.
type CommunityEditionNamespaceManager struct{}

// GetParentNamespace always returns namespace.RootNamespace.
func (n *CommunityEditionNamespaceManager) GetParentNamespace(_ string) *namespace.Namespace {
	return namespace.RootNamespace
}

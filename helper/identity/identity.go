// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/sdk/logical"
)

func (g *Group) Clone() (*Group, error) {
	if g == nil {
		return nil, fmt.Errorf("nil group")
	}

	marshaledGroup, err := proto.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group: %w", err)
	}

	var clonedGroup Group
	err = proto.Unmarshal(marshaledGroup, &clonedGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group: %w", err)
	}

	return &clonedGroup, nil
}

func (e *Entity) Clone() (*Entity, error) {
	if e == nil {
		return nil, fmt.Errorf("nil entity")
	}

	marshaledEntity, err := proto.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity: %w", err)
	}

	var clonedEntity Entity
	err = proto.Unmarshal(marshaledEntity, &clonedEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity: %w", err)
	}

	return &clonedEntity, nil
}

func (e *Entity) UpsertAlias(alias *Alias) {
	for i, item := range e.Aliases {
		if item.ID == alias.ID {
			e.Aliases[i] = alias
			return
		}
	}
	e.Aliases = append(e.Aliases, alias)
}

func (e *Entity) DeleteAliasByID(aliasID string) {
	idx := -1
	for i, item := range e.Aliases {
		if item.ID == aliasID {
			idx = i
			break
		}
	}

	if idx < 0 {
		return
	}

	e.Aliases = append(e.Aliases[:idx], e.Aliases[idx+1:]...)
}

func (p *Alias) Clone() (*Alias, error) {
	if p == nil {
		return nil, fmt.Errorf("nil alias")
	}

	marshaledAlias, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal alias: %w", err)
	}

	var clonedAlias Alias
	err = proto.Unmarshal(marshaledAlias, &clonedAlias)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal alias: %w", err)
	}

	return &clonedAlias, nil
}

// ToSDKAlias converts the provided alias to an SDK compatible alias.
func ToSDKAlias(a *Alias) *logical.Alias {
	if a == nil {
		return nil
	}
	metadata := make(map[string]string, len(a.Metadata))
	for k, v := range a.Metadata {
		metadata[k] = v
	}

	return &logical.Alias{
		Name:           a.Name,
		ID:             a.ID,
		MountAccessor:  a.MountAccessor,
		MountType:      a.MountType,
		Metadata:       metadata,
		NamespaceID:    a.NamespaceID,
		CustomMetadata: a.CustomMetadata,
	}
}

// ToSDKEntity converts the provided entity to an SDK compatible entity.
func ToSDKEntity(e *Entity) *logical.Entity {
	if e == nil {
		return nil
	}

	aliases := make([]*logical.Alias, len(e.Aliases))

	for i, a := range e.Aliases {
		aliases[i] = ToSDKAlias(a)
	}

	metadata := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		metadata[k] = v
	}

	return &logical.Entity{
		ID:          e.ID,
		Name:        e.Name,
		Disabled:    e.Disabled,
		Aliases:     aliases,
		Metadata:    metadata,
		NamespaceID: e.NamespaceID,
	}
}

// ToSDKGroup converts the provided group to an SDK compatible group.
func ToSDKGroup(g *Group) *logical.Group {
	if g == nil {
		return nil
	}

	metadata := make(map[string]string, len(g.Metadata))
	for k, v := range g.Metadata {
		metadata[k] = v
	}

	return &logical.Group{
		ID:          g.ID,
		Name:        g.Name,
		Metadata:    metadata,
		NamespaceID: g.NamespaceID,
	}
}

// ToSDKGroups converts the provided group list to an SDK compatible group list.
func ToSDKGroups(groups []*Group) []*logical.Group {
	if groups == nil {
		return nil
	}

	ret := make([]*logical.Group, len(groups))

	for i, g := range groups {
		ret[i] = ToSDKGroup(g)
	}
	return ret
}

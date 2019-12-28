package identity

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
)

func (g *Group) Clone() (*Group, error) {
	if g == nil {
		return nil, fmt.Errorf("nil group")
	}

	marshaledGroup, err := proto.Marshal(g)
	if err != nil {
		return nil, errwrap.Wrapf("failed to marshal group: {{err}}", err)
	}

	var clonedGroup Group
	err = proto.Unmarshal(marshaledGroup, &clonedGroup)
	if err != nil {
		return nil, errwrap.Wrapf("failed to unmarshal group: {{err}}", err)
	}

	return &clonedGroup, nil
}

func (e *Entity) Clone() (*Entity, error) {
	if e == nil {
		return nil, fmt.Errorf("nil entity")
	}

	marshaledEntity, err := proto.Marshal(e)
	if err != nil {
		return nil, errwrap.Wrapf("failed to marshal entity: {{err}}", err)
	}

	var clonedEntity Entity
	err = proto.Unmarshal(marshaledEntity, &clonedEntity)
	if err != nil {
		return nil, errwrap.Wrapf("failed to unmarshal entity: {{err}}", err)
	}

	return &clonedEntity, nil
}

func (p *Alias) Clone() (*Alias, error) {
	if p == nil {
		return nil, fmt.Errorf("nil alias")
	}

	marshaledAlias, err := proto.Marshal(p)
	if err != nil {
		return nil, errwrap.Wrapf("failed to marshal alias: {{err}}", err)
	}

	var clonedAlias Alias
	err = proto.Unmarshal(marshaledAlias, &clonedAlias)
	if err != nil {
		return nil, errwrap.Wrapf("failed to unmarshal alias: {{err}}", err)
	}

	return &clonedAlias, nil
}

func FromLogicalAlias(a *logical.Alias) *Alias {
	metadata := make(map[string]string, len(a.Metadata))
	for k, v := range a.Metadata {
		metadata[k] = v
	}

	return &Alias{
		Name:          a.Name,
		ID:            a.ID,
		MountAccessor: a.MountAccessor,
		MountType:     a.MountType,
		Metadata:      metadata,
	}
}

func FromLogicalEntity(e *logical.Entity) *Entity {
	aliases := make([]*Alias, len(e.Aliases))

	for i, a := range e.Aliases {
		aliases[i] = FromLogicalAlias(a)
	}

	metadata := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		metadata[k] = v
	}

	return &Entity{
		ID:       e.ID,
		Name:     e.Name,
		Disabled: e.Disabled,
		Aliases:  aliases,
		Metadata: metadata,
	}
}

func FromLogicalGroup(g *logical.Group) *Group {
	metadata := make(map[string]string, len(g.Metadata))
	for k, v := range g.Metadata {
		metadata[k] = v
	}

	return &Group{
		ID:       g.ID,
		Name:     g.Name,
		Metadata: metadata,
	}
}

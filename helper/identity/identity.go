package identity

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
)

func (g *Group) Clone() (*Group, error) {
	if g == nil {
		return nil, fmt.Errorf("nil group")
	}

	marshaledGroup, err := proto.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group: %v", err)
	}

	var clonedGroup Group
	err = proto.Unmarshal(marshaledGroup, &clonedGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group: %v", err)
	}

	return &clonedGroup, nil
}

func (e *Entity) Clone() (*Entity, error) {
	if e == nil {
		return nil, fmt.Errorf("nil entity")
	}

	marshaledEntity, err := proto.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity: %v", err)
	}

	var clonedEntity Entity
	err = proto.Unmarshal(marshaledEntity, &clonedEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity: %v", err)
	}

	return &clonedEntity, nil
}

func (p *Alias) Clone() (*Alias, error) {
	if p == nil {
		return nil, fmt.Errorf("nil alias")
	}

	marshaledAlias, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal alias: %v", err)
	}

	var clonedAlias Alias
	err = proto.Unmarshal(marshaledAlias, &clonedAlias)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal alias: %v", err)
	}

	return &clonedAlias, nil
}

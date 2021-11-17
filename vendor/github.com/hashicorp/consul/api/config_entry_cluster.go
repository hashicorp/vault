package api

import "encoding/json"

type MeshConfigEntry struct {
	Namespace        string                     `json:",omitempty"`
	TransparentProxy TransparentProxyMeshConfig `alias:"transparent_proxy"`
	Meta             map[string]string          `json:",omitempty"`
	CreateIndex      uint64
	ModifyIndex      uint64
}

type TransparentProxyMeshConfig struct {
	MeshDestinationsOnly bool `alias:"mesh_destinations_only"`
}

func (e *MeshConfigEntry) GetKind() string {
	return MeshConfig
}

func (e *MeshConfigEntry) GetName() string {
	return MeshConfigMesh
}

func (e *MeshConfigEntry) GetNamespace() string {
	return e.Namespace
}

func (e *MeshConfigEntry) GetMeta() map[string]string {
	return e.Meta
}

func (e *MeshConfigEntry) GetCreateIndex() uint64 {
	return e.CreateIndex
}

func (e *MeshConfigEntry) GetModifyIndex() uint64 {
	return e.ModifyIndex
}

// MarshalJSON adds the Kind field so that the JSON can be decoded back into the
// correct type.
func (e *MeshConfigEntry) MarshalJSON() ([]byte, error) {
	type Alias MeshConfigEntry
	source := &struct {
		Kind string
		*Alias
	}{
		Kind:  MeshConfig,
		Alias: (*Alias)(e),
	}
	return json.Marshal(source)
}

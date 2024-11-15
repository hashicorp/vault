// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/nomad/api"
)

// NewNomadVarMeta is used to create a NomadVarMeta from a Nomad API
// VariableMetadata response.
func NewNomadVarMeta(in *api.VariableMetadata) *NomadVarMeta {
	return &NomadVarMeta{
		Namespace:   in.Namespace,
		Path:        in.Path,
		CreateIndex: in.CreateIndex,
		ModifyIndex: in.ModifyIndex,
		CreateTime:  nanoTime(in.CreateTime),
		ModifyTime:  nanoTime(in.ModifyTime),
	}
}

// NewNomadVariable is used to create a NomadVariable from a Nomad
// API Variable response.
func NewNomadVariable(in *api.Variable) *NomadVariable {
	out := NomadVariable{
		Namespace:   in.Namespace,
		Path:        in.Path,
		CreateIndex: in.CreateIndex,
		ModifyIndex: in.ModifyIndex,
		CreateTime:  nanoTime(in.CreateTime),
		ModifyTime:  nanoTime(in.ModifyTime),
		Items:       map[string]NomadVarItem{},
		nVar:        in,
	}

	items := make(NomadVarItems, len(in.Items))
	for k, v := range in.Items {
		items[k] = NomadVarItem{k, v, &out}
	}
	out.Items = items
	return &out
}

// NomadVariable is a template friendly container struct that allows for
// the NomadVar funcs to start inside of Items and have a rational way back up
// to the Variable that is JSON structurally equivalent to the API response.
// This struct's zero value is not trivially usable and should be created with
// NewNomadVariable--especially when outside of the dependency package as
// there is no access to nVar.
type NomadVariable struct {
	Namespace, Path          string
	CreateIndex, ModifyIndex uint64
	CreateTime, ModifyTime   nanoTime
	Items                    NomadVarItems
	nVar                     *api.Variable
}

func (cv NomadVariable) Metadata() *NomadVarMeta {
	return NewNomadVarMeta(cv.nVar.Metadata())
}

type NomadVarItems map[string]NomadVarItem

// Keys returns a sorted list of the Item map's keys.
func (v NomadVarItems) Keys() []string {
	out := make([]string, 0, len(v))
	for k := range v {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Values produces a key-sorted list of the Items map's values
func (v NomadVarItems) Values() []string {
	out := make([]string, 0, len(v))
	for _, k := range v.Keys() {
		out = append(out, v[k].String())
	}
	return out
}

// Tuples produces a key-sorted list of K,V tuple structs from the Items map's
// values
func (v NomadVarItems) Tuples() []struct{ K, V string } {
	out := make([]struct{ K, V string }, 0, len(v))
	for _, k := range v.Keys() {
		out = append(out, struct{ K, V string }{K: k, V: v[k].String()})
	}
	return out
}

// Metadata returns this item's parent's metadata
func (i NomadVarItems) Metadata() *NomadVarMeta {
	for _, v := range i {
		return v.parent.Metadata()
	}
	return nil
}

// Parent returns the item's container object
func (i NomadVarItems) Parent() *NomadVariable {
	for _, v := range i {
		return v.Parent()
	}
	return nil
}

func (i NomadVarItems) ItemsMap() map[string]interface{} {
	if len(i) == 0 {
		return nil
	}
	out := make(map[string]interface{})
	for k, v := range i {
		out[k] = v.String()
	}
	return out
}

// NomadVarItem enriches the basic string values in a api.Variable's Items
// map with additional helper funcs for formatting and access to its parent
// item. This enables us to have the template funcs start at the Items
// collection without the user having to delve to it themselves and to minimize
// the number of template funcs that we have to provide for coverage.
type NomadVarItem struct {
	Key, Value string
	parent     *NomadVariable
}

func (v NomadVarItem) String() string               { return v.Value }
func (v NomadVarItem) Metadata() *NomadVarMeta      { return v.parent.Metadata() }
func (v NomadVarItem) Parent() *NomadVariable       { return v.parent }
func (v NomadVarItem) MarshalJSON() ([]byte, error) { return []byte(fmt.Sprintf("%q", v.Value)), nil }

// NomadVarMeta provides the same fields as api.VariableMetadata
// but aliases the times into a more template friendly alternative.
type NomadVarMeta struct {
	Namespace, Path          string
	CreateIndex, ModifyIndex uint64
	CreateTime, ModifyTime   nanoTime
}

func (s NomadVarMeta) String() string { return s.Path }

// nanoTime is the typical storage encoding for times in Nomad's backend. They
// are not pretty for consul-template consumption, so this gives us a type to
// add receivers on.
type nanoTime int64

func (t nanoTime) String() string  { return fmt.Sprintf("%v", time.Unix(0, int64(t))) }
func (t nanoTime) Time() time.Time { return time.Unix(0, int64(t)) }

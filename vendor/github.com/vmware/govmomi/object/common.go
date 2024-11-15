/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package object

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	ErrNotSupported = errors.New("product/version specific feature not supported by target")
)

// Common contains the fields and functions common to all objects.
type Common struct {
	InventoryPath string

	c *vim25.Client
	r types.ManagedObjectReference
}

func (c Common) String() string {
	ref := fmt.Sprintf("%v", c.Reference())

	if c.InventoryPath == "" {
		return ref
	}

	return fmt.Sprintf("%s @ %s", ref, c.InventoryPath)
}

func NewCommon(c *vim25.Client, r types.ManagedObjectReference) Common {
	return Common{c: c, r: r}
}

func (c Common) Reference() types.ManagedObjectReference {
	return c.r
}

func (c Common) Client() *vim25.Client {
	return c.c
}

// Name returns the base name of the InventoryPath field
func (c Common) Name() string {
	if c.InventoryPath == "" {
		return ""
	}
	return path.Base(c.InventoryPath)
}

func (c *Common) SetInventoryPath(p string) {
	c.InventoryPath = p
}

// ObjectName fetches the mo.ManagedEntity.Name field via the property collector.
func (c Common) ObjectName(ctx context.Context) (string, error) {
	var content []types.ObjectContent

	err := c.Properties(ctx, c.Reference(), []string{"name"}, &content)
	if err != nil {
		return "", err
	}

	for i := range content {
		for _, prop := range content[i].PropSet {
			return prop.Val.(string), nil
		}
	}

	return "", nil
}

// Properties is a wrapper for property.DefaultCollector().RetrieveOne()
func (c Common) Properties(ctx context.Context, r types.ManagedObjectReference, ps []string, dst interface{}) error {
	return property.DefaultCollector(c.c).RetrieveOne(ctx, r, ps, dst)
}

func (c Common) Destroy(ctx context.Context) (*Task, error) {
	req := types.Destroy_Task{
		This: c.Reference(),
	}

	res, err := methods.Destroy_Task(ctx, c.c, &req)
	if err != nil {
		return nil, err
	}

	return NewTask(c.c, res.Returnval), nil
}

func (c Common) Rename(ctx context.Context, name string) (*Task, error) {
	req := types.Rename_Task{
		This:    c.Reference(),
		NewName: name,
	}

	res, err := methods.Rename_Task(ctx, c.c, &req)
	if err != nil {
		return nil, err
	}

	return NewTask(c.c, res.Returnval), nil
}

func (c Common) SetCustomValue(ctx context.Context, key string, value string) error {
	req := types.SetCustomValue{
		This:  c.Reference(),
		Key:   key,
		Value: value,
	}

	_, err := methods.SetCustomValue(ctx, c.c, &req)
	return err
}

var refTypeMap = map[string]string{
	"datacenter":  "Datacenter",
	"datastore":   "Datastore",
	"domain":      "ComputeResource",
	"dvportgroup": "DistributedVirtualPortgroup",
	"dvs":         "DistributedVirtualSwitch",
	"group":       "Folder",
	"host":        "HostSystem",
	"network":     "Network",
	"resgroup":    "ResourcePool",
	"vm":          "VirtualMachine",
}

// sub types
var prefixTypeMap = map[string]struct{ prefix, kind string }{
	"domain":   {"c", "ClusterComputeResource"}, // extends ComputeResource
	"group":    {"p", "StoragePod"},             // extends Folder
	"resgroup": {"v", "VirtualApp"},             // extends ResourcePool
}

// ReferenceFromString converts a string to ManagedObjectReference.
// First checks for ManagedObjectReference (MOR), in the format of:
// "$Type:$ID", e.g. "Datacenter:datacenter-3"
// Next checks for Managed Object ID (MOID), where type is derived from the ID.
// For example, "datacenter-3" is converted to a MOR "Datacenter:datacenter-3"
// Returns nil if string is not in either format.
func ReferenceFromString(s string) *types.ManagedObjectReference {
	var ref types.ManagedObjectReference
	if ref.FromString(s) && mo.IsManagedObjectType(ref.Type) {
		return &ref
	}

	id := strings.SplitN(s, "-", 2)
	if len(id) != 2 {
		return nil
	}

	if kind, ok := refTypeMap[id[0]]; ok {
		if p, ok := prefixTypeMap[id[0]]; ok {
			if strings.HasPrefix(id[1], p.prefix) {
				return &types.ManagedObjectReference{
					Type:  p.kind,
					Value: s,
				}
			}
		}

		return &types.ManagedObjectReference{
			Type:  kind,
			Value: s,
		}
	}

	return nil
}

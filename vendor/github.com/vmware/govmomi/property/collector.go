/*
Copyright (c) 2015-2024 VMware, Inc. All Rights Reserved.

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

package property

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrConcurrentCollector is returned from WaitForUpdates, WaitForUpdatesEx,
// or CheckForUpdates if any of those calls are unable to obtain an exclusive
// lock for the property collector.
var ErrConcurrentCollector = fmt.Errorf(
	"only one goroutine may invoke WaitForUpdates, WaitForUpdatesEx, " +
		"or CheckForUpdates on a given PropertyCollector")

// Collector models the PropertyCollector managed object.
//
// For more information, see:
// http://pubs.vmware.com/vsphere-60/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvmodl.query.PropertyCollector.html
type Collector struct {
	mu           sync.Mutex
	roundTripper soap.RoundTripper
	reference    types.ManagedObjectReference
}

// DefaultCollector returns the session's default property collector.
func DefaultCollector(c *vim25.Client) *Collector {
	p := Collector{
		roundTripper: c,
		reference:    c.ServiceContent.PropertyCollector,
	}

	return &p
}

func (p *Collector) Reference() types.ManagedObjectReference {
	return p.reference
}

// Create creates a new session-specific Collector that can be used to
// retrieve property updates independent of any other Collector.
func (p *Collector) Create(ctx context.Context) (*Collector, error) {
	req := types.CreatePropertyCollector{
		This: p.Reference(),
	}

	res, err := methods.CreatePropertyCollector(ctx, p.roundTripper, &req)
	if err != nil {
		return nil, err
	}

	newp := Collector{
		roundTripper: p.roundTripper,
		reference:    res.Returnval,
	}

	return &newp, nil
}

// Destroy destroys this Collector.
func (p *Collector) Destroy(ctx context.Context) error {
	req := types.DestroyPropertyCollector{
		This: p.Reference(),
	}

	_, err := methods.DestroyPropertyCollector(ctx, p.roundTripper, &req)
	if err != nil {
		return err
	}

	p.reference = types.ManagedObjectReference{}
	return nil
}

func (p *Collector) CreateFilter(ctx context.Context, req types.CreateFilter) (*Filter, error) {
	req.This = p.Reference()

	resp, err := methods.CreateFilter(ctx, p.roundTripper, &req)
	if err != nil {
		return nil, err
	}

	return &Filter{roundTripper: p.roundTripper, reference: resp.Returnval}, nil
}

// Deprecated: Please use WaitForUpdatesEx instead.
func (p *Collector) WaitForUpdates(
	ctx context.Context,
	version string,
	opts ...*types.WaitOptions) (*types.UpdateSet, error) {

	if !p.mu.TryLock() {
		return nil, ErrConcurrentCollector
	}
	defer p.mu.Unlock()

	req := types.WaitForUpdatesEx{
		This:    p.Reference(),
		Version: version,
	}

	if len(opts) == 1 {
		req.Options = opts[0]
	} else if len(opts) > 1 {
		panic("only one option may be specified")
	}

	res, err := methods.WaitForUpdatesEx(ctx, p.roundTripper, &req)
	if err != nil {
		return nil, err
	}

	return res.Returnval, nil
}

func (p *Collector) CancelWaitForUpdates(ctx context.Context) error {
	req := &types.CancelWaitForUpdates{This: p.Reference()}
	_, err := methods.CancelWaitForUpdates(ctx, p.roundTripper, req)
	return err
}

// RetrieveProperties wraps RetrievePropertiesEx and ContinueRetrievePropertiesEx to collect properties in batches.
func (p *Collector) RetrieveProperties(
	ctx context.Context,
	req types.RetrieveProperties,
	maxObjectsArgs ...int32) (*types.RetrievePropertiesResponse, error) {

	var opts types.RetrieveOptions
	if l := len(maxObjectsArgs); l > 1 {
		return nil, fmt.Errorf("maxObjectsArgs accepts a single value")
	} else if l == 1 {
		opts.MaxObjects = maxObjectsArgs[0]
	}

	objects, err := mo.RetrievePropertiesEx(ctx, p.roundTripper, types.RetrievePropertiesEx{
		This:    p.Reference(),
		SpecSet: req.SpecSet,
		Options: opts,
	})
	if err != nil {
		return nil, err
	}

	return &types.RetrievePropertiesResponse{Returnval: objects}, nil
}

// Retrieve loads properties for a slice of managed objects. The dst argument
// must be a pointer to a []interface{}, which is populated with the instances
// of the specified managed objects, with the relevant properties filled in. If
// the properties slice is nil, all properties are loaded.
// Note that pointer types are optional fields that may be left as a nil value.
// The caller should check such fields for a nil value before dereferencing.
func (p *Collector) Retrieve(ctx context.Context, objs []types.ManagedObjectReference, ps []string, dst interface{}) error {
	if len(objs) == 0 {
		return errors.New("object references is empty")
	}

	kinds := make(map[string]bool)

	var propSet []types.PropertySpec
	var objectSet []types.ObjectSpec

	for _, obj := range objs {
		if _, ok := kinds[obj.Type]; !ok {
			spec := types.PropertySpec{
				Type: obj.Type,
			}
			if len(ps) == 0 {
				spec.All = types.NewBool(true)
			} else {
				spec.PathSet = ps
			}
			propSet = append(propSet, spec)
			kinds[obj.Type] = true
		}

		objectSpec := types.ObjectSpec{
			Obj:  obj,
			Skip: types.NewBool(false),
		}

		objectSet = append(objectSet, objectSpec)
	}

	req := types.RetrieveProperties{
		SpecSet: []types.PropertyFilterSpec{
			{
				ObjectSet: objectSet,
				PropSet:   propSet,
			},
		},
	}

	res, err := p.RetrieveProperties(ctx, req)
	if err != nil {
		return err
	}

	if d, ok := dst.(*[]types.ObjectContent); ok {
		*d = res.Returnval
		return nil
	}

	return mo.LoadObjectContent(res.Returnval, dst)
}

// RetrieveWithFilter populates dst as Retrieve does, but only for entities
// that match the specified filter.
func (p *Collector) RetrieveWithFilter(
	ctx context.Context,
	objs []types.ManagedObjectReference,
	ps []string,
	dst interface{},
	filter Match) error {

	if len(filter) == 0 {
		return p.Retrieve(ctx, objs, ps, dst)
	}

	var content []types.ObjectContent

	err := p.Retrieve(ctx, objs, filter.Keys(), &content)
	if err != nil {
		return err
	}

	objs = filter.ObjectContent(content)

	if len(objs) == 0 {
		return nil
	}

	return p.Retrieve(ctx, objs, ps, dst)
}

// RetrieveOne calls Retrieve with a single managed object reference via Collector.Retrieve().
func (p *Collector) RetrieveOne(ctx context.Context, obj types.ManagedObjectReference, ps []string, dst interface{}) error {
	var objs = []types.ManagedObjectReference{obj}
	return p.Retrieve(ctx, objs, ps, dst)
}

// WaitForUpdatesEx waits for any of the specified properties of the specified
// managed object to change. It calls the specified function for every update it
// receives. If this function returns false, it continues waiting for
// subsequent updates. If this function returns true, it stops waiting and
// returns.
//
// If the Context is canceled, a call to CancelWaitForUpdates() is made and its
// error value is returned.
//
// By default, ObjectUpdate.MissingSet faults are not propagated to the returned
// error, set WaitFilter.PropagateMissing=true to enable MissingSet fault
// propagation.
func (p *Collector) WaitForUpdatesEx(
	ctx context.Context,
	opts *WaitOptions,
	onUpdatesFn func([]types.ObjectUpdate) bool) error {

	if !p.mu.TryLock() {
		return ErrConcurrentCollector
	}
	defer p.mu.Unlock()

	req := types.WaitForUpdatesEx{
		This:    p.Reference(),
		Options: opts.Options,
	}

	for {
		res, err := methods.WaitForUpdatesEx(ctx, p.roundTripper, &req)
		if err != nil {
			if ctx.Err() == context.Canceled {
				return p.CancelWaitForUpdates(context.Background())
			}
			return err
		}

		set := res.Returnval
		if set == nil {
			if req.Options != nil && req.Options.MaxWaitSeconds != nil {
				return nil // WaitOptions.MaxWaitSeconds exceeded
			}
			// Retry if the result came back empty
			continue
		}

		req.Version = set.Version
		opts.Truncated = false
		if set.Truncated != nil {
			opts.Truncated = *set.Truncated
		}

		for _, fs := range set.FilterSet {
			if opts.PropagateMissing {
				for i := range fs.ObjectSet {
					for _, p := range fs.ObjectSet[i].MissingSet {
						// Same behavior as mo.ObjectContentToType()
						return soap.WrapVimFault(p.Fault.Fault)
					}
				}
			}

			if onUpdatesFn(fs.ObjectSet) {
				return nil
			}
		}
	}
}

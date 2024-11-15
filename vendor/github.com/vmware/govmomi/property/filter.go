/*
Copyright (c) 2017-2024 VMware, Inc. All Rights Reserved.

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

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

// Filter models the Filter managed object.
//
// For more information, see:
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/184bb3ba-6fa8-4574-a767-d0c96e2a38f4/ba9422ef-405c-47dd-8553-e11b619185b2/SDK/vsphere-ws/docs/ReferenceGuide/vmodl.query.PropertyCollector.Filter.html.
type Filter struct {
	roundTripper soap.RoundTripper
	reference    types.ManagedObjectReference
}

func (f Filter) Reference() types.ManagedObjectReference {
	return f.reference
}

// Destroy destroys this filter.
//
// This operation can be called explicitly, or it can take place implicitly when
// the session that created the filter is closed.
func (f *Filter) Destroy(ctx context.Context) error {
	if _, err := methods.DestroyPropertyFilter(
		ctx,
		f.roundTripper,
		&types.DestroyPropertyFilter{This: f.Reference()}); err != nil {

		return err
	}
	f.reference = types.ManagedObjectReference{}
	return nil
}

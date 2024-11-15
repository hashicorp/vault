/**
 * Copyright 2016-2024 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 * on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// AUTOMATICALLY GENERATED CODE - DO NOT MODIFY

package datatypes

// SoftLayer_Device_Status is used to indicate the current status of a device
type Device_Status struct {
	Entity

	// The device status's associated unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The device status's unique string identifier.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of the status.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * AUTOMATICALLY GENERATED CODE - DO NOT MODIFY
 */

package datatypes

// The SoftLayer_Service_External_Resource is a placeholder that references a service being provided outside of the standard SoftLayer system.
type Service_External_Resource struct {
	Entity

	// The customer account that is consuming the service.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The customer account that is consuming the related service.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The unique identifier in the service provider's system.
	ExternalIdentifier *string `json:"externalIdentifier,omitempty" xmlrpc:"externalIdentifier,omitempty"`

	// An external resource's unique identifier in the SoftLayer system.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`
}

// no documentation yet
type Service_Provider struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

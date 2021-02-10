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

// Throw this exception if there are validation errors. The types are specified in SoftLayer_Brand_Creation_Input including: KEY_NAME, PREFIX, NAME, LONG_NAME, SUPPORT_POLICY, POLICY_ACKNOWLEDGEMENT_FLAG, etc.
type Exception_Brand_Creation struct {
	Entity

	// no documentation yet
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// This exception is thrown if the component locator client cannot find or communicate with the component locator service.
type Exception_Hardware_Component_Locator_ComponentLocatorException struct {
	Entity
}

// This exception is thrown if the argument is of incorrect type.
type Exception_Hardware_Component_Locator_InvalidGenericComponentArgument struct {
	Entity
}

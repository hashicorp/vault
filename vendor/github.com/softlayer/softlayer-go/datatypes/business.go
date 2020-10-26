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

// Contains business partner channel information
type Business_Partner_Channel struct {
	Entity

	// Business partner channel description
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Business partner channel name
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`
}

// Contains business partner segment information
type Business_Partner_Segment struct {
	Entity

	// Business partner segment description
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Business partner segment name
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`
}

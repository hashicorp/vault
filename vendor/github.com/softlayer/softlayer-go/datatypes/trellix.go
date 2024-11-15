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

// The Trellix_Epolicy_Orchestrator_Version51_Agent_Details data type represents a virus scan agent and contains details about its version.
type Trellix_Epolicy_Orchestrator_Version51_Agent_Details struct {
	Entity

	// Version number of the anti-virus scan agent.
	AgentVersion *string `json:"agentVersion,omitempty" xmlrpc:"agentVersion,omitempty"`

	// The date of the last time the anti-virus agent checked in.
	LastUpdate *Time `json:"lastUpdate,omitempty" xmlrpc:"lastUpdate,omitempty"`
}

// The Trellix_Epolicy_Orchestrator_Version51_Policy_Object data type represents a virus scan agent and contains details about its version.
type Trellix_Epolicy_Orchestrator_Version51_Policy_Object struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The Trellix_Epolicy_Orchestrator_Version51_Product_Properties data type represents the version of the virus data file
type Trellix_Epolicy_Orchestrator_Version51_Product_Properties struct {
	Entity

	// no documentation yet
	DatVersion *string `json:"datVersion,omitempty" xmlrpc:"datVersion,omitempty"`
}

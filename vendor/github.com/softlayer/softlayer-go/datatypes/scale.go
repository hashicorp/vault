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

// no documentation yet
type Scale_Asset struct {
	Entity
}

// no documentation yet
type Scale_Asset_Hardware struct {
	Scale_Asset
}

// no documentation yet
type Scale_Asset_Virtual_Guest struct {
	Scale_Asset
}

// no documentation yet
type Scale_Group struct {
	Entity

	// The identifier of the account assigned to this group.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`
}

// no documentation yet
type Scale_LoadBalancer struct {
	Entity

	// The identifier for the health check of this load balancer configuration
	HealthCheckId *int `json:"healthCheckId,omitempty" xmlrpc:"healthCheckId,omitempty"`
}

// no documentation yet
type Scale_Member struct {
	Entity
}

// no documentation yet
type Scale_Member_Virtual_Guest struct {
	Scale_Member
}

// no documentation yet
type Scale_Network_Vlan struct {
	Entity

	// The identifier for the VLAN to scale with.
	NetworkVlanId *int `json:"networkVlanId,omitempty" xmlrpc:"networkVlanId,omitempty"`
}

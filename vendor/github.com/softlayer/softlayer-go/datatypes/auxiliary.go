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

// no documentation yet
type Auxiliary_Network_Status struct {
	Entity
}

// A SoftLayer_Auxiliary_Notification_Emergency data object represents a notification event being broadcast to the SoftLayer customer base. It is used to provide information regarding outages or current known issues.
type Auxiliary_Notification_Emergency struct {
	Entity

	// The date this event was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The device (if any) effected by this event.
	Device *string `json:"device,omitempty" xmlrpc:"device,omitempty"`

	// The duration of this event.
	Duration *string `json:"duration,omitempty" xmlrpc:"duration,omitempty"`

	// The device (if any) effected by this event.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The location effected by this event.
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// A message describing this event.
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// The last date this event was modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The service(s) (if any) effected by this event.
	ServicesAffected *string `json:"servicesAffected,omitempty" xmlrpc:"servicesAffected,omitempty"`

	// The signature of the SoftLayer employee department associated with this notification.
	Signature *Auxiliary_Notification_Emergency_Signature `json:"signature,omitempty" xmlrpc:"signature,omitempty"`

	// The date this event will start.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`

	// The status of this notification.
	Status *Auxiliary_Notification_Emergency_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Current status record for this event.
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`
}

// Every SoftLayer_Auxiliary_Notification_Emergency has a signatureId that references a SoftLayer_Auxiliary_Notification_Emergency_Signature data type.  The signature is the user or group  responsible for the current event.
type Auxiliary_Notification_Emergency_Signature struct {
	Entity

	// The name or signature for the current Emergency Notification.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// Every SoftLayer_Auxiliary_Notification_Emergency has a statusId that references a SoftLayer_Auxiliary_Notification_Emergency_Status data type.  The status is used to determine the current state of the event.
type Auxiliary_Notification_Emergency_Status struct {
	Entity

	// A name describing the status of the current Emergency Notification.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Auxiliary_Shipping_Courier data type contains general information relating the different (major) couriers that SoftLayer may use for shipping.
type Auxiliary_Shipping_Courier struct {
	Entity

	// The unique id of the shipping courier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The unique keyname of the shipping courier.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of the shipping courier.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The url to shipping courier's website.
	Url *string `json:"url,omitempty" xmlrpc:"url,omitempty"`
}

// no documentation yet
type Auxiliary_Shipping_Courier_Type struct {
	Entity

	// no documentation yet
	Courier []Auxiliary_Shipping_Courier `json:"courier,omitempty" xmlrpc:"courier,omitempty"`

	// A count of
	CourierCount *uint `json:"courierCount,omitempty" xmlrpc:"courierCount,omitempty"`

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

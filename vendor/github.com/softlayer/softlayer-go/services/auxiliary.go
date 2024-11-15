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

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// no documentation yet
type Auxiliary_Network_Status struct {
	Session session.SLSession
	Options sl.Options
}

// GetAuxiliaryNetworkStatusService returns an instance of the Auxiliary_Network_Status SoftLayer service
func GetAuxiliaryNetworkStatusService(sess session.SLSession) Auxiliary_Network_Status {
	return Auxiliary_Network_Status{Session: sess}
}

func (r Auxiliary_Network_Status) Id(id int) Auxiliary_Network_Status {
	r.Options.Id = &id
	return r
}

func (r Auxiliary_Network_Status) Mask(mask string) Auxiliary_Network_Status {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Auxiliary_Network_Status) Filter(filter string) Auxiliary_Network_Status {
	r.Options.Filter = filter
	return r
}

func (r Auxiliary_Network_Status) Limit(limit int) Auxiliary_Network_Status {
	r.Options.Limit = &limit
	return r
}

func (r Auxiliary_Network_Status) Offset(offset int) Auxiliary_Network_Status {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
// Deprecated: This function has been marked as deprecated.
func (r Auxiliary_Network_Status) GetNetworkStatus(target *string) (resp []datatypes.Container_Auxiliary_Network_Status_Reading, err error) {
	params := []interface{}{
		target,
	}
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Network_Status", "getNetworkStatus", params, &r.Options, &resp)
	return
}

// A SoftLayer_Auxiliary_Notification_Emergency data object represents a notification event being broadcast to the SoftLayer customer base. It is used to provide information regarding outages or current known issues.
type Auxiliary_Notification_Emergency struct {
	Session session.SLSession
	Options sl.Options
}

// GetAuxiliaryNotificationEmergencyService returns an instance of the Auxiliary_Notification_Emergency SoftLayer service
func GetAuxiliaryNotificationEmergencyService(sess session.SLSession) Auxiliary_Notification_Emergency {
	return Auxiliary_Notification_Emergency{Session: sess}
}

func (r Auxiliary_Notification_Emergency) Id(id int) Auxiliary_Notification_Emergency {
	r.Options.Id = &id
	return r
}

func (r Auxiliary_Notification_Emergency) Mask(mask string) Auxiliary_Notification_Emergency {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Auxiliary_Notification_Emergency) Filter(filter string) Auxiliary_Notification_Emergency {
	r.Options.Filter = filter
	return r
}

func (r Auxiliary_Notification_Emergency) Limit(limit int) Auxiliary_Notification_Emergency {
	r.Options.Limit = &limit
	return r
}

func (r Auxiliary_Notification_Emergency) Offset(offset int) Auxiliary_Notification_Emergency {
	r.Options.Offset = &offset
	return r
}

// Retrieve an array of SoftLayer_Auxiliary_Notification_Emergency data types, which contain all notification events regardless of status.
func (r Auxiliary_Notification_Emergency) GetAllObjects() (resp []datatypes.Auxiliary_Notification_Emergency, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Notification_Emergency", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve an array of SoftLayer_Auxiliary_Notification_Emergency data types, which contain all current notification events.
func (r Auxiliary_Notification_Emergency) GetCurrentNotifications() (resp []datatypes.Auxiliary_Notification_Emergency, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Notification_Emergency", "getCurrentNotifications", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Auxiliary_Notification_Emergency object, it can be used to check for current notifications being broadcast by SoftLayer.
func (r Auxiliary_Notification_Emergency) GetObject() (resp datatypes.Auxiliary_Notification_Emergency, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Notification_Emergency", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The signature of the SoftLayer employee department associated with this notification.
func (r Auxiliary_Notification_Emergency) GetSignature() (resp datatypes.Auxiliary_Notification_Emergency_Signature, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Notification_Emergency", "getSignature", nil, &r.Options, &resp)
	return
}

// Retrieve The status of this notification.
func (r Auxiliary_Notification_Emergency) GetStatus() (resp datatypes.Auxiliary_Notification_Emergency_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Notification_Emergency", "getStatus", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Auxiliary_Shipping_Courier_Type struct {
	Session session.SLSession
	Options sl.Options
}

// GetAuxiliaryShippingCourierTypeService returns an instance of the Auxiliary_Shipping_Courier_Type SoftLayer service
func GetAuxiliaryShippingCourierTypeService(sess session.SLSession) Auxiliary_Shipping_Courier_Type {
	return Auxiliary_Shipping_Courier_Type{Session: sess}
}

func (r Auxiliary_Shipping_Courier_Type) Id(id int) Auxiliary_Shipping_Courier_Type {
	r.Options.Id = &id
	return r
}

func (r Auxiliary_Shipping_Courier_Type) Mask(mask string) Auxiliary_Shipping_Courier_Type {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Auxiliary_Shipping_Courier_Type) Filter(filter string) Auxiliary_Shipping_Courier_Type {
	r.Options.Filter = filter
	return r
}

func (r Auxiliary_Shipping_Courier_Type) Limit(limit int) Auxiliary_Shipping_Courier_Type {
	r.Options.Limit = &limit
	return r
}

func (r Auxiliary_Shipping_Courier_Type) Offset(offset int) Auxiliary_Shipping_Courier_Type {
	r.Options.Offset = &offset
	return r
}

// Retrieve
func (r Auxiliary_Shipping_Courier_Type) GetCourier() (resp []datatypes.Auxiliary_Shipping_Courier, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Shipping_Courier_Type", "getCourier", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Auxiliary_Shipping_Courier_Type) GetObject() (resp datatypes.Auxiliary_Shipping_Courier_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Shipping_Courier_Type", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Auxiliary_Shipping_Courier_Type) GetTypeByKeyName(keyName *string) (resp datatypes.Auxiliary_Shipping_Courier_Type, err error) {
	params := []interface{}{
		keyName,
	}
	err = r.Session.DoRequest("SoftLayer_Auxiliary_Shipping_Courier_Type", "getTypeByKeyName", params, &r.Options, &resp)
	return
}

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

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// no documentation yet
type Workload_Citrix_Client struct {
	Session *session.Session
	Options sl.Options
}

// GetWorkloadCitrixClientService returns an instance of the Workload_Citrix_Client SoftLayer service
func GetWorkloadCitrixClientService(sess *session.Session) Workload_Citrix_Client {
	return Workload_Citrix_Client{Session: sess}
}

func (r Workload_Citrix_Client) Id(id int) Workload_Citrix_Client {
	r.Options.Id = &id
	return r
}

func (r Workload_Citrix_Client) Mask(mask string) Workload_Citrix_Client {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Workload_Citrix_Client) Filter(filter string) Workload_Citrix_Client {
	r.Options.Filter = filter
	return r
}

func (r Workload_Citrix_Client) Limit(limit int) Workload_Citrix_Client {
	r.Options.Limit = &limit
	return r
}

func (r Workload_Citrix_Client) Offset(offset int) Workload_Citrix_Client {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Workload_Citrix_Client) CreateResourceLocation(citrixCredentials *datatypes.Workload_Citrix_Request) (resp datatypes.Workload_Citrix_Client_Response, err error) {
	params := []interface{}{
		citrixCredentials,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Client", "createResourceLocation", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Workload_Citrix_Client) GetResourceLocations(citrixCredentials *datatypes.Workload_Citrix_Request) (resp datatypes.Workload_Citrix_Client_Response, err error) {
	params := []interface{}{
		citrixCredentials,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Client", "getResourceLocations", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Workload_Citrix_Client) ValidateCitrixCredentials(citrixCredentials *datatypes.Workload_Citrix_Request) (resp datatypes.Workload_Citrix_Client_Response, err error) {
	params := []interface{}{
		citrixCredentials,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Client", "validateCitrixCredentials", params, &r.Options, &resp)
	return
}

// no documentation yet
type Workload_Citrix_Workspace_Order struct {
	Session *session.Session
	Options sl.Options
}

// GetWorkloadCitrixWorkspaceOrderService returns an instance of the Workload_Citrix_Workspace_Order SoftLayer service
func GetWorkloadCitrixWorkspaceOrderService(sess *session.Session) Workload_Citrix_Workspace_Order {
	return Workload_Citrix_Workspace_Order{Session: sess}
}

func (r Workload_Citrix_Workspace_Order) Id(id int) Workload_Citrix_Workspace_Order {
	r.Options.Id = &id
	return r
}

func (r Workload_Citrix_Workspace_Order) Mask(mask string) Workload_Citrix_Workspace_Order {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Workload_Citrix_Workspace_Order) Filter(filter string) Workload_Citrix_Workspace_Order {
	r.Options.Filter = filter
	return r
}

func (r Workload_Citrix_Workspace_Order) Limit(limit int) Workload_Citrix_Workspace_Order {
	r.Options.Limit = &limit
	return r
}

func (r Workload_Citrix_Workspace_Order) Offset(offset int) Workload_Citrix_Workspace_Order {
	r.Options.Offset = &offset
	return r
}

// This method will cancel the resources associated with the provided VLAN and have a 'cvad' tag reference.
func (r Workload_Citrix_Workspace_Order) CancelWorkspaceResources(vlanIdentifier *string, cancelImmediately *bool, customerNote *string) (resp datatypes.Workload_Citrix_Workspace_Response_Result, err error) {
	params := []interface{}{
		vlanIdentifier,
		cancelImmediately,
		customerNote,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Workspace_Order", "cancelWorkspaceResources", params, &r.Options, &resp)
	return
}

// This method will return the list of names of VLANs which have a 'cvad' tag reference.  This name can be used with the cancelWorkspaceOrders method.
func (r Workload_Citrix_Workspace_Order) GetWorkspaceNames() (resp datatypes.Workload_Citrix_Workspace_Response_Result, err error) {
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Workspace_Order", "getWorkspaceNames", nil, &r.Options, &resp)
	return
}

// This method will return the list of resources which could be cancelled using cancelWorkspaceResources
func (r Workload_Citrix_Workspace_Order) GetWorkspaceResources(vlanIdentifier *string) (resp datatypes.Workload_Citrix_Workspace_Response_Result, err error) {
	params := []interface{}{
		vlanIdentifier,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Workspace_Order", "getWorkspaceResources", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Workload_Citrix_Workspace_Order) PlaceWorkspaceOrder(orderContainer *datatypes.Workload_Citrix_Workspace_Order_Container) (resp datatypes.Workload_Citrix_Workspace_Response, err error) {
	params := []interface{}{
		orderContainer,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Workspace_Order", "placeWorkspaceOrder", params, &r.Options, &resp)
	return
}

// This service is used to verify that an order meets all the necessary requirements to purchase Citrix Virtual Apps and Desktops running on IBM cloud.
func (r Workload_Citrix_Workspace_Order) VerifyWorkspaceOrder(orderContainer *datatypes.Workload_Citrix_Workspace_Order_Container) (resp datatypes.Workload_Citrix_Workspace_Response, err error) {
	params := []interface{}{
		orderContainer,
	}
	err = r.Session.DoRequest("SoftLayer_Workload_Citrix_Workspace_Order", "verifyWorkspaceOrder", params, &r.Options, &resp)
	return
}

// no documentation yet
type Workload_Citrix_Workspace_Response struct {
	Session *session.Session
	Options sl.Options
}

// GetWorkloadCitrixWorkspaceResponseService returns an instance of the Workload_Citrix_Workspace_Response SoftLayer service
func GetWorkloadCitrixWorkspaceResponseService(sess *session.Session) Workload_Citrix_Workspace_Response {
	return Workload_Citrix_Workspace_Response{Session: sess}
}

func (r Workload_Citrix_Workspace_Response) Id(id int) Workload_Citrix_Workspace_Response {
	r.Options.Id = &id
	return r
}

func (r Workload_Citrix_Workspace_Response) Mask(mask string) Workload_Citrix_Workspace_Response {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Workload_Citrix_Workspace_Response) Filter(filter string) Workload_Citrix_Workspace_Response {
	r.Options.Filter = filter
	return r
}

func (r Workload_Citrix_Workspace_Response) Limit(limit int) Workload_Citrix_Workspace_Response {
	r.Options.Limit = &limit
	return r
}

func (r Workload_Citrix_Workspace_Response) Offset(offset int) Workload_Citrix_Workspace_Response {
	r.Options.Offset = &offset
	return r
}

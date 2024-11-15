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
type Verify_Api_HttpObj struct {
	Session session.SLSession
	Options sl.Options
}

// GetVerifyApiHttpObjService returns an instance of the Verify_Api_HttpObj SoftLayer service
func GetVerifyApiHttpObjService(sess session.SLSession) Verify_Api_HttpObj {
	return Verify_Api_HttpObj{Session: sess}
}

func (r Verify_Api_HttpObj) Id(id int) Verify_Api_HttpObj {
	r.Options.Id = &id
	return r
}

func (r Verify_Api_HttpObj) Mask(mask string) Verify_Api_HttpObj {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Verify_Api_HttpObj) Filter(filter string) Verify_Api_HttpObj {
	r.Options.Filter = filter
	return r
}

func (r Verify_Api_HttpObj) Limit(limit int) Verify_Api_HttpObj {
	r.Options.Limit = &limit
	return r
}

func (r Verify_Api_HttpObj) Offset(offset int) Verify_Api_HttpObj {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Verify_Api_HttpObj) CreateObject(templateObject *datatypes.Verify_Api_HttpObj) (resp datatypes.Verify_Api_HttpObj, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpObj", "createObject", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpObj) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpObj", "deleteObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpObj) GetAllObjects() (resp []datatypes.Verify_Api_HttpObj, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpObj", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpObj) GetObject() (resp datatypes.Verify_Api_HttpObj, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpObj", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Verify_Api_HttpsObj struct {
	Session session.SLSession
	Options sl.Options
}

// GetVerifyApiHttpsObjService returns an instance of the Verify_Api_HttpsObj SoftLayer service
func GetVerifyApiHttpsObjService(sess session.SLSession) Verify_Api_HttpsObj {
	return Verify_Api_HttpsObj{Session: sess}
}

func (r Verify_Api_HttpsObj) Id(id int) Verify_Api_HttpsObj {
	r.Options.Id = &id
	return r
}

func (r Verify_Api_HttpsObj) Mask(mask string) Verify_Api_HttpsObj {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Verify_Api_HttpsObj) Filter(filter string) Verify_Api_HttpsObj {
	r.Options.Filter = filter
	return r
}

func (r Verify_Api_HttpsObj) Limit(limit int) Verify_Api_HttpsObj {
	r.Options.Limit = &limit
	return r
}

func (r Verify_Api_HttpsObj) Offset(offset int) Verify_Api_HttpsObj {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Verify_Api_HttpsObj) CreateObject(templateObject *datatypes.Verify_Api_HttpsObj) (resp datatypes.Verify_Api_HttpsObj, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpsObj", "createObject", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpsObj) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpsObj", "deleteObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpsObj) GetAllObjects() (resp []datatypes.Verify_Api_HttpsObj, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpsObj", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Verify_Api_HttpsObj) GetObject() (resp datatypes.Verify_Api_HttpsObj, err error) {
	err = r.Session.DoRequest("SoftLayer_Verify_Api_HttpsObj", "getObject", nil, &r.Options, &resp)
	return
}

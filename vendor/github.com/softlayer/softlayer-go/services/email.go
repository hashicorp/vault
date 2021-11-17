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
type Email_Subscription struct {
	Session *session.Session
	Options sl.Options
}

// GetEmailSubscriptionService returns an instance of the Email_Subscription SoftLayer service
func GetEmailSubscriptionService(sess *session.Session) Email_Subscription {
	return Email_Subscription{Session: sess}
}

func (r Email_Subscription) Id(id int) Email_Subscription {
	r.Options.Id = &id
	return r
}

func (r Email_Subscription) Mask(mask string) Email_Subscription {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Email_Subscription) Filter(filter string) Email_Subscription {
	r.Options.Filter = filter
	return r
}

func (r Email_Subscription) Limit(limit int) Email_Subscription {
	r.Options.Limit = &limit
	return r
}

func (r Email_Subscription) Offset(offset int) Email_Subscription {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Email_Subscription) Disable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription", "disable", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Email_Subscription) Enable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription", "enable", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Email_Subscription) GetAllObjects() (resp []datatypes.Email_Subscription, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Email_Subscription) GetEnabled() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription", "getEnabled", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Email_Subscription) GetObject() (resp datatypes.Email_Subscription, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Email_Subscription_Group struct {
	Session *session.Session
	Options sl.Options
}

// GetEmailSubscriptionGroupService returns an instance of the Email_Subscription_Group SoftLayer service
func GetEmailSubscriptionGroupService(sess *session.Session) Email_Subscription_Group {
	return Email_Subscription_Group{Session: sess}
}

func (r Email_Subscription_Group) Id(id int) Email_Subscription_Group {
	r.Options.Id = &id
	return r
}

func (r Email_Subscription_Group) Mask(mask string) Email_Subscription_Group {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Email_Subscription_Group) Filter(filter string) Email_Subscription_Group {
	r.Options.Filter = filter
	return r
}

func (r Email_Subscription_Group) Limit(limit int) Email_Subscription_Group {
	r.Options.Limit = &limit
	return r
}

func (r Email_Subscription_Group) Offset(offset int) Email_Subscription_Group {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Email_Subscription_Group) GetAllObjects() (resp []datatypes.Email_Subscription_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription_Group", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Email_Subscription_Group) GetObject() (resp datatypes.Email_Subscription_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription_Group", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve All email subscriptions associated with this group.
func (r Email_Subscription_Group) GetSubscriptions() (resp []datatypes.Email_Subscription, err error) {
	err = r.Session.DoRequest("SoftLayer_Email_Subscription_Group", "getSubscriptions", nil, &r.Options, &resp)
	return
}

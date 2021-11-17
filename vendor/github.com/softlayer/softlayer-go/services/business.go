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

// Contains business partner channel information
type Business_Partner_Channel struct {
	Session *session.Session
	Options sl.Options
}

// GetBusinessPartnerChannelService returns an instance of the Business_Partner_Channel SoftLayer service
func GetBusinessPartnerChannelService(sess *session.Session) Business_Partner_Channel {
	return Business_Partner_Channel{Session: sess}
}

func (r Business_Partner_Channel) Id(id int) Business_Partner_Channel {
	r.Options.Id = &id
	return r
}

func (r Business_Partner_Channel) Mask(mask string) Business_Partner_Channel {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Business_Partner_Channel) Filter(filter string) Business_Partner_Channel {
	r.Options.Filter = filter
	return r
}

func (r Business_Partner_Channel) Limit(limit int) Business_Partner_Channel {
	r.Options.Limit = &limit
	return r
}

func (r Business_Partner_Channel) Offset(offset int) Business_Partner_Channel {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Business_Partner_Channel) GetObject() (resp datatypes.Business_Partner_Channel, err error) {
	err = r.Session.DoRequest("SoftLayer_Business_Partner_Channel", "getObject", nil, &r.Options, &resp)
	return
}

// Contains business partner segment information
type Business_Partner_Segment struct {
	Session *session.Session
	Options sl.Options
}

// GetBusinessPartnerSegmentService returns an instance of the Business_Partner_Segment SoftLayer service
func GetBusinessPartnerSegmentService(sess *session.Session) Business_Partner_Segment {
	return Business_Partner_Segment{Session: sess}
}

func (r Business_Partner_Segment) Id(id int) Business_Partner_Segment {
	r.Options.Id = &id
	return r
}

func (r Business_Partner_Segment) Mask(mask string) Business_Partner_Segment {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Business_Partner_Segment) Filter(filter string) Business_Partner_Segment {
	r.Options.Filter = filter
	return r
}

func (r Business_Partner_Segment) Limit(limit int) Business_Partner_Segment {
	r.Options.Limit = &limit
	return r
}

func (r Business_Partner_Segment) Offset(offset int) Business_Partner_Segment {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Business_Partner_Segment) GetObject() (resp datatypes.Business_Partner_Segment, err error) {
	err = r.Session.DoRequest("SoftLayer_Business_Partner_Segment", "getObject", nil, &r.Options, &resp)
	return
}

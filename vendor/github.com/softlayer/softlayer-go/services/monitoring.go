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

// DEPRECATED. The SoftLayer_Monitoring_Robot data type contains general information relating to a monitoring robot.
type Monitoring_Robot struct {
	Session *session.Session
	Options sl.Options
}

// GetMonitoringRobotService returns an instance of the Monitoring_Robot SoftLayer service
func GetMonitoringRobotService(sess *session.Session) Monitoring_Robot {
	return Monitoring_Robot{Session: sess}
}

func (r Monitoring_Robot) Id(id int) Monitoring_Robot {
	r.Options.Id = &id
	return r
}

func (r Monitoring_Robot) Mask(mask string) Monitoring_Robot {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Monitoring_Robot) Filter(filter string) Monitoring_Robot {
	r.Options.Filter = filter
	return r
}

func (r Monitoring_Robot) Limit(limit int) Monitoring_Robot {
	r.Options.Limit = &limit
	return r
}

func (r Monitoring_Robot) Offset(offset int) Monitoring_Robot {
	r.Options.Offset = &offset
	return r
}

// DEPRECATED. Checks if a monitoring robot can communicate with SoftLayer monitoring management system via the private network.
//
// TCP port 48000 - 48002 must be open on your server or your virtual server in order for this test to succeed.
func (r Monitoring_Robot) CheckConnection() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Monitoring_Robot", "checkConnection", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Monitoring_Robot) GetObject() (resp datatypes.Monitoring_Robot, err error) {
	err = r.Session.DoRequest("SoftLayer_Monitoring_Robot", "getObject", nil, &r.Options, &resp)
	return
}

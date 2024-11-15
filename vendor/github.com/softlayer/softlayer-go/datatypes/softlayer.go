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

package datatypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Void is a dummy type for identifying void return values from methods
type Void int

// Time type overrides the default json marshaler with the SoftLayer custom format
type Time struct {
	time.Time
}

func (r Time) String() string {
	return r.Time.Format(time.RFC3339)
}

// MarshalJSON returns the json encoding of the datatypes.Time receiver.  This
// override is necessary to ensure datetimes are formatted in the way SoftLayer
// expects - that is, using the RFC3339 format, without nanoseconds.
func (r Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

// MarshalText returns a text encoding of the datatypes.Time receiver.  This
// is mainly provided to complete what might be expected of a type that
// implements the Marshaler interface.
func (r Time) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// FIXME: Need to have special unmarshaling of some values defined as float type
// in the metadata that actually come down as strings in the api.
// e.g. SoftLayer_Product_Item.capacity
// Float64 is a float type that deals with some of the oddities when
// unmarshalling from the SLAPI
//
// Code borrowed from https://github.com/sudorandom/softlayer-go/blob/master/slapi/types/float.go
type Float64 float64

// UnmarshalJSON statisied the json.Unmarshaler interface
func (f *Float64) UnmarshalJSON(data []byte) error {

	// Attempt parsing the float normally
	v, err := strconv.ParseFloat(string(data), 64)

	// Attempt parsing the float as a string
	if err != nil {
		if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
			return fmt.Errorf("malformed data")
		}

		v, err = strconv.ParseFloat(string(data[1:len(data)-1]), 64)
		if err != nil {
			return err
		}
	}
	*f = Float64(v)
	return nil
}

// Used to set the appropriate complexType field in the passed product order.
// Employs reflection to determine the type of the passed value and use it
// to derive the complexType to send to SoftLayer.
func SetComplexType(v interface{}) error {
	orderDataPtr := reflect.ValueOf(v)
	if orderDataPtr.Type().Name() != "" {
		return errors.New("Did not pass a pointer to a product order.")
	}

	orderDataValue := reflect.Indirect(reflect.ValueOf(v))
	orderDataType := orderDataValue.Type().Name()
	if !strings.HasPrefix(orderDataType, "Container_Product_Order") {
		return fmt.Errorf("Did not pass a pointer to a product order: %s", orderDataType)
	}

	complexTypeField := orderDataValue.FieldByName("ComplexType")
	complexType := "SoftLayer_" + orderDataType
	complexTypeField.Set(reflect.ValueOf(&complexType))

	return nil
}

// Used for detecting the search result before the real unmarshalling can begin
type Result_Detector struct {
	// An array of terms that were matched in the resource object.
	MatchedTerms []string `json:"matchedTerms,omitempty" xmlrpc:"matchedTerms,omitempty"`
	// The score ratio of the result for relevance to the search criteria.
	RelevanceScore *Float64 `json:"relevanceScore,omitempty" xmlrpc:"relevanceScore,omitempty"`
	// A search results resource object that matched search criteria.
	Resource interface{} `json:"resource,omitempty" xmlrpc:"resource,omitempty"`
	// The type of the resource object that matched search criteria.
	ResourceType *string `json:"resourceType,omitempty" xmlrpc:"resourceType,omitempty"`
}

func (result *Container_Search_Result) UnmarshalJSON(data []byte) error {
	var detector Result_Detector
	err := json.Unmarshal(data, &detector)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(*detector.ResourceType, "SoftLayer_") {
		return fmt.Errorf("Resource is not a SoftLayer_ Type, but instead is %v", *detector.ResourceType)
	}

	var theType interface{}
	switch *detector.ResourceType {
	case "SoftLayer_Virtual_Guest":
		theType = new(Virtual_Guest)
		// theType = reflect.New(reflect.TypeOf(Virtual_Guest{})).Interface().(*Virtual_Guest)
	case "SoftLayer_Event_Log":
		theType = new(Event_Log)
	case "SoftLayer_Virtual_DedicatedHost":
		theType = new(Virtual_DedicatedHost)
	case "SoftLayer_Hardware":
		theType = new(Hardware)
	case "SoftLayer_Network_Application_Delivery_Controller":
		theType = new(Network_Application_Delivery_Controller)
	case "SoftLayer_Network_Subnet_IpAddress":
		theType = new(Network_Subnet_IpAddress)
	case "SoftLayer_Network_Vlan":
		theType = new(Network_Vlan)
	case "SoftLayer_Network_Vlan_Firewall":
		theType = new(Network_Vlan_Firewall)
	case "SoftLayer_Ticket":
		theType = new(Ticket)
	}
	detector.Resource = theType
	err = json.Unmarshal(data, &detector)

	result.MatchedTerms = detector.MatchedTerms
	result.RelevanceScore = detector.RelevanceScore
	result.Resource = detector.Resource
	result.ResourceType = detector.ResourceType

	return nil
}

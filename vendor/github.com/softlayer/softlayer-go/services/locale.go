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
type Locale struct {
	Session session.SLSession
	Options sl.Options
}

// GetLocaleService returns an instance of the Locale SoftLayer service
func GetLocaleService(sess session.SLSession) Locale {
	return Locale{Session: sess}
}

func (r Locale) Id(id int) Locale {
	r.Options.Id = &id
	return r
}

func (r Locale) Mask(mask string) Locale {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Locale) Filter(filter string) Locale {
	r.Options.Filter = filter
	return r
}

func (r Locale) Limit(limit int) Locale {
	r.Options.Limit = &limit
	return r
}

func (r Locale) Offset(offset int) Locale {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Locale) GetClosestToLanguageTag(languageTag *string) (resp datatypes.Locale, err error) {
	params := []interface{}{
		languageTag,
	}
	err = r.Session.DoRequest("SoftLayer_Locale", "getClosestToLanguageTag", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Locale) GetObject() (resp datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Locale_Country struct {
	Session session.SLSession
	Options sl.Options
}

// GetLocaleCountryService returns an instance of the Locale_Country SoftLayer service
func GetLocaleCountryService(sess session.SLSession) Locale_Country {
	return Locale_Country{Session: sess}
}

func (r Locale_Country) Id(id int) Locale_Country {
	r.Options.Id = &id
	return r
}

func (r Locale_Country) Mask(mask string) Locale_Country {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Locale_Country) Filter(filter string) Locale_Country {
	r.Options.Filter = filter
	return r
}

func (r Locale_Country) Limit(limit int) Locale_Country {
	r.Options.Limit = &limit
	return r
}

func (r Locale_Country) Offset(offset int) Locale_Country {
	r.Options.Offset = &offset
	return r
}

// This method is to get the collection of VAT country codes and VAT ID Regexes.
func (r Locale_Country) GetAllVatCountryCodesAndVatIdRegexes() (resp []datatypes.Container_Collection_Locale_VatCountryCodeAndFormat, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getAllVatCountryCodesAndVatIdRegexes", nil, &r.Options, &resp)
	return
}

// Use this method to retrieve a list of countries and locale information available to the current user.
func (r Locale_Country) GetAvailableCountries() (resp []datatypes.Locale_Country, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getAvailableCountries", nil, &r.Options, &resp)
	return
}

// Use this method to retrieve a list of countries and locale information such as country code and state/provinces.
func (r Locale_Country) GetCountries() (resp []datatypes.Locale_Country, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getCountries", nil, &r.Options, &resp)
	return
}

// This method will return a collection of [[SoftLayer_Container_Collection_Locale_CountryCode]] objects. If the country has states, a [[SoftLayer_Container_Collection_Locale_StateCode]] collection will be provided with the country.
func (r Locale_Country) GetCountriesAndStates(usFirstFlag *bool) (resp []datatypes.Container_Collection_Locale_CountryCode, err error) {
	params := []interface{}{
		usFirstFlag,
	}
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getCountriesAndStates", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Locale_Country) GetObject() (resp datatypes.Locale_Country, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getObject", nil, &r.Options, &resp)
	return
}

// This method will return an array of country codes that require postal code
func (r Locale_Country) GetPostalCodeRequiredCountryCodes() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getPostalCodeRequiredCountryCodes", nil, &r.Options, &resp)
	return
}

// Retrieve States that belong to this country.
func (r Locale_Country) GetStates() (resp []datatypes.Locale_StateProvince, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getStates", nil, &r.Options, &resp)
	return
}

// This method will return an array of ISO 3166 Alpha-2 country codes that use a Value-Added Tax (VAT) ID. Note the difference between [[SoftLayer_Locale_Country/getVatRequiredCountryCodes]] - this method will provide <strong>all</strong> country codes that use VAT ID, including those which are required.
func (r Locale_Country) GetVatCountries() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getVatCountries", nil, &r.Options, &resp)
	return
}

// This method will return an array of ISO 3166 Alpha-2 country codes that use a Value-Added Tax (VAT) ID. Note the difference between [[SoftLayer_Locale_Country/getVatCountries]] - this method will provide country codes where a VAT ID is required for onboarding to IBM Cloud.
func (r Locale_Country) GetVatRequiredCountryCodes() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "getVatRequiredCountryCodes", nil, &r.Options, &resp)
	return
}

// Returns true if the country code is in the European Union (EU), false otherwise.
func (r Locale_Country) IsEuropeanUnionCountry(iso2CountryCode *string) (resp bool, err error) {
	params := []interface{}{
		iso2CountryCode,
	}
	err = r.Session.DoRequest("SoftLayer_Locale_Country", "isEuropeanUnionCountry", params, &r.Options, &resp)
	return
}

// Each User is assigned a timezone allowing for a precise local timestamp.
type Locale_Timezone struct {
	Session session.SLSession
	Options sl.Options
}

// GetLocaleTimezoneService returns an instance of the Locale_Timezone SoftLayer service
func GetLocaleTimezoneService(sess session.SLSession) Locale_Timezone {
	return Locale_Timezone{Session: sess}
}

func (r Locale_Timezone) Id(id int) Locale_Timezone {
	r.Options.Id = &id
	return r
}

func (r Locale_Timezone) Mask(mask string) Locale_Timezone {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Locale_Timezone) Filter(filter string) Locale_Timezone {
	r.Options.Filter = filter
	return r
}

func (r Locale_Timezone) Limit(limit int) Locale_Timezone {
	r.Options.Limit = &limit
	return r
}

func (r Locale_Timezone) Offset(offset int) Locale_Timezone {
	r.Options.Offset = &offset
	return r
}

// Retrieve all timezone objects.
func (r Locale_Timezone) GetAllObjects() (resp []datatypes.Locale_Timezone, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Timezone", "getAllObjects", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Locale_Timezone object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_Locale_Timezone service.
func (r Locale_Timezone) GetObject() (resp datatypes.Locale_Timezone, err error) {
	err = r.Session.DoRequest("SoftLayer_Locale_Timezone", "getObject", nil, &r.Options, &resp)
	return
}

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

// The SoftLayer_Brand data type contains brand information relating to the single SoftLayer customer account.
//
// IBM Cloud Infrastructure customers are unable to change their brand information in the portal or the API.
type Brand struct {
	Session session.SLSession
	Options sl.Options
}

// GetBrandService returns an instance of the Brand SoftLayer service
func GetBrandService(sess session.SLSession) Brand {
	return Brand{Session: sess}
}

func (r Brand) Id(id int) Brand {
	r.Options.Id = &id
	return r
}

func (r Brand) Mask(mask string) Brand {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Brand) Filter(filter string) Brand {
	r.Options.Filter = filter
	return r
}

func (r Brand) Limit(limit int) Brand {
	r.Options.Limit = &limit
	return r
}

func (r Brand) Offset(offset int) Brand {
	r.Options.Offset = &offset
	return r
}

// Create a new customer account record. By default, the newly created account will be associated to a platform (PaaS) account. To skip the automatic creation and linking to a new platform account, set the <em>bluemixLinkedFlag</em> to <strong>false</strong> on the account template.
func (r Brand) CreateCustomerAccount(account *datatypes.Account, bypassDuplicateAccountCheck *bool) (resp datatypes.Account, err error) {
	params := []interface{}{
		account,
		bypassDuplicateAccountCheck,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "createCustomerAccount", params, &r.Options, &resp)
	return
}

// createObject() allows the creation of a new brand. This will also create an `account`
// to serve as the owner of the brand.
//
// In order to create a brand, a template object must be sent in with several required values.
//
// ### Input [[SoftLayer_Brand]]
//
// - `name`
//   - Name of brand
//   - Required
//   - Type: string
//
// - `keyName`
//   - Reference key name
//   - Required
//   - Type: string
//
// - `longName`
//   - More descriptive name of brand
//   - Required
//   - Type: string
//
// - `account.firstName`
//   - First Name of account contact
//   - Required
//   - Type: string
//
// - `account.lastName`
//   - Last Name of account contact
//   - Required
//   - Type: string
//
// - `account.address1`
//   - Street Address of company
//   - Required
//   - Type: string
//
// - `account.address2`
//   - Street Address of company
//   - Optional
//   - Type: string
//
// - `account.city`
//   - City of company
//   - Required
//   - Type: string
//
// - `account.state`
//   - State of company (if applicable)
//   - Conditionally Required
//   - Type: string
//
// - `account.postalCode`
//   - Postal Code of company
//   - Required
//   - Type: string
//
// - `account.country`
//   - Country of company
//   - Required
//   - Type: string
//
// - `account.officePhone`
//   - Office Phone number of Company
//   - Required
//   - Type: string
//
// - `account.alternatePhone`
//   - Alternate Phone number of Company
//   - Optional
//   - Type: string
//
// - `account.companyName`
//   - Name of company
//   - Required
//   - Type: string
//
// - `account.email`
//   - Email address of account contact
//   - Required
//   - Type: string
//
// REST Example:
// ```
//
//	curl -X POST -d '{
//	    "parameters":[{
//	        "name": "Brand Corp",
//	        "keyName": "BRAND_CORP",
//	        "longName": "Brand Corporation",
//	        "account": {
//	            "firstName": "Gloria",
//	            "lastName": "Brand",
//	            "address1": "123 Drive",
//	            "city": "Boston",
//	            "state": "MA",
//	            "postalCode": "02107",
//	            "country": "US",
//	            "companyName": "Brand Corp",
//	            "officePhone": "857-111-1111",
//	            "email": "noreply@example.com"
//	        }
//	    }]
//	}' https://api.softlayer.com/rest/v3.1/SoftLayer_Brand/createObject.json
//
// ```
func (r Brand) CreateObject(templateObject *datatypes.Brand) (resp datatypes.Brand, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "createObject", params, &r.Options, &resp)
	return
}

// Disable an account associated with this Brand.  Anything that would disqualify the account from being disabled will cause an exception to be raised.
func (r Brand) DisableAccount(accountId *int) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "disableAccount", params, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve All accounts owned by the brand.
func (r Brand) GetAllOwnedAccounts() (resp []datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getAllOwnedAccounts", nil, &r.Options, &resp)
	return
}

// (DEPRECATED) Use [[SoftLayer_Ticket_Subject::getAllObjects]] method.
// Deprecated: This function has been marked as deprecated.
func (r Brand) GetAllTicketSubjects(account *datatypes.Account) (resp []datatypes.Ticket_Subject, err error) {
	params := []interface{}{
		account,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "getAllTicketSubjects", params, &r.Options, &resp)
	return
}

// Retrieve This flag indicates if creation of accounts is allowed.
func (r Brand) GetAllowAccountCreationFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getAllowAccountCreationFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Returns snapshots of billing items recorded periodically given an account ID owned by the brand those billing items belong to. Retrieving billing item snapshots is more performant than retrieving billing items directly and performs less relational joins improving retrieval efficiency. The downside is, they are not real time, and do not share relational parity with the original billing item.
func (r Brand) GetBillingItemSnapshots() (resp []datatypes.Billing_Item_Chronicle, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getBillingItemSnapshots", nil, &r.Options, &resp)
	return
}

// This service returns the snapshots of billing items recorded periodically given an account ID. The provided account ID must be owned by the brand that calls this service. In this context, it can be interpreted that the billing items snapshots belong to both the account and that accounts brand. Retrieving billing item snapshots is more performant than retrieving billing items directly and performs less relational joins improving retrieval efficiency.
//
// The downside is, they are not real time, and do not share relational parity with the original billing item.
func (r Brand) GetBillingItemSnapshotsForSingleOwnedAccount(accountId *int) (resp []datatypes.Billing_Item_Chronicle, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "getBillingItemSnapshotsForSingleOwnedAccount", params, &r.Options, &resp)
	return
}

// This service returns the snapshots of billing items recorded periodically given an account ID owned by the brand those billing items belong to. Retrieving billing item snapshots is more performant than retrieving billing items directly and performs less relational joins improving retrieval efficiency.
//
// The downside is, they are not real time, and do not share relational parity with the original billing item.
func (r Brand) GetBillingItemSnapshotsWithExternalAccountId(externalAccountId *string) (resp []datatypes.Billing_Item_Chronicle, err error) {
	params := []interface{}{
		externalAccountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "getBillingItemSnapshotsWithExternalAccountId", params, &r.Options, &resp)
	return
}

// Retrieve Business Partner details for the brand. Country Enterprise Code, Channel, Segment, Reseller Level.
func (r Brand) GetBusinessPartner() (resp datatypes.Brand_Business_Partner, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getBusinessPartner", nil, &r.Options, &resp)
	return
}

// Retrieve Flag indicating if the brand is a business partner.
func (r Brand) GetBusinessPartnerFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getBusinessPartnerFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The Product Catalog for the Brand
func (r Brand) GetCatalog() (resp datatypes.Product_Catalog, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getCatalog", nil, &r.Options, &resp)
	return
}

// Retrieve the contact information for the brand such as the corporate or support contact.  This will include the contact name, telephone number, fax number, email address, and mailing address of the contact.
func (r Brand) GetContactInformation() (resp []datatypes.Brand_Contact, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getContactInformation", nil, &r.Options, &resp)
	return
}

// Retrieve The contacts for the brand.
func (r Brand) GetContacts() (resp []datatypes.Brand_Contact, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getContacts", nil, &r.Options, &resp)
	return
}

// Retrieve This references relationship between brands, locations and countries associated with a user's account that are ineligible when ordering products. For example, the India datacenter may not be available on this brand for customers that live in Great Britain.
func (r Brand) GetCustomerCountryLocationRestrictions() (resp []datatypes.Brand_Restriction_Location_CustomerCountry, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getCustomerCountryLocationRestrictions", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetDistributor() (resp datatypes.Brand, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getDistributor", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetDistributorChildFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getDistributorChildFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetDistributorFlag() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getDistributorFlag", nil, &r.Options, &resp)
	return
}

// Retrieve An account's associated hardware objects.
func (r Brand) GetHardware() (resp []datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getHardware", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetHasAgentAdvancedSupportFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getHasAgentAdvancedSupportFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetHasAgentSupportFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getHasAgentSupportFlag", nil, &r.Options, &resp)
	return
}

// Get the payment processor merchant name.
func (r Brand) GetMerchantName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getMerchantName", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Brand) GetObject() (resp datatypes.Brand, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetOpenTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getOpenTickets", nil, &r.Options, &resp)
	return
}

// Retrieve Active accounts owned by the brand.
func (r Brand) GetOwnedAccounts() (resp []datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getOwnedAccounts", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetSecurityLevel() (resp datatypes.Security_Level, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getSecurityLevel", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetTicketGroups() (resp []datatypes.Ticket_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getTicketGroups", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getTickets", nil, &r.Options, &resp)
	return
}

// (DEPRECATED) Use [[SoftLayer_User_Customer::getImpersonationToken]] method.
func (r Brand) GetToken(userId *int) (resp string, err error) {
	params := []interface{}{
		userId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "getToken", params, &r.Options, &resp)
	return
}

// Retrieve
func (r Brand) GetUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getUsers", nil, &r.Options, &resp)
	return
}

// Retrieve An account's associated virtual guest objects.
func (r Brand) GetVirtualGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "getVirtualGuests", nil, &r.Options, &resp)
	return
}

// Check if the brand is IBM SLIC top level brand or sub brand.
func (r Brand) IsIbmSlicBrand() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "isIbmSlicBrand", nil, &r.Options, &resp)
	return
}

// Check if the alternate billing system of brand is Bluemix.
func (r Brand) IsPlatformServicesBrand() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand", "isPlatformServicesBrand", nil, &r.Options, &resp)
	return
}

// Will attempt to migrate an external account to the brand in context.
func (r Brand) MigrateExternalAccount(accountId *int) (resp datatypes.Account_Brand_Migration_Request, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "migrateExternalAccount", params, &r.Options, &resp)
	return
}

// Reactivate an account associated with this Brand.  Anything that would disqualify the account from being reactivated will cause an exception to be raised.
func (r Brand) ReactivateAccount(accountId *int) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "reactivateAccount", params, &r.Options, &resp)
	return
}

// When this service is called given an IBM Cloud infrastructure account ID owned by the calling brand, the process is started to refresh the billing item snapshots belonging to that account. This refresh is async and can take an undetermined amount of time. Even if this endpoint returns an OK, it doesn't guarantee that refresh did not fail or encounter issues.
func (r Brand) RefreshBillingItemSnapshot(accountId *int) (resp bool, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "refreshBillingItemSnapshot", params, &r.Options, &resp)
	return
}

// Verify that an account may be disabled by a Brand Agent.  Anything that would disqualify the account from being disabled will cause an exception to be raised.
func (r Brand) VerifyCanDisableAccount(accountId *int) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "verifyCanDisableAccount", params, &r.Options, &resp)
	return
}

// Verify that an account may be reactivated by a Brand Agent.  Anything that would disqualify the account from being reactivated will cause an exception to be raised.
func (r Brand) VerifyCanReactivateAccount(accountId *int) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Brand", "verifyCanReactivateAccount", params, &r.Options, &resp)
	return
}

// Contains business partner details associated with a brand. Country Enterprise Identifier (CEID), Channel ID, Segment ID and Reseller Level.
type Brand_Business_Partner struct {
	Session session.SLSession
	Options sl.Options
}

// GetBrandBusinessPartnerService returns an instance of the Brand_Business_Partner SoftLayer service
func GetBrandBusinessPartnerService(sess session.SLSession) Brand_Business_Partner {
	return Brand_Business_Partner{Session: sess}
}

func (r Brand_Business_Partner) Id(id int) Brand_Business_Partner {
	r.Options.Id = &id
	return r
}

func (r Brand_Business_Partner) Mask(mask string) Brand_Business_Partner {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Brand_Business_Partner) Filter(filter string) Brand_Business_Partner {
	r.Options.Filter = filter
	return r
}

func (r Brand_Business_Partner) Limit(limit int) Brand_Business_Partner {
	r.Options.Limit = &limit
	return r
}

func (r Brand_Business_Partner) Offset(offset int) Brand_Business_Partner {
	r.Options.Offset = &offset
	return r
}

// Retrieve Brand associated with the business partner data
func (r Brand_Business_Partner) GetBrand() (resp datatypes.Brand, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Business_Partner", "getBrand", nil, &r.Options, &resp)
	return
}

// Retrieve Channel indicator used to categorize business partner revenue.
func (r Brand_Business_Partner) GetChannel() (resp datatypes.Business_Partner_Channel, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Business_Partner", "getChannel", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Brand_Business_Partner) GetObject() (resp datatypes.Brand_Business_Partner, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Business_Partner", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve Segment indicator used to categorize business partner revenue.
func (r Brand_Business_Partner) GetSegment() (resp datatypes.Business_Partner_Segment, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Business_Partner", "getSegment", nil, &r.Options, &resp)
	return
}

// The [[SoftLayer_Brand_Restriction_Location_CustomerCountry]] data type defines the relationship between brands, locations and countries associated with a user's account that are ineligible when ordering products. For example, the India datacenter may not be available on the SoftLayer US brand for customers that live in Great Britain.
type Brand_Restriction_Location_CustomerCountry struct {
	Session session.SLSession
	Options sl.Options
}

// GetBrandRestrictionLocationCustomerCountryService returns an instance of the Brand_Restriction_Location_CustomerCountry SoftLayer service
func GetBrandRestrictionLocationCustomerCountryService(sess session.SLSession) Brand_Restriction_Location_CustomerCountry {
	return Brand_Restriction_Location_CustomerCountry{Session: sess}
}

func (r Brand_Restriction_Location_CustomerCountry) Id(id int) Brand_Restriction_Location_CustomerCountry {
	r.Options.Id = &id
	return r
}

func (r Brand_Restriction_Location_CustomerCountry) Mask(mask string) Brand_Restriction_Location_CustomerCountry {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Brand_Restriction_Location_CustomerCountry) Filter(filter string) Brand_Restriction_Location_CustomerCountry {
	r.Options.Filter = filter
	return r
}

func (r Brand_Restriction_Location_CustomerCountry) Limit(limit int) Brand_Restriction_Location_CustomerCountry {
	r.Options.Limit = &limit
	return r
}

func (r Brand_Restriction_Location_CustomerCountry) Offset(offset int) Brand_Restriction_Location_CustomerCountry {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Brand_Restriction_Location_CustomerCountry) GetAllObjects() (resp []datatypes.Brand_Restriction_Location_CustomerCountry, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Restriction_Location_CustomerCountry", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve This references the brand that has a brand-location-country restriction setup.
func (r Brand_Restriction_Location_CustomerCountry) GetBrand() (resp datatypes.Brand, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Restriction_Location_CustomerCountry", "getBrand", nil, &r.Options, &resp)
	return
}

// Retrieve This references the datacenter that has a brand-location-country restriction setup. For example, if a datacenter is listed with a restriction for Canada, a Canadian customer may not be eligible to order services at that location.
func (r Brand_Restriction_Location_CustomerCountry) GetLocation() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Restriction_Location_CustomerCountry", "getLocation", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Brand_Restriction_Location_CustomerCountry) GetObject() (resp datatypes.Brand_Restriction_Location_CustomerCountry, err error) {
	err = r.Session.DoRequest("SoftLayer_Brand_Restriction_Location_CustomerCountry", "getObject", nil, &r.Options, &resp)
	return
}

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

// The SoftLayer_Brand data type contains brand information relating to the single SoftLayer customer account.
//
// IBM Cloud Infrastructure customers are unable to change their brand information in the portal or the API.
type Brand struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A count of all accounts owned by the brand.
	AllOwnedAccountCount *uint `json:"allOwnedAccountCount,omitempty" xmlrpc:"allOwnedAccountCount,omitempty"`

	// All accounts owned by the brand.
	AllOwnedAccounts []Account `json:"allOwnedAccounts,omitempty" xmlrpc:"allOwnedAccounts,omitempty"`

	// This flag indicates if creation of accounts is allowed.
	AllowAccountCreationFlag *bool `json:"allowAccountCreationFlag,omitempty" xmlrpc:"allowAccountCreationFlag,omitempty"`

	// A count of returns snapshots of billing items recorded periodically given an account ID owned by the brand those billing items belong to. Retrieving billing item snapshots is more performant than retrieving billing items directly and performs less relational joins improving retrieval efficiency. The downside is, they are not real time, and do not share relational parity with the original billing item.
	BillingItemSnapshotCount *uint `json:"billingItemSnapshotCount,omitempty" xmlrpc:"billingItemSnapshotCount,omitempty"`

	// Returns snapshots of billing items recorded periodically given an account ID owned by the brand those billing items belong to. Retrieving billing item snapshots is more performant than retrieving billing items directly and performs less relational joins improving retrieval efficiency. The downside is, they are not real time, and do not share relational parity with the original billing item.
	BillingItemSnapshots []Billing_Item_Chronicle `json:"billingItemSnapshots,omitempty" xmlrpc:"billingItemSnapshots,omitempty"`

	// Business Partner details for the brand. Country Enterprise Code, Channel, Segment, Reseller Level.
	BusinessPartner *Brand_Business_Partner `json:"businessPartner,omitempty" xmlrpc:"businessPartner,omitempty"`

	// Flag indicating if the brand is a business partner.
	BusinessPartnerFlag *bool `json:"businessPartnerFlag,omitempty" xmlrpc:"businessPartnerFlag,omitempty"`

	// The Product Catalog for the Brand
	Catalog *Product_Catalog `json:"catalog,omitempty" xmlrpc:"catalog,omitempty"`

	// ID of the Catalog used by this Brand
	CatalogId *int `json:"catalogId,omitempty" xmlrpc:"catalogId,omitempty"`

	// A count of the contacts for the brand.
	ContactCount *uint `json:"contactCount,omitempty" xmlrpc:"contactCount,omitempty"`

	// The contacts for the brand.
	Contacts []Brand_Contact `json:"contacts,omitempty" xmlrpc:"contacts,omitempty"`

	// A count of this references relationship between brands, locations and countries associated with a user's account that are ineligible when ordering products. For example, the India datacenter may not be available on this brand for customers that live in Great Britain.
	CustomerCountryLocationRestrictionCount *uint `json:"customerCountryLocationRestrictionCount,omitempty" xmlrpc:"customerCountryLocationRestrictionCount,omitempty"`

	// This references relationship between brands, locations and countries associated with a user's account that are ineligible when ordering products. For example, the India datacenter may not be available on this brand for customers that live in Great Britain.
	CustomerCountryLocationRestrictions []Brand_Restriction_Location_CustomerCountry `json:"customerCountryLocationRestrictions,omitempty" xmlrpc:"customerCountryLocationRestrictions,omitempty"`

	// no documentation yet
	Distributor *Brand `json:"distributor,omitempty" xmlrpc:"distributor,omitempty"`

	// no documentation yet
	DistributorChildFlag *bool `json:"distributorChildFlag,omitempty" xmlrpc:"distributorChildFlag,omitempty"`

	// no documentation yet
	DistributorFlag *string `json:"distributorFlag,omitempty" xmlrpc:"distributorFlag,omitempty"`

	// An account's associated hardware objects.
	Hardware []Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// A count of an account's associated hardware objects.
	HardwareCount *uint `json:"hardwareCount,omitempty" xmlrpc:"hardwareCount,omitempty"`

	// no documentation yet
	HasAgentAdvancedSupportFlag *bool `json:"hasAgentAdvancedSupportFlag,omitempty" xmlrpc:"hasAgentAdvancedSupportFlag,omitempty"`

	// no documentation yet
	HasAgentSupportFlag *bool `json:"hasAgentSupportFlag,omitempty" xmlrpc:"hasAgentSupportFlag,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The brand key name.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The brand long name.
	LongName *string `json:"longName,omitempty" xmlrpc:"longName,omitempty"`

	// The brand name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of
	OpenTicketCount *uint `json:"openTicketCount,omitempty" xmlrpc:"openTicketCount,omitempty"`

	// no documentation yet
	OpenTickets []Ticket `json:"openTickets,omitempty" xmlrpc:"openTickets,omitempty"`

	// A count of active accounts owned by the brand.
	OwnedAccountCount *uint `json:"ownedAccountCount,omitempty" xmlrpc:"ownedAccountCount,omitempty"`

	// Active accounts owned by the brand.
	OwnedAccounts []Account `json:"ownedAccounts,omitempty" xmlrpc:"ownedAccounts,omitempty"`

	// no documentation yet
	SecurityLevel *Security_Level `json:"securityLevel,omitempty" xmlrpc:"securityLevel,omitempty"`

	// A count of
	TicketCount *uint `json:"ticketCount,omitempty" xmlrpc:"ticketCount,omitempty"`

	// A count of
	TicketGroupCount *uint `json:"ticketGroupCount,omitempty" xmlrpc:"ticketGroupCount,omitempty"`

	// no documentation yet
	TicketGroups []Ticket_Group `json:"ticketGroups,omitempty" xmlrpc:"ticketGroups,omitempty"`

	// no documentation yet
	Tickets []Ticket `json:"tickets,omitempty" xmlrpc:"tickets,omitempty"`

	// A count of
	UserCount *uint `json:"userCount,omitempty" xmlrpc:"userCount,omitempty"`

	// no documentation yet
	Users []User_Customer `json:"users,omitempty" xmlrpc:"users,omitempty"`

	// A count of an account's associated virtual guest objects.
	VirtualGuestCount *uint `json:"virtualGuestCount,omitempty" xmlrpc:"virtualGuestCount,omitempty"`

	// An account's associated virtual guest objects.
	VirtualGuests []Virtual_Guest `json:"virtualGuests,omitempty" xmlrpc:"virtualGuests,omitempty"`
}

// no documentation yet
type Brand_Attribute struct {
	Entity

	// no documentation yet
	Brand *Brand `json:"brand,omitempty" xmlrpc:"brand,omitempty"`
}

// Contains business partner details associated with a brand. Country Enterprise Identifier (CEID), Channel ID, Segment ID and Reseller Level.
type Brand_Business_Partner struct {
	Entity

	// Brand associated with the business partner data
	Brand *Brand `json:"brand,omitempty" xmlrpc:"brand,omitempty"`

	// Channel indicator used to categorize business partner revenue.
	Channel *Business_Partner_Channel `json:"channel,omitempty" xmlrpc:"channel,omitempty"`

	// Brand business partner channel identifier
	ChannelId *int `json:"channelId,omitempty" xmlrpc:"channelId,omitempty"`

	// Brand business partner country enterprise code
	CountryEnterpriseCode *string `json:"countryEnterpriseCode,omitempty" xmlrpc:"countryEnterpriseCode,omitempty"`

	// Reseller level of a brand business partner
	ResellerLevel *int `json:"resellerLevel,omitempty" xmlrpc:"resellerLevel,omitempty"`

	// Segment indicator used to categorize business partner revenue.
	Segment *Business_Partner_Segment `json:"segment,omitempty" xmlrpc:"segment,omitempty"`

	// Brand business partner segment identifier
	SegmentId *int `json:"segmentId,omitempty" xmlrpc:"segmentId,omitempty"`
}

// SoftLayer_Brand_Contact contains the contact information for the brand such as Corporate or Support contact information
type Brand_Contact struct {
	Entity

	// The contact's address 1.
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// The contact's address 2.
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// The contact's alternate phone number.
	AlternatePhone *string `json:"alternatePhone,omitempty" xmlrpc:"alternatePhone,omitempty"`

	// no documentation yet
	Brand *Brand `json:"brand,omitempty" xmlrpc:"brand,omitempty"`

	// no documentation yet
	BrandContactType *Brand_Contact_Type `json:"brandContactType,omitempty" xmlrpc:"brandContactType,omitempty"`

	// The contact's type identifier.
	BrandContactTypeId *int `json:"brandContactTypeId,omitempty" xmlrpc:"brandContactTypeId,omitempty"`

	// The contact's city.
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// The contact's country.
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// The contact's email address.
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// The contact's fax number.
	FaxPhone *string `json:"faxPhone,omitempty" xmlrpc:"faxPhone,omitempty"`

	// The contact's first name.
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// The contact's last name.
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// The contact's phone number.
	OfficePhone *string `json:"officePhone,omitempty" xmlrpc:"officePhone,omitempty"`

	// The contact's postal code.
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// The contact's state.
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// SoftLayer_Brand_Contact_Type contains the contact type information for the brand contacts such as Corporate or Support contact type
type Brand_Contact_Type struct {
	Entity

	// Contact type description.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Contact type key name.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// Contact type name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Brand_Payment_Processor struct {
	Entity

	// no documentation yet
	Brand *Brand `json:"brand,omitempty" xmlrpc:"brand,omitempty"`

	// no documentation yet
	PaymentProcessor *Billing_Payment_Processor `json:"paymentProcessor,omitempty" xmlrpc:"paymentProcessor,omitempty"`
}

// The [[SoftLayer_Brand_Restriction_Location_CustomerCountry]] data type defines the relationship between brands, locations and countries associated with a user's account that are ineligible when ordering products. For example, the India datacenter may not be available on the SoftLayer US brand for customers that live in Great Britain.
type Brand_Restriction_Location_CustomerCountry struct {
	Entity

	// This references the brand that has a brand-location-country restriction setup.
	Brand *Brand `json:"brand,omitempty" xmlrpc:"brand,omitempty"`

	// The brand associated with customer's account.
	BrandId *int `json:"brandId,omitempty" xmlrpc:"brandId,omitempty"`

	// country code associated with customer's account.
	CustomerCountryCode *string `json:"customerCountryCode,omitempty" xmlrpc:"customerCountryCode,omitempty"`

	// This references the datacenter that has a brand-location-country restriction setup. For example, if a datacenter is listed with a restriction for Canada, a Canadian customer may not be eligible to order services at that location.
	Location *Location `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// The id for datacenter location.
	LocationId *int `json:"locationId,omitempty" xmlrpc:"locationId,omitempty"`
}

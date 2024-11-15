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

// no documentation yet
type Container_Account_Authentication_OpenIdConnect_UsernameLookupContainer struct {
	Entity

	// no documentation yet
	Active *bool `json:"active,omitempty" xmlrpc:"active,omitempty"`

	// no documentation yet
	EmailAddress *string `json:"emailAddress,omitempty" xmlrpc:"emailAddress,omitempty"`

	// no documentation yet
	FamilyName *string `json:"familyName,omitempty" xmlrpc:"familyName,omitempty"`

	// no documentation yet
	Federated *bool `json:"federated,omitempty" xmlrpc:"federated,omitempty"`

	// no documentation yet
	FoundAs *string `json:"foundAs,omitempty" xmlrpc:"foundAs,omitempty"`

	// no documentation yet
	GivenName *string `json:"givenName,omitempty" xmlrpc:"givenName,omitempty"`

	// no documentation yet
	NumberOfIbmIdsWithEmailAddress *int `json:"numberOfIbmIdsWithEmailAddress,omitempty" xmlrpc:"numberOfIbmIdsWithEmailAddress,omitempty"`

	// no documentation yet
	Realm *string `json:"realm,omitempty" xmlrpc:"realm,omitempty"`

	// no documentation yet
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`

	// no documentation yet
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// SoftLayer_Container_Account_Discount_Program models a single outbound object for a graph of given data sets.
type Container_Account_Discount_Program struct {
	Entity

	// The credit allowance that has already been applied during the current billing cycle. If the lifetime limit has been or soon will be reached, this amount may included credit applied in previous billing cycles.
	AppliedCredit *Float64 `json:"appliedCredit,omitempty" xmlrpc:"appliedCredit,omitempty"`

	// Flag to signify whether the account is a participant in the discount program.
	IsParticipant *bool `json:"isParticipant,omitempty" xmlrpc:"isParticipant,omitempty"`

	// Credit allowance applied over the course of the entire program enrollment. For enrollments without a lifetime restriction, this property will not be populated as credit will be tracked on a purely monthly basis.
	LifetimeAppliedCredit *Float64 `json:"lifetimeAppliedCredit,omitempty" xmlrpc:"lifetimeAppliedCredit,omitempty"`

	// Credit allowance available over the course of the entire program enrollment. If null, enrollment credit is applied on a strictly monthly basis and there is no lifetime maximum. Enrollments with non-null lifetime credit will receive the lesser of the remaining monthly credit or the remaining lifetime credit.
	LifetimeCredit *Float64 `json:"lifetimeCredit,omitempty" xmlrpc:"lifetimeCredit,omitempty"`

	// Remaining credit allowance available over the remaining duration of the program enrollment. If null, enrollment credit is applied on a strictly monthly basis and there is no lifetime maximum. Enrollments with non-null remaining lifetime credit will receive the lesser of the remaining monthly credit or the remaining lifetime credit.
	LifetimeRemainingCredit *Float64 `json:"lifetimeRemainingCredit,omitempty" xmlrpc:"lifetimeRemainingCredit,omitempty"`

	// Maximum number of orders the enrolled account is allowed to have open at one time. If null, then the Flexible Credit Program does not impose an order limit.
	MaximumActiveOrders *Float64 `json:"maximumActiveOrders,omitempty" xmlrpc:"maximumActiveOrders,omitempty"`

	// The monthly credit allowance that is available at the beginning of the billing cycle.
	MonthlyCredit *Float64 `json:"monthlyCredit,omitempty" xmlrpc:"monthlyCredit,omitempty"`

	// DEPRECATED: Taxes are calculated in real time and discount amounts are shown pre-tax in all cases. Tax values in the SoftLayer_Container_Account_Discount_Program container are now populated with the related pre-tax values.
	PostTaxRemainingCredit *Float64 `json:"postTaxRemainingCredit,omitempty" xmlrpc:"postTaxRemainingCredit,omitempty"`

	// The date at which the program expires in MM/DD/YYYY format.
	ProgramEndDate *Time `json:"programEndDate,omitempty" xmlrpc:"programEndDate,omitempty"`

	// Name of the Flexible Credit Program the account is enrolled in.
	ProgramName *string `json:"programName,omitempty" xmlrpc:"programName,omitempty"`

	// The credit allowance that is available during the current billing cycle. If the lifetime limit has been or soon will be reached, this amount may be reduced by credit applied in previous billing cycles.
	RemainingCredit *Float64 `json:"remainingCredit,omitempty" xmlrpc:"remainingCredit,omitempty"`

	// DEPRECATED: Taxes are calculated in real time and discount amounts are shown pre-tax in all cases. Tax values in the SoftLayer_Container_Account_Discount_Program container are now populated with the related pre-tax values.
	RemainingCreditTax *Float64 `json:"remainingCreditTax,omitempty" xmlrpc:"remainingCreditTax,omitempty"`
}

// no documentation yet
type Container_Account_Discount_Program_Collection struct {
	Entity

	// The amount of credit that has been used by all account level enrollments in the billing cycle.
	AccountLevelAppliedCredit *Float64 `json:"accountLevelAppliedCredit,omitempty" xmlrpc:"accountLevelAppliedCredit,omitempty"`

	// Account level credit allowance applied over the course of entire active program enrollments. For enrollments without a lifetime restriction, this property will not be populated as credit will be tracked on a purely monthly basis.
	AccountLevelLifetimeAppliedCredit *Float64 `json:"accountLevelLifetimeAppliedCredit,omitempty" xmlrpc:"accountLevelLifetimeAppliedCredit,omitempty"`

	// The total account level credit over the course of an entire program enrollment. This value may be null, in which case the enrollment credit is applied on a monthly basis and there is no lifetime maximum.
	AccountLevelLifetimeCredit *Float64 `json:"accountLevelLifetimeCredit,omitempty" xmlrpc:"accountLevelLifetimeCredit,omitempty"`

	// Remaining account level credit allowance available over the remaining duration of the program enrollments. If null, enrollment credit is applied on a strictly monthly basis and there is no lifetime maximum. Enrollments with non-null remaining lifetime credit will receive the lesser of the remaining monthly credit or the remaining lifetime credit.
	AccountLevelLifetimeRemainingCredit *Float64 `json:"accountLevelLifetimeRemainingCredit,omitempty" xmlrpc:"accountLevelLifetimeRemainingCredit,omitempty"`

	// The total account level monthly credit allowance available at the beginning of a billing cycle.
	AccountLevelMonthlyCredit *Float64 `json:"accountLevelMonthlyCredit,omitempty" xmlrpc:"accountLevelMonthlyCredit,omitempty"`

	// The total account level credit allowance still available during the current billing cycle.
	AccountLevelRemainingCredit *Float64 `json:"accountLevelRemainingCredit,omitempty" xmlrpc:"accountLevelRemainingCredit,omitempty"`

	// The active enrollments for this account.
	Enrollments []FlexibleCredit_Enrollment `json:"enrollments,omitempty" xmlrpc:"enrollments,omitempty"`

	// Indicates whether or not the account is participating in any account level Flexible Credit programs.
	IsAccountLevelParticipantFlag *bool `json:"isAccountLevelParticipantFlag,omitempty" xmlrpc:"isAccountLevelParticipantFlag,omitempty"`

	// Indicates whether or not the account is participating in any Flexible Credit programs.
	IsParticipantFlag *bool `json:"isParticipantFlag,omitempty" xmlrpc:"isParticipantFlag,omitempty"`

	// Indicates whether or not the account is participating in any product specific level Flexible Credit programs.
	IsProductSpecificParticipantFlag *bool `json:"isProductSpecificParticipantFlag,omitempty" xmlrpc:"isProductSpecificParticipantFlag,omitempty"`

	// The amount of credit that has been used by all product specific enrollments in the billing cycle.
	ProductSpecificAppliedCredit *Float64 `json:"productSpecificAppliedCredit,omitempty" xmlrpc:"productSpecificAppliedCredit,omitempty"`

	// Product specific credit allowance applied over the course of entire active program enrollments. For enrollments without a lifetime restriction, this property will not be populated as credit will be tracked on a purely monthly basis.
	ProductSpecificLifetimeAppliedCredit *Float64 `json:"productSpecificLifetimeAppliedCredit,omitempty" xmlrpc:"productSpecificLifetimeAppliedCredit,omitempty"`

	// The total product specific credit over the course of an entire program enrollment. This value may be null, in which case the enrollment credit is applied on a monthly basis and there is no lifetime maximum.
	ProductSpecificLifetimeCredit *Float64 `json:"productSpecificLifetimeCredit,omitempty" xmlrpc:"productSpecificLifetimeCredit,omitempty"`

	// Remaining product specific level credit allowance available over the remaining duration of the program enrollments. If null, enrollment credit is applied on a strictly monthly basis and there is no lifetime maximum. Enrollments with non-null remaining lifetime credit will receive the lesser of the remaining monthly credit or the remaining lifetime credit.
	ProductSpecificLifetimeRemainingCredit *Float64 `json:"productSpecificLifetimeRemainingCredit,omitempty" xmlrpc:"productSpecificLifetimeRemainingCredit,omitempty"`

	// The total product specific monthly credit allowance available at the beginning of a billing cycle.
	ProductSpecificMonthlyCredit *Float64 `json:"productSpecificMonthlyCredit,omitempty" xmlrpc:"productSpecificMonthlyCredit,omitempty"`

	// The total product specific credit allowance still available during the current billing cycle.
	ProductSpecificRemainingCredit *Float64 `json:"productSpecificRemainingCredit,omitempty" xmlrpc:"productSpecificRemainingCredit,omitempty"`

	// The credit allowance that has already been applied during the current billing cycle from all enrollments. If the lifetime limit has been or soon will be reached, this amount may included credit applied in previous billing cycles.
	TotalAppliedCredit *Float64 `json:"totalAppliedCredit,omitempty" xmlrpc:"totalAppliedCredit,omitempty"`

	// The credit allowance that is available during the current billing cycle from all enrollments. If the lifetime limit has been or soon will be reached, this amount may be reduced by credit applied in previous billing cycles.
	TotalRemainingCredit *Float64 `json:"totalRemainingCredit,omitempty" xmlrpc:"totalRemainingCredit,omitempty"`
}

// no documentation yet
type Container_Account_External_Setup_ProvisioningHoldLifted struct {
	Entity

	// no documentation yet
	AdditionalAttributes *Container_Account_External_Setup_ProvisioningHoldLifted_Attributes `json:"additionalAttributes,omitempty" xmlrpc:"additionalAttributes,omitempty"`

	// no documentation yet
	Code *string `json:"code,omitempty" xmlrpc:"code,omitempty"`

	// no documentation yet
	Error *string `json:"error,omitempty" xmlrpc:"error,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// no documentation yet
type Container_Account_External_Setup_ProvisioningHoldLifted_Attributes struct {
	Entity

	// no documentation yet
	BrandKeyName *string `json:"brandKeyName,omitempty" xmlrpc:"brandKeyName,omitempty"`

	// no documentation yet
	SoftLayerBrandMoveDate *Time `json:"softLayerBrandMoveDate,omitempty" xmlrpc:"softLayerBrandMoveDate,omitempty"`
}

// Historical Summary Container for account resource details
type Container_Account_Historical_Summary struct {
	Entity

	// Array of server uptime detail containers
	Details []Container_Account_Historical_Summary_Detail `json:"details,omitempty" xmlrpc:"details,omitempty"`

	// The maximum date included in the summary.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The minimum date included in the summary.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// Historical Summary Details Container for a resource's data
type Container_Account_Historical_Summary_Detail struct {
	Entity

	// The maximum date included in the detail.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The minimum date included in the detail.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// Historical Summary Details Container for a host resource uptime
type Container_Account_Historical_Summary_Detail_Uptime struct {
	Container_Account_Historical_Summary_Detail

	// The hardware for uptime details.
	CloudComputingInstance *Virtual_Guest `json:"cloudComputingInstance,omitempty" xmlrpc:"cloudComputingInstance,omitempty"`

	// The data associated with a host uptime details.
	Data []Metric_Tracking_Object_Data `json:"data,omitempty" xmlrpc:"data,omitempty"`

	// The hardware for uptime details.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`
}

// Historical Summary Container for account host's resource uptime details
type Container_Account_Historical_Summary_Uptime struct {
	Container_Account_Historical_Summary
}

// no documentation yet
type Container_Account_Internal_Ibm_CostRecovery struct {
	Entity

	// no documentation yet
	AccountId *string `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	CountryId *string `json:"countryId,omitempty" xmlrpc:"countryId,omitempty"`
}

// Contains data required to both request a new IaaS account for active IBM employees and review pending requests. Fields used exclusively in the review process are scrubbed of user input.
type Container_Account_Internal_Ibm_Request struct {
	Entity

	// Purpose of the internal IBM account chosen from the list of available
	AccountType *string `json:"accountType,omitempty" xmlrpc:"accountType,omitempty"`

	// If not provided, will attempt to retrieve from BluePages
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// If no address provided, will attempt to retrieve from BluePages
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// If not provided, will attempt to retrieve from BluePages
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Name of the company displayed on the IaaS account
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// no documentation yet
	CostRecoveryAccountId *string `json:"costRecoveryAccountId,omitempty" xmlrpc:"costRecoveryAccountId,omitempty"`

	// no documentation yet
	CostRecoveryCountryId *string `json:"costRecoveryCountryId,omitempty" xmlrpc:"costRecoveryCountryId,omitempty"`

	// If not provided, will attempt to retrieve from BluePages
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// True if the request has been denied by either the IaaS team or the
	DeniedFlag *bool `json:"deniedFlag,omitempty" xmlrpc:"deniedFlag,omitempty"`

	// Department within the division which will be changed during cost recovery. [DEPRECATED]
	// Deprecated: This function has been marked as deprecated.
	DepartmentCode *string `json:"departmentCode,omitempty" xmlrpc:"departmentCode,omitempty"`

	// Country code assigned to the department for cost recovery. [DEPRECATED]
	// Deprecated: This function has been marked as deprecated.
	DepartmentCountry *string `json:"departmentCountry,omitempty" xmlrpc:"departmentCountry,omitempty"`

	// Division code used for cost recovery. [DEPRECATED]
	// Deprecated: This function has been marked as deprecated.
	DivisionCode *string `json:"divisionCode,omitempty" xmlrpc:"divisionCode,omitempty"`

	// Account owner's IBM email address. Must be a discoverable email
	EmailAddress *string `json:"emailAddress,omitempty" xmlrpc:"emailAddress,omitempty"`

	// Applicant's first name, as provided by IBM BluePages API.
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// Applicant's last name, as provided by IBM BluePages API.
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// APPROVED if the request has been approved by the first-line manager,
	ManagerApprovalStatus *string `json:"managerApprovalStatus,omitempty" xmlrpc:"managerApprovalStatus,omitempty"`

	// True for accounts intended to be multi-tenant and false otherwise
	MultiTenantFlag *bool `json:"multiTenantFlag,omitempty" xmlrpc:"multiTenantFlag,omitempty"`

	// Account owner's primary phone number. If no phone number is available
	OfficePhone *string `json:"officePhone,omitempty" xmlrpc:"officePhone,omitempty"`

	// Bluemix PaaS 32 digit hexadecimal account id being automatically linked
	PaasAccountId *string `json:"paasAccountId,omitempty" xmlrpc:"paasAccountId,omitempty"`

	// If not provided, will attempt to retrieve from BluePages
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Stated purpose of the new account this request would create
	Purpose *string `json:"purpose,omitempty" xmlrpc:"purpose,omitempty"`

	// Division's security SME's email address, if available
	SecuritySubjectMatterExpertEmail *string `json:"securitySubjectMatterExpertEmail,omitempty" xmlrpc:"securitySubjectMatterExpertEmail,omitempty"`

	// Division's security SME's name, if available
	SecuritySubjectMatterExpertName *string `json:"securitySubjectMatterExpertName,omitempty" xmlrpc:"securitySubjectMatterExpertName,omitempty"`

	// Division's security SME's phone, if available
	SecuritySubjectMatterExpertPhone *string `json:"securitySubjectMatterExpertPhone,omitempty" xmlrpc:"securitySubjectMatterExpertPhone,omitempty"`

	// If required for chosen country and not provided, will attempt
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// no documentation yet
type Container_Account_Payment_Method_CreditCard struct {
	Entity

	// no documentation yet
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// no documentation yet
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	CurrencyShortName *string `json:"currencyShortName,omitempty" xmlrpc:"currencyShortName,omitempty"`

	// no documentation yet
	CybersourceAssignedCardType *string `json:"cybersourceAssignedCardType,omitempty" xmlrpc:"cybersourceAssignedCardType,omitempty"`

	// no documentation yet
	ExpireMonth *string `json:"expireMonth,omitempty" xmlrpc:"expireMonth,omitempty"`

	// no documentation yet
	ExpireYear *string `json:"expireYear,omitempty" xmlrpc:"expireYear,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastFourDigits *string `json:"lastFourDigits,omitempty" xmlrpc:"lastFourDigits,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	Nickname *string `json:"nickname,omitempty" xmlrpc:"nickname,omitempty"`

	// no documentation yet
	PaymentMethodRoleName *string `json:"paymentMethodRoleName,omitempty" xmlrpc:"paymentMethodRoleName,omitempty"`

	// no documentation yet
	PaymentTypeId *string `json:"paymentTypeId,omitempty" xmlrpc:"paymentTypeId,omitempty"`

	// no documentation yet
	PaymentTypeName *string `json:"paymentTypeName,omitempty" xmlrpc:"paymentTypeName,omitempty"`

	// no documentation yet
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// no documentation yet
type Container_Account_PersonalInformation struct {
	Entity

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// no documentation yet
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// no documentation yet
	AlternatePhone *string `json:"alternatePhone,omitempty" xmlrpc:"alternatePhone,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	OfficePhone *string `json:"officePhone,omitempty" xmlrpc:"officePhone,omitempty"`

	// no documentation yet
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// no documentation yet
	RequestDate *Time `json:"requestDate,omitempty" xmlrpc:"requestDate,omitempty"`

	// no documentation yet
	RequestId *int `json:"requestId,omitempty" xmlrpc:"requestId,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// The customer and prospective owner of a proof of concept account established by an IBMer.
type Container_Account_ProofOfConcept_Contact_Customer struct {
	Entity

	// Customer's address
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// Customer's address
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// Customer's city
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Customer's ISO country code
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// Customer's email address
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// Customer's first name
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// Customer's last name
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// Customer's primary phone number
	Phone *string `json:"phone,omitempty" xmlrpc:"phone,omitempty"`

	// Customer's postal code
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Customer's state
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// Customer's VAT ID
	VatId *string `json:"vatId,omitempty" xmlrpc:"vatId,omitempty"`
}

// IBMer who is submitting a proof of concept request on behalf of a prospective customer.
type Container_Account_ProofOfConcept_Contact_Ibmer_Requester struct {
	Entity

	// Customer's address
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// Customer's address
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// no documentation yet
	BusinessUnit *string `json:"businessUnit,omitempty" xmlrpc:"businessUnit,omitempty"`

	// Customer's city
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Customer's ISO country code
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// Customer's email address
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// Customer's first name
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// Customer's last name
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	OrganizationCountry *string `json:"organizationCountry,omitempty" xmlrpc:"organizationCountry,omitempty"`

	// no documentation yet
	PaasAccountId *string `json:"paasAccountId,omitempty" xmlrpc:"paasAccountId,omitempty"`

	// Customer's primary phone number
	Phone *string `json:"phone,omitempty" xmlrpc:"phone,omitempty"`

	// Customer's postal code
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Customer's state
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	SubOrganization *string `json:"subOrganization,omitempty" xmlrpc:"subOrganization,omitempty"`

	// no documentation yet
	Uid *string `json:"uid,omitempty" xmlrpc:"uid,omitempty"`

	// Customer's VAT ID
	VatId *string `json:"vatId,omitempty" xmlrpc:"vatId,omitempty"`
}

// IBMer who will assist the requester with technical aspects of configuring the proof of concept account.
type Container_Account_ProofOfConcept_Contact_Ibmer_Technical struct {
	Entity

	// Customer's address
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// Customer's address
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// Customer's city
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Customer's ISO country code
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// Customer's email address
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// Customer's first name
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// Customer's last name
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// Customer's primary phone number
	Phone *string `json:"phone,omitempty" xmlrpc:"phone,omitempty"`

	// Customer's postal code
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Customer's state
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	Uid *string `json:"uid,omitempty" xmlrpc:"uid,omitempty"`

	// Customer's VAT ID
	VatId *string `json:"vatId,omitempty" xmlrpc:"vatId,omitempty"`
}

// Proof of concept request using the account team funding model. Note that proof of concept account request are available only to internal IBM employees.
type Container_Account_ProofOfConcept_Request_AccountFunded struct {
	Container_Account_ProofOfConcept_Request_GlobalFunded

	// Billing codes for the department paying for the proof of concept account
	CostRecoveryRequest *Container_Account_ProofOfConcept_Request_CostRecovery `json:"costRecoveryRequest,omitempty" xmlrpc:"costRecoveryRequest,omitempty"`
}

// Funding codes for the department paying for the proof of concept account.
type Container_Account_ProofOfConcept_Request_CostRecovery struct {
	Entity

	// Internal billing system country code
	CountryCode *string `json:"countryCode,omitempty" xmlrpc:"countryCode,omitempty"`

	// Customer's Internal billing system department code
	DepartmentCode *string `json:"departmentCode,omitempty" xmlrpc:"departmentCode,omitempty"`

	// Internal billing system division code
	DivisionCode *string `json:"divisionCode,omitempty" xmlrpc:"divisionCode,omitempty"`
}

// Proof of concept request using the global funding model. Note that proof of concept account request are available only to internal IBM employees.
type Container_Account_ProofOfConcept_Request_GlobalFunded struct {
	Entity

	// Dollar amount of funding requested for the proof of concept period
	Amount *Float64 `json:"amount,omitempty" xmlrpc:"amount,omitempty"`

	// Customer intended to take over ownership and and billing of the account
	Customer *Container_Account_ProofOfConcept_Contact_Customer `json:"customer,omitempty" xmlrpc:"customer,omitempty"`

	// Explanation of the purpose of the proof of concept request
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// End date for the proof of concept period
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// Internal opportunity system details
	Opportunity *Container_Account_ProofOfConcept_Request_Opportunity `json:"opportunity,omitempty" xmlrpc:"opportunity,omitempty"`

	// Name of the project or company and will become the account companyName
	ProjectName *string `json:"projectName,omitempty" xmlrpc:"projectName,omitempty"`

	// IBM region responsible for overseeing the proof of concept account
	RegionKeyName *string `json:"regionKeyName,omitempty" xmlrpc:"regionKeyName,omitempty"`

	// IBMer requesting the proof of concept account
	Requester *Container_Account_ProofOfConcept_Contact_Ibmer_Requester `json:"requester,omitempty" xmlrpc:"requester,omitempty"`

	// Start date for the proof of concept period
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`

	// IBMer assisting with technical aspects of account configuration
	TechnicalContact *Container_Account_ProofOfConcept_Contact_Ibmer_Technical `json:"technicalContact,omitempty" xmlrpc:"technicalContact,omitempty"`
}

// Internal IBM opportunity codes required when applying for a Proof of Concept account.
type Container_Account_ProofOfConcept_Request_Opportunity struct {
	Entity

	// The campaign or promotion code for this request, provided by Sales.
	CampaignCode *string `json:"campaignCode,omitempty" xmlrpc:"campaignCode,omitempty"`

	// Expected monthly revenue.
	MonthlyRecurringRevenue *Float64 `json:"monthlyRecurringRevenue,omitempty" xmlrpc:"monthlyRecurringRevenue,omitempty"`

	// Internal system identifier.
	OpportunityNumber *string `json:"opportunityNumber,omitempty" xmlrpc:"opportunityNumber,omitempty"`

	// Expected overall contract value.
	TotalContractValue *Float64 `json:"totalContractValue,omitempty" xmlrpc:"totalContractValue,omitempty"`
}

// Full details presented to reviewers when determining whether or not to accept a proof of concept request. Note that reviewers are internal IBM employees and reviews are not exposed to external users.
type Container_Account_ProofOfConcept_Review struct {
	Entity

	// Type of brand the account will use
	AccountType *string `json:"accountType,omitempty" xmlrpc:"accountType,omitempty"`

	// Internal billing codes
	CostRecoveryCodes *Container_Account_ProofOfConcept_Request_CostRecovery `json:"costRecoveryCodes,omitempty" xmlrpc:"costRecoveryCodes,omitempty"`

	// Customer intended to take over billing after the proof of concept period
	Customer *Container_Account_ProofOfConcept_Contact_Customer `json:"customer,omitempty" xmlrpc:"customer,omitempty"`

	// Describes the purpose and rationale of the request
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Expected end date of the proof of concept period
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// Dollar amount of funding requested
	FundingAmount *Float64 `json:"fundingAmount,omitempty" xmlrpc:"fundingAmount,omitempty"`

	// Funding option chosen for the request
	FundingType *string `json:"fundingType,omitempty" xmlrpc:"fundingType,omitempty"`

	// System id of the request
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Name of the integrated offering team lead reviewing the request
	IotLeadName *string `json:"iotLeadName,omitempty" xmlrpc:"iotLeadName,omitempty"`

	// Name of the integrated offering team region
	IotRegionName *string `json:"iotRegionName,omitempty" xmlrpc:"iotRegionName,omitempty"`

	// Name of requesting IBMer's manager
	ManagerName *string `json:"managerName,omitempty" xmlrpc:"managerName,omitempty"`

	// Internal opportunity tracking information
	Opportunity *Container_Account_ProofOfConcept_Request_Opportunity `json:"opportunity,omitempty" xmlrpc:"opportunity,omitempty"`

	// Project name chosen by the requesting IBMer
	ProjectName *string `json:"projectName,omitempty" xmlrpc:"projectName,omitempty"`

	// IBMer requesting the account on behalf of a customer
	Requester *Container_Account_ProofOfConcept_Contact_Ibmer_Requester `json:"requester,omitempty" xmlrpc:"requester,omitempty"`

	// Summary of request's review activity
	ReviewHistory *Container_Account_ProofOfConcept_Review_History `json:"reviewHistory,omitempty" xmlrpc:"reviewHistory,omitempty"`

	// URL for the individual review
	ReviewUrl *string `json:"reviewUrl,omitempty" xmlrpc:"reviewUrl,omitempty"`

	// Expected start date of the proof of concept period
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`

	// Additional IBMer responsible for configuring the cloud capabilities
	TechnicalContact *Container_Account_ProofOfConcept_Contact_Ibmer_Technical `json:"technicalContact,omitempty" xmlrpc:"technicalContact,omitempty"`
}

// Review event within proof of concept request review period.
type Container_Account_ProofOfConcept_Review_Event struct {
	Entity

	// Explanation of the event.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Reviewer's email address.
	ReviewerEmail *string `json:"reviewerEmail,omitempty" xmlrpc:"reviewerEmail,omitempty"`

	// Reviewer's BluePages UID.
	ReviewerUid *string `json:"reviewerUid,omitempty" xmlrpc:"reviewerUid,omitempty"`
}

// Summary of review activity for a proof of concept request.
type Container_Account_ProofOfConcept_Review_History struct {
	Entity

	// True for approved requests associated with a new account and false otherwise.
	AccountCreatedFlag *bool `json:"accountCreatedFlag,omitempty" xmlrpc:"accountCreatedFlag,omitempty"`

	// True for denied requests and false otherwise.
	DeniedFlag *bool `json:"deniedFlag,omitempty" xmlrpc:"deniedFlag,omitempty"`

	// List of events occurring during the review.
	Events []Container_Account_ProofOfConcept_Review_Event `json:"events,omitempty" xmlrpc:"events,omitempty"`

	// True for fully reviewed requests and false otherwise.
	ReviewCompleteFlag *bool `json:"reviewCompleteFlag,omitempty" xmlrpc:"reviewCompleteFlag,omitempty"`
}

// Summary presented to reviewers when determining whether or not to accept a proof of concept request. Note that reviewers are internal IBM employees and reviews are not exposed to external users.
type Container_Account_ProofOfConcept_Review_Summary struct {
	Entity

	// Account's companyName
	AccountName *string `json:"accountName,omitempty" xmlrpc:"accountName,omitempty"`

	// Current account owner
	AccountOwnerName *string `json:"accountOwnerName,omitempty" xmlrpc:"accountOwnerName,omitempty"`

	// Dollar amount requested
	Amount *Float64 `json:"amount,omitempty" xmlrpc:"amount,omitempty"`

	// Date the request was submitted
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Email of the customer receiving the proof of concept account
	CustomerEmail *string `json:"customerEmail,omitempty" xmlrpc:"customerEmail,omitempty"`

	// Name of the customer receiving the proof of concept account
	CustomerName *string `json:"customerName,omitempty" xmlrpc:"customerName,omitempty"`

	// Request record's id
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Date of the last state change on the request
	LastUpdate *Time `json:"lastUpdate,omitempty" xmlrpc:"lastUpdate,omitempty"`

	// Email address of the reviewer, if any, currently reviewing the request
	NextApproverEmail *string `json:"nextApproverEmail,omitempty" xmlrpc:"nextApproverEmail,omitempty"`

	// Email address of the requester
	RequesterEmail *string `json:"requesterEmail,omitempty" xmlrpc:"requesterEmail,omitempty"`

	// Requesting IBMer's full name
	RequesterName *string `json:"requesterName,omitempty" xmlrpc:"requesterName,omitempty"`

	// URL for the individual review
	ReviewUrl *string `json:"reviewUrl,omitempty" xmlrpc:"reviewUrl,omitempty"`

	// Request's current status (Pending, Denied, or Approved)
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// Contains data related to an account after editing its information.
type Container_Account_Update_Response struct {
	Entity

	// Whether or not the update was accepted and applied.
	AcceptedFlag *bool `json:"acceptedFlag,omitempty" xmlrpc:"acceptedFlag,omitempty"`

	// The updated SoftLayer_Account.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// If a manual review is required, this will be populated with the SoftLayer_Ticket for that review.
	Ticket *Ticket `json:"ticket,omitempty" xmlrpc:"ticket,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_Common data type contains common information for requests to the getPortalLogin API. This is an abstract class that serves as a base that more specialized classes will derive from. For example, a request class specific to SoftLayer Native IMS Login (username and password).
type Container_Authentication_Request_Common struct {
	Container_Authentication_Request_Contract

	// The answer to your security question.
	SecurityQuestionAnswer *string `json:"securityQuestionAnswer,omitempty" xmlrpc:"securityQuestionAnswer,omitempty"`

	// A security question you wish to answer when authenticating to the SoftLayer customer portal. This parameter isn't required if no security questions are set on your portal account or if your account is configured to not require answering a security account upon login.
	SecurityQuestionId *int `json:"securityQuestionId,omitempty" xmlrpc:"securityQuestionId,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_Contract provides a common set of operations for implementing classes.
type Container_Authentication_Request_Contract struct {
	Entity
}

// The SoftLayer_Container_Authentication_Request_Native data type contains information for requests to the getPortalLogin API. This class is specific to the SoftLayer Native login (username/password). The request information will be verified to ensure it is valid, and then there will be an attempt to obtain a portal login token in authenticating the user with the provided information.
type Container_Authentication_Request_Native struct {
	Container_Authentication_Request_Common

	// no documentation yet
	AuxiliaryClaimsMiniToken *string `json:"auxiliaryClaimsMiniToken,omitempty" xmlrpc:"auxiliaryClaimsMiniToken,omitempty"`

	// Your SoftLayer customer portal user's portal password.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// The username you wish to authenticate to the SoftLayer customer portal with.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_Native_External data type contains information for requests to the getPortalLogin API. This class serves as a base class for more specialized external authentication classes to the SoftLayer Native login (username/password).
type Container_Authentication_Request_Native_External struct {
	Container_Authentication_Request_Native
}

// The SoftLayer_Container_Authentication_Request_Native_External_Totp data type contains information for requests to the getPortalLogin API. This class provides information to allow the user to submit a request to the native SoftLayer (username/password) login service for a portal login token, as well as submitting a request to the TOTP 2 factor authentication service.
type Container_Authentication_Request_Native_External_Totp struct {
	Container_Authentication_Request_Native_External

	// no documentation yet
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// no documentation yet
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`

	// no documentation yet
	Vendor *string `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_Native_External_Verisign data type contains information for requests to the getPortalLogin API. This class provides information to allow the user to submit a request to the native SoftLayer (username/password) login service for a portal login token, as well as submitting a request to the Verisign 2 factor authentication service.
type Container_Authentication_Request_Native_External_Verisign struct {
	Container_Authentication_Request_Native_External

	// no documentation yet
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// no documentation yet
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`

	// no documentation yet
	Vendor *string `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_OpenIdConnect data type contains information for requests to the getPortalLogin API. This class is specific to the SoftLayer Cloud Token login. The request information will be verified to ensure it is valid, and then there will be an attempt to obtain a portal login token in authenticating the user with the provided information.
type Container_Authentication_Request_OpenIdConnect struct {
	Container_Authentication_Request_Common

	// no documentation yet
	OpenIdConnectAccessToken *string `json:"openIdConnectAccessToken,omitempty" xmlrpc:"openIdConnectAccessToken,omitempty"`

	// no documentation yet
	OpenIdConnectAccountId *int `json:"openIdConnectAccountId,omitempty" xmlrpc:"openIdConnectAccountId,omitempty"`

	// no documentation yet
	OpenIdConnectProvider *string `json:"openIdConnectProvider,omitempty" xmlrpc:"openIdConnectProvider,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_OpenIdConnect_External data type contains information for requests to the getPortalLogin API. This class serves as a base class for more specialized external authentication classes to the SoftLayer OpenIdConnect login service.
type Container_Authentication_Request_OpenIdConnect_External struct {
	Container_Authentication_Request_OpenIdConnect
}

// The SoftLayer_Container_Authentication_Request_OpenIdConnect_External_Totp data type contains information for requests to the getPortalLogin API. This class provides information to allow the user to submit a request to the SoftLayer OpenIdConnect (token) login service for a portal login token, as well as submitting a request to the TOTP 2 factor authentication service.
type Container_Authentication_Request_OpenIdConnect_External_Totp struct {
	Container_Authentication_Request_OpenIdConnect_External

	// no documentation yet
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// no documentation yet
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`

	// no documentation yet
	Vendor *string `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`
}

// The SoftLayer_Container_Authentication_Request_OpenIdConnect_External_Verisign data type contains information for requests to the getPortalLogin API. This class provides information to allow the user to submit a request to the SoftLayer OpenIdConnect (token) login service for a portal login token, as well as submitting a request to the Verisign 2 factor authentication service.
type Container_Authentication_Request_OpenIdConnect_External_Verisign struct {
	Container_Authentication_Request_OpenIdConnect_External

	// no documentation yet
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// no documentation yet
	SecurityCode *int `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`

	// no documentation yet
	Vendor *string `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_2FactorAuthenticationNeeded data type contains information for specific responses from the getPortalLogin API. This class is indicative of a request that is missing the appropriate 2FA information.
type Container_Authentication_Response_2FactorAuthenticationNeeded struct {
	Container_Authentication_Response_Common

	// no documentation yet
	AdditionalData *Container_Authentication_Response_Common `json:"additionalData,omitempty" xmlrpc:"additionalData,omitempty"`

	// no documentation yet
	StatusKeyName *string `json:"statusKeyName,omitempty" xmlrpc:"statusKeyName,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_Account data type contains account information for responses from the getPortalLogin API.
type Container_Authentication_Response_Account struct {
	Entity

	// no documentation yet
	AccountCompanyName *string `json:"accountCompanyName,omitempty" xmlrpc:"accountCompanyName,omitempty"`

	// no documentation yet
	AccountCountry *string `json:"accountCountry,omitempty" xmlrpc:"accountCountry,omitempty"`

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	AccountStatusName *string `json:"accountStatusName,omitempty" xmlrpc:"accountStatusName,omitempty"`

	// no documentation yet
	BluemixAccountId *string `json:"bluemixAccountId,omitempty" xmlrpc:"bluemixAccountId,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	DefaultAccount *bool `json:"defaultAccount,omitempty" xmlrpc:"defaultAccount,omitempty"`

	// no documentation yet
	IpAddressCheckRequired *bool `json:"ipAddressCheckRequired,omitempty" xmlrpc:"ipAddressCheckRequired,omitempty"`

	// no documentation yet
	IsMasterUserFlag *bool `json:"isMasterUserFlag,omitempty" xmlrpc:"isMasterUserFlag,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	SecurityQuestionRequired *bool `json:"securityQuestionRequired,omitempty" xmlrpc:"securityQuestionRequired,omitempty"`

	// no documentation yet
	TotpExternalAuthenticationRequired *bool `json:"totpExternalAuthenticationRequired,omitempty" xmlrpc:"totpExternalAuthenticationRequired,omitempty"`

	// no documentation yet
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`

	// no documentation yet
	VerisignExternalAuthenticationRequired *bool `json:"verisignExternalAuthenticationRequired,omitempty" xmlrpc:"verisignExternalAuthenticationRequired,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_AccountIdMissing data type contains information for specific responses from the getPortalLogin API. This class is indicative of a request that is missing the account id.
type Container_Authentication_Response_AccountIdMissing struct {
	Container_Authentication_Response_Common

	// no documentation yet
	StatusKeyName *string `json:"statusKeyName,omitempty" xmlrpc:"statusKeyName,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_Common data type contains common information for responses from the getPortalLogin API. This is an abstract class that serves as a base that more specialized classes will derive from. For example, a response class that is specific to a successful response from the getPortalLogin API.
type Container_Authentication_Response_Common struct {
	Entity

	// The list of linked accounts for the authenticated SoftLayer customer portal user.
	Accounts []Container_Authentication_Response_Account `json:"accounts,omitempty" xmlrpc:"accounts,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_IpAddressRestrictionCheckNeeded data type indicates that the caller (IAM presumably) needs to do an IP address check of the logging-in user against the restricted IP list kept in BSS.  We don't know the IP address of the user here (only IAM does) so we return an indicator of which user matched the username and expect IAM to come back with another login call that will include a mini-JWT token that contains an assertion that the IP address was checked.
type Container_Authentication_Response_IpAddressRestrictionCheckNeeded struct {
	Container_Authentication_Response_Common

	// no documentation yet
	StatusKeyName *string `json:"statusKeyName,omitempty" xmlrpc:"statusKeyName,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_LOGIN_FAILED data type contains information for specific responses from the getPortalLogin API. This class is indicative of a request where there was an inability to login based on the information that was provided.
type Container_Authentication_Response_LoginFailed struct {
	Container_Authentication_Response_Common

	// no documentation yet
	ErrorMessage *string `json:"errorMessage,omitempty" xmlrpc:"errorMessage,omitempty"`

	// no documentation yet
	StatusKeyName *string `json:"statusKeyName,omitempty" xmlrpc:"statusKeyName,omitempty"`
}

// The SoftLayer_Container_Authentication_Response_SUCCESS data type contains information for specific responses from the getPortalLogin API. This class is indicative of a request that was successful in obtaining a portal login token from the getPortalLogin API.
type Container_Authentication_Response_Success struct {
	Container_Authentication_Response_Common

	// no documentation yet
	StatusKeyName *string `json:"statusKeyName,omitempty" xmlrpc:"statusKeyName,omitempty"`

	// The token for interacting with the SoftLayer customer portal.
	Token *Container_User_Authentication_Token `json:"token,omitempty" xmlrpc:"token,omitempty"`
}

// no documentation yet
type Container_Auxiliary_Network_Status_Reading struct {
	Entity

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	AveragePing *Float64 `json:"averagePing,omitempty" xmlrpc:"averagePing,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Fails *int `json:"fails,omitempty" xmlrpc:"fails,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Frequency *int `json:"frequency,omitempty" xmlrpc:"frequency,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Label *string `json:"label,omitempty" xmlrpc:"label,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	LastCheckDate *Time `json:"lastCheckDate,omitempty" xmlrpc:"lastCheckDate,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	LastDownDate *Time `json:"lastDownDate,omitempty" xmlrpc:"lastDownDate,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Latency *Float64 `json:"latency,omitempty" xmlrpc:"latency,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	MaximumPing *Float64 `json:"maximumPing,omitempty" xmlrpc:"maximumPing,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	MinimumPing *Float64 `json:"minimumPing,omitempty" xmlrpc:"minimumPing,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	PingLoss *Float64 `json:"pingLoss,omitempty" xmlrpc:"pingLoss,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	StatusCode *string `json:"statusCode,omitempty" xmlrpc:"statusCode,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	StatusMessage *string `json:"statusMessage,omitempty" xmlrpc:"statusMessage,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	Target *string `json:"target,omitempty" xmlrpc:"target,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	TargetType *string `json:"targetType,omitempty" xmlrpc:"targetType,omitempty"`
}

// SoftLayer_Container_Bandwidth_GraphInputs models a single inbound object for a given bandwidth graph.
type Container_Bandwidth_GraphInputs struct {
	Entity

	// This is a unix timestamp that represents the stop date/time for a graph.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The front-end or back-end network uplink interface associated with this server.
	NetworkInterfaceId *int `json:"networkInterfaceId,omitempty" xmlrpc:"networkInterfaceId,omitempty"`

	// *
	Pod *int `json:"pod,omitempty" xmlrpc:"pod,omitempty"`

	// This is a human readable name for the server or rack being graphed.
	ServerName *string `json:"serverName,omitempty" xmlrpc:"serverName,omitempty"`

	// This is a unix timestamp that represents the begin date/time for a graph.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// SoftLayer_Container_Bandwidth_GraphOutputs models a single outbound object for a given bandwidth graph.
type Container_Bandwidth_GraphOutputs struct {
	Entity

	// The raw PNG binary data to be displayed once the graph is drawn.
	GraphImage *[]byte `json:"graphImage,omitempty" xmlrpc:"graphImage,omitempty"`

	// The title that ended up being displayed as part of the graph image.
	GraphTitle *string `json:"graphTitle,omitempty" xmlrpc:"graphTitle,omitempty"`

	// The maximum date included in this graph.
	MaxEndDate *Time `json:"maxEndDate,omitempty" xmlrpc:"maxEndDate,omitempty"`

	// The minimum date included in this graph.
	MinStartDate *Time `json:"minStartDate,omitempty" xmlrpc:"minStartDate,omitempty"`
}

// SoftLayer_Container_Bandwidth_Projection models projected bandwidth use over a time range.
type Container_Bandwidth_Projection struct {
	Entity

	// Bandwidth limit for this hardware.
	AllowedUsage *string `json:"allowedUsage,omitempty" xmlrpc:"allowedUsage,omitempty"`

	// Estimated bandwidth usage so far this billing cycle.
	EstimatedUsage *string `json:"estimatedUsage,omitempty" xmlrpc:"estimatedUsage,omitempty"`

	// Hardware ID of server to monitor.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Projected usage for this hardware based on previous usage this billing cycle.
	ProjectedUsage *string `json:"projectedUsage,omitempty" xmlrpc:"projectedUsage,omitempty"`

	// the text name of the server being monitored.
	ServerName *string `json:"serverName,omitempty" xmlrpc:"serverName,omitempty"`

	// The minimum date included in this list.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// When a customer uses SoftLayer_Account::getBandwidthUsage, this container is used to return their usage information in bytes
type Container_Bandwidth_Usage struct {
	Entity

	// no documentation yet
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// no documentation yet
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// no documentation yet
	PrivateInUsage *Float64 `json:"privateInUsage,omitempty" xmlrpc:"privateInUsage,omitempty"`

	// no documentation yet
	PrivateOutUsage *Float64 `json:"privateOutUsage,omitempty" xmlrpc:"privateOutUsage,omitempty"`

	// no documentation yet
	PublicInUsage *Float64 `json:"publicInUsage,omitempty" xmlrpc:"publicInUsage,omitempty"`

	// no documentation yet
	PublicOutUsage *Float64 `json:"publicOutUsage,omitempty" xmlrpc:"publicOutUsage,omitempty"`

	// no documentation yet
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// no documentation yet
type Container_Billing_Currency_Country struct {
	Entity

	// no documentation yet
	AvailableCurrencies []Billing_Currency `json:"availableCurrencies,omitempty" xmlrpc:"availableCurrencies,omitempty"`

	// no documentation yet
	Country *Locale_Country `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	CurrencyCountryLocales []Billing_Currency_Country `json:"currencyCountryLocales,omitempty" xmlrpc:"currencyCountryLocales,omitempty"`
}

// no documentation yet
type Container_Billing_Currency_Format struct {
	Entity

	// no documentation yet
	Currency *string `json:"currency,omitempty" xmlrpc:"currency,omitempty"`

	// no documentation yet
	Display *int `json:"display,omitempty" xmlrpc:"display,omitempty"`

	// no documentation yet
	Format *string `json:"format,omitempty" xmlrpc:"format,omitempty"`

	// no documentation yet
	Locale *string `json:"locale,omitempty" xmlrpc:"locale,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Position *int `json:"position,omitempty" xmlrpc:"position,omitempty"`

	// no documentation yet
	Precision *int `json:"precision,omitempty" xmlrpc:"precision,omitempty"`

	// no documentation yet
	Script *string `json:"script,omitempty" xmlrpc:"script,omitempty"`

	// no documentation yet
	Service *string `json:"service,omitempty" xmlrpc:"service,omitempty"`

	// no documentation yet
	Symbol *string `json:"symbol,omitempty" xmlrpc:"symbol,omitempty"`

	// no documentation yet
	Tag *string `json:"tag,omitempty" xmlrpc:"tag,omitempty"`

	// no documentation yet
	Value *Float64 `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Container_Billing_Info_Ach struct {
	Entity

	// no documentation yet
	AccountNumber *string `json:"accountNumber,omitempty" xmlrpc:"accountNumber,omitempty"`

	// no documentation yet
	AccountType *string `json:"accountType,omitempty" xmlrpc:"accountType,omitempty"`

	// no documentation yet
	BankTransitNumber *string `json:"bankTransitNumber,omitempty" xmlrpc:"bankTransitNumber,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	FederalTaxId *string `json:"federalTaxId,omitempty" xmlrpc:"federalTaxId,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	PhoneNumber *string `json:"phoneNumber,omitempty" xmlrpc:"phoneNumber,omitempty"`

	// no documentation yet
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	Street1 *string `json:"street1,omitempty" xmlrpc:"street1,omitempty"`

	// no documentation yet
	Street2 *string `json:"street2,omitempty" xmlrpc:"street2,omitempty"`
}

// This container is used to provide all the options for [[SoftLayer_Billing_Invoice/emailInvoices|emailInvoices]] in order to have the necessary invoices generated and links sent to the user's email.
type Container_Billing_Invoice_Email struct {
	Entity

	// Excel Invoices to email
	ExcelInvoiceIds []int `json:"excelInvoiceIds,omitempty" xmlrpc:"excelInvoiceIds,omitempty"`

	// PDF Invoice Details to email
	PdfDetailedInvoiceIds []int `json:"pdfDetailedInvoiceIds,omitempty" xmlrpc:"pdfDetailedInvoiceIds,omitempty"`

	// PDF Invoices to email
	PdfInvoiceIds []int `json:"pdfInvoiceIds,omitempty" xmlrpc:"pdfInvoiceIds,omitempty"`

	// The type of Invoices to be emailed [current|next]. If next is selected, the account id will be used.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// SoftLayer_Container_Billing_Order_Status models an order status.
type Container_Billing_Order_Status struct {
	Entity

	// The description of the status.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The keyname of the status.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// Contains user information used to request a manual Catalyst enrollment.
type Container_Catalyst_ManualEnrollmentRequest struct {
	Entity

	// Applicant's email address
	CustomerEmail *string `json:"customerEmail,omitempty" xmlrpc:"customerEmail,omitempty"`

	// Applicant's first and last name
	CustomerName *string `json:"customerName,omitempty" xmlrpc:"customerName,omitempty"`

	// Name of applicant's startup company
	StartupName *string `json:"startupName,omitempty" xmlrpc:"startupName,omitempty"`

	// Flag indicating whether (true) or not (false) and applicant is
	VentureAffiliationFlag *bool `json:"ventureAffiliationFlag,omitempty" xmlrpc:"ventureAffiliationFlag,omitempty"`

	// Name of the venture capital fund, if any, applicant is affiliated with
	VentureFundName *string `json:"ventureFundName,omitempty" xmlrpc:"ventureFundName,omitempty"`
}

// This container is used to hold country locale information.
type Container_Collection_Locale_CountryCode struct {
	Entity

	// no documentation yet
	LongName *string `json:"longName,omitempty" xmlrpc:"longName,omitempty"`

	// no documentation yet
	ShortName *string `json:"shortName,omitempty" xmlrpc:"shortName,omitempty"`

	// no documentation yet
	StateCodes []Container_Collection_Locale_StateCode `json:"stateCodes,omitempty" xmlrpc:"stateCodes,omitempty"`
}

// This container is used to hold information regarding a state or province.
type Container_Collection_Locale_StateCode struct {
	Entity

	// no documentation yet
	LongName *string `json:"longName,omitempty" xmlrpc:"longName,omitempty"`

	// no documentation yet
	ShortName *string `json:"shortName,omitempty" xmlrpc:"shortName,omitempty"`
}

// This container is used to hold VAT information.
type Container_Collection_Locale_VatCountryCodeAndFormat struct {
	Entity

	// no documentation yet
	CountryCode *string `json:"countryCode,omitempty" xmlrpc:"countryCode,omitempty"`

	// no documentation yet
	Regex *string `json:"regex,omitempty" xmlrpc:"regex,omitempty"`
}

// no documentation yet
type Container_Disk_Image_Capture_Template struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Summary *string `json:"summary,omitempty" xmlrpc:"summary,omitempty"`

	// no documentation yet
	Volumes []Container_Disk_Image_Capture_Template_Volume `json:"volumes,omitempty" xmlrpc:"volumes,omitempty"`
}

// no documentation yet
type Container_Disk_Image_Capture_Template_Volume struct {
	Entity

	// A customer provided flag to indicate that the current volume is the boot drive
	BootVolumeFlag *bool `json:"bootVolumeFlag,omitempty" xmlrpc:"bootVolumeFlag,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Partitions []Container_Disk_Image_Capture_Template_Volume_Partition `json:"partitions,omitempty" xmlrpc:"partitions,omitempty"`

	// The storage group to capture
	StorageGroupId *int `json:"storageGroupId,omitempty" xmlrpc:"storageGroupId,omitempty"`
}

// no documentation yet
type Container_Disk_Image_Capture_Template_Volume_Partition struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Container_Exception data type represents a SoftLayer_Exception.
type Container_Exception struct {
	Entity

	// The SoftLayer_Exception class that the error is.
	ExceptionClass *string `json:"exceptionClass,omitempty" xmlrpc:"exceptionClass,omitempty"`

	// The exception message.
	ExceptionMessage *string `json:"exceptionMessage,omitempty" xmlrpc:"exceptionMessage,omitempty"`
}

// no documentation yet
type Container_Graph struct {
	Entity

	// base units associated with the graph.
	BaseUnit *string `json:"baseUnit,omitempty" xmlrpc:"baseUnit,omitempty"`

	// Graph range end datetime.
	EndDatetime *string `json:"endDatetime,omitempty" xmlrpc:"endDatetime,omitempty"`

	// The height of the graph image.
	Height *int `json:"height,omitempty" xmlrpc:"height,omitempty"`

	// The graph image.
	Image *[]byte `json:"image,omitempty" xmlrpc:"image,omitempty"`

	// The graph interval in seconds.
	Interval *int `json:"interval,omitempty" xmlrpc:"interval,omitempty"`

	// Metric types associated with the graph.
	Metrics []Container_Metric_Data_Type `json:"metrics,omitempty" xmlrpc:"metrics,omitempty"`

	// Indicator to control whether the graph data is normalized.
	NormalizeFlag *[]byte `json:"normalizeFlag,omitempty" xmlrpc:"normalizeFlag,omitempty"`

	// The options used to control the graph appearance.
	Options []Container_Graph_Option `json:"options,omitempty" xmlrpc:"options,omitempty"`

	// A collection of graph plots.
	Plots []Container_Graph_Plot `json:"plots,omitempty" xmlrpc:"plots,omitempty"`

	// Graph range start datetime.
	StartDatetime *string `json:"startDatetime,omitempty" xmlrpc:"startDatetime,omitempty"`

	// The name of the template to use; may be null.
	Template *string `json:"template,omitempty" xmlrpc:"template,omitempty"`

	// The title of the graph image.
	Title *string `json:"title,omitempty" xmlrpc:"title,omitempty"`

	// The width of the graph image.
	Width *int `json:"width,omitempty" xmlrpc:"width,omitempty"`
}

// no documentation yet
type Container_Graph_Option struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Container_Graph_Plot struct {
	Entity

	// no documentation yet
	Data []Container_Graph_Plot_Coordinate `json:"data,omitempty" xmlrpc:"data,omitempty"`

	// no documentation yet
	Metric *Container_Metric_Data_Type `json:"metric,omitempty" xmlrpc:"metric,omitempty"`

	// no documentation yet
	Unit *string `json:"unit,omitempty" xmlrpc:"unit,omitempty"`
}

// no documentation yet
type Container_Graph_Plot_Coordinate struct {
	Entity

	// no documentation yet
	XValue *Float64 `json:"xValue,omitempty" xmlrpc:"xValue,omitempty"`

	// no documentation yet
	YValue *Float64 `json:"yValue,omitempty" xmlrpc:"yValue,omitempty"`

	// no documentation yet
	ZValue *Float64 `json:"zValue,omitempty" xmlrpc:"zValue,omitempty"`
}

// no documentation yet
type Container_Hardware_CaptureEnabled struct {
	Entity

	// no documentation yet
	Enabled *bool `json:"enabled,omitempty" xmlrpc:"enabled,omitempty"`

	// no documentation yet
	Reasons []string `json:"reasons,omitempty" xmlrpc:"reasons,omitempty"`
}

// The hardware configuration container is used to provide configuration options for servers.
//
// Each configuration option will include both an <code>itemPrice</code> and a <code>template</code>.
//
// The <code>itemPrice</code> value will provide hourly and monthly costs (if either are applicable), and a description of the option.
//
// The <code>template</code> will provide a fragment of the request with the properties and values that must be sent when creating a server with the option.
//
// The [[SoftLayer_Hardware/getCreateObjectOptions|getCreateObjectOptions]] method returns this data structure.
//
// <style type="text/css">#properties .views-field-body p { margin-top: 1.5em; };</style>
type Container_Hardware_Configuration struct {
	Entity

	//
	// <div style="width: 200%">
	// Available datacenter options.
	//
	//
	// The <code>datacenter.name</code> value in the template represents which datacenter the server will be provisioned in.
	// </div>
	Datacenters []Container_Hardware_Configuration_Option `json:"datacenters,omitempty" xmlrpc:"datacenters,omitempty"`

	//
	// <div style="width: 200%">
	// Available fixed configuration preset options.
	//
	//
	// The <code>fixedConfigurationPreset.keyName</code> value in the template is an identifier for a particular fixed configuration. When provided exactly as shown in the template, that fixed configuration will be used.
	//
	//
	// When providing a <code>fixedConfigurationPreset.keyName</code> while ordering a server the <code>processors</code> and <code>hardDrives</code> configuration options cannot be used.
	// </div>
	FixedConfigurationPresets []Container_Hardware_Configuration_Option `json:"fixedConfigurationPresets,omitempty" xmlrpc:"fixedConfigurationPresets,omitempty"`

	//
	// <div style="width: 200%">
	// Available hard drive options.
	//
	//
	// A server will have at least one hard drive.
	//
	//
	// The <code>hardDrives.capacity</code> value in the template represents the size, in gigabytes, of the disk.
	// </div>
	HardDrives []Container_Hardware_Configuration_Option `json:"hardDrives,omitempty" xmlrpc:"hardDrives,omitempty"`

	//
	// <div style="width: 200%">
	// Available network component options.
	//
	//
	// The <code>networkComponent.maxSpeed</code> value in the template represents the link speed, in megabits per second, of the network connections for a server.
	// </div>
	NetworkComponents []Container_Hardware_Configuration_Option `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`

	//
	// <div style="width: 200%">
	// Available operating system options.
	//
	//
	// The <code>operatingSystemReferenceCode</code> value in the template is an identifier for a particular operating system. When provided exactly as shown in the template, that operating system will be used.
	//
	//
	// A reference code is structured as three tokens separated by underscores. The first token represents the product, the second is the version of the product, and the third is whether the OS is 32 or 64bit.
	//
	//
	// When providing an <code>operatingSystemReferenceCode</code> while ordering a server the only token required to match exactly is the product. The version token may be given as 'LATEST', else it will require an exact match as well. When the bits token is not provided, 64 bits will be assumed.
	//
	//
	// Providing the value of 'LATEST' for a version will select the latest release of that product for the operating system. As this may change over time, you should be sure that the release version is irrelevant for your applications.
	//
	//
	// For Windows based operating systems the version will represent both the release version (2008, 2012, etc) and the edition (Standard, Enterprise, etc). For all other operating systems the version will represent the major version (Centos 6, Ubuntu 12, etc) of that operating system, minor versions are represented in few reference codes where they are significant.
	// </div>
	OperatingSystems []Container_Hardware_Configuration_Option `json:"operatingSystems,omitempty" xmlrpc:"operatingSystems,omitempty"`

	//
	// <div style="width: 200%">
	// Available processor options.
	//
	//
	// The <code>processorCoreAmount</code> value in the template represents the number of cores allocated to the server.
	// The <code>memoryCapacity</code> value in the template represents the amount of memory, in gigabytes, allocated to the server.
	// </div>
	Processors []Container_Hardware_Configuration_Option `json:"processors,omitempty" xmlrpc:"processors,omitempty"`
}

// An option found within a [[SoftLayer_Container_Hardware_Configuration (type)]] structure.
type Container_Hardware_Configuration_Option struct {
	Entity

	//
	// Provides hourly and monthly costs (if either are applicable), and a description of the option.
	ItemPrice *Product_Item_Price `json:"itemPrice,omitempty" xmlrpc:"itemPrice,omitempty"`

	//
	// Provides a description of a fixed configuration preset with monthly and hourly costs.
	Preset *Product_Package_Preset `json:"preset,omitempty" xmlrpc:"preset,omitempty"`

	//
	// Provides a fragment of the request with the properties and values that must be sent when creating a server with the option.
	Template *Hardware `json:"template,omitempty" xmlrpc:"template,omitempty"`
}

// no documentation yet
type Container_Hardware_DiskImageMap struct {
	Entity

	// no documentation yet
	BootFlag *int `json:"bootFlag,omitempty" xmlrpc:"bootFlag,omitempty"`

	// no documentation yet
	DiskImageUUID *string `json:"diskImageUUID,omitempty" xmlrpc:"diskImageUUID,omitempty"`

	// no documentation yet
	DiskSerialNumber *string `json:"diskSerialNumber,omitempty" xmlrpc:"diskSerialNumber,omitempty"`
}

// no documentation yet
type Container_Hardware_MassUpdate struct {
	Entity

	// The hardwares updated by the mass update tool
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Errors encountered while mass updating hardwares
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// The hardwares that failed to update
	SuccessFlag *string `json:"successFlag,omitempty" xmlrpc:"successFlag,omitempty"`
}

// no documentation yet
type Container_Hardware_Pool_Details struct {
	Entity

	// no documentation yet
	PendingOrders *int `json:"pendingOrders,omitempty" xmlrpc:"pendingOrders,omitempty"`

	// no documentation yet
	PendingTransactions *int `json:"pendingTransactions,omitempty" xmlrpc:"pendingTransactions,omitempty"`

	// no documentation yet
	PoolDescription *string `json:"poolDescription,omitempty" xmlrpc:"poolDescription,omitempty"`

	// no documentation yet
	PoolKeyName *string `json:"poolKeyName,omitempty" xmlrpc:"poolKeyName,omitempty"`

	// no documentation yet
	PoolName *string `json:"poolName,omitempty" xmlrpc:"poolName,omitempty"`

	// no documentation yet
	Routers []Container_Hardware_Pool_Details_Router `json:"routers,omitempty" xmlrpc:"routers,omitempty"`

	// no documentation yet
	TotalHardware *int `json:"totalHardware,omitempty" xmlrpc:"totalHardware,omitempty"`

	// no documentation yet
	TotalInventoryHardware *int `json:"totalInventoryHardware,omitempty" xmlrpc:"totalInventoryHardware,omitempty"`

	// no documentation yet
	TotalProvisionedHardware *int `json:"totalProvisionedHardware,omitempty" xmlrpc:"totalProvisionedHardware,omitempty"`

	// no documentation yet
	TotalTestedHardware *int `json:"totalTestedHardware,omitempty" xmlrpc:"totalTestedHardware,omitempty"`

	// no documentation yet
	TotalTestingHardware *int `json:"totalTestingHardware,omitempty" xmlrpc:"totalTestingHardware,omitempty"`
}

// no documentation yet
type Container_Hardware_Pool_Details_Router struct {
	Entity

	// no documentation yet
	PoolThreshold *int `json:"poolThreshold,omitempty" xmlrpc:"poolThreshold,omitempty"`

	// no documentation yet
	RouterId *int `json:"routerId,omitempty" xmlrpc:"routerId,omitempty"`

	// no documentation yet
	RouterName *string `json:"routerName,omitempty" xmlrpc:"routerName,omitempty"`

	// no documentation yet
	TotalHardware *int `json:"totalHardware,omitempty" xmlrpc:"totalHardware,omitempty"`

	// no documentation yet
	TotalInventoryHardware *int `json:"totalInventoryHardware,omitempty" xmlrpc:"totalInventoryHardware,omitempty"`

	// no documentation yet
	TotalProvisionedHardware *int `json:"totalProvisionedHardware,omitempty" xmlrpc:"totalProvisionedHardware,omitempty"`

	// no documentation yet
	TotalTestedHardware *int `json:"totalTestedHardware,omitempty" xmlrpc:"totalTestedHardware,omitempty"`

	// no documentation yet
	TotalTestingHardware *int `json:"totalTestingHardware,omitempty" xmlrpc:"totalTestingHardware,omitempty"`
}

// The SoftLayer_Container_Hardware_Server_Configuration data type contains information relating to a server's item price information, and hard drive partition information.
type Container_Hardware_Server_Configuration struct {
	Entity

	// A flag indicating that the server will be moved into the spare pool after an Operating system reload.
	AddToSparePoolAfterOsReload *int `json:"addToSparePoolAfterOsReload,omitempty" xmlrpc:"addToSparePoolAfterOsReload,omitempty"`

	// The customer provision uri will be used to download and execute a customer defined script on the host at the end of provisioning.
	CustomProvisionScriptUri *string `json:"customProvisionScriptUri,omitempty" xmlrpc:"customProvisionScriptUri,omitempty"`

	// A flag indicating that the primary drive will be converted to a portable storage volume during an Operating System reload.
	DriveRetentionFlag *bool `json:"driveRetentionFlag,omitempty" xmlrpc:"driveRetentionFlag,omitempty"`

	// A flag indicating that all data will be erased from drives during an Operating System reload.
	EraseHardDrives *int `json:"eraseHardDrives,omitempty" xmlrpc:"eraseHardDrives,omitempty"`

	// The hard drive partitions that a server can be partitioned with.
	HardDrives []Hardware_Component `json:"hardDrives,omitempty" xmlrpc:"hardDrives,omitempty"`

	// An Image Template ID [[SoftLayer_Virtual_Guest_Block_Device_Template_Group]] that will be deployed to the host.  If provided no item prices are required.
	ImageTemplateId *int `json:"imageTemplateId,omitempty" xmlrpc:"imageTemplateId,omitempty"`

	// Whether the OS reload will be in-place for accounts that support it.
	InPlaceFlag *bool `json:"inPlaceFlag,omitempty" xmlrpc:"inPlaceFlag,omitempty"`

	// The item prices that a server can be configured with.
	ItemPrices []Product_Item_Price `json:"itemPrices,omitempty" xmlrpc:"itemPrices,omitempty"`

	// A flag indicating that the provision should use LVM for all logical drives.
	LvmFlag *int `json:"lvmFlag,omitempty" xmlrpc:"lvmFlag,omitempty"`

	// A flag indicating that the remote management cards password will be reset.
	ResetIpmiPassword *int `json:"resetIpmiPassword,omitempty" xmlrpc:"resetIpmiPassword,omitempty"`

	// The token of the requesting service. Do not set.
	ServiceToken *string `json:"serviceToken,omitempty" xmlrpc:"serviceToken,omitempty"`

	// IDs to SoftLayer_Security_Ssh_Key objects on the current account which will be added to the server for authentication. SSH Keys will not be added to servers with Microsoft Windows.
	SshKeyIds []int `json:"sshKeyIds,omitempty" xmlrpc:"sshKeyIds,omitempty"`

	// A flag indicating that the BIOS will be updated when installing the operating system.
	UpgradeBios *int `json:"upgradeBios,omitempty" xmlrpc:"upgradeBios,omitempty"`

	// A flag indicating that the firmware on all hard drives will be updated when installing the operating system.
	UpgradeHardDriveFirmware *int `json:"upgradeHardDriveFirmware,omitempty" xmlrpc:"upgradeHardDriveFirmware,omitempty"`
}

// The SoftLayer_Container_Hardware_Server_Details data type contains information relating to a server's component information, network information, and software information.
type Container_Hardware_Server_Details struct {
	Entity

	// The components that belong to a piece of hardware.
	Components []Hardware_Component `json:"components,omitempty" xmlrpc:"components,omitempty"`

	// The network components that belong to a piece of hardware.
	NetworkComponents []Network_Component `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`

	// The software that belong to a piece of hardware.
	Software []Software_Component `json:"software,omitempty" xmlrpc:"software,omitempty"`
}

// no documentation yet
type Container_Hardware_Server_Request struct {
	Entity

	// no documentation yet
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// no documentation yet
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// no documentation yet
	SuccessFlag *bool `json:"successFlag,omitempty" xmlrpc:"successFlag,omitempty"`
}

// no documentation yet
type Container_Image_StorageGroupDetails struct {
	Entity

	// no documentation yet
	Drives []Container_Image_StorageGroupDetails_Drives `json:"drives,omitempty" xmlrpc:"drives,omitempty"`

	// no documentation yet
	StorageGroupName *string `json:"storageGroupName,omitempty" xmlrpc:"storageGroupName,omitempty"`

	// no documentation yet
	StorageGroupType *string `json:"storageGroupType,omitempty" xmlrpc:"storageGroupType,omitempty"`
}

// no documentation yet
type Container_Image_StorageGroupDetails_Drives struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	DiskSpace *string `json:"diskSpace,omitempty" xmlrpc:"diskSpace,omitempty"`

	// no documentation yet
	Units *string `json:"units,omitempty" xmlrpc:"units,omitempty"`
}

// SoftLayer_Container_KnowledgeLayer_QuestionAnswer models a single question and answer pair from SoftLayer's KnowledgeLayer knowledge base. SoftLayer's backend network interfaces with the KnowledgeLayer to recommend helpful articles when support tickets are created.
type Container_KnowledgeLayer_QuestionAnswer struct {
	Entity

	// The answer to a question asked on the SoftLayer KnowledgeLayer.
	Answer *string `json:"answer,omitempty" xmlrpc:"answer,omitempty"`

	// The link to a question asked on the SoftLayer KnowledgeLayer.
	Link *string `json:"link,omitempty" xmlrpc:"link,omitempty"`

	// A question asked on the SoftLayer KnowledgeLayer.
	Question *string `json:"question,omitempty" xmlrpc:"question,omitempty"`
}

// no documentation yet
type Container_Message struct {
	Entity

	// no documentation yet
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// no documentation yet
type Container_Metric_Data_Type struct {
	Entity

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	SummaryType *string `json:"summaryType,omitempty" xmlrpc:"summaryType,omitempty"`

	// no documentation yet
	Unit *string `json:"unit,omitempty" xmlrpc:"unit,omitempty"`
}

// SoftLayer_Container_Metric_Tracking_Object_Details This container is a parent class for detailing diverse metrics.
type Container_Metric_Tracking_Object_Details struct {
	Entity

	// The name that best describes the metric being collected.
	MetricName *string `json:"metricName,omitempty" xmlrpc:"metricName,omitempty"`
}

// SoftLayer_Container_Metric_Tracking_Object_Summary This container is a parent class for summarizing diverse metrics.
type Container_Metric_Tracking_Object_Summary struct {
	Entity

	// The name that best describes the metric being collected.
	MetricName *string `json:"metricName,omitempty" xmlrpc:"metricName,omitempty"`
}

// SoftLayer_Container_Metric_Tracking_Object_Virtual_Host_Details This container details a virtual host's metric data.
type Container_Metric_Tracking_Object_Virtual_Host_Details struct {
	Container_Metric_Tracking_Object_Details

	// The day this metric was collected.
	Day *Time `json:"day,omitempty" xmlrpc:"day,omitempty"`

	// The maximum number of guests hosted by this platform for the given day.
	MaxInstances *int `json:"maxInstances,omitempty" xmlrpc:"maxInstances,omitempty"`

	// The maximum amount of memory utilized by this platform for the given day.
	MaxMemoryUsage *int `json:"maxMemoryUsage,omitempty" xmlrpc:"maxMemoryUsage,omitempty"`

	// The mean number of guests hosted by this platform for the given day.
	MeanInstances *Float64 `json:"meanInstances,omitempty" xmlrpc:"meanInstances,omitempty"`

	// The mean amount of memory utilized by this platform for the given day.
	MeanMemoryUsage *Float64 `json:"meanMemoryUsage,omitempty" xmlrpc:"meanMemoryUsage,omitempty"`

	// The minimum number of guests hosted by this platform for the given day.
	MinInstances *int `json:"minInstances,omitempty" xmlrpc:"minInstances,omitempty"`

	// The minimum amount of memory utilized by this platform for the given day.
	MinMemoryUsage *int `json:"minMemoryUsage,omitempty" xmlrpc:"minMemoryUsage,omitempty"`
}

// SoftLayer_Container_Metric_Tracking_Object_Virtual_Host_Summary This container summarizes a virtual host's metric data.
type Container_Metric_Tracking_Object_Virtual_Host_Summary struct {
	Container_Metric_Tracking_Object_Summary

	// The average amount of memory usage thus far in this billing cycle.
	AvgMemoryUsageInBillingCycle *int `json:"avgMemoryUsageInBillingCycle,omitempty" xmlrpc:"avgMemoryUsageInBillingCycle,omitempty"`

	// Current bill cycle end date.
	CurrentBillCycleEnd *Time `json:"currentBillCycleEnd,omitempty" xmlrpc:"currentBillCycleEnd,omitempty"`

	// Current bill cycle start date.
	CurrentBillCycleStart *Time `json:"currentBillCycleStart,omitempty" xmlrpc:"currentBillCycleStart,omitempty"`

	// The last count of instances this platform was hosting.
	LastInstanceCount *int `json:"lastInstanceCount,omitempty" xmlrpc:"lastInstanceCount,omitempty"`

	// The last amount of memory this platform was using.
	LastMemoryUsageAmount *int `json:"lastMemoryUsageAmount,omitempty" xmlrpc:"lastMemoryUsageAmount,omitempty"`

	// The last time this virtual host was polled for metrics.
	LastPollTime *Time `json:"lastPollTime,omitempty" xmlrpc:"lastPollTime,omitempty"`

	// The max number of instances hosted thus far in this billing cycle.
	MaxInstanceInBillingCycle *int `json:"maxInstanceInBillingCycle,omitempty" xmlrpc:"maxInstanceInBillingCycle,omitempty"`

	// Previous bill cycle end date.
	PreviousBillCycleEnd *Time `json:"previousBillCycleEnd,omitempty" xmlrpc:"previousBillCycleEnd,omitempty"`

	// Previous bill cycle start date.
	PreviousBillCycleStart *Time `json:"previousBillCycleStart,omitempty" xmlrpc:"previousBillCycleStart,omitempty"`

	// This virtual hosting platform name.
	VirtualPlatformName *string `json:"virtualPlatformName,omitempty" xmlrpc:"virtualPlatformName,omitempty"`
}

// The SoftLayer_Container_Monitoring_Alarm_History data type contains information relating to SoftLayer monitoring alarm history.
type Container_Monitoring_Alarm_History struct {
	Entity

	// Account ID that this alarm belongs to
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// DEPRECATED. ID of the monitoring agent that triggered this alarm
	// Deprecated: This function has been marked as deprecated.
	AgentId *int `json:"agentId,omitempty" xmlrpc:"agentId,omitempty"`

	// Alarm ID
	AlarmId *string `json:"alarmId,omitempty" xmlrpc:"alarmId,omitempty"`

	// Time that an alarm was closed.
	ClosedDate *Time `json:"closedDate,omitempty" xmlrpc:"closedDate,omitempty"`

	// Time that an alarm was triggered
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Alarm message
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// DEPRECATED. Robot ID
	// Deprecated: This function has been marked as deprecated.
	RobotId *int `json:"robotId,omitempty" xmlrpc:"robotId,omitempty"`

	// Severity of an alarm
	Severity *string `json:"severity,omitempty" xmlrpc:"severity,omitempty"`
}

// This object holds authentication data to a server.
type Container_Network_Authentication_Data struct {
	Entity

	// The name of a host
	Host *string `json:"host,omitempty" xmlrpc:"host,omitempty"`

	// The authentication password
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// The port number
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// The type of network protocol. This can be ftp, ssh and so on.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The authentication username
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// SoftLayer_Container_Network_Bandwidth_Data_Summary models an interface's overall bandwidth usage during it's current billing cycle.
type Container_Network_Bandwidth_Data_Summary struct {
	Entity

	// The amount of bandwidth a server has allocated to it in it's current billing period.
	AllowedUsage *Float64 `json:"allowedUsage,omitempty" xmlrpc:"allowedUsage,omitempty"`

	// The amount of bandwidth that a server has used within it's current billing period.
	EstimatedUsage *Float64 `json:"estimatedUsage,omitempty" xmlrpc:"estimatedUsage,omitempty"`

	// The amount of bandwidth a server is projected to use within its billing period, based on it's current usage.
	ProjectedUsage *Float64 `json:"projectedUsage,omitempty" xmlrpc:"projectedUsage,omitempty"`

	// The unit of measurement used in a bandwidth data summary.
	UsageUnits *string `json:"usageUnits,omitempty" xmlrpc:"usageUnits,omitempty"`
}

// SoftLayer_Container_Network_Bandwidth_Version1_Usage models an hourly bandwidth record.
type Container_Network_Bandwidth_Version1_Usage struct {
	Entity

	// The amount of incoming bandwidth that a server has used within the hour of the recordedDate.
	IncomingAmount *Float64 `json:"incomingAmount,omitempty" xmlrpc:"incomingAmount,omitempty"`

	// The amount of outgoing bandwidth that a server has used within the hour of the recordedDate.
	OutgoingAmount *Float64 `json:"outgoingAmount,omitempty" xmlrpc:"outgoingAmount,omitempty"`

	// The date and time that the bandwidth was used by a piece of hardware
	RecordedDate *Time `json:"recordedDate,omitempty" xmlrpc:"recordedDate,omitempty"`
}

// The SoftLayer_Container_Network_CdnMarketplace_Configuration_Behavior_ModifyResponseHeader data type contains information for specific responses from the modify response header API.
type Container_Network_CdnMarketplace_Configuration_Behavior_ModifyResponseHeader struct {
	Entity

	// Specifies the delimiter to be used when indicating multiple values for a header. Valid delimiter is, a <space>, , (comma), ; (semicolon), ,<space> (comma and space), or ;<space> (semicolon and space).
	Delimiter *string `json:"delimiter,omitempty" xmlrpc:"delimiter,omitempty"`

	// The description of modify response header.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A collection of key value pairs that specify the headers and associated values to be modified. The header name and header value must be separated by colon (:). Example: ['header1:value1','header2:Value2']
	Headers []string `json:"headers,omitempty" xmlrpc:"headers,omitempty"`

	// The uniqueId of the modify response header to which the existing behavior belongs.
	ModResHeaderUniqueId *string `json:"modResHeaderUniqueId,omitempty" xmlrpc:"modResHeaderUniqueId,omitempty"`

	// The path, relative to the domain that is accessed via modify response header.
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// The type of the modify response header, could be append/modify/delete. Set this to append to add a given header value to a header name set in the headerList. Set this to delete to remove a given header value from a header name set in the headerList. Set this to overwrite to match on a specified header name and replace its existing header value with a new one you specify.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The uniqueId of the mapping to which the existing behavior belongs.
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`
}

// The SoftLayer_Container_Network_CdnMarketplace_Configuration_Behavior_TokenAuth data type contains information for specific responses from the Token Authentication API.
type Container_Network_CdnMarketplace_Configuration_Behavior_TokenAuth struct {
	Entity

	// Specifies a single character to separate access control list (ACL) fields. The default value is '!'.
	AclDelimiter *string `json:"aclDelimiter,omitempty" xmlrpc:"aclDelimiter,omitempty"`

	// Possible values '0' and '1'. If set to '1', input values are escaped before adding them to the token. Default value is '1'.
	EscapeTokenInputs *string `json:"escapeTokenInputs,omitempty" xmlrpc:"escapeTokenInputs,omitempty"`

	// Specifies the algorithm to use for the token's hash-based message authentication code (HMAC) field. Valid entries are 'SHA256', 'SHA1', or 'MD5'. The default value is 'SHA256'.
	HmacAlgorithm *string `json:"hmacAlgorithm,omitempty" xmlrpc:"hmacAlgorithm,omitempty"`

	// Possible values '0' and '1'. If set to '1', query strings are removed from a URL when computing the token's HMAC algorithm. Default value is '0'.
	IgnoreQueryString *string `json:"ignoreQueryString,omitempty" xmlrpc:"ignoreQueryString,omitempty"`

	// The token name. If this value is empty, then it is set to the default value '__token__'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The path, relative to the domain that is accessed via token authentication.
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// Specifies a single character to separate the individual token fields. The default value is '~'.
	TokenDelimiter *string `json:"tokenDelimiter,omitempty" xmlrpc:"tokenDelimiter,omitempty"`

	// The token encryption key, which specifies an even number of hex digits for the token key. An entry can be up to 64 characters in length.
	TokenKey *string `json:"tokenKey,omitempty" xmlrpc:"tokenKey,omitempty"`

	// The token transition key, which specifies an even number of hex digits for the token transition key. An entry can be up to 64 characters in length.
	TransitionKey *string `json:"transitionKey,omitempty" xmlrpc:"transitionKey,omitempty"`

	// The uniqueId of the mapping to which the existing behavior belongs.
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Configuration_Cache_Purge struct {
	Entity

	// no documentation yet
	Date *string `json:"date,omitempty" xmlrpc:"date,omitempty"`

	// no documentation yet
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// no documentation yet
	Saved *string `json:"saved,omitempty" xmlrpc:"saved,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// The SoftLayer_Container_Network_CdnMarketplace_Configuration_Cache_PurgeGroup data type contains information for specific responses from the Purge Group API. Each of the Purge Group APIs returns a collection of this type
type Container_Network_CdnMarketplace_Configuration_Cache_PurgeGroup struct {
	Entity

	// Date in which record is created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A identifier that is unique to purge group.
	GroupUniqueId *string `json:"groupUniqueId,omitempty" xmlrpc:"groupUniqueId,omitempty"`

	// The Unix timestamp of the last purge.
	LastPurgeDate *Time `json:"lastPurgeDate,omitempty" xmlrpc:"lastPurgeDate,omitempty"`

	// Purge Group name. The favorite group name must be unique, but non-favorite groups do not have this limitation
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The following options are available to create a Purge Group: option 1: only purge the paths in the group, but don't save as favorite. option 2: only save the purge group as favorite, but don't purge paths. option 3: save the purge group as favorite and also purge paths.
	Option *int `json:"option,omitempty" xmlrpc:"option,omitempty"`

	// Total number of purge paths.
	PathCount *int `json:"pathCount,omitempty" xmlrpc:"pathCount,omitempty"`

	// A collection of purge paths.
	Paths []string `json:"paths,omitempty" xmlrpc:"paths,omitempty"`

	// The purge's status when the input option field is 1 or 3. Status can be SUCCESS, FAILED, or IN_PROGRESS.
	PurgeStatus *string `json:"purgeStatus,omitempty" xmlrpc:"purgeStatus,omitempty"`

	// Type of the Purge Group, currently SAVED or UNSAVED.
	Saved *string `json:"saved,omitempty" xmlrpc:"saved,omitempty"`

	// A identifier that is unique to domain mapping.
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`
}

// The SoftLayer_Container_Network_CdnMarketplace_Configuration_Cache_PurgeGroupHistory data type contains information for specific responses from the Purge Group API and Purge History API.
type Container_Network_CdnMarketplace_Configuration_Cache_PurgeGroupHistory struct {
	Entity

	// Date in which record is created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Purge Group name. The favorite group name must be unique, but un-favorite groups do not have this limitation
	GroupName *string `json:"groupName,omitempty" xmlrpc:"groupName,omitempty"`

	// Purge group unique ID
	GroupUniqueId *string `json:"groupUniqueId,omitempty" xmlrpc:"groupUniqueId,omitempty"`

	// The purge's status. Status can be SUCCESS, FAILED, or IN_PROGRESS.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Domain mapping unique ID.
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Configuration_Input struct {
	Entity

	// no documentation yet
	BucketName *string `json:"bucketName,omitempty" xmlrpc:"bucketName,omitempty"`

	// no documentation yet
	CacheKeyQueryRule *string `json:"cacheKeyQueryRule,omitempty" xmlrpc:"cacheKeyQueryRule,omitempty"`

	// no documentation yet
	CertificateType *string `json:"certificateType,omitempty" xmlrpc:"certificateType,omitempty"`

	// no documentation yet
	Cname *string `json:"cname,omitempty" xmlrpc:"cname,omitempty"`

	// no documentation yet
	Domain *string `json:"domain,omitempty" xmlrpc:"domain,omitempty"`

	// no documentation yet
	DynamicContentAcceleration *Container_Network_CdnMarketplace_Configuration_Performance_DynamicContentAcceleration `json:"dynamicContentAcceleration,omitempty" xmlrpc:"dynamicContentAcceleration,omitempty"`

	// no documentation yet
	FileExtension *string `json:"fileExtension,omitempty" xmlrpc:"fileExtension,omitempty"`

	// no documentation yet
	GeoblockingRule *Network_CdnMarketplace_Configuration_Behavior_Geoblocking `json:"geoblockingRule,omitempty" xmlrpc:"geoblockingRule,omitempty"`

	// no documentation yet
	Header *string `json:"header,omitempty" xmlrpc:"header,omitempty"`

	// no documentation yet
	HotlinkProtection *Network_CdnMarketplace_Configuration_Behavior_HotlinkProtection `json:"hotlinkProtection,omitempty" xmlrpc:"hotlinkProtection,omitempty"`

	// no documentation yet
	HttpPort *int `json:"httpPort,omitempty" xmlrpc:"httpPort,omitempty"`

	// no documentation yet
	HttpsPort *int `json:"httpsPort,omitempty" xmlrpc:"httpsPort,omitempty"`

	// Used by the following method: updateOriginPath(). This property will store the path of the path record to be saved. The $path attribute stores the new path.
	OldPath *string `json:"oldPath,omitempty" xmlrpc:"oldPath,omitempty"`

	// no documentation yet
	Origin *string `json:"origin,omitempty" xmlrpc:"origin,omitempty"`

	// no documentation yet
	OriginType *string `json:"originType,omitempty" xmlrpc:"originType,omitempty"`

	// no documentation yet
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// no documentation yet
	PerformanceConfiguration *string `json:"performanceConfiguration,omitempty" xmlrpc:"performanceConfiguration,omitempty"`

	// no documentation yet
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// no documentation yet
	RespectHeaders *string `json:"respectHeaders,omitempty" xmlrpc:"respectHeaders,omitempty"`

	// no documentation yet
	ServeStale *string `json:"serveStale,omitempty" xmlrpc:"serveStale,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// no documentation yet
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`

	// no documentation yet
	VendorName *string `json:"vendorName,omitempty" xmlrpc:"vendorName,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Configuration_Mapping struct {
	Entity

	// no documentation yet
	AkamaiCname *string `json:"akamaiCname,omitempty" xmlrpc:"akamaiCname,omitempty"`

	// no documentation yet
	BucketName *string `json:"bucketName,omitempty" xmlrpc:"bucketName,omitempty"`

	// no documentation yet
	CacheKeyQueryRule *string `json:"cacheKeyQueryRule,omitempty" xmlrpc:"cacheKeyQueryRule,omitempty"`

	// no documentation yet
	CertificateType *string `json:"certificateType,omitempty" xmlrpc:"certificateType,omitempty"`

	// no documentation yet
	Cname *string `json:"cname,omitempty" xmlrpc:"cname,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Domain *string `json:"domain,omitempty" xmlrpc:"domain,omitempty"`

	// no documentation yet
	DynamicContentAcceleration *Container_Network_CdnMarketplace_Configuration_Performance_DynamicContentAcceleration `json:"dynamicContentAcceleration,omitempty" xmlrpc:"dynamicContentAcceleration,omitempty"`

	// no documentation yet
	FileExtension *string `json:"fileExtension,omitempty" xmlrpc:"fileExtension,omitempty"`

	// no documentation yet
	Header *string `json:"header,omitempty" xmlrpc:"header,omitempty"`

	// no documentation yet
	HttpPort *int `json:"httpPort,omitempty" xmlrpc:"httpPort,omitempty"`

	// no documentation yet
	HttpsChallengeRedirectUrl *string `json:"httpsChallengeRedirectUrl,omitempty" xmlrpc:"httpsChallengeRedirectUrl,omitempty"`

	// no documentation yet
	HttpsChallengeResponse *string `json:"httpsChallengeResponse,omitempty" xmlrpc:"httpsChallengeResponse,omitempty"`

	// no documentation yet
	HttpsChallengeUrl *string `json:"httpsChallengeUrl,omitempty" xmlrpc:"httpsChallengeUrl,omitempty"`

	// no documentation yet
	HttpsPort *int `json:"httpsPort,omitempty" xmlrpc:"httpsPort,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	OriginHost *string `json:"originHost,omitempty" xmlrpc:"originHost,omitempty"`

	// no documentation yet
	OriginType *string `json:"originType,omitempty" xmlrpc:"originType,omitempty"`

	// no documentation yet
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// no documentation yet
	PerformanceConfiguration *string `json:"performanceConfiguration,omitempty" xmlrpc:"performanceConfiguration,omitempty"`

	// no documentation yet
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// no documentation yet
	RespectHeaders *bool `json:"respectHeaders,omitempty" xmlrpc:"respectHeaders,omitempty"`

	// no documentation yet
	ServeStale *bool `json:"serveStale,omitempty" xmlrpc:"serveStale,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// no documentation yet
	UniqueId *string `json:"uniqueId,omitempty" xmlrpc:"uniqueId,omitempty"`

	// no documentation yet
	VendorName *string `json:"vendorName,omitempty" xmlrpc:"vendorName,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Configuration_Mapping_Path struct {
	Entity

	// no documentation yet
	BucketName *string `json:"bucketName,omitempty" xmlrpc:"bucketName,omitempty"`

	// no documentation yet
	CacheKeyQueryRule *string `json:"cacheKeyQueryRule,omitempty" xmlrpc:"cacheKeyQueryRule,omitempty"`

	// no documentation yet
	DynamicContentAcceleration *Container_Network_CdnMarketplace_Configuration_Performance_DynamicContentAcceleration `json:"dynamicContentAcceleration,omitempty" xmlrpc:"dynamicContentAcceleration,omitempty"`

	// no documentation yet
	FileExtension *string `json:"fileExtension,omitempty" xmlrpc:"fileExtension,omitempty"`

	// no documentation yet
	Header *string `json:"header,omitempty" xmlrpc:"header,omitempty"`

	// no documentation yet
	HttpPort *int `json:"httpPort,omitempty" xmlrpc:"httpPort,omitempty"`

	// no documentation yet
	HttpsPort *int `json:"httpsPort,omitempty" xmlrpc:"httpsPort,omitempty"`

	// no documentation yet
	MappingUniqueId *string `json:"mappingUniqueId,omitempty" xmlrpc:"mappingUniqueId,omitempty"`

	// no documentation yet
	Origin *string `json:"origin,omitempty" xmlrpc:"origin,omitempty"`

	// no documentation yet
	OriginType *string `json:"originType,omitempty" xmlrpc:"originType,omitempty"`

	// no documentation yet
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// no documentation yet
	PerformanceConfiguration *string `json:"performanceConfiguration,omitempty" xmlrpc:"performanceConfiguration,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Configuration_Performance_DynamicContentAcceleration struct {
	Entity

	// The detectionPath is used by CDN edge servers to find the best optimized route from edge to the origin server. The Akamai edge servers fetch the test object from the origin to know the network condition to your origin server, and then calculate the best optimized route with the network condition. The best path to origin must be known at the time a user’s request arrives at an edge server, since any in-line analysis or probing would defeat the purpose of speeding things up.
	DetectionPath *string `json:"detectionPath,omitempty" xmlrpc:"detectionPath,omitempty"`

	// Serving compressed images reduces the amount of content required to load a page. This feature helps offset less robust connections, such as those formed with mobile devices. Basically, if your site visitors have slow network speeds, MobileImageCompression technology can automatically increase compression of JPEG images to speed up loading. On the other hand, this feature results in lossy compression or irreversible compression, and may affect the quality of the images on your site.
	//
	// JPG supported file extensions: .jpg, .jpeg, .jpe, .jig, .jgig, .jgi The default is enabled.
	MobileImageCompressionEnabled *bool `json:"mobileImageCompressionEnabled,omitempty" xmlrpc:"mobileImageCompressionEnabled,omitempty"`

	// Inspects HTML responses and prefetches embedded objects in HTML files. Prefetching works on any page that includes <img>, <script>, or <link> tags that specify relative paths. It also works when the resource hostname matches the request domain in the HTML file, and it is part of a fully qualified URI. When set to true, edge servers prefetch objects with the following file extensions:
	//
	// aif, aiff, au, avi, bin, bmp, cab, carb, cct, cdf, class, css, doc, dcr, dtd, exe, flv, gcf, gff, gif, grv, hdml, hqx, ico, ini, jpeg, jpg, js, mov, mp3, nc, pct, pdf, png, ppc, pws, swa, swf, txt, vbs, w32, wav, wbmp, wml, wmlc, wmls, wmlsc, xsd, and zip.
	//
	// The default is enabled.
	PrefetchEnabled *bool `json:"prefetchEnabled,omitempty" xmlrpc:"prefetchEnabled,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Metrics struct {
	Entity

	// no documentation yet
	Descriptions []string `json:"descriptions,omitempty" xmlrpc:"descriptions,omitempty"`

	// no documentation yet
	Names []string `json:"names,omitempty" xmlrpc:"names,omitempty"`

	// no documentation yet
	Percentage []string `json:"percentage,omitempty" xmlrpc:"percentage,omitempty"`

	// no documentation yet
	Time []int `json:"time,omitempty" xmlrpc:"time,omitempty"`

	// no documentation yet
	Totals []string `json:"totals,omitempty" xmlrpc:"totals,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	Xaxis []string `json:"xaxis,omitempty" xmlrpc:"xaxis,omitempty"`

	// no documentation yet
	Yaxis1 []string `json:"yaxis1,omitempty" xmlrpc:"yaxis1,omitempty"`

	// no documentation yet
	Yaxis10 []string `json:"yaxis10,omitempty" xmlrpc:"yaxis10,omitempty"`

	// no documentation yet
	Yaxis11 []string `json:"yaxis11,omitempty" xmlrpc:"yaxis11,omitempty"`

	// no documentation yet
	Yaxis12 []string `json:"yaxis12,omitempty" xmlrpc:"yaxis12,omitempty"`

	// no documentation yet
	Yaxis13 []string `json:"yaxis13,omitempty" xmlrpc:"yaxis13,omitempty"`

	// no documentation yet
	Yaxis14 []string `json:"yaxis14,omitempty" xmlrpc:"yaxis14,omitempty"`

	// no documentation yet
	Yaxis15 []string `json:"yaxis15,omitempty" xmlrpc:"yaxis15,omitempty"`

	// no documentation yet
	Yaxis16 []string `json:"yaxis16,omitempty" xmlrpc:"yaxis16,omitempty"`

	// no documentation yet
	Yaxis17 []string `json:"yaxis17,omitempty" xmlrpc:"yaxis17,omitempty"`

	// no documentation yet
	Yaxis18 []string `json:"yaxis18,omitempty" xmlrpc:"yaxis18,omitempty"`

	// no documentation yet
	Yaxis19 []string `json:"yaxis19,omitempty" xmlrpc:"yaxis19,omitempty"`

	// no documentation yet
	Yaxis2 []string `json:"yaxis2,omitempty" xmlrpc:"yaxis2,omitempty"`

	// no documentation yet
	Yaxis20 []string `json:"yaxis20,omitempty" xmlrpc:"yaxis20,omitempty"`

	// no documentation yet
	Yaxis3 []string `json:"yaxis3,omitempty" xmlrpc:"yaxis3,omitempty"`

	// no documentation yet
	Yaxis4 []string `json:"yaxis4,omitempty" xmlrpc:"yaxis4,omitempty"`

	// no documentation yet
	Yaxis5 []string `json:"yaxis5,omitempty" xmlrpc:"yaxis5,omitempty"`

	// no documentation yet
	Yaxis6 []string `json:"yaxis6,omitempty" xmlrpc:"yaxis6,omitempty"`

	// no documentation yet
	Yaxis7 []string `json:"yaxis7,omitempty" xmlrpc:"yaxis7,omitempty"`

	// no documentation yet
	Yaxis8 []string `json:"yaxis8,omitempty" xmlrpc:"yaxis8,omitempty"`

	// no documentation yet
	Yaxis9 []string `json:"yaxis9,omitempty" xmlrpc:"yaxis9,omitempty"`
}

// no documentation yet
type Container_Network_CdnMarketplace_Vendor struct {
	Entity

	// no documentation yet
	FeatureSummary *string `json:"featureSummary,omitempty" xmlrpc:"featureSummary,omitempty"`

	// no documentation yet
	Features *string `json:"features,omitempty" xmlrpc:"features,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// no documentation yet
	VendorName *string `json:"vendorName,omitempty" xmlrpc:"vendorName,omitempty"`
}

// SoftLayer_Container_Network_Directory_Listing represents a single entry in a listing of files within a remote directory. API methods that return remote directory listings typically return arrays of SoftLayer_Container_Network_Directory_Listing objects.
type Container_Network_Directory_Listing struct {
	Entity

	// If the file in a directory listing is a directory itself then fileCount is the number of files within the directory.
	FileCount *int `json:"fileCount,omitempty" xmlrpc:"fileCount,omitempty"`

	// The name of a directory or a file within a directory listing.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The type of file in a directory listing. If a directory listing entry is a directory itself then type is set to "directory". Otherwise, type is a blank string.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The LoadBalancer_StatusEntry object stores information about the current status of a particular load balancer service.
//
// It is a data container that cannot be edited, deleted, or saved.
//
// It is returned exclusively by the getStatus method on the [[SoftLayer_Network_LoadBalancer_Service]] service
type Container_Network_LoadBalancer_StatusEntry struct {
	Entity

	// The value of the entry.
	Content *string `json:"content,omitempty" xmlrpc:"content,omitempty"`

	// Text description of the status entry
	Label *string `json:"label,omitempty" xmlrpc:"label,omitempty"`
}

// This datatype is deprecated and will be removed in API version 3.2.
type Container_Network_Message_Delivery_Email struct {
	Entity

	// no documentation yet
	Body *string `json:"body,omitempty" xmlrpc:"body,omitempty"`

	// no documentation yet
	ContainsHtml *bool `json:"containsHtml,omitempty" xmlrpc:"containsHtml,omitempty"`

	// no documentation yet
	From *string `json:"from,omitempty" xmlrpc:"from,omitempty"`

	// no documentation yet
	Subject *string `json:"subject,omitempty" xmlrpc:"subject,omitempty"`

	// no documentation yet
	To *string `json:"to,omitempty" xmlrpc:"to,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Account struct {
	Entity

	// no documentation yet
	Offerings []Container_Network_Message_Delivery_Email_Sendgrid_Account_Offering `json:"offerings,omitempty" xmlrpc:"offerings,omitempty"`

	// no documentation yet
	Profile *Container_Network_Message_Delivery_Email_Sendgrid_Account_Profile `json:"profile,omitempty" xmlrpc:"profile,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Account_Offering struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Quantity *int `json:"quantity,omitempty" xmlrpc:"quantity,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Account_Overview struct {
	Entity

	// no documentation yet
	CreditsAllowed *int `json:"creditsAllowed,omitempty" xmlrpc:"creditsAllowed,omitempty"`

	// no documentation yet
	CreditsOverage *int `json:"creditsOverage,omitempty" xmlrpc:"creditsOverage,omitempty"`

	// no documentation yet
	CreditsRemain *int `json:"creditsRemain,omitempty" xmlrpc:"creditsRemain,omitempty"`

	// no documentation yet
	CreditsUsed *int `json:"creditsUsed,omitempty" xmlrpc:"creditsUsed,omitempty"`

	// no documentation yet
	Email *int `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Package *string `json:"package,omitempty" xmlrpc:"package,omitempty"`

	// no documentation yet
	Reputation *int `json:"reputation,omitempty" xmlrpc:"reputation,omitempty"`

	// no documentation yet
	Requests *int `json:"requests,omitempty" xmlrpc:"requests,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Account_Profile struct {
	Entity

	// no documentation yet
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// no documentation yet
	CompanyWebsite *string `json:"companyWebsite,omitempty" xmlrpc:"companyWebsite,omitempty"`

	// no documentation yet
	CreatedAt *string `json:"createdAt,omitempty" xmlrpc:"createdAt,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	Overage *int `json:"overage,omitempty" xmlrpc:"overage,omitempty"`

	// no documentation yet
	Package *string `json:"package,omitempty" xmlrpc:"package,omitempty"`

	// no documentation yet
	Remain *int `json:"remain,omitempty" xmlrpc:"remain,omitempty"`

	// no documentation yet
	Reputation *int `json:"reputation,omitempty" xmlrpc:"reputation,omitempty"`

	// no documentation yet
	Total *int `json:"total,omitempty" xmlrpc:"total,omitempty"`

	// no documentation yet
	UpdatedAt *string `json:"updatedAt,omitempty" xmlrpc:"updatedAt,omitempty"`

	// no documentation yet
	Used *int `json:"used,omitempty" xmlrpc:"used,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Catalog_Item struct {
	Entity

	// no documentation yet
	Entitlements *Container_Network_Message_Delivery_Email_Sendgrid_Catalog_Item_Entitlements `json:"entitlements,omitempty" xmlrpc:"entitlements,omitempty"`

	// no documentation yet
	Offering *Container_Network_Message_Delivery_Email_Sendgrid_Catalog_Item_Offering `json:"offering,omitempty" xmlrpc:"offering,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Catalog_Item_Entitlements struct {
	Entity

	// no documentation yet
	EmailSendsMaxMonthly *int `json:"emailSendsMaxMonthly,omitempty" xmlrpc:"emailSendsMaxMonthly,omitempty"`

	// no documentation yet
	IpCount *int `json:"ipCount,omitempty" xmlrpc:"ipCount,omitempty"`

	// no documentation yet
	TeammatesMaxTotal *int `json:"teammatesMaxTotal,omitempty" xmlrpc:"teammatesMaxTotal,omitempty"`

	// no documentation yet
	UsersMaxTotal *int `json:"usersMaxTotal,omitempty" xmlrpc:"usersMaxTotal,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Catalog_Item_Offering struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Quantity *int `json:"quantity,omitempty" xmlrpc:"quantity,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Customer_Profile struct {
	Entity

	// no documentation yet
	Address *string `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	Phone *string `json:"phone,omitempty" xmlrpc:"phone,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	Website *string `json:"website,omitempty" xmlrpc:"website,omitempty"`

	// no documentation yet
	Zip *string `json:"zip,omitempty" xmlrpc:"zip,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_List_Entry struct {
	Entity

	// no documentation yet
	Created *string `json:"created,omitempty" xmlrpc:"created,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	Reason *string `json:"reason,omitempty" xmlrpc:"reason,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Statistics struct {
	Entity

	// no documentation yet
	Blocks *int `json:"blocks,omitempty" xmlrpc:"blocks,omitempty"`

	// no documentation yet
	Bounces *int `json:"bounces,omitempty" xmlrpc:"bounces,omitempty"`

	// no documentation yet
	Clicks *int `json:"clicks,omitempty" xmlrpc:"clicks,omitempty"`

	// no documentation yet
	Date *string `json:"date,omitempty" xmlrpc:"date,omitempty"`

	// no documentation yet
	Delivered *int `json:"delivered,omitempty" xmlrpc:"delivered,omitempty"`

	// no documentation yet
	InvalidEmail *int `json:"invalidEmail,omitempty" xmlrpc:"invalidEmail,omitempty"`

	// no documentation yet
	Opens *int `json:"opens,omitempty" xmlrpc:"opens,omitempty"`

	// no documentation yet
	RepeatBounces *int `json:"repeatBounces,omitempty" xmlrpc:"repeatBounces,omitempty"`

	// no documentation yet
	RepeatSpamReports *int `json:"repeatSpamReports,omitempty" xmlrpc:"repeatSpamReports,omitempty"`

	// no documentation yet
	RepeatUnsubscribes *int `json:"repeatUnsubscribes,omitempty" xmlrpc:"repeatUnsubscribes,omitempty"`

	// no documentation yet
	Requests *int `json:"requests,omitempty" xmlrpc:"requests,omitempty"`

	// no documentation yet
	SpamReports *int `json:"spamReports,omitempty" xmlrpc:"spamReports,omitempty"`

	// no documentation yet
	UniqueClicks *int `json:"uniqueClicks,omitempty" xmlrpc:"uniqueClicks,omitempty"`

	// no documentation yet
	UniqueOpens *int `json:"uniqueOpens,omitempty" xmlrpc:"uniqueOpens,omitempty"`

	// no documentation yet
	Unsubscribes *int `json:"unsubscribes,omitempty" xmlrpc:"unsubscribes,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Statistics_Graph struct {
	Entity

	// no documentation yet
	GraphError *string `json:"graphError,omitempty" xmlrpc:"graphError,omitempty"`

	// no documentation yet
	GraphImage *[]byte `json:"graphImage,omitempty" xmlrpc:"graphImage,omitempty"`

	// no documentation yet
	GraphTitle *string `json:"graphTitle,omitempty" xmlrpc:"graphTitle,omitempty"`
}

// no documentation yet
type Container_Network_Message_Delivery_Email_Sendgrid_Statistics_Options struct {
	Entity

	// no documentation yet
	AggregatedBy *bool `json:"aggregatedBy,omitempty" xmlrpc:"aggregatedBy,omitempty"`

	// no documentation yet
	AggregatesOnly *bool `json:"aggregatesOnly,omitempty" xmlrpc:"aggregatesOnly,omitempty"`

	// no documentation yet
	Category *string `json:"category,omitempty" xmlrpc:"category,omitempty"`

	// no documentation yet
	Days *int `json:"days,omitempty" xmlrpc:"days,omitempty"`

	// no documentation yet
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// no documentation yet
	SelectedStatistics []string `json:"selectedStatistics,omitempty" xmlrpc:"selectedStatistics,omitempty"`

	// no documentation yet
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// no documentation yet
type Container_Network_Port_Statistic struct {
	Entity

	// no documentation yet
	AdministrativeStatus *int `json:"administrativeStatus,omitempty" xmlrpc:"administrativeStatus,omitempty"`

	// no documentation yet
	InDiscardPackets *uint `json:"inDiscardPackets,omitempty" xmlrpc:"inDiscardPackets,omitempty"`

	// no documentation yet
	InErrorPackets *uint `json:"inErrorPackets,omitempty" xmlrpc:"inErrorPackets,omitempty"`

	// no documentation yet
	InOctets *uint `json:"inOctets,omitempty" xmlrpc:"inOctets,omitempty"`

	// no documentation yet
	InUnicastPackets *uint `json:"inUnicastPackets,omitempty" xmlrpc:"inUnicastPackets,omitempty"`

	// no documentation yet
	MaximumTransmissionUnit *uint `json:"maximumTransmissionUnit,omitempty" xmlrpc:"maximumTransmissionUnit,omitempty"`

	// no documentation yet
	OperationalStatus *int `json:"operationalStatus,omitempty" xmlrpc:"operationalStatus,omitempty"`

	// no documentation yet
	OutDiscardPackets *uint `json:"outDiscardPackets,omitempty" xmlrpc:"outDiscardPackets,omitempty"`

	// no documentation yet
	OutErrorPackets *uint `json:"outErrorPackets,omitempty" xmlrpc:"outErrorPackets,omitempty"`

	// no documentation yet
	OutOctets *uint `json:"outOctets,omitempty" xmlrpc:"outOctets,omitempty"`

	// no documentation yet
	OutUnicastPackets *uint `json:"outUnicastPackets,omitempty" xmlrpc:"outUnicastPackets,omitempty"`

	// no documentation yet
	PortDuplex *uint `json:"portDuplex,omitempty" xmlrpc:"portDuplex,omitempty"`

	// no documentation yet
	Speed *uint `json:"speed,omitempty" xmlrpc:"speed,omitempty"`
}

// no documentation yet
type Container_Network_SecurityGroup_Limit struct {
	Entity

	// A key value describing what type of limit.
	TypeKey *string `json:"typeKey,omitempty" xmlrpc:"typeKey,omitempty"`

	// The value of the security group limit.
	Value *int `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Container_Network_Service_Resource_ObjectStorage_ConnectionInformation struct {
	Entity

	// no documentation yet
	Datacenter *string `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// no documentation yet
	DatacenterShortName *string `json:"datacenterShortName,omitempty" xmlrpc:"datacenterShortName,omitempty"`

	// no documentation yet
	PrivateEndpoint *string `json:"privateEndpoint,omitempty" xmlrpc:"privateEndpoint,omitempty"`

	// no documentation yet
	PublicEndpoint *string `json:"publicEndpoint,omitempty" xmlrpc:"publicEndpoint,omitempty"`
}

// no documentation yet
type Container_Network_Storage_Backup_Evault_WebCc_Authentication_Details struct {
	Entity

	// no documentation yet
	EventValidation *string `json:"eventValidation,omitempty" xmlrpc:"eventValidation,omitempty"`

	// no documentation yet
	ViewState *string `json:"viewState,omitempty" xmlrpc:"viewState,omitempty"`

	// no documentation yet
	WebCcFormName *string `json:"webCcFormName,omitempty" xmlrpc:"webCcFormName,omitempty"`

	// no documentation yet
	WebCcUrl *string `json:"webCcUrl,omitempty" xmlrpc:"webCcUrl,omitempty"`

	// no documentation yet
	WebCcUserId *string `json:"webCcUserId,omitempty" xmlrpc:"webCcUserId,omitempty"`

	// no documentation yet
	WebCcUserPassword *string `json:"webCcUserPassword,omitempty" xmlrpc:"webCcUserPassword,omitempty"`
}

// no documentation yet
type Container_Network_Storage_DataCenterLimits_VolumeCountLimitContainer struct {
	Entity

	// no documentation yet
	DatacenterName *string `json:"datacenterName,omitempty" xmlrpc:"datacenterName,omitempty"`

	// no documentation yet
	MaximumAvailableCount *int `json:"maximumAvailableCount,omitempty" xmlrpc:"maximumAvailableCount,omitempty"`

	// no documentation yet
	ProvisionedCount *int `json:"provisionedCount,omitempty" xmlrpc:"provisionedCount,omitempty"`
}

// no documentation yet
type Container_Network_Storage_DuplicateConversionStatusInformation struct {
	Entity

	// This represents the timestamp when current conversion process started.
	ActiveConversionStartTime *string `json:"activeConversionStartTime,omitempty" xmlrpc:"activeConversionStartTime,omitempty"`

	// This represents the percentage progress of conversion of a dependent
	DeDuplicateConversionPercentage *int `json:"deDuplicateConversionPercentage,omitempty" xmlrpc:"deDuplicateConversionPercentage,omitempty"`

	// This represents the volume username.
	VolumeUsername *string `json:"volumeUsername,omitempty" xmlrpc:"volumeUsername,omitempty"`
}

// SoftLayer's StorageLayer Evault services provides details regarding the the purchased vault.
//
// When a job is created using the Webcc Console, the job created is identified as a task on the vault. Using this service, information regarding the task can be retrieved.
type Container_Network_Storage_Evault_Vault_Task struct {
	Entity

	// Unique identifier for the task.
	Id *uint `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The hostname provided at time of agent registration.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Total compressed bytes used for the task.
	UsedPoolsize *uint `json:"usedPoolsize,omitempty" xmlrpc:"usedPoolsize,omitempty"`
}

// The SoftLayer_Container_Network_Storage_Evault_WebCc_AgentStatus will contain the timestamp of the last backup performed by the EVault agent.  The agent status will also be returned.
type Container_Network_Storage_Evault_WebCc_AgentStatus struct {
	Entity

	// Timestamp of last backup performed by the EVault backup agent
	LastBackup *Time `json:"lastBackup,omitempty" xmlrpc:"lastBackup,omitempty"`

	// Status indicating the accumulative status result of all jobs performed by the evault agent.  For example, if one job out three jobs failed agent status will by "Failed Backup(s)".
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// The SoftLayer_Container_Network_Storage_Evault_WebCc_BackupResults will contain the timeframe of backups and the results will also be returned.
type Container_Network_Storage_Evault_WebCc_BackupResults struct {
	Entity

	// Timestamp of begin time
	BeginTime *Time `json:"beginTime,omitempty" xmlrpc:"beginTime,omitempty"`

	// Count of backups with conflicts.
	Conflict *string `json:"conflict,omitempty" xmlrpc:"conflict,omitempty"`

	// Timestamp of end time
	EndTime *Time `json:"endTime,omitempty" xmlrpc:"endTime,omitempty"`

	// Count of failed backups.
	Failed *string `json:"failed,omitempty" xmlrpc:"failed,omitempty"`

	// Count of successfull backups.
	Success *string `json:"success,omitempty" xmlrpc:"success,omitempty"`
}

// The SoftLayer_Container_Network_Storage_Evault_WebCc_JobDetails will contain basic details for all backup and restore jobs performed by the StorageLayer EVault service offering.
type Container_Network_Storage_Evault_WebCc_JobDetails struct {
	Entity

	// The number of bytes currently used by the backup job. (provided only for backup jobs)
	BytesUsed *uint `json:"bytesUsed,omitempty" xmlrpc:"bytesUsed,omitempty"`

	// Description of the backup/restore job
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// hardware id
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Date of the last jobrun.
	LastRunDate *Time `json:"lastRunDate,omitempty" xmlrpc:"lastRunDate,omitempty"`

	// Name of the backup/restore job
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Size of backup job when it was first run. (provided only for backup jobs)
	OriginalSize *uint `json:"originalSize,omitempty" xmlrpc:"originalSize,omitempty"`

	// Percentage of overall used space allocated by the job. (provided only for backup jobs)
	PercentageOfTotalUsage *int `json:"percentageOfTotalUsage,omitempty" xmlrpc:"percentageOfTotalUsage,omitempty"`

	// Result of the latest jobrun.
	Result *string `json:"result,omitempty" xmlrpc:"result,omitempty"`

	// virtual guest id
	VirtualGuestId *int `json:"virtualGuestId,omitempty" xmlrpc:"virtualGuestId,omitempty"`
}

// The SoftLayer_Container_Network_Storage_Host will contain the reference id field for the object associated with the host.  The host object type will also be returned.
type Container_Network_Storage_Host struct {
	Entity

	// Reference id field for object associated with host.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Type for the object associated with host
	ObjectType *string `json:"objectType,omitempty" xmlrpc:"objectType,omitempty"`
}

// The SoftLayer_Container_Network_Storage_HostsGatewayInformation will contain the reference id field for the object associated with the host. The host object type will also be returned.
type Container_Network_Storage_HostsGatewayInformation struct {
	Entity

	// Reference id field for object associated with host.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	IsBehindGatewayDevice *bool `json:"isBehindGatewayDevice,omitempty" xmlrpc:"isBehindGatewayDevice,omitempty"`

	// Type for the object associated with host
	ObjectType *string `json:"objectType,omitempty" xmlrpc:"objectType,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_Bucket provides description of a bucket
type Container_Network_Storage_Hub_ObjectStorage_Bucket struct {
	Entity

	// no documentation yet
	BytesUsed *int `json:"bytesUsed,omitempty" xmlrpc:"bytesUsed,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	ObjectCount *int `json:"objectCount,omitempty" xmlrpc:"objectCount,omitempty"`

	// no documentation yet
	StorageLocation *string `json:"storageLocation,omitempty" xmlrpc:"storageLocation,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_ContentDeliveryUrl provides specific details is a container which contains the cdn urls associated with an object storage account
type Container_Network_Storage_Hub_ObjectStorage_ContentDeliveryUrl struct {
	Entity

	// no documentation yet
	Datacenter *string `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// no documentation yet
	FlashUrl *string `json:"flashUrl,omitempty" xmlrpc:"flashUrl,omitempty"`

	// no documentation yet
	HttpUrl *string `json:"httpUrl,omitempty" xmlrpc:"httpUrl,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_Endpoint provides specific details on available endpoint URLs and locations.
type Container_Network_Storage_Hub_ObjectStorage_Endpoint struct {
	Entity

	// no documentation yet
	Legacy *bool `json:"legacy,omitempty" xmlrpc:"legacy,omitempty"`

	// no documentation yet
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// no documentation yet
	Region *string `json:"region,omitempty" xmlrpc:"region,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	Url *string `json:"url,omitempty" xmlrpc:"url,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_File provides specific details that only apply to files that are sent or received from CloudLayer storage resources.
type Container_Network_Storage_Hub_ObjectStorage_File struct {
	Container_Utility_File_Entity

	// no documentation yet
	Folder *string `json:"folder,omitempty" xmlrpc:"folder,omitempty"`

	// no documentation yet
	Hash *string `json:"hash,omitempty" xmlrpc:"hash,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_Container provides details about containers which store collections of files.
type Container_Network_Storage_Hub_ObjectStorage_Folder struct {
	Entity

	// no documentation yet
	Bytes *uint `json:"bytes,omitempty" xmlrpc:"bytes,omitempty"`

	// no documentation yet
	Count *uint `json:"count,omitempty" xmlrpc:"count,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_Node provides detailed information for a particular object storage node
type Container_Network_Storage_Hub_ObjectStorage_Node struct {
	Entity

	// no documentation yet
	DeviceName *string `json:"deviceName,omitempty" xmlrpc:"deviceName,omitempty"`

	// no documentation yet
	ResourceName *string `json:"resourceName,omitempty" xmlrpc:"resourceName,omitempty"`

	// no documentation yet
	UserAuthUrl *string `json:"userAuthUrl,omitempty" xmlrpc:"userAuthUrl,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_Policy provides specific details on available storage policies.
type Container_Network_Storage_Hub_ObjectStorage_Policy struct {
	Entity

	// no documentation yet
	PolicyCode *string `json:"policyCode,omitempty" xmlrpc:"policyCode,omitempty"`
}

// SoftLayer_Container_Network_Storage_Hub_ObjectStorage_Provision provides description of a provision
type Container_Network_Storage_Hub_ObjectStorage_Provision struct {
	Entity

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	Provision *string `json:"provision,omitempty" xmlrpc:"provision,omitempty"`

	// no documentation yet
	ProvisionCreateDate *Time `json:"provisionCreateDate,omitempty" xmlrpc:"provisionCreateDate,omitempty"`

	// no documentation yet
	ProvisionModifyDate *Time `json:"provisionModifyDate,omitempty" xmlrpc:"provisionModifyDate,omitempty"`

	// no documentation yet
	ProvisionTime *int `json:"provisionTime,omitempty" xmlrpc:"provisionTime,omitempty"`
}

// no documentation yet
type Container_Network_Storage_MassDataMigration_Request_Address struct {
	Entity

	// Line 1 of the address - typically the number and street address the MDMS device will be delivered to
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// Line 2 of the address
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// First and last name of the customer on the shipping address
	AddressAttention *string `json:"addressAttention,omitempty" xmlrpc:"addressAttention,omitempty"`

	// The datacenter name where the MDMS device will be shipped to
	AddressNickname *string `json:"addressNickname,omitempty" xmlrpc:"addressNickname,omitempty"`

	// The shipping address city
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Name of the company device is being shipped to
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// The shipping address country
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// The shipping address postal code
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// The shipping address state
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// no documentation yet
type Container_Network_Storage_NetworkConnectionInformation struct {
	Entity

	// no documentation yet
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// no documentation yet
	StorageType *string `json:"storageType,omitempty" xmlrpc:"storageType,omitempty"`
}

// Container for Volume Duplicate Information
type Container_Network_Storage_VolumeDuplicateParameters struct {
	Entity

	// The iopsPerGB of the volume
	IopsPerGb *Float64 `json:"iopsPerGb,omitempty" xmlrpc:"iopsPerGb,omitempty"`

	// Returns true if volume can be duplicated; false otherwise
	IsDuplicatable *bool `json:"isDuplicatable,omitempty" xmlrpc:"isDuplicatable,omitempty"`

	// This represents the location id
	LocationId *int `json:"locationId,omitempty" xmlrpc:"locationId,omitempty"`

	// This represents the location name
	LocationName *string `json:"locationName,omitempty" xmlrpc:"locationName,omitempty"`

	// The maximumIopsPerGb allowed for a duplicated volume
	MaximumIopsPerGb *Float64 `json:"maximumIopsPerGb,omitempty" xmlrpc:"maximumIopsPerGb,omitempty"`

	// The maximumIopsTier allowed for a duplicated volume
	MaximumIopsTier *Float64 `json:"maximumIopsTier,omitempty" xmlrpc:"maximumIopsTier,omitempty"`

	// The maximumVolumeSize allowed for a duplicated volume
	MaximumVolumeSize *int `json:"maximumVolumeSize,omitempty" xmlrpc:"maximumVolumeSize,omitempty"`

	// The minimumIopsPerGb allowed for a duplicated volume
	MinimumIopsPerGb *Float64 `json:"minimumIopsPerGb,omitempty" xmlrpc:"minimumIopsPerGb,omitempty"`

	// The minimumIopsTier allowed for a duplicated volume
	MinimumIopsTier *Float64 `json:"minimumIopsTier,omitempty" xmlrpc:"minimumIopsTier,omitempty"`

	// The minimumVolumeSize allowed for a duplicated volume
	MinimumVolumeSize *int `json:"minimumVolumeSize,omitempty" xmlrpc:"minimumVolumeSize,omitempty"`

	// The volume duplicate status description
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// This represents the volume username
	VolumeUsername *string `json:"volumeUsername,omitempty" xmlrpc:"volumeUsername,omitempty"`
}

// SoftLayer_Container_Subnet_IPAddress models an IP v4 address as it exists as a member of it's subnet, letting the user know if it is a network identifier, gateway, broadcast, or useable address. Addresses that are neither the network identifier nor the gateway nor the broadcast addresses are usable by SoftLayer servers.
type Container_Network_Subnet_IpAddress struct {
	Entity

	// The hardware that an IP address is associated with.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// An IP address expressed in dotted-quad notation.
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// Whether an IP address is its subnet's broadcast address.
	IsBroadcastAddress *bool `json:"isBroadcastAddress,omitempty" xmlrpc:"isBroadcastAddress,omitempty"`

	// Whether an IP address is its subnet's gateway address. Gateway addresses exist on SoftLayer's routers and many not be assigned to servers.
	IsGatewayAddress *bool `json:"isGatewayAddress,omitempty" xmlrpc:"isGatewayAddress,omitempty"`

	// Whether an IP address is its subnet's network identifier address.
	IsNetworkAddress *bool `json:"isNetworkAddress,omitempty" xmlrpc:"isNetworkAddress,omitempty"`
}

// SoftLayer_Container_Network_Subnet_Registration_SubnetReference is provided to reference [[SoftLayer_Network_Subnet_Registration]] object and the [[SoftLayer_Network_Subnet]] it references, in CIDR form.
type Container_Network_Subnet_Registration_SubnetReference struct {
	Entity

	// The ID of the [[SoftLayer_Network_Subnet_Registration]] object.
	RegistrationId *int `json:"registrationId,omitempty" xmlrpc:"registrationId,omitempty"`

	// The subnet address in CIDR form.
	SubnetCidr *string `json:"subnetCidr,omitempty" xmlrpc:"subnetCidr,omitempty"`
}

// SoftLayer_Container_Subnet_Registration_TransactionDetails is provided to return details of a newly created Subnet Registration Transaction.
type Container_Network_Subnet_Registration_TransactionDetails struct {
	Entity

	// The IDs and Subnets of the [[SoftLayer_Network_Subnet_Registration]] object.
	SubnetReferences []Container_Network_Subnet_Registration_SubnetReference `json:"subnetReferences,omitempty" xmlrpc:"subnetReferences,omitempty"`

	// The ID of the Transaction object.
	TransactionId *int `json:"transactionId,omitempty" xmlrpc:"transactionId,omitempty"`
}

// Represents the acceptance status of a Policy.
type Container_Policy_Acceptance struct {
	Entity

	// Flag to indicate if a policy has been previously accepted.
	AcceptedFlag *bool `json:"acceptedFlag,omitempty" xmlrpc:"acceptedFlag,omitempty"`

	// Name of the policy for which we are representing it's acceptance status.
	PolicyName *string `json:"policyName,omitempty" xmlrpc:"policyName,omitempty"`

	// ID of the [[SoftLayer_Product_Item_Policy_Assignment]].
	ProductPolicyAssignmentId *int `json:"productPolicyAssignmentId,omitempty" xmlrpc:"productPolicyAssignmentId,omitempty"`
}

// The SoftLayer_Container_Product_Item_Category data type represents a single product item category.
type Container_Product_Item_Category struct {
	Entity

	// identifier for category.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`
}

// The SoftLayer_Container_Product_Item_Category_Question_Answer data type represents an answer to an item category question.  It contains the category, the question being answered, and the answer.
type Container_Product_Item_Category_Question_Answer struct {
	Entity

	// The answer to the question.
	Answer *string `json:"answer,omitempty" xmlrpc:"answer,omitempty"`

	// The product item category code.
	CategoryCode *string `json:"categoryCode,omitempty" xmlrpc:"categoryCode,omitempty"`

	// The product item category id.
	CategoryId *int `json:"categoryId,omitempty" xmlrpc:"categoryId,omitempty"`

	// The product item category question id.
	QuestionId *int `json:"questionId,omitempty" xmlrpc:"questionId,omitempty"`
}

// The SoftLayer_Container_Product_Item_Category_ZeroFee_Count data type represents a count of zero fee billing/invoice items.
type Container_Product_Item_Category_ZeroFee_Count struct {
	Entity

	// The product item category code.
	CategoryCode *string `json:"categoryCode,omitempty" xmlrpc:"categoryCode,omitempty"`

	// The product item category id.
	CategoryId *int `json:"categoryId,omitempty" xmlrpc:"categoryId,omitempty"`

	// The product item category name.
	CategoryName *string `json:"categoryName,omitempty" xmlrpc:"categoryName,omitempty"`

	// The count of zero fee items for this category.
	Count *int `json:"count,omitempty" xmlrpc:"count,omitempty"`
}

// The SoftLayer_Container_Product_Item_Discount_Program data type represents the information about a discount that is related to a specific product item.
type Container_Product_Item_Discount_Program struct {
	Entity

	// The number of times the item discount(s) may be applied for that order container.  At a minimum the number will be 1 and at most, it will match the quantity of the order container.
	ApplicableQuantity *int `json:"applicableQuantity,omitempty" xmlrpc:"applicableQuantity,omitempty"`

	// The product item that the discount applies to.
	Item *Product_Item `json:"item,omitempty" xmlrpc:"item,omitempty"`

	// The sum of the one time fees (one time, setup and labor) of the prices of this container multiplied by the applicable quantity of this container.
	OneTimeAmount *Float64 `json:"oneTimeAmount,omitempty" xmlrpc:"oneTimeAmount,omitempty"`

	// The tax amount on the one time fees (one time, setup and labor) of the prices of this container mulitiplied by the applicable quantity of this container.
	OneTimeTax *Float64 `json:"oneTimeTax,omitempty" xmlrpc:"oneTimeTax,omitempty"`

	// The item prices that contain the amount of the discount in the recurringFee field.  There may be one or more prices.
	Prices []Product_Item_Price `json:"prices,omitempty" xmlrpc:"prices,omitempty"`

	// The sum of the one time fees (one time, setup and labor) of the prices of this container multiplied by the applicable quantity of this container with the proration factor applied.
	ProratedOneTimeAmount *Float64 `json:"proratedOneTimeAmount,omitempty" xmlrpc:"proratedOneTimeAmount,omitempty"`

	// The tax amount on the one time fees (one time, setup and labor) of the prices of this container mulitiplied by the applicable quantity of this container with the proration factor applied.
	ProratedOneTimeTax *Float64 `json:"proratedOneTimeTax,omitempty" xmlrpc:"proratedOneTimeTax,omitempty"`

	// The sum of the recurring fees of the prices of this container multiplied by the applicable quantity of this container with the proration factor applied.
	ProratedRecurringAmount *Float64 `json:"proratedRecurringAmount,omitempty" xmlrpc:"proratedRecurringAmount,omitempty"`

	// The tax amount on the recurring fees of the prices of this container mulitiplied by the applicable quantity of this container with the proration factor applied.
	ProratedRecurringTax *Float64 `json:"proratedRecurringTax,omitempty" xmlrpc:"proratedRecurringTax,omitempty"`

	// The sum of the recurring fees of the prices of this container multiplied by the applicable quantity of this container.
	RecurringAmount *Float64 `json:"recurringAmount,omitempty" xmlrpc:"recurringAmount,omitempty"`

	// The tax amount on the recurring fees of the prices of this container mulitiplied by the applicable quantity of this container.
	RecurringTax *Float64 `json:"recurringTax,omitempty" xmlrpc:"recurringTax,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order struct {
	Entity

	// Deprecated.
	// Deprecated: This function has been marked as deprecated.
	BigDataOrderFlag *bool `json:"bigDataOrderFlag,omitempty" xmlrpc:"bigDataOrderFlag,omitempty"`

	// Billing Information associated with an order. For existing customers this information is completely ignored. Do not send this information for existing customers.
	BillingInformation *Container_Product_Order_Billing_Information `json:"billingInformation,omitempty" xmlrpc:"billingInformation,omitempty"`

	// This is the ID of the [[SoftLayer_Billing_Order_Item]] of this configuration/container. It is used for rebuilding an order container from a quote and is set automatically.
	BillingOrderItemId *int `json:"billingOrderItemId,omitempty" xmlrpc:"billingOrderItemId,omitempty"`

	// The URL to which PayPal redirects browser after checkout has been canceled before completion of a payment.
	CancelUrl *string `json:"cancelUrl,omitempty" xmlrpc:"cancelUrl,omitempty"`

	// Added by softlayer-go. This hints to the API what kind of product order this is.
	ComplexType *string `json:"complexType,omitempty" xmlrpc:"complexType,omitempty"`

	// User-specified description to identify a particular order container. This is useful if you have a multi-configuration order (multiple <code>orderContainers</code>) and you want to be able to easily determine one from another. Populating this value may be helpful if an exception is thrown when placing an order and it's tied to a specific order container.
	ContainerIdentifier *string `json:"containerIdentifier,omitempty" xmlrpc:"containerIdentifier,omitempty"`

	// This hash is internally-generated and is used to for tracking order containers.
	ContainerSplHash *string `json:"containerSplHash,omitempty" xmlrpc:"containerSplHash,omitempty"`

	// The currency type chosen at checkout.
	CurrencyShortName *string `json:"currencyShortName,omitempty" xmlrpc:"currencyShortName,omitempty"`

	// Device Fingerprint Identifier - Optional.
	DeviceFingerprintId *string `json:"deviceFingerprintId,omitempty" xmlrpc:"deviceFingerprintId,omitempty"`

	// This has been deprecated. It is the identifier used to track configurations in legacy order forms.
	// Deprecated: This function has been marked as deprecated.
	DisplayLayerSessionId *string `json:"displayLayerSessionId,omitempty" xmlrpc:"displayLayerSessionId,omitempty"`

	// no documentation yet
	ExtendedHardwareTesting *bool `json:"extendedHardwareTesting,omitempty" xmlrpc:"extendedHardwareTesting,omitempty"`

	// The [[SoftLayer_Product_Item_Price]] for the Flexible Credit Program discount.  The <code>oneTimeFee</code> field contains the calculated discount being applied to the order.
	FlexibleCreditProgramPrice *Product_Item_Price `json:"flexibleCreditProgramPrice,omitempty" xmlrpc:"flexibleCreditProgramPrice,omitempty"`

	// This flag indicates that the customer consented to the GDPR terms for the quote.
	GdprConsentFlag *bool `json:"gdprConsentFlag,omitempty" xmlrpc:"gdprConsentFlag,omitempty"`

	// For orders that contain servers (bare metal, virtual server, big data, etc.), the hardware property is required. This property is an array of [[SoftLayer_Hardware]] objects. The <code>hostname</code> and <code>domain</code> properties are required for each hardware object. Note that virtual server ([[SoftLayer_Container_Product_Order_Virtual_Guest]]) orders may populate this field instead of the <code>virtualGuests</code> property.
	Hardware []Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// An optional virtual disk image template identifier to be used as an installation base for a computing instance order
	ImageTemplateGlobalIdentifier *string `json:"imageTemplateGlobalIdentifier,omitempty" xmlrpc:"imageTemplateGlobalIdentifier,omitempty"`

	// An optional virtual disk image template identifier to be used as an installation base for a computing instance order
	ImageTemplateId *int `json:"imageTemplateId,omitempty" xmlrpc:"imageTemplateId,omitempty"`

	// Flag to identify a "managed" order. This value is set internally.
	IsManagedOrder *int `json:"isManagedOrder,omitempty" xmlrpc:"isManagedOrder,omitempty"`

	// The collection of [[SoftLayer_Container_Product_Item_Category_Question_Answer]] for any product category that has additional questions requiring user input.
	ItemCategoryQuestionAnswers []Container_Product_Item_Category_Question_Answer `json:"itemCategoryQuestionAnswers,omitempty" xmlrpc:"itemCategoryQuestionAnswers,omitempty"`

	// The [[SoftLayer_Location_Region]] keyname or specific [[SoftLayer_Location_Datacenter]] id where the order should be provisioned. If this value is provided and the <code>regionalGroup</code> property is also specified, an exception will be thrown indicating that only 1 is allowed.
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// This [[SoftLayer_Location]] object will be determined from the <code>location</code> property and will be returned in the order verification or placement response. Any value specified here will get overwritten by the verification process.
	LocationObject *Location `json:"locationObject,omitempty" xmlrpc:"locationObject,omitempty"`

	// A generic message about the order. Does not need to be sent in with any orders.
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// Orders may contain an array of configurations. Populating this property allows you to purchase multiple configurations within a single order. Each order container will have its own individual settings independent of the other order containers. For example, it is possible to order a bare metal server in one configuration and a virtual server in another.
	//
	// If <code>orderContainers</code> is populated on the base order container, most of the configuration-specific properties are ignored on the base container. For example, <code>prices</code>, <code>location</code> and <code>packageId</code> will be ignored on the base container, but since the <code>billingInformation</code> is a property that's not specific to a single order container (but the order as a whole) it must be populated on the base container.
	OrderContainers []Container_Product_Order `json:"orderContainers,omitempty" xmlrpc:"orderContainers,omitempty"`

	// This is deprecated and does not do anything.
	OrderHostnames []string `json:"orderHostnames,omitempty" xmlrpc:"orderHostnames,omitempty"`

	// Collection of exceptions resulting from the verification of the order. This value is set internally and is not required for end users when placing an order. When placing API orders, users can use this value to determine the container-specific exception that was thrown.
	OrderVerificationExceptions []Container_Exception `json:"orderVerificationExceptions,omitempty" xmlrpc:"orderVerificationExceptions,omitempty"`

	// The [[SoftLayer_Product_Package]] id for an order container. This is required to place an order.
	PackageId *int `json:"packageId,omitempty" xmlrpc:"packageId,omitempty"`

	// The Payment Type is Optional. If nothing is sent in, then the normal method of payment will be used. For paypal customers, this means a paypalToken will be returned in the receipt. This token is to be used on the paypal website to complete the order. For Credit Card customers, the card on file in our system will be used to make an initial authorization. To force the order to use a payment type, use one of the following: CARD_ON_FILE or PAYPAL
	PaymentType *string `json:"paymentType,omitempty" xmlrpc:"paymentType,omitempty"`

	// The post-tax recurring charge for the order. This is the sum of preTaxRecurring + totalRecurringTax.
	PostTaxRecurring *Float64 `json:"postTaxRecurring,omitempty" xmlrpc:"postTaxRecurring,omitempty"`

	// The post-tax recurring hourly charge for the order. Since taxes are not calculated for hourly orders, this value will be the same as preTaxRecurringHourly.
	PostTaxRecurringHourly *Float64 `json:"postTaxRecurringHourly,omitempty" xmlrpc:"postTaxRecurringHourly,omitempty"`

	// The post-tax recurring monthly charge for the order. This is the sum of preTaxRecurringMonthly + totalRecurringTax.
	PostTaxRecurringMonthly *Float64 `json:"postTaxRecurringMonthly,omitempty" xmlrpc:"postTaxRecurringMonthly,omitempty"`

	// The post-tax setup fees of the order. This is the sum of preTaxSetup + totalSetupTax;
	PostTaxSetup *Float64 `json:"postTaxSetup,omitempty" xmlrpc:"postTaxSetup,omitempty"`

	// The pre-tax recurring total of the order. If there are mixed monthly and hourly prices on the order, this will be the sum of preTaxRecurringHourly and preTaxRecurringMonthly.
	PreTaxRecurring *Float64 `json:"preTaxRecurring,omitempty" xmlrpc:"preTaxRecurring,omitempty"`

	// The pre-tax hourly recurring total of the order. If there are only monthly prices on the order, this value will be 0.
	PreTaxRecurringHourly *Float64 `json:"preTaxRecurringHourly,omitempty" xmlrpc:"preTaxRecurringHourly,omitempty"`

	// The pre-tax monthly recurring total of the order. If there are only hourly prices on the order, this value will be 0.
	PreTaxRecurringMonthly *Float64 `json:"preTaxRecurringMonthly,omitempty" xmlrpc:"preTaxRecurringMonthly,omitempty"`

	// The pre-tax setup fee total of the order.
	PreTaxSetup *Float64 `json:"preTaxSetup,omitempty" xmlrpc:"preTaxSetup,omitempty"`

	// If there are any presale events available for an order, this value will be populated. It is set internally and is not required for end users when placing an order. See [[SoftLayer_Sales_Presale_Event]] for more info.
	PresaleEvent *Sales_Presale_Event `json:"presaleEvent,omitempty" xmlrpc:"presaleEvent,omitempty"`

	// A preset configuration id for the package. Is required if not submitting any prices.
	PresetId *int `json:"presetId,omitempty" xmlrpc:"presetId,omitempty"`

	// This is a collection of [[SoftLayer_Product_Item_Price]] objects. The only required property to populate for an item price object when ordering is its <code>id</code> - all other supplied information about the price (e.g., recurringFee, setupFee, etc.) will be ignored. Unless the [[SoftLayer_Product_Package]] associated with the order allows for preset prices, this property is required to place an order.
	Prices []Product_Item_Price `json:"prices,omitempty" xmlrpc:"prices,omitempty"`

	// The id of a [[SoftLayer_Hardware_Component_Partition_Template]]. This property is optional. If no partition template is provided, a default will be used according to the operating system chosen with the order. Using the [[SoftLayer_Hardware_Component_Partition_OperatingSystem]] service, getPartitionTemplates will return those available for the particular operating system.
	PrimaryDiskPartitionId *int `json:"primaryDiskPartitionId,omitempty" xmlrpc:"primaryDiskPartitionId,omitempty"`

	// Priorities to set on replication set servers.
	Priorities []string `json:"priorities,omitempty" xmlrpc:"priorities,omitempty"`

	// Deprecated.
	// Deprecated: This function has been marked as deprecated.
	PrivateCloudOrderFlag *bool `json:"privateCloudOrderFlag,omitempty" xmlrpc:"privateCloudOrderFlag,omitempty"`

	// Deprecated.
	// Deprecated: This function has been marked as deprecated.
	PrivateCloudOrderType *string `json:"privateCloudOrderType,omitempty" xmlrpc:"privateCloudOrderType,omitempty"`

	// Optional promotion code for an order.
	PromotionCode *string `json:"promotionCode,omitempty" xmlrpc:"promotionCode,omitempty"`

	// Generic properties.
	Properties []Container_Product_Order_Property `json:"properties,omitempty" xmlrpc:"properties,omitempty"`

	// The Prorated Initial Charge plus the balance on the account. Only the recurring fees are prorated. Here's how the calculation works: We take the postTaxRecurring value and we prorate it based on the time between now and the next bill date for this account. After this, we add in the setup fee since this is not prorated. Then, if there is a balance on the account, we add that to the account. In the event that there is a credit balance on the account, we will subtract this amount from the order total. If the credit balance on the account is greater than the prorated initial charge, the order will go through without a charge to the credit card on the account or requiring a paypal payment. The credit on the account will be reduced by the order total, and the order will await approval from sales, as normal. If there is a pending order already in the system, We will ignore the balance on the account completely, in the calculation of the initial charge. This is to protect against two orders coming into the system and getting the benefit of a credit balance, or worse, both orders being charged the order amount + the balance on the account.
	ProratedInitialCharge *Float64 `json:"proratedInitialCharge,omitempty" xmlrpc:"proratedInitialCharge,omitempty"`

	// This is the same as the proratedInitialCharge, except the balance on the account is ignored. This is the prorated total amount of the order.
	ProratedOrderTotal *Float64 `json:"proratedOrderTotal,omitempty" xmlrpc:"proratedOrderTotal,omitempty"`

	// The URLs for scripts to execute on their respective servers after they have been provisioned. Provision scripts are not available for Microsoft Windows servers.
	ProvisionScripts []string `json:"provisionScripts,omitempty" xmlrpc:"provisionScripts,omitempty"`

	// The quantity of the item being ordered
	Quantity *int `json:"quantity,omitempty" xmlrpc:"quantity,omitempty"`

	// A custom name to be assigned to the quote.
	QuoteName *string `json:"quoteName,omitempty" xmlrpc:"quoteName,omitempty"`

	// Specifying a regional group name allows you to not worry about placing your server or service at a specific datacenter, but to any datacenter within that regional group. See [[SoftLayer_Location_Group_Regional]] to get a list of available regional group names.
	//
	// <code>location</code> and <code>regionalGroup</code> are mutually exclusive on an order container. If both location and regionalGroup are provided, an exception will be thrown indicating that only 1 is allowed.
	//
	// If a regional group is provided and VLANs are specified (within the <code>hardware</code> or <code>virtualGuests</code> properties), we will use the datacenter where the VLANs are located. If no VLANs are specified, we will use the preferred datacenter on the regional group object.
	RegionalGroup *string `json:"regionalGroup,omitempty" xmlrpc:"regionalGroup,omitempty"`

	// Deprecated.
	// Deprecated: This function has been marked as deprecated.
	ResourceGroupId *int `json:"resourceGroupId,omitempty" xmlrpc:"resourceGroupId,omitempty"`

	// Deprecated.
	// Deprecated: This function has been marked as deprecated.
	ResourceGroupName *string `json:"resourceGroupName,omitempty" xmlrpc:"resourceGroupName,omitempty"`

	// An optional resource group template identifier to be used as a deployment base for a Virtual Server (Private Node) order.
	ResourceGroupTemplateId *int `json:"resourceGroupTemplateId,omitempty" xmlrpc:"resourceGroupTemplateId,omitempty"`

	// The URL to which PayPal redirects browser after a payment is completed.
	ReturnUrl *string `json:"returnUrl,omitempty" xmlrpc:"returnUrl,omitempty"`

	// This flag indicates that the quote should be sent to the email address associated with the account or order.
	SendQuoteEmailFlag *bool `json:"sendQuoteEmailFlag,omitempty" xmlrpc:"sendQuoteEmailFlag,omitempty"`

	// The number of cores for the server being ordered. This value is set internally.
	ServerCoreCount *int `json:"serverCoreCount,omitempty" xmlrpc:"serverCoreCount,omitempty"`

	// The token of a requesting service. Do not set.
	ServiceToken *string `json:"serviceToken,omitempty" xmlrpc:"serviceToken,omitempty"`

	// An optional computing instance identifier to be used as an installation base for a computing instance order
	SourceVirtualGuestId *int `json:"sourceVirtualGuestId,omitempty" xmlrpc:"sourceVirtualGuestId,omitempty"`

	// The containers which hold SoftLayer_Security_Ssh_Key IDs to add to their respective servers. The order of containers passed in needs to match the order they are assigned to either hardware or virtualGuests. SSH Keys will not be assigned for servers with Microsoft Windows.
	SshKeys []Container_Product_Order_SshKeys `json:"sshKeys,omitempty" xmlrpc:"sshKeys,omitempty"`

	// An optional parameter for step-based order processing.
	StepId *int `json:"stepId,omitempty" xmlrpc:"stepId,omitempty"`

	//
	//
	// For orders that want to add storage groups such as RAID across multiple disks, simply add [[SoftLayer_Container_Product_Order_Storage_Group]] objects to this array. Storage groups will only be used if the 'RAID' disk controller price is selected. Any other disk controller types will ignore the storage groups set here.
	//
	// The first storage group in this array will be considered the primary storage group, which is used for the OS. Any other storage groups will act as data storage.
	//
	//
	StorageGroups []Container_Product_Order_Storage_Group `json:"storageGroups,omitempty" xmlrpc:"storageGroups,omitempty"`

	// The order container may not contain the final tax rates when it is returned from [[SoftLayer_Product_Order/verifyOrder|verifyOrder]]. This hash will facilitate checking if the tax rates have finished being calculated and retrieving the accurate tax rate values.
	TaxCacheHash *string `json:"taxCacheHash,omitempty" xmlrpc:"taxCacheHash,omitempty"`

	// Flag to indicate if the order container has the final tax rates for the order. Some tax rates are calculated in the background because they take longer, and they might not be finished when the container is returned from [[SoftLayer_Product_Order/verifyOrder|verifyOrder]].
	TaxCompletedFlag *bool `json:"taxCompletedFlag,omitempty" xmlrpc:"taxCompletedFlag,omitempty"`

	// The SoftLayer_Product_Item_Price for the Tech Incubator discount.  The oneTimeFee field contain the calculated discount being applied to the order.
	TechIncubatorItemPrice *Product_Item_Price `json:"techIncubatorItemPrice,omitempty" xmlrpc:"techIncubatorItemPrice,omitempty"`

	// The total tax portion of the recurring fees.
	TotalRecurringTax *Float64 `json:"totalRecurringTax,omitempty" xmlrpc:"totalRecurringTax,omitempty"`

	// The tax amount of the setup fees.
	TotalSetupTax *Float64 `json:"totalSetupTax,omitempty" xmlrpc:"totalSetupTax,omitempty"`

	// This is a collection of [[SoftLayer_Product_Item_Price]] objects which will be used when the service offering being ordered generates usage. This is a read-only property. Setting this property will not change the order.
	UsagePrices []Product_Item_Price `json:"usagePrices,omitempty" xmlrpc:"usagePrices,omitempty"`

	// An optional flag to use hourly pricing instead of standard monthly pricing.
	UseHourlyPricing *bool `json:"useHourlyPricing,omitempty" xmlrpc:"useHourlyPricing,omitempty"`

	// For virtual guest (virtual server) orders, this property is required if you did not specify data in the <code>hardware</code> property. This is an array of [[SoftLayer_Virtual_Guest]] objects. The <code>hostname</code> and <code>domain</code> properties are required for each virtual guest object. There is no need to specify data in this property and the <code>hardware</code> property - only one is required for virtual server orders.
	VirtualGuests []Virtual_Guest `json:"virtualGuests,omitempty" xmlrpc:"virtualGuests,omitempty"`
}

// This datatype is to be used for data transfer requests.
type Container_Product_Order_Account_Media_Data_Transfer_Request struct {
	Container_Product_Order

	// An instance of [[SoftLayer_Account_Media_Data_Transfer_Request]]
	Request *Account_Media_Data_Transfer_Request `json:"request,omitempty" xmlrpc:"request,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. The SoftLayer_Container_Product_Order_Attribute_Address datatype contains the address information.
type Container_Product_Order_Attribute_Address struct {
	Entity

	// The physical street address.
	AddressLine1 *string `json:"addressLine1,omitempty" xmlrpc:"addressLine1,omitempty"`

	// The second line in the address. Information such as suite number goes here.
	AddressLine2 *string `json:"addressLine2,omitempty" xmlrpc:"addressLine2,omitempty"`

	// The city name
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// The 2-character Country code. (i.e. US)
	CountryCode *string `json:"countryCode,omitempty" xmlrpc:"countryCode,omitempty"`

	// State, Region or Province not part of the U.S. or Canada.
	NonUsState *string `json:"nonUsState,omitempty" xmlrpc:"nonUsState,omitempty"`

	// The Zip or Postal Code.
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// U.S. State, Region or Canadian Province.
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. The SoftLayer_Container_Product_Order_Attribute_Contact datatype contains the contact information.
type Container_Product_Order_Attribute_Contact struct {
	Entity

	// The address information of the contact.
	Address *Container_Product_Order_Attribute_Address `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// The email address of the contact.
	EmailAddress *string `json:"emailAddress,omitempty" xmlrpc:"emailAddress,omitempty"`

	// The fax number associated with a contact. This is an optional value.
	FaxNumber *string `json:"faxNumber,omitempty" xmlrpc:"faxNumber,omitempty"`

	// The first name of the contact.
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// The last name of the contact.
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// The organization name of the contact.
	OrganizationName *string `json:"organizationName,omitempty" xmlrpc:"organizationName,omitempty"`

	// The phone number associated with a contact.
	PhoneNumber *string `json:"phoneNumber,omitempty" xmlrpc:"phoneNumber,omitempty"`

	// The title of the contact.
	Title *string `json:"title,omitempty" xmlrpc:"title,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. The SoftLayer_Container_Product_Order_Attribute_Organization datatype contains the organization information.
type Container_Product_Order_Attribute_Organization struct {
	Entity

	// The address information of the contact.
	Address *Container_Product_Order_Attribute_Address `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// The fax number associated with an organization. This is an optional value.
	FaxNumber *string `json:"faxNumber,omitempty" xmlrpc:"faxNumber,omitempty"`

	// The name of an organization.
	OrganizationName *string `json:"organizationName,omitempty" xmlrpc:"organizationName,omitempty"`

	// The phone number associated with an organization.
	PhoneNumber *string `json:"phoneNumber,omitempty" xmlrpc:"phoneNumber,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Billing_Information struct {
	Entity

	// The physical street address. Reserve information such as "apartment #123" or "Suite 2" for line 1.
	BillingAddressLine1 *string `json:"billingAddressLine1,omitempty" xmlrpc:"billingAddressLine1,omitempty"`

	// The second line in the address. Information such as suite number goes here.
	BillingAddressLine2 *string `json:"billingAddressLine2,omitempty" xmlrpc:"billingAddressLine2,omitempty"`

	// The city in which a customer's account resides.
	BillingCity *string `json:"billingCity,omitempty" xmlrpc:"billingCity,omitempty"`

	// The 2-character Country code for an account's address. (i.e. US)
	BillingCountryCode *string `json:"billingCountryCode,omitempty" xmlrpc:"billingCountryCode,omitempty"`

	// The email address associated with a customer account.
	BillingEmail *string `json:"billingEmail,omitempty" xmlrpc:"billingEmail,omitempty"`

	// the company name for an account.
	BillingNameCompany *string `json:"billingNameCompany,omitempty" xmlrpc:"billingNameCompany,omitempty"`

	// The first name of the customer account owner.
	BillingNameFirst *string `json:"billingNameFirst,omitempty" xmlrpc:"billingNameFirst,omitempty"`

	// The last name of the customer account owner
	BillingNameLast *string `json:"billingNameLast,omitempty" xmlrpc:"billingNameLast,omitempty"`

	// The fax number associated with a customer account.
	BillingPhoneFax *string `json:"billingPhoneFax,omitempty" xmlrpc:"billingPhoneFax,omitempty"`

	// The phone number associated with a customer account.
	BillingPhoneVoice *string `json:"billingPhoneVoice,omitempty" xmlrpc:"billingPhoneVoice,omitempty"`

	// The Zip or Postal Code for the billing address on an account.
	BillingPostalCode *string `json:"billingPostalCode,omitempty" xmlrpc:"billingPostalCode,omitempty"`

	// The State for the account.
	BillingState *string `json:"billingState,omitempty" xmlrpc:"billingState,omitempty"`

	// Total height of browser screen in pixels.
	BrowserScreenHeight *string `json:"browserScreenHeight,omitempty" xmlrpc:"browserScreenHeight,omitempty"`

	// Total width of browser screen in pixels.
	BrowserScreenWidth *string `json:"browserScreenWidth,omitempty" xmlrpc:"browserScreenWidth,omitempty"`

	// The credit card number to use.
	CardAccountNumber *string `json:"cardAccountNumber,omitempty" xmlrpc:"cardAccountNumber,omitempty"`

	// The payment card expiration month
	CardExpirationMonth *int `json:"cardExpirationMonth,omitempty" xmlrpc:"cardExpirationMonth,omitempty"`

	// The payment card expiration year
	CardExpirationYear *int `json:"cardExpirationYear,omitempty" xmlrpc:"cardExpirationYear,omitempty"`

	// The Card Verification Value Code (CVV) number
	CreditCardVerificationNumber *string `json:"creditCardVerificationNumber,omitempty" xmlrpc:"creditCardVerificationNumber,omitempty"`

	// 1 = opted in,  0 = not opted in. Select the EU Supported option if you use IBM Bluemix Infrastructure services to process EU citizens' personal data. This option limits Level 1 and Level 2 support to the EU. However, IBM Bluemix and SoftLayer teams outside the EU perform processing activities when they are not resolved at Level 1 or 2. These activities are always at your instruction and do not impact the security or privacy of your data. As with our standard services, you must review the impact these cross-border processing activities have on your services and take any necessary measures, including review of IBM's US-EU Privacy Shield registration and Data Processing Addendum.  If you select products, services, or locations outside the EU, all processing activities will be performed outside of the EU. If you select other IBM services in addition to Bluemix IaaS (IBM or a third party), determine the service location in order to meet any additional data protection or processing requirements that permit cross-border transfers.
	EuSupported *bool `json:"euSupported,omitempty" xmlrpc:"euSupported,omitempty"`

	// If true, order is being placed by a business.
	IsBusinessFlag *bool `json:"isBusinessFlag,omitempty" xmlrpc:"isBusinessFlag,omitempty"`

	// The purpose of this property is to allow enablement of 3D Secure (3DS). This is the Reference ID that corresponds to the device data for Payer Authentication. In order to properly enable 3DS, this will require implementation of Cardinal Cruise Hybrid.
	//
	// Please refer to https://cardinaldocs.atlassian.net/wiki/spaces/CC/pages/360668/Cardinal+Cruise+Hybrid and view section under "DFReferenceId / ReferenceId" to populate this property accordingly.
	PayerAuthenticationEnrollmentReferenceId *string `json:"payerAuthenticationEnrollmentReferenceId,omitempty" xmlrpc:"payerAuthenticationEnrollmentReferenceId,omitempty"`

	// The URL where the issuing bank will redirect.
	PayerAuthenticationEnrollmentReturnUrl *string `json:"payerAuthenticationEnrollmentReturnUrl,omitempty" xmlrpc:"payerAuthenticationEnrollmentReturnUrl,omitempty"`

	// "Continue with Consumer Authentication" decoded response JWT (JSON Web Token) after successful authentication. The response is part of the implementation of Cardinal Cruise Hybrid.
	//
	// Please refer to https://cardinaldocs.atlassian.net/wiki/spaces/CC/pages/360668/Cardinal+Cruise+Hybrid and view section under "Continue with Consumer Authentication" to populate this property accordingly based on the CCA response.
	PayerAuthenticationWebToken *string `json:"payerAuthenticationWebToken,omitempty" xmlrpc:"payerAuthenticationWebToken,omitempty"`

	// Tax exempt status. 1 = exempt (not taxable),  0 = not exempt (taxable)
	TaxExempt *int `json:"taxExempt,omitempty" xmlrpc:"taxExempt,omitempty"`

	// The VAT ID entered at checkout
	VatId *string `json:"vatId,omitempty" xmlrpc:"vatId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a Gateway Appliance Cluster order with SoftLayer.
type Container_Product_Order_Gateway_Appliance_Cluster struct {
	Container_Product_Order

	// Used to identify which items on an order belong in the same cluster.
	ClusterIdentifier *string `json:"clusterIdentifier,omitempty" xmlrpc:"clusterIdentifier,omitempty"`

	// Indicates what type of cluster order is being placed (HA, Provision).
	ClusterOrderType *string `json:"clusterOrderType,omitempty" xmlrpc:"clusterOrderType,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to upgrade a [[SoftLayer_Network_Gateway (type)|network gateway]].
type Container_Product_Order_Gateway_Appliance_Upgrade struct {
	Container_Product_Order

	// Identifier for the [[SoftLayer_Network_Gateway (type)|network gateway]] being upgraded.
	GatewayId *int `json:"gatewayId,omitempty" xmlrpc:"gatewayId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a hardware security module order with SoftLayer.
type Container_Product_Order_Hardware_Security_Module struct {
	Container_Product_Order_Hardware_Server
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Hardware_Server struct {
	Container_Product_Order

	// Used to identify which category should be used for the boot disk.
	BootCategoryCode *string `json:"bootCategoryCode,omitempty" xmlrpc:"bootCategoryCode,omitempty"`

	// Used to identify which items on an order belong in the same cluster.
	ClusterIdentifier *string `json:"clusterIdentifier,omitempty" xmlrpc:"clusterIdentifier,omitempty"`

	// Indicates what type of cluster order is being placed (HA, Provision).
	ClusterOrderType *string `json:"clusterOrderType,omitempty" xmlrpc:"clusterOrderType,omitempty"`

	// Used to identify which gateway is being upgraded to HA.
	ClusterResourceId *int `json:"clusterResourceId,omitempty" xmlrpc:"clusterResourceId,omitempty"`

	// Array of disk drive slot categories to destroy on reclaim. For example: ['disk0', 'disk1', 'disk2']. One drive_destruction price must be included for each slot provided. Note that once the initial order or upgrade order are approved, the destruction property <strong>is not removable</strong> and the drives will be destroyed at the end of the server's lifecycle. Not all drive slots are required, but all can be provided.
	DriveDestructionDisks []string `json:"driveDestructionDisks,omitempty" xmlrpc:"driveDestructionDisks,omitempty"`

	// Id used with the monitoring package. (Deprecated)
	// Deprecated: This function has been marked as deprecated.
	MonitoringAgentConfigurationTemplateGroupId *int `json:"monitoringAgentConfigurationTemplateGroupId,omitempty" xmlrpc:"monitoringAgentConfigurationTemplateGroupId,omitempty"`

	// When ordering Virtual Server (Private Node), this variable specifies the role of the server configuration. (Deprecated)
	PrivateCloudServerRole *string `json:"privateCloudServerRole,omitempty" xmlrpc:"privateCloudServerRole,omitempty"`

	// Used to identify which device the new server should be attached to.
	RequiredUpstreamDeviceId *int `json:"requiredUpstreamDeviceId,omitempty" xmlrpc:"requiredUpstreamDeviceId,omitempty"`

	// tags (used in MongoDB deployments). (Deprecated)
	Tags []Container_Product_Order_Property `json:"tags,omitempty" xmlrpc:"tags,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Hardware_Server_Colocation struct {
	Container_Product_Order_Hardware_Server
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a Gateway Appliance order.
type Container_Product_Order_Hardware_Server_Gateway_Appliance struct {
	Container_Product_Order_Hardware_Server
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a hardware upgrade.
type Container_Product_Order_Hardware_Server_Upgrade struct {
	Container_Product_Order_Hardware_Server
}

// Use this datatype to upgrade your existing monthly-billed server to term based pricing. Only monthly to 1 year, and 1 year to 3 year migrations are available. A new billing agreement contract will be created upon order approval, starting at the next billing cycle. A price is required for each existing billing item and all term-based prices must match in length. Hourly billed servers are not eligible for this upgrade. Downgrading to a shorter term is not available. Multiple term upgrades per billing cycle are not allowed.
type Container_Product_Order_Hardware_Server_Upgrade_MigrateToReserved struct {
	Container_Product_Order_Hardware_Server_Upgrade

	// no documentation yet
	TermLength *int `json:"termLength,omitempty" xmlrpc:"termLength,omitempty"`

	// no documentation yet
	TermStartDate *Time `json:"termStartDate,omitempty" xmlrpc:"termStartDate,omitempty"`
}

// no documentation yet
type Container_Product_Order_Hardware_Server_Vpc struct {
	Container_Product_Order_Hardware_Server

	// no documentation yet
	Crn *string `json:"crn,omitempty" xmlrpc:"crn,omitempty"`

	// no documentation yet
	InstanceProfile *string `json:"instanceProfile,omitempty" xmlrpc:"instanceProfile,omitempty"`

	// no documentation yet
	IpAllocations []Container_Product_Order_Vpc_IpAllocation `json:"ipAllocations,omitempty" xmlrpc:"ipAllocations,omitempty"`

	// no documentation yet
	ResourceGroup *string `json:"resourceGroup,omitempty" xmlrpc:"resourceGroup,omitempty"`

	// no documentation yet
	ServerId *string `json:"serverId,omitempty" xmlrpc:"serverId,omitempty"`

	// no documentation yet
	ServicePortInterfaceId *string `json:"servicePortInterfaceId,omitempty" xmlrpc:"servicePortInterfaceId,omitempty"`

	// no documentation yet
	ServicePortIpAllocationId *string `json:"servicePortIpAllocationId,omitempty" xmlrpc:"servicePortIpAllocationId,omitempty"`

	// no documentation yet
	ServicePortVpcId *string `json:"servicePortVpcId,omitempty" xmlrpc:"servicePortVpcId,omitempty"`

	// no documentation yet
	Subnets []Container_Product_Order_Vpc_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// no documentation yet
	Zone *string `json:"zone,omitempty" xmlrpc:"zone,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a Monitoring Package order with SoftLayer. This class is no longer available.
type Container_Product_Order_Monitoring_Package struct {
	Container_Product_Order

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	ServerType *string `json:"serverType,omitempty" xmlrpc:"serverType,omitempty"`
}

// This is a datatype used with multi-configuration deployments. Multi-configuration deployments also have a deployment specific datatype that should be used in lieu of this one.
type Container_Product_Order_MultiConfiguration struct {
	Container_Product_Order
}

// no documentation yet
type Container_Product_Order_MultiConfiguration_Tornado struct {
	Container_Product_Order_MultiConfiguration
}

// (DEPRECATED) This type contains the structure of network-related objects that may be specified when ordering services.
type Container_Product_Order_Network struct {
	Entity

	// The [[SoftLayer_Network]] object.
	Network *Network `json:"network,omitempty" xmlrpc:"network,omitempty"`

	// The list of public [[SoftLayer_Container_Product_Order_Network_Vlan|vlans]] available for ordering. Each VLAN will have list of public subnets that are accessible to the VLAN.
	PublicVlans []Container_Product_Order `json:"publicVlans,omitempty" xmlrpc:"publicVlans,omitempty"`

	// The list of private [[SoftLayer_Container_Product_Order_Network_Subnet|subnets]] available for ordering with a description of their available IP space.
	Subnets []Container_Product_Order `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an application delivery controller order with SoftLayer.
type Container_Product_Order_Network_Application_Delivery_Controller struct {
	Container_Product_Order

	// An optional [[SoftLayer_Network_Application_Delivery_Controller]] identifier that is used for upgrading an existing application delivery controller.
	ApplicationDeliveryControllerId *int `json:"applicationDeliveryControllerId,omitempty" xmlrpc:"applicationDeliveryControllerId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder when purchasing a Network Interconnect.
type Container_Product_Order_Network_Interconnect struct {
	Container_Product_Order

	// The BGP ASN.
	BgpAsn *string `json:"bgpAsn,omitempty" xmlrpc:"bgpAsn,omitempty"`

	// The [[SoftLayer_Network_Interconnect]] for this order, ID must be provided.
	InterconnectId *int `json:"interconnectId,omitempty" xmlrpc:"interconnectId,omitempty"`

	// The [[SoftLayer_Network_DirectLink_Location]] for this order, ID must be provided.
	InterconnectLocation *Network_DirectLink_Location `json:"interconnectLocation,omitempty" xmlrpc:"interconnectLocation,omitempty"`

	// The [[SoftLayer_Network_Interconnect_Tenant]] being ordered. Only the ID is required. If this ID is specified, then properties such as networkIdentifier, ipAddressRange, and interconnectId do not need to be specified.
	InterconnectTenant *Network_Interconnect_Tenant `json:"interconnectTenant,omitempty" xmlrpc:"interconnectTenant,omitempty"`

	// Optional IP address for this link.
	IpAddressRange *string `json:"ipAddressRange,omitempty" xmlrpc:"ipAddressRange,omitempty"`

	// A name to identify this Direct Link resource.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Optional network identifier for this link.
	NetworkIdentifier *string `json:"networkIdentifier,omitempty" xmlrpc:"networkIdentifier,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an upgrade order for Direct Link.
type Container_Product_Order_Network_Interconnect_Upgrade struct {
	Container_Product_Order_Network_Interconnect
}

// This is the default container type for network load balancer orders.
type Container_Product_Order_Network_LoadBalancer struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for a Load Balancer as a Service.
type Container_Product_Order_Network_LoadBalancer_AsAService struct {
	Container_Product_Order

	// A description of this Load Balancer.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The [[SoftLayer_Network_LBaaS_LoadBalancerHealthMonitorConfiguration]]s for this Load Balancer.
	HealthMonitorConfigurations []Network_LBaaS_LoadBalancerHealthMonitorConfiguration `json:"healthMonitorConfigurations,omitempty" xmlrpc:"healthMonitorConfigurations,omitempty"`

	// Specify whether this load balancer is a public or internal facing load balancer. If this value is omitted, the value will default to true.
	IsPublic *bool `json:"isPublic,omitempty" xmlrpc:"isPublic,omitempty"`

	// A name to identify this Load Balancer.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The [[SoftLayer_Network_LBaaS_LoadBalancerProtocolConfiguration]]s for this Load Balancer.
	ProtocolConfigurations []Network_LBaaS_LoadBalancerProtocolConfiguration `json:"protocolConfigurations,omitempty" xmlrpc:"protocolConfigurations,omitempty"`

	// Specify the public subnet where this load balancer will be provisioned when useSystemPublicIpPool is false. This is valid only for public(1) load balancer. The public subnet should match the private subnet.
	PublicSubnets []Network_Subnet `json:"publicSubnets,omitempty" xmlrpc:"publicSubnets,omitempty"`

	// The [[SoftLayer_Network_LBaaS_LoadBalancerServerInstanceInfo]]s for this Load Balancer.
	ServerInstancesInformation []Network_LBaaS_LoadBalancerServerInstanceInfo `json:"serverInstancesInformation,omitempty" xmlrpc:"serverInstancesInformation,omitempty"`

	// The [[SoftLayer_Network_Subnet]]s where this Load Balancer will be provisioned.
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// Specify the type of this load balancer. If isPublic is omitted, it specifies the load balacner as private(0), public(1) or public to public(2). If isPublic is set as True, only public(1) or public to public(2) is valid. If isPublic is set as False, this value is ignored. If this value is omitted, the value will be set according to isPublic value.
	Type *int `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// Specify if this load balancer uses system IP pool (true, default) or customer's (null|false) public subnet to allocate IP addresses.
	UseSystemPublicIpPool *bool `json:"useSystemPublicIpPool,omitempty" xmlrpc:"useSystemPublicIpPool,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a network message delivery order with SoftLayer.
type Container_Product_Order_Network_Message_Delivery struct {
	Container_Product_Order

	// This property has been deprecated and should no longer be used.
	//
	// The account password for SendGrid enrollment.
	// Deprecated: This function has been marked as deprecated.
	AccountPassword *string `json:"accountPassword,omitempty" xmlrpc:"accountPassword,omitempty"`

	// This property has been deprecated and should no longer be used.
	//
	// The username for SendGrid enrollment.
	// Deprecated: This function has been marked as deprecated.
	AccountUsername *string `json:"accountUsername,omitempty" xmlrpc:"accountUsername,omitempty"`

	// The email address for SendGrid enrollment.
	EmailAddress *string `json:"emailAddress,omitempty" xmlrpc:"emailAddress,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an upgrade order for network message delivery.
type Container_Product_Order_Network_Message_Delivery_Upgrade struct {
	Container_Product_Order_Network_Message_Delivery

	// The ID of the [[SoftLayer_Network_Message_Delivery]] being upgraded.
	MessageDeliveryId *int `json:"messageDeliveryId,omitempty" xmlrpc:"messageDeliveryId,omitempty"`
}

// This is the base data type for Performance storage order containers. If you wish to place an order you must not use this class and instead use the appropriate child container for the type of storage you would like to order: [[SoftLayer_Container_Product_Order_Network_PerformanceStorage_Nfs]] for File and [[SoftLayer_Container_Product_Order_Network_PerformanceStorage_Iscsi]] for Block storage.
type Container_Product_Order_Network_PerformanceStorage struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for iSCSI (Block) Performance Storage
type Container_Product_Order_Network_PerformanceStorage_Iscsi struct {
	Container_Product_Order_Network_PerformanceStorage

	// OS Type to be used when formatting the storage space, this should match the OS type that will be connecting to the LUN. The only required property its the keyName of the OS type.
	OsFormatType *Network_Storage_Iscsi_OS_Type `json:"osFormatType,omitempty" xmlrpc:"osFormatType,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for NFS (File) Performance Storage
type Container_Product_Order_Network_PerformanceStorage_Nfs struct {
	Container_Product_Order_Network_PerformanceStorage
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a hardware firewall order with SoftLayer.
type Container_Product_Order_Network_Protection_Firewall struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a hardware (dedicated) firewall order with SoftLayer.
type Container_Product_Order_Network_Protection_Firewall_Dedicated struct {
	Container_Product_Order

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	RouterId *int `json:"routerId,omitempty" xmlrpc:"routerId,omitempty"`

	// generic properties.
	Vlan *Network_Vlan `json:"vlan,omitempty" xmlrpc:"vlan,omitempty"`

	// generic properties.
	VlanId *int `json:"vlanId,omitempty" xmlrpc:"vlanId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Network_Protection_Firewall_Dedicated_Upgrade struct {
	Container_Product_Order_Network_Protection_Firewall_Dedicated

	// no documentation yet
	FirewallId *int `json:"firewallId,omitempty" xmlrpc:"firewallId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for Storage as a Service.
type Container_Product_Order_Network_Storage_AsAService struct {
	Container_Product_Order

	// Optional property to specify provisioning to a dedicated cluster at order time. The `id` property of the [[SoftLayer_Network_Storage_DedicatedCluster]] should be provided to dictate where to provision storage to. Note your account must be enabled to order into the desired location(s) prior to placing the order.
	DedicatedCluster *Network_Storage_DedicatedCluster `json:"dedicatedCluster,omitempty" xmlrpc:"dedicatedCluster,omitempty"`

	// This must be populated only for duplicating a specific snapshot for volume duplicating. It represents the identifier of the origin [[SoftLayer_Network_Storage_Snapshot]]
	DuplicateOriginSnapshotId *int `json:"duplicateOriginSnapshotId,omitempty" xmlrpc:"duplicateOriginSnapshotId,omitempty"`

	// This must be populated only for duplicate volume ordering. It represents the identifier of the origin [[SoftLayer_Network_Storage]].
	DuplicateOriginVolumeId *int `json:"duplicateOriginVolumeId,omitempty" xmlrpc:"duplicateOriginVolumeId,omitempty"`

	// When ordering performance by IOPS, populate this property with how many.
	Iops *int `json:"iops,omitempty" xmlrpc:"iops,omitempty"`

	// This can be optionally populated only for duplicate volume ordering. When set, this flag denotes that the duplicate volume being ordered can refresh its data using snapshots from the specified origin volume.
	IsDependentDuplicateFlag *bool `json:"isDependentDuplicateFlag,omitempty" xmlrpc:"isDependentDuplicateFlag,omitempty"`

	// This must be populated only for replicant volume ordering. It represents the identifier of the origin [[SoftLayer_Network_Storage]].
	OriginVolumeId *int `json:"originVolumeId,omitempty" xmlrpc:"originVolumeId,omitempty"`

	// This must be populated only for replicant volume ordering. It represents the [[SoftLayer_Network_Storage_Schedule]] that will be be used to replicate the origin [[SoftLayer_Network_Storage]] volume.
	OriginVolumeScheduleId *int `json:"originVolumeScheduleId,omitempty" xmlrpc:"originVolumeScheduleId,omitempty"`

	// This must be populated for block storage orders. This should match the OS type of the host(s) that will connect to the volume. The only required property is the keyName of the OS type. This property is ignored for file storage orders.
	OsFormatType *Network_Storage_Iscsi_OS_Type `json:"osFormatType,omitempty" xmlrpc:"osFormatType,omitempty"`

	// Volume size in GB's.
	VolumeSize *int `json:"volumeSize,omitempty" xmlrpc:"volumeSize,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an upgrade order for Storage as a Service.
type Container_Product_Order_Network_Storage_AsAService_Upgrade struct {
	Container_Product_Order_Network_Storage_AsAService

	// The [[SoftLayer_Network_Storage]] being upgraded. Only it's ID is required.
	Volume *Network_Storage `json:"volume,omitempty" xmlrpc:"volume,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for additional Evault plugins.
type Container_Product_Order_Network_Storage_Backup_Evault_Plugin struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an Evault order with SoftLayer.
type Container_Product_Order_Network_Storage_Backup_Evault_Vault struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for Enterprise Storage
type Container_Product_Order_Network_Storage_Enterprise struct {
	Container_Product_Order

	// This must be populated only for replicant volume ordering. It represents the identifier of the origin [[SoftLayer_Network_Storage]].
	OriginVolumeId *int `json:"originVolumeId,omitempty" xmlrpc:"originVolumeId,omitempty"`

	// This must be populated only for replicant volume ordering. It represents the [[SoftLayer_Network_Storage_Schedule]] that will be be used to replicate the origin [[SoftLayer_Network_Storage]] volume.
	OriginVolumeScheduleId *int `json:"originVolumeScheduleId,omitempty" xmlrpc:"originVolumeScheduleId,omitempty"`

	// This must be populated for block storage orders. This should match the OS type of the host(s) that will connect to the volume. The only required property is the keyName of the OS type. This property is ignored for file storage orders.
	OsFormatType *Network_Storage_Iscsi_OS_Type `json:"osFormatType,omitempty" xmlrpc:"osFormatType,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for Enterprise Storage Snapshot Space.
type Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace struct {
	Container_Product_Order

	// The [[SoftLayer_Network_Storage]] id for which snapshot space is being ordered for.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an upgrade order for Enterprise Storage Snapshot Space.
type Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade struct {
	Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace
}

// This datatype is to be used for object storage orders.
type Container_Product_Order_Network_Storage_Hub struct {
	Container_Product_Order
}

// This class is used to contain a datacenter location and its associated active usage rate prices for object storage ordering.
type Container_Product_Order_Network_Storage_Hub_Datacenter struct {
	Entity

	// The datacenter location where object storage is available.
	Location *Location `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// The collection of active usage rate item prices.
	UsageRatePrices []Product_Item_Price `json:"usageRatePrices,omitempty" xmlrpc:"usageRatePrices,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an ISCSI order with SoftLayer.
type Container_Product_Order_Network_Storage_Iscsi struct {
	Container_Product_Order
}

// This datatype is to be used for mass data migration requests.
type Container_Product_Order_Network_Storage_MassDataMigration_Request struct {
	Container_Product_Order

	// Line 1 of the address - typically the number and street address the MDMS device will be delivered to
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// Line 2 of the address
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// First and last name of the customer on the shipping address
	AddressAttention *string `json:"addressAttention,omitempty" xmlrpc:"addressAttention,omitempty"`

	// The datacenter name where the MDMS device will be shipped to
	AddressNickname *string `json:"addressNickname,omitempty" xmlrpc:"addressNickname,omitempty"`

	// The shipping address city
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// Name of the company device is being shipped to
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// Cloud Object Storage Account ID for the data offload destination
	CosAccountId *string `json:"cosAccountId,omitempty" xmlrpc:"cosAccountId,omitempty"`

	// Cloud Object Storage Bucket for the data offload destination
	CosBucketName *string `json:"cosBucketName,omitempty" xmlrpc:"cosBucketName,omitempty"`

	// The shipping address country
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// Default Gateway used for preconfiguring the Eth1 port on the MDMS device to access the user interface
	Eth1DefaultGateway *string `json:"eth1DefaultGateway,omitempty" xmlrpc:"eth1DefaultGateway,omitempty"`

	// Netmask used for preconfiguring the Eth1 port on the MDMS device to access the user interface
	Eth1Netmask *string `json:"eth1Netmask,omitempty" xmlrpc:"eth1Netmask,omitempty"`

	// Static IP Address used for preconfiguring the Eth1 port on the MDMS device to access the user interface
	Eth1StaticIp *string `json:"eth1StaticIp,omitempty" xmlrpc:"eth1StaticIp,omitempty"`

	// Netmask used for preconfiguring the Eth3 port on the MDMS device to enable data transfer
	Eth3Netmask *string `json:"eth3Netmask,omitempty" xmlrpc:"eth3Netmask,omitempty"`

	// Static IP Address used for preconfiguring the Eth3 port on the MDMS device to enable data transfer
	Eth3StaticIp *string `json:"eth3StaticIp,omitempty" xmlrpc:"eth3StaticIp,omitempty"`

	// The e-mails of the MDMS key contacts
	KeyContactEmails []string `json:"keyContactEmails,omitempty" xmlrpc:"keyContactEmails,omitempty"`

	// The names of the MDMS key contacts
	KeyContactNames []string `json:"keyContactNames,omitempty" xmlrpc:"keyContactNames,omitempty"`

	// The phone numbers of the MDMS key contacts
	KeyContactPhoneNumbers []string `json:"keyContactPhoneNumbers,omitempty" xmlrpc:"keyContactPhoneNumbers,omitempty"`

	// The roles of the MDMS key contacts
	KeyContactRoles []string `json:"keyContactRoles,omitempty" xmlrpc:"keyContactRoles,omitempty"`

	// The shipping address postal code
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Name of the Mass Data Migration Service job request
	RequestName *string `json:"requestName,omitempty" xmlrpc:"requestName,omitempty"`

	// Shipping address and information where device will be shipped to
	ShippingAddress *Container_Network_Storage_MassDataMigration_Request_Address `json:"shippingAddress,omitempty" xmlrpc:"shippingAddress,omitempty"`

	// The shipping address state
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// The SoftLayer_Container_Product_Order_Network_Storage_Modification datatype has everything required to place a modification to an existing StorageLayer account with SoftLayer. Modifications, at present time, include upgrade and downgrades only. The ”volumeId” property must be set to the network storage volume id to be upgraded. Once populated send this container to the [[SoftLayer_Product_Order::placeOrder]] method.
//
// The ”packageId” property passed in for CloudLayer storage accounts must be set to 0 (zero) and the ”quantity” property must be set to 1. The location does not have to be set. Please use the [[SoftLayer_Product_Package]] service to retrieve a list of CloudLayer items.
//
// NOTE: When upgrading CloudLayer storage service from a metered plan (pay as you go) to a non-metered plan, make sure the chosen plan's storage allotment has enough space to cover the current usage. If the chosen plan's usage allotment is less than the CloudLayer storage's usage the order will be rejected.
type Container_Product_Order_Network_Storage_Modification struct {
	Container_Product_Order

	// The id of the StorageLayer account to modify.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder when placing network attached storage orders.
type Container_Product_Order_Network_Storage_Nas struct {
	Container_Product_Order
}

// This datatype is to be used for ordering object storage products using the object_storage [[SoftLayer_Product_Item_Category|category]]. For object storage products using hub [[SoftLayer_Product_Item_Category|category]] use the [[SoftLayer_Container_Product_Order_Network_Storage_Hub]] order container.
type Container_Product_Order_Network_Storage_Object struct {
	Container_Product_Order
}

// This class is used to contain a location group and its associated active usage rate prices for object storage ordering.
type Container_Product_Order_Network_Storage_ObjectStorage_LocationGroup struct {
	Entity

	// The datacenter location where object storage is available.
	ClusterGeolocationType *string `json:"clusterGeolocationType,omitempty" xmlrpc:"clusterGeolocationType,omitempty"`

	// The datacenter location where object storage is available.
	LocationGroup *Location_Group `json:"locationGroup,omitempty" xmlrpc:"locationGroup,omitempty"`

	// The collection of active usage rate item prices.
	UsageRatePrices []Product_Item_Price `json:"usageRatePrices,omitempty" xmlrpc:"usageRatePrices,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a subnet order with SoftLayer.
type Container_Product_Order_Network_Subnet struct {
	Container_Product_Order

	// The description which includes the network identifier, Classless Inter-Domain Routing prefix and the available slot count.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The [[SoftLayer_Network_Subnet_IpAddress]] id.
	EndPointIpAddressId *int `json:"endPointIpAddressId,omitempty" xmlrpc:"endPointIpAddressId,omitempty"`

	// The [[SoftLayer_Network_Vlan]] id.
	EndPointVlanId *int `json:"endPointVlanId,omitempty" xmlrpc:"endPointVlanId,omitempty"`

	// The [[SoftLayer_Network_Subnet]] id.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// This is the hostname for the router associated with the [[SoftLayer_Network_Subnet|subnet]]. This is a readonly property.
	RouterHostname *string `json:"routerHostname,omitempty" xmlrpc:"routerHostname,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a network ipsec vpn order with SoftLayer.
type Container_Product_Order_Network_Tunnel_Ipsec struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a network vlan order with SoftLayer.
type Container_Product_Order_Network_Vlan struct {
	Container_Product_Order

	// The description which includes the primary router's hostname plus the vlan number.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The datacenter portion of the hostname.
	HostnameDatacenter *string `json:"hostnameDatacenter,omitempty" xmlrpc:"hostnameDatacenter,omitempty"`

	// The router portion of the hostname.
	HostnameRouter *string `json:"hostnameRouter,omitempty" xmlrpc:"hostnameRouter,omitempty"`

	// The [[SoftLayer_Network_Vlan]] id.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The optional name for this VLAN
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The router object on which the new VLAN should be created.
	Router *Hardware `json:"router,omitempty" xmlrpc:"router,omitempty"`

	// The ID of the [[SoftLayer_Hardware_Router]] object on which the new VLAN should be created.
	RouterId *int `json:"routerId,omitempty" xmlrpc:"routerId,omitempty"`

	// The collection of subnets associated with this vlan.
	Subnets []Container_Product_Order `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// The vlan number.
	VlanNumber *int `json:"vlanNumber,omitempty" xmlrpc:"vlanNumber,omitempty"`
}

// This class contains the collections of public and private VLANs that are available during the ordering process.
type Container_Product_Order_Network_Vlans struct {
	Entity

	// The collection of private vlans available during ordering.
	PrivateVlans []Container_Product_Order `json:"privateVlans,omitempty" xmlrpc:"privateVlans,omitempty"`

	// The collection of public vlans available during ordering.
	PublicVlans []Container_Product_Order `json:"publicVlans,omitempty" xmlrpc:"publicVlans,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder when linking a Bluemix account to a newly created SoftLayer account.
type Container_Product_Order_NewCustomerSetup struct {
	Container_Product_Order

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	AuthorizationToken *string `json:"authorizationToken,omitempty" xmlrpc:"authorizationToken,omitempty"`

	// no documentation yet
	ExternalAccountId *string `json:"externalAccountId,omitempty" xmlrpc:"externalAccountId,omitempty"`

	// no documentation yet
	ExternalServiceProviderKey *string `json:"externalServiceProviderKey,omitempty" xmlrpc:"externalServiceProviderKey,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order for Private Cloud.
type Container_Product_Order_Private_Cloud struct {
	Container_Product_Order
}

// This is used for storing various items about the order. Currently used for storing additional raid information when ordering servers. This is optional
type Container_Product_Order_Property struct {
	Entity

	// The property name
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The property value
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// When an order is placed (SoftLayer_Product_Order::placeOrder), a receipt is returned when the order is created successfully. The information in the receipt helps explain information about the order. It's order ID, and all the data within the order as well.
//
// For PayPal Orders, an URL is also returned to the user so that the user can complete the transaction. Users paying with PayPal must continue on to this URL, login and pay. When doing this, PayPal will redirect the user back to a SoftLayer page which will then "finalize" the authorization process. From here, Sales will verify the order by contacting the user in some way, unless sales has already spoken to the user about approving the order.
//
// For users paying with a credit card, a receipt means the order has gone to sales and is awaiting approval.
type Container_Product_Order_Receipt struct {
	Entity

	// This URL refers to the location where you will visit to complete the payment authorization for an external service, such as PayPal. This property is associated with <code>externalPaymentToken</code> and will only be populated when purchasing products with an external service.
	//
	// Once you visit this location, you will be presented with the options to confirm payment or deny payment. If you confirm payment, you will be redirected back to the receipt for your order. If you deny, you will be redirected back to the cancel order page where you do not need to take any additional action.
	//
	// Until you confirm payment with the external service, your products will not be provisioned or accessible for your consumption. Upon successfully confirming payment, our system will be notified and the order approval and provisioning systems will begin processing. After provisioning is complete, your services will be available.
	ExternalPaymentCheckoutUrl *string `json:"externalPaymentCheckoutUrl,omitempty" xmlrpc:"externalPaymentCheckoutUrl,omitempty"`

	// This token refers to the identifier for the external payment authorization. This token is associated with the <code>externalPaymentCheckoutUrl</code> and is only populated when purchasing products with an external service like PayPal.
	ExternalPaymentToken *string `json:"externalPaymentToken,omitempty" xmlrpc:"externalPaymentToken,omitempty"`

	// The date when SoftLayer received the order.
	OrderDate *Time `json:"orderDate,omitempty" xmlrpc:"orderDate,omitempty"`

	// This is a copy of the order container (SoftLayer_Container_Product_Order) which holds all the data related to an order. This will only return when an order is processed successfully. It will contain all the items in an order as well as the order totals.
	OrderDetails *Container_Product_Order `json:"orderDetails,omitempty" xmlrpc:"orderDetails,omitempty"`

	// SoftLayer's unique identifier for the order.
	OrderId *int `json:"orderId,omitempty" xmlrpc:"orderId,omitempty"`

	// Deprecation notice: use <code>externalPaymentCheckoutUrl</code> instead of this property.
	//
	// This URL refers to the location where you will visit to complete the payment authorization for PayPal. This property is associated with <code>paypalToken</code> and will only be populated when purchasing products with PayPal.
	//
	// Once you visit PayPal's site, you will be presented with the options to confirm payment or deny payment. If you confirm payment, you will be redirected back to the receipt for your order. If you deny, you will be redirected back to the cancel order page where you do not need to take any additional action.
	//
	// Until you confirm payment with PayPal, your products will not be provisioned or accessible for your consumption. Upon successfully confirming payment, our system will be notified and the order approval and provisioning systems will begin processing. After provisioning is complete, your services will be available.
	PaypalCheckoutUrl *string `json:"paypalCheckoutUrl,omitempty" xmlrpc:"paypalCheckoutUrl,omitempty"`

	// Deprecation notice: use <code>externalPaymentToken</code> instead of this property.
	//
	// This token refers to the identifier provided when payment is processed via PayPal. This token is associated with the <code>paypalCheckoutUrl</code>.
	PaypalToken *string `json:"paypalToken,omitempty" xmlrpc:"paypalToken,omitempty"`

	// This is a copy of the order that was successfully placed (SoftLayer_Billing_Order). This will only return when an order is processed successfully.
	PlacedOrder *Billing_Order `json:"placedOrder,omitempty" xmlrpc:"placedOrder,omitempty"`

	// This is a copy of the quote container (SoftLayer_Billing_Order_Quote) which holds all the data related to a quote. This will only return when a quote is processed successfully.
	Quote *Billing_Order_Quote `json:"quote,omitempty" xmlrpc:"quote,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype contains everything required to place a secure certificate order with SoftLayer.
type Container_Product_Order_Security_Certificate struct {
	Container_Product_Order

	// The administrator contact associated with a SSL certificate. If the contact is not provided the technical contact will be used. If the address is not provided the organization information address will be used.
	AdministrativeContact *Container_Product_Order_Attribute_Contact `json:"administrativeContact,omitempty" xmlrpc:"administrativeContact,omitempty"`

	// The billing contact associated with a SSL certificate. If the contact is not provided the technical contact will be used. If the address is not provided the organization information address will be used.
	BillingContact *Container_Product_Order_Attribute_Contact `json:"billingContact,omitempty" xmlrpc:"billingContact,omitempty"`

	// The base64 encoded string that sent from an applicant to a certificate authority. The CSR contains information identifying the applicant and the public key chosen by the applicant. The corresponding private key should not be included.
	CertificateSigningRequest *string `json:"certificateSigningRequest,omitempty" xmlrpc:"certificateSigningRequest,omitempty"`

	// The email address that can approve a secure certificate order.
	OrderApproverEmailAddress *string `json:"orderApproverEmailAddress,omitempty" xmlrpc:"orderApproverEmailAddress,omitempty"`

	// The organization information associated with a SSL certificate.
	OrganizationInformation *Container_Product_Order_Attribute_Organization `json:"organizationInformation,omitempty" xmlrpc:"organizationInformation,omitempty"`

	// Indicates if it is an renewal order of an existing SSL certificate.
	RenewalFlag *bool `json:"renewalFlag,omitempty" xmlrpc:"renewalFlag,omitempty"`

	// (DEPRECATED) Do not set this property, as it will always be set to 1.
	// Deprecated: This function has been marked as deprecated.
	ServerCount *int `json:"serverCount,omitempty" xmlrpc:"serverCount,omitempty"`

	// The server type. This is the name from a [[SoftLayer_Security_Certificate_Request_ServerType]] object.
	ServerType *string `json:"serverType,omitempty" xmlrpc:"serverType,omitempty"`

	// The technical contact associated with a SSL certificate. If the address is not provided the organization information address will be used.
	TechnicalContact *Container_Product_Order_Attribute_Contact `json:"technicalContact,omitempty" xmlrpc:"technicalContact,omitempty"`

	// (DEPRECATED) The period that a SSL certificate is valid for.  For example, 12, 24, 36. This property will be set automatically based on the certificate product ordered when verifying or placing orders.
	ValidityMonths *int `json:"validityMonths,omitempty" xmlrpc:"validityMonths,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder.
type Container_Product_Order_Service struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder.
type Container_Product_Order_Service_External struct {
	Container_Product_Order

	// For orders that contain servers (bare metal, virtual server, big data, etc.), the hardware property is required. This property is an array of [[SoftLayer_Hardware]] objects. The <code>hostname</code> and <code>domain</code> properties are required for each hardware object. Note that virtual server ([[SoftLayer_Container_Product_Order_Virtual_Guest]]) orders may populate this field instead of the <code>virtualGuests</code> property.
	ExternalResources []Service_External_Resource `json:"externalResources,omitempty" xmlrpc:"externalResources,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a virtual license order with SoftLayer.
type Container_Product_Order_Software_Component_Virtual struct {
	Container_Product_Order

	// array of ip address ids for which a license should be allocated for.
	EndPointIpAddressIds []int `json:"endPointIpAddressIds,omitempty" xmlrpc:"endPointIpAddressIds,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a hardware security module order with SoftLayer.
type Container_Product_Order_Software_License struct {
	Container_Product_Order
}

// This object holds all of the ssh key ids that will allow authentication to a single server.
type Container_Product_Order_SshKeys struct {
	Entity

	// An array of SoftLayer_Security_Ssh_Key IDs to assign to a server.
	SshKeyIds []int `json:"sshKeyIds,omitempty" xmlrpc:"sshKeyIds,omitempty"`
}

// A single storage group container used for a hardware server order.
//
// This object describes a single storage group that can be added to an order container.
type Container_Product_Order_Storage_Group struct {
	Entity

	// Size of the array in gigabytes. Must be within limitations of the smallest drive assigned to the storage group and the storage group type.
	ArraySize *Float64 `json:"arraySize,omitempty" xmlrpc:"arraySize,omitempty"`

	// The array type id from a [[SoftLayer_Configuration_Storage_Group_Array_Type]] object.
	ArrayTypeId *int `json:"arrayTypeId,omitempty" xmlrpc:"arrayTypeId,omitempty"`

	// Defines the disk controller to put the storage group and the hard drives on.
	//
	// This must match a disk controller price on the order. The disk controller index is 0-indexed. 'disk_controller' = 0 'disk_controller1' = 1 'disk_controller2' = 2
	DiskControllerIndex *int `json:"diskControllerIndex,omitempty" xmlrpc:"diskControllerIndex,omitempty"`

	// String array of category codes for drives to use in the storage group as an alternative to their index positions.
	//
	// This must be specified if ordering a storage group with PCIe drives.
	HardDriveCategoryCodes []string `json:"hardDriveCategoryCodes,omitempty" xmlrpc:"hardDriveCategoryCodes,omitempty"`

	// Integer array of drive indexes to use in the storage group.
	HardDrives []int `json:"hardDrives,omitempty" xmlrpc:"hardDrives,omitempty"`

	// If an array should be protected by an hotspare, the drive index of the hotspares should be here.
	//
	// If a drive is a hotspare for all arrays then a separate storage group with array type GLOBAL_HOT_SPARE should be used
	HotSpareDrives []int `json:"hotSpareDrives,omitempty" xmlrpc:"hotSpareDrives,omitempty"`

	// << EOT
	LvmFlag *bool `json:"lvmFlag,omitempty" xmlrpc:"lvmFlag,omitempty"`

	// The id for a [[SoftLayer_Hardware_Component_Partition_Template]] object, which will determine the partitions to add to the storage group.
	//
	// If this storage group is not a primary storage group, then this will not be used.
	PartitionTemplateId *int `json:"partitionTemplateId,omitempty" xmlrpc:"partitionTemplateId,omitempty"`

	// Defines the partitions for the storage group.
	//
	// If this storage group is not a secondary storage group, then this will not be used.
	Partitions []Container_Product_Order_Storage_Group_Partition `json:"partitions,omitempty" xmlrpc:"partitions,omitempty"`
}

// A storage group partition container used for a hardware server order.
//
// This object describes the partitions for a single storage group that can be added to an order container.
type Container_Product_Order_Storage_Group_Partition struct {
	Entity

	// Is this a grow partition
	IsGrow *bool `json:"isGrow,omitempty" xmlrpc:"isGrow,omitempty"`

	// The name of this partition
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The size of this partition
	Size *Float64 `json:"size,omitempty" xmlrpc:"size,omitempty"`
}

// When ordering paid support this datatype needs to be populated and sent to SoftLayer_Product_Order::placeOrder.
type Container_Product_Order_Support struct {
	Container_Product_Order
}

// This container type is used for placing orders for external authentication, such as phone-based authentication.
type Container_Product_Order_User_Customer_External_Binding struct {
	Container_Product_Order

	// The external id that access to external authentication is being purchased for.
	ExternalId *string `json:"externalId,omitempty" xmlrpc:"externalId,omitempty"`

	// The SoftLayer [[SoftLayer_User_Customer|user]] identifier that an external binding is being purchased for.
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`

	// The [[SoftLayer_User_Customer_External_Binding_Vendor|vendor]] identifier for the external binding being purchased.
	VendorId *int `json:"vendorId,omitempty" xmlrpc:"vendorId,omitempty"`
}

// This is the default container type for Dedicated Virtual Host orders.
type Container_Product_Order_Virtual_DedicatedHost struct {
	Container_Product_Order
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place a Portable Storage order with SoftLayer.
type Container_Product_Order_Virtual_Disk_Image struct {
	Container_Product_Order

	// Label for the portable storage volume.
	DiskDescription *string `json:"diskDescription,omitempty" xmlrpc:"diskDescription,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Virtual_Guest struct {
	Container_Product_Order_Hardware_Server

	// The mode used to boot the [[SoftLayer_Virtual_Guest]].  Supported values are 'PV' and 'HVM'.
	BootMode *string `json:"bootMode,omitempty" xmlrpc:"bootMode,omitempty"`

	// Identifier of the [[SoftLayer_Virtual_Disk_Image]] to boot from.
	BootableDiskId *int `json:"bootableDiskId,omitempty" xmlrpc:"bootableDiskId,omitempty"`

	// Identifier of [[SoftLayer_Virtual_DedicatedHost]] to order
	HostId *int `json:"hostId,omitempty" xmlrpc:"hostId,omitempty"`

	// Identifier of [[SoftLayer_Virtual_ReservedCapacityGroup]] to order
	ReservedCapacityId *int `json:"reservedCapacityId,omitempty" xmlrpc:"reservedCapacityId,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Product_Order::placeOrder. This datatype has everything required to place an order with SoftLayer.
type Container_Product_Order_Virtual_Guest_Upgrade struct {
	Container_Product_Order_Virtual_Guest
}

// no documentation yet
type Container_Product_Order_Virtual_Guest_Vpc struct {
	Container_Product_Order_Virtual_Guest

	// no documentation yet
	AdditionalNetworkInterfaces []Container_Product_Order_Virtual_Guest_Vpc_NetworkInterface `json:"additionalNetworkInterfaces,omitempty" xmlrpc:"additionalNetworkInterfaces,omitempty"`

	// no documentation yet
	Crn *string `json:"crn,omitempty" xmlrpc:"crn,omitempty"`

	// no documentation yet
	InstanceProfile *string `json:"instanceProfile,omitempty" xmlrpc:"instanceProfile,omitempty"`

	// no documentation yet
	IpAllocations []Container_Product_Order_Vpc_IpAllocation `json:"ipAllocations,omitempty" xmlrpc:"ipAllocations,omitempty"`

	// no documentation yet
	OverlayNetworkFlag *bool `json:"overlayNetworkFlag,omitempty" xmlrpc:"overlayNetworkFlag,omitempty"`

	// no documentation yet
	ResourceGroup *string `json:"resourceGroup,omitempty" xmlrpc:"resourceGroup,omitempty"`

	// no documentation yet
	ServerId *string `json:"serverId,omitempty" xmlrpc:"serverId,omitempty"`

	// no documentation yet
	ServicePortCidr *string `json:"servicePortCidr,omitempty" xmlrpc:"servicePortCidr,omitempty"`

	// no documentation yet
	ServicePortDns []string `json:"servicePortDns,omitempty" xmlrpc:"servicePortDns,omitempty"`

	// no documentation yet
	ServicePortGateway *string `json:"servicePortGateway,omitempty" xmlrpc:"servicePortGateway,omitempty"`

	// no documentation yet
	ServicePortInterfaceId *string `json:"servicePortInterfaceId,omitempty" xmlrpc:"servicePortInterfaceId,omitempty"`

	// no documentation yet
	ServicePortIpAddress *string `json:"servicePortIpAddress,omitempty" xmlrpc:"servicePortIpAddress,omitempty"`

	// no documentation yet
	ServicePortIpAllocationId *string `json:"servicePortIpAllocationId,omitempty" xmlrpc:"servicePortIpAllocationId,omitempty"`

	// no documentation yet
	ServicePortVpcId *string `json:"servicePortVpcId,omitempty" xmlrpc:"servicePortVpcId,omitempty"`

	// no documentation yet
	StorageVolumes []Container_Product_Order_Virtual_Guest_Vpc_StorageVolume `json:"storageVolumes,omitempty" xmlrpc:"storageVolumes,omitempty"`

	// no documentation yet
	Subnets []Container_Product_Order_Vpc_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// no documentation yet
	Zone *string `json:"zone,omitempty" xmlrpc:"zone,omitempty"`
}

// no documentation yet
type Container_Product_Order_Virtual_Guest_Vpc_NetworkInterface struct {
	Entity

	// no documentation yet
	Cidr *string `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// no documentation yet
	Dns []string `json:"dns,omitempty" xmlrpc:"dns,omitempty"`

	// no documentation yet
	Gateway *string `json:"gateway,omitempty" xmlrpc:"gateway,omitempty"`

	// no documentation yet
	InterfaceId *string `json:"interfaceId,omitempty" xmlrpc:"interfaceId,omitempty"`

	// no documentation yet
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// no documentation yet
	IpAllocationId *string `json:"ipAllocationId,omitempty" xmlrpc:"ipAllocationId,omitempty"`

	// no documentation yet
	SecurityGroupIds []int `json:"securityGroupIds,omitempty" xmlrpc:"securityGroupIds,omitempty"`

	// no documentation yet
	SubnetId *string `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`

	// no documentation yet
	VpcId *string `json:"vpcId,omitempty" xmlrpc:"vpcId,omitempty"`
}

// no documentation yet
type Container_Product_Order_Virtual_Guest_Vpc_StorageVolume struct {
	Entity

	// no documentation yet
	AttachmentName *string `json:"attachmentName,omitempty" xmlrpc:"attachmentName,omitempty"`

	// no documentation yet
	Capacity *int `json:"capacity,omitempty" xmlrpc:"capacity,omitempty"`

	// no documentation yet
	DeleteOnReclaim *bool `json:"deleteOnReclaim,omitempty" xmlrpc:"deleteOnReclaim,omitempty"`

	// no documentation yet
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Index *int `json:"index,omitempty" xmlrpc:"index,omitempty"`

	// no documentation yet
	Iops *int `json:"iops,omitempty" xmlrpc:"iops,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Profile *string `json:"profile,omitempty" xmlrpc:"profile,omitempty"`

	// no documentation yet
	ResourceGroup *string `json:"resourceGroup,omitempty" xmlrpc:"resourceGroup,omitempty"`

	// no documentation yet
	RootKeyCrn *string `json:"rootKeyCrn,omitempty" xmlrpc:"rootKeyCrn,omitempty"`
}

// no documentation yet
type Container_Product_Order_Virtual_Guest_Vpc_Upgrade struct {
	Container_Product_Order_Virtual_Guest_Vpc
}

// This is the default container type for Reserved Capacity orders.
type Container_Product_Order_Virtual_ReservedCapacity struct {
	Container_Product_Order

	// Identifier of [[SoftLayer_Hardware_Router]] on which the capacity will be
	BackendRouterId *int `json:"backendRouterId,omitempty" xmlrpc:"backendRouterId,omitempty"`

	// Name for the [[SoftLayer_Virtual_ReservedCapacityGroup]] being ordered.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Container_Product_Order_Vpc_IpAllocation struct {
	Entity

	// no documentation yet
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Ip *string `json:"ip,omitempty" xmlrpc:"ip,omitempty"`
}

// no documentation yet
type Container_Product_Order_Vpc_Subnet struct {
	Entity

	// no documentation yet
	Cidr *string `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// no documentation yet
	Dns *string `json:"dns,omitempty" xmlrpc:"dns,omitempty"`

	// no documentation yet
	Gateway *string `json:"gateway,omitempty" xmlrpc:"gateway,omitempty"`

	// no documentation yet
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Vlan *int `json:"vlan,omitempty" xmlrpc:"vlan,omitempty"`
}

// The SoftLayer_Container_Product_Promotion data type contains information about a promotion and its requirements.
type Container_Product_Promotion struct {
	Entity

	// no documentation yet
	Code *string `json:"code,omitempty" xmlrpc:"code,omitempty"`

	// no documentation yet
	ExpirationDate *Time `json:"expirationDate,omitempty" xmlrpc:"expirationDate,omitempty"`

	// no documentation yet
	Locations []Location `json:"locations,omitempty" xmlrpc:"locations,omitempty"`

	// no documentation yet
	RequirementGroups []Container_Product_Promotion_RequirementGroup `json:"requirementGroups,omitempty" xmlrpc:"requirementGroups,omitempty"`
}

// The SoftLayer_Container_Product_Promotion_RequirementGroup data type contains the required options that must be present on an order for the promotion to be applied. At least one of the categories, presets, or prices must be on the order.
type Container_Product_Promotion_RequirementGroup struct {
	Entity

	// The category options to choose from for this requirement group
	Categories []Product_Item_Category `json:"categories,omitempty" xmlrpc:"categories,omitempty"`

	// The preset options to choose from for this requirement group
	Presets []Product_Package_Preset `json:"presets,omitempty" xmlrpc:"presets,omitempty"`

	// The price options to choose from for this requirement group
	Prices []Product_Item_Price `json:"prices,omitempty" xmlrpc:"prices,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Provisioning_Maintenance_Window::addCustomerUpgradeWindow. This datatype has everything required to place an order with SoftLayer.
type Container_Provisioning_Maintenance_Window struct {
	Entity

	// Maintenance classifications.
	ClassificationIds []Provisioning_Maintenance_Classification `json:"classificationIds,omitempty" xmlrpc:"classificationIds,omitempty"`

	// Maintenance classifications.
	ItemCategoryIds []Product_Item_Category `json:"itemCategoryIds,omitempty" xmlrpc:"itemCategoryIds,omitempty"`

	// The maintenance window id
	MaintenanceWindowId *int `json:"maintenanceWindowId,omitempty" xmlrpc:"maintenanceWindowId,omitempty"`

	// Maintenance window ticket id
	TicketId *int `json:"ticketId,omitempty" xmlrpc:"ticketId,omitempty"`

	// Maintenance window date
	WindowMaintenanceDate *Time `json:"windowMaintenanceDate,omitempty" xmlrpc:"windowMaintenanceDate,omitempty"`
}

// no documentation yet
type Container_Referral_Partner_Commission struct {
	Entity

	// no documentation yet
	CommissionAmount *Float64 `json:"commissionAmount,omitempty" xmlrpc:"commissionAmount,omitempty"`

	// no documentation yet
	CommissionRate *Float64 `json:"commissionRate,omitempty" xmlrpc:"commissionRate,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	ReferralAccountId *int `json:"referralAccountId,omitempty" xmlrpc:"referralAccountId,omitempty"`

	// no documentation yet
	ReferralCompanyName *string `json:"referralCompanyName,omitempty" xmlrpc:"referralCompanyName,omitempty"`

	// no documentation yet
	ReferralPartnerAccountId *int `json:"referralPartnerAccountId,omitempty" xmlrpc:"referralPartnerAccountId,omitempty"`

	// no documentation yet
	ReferralRevenue *Float64 `json:"referralRevenue,omitempty" xmlrpc:"referralRevenue,omitempty"`
}

// no documentation yet
type Container_Referral_Partner_Payment_Option struct {
	Entity

	// no documentation yet
	AccountNumber *string `json:"accountNumber,omitempty" xmlrpc:"accountNumber,omitempty"`

	// no documentation yet
	AccountType *string `json:"accountType,omitempty" xmlrpc:"accountType,omitempty"`

	// no documentation yet
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// no documentation yet
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// no documentation yet
	BankTransitNumber *string `json:"bankTransitNumber,omitempty" xmlrpc:"bankTransitNumber,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	FederalTaxId *string `json:"federalTaxId,omitempty" xmlrpc:"federalTaxId,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	PaymentType *string `json:"paymentType,omitempty" xmlrpc:"paymentType,omitempty"`

	// no documentation yet
	PaypalEmail *string `json:"paypalEmail,omitempty" xmlrpc:"paypalEmail,omitempty"`

	// no documentation yet
	PhoneNumber *string `json:"phoneNumber,omitempty" xmlrpc:"phoneNumber,omitempty"`

	// no documentation yet
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// no documentation yet
type Container_Referral_Partner_Prospect struct {
	Entity

	// no documentation yet
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// no documentation yet
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// no documentation yet
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// no documentation yet
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// no documentation yet
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// no documentation yet
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// no documentation yet
	OfficePhone *string `json:"officePhone,omitempty" xmlrpc:"officePhone,omitempty"`

	// no documentation yet
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// no documentation yet
	Questions []string `json:"questions,omitempty" xmlrpc:"questions,omitempty"`

	// no documentation yet
	Responses []Survey_Response `json:"responses,omitempty" xmlrpc:"responses,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	SurveyId *string `json:"surveyId,omitempty" xmlrpc:"surveyId,omitempty"`
}

// The SoftLayer_Container_RemoteManagement_Graphs_SensorSpeed contains graphs to  display speed for each of the server's fans.  Fan speeds are gathered from the server's remote management card.
type Container_RemoteManagement_Graphs_SensorSpeed struct {
	Entity

	// The graph to display the server's fan speed.
	Graph *[]byte `json:"graph,omitempty" xmlrpc:"graph,omitempty"`

	// A title that may be used to display for the graph.
	Title *string `json:"title,omitempty" xmlrpc:"title,omitempty"`
}

// The SoftLayer_Container_RemoteManagement_Graphs_SensorTemperature contains graphs to display the cpu(s) and system temperatures retrieved from the management card using thermometer graphs.
type Container_RemoteManagement_Graphs_SensorTemperature struct {
	Entity

	// The graph to display the server's cpu(s) and system temperatures.
	Graph *[]byte `json:"graph,omitempty" xmlrpc:"graph,omitempty"`

	// A title that may be used to display for the graph.
	Title *string `json:"title,omitempty" xmlrpc:"title,omitempty"`
}

// The SoftLayer_Container_RemoteManagement_PmInfo contains pminfo information retrieved from a server's remote management card.
type Container_RemoteManagement_PmInfo struct {
	Entity

	// PmInfo ID
	PmInfoId *string `json:"pmInfoId,omitempty" xmlrpc:"pmInfoId,omitempty"`

	// PmInfo Reading
	PmInfoReading *string `json:"pmInfoReading,omitempty" xmlrpc:"pmInfoReading,omitempty"`
}

// The SoftLayer_Container_RemoteManagement_SensorReadings contains sensor information retrieved from a server's remote management card.
type Container_RemoteManagement_SensorReading struct {
	Entity

	// Lower Non-Recoverable threshold
	LowerCritical *string `json:"lowerCritical,omitempty" xmlrpc:"lowerCritical,omitempty"`

	// Lower Non-Critical threshold
	LowerNonCritical *string `json:"lowerNonCritical,omitempty" xmlrpc:"lowerNonCritical,omitempty"`

	// Lower Non-Recoverable threshold
	LowerNonRecoverable *string `json:"lowerNonRecoverable,omitempty" xmlrpc:"lowerNonRecoverable,omitempty"`

	// Sensor ID
	SensorId *string `json:"sensorId,omitempty" xmlrpc:"sensorId,omitempty"`

	// Sensor Reading
	SensorReading *string `json:"sensorReading,omitempty" xmlrpc:"sensorReading,omitempty"`

	// Sensor Units
	SensorUnits *string `json:"sensorUnits,omitempty" xmlrpc:"sensorUnits,omitempty"`

	// Sensor Status
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Upper Critical threshold
	UpperCritical *string `json:"upperCritical,omitempty" xmlrpc:"upperCritical,omitempty"`

	// Upper Non-Critical threshold
	UpperNonCritical *string `json:"upperNonCritical,omitempty" xmlrpc:"upperNonCritical,omitempty"`

	// Upper Non-Recoverable threshold
	UpperNonRecoverable *string `json:"upperNonRecoverable,omitempty" xmlrpc:"upperNonRecoverable,omitempty"`
}

// The SoftLayer_Container_RemoteManagement_SensorReadingsWithGraphs contains the raw data retrieved from a server's remote management card.  Along with the raw data, two sets of graphs will be returned.  One set of graphs is used to display, using thermometer graphs, the temperatures (cpu(s) and system) retrieved from the management card.  The other set is used to display speed for each of the server's fans.
type Container_RemoteManagement_SensorReadingsWithGraphs struct {
	Entity

	// The raw data returned from the server's remote management card.
	RawData []Container_RemoteManagement_SensorReading `json:"rawData,omitempty" xmlrpc:"rawData,omitempty"`

	// The graph(s) to display the server's fan speeds.
	SpeedGraphs []Container_RemoteManagement_Graphs_SensorSpeed `json:"speedGraphs,omitempty" xmlrpc:"speedGraphs,omitempty"`

	// The graph(s) to display the server's cpu(s) and system temperatures.
	TemperatureGraphs []Container_RemoteManagement_Graphs_SensorTemperature `json:"temperatureGraphs,omitempty" xmlrpc:"temperatureGraphs,omitempty"`
}

// The metadata service resource container is used to store information about a single service resource.
type Container_Resource_Metadata_ServiceResource struct {
	Entity

	// The backend IP address for this resource
	BackendIpAddress *string `json:"backendIpAddress,omitempty" xmlrpc:"backendIpAddress,omitempty"`

	// The type for this resource
	Type *Network_Service_Resource_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// This data type is a container that stores information about a single indexed object type.  Object type information can be used for discovery of searchable data and for creation or validation of object index search strings.  Each of these containers holds a collection of <b>[[SoftLayer_Container_Search_ObjectType_Property (type)|SoftLayer_Container_Search_ObjectType_Property]]</b> objects, specifying which object properties are exposed for the current user.  Refer to the the documentation for the <b>[[SoftLayer_Search/search|search()]]</b> method for information on using object types in search strings.
type Container_Search_ObjectType struct {
	Entity

	// Name of object type.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A collection of [[SoftLayer_Container_Search_ObjectType_Property|object properties]].
	Properties []Container_Search_ObjectType_Property `json:"properties,omitempty" xmlrpc:"properties,omitempty"`
}

// This data type is a container that stores information about a single property of a searchable object type.  Each <b>[[SoftLayer_Container_Search_ObjectType (type)|SoftLayer_Container_Search_ObjectType]]</b> object holds a collection of these properties.  Property information can be used for discovery of searchable data and for the creation or validation of object index search strings.  Note that properties are only understood by the <b>[[SoftLayer_Search/advancedSearch|advancedSearch()]]</b> method.  Refer to the <b>advancedSearch()</b> method for information on using properties in search strings.
type Container_Search_ObjectType_Property struct {
	Entity

	// Name of property.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Indicates if this property can be sorted.
	SortableFlag *bool `json:"sortableFlag,omitempty" xmlrpc:"sortableFlag,omitempty"`

	// Property data type.  Valid values include 'boolean', 'integer', 'date', 'string' or 'text'.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Container_Search_Result data type represents a result row from an execution of Search service.
type Container_Search_Result struct {
	Entity

	// An array of terms that were matched in the resource object.
	MatchedTerms []string `json:"matchedTerms,omitempty" xmlrpc:"matchedTerms,omitempty"`

	// The score ratio of the result for relevance to the search criteria.
	RelevanceScore *Float64 `json:"relevanceScore,omitempty" xmlrpc:"relevanceScore,omitempty"`

	// A search results resource object that matched search criteria.
	Resource interface{} `json:"resource,omitempty" xmlrpc:"resource,omitempty"`

	// The type of the resource object that matched search criteria.
	ResourceType *string `json:"resourceType,omitempty" xmlrpc:"resourceType,omitempty"`
}

// The SoftLayer_Container_Software_Component_HostIps_Policy container holds the title and value of a current host ips policy.
type Container_Software_Component_HostIps_Policy struct {
	Entity

	// The value of a host ips category.
	Policy *string `json:"policy,omitempty" xmlrpc:"policy,omitempty"`

	// The category title of a host ips policy.
	PolicyTitle *string `json:"policyTitle,omitempty" xmlrpc:"policyTitle,omitempty"`
}

// These are the results of a tax calculation. The tax calculation was kicked off but allowed to run in the background. This type stores the results so that an interface can be updated with up-to-date information.
type Container_Tax_Cache struct {
	Entity

	// The percentage of the final total that should be tax.
	EffectiveTaxRate *Float64 `json:"effectiveTaxRate,omitempty" xmlrpc:"effectiveTaxRate,omitempty"`

	// no documentation yet
	FailureMessage *string `json:"failureMessage,omitempty" xmlrpc:"failureMessage,omitempty"`

	// The container that holds the four actual tax rates, one for each fee type.
	Items []Container_Tax_Cache_Item `json:"items,omitempty" xmlrpc:"items,omitempty"`

	// The status of the tax request. This should be PENDING, FAILED, or COMPLETED.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The final amount of tax for the order.
	TotalTaxAmount *Float64 `json:"totalTaxAmount,omitempty" xmlrpc:"totalTaxAmount,omitempty"`
}

// This represents one order item in a tax calculation.
type Container_Tax_Cache_Item struct {
	Entity

	// The category code for the referenced product.
	CategoryCode *string `json:"categoryCode,omitempty" xmlrpc:"categoryCode,omitempty"`

	// This hash will match to the hash on an order container.
	ContainerHash *string `json:"containerHash,omitempty" xmlrpc:"containerHash,omitempty"`

	// The reference to the price for this order item.
	ItemPriceId *int `json:"itemPriceId,omitempty" xmlrpc:"itemPriceId,omitempty"`

	// This is the container containing the individual tax rates.
	TaxRates *Container_Tax_Rates `json:"taxRates,omitempty" xmlrpc:"taxRates,omitempty"`
}

// This contains the four tax rates, one for each fee type.
type Container_Tax_Rates struct {
	Entity

	// The tax rate associated with the labor fee.
	LaborTaxRate *Float64 `json:"laborTaxRate,omitempty" xmlrpc:"laborTaxRate,omitempty"`

	// A reference to a location.
	LocationId *Float64 `json:"locationId,omitempty" xmlrpc:"locationId,omitempty"`

	// The tax rate associated with the one-time fee.
	OneTimeTaxRate *Float64 `json:"oneTimeTaxRate,omitempty" xmlrpc:"oneTimeTaxRate,omitempty"`

	// The tax rate associated with the recurring fee.
	RecurringTaxRate *Float64 `json:"recurringTaxRate,omitempty" xmlrpc:"recurringTaxRate,omitempty"`

	// The tax rate associated with the setup fee.
	SetupTaxRate *Float64 `json:"setupTaxRate,omitempty" xmlrpc:"setupTaxRate,omitempty"`
}

// SoftLayer_Container_Ticket_GraphInputs models a single inbound object for a given ticket graph.
type Container_Ticket_GraphInputs struct {
	Entity

	// This is a unix timestamp that represents the stop date/time for a graph.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The front-end or back-end network uplink interface associated with this server.
	NetworkInterfaceId *int `json:"networkInterfaceId,omitempty" xmlrpc:"networkInterfaceId,omitempty"`

	// *
	Pod *int `json:"pod,omitempty" xmlrpc:"pod,omitempty"`

	// This is a human readable name for the server or rack being graphed.
	ServerName *string `json:"serverName,omitempty" xmlrpc:"serverName,omitempty"`

	// This is a unix timestamp that represents the begin date/time for a graph.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`
}

// SoftLayer_Container_Ticket_GraphOutputs models a single outbound object for a given bandwidth graph.
type Container_Ticket_GraphOutputs struct {
	Entity

	// The raw PNG binary data to be displayed once the graph is drawn.
	GraphImage *[]byte `json:"graphImage,omitempty" xmlrpc:"graphImage,omitempty"`

	// The title that ended up being displayed as part of the graph image.
	GraphTitle *string `json:"graphTitle,omitempty" xmlrpc:"graphTitle,omitempty"`

	// The maximum date included in this graph.
	MaxEndDate *Time `json:"maxEndDate,omitempty" xmlrpc:"maxEndDate,omitempty"`

	// The minimum date included in this graph.
	MinStartDate *Time `json:"minStartDate,omitempty" xmlrpc:"minStartDate,omitempty"`
}

// no documentation yet
type Container_Ticket_Priority struct {
	Entity

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Value *int `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Container_Ticket_Survey_Preference struct {
	Entity

	// no documentation yet
	Applicable *bool `json:"applicable,omitempty" xmlrpc:"applicable,omitempty"`

	// no documentation yet
	OptedOut *bool `json:"optedOut,omitempty" xmlrpc:"optedOut,omitempty"`

	// no documentation yet
	OptedOutDate *Time `json:"optedOutDate,omitempty" xmlrpc:"optedOutDate,omitempty"`
}

// Container class used to hold user authentication token
type Container_User_Authentication_Token struct {
	Entity

	// hash that gets populated for user authentication
	Hash *string `json:"hash,omitempty" xmlrpc:"hash,omitempty"`

	// the user authenticated object
	User *User_Customer `json:"user,omitempty" xmlrpc:"user,omitempty"`

	// the id of the user to authenticate
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`
}

// Container classed used to hold external authentication information
type Container_User_Customer_External_Binding struct {
	Entity

	// The unique token that is created by an external authentication request.
	AuthenticationToken *string `json:"authenticationToken,omitempty" xmlrpc:"authenticationToken,omitempty"`

	// Added by softlayer-go. This hints to the API what kind of binding this is.
	ComplexType *string `json:"complexType,omitempty" xmlrpc:"complexType,omitempty"`

	// The OpenID Connect access token which provides access to a resource by the OpenID Connect provider.
	OpenIdConnectAccessToken *string `json:"openIdConnectAccessToken,omitempty" xmlrpc:"openIdConnectAccessToken,omitempty"`

	// The account to login to, if not provided a default will be used.
	OpenIdConnectAccountId *int `json:"openIdConnectAccountId,omitempty" xmlrpc:"openIdConnectAccountId,omitempty"`

	// The OpenID Connect provider type, as a string.
	OpenIdConnectProvider *string `json:"openIdConnectProvider,omitempty" xmlrpc:"openIdConnectProvider,omitempty"`

	// Your SoftLayer customer portal user's portal password.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// A second security code that is only required if your credential has become unsynchronized.
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// The security code used to validate a VeriSign credential.
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`

	// The answer to your security question.
	SecurityQuestionAnswer *string `json:"securityQuestionAnswer,omitempty" xmlrpc:"securityQuestionAnswer,omitempty"`

	// A security question you wish to answer when authenticating to the SoftLayer customer portal. This parameter isn't required if no security questions are set on your portal account or if your account is configured to not require answering a security account upon login.
	SecurityQuestionId *int `json:"securityQuestionId,omitempty" xmlrpc:"securityQuestionId,omitempty"`

	// The username you wish to authenticate to the SoftLayer customer portal with.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// The name of the vendor that will be used for external authentication
	Vendor *string `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`
}

// Container classed used to hold portal token
type Container_User_Customer_External_Binding_Totp struct {
	Container_User_Customer_External_Binding

	// The security code used to validate a Totp credential.
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`
}

// Container classed used to hold details about an external authentication vendor.
type Container_User_Customer_External_Binding_Vendor struct {
	Entity

	// The keyname used to identify an external authentication vendor.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of an external authentication vendor.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// Container classed used to hold portal token
type Container_User_Customer_External_Binding_Verisign struct {
	Container_User_Customer_External_Binding

	// A second security code that is only required if your credential has become unsynchronized.
	SecondSecurityCode *string `json:"secondSecurityCode,omitempty" xmlrpc:"secondSecurityCode,omitempty"`

	// The security code used to validate a VeriSign credential.
	SecurityCode *string `json:"securityCode,omitempty" xmlrpc:"securityCode,omitempty"`
}

// no documentation yet
type Container_User_Customer_OpenIdConnect_LoginAccountInfo struct {
	Entity

	// The customer account's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The company name associated with an account.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Container_User_Customer_OpenIdConnect_MigrationState struct {
	Entity

	// The number of days remaining in the grace period for this user's account to
	DaysToGracePeriodEnd *int `json:"daysToGracePeriodEnd,omitempty" xmlrpc:"daysToGracePeriodEnd,omitempty"`

	// Flag for whether the email address inside this SoftLayer_User_Customer object
	EmailAlreadyUsedForInvitationToAccount *bool `json:"emailAlreadyUsedForInvitationToAccount,omitempty" xmlrpc:"emailAlreadyUsedForInvitationToAccount,omitempty"`

	// Flag for whether the email address inside this SoftLayer_User_Customer object
	EmailAlreadyUsedForLinkToAccount *bool `json:"emailAlreadyUsedForLinkToAccount,omitempty" xmlrpc:"emailAlreadyUsedForLinkToAccount,omitempty"`

	// The IBMid email address where an invitation was sent.
	ExistingInvitationOpenIdConnectName *string `json:"existingInvitationOpenIdConnectName,omitempty" xmlrpc:"existingInvitationOpenIdConnectName,omitempty"`

	// Flag for whether the account is OpenIdConnect authenticated or not.
	IsAccountOpenIdConnectAuthenticated *bool `json:"isAccountOpenIdConnectAuthenticated,omitempty" xmlrpc:"isAccountOpenIdConnectAuthenticated,omitempty"`
}

// Container for holding information necessary for the setting and resetting of customer passwords
type Container_User_Customer_PasswordSet struct {
	Entity

	// Id of SoftLayer_User_Security_Question.
	AnsweredSecurityQuestionId *int `json:"answeredSecurityQuestionId,omitempty" xmlrpc:"answeredSecurityQuestionId,omitempty"`

	// The authentication methods required.
	AuthenticationMethods []int `json:"authenticationMethods,omitempty" xmlrpc:"authenticationMethods,omitempty"`

	// The number of digits required.
	DigitCountRequirement *int `json:"digitCountRequirement,omitempty" xmlrpc:"digitCountRequirement,omitempty"`

	// The password key provided to user in the password set url link sent via email.
	Key *string `json:"key,omitempty" xmlrpc:"key,omitempty"`

	// The number of lowercase letters required.
	LowercaseCountRequirement *int `json:"lowercaseCountRequirement,omitempty" xmlrpc:"lowercaseCountRequirement,omitempty"`

	// The maximum password length requirement.
	MaximumPasswordLengthRequirement *int `json:"maximumPasswordLengthRequirement,omitempty" xmlrpc:"maximumPasswordLengthRequirement,omitempty"`

	// The minimum password length requirement.
	MinimumPasswordLengthRequirement *int `json:"minimumPasswordLengthRequirement,omitempty" xmlrpc:"minimumPasswordLengthRequirement,omitempty"`

	// The user's new password.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// Answer to security question provided by the user.
	SecurityAnswer *string `json:"securityAnswer,omitempty" xmlrpc:"securityAnswer,omitempty"`

	// Array of SoftLayer_User_Security_Question.
	SecurityQuestions []User_Security_Question `json:"securityQuestions,omitempty" xmlrpc:"securityQuestions,omitempty"`

	// The number of special characters required.
	SpecialCharacterCountRequirement *int `json:"specialCharacterCountRequirement,omitempty" xmlrpc:"specialCharacterCountRequirement,omitempty"`

	// List of the allowed special characters.
	SpecialCharactersAllowed *string `json:"specialCharactersAllowed,omitempty" xmlrpc:"specialCharactersAllowed,omitempty"`

	// The number of uppercase letters required.
	UppercaseCountRequirement *int `json:"uppercaseCountRequirement,omitempty" xmlrpc:"uppercaseCountRequirement,omitempty"`

	// The id of the user to authenticate.
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`
}

// Container classed used to hold mobile portal token
type Container_User_Customer_Portal_MobileToken struct {
	Container_User_Customer_Portal_Token

	// True if this user login required an external binding.
	HasExternalBinding *bool `json:"hasExternalBinding,omitempty" xmlrpc:"hasExternalBinding,omitempty"`
}

// Container classed used to hold portal token
type Container_User_Customer_Portal_Token struct {
	Entity

	// hash of logged in user session id
	Hash *string `json:"hash,omitempty" xmlrpc:"hash,omitempty"`

	// the logged in user data
	User *User_Customer `json:"user,omitempty" xmlrpc:"user,omitempty"`

	// the id of the logged in user
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`
}

// no documentation yet
type Container_User_Customer_Profile_Event_HyperWarp_ProfileChange struct {
	Entity

	// no documentation yet
	Account_id *string `json:"account_id,omitempty" xmlrpc:"account_id,omitempty"`

	// no documentation yet
	Context *Container_User_Customer_Profile_Event_HyperWarp_ProfileChange_Context `json:"context,omitempty" xmlrpc:"context,omitempty"`

	// no documentation yet
	Event_id *string `json:"event_id,omitempty" xmlrpc:"event_id,omitempty"`

	// no documentation yet
	Event_properties *Container_User_Customer_Profile_Event_HyperWarp_ProfileChange_EventProperties `json:"event_properties,omitempty" xmlrpc:"event_properties,omitempty"`

	// no documentation yet
	Event_type *string `json:"event_type,omitempty" xmlrpc:"event_type,omitempty"`

	// no documentation yet
	Publisher *string `json:"publisher,omitempty" xmlrpc:"publisher,omitempty"`

	// no documentation yet
	Timestamp *string `json:"timestamp,omitempty" xmlrpc:"timestamp,omitempty"`

	// no documentation yet
	Version *string `json:"version,omitempty" xmlrpc:"version,omitempty"`
}

// no documentation yet
type Container_User_Customer_Profile_Event_HyperWarp_ProfileChange_Context struct {
	Entity

	// no documentation yet
	Previous_values *Container_User_Customer_Profile_Event_HyperWarp_ProfileChange_EventProperties `json:"previous_values,omitempty" xmlrpc:"previous_values,omitempty"`
}

// no documentation yet
type Container_User_Customer_Profile_Event_HyperWarp_ProfileChange_EventProperties struct {
	Entity

	// no documentation yet
	Allowed_ip_addresses *string `json:"allowed_ip_addresses,omitempty" xmlrpc:"allowed_ip_addresses,omitempty"`

	// no documentation yet
	Altphonenumber *string `json:"altphonenumber,omitempty" xmlrpc:"altphonenumber,omitempty"`

	// no documentation yet
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// no documentation yet
	Firstname *string `json:"firstname,omitempty" xmlrpc:"firstname,omitempty"`

	// no documentation yet
	Iam_id *string `json:"iam_id,omitempty" xmlrpc:"iam_id,omitempty"`

	// no documentation yet
	Language *string `json:"language,omitempty" xmlrpc:"language,omitempty"`

	// no documentation yet
	Lastname *string `json:"lastname,omitempty" xmlrpc:"lastname,omitempty"`

	// no documentation yet
	Notification_language *string `json:"notification_language,omitempty" xmlrpc:"notification_language,omitempty"`

	// no documentation yet
	Origin *string `json:"origin,omitempty" xmlrpc:"origin,omitempty"`

	// no documentation yet
	Phonenumber *string `json:"phonenumber,omitempty" xmlrpc:"phonenumber,omitempty"`

	// no documentation yet
	Photo *string `json:"photo,omitempty" xmlrpc:"photo,omitempty"`

	// no documentation yet
	Realm *string `json:"realm,omitempty" xmlrpc:"realm,omitempty"`

	// no documentation yet
	Self_manage *bool `json:"self_manage,omitempty" xmlrpc:"self_manage,omitempty"`

	// no documentation yet
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`

	// no documentation yet
	Substate *string `json:"substate,omitempty" xmlrpc:"substate,omitempty"`

	// no documentation yet
	User_id *string `json:"user_id,omitempty" xmlrpc:"user_id,omitempty"`
}

// Container classed used to hold portal token
type Container_User_Employee_External_Binding_Verisign struct {
	Entity
}

// At times,such as when attaching files to tickets, it is necessary to send files to SoftLayer API methods. The SoftLayer_Container_Utility_File_Attachment data type models a single file to upload to the API.
type Container_Utility_File_Attachment struct {
	Entity

	// The contents of a file that is uploaded to the SoftLayer API.
	Data *[]byte `json:"data,omitempty" xmlrpc:"data,omitempty"`

	// The name of a file that is uploaded to the SoftLayer API.
	Filename *string `json:"filename,omitempty" xmlrpc:"filename,omitempty"`
}

// SoftLayer_Container_Utility_File_Entity data type models a single entity on a storage resource. Entities can include anything within a storage volume including files, folders, directories, and CloudLayer storage projects.
type Container_Utility_File_Entity struct {
	Entity

	// A file entity's raw content.
	Content *[]byte `json:"content,omitempty" xmlrpc:"content,omitempty"`

	// A file entity's MIME content type.
	ContentType *string `json:"contentType,omitempty" xmlrpc:"contentType,omitempty"`

	// The date a file entity was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The date a CloudLayer storage file entity was moved into the recycle bin. This field applies to files that are pending deletion in the recycle bin.
	DeleteDate *Time `json:"deleteDate,omitempty" xmlrpc:"deleteDate,omitempty"`

	// Unique identifier for the file. This can be either a number or guid.
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Whether a CloudLayer storage file entity is shared with another CloudLayer user.
	IsShared *int `json:"isShared,omitempty" xmlrpc:"isShared,omitempty"`

	// The date a file entity was last changed.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A file entity's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The owner is usually the account who first upload or created the file on the resource or the account who is responsible for the file at the moment.
	Owner *string `json:"owner,omitempty" xmlrpc:"owner,omitempty"`

	// The size of a file entity in bytes.
	Size *uint `json:"size,omitempty" xmlrpc:"size,omitempty"`

	// A CloudLayer storage file entity's type. Types can include "file", "folder", "dir", and "project".
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The latest revision of a file on a CloudLayer storage volume. This number increments each time a new revision of the file is uploaded.
	Version *int `json:"version,omitempty" xmlrpc:"version,omitempty"`
}

// no documentation yet
type Container_Utility_Message struct {
	Entity

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	Summary *string `json:"summary,omitempty" xmlrpc:"summary,omitempty"`
}

// SoftLayer customer servers that are purchased with the Microsoft Windows operating system are configured by default to retrieve updates from SoftLayer's local Windows Server Update Services (WSUS) server. Periodically, these servers synchronize and check for new updates from their local WSUS server. SoftLayer_Container_Utility_Microsoft_Windows_UpdateServices_Status models the results of a server's last synchronization attempt as queried from SoftLayer's WSUS servers.
type Container_Utility_Microsoft_Windows_UpdateServices_Status struct {
	Entity

	// The last time a server rebooted due to a Windows Update.
	LastRebootDate *Time `json:"lastRebootDate,omitempty" xmlrpc:"lastRebootDate,omitempty"`

	// The last time that SoftLayer's local WSUS server received a status update from a customer server.
	LastStatusDate *Time `json:"lastStatusDate,omitempty" xmlrpc:"lastStatusDate,omitempty"`

	// The last time a server synchronized with SoftLayer's local WSUS server.
	LastSyncDate *Time `json:"lastSyncDate,omitempty" xmlrpc:"lastSyncDate,omitempty"`

	// This is the private IP address for this server.
	PrivateIPAddress *string `json:"privateIPAddress,omitempty" xmlrpc:"privateIPAddress,omitempty"`

	// The status message returned from a server's last synchronization with SoftLayer's local WSUS server.
	SyncStatus *string `json:"syncStatus,omitempty" xmlrpc:"syncStatus,omitempty"`

	// A server's update status, as retrieved form SoftLayer's local WSUS server.
	UpdateStatus *string `json:"updateStatus,omitempty" xmlrpc:"updateStatus,omitempty"`
}

// SoftLayer_Container_Utility_Microsoft_Windows_UpdateServices_UpdateItem models a single Microsoft Update as reported by SoftLayer's private Windows Server Update Services (WSUS) services. All servers purchased with Microsoft Windows retrieve updates from SoftLayer's WSUS servers by default.
type Container_Utility_Microsoft_Windows_UpdateServices_UpdateItem struct {
	Entity

	// A short description of a Microsoft Windows Update.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Flag indicating that this patch failed to properly install
	Failed *bool `json:"failed,omitempty" xmlrpc:"failed,omitempty"`

	// A Windows Update's knowledge base article number. Every Windows Update can be referenced on the Microsoft Help and Support site at the URL <nowiki>http://support.microsoft.com/kb/<article number></nowiki>.
	KbArticleNumber *int `json:"kbArticleNumber,omitempty" xmlrpc:"kbArticleNumber,omitempty"`

	// Flag indicating that the update is entirely optionals
	Optional *bool `json:"optional,omitempty" xmlrpc:"optional,omitempty"`

	// Flag indicating that a reboot is needed for this update to be fully applied
	RequiresReboot *bool `json:"requiresReboot,omitempty" xmlrpc:"requiresReboot,omitempty"`
}

// The SoftLayer_Container_Utility_Network_Firewall_Rule_Attribute data type contains information relating to a single firewall rule.
type Container_Utility_Network_Firewall_Rule_Attribute struct {
	Entity

	// The valid actions for use with rules.
	Actions []string `json:"actions,omitempty" xmlrpc:"actions,omitempty"`

	// Maximum allowed number of rules.
	MaximumRuleCount *int `json:"maximumRuleCount,omitempty" xmlrpc:"maximumRuleCount,omitempty"`

	// The valid protocols for use with rules.
	Protocols []string `json:"protocols,omitempty" xmlrpc:"protocols,omitempty"`

	// The valid source ip subnet masks for use with rules.
	SourceIpSubnetMasks []Container_Utility_Network_Subnet_Mask_Generic_Detail `json:"sourceIpSubnetMasks,omitempty" xmlrpc:"sourceIpSubnetMasks,omitempty"`
}

// The SoftLayer_Container_Utility_Network_Subnet_Mask_Generic_Detail data type contains information relating to a subnet mask and details associated with that object.
type Container_Utility_Network_Subnet_Mask_Generic_Detail struct {
	Entity

	// The subnet cidr prefix.
	Cidr *string `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// The subnet mask description.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The subnet mask.
	Mask *string `json:"mask,omitempty" xmlrpc:"mask,omitempty"`
}

// The SoftLayer_Container_Virtual_ConsoleData data type contains information used to access a VSIs console
type Container_Virtual_ConsoleData struct {
	Entity

	// The websocket host address used to access the console
	WebsocketHost *string `json:"websocketHost,omitempty" xmlrpc:"websocketHost,omitempty"`

	// The path to the websocket
	WebsocketPath *string `json:"websocketPath,omitempty" xmlrpc:"websocketPath,omitempty"`

	// The websocket port used to access the console
	WebsocketPort *string `json:"websocketPort,omitempty" xmlrpc:"websocketPort,omitempty"`

	// The token used to authenticate with the console websocket
	WebsocketToken *string `json:"websocketToken,omitempty" xmlrpc:"websocketToken,omitempty"`
}

// This data type represents the structure to hold the allocation properties of a [[SoftLayer_Virtual_DedicatedHost]].
type Container_Virtual_DedicatedHost_AllocationStatus struct {
	Entity

	// Number of allocated CPU cores on the specified dedicated host.
	CpuAllocated *int `json:"cpuAllocated,omitempty" xmlrpc:"cpuAllocated,omitempty"`

	// Number of available CPU cores on the specified dedicated host.
	CpuAvailable *int `json:"cpuAvailable,omitempty" xmlrpc:"cpuAvailable,omitempty"`

	// Total number of CPU cores on the dedicated host.
	CpuCount *int `json:"cpuCount,omitempty" xmlrpc:"cpuCount,omitempty"`

	// Amount of allocated disk space on the specified dedicated host.
	DiskAllocated *int `json:"diskAllocated,omitempty" xmlrpc:"diskAllocated,omitempty"`

	// Amount of available disk space on the specified dedicated host.
	DiskAvailable *int `json:"diskAvailable,omitempty" xmlrpc:"diskAvailable,omitempty"`

	// Total amount of disk capacity on the dedicated host.
	DiskCapacity *int `json:"diskCapacity,omitempty" xmlrpc:"diskCapacity,omitempty"`

	// Number of allocated guests on the specified dedicated host.
	GuestCount *int `json:"guestCount,omitempty" xmlrpc:"guestCount,omitempty"`

	// Amount of allocated memory on the specified dedicated host.
	MemoryAllocated *int `json:"memoryAllocated,omitempty" xmlrpc:"memoryAllocated,omitempty"`

	// Amount of available memory on the specified dedicated host.
	MemoryAvailable *int `json:"memoryAvailable,omitempty" xmlrpc:"memoryAvailable,omitempty"`

	// Total amount of memory capacity on the dedicated host.
	MemoryCapacity *int `json:"memoryCapacity,omitempty" xmlrpc:"memoryCapacity,omitempty"`
}

// This data type represents PCI device allocation properties of a [[SoftLayer_Virtual_DedicatedHost]].
type Container_Virtual_DedicatedHost_Pci_Device_AllocationStatus struct {
	Entity

	// The number of PCI devices on the host.
	DeviceCount *int `json:"deviceCount,omitempty" xmlrpc:"deviceCount,omitempty"`

	// The name of the PCI devices on the host.
	DeviceName *string `json:"deviceName,omitempty" xmlrpc:"deviceName,omitempty"`

	// The number of PCI devices currently allocated to guests.
	DevicesAllocated *int `json:"devicesAllocated,omitempty" xmlrpc:"devicesAllocated,omitempty"`

	// The number of PCI devices available for allocation.
	DevicesAvailable *int `json:"devicesAvailable,omitempty" xmlrpc:"devicesAvailable,omitempty"`

	// The generic component model ID of the PCI device.
	HardwareComponentModelGenericId *int `json:"hardwareComponentModelGenericId,omitempty" xmlrpc:"hardwareComponentModelGenericId,omitempty"`

	// The ID of the host that the dedicated host is on.
	HostId *int `json:"hostId,omitempty" xmlrpc:"hostId,omitempty"`
}

// The SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration data type contains information relating to a template's external location for importing and exporting
type Container_Virtual_Guest_Block_Device_Template_Configuration struct {
	Entity

	//
	// Optional virtualization boot mode parameter, if set, can mark a template to boot specifically into PV or HVM.
	BootMode *string `json:"bootMode,omitempty" xmlrpc:"bootMode,omitempty"`

	//
	// Specifies if image is using a customer's software license.
	Byol *bool `json:"byol,omitempty" xmlrpc:"byol,omitempty"`

	//
	// Specifies if image requires cloud-init.
	CloudInit *bool `json:"cloudInit,omitempty" xmlrpc:"cloudInit,omitempty"`

	//
	// CRN to customer root key
	CrkCrn *string `json:"crkCrn,omitempty" xmlrpc:"crkCrn,omitempty"`

	//
	// For future use; not currently defined.
	EnvironmentType []string `json:"environmentType,omitempty" xmlrpc:"environmentType,omitempty"`

	//
	// IBM Cloud HMAC Access Key
	IbmAccessKey *string `json:"ibmAccessKey,omitempty" xmlrpc:"ibmAccessKey,omitempty"`

	//
	// IBM Cloud (Bluemix) API Key
	IbmApiKey *string `json:"ibmApiKey,omitempty" xmlrpc:"ibmApiKey,omitempty"`

	//
	// IBM HMAC Secret Key
	IbmSecretKey *string `json:"ibmSecretKey,omitempty" xmlrpc:"ibmSecretKey,omitempty"`

	//
	// Specifies if image is encrypted or not.
	IsEncrypted *bool `json:"isEncrypted,omitempty" xmlrpc:"isEncrypted,omitempty"`

	// The group name to be applied to the imported template
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The note to be applied to the imported template
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	//
	// The referenceCode of the operating system software description for the imported VHD
	OperatingSystemReferenceCode *string `json:"operatingSystemReferenceCode,omitempty" xmlrpc:"operatingSystemReferenceCode,omitempty"`

	//
	// Name of the IBM Key Protect Key Name. Required if using an encrypted image.
	RootKeyId *string `json:"rootKeyId,omitempty" xmlrpc:"rootKeyId,omitempty"`

	//
	// Optional Collection of modes that this template supports booting into.
	SupportedBootModes []string `json:"supportedBootModes,omitempty" xmlrpc:"supportedBootModes,omitempty"`

	//
	// The URI for an object storage object (.vhd/.iso file)
	// <code>swift://<ObjectStorageAccountName>@<clusterName>/<containerName>/<fileName.(vhd|iso)></code>
	Uri *string `json:"uri,omitempty" xmlrpc:"uri,omitempty"`

	//
	// Wrapped Decryption Key provided by IBM Key Protect
	WrappedDek *string `json:"wrappedDek,omitempty" xmlrpc:"wrappedDek,omitempty"`
}

// no documentation yet
type Container_Virtual_Guest_Block_Device_Template_Group_RiasAccount struct {
	Entity

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	MasterUserId *int `json:"masterUserId,omitempty" xmlrpc:"masterUserId,omitempty"`

	// no documentation yet
	Token *string `json:"token,omitempty" xmlrpc:"token,omitempty"`
}

// The guest configuration container is used to provide configuration options for creating computing instances.
//
// Each configuration option will include both an <code>itemPrice</code> and a <code>template</code>.
//
// The <code>itemPrice</code> value will provide hourly and monthly costs (if either are applicable), and a description of the option.
//
// The <code>template</code> will provide a fragment of the request with the properties and values that must be sent when creating a computing instance with the option.
//
// The [[SoftLayer_Virtual_Guest/getCreateObjectOptions|getCreateObjectOptions]] method returns this data structure.
//
// <style type="text/css">#properties .views-field-body p { margin-top: 1.5em; };</style>
type Container_Virtual_Guest_Configuration struct {
	Entity

	//
	// <div style="width: 200%">
	// Available block device options.
	//
	//
	// A computing instance will have at least one block device represented by a <code>device</code> number of <code>'0'</code>.
	//
	//
	// The <code>blockDevices.device</code> value in the template represents which device the option is for.
	// The <code>blockDevices.diskImage.capacity</code> value in the template represents the size, in gigabytes, of the disk.
	// The <code>localDiskFlag</code> value in the template represents whether the option is a local or SAN based disk.
	//
	//
	// Note: The block device number <code>'1'</code> is reserved for the SWAP disk attached to the computing instance.
	// </div>
	BlockDevices []Container_Virtual_Guest_Configuration_Option `json:"blockDevices,omitempty" xmlrpc:"blockDevices,omitempty"`

	//
	// <div style="width: 200%">
	// Available datacenter options.
	//
	//
	// The <code>datacenter.name</code> value in the template represents which datacenter the computing instance will be provisioned in.
	// </div>
	Datacenters []Container_Virtual_Guest_Configuration_Option `json:"datacenters,omitempty" xmlrpc:"datacenters,omitempty"`

	//
	// <div style="width: 200%">
	//
	//
	// Available flavor options.
	//
	//
	// The <code>supplementalCreateObjectOptions.flavorKeyName</code> value in the template is an identifier for a particular core, ram, and primary disk configuration.
	//
	//
	// When providing a <code>supplementalCreateObjectOptions.flavorKeyName</code> option the core, ram, and primary disk options are not needed. If those options are provided they are validated against the flavor.
	// </div>
	Flavors []Container_Virtual_Guest_Configuration_Option `json:"flavors,omitempty" xmlrpc:"flavors,omitempty"`

	//
	// <div style="width: 200%">
	// Available memory options.
	//
	//
	// The <code>maxMemory</code> value in the template represents the amount of memory, in megabytes, allocated to the computing instance.
	// </div>
	Memory []Container_Virtual_Guest_Configuration_Option `json:"memory,omitempty" xmlrpc:"memory,omitempty"`

	//
	// <div style="width: 200%">
	// Available network component options.
	//
	//
	// The <code>networkComponent.maxSpeed</code> value in the template represents the link speed, in megabits per second, of the network connections for a computing instance.
	// </div>
	NetworkComponents []Container_Virtual_Guest_Configuration_Option `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`

	//
	// <div style="width: 200%">
	// Available operating system options.
	//
	//
	// The <code>operatingSystemReferenceCode</code> value in the template is an identifier for a particular operating system. When provided exactly as shown in the template, that operating system will be used.
	//
	//
	// A reference code is structured as three tokens separated by underscores. The first token represents the product, the second is the version of the product, and the third is whether the OS is 32 or 64bit.
	//
	//
	// When providing an <code>operatingSystemReferenceCode</code> while ordering a computing instance the only token required to match exactly is the product. The version token may be given as 'LATEST', else it will require an exact match as well. When the bits token is not provided, 64 bits will be assumed.
	//
	//
	// Providing the value of 'LATEST' for a version will select the latest release of that product for the operating system. As this may change over time, you should be sure that the release version is irrelevant for your applications.
	//
	//
	// For Windows based operating systems the version will represent both the release version (2008, 2012, etc) and the edition (Standard, Enterprise, etc). For all other operating systems the version will represent the major version (Centos 6, Ubuntu 12, etc) of that operating system, minor versions are not represented in a reference code.
	//
	//
	// <b>Notice</b> - Some operating systems are charged based on the value specified in <code>startCpus</code>. The price which is used can be determined by calling [[SoftLayer_Virtual_Guest/generateOrderTemplate|generateOrderTemplate]] with your desired device specifications.
	// </div>
	OperatingSystems []Container_Virtual_Guest_Configuration_Option `json:"operatingSystems,omitempty" xmlrpc:"operatingSystems,omitempty"`

	//
	// <div style="width: 200%">
	// Available processor options.
	//
	//
	// The <code>startCpus</code> value in the template represents the number of cores allocated to the computing instance.
	// The <code>dedicatedAccountHostOnlyFlag</code> value in the template represents whether the instance will run on hosts with instances belonging to other accounts.
	// </div>
	Processors []Container_Virtual_Guest_Configuration_Option `json:"processors,omitempty" xmlrpc:"processors,omitempty"`
}

// An option found within a [[SoftLayer_Container_Virtual_Guest_Configuration (type)]] structure.
type Container_Virtual_Guest_Configuration_Option struct {
	Entity

	//
	// Provides a description of a pre-defined configuration with monthly and hourly costs.
	Flavor *Product_Package_Preset `json:"flavor,omitempty" xmlrpc:"flavor,omitempty"`

	//
	// Provides hourly and monthly costs (if either are applicable), and a description of the option.
	ItemPrice *Product_Item_Price `json:"itemPrice,omitempty" xmlrpc:"itemPrice,omitempty"`

	//
	// Provides a fragment of the request with the properties and values that must be sent when creating a computing instance with the option.
	Template *Virtual_Guest `json:"template,omitempty" xmlrpc:"template,omitempty"`
}

// The SoftLayer_Container_Virtual_Guest_PendingMaintenanceAction data type contains information relating to a SoftLayer_Virtual_Guest's pending maintenance actions.
type Container_Virtual_Guest_PendingMaintenanceAction struct {
	Entity

	// The ID of the associated action.
	ActionId *int `json:"actionId,omitempty" xmlrpc:"actionId,omitempty"`

	// The datetime at which this action will be initiated regardless of customer action (if it has not already been completed).
	DueDate *Time `json:"dueDate,omitempty" xmlrpc:"dueDate,omitempty"`

	// User-friendly status.
	//
	// The <code>Completed</code> status means that it is done, no further action is required. The <code>Scheduled</code> status means that the action is pending and will start on the <code>dueDate</code> if no customer action is taken before such time. The <code>In Progress</code> status means the action is currently being executed.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The ticket associated with this maintenance action.
	Ticket *Ticket `json:"ticket,omitempty" xmlrpc:"ticket,omitempty"`

	// The Title for the associated action.
	Title *string `json:"title,omitempty" xmlrpc:"title,omitempty"`

	// The Trigger Explanation for the associated action.
	TriggerExplanation *string `json:"triggerExplanation,omitempty" xmlrpc:"triggerExplanation,omitempty"`
}

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

// The SoftLayer_User_Customer data type contains general information relating to a single SoftLayer customer portal user. Personal information in this type such as names, addresses, and phone numbers are not necessarily associated with the customer account the user is assigned to.
type User_Customer struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerService returns an instance of the User_Customer SoftLayer service
func GetUserCustomerService(sess session.SLSession) User_Customer {
	return User_Customer{Session: sess}
}

func (r User_Customer) Id(id int) User_Customer {
	r.Options.Id = &id
	return r
}

func (r User_Customer) Mask(mask string) User_Customer {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer) Filter(filter string) User_Customer {
	r.Options.Filter = filter
	return r
}

func (r User_Customer) Limit(limit int) User_Customer {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer) Offset(offset int) User_Customer {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer) AcknowledgeSupportPolicy() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer", "acknowledgeSupportPolicy", nil, &r.Options, &resp)
	return
}

// Create a user's API authentication key, allowing that user access to query the SoftLayer API. addApiAuthenticationKey() returns the user's new API key. Each portal user is allowed only one API key.
func (r User_Customer) AddApiAuthenticationKey() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addApiAuthenticationKey", nil, &r.Options, &resp)
	return
}

// Grants the user access to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer) AddBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Add multiple hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkHardwareAccess() does not attempt to add hardware access if the given user already has access to that hardware object.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer) AddBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Add multiple permissions to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. addBulkPortalPermission() does not attempt to add permissions already assigned to the user.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer) AddBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permissions,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) AddBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addBulkRoles", params, &r.Options, &resp)
	return
}

// Add multiple CloudLayer Computing Instances to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkVirtualGuestAccess() does not attempt to add CloudLayer Computing Instance access if the given user already has access to that CloudLayer Computing Instance object.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
func (r User_Customer) AddBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Grants the user access to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer) AddDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) AddExternalBinding(externalBinding *datatypes.User_External_Binding) (resp datatypes.User_Customer_External_Binding, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addExternalBinding", params, &r.Options, &resp)
	return
}

// Add hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the hardware you're attempting to add then addHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer) AddHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addHardwareAccess", params, &r.Options, &resp)
	return
}

// Create a notification subscription record for the user. If a subscription record exists for the notification, the record will be set to active, if currently inactive.
func (r User_Customer) AddNotificationSubscriber(notificationKeyName *string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Add a permission to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. If the user already has the permission you're attempting to add then addPortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are added based on the keyName property of the permission parameter.
func (r User_Customer) AddPortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permission,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) AddRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addRole", params, &r.Options, &resp)
	return
}

// Add a CloudLayer Computing Instance to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the CloudLayer Computing Instance you're attempting to add then addVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer) AddVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "addVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method can be used in place of [[SoftLayer_User_Customer::editObject]] to change the parent user of this user.
//
// The new parent must be a user on the same account, and must not be a child of this user.  A user is not allowed to change their own parent.
//
// If the cascadeFlag is set to false, then an exception will be thrown if the new parent does not have all of the permissions that this user possesses.  If the cascadeFlag is set to true, then permissions will be removed from this user and the descendants of this user as necessary so that no children of the parent will have permissions that the parent does not possess. However, setting the cascadeFlag to true will not remove the access all device permissions from this user. The customer portal will need to be used to remove these permissions.
func (r User_Customer) AssignNewParentId(parentId *int, cascadePermissionsFlag *bool) (resp datatypes.User_Customer, err error) {
	params := []interface{}{
		parentId,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "assignNewParentId", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to modify using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer) ChangePreference(preferenceTypeKeyName *string, value *string) (resp []datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
		value,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "changePreference", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer) CreateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "createNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Create a new user in the SoftLayer customer portal. It is not possible to set up SLL enable flags during object creation. These flags are ignored during object creation. You will need to make a subsequent call to edit object in order to enable VPN access.
//
// An account's master user and sub-users who have the User Manage permission can add new users.
//
// Users are created with a default permission set. After adding a user it may be helpful to set their permissions and device access.
//
// secondaryPasswordTimeoutDays will be set to the system configured default value if the attribute is not provided or the attribute is not a valid value.
//
// Note, neither password nor vpnPassword parameters are required.
//
// Password When a new user is created, an email will be sent to the new user's email address with a link to a url that will allow the new user to create or change their password for the SoftLayer customer portal.
//
// If the password parameter is provided and is not null, then that value will be validated. If it is a valid password, then the user will be created with this password.  This user will still receive a portal password email.  It can be used within 24 hours to change their password, or it can be allowed to expire, and the password provided during user creation will remain as the user's password.
//
// If the password parameter is not provided or the value is null, the user must set their portal password using the link sent in email within 24 hours.  If the user fails to set their password within 24 hours, then a non-master user can use the "Reset Password" link on the login page of the portal to request a new email.  A master user can use the link to retrieve a phone number to call to assist in resetting their password.
//
// The password parameter is ignored for VPN_ONLY users or for IBMid authenticated users.
//
// vpnPassword If the vpnPassword is provided, then the user's vpnPassword will be set to the provided password.  When creating a vpn only user, the vpnPassword MUST be supplied.  If the vpnPassword is not provided, then the user will need to use the portal to edit their profile and set the vpnPassword.
//
// IBMid considerations When a SoftLayer account is linked to a Platform Services (PaaS, formerly Bluemix) account, AND the trait on the SoftLayer Account indicating IBMid authentication is set, then SoftLayer will delegate the creation of an ACTIVE user to PaaS. This means that even though the request to create a new user in such an account may start at the IMS API, via this delegation we effectively turn it into a request that is driven by PaaS. In particular this means that any "invitation email" that comes to the user, will come from PaaS, not from IMS via IBMid.
//
// Users created in states other than ACTIVE (for example, a VPN_ONLY user) will be created directly in IMS without delegation (but note that no invitation is sent for a user created in any state other than ACTIVE).
func (r User_Customer) CreateObject(templateObject *datatypes.User_Customer, password *string, vpnPassword *string) (resp datatypes.User_Customer, err error) {
	params := []interface{}{
		templateObject,
		password,
		vpnPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "createObject", params, &r.Options, &resp)
	return
}

// Create delivery methods for a notification that the user is subscribed to. Multiple delivery method keyNames can be supplied to create multiple delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer) CreateSubscriberDeliveryMethods(notificationKeyName *string, deliveryMethodKeyNames []string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "createSubscriberDeliveryMethods", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer) DeactivateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "deactivateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObject() if you wish to edit a single user account. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer) EditObject(templateObject *datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "editObject", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObjects() if you wish to edit multiple users at once. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer) EditObjects(templateObjects []datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "editObjects", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) FindUserPreference(profileName *string, containerKeyname *string, preferenceKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
		preferenceKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "findUserPreference", params, &r.Options, &resp)
	return
}

// Retrieve The customer account that a user belongs to.
func (r User_Customer) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer) GetActions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getActions", nil, &r.Options, &resp)
	return
}

// The getActiveExternalAuthenticationVendors method will return a list of available external vendors that a SoftLayer user can authenticate against.  The list will only contain vendors for which the user has at least one active external binding.
func (r User_Customer) GetActiveExternalAuthenticationVendors() (resp []datatypes.Container_User_Customer_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getActiveExternalAuthenticationVendors", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's additional email addresses. These email addresses are contacted when updates are made to support tickets.
func (r User_Customer) GetAdditionalEmails() (resp []datatypes.User_Customer_AdditionalEmail, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAdditionalEmails", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetAgentImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAgentImpersonationToken", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetAllowedDedicatedHostIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAllowedDedicatedHostIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetAllowedHardwareIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAllowedHardwareIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetAllowedVirtualGuestIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAllowedVirtualGuestIds", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's API Authentication keys. There is a max limit of one API key per user.
func (r User_Customer) GetApiAuthenticationKeys() (resp []datatypes.User_Customer_ApiAuthentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getApiAuthenticationKeys", nil, &r.Options, &resp)
	return
}

// This method generate user authentication token and return [[SoftLayer_Container_User_Authentication_Token]] object which will be used to authenticate user to login to SoftLayer customer portal.
func (r User_Customer) GetAuthenticationToken(token *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Authentication_Token, err error) {
	params := []interface{}{
		token,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getAuthenticationToken", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's child users. Some portal users may not have child users.
func (r User_Customer) GetChildUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getChildUsers", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated closed tickets.
func (r User_Customer) GetClosedTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getClosedTickets", nil, &r.Options, &resp)
	return
}

// Retrieve The dedicated hosts to which the user has been granted access.
func (r User_Customer) GetDedicatedHosts() (resp []datatypes.Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getDedicatedHosts", nil, &r.Options, &resp)
	return
}

// This method is not applicable to legacy SoftLayer-authenticated users and can only be invoked for IBMid-authenticated users.
func (r User_Customer) GetDefaultAccount(providerType *string) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getDefaultAccount", params, &r.Options, &resp)
	return
}

// Retrieve The external authentication bindings that link an external identifier to a SoftLayer user.
func (r User_Customer) GetExternalBindings() (resp []datatypes.User_External_Binding, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getExternalBindings", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible hardware. These permissions control which hardware a user has access to in the SoftLayer customer portal.
func (r User_Customer) GetHardware() (resp []datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHardware", nil, &r.Options, &resp)
	return
}

// Retrieve the number of servers that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on hardware. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer) GetHardwareCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHardwareCount", nil, &r.Options, &resp)
	return
}

// Retrieve Hardware notifications associated with this user. A hardware notification links a user to a piece of hardware, and that user will be notified if any monitors on that hardware fail, if the monitors have a status of 'Notify User'.
func (r User_Customer) GetHardwareNotifications() (resp []datatypes.User_Customer_Notification_Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHardwareNotifications", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user has acknowledged the support policy.
func (r User_Customer) GetHasAcknowledgedSupportPolicyFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHasAcknowledgedSupportPolicyFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Permission granting the user access to all Dedicated Host devices on the account.
func (r User_Customer) GetHasFullDedicatedHostAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHasFullDedicatedHostAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all hardware on their account.
func (r User_Customer) GetHasFullHardwareAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHasFullHardwareAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all virtual guests on their account.
func (r User_Customer) GetHasFullVirtualGuestAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getHasFullVirtualGuestAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Specifically relating the Customer instance to an IBMid. A Customer instance may or may not have an IBMid link.
func (r User_Customer) GetIbmIdLink() (resp datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getIbmIdLink", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getImpersonationToken", nil, &r.Options, &resp)
	return
}

// Retrieve Contains the definition of the layout profile.
func (r User_Customer) GetLayoutProfiles() (resp []datatypes.Layout_Profile, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getLayoutProfiles", nil, &r.Options, &resp)
	return
}

// Retrieve A user's locale. Locale holds user's language and region information.
func (r User_Customer) GetLocale() (resp datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getLocale", nil, &r.Options, &resp)
	return
}

// Retrieve A user's attempts to log into the SoftLayer customer portal.
func (r User_Customer) GetLoginAttempts() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getLoginAttempts", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a user to the SoftLayer customer portal using the provided authentication container. Depending on the specific type of authentication container that is used, this API will leverage the appropriate authentication protocol. If authentication is successful then the API returns a list of linked accounts for the user, a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer) GetLoginToken(request *datatypes.Container_Authentication_Request_Contract) (resp datatypes.Container_Authentication_Response_Common, err error) {
	params := []interface{}{
		request,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getLoginToken", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IBMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one SoftLayer user per account. This effectively links the OpenIdConnect identity to those accounts. This API returns a list of all the accounts for which there is a link between the OpenIdConnect identity and a SoftLayer user. Invoke this only on IBMid-authenticated users.
func (r User_Customer) GetMappedAccounts(providerType *string) (resp []datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getMappedAccounts", params, &r.Options, &resp)
	return
}

// Retrieve Notification subscription records for the user.
func (r User_Customer) GetNotificationSubscribers() (resp []datatypes.Notification_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getNotificationSubscribers", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer service. You can only retrieve users that are assigned to the customer account belonging to the user making the API call.
func (r User_Customer) GetObject() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getObject", nil, &r.Options, &resp)
	return
}

// This API returns a SoftLayer_Container_User_Customer_OpenIdConnect_MigrationState object containing the necessary information to determine what migration state the user is in. If the account is not OpenIdConnect authenticated, then an exception is thrown.
func (r User_Customer) GetOpenIdConnectMigrationState() (resp datatypes.Container_User_Customer_OpenIdConnect_MigrationState, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getOpenIdConnectMigrationState", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated open tickets.
func (r User_Customer) GetOpenTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getOpenTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's vpn accessible subnets.
func (r User_Customer) GetOverrides() (resp []datatypes.Network_Service_Vpn_Overrides, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getOverrides", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's parent user. If a SoftLayer_User_Customer has a null parentId property then it doesn't have a parent user.
func (r User_Customer) GetParent() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getParent", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetPasswordRequirements(isVpn *bool) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		isVpn,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPasswordRequirements", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's permissions. These permissions control that user's access to functions within the SoftLayer customer portal and API.
func (r User_Customer) GetPermissions() (resp []datatypes.User_Customer_CustomerPermission_Permission, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPermissions", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a username and password to the SoftLayer customer portal. Many portal user accounts are configured to require answering a security question on login. In this case getPortalLoginToken() also verifies the given security question ID and answer. If authentication is successful then the API returns a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer) GetPortalLoginToken(username *string, password *string, securityQuestionId *int, securityQuestionAnswer *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		username,
		password,
		securityQuestionId,
		securityQuestionAnswer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPortalLoginToken", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to get using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer) GetPreference(preferenceTypeKeyName *string) (resp datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPreference", params, &r.Options, &resp)
	return
}

// Use any of the preference types to fetch or modify user preferences using [[SoftLayer_User_Customer::getPreference|getPreference]] or [[SoftLayer_User_Customer::changePreference|changePreference]], respectively.
func (r User_Customer) GetPreferenceTypes() (resp []datatypes.User_Preference_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPreferenceTypes", nil, &r.Options, &resp)
	return
}

// Retrieve Data type contains a single user preference to a specific preference type.
func (r User_Customer) GetPreferences() (resp []datatypes.User_Preference, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getPreferences", nil, &r.Options, &resp)
	return
}

// Retrieve the authentication requirements for an outstanding password set/reset request.  The requirements returned in the same SoftLayer_Container_User_Customer_PasswordSet container which is provided as a parameter into this request.  The SoftLayer_Container_User_Customer_PasswordSet::authenticationMethods array will contain an entry for each authentication method required for the user.  See SoftLayer_Container_User_Customer_PasswordSet for more details.
//
// If the user has required authentication methods, then authentication information will be supplied to the SoftLayer_User_Customer::processPasswordSetRequest method within this same SoftLayer_Container_User_Customer_PasswordSet container.  All existing information in the container must continue to exist in the container to complete the password set/reset process.
func (r User_Customer) GetRequirementsForPasswordSet(passwordSet *datatypes.Container_User_Customer_PasswordSet) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		passwordSet,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getRequirementsForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer) GetRoles() (resp []datatypes.User_Permission_Role, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getRoles", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's security question answers. Some portal users may not have security answers or may not be configured to require answering a security question on login.
func (r User_Customer) GetSecurityAnswers() (resp []datatypes.User_Customer_Security_Answer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's notification subscription records.
func (r User_Customer) GetSubscribers() (resp []datatypes.Notification_User_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSubscribers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's successful attempts to log into the SoftLayer customer portal.
func (r User_Customer) GetSuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user is required to acknowledge the support policy for portal access.
func (r User_Customer) GetSupportPolicyAcknowledgementRequiredFlag() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSupportPolicyAcknowledgementRequiredFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetSupportPolicyDocument() (resp []byte, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSupportPolicyDocument", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetSupportPolicyName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSupportPolicyName", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetSupportedLocales() (resp []datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSupportedLocales", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user must take a brief survey the next time they log into the SoftLayer customer portal.
func (r User_Customer) GetSurveyRequiredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSurveyRequiredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The surveys that a user has taken in the SoftLayer customer portal.
func (r User_Customer) GetSurveys() (resp []datatypes.Survey, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getSurveys", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated tickets.
func (r User_Customer) GetTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's time zone.
func (r User_Customer) GetTimezone() (resp datatypes.Locale_Timezone, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getTimezone", nil, &r.Options, &resp)
	return
}

// Retrieve A user's unsuccessful attempts to log into the SoftLayer customer portal.
func (r User_Customer) GetUnsuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getUnsuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Retrieve a user id using a password token provided to the user in an email generated by the SoftLayer_User_Customer::initiatePortalPasswordChange request. Password recovery keys are valid for 24 hours after they're generated.
//
// When a new user is created or when a user has requested a password change using initiatePortalPasswordChange, they will have received an email that contains a url with a token.  That token is used as the parameter for getUserIdForPasswordSet.  Once the user id is known, then the SoftLayer_User_Customer object can be retrieved which is necessary to complete the process to set or reset a user's password.
func (r User_Customer) GetUserIdForPasswordSet(key *string) (resp int, err error) {
	params := []interface{}{
		key,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getUserIdForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve User customer link with IBMid and IAMid.
func (r User_Customer) GetUserLinks() (resp []datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getUserLinks", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) GetUserPreferences(profileName *string, containerKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getUserPreferences", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's status, which controls overall access to the SoftLayer customer portal and VPN access to the private network.
func (r User_Customer) GetUserStatus() (resp datatypes.User_Customer_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getUserStatus", nil, &r.Options, &resp)
	return
}

// Retrieve the number of CloudLayer Computing Instances that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on CloudLayer Computing Instances. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer) GetVirtualGuestCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getVirtualGuestCount", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible CloudLayer Computing Instances. These permissions control which CloudLayer Computing Instances a user has access to in the SoftLayer customer portal.
func (r User_Customer) GetVirtualGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "getVirtualGuests", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) InTerminalStatus() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "inTerminalStatus", nil, &r.Options, &resp)
	return
}

// Sends password change email to the user containing url that allows the user the change their password. This is the first step when a user wishes to change their password.  The url that is generated contains a one-time use token that is valid for only 24-hours.
//
// If this is a new master user who has never logged into the portal, then password reset will be initiated. Once a master user has logged into the portal, they must setup their security questions prior to logging out because master users are required to answer a security question during the password reset process.  Should a master user not have security questions defined and not remember their password in order to define the security questions, then they will need to contact support at live chat or Revenue Services for assistance.
//
// Due to security reasons, the number of reset requests per username are limited within a undisclosed timeframe.
func (r User_Customer) InitiatePortalPasswordChange(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "initiatePortalPasswordChange", params, &r.Options, &resp)
	return
}

// A Brand Agent that has permissions to Add Customer Accounts will be able to request the password email be sent to the Master User of a Customer Account created by the same Brand as the agent making the request. Due to security reasons, the number of reset requests are limited within an undisclosed timeframe.
func (r User_Customer) InitiatePortalPasswordChangeByBrandAgent(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "initiatePortalPasswordChangeByBrandAgent", params, &r.Options, &resp)
	return
}

// Send email invitation to a user to join a SoftLayer account and authenticate with OpenIdConnect. Throws an exception on error.
func (r User_Customer) InviteUserToLinkOpenIdConnect(providerType *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "inviteUserToLinkOpenIdConnect", params, &r.Options, &resp)
	return
}

// Portal users are considered master users if they don't have an associated parent user. The only users who don't have parent users are users whose username matches their SoftLayer account name. Master users have special permissions throughout the SoftLayer customer portal.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer) IsMasterUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "isMasterUser", nil, &r.Options, &resp)
	return
}

// Determine if a string is the given user's login password to the SoftLayer customer portal.
func (r User_Customer) IsValidPortalPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "isValidPortalPassword", params, &r.Options, &resp)
	return
}

// The perform external authentication method will authenticate the given external authentication container with an external vendor.  The authentication container and its contents will be verified before an attempt is made to authenticate the contents of the container with an external vendor.
func (r User_Customer) PerformExternalAuthentication(authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "performExternalAuthentication", params, &r.Options, &resp)
	return
}

// Set the password for a user who has an outstanding password request. A user with an outstanding password request will have an unused and unexpired password key.  The password key is part of the url provided to the user in the email sent to the user with information on how to set their password.  The email was generated by the SoftLayer_User_Customer::initiatePortalPasswordRequest request. Password recovery keys are valid for 24 hours after they're generated.
//
// If the user has required authentication methods as specified by in the SoftLayer_Container_User_Customer_PasswordSet container returned from the SoftLayer_User_Customer::getRequirementsForPasswordSet request, then additional requests must be made to processPasswordSetRequest to authenticate the user before changing the password.  First, if the user has security questions set on their profile, they will be required to answer one of their questions correctly. Next, if the user has Verisign or Google Authentication on their account, they must authenticate according to the two-factor provider.  All of this authentication is done using the SoftLayer_Container_User_Customer_PasswordSet container.
//
// User portal passwords must match the following restrictions. Portal passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ + =
// * ...not match your username
func (r User_Customer) ProcessPasswordSetRequest(passwordSet *datatypes.Container_User_Customer_PasswordSet, authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp bool, err error) {
	params := []interface{}{
		passwordSet,
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "processPasswordSetRequest", params, &r.Options, &resp)
	return
}

// Revoke access to all dedicated hosts on the account for this user. The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer) RemoveAllDedicatedHostAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeAllDedicatedHostAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer) RemoveAllHardwareAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeAllHardwareAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all cloud computing instances from a portal user's instance access list. A user's instance access list controls which of an account's computing instance objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer) RemoveAllVirtualAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeAllVirtualAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove a user's API authentication key, removing that user's access to query the SoftLayer API.
func (r User_Customer) RemoveApiAuthenticationKey(keyId *int) (resp bool, err error) {
	params := []interface{}{
		keyId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeApiAuthenticationKey", params, &r.Options, &resp)
	return
}

// Revokes access for the user to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// If the user has full dedicatedHost access, then it will provide access to "ALL but passed in" dedicatedHost ids.
func (r User_Customer) RemoveBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Remove multiple hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting to remove then removeBulkHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// If the user has full hardware access, then it will provide access to "ALL but passed in" hardware ids.
func (r User_Customer) RemoveBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) multiple permissions from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. removePortalPermission() does not attempt to remove permissions that are not assigned to the user.
//
// Users can grant or revoke permissions to their child users, but not to themselves. An account's master has all portal permissions and can grant permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permissions from a user will cascade down the child hierarchy and remove the permissions from this user along with all child users who also have the permission.
//
// If the cascadePermissionsFlag is not provided or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer) RemoveBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permissions,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) RemoveBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeBulkRoles", params, &r.Options, &resp)
	return
}

// Remove multiple CloudLayer Computing Instances from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeBulkVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer) RemoveBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Revokes access for the user to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer) RemoveDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) RemoveExternalBinding(externalBinding *datatypes.User_External_Binding) (resp bool, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeExternalBinding", params, &r.Options, &resp)
	return
}

// Remove hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting remove add then removeHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer) RemoveHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) a permission from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. If the user does not have the permission you're attempting to remove then removePortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permission from a user will cascade down the child hierarchy and remove the permission from this user and all child users who also have the permission.
//
// If the cascadePermissionsFlag is not set or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission parameter.
func (r User_Customer) RemovePortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permission,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removePortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) RemoveRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeRole", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) RemoveSecurityAnswers() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Remove a CloudLayer Computing Instance from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's computing instances a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set instance access for any of the other users on their account.
func (r User_Customer) RemoveVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "removeVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer user is linked to, if we need to do that for some reason. It will do this by modifying the link to the desired new IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is already Bluemix linked. To reset a link for the Bluemix-linked user account, use resetOpenIdConnectLinkUnifiedUserManagementMode.
func (r User_Customer) ResetOpenIdConnectLink(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "resetOpenIdConnectLink", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer master user is linked to, if we need to do that for some reason. It will do this by unlinking the new owner IBMid from its current user association in this account, if there is one (note that the new owner IBMid is not required to already be a member of the IMS account). Then it will modify the existing IBMid link for the master user to use the new owner IBMid-realm IAMid. At this point, if the new owner IBMid isn't already a member of the PaaS account, it will attempt to add it. As a last step, it will call PaaS to modify the owner on that side, if necessary.  Only when all those steps are complete, it will commit the IMS-side DB changes.  Then, it will clean up the SoftLayer user that was linked to the new owner IBMid (this user became unlinked as the first step in this process).  It will also call BSS to delete the old owner IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is not Bluemix linked. To reset a link for the user account not linked to Bluemix, use resetOpenIdConnectLink.
func (r User_Customer) ResetOpenIdConnectLinkUnifiedUserManagementMode(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "resetOpenIdConnectLinkUnifiedUserManagementMode", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) SamlAuthenticate(accountId *string, samlResponse *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		accountId,
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "samlAuthenticate", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) SamlBeginAuthentication(accountId *int) (resp string, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "samlBeginAuthentication", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) SamlBeginLogout() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "samlBeginLogout", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) SamlLogout(samlResponse *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "samlLogout", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer) SelfPasswordChange(currentPassword *string, newPassword *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		currentPassword,
		newPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "selfPasswordChange", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IBMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one per account. If an OpenIdConnect identity is mapped to multiple accounts in this manner, one such account should be identified as the default account for that identity. Invoke this only on IBMid-authenticated users.
func (r User_Customer) SetDefaultAccount(providerType *string, accountId *int) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "setDefaultAccount", params, &r.Options, &resp)
	return
}

// As master user, calling this api for the IBMid provider type when there is an existing IBMid for the email on the SL account will silently (without sending an invitation email) create a link for the IBMid. NOTE: If the SoftLayer user is already linked to IBMid, this call will fail. If the IBMid specified by the email of this user, is already used in a link to another user in this account, this call will fail. If there is already an open invitation from this SoftLayer user to this or any IBMid, this call will fail. If there is already an open invitation from some other SoftLayer user in this account to this IBMid, then this call will fail.
func (r User_Customer) SilentlyMigrateUserOpenIdConnect(providerType *string) (resp bool, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "silentlyMigrateUserOpenIdConnect", params, &r.Options, &resp)
	return
}

// This method allows the master user of an account to undo the designation of this user as an alternate master user.  This can not be applied to the true master user of the account.
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted this user.  This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner turns off an "alternate/auxiliary master user / account owner".
func (r User_Customer) TurnOffMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer", "turnOffMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// This method allows the master user of an account to designate this user as an alternate master user.  Effectively this means that this user should have "all the same IMS permissions as a master user".
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted to this user. This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner designates an "alternate/auxiliary master user / account owner".
func (r User_Customer) TurnOnMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer", "turnOnMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// Update the active status for a notification that the user is subscribed to. A notification along with an active flag can be supplied to update the active status for a particular notification subscription.
func (r User_Customer) UpdateNotificationSubscriber(notificationKeyName *string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "updateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Update a user's login security questions and answers on the SoftLayer customer portal. These questions and answers are used to optionally log into the SoftLayer customer portal using two-factor authentication. Each user must have three distinct questions set with a unique answer for each question, and each answer may only contain alphanumeric or the . , - _ ( ) [ ] : ; > < characters. Existing user security questions and answers are deleted before new ones are set, and users may only update their own security questions and answers.
func (r User_Customer) UpdateSecurityAnswers(questions []datatypes.User_Security_Question, answers []string) (resp bool, err error) {
	params := []interface{}{
		questions,
		answers,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "updateSecurityAnswers", params, &r.Options, &resp)
	return
}

// Update a delivery method for a notification that the user is subscribed to. A delivery method keyName along with an active flag can be supplied to update the active status of the delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer) UpdateSubscriberDeliveryMethod(notificationKeyName *string, deliveryMethodKeyNames []string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "updateSubscriberDeliveryMethod", params, &r.Options, &resp)
	return
}

// Update a user's VPN password on the SoftLayer customer portal. As with portal passwords, VPN passwords must match the following restrictions. VPN passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ =
// * ...not match your username
// Finally, users can only update their own VPN password. An account's master user can update any of their account users' VPN passwords.
func (r User_Customer) UpdateVpnPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "updateVpnPassword", params, &r.Options, &resp)
	return
}

// Always call this function to enable changes when manually configuring VPN subnet access.
func (r User_Customer) UpdateVpnUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer", "updateVpnUser", nil, &r.Options, &resp)
	return
}

// This method validate the given authentication token using the user id by comparing it with the actual user authentication token and return [[SoftLayer_Container_User_Customer_Portal_Token]] object
func (r User_Customer) ValidateAuthenticationToken(authenticationToken *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationToken,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer", "validateAuthenticationToken", params, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_ApiAuthentication type contains user's authentication key(s).
type User_Customer_ApiAuthentication struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerApiAuthenticationService returns an instance of the User_Customer_ApiAuthentication SoftLayer service
func GetUserCustomerApiAuthenticationService(sess session.SLSession) User_Customer_ApiAuthentication {
	return User_Customer_ApiAuthentication{Session: sess}
}

func (r User_Customer_ApiAuthentication) Id(id int) User_Customer_ApiAuthentication {
	r.Options.Id = &id
	return r
}

func (r User_Customer_ApiAuthentication) Mask(mask string) User_Customer_ApiAuthentication {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_ApiAuthentication) Filter(filter string) User_Customer_ApiAuthentication {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_ApiAuthentication) Limit(limit int) User_Customer_ApiAuthentication {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_ApiAuthentication) Offset(offset int) User_Customer_ApiAuthentication {
	r.Options.Offset = &offset
	return r
}

// Edit the properties of customer ApiAuthentication record by passing in a modified instance of a SoftLayer_User_Customer_ApiAuthentication object. Only the ipAddressRestriction property can be modified.
func (r User_Customer_ApiAuthentication) EditObject(templateObject *datatypes.User_Customer_ApiAuthentication) (resp datatypes.User_Customer_ApiAuthentication, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_ApiAuthentication", "editObject", params, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer_ApiAuthentication object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_ApiAuthentication service.
func (r User_Customer_ApiAuthentication) GetObject() (resp datatypes.User_Customer_ApiAuthentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_ApiAuthentication", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The user who owns the api authentication key.
func (r User_Customer_ApiAuthentication) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_ApiAuthentication", "getUser", nil, &r.Options, &resp)
	return
}

// Each SoftLayer portal account is assigned a series of permissions that determine what access the user has to functions within the SoftLayer customer portal. This status is reflected in the SoftLayer_User_Customer_Status data type. Permissions differ from user status in that user status applies globally to the portal while user permissions are applied to specific portal functions.
type User_Customer_CustomerPermission_Permission struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerCustomerPermissionPermissionService returns an instance of the User_Customer_CustomerPermission_Permission SoftLayer service
func GetUserCustomerCustomerPermissionPermissionService(sess session.SLSession) User_Customer_CustomerPermission_Permission {
	return User_Customer_CustomerPermission_Permission{Session: sess}
}

func (r User_Customer_CustomerPermission_Permission) Id(id int) User_Customer_CustomerPermission_Permission {
	r.Options.Id = &id
	return r
}

func (r User_Customer_CustomerPermission_Permission) Mask(mask string) User_Customer_CustomerPermission_Permission {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_CustomerPermission_Permission) Filter(filter string) User_Customer_CustomerPermission_Permission {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_CustomerPermission_Permission) Limit(limit int) User_Customer_CustomerPermission_Permission {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_CustomerPermission_Permission) Offset(offset int) User_Customer_CustomerPermission_Permission {
	r.Options.Offset = &offset
	return r
}

// Retrieve all available permissions.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer_CustomerPermission_Permission) GetAllObjects() (resp []datatypes.User_Customer_CustomerPermission_Permission, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_CustomerPermission_Permission", "getAllObjects", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer_CustomerPermission_Permission object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_CustomerPermission_Permission service.
func (r User_Customer_CustomerPermission_Permission) GetObject() (resp datatypes.User_Customer_CustomerPermission_Permission, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_CustomerPermission_Permission", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_External_Binding data type contains general information for a single external binding.  This includes the 3rd party vendor, type of binding, and a unique identifier and password that is used to authenticate against the 3rd party service.
type User_Customer_External_Binding struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerExternalBindingService returns an instance of the User_Customer_External_Binding SoftLayer service
func GetUserCustomerExternalBindingService(sess session.SLSession) User_Customer_External_Binding {
	return User_Customer_External_Binding{Session: sess}
}

func (r User_Customer_External_Binding) Id(id int) User_Customer_External_Binding {
	r.Options.Id = &id
	return r
}

func (r User_Customer_External_Binding) Mask(mask string) User_Customer_External_Binding {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_External_Binding) Filter(filter string) User_Customer_External_Binding {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_External_Binding) Limit(limit int) User_Customer_External_Binding {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_External_Binding) Offset(offset int) User_Customer_External_Binding {
	r.Options.Offset = &offset
	return r
}

// Delete an external authentication binding.  If the external binding currently has an active billing item associated you will be prevented from deleting the binding.  The alternative method to remove an external authentication binding is to use the service cancellation form.
func (r User_Customer_External_Binding) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "deleteObject", nil, &r.Options, &resp)
	return
}

// Disabling an external binding will allow you to keep the external binding on your SoftLayer account, but will not require you to authentication with our trusted 2 form factor vendor when logging into the SoftLayer customer portal.
//
// You may supply one of the following reason when you disable an external binding:
// *Unspecified
// *TemporarilyUnavailable
// *Lost
// *Stolen
func (r User_Customer_External_Binding) Disable(reason *string) (resp bool, err error) {
	params := []interface{}{
		reason,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "disable", params, &r.Options, &resp)
	return
}

// Enabling an external binding will activate the binding on your account and require you to authenticate with our trusted 3rd party 2 form factor vendor when logging into the SoftLayer customer portal.
//
// Please note that API access will be disabled for users that have an active external binding.
func (r User_Customer_External_Binding) Enable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "enable", nil, &r.Options, &resp)
	return
}

// Retrieve Attributes of an external authentication binding.
func (r User_Customer_External_Binding) GetAttributes() (resp []datatypes.User_External_Binding_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Information regarding the billing item for external authentication.
func (r User_Customer_External_Binding) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve An optional note for identifying the external binding.
func (r User_Customer_External_Binding) GetNote() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getNote", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding) GetObject() (resp datatypes.User_Customer_External_Binding, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The type of external authentication binding.
func (r User_Customer_External_Binding) GetType() (resp datatypes.User_External_Binding_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getType", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer user that the external authentication binding belongs to.
func (r User_Customer_External_Binding) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getUser", nil, &r.Options, &resp)
	return
}

// Retrieve The vendor of an external authentication binding.
func (r User_Customer_External_Binding) GetVendor() (resp datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "getVendor", nil, &r.Options, &resp)
	return
}

// Update the note of an external binding.  The note is an optional property that is used to store information about a binding.
func (r User_Customer_External_Binding) UpdateNote(text *string) (resp bool, err error) {
	params := []interface{}{
		text,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding", "updateNote", params, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_External_Binding_Totp data type contains information about a single time-based one time password external binding.  The external binding information is used when a SoftLayer customer logs into the SoftLayer customer portal to authenticate them.
//
// The information provided by this external binding data type includes:
// * The type of credential
// * The current state of the credential
// ** Active
// ** Inactive
//
// SoftLayer users with an active external binding will be prohibited from using the API for security reasons.
type User_Customer_External_Binding_Totp struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerExternalBindingTotpService returns an instance of the User_Customer_External_Binding_Totp SoftLayer service
func GetUserCustomerExternalBindingTotpService(sess session.SLSession) User_Customer_External_Binding_Totp {
	return User_Customer_External_Binding_Totp{Session: sess}
}

func (r User_Customer_External_Binding_Totp) Id(id int) User_Customer_External_Binding_Totp {
	r.Options.Id = &id
	return r
}

func (r User_Customer_External_Binding_Totp) Mask(mask string) User_Customer_External_Binding_Totp {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_External_Binding_Totp) Filter(filter string) User_Customer_External_Binding_Totp {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_External_Binding_Totp) Limit(limit int) User_Customer_External_Binding_Totp {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_External_Binding_Totp) Offset(offset int) User_Customer_External_Binding_Totp {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer_External_Binding_Totp) Activate() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "activate", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding_Totp) Deactivate() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "deactivate", nil, &r.Options, &resp)
	return
}

// Delete an external authentication binding.  If the external binding currently has an active billing item associated you will be prevented from deleting the binding.  The alternative method to remove an external authentication binding is to use the service cancellation form.
func (r User_Customer_External_Binding_Totp) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "deleteObject", nil, &r.Options, &resp)
	return
}

// Disabling an external binding will allow you to keep the external binding on your SoftLayer account, but will not require you to authentication with our trusted 2 form factor vendor when logging into the SoftLayer customer portal.
//
// You may supply one of the following reason when you disable an external binding:
// *Unspecified
// *TemporarilyUnavailable
// *Lost
// *Stolen
func (r User_Customer_External_Binding_Totp) Disable(reason *string) (resp bool, err error) {
	params := []interface{}{
		reason,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "disable", params, &r.Options, &resp)
	return
}

// Enabling an external binding will activate the binding on your account and require you to authenticate with our trusted 3rd party 2 form factor vendor when logging into the SoftLayer customer portal.
//
// Please note that API access will be disabled for users that have an active external binding.
func (r User_Customer_External_Binding_Totp) Enable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "enable", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding_Totp) GenerateSecretKey() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "generateSecretKey", nil, &r.Options, &resp)
	return
}

// Retrieve Attributes of an external authentication binding.
func (r User_Customer_External_Binding_Totp) GetAttributes() (resp []datatypes.User_External_Binding_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Information regarding the billing item for external authentication.
func (r User_Customer_External_Binding_Totp) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve An optional note for identifying the external binding.
func (r User_Customer_External_Binding_Totp) GetNote() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getNote", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding_Totp) GetObject() (resp datatypes.User_Customer_External_Binding_Totp, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The type of external authentication binding.
func (r User_Customer_External_Binding_Totp) GetType() (resp datatypes.User_External_Binding_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getType", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer user that the external authentication binding belongs to.
func (r User_Customer_External_Binding_Totp) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getUser", nil, &r.Options, &resp)
	return
}

// Retrieve The vendor of an external authentication binding.
func (r User_Customer_External_Binding_Totp) GetVendor() (resp datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "getVendor", nil, &r.Options, &resp)
	return
}

// Update the note of an external binding.  The note is an optional property that is used to store information about a binding.
func (r User_Customer_External_Binding_Totp) UpdateNote(text *string) (resp bool, err error) {
	params := []interface{}{
		text,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Totp", "updateNote", params, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_External_Binding_Vendor data type contains information for a single external binding vendor.  This information includes a user friendly vendor name, a unique version of the vendor name, and a unique internal identifier that can be used when creating a new external binding.
type User_Customer_External_Binding_Vendor struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerExternalBindingVendorService returns an instance of the User_Customer_External_Binding_Vendor SoftLayer service
func GetUserCustomerExternalBindingVendorService(sess session.SLSession) User_Customer_External_Binding_Vendor {
	return User_Customer_External_Binding_Vendor{Session: sess}
}

func (r User_Customer_External_Binding_Vendor) Id(id int) User_Customer_External_Binding_Vendor {
	r.Options.Id = &id
	return r
}

func (r User_Customer_External_Binding_Vendor) Mask(mask string) User_Customer_External_Binding_Vendor {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_External_Binding_Vendor) Filter(filter string) User_Customer_External_Binding_Vendor {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_External_Binding_Vendor) Limit(limit int) User_Customer_External_Binding_Vendor {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_External_Binding_Vendor) Offset(offset int) User_Customer_External_Binding_Vendor {
	r.Options.Offset = &offset
	return r
}

// getAllObjects() will return a list of the available external binding vendors that SoftLayer supports.  Use this list to select the appropriate vendor when creating a new external binding.
func (r User_Customer_External_Binding_Vendor) GetAllObjects() (resp []datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Vendor", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding_Vendor) GetObject() (resp datatypes.User_Customer_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Vendor", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_External_Binding_Verisign data type contains information about a single VeriSign external binding.  The external binding information is used when a SoftLayer customer logs into the SoftLayer customer portal to authenticate them against a 3rd party, in this case VeriSign.
//
// The information provided by the VeriSign external binding data type includes:
// * The type of credential
// * The current state of the credential
// ** Enabled
// ** Disabled
// ** Locked
// * The credential's expiration date
// * The last time the credential was updated
//
// SoftLayer users with an active external binding will be prohibited from using the API for security reasons.
type User_Customer_External_Binding_Verisign struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerExternalBindingVerisignService returns an instance of the User_Customer_External_Binding_Verisign SoftLayer service
func GetUserCustomerExternalBindingVerisignService(sess session.SLSession) User_Customer_External_Binding_Verisign {
	return User_Customer_External_Binding_Verisign{Session: sess}
}

func (r User_Customer_External_Binding_Verisign) Id(id int) User_Customer_External_Binding_Verisign {
	r.Options.Id = &id
	return r
}

func (r User_Customer_External_Binding_Verisign) Mask(mask string) User_Customer_External_Binding_Verisign {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_External_Binding_Verisign) Filter(filter string) User_Customer_External_Binding_Verisign {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_External_Binding_Verisign) Limit(limit int) User_Customer_External_Binding_Verisign {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_External_Binding_Verisign) Offset(offset int) User_Customer_External_Binding_Verisign {
	r.Options.Offset = &offset
	return r
}

// Delete a VeriSign external binding.  The only VeriSign external binding that can be deleted through this method is the free VeriSign external binding for the master user of a SoftLayer account. All other external bindings must be canceled using the SoftLayer service cancellation form.
//
// When a VeriSign external binding is deleted the credential is deactivated in VeriSign's system for use on the SoftLayer site and the $0 billing item associated with the free VeriSign external binding is cancelled.
func (r User_Customer_External_Binding_Verisign) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "deleteObject", nil, &r.Options, &resp)
	return
}

// Disabling an external binding will allow you to keep the external binding on your SoftLayer account, but will not require you to authentication with our trusted 2 form factor vendor when logging into the SoftLayer customer portal.
//
// You may supply one of the following reason when you disable an external binding:
// *Unspecified
// *TemporarilyUnavailable
// *Lost
// *Stolen
func (r User_Customer_External_Binding_Verisign) Disable(reason *string) (resp bool, err error) {
	params := []interface{}{
		reason,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "disable", params, &r.Options, &resp)
	return
}

// Enabling an external binding will activate the binding on your account and require you to authenticate with our trusted 3rd party 2 form factor vendor when logging into the SoftLayer customer portal.
//
// Please note that API access will be disabled for users that have an active external binding.
func (r User_Customer_External_Binding_Verisign) Enable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "enable", nil, &r.Options, &resp)
	return
}

// An activation code is required when provisioning a new mobile credential from Verisign.  This method will return the required activation code.
func (r User_Customer_External_Binding_Verisign) GetActivationCodeForMobileClient() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getActivationCodeForMobileClient", nil, &r.Options, &resp)
	return
}

// Retrieve Attributes of an external authentication binding.
func (r User_Customer_External_Binding_Verisign) GetAttributes() (resp []datatypes.User_External_Binding_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Information regarding the billing item for external authentication.
func (r User_Customer_External_Binding_Verisign) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve The date that a VeriSign credential expires.
func (r User_Customer_External_Binding_Verisign) GetCredentialExpirationDate() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getCredentialExpirationDate", nil, &r.Options, &resp)
	return
}

// Retrieve The last time a VeriSign credential was updated.
func (r User_Customer_External_Binding_Verisign) GetCredentialLastUpdateDate() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getCredentialLastUpdateDate", nil, &r.Options, &resp)
	return
}

// Retrieve The current state of a VeriSign credential. This can be 'Enabled', 'Disabled', or 'Locked'.
func (r User_Customer_External_Binding_Verisign) GetCredentialState() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getCredentialState", nil, &r.Options, &resp)
	return
}

// Retrieve The type of VeriSign credential. This can be either 'Hardware' or 'Software'.
func (r User_Customer_External_Binding_Verisign) GetCredentialType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getCredentialType", nil, &r.Options, &resp)
	return
}

// Retrieve An optional note for identifying the external binding.
func (r User_Customer_External_Binding_Verisign) GetNote() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getNote", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_External_Binding_Verisign) GetObject() (resp datatypes.User_Customer_External_Binding_Verisign, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The type of external authentication binding.
func (r User_Customer_External_Binding_Verisign) GetType() (resp datatypes.User_External_Binding_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getType", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer user that the external authentication binding belongs to.
func (r User_Customer_External_Binding_Verisign) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getUser", nil, &r.Options, &resp)
	return
}

// Retrieve The vendor of an external authentication binding.
func (r User_Customer_External_Binding_Verisign) GetVendor() (resp datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "getVendor", nil, &r.Options, &resp)
	return
}

// If a VeriSign credential becomes locked because of too many failed login attempts the unlock method can be used to unlock a VeriSign credential. As a security precaution a valid security code generated by the credential will be required before the credential is unlocked.
func (r User_Customer_External_Binding_Verisign) Unlock(securityCode *string) (resp bool, err error) {
	params := []interface{}{
		securityCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "unlock", params, &r.Options, &resp)
	return
}

// Update the note of an external binding.  The note is an optional property that is used to store information about a binding.
func (r User_Customer_External_Binding_Verisign) UpdateNote(text *string) (resp bool, err error) {
	params := []interface{}{
		text,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "updateNote", params, &r.Options, &resp)
	return
}

// Validate the user id and VeriSign credential id used to create an external authentication binding.
func (r User_Customer_External_Binding_Verisign) ValidateCredentialId(userId *int, externalId *string) (resp bool, err error) {
	params := []interface{}{
		userId,
		externalId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_External_Binding_Verisign", "validateCredentialId", params, &r.Options, &resp)
	return
}

// no documentation yet
type User_Customer_Invitation struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerInvitationService returns an instance of the User_Customer_Invitation SoftLayer service
func GetUserCustomerInvitationService(sess session.SLSession) User_Customer_Invitation {
	return User_Customer_Invitation{Session: sess}
}

func (r User_Customer_Invitation) Id(id int) User_Customer_Invitation {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Invitation) Mask(mask string) User_Customer_Invitation {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Invitation) Filter(filter string) User_Customer_Invitation {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Invitation) Limit(limit int) User_Customer_Invitation {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Invitation) Offset(offset int) User_Customer_Invitation {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer_Invitation) GetObject() (resp datatypes.User_Customer_Invitation, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Invitation", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer_Invitation) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Invitation", "getUser", nil, &r.Options, &resp)
	return
}

// The Customer_Notification_Hardware object stores links between customers and the hardware devices they wish to monitor.  This link is not enough, the user must be sure to also create SoftLayer_Network_Monitor_Version1_Query_Host instance with the response action set to "notify users" in order for the users linked to that hardware object to be notified on failure.
type User_Customer_Notification_Hardware struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerNotificationHardwareService returns an instance of the User_Customer_Notification_Hardware SoftLayer service
func GetUserCustomerNotificationHardwareService(sess session.SLSession) User_Customer_Notification_Hardware {
	return User_Customer_Notification_Hardware{Session: sess}
}

func (r User_Customer_Notification_Hardware) Id(id int) User_Customer_Notification_Hardware {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Notification_Hardware) Mask(mask string) User_Customer_Notification_Hardware {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Notification_Hardware) Filter(filter string) User_Customer_Notification_Hardware {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Notification_Hardware) Limit(limit int) User_Customer_Notification_Hardware {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Notification_Hardware) Offset(offset int) User_Customer_Notification_Hardware {
	r.Options.Offset = &offset
	return r
}

// Passing in an unsaved instances of a Customer_Notification_Hardware object into this function will create the object and return the results to the user.
func (r User_Customer_Notification_Hardware) CreateObject(templateObject *datatypes.User_Customer_Notification_Hardware) (resp datatypes.User_Customer_Notification_Hardware, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "createObject", params, &r.Options, &resp)
	return
}

// Passing in a collection of unsaved instances of Customer_Notification_Hardware objects into this function will create all objects and return the results to the user.
func (r User_Customer_Notification_Hardware) CreateObjects(templateObjects []datatypes.User_Customer_Notification_Hardware) (resp []datatypes.User_Customer_Notification_Hardware, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "createObjects", params, &r.Options, &resp)
	return
}

// Like any other API object, the customer notification objects can be deleted by passing an instance of them into this function.  The ID on the object must be set.
func (r User_Customer_Notification_Hardware) DeleteObjects(templateObjects []datatypes.User_Customer_Notification_Hardware) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "deleteObjects", params, &r.Options, &resp)
	return
}

// This method returns all Customer_Notification_Hardware objects associated with the passed in hardware ID as long as that hardware ID is owned by the current user's account.
//
// This behavior can also be accomplished by simply tapping monitoringUserNotification on the Hardware_Server object.
func (r User_Customer_Notification_Hardware) FindByHardwareId(hardwareId *int) (resp []datatypes.User_Customer_Notification_Hardware, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "findByHardwareId", params, &r.Options, &resp)
	return
}

// Retrieve The hardware object that will be monitored.
func (r User_Customer_Notification_Hardware) GetHardware() (resp datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "getHardware", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer_Notification_Hardware object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_Notification_Hardware service. You can only retrieve hardware notifications attached to hardware and users that belong to your account
func (r User_Customer_Notification_Hardware) GetObject() (resp datatypes.User_Customer_Notification_Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The user that will be notified when the associated hardware object fails a monitoring instance.
func (r User_Customer_Notification_Hardware) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Hardware", "getUser", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_Notification_Virtual_Guest object stores links between customers and the virtual guests they wish to monitor.  This link is not enough, the user must be sure to also create SoftLayer_Network_Monitor_Version1_Query_Host instance with the response action set to "notify users" in order for the users linked to that Virtual Guest object to be notified on failure.
type User_Customer_Notification_Virtual_Guest struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerNotificationVirtualGuestService returns an instance of the User_Customer_Notification_Virtual_Guest SoftLayer service
func GetUserCustomerNotificationVirtualGuestService(sess session.SLSession) User_Customer_Notification_Virtual_Guest {
	return User_Customer_Notification_Virtual_Guest{Session: sess}
}

func (r User_Customer_Notification_Virtual_Guest) Id(id int) User_Customer_Notification_Virtual_Guest {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Notification_Virtual_Guest) Mask(mask string) User_Customer_Notification_Virtual_Guest {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Notification_Virtual_Guest) Filter(filter string) User_Customer_Notification_Virtual_Guest {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Notification_Virtual_Guest) Limit(limit int) User_Customer_Notification_Virtual_Guest {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Notification_Virtual_Guest) Offset(offset int) User_Customer_Notification_Virtual_Guest {
	r.Options.Offset = &offset
	return r
}

// Passing in an unsaved instance of a SoftLayer_Customer_Notification_Virtual_Guest object into this function will create the object and return the results to the user.
func (r User_Customer_Notification_Virtual_Guest) CreateObject(templateObject *datatypes.User_Customer_Notification_Virtual_Guest) (resp datatypes.User_Customer_Notification_Virtual_Guest, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "createObject", params, &r.Options, &resp)
	return
}

// Passing in a collection of unsaved instances of SoftLayer_Customer_Notification_Virtual_Guest objects into this function will create all objects and return the results to the user.
func (r User_Customer_Notification_Virtual_Guest) CreateObjects(templateObjects []datatypes.User_Customer_Notification_Virtual_Guest) (resp []datatypes.User_Customer_Notification_Virtual_Guest, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "createObjects", params, &r.Options, &resp)
	return
}

// Like any other API object, the customer notification objects can be deleted by passing an instance of them into this function.  The ID on the object must be set.
func (r User_Customer_Notification_Virtual_Guest) DeleteObjects(templateObjects []datatypes.User_Customer_Notification_Virtual_Guest) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "deleteObjects", params, &r.Options, &resp)
	return
}

// This method returns all SoftLayer_User_Customer_Notification_Virtual_Guest objects associated with the passed in ID as long as that Virtual Guest ID is owned by the current user's account.
//
// This behavior can also be accomplished by simply tapping monitoringUserNotification on the Virtual_Guest object.
func (r User_Customer_Notification_Virtual_Guest) FindByGuestId(id *int) (resp []datatypes.User_Customer_Notification_Virtual_Guest, err error) {
	params := []interface{}{
		id,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "findByGuestId", params, &r.Options, &resp)
	return
}

// Retrieve The virtual guest object that will be monitored.
func (r User_Customer_Notification_Virtual_Guest) GetGuest() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "getGuest", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer_Notification_Virtual_Guest object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_Notification_Virtual_Guest service. You can only retrieve guest notifications attached to virtual guests and users that belong to your account
func (r User_Customer_Notification_Virtual_Guest) GetObject() (resp datatypes.User_Customer_Notification_Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The user that will be notified when the associated virtual guest object fails a monitoring instance.
func (r User_Customer_Notification_Virtual_Guest) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Notification_Virtual_Guest", "getUser", nil, &r.Options, &resp)
	return
}

// no documentation yet
type User_Customer_OpenIdConnect struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerOpenIdConnectService returns an instance of the User_Customer_OpenIdConnect SoftLayer service
func GetUserCustomerOpenIdConnectService(sess session.SLSession) User_Customer_OpenIdConnect {
	return User_Customer_OpenIdConnect{Session: sess}
}

func (r User_Customer_OpenIdConnect) Id(id int) User_Customer_OpenIdConnect {
	r.Options.Id = &id
	return r
}

func (r User_Customer_OpenIdConnect) Mask(mask string) User_Customer_OpenIdConnect {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_OpenIdConnect) Filter(filter string) User_Customer_OpenIdConnect {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_OpenIdConnect) Limit(limit int) User_Customer_OpenIdConnect {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_OpenIdConnect) Offset(offset int) User_Customer_OpenIdConnect {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer_OpenIdConnect) AcknowledgeSupportPolicy() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "acknowledgeSupportPolicy", nil, &r.Options, &resp)
	return
}

// Completes invitation process for an OpenIdConnect user created by Bluemix Unified User Console.
func (r User_Customer_OpenIdConnect) ActivateOpenIdConnectUser(verificationCode *string, userInfo *datatypes.User_Customer, iamId *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		verificationCode,
		userInfo,
		iamId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "activateOpenIdConnectUser", params, &r.Options, &resp)
	return
}

// Create a user's API authentication key, allowing that user access to query the SoftLayer API. addApiAuthenticationKey() returns the user's new API key. Each portal user is allowed only one API key.
func (r User_Customer_OpenIdConnect) AddApiAuthenticationKey() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addApiAuthenticationKey", nil, &r.Options, &resp)
	return
}

// Grants the user access to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) AddBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Add multiple hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkHardwareAccess() does not attempt to add hardware access if the given user already has access to that hardware object.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) AddBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Add multiple permissions to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. addBulkPortalPermission() does not attempt to add permissions already assigned to the user.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer_OpenIdConnect) AddBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permissions,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) AddBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addBulkRoles", params, &r.Options, &resp)
	return
}

// Add multiple CloudLayer Computing Instances to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkVirtualGuestAccess() does not attempt to add CloudLayer Computing Instance access if the given user already has access to that CloudLayer Computing Instance object.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) AddBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Grants the user access to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect) AddDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) AddExternalBinding(externalBinding *datatypes.User_External_Binding) (resp datatypes.User_Customer_External_Binding, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addExternalBinding", params, &r.Options, &resp)
	return
}

// Add hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the hardware you're attempting to add then addHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect) AddHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addHardwareAccess", params, &r.Options, &resp)
	return
}

// Create a notification subscription record for the user. If a subscription record exists for the notification, the record will be set to active, if currently inactive.
func (r User_Customer_OpenIdConnect) AddNotificationSubscriber(notificationKeyName *string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Add a permission to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. If the user already has the permission you're attempting to add then addPortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are added based on the keyName property of the permission parameter.
func (r User_Customer_OpenIdConnect) AddPortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permission,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) AddRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addRole", params, &r.Options, &resp)
	return
}

// Add a CloudLayer Computing Instance to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the CloudLayer Computing Instance you're attempting to add then addVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect) AddVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "addVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method can be used in place of [[SoftLayer_User_Customer::editObject]] to change the parent user of this user.
//
// The new parent must be a user on the same account, and must not be a child of this user.  A user is not allowed to change their own parent.
//
// If the cascadeFlag is set to false, then an exception will be thrown if the new parent does not have all of the permissions that this user possesses.  If the cascadeFlag is set to true, then permissions will be removed from this user and the descendants of this user as necessary so that no children of the parent will have permissions that the parent does not possess. However, setting the cascadeFlag to true will not remove the access all device permissions from this user. The customer portal will need to be used to remove these permissions.
func (r User_Customer_OpenIdConnect) AssignNewParentId(parentId *int, cascadePermissionsFlag *bool) (resp datatypes.User_Customer, err error) {
	params := []interface{}{
		parentId,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "assignNewParentId", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to modify using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer_OpenIdConnect) ChangePreference(preferenceTypeKeyName *string, value *string) (resp []datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
		value,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "changePreference", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) CompleteInvitationAfterLogin(providerType *string, accessToken *string, emailRegistrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		accessToken,
		emailRegistrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "completeInvitationAfterLogin", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer_OpenIdConnect) CreateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "createNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Create a new user in the SoftLayer customer portal. It is not possible to set up SLL enable flags during object creation. These flags are ignored during object creation. You will need to make a subsequent call to edit object in order to enable VPN access.
//
// An account's master user and sub-users who have the User Manage permission can add new users.
//
// Users are created with a default permission set. After adding a user it may be helpful to set their permissions and device access.
//
// secondaryPasswordTimeoutDays will be set to the system configured default value if the attribute is not provided or the attribute is not a valid value.
//
// Note, neither password nor vpnPassword parameters are required.
//
// Password When a new user is created, an email will be sent to the new user's email address with a link to a url that will allow the new user to create or change their password for the SoftLayer customer portal.
//
// If the password parameter is provided and is not null, then that value will be validated. If it is a valid password, then the user will be created with this password.  This user will still receive a portal password email.  It can be used within 24 hours to change their password, or it can be allowed to expire, and the password provided during user creation will remain as the user's password.
//
// If the password parameter is not provided or the value is null, the user must set their portal password using the link sent in email within 24 hours.  If the user fails to set their password within 24 hours, then a non-master user can use the "Reset Password" link on the login page of the portal to request a new email.  A master user can use the link to retrieve a phone number to call to assist in resetting their password.
//
// The password parameter is ignored for VPN_ONLY users or for IBMid authenticated users.
//
// vpnPassword If the vpnPassword is provided, then the user's vpnPassword will be set to the provided password.  When creating a vpn only user, the vpnPassword MUST be supplied.  If the vpnPassword is not provided, then the user will need to use the portal to edit their profile and set the vpnPassword.
//
// IBMid considerations When a SoftLayer account is linked to a Platform Services (PaaS, formerly Bluemix) account, AND the trait on the SoftLayer Account indicating IBMid authentication is set, then SoftLayer will delegate the creation of an ACTIVE user to PaaS. This means that even though the request to create a new user in such an account may start at the IMS API, via this delegation we effectively turn it into a request that is driven by PaaS. In particular this means that any "invitation email" that comes to the user, will come from PaaS, not from IMS via IBMid.
//
// Users created in states other than ACTIVE (for example, a VPN_ONLY user) will be created directly in IMS without delegation (but note that no invitation is sent for a user created in any state other than ACTIVE).
func (r User_Customer_OpenIdConnect) CreateObject(templateObject *datatypes.User_Customer_OpenIdConnect, password *string, vpnPassword *string) (resp datatypes.User_Customer_OpenIdConnect, err error) {
	params := []interface{}{
		templateObject,
		password,
		vpnPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "createObject", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) CreateOpenIdConnectUserAndCompleteInvitation(providerType *string, user *datatypes.User_Customer, password *string, registrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		user,
		password,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "createOpenIdConnectUserAndCompleteInvitation", params, &r.Options, &resp)
	return
}

// Create delivery methods for a notification that the user is subscribed to. Multiple delivery method keyNames can be supplied to create multiple delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer_OpenIdConnect) CreateSubscriberDeliveryMethods(notificationKeyName *string, deliveryMethodKeyNames []string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "createSubscriberDeliveryMethods", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer_OpenIdConnect) DeactivateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "deactivateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Declines an invitation to link an OpenIdConnect identity to a SoftLayer (Atlas) identity and account. Note that this uses a registration code that is likely a one-time-use-only token, so if an invitation has already been processed (accepted or previously declined) it will not be possible to process it a second time.
func (r User_Customer_OpenIdConnect) DeclineInvitation(providerType *string, registrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "declineInvitation", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObject() if you wish to edit a single user account. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer_OpenIdConnect) EditObject(templateObject *datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "editObject", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObjects() if you wish to edit multiple users at once. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer_OpenIdConnect) EditObjects(templateObjects []datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "editObjects", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) FindUserPreference(profileName *string, containerKeyname *string, preferenceKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
		preferenceKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "findUserPreference", params, &r.Options, &resp)
	return
}

// Retrieve The customer account that a user belongs to.
func (r User_Customer_OpenIdConnect) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer_OpenIdConnect) GetActions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getActions", nil, &r.Options, &resp)
	return
}

// The getActiveExternalAuthenticationVendors method will return a list of available external vendors that a SoftLayer user can authenticate against.  The list will only contain vendors for which the user has at least one active external binding.
func (r User_Customer_OpenIdConnect) GetActiveExternalAuthenticationVendors() (resp []datatypes.Container_User_Customer_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getActiveExternalAuthenticationVendors", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's additional email addresses. These email addresses are contacted when updates are made to support tickets.
func (r User_Customer_OpenIdConnect) GetAdditionalEmails() (resp []datatypes.User_Customer_AdditionalEmail, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAdditionalEmails", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetAgentImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAgentImpersonationToken", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetAllowedDedicatedHostIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAllowedDedicatedHostIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetAllowedHardwareIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAllowedHardwareIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetAllowedVirtualGuestIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAllowedVirtualGuestIds", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's API Authentication keys. There is a max limit of one API key per user.
func (r User_Customer_OpenIdConnect) GetApiAuthenticationKeys() (resp []datatypes.User_Customer_ApiAuthentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getApiAuthenticationKeys", nil, &r.Options, &resp)
	return
}

// This method generate user authentication token and return [[SoftLayer_Container_User_Authentication_Token]] object which will be used to authenticate user to login to SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetAuthenticationToken(token *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Authentication_Token, err error) {
	params := []interface{}{
		token,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getAuthenticationToken", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's child users. Some portal users may not have child users.
func (r User_Customer_OpenIdConnect) GetChildUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getChildUsers", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated closed tickets.
func (r User_Customer_OpenIdConnect) GetClosedTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getClosedTickets", nil, &r.Options, &resp)
	return
}

// Retrieve The dedicated hosts to which the user has been granted access.
func (r User_Customer_OpenIdConnect) GetDedicatedHosts() (resp []datatypes.Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getDedicatedHosts", nil, &r.Options, &resp)
	return
}

// This API gets the account associated with the default user for the OpenIdConnect identity that is linked to the current active SoftLayer user identity. When a single active user is found for that IAMid, it becomes the default user and the associated account is returned. When multiple default users are found only the first is preserved and the associated account is returned (remaining defaults see their default flag unset). If the current SoftLayer user identity isn't linked to any OpenIdConnect identity, or if none of the linked users were found as defaults, the API returns null. Invoke this only on IAMid-authenticated users.
func (r User_Customer_OpenIdConnect) GetDefaultAccount(providerType *string) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getDefaultAccount", params, &r.Options, &resp)
	return
}

// Retrieve The external authentication bindings that link an external identifier to a SoftLayer user.
func (r User_Customer_OpenIdConnect) GetExternalBindings() (resp []datatypes.User_External_Binding, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getExternalBindings", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible hardware. These permissions control which hardware a user has access to in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetHardware() (resp []datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHardware", nil, &r.Options, &resp)
	return
}

// Retrieve the number of servers that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on hardware. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer_OpenIdConnect) GetHardwareCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHardwareCount", nil, &r.Options, &resp)
	return
}

// Retrieve Hardware notifications associated with this user. A hardware notification links a user to a piece of hardware, and that user will be notified if any monitors on that hardware fail, if the monitors have a status of 'Notify User'.
func (r User_Customer_OpenIdConnect) GetHardwareNotifications() (resp []datatypes.User_Customer_Notification_Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHardwareNotifications", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user has acknowledged the support policy.
func (r User_Customer_OpenIdConnect) GetHasAcknowledgedSupportPolicyFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHasAcknowledgedSupportPolicyFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Permission granting the user access to all Dedicated Host devices on the account.
func (r User_Customer_OpenIdConnect) GetHasFullDedicatedHostAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHasFullDedicatedHostAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all hardware on their account.
func (r User_Customer_OpenIdConnect) GetHasFullHardwareAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHasFullHardwareAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all virtual guests on their account.
func (r User_Customer_OpenIdConnect) GetHasFullVirtualGuestAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getHasFullVirtualGuestAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Specifically relating the Customer instance to an IBMid. A Customer instance may or may not have an IBMid link.
func (r User_Customer_OpenIdConnect) GetIbmIdLink() (resp datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getIbmIdLink", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getImpersonationToken", nil, &r.Options, &resp)
	return
}

// Retrieve Contains the definition of the layout profile.
func (r User_Customer_OpenIdConnect) GetLayoutProfiles() (resp []datatypes.Layout_Profile, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getLayoutProfiles", nil, &r.Options, &resp)
	return
}

// Retrieve A user's locale. Locale holds user's language and region information.
func (r User_Customer_OpenIdConnect) GetLocale() (resp datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getLocale", nil, &r.Options, &resp)
	return
}

// Validates a supplied OpenIdConnect access token to the SoftLayer customer portal and returns the default account name and id for the active user. An exception will be thrown if no matching customer is found.
func (r User_Customer_OpenIdConnect) GetLoginAccountInfoOpenIdConnect(providerType *string, accessToken *string) (resp datatypes.Container_User_Customer_OpenIdConnect_LoginAccountInfo, err error) {
	params := []interface{}{
		providerType,
		accessToken,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getLoginAccountInfoOpenIdConnect", params, &r.Options, &resp)
	return
}

// Retrieve A user's attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetLoginAttempts() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getLoginAttempts", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a user to the SoftLayer customer portal using the provided authentication container. Depending on the specific type of authentication container that is used, this API will leverage the appropriate authentication protocol. If authentication is successful then the API returns a list of linked accounts for the user, a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer_OpenIdConnect) GetLoginToken(request *datatypes.Container_Authentication_Request_Contract) (resp datatypes.Container_Authentication_Response_Common, err error) {
	params := []interface{}{
		request,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getLoginToken", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IAMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one SoftLayer user per account. This effectively links the OpenIdConnect identity to those accounts. This API returns a list of all active accounts for which there is a link between the OpenIdConnect identity and a SoftLayer user. Invoke this only on IAMid-authenticated users.
func (r User_Customer_OpenIdConnect) GetMappedAccounts(providerType *string) (resp []datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getMappedAccounts", params, &r.Options, &resp)
	return
}

// Retrieve Notification subscription records for the user.
func (r User_Customer_OpenIdConnect) GetNotificationSubscribers() (resp []datatypes.Notification_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getNotificationSubscribers", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetObject() (resp datatypes.User_Customer_OpenIdConnect, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getObject", nil, &r.Options, &resp)
	return
}

// This API returns a SoftLayer_Container_User_Customer_OpenIdConnect_MigrationState object containing the necessary information to determine what migration state the user is in. If the account is not OpenIdConnect authenticated, then an exception is thrown.
func (r User_Customer_OpenIdConnect) GetOpenIdConnectMigrationState() (resp datatypes.Container_User_Customer_OpenIdConnect_MigrationState, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getOpenIdConnectMigrationState", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetOpenIdRegistrationInfoFromCode(providerType *string, registrationCode *string) (resp datatypes.Account_Authentication_OpenIdConnect_RegistrationInformation, err error) {
	params := []interface{}{
		providerType,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getOpenIdRegistrationInfoFromCode", params, &r.Options, &resp)
	return
}

// Retrieve An user's associated open tickets.
func (r User_Customer_OpenIdConnect) GetOpenTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getOpenTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's vpn accessible subnets.
func (r User_Customer_OpenIdConnect) GetOverrides() (resp []datatypes.Network_Service_Vpn_Overrides, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getOverrides", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's parent user. If a SoftLayer_User_Customer has a null parentId property then it doesn't have a parent user.
func (r User_Customer_OpenIdConnect) GetParent() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getParent", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetPasswordRequirements(isVpn *bool) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		isVpn,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPasswordRequirements", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's permissions. These permissions control that user's access to functions within the SoftLayer customer portal and API.
func (r User_Customer_OpenIdConnect) GetPermissions() (resp []datatypes.User_Customer_CustomerPermission_Permission, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPermissions", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a username and password to the SoftLayer customer portal. Many portal user accounts are configured to require answering a security question on login. In this case getPortalLoginToken() also verifies the given security question ID and answer. If authentication is successful then the API returns a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer_OpenIdConnect) GetPortalLoginToken(username *string, password *string, securityQuestionId *int, securityQuestionAnswer *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		username,
		password,
		securityQuestionId,
		securityQuestionAnswer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPortalLoginToken", params, &r.Options, &resp)
	return
}

// Attempt to authenticate a supplied OpenIdConnect access token to the SoftLayer customer portal. If authentication is successful then the API returns a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer_OpenIdConnect) GetPortalLoginTokenOpenIdConnect(providerType *string, accessToken *string, accountId *int, securityQuestionId *int, securityQuestionAnswer *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		providerType,
		accessToken,
		accountId,
		securityQuestionId,
		securityQuestionAnswer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPortalLoginTokenOpenIdConnect", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to get using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer_OpenIdConnect) GetPreference(preferenceTypeKeyName *string) (resp datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPreference", params, &r.Options, &resp)
	return
}

// Use any of the preference types to fetch or modify user preferences using [[SoftLayer_User_Customer::getPreference|getPreference]] or [[SoftLayer_User_Customer::changePreference|changePreference]], respectively.
func (r User_Customer_OpenIdConnect) GetPreferenceTypes() (resp []datatypes.User_Preference_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPreferenceTypes", nil, &r.Options, &resp)
	return
}

// Retrieve Data type contains a single user preference to a specific preference type.
func (r User_Customer_OpenIdConnect) GetPreferences() (resp []datatypes.User_Preference, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getPreferences", nil, &r.Options, &resp)
	return
}

// Retrieve the authentication requirements for an outstanding password set/reset request.  The requirements returned in the same SoftLayer_Container_User_Customer_PasswordSet container which is provided as a parameter into this request.  The SoftLayer_Container_User_Customer_PasswordSet::authenticationMethods array will contain an entry for each authentication method required for the user.  See SoftLayer_Container_User_Customer_PasswordSet for more details.
//
// If the user has required authentication methods, then authentication information will be supplied to the SoftLayer_User_Customer::processPasswordSetRequest method within this same SoftLayer_Container_User_Customer_PasswordSet container.  All existing information in the container must continue to exist in the container to complete the password set/reset process.
func (r User_Customer_OpenIdConnect) GetRequirementsForPasswordSet(passwordSet *datatypes.Container_User_Customer_PasswordSet) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		passwordSet,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getRequirementsForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer_OpenIdConnect) GetRoles() (resp []datatypes.User_Permission_Role, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getRoles", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's security question answers. Some portal users may not have security answers or may not be configured to require answering a security question on login.
func (r User_Customer_OpenIdConnect) GetSecurityAnswers() (resp []datatypes.User_Customer_Security_Answer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's notification subscription records.
func (r User_Customer_OpenIdConnect) GetSubscribers() (resp []datatypes.Notification_User_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSubscribers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's successful attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetSuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user is required to acknowledge the support policy for portal access.
func (r User_Customer_OpenIdConnect) GetSupportPolicyAcknowledgementRequiredFlag() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSupportPolicyAcknowledgementRequiredFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetSupportPolicyDocument() (resp []byte, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSupportPolicyDocument", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetSupportPolicyName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSupportPolicyName", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetSupportedLocales() (resp []datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSupportedLocales", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user must take a brief survey the next time they log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetSurveyRequiredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSurveyRequiredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The surveys that a user has taken in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetSurveys() (resp []datatypes.Survey, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getSurveys", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated tickets.
func (r User_Customer_OpenIdConnect) GetTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's time zone.
func (r User_Customer_OpenIdConnect) GetTimezone() (resp datatypes.Locale_Timezone, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getTimezone", nil, &r.Options, &resp)
	return
}

// Retrieve A user's unsuccessful attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetUnsuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUnsuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Returns an IMS User Object from the provided OpenIdConnect User ID or IBMid Unique Identifier for the Account of the active user. Enforces the User Management permissions for the Active User. An exception will be thrown if no matching IMS User is found. NOTE that providing IBMid Unique Identifier is optional, but it will be preferred over OpenIdConnect User ID if provided.
func (r User_Customer_OpenIdConnect) GetUserForUnifiedInvitation(openIdConnectUserId *string, uniqueIdentifier *string, searchInvitationsNotLinksFlag *string, accountId *string) (resp datatypes.User_Customer_OpenIdConnect, err error) {
	params := []interface{}{
		openIdConnectUserId,
		uniqueIdentifier,
		searchInvitationsNotLinksFlag,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUserForUnifiedInvitation", params, &r.Options, &resp)
	return
}

// Retrieve a user id using a password token provided to the user in an email generated by the SoftLayer_User_Customer::initiatePortalPasswordChange request. Password recovery keys are valid for 24 hours after they're generated.
//
// When a new user is created or when a user has requested a password change using initiatePortalPasswordChange, they will have received an email that contains a url with a token.  That token is used as the parameter for getUserIdForPasswordSet.  Once the user id is known, then the SoftLayer_User_Customer object can be retrieved which is necessary to complete the process to set or reset a user's password.
func (r User_Customer_OpenIdConnect) GetUserIdForPasswordSet(key *string) (resp int, err error) {
	params := []interface{}{
		key,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUserIdForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve User customer link with IBMid and IAMid.
func (r User_Customer_OpenIdConnect) GetUserLinks() (resp []datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUserLinks", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) GetUserPreferences(profileName *string, containerKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUserPreferences", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's status, which controls overall access to the SoftLayer customer portal and VPN access to the private network.
func (r User_Customer_OpenIdConnect) GetUserStatus() (resp datatypes.User_Customer_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getUserStatus", nil, &r.Options, &resp)
	return
}

// Retrieve the number of CloudLayer Computing Instances that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on CloudLayer Computing Instances. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer_OpenIdConnect) GetVirtualGuestCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getVirtualGuestCount", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible CloudLayer Computing Instances. These permissions control which CloudLayer Computing Instances a user has access to in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) GetVirtualGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "getVirtualGuests", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) InTerminalStatus() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "inTerminalStatus", nil, &r.Options, &resp)
	return
}

// Sends password change email to the user containing url that allows the user the change their password. This is the first step when a user wishes to change their password.  The url that is generated contains a one-time use token that is valid for only 24-hours.
//
// If this is a new master user who has never logged into the portal, then password reset will be initiated. Once a master user has logged into the portal, they must setup their security questions prior to logging out because master users are required to answer a security question during the password reset process.  Should a master user not have security questions defined and not remember their password in order to define the security questions, then they will need to contact support at live chat or Revenue Services for assistance.
//
// Due to security reasons, the number of reset requests per username are limited within a undisclosed timeframe.
func (r User_Customer_OpenIdConnect) InitiatePortalPasswordChange(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "initiatePortalPasswordChange", params, &r.Options, &resp)
	return
}

// A Brand Agent that has permissions to Add Customer Accounts will be able to request the password email be sent to the Master User of a Customer Account created by the same Brand as the agent making the request. Due to security reasons, the number of reset requests are limited within an undisclosed timeframe.
func (r User_Customer_OpenIdConnect) InitiatePortalPasswordChangeByBrandAgent(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "initiatePortalPasswordChangeByBrandAgent", params, &r.Options, &resp)
	return
}

// Send email invitation to a user to join a SoftLayer account and authenticate with OpenIdConnect. Throws an exception on error.
func (r User_Customer_OpenIdConnect) InviteUserToLinkOpenIdConnect(providerType *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "inviteUserToLinkOpenIdConnect", params, &r.Options, &resp)
	return
}

// Portal users are considered master users if they don't have an associated parent user. The only users who don't have parent users are users whose username matches their SoftLayer account name. Master users have special permissions throughout the SoftLayer customer portal.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer_OpenIdConnect) IsMasterUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "isMasterUser", nil, &r.Options, &resp)
	return
}

// Determine if a string is the given user's login password to the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect) IsValidPortalPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "isValidPortalPassword", params, &r.Options, &resp)
	return
}

// The perform external authentication method will authenticate the given external authentication container with an external vendor.  The authentication container and its contents will be verified before an attempt is made to authenticate the contents of the container with an external vendor.
func (r User_Customer_OpenIdConnect) PerformExternalAuthentication(authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "performExternalAuthentication", params, &r.Options, &resp)
	return
}

// Set the password for a user who has an outstanding password request. A user with an outstanding password request will have an unused and unexpired password key.  The password key is part of the url provided to the user in the email sent to the user with information on how to set their password.  The email was generated by the SoftLayer_User_Customer::initiatePortalPasswordRequest request. Password recovery keys are valid for 24 hours after they're generated.
//
// If the user has required authentication methods as specified by in the SoftLayer_Container_User_Customer_PasswordSet container returned from the SoftLayer_User_Customer::getRequirementsForPasswordSet request, then additional requests must be made to processPasswordSetRequest to authenticate the user before changing the password.  First, if the user has security questions set on their profile, they will be required to answer one of their questions correctly. Next, if the user has Verisign or Google Authentication on their account, they must authenticate according to the two-factor provider.  All of this authentication is done using the SoftLayer_Container_User_Customer_PasswordSet container.
//
// User portal passwords must match the following restrictions. Portal passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ + =
// * ...not match your username
func (r User_Customer_OpenIdConnect) ProcessPasswordSetRequest(passwordSet *datatypes.Container_User_Customer_PasswordSet, authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp bool, err error) {
	params := []interface{}{
		passwordSet,
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "processPasswordSetRequest", params, &r.Options, &resp)
	return
}

// Revoke access to all dedicated hosts on the account for this user. The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect) RemoveAllDedicatedHostAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeAllDedicatedHostAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect) RemoveAllHardwareAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeAllHardwareAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all cloud computing instances from a portal user's instance access list. A user's instance access list controls which of an account's computing instance objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect) RemoveAllVirtualAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeAllVirtualAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove a user's API authentication key, removing that user's access to query the SoftLayer API.
func (r User_Customer_OpenIdConnect) RemoveApiAuthenticationKey(keyId *int) (resp bool, err error) {
	params := []interface{}{
		keyId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeApiAuthenticationKey", params, &r.Options, &resp)
	return
}

// Revokes access for the user to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// If the user has full dedicatedHost access, then it will provide access to "ALL but passed in" dedicatedHost ids.
func (r User_Customer_OpenIdConnect) RemoveBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Remove multiple hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting to remove then removeBulkHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// If the user has full hardware access, then it will provide access to "ALL but passed in" hardware ids.
func (r User_Customer_OpenIdConnect) RemoveBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) multiple permissions from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. removePortalPermission() does not attempt to remove permissions that are not assigned to the user.
//
// Users can grant or revoke permissions to their child users, but not to themselves. An account's master has all portal permissions and can grant permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permissions from a user will cascade down the child hierarchy and remove the permissions from this user along with all child users who also have the permission.
//
// If the cascadePermissionsFlag is not provided or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer_OpenIdConnect) RemoveBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permissions,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) RemoveBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeBulkRoles", params, &r.Options, &resp)
	return
}

// Remove multiple CloudLayer Computing Instances from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeBulkVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) RemoveBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Revokes access for the user to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) RemoveDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) RemoveExternalBinding(externalBinding *datatypes.User_External_Binding) (resp bool, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeExternalBinding", params, &r.Options, &resp)
	return
}

// Remove hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting remove add then removeHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) RemoveHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) a permission from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. If the user does not have the permission you're attempting to remove then removePortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permission from a user will cascade down the child hierarchy and remove the permission from this user and all child users who also have the permission.
//
// If the cascadePermissionsFlag is not set or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission parameter.
func (r User_Customer_OpenIdConnect) RemovePortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permission,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removePortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) RemoveRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeRole", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) RemoveSecurityAnswers() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Remove a CloudLayer Computing Instance from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's computing instances a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set instance access for any of the other users on their account.
func (r User_Customer_OpenIdConnect) RemoveVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "removeVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer user is linked to, if we need to do that for some reason. It will do this by modifying the link to the desired new IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is already Bluemix linked. To reset a link for the Bluemix-linked user account, use resetOpenIdConnectLinkUnifiedUserManagementMode.
func (r User_Customer_OpenIdConnect) ResetOpenIdConnectLink(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "resetOpenIdConnectLink", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer master user is linked to, if we need to do that for some reason. It will do this by unlinking the new owner IBMid from its current user association in this account, if there is one (note that the new owner IBMid is not required to already be a member of the IMS account). Then it will modify the existing IBMid link for the master user to use the new owner IBMid-realm IAMid. At this point, if the new owner IBMid isn't already a member of the PaaS account, it will attempt to add it. As a last step, it will call PaaS to modify the owner on that side, if necessary.  Only when all those steps are complete, it will commit the IMS-side DB changes.  Then, it will clean up the SoftLayer user that was linked to the new owner IBMid (this user became unlinked as the first step in this process).  It will also call BSS to delete the old owner IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is not Bluemix linked. To reset a link for the user account not linked to Bluemix, use resetOpenIdConnectLink.
func (r User_Customer_OpenIdConnect) ResetOpenIdConnectLinkUnifiedUserManagementMode(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "resetOpenIdConnectLinkUnifiedUserManagementMode", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) SamlAuthenticate(accountId *string, samlResponse *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		accountId,
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "samlAuthenticate", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) SamlBeginAuthentication(accountId *int) (resp string, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "samlBeginAuthentication", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) SamlBeginLogout() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "samlBeginLogout", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) SamlLogout(samlResponse *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "samlLogout", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect) SelfPasswordChange(currentPassword *string, newPassword *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		currentPassword,
		newPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "selfPasswordChange", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IAMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one per account. If an OpenIdConnect identity is mapped to multiple accounts in this manner, one such account should be identified as the default account for that identity. Invoke this only on IBMid-authenticated users.
func (r User_Customer_OpenIdConnect) SetDefaultAccount(providerType *string, accountId *int) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "setDefaultAccount", params, &r.Options, &resp)
	return
}

// As master user, calling this api for the IBMid provider type when there is an existing IBMid for the email on the SL account will silently (without sending an invitation email) create a link for the IBMid. NOTE: If the SoftLayer user is already linked to IBMid, this call will fail. If the IBMid specified by the email of this user, is already used in a link to another user in this account, this call will fail. If there is already an open invitation from this SoftLayer user to this or any IBMid, this call will fail. If there is already an open invitation from some other SoftLayer user in this account to this IBMid, then this call will fail.
func (r User_Customer_OpenIdConnect) SilentlyMigrateUserOpenIdConnect(providerType *string) (resp bool, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "silentlyMigrateUserOpenIdConnect", params, &r.Options, &resp)
	return
}

// This method allows the master user of an account to undo the designation of this user as an alternate master user.  This can not be applied to the true master user of the account.
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted this user.  This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner turns off an "alternate/auxiliary master user / account owner".
func (r User_Customer_OpenIdConnect) TurnOffMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "turnOffMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// This method allows the master user of an account to designate this user as an alternate master user.  Effectively this means that this user should have "all the same IMS permissions as a master user".
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted to this user. This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner designates an "alternate/auxiliary master user / account owner".
func (r User_Customer_OpenIdConnect) TurnOnMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "turnOnMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// Update the active status for a notification that the user is subscribed to. A notification along with an active flag can be supplied to update the active status for a particular notification subscription.
func (r User_Customer_OpenIdConnect) UpdateNotificationSubscriber(notificationKeyName *string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "updateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Update a user's login security questions and answers on the SoftLayer customer portal. These questions and answers are used to optionally log into the SoftLayer customer portal using two-factor authentication. Each user must have three distinct questions set with a unique answer for each question, and each answer may only contain alphanumeric or the . , - _ ( ) [ ] : ; > < characters. Existing user security questions and answers are deleted before new ones are set, and users may only update their own security questions and answers.
func (r User_Customer_OpenIdConnect) UpdateSecurityAnswers(questions []datatypes.User_Security_Question, answers []string) (resp bool, err error) {
	params := []interface{}{
		questions,
		answers,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "updateSecurityAnswers", params, &r.Options, &resp)
	return
}

// Update a delivery method for a notification that the user is subscribed to. A delivery method keyName along with an active flag can be supplied to update the active status of the delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer_OpenIdConnect) UpdateSubscriberDeliveryMethod(notificationKeyName *string, deliveryMethodKeyNames []string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "updateSubscriberDeliveryMethod", params, &r.Options, &resp)
	return
}

// Update a user's VPN password on the SoftLayer customer portal. As with portal passwords, VPN passwords must match the following restrictions. VPN passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ =
// * ...not match your username
// Finally, users can only update their own VPN password. An account's master user can update any of their account users' VPN passwords.
func (r User_Customer_OpenIdConnect) UpdateVpnPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "updateVpnPassword", params, &r.Options, &resp)
	return
}

// Always call this function to enable changes when manually configuring VPN subnet access.
func (r User_Customer_OpenIdConnect) UpdateVpnUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "updateVpnUser", nil, &r.Options, &resp)
	return
}

// This method validate the given authentication token using the user id by comparing it with the actual user authentication token and return [[SoftLayer_Container_User_Customer_Portal_Token]] object
func (r User_Customer_OpenIdConnect) ValidateAuthenticationToken(authenticationToken *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationToken,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect", "validateAuthenticationToken", params, &r.Options, &resp)
	return
}

// no documentation yet
type User_Customer_OpenIdConnect_TrustedProfile struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerOpenIdConnectTrustedProfileService returns an instance of the User_Customer_OpenIdConnect_TrustedProfile SoftLayer service
func GetUserCustomerOpenIdConnectTrustedProfileService(sess session.SLSession) User_Customer_OpenIdConnect_TrustedProfile {
	return User_Customer_OpenIdConnect_TrustedProfile{Session: sess}
}

func (r User_Customer_OpenIdConnect_TrustedProfile) Id(id int) User_Customer_OpenIdConnect_TrustedProfile {
	r.Options.Id = &id
	return r
}

func (r User_Customer_OpenIdConnect_TrustedProfile) Mask(mask string) User_Customer_OpenIdConnect_TrustedProfile {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_OpenIdConnect_TrustedProfile) Filter(filter string) User_Customer_OpenIdConnect_TrustedProfile {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_OpenIdConnect_TrustedProfile) Limit(limit int) User_Customer_OpenIdConnect_TrustedProfile {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_OpenIdConnect_TrustedProfile) Offset(offset int) User_Customer_OpenIdConnect_TrustedProfile {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) AcknowledgeSupportPolicy() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "acknowledgeSupportPolicy", nil, &r.Options, &resp)
	return
}

// Completes invitation process for an OpenIdConnect user created by Bluemix Unified User Console.
func (r User_Customer_OpenIdConnect_TrustedProfile) ActivateOpenIdConnectUser(verificationCode *string, userInfo *datatypes.User_Customer, iamId *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		verificationCode,
		userInfo,
		iamId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "activateOpenIdConnectUser", params, &r.Options, &resp)
	return
}

// Create a user's API authentication key, allowing that user access to query the SoftLayer API. addApiAuthenticationKey() returns the user's new API key. Each portal user is allowed only one API key.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddApiAuthenticationKey() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addApiAuthenticationKey", nil, &r.Options, &resp)
	return
}

// Grants the user access to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Add multiple hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkHardwareAccess() does not attempt to add hardware access if the given user already has access to that hardware object.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Add multiple permissions to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. addBulkPortalPermission() does not attempt to add permissions already assigned to the user.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permissions,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) AddBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addBulkRoles", params, &r.Options, &resp)
	return
}

// Add multiple CloudLayer Computing Instances to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. addBulkVirtualGuestAccess() does not attempt to add CloudLayer Computing Instance access if the given user already has access to that CloudLayer Computing Instance object.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Grants the user access to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) AddExternalBinding(externalBinding *datatypes.User_External_Binding) (resp datatypes.User_Customer_External_Binding, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addExternalBinding", params, &r.Options, &resp)
	return
}

// Add hardware to a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the hardware you're attempting to add then addHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addHardwareAccess", params, &r.Options, &resp)
	return
}

// Create a notification subscription record for the user. If a subscription record exists for the notification, the record will be set to active, if currently inactive.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddNotificationSubscriber(notificationKeyName *string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Add a permission to a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. If the user already has the permission you're attempting to add then addPortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are added based on the keyName property of the permission parameter.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddPortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission) (resp bool, err error) {
	params := []interface{}{
		permission,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) AddRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addRole", params, &r.Options, &resp)
	return
}

// Add a CloudLayer Computing Instance to a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user already has access to the CloudLayer Computing Instance you're attempting to add then addVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set CloudLayer Computing Instance access for any of the other users on their account.
//
// Only the USER_MANAGE permission is required to execute this.
func (r User_Customer_OpenIdConnect_TrustedProfile) AddVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "addVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method can be used in place of [[SoftLayer_User_Customer::editObject]] to change the parent user of this user.
//
// The new parent must be a user on the same account, and must not be a child of this user.  A user is not allowed to change their own parent.
//
// If the cascadeFlag is set to false, then an exception will be thrown if the new parent does not have all of the permissions that this user possesses.  If the cascadeFlag is set to true, then permissions will be removed from this user and the descendants of this user as necessary so that no children of the parent will have permissions that the parent does not possess. However, setting the cascadeFlag to true will not remove the access all device permissions from this user. The customer portal will need to be used to remove these permissions.
func (r User_Customer_OpenIdConnect_TrustedProfile) AssignNewParentId(parentId *int, cascadePermissionsFlag *bool) (resp datatypes.User_Customer, err error) {
	params := []interface{}{
		parentId,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "assignNewParentId", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to modify using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer_OpenIdConnect_TrustedProfile) ChangePreference(preferenceTypeKeyName *string, value *string) (resp []datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
		value,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "changePreference", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) CompleteInvitationAfterLogin(providerType *string, accessToken *string, emailRegistrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		accessToken,
		emailRegistrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "completeInvitationAfterLogin", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer_OpenIdConnect_TrustedProfile) CreateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "createNotificationSubscriber", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) CreateObject(templateObject *datatypes.User_Customer_OpenIdConnect_TrustedProfile, password *string, vpnPassword *string) (resp datatypes.User_Customer_OpenIdConnect_TrustedProfile, err error) {
	params := []interface{}{
		templateObject,
		password,
		vpnPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "createObject", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) CreateOpenIdConnectUserAndCompleteInvitation(providerType *string, user *datatypes.User_Customer, password *string, registrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		user,
		password,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "createOpenIdConnectUserAndCompleteInvitation", params, &r.Options, &resp)
	return
}

// Create delivery methods for a notification that the user is subscribed to. Multiple delivery method keyNames can be supplied to create multiple delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer_OpenIdConnect_TrustedProfile) CreateSubscriberDeliveryMethods(notificationKeyName *string, deliveryMethodKeyNames []string) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "createSubscriberDeliveryMethods", params, &r.Options, &resp)
	return
}

// Create a new subscriber for a given resource.
func (r User_Customer_OpenIdConnect_TrustedProfile) DeactivateNotificationSubscriber(keyName *string, resourceTableId *int) (resp bool, err error) {
	params := []interface{}{
		keyName,
		resourceTableId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "deactivateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Declines an invitation to link an OpenIdConnect identity to a SoftLayer (Atlas) identity and account. Note that this uses a registration code that is likely a one-time-use-only token, so if an invitation has already been processed (accepted or previously declined) it will not be possible to process it a second time.
func (r User_Customer_OpenIdConnect_TrustedProfile) DeclineInvitation(providerType *string, registrationCode *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "declineInvitation", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObject() if you wish to edit a single user account. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer_OpenIdConnect_TrustedProfile) EditObject(templateObject *datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "editObject", params, &r.Options, &resp)
	return
}

// Account master users and sub-users who have the User Manage permission in the SoftLayer customer portal can update other user's information. Use editObjects() if you wish to edit multiple users at once. Users who do not have the User Manage permission can only update their own information.
func (r User_Customer_OpenIdConnect_TrustedProfile) EditObjects(templateObjects []datatypes.User_Customer) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "editObjects", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) FindUserPreference(profileName *string, containerKeyname *string, preferenceKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
		preferenceKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "findUserPreference", params, &r.Options, &resp)
	return
}

// Retrieve The customer account that a user belongs to.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer_OpenIdConnect_TrustedProfile) GetActions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getActions", nil, &r.Options, &resp)
	return
}

// The getActiveExternalAuthenticationVendors method will return a list of available external vendors that a SoftLayer user can authenticate against.  The list will only contain vendors for which the user has at least one active external binding.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetActiveExternalAuthenticationVendors() (resp []datatypes.Container_User_Customer_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getActiveExternalAuthenticationVendors", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's additional email addresses. These email addresses are contacted when updates are made to support tickets.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAdditionalEmails() (resp []datatypes.User_Customer_AdditionalEmail, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAdditionalEmails", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAgentImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAgentImpersonationToken", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAllowedDedicatedHostIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAllowedDedicatedHostIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAllowedHardwareIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAllowedHardwareIds", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAllowedVirtualGuestIds() (resp []int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAllowedVirtualGuestIds", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's API Authentication keys. There is a max limit of one API key per user.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetApiAuthenticationKeys() (resp []datatypes.User_Customer_ApiAuthentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getApiAuthenticationKeys", nil, &r.Options, &resp)
	return
}

// This method generate user authentication token and return [[SoftLayer_Container_User_Authentication_Token]] object which will be used to authenticate user to login to SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetAuthenticationToken(token *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Authentication_Token, err error) {
	params := []interface{}{
		token,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getAuthenticationToken", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's child users. Some portal users may not have child users.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetChildUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getChildUsers", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated closed tickets.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetClosedTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getClosedTickets", nil, &r.Options, &resp)
	return
}

// Retrieve The dedicated hosts to which the user has been granted access.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetDedicatedHosts() (resp []datatypes.Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getDedicatedHosts", nil, &r.Options, &resp)
	return
}

// This API gets the account associated with the default user for the OpenIdConnect identity that is linked to the current active SoftLayer user identity. When a single active user is found for that IAMid, it becomes the default user and the associated account is returned. When multiple default users are found only the first is preserved and the associated account is returned (remaining defaults see their default flag unset). If the current SoftLayer user identity isn't linked to any OpenIdConnect identity, or if none of the linked users were found as defaults, the API returns null. Invoke this only on IAMid-authenticated users.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetDefaultAccount(providerType *string) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getDefaultAccount", params, &r.Options, &resp)
	return
}

// Retrieve The external authentication bindings that link an external identifier to a SoftLayer user.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetExternalBindings() (resp []datatypes.User_External_Binding, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getExternalBindings", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible hardware. These permissions control which hardware a user has access to in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHardware() (resp []datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHardware", nil, &r.Options, &resp)
	return
}

// Retrieve the number of servers that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on hardware. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHardwareCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHardwareCount", nil, &r.Options, &resp)
	return
}

// Retrieve Hardware notifications associated with this user. A hardware notification links a user to a piece of hardware, and that user will be notified if any monitors on that hardware fail, if the monitors have a status of 'Notify User'.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHardwareNotifications() (resp []datatypes.User_Customer_Notification_Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHardwareNotifications", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user has acknowledged the support policy.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHasAcknowledgedSupportPolicyFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHasAcknowledgedSupportPolicyFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Permission granting the user access to all Dedicated Host devices on the account.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHasFullDedicatedHostAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHasFullDedicatedHostAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all hardware on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHasFullHardwareAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHasFullHardwareAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a portal user has access to all virtual guests on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetHasFullVirtualGuestAccessFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getHasFullVirtualGuestAccessFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Specifically relating the Customer instance to an IBMid. A Customer instance may or may not have an IBMid link.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetIbmIdLink() (resp datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getIbmIdLink", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetImpersonationToken() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getImpersonationToken", nil, &r.Options, &resp)
	return
}

// Retrieve Contains the definition of the layout profile.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetLayoutProfiles() (resp []datatypes.Layout_Profile, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getLayoutProfiles", nil, &r.Options, &resp)
	return
}

// Retrieve A user's locale. Locale holds user's language and region information.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetLocale() (resp datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getLocale", nil, &r.Options, &resp)
	return
}

// Validates a supplied OpenIdConnect access token to the SoftLayer customer portal and returns the default account name and id for the active user. An exception will be thrown if no matching customer is found.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetLoginAccountInfoOpenIdConnect(providerType *string, accessToken *string) (resp datatypes.Container_User_Customer_OpenIdConnect_LoginAccountInfo, err error) {
	params := []interface{}{
		providerType,
		accessToken,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getLoginAccountInfoOpenIdConnect", params, &r.Options, &resp)
	return
}

// Retrieve A user's attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetLoginAttempts() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getLoginAttempts", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a user to the SoftLayer customer portal using the provided authentication container. Depending on the specific type of authentication container that is used, this API will leverage the appropriate authentication protocol. If authentication is successful then the API returns a list of linked accounts for the user, a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetLoginToken(request *datatypes.Container_Authentication_Request_Contract) (resp datatypes.Container_Authentication_Response_Common, err error) {
	params := []interface{}{
		request,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getLoginToken", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IAMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one SoftLayer user per account. This effectively links the OpenIdConnect identity to those accounts. This API returns a list of all active accounts for which there is a link between the OpenIdConnect identity and a SoftLayer user. Invoke this only on IAMid-authenticated users.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetMappedAccounts(providerType *string) (resp []datatypes.Account, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getMappedAccounts", params, &r.Options, &resp)
	return
}

// Retrieve Notification subscription records for the user.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetNotificationSubscribers() (resp []datatypes.Notification_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getNotificationSubscribers", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetObject() (resp datatypes.User_Customer_OpenIdConnect_TrustedProfile, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getObject", nil, &r.Options, &resp)
	return
}

// This API returns a SoftLayer_Container_User_Customer_OpenIdConnect_MigrationState object containing the necessary information to determine what migration state the user is in. If the account is not OpenIdConnect authenticated, then an exception is thrown.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetOpenIdConnectMigrationState() (resp datatypes.Container_User_Customer_OpenIdConnect_MigrationState, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getOpenIdConnectMigrationState", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetOpenIdRegistrationInfoFromCode(providerType *string, registrationCode *string) (resp datatypes.Account_Authentication_OpenIdConnect_RegistrationInformation, err error) {
	params := []interface{}{
		providerType,
		registrationCode,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getOpenIdRegistrationInfoFromCode", params, &r.Options, &resp)
	return
}

// Retrieve An user's associated open tickets.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetOpenTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getOpenTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's vpn accessible subnets.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetOverrides() (resp []datatypes.Network_Service_Vpn_Overrides, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getOverrides", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's parent user. If a SoftLayer_User_Customer has a null parentId property then it doesn't have a parent user.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetParent() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getParent", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPasswordRequirements(isVpn *bool) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		isVpn,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPasswordRequirements", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's permissions. These permissions control that user's access to functions within the SoftLayer customer portal and API.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPermissions() (resp []datatypes.User_Customer_CustomerPermission_Permission, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPermissions", nil, &r.Options, &resp)
	return
}

// Attempt to authenticate a username and password to the SoftLayer customer portal. Many portal user accounts are configured to require answering a security question on login. In this case getPortalLoginToken() also verifies the given security question ID and answer. If authentication is successful then the API returns a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPortalLoginToken(username *string, password *string, securityQuestionId *int, securityQuestionAnswer *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		username,
		password,
		securityQuestionId,
		securityQuestionAnswer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPortalLoginToken", params, &r.Options, &resp)
	return
}

// Attempt to authenticate a supplied OpenIdConnect access token to the SoftLayer customer portal. If authentication is successful then the API returns a token containing the ID of the authenticated user and a hash key used by the SoftLayer customer portal to maintain authentication.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPortalLoginTokenOpenIdConnect(providerType *string, accessToken *string, accountId *int, securityQuestionId *int, securityQuestionAnswer *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		providerType,
		accessToken,
		accountId,
		securityQuestionId,
		securityQuestionAnswer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPortalLoginTokenOpenIdConnect", params, &r.Options, &resp)
	return
}

// Select a type of preference you would like to get using [[SoftLayer_User_Customer::getPreferenceTypes|getPreferenceTypes]] and invoke this method using that preference type key name.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPreference(preferenceTypeKeyName *string) (resp datatypes.User_Preference, err error) {
	params := []interface{}{
		preferenceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPreference", params, &r.Options, &resp)
	return
}

// Use any of the preference types to fetch or modify user preferences using [[SoftLayer_User_Customer::getPreference|getPreference]] or [[SoftLayer_User_Customer::changePreference|changePreference]], respectively.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPreferenceTypes() (resp []datatypes.User_Preference_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPreferenceTypes", nil, &r.Options, &resp)
	return
}

// Retrieve Data type contains a single user preference to a specific preference type.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetPreferences() (resp []datatypes.User_Preference, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getPreferences", nil, &r.Options, &resp)
	return
}

// Retrieve the authentication requirements for an outstanding password set/reset request.  The requirements returned in the same SoftLayer_Container_User_Customer_PasswordSet container which is provided as a parameter into this request.  The SoftLayer_Container_User_Customer_PasswordSet::authenticationMethods array will contain an entry for each authentication method required for the user.  See SoftLayer_Container_User_Customer_PasswordSet for more details.
//
// If the user has required authentication methods, then authentication information will be supplied to the SoftLayer_User_Customer::processPasswordSetRequest method within this same SoftLayer_Container_User_Customer_PasswordSet container.  All existing information in the container must continue to exist in the container to complete the password set/reset process.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetRequirementsForPasswordSet(passwordSet *datatypes.Container_User_Customer_PasswordSet) (resp datatypes.Container_User_Customer_PasswordSet, err error) {
	params := []interface{}{
		passwordSet,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getRequirementsForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Customer_OpenIdConnect_TrustedProfile) GetRoles() (resp []datatypes.User_Permission_Role, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getRoles", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's security question answers. Some portal users may not have security answers or may not be configured to require answering a security question on login.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSecurityAnswers() (resp []datatypes.User_Customer_Security_Answer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's notification subscription records.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSubscribers() (resp []datatypes.Notification_User_Subscriber, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSubscribers", nil, &r.Options, &resp)
	return
}

// Retrieve A user's successful attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user is required to acknowledge the support policy for portal access.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSupportPolicyAcknowledgementRequiredFlag() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSupportPolicyAcknowledgementRequiredFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSupportPolicyDocument() (resp []byte, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSupportPolicyDocument", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSupportPolicyName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSupportPolicyName", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSupportedLocales() (resp []datatypes.Locale, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSupportedLocales", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a user must take a brief survey the next time they log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSurveyRequiredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSurveyRequiredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The surveys that a user has taken in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetSurveys() (resp []datatypes.Survey, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getSurveys", nil, &r.Options, &resp)
	return
}

// Retrieve An user's associated tickets.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's time zone.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetTimezone() (resp datatypes.Locale_Timezone, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getTimezone", nil, &r.Options, &resp)
	return
}

// Retrieve A user's unsuccessful attempts to log into the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUnsuccessfulLogins() (resp []datatypes.User_Customer_Access_Authentication, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUnsuccessfulLogins", nil, &r.Options, &resp)
	return
}

// Returns an IMS User Object from the provided OpenIdConnect User ID or IBMid Unique Identifier for the Account of the active user. Enforces the User Management permissions for the Active User. An exception will be thrown if no matching IMS User is found. NOTE that providing IBMid Unique Identifier is optional, but it will be preferred over OpenIdConnect User ID if provided.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUserForUnifiedInvitation(openIdConnectUserId *string, uniqueIdentifier *string, searchInvitationsNotLinksFlag *string, accountId *string) (resp datatypes.User_Customer_OpenIdConnect, err error) {
	params := []interface{}{
		openIdConnectUserId,
		uniqueIdentifier,
		searchInvitationsNotLinksFlag,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUserForUnifiedInvitation", params, &r.Options, &resp)
	return
}

// Retrieve a user id using a password token provided to the user in an email generated by the SoftLayer_User_Customer::initiatePortalPasswordChange request. Password recovery keys are valid for 24 hours after they're generated.
//
// When a new user is created or when a user has requested a password change using initiatePortalPasswordChange, they will have received an email that contains a url with a token.  That token is used as the parameter for getUserIdForPasswordSet.  Once the user id is known, then the SoftLayer_User_Customer object can be retrieved which is necessary to complete the process to set or reset a user's password.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUserIdForPasswordSet(key *string) (resp int, err error) {
	params := []interface{}{
		key,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUserIdForPasswordSet", params, &r.Options, &resp)
	return
}

// Retrieve User customer link with IBMid and IAMid.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUserLinks() (resp []datatypes.User_Customer_Link, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUserLinks", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUserPreferences(profileName *string, containerKeyname *string) (resp []datatypes.Layout_Profile, err error) {
	params := []interface{}{
		profileName,
		containerKeyname,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUserPreferences", params, &r.Options, &resp)
	return
}

// Retrieve A portal user's status, which controls overall access to the SoftLayer customer portal and VPN access to the private network.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetUserStatus() (resp datatypes.User_Customer_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getUserStatus", nil, &r.Options, &resp)
	return
}

// Retrieve the number of CloudLayer Computing Instances that a portal user has access to. Portal users can have restrictions set to limit services for and to perform actions on CloudLayer Computing Instances. You can set these permissions in the portal by clicking the "administrative" then "user admin" links.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetVirtualGuestCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getVirtualGuestCount", nil, &r.Options, &resp)
	return
}

// Retrieve A portal user's accessible CloudLayer Computing Instances. These permissions control which CloudLayer Computing Instances a user has access to in the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) GetVirtualGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "getVirtualGuests", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) InTerminalStatus() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "inTerminalStatus", nil, &r.Options, &resp)
	return
}

// Sends password change email to the user containing url that allows the user the change their password. This is the first step when a user wishes to change their password.  The url that is generated contains a one-time use token that is valid for only 24-hours.
//
// If this is a new master user who has never logged into the portal, then password reset will be initiated. Once a master user has logged into the portal, they must setup their security questions prior to logging out because master users are required to answer a security question during the password reset process.  Should a master user not have security questions defined and not remember their password in order to define the security questions, then they will need to contact support at live chat or Revenue Services for assistance.
//
// Due to security reasons, the number of reset requests per username are limited within a undisclosed timeframe.
func (r User_Customer_OpenIdConnect_TrustedProfile) InitiatePortalPasswordChange(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "initiatePortalPasswordChange", params, &r.Options, &resp)
	return
}

// A Brand Agent that has permissions to Add Customer Accounts will be able to request the password email be sent to the Master User of a Customer Account created by the same Brand as the agent making the request. Due to security reasons, the number of reset requests are limited within an undisclosed timeframe.
func (r User_Customer_OpenIdConnect_TrustedProfile) InitiatePortalPasswordChangeByBrandAgent(username *string) (resp bool, err error) {
	params := []interface{}{
		username,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "initiatePortalPasswordChangeByBrandAgent", params, &r.Options, &resp)
	return
}

// Send email invitation to a user to join a SoftLayer account and authenticate with OpenIdConnect. Throws an exception on error.
func (r User_Customer_OpenIdConnect_TrustedProfile) InviteUserToLinkOpenIdConnect(providerType *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "inviteUserToLinkOpenIdConnect", params, &r.Options, &resp)
	return
}

// Portal users are considered master users if they don't have an associated parent user. The only users who don't have parent users are users whose username matches their SoftLayer account name. Master users have special permissions throughout the SoftLayer customer portal.
// Deprecated: This function has been marked as deprecated.
func (r User_Customer_OpenIdConnect_TrustedProfile) IsMasterUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "isMasterUser", nil, &r.Options, &resp)
	return
}

// Determine if a string is the given user's login password to the SoftLayer customer portal.
func (r User_Customer_OpenIdConnect_TrustedProfile) IsValidPortalPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "isValidPortalPassword", params, &r.Options, &resp)
	return
}

// The perform external authentication method will authenticate the given external authentication container with an external vendor.  The authentication container and its contents will be verified before an attempt is made to authenticate the contents of the container with an external vendor.
func (r User_Customer_OpenIdConnect_TrustedProfile) PerformExternalAuthentication(authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "performExternalAuthentication", params, &r.Options, &resp)
	return
}

// Set the password for a user who has an outstanding password request. A user with an outstanding password request will have an unused and unexpired password key.  The password key is part of the url provided to the user in the email sent to the user with information on how to set their password.  The email was generated by the SoftLayer_User_Customer::initiatePortalPasswordRequest request. Password recovery keys are valid for 24 hours after they're generated.
//
// If the user has required authentication methods as specified by in the SoftLayer_Container_User_Customer_PasswordSet container returned from the SoftLayer_User_Customer::getRequirementsForPasswordSet request, then additional requests must be made to processPasswordSetRequest to authenticate the user before changing the password.  First, if the user has security questions set on their profile, they will be required to answer one of their questions correctly. Next, if the user has Verisign or Google Authentication on their account, they must authenticate according to the two-factor provider.  All of this authentication is done using the SoftLayer_Container_User_Customer_PasswordSet container.
//
// User portal passwords must match the following restrictions. Portal passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ + =
// * ...not match your username
func (r User_Customer_OpenIdConnect_TrustedProfile) ProcessPasswordSetRequest(passwordSet *datatypes.Container_User_Customer_PasswordSet, authenticationContainer *datatypes.Container_User_Customer_External_Binding) (resp bool, err error) {
	params := []interface{}{
		passwordSet,
		authenticationContainer,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "processPasswordSetRequest", params, &r.Options, &resp)
	return
}

// Revoke access to all dedicated hosts on the account for this user. The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveAllDedicatedHostAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeAllDedicatedHostAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveAllHardwareAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeAllHardwareAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove all cloud computing instances from a portal user's instance access list. A user's instance access list controls which of an account's computing instance objects a user has access to in the SoftLayer customer portal and API. If the current user does not have administrative privileges over this user, an inadequate permissions exception will get thrown.
//
// Users can call this function on child users, but not to themselves. An account's master has access to all users permissions on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveAllVirtualAccessForThisUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeAllVirtualAccessForThisUser", nil, &r.Options, &resp)
	return
}

// Remove a user's API authentication key, removing that user's access to query the SoftLayer API.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveApiAuthenticationKey(keyId *int) (resp bool, err error) {
	params := []interface{}{
		keyId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeApiAuthenticationKey", params, &r.Options, &resp)
	return
}

// Revokes access for the user to one or more dedicated host devices.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
//
// If the user has full dedicatedHost access, then it will provide access to "ALL but passed in" dedicatedHost ids.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveBulkDedicatedHostAccess(dedicatedHostIds []int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeBulkDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// Remove multiple hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting to remove then removeBulkHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
//
// If the user has full hardware access, then it will provide access to "ALL but passed in" hardware ids.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveBulkHardwareAccess(hardwareIds []int) (resp bool, err error) {
	params := []interface{}{
		hardwareIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeBulkHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) multiple permissions from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. removePortalPermission() does not attempt to remove permissions that are not assigned to the user.
//
// Users can grant or revoke permissions to their child users, but not to themselves. An account's master has all portal permissions and can grant permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permissions from a user will cascade down the child hierarchy and remove the permissions from this user along with all child users who also have the permission.
//
// If the cascadePermissionsFlag is not provided or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission objects within the permissions parameter.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveBulkPortalPermission(permissions []datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permissions,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeBulkPortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveBulkRoles(roles []datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		roles,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeBulkRoles", params, &r.Options, &resp)
	return
}

// Remove multiple CloudLayer Computing Instances from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's CloudLayer Computing Instance objects a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeBulkVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveBulkVirtualGuestAccess(virtualGuestIds []int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestIds,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeBulkVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// Revokes access for the user to a single dedicated host device.  The user will only be allowed to see and access devices in both the portal and the API to which they have been granted access.  If the user's account has devices to which the user has not been granted access or the access has been revoked, then "not found" exceptions are thrown if the user attempts to access any of these devices.
//
// Users can assign device access to their child users, but not to themselves. An account's master has access to all devices on their customer account and can set dedicated host access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveDedicatedHostAccess(dedicatedHostId *int) (resp bool, err error) {
	params := []interface{}{
		dedicatedHostId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeDedicatedHostAccess", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveExternalBinding(externalBinding *datatypes.User_External_Binding) (resp bool, err error) {
	params := []interface{}{
		externalBinding,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeExternalBinding", params, &r.Options, &resp)
	return
}

// Remove hardware from a portal user's hardware access list. A user's hardware access list controls which of an account's hardware objects a user has access to in the SoftLayer customer portal and API. Hardware does not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the hardware you're attempting remove add then removeHardwareAccess() returns true.
//
// Users can assign hardware access to their child users, but not to themselves. An account's master has access to all hardware on their customer account and can set hardware access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveHardwareAccess(hardwareId *int) (resp bool, err error) {
	params := []interface{}{
		hardwareId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeHardwareAccess", params, &r.Options, &resp)
	return
}

// Remove (revoke) a permission from a portal user's permission set. [[SoftLayer_User_Customer_CustomerPermission_Permission]] control which features in the SoftLayer customer portal and API a user may use. Removing a user's permission will affect that user's portal and API access. If the user does not have the permission you're attempting to remove then removePortalPermission() returns true.
//
// Users can assign permissions to their child users, but not to themselves. An account's master has all portal permissions and can set permissions for any of the other users on their account.
//
// If the cascadePermissionsFlag is set to true, then removing the permission from a user will cascade down the child hierarchy and remove the permission from this user and all child users who also have the permission.
//
// If the cascadePermissionsFlag is not set or is set to false and the user has children users who have the permission, then an exception will be thrown, and the permission will not be removed from this user.
//
// Use the [[SoftLayer_User_Customer_CustomerPermission_Permission::getAllObjects]] method to retrieve a list of all permissions available in the SoftLayer customer portal and API. Permissions are removed based on the keyName property of the permission parameter.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemovePortalPermission(permission *datatypes.User_Customer_CustomerPermission_Permission, cascadePermissionsFlag *bool) (resp bool, err error) {
	params := []interface{}{
		permission,
		cascadePermissionsFlag,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removePortalPermission", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeRole", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveSecurityAnswers() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeSecurityAnswers", nil, &r.Options, &resp)
	return
}

// Remove a CloudLayer Computing Instance from a portal user's access list. A user's CloudLayer Computing Instance access list controls which of an account's computing instances a user has access to in the SoftLayer customer portal and API. CloudLayer Computing Instances do not exist in the SoftLayer portal and returns "not found" exceptions in the API if the user doesn't have access to it. If a user does not has access to the CloudLayer Computing Instance you're attempting remove add then removeVirtualGuestAccess() returns true.
//
// Users can assign CloudLayer Computing Instance access to their child users, but not to themselves. An account's master has access to all CloudLayer Computing Instances on their customer account and can set instance access for any of the other users on their account.
func (r User_Customer_OpenIdConnect_TrustedProfile) RemoveVirtualGuestAccess(virtualGuestId *int) (resp bool, err error) {
	params := []interface{}{
		virtualGuestId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "removeVirtualGuestAccess", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer user is linked to, if we need to do that for some reason. It will do this by modifying the link to the desired new IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is already Bluemix linked. To reset a link for the Bluemix-linked user account, use resetOpenIdConnectLinkUnifiedUserManagementMode.
func (r User_Customer_OpenIdConnect_TrustedProfile) ResetOpenIdConnectLink(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "resetOpenIdConnectLink", params, &r.Options, &resp)
	return
}

// This method will change the IBMid that a SoftLayer master user is linked to, if we need to do that for some reason. It will do this by unlinking the new owner IBMid from its current user association in this account, if there is one (note that the new owner IBMid is not required to already be a member of the IMS account). Then it will modify the existing IBMid link for the master user to use the new owner IBMid-realm IAMid. At this point, if the new owner IBMid isn't already a member of the PaaS account, it will attempt to add it. As a last step, it will call PaaS to modify the owner on that side, if necessary.  Only when all those steps are complete, it will commit the IMS-side DB changes.  Then, it will clean up the SoftLayer user that was linked to the new owner IBMid (this user became unlinked as the first step in this process).  It will also call BSS to delete the old owner IBMid. NOTE:  This method cannot be used to "un-link" a SoftLayer user.  Once linked, a SoftLayer user can never be un-linked. Also, this method cannot be used to reset the link if the user account is not Bluemix linked. To reset a link for the user account not linked to Bluemix, use resetOpenIdConnectLink.
func (r User_Customer_OpenIdConnect_TrustedProfile) ResetOpenIdConnectLinkUnifiedUserManagementMode(providerType *string, newIbmIdUsername *string, removeSecuritySettings *bool) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		providerType,
		newIbmIdUsername,
		removeSecuritySettings,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "resetOpenIdConnectLinkUnifiedUserManagementMode", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) SamlAuthenticate(accountId *string, samlResponse *string) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		accountId,
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "samlAuthenticate", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) SamlBeginAuthentication(accountId *int) (resp string, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "samlBeginAuthentication", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) SamlBeginLogout() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "samlBeginLogout", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) SamlLogout(samlResponse *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		samlResponse,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "samlLogout", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_OpenIdConnect_TrustedProfile) SelfPasswordChange(currentPassword *string, newPassword *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		currentPassword,
		newPassword,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "selfPasswordChange", params, &r.Options, &resp)
	return
}

// An OpenIdConnect identity, for example an IAMid, can be linked or mapped to one or more individual SoftLayer users, but no more than one per account. If an OpenIdConnect identity is mapped to multiple accounts in this manner, one such account should be identified as the default account for that identity. Invoke this only on IBMid-authenticated users.
func (r User_Customer_OpenIdConnect_TrustedProfile) SetDefaultAccount(providerType *string, accountId *int) (resp datatypes.Account, err error) {
	params := []interface{}{
		providerType,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "setDefaultAccount", params, &r.Options, &resp)
	return
}

// As master user, calling this api for the IBMid provider type when there is an existing IBMid for the email on the SL account will silently (without sending an invitation email) create a link for the IBMid. NOTE: If the SoftLayer user is already linked to IBMid, this call will fail. If the IBMid specified by the email of this user, is already used in a link to another user in this account, this call will fail. If there is already an open invitation from this SoftLayer user to this or any IBMid, this call will fail. If there is already an open invitation from some other SoftLayer user in this account to this IBMid, then this call will fail.
func (r User_Customer_OpenIdConnect_TrustedProfile) SilentlyMigrateUserOpenIdConnect(providerType *string) (resp bool, err error) {
	params := []interface{}{
		providerType,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "silentlyMigrateUserOpenIdConnect", params, &r.Options, &resp)
	return
}

// This method allows the master user of an account to undo the designation of this user as an alternate master user.  This can not be applied to the true master user of the account.
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted this user.  This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner turns off an "alternate/auxiliary master user / account owner".
func (r User_Customer_OpenIdConnect_TrustedProfile) TurnOffMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "turnOffMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// This method allows the master user of an account to designate this user as an alternate master user.  Effectively this means that this user should have "all the same IMS permissions as a master user".
//
// Note that this method, by itself, WILL NOT affect the IAM Policies granted to this user. This API is not intended for general customer use.  It is intended to be called by IAM, in concert with other actions taken by IAM when the master user / account owner designates an "alternate/auxiliary master user / account owner".
func (r User_Customer_OpenIdConnect_TrustedProfile) TurnOnMasterUserPermissionCheckMode() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "turnOnMasterUserPermissionCheckMode", nil, &r.Options, &resp)
	return
}

// Update the active status for a notification that the user is subscribed to. A notification along with an active flag can be supplied to update the active status for a particular notification subscription.
func (r User_Customer_OpenIdConnect_TrustedProfile) UpdateNotificationSubscriber(notificationKeyName *string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "updateNotificationSubscriber", params, &r.Options, &resp)
	return
}

// Update a user's login security questions and answers on the SoftLayer customer portal. These questions and answers are used to optionally log into the SoftLayer customer portal using two-factor authentication. Each user must have three distinct questions set with a unique answer for each question, and each answer may only contain alphanumeric or the . , - _ ( ) [ ] : ; > < characters. Existing user security questions and answers are deleted before new ones are set, and users may only update their own security questions and answers.
func (r User_Customer_OpenIdConnect_TrustedProfile) UpdateSecurityAnswers(questions []datatypes.User_Security_Question, answers []string) (resp bool, err error) {
	params := []interface{}{
		questions,
		answers,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "updateSecurityAnswers", params, &r.Options, &resp)
	return
}

// Update a delivery method for a notification that the user is subscribed to. A delivery method keyName along with an active flag can be supplied to update the active status of the delivery methods for the specified notification. Available delivery methods - 'EMAIL'. Available notifications - 'PLANNED_MAINTENANCE', 'UNPLANNED_INCIDENT'.
func (r User_Customer_OpenIdConnect_TrustedProfile) UpdateSubscriberDeliveryMethod(notificationKeyName *string, deliveryMethodKeyNames []string, active *int) (resp bool, err error) {
	params := []interface{}{
		notificationKeyName,
		deliveryMethodKeyNames,
		active,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "updateSubscriberDeliveryMethod", params, &r.Options, &resp)
	return
}

// Update a user's VPN password on the SoftLayer customer portal. As with portal passwords, VPN passwords must match the following restrictions. VPN passwords must...
// * ...be over eight characters long.
// * ...be under twenty characters long.
// * ...contain at least one uppercase letter
// * ...contain at least one lowercase letter
// * ...contain at least one number
// * ...contain one of the special characters _ - | @ . , ? / ! ~ # $ % ^ & * ( ) { } [ ] \ =
// * ...not match your username
// Finally, users can only update their own VPN password. An account's master user can update any of their account users' VPN passwords.
func (r User_Customer_OpenIdConnect_TrustedProfile) UpdateVpnPassword(password *string) (resp bool, err error) {
	params := []interface{}{
		password,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "updateVpnPassword", params, &r.Options, &resp)
	return
}

// Always call this function to enable changes when manually configuring VPN subnet access.
func (r User_Customer_OpenIdConnect_TrustedProfile) UpdateVpnUser() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "updateVpnUser", nil, &r.Options, &resp)
	return
}

// This method validate the given authentication token using the user id by comparing it with the actual user authentication token and return [[SoftLayer_Container_User_Customer_Portal_Token]] object
func (r User_Customer_OpenIdConnect_TrustedProfile) ValidateAuthenticationToken(authenticationToken *datatypes.Container_User_Authentication_Token) (resp datatypes.Container_User_Customer_Portal_Token, err error) {
	params := []interface{}{
		authenticationToken,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_OpenIdConnect_TrustedProfile", "validateAuthenticationToken", params, &r.Options, &resp)
	return
}

// no documentation yet
type User_Customer_Profile_Event_HyperWarp struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerProfileEventHyperWarpService returns an instance of the User_Customer_Profile_Event_HyperWarp SoftLayer service
func GetUserCustomerProfileEventHyperWarpService(sess session.SLSession) User_Customer_Profile_Event_HyperWarp {
	return User_Customer_Profile_Event_HyperWarp{Session: sess}
}

func (r User_Customer_Profile_Event_HyperWarp) Id(id int) User_Customer_Profile_Event_HyperWarp {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Profile_Event_HyperWarp) Mask(mask string) User_Customer_Profile_Event_HyperWarp {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Profile_Event_HyperWarp) Filter(filter string) User_Customer_Profile_Event_HyperWarp {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Profile_Event_HyperWarp) Limit(limit int) User_Customer_Profile_Event_HyperWarp {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Profile_Event_HyperWarp) Offset(offset int) User_Customer_Profile_Event_HyperWarp {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Customer_Profile_Event_HyperWarp) ReceiveEventDirect(eventJson *datatypes.Container_User_Customer_Profile_Event_HyperWarp_ProfileChange) (resp bool, err error) {
	params := []interface{}{
		eventJson,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Profile_Event_HyperWarp", "receiveEventDirect", params, &r.Options, &resp)
	return
}

// Contains user information for Service Provider Enrollment.
type User_Customer_Prospect_ServiceProvider_EnrollRequest struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerProspectServiceProviderEnrollRequestService returns an instance of the User_Customer_Prospect_ServiceProvider_EnrollRequest SoftLayer service
func GetUserCustomerProspectServiceProviderEnrollRequestService(sess session.SLSession) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	return User_Customer_Prospect_ServiceProvider_EnrollRequest{Session: sess}
}

func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Id(id int) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Mask(mask string) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Filter(filter string) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Limit(limit int) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Offset(offset int) User_Customer_Prospect_ServiceProvider_EnrollRequest {
	r.Options.Offset = &offset
	return r
}

// Create a new Service Provider Enrollment
func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) Enroll(templateObject *datatypes.User_Customer_Prospect_ServiceProvider_EnrollRequest) (resp datatypes.User_Customer_Prospect_ServiceProvider_EnrollRequest, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Customer_Prospect_ServiceProvider_EnrollRequest", "enroll", params, &r.Options, &resp)
	return
}

// Retrieve Catalyst company types.
func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) GetCompanyType() (resp datatypes.Catalyst_Company_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Prospect_ServiceProvider_EnrollRequest", "getCompanyType", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Customer_Prospect_ServiceProvider_EnrollRequest) GetObject() (resp datatypes.User_Customer_Prospect_ServiceProvider_EnrollRequest, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Prospect_ServiceProvider_EnrollRequest", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Customer_Security_Answer type contains user's answers to security questions.
type User_Customer_Security_Answer struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerSecurityAnswerService returns an instance of the User_Customer_Security_Answer SoftLayer service
func GetUserCustomerSecurityAnswerService(sess session.SLSession) User_Customer_Security_Answer {
	return User_Customer_Security_Answer{Session: sess}
}

func (r User_Customer_Security_Answer) Id(id int) User_Customer_Security_Answer {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Security_Answer) Mask(mask string) User_Customer_Security_Answer {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Security_Answer) Filter(filter string) User_Customer_Security_Answer {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Security_Answer) Limit(limit int) User_Customer_Security_Answer {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Security_Answer) Offset(offset int) User_Customer_Security_Answer {
	r.Options.Offset = &offset
	return r
}

// getObject retrieves the SoftLayer_User_Customer_Security_Answer object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_Security_Answer service.
func (r User_Customer_Security_Answer) GetObject() (resp datatypes.User_Customer_Security_Answer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Security_Answer", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The question the security answer is associated with.
func (r User_Customer_Security_Answer) GetQuestion() (resp datatypes.User_Security_Question, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Security_Answer", "getQuestion", nil, &r.Options, &resp)
	return
}

// Retrieve The user who the security answer belongs to.
func (r User_Customer_Security_Answer) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Security_Answer", "getUser", nil, &r.Options, &resp)
	return
}

// Each SoftLayer User Customer instance is assigned a status code that determines how it's treated in the customer portal. This status is reflected in the SoftLayer_User_Customer_Status data type. Status differs from user permissions in that user status applies globally to the portal while user permissions are applied to specific portal functions.
//
// Note that a status of "PENDING" also has been added. This status is specific to users that are configured to use IBMid authentication. This would include some (not all) users on accounts that are linked to Platform Services (PaaS, formerly Bluemix) accounts, but is not limited to users in such accounts. Using IBMid authentication is optional for active users even if it is not required by the account type. PENDING status indicates that a relationship between an IBMid and a user is being set up but is not complete. To be complete, PENDING users need to perform an action ("accepting the invitation") before becoming an active user within IBM Cloud and/or IMS. PENDING is a system state, and can not be administered by users (including the account master user). SoftLayer Commercial is the only environment where IBMid and/or account linking are used.
type User_Customer_Status struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserCustomerStatusService returns an instance of the User_Customer_Status SoftLayer service
func GetUserCustomerStatusService(sess session.SLSession) User_Customer_Status {
	return User_Customer_Status{Session: sess}
}

func (r User_Customer_Status) Id(id int) User_Customer_Status {
	r.Options.Id = &id
	return r
}

func (r User_Customer_Status) Mask(mask string) User_Customer_Status {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Customer_Status) Filter(filter string) User_Customer_Status {
	r.Options.Filter = filter
	return r
}

func (r User_Customer_Status) Limit(limit int) User_Customer_Status {
	r.Options.Limit = &limit
	return r
}

func (r User_Customer_Status) Offset(offset int) User_Customer_Status {
	r.Options.Offset = &offset
	return r
}

// Retrieve all user status objects.
func (r User_Customer_Status) GetAllObjects() (resp []datatypes.User_Customer_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Status", "getAllObjects", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_User_Customer_Status object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_User_Customer_Status service.
func (r User_Customer_Status) GetObject() (resp datatypes.User_Customer_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Customer_Status", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_External_Binding data type contains general information for a single external binding.  This includes the 3rd party vendor, type of binding, and a unique identifier and password that is used to authenticate against the 3rd party service.
type User_External_Binding struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserExternalBindingService returns an instance of the User_External_Binding SoftLayer service
func GetUserExternalBindingService(sess session.SLSession) User_External_Binding {
	return User_External_Binding{Session: sess}
}

func (r User_External_Binding) Id(id int) User_External_Binding {
	r.Options.Id = &id
	return r
}

func (r User_External_Binding) Mask(mask string) User_External_Binding {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_External_Binding) Filter(filter string) User_External_Binding {
	r.Options.Filter = filter
	return r
}

func (r User_External_Binding) Limit(limit int) User_External_Binding {
	r.Options.Limit = &limit
	return r
}

func (r User_External_Binding) Offset(offset int) User_External_Binding {
	r.Options.Offset = &offset
	return r
}

// Delete an external authentication binding.  If the external binding currently has an active billing item associated you will be prevented from deleting the binding.  The alternative method to remove an external authentication binding is to use the service cancellation form.
func (r User_External_Binding) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "deleteObject", nil, &r.Options, &resp)
	return
}

// Retrieve Attributes of an external authentication binding.
func (r User_External_Binding) GetAttributes() (resp []datatypes.User_External_Binding_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Information regarding the billing item for external authentication.
func (r User_External_Binding) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve An optional note for identifying the external binding.
func (r User_External_Binding) GetNote() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getNote", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_External_Binding) GetObject() (resp datatypes.User_External_Binding, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The type of external authentication binding.
func (r User_External_Binding) GetType() (resp datatypes.User_External_Binding_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getType", nil, &r.Options, &resp)
	return
}

// Retrieve The vendor of an external authentication binding.
func (r User_External_Binding) GetVendor() (resp datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "getVendor", nil, &r.Options, &resp)
	return
}

// Update the note of an external binding.  The note is an optional property that is used to store information about a binding.
func (r User_External_Binding) UpdateNote(text *string) (resp bool, err error) {
	params := []interface{}{
		text,
	}
	err = r.Session.DoRequest("SoftLayer_User_External_Binding", "updateNote", params, &r.Options, &resp)
	return
}

// The SoftLayer_User_External_Binding_Vendor data type contains information for a single external binding vendor.  This information includes a user friendly vendor name, a unique version of the vendor name, and a unique internal identifier that can be used when creating a new external binding.
type User_External_Binding_Vendor struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserExternalBindingVendorService returns an instance of the User_External_Binding_Vendor SoftLayer service
func GetUserExternalBindingVendorService(sess session.SLSession) User_External_Binding_Vendor {
	return User_External_Binding_Vendor{Session: sess}
}

func (r User_External_Binding_Vendor) Id(id int) User_External_Binding_Vendor {
	r.Options.Id = &id
	return r
}

func (r User_External_Binding_Vendor) Mask(mask string) User_External_Binding_Vendor {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_External_Binding_Vendor) Filter(filter string) User_External_Binding_Vendor {
	r.Options.Filter = filter
	return r
}

func (r User_External_Binding_Vendor) Limit(limit int) User_External_Binding_Vendor {
	r.Options.Limit = &limit
	return r
}

func (r User_External_Binding_Vendor) Offset(offset int) User_External_Binding_Vendor {
	r.Options.Offset = &offset
	return r
}

// getAllObjects() will return a list of the available external binding vendors that SoftLayer supports.  Use this list to select the appropriate vendor when creating a new external binding.
func (r User_External_Binding_Vendor) GetAllObjects() (resp []datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding_Vendor", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_External_Binding_Vendor) GetObject() (resp datatypes.User_External_Binding_Vendor, err error) {
	err = r.Session.DoRequest("SoftLayer_User_External_Binding_Vendor", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Permission_Action data type contains local attributes to identify and describe the valid actions a customer user can perform within IMS.  This includes a name, key name, and description.  This data can not be modified by users of IMS.
//
// It also contains relational attributes that indicate which SoftLayer_User_Permission_Group's include the action.
type User_Permission_Action struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionActionService returns an instance of the User_Permission_Action SoftLayer service
func GetUserPermissionActionService(sess session.SLSession) User_Permission_Action {
	return User_Permission_Action{Session: sess}
}

func (r User_Permission_Action) Id(id int) User_Permission_Action {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Action) Mask(mask string) User_Permission_Action {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Action) Filter(filter string) User_Permission_Action {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Action) Limit(limit int) User_Permission_Action {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Action) Offset(offset int) User_Permission_Action {
	r.Options.Offset = &offset
	return r
}

// Object filters and result limits are enabled on this method.
func (r User_Permission_Action) GetAllObjects() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Action", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Action) GetDepartment() (resp datatypes.User_Permission_Department, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Action", "getDepartment", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Action) GetObject() (resp datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Action", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
type User_Permission_Department struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionDepartmentService returns an instance of the User_Permission_Department SoftLayer service
func GetUserPermissionDepartmentService(sess session.SLSession) User_Permission_Department {
	return User_Permission_Department{Session: sess}
}

func (r User_Permission_Department) Id(id int) User_Permission_Department {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Department) Mask(mask string) User_Permission_Department {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Department) Filter(filter string) User_Permission_Department {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Department) Limit(limit int) User_Permission_Department {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Department) Offset(offset int) User_Permission_Department {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r User_Permission_Department) GetAllObjects() (resp []datatypes.User_Permission_Department, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Department", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Department) GetObject() (resp datatypes.User_Permission_Department, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Department", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Department) GetPermissions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Department", "getPermissions", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Permission_Group data type contains local attributes to identify and describe the permission groups that have been created within IMS.  These includes a name, description, and account id.  Permission groups are defined specifically for a single [[SoftLayer_Account]].
//
// It also contains relational attributes that indicate what SoftLayer_User_Permission_Action objects belong to a particular group, and what SoftLayer_User_Permission_Role objects the group is linked.
type User_Permission_Group struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionGroupService returns an instance of the User_Permission_Group SoftLayer service
func GetUserPermissionGroupService(sess session.SLSession) User_Permission_Group {
	return User_Permission_Group{Session: sess}
}

func (r User_Permission_Group) Id(id int) User_Permission_Group {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Group) Mask(mask string) User_Permission_Group {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Group) Filter(filter string) User_Permission_Group {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Group) Limit(limit int) User_Permission_Group {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Group) Offset(offset int) User_Permission_Group {
	r.Options.Offset = &offset
	return r
}

// Assigns a SoftLayer_User_Permission_Action object to the group.
func (r User_Permission_Group) AddAction(action *datatypes.User_Permission_Action) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		action,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "addAction", params, &r.Options, &resp)
	return
}

// Assigns multiple SoftLayer_User_Permission_Action objects to the group.
func (r User_Permission_Group) AddBulkActions(actions []datatypes.User_Permission_Action) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		actions,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "addBulkActions", params, &r.Options, &resp)
	return
}

// Links multiple SoftLayer_Hardware_Server, SoftLayer_Virtual_Guest, or SoftLayer_Virtual_DedicatedHost objects to the group. All objects must be of the same type.
func (r User_Permission_Group) AddBulkResourceObjects(resourceObjects []datatypes.Entity, resourceTypeKeyName *string) (resp bool, err error) {
	params := []interface{}{
		resourceObjects,
		resourceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "addBulkResourceObjects", params, &r.Options, &resp)
	return
}

// Links a SoftLayer_Hardware_Server, SoftLayer_Virtual_Guest, or SoftLayer_Virtual_DedicatedHost object to the group.
func (r User_Permission_Group) AddResourceObject(resourceObject *datatypes.Entity, resourceTypeKeyName *string) (resp bool, err error) {
	params := []interface{}{
		resourceObject,
		resourceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "addResourceObject", params, &r.Options, &resp)
	return
}

// Customer created permission groups must be of type NORMAL.  The SYSTEM type is reserved for internal use. The account id supplied in the template permission group must match account id of the user who is creating the permission group.  The user who is creating the permission group must have the permission to manage users.
func (r User_Permission_Group) CreateObject(templateObject *datatypes.User_Permission_Group) (resp datatypes.User_Permission_Group, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "createObject", params, &r.Options, &resp)
	return
}

// Customer users can only delete permission groups of type NORMAL.  The SYSTEM type is reserved for internal use. The user who is creating the permission group must have the permission to manage users.
func (r User_Permission_Group) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "deleteObject", nil, &r.Options, &resp)
	return
}

// Allows a user to modify the name and description of an existing customer permission group. Customer permission groups must be of type NORMAL.  The SYSTEM type is reserved for internal use. The account id supplied in the template permission group must match account id of the user who is creating the permission group.  The user who is creating the permission group must have the permission to manage users.
func (r User_Permission_Group) EditObject(templateObject *datatypes.User_Permission_Group) (resp datatypes.User_Permission_Group, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Group) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Group) GetActions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "getActions", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Group) GetObject() (resp datatypes.User_Permission_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Group) GetRoles() (resp []datatypes.User_Permission_Role, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "getRoles", nil, &r.Options, &resp)
	return
}

// Retrieve The type of the permission group.
func (r User_Permission_Group) GetType() (resp datatypes.User_Permission_Group_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "getType", nil, &r.Options, &resp)
	return
}

// Links a SoftLayer_User_Permission_Role object to the group.
func (r User_Permission_Group) LinkRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "linkRole", params, &r.Options, &resp)
	return
}

// Unassigns a SoftLayer_User_Permission_Action object from the group.
func (r User_Permission_Group) RemoveAction(action *datatypes.User_Permission_Action) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		action,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "removeAction", params, &r.Options, &resp)
	return
}

// Unassigns multiple SoftLayer_User_Permission_Action objects from the group.
func (r User_Permission_Group) RemoveBulkActions(actions []datatypes.User_Permission_Action) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		actions,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "removeBulkActions", params, &r.Options, &resp)
	return
}

// Unlinks multiple SoftLayer_Hardware_Server, SoftLayer_Virtual_Guest, or SoftLayer_Virtual_DedicatedHost objects from the group. All objects must be of the same type.
func (r User_Permission_Group) RemoveBulkResourceObjects(resourceObjects []datatypes.Entity, resourceTypeKeyName *string) (resp bool, err error) {
	params := []interface{}{
		resourceObjects,
		resourceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "removeBulkResourceObjects", params, &r.Options, &resp)
	return
}

// Unlinks a SoftLayer_Hardware_Server, SoftLayer_Virtual_Guest, or SoftLayer_Virtual_DedicatedHost object from the group.
func (r User_Permission_Group) RemoveResourceObject(resourceObject *datatypes.Entity, resourceTypeKeyName *string) (resp bool, err error) {
	params := []interface{}{
		resourceObject,
		resourceTypeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "removeResourceObject", params, &r.Options, &resp)
	return
}

// Removes a link from SoftLayer_User_Permission_Role object to the group.
func (r User_Permission_Group) UnlinkRole(role *datatypes.User_Permission_Role) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		role,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group", "unlinkRole", params, &r.Options, &resp)
	return
}

// These are the attributes which describe a SoftLayer_User_Permission_Group_Type. All SoftLayer_User_Permission_Group objects must be linked to one of these types.
//
// For further information see: [[SoftLayer_User_Permission_Group]].
type User_Permission_Group_Type struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionGroupTypeService returns an instance of the User_Permission_Group_Type SoftLayer service
func GetUserPermissionGroupTypeService(sess session.SLSession) User_Permission_Group_Type {
	return User_Permission_Group_Type{Session: sess}
}

func (r User_Permission_Group_Type) Id(id int) User_Permission_Group_Type {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Group_Type) Mask(mask string) User_Permission_Group_Type {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Group_Type) Filter(filter string) User_Permission_Group_Type {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Group_Type) Limit(limit int) User_Permission_Group_Type {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Group_Type) Offset(offset int) User_Permission_Group_Type {
	r.Options.Offset = &offset
	return r
}

// Retrieve The groups that are of this type.
func (r User_Permission_Group_Type) GetGroups() (resp []datatypes.User_Permission_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group_Type", "getGroups", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Group_Type) GetObject() (resp datatypes.User_Permission_Group_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Group_Type", "getObject", nil, &r.Options, &resp)
	return
}

// These are the variables relating to SoftLayer_User_Permission_Resource_Type. Collectively they describe the types of resources which can be linked to [[SoftLayer_User_Permission_Group]].
type User_Permission_Resource_Type struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionResourceTypeService returns an instance of the User_Permission_Resource_Type SoftLayer service
func GetUserPermissionResourceTypeService(sess session.SLSession) User_Permission_Resource_Type {
	return User_Permission_Resource_Type{Session: sess}
}

func (r User_Permission_Resource_Type) Id(id int) User_Permission_Resource_Type {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Resource_Type) Mask(mask string) User_Permission_Resource_Type {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Resource_Type) Filter(filter string) User_Permission_Resource_Type {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Resource_Type) Limit(limit int) User_Permission_Resource_Type {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Resource_Type) Offset(offset int) User_Permission_Resource_Type {
	r.Options.Offset = &offset
	return r
}

// Retrieve an array of SoftLayer_User_Permission_Resource_Type objects.
func (r User_Permission_Resource_Type) GetAllObjects() (resp []datatypes.User_Permission_Resource_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Resource_Type", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Resource_Type) GetObject() (resp datatypes.User_Permission_Resource_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Resource_Type", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_User_Permission_Role data type contains local attributes to identify and describe the permission roles that have been created within IMS.  These includes a name, description, and account id.  Permission groups are defined specifically for a single [[SoftLayer_Account]].
//
// It also contains relational attributes that indicate what SoftLayer_User_Permission_Group objects are linked to a particular role, and the SoftLayer_User_Customer objects assigned to the role.
type User_Permission_Role struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserPermissionRoleService returns an instance of the User_Permission_Role SoftLayer service
func GetUserPermissionRoleService(sess session.SLSession) User_Permission_Role {
	return User_Permission_Role{Session: sess}
}

func (r User_Permission_Role) Id(id int) User_Permission_Role {
	r.Options.Id = &id
	return r
}

func (r User_Permission_Role) Mask(mask string) User_Permission_Role {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Permission_Role) Filter(filter string) User_Permission_Role {
	r.Options.Filter = filter
	return r
}

func (r User_Permission_Role) Limit(limit int) User_Permission_Role {
	r.Options.Limit = &limit
	return r
}

func (r User_Permission_Role) Offset(offset int) User_Permission_Role {
	r.Options.Offset = &offset
	return r
}

// Assigns a SoftLayer_User_Customer object to the role.
func (r User_Permission_Role) AddUser(user *datatypes.User_Customer) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		user,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "addUser", params, &r.Options, &resp)
	return
}

// Customer created permission roles must set the systemFlag attribute to false.  The SYSTEM type is reserved for internal use. The account id supplied in the template permission group must match account id of the user who is creating the permission group.  The user who is creating the permission group must have the permission to manage users.
func (r User_Permission_Role) CreateObject(templateObject *datatypes.User_Permission_Role) (resp datatypes.User_Permission_Role, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "createObject", params, &r.Options, &resp)
	return
}

// Customer users can only delete permission roles with systemFlag set to false.  The SYSTEM type is reserved for internal use. The user who is creating the permission role must have the permission to manage users.
func (r User_Permission_Role) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "deleteObject", nil, &r.Options, &resp)
	return
}

// Allows a user to modify the name and description of an existing customer permission role. Customer permission roles must set the systemFlag attribute to false.  The SYSTEM type is reserved for internal use. The account id supplied in the template permission role must match account id of the user who is creating the permission role.  The user who is creating the permission role must have the permission to manage users.
func (r User_Permission_Role) EditObject(templateObject *datatypes.User_Permission_Role) (resp datatypes.User_Permission_Role, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Role) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Role) GetActions() (resp []datatypes.User_Permission_Action, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "getActions", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Role) GetGroups() (resp []datatypes.User_Permission_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "getGroups", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r User_Permission_Role) GetObject() (resp datatypes.User_Permission_Role, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r User_Permission_Role) GetUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "getUsers", nil, &r.Options, &resp)
	return
}

// Links a SoftLayer_User_Permission_Group object to the role.
func (r User_Permission_Role) LinkGroup(group *datatypes.User_Permission_Group) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		group,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "linkGroup", params, &r.Options, &resp)
	return
}

// Unassigns a SoftLayer_User_Customer object from the role.
func (r User_Permission_Role) RemoveUser(user *datatypes.User_Customer) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		user,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "removeUser", params, &r.Options, &resp)
	return
}

// Unlinks a SoftLayer_User_Permission_Group object to the role.
func (r User_Permission_Role) UnlinkGroup(group *datatypes.User_Permission_Group) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		group,
	}
	err = r.Session.DoRequest("SoftLayer_User_Permission_Role", "unlinkGroup", params, &r.Options, &resp)
	return
}

// The SoftLayer_User_Security_Question data type contains questions.
type User_Security_Question struct {
	Session session.SLSession
	Options sl.Options
}

// GetUserSecurityQuestionService returns an instance of the User_Security_Question SoftLayer service
func GetUserSecurityQuestionService(sess session.SLSession) User_Security_Question {
	return User_Security_Question{Session: sess}
}

func (r User_Security_Question) Id(id int) User_Security_Question {
	r.Options.Id = &id
	return r
}

func (r User_Security_Question) Mask(mask string) User_Security_Question {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r User_Security_Question) Filter(filter string) User_Security_Question {
	r.Options.Filter = filter
	return r
}

func (r User_Security_Question) Limit(limit int) User_Security_Question {
	r.Options.Limit = &limit
	return r
}

func (r User_Security_Question) Offset(offset int) User_Security_Question {
	r.Options.Offset = &offset
	return r
}

// Retrieve all viewable security questions.
func (r User_Security_Question) GetAllObjects() (resp []datatypes.User_Security_Question, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Security_Question", "getAllObjects", nil, &r.Options, &resp)
	return
}

// getAllObjects retrieves all the SoftLayer_User_Security_Question objects where it is set to be viewable.
func (r User_Security_Question) GetObject() (resp datatypes.User_Security_Question, err error) {
	err = r.Session.DoRequest("SoftLayer_User_Security_Question", "getObject", nil, &r.Options, &resp)
	return
}

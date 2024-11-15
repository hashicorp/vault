package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// User represents a Microsoft Entra user account.
type User struct {
    DirectoryObject
}
// NewUser instantiates a new User and sets the default values.
func NewUser()(*User) {
    m := &User{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.user"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUser(), nil
}
// GetAboutMe gets the aboutMe property value. A freeform text entry field for the user to describe themselves. Returned only on $select.
// returns a *string when successful
func (m *User) GetAboutMe()(*string) {
    val, err := m.GetBackingStore().Get("aboutMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAccountEnabled gets the accountEnabled property value. true if the account is enabled; otherwise, false. This property is required when a user is created. Returned only on $select. Supports $filter (eq, ne, not, and in).
// returns a *bool when successful
func (m *User) GetAccountEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("accountEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetActivities gets the activities property value. The user's activities across devices. Read-only. Nullable.
// returns a []UserActivityable when successful
func (m *User) GetActivities()([]UserActivityable) {
    val, err := m.GetBackingStore().Get("activities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserActivityable)
    }
    return nil
}
// GetAgeGroup gets the ageGroup property value. Sets the age group of the user. Allowed values: null, Minor, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).
// returns a *string when successful
func (m *User) GetAgeGroup()(*string) {
    val, err := m.GetBackingStore().Get("ageGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAgreementAcceptances gets the agreementAcceptances property value. The user's terms of use acceptance statuses. Read-only. Nullable.
// returns a []AgreementAcceptanceable when successful
func (m *User) GetAgreementAcceptances()([]AgreementAcceptanceable) {
    val, err := m.GetBackingStore().Get("agreementAcceptances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementAcceptanceable)
    }
    return nil
}
// GetAppRoleAssignments gets the appRoleAssignments property value. Represents the app roles a user is granted for an application. Supports $expand.
// returns a []AppRoleAssignmentable when successful
func (m *User) GetAppRoleAssignments()([]AppRoleAssignmentable) {
    val, err := m.GetBackingStore().Get("appRoleAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppRoleAssignmentable)
    }
    return nil
}
// GetAssignedLicenses gets the assignedLicenses property value. The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate between directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses. Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).
// returns a []AssignedLicenseable when successful
func (m *User) GetAssignedLicenses()([]AssignedLicenseable) {
    val, err := m.GetBackingStore().Get("assignedLicenses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AssignedLicenseable)
    }
    return nil
}
// GetAssignedPlans gets the assignedPlans property value. The plans that are assigned to the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq and not).
// returns a []AssignedPlanable when successful
func (m *User) GetAssignedPlans()([]AssignedPlanable) {
    val, err := m.GetBackingStore().Get("assignedPlans")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AssignedPlanable)
    }
    return nil
}
// GetAuthentication gets the authentication property value. The authentication methods that are supported for the user.
// returns a Authenticationable when successful
func (m *User) GetAuthentication()(Authenticationable) {
    val, err := m.GetBackingStore().Get("authentication")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Authenticationable)
    }
    return nil
}
// GetAuthorizationInfo gets the authorizationInfo property value. The authorizationInfo property
// returns a AuthorizationInfoable when successful
func (m *User) GetAuthorizationInfo()(AuthorizationInfoable) {
    val, err := m.GetBackingStore().Get("authorizationInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthorizationInfoable)
    }
    return nil
}
// GetBirthday gets the birthday property value. The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.
// returns a *Time when successful
func (m *User) GetBirthday()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("birthday")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetBusinessPhones gets the businessPhones property value. The telephone numbers for the user. NOTE: Although it's a string collection, only one number can be set for this property. Read-only for users synced from the on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).
// returns a []string when successful
func (m *User) GetBusinessPhones()([]string) {
    val, err := m.GetBackingStore().Get("businessPhones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCalendar gets the calendar property value. The user's primary calendar. Read-only.
// returns a Calendarable when successful
func (m *User) GetCalendar()(Calendarable) {
    val, err := m.GetBackingStore().Get("calendar")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Calendarable)
    }
    return nil
}
// GetCalendarGroups gets the calendarGroups property value. The user's calendar groups. Read-only. Nullable.
// returns a []CalendarGroupable when successful
func (m *User) GetCalendarGroups()([]CalendarGroupable) {
    val, err := m.GetBackingStore().Get("calendarGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CalendarGroupable)
    }
    return nil
}
// GetCalendars gets the calendars property value. The user's calendars. Read-only. Nullable.
// returns a []Calendarable when successful
func (m *User) GetCalendars()([]Calendarable) {
    val, err := m.GetBackingStore().Get("calendars")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Calendarable)
    }
    return nil
}
// GetCalendarView gets the calendarView property value. The calendar view for the calendar. Read-only. Nullable.
// returns a []Eventable when successful
func (m *User) GetCalendarView()([]Eventable) {
    val, err := m.GetBackingStore().Get("calendarView")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Eventable)
    }
    return nil
}
// GetChats gets the chats property value. The chats property
// returns a []Chatable when successful
func (m *User) GetChats()([]Chatable) {
    val, err := m.GetBackingStore().Get("chats")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Chatable)
    }
    return nil
}
// GetCity gets the city property value. The city where the user is located. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetCity()(*string) {
    val, err := m.GetBackingStore().Get("city")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCloudClipboard gets the cloudClipboard property value. The cloudClipboard property
// returns a CloudClipboardRootable when successful
func (m *User) GetCloudClipboard()(CloudClipboardRootable) {
    val, err := m.GetBackingStore().Get("cloudClipboard")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CloudClipboardRootable)
    }
    return nil
}
// GetCompanyName gets the companyName property value. The name of the company that the user is associated with. This property can be useful for describing the company that a guest comes from. The maximum length is 64 characters.Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetCompanyName()(*string) {
    val, err := m.GetBackingStore().Get("companyName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConsentProvidedForMinor gets the consentProvidedForMinor property value. Sets whether consent was obtained for minors. Allowed values: null, Granted, Denied, and NotRequired. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).
// returns a *string when successful
func (m *User) GetConsentProvidedForMinor()(*string) {
    val, err := m.GetBackingStore().Get("consentProvidedForMinor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContactFolders gets the contactFolders property value. The user's contacts folders. Read-only. Nullable.
// returns a []ContactFolderable when successful
func (m *User) GetContactFolders()([]ContactFolderable) {
    val, err := m.GetBackingStore().Get("contactFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContactFolderable)
    }
    return nil
}
// GetContacts gets the contacts property value. The user's contacts. Read-only. Nullable.
// returns a []Contactable when successful
func (m *User) GetContacts()([]Contactable) {
    val, err := m.GetBackingStore().Get("contacts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Contactable)
    }
    return nil
}
// GetCountry gets the country property value. The country/region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetCountry()(*string) {
    val, err := m.GetBackingStore().Get("country")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time the user was created, in ISO 8601 format and UTC. The value can't be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
// returns a *Time when successful
func (m *User) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedObjects gets the createdObjects property value. Directory objects that the user created. Read-only. Nullable.
// returns a []DirectoryObjectable when successful
func (m *User) GetCreatedObjects()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("createdObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetCreationType gets the creationType property value. Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by a guest signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).
// returns a *string when successful
func (m *User) GetCreationType()(*string) {
    val, err := m.GetBackingStore().Get("creationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCustomSecurityAttributes gets the customSecurityAttributes property value. An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). The filter value is case-sensitive. To read this property, the calling app must be assigned the CustomSecAttributeAssignment.Read.All permission. To write this property, the calling app must be assigned the CustomSecAttributeAssignment.ReadWrite.All permissions. To read or write this property in delegated scenarios, the admin must be assigned the Attribute Assignment Administrator role.
// returns a CustomSecurityAttributeValueable when successful
func (m *User) GetCustomSecurityAttributes()(CustomSecurityAttributeValueable) {
    val, err := m.GetBackingStore().Get("customSecurityAttributes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CustomSecurityAttributeValueable)
    }
    return nil
}
// GetDepartment gets the department property value. The name of the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).
// returns a *string when successful
func (m *User) GetDepartment()(*string) {
    val, err := m.GetBackingStore().Get("department")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceEnrollmentLimit gets the deviceEnrollmentLimit property value. The limit on the maximum number of devices that the user is permitted to enroll. Allowed values are 5 or 1000.
// returns a *int32 when successful
func (m *User) GetDeviceEnrollmentLimit()(*int32) {
    val, err := m.GetBackingStore().Get("deviceEnrollmentLimit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeviceManagementTroubleshootingEvents gets the deviceManagementTroubleshootingEvents property value. The list of troubleshooting events for this user.
// returns a []DeviceManagementTroubleshootingEventable when successful
func (m *User) GetDeviceManagementTroubleshootingEvents()([]DeviceManagementTroubleshootingEventable) {
    val, err := m.GetBackingStore().Get("deviceManagementTroubleshootingEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementTroubleshootingEventable)
    }
    return nil
}
// GetDirectReports gets the directReports property value. The users and contacts that report to the user. (The users and contacts that have their manager property set to this user.) Read-only. Nullable. Supports $expand.
// returns a []DirectoryObjectable when successful
func (m *User) GetDirectReports()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("directReports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name displayed in the address book for the user. This value is usually the combination of the user's first name, middle initial, and family name. This property is required when a user is created and it can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.
// returns a *string when successful
func (m *User) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDrive gets the drive property value. The user's OneDrive. Read-only.
// returns a Driveable when successful
func (m *User) GetDrive()(Driveable) {
    val, err := m.GetBackingStore().Get("drive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Driveable)
    }
    return nil
}
// GetDrives gets the drives property value. A collection of drives available for this user. Read-only.
// returns a []Driveable when successful
func (m *User) GetDrives()([]Driveable) {
    val, err := m.GetBackingStore().Get("drives")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Driveable)
    }
    return nil
}
// GetEmployeeExperience gets the employeeExperience property value. The employeeExperience property
// returns a EmployeeExperienceUserable when successful
func (m *User) GetEmployeeExperience()(EmployeeExperienceUserable) {
    val, err := m.GetBackingStore().Get("employeeExperience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmployeeExperienceUserable)
    }
    return nil
}
// GetEmployeeHireDate gets the employeeHireDate property value. The date and time when the user was hired or will start work in a future hire. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
// returns a *Time when successful
func (m *User) GetEmployeeHireDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("employeeHireDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEmployeeId gets the employeeId property value. The employee identifier assigned to the user by the organization. The maximum length is 16 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetEmployeeId()(*string) {
    val, err := m.GetBackingStore().Get("employeeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmployeeLeaveDateTime gets the employeeLeaveDateTime property value. The date and time when the user left or will leave the organization. To read this property, the calling app must be assigned the User-LifeCycleInfo.Read.All permission. To write this property, the calling app must be assigned the User.Read.All and User-LifeCycleInfo.ReadWrite.All permissions. To read this property in delegated scenarios, the admin needs at least one of the following Microsoft Entra roles: Lifecycle Workflows Administrator (least privilege), Global Reader. To write this property in delegated scenarios, the admin needs the Global Administrator role. Supports $filter (eq, ne, not , ge, le, in). For more information, see Configure the employeeLeaveDateTime property for a user.
// returns a *Time when successful
func (m *User) GetEmployeeLeaveDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("employeeLeaveDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEmployeeOrgData gets the employeeOrgData property value. Represents organization data (for example, division and costCenter) associated with a user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
// returns a EmployeeOrgDataable when successful
func (m *User) GetEmployeeOrgData()(EmployeeOrgDataable) {
    val, err := m.GetBackingStore().Get("employeeOrgData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmployeeOrgDataable)
    }
    return nil
}
// GetEmployeeType gets the employeeType property value. Captures enterprise worker type. For example, Employee, Contractor, Consultant, or Vendor. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith).
// returns a *string when successful
func (m *User) GetEmployeeType()(*string) {
    val, err := m.GetBackingStore().Get("employeeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEvents gets the events property value. The user's events. Default is to show Events under the Default Calendar. Read-only. Nullable.
// returns a []Eventable when successful
func (m *User) GetEvents()([]Eventable) {
    val, err := m.GetBackingStore().Get("events")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Eventable)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the user. Read-only. Supports $expand. Nullable.
// returns a []Extensionable when successful
func (m *User) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetExternalUserState gets the externalUserState property value. For a guest invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be PendingAcceptance or Accepted, or null for all other users. Returned only on $select. Supports $filter (eq, ne, not , in).
// returns a *string when successful
func (m *User) GetExternalUserState()(*string) {
    val, err := m.GetBackingStore().Get("externalUserState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalUserStateChangeDateTime gets the externalUserStateChangeDateTime property value. Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. Supports $filter (eq, ne, not , in).
// returns a *Time when successful
func (m *User) GetExternalUserStateChangeDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("externalUserStateChangeDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFaxNumber gets the faxNumber property value. The fax number of the user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetFaxNumber()(*string) {
    val, err := m.GetBackingStore().Get("faxNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *User) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["aboutMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAboutMe(val)
        }
        return nil
    }
    res["accountEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountEnabled(val)
        }
        return nil
    }
    res["activities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserActivityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserActivityable)
                }
            }
            m.SetActivities(res)
        }
        return nil
    }
    res["ageGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAgeGroup(val)
        }
        return nil
    }
    res["agreementAcceptances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAgreementAcceptanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AgreementAcceptanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AgreementAcceptanceable)
                }
            }
            m.SetAgreementAcceptances(res)
        }
        return nil
    }
    res["appRoleAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppRoleAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppRoleAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppRoleAssignmentable)
                }
            }
            m.SetAppRoleAssignments(res)
        }
        return nil
    }
    res["assignedLicenses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAssignedLicenseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AssignedLicenseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AssignedLicenseable)
                }
            }
            m.SetAssignedLicenses(res)
        }
        return nil
    }
    res["assignedPlans"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAssignedPlanFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AssignedPlanable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AssignedPlanable)
                }
            }
            m.SetAssignedPlans(res)
        }
        return nil
    }
    res["authentication"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthentication(val.(Authenticationable))
        }
        return nil
    }
    res["authorizationInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthorizationInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthorizationInfo(val.(AuthorizationInfoable))
        }
        return nil
    }
    res["birthday"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBirthday(val)
        }
        return nil
    }
    res["businessPhones"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetBusinessPhones(res)
        }
        return nil
    }
    res["calendar"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCalendarFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalendar(val.(Calendarable))
        }
        return nil
    }
    res["calendarGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCalendarGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CalendarGroupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CalendarGroupable)
                }
            }
            m.SetCalendarGroups(res)
        }
        return nil
    }
    res["calendars"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCalendarFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Calendarable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Calendarable)
                }
            }
            m.SetCalendars(res)
        }
        return nil
    }
    res["calendarView"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Eventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Eventable)
                }
            }
            m.SetCalendarView(res)
        }
        return nil
    }
    res["chats"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Chatable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Chatable)
                }
            }
            m.SetChats(res)
        }
        return nil
    }
    res["city"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCity(val)
        }
        return nil
    }
    res["cloudClipboard"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCloudClipboardRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudClipboard(val.(CloudClipboardRootable))
        }
        return nil
    }
    res["companyName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompanyName(val)
        }
        return nil
    }
    res["consentProvidedForMinor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConsentProvidedForMinor(val)
        }
        return nil
    }
    res["contactFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContactFolderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContactFolderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContactFolderable)
                }
            }
            m.SetContactFolders(res)
        }
        return nil
    }
    res["contacts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContactFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Contactable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Contactable)
                }
            }
            m.SetContacts(res)
        }
        return nil
    }
    res["country"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountry(val)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["createdObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetCreatedObjects(res)
        }
        return nil
    }
    res["creationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreationType(val)
        }
        return nil
    }
    res["customSecurityAttributes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCustomSecurityAttributeValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomSecurityAttributes(val.(CustomSecurityAttributeValueable))
        }
        return nil
    }
    res["department"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDepartment(val)
        }
        return nil
    }
    res["deviceEnrollmentLimit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceEnrollmentLimit(val)
        }
        return nil
    }
    res["deviceManagementTroubleshootingEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementTroubleshootingEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementTroubleshootingEventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceManagementTroubleshootingEventable)
                }
            }
            m.SetDeviceManagementTroubleshootingEvents(res)
        }
        return nil
    }
    res["directReports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetDirectReports(res)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["drive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDrive(val.(Driveable))
        }
        return nil
    }
    res["drives"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Driveable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Driveable)
                }
            }
            m.SetDrives(res)
        }
        return nil
    }
    res["employeeExperience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmployeeExperienceUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeExperience(val.(EmployeeExperienceUserable))
        }
        return nil
    }
    res["employeeHireDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeHireDate(val)
        }
        return nil
    }
    res["employeeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeId(val)
        }
        return nil
    }
    res["employeeLeaveDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeLeaveDateTime(val)
        }
        return nil
    }
    res["employeeOrgData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmployeeOrgDataFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeOrgData(val.(EmployeeOrgDataable))
        }
        return nil
    }
    res["employeeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmployeeType(val)
        }
        return nil
    }
    res["events"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Eventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Eventable)
                }
            }
            m.SetEvents(res)
        }
        return nil
    }
    res["extensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Extensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Extensionable)
                }
            }
            m.SetExtensions(res)
        }
        return nil
    }
    res["externalUserState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalUserState(val)
        }
        return nil
    }
    res["externalUserStateChangeDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalUserStateChangeDateTime(val)
        }
        return nil
    }
    res["faxNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFaxNumber(val)
        }
        return nil
    }
    res["followedSites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Siteable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Siteable)
                }
            }
            m.SetFollowedSites(res)
        }
        return nil
    }
    res["givenName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGivenName(val)
        }
        return nil
    }
    res["hireDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHireDate(val)
        }
        return nil
    }
    res["identities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateObjectIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ObjectIdentityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ObjectIdentityable)
                }
            }
            m.SetIdentities(res)
        }
        return nil
    }
    res["imAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetImAddresses(res)
        }
        return nil
    }
    res["inferenceClassification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInferenceClassificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInferenceClassification(val.(InferenceClassificationable))
        }
        return nil
    }
    res["insights"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemInsightsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInsights(val.(ItemInsightsable))
        }
        return nil
    }
    res["interests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetInterests(res)
        }
        return nil
    }
    res["isManagementRestricted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsManagementRestricted(val)
        }
        return nil
    }
    res["isResourceAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsResourceAccount(val)
        }
        return nil
    }
    res["jobTitle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJobTitle(val)
        }
        return nil
    }
    res["joinedTeams"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Teamable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Teamable)
                }
            }
            m.SetJoinedTeams(res)
        }
        return nil
    }
    res["lastPasswordChangeDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastPasswordChangeDateTime(val)
        }
        return nil
    }
    res["legalAgeGroupClassification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLegalAgeGroupClassification(val)
        }
        return nil
    }
    res["licenseAssignmentStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLicenseAssignmentStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LicenseAssignmentStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LicenseAssignmentStateable)
                }
            }
            m.SetLicenseAssignmentStates(res)
        }
        return nil
    }
    res["licenseDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLicenseDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LicenseDetailsable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LicenseDetailsable)
                }
            }
            m.SetLicenseDetails(res)
        }
        return nil
    }
    res["mail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMail(val)
        }
        return nil
    }
    res["mailboxSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMailboxSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailboxSettings(val.(MailboxSettingsable))
        }
        return nil
    }
    res["mailFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailFolderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailFolderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailFolderable)
                }
            }
            m.SetMailFolders(res)
        }
        return nil
    }
    res["mailNickname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailNickname(val)
        }
        return nil
    }
    res["managedAppRegistrations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppRegistrationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppRegistrationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppRegistrationable)
                }
            }
            m.SetManagedAppRegistrations(res)
        }
        return nil
    }
    res["managedDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedDeviceable)
                }
            }
            m.SetManagedDevices(res)
        }
        return nil
    }
    res["manager"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManager(val.(DirectoryObjectable))
        }
        return nil
    }
    res["memberOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetMemberOf(res)
        }
        return nil
    }
    res["messages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Messageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Messageable)
                }
            }
            m.SetMessages(res)
        }
        return nil
    }
    res["mobilePhone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMobilePhone(val)
        }
        return nil
    }
    res["mySite"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMySite(val)
        }
        return nil
    }
    res["oauth2PermissionGrants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOAuth2PermissionGrantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OAuth2PermissionGrantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OAuth2PermissionGrantable)
                }
            }
            m.SetOauth2PermissionGrants(res)
        }
        return nil
    }
    res["officeLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOfficeLocation(val)
        }
        return nil
    }
    res["onenote"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnenoteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnenote(val.(Onenoteable))
        }
        return nil
    }
    res["onlineMeetings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnlineMeetingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnlineMeetingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnlineMeetingable)
                }
            }
            m.SetOnlineMeetings(res)
        }
        return nil
    }
    res["onPremisesDistinguishedName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesDistinguishedName(val)
        }
        return nil
    }
    res["onPremisesDomainName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesDomainName(val)
        }
        return nil
    }
    res["onPremisesExtensionAttributes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnPremisesExtensionAttributesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesExtensionAttributes(val.(OnPremisesExtensionAttributesable))
        }
        return nil
    }
    res["onPremisesImmutableId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesImmutableId(val)
        }
        return nil
    }
    res["onPremisesLastSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesLastSyncDateTime(val)
        }
        return nil
    }
    res["onPremisesProvisioningErrors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnPremisesProvisioningErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnPremisesProvisioningErrorable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnPremisesProvisioningErrorable)
                }
            }
            m.SetOnPremisesProvisioningErrors(res)
        }
        return nil
    }
    res["onPremisesSamAccountName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesSamAccountName(val)
        }
        return nil
    }
    res["onPremisesSecurityIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesSecurityIdentifier(val)
        }
        return nil
    }
    res["onPremisesSyncEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesSyncEnabled(val)
        }
        return nil
    }
    res["onPremisesUserPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesUserPrincipalName(val)
        }
        return nil
    }
    res["otherMails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetOtherMails(res)
        }
        return nil
    }
    res["outlook"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOutlookUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOutlook(val.(OutlookUserable))
        }
        return nil
    }
    res["ownedDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetOwnedDevices(res)
        }
        return nil
    }
    res["ownedObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetOwnedObjects(res)
        }
        return nil
    }
    res["passwordPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordPolicies(val)
        }
        return nil
    }
    res["passwordProfile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePasswordProfileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordProfile(val.(PasswordProfileable))
        }
        return nil
    }
    res["pastProjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetPastProjects(res)
        }
        return nil
    }
    res["people"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePersonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Personable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Personable)
                }
            }
            m.SetPeople(res)
        }
        return nil
    }
    res["permissionGrants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateResourceSpecificPermissionGrantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ResourceSpecificPermissionGrantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ResourceSpecificPermissionGrantable)
                }
            }
            m.SetPermissionGrants(res)
        }
        return nil
    }
    res["photo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateProfilePhotoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoto(val.(ProfilePhotoable))
        }
        return nil
    }
    res["photos"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProfilePhotoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProfilePhotoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProfilePhotoable)
                }
            }
            m.SetPhotos(res)
        }
        return nil
    }
    res["planner"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlanner(val.(PlannerUserable))
        }
        return nil
    }
    res["postalCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPostalCode(val)
        }
        return nil
    }
    res["preferredDataLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreferredDataLocation(val)
        }
        return nil
    }
    res["preferredLanguage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreferredLanguage(val)
        }
        return nil
    }
    res["preferredName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreferredName(val)
        }
        return nil
    }
    res["presence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePresenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPresence(val.(Presenceable))
        }
        return nil
    }
    res["print"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserPrintFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrint(val.(UserPrintable))
        }
        return nil
    }
    res["provisionedPlans"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProvisionedPlanFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProvisionedPlanable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProvisionedPlanable)
                }
            }
            m.SetProvisionedPlans(res)
        }
        return nil
    }
    res["proxyAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetProxyAddresses(res)
        }
        return nil
    }
    res["registeredDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetRegisteredDevices(res)
        }
        return nil
    }
    res["responsibilities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetResponsibilities(res)
        }
        return nil
    }
    res["schools"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetSchools(res)
        }
        return nil
    }
    res["scopedRoleMemberOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateScopedRoleMembershipFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ScopedRoleMembershipable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ScopedRoleMembershipable)
                }
            }
            m.SetScopedRoleMemberOf(res)
        }
        return nil
    }
    res["securityIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecurityIdentifier(val)
        }
        return nil
    }
    res["serviceProvisioningErrors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceProvisioningErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceProvisioningErrorable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceProvisioningErrorable)
                }
            }
            m.SetServiceProvisioningErrors(res)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(UserSettingsable))
        }
        return nil
    }
    res["showInAddressList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowInAddressList(val)
        }
        return nil
    }
    res["signInActivity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSignInActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInActivity(val.(SignInActivityable))
        }
        return nil
    }
    res["signInSessionsValidFromDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInSessionsValidFromDateTime(val)
        }
        return nil
    }
    res["skills"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetSkills(res)
        }
        return nil
    }
    res["solutions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserSolutionRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSolutions(val.(UserSolutionRootable))
        }
        return nil
    }
    res["sponsors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetSponsors(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val)
        }
        return nil
    }
    res["streetAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStreetAddress(val)
        }
        return nil
    }
    res["surname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSurname(val)
        }
        return nil
    }
    res["teamwork"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserTeamworkFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamwork(val.(UserTeamworkable))
        }
        return nil
    }
    res["todo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTodoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTodo(val.(Todoable))
        }
        return nil
    }
    res["transitiveMemberOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetTransitiveMemberOf(res)
        }
        return nil
    }
    res["usageLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsageLocation(val)
        }
        return nil
    }
    res["userPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPrincipalName(val)
        }
        return nil
    }
    res["userType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserType(val)
        }
        return nil
    }
    return res
}
// GetFollowedSites gets the followedSites property value. The followedSites property
// returns a []Siteable when successful
func (m *User) GetFollowedSites()([]Siteable) {
    val, err := m.GetBackingStore().Get("followedSites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Siteable)
    }
    return nil
}
// GetGivenName gets the givenName property value. The given name (first name) of the user. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetGivenName()(*string) {
    val, err := m.GetBackingStore().Get("givenName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHireDate gets the hireDate property value. The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.  Note: This property is specific to SharePoint in Microsoft 365. We recommend using the native employeeHireDate property to set and update hire date values using Microsoft Graph APIs.
// returns a *Time when successful
func (m *User) GetHireDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("hireDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetIdentities gets the identities property value. Represents the identities that can be used to sign in to this user account. Microsoft (also known as a local account), organizations, or social identity providers such as Facebook, Google, and Microsoft can provide identity and tie it to a user account. It might contain multiple items with the same signInType value. Returned only on $select.  Supports $filter (eq) with limitations.
// returns a []ObjectIdentityable when successful
func (m *User) GetIdentities()([]ObjectIdentityable) {
    val, err := m.GetBackingStore().Get("identities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ObjectIdentityable)
    }
    return nil
}
// GetImAddresses gets the imAddresses property value. The instant message voice-over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).
// returns a []string when successful
func (m *User) GetImAddresses()([]string) {
    val, err := m.GetBackingStore().Get("imAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetInferenceClassification gets the inferenceClassification property value. Relevance classification of the user's messages based on explicit designations that override inferred relevance or importance.
// returns a InferenceClassificationable when successful
func (m *User) GetInferenceClassification()(InferenceClassificationable) {
    val, err := m.GetBackingStore().Get("inferenceClassification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InferenceClassificationable)
    }
    return nil
}
// GetInsights gets the insights property value. Represents relationships between a user and items such as OneDrive for work or school documents, calculated using advanced analytics and machine learning techniques. Read-only. Nullable.
// returns a ItemInsightsable when successful
func (m *User) GetInsights()(ItemInsightsable) {
    val, err := m.GetBackingStore().Get("insights")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemInsightsable)
    }
    return nil
}
// GetInterests gets the interests property value. A list for the user to describe their interests. Returned only on $select.
// returns a []string when successful
func (m *User) GetInterests()([]string) {
    val, err := m.GetBackingStore().Get("interests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetIsManagementRestricted gets the isManagementRestricted property value. The isManagementRestricted property
// returns a *bool when successful
func (m *User) GetIsManagementRestricted()(*bool) {
    val, err := m.GetBackingStore().Get("isManagementRestricted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsResourceAccount gets the isResourceAccount property value. Don't use – reserved for future use.
// returns a *bool when successful
func (m *User) GetIsResourceAccount()(*bool) {
    val, err := m.GetBackingStore().Get("isResourceAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJobTitle gets the jobTitle property value. The user's job title. Maximum length is 128 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetJobTitle()(*string) {
    val, err := m.GetBackingStore().Get("jobTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetJoinedTeams gets the joinedTeams property value. The joinedTeams property
// returns a []Teamable when successful
func (m *User) GetJoinedTeams()([]Teamable) {
    val, err := m.GetBackingStore().Get("joinedTeams")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Teamable)
    }
    return nil
}
// GetLastPasswordChangeDateTime gets the lastPasswordChangeDateTime property value. The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.
// returns a *Time when successful
func (m *User) GetLastPasswordChangeDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastPasswordChangeDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLegalAgeGroupClassification gets the legalAgeGroupClassification property value. Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select.
// returns a *string when successful
func (m *User) GetLegalAgeGroupClassification()(*string) {
    val, err := m.GetBackingStore().Get("legalAgeGroupClassification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLicenseAssignmentStates gets the licenseAssignmentStates property value. State of license assignments for this user. Also indicates licenses that are directly assigned or the user inherited through group memberships. Read-only. Returned only on $select.
// returns a []LicenseAssignmentStateable when successful
func (m *User) GetLicenseAssignmentStates()([]LicenseAssignmentStateable) {
    val, err := m.GetBackingStore().Get("licenseAssignmentStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LicenseAssignmentStateable)
    }
    return nil
}
// GetLicenseDetails gets the licenseDetails property value. A collection of this user's license details. Read-only.
// returns a []LicenseDetailsable when successful
func (m *User) GetLicenseDetails()([]LicenseDetailsable) {
    val, err := m.GetBackingStore().Get("licenseDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LicenseDetailsable)
    }
    return nil
}
// GetMail gets the mail property value. The SMTP address for the user, for example, jeff@contoso.com. Changes to this property update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetMail()(*string) {
    val, err := m.GetBackingStore().Get("mail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMailboxSettings gets the mailboxSettings property value. Settings for the primary mailbox of the signed-in user. You can get or update settings for sending automatic replies to incoming messages, locale, and time zone. Returned only on $select.
// returns a MailboxSettingsable when successful
func (m *User) GetMailboxSettings()(MailboxSettingsable) {
    val, err := m.GetBackingStore().Get("mailboxSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MailboxSettingsable)
    }
    return nil
}
// GetMailFolders gets the mailFolders property value. The user's mail folders. Read-only. Nullable.
// returns a []MailFolderable when successful
func (m *User) GetMailFolders()([]MailFolderable) {
    val, err := m.GetBackingStore().Get("mailFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailFolderable)
    }
    return nil
}
// GetMailNickname gets the mailNickname property value. The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetMailNickname()(*string) {
    val, err := m.GetBackingStore().Get("mailNickname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetManagedAppRegistrations gets the managedAppRegistrations property value. Zero or more managed app registrations that belong to the user.
// returns a []ManagedAppRegistrationable when successful
func (m *User) GetManagedAppRegistrations()([]ManagedAppRegistrationable) {
    val, err := m.GetBackingStore().Get("managedAppRegistrations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppRegistrationable)
    }
    return nil
}
// GetManagedDevices gets the managedDevices property value. The managed devices associated with the user.
// returns a []ManagedDeviceable when successful
func (m *User) GetManagedDevices()([]ManagedDeviceable) {
    val, err := m.GetBackingStore().Get("managedDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceable)
    }
    return nil
}
// GetManager gets the manager property value. The user or contact that is this user's manager. Read-only. Supports $expand.
// returns a DirectoryObjectable when successful
func (m *User) GetManager()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("manager")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetMemberOf gets the memberOf property value. The groups and directory roles that the user is a member of. Read-only. Nullable. Supports $expand.
// returns a []DirectoryObjectable when successful
func (m *User) GetMemberOf()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("memberOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetMessages gets the messages property value. The messages in a mailbox or folder. Read-only. Nullable.
// returns a []Messageable when successful
func (m *User) GetMessages()([]Messageable) {
    val, err := m.GetBackingStore().Get("messages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Messageable)
    }
    return nil
}
// GetMobilePhone gets the mobilePhone property value. The primary cellular telephone number for the user. Read-only for users synced from the on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.
// returns a *string when successful
func (m *User) GetMobilePhone()(*string) {
    val, err := m.GetBackingStore().Get("mobilePhone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMySite gets the mySite property value. The URL for the user's site. Returned only on $select.
// returns a *string when successful
func (m *User) GetMySite()(*string) {
    val, err := m.GetBackingStore().Get("mySite")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOauth2PermissionGrants gets the oauth2PermissionGrants property value. The oauth2PermissionGrants property
// returns a []OAuth2PermissionGrantable when successful
func (m *User) GetOauth2PermissionGrants()([]OAuth2PermissionGrantable) {
    val, err := m.GetBackingStore().Get("oauth2PermissionGrants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OAuth2PermissionGrantable)
    }
    return nil
}
// GetOfficeLocation gets the officeLocation property value. The office location in the user's place of business. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetOfficeLocation()(*string) {
    val, err := m.GetBackingStore().Get("officeLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnenote gets the onenote property value. The onenote property
// returns a Onenoteable when successful
func (m *User) GetOnenote()(Onenoteable) {
    val, err := m.GetBackingStore().Get("onenote")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Onenoteable)
    }
    return nil
}
// GetOnlineMeetings gets the onlineMeetings property value. Information about a meeting, including the URL used to join a meeting, the attendees list, and the description.
// returns a []OnlineMeetingable when successful
func (m *User) GetOnlineMeetings()([]OnlineMeetingable) {
    val, err := m.GetBackingStore().Get("onlineMeetings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnlineMeetingable)
    }
    return nil
}
// GetOnPremisesDistinguishedName gets the onPremisesDistinguishedName property value. Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.
// returns a *string when successful
func (m *User) GetOnPremisesDistinguishedName()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesDistinguishedName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesDomainName gets the onPremisesDomainName property value. Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.
// returns a *string when successful
func (m *User) GetOnPremisesDomainName()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesDomainName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesExtensionAttributes gets the onPremisesExtensionAttributes property value. Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. Each attribute can store up to 1024 characters. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during the creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).
// returns a OnPremisesExtensionAttributesable when successful
func (m *User) GetOnPremisesExtensionAttributes()(OnPremisesExtensionAttributesable) {
    val, err := m.GetBackingStore().Get("onPremisesExtensionAttributes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnPremisesExtensionAttributesable)
    }
    return nil
}
// GetOnPremisesImmutableId gets the onPremisesImmutableId property value. This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property. NOTE: The $ and _ characters can't be used when specifying this property. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).
// returns a *string when successful
func (m *User) GetOnPremisesImmutableId()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesImmutableId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesLastSyncDateTime gets the onPremisesLastSyncDateTime property value. Indicates the last time at which the object was synced with the on-premises directory; for example: 2013-02-16T03:04:54Z. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).
// returns a *Time when successful
func (m *User) GetOnPremisesLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("onPremisesLastSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetOnPremisesProvisioningErrors gets the onPremisesProvisioningErrors property value. Errors when using Microsoft synchronization product during provisioning. Returned only on $select. Supports $filter (eq, not, ge, le).
// returns a []OnPremisesProvisioningErrorable when successful
func (m *User) GetOnPremisesProvisioningErrors()([]OnPremisesProvisioningErrorable) {
    val, err := m.GetBackingStore().Get("onPremisesProvisioningErrors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnPremisesProvisioningErrorable)
    }
    return nil
}
// GetOnPremisesSamAccountName gets the onPremisesSamAccountName property value. Contains the on-premises samAccountName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).
// returns a *string when successful
func (m *User) GetOnPremisesSamAccountName()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesSamAccountName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesSecurityIdentifier gets the onPremisesSecurityIdentifier property value. Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq including on null values).
// returns a *string when successful
func (m *User) GetOnPremisesSecurityIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesSecurityIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesSyncEnabled gets the onPremisesSyncEnabled property value. true if this user object is currently being synced from an on-premises Active Directory (AD); otherwise the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values).
// returns a *bool when successful
func (m *User) GetOnPremisesSyncEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("onPremisesSyncEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOnPremisesUserPrincipalName gets the onPremisesUserPrincipalName property value. Contains the on-premises userPrincipalName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).
// returns a *string when successful
func (m *User) GetOnPremisesUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesUserPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOtherMails gets the otherMails property value. A list of other email addresses for the user; for example: ['bob@contoso.com', 'Robert@fabrikam.com']. NOTE: This property can't contain accent characters. Returned only on $select. Supports $filter (eq, not, ge, le, in, startsWith, endsWith, /$count eq 0, /$count ne 0).
// returns a []string when successful
func (m *User) GetOtherMails()([]string) {
    val, err := m.GetBackingStore().Get("otherMails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetOutlook gets the outlook property value. The outlook property
// returns a OutlookUserable when successful
func (m *User) GetOutlook()(OutlookUserable) {
    val, err := m.GetBackingStore().Get("outlook")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OutlookUserable)
    }
    return nil
}
// GetOwnedDevices gets the ownedDevices property value. Devices the user owns. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
// returns a []DirectoryObjectable when successful
func (m *User) GetOwnedDevices()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("ownedDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetOwnedObjects gets the ownedObjects property value. Directory objects the user owns. Read-only. Nullable. Supports $expand, $select nested in $expand, and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
// returns a []DirectoryObjectable when successful
func (m *User) GetOwnedObjects()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("ownedObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetPasswordPolicies gets the passwordPolicies property value. Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two might be specified together; for example: DisablePasswordExpiration, DisableStrongPassword. Returned only on $select. For more information on the default password policies, see Microsoft Entra password policies. Supports $filter (ne, not, and eq on null values).
// returns a *string when successful
func (m *User) GetPasswordPolicies()(*string) {
    val, err := m.GetBackingStore().Get("passwordPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasswordProfile gets the passwordProfile property value. Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). To update this property:  In delegated access, the calling app must be assigned the Directory.AccessAsUser.All delegated permission on behalf of the signed-in user.  In application-only access, the calling app must be assigned the User.ReadWrite.All (least privilege) or Directory.ReadWrite.All (higher privilege) application permission and at least the User Administrator Microsoft Entra role.
// returns a PasswordProfileable when successful
func (m *User) GetPasswordProfile()(PasswordProfileable) {
    val, err := m.GetBackingStore().Get("passwordProfile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PasswordProfileable)
    }
    return nil
}
// GetPastProjects gets the pastProjects property value. A list for the user to enumerate their past projects. Returned only on $select.
// returns a []string when successful
func (m *User) GetPastProjects()([]string) {
    val, err := m.GetBackingStore().Get("pastProjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPeople gets the people property value. People that are relevant to the user. Read-only. Nullable.
// returns a []Personable when successful
func (m *User) GetPeople()([]Personable) {
    val, err := m.GetBackingStore().Get("people")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Personable)
    }
    return nil
}
// GetPermissionGrants gets the permissionGrants property value. List all resource-specific permission grants of a user.
// returns a []ResourceSpecificPermissionGrantable when successful
func (m *User) GetPermissionGrants()([]ResourceSpecificPermissionGrantable) {
    val, err := m.GetBackingStore().Get("permissionGrants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ResourceSpecificPermissionGrantable)
    }
    return nil
}
// GetPhoto gets the photo property value. The user's profile photo. Read-only.
// returns a ProfilePhotoable when successful
func (m *User) GetPhoto()(ProfilePhotoable) {
    val, err := m.GetBackingStore().Get("photo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ProfilePhotoable)
    }
    return nil
}
// GetPhotos gets the photos property value. The collection of the user's profile photos in different sizes. Read-only.
// returns a []ProfilePhotoable when successful
func (m *User) GetPhotos()([]ProfilePhotoable) {
    val, err := m.GetBackingStore().Get("photos")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProfilePhotoable)
    }
    return nil
}
// GetPlanner gets the planner property value. Entry-point to the Planner resource that might exist for a user. Read-only.
// returns a PlannerUserable when successful
func (m *User) GetPlanner()(PlannerUserable) {
    val, err := m.GetBackingStore().Get("planner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerUserable)
    }
    return nil
}
// GetPostalCode gets the postalCode property value. The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetPostalCode()(*string) {
    val, err := m.GetBackingStore().Get("postalCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreferredDataLocation gets the preferredDataLocation property value. The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.
// returns a *string when successful
func (m *User) GetPreferredDataLocation()(*string) {
    val, err := m.GetBackingStore().Get("preferredDataLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreferredLanguage gets the preferredLanguage property value. The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language, and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)
// returns a *string when successful
func (m *User) GetPreferredLanguage()(*string) {
    val, err := m.GetBackingStore().Get("preferredLanguage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreferredName gets the preferredName property value. The preferred name for the user. Not Supported. This attribute returns an empty string.Returned only on $select.
// returns a *string when successful
func (m *User) GetPreferredName()(*string) {
    val, err := m.GetBackingStore().Get("preferredName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPresence gets the presence property value. The presence property
// returns a Presenceable when successful
func (m *User) GetPresence()(Presenceable) {
    val, err := m.GetBackingStore().Get("presence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Presenceable)
    }
    return nil
}
// GetPrint gets the print property value. The print property
// returns a UserPrintable when successful
func (m *User) GetPrint()(UserPrintable) {
    val, err := m.GetBackingStore().Get("print")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserPrintable)
    }
    return nil
}
// GetProvisionedPlans gets the provisionedPlans property value. The plans that are provisioned for the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le).
// returns a []ProvisionedPlanable when successful
func (m *User) GetProvisionedPlans()([]ProvisionedPlanable) {
    val, err := m.GetBackingStore().Get("provisionedPlans")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProvisionedPlanable)
    }
    return nil
}
// GetProxyAddresses gets the proxyAddresses property value. For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address, while those addresses prefixed with smtp are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of 10 unique addresses. Read-only in Microsoft Graph; you can update this property only through the Microsoft 365 admin center. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).
// returns a []string when successful
func (m *User) GetProxyAddresses()([]string) {
    val, err := m.GetBackingStore().Get("proxyAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetRegisteredDevices gets the registeredDevices property value. Devices that are registered for the user. Read-only. Nullable. Supports $expand and returns up to 100 objects.
// returns a []DirectoryObjectable when successful
func (m *User) GetRegisteredDevices()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("registeredDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetResponsibilities gets the responsibilities property value. A list for the user to enumerate their responsibilities. Returned only on $select.
// returns a []string when successful
func (m *User) GetResponsibilities()([]string) {
    val, err := m.GetBackingStore().Get("responsibilities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSchools gets the schools property value. A list for the user to enumerate the schools they attended. Returned only on $select.
// returns a []string when successful
func (m *User) GetSchools()([]string) {
    val, err := m.GetBackingStore().Get("schools")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetScopedRoleMemberOf gets the scopedRoleMemberOf property value. The scopedRoleMemberOf property
// returns a []ScopedRoleMembershipable when successful
func (m *User) GetScopedRoleMemberOf()([]ScopedRoleMembershipable) {
    val, err := m.GetBackingStore().Get("scopedRoleMemberOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ScopedRoleMembershipable)
    }
    return nil
}
// GetSecurityIdentifier gets the securityIdentifier property value. Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).
// returns a *string when successful
func (m *User) GetSecurityIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("securityIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServiceProvisioningErrors gets the serviceProvisioningErrors property value. Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a user object.  Supports $filter (eq, not, for isResolved and serviceInstance).
// returns a []ServiceProvisioningErrorable when successful
func (m *User) GetServiceProvisioningErrors()([]ServiceProvisioningErrorable) {
    val, err := m.GetBackingStore().Get("serviceProvisioningErrors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceProvisioningErrorable)
    }
    return nil
}
// GetSettings gets the settings property value. The settings property
// returns a UserSettingsable when successful
func (m *User) GetSettings()(UserSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserSettingsable)
    }
    return nil
}
// GetShowInAddressList gets the showInAddressList property value. Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead. Represents whether the user should be included in the Outlook global address list. See Known issue.
// returns a *bool when successful
func (m *User) GetShowInAddressList()(*bool) {
    val, err := m.GetBackingStore().Get("showInAddressList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSignInActivity gets the signInActivity property value. Get the last signed-in date and request ID of the sign-in for a given user. Read-only.Returned only on $select. Supports $filter (eq, ne, not, ge, le) but not with any other filterable properties. Note: Details for this property require a Microsoft Entra ID P1 or P2 license and the AuditLog.Read.All permission.This property isn't returned for a user who never signed in or last signed in before April 2020.
// returns a SignInActivityable when successful
func (m *User) GetSignInActivity()(SignInActivityable) {
    val, err := m.GetBackingStore().Get("signInActivity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SignInActivityable)
    }
    return nil
}
// GetSignInSessionsValidFromDateTime gets the signInSessionsValidFromDateTime property value. Any refresh tokens or session tokens (session cookies) issued before this time are invalid. Applications get an error when using an invalid refresh or session token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.
// returns a *Time when successful
func (m *User) GetSignInSessionsValidFromDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("signInSessionsValidFromDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSkills gets the skills property value. A list for the user to enumerate their skills. Returned only on $select.
// returns a []string when successful
func (m *User) GetSkills()([]string) {
    val, err := m.GetBackingStore().Get("skills")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSolutions gets the solutions property value. The identifier that relates the user to the working time schedule triggers. Read-Only. Nullable
// returns a UserSolutionRootable when successful
func (m *User) GetSolutions()(UserSolutionRootable) {
    val, err := m.GetBackingStore().Get("solutions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserSolutionRootable)
    }
    return nil
}
// GetSponsors gets the sponsors property value. The users and groups responsible for this guest's privileges in the tenant and keeping the guest's information and access updated. (HTTP Methods: GET, POST, DELETE.). Supports $expand.
// returns a []DirectoryObjectable when successful
func (m *User) GetSponsors()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("sponsors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetState gets the state property value. The state or province in the user's address. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetState()(*string) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStreetAddress gets the streetAddress property value. The street address of the user's place of business. Maximum length is 1,024 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetStreetAddress()(*string) {
    val, err := m.GetBackingStore().Get("streetAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSurname gets the surname property value. The user's surname (family name or last name). Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetSurname()(*string) {
    val, err := m.GetBackingStore().Get("surname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTeamwork gets the teamwork property value. A container for Microsoft Teams features available for the user. Read-only. Nullable.
// returns a UserTeamworkable when successful
func (m *User) GetTeamwork()(UserTeamworkable) {
    val, err := m.GetBackingStore().Get("teamwork")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserTeamworkable)
    }
    return nil
}
// GetTodo gets the todo property value. Represents the To Do services available to a user.
// returns a Todoable when successful
func (m *User) GetTodo()(Todoable) {
    val, err := m.GetBackingStore().Get("todo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Todoable)
    }
    return nil
}
// GetTransitiveMemberOf gets the transitiveMemberOf property value. The groups, including nested groups, and directory roles that a user is a member of. Nullable.
// returns a []DirectoryObjectable when successful
func (m *User) GetTransitiveMemberOf()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("transitiveMemberOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetUsageLocation gets the usageLocation property value. A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
// returns a *string when successful
func (m *User) GetUsageLocation()(*string) {
    val, err := m.GetBackingStore().Get("usageLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this value should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.
// returns a *string when successful
func (m *User) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserType gets the userType property value. A string value that can be used to classify user types in your directory. The possible values are Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for members and guests, see What are the default user permissions in Microsoft Entra ID?
// returns a *string when successful
func (m *User) GetUserType()(*string) {
    val, err := m.GetBackingStore().Get("userType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *User) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("aboutMe", m.GetAboutMe())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("accountEnabled", m.GetAccountEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetActivities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetActivities()))
        for i, v := range m.GetActivities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("activities", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ageGroup", m.GetAgeGroup())
        if err != nil {
            return err
        }
    }
    if m.GetAgreementAcceptances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAgreementAcceptances()))
        for i, v := range m.GetAgreementAcceptances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("agreementAcceptances", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppRoleAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppRoleAssignments()))
        for i, v := range m.GetAppRoleAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appRoleAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignedLicenses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignedLicenses()))
        for i, v := range m.GetAssignedLicenses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignedLicenses", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignedPlans() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignedPlans()))
        for i, v := range m.GetAssignedPlans() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignedPlans", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("authentication", m.GetAuthentication())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("authorizationInfo", m.GetAuthorizationInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("birthday", m.GetBirthday())
        if err != nil {
            return err
        }
    }
    if m.GetBusinessPhones() != nil {
        err = writer.WriteCollectionOfStringValues("businessPhones", m.GetBusinessPhones())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("calendar", m.GetCalendar())
        if err != nil {
            return err
        }
    }
    if m.GetCalendarGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCalendarGroups()))
        for i, v := range m.GetCalendarGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("calendarGroups", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCalendars() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCalendars()))
        for i, v := range m.GetCalendars() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("calendars", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCalendarView() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCalendarView()))
        for i, v := range m.GetCalendarView() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("calendarView", cast)
        if err != nil {
            return err
        }
    }
    if m.GetChats() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChats()))
        for i, v := range m.GetChats() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("chats", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("city", m.GetCity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("cloudClipboard", m.GetCloudClipboard())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("companyName", m.GetCompanyName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("consentProvidedForMinor", m.GetConsentProvidedForMinor())
        if err != nil {
            return err
        }
    }
    if m.GetContactFolders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContactFolders()))
        for i, v := range m.GetContactFolders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("contactFolders", cast)
        if err != nil {
            return err
        }
    }
    if m.GetContacts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContacts()))
        for i, v := range m.GetContacts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("contacts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("country", m.GetCountry())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetCreatedObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCreatedObjects()))
        for i, v := range m.GetCreatedObjects() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("createdObjects", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("creationType", m.GetCreationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("customSecurityAttributes", m.GetCustomSecurityAttributes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("department", m.GetDepartment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deviceEnrollmentLimit", m.GetDeviceEnrollmentLimit())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceManagementTroubleshootingEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceManagementTroubleshootingEvents()))
        for i, v := range m.GetDeviceManagementTroubleshootingEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceManagementTroubleshootingEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDirectReports() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDirectReports()))
        for i, v := range m.GetDirectReports() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("directReports", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("drive", m.GetDrive())
        if err != nil {
            return err
        }
    }
    if m.GetDrives() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDrives()))
        for i, v := range m.GetDrives() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("drives", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("employeeExperience", m.GetEmployeeExperience())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("employeeHireDate", m.GetEmployeeHireDate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("employeeId", m.GetEmployeeId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("employeeLeaveDateTime", m.GetEmployeeLeaveDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("employeeOrgData", m.GetEmployeeOrgData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("employeeType", m.GetEmployeeType())
        if err != nil {
            return err
        }
    }
    if m.GetEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEvents()))
        for i, v := range m.GetEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("events", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensions()))
        for i, v := range m.GetExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("extensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalUserState", m.GetExternalUserState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("externalUserStateChangeDateTime", m.GetExternalUserStateChangeDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("faxNumber", m.GetFaxNumber())
        if err != nil {
            return err
        }
    }
    if m.GetFollowedSites() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFollowedSites()))
        for i, v := range m.GetFollowedSites() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("followedSites", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("givenName", m.GetGivenName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("hireDate", m.GetHireDate())
        if err != nil {
            return err
        }
    }
    if m.GetIdentities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIdentities()))
        for i, v := range m.GetIdentities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("identities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetImAddresses() != nil {
        err = writer.WriteCollectionOfStringValues("imAddresses", m.GetImAddresses())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("inferenceClassification", m.GetInferenceClassification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("insights", m.GetInsights())
        if err != nil {
            return err
        }
    }
    if m.GetInterests() != nil {
        err = writer.WriteCollectionOfStringValues("interests", m.GetInterests())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isManagementRestricted", m.GetIsManagementRestricted())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isResourceAccount", m.GetIsResourceAccount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("jobTitle", m.GetJobTitle())
        if err != nil {
            return err
        }
    }
    if m.GetJoinedTeams() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetJoinedTeams()))
        for i, v := range m.GetJoinedTeams() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("joinedTeams", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastPasswordChangeDateTime", m.GetLastPasswordChangeDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("legalAgeGroupClassification", m.GetLegalAgeGroupClassification())
        if err != nil {
            return err
        }
    }
    if m.GetLicenseAssignmentStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLicenseAssignmentStates()))
        for i, v := range m.GetLicenseAssignmentStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("licenseAssignmentStates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLicenseDetails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLicenseDetails()))
        for i, v := range m.GetLicenseDetails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("licenseDetails", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mail", m.GetMail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mailboxSettings", m.GetMailboxSettings())
        if err != nil {
            return err
        }
    }
    if m.GetMailFolders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMailFolders()))
        for i, v := range m.GetMailFolders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mailFolders", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mailNickname", m.GetMailNickname())
        if err != nil {
            return err
        }
    }
    if m.GetManagedAppRegistrations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedAppRegistrations()))
        for i, v := range m.GetManagedAppRegistrations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedAppRegistrations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetManagedDevices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedDevices()))
        for i, v := range m.GetManagedDevices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedDevices", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("manager", m.GetManager())
        if err != nil {
            return err
        }
    }
    if m.GetMemberOf() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMemberOf()))
        for i, v := range m.GetMemberOf() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("memberOf", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessages()))
        for i, v := range m.GetMessages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messages", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mobilePhone", m.GetMobilePhone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mySite", m.GetMySite())
        if err != nil {
            return err
        }
    }
    if m.GetOauth2PermissionGrants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOauth2PermissionGrants()))
        for i, v := range m.GetOauth2PermissionGrants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("oauth2PermissionGrants", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("officeLocation", m.GetOfficeLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onenote", m.GetOnenote())
        if err != nil {
            return err
        }
    }
    if m.GetOnlineMeetings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOnlineMeetings()))
        for i, v := range m.GetOnlineMeetings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("onlineMeetings", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesDistinguishedName", m.GetOnPremisesDistinguishedName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesDomainName", m.GetOnPremisesDomainName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onPremisesExtensionAttributes", m.GetOnPremisesExtensionAttributes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesImmutableId", m.GetOnPremisesImmutableId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("onPremisesLastSyncDateTime", m.GetOnPremisesLastSyncDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetOnPremisesProvisioningErrors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOnPremisesProvisioningErrors()))
        for i, v := range m.GetOnPremisesProvisioningErrors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("onPremisesProvisioningErrors", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesSamAccountName", m.GetOnPremisesSamAccountName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesSecurityIdentifier", m.GetOnPremisesSecurityIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("onPremisesSyncEnabled", m.GetOnPremisesSyncEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onPremisesUserPrincipalName", m.GetOnPremisesUserPrincipalName())
        if err != nil {
            return err
        }
    }
    if m.GetOtherMails() != nil {
        err = writer.WriteCollectionOfStringValues("otherMails", m.GetOtherMails())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("outlook", m.GetOutlook())
        if err != nil {
            return err
        }
    }
    if m.GetOwnedDevices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOwnedDevices()))
        for i, v := range m.GetOwnedDevices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ownedDevices", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOwnedObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOwnedObjects()))
        for i, v := range m.GetOwnedObjects() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ownedObjects", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("passwordPolicies", m.GetPasswordPolicies())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("passwordProfile", m.GetPasswordProfile())
        if err != nil {
            return err
        }
    }
    if m.GetPastProjects() != nil {
        err = writer.WriteCollectionOfStringValues("pastProjects", m.GetPastProjects())
        if err != nil {
            return err
        }
    }
    if m.GetPeople() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPeople()))
        for i, v := range m.GetPeople() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("people", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPermissionGrants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPermissionGrants()))
        for i, v := range m.GetPermissionGrants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("permissionGrants", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    if m.GetPhotos() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPhotos()))
        for i, v := range m.GetPhotos() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("photos", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("planner", m.GetPlanner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("postalCode", m.GetPostalCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preferredDataLocation", m.GetPreferredDataLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preferredLanguage", m.GetPreferredLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preferredName", m.GetPreferredName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("presence", m.GetPresence())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("print", m.GetPrint())
        if err != nil {
            return err
        }
    }
    if m.GetProvisionedPlans() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProvisionedPlans()))
        for i, v := range m.GetProvisionedPlans() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("provisionedPlans", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProxyAddresses() != nil {
        err = writer.WriteCollectionOfStringValues("proxyAddresses", m.GetProxyAddresses())
        if err != nil {
            return err
        }
    }
    if m.GetRegisteredDevices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRegisteredDevices()))
        for i, v := range m.GetRegisteredDevices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("registeredDevices", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResponsibilities() != nil {
        err = writer.WriteCollectionOfStringValues("responsibilities", m.GetResponsibilities())
        if err != nil {
            return err
        }
    }
    if m.GetSchools() != nil {
        err = writer.WriteCollectionOfStringValues("schools", m.GetSchools())
        if err != nil {
            return err
        }
    }
    if m.GetScopedRoleMemberOf() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScopedRoleMemberOf()))
        for i, v := range m.GetScopedRoleMemberOf() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scopedRoleMemberOf", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("securityIdentifier", m.GetSecurityIdentifier())
        if err != nil {
            return err
        }
    }
    if m.GetServiceProvisioningErrors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServiceProvisioningErrors()))
        for i, v := range m.GetServiceProvisioningErrors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("serviceProvisioningErrors", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showInAddressList", m.GetShowInAddressList())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("signInActivity", m.GetSignInActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("signInSessionsValidFromDateTime", m.GetSignInSessionsValidFromDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetSkills() != nil {
        err = writer.WriteCollectionOfStringValues("skills", m.GetSkills())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("solutions", m.GetSolutions())
        if err != nil {
            return err
        }
    }
    if m.GetSponsors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSponsors()))
        for i, v := range m.GetSponsors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sponsors", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("streetAddress", m.GetStreetAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("surname", m.GetSurname())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("teamwork", m.GetTeamwork())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("todo", m.GetTodo())
        if err != nil {
            return err
        }
    }
    if m.GetTransitiveMemberOf() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTransitiveMemberOf()))
        for i, v := range m.GetTransitiveMemberOf() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("transitiveMemberOf", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("usageLocation", m.GetUsageLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userType", m.GetUserType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAboutMe sets the aboutMe property value. A freeform text entry field for the user to describe themselves. Returned only on $select.
func (m *User) SetAboutMe(value *string)() {
    err := m.GetBackingStore().Set("aboutMe", value)
    if err != nil {
        panic(err)
    }
}
// SetAccountEnabled sets the accountEnabled property value. true if the account is enabled; otherwise, false. This property is required when a user is created. Returned only on $select. Supports $filter (eq, ne, not, and in).
func (m *User) SetAccountEnabled(value *bool)() {
    err := m.GetBackingStore().Set("accountEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetActivities sets the activities property value. The user's activities across devices. Read-only. Nullable.
func (m *User) SetActivities(value []UserActivityable)() {
    err := m.GetBackingStore().Set("activities", value)
    if err != nil {
        panic(err)
    }
}
// SetAgeGroup sets the ageGroup property value. Sets the age group of the user. Allowed values: null, Minor, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).
func (m *User) SetAgeGroup(value *string)() {
    err := m.GetBackingStore().Set("ageGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetAgreementAcceptances sets the agreementAcceptances property value. The user's terms of use acceptance statuses. Read-only. Nullable.
func (m *User) SetAgreementAcceptances(value []AgreementAcceptanceable)() {
    err := m.GetBackingStore().Set("agreementAcceptances", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoleAssignments sets the appRoleAssignments property value. Represents the app roles a user is granted for an application. Supports $expand.
func (m *User) SetAppRoleAssignments(value []AppRoleAssignmentable)() {
    err := m.GetBackingStore().Set("appRoleAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedLicenses sets the assignedLicenses property value. The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate between directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses. Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).
func (m *User) SetAssignedLicenses(value []AssignedLicenseable)() {
    err := m.GetBackingStore().Set("assignedLicenses", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedPlans sets the assignedPlans property value. The plans that are assigned to the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq and not).
func (m *User) SetAssignedPlans(value []AssignedPlanable)() {
    err := m.GetBackingStore().Set("assignedPlans", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthentication sets the authentication property value. The authentication methods that are supported for the user.
func (m *User) SetAuthentication(value Authenticationable)() {
    err := m.GetBackingStore().Set("authentication", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthorizationInfo sets the authorizationInfo property value. The authorizationInfo property
func (m *User) SetAuthorizationInfo(value AuthorizationInfoable)() {
    err := m.GetBackingStore().Set("authorizationInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetBirthday sets the birthday property value. The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.
func (m *User) SetBirthday(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("birthday", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessPhones sets the businessPhones property value. The telephone numbers for the user. NOTE: Although it's a string collection, only one number can be set for this property. Read-only for users synced from the on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).
func (m *User) SetBusinessPhones(value []string)() {
    err := m.GetBackingStore().Set("businessPhones", value)
    if err != nil {
        panic(err)
    }
}
// SetCalendar sets the calendar property value. The user's primary calendar. Read-only.
func (m *User) SetCalendar(value Calendarable)() {
    err := m.GetBackingStore().Set("calendar", value)
    if err != nil {
        panic(err)
    }
}
// SetCalendarGroups sets the calendarGroups property value. The user's calendar groups. Read-only. Nullable.
func (m *User) SetCalendarGroups(value []CalendarGroupable)() {
    err := m.GetBackingStore().Set("calendarGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetCalendars sets the calendars property value. The user's calendars. Read-only. Nullable.
func (m *User) SetCalendars(value []Calendarable)() {
    err := m.GetBackingStore().Set("calendars", value)
    if err != nil {
        panic(err)
    }
}
// SetCalendarView sets the calendarView property value. The calendar view for the calendar. Read-only. Nullable.
func (m *User) SetCalendarView(value []Eventable)() {
    err := m.GetBackingStore().Set("calendarView", value)
    if err != nil {
        panic(err)
    }
}
// SetChats sets the chats property value. The chats property
func (m *User) SetChats(value []Chatable)() {
    err := m.GetBackingStore().Set("chats", value)
    if err != nil {
        panic(err)
    }
}
// SetCity sets the city property value. The city where the user is located. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetCity(value *string)() {
    err := m.GetBackingStore().Set("city", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudClipboard sets the cloudClipboard property value. The cloudClipboard property
func (m *User) SetCloudClipboard(value CloudClipboardRootable)() {
    err := m.GetBackingStore().Set("cloudClipboard", value)
    if err != nil {
        panic(err)
    }
}
// SetCompanyName sets the companyName property value. The name of the company that the user is associated with. This property can be useful for describing the company that a guest comes from. The maximum length is 64 characters.Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetCompanyName(value *string)() {
    err := m.GetBackingStore().Set("companyName", value)
    if err != nil {
        panic(err)
    }
}
// SetConsentProvidedForMinor sets the consentProvidedForMinor property value. Sets whether consent was obtained for minors. Allowed values: null, Granted, Denied, and NotRequired. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).
func (m *User) SetConsentProvidedForMinor(value *string)() {
    err := m.GetBackingStore().Set("consentProvidedForMinor", value)
    if err != nil {
        panic(err)
    }
}
// SetContactFolders sets the contactFolders property value. The user's contacts folders. Read-only. Nullable.
func (m *User) SetContactFolders(value []ContactFolderable)() {
    err := m.GetBackingStore().Set("contactFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetContacts sets the contacts property value. The user's contacts. Read-only. Nullable.
func (m *User) SetContacts(value []Contactable)() {
    err := m.GetBackingStore().Set("contacts", value)
    if err != nil {
        panic(err)
    }
}
// SetCountry sets the country property value. The country/region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetCountry(value *string)() {
    err := m.GetBackingStore().Set("country", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time the user was created, in ISO 8601 format and UTC. The value can't be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
func (m *User) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedObjects sets the createdObjects property value. Directory objects that the user created. Read-only. Nullable.
func (m *User) SetCreatedObjects(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("createdObjects", value)
    if err != nil {
        panic(err)
    }
}
// SetCreationType sets the creationType property value. Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by a guest signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).
func (m *User) SetCreationType(value *string)() {
    err := m.GetBackingStore().Set("creationType", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomSecurityAttributes sets the customSecurityAttributes property value. An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). The filter value is case-sensitive. To read this property, the calling app must be assigned the CustomSecAttributeAssignment.Read.All permission. To write this property, the calling app must be assigned the CustomSecAttributeAssignment.ReadWrite.All permissions. To read or write this property in delegated scenarios, the admin must be assigned the Attribute Assignment Administrator role.
func (m *User) SetCustomSecurityAttributes(value CustomSecurityAttributeValueable)() {
    err := m.GetBackingStore().Set("customSecurityAttributes", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartment sets the department property value. The name of the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).
func (m *User) SetDepartment(value *string)() {
    err := m.GetBackingStore().Set("department", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceEnrollmentLimit sets the deviceEnrollmentLimit property value. The limit on the maximum number of devices that the user is permitted to enroll. Allowed values are 5 or 1000.
func (m *User) SetDeviceEnrollmentLimit(value *int32)() {
    err := m.GetBackingStore().Set("deviceEnrollmentLimit", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceManagementTroubleshootingEvents sets the deviceManagementTroubleshootingEvents property value. The list of troubleshooting events for this user.
func (m *User) SetDeviceManagementTroubleshootingEvents(value []DeviceManagementTroubleshootingEventable)() {
    err := m.GetBackingStore().Set("deviceManagementTroubleshootingEvents", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectReports sets the directReports property value. The users and contacts that report to the user. (The users and contacts that have their manager property set to this user.) Read-only. Nullable. Supports $expand.
func (m *User) SetDirectReports(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("directReports", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name displayed in the address book for the user. This value is usually the combination of the user's first name, middle initial, and family name. This property is required when a user is created and it can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.
func (m *User) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDrive sets the drive property value. The user's OneDrive. Read-only.
func (m *User) SetDrive(value Driveable)() {
    err := m.GetBackingStore().Set("drive", value)
    if err != nil {
        panic(err)
    }
}
// SetDrives sets the drives property value. A collection of drives available for this user. Read-only.
func (m *User) SetDrives(value []Driveable)() {
    err := m.GetBackingStore().Set("drives", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeExperience sets the employeeExperience property value. The employeeExperience property
func (m *User) SetEmployeeExperience(value EmployeeExperienceUserable)() {
    err := m.GetBackingStore().Set("employeeExperience", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeHireDate sets the employeeHireDate property value. The date and time when the user was hired or will start work in a future hire. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
func (m *User) SetEmployeeHireDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("employeeHireDate", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeId sets the employeeId property value. The employee identifier assigned to the user by the organization. The maximum length is 16 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
func (m *User) SetEmployeeId(value *string)() {
    err := m.GetBackingStore().Set("employeeId", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeLeaveDateTime sets the employeeLeaveDateTime property value. The date and time when the user left or will leave the organization. To read this property, the calling app must be assigned the User-LifeCycleInfo.Read.All permission. To write this property, the calling app must be assigned the User.Read.All and User-LifeCycleInfo.ReadWrite.All permissions. To read this property in delegated scenarios, the admin needs at least one of the following Microsoft Entra roles: Lifecycle Workflows Administrator (least privilege), Global Reader. To write this property in delegated scenarios, the admin needs the Global Administrator role. Supports $filter (eq, ne, not , ge, le, in). For more information, see Configure the employeeLeaveDateTime property for a user.
func (m *User) SetEmployeeLeaveDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("employeeLeaveDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeOrgData sets the employeeOrgData property value. Represents organization data (for example, division and costCenter) associated with a user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).
func (m *User) SetEmployeeOrgData(value EmployeeOrgDataable)() {
    err := m.GetBackingStore().Set("employeeOrgData", value)
    if err != nil {
        panic(err)
    }
}
// SetEmployeeType sets the employeeType property value. Captures enterprise worker type. For example, Employee, Contractor, Consultant, or Vendor. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith).
func (m *User) SetEmployeeType(value *string)() {
    err := m.GetBackingStore().Set("employeeType", value)
    if err != nil {
        panic(err)
    }
}
// SetEvents sets the events property value. The user's events. Default is to show Events under the Default Calendar. Read-only. Nullable.
func (m *User) SetEvents(value []Eventable)() {
    err := m.GetBackingStore().Set("events", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the user. Read-only. Supports $expand. Nullable.
func (m *User) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalUserState sets the externalUserState property value. For a guest invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be PendingAcceptance or Accepted, or null for all other users. Returned only on $select. Supports $filter (eq, ne, not , in).
func (m *User) SetExternalUserState(value *string)() {
    err := m.GetBackingStore().Set("externalUserState", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalUserStateChangeDateTime sets the externalUserStateChangeDateTime property value. Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. Supports $filter (eq, ne, not , in).
func (m *User) SetExternalUserStateChangeDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("externalUserStateChangeDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFaxNumber sets the faxNumber property value. The fax number of the user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
func (m *User) SetFaxNumber(value *string)() {
    err := m.GetBackingStore().Set("faxNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetFollowedSites sets the followedSites property value. The followedSites property
func (m *User) SetFollowedSites(value []Siteable)() {
    err := m.GetBackingStore().Set("followedSites", value)
    if err != nil {
        panic(err)
    }
}
// SetGivenName sets the givenName property value. The given name (first name) of the user. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
func (m *User) SetGivenName(value *string)() {
    err := m.GetBackingStore().Set("givenName", value)
    if err != nil {
        panic(err)
    }
}
// SetHireDate sets the hireDate property value. The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.  Note: This property is specific to SharePoint in Microsoft 365. We recommend using the native employeeHireDate property to set and update hire date values using Microsoft Graph APIs.
func (m *User) SetHireDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("hireDate", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentities sets the identities property value. Represents the identities that can be used to sign in to this user account. Microsoft (also known as a local account), organizations, or social identity providers such as Facebook, Google, and Microsoft can provide identity and tie it to a user account. It might contain multiple items with the same signInType value. Returned only on $select.  Supports $filter (eq) with limitations.
func (m *User) SetIdentities(value []ObjectIdentityable)() {
    err := m.GetBackingStore().Set("identities", value)
    if err != nil {
        panic(err)
    }
}
// SetImAddresses sets the imAddresses property value. The instant message voice-over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).
func (m *User) SetImAddresses(value []string)() {
    err := m.GetBackingStore().Set("imAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetInferenceClassification sets the inferenceClassification property value. Relevance classification of the user's messages based on explicit designations that override inferred relevance or importance.
func (m *User) SetInferenceClassification(value InferenceClassificationable)() {
    err := m.GetBackingStore().Set("inferenceClassification", value)
    if err != nil {
        panic(err)
    }
}
// SetInsights sets the insights property value. Represents relationships between a user and items such as OneDrive for work or school documents, calculated using advanced analytics and machine learning techniques. Read-only. Nullable.
func (m *User) SetInsights(value ItemInsightsable)() {
    err := m.GetBackingStore().Set("insights", value)
    if err != nil {
        panic(err)
    }
}
// SetInterests sets the interests property value. A list for the user to describe their interests. Returned only on $select.
func (m *User) SetInterests(value []string)() {
    err := m.GetBackingStore().Set("interests", value)
    if err != nil {
        panic(err)
    }
}
// SetIsManagementRestricted sets the isManagementRestricted property value. The isManagementRestricted property
func (m *User) SetIsManagementRestricted(value *bool)() {
    err := m.GetBackingStore().Set("isManagementRestricted", value)
    if err != nil {
        panic(err)
    }
}
// SetIsResourceAccount sets the isResourceAccount property value. Don't use – reserved for future use.
func (m *User) SetIsResourceAccount(value *bool)() {
    err := m.GetBackingStore().Set("isResourceAccount", value)
    if err != nil {
        panic(err)
    }
}
// SetJobTitle sets the jobTitle property value. The user's job title. Maximum length is 128 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).
func (m *User) SetJobTitle(value *string)() {
    err := m.GetBackingStore().Set("jobTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetJoinedTeams sets the joinedTeams property value. The joinedTeams property
func (m *User) SetJoinedTeams(value []Teamable)() {
    err := m.GetBackingStore().Set("joinedTeams", value)
    if err != nil {
        panic(err)
    }
}
// SetLastPasswordChangeDateTime sets the lastPasswordChangeDateTime property value. The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.
func (m *User) SetLastPasswordChangeDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastPasswordChangeDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLegalAgeGroupClassification sets the legalAgeGroupClassification property value. Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select.
func (m *User) SetLegalAgeGroupClassification(value *string)() {
    err := m.GetBackingStore().Set("legalAgeGroupClassification", value)
    if err != nil {
        panic(err)
    }
}
// SetLicenseAssignmentStates sets the licenseAssignmentStates property value. State of license assignments for this user. Also indicates licenses that are directly assigned or the user inherited through group memberships. Read-only. Returned only on $select.
func (m *User) SetLicenseAssignmentStates(value []LicenseAssignmentStateable)() {
    err := m.GetBackingStore().Set("licenseAssignmentStates", value)
    if err != nil {
        panic(err)
    }
}
// SetLicenseDetails sets the licenseDetails property value. A collection of this user's license details. Read-only.
func (m *User) SetLicenseDetails(value []LicenseDetailsable)() {
    err := m.GetBackingStore().Set("licenseDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetMail sets the mail property value. The SMTP address for the user, for example, jeff@contoso.com. Changes to this property update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).
func (m *User) SetMail(value *string)() {
    err := m.GetBackingStore().Set("mail", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxSettings sets the mailboxSettings property value. Settings for the primary mailbox of the signed-in user. You can get or update settings for sending automatic replies to incoming messages, locale, and time zone. Returned only on $select.
func (m *User) SetMailboxSettings(value MailboxSettingsable)() {
    err := m.GetBackingStore().Set("mailboxSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetMailFolders sets the mailFolders property value. The user's mail folders. Read-only. Nullable.
func (m *User) SetMailFolders(value []MailFolderable)() {
    err := m.GetBackingStore().Set("mailFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetMailNickname sets the mailNickname property value. The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetMailNickname(value *string)() {
    err := m.GetBackingStore().Set("mailNickname", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedAppRegistrations sets the managedAppRegistrations property value. Zero or more managed app registrations that belong to the user.
func (m *User) SetManagedAppRegistrations(value []ManagedAppRegistrationable)() {
    err := m.GetBackingStore().Set("managedAppRegistrations", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDevices sets the managedDevices property value. The managed devices associated with the user.
func (m *User) SetManagedDevices(value []ManagedDeviceable)() {
    err := m.GetBackingStore().Set("managedDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetManager sets the manager property value. The user or contact that is this user's manager. Read-only. Supports $expand.
func (m *User) SetManager(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("manager", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberOf sets the memberOf property value. The groups and directory roles that the user is a member of. Read-only. Nullable. Supports $expand.
func (m *User) SetMemberOf(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("memberOf", value)
    if err != nil {
        panic(err)
    }
}
// SetMessages sets the messages property value. The messages in a mailbox or folder. Read-only. Nullable.
func (m *User) SetMessages(value []Messageable)() {
    err := m.GetBackingStore().Set("messages", value)
    if err != nil {
        panic(err)
    }
}
// SetMobilePhone sets the mobilePhone property value. The primary cellular telephone number for the user. Read-only for users synced from the on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.
func (m *User) SetMobilePhone(value *string)() {
    err := m.GetBackingStore().Set("mobilePhone", value)
    if err != nil {
        panic(err)
    }
}
// SetMySite sets the mySite property value. The URL for the user's site. Returned only on $select.
func (m *User) SetMySite(value *string)() {
    err := m.GetBackingStore().Set("mySite", value)
    if err != nil {
        panic(err)
    }
}
// SetOauth2PermissionGrants sets the oauth2PermissionGrants property value. The oauth2PermissionGrants property
func (m *User) SetOauth2PermissionGrants(value []OAuth2PermissionGrantable)() {
    err := m.GetBackingStore().Set("oauth2PermissionGrants", value)
    if err != nil {
        panic(err)
    }
}
// SetOfficeLocation sets the officeLocation property value. The office location in the user's place of business. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetOfficeLocation(value *string)() {
    err := m.GetBackingStore().Set("officeLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetOnenote sets the onenote property value. The onenote property
func (m *User) SetOnenote(value Onenoteable)() {
    err := m.GetBackingStore().Set("onenote", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeetings sets the onlineMeetings property value. Information about a meeting, including the URL used to join a meeting, the attendees list, and the description.
func (m *User) SetOnlineMeetings(value []OnlineMeetingable)() {
    err := m.GetBackingStore().Set("onlineMeetings", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesDistinguishedName sets the onPremisesDistinguishedName property value. Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.
func (m *User) SetOnPremisesDistinguishedName(value *string)() {
    err := m.GetBackingStore().Set("onPremisesDistinguishedName", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesDomainName sets the onPremisesDomainName property value. Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.
func (m *User) SetOnPremisesDomainName(value *string)() {
    err := m.GetBackingStore().Set("onPremisesDomainName", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesExtensionAttributes sets the onPremisesExtensionAttributes property value. Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. Each attribute can store up to 1024 characters. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during the creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).
func (m *User) SetOnPremisesExtensionAttributes(value OnPremisesExtensionAttributesable)() {
    err := m.GetBackingStore().Set("onPremisesExtensionAttributes", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesImmutableId sets the onPremisesImmutableId property value. This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property. NOTE: The $ and _ characters can't be used when specifying this property. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).
func (m *User) SetOnPremisesImmutableId(value *string)() {
    err := m.GetBackingStore().Set("onPremisesImmutableId", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesLastSyncDateTime sets the onPremisesLastSyncDateTime property value. Indicates the last time at which the object was synced with the on-premises directory; for example: 2013-02-16T03:04:54Z. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).
func (m *User) SetOnPremisesLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("onPremisesLastSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesProvisioningErrors sets the onPremisesProvisioningErrors property value. Errors when using Microsoft synchronization product during provisioning. Returned only on $select. Supports $filter (eq, not, ge, le).
func (m *User) SetOnPremisesProvisioningErrors(value []OnPremisesProvisioningErrorable)() {
    err := m.GetBackingStore().Set("onPremisesProvisioningErrors", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesSamAccountName sets the onPremisesSamAccountName property value. Contains the on-premises samAccountName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).
func (m *User) SetOnPremisesSamAccountName(value *string)() {
    err := m.GetBackingStore().Set("onPremisesSamAccountName", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesSecurityIdentifier sets the onPremisesSecurityIdentifier property value. Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq including on null values).
func (m *User) SetOnPremisesSecurityIdentifier(value *string)() {
    err := m.GetBackingStore().Set("onPremisesSecurityIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesSyncEnabled sets the onPremisesSyncEnabled property value. true if this user object is currently being synced from an on-premises Active Directory (AD); otherwise the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values).
func (m *User) SetOnPremisesSyncEnabled(value *bool)() {
    err := m.GetBackingStore().Set("onPremisesSyncEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesUserPrincipalName sets the onPremisesUserPrincipalName property value. Contains the on-premises userPrincipalName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).
func (m *User) SetOnPremisesUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("onPremisesUserPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetOtherMails sets the otherMails property value. A list of other email addresses for the user; for example: ['bob@contoso.com', 'Robert@fabrikam.com']. NOTE: This property can't contain accent characters. Returned only on $select. Supports $filter (eq, not, ge, le, in, startsWith, endsWith, /$count eq 0, /$count ne 0).
func (m *User) SetOtherMails(value []string)() {
    err := m.GetBackingStore().Set("otherMails", value)
    if err != nil {
        panic(err)
    }
}
// SetOutlook sets the outlook property value. The outlook property
func (m *User) SetOutlook(value OutlookUserable)() {
    err := m.GetBackingStore().Set("outlook", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnedDevices sets the ownedDevices property value. Devices the user owns. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *User) SetOwnedDevices(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("ownedDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnedObjects sets the ownedObjects property value. Directory objects the user owns. Read-only. Nullable. Supports $expand, $select nested in $expand, and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *User) SetOwnedObjects(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("ownedObjects", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPolicies sets the passwordPolicies property value. Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two might be specified together; for example: DisablePasswordExpiration, DisableStrongPassword. Returned only on $select. For more information on the default password policies, see Microsoft Entra password policies. Supports $filter (ne, not, and eq on null values).
func (m *User) SetPasswordPolicies(value *string)() {
    err := m.GetBackingStore().Set("passwordPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordProfile sets the passwordProfile property value. Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). To update this property:  In delegated access, the calling app must be assigned the Directory.AccessAsUser.All delegated permission on behalf of the signed-in user.  In application-only access, the calling app must be assigned the User.ReadWrite.All (least privilege) or Directory.ReadWrite.All (higher privilege) application permission and at least the User Administrator Microsoft Entra role.
func (m *User) SetPasswordProfile(value PasswordProfileable)() {
    err := m.GetBackingStore().Set("passwordProfile", value)
    if err != nil {
        panic(err)
    }
}
// SetPastProjects sets the pastProjects property value. A list for the user to enumerate their past projects. Returned only on $select.
func (m *User) SetPastProjects(value []string)() {
    err := m.GetBackingStore().Set("pastProjects", value)
    if err != nil {
        panic(err)
    }
}
// SetPeople sets the people property value. People that are relevant to the user. Read-only. Nullable.
func (m *User) SetPeople(value []Personable)() {
    err := m.GetBackingStore().Set("people", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionGrants sets the permissionGrants property value. List all resource-specific permission grants of a user.
func (m *User) SetPermissionGrants(value []ResourceSpecificPermissionGrantable)() {
    err := m.GetBackingStore().Set("permissionGrants", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoto sets the photo property value. The user's profile photo. Read-only.
func (m *User) SetPhoto(value ProfilePhotoable)() {
    err := m.GetBackingStore().Set("photo", value)
    if err != nil {
        panic(err)
    }
}
// SetPhotos sets the photos property value. The collection of the user's profile photos in different sizes. Read-only.
func (m *User) SetPhotos(value []ProfilePhotoable)() {
    err := m.GetBackingStore().Set("photos", value)
    if err != nil {
        panic(err)
    }
}
// SetPlanner sets the planner property value. Entry-point to the Planner resource that might exist for a user. Read-only.
func (m *User) SetPlanner(value PlannerUserable)() {
    err := m.GetBackingStore().Set("planner", value)
    if err != nil {
        panic(err)
    }
}
// SetPostalCode sets the postalCode property value. The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetPostalCode(value *string)() {
    err := m.GetBackingStore().Set("postalCode", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredDataLocation sets the preferredDataLocation property value. The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.
func (m *User) SetPreferredDataLocation(value *string)() {
    err := m.GetBackingStore().Set("preferredDataLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredLanguage sets the preferredLanguage property value. The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language, and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)
func (m *User) SetPreferredLanguage(value *string)() {
    err := m.GetBackingStore().Set("preferredLanguage", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredName sets the preferredName property value. The preferred name for the user. Not Supported. This attribute returns an empty string.Returned only on $select.
func (m *User) SetPreferredName(value *string)() {
    err := m.GetBackingStore().Set("preferredName", value)
    if err != nil {
        panic(err)
    }
}
// SetPresence sets the presence property value. The presence property
func (m *User) SetPresence(value Presenceable)() {
    err := m.GetBackingStore().Set("presence", value)
    if err != nil {
        panic(err)
    }
}
// SetPrint sets the print property value. The print property
func (m *User) SetPrint(value UserPrintable)() {
    err := m.GetBackingStore().Set("print", value)
    if err != nil {
        panic(err)
    }
}
// SetProvisionedPlans sets the provisionedPlans property value. The plans that are provisioned for the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le).
func (m *User) SetProvisionedPlans(value []ProvisionedPlanable)() {
    err := m.GetBackingStore().Set("provisionedPlans", value)
    if err != nil {
        panic(err)
    }
}
// SetProxyAddresses sets the proxyAddresses property value. For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address, while those addresses prefixed with smtp are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of 10 unique addresses. Read-only in Microsoft Graph; you can update this property only through the Microsoft 365 admin center. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).
func (m *User) SetProxyAddresses(value []string)() {
    err := m.GetBackingStore().Set("proxyAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetRegisteredDevices sets the registeredDevices property value. Devices that are registered for the user. Read-only. Nullable. Supports $expand and returns up to 100 objects.
func (m *User) SetRegisteredDevices(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("registeredDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetResponsibilities sets the responsibilities property value. A list for the user to enumerate their responsibilities. Returned only on $select.
func (m *User) SetResponsibilities(value []string)() {
    err := m.GetBackingStore().Set("responsibilities", value)
    if err != nil {
        panic(err)
    }
}
// SetSchools sets the schools property value. A list for the user to enumerate the schools they attended. Returned only on $select.
func (m *User) SetSchools(value []string)() {
    err := m.GetBackingStore().Set("schools", value)
    if err != nil {
        panic(err)
    }
}
// SetScopedRoleMemberOf sets the scopedRoleMemberOf property value. The scopedRoleMemberOf property
func (m *User) SetScopedRoleMemberOf(value []ScopedRoleMembershipable)() {
    err := m.GetBackingStore().Set("scopedRoleMemberOf", value)
    if err != nil {
        panic(err)
    }
}
// SetSecurityIdentifier sets the securityIdentifier property value. Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).
func (m *User) SetSecurityIdentifier(value *string)() {
    err := m.GetBackingStore().Set("securityIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceProvisioningErrors sets the serviceProvisioningErrors property value. Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a user object.  Supports $filter (eq, not, for isResolved and serviceInstance).
func (m *User) SetServiceProvisioningErrors(value []ServiceProvisioningErrorable)() {
    err := m.GetBackingStore().Set("serviceProvisioningErrors", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. The settings property
func (m *User) SetSettings(value UserSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetShowInAddressList sets the showInAddressList property value. Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead. Represents whether the user should be included in the Outlook global address list. See Known issue.
func (m *User) SetShowInAddressList(value *bool)() {
    err := m.GetBackingStore().Set("showInAddressList", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInActivity sets the signInActivity property value. Get the last signed-in date and request ID of the sign-in for a given user. Read-only.Returned only on $select. Supports $filter (eq, ne, not, ge, le) but not with any other filterable properties. Note: Details for this property require a Microsoft Entra ID P1 or P2 license and the AuditLog.Read.All permission.This property isn't returned for a user who never signed in or last signed in before April 2020.
func (m *User) SetSignInActivity(value SignInActivityable)() {
    err := m.GetBackingStore().Set("signInActivity", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInSessionsValidFromDateTime sets the signInSessionsValidFromDateTime property value. Any refresh tokens or session tokens (session cookies) issued before this time are invalid. Applications get an error when using an invalid refresh or session token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.
func (m *User) SetSignInSessionsValidFromDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("signInSessionsValidFromDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSkills sets the skills property value. A list for the user to enumerate their skills. Returned only on $select.
func (m *User) SetSkills(value []string)() {
    err := m.GetBackingStore().Set("skills", value)
    if err != nil {
        panic(err)
    }
}
// SetSolutions sets the solutions property value. The identifier that relates the user to the working time schedule triggers. Read-Only. Nullable
func (m *User) SetSolutions(value UserSolutionRootable)() {
    err := m.GetBackingStore().Set("solutions", value)
    if err != nil {
        panic(err)
    }
}
// SetSponsors sets the sponsors property value. The users and groups responsible for this guest's privileges in the tenant and keeping the guest's information and access updated. (HTTP Methods: GET, POST, DELETE.). Supports $expand.
func (m *User) SetSponsors(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("sponsors", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state or province in the user's address. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetState(value *string)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetStreetAddress sets the streetAddress property value. The street address of the user's place of business. Maximum length is 1,024 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetStreetAddress(value *string)() {
    err := m.GetBackingStore().Set("streetAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetSurname sets the surname property value. The user's surname (family name or last name). Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetSurname(value *string)() {
    err := m.GetBackingStore().Set("surname", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamwork sets the teamwork property value. A container for Microsoft Teams features available for the user. Read-only. Nullable.
func (m *User) SetTeamwork(value UserTeamworkable)() {
    err := m.GetBackingStore().Set("teamwork", value)
    if err != nil {
        panic(err)
    }
}
// SetTodo sets the todo property value. Represents the To Do services available to a user.
func (m *User) SetTodo(value Todoable)() {
    err := m.GetBackingStore().Set("todo", value)
    if err != nil {
        panic(err)
    }
}
// SetTransitiveMemberOf sets the transitiveMemberOf property value. The groups, including nested groups, and directory roles that a user is a member of. Nullable.
func (m *User) SetTransitiveMemberOf(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("transitiveMemberOf", value)
    if err != nil {
        panic(err)
    }
}
// SetUsageLocation sets the usageLocation property value. A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).
func (m *User) SetUsageLocation(value *string)() {
    err := m.GetBackingStore().Set("usageLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this value should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.
func (m *User) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserType sets the userType property value. A string value that can be used to classify user types in your directory. The possible values are Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for members and guests, see What are the default user permissions in Microsoft Entra ID?
func (m *User) SetUserType(value *string)() {
    err := m.GetBackingStore().Set("userType", value)
    if err != nil {
        panic(err)
    }
}
type Userable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAboutMe()(*string)
    GetAccountEnabled()(*bool)
    GetActivities()([]UserActivityable)
    GetAgeGroup()(*string)
    GetAgreementAcceptances()([]AgreementAcceptanceable)
    GetAppRoleAssignments()([]AppRoleAssignmentable)
    GetAssignedLicenses()([]AssignedLicenseable)
    GetAssignedPlans()([]AssignedPlanable)
    GetAuthentication()(Authenticationable)
    GetAuthorizationInfo()(AuthorizationInfoable)
    GetBirthday()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetBusinessPhones()([]string)
    GetCalendar()(Calendarable)
    GetCalendarGroups()([]CalendarGroupable)
    GetCalendars()([]Calendarable)
    GetCalendarView()([]Eventable)
    GetChats()([]Chatable)
    GetCity()(*string)
    GetCloudClipboard()(CloudClipboardRootable)
    GetCompanyName()(*string)
    GetConsentProvidedForMinor()(*string)
    GetContactFolders()([]ContactFolderable)
    GetContacts()([]Contactable)
    GetCountry()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedObjects()([]DirectoryObjectable)
    GetCreationType()(*string)
    GetCustomSecurityAttributes()(CustomSecurityAttributeValueable)
    GetDepartment()(*string)
    GetDeviceEnrollmentLimit()(*int32)
    GetDeviceManagementTroubleshootingEvents()([]DeviceManagementTroubleshootingEventable)
    GetDirectReports()([]DirectoryObjectable)
    GetDisplayName()(*string)
    GetDrive()(Driveable)
    GetDrives()([]Driveable)
    GetEmployeeExperience()(EmployeeExperienceUserable)
    GetEmployeeHireDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEmployeeId()(*string)
    GetEmployeeLeaveDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEmployeeOrgData()(EmployeeOrgDataable)
    GetEmployeeType()(*string)
    GetEvents()([]Eventable)
    GetExtensions()([]Extensionable)
    GetExternalUserState()(*string)
    GetExternalUserStateChangeDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFaxNumber()(*string)
    GetFollowedSites()([]Siteable)
    GetGivenName()(*string)
    GetHireDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIdentities()([]ObjectIdentityable)
    GetImAddresses()([]string)
    GetInferenceClassification()(InferenceClassificationable)
    GetInsights()(ItemInsightsable)
    GetInterests()([]string)
    GetIsManagementRestricted()(*bool)
    GetIsResourceAccount()(*bool)
    GetJobTitle()(*string)
    GetJoinedTeams()([]Teamable)
    GetLastPasswordChangeDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLegalAgeGroupClassification()(*string)
    GetLicenseAssignmentStates()([]LicenseAssignmentStateable)
    GetLicenseDetails()([]LicenseDetailsable)
    GetMail()(*string)
    GetMailboxSettings()(MailboxSettingsable)
    GetMailFolders()([]MailFolderable)
    GetMailNickname()(*string)
    GetManagedAppRegistrations()([]ManagedAppRegistrationable)
    GetManagedDevices()([]ManagedDeviceable)
    GetManager()(DirectoryObjectable)
    GetMemberOf()([]DirectoryObjectable)
    GetMessages()([]Messageable)
    GetMobilePhone()(*string)
    GetMySite()(*string)
    GetOauth2PermissionGrants()([]OAuth2PermissionGrantable)
    GetOfficeLocation()(*string)
    GetOnenote()(Onenoteable)
    GetOnlineMeetings()([]OnlineMeetingable)
    GetOnPremisesDistinguishedName()(*string)
    GetOnPremisesDomainName()(*string)
    GetOnPremisesExtensionAttributes()(OnPremisesExtensionAttributesable)
    GetOnPremisesImmutableId()(*string)
    GetOnPremisesLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetOnPremisesProvisioningErrors()([]OnPremisesProvisioningErrorable)
    GetOnPremisesSamAccountName()(*string)
    GetOnPremisesSecurityIdentifier()(*string)
    GetOnPremisesSyncEnabled()(*bool)
    GetOnPremisesUserPrincipalName()(*string)
    GetOtherMails()([]string)
    GetOutlook()(OutlookUserable)
    GetOwnedDevices()([]DirectoryObjectable)
    GetOwnedObjects()([]DirectoryObjectable)
    GetPasswordPolicies()(*string)
    GetPasswordProfile()(PasswordProfileable)
    GetPastProjects()([]string)
    GetPeople()([]Personable)
    GetPermissionGrants()([]ResourceSpecificPermissionGrantable)
    GetPhoto()(ProfilePhotoable)
    GetPhotos()([]ProfilePhotoable)
    GetPlanner()(PlannerUserable)
    GetPostalCode()(*string)
    GetPreferredDataLocation()(*string)
    GetPreferredLanguage()(*string)
    GetPreferredName()(*string)
    GetPresence()(Presenceable)
    GetPrint()(UserPrintable)
    GetProvisionedPlans()([]ProvisionedPlanable)
    GetProxyAddresses()([]string)
    GetRegisteredDevices()([]DirectoryObjectable)
    GetResponsibilities()([]string)
    GetSchools()([]string)
    GetScopedRoleMemberOf()([]ScopedRoleMembershipable)
    GetSecurityIdentifier()(*string)
    GetServiceProvisioningErrors()([]ServiceProvisioningErrorable)
    GetSettings()(UserSettingsable)
    GetShowInAddressList()(*bool)
    GetSignInActivity()(SignInActivityable)
    GetSignInSessionsValidFromDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSkills()([]string)
    GetSolutions()(UserSolutionRootable)
    GetSponsors()([]DirectoryObjectable)
    GetState()(*string)
    GetStreetAddress()(*string)
    GetSurname()(*string)
    GetTeamwork()(UserTeamworkable)
    GetTodo()(Todoable)
    GetTransitiveMemberOf()([]DirectoryObjectable)
    GetUsageLocation()(*string)
    GetUserPrincipalName()(*string)
    GetUserType()(*string)
    SetAboutMe(value *string)()
    SetAccountEnabled(value *bool)()
    SetActivities(value []UserActivityable)()
    SetAgeGroup(value *string)()
    SetAgreementAcceptances(value []AgreementAcceptanceable)()
    SetAppRoleAssignments(value []AppRoleAssignmentable)()
    SetAssignedLicenses(value []AssignedLicenseable)()
    SetAssignedPlans(value []AssignedPlanable)()
    SetAuthentication(value Authenticationable)()
    SetAuthorizationInfo(value AuthorizationInfoable)()
    SetBirthday(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetBusinessPhones(value []string)()
    SetCalendar(value Calendarable)()
    SetCalendarGroups(value []CalendarGroupable)()
    SetCalendars(value []Calendarable)()
    SetCalendarView(value []Eventable)()
    SetChats(value []Chatable)()
    SetCity(value *string)()
    SetCloudClipboard(value CloudClipboardRootable)()
    SetCompanyName(value *string)()
    SetConsentProvidedForMinor(value *string)()
    SetContactFolders(value []ContactFolderable)()
    SetContacts(value []Contactable)()
    SetCountry(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedObjects(value []DirectoryObjectable)()
    SetCreationType(value *string)()
    SetCustomSecurityAttributes(value CustomSecurityAttributeValueable)()
    SetDepartment(value *string)()
    SetDeviceEnrollmentLimit(value *int32)()
    SetDeviceManagementTroubleshootingEvents(value []DeviceManagementTroubleshootingEventable)()
    SetDirectReports(value []DirectoryObjectable)()
    SetDisplayName(value *string)()
    SetDrive(value Driveable)()
    SetDrives(value []Driveable)()
    SetEmployeeExperience(value EmployeeExperienceUserable)()
    SetEmployeeHireDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEmployeeId(value *string)()
    SetEmployeeLeaveDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEmployeeOrgData(value EmployeeOrgDataable)()
    SetEmployeeType(value *string)()
    SetEvents(value []Eventable)()
    SetExtensions(value []Extensionable)()
    SetExternalUserState(value *string)()
    SetExternalUserStateChangeDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFaxNumber(value *string)()
    SetFollowedSites(value []Siteable)()
    SetGivenName(value *string)()
    SetHireDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIdentities(value []ObjectIdentityable)()
    SetImAddresses(value []string)()
    SetInferenceClassification(value InferenceClassificationable)()
    SetInsights(value ItemInsightsable)()
    SetInterests(value []string)()
    SetIsManagementRestricted(value *bool)()
    SetIsResourceAccount(value *bool)()
    SetJobTitle(value *string)()
    SetJoinedTeams(value []Teamable)()
    SetLastPasswordChangeDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLegalAgeGroupClassification(value *string)()
    SetLicenseAssignmentStates(value []LicenseAssignmentStateable)()
    SetLicenseDetails(value []LicenseDetailsable)()
    SetMail(value *string)()
    SetMailboxSettings(value MailboxSettingsable)()
    SetMailFolders(value []MailFolderable)()
    SetMailNickname(value *string)()
    SetManagedAppRegistrations(value []ManagedAppRegistrationable)()
    SetManagedDevices(value []ManagedDeviceable)()
    SetManager(value DirectoryObjectable)()
    SetMemberOf(value []DirectoryObjectable)()
    SetMessages(value []Messageable)()
    SetMobilePhone(value *string)()
    SetMySite(value *string)()
    SetOauth2PermissionGrants(value []OAuth2PermissionGrantable)()
    SetOfficeLocation(value *string)()
    SetOnenote(value Onenoteable)()
    SetOnlineMeetings(value []OnlineMeetingable)()
    SetOnPremisesDistinguishedName(value *string)()
    SetOnPremisesDomainName(value *string)()
    SetOnPremisesExtensionAttributes(value OnPremisesExtensionAttributesable)()
    SetOnPremisesImmutableId(value *string)()
    SetOnPremisesLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetOnPremisesProvisioningErrors(value []OnPremisesProvisioningErrorable)()
    SetOnPremisesSamAccountName(value *string)()
    SetOnPremisesSecurityIdentifier(value *string)()
    SetOnPremisesSyncEnabled(value *bool)()
    SetOnPremisesUserPrincipalName(value *string)()
    SetOtherMails(value []string)()
    SetOutlook(value OutlookUserable)()
    SetOwnedDevices(value []DirectoryObjectable)()
    SetOwnedObjects(value []DirectoryObjectable)()
    SetPasswordPolicies(value *string)()
    SetPasswordProfile(value PasswordProfileable)()
    SetPastProjects(value []string)()
    SetPeople(value []Personable)()
    SetPermissionGrants(value []ResourceSpecificPermissionGrantable)()
    SetPhoto(value ProfilePhotoable)()
    SetPhotos(value []ProfilePhotoable)()
    SetPlanner(value PlannerUserable)()
    SetPostalCode(value *string)()
    SetPreferredDataLocation(value *string)()
    SetPreferredLanguage(value *string)()
    SetPreferredName(value *string)()
    SetPresence(value Presenceable)()
    SetPrint(value UserPrintable)()
    SetProvisionedPlans(value []ProvisionedPlanable)()
    SetProxyAddresses(value []string)()
    SetRegisteredDevices(value []DirectoryObjectable)()
    SetResponsibilities(value []string)()
    SetSchools(value []string)()
    SetScopedRoleMemberOf(value []ScopedRoleMembershipable)()
    SetSecurityIdentifier(value *string)()
    SetServiceProvisioningErrors(value []ServiceProvisioningErrorable)()
    SetSettings(value UserSettingsable)()
    SetShowInAddressList(value *bool)()
    SetSignInActivity(value SignInActivityable)()
    SetSignInSessionsValidFromDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSkills(value []string)()
    SetSolutions(value UserSolutionRootable)()
    SetSponsors(value []DirectoryObjectable)()
    SetState(value *string)()
    SetStreetAddress(value *string)()
    SetSurname(value *string)()
    SetTeamwork(value UserTeamworkable)()
    SetTodo(value Todoable)()
    SetTransitiveMemberOf(value []DirectoryObjectable)()
    SetUsageLocation(value *string)()
    SetUserPrincipalName(value *string)()
    SetUserType(value *string)()
}

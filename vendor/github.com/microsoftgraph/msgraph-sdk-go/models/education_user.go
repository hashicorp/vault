package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationUser struct {
    Entity
}
// NewEducationUser instantiates a new EducationUser and sets the default values.
func NewEducationUser()(*EducationUser) {
    m := &EducationUser{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationUser(), nil
}
// GetAccountEnabled gets the accountEnabled property value. True if the account is enabled; otherwise, false. This property is required when a user is created. Supports $filter.
// returns a *bool when successful
func (m *EducationUser) GetAccountEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("accountEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAssignedLicenses gets the assignedLicenses property value. The licenses that are assigned to the user. Not nullable.
// returns a []AssignedLicenseable when successful
func (m *EducationUser) GetAssignedLicenses()([]AssignedLicenseable) {
    val, err := m.GetBackingStore().Get("assignedLicenses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AssignedLicenseable)
    }
    return nil
}
// GetAssignedPlans gets the assignedPlans property value. The plans that are assigned to the user. Read-only. Not nullable.
// returns a []AssignedPlanable when successful
func (m *EducationUser) GetAssignedPlans()([]AssignedPlanable) {
    val, err := m.GetBackingStore().Get("assignedPlans")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AssignedPlanable)
    }
    return nil
}
// GetAssignments gets the assignments property value. Assignments belonging to the user.
// returns a []EducationAssignmentable when successful
func (m *EducationUser) GetAssignments()([]EducationAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationAssignmentable)
    }
    return nil
}
// GetBusinessPhones gets the businessPhones property value. The telephone numbers for the user. Note: Although this is a string collection, only one number can be set for this property.
// returns a []string when successful
func (m *EducationUser) GetBusinessPhones()([]string) {
    val, err := m.GetBackingStore().Get("businessPhones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetClasses gets the classes property value. Classes to which the user belongs. Nullable.
// returns a []EducationClassable when successful
func (m *EducationUser) GetClasses()([]EducationClassable) {
    val, err := m.GetBackingStore().Get("classes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationClassable)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. The entity who created the user.
// returns a IdentitySetable when successful
func (m *EducationUser) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetDepartment gets the department property value. The name for the department in which the user works. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetDepartment()(*string) {
    val, err := m.GetBackingStore().Get("department")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial, and last name. This property is required when a user is created and it cannot be cleared during updates. Supports $filter and $orderby.
// returns a *string when successful
func (m *EducationUser) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalSource gets the externalSource property value. Where this user was created from. Possible values are: sis, manual.
// returns a *EducationExternalSource when successful
func (m *EducationUser) GetExternalSource()(*EducationExternalSource) {
    val, err := m.GetBackingStore().Get("externalSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationExternalSource)
    }
    return nil
}
// GetExternalSourceDetail gets the externalSourceDetail property value. The name of the external source this resource was generated from.
// returns a *string when successful
func (m *EducationUser) GetExternalSourceDetail()(*string) {
    val, err := m.GetBackingStore().Get("externalSourceDetail")
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
func (m *EducationUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationAssignmentable)
                }
            }
            m.SetAssignments(res)
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
    res["classes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationClassFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationClassable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationClassable)
                }
            }
            m.SetClasses(res)
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(IdentitySetable))
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
    res["externalSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationExternalSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSource(val.(*EducationExternalSource))
        }
        return nil
    }
    res["externalSourceDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSourceDetail(val)
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
    res["mailingAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailingAddress(val.(PhysicalAddressable))
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
    res["middleName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMiddleName(val)
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
    res["onPremisesInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationOnPremisesInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnPremisesInfo(val.(EducationOnPremisesInfoable))
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
    res["primaryRole"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationUserRole)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrimaryRole(val.(*EducationUserRole))
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
    res["refreshTokensValidFromDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRefreshTokensValidFromDateTime(val)
        }
        return nil
    }
    res["relatedContacts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRelatedContactFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RelatedContactable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RelatedContactable)
                }
            }
            m.SetRelatedContacts(res)
        }
        return nil
    }
    res["residenceAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResidenceAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["rubrics"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationRubricFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationRubricable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationRubricable)
                }
            }
            m.SetRubrics(res)
        }
        return nil
    }
    res["schools"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationSchoolFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationSchoolable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationSchoolable)
                }
            }
            m.SetSchools(res)
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
    res["student"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationStudentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStudent(val.(EducationStudentable))
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
    res["taughtClasses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationClassFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationClassable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationClassable)
                }
            }
            m.SetTaughtClasses(res)
        }
        return nil
    }
    res["teacher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationTeacherFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeacher(val.(EducationTeacherable))
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
    res["user"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUser(val.(Userable))
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
// GetGivenName gets the givenName property value. The given name (first name) of the user. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetGivenName()(*string) {
    val, err := m.GetBackingStore().Get("givenName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMail gets the mail property value. The SMTP address for the user, for example, jeff@contoso.com. Read-Only. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetMail()(*string) {
    val, err := m.GetBackingStore().Get("mail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMailingAddress gets the mailingAddress property value. The mail address of the user.
// returns a PhysicalAddressable when successful
func (m *EducationUser) GetMailingAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("mailingAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetMailNickname gets the mailNickname property value. The mail alias for the user. This property must be specified when a user is created. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetMailNickname()(*string) {
    val, err := m.GetBackingStore().Get("mailNickname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMiddleName gets the middleName property value. The middle name of the user.
// returns a *string when successful
func (m *EducationUser) GetMiddleName()(*string) {
    val, err := m.GetBackingStore().Get("middleName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMobilePhone gets the mobilePhone property value. The primary cellular telephone number for the user.
// returns a *string when successful
func (m *EducationUser) GetMobilePhone()(*string) {
    val, err := m.GetBackingStore().Get("mobilePhone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOfficeLocation gets the officeLocation property value. The office location for the user.
// returns a *string when successful
func (m *EducationUser) GetOfficeLocation()(*string) {
    val, err := m.GetBackingStore().Get("officeLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesInfo gets the onPremisesInfo property value. Additional information used to associate the Microsoft Entra user with its Active Directory counterpart.
// returns a EducationOnPremisesInfoable when successful
func (m *EducationUser) GetOnPremisesInfo()(EducationOnPremisesInfoable) {
    val, err := m.GetBackingStore().Get("onPremisesInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationOnPremisesInfoable)
    }
    return nil
}
// GetPasswordPolicies gets the passwordPolicies property value. Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two can be specified together; for example: DisablePasswordExpiration, DisableStrongPassword.
// returns a *string when successful
func (m *EducationUser) GetPasswordPolicies()(*string) {
    val, err := m.GetBackingStore().Get("passwordPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasswordProfile gets the passwordProfile property value. Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required.
// returns a PasswordProfileable when successful
func (m *EducationUser) GetPasswordProfile()(PasswordProfileable) {
    val, err := m.GetBackingStore().Get("passwordProfile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PasswordProfileable)
    }
    return nil
}
// GetPreferredLanguage gets the preferredLanguage property value. The preferred language for the user that should follow the ISO 639-1 code, for example, en-US.
// returns a *string when successful
func (m *EducationUser) GetPreferredLanguage()(*string) {
    val, err := m.GetBackingStore().Get("preferredLanguage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrimaryRole gets the primaryRole property value. The primaryRole property
// returns a *EducationUserRole when successful
func (m *EducationUser) GetPrimaryRole()(*EducationUserRole) {
    val, err := m.GetBackingStore().Get("primaryRole")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationUserRole)
    }
    return nil
}
// GetProvisionedPlans gets the provisionedPlans property value. The plans that are provisioned for the user. Read-only. Not nullable.
// returns a []ProvisionedPlanable when successful
func (m *EducationUser) GetProvisionedPlans()([]ProvisionedPlanable) {
    val, err := m.GetBackingStore().Get("provisionedPlans")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProvisionedPlanable)
    }
    return nil
}
// GetRefreshTokensValidFromDateTime gets the refreshTokensValidFromDateTime property value. Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).  If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Returned only on $select. Read-only.
// returns a *Time when successful
func (m *EducationUser) GetRefreshTokensValidFromDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("refreshTokensValidFromDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRelatedContacts gets the relatedContacts property value. Related records associated with the user. Read-only.
// returns a []RelatedContactable when successful
func (m *EducationUser) GetRelatedContacts()([]RelatedContactable) {
    val, err := m.GetBackingStore().Get("relatedContacts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RelatedContactable)
    }
    return nil
}
// GetResidenceAddress gets the residenceAddress property value. The address where the user lives.
// returns a PhysicalAddressable when successful
func (m *EducationUser) GetResidenceAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("residenceAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetRubrics gets the rubrics property value. When set, the grading rubric attached to the assignment.
// returns a []EducationRubricable when successful
func (m *EducationUser) GetRubrics()([]EducationRubricable) {
    val, err := m.GetBackingStore().Get("rubrics")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationRubricable)
    }
    return nil
}
// GetSchools gets the schools property value. Schools to which the user belongs. Nullable.
// returns a []EducationSchoolable when successful
func (m *EducationUser) GetSchools()([]EducationSchoolable) {
    val, err := m.GetBackingStore().Get("schools")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationSchoolable)
    }
    return nil
}
// GetShowInAddressList gets the showInAddressList property value. True if the Outlook Global Address List should contain this user; otherwise, false. If not set, this will be treated as true. For users invited through the invitation manager, this property will be set to false.
// returns a *bool when successful
func (m *EducationUser) GetShowInAddressList()(*bool) {
    val, err := m.GetBackingStore().Get("showInAddressList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStudent gets the student property value. If the primary role is student, this block will contain student specific data.
// returns a EducationStudentable when successful
func (m *EducationUser) GetStudent()(EducationStudentable) {
    val, err := m.GetBackingStore().Get("student")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationStudentable)
    }
    return nil
}
// GetSurname gets the surname property value. The user's surname (family name or last name). Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetSurname()(*string) {
    val, err := m.GetBackingStore().Get("surname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTaughtClasses gets the taughtClasses property value. Classes for which the user is a teacher.
// returns a []EducationClassable when successful
func (m *EducationUser) GetTaughtClasses()([]EducationClassable) {
    val, err := m.GetBackingStore().Get("taughtClasses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationClassable)
    }
    return nil
}
// GetTeacher gets the teacher property value. If the primary role is teacher, this block will contain teacher specific data.
// returns a EducationTeacherable when successful
func (m *EducationUser) GetTeacher()(EducationTeacherable) {
    val, err := m.GetBackingStore().Get("teacher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationTeacherable)
    }
    return nil
}
// GetUsageLocation gets the usageLocation property value. A two-letter country code (ISO standard 3166). Required for users who will be assigned licenses due to a legal requirement to check for availability of services in countries or regions. Examples include: US, JP, and GB. Not nullable. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetUsageLocation()(*string) {
    val, err := m.GetBackingStore().Get("usageLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUser gets the user property value. The directory user that corresponds to this user.
// returns a Userable when successful
func (m *EducationUser) GetUser()(Userable) {
    val, err := m.GetBackingStore().Get("user")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Userable)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. The user principal name (UPN) of the user. The UPN is an internet-style login name for the user based on the internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of the organization. Supports $filter and $orderby.
// returns a *string when successful
func (m *EducationUser) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserType gets the userType property value. A string value that can be used to classify user types in your directory, such as Member and Guest. Supports $filter.
// returns a *string when successful
func (m *EducationUser) GetUserType()(*string) {
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
func (m *EducationUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("accountEnabled", m.GetAccountEnabled())
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
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
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
    if m.GetClasses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetClasses()))
        for i, v := range m.GetClasses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("classes", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
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
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetExternalSource() != nil {
        cast := (*m.GetExternalSource()).String()
        err = writer.WriteStringValue("externalSource", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalSourceDetail", m.GetExternalSourceDetail())
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
        err = writer.WriteStringValue("mail", m.GetMail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mailingAddress", m.GetMailingAddress())
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
    {
        err = writer.WriteStringValue("middleName", m.GetMiddleName())
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
        err = writer.WriteStringValue("officeLocation", m.GetOfficeLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onPremisesInfo", m.GetOnPremisesInfo())
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
    {
        err = writer.WriteStringValue("preferredLanguage", m.GetPreferredLanguage())
        if err != nil {
            return err
        }
    }
    if m.GetPrimaryRole() != nil {
        cast := (*m.GetPrimaryRole()).String()
        err = writer.WriteStringValue("primaryRole", &cast)
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
    {
        err = writer.WriteTimeValue("refreshTokensValidFromDateTime", m.GetRefreshTokensValidFromDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetRelatedContacts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRelatedContacts()))
        for i, v := range m.GetRelatedContacts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("relatedContacts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("residenceAddress", m.GetResidenceAddress())
        if err != nil {
            return err
        }
    }
    if m.GetRubrics() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRubrics()))
        for i, v := range m.GetRubrics() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rubrics", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSchools() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSchools()))
        for i, v := range m.GetSchools() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("schools", cast)
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
        err = writer.WriteObjectValue("student", m.GetStudent())
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
    if m.GetTaughtClasses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTaughtClasses()))
        for i, v := range m.GetTaughtClasses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("taughtClasses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("teacher", m.GetTeacher())
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
        err = writer.WriteObjectValue("user", m.GetUser())
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
// SetAccountEnabled sets the accountEnabled property value. True if the account is enabled; otherwise, false. This property is required when a user is created. Supports $filter.
func (m *EducationUser) SetAccountEnabled(value *bool)() {
    err := m.GetBackingStore().Set("accountEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedLicenses sets the assignedLicenses property value. The licenses that are assigned to the user. Not nullable.
func (m *EducationUser) SetAssignedLicenses(value []AssignedLicenseable)() {
    err := m.GetBackingStore().Set("assignedLicenses", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedPlans sets the assignedPlans property value. The plans that are assigned to the user. Read-only. Not nullable.
func (m *EducationUser) SetAssignedPlans(value []AssignedPlanable)() {
    err := m.GetBackingStore().Set("assignedPlans", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignments sets the assignments property value. Assignments belonging to the user.
func (m *EducationUser) SetAssignments(value []EducationAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessPhones sets the businessPhones property value. The telephone numbers for the user. Note: Although this is a string collection, only one number can be set for this property.
func (m *EducationUser) SetBusinessPhones(value []string)() {
    err := m.GetBackingStore().Set("businessPhones", value)
    if err != nil {
        panic(err)
    }
}
// SetClasses sets the classes property value. Classes to which the user belongs. Nullable.
func (m *EducationUser) SetClasses(value []EducationClassable)() {
    err := m.GetBackingStore().Set("classes", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. The entity who created the user.
func (m *EducationUser) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartment sets the department property value. The name for the department in which the user works. Supports $filter.
func (m *EducationUser) SetDepartment(value *string)() {
    err := m.GetBackingStore().Set("department", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial, and last name. This property is required when a user is created and it cannot be cleared during updates. Supports $filter and $orderby.
func (m *EducationUser) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSource sets the externalSource property value. Where this user was created from. Possible values are: sis, manual.
func (m *EducationUser) SetExternalSource(value *EducationExternalSource)() {
    err := m.GetBackingStore().Set("externalSource", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSourceDetail sets the externalSourceDetail property value. The name of the external source this resource was generated from.
func (m *EducationUser) SetExternalSourceDetail(value *string)() {
    err := m.GetBackingStore().Set("externalSourceDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetGivenName sets the givenName property value. The given name (first name) of the user. Supports $filter.
func (m *EducationUser) SetGivenName(value *string)() {
    err := m.GetBackingStore().Set("givenName", value)
    if err != nil {
        panic(err)
    }
}
// SetMail sets the mail property value. The SMTP address for the user, for example, jeff@contoso.com. Read-Only. Supports $filter.
func (m *EducationUser) SetMail(value *string)() {
    err := m.GetBackingStore().Set("mail", value)
    if err != nil {
        panic(err)
    }
}
// SetMailingAddress sets the mailingAddress property value. The mail address of the user.
func (m *EducationUser) SetMailingAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("mailingAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetMailNickname sets the mailNickname property value. The mail alias for the user. This property must be specified when a user is created. Supports $filter.
func (m *EducationUser) SetMailNickname(value *string)() {
    err := m.GetBackingStore().Set("mailNickname", value)
    if err != nil {
        panic(err)
    }
}
// SetMiddleName sets the middleName property value. The middle name of the user.
func (m *EducationUser) SetMiddleName(value *string)() {
    err := m.GetBackingStore().Set("middleName", value)
    if err != nil {
        panic(err)
    }
}
// SetMobilePhone sets the mobilePhone property value. The primary cellular telephone number for the user.
func (m *EducationUser) SetMobilePhone(value *string)() {
    err := m.GetBackingStore().Set("mobilePhone", value)
    if err != nil {
        panic(err)
    }
}
// SetOfficeLocation sets the officeLocation property value. The office location for the user.
func (m *EducationUser) SetOfficeLocation(value *string)() {
    err := m.GetBackingStore().Set("officeLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesInfo sets the onPremisesInfo property value. Additional information used to associate the Microsoft Entra user with its Active Directory counterpart.
func (m *EducationUser) SetOnPremisesInfo(value EducationOnPremisesInfoable)() {
    err := m.GetBackingStore().Set("onPremisesInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPolicies sets the passwordPolicies property value. Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two can be specified together; for example: DisablePasswordExpiration, DisableStrongPassword.
func (m *EducationUser) SetPasswordPolicies(value *string)() {
    err := m.GetBackingStore().Set("passwordPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordProfile sets the passwordProfile property value. Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required.
func (m *EducationUser) SetPasswordProfile(value PasswordProfileable)() {
    err := m.GetBackingStore().Set("passwordProfile", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredLanguage sets the preferredLanguage property value. The preferred language for the user that should follow the ISO 639-1 code, for example, en-US.
func (m *EducationUser) SetPreferredLanguage(value *string)() {
    err := m.GetBackingStore().Set("preferredLanguage", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryRole sets the primaryRole property value. The primaryRole property
func (m *EducationUser) SetPrimaryRole(value *EducationUserRole)() {
    err := m.GetBackingStore().Set("primaryRole", value)
    if err != nil {
        panic(err)
    }
}
// SetProvisionedPlans sets the provisionedPlans property value. The plans that are provisioned for the user. Read-only. Not nullable.
func (m *EducationUser) SetProvisionedPlans(value []ProvisionedPlanable)() {
    err := m.GetBackingStore().Set("provisionedPlans", value)
    if err != nil {
        panic(err)
    }
}
// SetRefreshTokensValidFromDateTime sets the refreshTokensValidFromDateTime property value. Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).  If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Returned only on $select. Read-only.
func (m *EducationUser) SetRefreshTokensValidFromDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("refreshTokensValidFromDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRelatedContacts sets the relatedContacts property value. Related records associated with the user. Read-only.
func (m *EducationUser) SetRelatedContacts(value []RelatedContactable)() {
    err := m.GetBackingStore().Set("relatedContacts", value)
    if err != nil {
        panic(err)
    }
}
// SetResidenceAddress sets the residenceAddress property value. The address where the user lives.
func (m *EducationUser) SetResidenceAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("residenceAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetRubrics sets the rubrics property value. When set, the grading rubric attached to the assignment.
func (m *EducationUser) SetRubrics(value []EducationRubricable)() {
    err := m.GetBackingStore().Set("rubrics", value)
    if err != nil {
        panic(err)
    }
}
// SetSchools sets the schools property value. Schools to which the user belongs. Nullable.
func (m *EducationUser) SetSchools(value []EducationSchoolable)() {
    err := m.GetBackingStore().Set("schools", value)
    if err != nil {
        panic(err)
    }
}
// SetShowInAddressList sets the showInAddressList property value. True if the Outlook Global Address List should contain this user; otherwise, false. If not set, this will be treated as true. For users invited through the invitation manager, this property will be set to false.
func (m *EducationUser) SetShowInAddressList(value *bool)() {
    err := m.GetBackingStore().Set("showInAddressList", value)
    if err != nil {
        panic(err)
    }
}
// SetStudent sets the student property value. If the primary role is student, this block will contain student specific data.
func (m *EducationUser) SetStudent(value EducationStudentable)() {
    err := m.GetBackingStore().Set("student", value)
    if err != nil {
        panic(err)
    }
}
// SetSurname sets the surname property value. The user's surname (family name or last name). Supports $filter.
func (m *EducationUser) SetSurname(value *string)() {
    err := m.GetBackingStore().Set("surname", value)
    if err != nil {
        panic(err)
    }
}
// SetTaughtClasses sets the taughtClasses property value. Classes for which the user is a teacher.
func (m *EducationUser) SetTaughtClasses(value []EducationClassable)() {
    err := m.GetBackingStore().Set("taughtClasses", value)
    if err != nil {
        panic(err)
    }
}
// SetTeacher sets the teacher property value. If the primary role is teacher, this block will contain teacher specific data.
func (m *EducationUser) SetTeacher(value EducationTeacherable)() {
    err := m.GetBackingStore().Set("teacher", value)
    if err != nil {
        panic(err)
    }
}
// SetUsageLocation sets the usageLocation property value. A two-letter country code (ISO standard 3166). Required for users who will be assigned licenses due to a legal requirement to check for availability of services in countries or regions. Examples include: US, JP, and GB. Not nullable. Supports $filter.
func (m *EducationUser) SetUsageLocation(value *string)() {
    err := m.GetBackingStore().Set("usageLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetUser sets the user property value. The directory user that corresponds to this user.
func (m *EducationUser) SetUser(value Userable)() {
    err := m.GetBackingStore().Set("user", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. The user principal name (UPN) of the user. The UPN is an internet-style login name for the user based on the internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of the organization. Supports $filter and $orderby.
func (m *EducationUser) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserType sets the userType property value. A string value that can be used to classify user types in your directory, such as Member and Guest. Supports $filter.
func (m *EducationUser) SetUserType(value *string)() {
    err := m.GetBackingStore().Set("userType", value)
    if err != nil {
        panic(err)
    }
}
type EducationUserable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountEnabled()(*bool)
    GetAssignedLicenses()([]AssignedLicenseable)
    GetAssignedPlans()([]AssignedPlanable)
    GetAssignments()([]EducationAssignmentable)
    GetBusinessPhones()([]string)
    GetClasses()([]EducationClassable)
    GetCreatedBy()(IdentitySetable)
    GetDepartment()(*string)
    GetDisplayName()(*string)
    GetExternalSource()(*EducationExternalSource)
    GetExternalSourceDetail()(*string)
    GetGivenName()(*string)
    GetMail()(*string)
    GetMailingAddress()(PhysicalAddressable)
    GetMailNickname()(*string)
    GetMiddleName()(*string)
    GetMobilePhone()(*string)
    GetOfficeLocation()(*string)
    GetOnPremisesInfo()(EducationOnPremisesInfoable)
    GetPasswordPolicies()(*string)
    GetPasswordProfile()(PasswordProfileable)
    GetPreferredLanguage()(*string)
    GetPrimaryRole()(*EducationUserRole)
    GetProvisionedPlans()([]ProvisionedPlanable)
    GetRefreshTokensValidFromDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRelatedContacts()([]RelatedContactable)
    GetResidenceAddress()(PhysicalAddressable)
    GetRubrics()([]EducationRubricable)
    GetSchools()([]EducationSchoolable)
    GetShowInAddressList()(*bool)
    GetStudent()(EducationStudentable)
    GetSurname()(*string)
    GetTaughtClasses()([]EducationClassable)
    GetTeacher()(EducationTeacherable)
    GetUsageLocation()(*string)
    GetUser()(Userable)
    GetUserPrincipalName()(*string)
    GetUserType()(*string)
    SetAccountEnabled(value *bool)()
    SetAssignedLicenses(value []AssignedLicenseable)()
    SetAssignedPlans(value []AssignedPlanable)()
    SetAssignments(value []EducationAssignmentable)()
    SetBusinessPhones(value []string)()
    SetClasses(value []EducationClassable)()
    SetCreatedBy(value IdentitySetable)()
    SetDepartment(value *string)()
    SetDisplayName(value *string)()
    SetExternalSource(value *EducationExternalSource)()
    SetExternalSourceDetail(value *string)()
    SetGivenName(value *string)()
    SetMail(value *string)()
    SetMailingAddress(value PhysicalAddressable)()
    SetMailNickname(value *string)()
    SetMiddleName(value *string)()
    SetMobilePhone(value *string)()
    SetOfficeLocation(value *string)()
    SetOnPremisesInfo(value EducationOnPremisesInfoable)()
    SetPasswordPolicies(value *string)()
    SetPasswordProfile(value PasswordProfileable)()
    SetPreferredLanguage(value *string)()
    SetPrimaryRole(value *EducationUserRole)()
    SetProvisionedPlans(value []ProvisionedPlanable)()
    SetRefreshTokensValidFromDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRelatedContacts(value []RelatedContactable)()
    SetResidenceAddress(value PhysicalAddressable)()
    SetRubrics(value []EducationRubricable)()
    SetSchools(value []EducationSchoolable)()
    SetShowInAddressList(value *bool)()
    SetStudent(value EducationStudentable)()
    SetSurname(value *string)()
    SetTaughtClasses(value []EducationClassable)()
    SetTeacher(value EducationTeacherable)()
    SetUsageLocation(value *string)()
    SetUser(value Userable)()
    SetUserPrincipalName(value *string)()
    SetUserType(value *string)()
}

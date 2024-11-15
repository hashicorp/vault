package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Contact struct {
    OutlookItem
}
// NewContact instantiates a new Contact and sets the default values.
func NewContact()(*Contact) {
    m := &Contact{
        OutlookItem: *NewOutlookItem(),
    }
    odataTypeValue := "#microsoft.graph.contact"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateContactFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateContactFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContact(), nil
}
// GetAssistantName gets the assistantName property value. The name of the contact's assistant.
// returns a *string when successful
func (m *Contact) GetAssistantName()(*string) {
    val, err := m.GetBackingStore().Get("assistantName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBirthday gets the birthday property value. The contact's birthday. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *Contact) GetBirthday()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("birthday")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetBusinessAddress gets the businessAddress property value. The contact's business address.
// returns a PhysicalAddressable when successful
func (m *Contact) GetBusinessAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("businessAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetBusinessHomePage gets the businessHomePage property value. The business home page of the contact.
// returns a *string when successful
func (m *Contact) GetBusinessHomePage()(*string) {
    val, err := m.GetBackingStore().Get("businessHomePage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBusinessPhones gets the businessPhones property value. The contact's business phone numbers.
// returns a []string when successful
func (m *Contact) GetBusinessPhones()([]string) {
    val, err := m.GetBackingStore().Get("businessPhones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetChildren gets the children property value. The names of the contact's children.
// returns a []string when successful
func (m *Contact) GetChildren()([]string) {
    val, err := m.GetBackingStore().Get("children")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCompanyName gets the companyName property value. The name of the contact's company.
// returns a *string when successful
func (m *Contact) GetCompanyName()(*string) {
    val, err := m.GetBackingStore().Get("companyName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDepartment gets the department property value. The contact's department.
// returns a *string when successful
func (m *Contact) GetDepartment()(*string) {
    val, err := m.GetBackingStore().Get("department")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The contact's display name. You can specify the display name in a create or update operation. Note that later updates to other properties may cause an automatically generated value to overwrite the displayName value you have specified. To preserve a pre-existing value, always include it as displayName in an update operation.
// returns a *string when successful
func (m *Contact) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmailAddresses gets the emailAddresses property value. The contact's email addresses.
// returns a []EmailAddressable when successful
func (m *Contact) GetEmailAddresses()([]EmailAddressable) {
    val, err := m.GetBackingStore().Get("emailAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EmailAddressable)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the contact. Read-only. Nullable.
// returns a []Extensionable when successful
func (m *Contact) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Contact) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OutlookItem.GetFieldDeserializers()
    res["assistantName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssistantName(val)
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
    res["businessAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBusinessAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["businessHomePage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBusinessHomePage(val)
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
    res["children"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetChildren(res)
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
    res["emailAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEmailAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EmailAddressable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EmailAddressable)
                }
            }
            m.SetEmailAddresses(res)
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
    res["fileAs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileAs(val)
        }
        return nil
    }
    res["generation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGeneration(val)
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
    res["homeAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHomeAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["homePhones"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetHomePhones(res)
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
    res["initials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitials(val)
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
    res["manager"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManager(val)
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
    res["multiValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MultiValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MultiValueLegacyExtendedPropertyable)
                }
            }
            m.SetMultiValueExtendedProperties(res)
        }
        return nil
    }
    res["nickName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNickName(val)
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
    res["otherAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOtherAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["parentFolderId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentFolderId(val)
        }
        return nil
    }
    res["personalNotes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonalNotes(val)
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
    res["profession"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProfession(val)
        }
        return nil
    }
    res["singleValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SingleValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SingleValueLegacyExtendedPropertyable)
                }
            }
            m.SetSingleValueExtendedProperties(res)
        }
        return nil
    }
    res["spouseName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpouseName(val)
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
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    res["yomiCompanyName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetYomiCompanyName(val)
        }
        return nil
    }
    res["yomiGivenName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetYomiGivenName(val)
        }
        return nil
    }
    res["yomiSurname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetYomiSurname(val)
        }
        return nil
    }
    return res
}
// GetFileAs gets the fileAs property value. The name the contact is filed under.
// returns a *string when successful
func (m *Contact) GetFileAs()(*string) {
    val, err := m.GetBackingStore().Get("fileAs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetGeneration gets the generation property value. The contact's suffix.
// returns a *string when successful
func (m *Contact) GetGeneration()(*string) {
    val, err := m.GetBackingStore().Get("generation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetGivenName gets the givenName property value. The contact's given name.
// returns a *string when successful
func (m *Contact) GetGivenName()(*string) {
    val, err := m.GetBackingStore().Get("givenName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHomeAddress gets the homeAddress property value. The contact's home address.
// returns a PhysicalAddressable when successful
func (m *Contact) GetHomeAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("homeAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetHomePhones gets the homePhones property value. The contact's home phone numbers.
// returns a []string when successful
func (m *Contact) GetHomePhones()([]string) {
    val, err := m.GetBackingStore().Get("homePhones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetImAddresses gets the imAddresses property value. The contact's instant messaging (IM) addresses.
// returns a []string when successful
func (m *Contact) GetImAddresses()([]string) {
    val, err := m.GetBackingStore().Get("imAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetInitials gets the initials property value. The contact's initials.
// returns a *string when successful
func (m *Contact) GetInitials()(*string) {
    val, err := m.GetBackingStore().Get("initials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetJobTitle gets the jobTitle property value. The contact’s job title.
// returns a *string when successful
func (m *Contact) GetJobTitle()(*string) {
    val, err := m.GetBackingStore().Get("jobTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetManager gets the manager property value. The name of the contact's manager.
// returns a *string when successful
func (m *Contact) GetManager()(*string) {
    val, err := m.GetBackingStore().Get("manager")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMiddleName gets the middleName property value. The contact's middle name.
// returns a *string when successful
func (m *Contact) GetMiddleName()(*string) {
    val, err := m.GetBackingStore().Get("middleName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMobilePhone gets the mobilePhone property value. The contact's mobile phone number.
// returns a *string when successful
func (m *Contact) GetMobilePhone()(*string) {
    val, err := m.GetBackingStore().Get("mobilePhone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the contact. Read-only. Nullable.
// returns a []MultiValueLegacyExtendedPropertyable when successful
func (m *Contact) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("multiValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MultiValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetNickName gets the nickName property value. The contact's nickname.
// returns a *string when successful
func (m *Contact) GetNickName()(*string) {
    val, err := m.GetBackingStore().Get("nickName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOfficeLocation gets the officeLocation property value. The location of the contact's office.
// returns a *string when successful
func (m *Contact) GetOfficeLocation()(*string) {
    val, err := m.GetBackingStore().Get("officeLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOtherAddress gets the otherAddress property value. Other addresses for the contact.
// returns a PhysicalAddressable when successful
func (m *Contact) GetOtherAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("otherAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetParentFolderId gets the parentFolderId property value. The ID of the contact's parent folder.
// returns a *string when successful
func (m *Contact) GetParentFolderId()(*string) {
    val, err := m.GetBackingStore().Get("parentFolderId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersonalNotes gets the personalNotes property value. The user's notes about the contact.
// returns a *string when successful
func (m *Contact) GetPersonalNotes()(*string) {
    val, err := m.GetBackingStore().Get("personalNotes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhoto gets the photo property value. Optional contact picture. You can get or set a photo for a contact.
// returns a ProfilePhotoable when successful
func (m *Contact) GetPhoto()(ProfilePhotoable) {
    val, err := m.GetBackingStore().Get("photo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ProfilePhotoable)
    }
    return nil
}
// GetProfession gets the profession property value. The contact's profession.
// returns a *string when successful
func (m *Contact) GetProfession()(*string) {
    val, err := m.GetBackingStore().Get("profession")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the contact. Read-only. Nullable.
// returns a []SingleValueLegacyExtendedPropertyable when successful
func (m *Contact) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("singleValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SingleValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetSpouseName gets the spouseName property value. The name of the contact's spouse/partner.
// returns a *string when successful
func (m *Contact) GetSpouseName()(*string) {
    val, err := m.GetBackingStore().Get("spouseName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSurname gets the surname property value. The contact's surname.
// returns a *string when successful
func (m *Contact) GetSurname()(*string) {
    val, err := m.GetBackingStore().Get("surname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTitle gets the title property value. The contact's title.
// returns a *string when successful
func (m *Contact) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetYomiCompanyName gets the yomiCompanyName property value. The phonetic Japanese company name of the contact.
// returns a *string when successful
func (m *Contact) GetYomiCompanyName()(*string) {
    val, err := m.GetBackingStore().Get("yomiCompanyName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetYomiGivenName gets the yomiGivenName property value. The phonetic Japanese given name (first name) of the contact.
// returns a *string when successful
func (m *Contact) GetYomiGivenName()(*string) {
    val, err := m.GetBackingStore().Get("yomiGivenName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetYomiSurname gets the yomiSurname property value. The phonetic Japanese surname (last name)  of the contact.
// returns a *string when successful
func (m *Contact) GetYomiSurname()(*string) {
    val, err := m.GetBackingStore().Get("yomiSurname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Contact) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OutlookItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("assistantName", m.GetAssistantName())
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
    {
        err = writer.WriteObjectValue("businessAddress", m.GetBusinessAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("businessHomePage", m.GetBusinessHomePage())
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
    if m.GetChildren() != nil {
        err = writer.WriteCollectionOfStringValues("children", m.GetChildren())
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
    if m.GetEmailAddresses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEmailAddresses()))
        for i, v := range m.GetEmailAddresses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("emailAddresses", cast)
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
        err = writer.WriteStringValue("fileAs", m.GetFileAs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("generation", m.GetGeneration())
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
        err = writer.WriteObjectValue("homeAddress", m.GetHomeAddress())
        if err != nil {
            return err
        }
    }
    if m.GetHomePhones() != nil {
        err = writer.WriteCollectionOfStringValues("homePhones", m.GetHomePhones())
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
        err = writer.WriteStringValue("initials", m.GetInitials())
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
    {
        err = writer.WriteStringValue("manager", m.GetManager())
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
    if m.GetMultiValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMultiValueExtendedProperties()))
        for i, v := range m.GetMultiValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("multiValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("nickName", m.GetNickName())
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
        err = writer.WriteObjectValue("otherAddress", m.GetOtherAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentFolderId", m.GetParentFolderId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("personalNotes", m.GetPersonalNotes())
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
    {
        err = writer.WriteStringValue("profession", m.GetProfession())
        if err != nil {
            return err
        }
    }
    if m.GetSingleValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSingleValueExtendedProperties()))
        for i, v := range m.GetSingleValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("singleValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("spouseName", m.GetSpouseName())
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
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("yomiCompanyName", m.GetYomiCompanyName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("yomiGivenName", m.GetYomiGivenName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("yomiSurname", m.GetYomiSurname())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssistantName sets the assistantName property value. The name of the contact's assistant.
func (m *Contact) SetAssistantName(value *string)() {
    err := m.GetBackingStore().Set("assistantName", value)
    if err != nil {
        panic(err)
    }
}
// SetBirthday sets the birthday property value. The contact's birthday. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *Contact) SetBirthday(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("birthday", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessAddress sets the businessAddress property value. The contact's business address.
func (m *Contact) SetBusinessAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("businessAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessHomePage sets the businessHomePage property value. The business home page of the contact.
func (m *Contact) SetBusinessHomePage(value *string)() {
    err := m.GetBackingStore().Set("businessHomePage", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessPhones sets the businessPhones property value. The contact's business phone numbers.
func (m *Contact) SetBusinessPhones(value []string)() {
    err := m.GetBackingStore().Set("businessPhones", value)
    if err != nil {
        panic(err)
    }
}
// SetChildren sets the children property value. The names of the contact's children.
func (m *Contact) SetChildren(value []string)() {
    err := m.GetBackingStore().Set("children", value)
    if err != nil {
        panic(err)
    }
}
// SetCompanyName sets the companyName property value. The name of the contact's company.
func (m *Contact) SetCompanyName(value *string)() {
    err := m.GetBackingStore().Set("companyName", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartment sets the department property value. The contact's department.
func (m *Contact) SetDepartment(value *string)() {
    err := m.GetBackingStore().Set("department", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The contact's display name. You can specify the display name in a create or update operation. Note that later updates to other properties may cause an automatically generated value to overwrite the displayName value you have specified. To preserve a pre-existing value, always include it as displayName in an update operation.
func (m *Contact) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddresses sets the emailAddresses property value. The contact's email addresses.
func (m *Contact) SetEmailAddresses(value []EmailAddressable)() {
    err := m.GetBackingStore().Set("emailAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the contact. Read-only. Nullable.
func (m *Contact) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetFileAs sets the fileAs property value. The name the contact is filed under.
func (m *Contact) SetFileAs(value *string)() {
    err := m.GetBackingStore().Set("fileAs", value)
    if err != nil {
        panic(err)
    }
}
// SetGeneration sets the generation property value. The contact's suffix.
func (m *Contact) SetGeneration(value *string)() {
    err := m.GetBackingStore().Set("generation", value)
    if err != nil {
        panic(err)
    }
}
// SetGivenName sets the givenName property value. The contact's given name.
func (m *Contact) SetGivenName(value *string)() {
    err := m.GetBackingStore().Set("givenName", value)
    if err != nil {
        panic(err)
    }
}
// SetHomeAddress sets the homeAddress property value. The contact's home address.
func (m *Contact) SetHomeAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("homeAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetHomePhones sets the homePhones property value. The contact's home phone numbers.
func (m *Contact) SetHomePhones(value []string)() {
    err := m.GetBackingStore().Set("homePhones", value)
    if err != nil {
        panic(err)
    }
}
// SetImAddresses sets the imAddresses property value. The contact's instant messaging (IM) addresses.
func (m *Contact) SetImAddresses(value []string)() {
    err := m.GetBackingStore().Set("imAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetInitials sets the initials property value. The contact's initials.
func (m *Contact) SetInitials(value *string)() {
    err := m.GetBackingStore().Set("initials", value)
    if err != nil {
        panic(err)
    }
}
// SetJobTitle sets the jobTitle property value. The contact’s job title.
func (m *Contact) SetJobTitle(value *string)() {
    err := m.GetBackingStore().Set("jobTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetManager sets the manager property value. The name of the contact's manager.
func (m *Contact) SetManager(value *string)() {
    err := m.GetBackingStore().Set("manager", value)
    if err != nil {
        panic(err)
    }
}
// SetMiddleName sets the middleName property value. The contact's middle name.
func (m *Contact) SetMiddleName(value *string)() {
    err := m.GetBackingStore().Set("middleName", value)
    if err != nil {
        panic(err)
    }
}
// SetMobilePhone sets the mobilePhone property value. The contact's mobile phone number.
func (m *Contact) SetMobilePhone(value *string)() {
    err := m.GetBackingStore().Set("mobilePhone", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the contact. Read-only. Nullable.
func (m *Contact) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("multiValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetNickName sets the nickName property value. The contact's nickname.
func (m *Contact) SetNickName(value *string)() {
    err := m.GetBackingStore().Set("nickName", value)
    if err != nil {
        panic(err)
    }
}
// SetOfficeLocation sets the officeLocation property value. The location of the contact's office.
func (m *Contact) SetOfficeLocation(value *string)() {
    err := m.GetBackingStore().Set("officeLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetOtherAddress sets the otherAddress property value. Other addresses for the contact.
func (m *Contact) SetOtherAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("otherAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetParentFolderId sets the parentFolderId property value. The ID of the contact's parent folder.
func (m *Contact) SetParentFolderId(value *string)() {
    err := m.GetBackingStore().Set("parentFolderId", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonalNotes sets the personalNotes property value. The user's notes about the contact.
func (m *Contact) SetPersonalNotes(value *string)() {
    err := m.GetBackingStore().Set("personalNotes", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoto sets the photo property value. Optional contact picture. You can get or set a photo for a contact.
func (m *Contact) SetPhoto(value ProfilePhotoable)() {
    err := m.GetBackingStore().Set("photo", value)
    if err != nil {
        panic(err)
    }
}
// SetProfession sets the profession property value. The contact's profession.
func (m *Contact) SetProfession(value *string)() {
    err := m.GetBackingStore().Set("profession", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the contact. Read-only. Nullable.
func (m *Contact) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("singleValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetSpouseName sets the spouseName property value. The name of the contact's spouse/partner.
func (m *Contact) SetSpouseName(value *string)() {
    err := m.GetBackingStore().Set("spouseName", value)
    if err != nil {
        panic(err)
    }
}
// SetSurname sets the surname property value. The contact's surname.
func (m *Contact) SetSurname(value *string)() {
    err := m.GetBackingStore().Set("surname", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. The contact's title.
func (m *Contact) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
// SetYomiCompanyName sets the yomiCompanyName property value. The phonetic Japanese company name of the contact.
func (m *Contact) SetYomiCompanyName(value *string)() {
    err := m.GetBackingStore().Set("yomiCompanyName", value)
    if err != nil {
        panic(err)
    }
}
// SetYomiGivenName sets the yomiGivenName property value. The phonetic Japanese given name (first name) of the contact.
func (m *Contact) SetYomiGivenName(value *string)() {
    err := m.GetBackingStore().Set("yomiGivenName", value)
    if err != nil {
        panic(err)
    }
}
// SetYomiSurname sets the yomiSurname property value. The phonetic Japanese surname (last name)  of the contact.
func (m *Contact) SetYomiSurname(value *string)() {
    err := m.GetBackingStore().Set("yomiSurname", value)
    if err != nil {
        panic(err)
    }
}
type Contactable interface {
    OutlookItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssistantName()(*string)
    GetBirthday()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetBusinessAddress()(PhysicalAddressable)
    GetBusinessHomePage()(*string)
    GetBusinessPhones()([]string)
    GetChildren()([]string)
    GetCompanyName()(*string)
    GetDepartment()(*string)
    GetDisplayName()(*string)
    GetEmailAddresses()([]EmailAddressable)
    GetExtensions()([]Extensionable)
    GetFileAs()(*string)
    GetGeneration()(*string)
    GetGivenName()(*string)
    GetHomeAddress()(PhysicalAddressable)
    GetHomePhones()([]string)
    GetImAddresses()([]string)
    GetInitials()(*string)
    GetJobTitle()(*string)
    GetManager()(*string)
    GetMiddleName()(*string)
    GetMobilePhone()(*string)
    GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable)
    GetNickName()(*string)
    GetOfficeLocation()(*string)
    GetOtherAddress()(PhysicalAddressable)
    GetParentFolderId()(*string)
    GetPersonalNotes()(*string)
    GetPhoto()(ProfilePhotoable)
    GetProfession()(*string)
    GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable)
    GetSpouseName()(*string)
    GetSurname()(*string)
    GetTitle()(*string)
    GetYomiCompanyName()(*string)
    GetYomiGivenName()(*string)
    GetYomiSurname()(*string)
    SetAssistantName(value *string)()
    SetBirthday(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetBusinessAddress(value PhysicalAddressable)()
    SetBusinessHomePage(value *string)()
    SetBusinessPhones(value []string)()
    SetChildren(value []string)()
    SetCompanyName(value *string)()
    SetDepartment(value *string)()
    SetDisplayName(value *string)()
    SetEmailAddresses(value []EmailAddressable)()
    SetExtensions(value []Extensionable)()
    SetFileAs(value *string)()
    SetGeneration(value *string)()
    SetGivenName(value *string)()
    SetHomeAddress(value PhysicalAddressable)()
    SetHomePhones(value []string)()
    SetImAddresses(value []string)()
    SetInitials(value *string)()
    SetJobTitle(value *string)()
    SetManager(value *string)()
    SetMiddleName(value *string)()
    SetMobilePhone(value *string)()
    SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)()
    SetNickName(value *string)()
    SetOfficeLocation(value *string)()
    SetOtherAddress(value PhysicalAddressable)()
    SetParentFolderId(value *string)()
    SetPersonalNotes(value *string)()
    SetPhoto(value ProfilePhotoable)()
    SetProfession(value *string)()
    SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)()
    SetSpouseName(value *string)()
    SetSurname(value *string)()
    SetTitle(value *string)()
    SetYomiCompanyName(value *string)()
    SetYomiGivenName(value *string)()
    SetYomiSurname(value *string)()
}

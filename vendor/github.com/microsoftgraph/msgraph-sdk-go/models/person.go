package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Person struct {
    Entity
}
// NewPerson instantiates a new Person and sets the default values.
func NewPerson()(*Person) {
    m := &Person{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePersonFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePersonFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPerson(), nil
}
// GetBirthday gets the birthday property value. The person's birthday.
// returns a *string when successful
func (m *Person) GetBirthday()(*string) {
    val, err := m.GetBackingStore().Get("birthday")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCompanyName gets the companyName property value. The name of the person's company.
// returns a *string when successful
func (m *Person) GetCompanyName()(*string) {
    val, err := m.GetBackingStore().Get("companyName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDepartment gets the department property value. The person's department.
// returns a *string when successful
func (m *Person) GetDepartment()(*string) {
    val, err := m.GetBackingStore().Get("department")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The person's display name.
// returns a *string when successful
func (m *Person) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *Person) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["birthday"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBirthday(val)
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
    res["imAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImAddress(val)
        }
        return nil
    }
    res["isFavorite"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFavorite(val)
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
    res["personNotes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonNotes(val)
        }
        return nil
    }
    res["personType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePersonTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonType(val.(PersonTypeable))
        }
        return nil
    }
    res["phones"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePhoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Phoneable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Phoneable)
                }
            }
            m.SetPhones(res)
        }
        return nil
    }
    res["postalAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Locationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Locationable)
                }
            }
            m.SetPostalAddresses(res)
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
    res["scoredEmailAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateScoredEmailAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ScoredEmailAddressable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ScoredEmailAddressable)
                }
            }
            m.SetScoredEmailAddresses(res)
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
    res["websites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWebsiteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Websiteable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Websiteable)
                }
            }
            m.SetWebsites(res)
        }
        return nil
    }
    res["yomiCompany"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetYomiCompany(val)
        }
        return nil
    }
    return res
}
// GetGivenName gets the givenName property value. The person's given name.
// returns a *string when successful
func (m *Person) GetGivenName()(*string) {
    val, err := m.GetBackingStore().Get("givenName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetImAddress gets the imAddress property value. The instant message voice over IP (VOIP) session initiation protocol (SIP) address for the user. Read-only.
// returns a *string when successful
func (m *Person) GetImAddress()(*string) {
    val, err := m.GetBackingStore().Get("imAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsFavorite gets the isFavorite property value. True if the user has flagged this person as a favorite.
// returns a *bool when successful
func (m *Person) GetIsFavorite()(*bool) {
    val, err := m.GetBackingStore().Get("isFavorite")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJobTitle gets the jobTitle property value. The person's job title.
// returns a *string when successful
func (m *Person) GetJobTitle()(*string) {
    val, err := m.GetBackingStore().Get("jobTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOfficeLocation gets the officeLocation property value. The location of the person's office.
// returns a *string when successful
func (m *Person) GetOfficeLocation()(*string) {
    val, err := m.GetBackingStore().Get("officeLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersonNotes gets the personNotes property value. Free-form notes that the user has taken about this person.
// returns a *string when successful
func (m *Person) GetPersonNotes()(*string) {
    val, err := m.GetBackingStore().Get("personNotes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersonType gets the personType property value. The type of person.
// returns a PersonTypeable when successful
func (m *Person) GetPersonType()(PersonTypeable) {
    val, err := m.GetBackingStore().Get("personType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PersonTypeable)
    }
    return nil
}
// GetPhones gets the phones property value. The person's phone numbers.
// returns a []Phoneable when successful
func (m *Person) GetPhones()([]Phoneable) {
    val, err := m.GetBackingStore().Get("phones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Phoneable)
    }
    return nil
}
// GetPostalAddresses gets the postalAddresses property value. The person's addresses.
// returns a []Locationable when successful
func (m *Person) GetPostalAddresses()([]Locationable) {
    val, err := m.GetBackingStore().Get("postalAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Locationable)
    }
    return nil
}
// GetProfession gets the profession property value. The person's profession.
// returns a *string when successful
func (m *Person) GetProfession()(*string) {
    val, err := m.GetBackingStore().Get("profession")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScoredEmailAddresses gets the scoredEmailAddresses property value. The person's email addresses.
// returns a []ScoredEmailAddressable when successful
func (m *Person) GetScoredEmailAddresses()([]ScoredEmailAddressable) {
    val, err := m.GetBackingStore().Get("scoredEmailAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ScoredEmailAddressable)
    }
    return nil
}
// GetSurname gets the surname property value. The person's surname.
// returns a *string when successful
func (m *Person) GetSurname()(*string) {
    val, err := m.GetBackingStore().Get("surname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. The user principal name (UPN) of the person. The UPN is an Internet-style login name for the person based on the Internet standard RFC 822. By convention, this should map to the person's email name. The general format is alias@domain.
// returns a *string when successful
func (m *Person) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebsites gets the websites property value. The person's websites.
// returns a []Websiteable when successful
func (m *Person) GetWebsites()([]Websiteable) {
    val, err := m.GetBackingStore().Get("websites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Websiteable)
    }
    return nil
}
// GetYomiCompany gets the yomiCompany property value. The phonetic Japanese name of the person's company.
// returns a *string when successful
func (m *Person) GetYomiCompany()(*string) {
    val, err := m.GetBackingStore().Get("yomiCompany")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Person) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("birthday", m.GetBirthday())
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
    {
        err = writer.WriteStringValue("givenName", m.GetGivenName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("imAddress", m.GetImAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFavorite", m.GetIsFavorite())
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
        err = writer.WriteStringValue("officeLocation", m.GetOfficeLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("personNotes", m.GetPersonNotes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("personType", m.GetPersonType())
        if err != nil {
            return err
        }
    }
    if m.GetPhones() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPhones()))
        for i, v := range m.GetPhones() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("phones", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPostalAddresses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPostalAddresses()))
        for i, v := range m.GetPostalAddresses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("postalAddresses", cast)
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
    if m.GetScoredEmailAddresses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScoredEmailAddresses()))
        for i, v := range m.GetScoredEmailAddresses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scoredEmailAddresses", cast)
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
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    if m.GetWebsites() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWebsites()))
        for i, v := range m.GetWebsites() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("websites", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("yomiCompany", m.GetYomiCompany())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBirthday sets the birthday property value. The person's birthday.
func (m *Person) SetBirthday(value *string)() {
    err := m.GetBackingStore().Set("birthday", value)
    if err != nil {
        panic(err)
    }
}
// SetCompanyName sets the companyName property value. The name of the person's company.
func (m *Person) SetCompanyName(value *string)() {
    err := m.GetBackingStore().Set("companyName", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartment sets the department property value. The person's department.
func (m *Person) SetDepartment(value *string)() {
    err := m.GetBackingStore().Set("department", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The person's display name.
func (m *Person) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetGivenName sets the givenName property value. The person's given name.
func (m *Person) SetGivenName(value *string)() {
    err := m.GetBackingStore().Set("givenName", value)
    if err != nil {
        panic(err)
    }
}
// SetImAddress sets the imAddress property value. The instant message voice over IP (VOIP) session initiation protocol (SIP) address for the user. Read-only.
func (m *Person) SetImAddress(value *string)() {
    err := m.GetBackingStore().Set("imAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFavorite sets the isFavorite property value. True if the user has flagged this person as a favorite.
func (m *Person) SetIsFavorite(value *bool)() {
    err := m.GetBackingStore().Set("isFavorite", value)
    if err != nil {
        panic(err)
    }
}
// SetJobTitle sets the jobTitle property value. The person's job title.
func (m *Person) SetJobTitle(value *string)() {
    err := m.GetBackingStore().Set("jobTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetOfficeLocation sets the officeLocation property value. The location of the person's office.
func (m *Person) SetOfficeLocation(value *string)() {
    err := m.GetBackingStore().Set("officeLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonNotes sets the personNotes property value. Free-form notes that the user has taken about this person.
func (m *Person) SetPersonNotes(value *string)() {
    err := m.GetBackingStore().Set("personNotes", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonType sets the personType property value. The type of person.
func (m *Person) SetPersonType(value PersonTypeable)() {
    err := m.GetBackingStore().Set("personType", value)
    if err != nil {
        panic(err)
    }
}
// SetPhones sets the phones property value. The person's phone numbers.
func (m *Person) SetPhones(value []Phoneable)() {
    err := m.GetBackingStore().Set("phones", value)
    if err != nil {
        panic(err)
    }
}
// SetPostalAddresses sets the postalAddresses property value. The person's addresses.
func (m *Person) SetPostalAddresses(value []Locationable)() {
    err := m.GetBackingStore().Set("postalAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetProfession sets the profession property value. The person's profession.
func (m *Person) SetProfession(value *string)() {
    err := m.GetBackingStore().Set("profession", value)
    if err != nil {
        panic(err)
    }
}
// SetScoredEmailAddresses sets the scoredEmailAddresses property value. The person's email addresses.
func (m *Person) SetScoredEmailAddresses(value []ScoredEmailAddressable)() {
    err := m.GetBackingStore().Set("scoredEmailAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetSurname sets the surname property value. The person's surname.
func (m *Person) SetSurname(value *string)() {
    err := m.GetBackingStore().Set("surname", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. The user principal name (UPN) of the person. The UPN is an Internet-style login name for the person based on the Internet standard RFC 822. By convention, this should map to the person's email name. The general format is alias@domain.
func (m *Person) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetWebsites sets the websites property value. The person's websites.
func (m *Person) SetWebsites(value []Websiteable)() {
    err := m.GetBackingStore().Set("websites", value)
    if err != nil {
        panic(err)
    }
}
// SetYomiCompany sets the yomiCompany property value. The phonetic Japanese name of the person's company.
func (m *Person) SetYomiCompany(value *string)() {
    err := m.GetBackingStore().Set("yomiCompany", value)
    if err != nil {
        panic(err)
    }
}
type Personable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBirthday()(*string)
    GetCompanyName()(*string)
    GetDepartment()(*string)
    GetDisplayName()(*string)
    GetGivenName()(*string)
    GetImAddress()(*string)
    GetIsFavorite()(*bool)
    GetJobTitle()(*string)
    GetOfficeLocation()(*string)
    GetPersonNotes()(*string)
    GetPersonType()(PersonTypeable)
    GetPhones()([]Phoneable)
    GetPostalAddresses()([]Locationable)
    GetProfession()(*string)
    GetScoredEmailAddresses()([]ScoredEmailAddressable)
    GetSurname()(*string)
    GetUserPrincipalName()(*string)
    GetWebsites()([]Websiteable)
    GetYomiCompany()(*string)
    SetBirthday(value *string)()
    SetCompanyName(value *string)()
    SetDepartment(value *string)()
    SetDisplayName(value *string)()
    SetGivenName(value *string)()
    SetImAddress(value *string)()
    SetIsFavorite(value *bool)()
    SetJobTitle(value *string)()
    SetOfficeLocation(value *string)()
    SetPersonNotes(value *string)()
    SetPersonType(value PersonTypeable)()
    SetPhones(value []Phoneable)()
    SetPostalAddresses(value []Locationable)()
    SetProfession(value *string)()
    SetScoredEmailAddresses(value []ScoredEmailAddressable)()
    SetSurname(value *string)()
    SetUserPrincipalName(value *string)()
    SetWebsites(value []Websiteable)()
    SetYomiCompany(value *string)()
}

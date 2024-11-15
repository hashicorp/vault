package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationSchool struct {
    EducationOrganization
}
// NewEducationSchool instantiates a new EducationSchool and sets the default values.
func NewEducationSchool()(*EducationSchool) {
    m := &EducationSchool{
        EducationOrganization: *NewEducationOrganization(),
    }
    odataTypeValue := "#microsoft.graph.educationSchool"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationSchoolFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationSchoolFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationSchool(), nil
}
// GetAddress gets the address property value. Address of the school.
// returns a PhysicalAddressable when successful
func (m *EducationSchool) GetAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("address")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetAdministrativeUnit gets the administrativeUnit property value. The underlying administrativeUnit for this school.
// returns a AdministrativeUnitable when successful
func (m *EducationSchool) GetAdministrativeUnit()(AdministrativeUnitable) {
    val, err := m.GetBackingStore().Get("administrativeUnit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AdministrativeUnitable)
    }
    return nil
}
// GetClasses gets the classes property value. Classes taught at the school. Nullable.
// returns a []EducationClassable when successful
func (m *EducationSchool) GetClasses()([]EducationClassable) {
    val, err := m.GetBackingStore().Get("classes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationClassable)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Entity who created the school.
// returns a IdentitySetable when successful
func (m *EducationSchool) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetExternalId gets the externalId property value. ID of school in syncing system.
// returns a *string when successful
func (m *EducationSchool) GetExternalId()(*string) {
    val, err := m.GetBackingStore().Get("externalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalPrincipalId gets the externalPrincipalId property value. ID of principal in syncing system.
// returns a *string when successful
func (m *EducationSchool) GetExternalPrincipalId()(*string) {
    val, err := m.GetBackingStore().Get("externalPrincipalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFax gets the fax property value. The fax property
// returns a *string when successful
func (m *EducationSchool) GetFax()(*string) {
    val, err := m.GetBackingStore().Get("fax")
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
func (m *EducationSchool) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationOrganization.GetFieldDeserializers()
    res["address"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["administrativeUnit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAdministrativeUnitFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAdministrativeUnit(val.(AdministrativeUnitable))
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
    res["externalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalId(val)
        }
        return nil
    }
    res["externalPrincipalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalPrincipalId(val)
        }
        return nil
    }
    res["fax"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFax(val)
        }
        return nil
    }
    res["highestGrade"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHighestGrade(val)
        }
        return nil
    }
    res["lowestGrade"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowestGrade(val)
        }
        return nil
    }
    res["phone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhone(val)
        }
        return nil
    }
    res["principalEmail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipalEmail(val)
        }
        return nil
    }
    res["principalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipalName(val)
        }
        return nil
    }
    res["schoolNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchoolNumber(val)
        }
        return nil
    }
    res["users"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationUserable)
                }
            }
            m.SetUsers(res)
        }
        return nil
    }
    return res
}
// GetHighestGrade gets the highestGrade property value. Highest grade taught.
// returns a *string when successful
func (m *EducationSchool) GetHighestGrade()(*string) {
    val, err := m.GetBackingStore().Get("highestGrade")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLowestGrade gets the lowestGrade property value. Lowest grade taught.
// returns a *string when successful
func (m *EducationSchool) GetLowestGrade()(*string) {
    val, err := m.GetBackingStore().Get("lowestGrade")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhone gets the phone property value. Phone number of school.
// returns a *string when successful
func (m *EducationSchool) GetPhone()(*string) {
    val, err := m.GetBackingStore().Get("phone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrincipalEmail gets the principalEmail property value. Email address of the principal.
// returns a *string when successful
func (m *EducationSchool) GetPrincipalEmail()(*string) {
    val, err := m.GetBackingStore().Get("principalEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrincipalName gets the principalName property value. Name of the principal.
// returns a *string when successful
func (m *EducationSchool) GetPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("principalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSchoolNumber gets the schoolNumber property value. School Number.
// returns a *string when successful
func (m *EducationSchool) GetSchoolNumber()(*string) {
    val, err := m.GetBackingStore().Get("schoolNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUsers gets the users property value. Users in the school. Nullable.
// returns a []EducationUserable when successful
func (m *EducationSchool) GetUsers()([]EducationUserable) {
    val, err := m.GetBackingStore().Get("users")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationUserable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationSchool) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationOrganization.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("address", m.GetAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("administrativeUnit", m.GetAdministrativeUnit())
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
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalPrincipalId", m.GetExternalPrincipalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("fax", m.GetFax())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("highestGrade", m.GetHighestGrade())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lowestGrade", m.GetLowestGrade())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("phone", m.GetPhone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalEmail", m.GetPrincipalEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalName", m.GetPrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("schoolNumber", m.GetSchoolNumber())
        if err != nil {
            return err
        }
    }
    if m.GetUsers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUsers()))
        for i, v := range m.GetUsers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("users", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAddress sets the address property value. Address of the school.
func (m *EducationSchool) SetAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("address", value)
    if err != nil {
        panic(err)
    }
}
// SetAdministrativeUnit sets the administrativeUnit property value. The underlying administrativeUnit for this school.
func (m *EducationSchool) SetAdministrativeUnit(value AdministrativeUnitable)() {
    err := m.GetBackingStore().Set("administrativeUnit", value)
    if err != nil {
        panic(err)
    }
}
// SetClasses sets the classes property value. Classes taught at the school. Nullable.
func (m *EducationSchool) SetClasses(value []EducationClassable)() {
    err := m.GetBackingStore().Set("classes", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Entity who created the school.
func (m *EducationSchool) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalId sets the externalId property value. ID of school in syncing system.
func (m *EducationSchool) SetExternalId(value *string)() {
    err := m.GetBackingStore().Set("externalId", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalPrincipalId sets the externalPrincipalId property value. ID of principal in syncing system.
func (m *EducationSchool) SetExternalPrincipalId(value *string)() {
    err := m.GetBackingStore().Set("externalPrincipalId", value)
    if err != nil {
        panic(err)
    }
}
// SetFax sets the fax property value. The fax property
func (m *EducationSchool) SetFax(value *string)() {
    err := m.GetBackingStore().Set("fax", value)
    if err != nil {
        panic(err)
    }
}
// SetHighestGrade sets the highestGrade property value. Highest grade taught.
func (m *EducationSchool) SetHighestGrade(value *string)() {
    err := m.GetBackingStore().Set("highestGrade", value)
    if err != nil {
        panic(err)
    }
}
// SetLowestGrade sets the lowestGrade property value. Lowest grade taught.
func (m *EducationSchool) SetLowestGrade(value *string)() {
    err := m.GetBackingStore().Set("lowestGrade", value)
    if err != nil {
        panic(err)
    }
}
// SetPhone sets the phone property value. Phone number of school.
func (m *EducationSchool) SetPhone(value *string)() {
    err := m.GetBackingStore().Set("phone", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalEmail sets the principalEmail property value. Email address of the principal.
func (m *EducationSchool) SetPrincipalEmail(value *string)() {
    err := m.GetBackingStore().Set("principalEmail", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalName sets the principalName property value. Name of the principal.
func (m *EducationSchool) SetPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("principalName", value)
    if err != nil {
        panic(err)
    }
}
// SetSchoolNumber sets the schoolNumber property value. School Number.
func (m *EducationSchool) SetSchoolNumber(value *string)() {
    err := m.GetBackingStore().Set("schoolNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetUsers sets the users property value. Users in the school. Nullable.
func (m *EducationSchool) SetUsers(value []EducationUserable)() {
    err := m.GetBackingStore().Set("users", value)
    if err != nil {
        panic(err)
    }
}
type EducationSchoolable interface {
    EducationOrganizationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAddress()(PhysicalAddressable)
    GetAdministrativeUnit()(AdministrativeUnitable)
    GetClasses()([]EducationClassable)
    GetCreatedBy()(IdentitySetable)
    GetExternalId()(*string)
    GetExternalPrincipalId()(*string)
    GetFax()(*string)
    GetHighestGrade()(*string)
    GetLowestGrade()(*string)
    GetPhone()(*string)
    GetPrincipalEmail()(*string)
    GetPrincipalName()(*string)
    GetSchoolNumber()(*string)
    GetUsers()([]EducationUserable)
    SetAddress(value PhysicalAddressable)()
    SetAdministrativeUnit(value AdministrativeUnitable)()
    SetClasses(value []EducationClassable)()
    SetCreatedBy(value IdentitySetable)()
    SetExternalId(value *string)()
    SetExternalPrincipalId(value *string)()
    SetFax(value *string)()
    SetHighestGrade(value *string)()
    SetLowestGrade(value *string)()
    SetPhone(value *string)()
    SetPrincipalEmail(value *string)()
    SetPrincipalName(value *string)()
    SetSchoolNumber(value *string)()
    SetUsers(value []EducationUserable)()
}

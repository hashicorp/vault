package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageSubject struct {
    Entity
}
// NewAccessPackageSubject instantiates a new AccessPackageSubject and sets the default values.
func NewAccessPackageSubject()(*AccessPackageSubject) {
    m := &AccessPackageSubject{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageSubjectFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageSubjectFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageSubject(), nil
}
// GetConnectedOrganization gets the connectedOrganization property value. The connected organization of the subject. Read-only. Nullable.
// returns a ConnectedOrganizationable when successful
func (m *AccessPackageSubject) GetConnectedOrganization()(ConnectedOrganizationable) {
    val, err := m.GetBackingStore().Get("connectedOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConnectedOrganizationable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the subject.
// returns a *string when successful
func (m *AccessPackageSubject) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmail gets the email property value. The email address of the subject.
// returns a *string when successful
func (m *AccessPackageSubject) GetEmail()(*string) {
    val, err := m.GetBackingStore().Get("email")
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
func (m *AccessPackageSubject) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["connectedOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConnectedOrganizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectedOrganization(val.(ConnectedOrganizationable))
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
    res["email"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmail(val)
        }
        return nil
    }
    res["objectId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetObjectId(val)
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
    res["subjectType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessPackageSubjectType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubjectType(val.(*AccessPackageSubjectType))
        }
        return nil
    }
    return res
}
// GetObjectId gets the objectId property value. The object identifier of the subject. null if the subject isn't yet a user in the tenant.
// returns a *string when successful
func (m *AccessPackageSubject) GetObjectId()(*string) {
    val, err := m.GetBackingStore().Get("objectId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnPremisesSecurityIdentifier gets the onPremisesSecurityIdentifier property value. A string representation of the principal's security identifier, if known, or null if the subject doesn't have a security identifier.
// returns a *string when successful
func (m *AccessPackageSubject) GetOnPremisesSecurityIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("onPremisesSecurityIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrincipalName gets the principalName property value. The principal name, if known, of the subject.
// returns a *string when successful
func (m *AccessPackageSubject) GetPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("principalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubjectType gets the subjectType property value. The resource type of the subject. The possible values are: notSpecified, user, servicePrincipal, unknownFutureValue.
// returns a *AccessPackageSubjectType when successful
func (m *AccessPackageSubject) GetSubjectType()(*AccessPackageSubjectType) {
    val, err := m.GetBackingStore().Get("subjectType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageSubjectType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageSubject) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("connectedOrganization", m.GetConnectedOrganization())
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
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("objectId", m.GetObjectId())
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
        err = writer.WriteStringValue("principalName", m.GetPrincipalName())
        if err != nil {
            return err
        }
    }
    if m.GetSubjectType() != nil {
        cast := (*m.GetSubjectType()).String()
        err = writer.WriteStringValue("subjectType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConnectedOrganization sets the connectedOrganization property value. The connected organization of the subject. Read-only. Nullable.
func (m *AccessPackageSubject) SetConnectedOrganization(value ConnectedOrganizationable)() {
    err := m.GetBackingStore().Set("connectedOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the subject.
func (m *AccessPackageSubject) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmail sets the email property value. The email address of the subject.
func (m *AccessPackageSubject) SetEmail(value *string)() {
    err := m.GetBackingStore().Set("email", value)
    if err != nil {
        panic(err)
    }
}
// SetObjectId sets the objectId property value. The object identifier of the subject. null if the subject isn't yet a user in the tenant.
func (m *AccessPackageSubject) SetObjectId(value *string)() {
    err := m.GetBackingStore().Set("objectId", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesSecurityIdentifier sets the onPremisesSecurityIdentifier property value. A string representation of the principal's security identifier, if known, or null if the subject doesn't have a security identifier.
func (m *AccessPackageSubject) SetOnPremisesSecurityIdentifier(value *string)() {
    err := m.GetBackingStore().Set("onPremisesSecurityIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalName sets the principalName property value. The principal name, if known, of the subject.
func (m *AccessPackageSubject) SetPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("principalName", value)
    if err != nil {
        panic(err)
    }
}
// SetSubjectType sets the subjectType property value. The resource type of the subject. The possible values are: notSpecified, user, servicePrincipal, unknownFutureValue.
func (m *AccessPackageSubject) SetSubjectType(value *AccessPackageSubjectType)() {
    err := m.GetBackingStore().Set("subjectType", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageSubjectable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConnectedOrganization()(ConnectedOrganizationable)
    GetDisplayName()(*string)
    GetEmail()(*string)
    GetObjectId()(*string)
    GetOnPremisesSecurityIdentifier()(*string)
    GetPrincipalName()(*string)
    GetSubjectType()(*AccessPackageSubjectType)
    SetConnectedOrganization(value ConnectedOrganizationable)()
    SetDisplayName(value *string)()
    SetEmail(value *string)()
    SetObjectId(value *string)()
    SetOnPremisesSecurityIdentifier(value *string)()
    SetPrincipalName(value *string)()
    SetSubjectType(value *AccessPackageSubjectType)()
}

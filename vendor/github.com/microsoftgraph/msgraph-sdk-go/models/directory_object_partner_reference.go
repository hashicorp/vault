package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DirectoryObjectPartnerReference struct {
    DirectoryObject
}
// NewDirectoryObjectPartnerReference instantiates a new DirectoryObjectPartnerReference and sets the default values.
func NewDirectoryObjectPartnerReference()(*DirectoryObjectPartnerReference) {
    m := &DirectoryObjectPartnerReference{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.directoryObjectPartnerReference"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDirectoryObjectPartnerReferenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectoryObjectPartnerReferenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryObjectPartnerReference(), nil
}
// GetDescription gets the description property value. Description of the object returned. Read-only.
// returns a *string when successful
func (m *DirectoryObjectPartnerReference) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of directory object being returned, like group or application. Read-only.
// returns a *string when successful
func (m *DirectoryObjectPartnerReference) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalPartnerTenantId gets the externalPartnerTenantId property value. The tenant identifier for the partner tenant. Read-only.
// returns a *UUID when successful
func (m *DirectoryObjectPartnerReference) GetExternalPartnerTenantId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("externalPartnerTenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DirectoryObjectPartnerReference) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["externalPartnerTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalPartnerTenantId(val)
        }
        return nil
    }
    res["objectType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetObjectType(val)
        }
        return nil
    }
    return res
}
// GetObjectType gets the objectType property value. The type of the referenced object in the partner tenant. Read-only.
// returns a *string when successful
func (m *DirectoryObjectPartnerReference) GetObjectType()(*string) {
    val, err := m.GetBackingStore().Get("objectType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DirectoryObjectPartnerReference) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteUUIDValue("externalPartnerTenantId", m.GetExternalPartnerTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("objectType", m.GetObjectType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Description of the object returned. Read-only.
func (m *DirectoryObjectPartnerReference) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of directory object being returned, like group or application. Read-only.
func (m *DirectoryObjectPartnerReference) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalPartnerTenantId sets the externalPartnerTenantId property value. The tenant identifier for the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) SetExternalPartnerTenantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("externalPartnerTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetObjectType sets the objectType property value. The type of the referenced object in the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) SetObjectType(value *string)() {
    err := m.GetBackingStore().Set("objectType", value)
    if err != nil {
        panic(err)
    }
}
type DirectoryObjectPartnerReferenceable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExternalPartnerTenantId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetObjectType()(*string)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExternalPartnerTenantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetObjectType(value *string)()
}

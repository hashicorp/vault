package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ResourceSpecificPermissionGrant struct {
    DirectoryObject
}
// NewResourceSpecificPermissionGrant instantiates a new ResourceSpecificPermissionGrant and sets the default values.
func NewResourceSpecificPermissionGrant()(*ResourceSpecificPermissionGrant) {
    m := &ResourceSpecificPermissionGrant{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.resourceSpecificPermissionGrant"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateResourceSpecificPermissionGrantFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateResourceSpecificPermissionGrantFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResourceSpecificPermissionGrant(), nil
}
// GetClientAppId gets the clientAppId property value. ID of the service principal of the Microsoft Entra app that has been granted access. Read-only.
// returns a *string when successful
func (m *ResourceSpecificPermissionGrant) GetClientAppId()(*string) {
    val, err := m.GetBackingStore().Get("clientAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetClientId gets the clientId property value. ID of the Microsoft Entra app that has been granted access. Read-only.
// returns a *string when successful
func (m *ResourceSpecificPermissionGrant) GetClientId()(*string) {
    val, err := m.GetBackingStore().Get("clientId")
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
func (m *ResourceSpecificPermissionGrant) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["clientAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientAppId(val)
        }
        return nil
    }
    res["clientId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientId(val)
        }
        return nil
    }
    res["permission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermission(val)
        }
        return nil
    }
    res["permissionType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermissionType(val)
        }
        return nil
    }
    res["resourceAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceAppId(val)
        }
        return nil
    }
    return res
}
// GetPermission gets the permission property value. The name of the resource-specific permission. Read-only.
// returns a *string when successful
func (m *ResourceSpecificPermissionGrant) GetPermission()(*string) {
    val, err := m.GetBackingStore().Get("permission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPermissionType gets the permissionType property value. The type of permission. Possible values are: Application, Delegated. Read-only.
// returns a *string when successful
func (m *ResourceSpecificPermissionGrant) GetPermissionType()(*string) {
    val, err := m.GetBackingStore().Get("permissionType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceAppId gets the resourceAppId property value. ID of the Microsoft Entra app that is hosting the resource. Read-only.
// returns a *string when successful
func (m *ResourceSpecificPermissionGrant) GetResourceAppId()(*string) {
    val, err := m.GetBackingStore().Get("resourceAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ResourceSpecificPermissionGrant) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("clientAppId", m.GetClientAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("clientId", m.GetClientId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("permission", m.GetPermission())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("permissionType", m.GetPermissionType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceAppId", m.GetResourceAppId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClientAppId sets the clientAppId property value. ID of the service principal of the Microsoft Entra app that has been granted access. Read-only.
func (m *ResourceSpecificPermissionGrant) SetClientAppId(value *string)() {
    err := m.GetBackingStore().Set("clientAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetClientId sets the clientId property value. ID of the Microsoft Entra app that has been granted access. Read-only.
func (m *ResourceSpecificPermissionGrant) SetClientId(value *string)() {
    err := m.GetBackingStore().Set("clientId", value)
    if err != nil {
        panic(err)
    }
}
// SetPermission sets the permission property value. The name of the resource-specific permission. Read-only.
func (m *ResourceSpecificPermissionGrant) SetPermission(value *string)() {
    err := m.GetBackingStore().Set("permission", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionType sets the permissionType property value. The type of permission. Possible values are: Application, Delegated. Read-only.
func (m *ResourceSpecificPermissionGrant) SetPermissionType(value *string)() {
    err := m.GetBackingStore().Set("permissionType", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceAppId sets the resourceAppId property value. ID of the Microsoft Entra app that is hosting the resource. Read-only.
func (m *ResourceSpecificPermissionGrant) SetResourceAppId(value *string)() {
    err := m.GetBackingStore().Set("resourceAppId", value)
    if err != nil {
        panic(err)
    }
}
type ResourceSpecificPermissionGrantable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClientAppId()(*string)
    GetClientId()(*string)
    GetPermission()(*string)
    GetPermissionType()(*string)
    GetResourceAppId()(*string)
    SetClientAppId(value *string)()
    SetClientId(value *string)()
    SetPermission(value *string)()
    SetPermissionType(value *string)()
    SetResourceAppId(value *string)()
}

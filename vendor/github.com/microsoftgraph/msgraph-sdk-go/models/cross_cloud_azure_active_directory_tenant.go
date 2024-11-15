package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CrossCloudAzureActiveDirectoryTenant struct {
    IdentitySource
}
// NewCrossCloudAzureActiveDirectoryTenant instantiates a new CrossCloudAzureActiveDirectoryTenant and sets the default values.
func NewCrossCloudAzureActiveDirectoryTenant()(*CrossCloudAzureActiveDirectoryTenant) {
    m := &CrossCloudAzureActiveDirectoryTenant{
        IdentitySource: *NewIdentitySource(),
    }
    odataTypeValue := "#microsoft.graph.crossCloudAzureActiveDirectoryTenant"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCrossCloudAzureActiveDirectoryTenantFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCrossCloudAzureActiveDirectoryTenantFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCrossCloudAzureActiveDirectoryTenant(), nil
}
// GetCloudInstance gets the cloudInstance property value. The ID of the cloud where the tenant is located, one of microsoftonline.com, microsoftonline.us or partner.microsoftonline.cn. Read only.
// returns a *string when successful
func (m *CrossCloudAzureActiveDirectoryTenant) GetCloudInstance()(*string) {
    val, err := m.GetBackingStore().Get("cloudInstance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the Microsoft Entra tenant. Read only.
// returns a *string when successful
func (m *CrossCloudAzureActiveDirectoryTenant) GetDisplayName()(*string) {
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
func (m *CrossCloudAzureActiveDirectoryTenant) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentitySource.GetFieldDeserializers()
    res["cloudInstance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudInstance(val)
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
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    return res
}
// GetTenantId gets the tenantId property value. The ID of the Microsoft Entra tenant. Read only.
// returns a *string when successful
func (m *CrossCloudAzureActiveDirectoryTenant) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CrossCloudAzureActiveDirectoryTenant) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentitySource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("cloudInstance", m.GetCloudInstance())
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
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCloudInstance sets the cloudInstance property value. The ID of the cloud where the tenant is located, one of microsoftonline.com, microsoftonline.us or partner.microsoftonline.cn. Read only.
func (m *CrossCloudAzureActiveDirectoryTenant) SetCloudInstance(value *string)() {
    err := m.GetBackingStore().Set("cloudInstance", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the Microsoft Entra tenant. Read only.
func (m *CrossCloudAzureActiveDirectoryTenant) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The ID of the Microsoft Entra tenant. Read only.
func (m *CrossCloudAzureActiveDirectoryTenant) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
type CrossCloudAzureActiveDirectoryTenantable interface {
    IdentitySourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCloudInstance()(*string)
    GetDisplayName()(*string)
    GetTenantId()(*string)
    SetCloudInstance(value *string)()
    SetDisplayName(value *string)()
    SetTenantId(value *string)()
}

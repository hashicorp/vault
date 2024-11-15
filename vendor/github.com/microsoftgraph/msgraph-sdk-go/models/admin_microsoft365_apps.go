package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AdminMicrosoft365Apps struct {
    Entity
}
// NewAdminMicrosoft365Apps instantiates a new AdminMicrosoft365Apps and sets the default values.
func NewAdminMicrosoft365Apps()(*AdminMicrosoft365Apps) {
    m := &AdminMicrosoft365Apps{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAdminMicrosoft365AppsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAdminMicrosoft365AppsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdminMicrosoft365Apps(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AdminMicrosoft365Apps) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["installationOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateM365AppsInstallationOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallationOptions(val.(M365AppsInstallationOptionsable))
        }
        return nil
    }
    return res
}
// GetInstallationOptions gets the installationOptions property value. A container for tenant-level settings for Microsoft 365 applications.
// returns a M365AppsInstallationOptionsable when successful
func (m *AdminMicrosoft365Apps) GetInstallationOptions()(M365AppsInstallationOptionsable) {
    val, err := m.GetBackingStore().Get("installationOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(M365AppsInstallationOptionsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AdminMicrosoft365Apps) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("installationOptions", m.GetInstallationOptions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInstallationOptions sets the installationOptions property value. A container for tenant-level settings for Microsoft 365 applications.
func (m *AdminMicrosoft365Apps) SetInstallationOptions(value M365AppsInstallationOptionsable)() {
    err := m.GetBackingStore().Set("installationOptions", value)
    if err != nil {
        panic(err)
    }
}
type AdminMicrosoft365Appsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInstallationOptions()(M365AppsInstallationOptionsable)
    SetInstallationOptions(value M365AppsInstallationOptionsable)()
}

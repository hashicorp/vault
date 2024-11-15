package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AdminReportSettings struct {
    Entity
}
// NewAdminReportSettings instantiates a new AdminReportSettings and sets the default values.
func NewAdminReportSettings()(*AdminReportSettings) {
    m := &AdminReportSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAdminReportSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAdminReportSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdminReportSettings(), nil
}
// GetDisplayConcealedNames gets the displayConcealedNames property value. If set to true, all reports conceal user information such as usernames, groups, and sites. If false, all reports show identifiable information. This property represents a setting in the Microsoft 365 admin center. Required.
// returns a *bool when successful
func (m *AdminReportSettings) GetDisplayConcealedNames()(*bool) {
    val, err := m.GetBackingStore().Get("displayConcealedNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AdminReportSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayConcealedNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayConcealedNames(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AdminReportSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("displayConcealedNames", m.GetDisplayConcealedNames())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayConcealedNames sets the displayConcealedNames property value. If set to true, all reports conceal user information such as usernames, groups, and sites. If false, all reports show identifiable information. This property represents a setting in the Microsoft 365 admin center. Required.
func (m *AdminReportSettings) SetDisplayConcealedNames(value *bool)() {
    err := m.GetBackingStore().Set("displayConcealedNames", value)
    if err != nil {
        panic(err)
    }
}
type AdminReportSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayConcealedNames()(*bool)
    SetDisplayConcealedNames(value *bool)()
}

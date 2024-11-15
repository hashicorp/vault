package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type InsightsSettings struct {
    Entity
}
// NewInsightsSettings instantiates a new InsightsSettings and sets the default values.
func NewInsightsSettings()(*InsightsSettings) {
    m := &InsightsSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateInsightsSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInsightsSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInsightsSettings(), nil
}
// GetDisabledForGroup gets the disabledForGroup property value. The ID of a Microsoft Entra group, of which the specified type of insights are disabled for its members. The default value is null. Optional.
// returns a *string when successful
func (m *InsightsSettings) GetDisabledForGroup()(*string) {
    val, err := m.GetBackingStore().Get("disabledForGroup")
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
func (m *InsightsSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["disabledForGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisabledForGroup(val)
        }
        return nil
    }
    res["isEnabledInOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabledInOrganization(val)
        }
        return nil
    }
    return res
}
// GetIsEnabledInOrganization gets the isEnabledInOrganization property value. true if insights of the specified type are enabled for the organization; false if insights of the specified type are disabled for all users without exceptions. The default value is true. Optional.
// returns a *bool when successful
func (m *InsightsSettings) GetIsEnabledInOrganization()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabledInOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *InsightsSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("disabledForGroup", m.GetDisabledForGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabledInOrganization", m.GetIsEnabledInOrganization())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisabledForGroup sets the disabledForGroup property value. The ID of a Microsoft Entra group, of which the specified type of insights are disabled for its members. The default value is null. Optional.
func (m *InsightsSettings) SetDisabledForGroup(value *string)() {
    err := m.GetBackingStore().Set("disabledForGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabledInOrganization sets the isEnabledInOrganization property value. true if insights of the specified type are enabled for the organization; false if insights of the specified type are disabled for all users without exceptions. The default value is true. Optional.
func (m *InsightsSettings) SetIsEnabledInOrganization(value *bool)() {
    err := m.GetBackingStore().Set("isEnabledInOrganization", value)
    if err != nil {
        panic(err)
    }
}
type InsightsSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisabledForGroup()(*string)
    GetIsEnabledInOrganization()(*bool)
    SetDisabledForGroup(value *string)()
    SetIsEnabledInOrganization(value *bool)()
}

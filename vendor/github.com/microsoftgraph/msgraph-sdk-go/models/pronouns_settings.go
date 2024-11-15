package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PronounsSettings struct {
    Entity
}
// NewPronounsSettings instantiates a new PronounsSettings and sets the default values.
func NewPronounsSettings()(*PronounsSettings) {
    m := &PronounsSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePronounsSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePronounsSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPronounsSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PronounsSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
// GetIsEnabledInOrganization gets the isEnabledInOrganization property value. true to enable pronouns in the organization; otherwise, false. The default value is false, and pronouns are disabled.
// returns a *bool when successful
func (m *PronounsSettings) GetIsEnabledInOrganization()(*bool) {
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
func (m *PronounsSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isEnabledInOrganization", m.GetIsEnabledInOrganization())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsEnabledInOrganization sets the isEnabledInOrganization property value. true to enable pronouns in the organization; otherwise, false. The default value is false, and pronouns are disabled.
func (m *PronounsSettings) SetIsEnabledInOrganization(value *bool)() {
    err := m.GetBackingStore().Set("isEnabledInOrganization", value)
    if err != nil {
        panic(err)
    }
}
type PronounsSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsEnabledInOrganization()(*bool)
    SetIsEnabledInOrganization(value *bool)()
}

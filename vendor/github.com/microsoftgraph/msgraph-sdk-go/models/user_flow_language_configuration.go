package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserFlowLanguageConfiguration struct {
    Entity
}
// NewUserFlowLanguageConfiguration instantiates a new UserFlowLanguageConfiguration and sets the default values.
func NewUserFlowLanguageConfiguration()(*UserFlowLanguageConfiguration) {
    m := &UserFlowLanguageConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserFlowLanguageConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserFlowLanguageConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserFlowLanguageConfiguration(), nil
}
// GetDefaultPages gets the defaultPages property value. Collection of pages with the default content to display in a user flow for a specified language. This collection doesn't allow any kind of modification.
// returns a []UserFlowLanguagePageable when successful
func (m *UserFlowLanguageConfiguration) GetDefaultPages()([]UserFlowLanguagePageable) {
    val, err := m.GetBackingStore().Get("defaultPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserFlowLanguagePageable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The language name to display. This property is read-only.
// returns a *string when successful
func (m *UserFlowLanguageConfiguration) GetDisplayName()(*string) {
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
func (m *UserFlowLanguageConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["defaultPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserFlowLanguagePageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserFlowLanguagePageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserFlowLanguagePageable)
                }
            }
            m.SetDefaultPages(res)
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
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
        }
        return nil
    }
    res["overridesPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserFlowLanguagePageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserFlowLanguagePageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserFlowLanguagePageable)
                }
            }
            m.SetOverridesPages(res)
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. Indicates whether the language is enabled within the user flow.
// returns a *bool when successful
func (m *UserFlowLanguageConfiguration) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOverridesPages gets the overridesPages property value. Collection of pages with the overrides messages to display in a user flow for a specified language. This collection only allows you to modify the content of the page, any other modification isn't allowed (creation or deletion of pages).
// returns a []UserFlowLanguagePageable when successful
func (m *UserFlowLanguageConfiguration) GetOverridesPages()([]UserFlowLanguagePageable) {
    val, err := m.GetBackingStore().Get("overridesPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserFlowLanguagePageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserFlowLanguageConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDefaultPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDefaultPages()))
        for i, v := range m.GetDefaultPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("defaultPages", cast)
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
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetOverridesPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOverridesPages()))
        for i, v := range m.GetOverridesPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("overridesPages", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefaultPages sets the defaultPages property value. Collection of pages with the default content to display in a user flow for a specified language. This collection doesn't allow any kind of modification.
func (m *UserFlowLanguageConfiguration) SetDefaultPages(value []UserFlowLanguagePageable)() {
    err := m.GetBackingStore().Set("defaultPages", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The language name to display. This property is read-only.
func (m *UserFlowLanguageConfiguration) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Indicates whether the language is enabled within the user flow.
func (m *UserFlowLanguageConfiguration) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOverridesPages sets the overridesPages property value. Collection of pages with the overrides messages to display in a user flow for a specified language. This collection only allows you to modify the content of the page, any other modification isn't allowed (creation or deletion of pages).
func (m *UserFlowLanguageConfiguration) SetOverridesPages(value []UserFlowLanguagePageable)() {
    err := m.GetBackingStore().Set("overridesPages", value)
    if err != nil {
        panic(err)
    }
}
type UserFlowLanguageConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDefaultPages()([]UserFlowLanguagePageable)
    GetDisplayName()(*string)
    GetIsEnabled()(*bool)
    GetOverridesPages()([]UserFlowLanguagePageable)
    SetDefaultPages(value []UserFlowLanguagePageable)()
    SetDisplayName(value *string)()
    SetIsEnabled(value *bool)()
    SetOverridesPages(value []UserFlowLanguagePageable)()
}

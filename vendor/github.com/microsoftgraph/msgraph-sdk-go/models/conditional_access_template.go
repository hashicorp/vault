package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConditionalAccessTemplate struct {
    Entity
}
// NewConditionalAccessTemplate instantiates a new ConditionalAccessTemplate and sets the default values.
func NewConditionalAccessTemplate()(*ConditionalAccessTemplate) {
    m := &ConditionalAccessTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConditionalAccessTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessTemplate(), nil
}
// GetDescription gets the description property value. The user-friendly name of the template.
// returns a *string when successful
func (m *ConditionalAccessTemplate) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetails gets the details property value. The details property
// returns a ConditionalAccessPolicyDetailable when successful
func (m *ConditionalAccessTemplate) GetDetails()(ConditionalAccessPolicyDetailable) {
    val, err := m.GetBackingStore().Get("details")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessPolicyDetailable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["details"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessPolicyDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetails(val.(ConditionalAccessPolicyDetailable))
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["scenarios"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTemplateScenarios)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScenarios(val.(*TemplateScenarios))
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The user-friendly name of the template.
// returns a *string when successful
func (m *ConditionalAccessTemplate) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScenarios gets the scenarios property value. The scenarios property
// returns a *TemplateScenarios when successful
func (m *ConditionalAccessTemplate) GetScenarios()(*TemplateScenarios) {
    val, err := m.GetBackingStore().Get("scenarios")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TemplateScenarios)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
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
        err = writer.WriteObjectValue("details", m.GetDetails())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetScenarios() != nil {
        cast := (*m.GetScenarios()).String()
        err = writer.WriteStringValue("scenarios", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. The user-friendly name of the template.
func (m *ConditionalAccessTemplate) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDetails sets the details property value. The details property
func (m *ConditionalAccessTemplate) SetDetails(value ConditionalAccessPolicyDetailable)() {
    err := m.GetBackingStore().Set("details", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The user-friendly name of the template.
func (m *ConditionalAccessTemplate) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetScenarios sets the scenarios property value. The scenarios property
func (m *ConditionalAccessTemplate) SetScenarios(value *TemplateScenarios)() {
    err := m.GetBackingStore().Set("scenarios", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessTemplateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDetails()(ConditionalAccessPolicyDetailable)
    GetName()(*string)
    GetScenarios()(*TemplateScenarios)
    SetDescription(value *string)()
    SetDetails(value ConditionalAccessPolicyDetailable)()
    SetName(value *string)()
    SetScenarios(value *TemplateScenarios)()
}

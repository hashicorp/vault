package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceComplianceScheduledActionForRule scheduled Action for Rule
type DeviceComplianceScheduledActionForRule struct {
    Entity
}
// NewDeviceComplianceScheduledActionForRule instantiates a new DeviceComplianceScheduledActionForRule and sets the default values.
func NewDeviceComplianceScheduledActionForRule()(*DeviceComplianceScheduledActionForRule) {
    m := &DeviceComplianceScheduledActionForRule{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceComplianceScheduledActionForRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceComplianceScheduledActionForRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceComplianceScheduledActionForRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceComplianceScheduledActionForRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["ruleName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRuleName(val)
        }
        return nil
    }
    res["scheduledActionConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceComplianceActionItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceComplianceActionItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceComplianceActionItemable)
                }
            }
            m.SetScheduledActionConfigurations(res)
        }
        return nil
    }
    return res
}
// GetRuleName gets the ruleName property value. Name of the rule which this scheduled action applies to. Currently scheduled actions are created per policy instead of per rule, thus RuleName is always set to default value PasswordRequired.
// returns a *string when successful
func (m *DeviceComplianceScheduledActionForRule) GetRuleName()(*string) {
    val, err := m.GetBackingStore().Get("ruleName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScheduledActionConfigurations gets the scheduledActionConfigurations property value. The list of scheduled action configurations for this compliance policy. Compliance policy must have one and only one block scheduled action.
// returns a []DeviceComplianceActionItemable when successful
func (m *DeviceComplianceScheduledActionForRule) GetScheduledActionConfigurations()([]DeviceComplianceActionItemable) {
    val, err := m.GetBackingStore().Get("scheduledActionConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceComplianceActionItemable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceComplianceScheduledActionForRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("ruleName", m.GetRuleName())
        if err != nil {
            return err
        }
    }
    if m.GetScheduledActionConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScheduledActionConfigurations()))
        for i, v := range m.GetScheduledActionConfigurations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scheduledActionConfigurations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRuleName sets the ruleName property value. Name of the rule which this scheduled action applies to. Currently scheduled actions are created per policy instead of per rule, thus RuleName is always set to default value PasswordRequired.
func (m *DeviceComplianceScheduledActionForRule) SetRuleName(value *string)() {
    err := m.GetBackingStore().Set("ruleName", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduledActionConfigurations sets the scheduledActionConfigurations property value. The list of scheduled action configurations for this compliance policy. Compliance policy must have one and only one block scheduled action.
func (m *DeviceComplianceScheduledActionForRule) SetScheduledActionConfigurations(value []DeviceComplianceActionItemable)() {
    err := m.GetBackingStore().Set("scheduledActionConfigurations", value)
    if err != nil {
        panic(err)
    }
}
type DeviceComplianceScheduledActionForRuleable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRuleName()(*string)
    GetScheduledActionConfigurations()([]DeviceComplianceActionItemable)
    SetRuleName(value *string)()
    SetScheduledActionConfigurations(value []DeviceComplianceActionItemable)()
}

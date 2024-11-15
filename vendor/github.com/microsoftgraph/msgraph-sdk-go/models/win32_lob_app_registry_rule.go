package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppRegistryRule a complex type to store registry rule data for a Win32 LOB app.
type Win32LobAppRegistryRule struct {
    Win32LobAppRule
}
// NewWin32LobAppRegistryRule instantiates a new Win32LobAppRegistryRule and sets the default values.
func NewWin32LobAppRegistryRule()(*Win32LobAppRegistryRule) {
    m := &Win32LobAppRegistryRule{
        Win32LobAppRule: *NewWin32LobAppRule(),
    }
    odataTypeValue := "#microsoft.graph.win32LobAppRegistryRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWin32LobAppRegistryRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppRegistryRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppRegistryRule(), nil
}
// GetCheck32BitOn64System gets the check32BitOn64System property value. A value indicating whether to search the 32-bit registry on 64-bit systems.
// returns a *bool when successful
func (m *Win32LobAppRegistryRule) GetCheck32BitOn64System()(*bool) {
    val, err := m.GetBackingStore().Get("check32BitOn64System")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetComparisonValue gets the comparisonValue property value. The registry comparison value.
// returns a *string when successful
func (m *Win32LobAppRegistryRule) GetComparisonValue()(*string) {
    val, err := m.GetBackingStore().Get("comparisonValue")
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
func (m *Win32LobAppRegistryRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Win32LobAppRule.GetFieldDeserializers()
    res["check32BitOn64System"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCheck32BitOn64System(val)
        }
        return nil
    }
    res["comparisonValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComparisonValue(val)
        }
        return nil
    }
    res["keyPath"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyPath(val)
        }
        return nil
    }
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppRegistryRuleOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val.(*Win32LobAppRegistryRuleOperationType))
        }
        return nil
    }
    res["operator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppRuleOperator)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperator(val.(*Win32LobAppRuleOperator))
        }
        return nil
    }
    res["valueName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValueName(val)
        }
        return nil
    }
    return res
}
// GetKeyPath gets the keyPath property value. The full path of the registry entry containing the value to detect.
// returns a *string when successful
func (m *Win32LobAppRegistryRule) GetKeyPath()(*string) {
    val, err := m.GetBackingStore().Get("keyPath")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperationType gets the operationType property value. A list of possible operations for rules used to make determinations about an application based on registry keys or values. Unless noted, the values can be used with either detection or requirement rules.
// returns a *Win32LobAppRegistryRuleOperationType when successful
func (m *Win32LobAppRegistryRule) GetOperationType()(*Win32LobAppRegistryRuleOperationType) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppRegistryRuleOperationType)
    }
    return nil
}
// GetOperator gets the operator property value. Contains properties for detection operator.
// returns a *Win32LobAppRuleOperator when successful
func (m *Win32LobAppRegistryRule) GetOperator()(*Win32LobAppRuleOperator) {
    val, err := m.GetBackingStore().Get("operator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppRuleOperator)
    }
    return nil
}
// GetValueName gets the valueName property value. The name of the registry value to detect.
// returns a *string when successful
func (m *Win32LobAppRegistryRule) GetValueName()(*string) {
    val, err := m.GetBackingStore().Get("valueName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobAppRegistryRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Win32LobAppRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("check32BitOn64System", m.GetCheck32BitOn64System())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("comparisonValue", m.GetComparisonValue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("keyPath", m.GetKeyPath())
        if err != nil {
            return err
        }
    }
    if m.GetOperationType() != nil {
        cast := (*m.GetOperationType()).String()
        err = writer.WriteStringValue("operationType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperator() != nil {
        cast := (*m.GetOperator()).String()
        err = writer.WriteStringValue("operator", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("valueName", m.GetValueName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCheck32BitOn64System sets the check32BitOn64System property value. A value indicating whether to search the 32-bit registry on 64-bit systems.
func (m *Win32LobAppRegistryRule) SetCheck32BitOn64System(value *bool)() {
    err := m.GetBackingStore().Set("check32BitOn64System", value)
    if err != nil {
        panic(err)
    }
}
// SetComparisonValue sets the comparisonValue property value. The registry comparison value.
func (m *Win32LobAppRegistryRule) SetComparisonValue(value *string)() {
    err := m.GetBackingStore().Set("comparisonValue", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyPath sets the keyPath property value. The full path of the registry entry containing the value to detect.
func (m *Win32LobAppRegistryRule) SetKeyPath(value *string)() {
    err := m.GetBackingStore().Set("keyPath", value)
    if err != nil {
        panic(err)
    }
}
// SetOperationType sets the operationType property value. A list of possible operations for rules used to make determinations about an application based on registry keys or values. Unless noted, the values can be used with either detection or requirement rules.
func (m *Win32LobAppRegistryRule) SetOperationType(value *Win32LobAppRegistryRuleOperationType)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperator sets the operator property value. Contains properties for detection operator.
func (m *Win32LobAppRegistryRule) SetOperator(value *Win32LobAppRuleOperator)() {
    err := m.GetBackingStore().Set("operator", value)
    if err != nil {
        panic(err)
    }
}
// SetValueName sets the valueName property value. The name of the registry value to detect.
func (m *Win32LobAppRegistryRule) SetValueName(value *string)() {
    err := m.GetBackingStore().Set("valueName", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppRegistryRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Win32LobAppRuleable
    GetCheck32BitOn64System()(*bool)
    GetComparisonValue()(*string)
    GetKeyPath()(*string)
    GetOperationType()(*Win32LobAppRegistryRuleOperationType)
    GetOperator()(*Win32LobAppRuleOperator)
    GetValueName()(*string)
    SetCheck32BitOn64System(value *bool)()
    SetComparisonValue(value *string)()
    SetKeyPath(value *string)()
    SetOperationType(value *Win32LobAppRegistryRuleOperationType)()
    SetOperator(value *Win32LobAppRuleOperator)()
    SetValueName(value *string)()
}

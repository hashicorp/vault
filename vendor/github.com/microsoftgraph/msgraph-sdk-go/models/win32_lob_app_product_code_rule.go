package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppProductCodeRule a complex type to store the product code and version rule data for a Win32 LOB app. This rule is not supported as a requirement rule.
type Win32LobAppProductCodeRule struct {
    Win32LobAppRule
}
// NewWin32LobAppProductCodeRule instantiates a new Win32LobAppProductCodeRule and sets the default values.
func NewWin32LobAppProductCodeRule()(*Win32LobAppProductCodeRule) {
    m := &Win32LobAppProductCodeRule{
        Win32LobAppRule: *NewWin32LobAppRule(),
    }
    odataTypeValue := "#microsoft.graph.win32LobAppProductCodeRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWin32LobAppProductCodeRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppProductCodeRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppProductCodeRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Win32LobAppProductCodeRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Win32LobAppRule.GetFieldDeserializers()
    res["productCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductCode(val)
        }
        return nil
    }
    res["productVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductVersion(val)
        }
        return nil
    }
    res["productVersionOperator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppRuleOperator)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductVersionOperator(val.(*Win32LobAppRuleOperator))
        }
        return nil
    }
    return res
}
// GetProductCode gets the productCode property value. The product code of the app.
// returns a *string when successful
func (m *Win32LobAppProductCodeRule) GetProductCode()(*string) {
    val, err := m.GetBackingStore().Get("productCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductVersion gets the productVersion property value. The product version comparison value.
// returns a *string when successful
func (m *Win32LobAppProductCodeRule) GetProductVersion()(*string) {
    val, err := m.GetBackingStore().Get("productVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductVersionOperator gets the productVersionOperator property value. Contains properties for detection operator.
// returns a *Win32LobAppRuleOperator when successful
func (m *Win32LobAppProductCodeRule) GetProductVersionOperator()(*Win32LobAppRuleOperator) {
    val, err := m.GetBackingStore().Get("productVersionOperator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppRuleOperator)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobAppProductCodeRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Win32LobAppRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("productCode", m.GetProductCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productVersion", m.GetProductVersion())
        if err != nil {
            return err
        }
    }
    if m.GetProductVersionOperator() != nil {
        cast := (*m.GetProductVersionOperator()).String()
        err = writer.WriteStringValue("productVersionOperator", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetProductCode sets the productCode property value. The product code of the app.
func (m *Win32LobAppProductCodeRule) SetProductCode(value *string)() {
    err := m.GetBackingStore().Set("productCode", value)
    if err != nil {
        panic(err)
    }
}
// SetProductVersion sets the productVersion property value. The product version comparison value.
func (m *Win32LobAppProductCodeRule) SetProductVersion(value *string)() {
    err := m.GetBackingStore().Set("productVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetProductVersionOperator sets the productVersionOperator property value. Contains properties for detection operator.
func (m *Win32LobAppProductCodeRule) SetProductVersionOperator(value *Win32LobAppRuleOperator)() {
    err := m.GetBackingStore().Set("productVersionOperator", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppProductCodeRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Win32LobAppRuleable
    GetProductCode()(*string)
    GetProductVersion()(*string)
    GetProductVersionOperator()(*Win32LobAppRuleOperator)
    SetProductCode(value *string)()
    SetProductVersion(value *string)()
    SetProductVersionOperator(value *Win32LobAppRuleOperator)()
}

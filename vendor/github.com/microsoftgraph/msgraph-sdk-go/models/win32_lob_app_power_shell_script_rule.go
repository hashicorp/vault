package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppPowerShellScriptRule a complex type to store the PowerShell script rule data for a Win32 LOB app.
type Win32LobAppPowerShellScriptRule struct {
    Win32LobAppRule
}
// NewWin32LobAppPowerShellScriptRule instantiates a new Win32LobAppPowerShellScriptRule and sets the default values.
func NewWin32LobAppPowerShellScriptRule()(*Win32LobAppPowerShellScriptRule) {
    m := &Win32LobAppPowerShellScriptRule{
        Win32LobAppRule: *NewWin32LobAppRule(),
    }
    odataTypeValue := "#microsoft.graph.win32LobAppPowerShellScriptRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWin32LobAppPowerShellScriptRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppPowerShellScriptRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppPowerShellScriptRule(), nil
}
// GetComparisonValue gets the comparisonValue property value. The script output comparison value. Do not specify a value if the rule is used for detection.
// returns a *string when successful
func (m *Win32LobAppPowerShellScriptRule) GetComparisonValue()(*string) {
    val, err := m.GetBackingStore().Get("comparisonValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the rule. Do not specify this value if the rule is used for detection.
// returns a *string when successful
func (m *Win32LobAppPowerShellScriptRule) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnforceSignatureCheck gets the enforceSignatureCheck property value. A value indicating whether a signature check is enforced.
// returns a *bool when successful
func (m *Win32LobAppPowerShellScriptRule) GetEnforceSignatureCheck()(*bool) {
    val, err := m.GetBackingStore().Get("enforceSignatureCheck")
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
func (m *Win32LobAppPowerShellScriptRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Win32LobAppRule.GetFieldDeserializers()
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
    res["enforceSignatureCheck"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnforceSignatureCheck(val)
        }
        return nil
    }
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppPowerShellScriptRuleOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val.(*Win32LobAppPowerShellScriptRuleOperationType))
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
    res["runAs32Bit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRunAs32Bit(val)
        }
        return nil
    }
    res["runAsAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRunAsAccountType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRunAsAccount(val.(*RunAsAccountType))
        }
        return nil
    }
    res["scriptContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScriptContent(val)
        }
        return nil
    }
    return res
}
// GetOperationType gets the operationType property value. Contains all supported Powershell Script output detection type.
// returns a *Win32LobAppPowerShellScriptRuleOperationType when successful
func (m *Win32LobAppPowerShellScriptRule) GetOperationType()(*Win32LobAppPowerShellScriptRuleOperationType) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppPowerShellScriptRuleOperationType)
    }
    return nil
}
// GetOperator gets the operator property value. Contains properties for detection operator.
// returns a *Win32LobAppRuleOperator when successful
func (m *Win32LobAppPowerShellScriptRule) GetOperator()(*Win32LobAppRuleOperator) {
    val, err := m.GetBackingStore().Get("operator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppRuleOperator)
    }
    return nil
}
// GetRunAs32Bit gets the runAs32Bit property value. A value indicating whether the script should run as 32-bit.
// returns a *bool when successful
func (m *Win32LobAppPowerShellScriptRule) GetRunAs32Bit()(*bool) {
    val, err := m.GetBackingStore().Get("runAs32Bit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRunAsAccount gets the runAsAccount property value. The execution context of the script. Do not specify this value if the rule is used for detection. Script detection rules will run in the same context as the associated app install context. Possible values are: system, user.
// returns a *RunAsAccountType when successful
func (m *Win32LobAppPowerShellScriptRule) GetRunAsAccount()(*RunAsAccountType) {
    val, err := m.GetBackingStore().Get("runAsAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RunAsAccountType)
    }
    return nil
}
// GetScriptContent gets the scriptContent property value. The base64-encoded script content.
// returns a *string when successful
func (m *Win32LobAppPowerShellScriptRule) GetScriptContent()(*string) {
    val, err := m.GetBackingStore().Get("scriptContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobAppPowerShellScriptRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Win32LobAppRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("comparisonValue", m.GetComparisonValue())
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
        err = writer.WriteBoolValue("enforceSignatureCheck", m.GetEnforceSignatureCheck())
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
        err = writer.WriteBoolValue("runAs32Bit", m.GetRunAs32Bit())
        if err != nil {
            return err
        }
    }
    if m.GetRunAsAccount() != nil {
        cast := (*m.GetRunAsAccount()).String()
        err = writer.WriteStringValue("runAsAccount", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("scriptContent", m.GetScriptContent())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetComparisonValue sets the comparisonValue property value. The script output comparison value. Do not specify a value if the rule is used for detection.
func (m *Win32LobAppPowerShellScriptRule) SetComparisonValue(value *string)() {
    err := m.GetBackingStore().Set("comparisonValue", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the rule. Do not specify this value if the rule is used for detection.
func (m *Win32LobAppPowerShellScriptRule) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEnforceSignatureCheck sets the enforceSignatureCheck property value. A value indicating whether a signature check is enforced.
func (m *Win32LobAppPowerShellScriptRule) SetEnforceSignatureCheck(value *bool)() {
    err := m.GetBackingStore().Set("enforceSignatureCheck", value)
    if err != nil {
        panic(err)
    }
}
// SetOperationType sets the operationType property value. Contains all supported Powershell Script output detection type.
func (m *Win32LobAppPowerShellScriptRule) SetOperationType(value *Win32LobAppPowerShellScriptRuleOperationType)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperator sets the operator property value. Contains properties for detection operator.
func (m *Win32LobAppPowerShellScriptRule) SetOperator(value *Win32LobAppRuleOperator)() {
    err := m.GetBackingStore().Set("operator", value)
    if err != nil {
        panic(err)
    }
}
// SetRunAs32Bit sets the runAs32Bit property value. A value indicating whether the script should run as 32-bit.
func (m *Win32LobAppPowerShellScriptRule) SetRunAs32Bit(value *bool)() {
    err := m.GetBackingStore().Set("runAs32Bit", value)
    if err != nil {
        panic(err)
    }
}
// SetRunAsAccount sets the runAsAccount property value. The execution context of the script. Do not specify this value if the rule is used for detection. Script detection rules will run in the same context as the associated app install context. Possible values are: system, user.
func (m *Win32LobAppPowerShellScriptRule) SetRunAsAccount(value *RunAsAccountType)() {
    err := m.GetBackingStore().Set("runAsAccount", value)
    if err != nil {
        panic(err)
    }
}
// SetScriptContent sets the scriptContent property value. The base64-encoded script content.
func (m *Win32LobAppPowerShellScriptRule) SetScriptContent(value *string)() {
    err := m.GetBackingStore().Set("scriptContent", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppPowerShellScriptRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Win32LobAppRuleable
    GetComparisonValue()(*string)
    GetDisplayName()(*string)
    GetEnforceSignatureCheck()(*bool)
    GetOperationType()(*Win32LobAppPowerShellScriptRuleOperationType)
    GetOperator()(*Win32LobAppRuleOperator)
    GetRunAs32Bit()(*bool)
    GetRunAsAccount()(*RunAsAccountType)
    GetScriptContent()(*string)
    SetComparisonValue(value *string)()
    SetDisplayName(value *string)()
    SetEnforceSignatureCheck(value *bool)()
    SetOperationType(value *Win32LobAppPowerShellScriptRuleOperationType)()
    SetOperator(value *Win32LobAppRuleOperator)()
    SetRunAs32Bit(value *bool)()
    SetRunAsAccount(value *RunAsAccountType)()
    SetScriptContent(value *string)()
}

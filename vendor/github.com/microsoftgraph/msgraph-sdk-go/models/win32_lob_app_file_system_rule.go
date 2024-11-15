package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppFileSystemRule a complex type to store file or folder rule data for a Win32 LOB app.
type Win32LobAppFileSystemRule struct {
    Win32LobAppRule
}
// NewWin32LobAppFileSystemRule instantiates a new Win32LobAppFileSystemRule and sets the default values.
func NewWin32LobAppFileSystemRule()(*Win32LobAppFileSystemRule) {
    m := &Win32LobAppFileSystemRule{
        Win32LobAppRule: *NewWin32LobAppRule(),
    }
    odataTypeValue := "#microsoft.graph.win32LobAppFileSystemRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWin32LobAppFileSystemRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppFileSystemRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppFileSystemRule(), nil
}
// GetCheck32BitOn64System gets the check32BitOn64System property value. A value indicating whether to expand environment variables in the 32-bit context on 64-bit systems.
// returns a *bool when successful
func (m *Win32LobAppFileSystemRule) GetCheck32BitOn64System()(*bool) {
    val, err := m.GetBackingStore().Get("check32BitOn64System")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetComparisonValue gets the comparisonValue property value. The file or folder comparison value.
// returns a *string when successful
func (m *Win32LobAppFileSystemRule) GetComparisonValue()(*string) {
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
func (m *Win32LobAppFileSystemRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["fileOrFolderName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileOrFolderName(val)
        }
        return nil
    }
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppFileSystemOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val.(*Win32LobAppFileSystemOperationType))
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
    res["path"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPath(val)
        }
        return nil
    }
    return res
}
// GetFileOrFolderName gets the fileOrFolderName property value. The file or folder name to look up.
// returns a *string when successful
func (m *Win32LobAppFileSystemRule) GetFileOrFolderName()(*string) {
    val, err := m.GetBackingStore().Get("fileOrFolderName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperationType gets the operationType property value. A list of possible operations for rules used to make determinations about an application based on files or folders. Unless noted, can be used with either detection or requirement rules.
// returns a *Win32LobAppFileSystemOperationType when successful
func (m *Win32LobAppFileSystemRule) GetOperationType()(*Win32LobAppFileSystemOperationType) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppFileSystemOperationType)
    }
    return nil
}
// GetOperator gets the operator property value. Contains properties for detection operator.
// returns a *Win32LobAppRuleOperator when successful
func (m *Win32LobAppFileSystemRule) GetOperator()(*Win32LobAppRuleOperator) {
    val, err := m.GetBackingStore().Get("operator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppRuleOperator)
    }
    return nil
}
// GetPath gets the path property value. The file or folder path to look up.
// returns a *string when successful
func (m *Win32LobAppFileSystemRule) GetPath()(*string) {
    val, err := m.GetBackingStore().Get("path")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobAppFileSystemRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("fileOrFolderName", m.GetFileOrFolderName())
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
        err = writer.WriteStringValue("path", m.GetPath())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCheck32BitOn64System sets the check32BitOn64System property value. A value indicating whether to expand environment variables in the 32-bit context on 64-bit systems.
func (m *Win32LobAppFileSystemRule) SetCheck32BitOn64System(value *bool)() {
    err := m.GetBackingStore().Set("check32BitOn64System", value)
    if err != nil {
        panic(err)
    }
}
// SetComparisonValue sets the comparisonValue property value. The file or folder comparison value.
func (m *Win32LobAppFileSystemRule) SetComparisonValue(value *string)() {
    err := m.GetBackingStore().Set("comparisonValue", value)
    if err != nil {
        panic(err)
    }
}
// SetFileOrFolderName sets the fileOrFolderName property value. The file or folder name to look up.
func (m *Win32LobAppFileSystemRule) SetFileOrFolderName(value *string)() {
    err := m.GetBackingStore().Set("fileOrFolderName", value)
    if err != nil {
        panic(err)
    }
}
// SetOperationType sets the operationType property value. A list of possible operations for rules used to make determinations about an application based on files or folders. Unless noted, can be used with either detection or requirement rules.
func (m *Win32LobAppFileSystemRule) SetOperationType(value *Win32LobAppFileSystemOperationType)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperator sets the operator property value. Contains properties for detection operator.
func (m *Win32LobAppFileSystemRule) SetOperator(value *Win32LobAppRuleOperator)() {
    err := m.GetBackingStore().Set("operator", value)
    if err != nil {
        panic(err)
    }
}
// SetPath sets the path property value. The file or folder path to look up.
func (m *Win32LobAppFileSystemRule) SetPath(value *string)() {
    err := m.GetBackingStore().Set("path", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppFileSystemRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Win32LobAppRuleable
    GetCheck32BitOn64System()(*bool)
    GetComparisonValue()(*string)
    GetFileOrFolderName()(*string)
    GetOperationType()(*Win32LobAppFileSystemOperationType)
    GetOperator()(*Win32LobAppRuleOperator)
    GetPath()(*string)
    SetCheck32BitOn64System(value *bool)()
    SetComparisonValue(value *string)()
    SetFileOrFolderName(value *string)()
    SetOperationType(value *Win32LobAppFileSystemOperationType)()
    SetOperator(value *Win32LobAppRuleOperator)()
    SetPath(value *string)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DriveProtectionRule struct {
    ProtectionRuleBase
}
// NewDriveProtectionRule instantiates a new DriveProtectionRule and sets the default values.
func NewDriveProtectionRule()(*DriveProtectionRule) {
    m := &DriveProtectionRule{
        ProtectionRuleBase: *NewProtectionRuleBase(),
    }
    odataTypeValue := "#microsoft.graph.driveProtectionRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDriveProtectionRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDriveProtectionRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDriveProtectionRule(), nil
}
// GetDriveExpression gets the driveExpression property value. Contains a drive expression. For examples, see driveExpression examples.
// returns a *string when successful
func (m *DriveProtectionRule) GetDriveExpression()(*string) {
    val, err := m.GetBackingStore().Get("driveExpression")
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
func (m *DriveProtectionRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionRuleBase.GetFieldDeserializers()
    res["driveExpression"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriveExpression(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *DriveProtectionRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionRuleBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("driveExpression", m.GetDriveExpression())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDriveExpression sets the driveExpression property value. Contains a drive expression. For examples, see driveExpression examples.
func (m *DriveProtectionRule) SetDriveExpression(value *string)() {
    err := m.GetBackingStore().Set("driveExpression", value)
    if err != nil {
        panic(err)
    }
}
type DriveProtectionRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionRuleBaseable
    GetDriveExpression()(*string)
    SetDriveExpression(value *string)()
}

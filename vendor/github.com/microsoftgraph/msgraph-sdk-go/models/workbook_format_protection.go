package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookFormatProtection struct {
    Entity
}
// NewWorkbookFormatProtection instantiates a new WorkbookFormatProtection and sets the default values.
func NewWorkbookFormatProtection()(*WorkbookFormatProtection) {
    m := &WorkbookFormatProtection{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFormatProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookFormatProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookFormatProtection(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookFormatProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["formulaHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormulaHidden(val)
        }
        return nil
    }
    res["locked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocked(val)
        }
        return nil
    }
    return res
}
// GetFormulaHidden gets the formulaHidden property value. Indicates whether Excel hides the formula for the cells in the range. A null value indicates that the entire range doesn't have uniform formula hidden setting.
// returns a *bool when successful
func (m *WorkbookFormatProtection) GetFormulaHidden()(*bool) {
    val, err := m.GetBackingStore().Get("formulaHidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocked gets the locked property value. Indicates whether Excel locks the cells in the object. A null value indicates that the entire range doesn't have uniform lock setting.
// returns a *bool when successful
func (m *WorkbookFormatProtection) GetLocked()(*bool) {
    val, err := m.GetBackingStore().Get("locked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookFormatProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("formulaHidden", m.GetFormulaHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("locked", m.GetLocked())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormulaHidden sets the formulaHidden property value. Indicates whether Excel hides the formula for the cells in the range. A null value indicates that the entire range doesn't have uniform formula hidden setting.
func (m *WorkbookFormatProtection) SetFormulaHidden(value *bool)() {
    err := m.GetBackingStore().Set("formulaHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetLocked sets the locked property value. Indicates whether Excel locks the cells in the object. A null value indicates that the entire range doesn't have uniform lock setting.
func (m *WorkbookFormatProtection) SetLocked(value *bool)() {
    err := m.GetBackingStore().Set("locked", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookFormatProtectionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFormulaHidden()(*bool)
    GetLocked()(*bool)
    SetFormulaHidden(value *bool)()
    SetLocked(value *bool)()
}

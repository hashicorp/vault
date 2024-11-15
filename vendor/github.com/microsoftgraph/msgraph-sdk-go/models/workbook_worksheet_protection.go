package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookWorksheetProtection struct {
    Entity
}
// NewWorkbookWorksheetProtection instantiates a new WorkbookWorksheetProtection and sets the default values.
func NewWorkbookWorksheetProtection()(*WorkbookWorksheetProtection) {
    m := &WorkbookWorksheetProtection{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookWorksheetProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookWorksheetProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookWorksheetProtection(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookWorksheetProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["options"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookWorksheetProtectionOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOptions(val.(WorkbookWorksheetProtectionOptionsable))
        }
        return nil
    }
    res["protected"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtected(val)
        }
        return nil
    }
    return res
}
// GetOptions gets the options property value. Worksheet protection options. Read-only.
// returns a WorkbookWorksheetProtectionOptionsable when successful
func (m *WorkbookWorksheetProtection) GetOptions()(WorkbookWorksheetProtectionOptionsable) {
    val, err := m.GetBackingStore().Get("options")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookWorksheetProtectionOptionsable)
    }
    return nil
}
// GetProtected gets the protected property value. Indicates whether the worksheet is protected.  Read-only.
// returns a *bool when successful
func (m *WorkbookWorksheetProtection) GetProtected()(*bool) {
    val, err := m.GetBackingStore().Get("protected")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookWorksheetProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("options", m.GetOptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("protected", m.GetProtected())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOptions sets the options property value. Worksheet protection options. Read-only.
func (m *WorkbookWorksheetProtection) SetOptions(value WorkbookWorksheetProtectionOptionsable)() {
    err := m.GetBackingStore().Set("options", value)
    if err != nil {
        panic(err)
    }
}
// SetProtected sets the protected property value. Indicates whether the worksheet is protected.  Read-only.
func (m *WorkbookWorksheetProtection) SetProtected(value *bool)() {
    err := m.GetBackingStore().Set("protected", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookWorksheetProtectionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOptions()(WorkbookWorksheetProtectionOptionsable)
    GetProtected()(*bool)
    SetOptions(value WorkbookWorksheetProtectionOptionsable)()
    SetProtected(value *bool)()
}

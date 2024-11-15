package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookOperation struct {
    Entity
}
// NewWorkbookOperation instantiates a new WorkbookOperation and sets the default values.
func NewWorkbookOperation()(*WorkbookOperation) {
    m := &WorkbookOperation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookOperation(), nil
}
// GetError gets the error property value. The error returned by the operation.
// returns a WorkbookOperationErrorable when successful
func (m *WorkbookOperation) GetError()(WorkbookOperationErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookOperationErrorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookOperationErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(WorkbookOperationErrorable))
        }
        return nil
    }
    res["resourceLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceLocation(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWorkbookOperationStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*WorkbookOperationStatus))
        }
        return nil
    }
    return res
}
// GetResourceLocation gets the resourceLocation property value. The resource URI for the result.
// returns a *string when successful
func (m *WorkbookOperation) GetResourceLocation()(*string) {
    val, err := m.GetBackingStore().Get("resourceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *WorkbookOperationStatus when successful
func (m *WorkbookOperation) GetStatus()(*WorkbookOperationStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WorkbookOperationStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("error", m.GetError())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceLocation", m.GetResourceLocation())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetError sets the error property value. The error returned by the operation.
func (m *WorkbookOperation) SetError(value WorkbookOperationErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceLocation sets the resourceLocation property value. The resource URI for the result.
func (m *WorkbookOperation) SetResourceLocation(value *string)() {
    err := m.GetBackingStore().Set("resourceLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *WorkbookOperation) SetStatus(value *WorkbookOperationStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookOperationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetError()(WorkbookOperationErrorable)
    GetResourceLocation()(*string)
    GetStatus()(*WorkbookOperationStatus)
    SetError(value WorkbookOperationErrorable)()
    SetResourceLocation(value *string)()
    SetStatus(value *WorkbookOperationStatus)()
}

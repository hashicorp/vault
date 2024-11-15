package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnenoteOperation struct {
    Operation
}
// NewOnenoteOperation instantiates a new OnenoteOperation and sets the default values.
func NewOnenoteOperation()(*OnenoteOperation) {
    m := &OnenoteOperation{
        Operation: *NewOperation(),
    }
    return m
}
// CreateOnenoteOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnenoteOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnenoteOperation(), nil
}
// GetError gets the error property value. The error returned by the operation.
// returns a OnenoteOperationErrorable when successful
func (m *OnenoteOperation) GetError()(OnenoteOperationErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnenoteOperationErrorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnenoteOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Operation.GetFieldDeserializers()
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnenoteOperationErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(OnenoteOperationErrorable))
        }
        return nil
    }
    res["percentComplete"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPercentComplete(val)
        }
        return nil
    }
    res["resourceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceId(val)
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
    return res
}
// GetPercentComplete gets the percentComplete property value. The operation percent complete if the operation is still in running status.
// returns a *string when successful
func (m *OnenoteOperation) GetPercentComplete()(*string) {
    val, err := m.GetBackingStore().Get("percentComplete")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceId gets the resourceId property value. The resource id.
// returns a *string when successful
func (m *OnenoteOperation) GetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("resourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceLocation gets the resourceLocation property value. The resource URI for the object. For example, the resource URI for a copied page or section.
// returns a *string when successful
func (m *OnenoteOperation) GetResourceLocation()(*string) {
    val, err := m.GetBackingStore().Get("resourceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnenoteOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Operation.Serialize(writer)
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
        err = writer.WriteStringValue("percentComplete", m.GetPercentComplete())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceId", m.GetResourceId())
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
    return nil
}
// SetError sets the error property value. The error returned by the operation.
func (m *OnenoteOperation) SetError(value OnenoteOperationErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetPercentComplete sets the percentComplete property value. The operation percent complete if the operation is still in running status.
func (m *OnenoteOperation) SetPercentComplete(value *string)() {
    err := m.GetBackingStore().Set("percentComplete", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceId sets the resourceId property value. The resource id.
func (m *OnenoteOperation) SetResourceId(value *string)() {
    err := m.GetBackingStore().Set("resourceId", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceLocation sets the resourceLocation property value. The resource URI for the object. For example, the resource URI for a copied page or section.
func (m *OnenoteOperation) SetResourceLocation(value *string)() {
    err := m.GetBackingStore().Set("resourceLocation", value)
    if err != nil {
        panic(err)
    }
}
type OnenoteOperationable interface {
    Operationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetError()(OnenoteOperationErrorable)
    GetPercentComplete()(*string)
    GetResourceId()(*string)
    GetResourceLocation()(*string)
    SetError(value OnenoteOperationErrorable)()
    SetPercentComplete(value *string)()
    SetResourceId(value *string)()
    SetResourceLocation(value *string)()
}

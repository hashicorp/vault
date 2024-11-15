package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RichLongRunningOperation the status of a long-running operation.
type RichLongRunningOperation struct {
    LongRunningOperation
}
// NewRichLongRunningOperation instantiates a new RichLongRunningOperation and sets the default values.
func NewRichLongRunningOperation()(*RichLongRunningOperation) {
    m := &RichLongRunningOperation{
        LongRunningOperation: *NewLongRunningOperation(),
    }
    return m
}
// CreateRichLongRunningOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRichLongRunningOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRichLongRunningOperation(), nil
}
// GetError gets the error property value. Error that caused the operation to fail.
// returns a PublicErrorable when successful
func (m *RichLongRunningOperation) GetError()(PublicErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicErrorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RichLongRunningOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.LongRunningOperation.GetFieldDeserializers()
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(PublicErrorable))
        }
        return nil
    }
    res["percentageComplete"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPercentageComplete(val)
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
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val)
        }
        return nil
    }
    return res
}
// GetPercentageComplete gets the percentageComplete property value. A value between 0 and 100 that indicates the progress of the operation.
// returns a *int32 when successful
func (m *RichLongRunningOperation) GetPercentageComplete()(*int32) {
    val, err := m.GetBackingStore().Get("percentageComplete")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetResourceId gets the resourceId property value. The unique identifier for the result.
// returns a *string when successful
func (m *RichLongRunningOperation) GetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("resourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The type of the operation.
// returns a *string when successful
func (m *RichLongRunningOperation) GetTypeEscaped()(*string) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RichLongRunningOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.LongRunningOperation.Serialize(writer)
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
        err = writer.WriteInt32Value("percentageComplete", m.GetPercentageComplete())
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
        err = writer.WriteStringValue("type", m.GetTypeEscaped())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetError sets the error property value. Error that caused the operation to fail.
func (m *RichLongRunningOperation) SetError(value PublicErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetPercentageComplete sets the percentageComplete property value. A value between 0 and 100 that indicates the progress of the operation.
func (m *RichLongRunningOperation) SetPercentageComplete(value *int32)() {
    err := m.GetBackingStore().Set("percentageComplete", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceId sets the resourceId property value. The unique identifier for the result.
func (m *RichLongRunningOperation) SetResourceId(value *string)() {
    err := m.GetBackingStore().Set("resourceId", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The type of the operation.
func (m *RichLongRunningOperation) SetTypeEscaped(value *string)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type RichLongRunningOperationable interface {
    LongRunningOperationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetError()(PublicErrorable)
    GetPercentageComplete()(*int32)
    GetResourceId()(*string)
    GetTypeEscaped()(*string)
    SetError(value PublicErrorable)()
    SetPercentageComplete(value *int32)()
    SetResourceId(value *string)()
    SetTypeEscaped(value *string)()
}

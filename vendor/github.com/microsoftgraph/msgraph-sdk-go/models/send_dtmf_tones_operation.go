package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SendDtmfTonesOperation struct {
    CommsOperation
}
// NewSendDtmfTonesOperation instantiates a new SendDtmfTonesOperation and sets the default values.
func NewSendDtmfTonesOperation()(*SendDtmfTonesOperation) {
    m := &SendDtmfTonesOperation{
        CommsOperation: *NewCommsOperation(),
    }
    return m
}
// CreateSendDtmfTonesOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSendDtmfTonesOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSendDtmfTonesOperation(), nil
}
// GetCompletionReason gets the completionReason property value. The results of the action. Possible values are: unknown, completedSuccessfully, mediaOperationCanceled, unknownfutureValue.
// returns a *SendDtmfCompletionReason when successful
func (m *SendDtmfTonesOperation) GetCompletionReason()(*SendDtmfCompletionReason) {
    val, err := m.GetBackingStore().Get("completionReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SendDtmfCompletionReason)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SendDtmfTonesOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CommsOperation.GetFieldDeserializers()
    res["completionReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSendDtmfCompletionReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletionReason(val.(*SendDtmfCompletionReason))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *SendDtmfTonesOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CommsOperation.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCompletionReason() != nil {
        cast := (*m.GetCompletionReason()).String()
        err = writer.WriteStringValue("completionReason", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletionReason sets the completionReason property value. The results of the action. Possible values are: unknown, completedSuccessfully, mediaOperationCanceled, unknownfutureValue.
func (m *SendDtmfTonesOperation) SetCompletionReason(value *SendDtmfCompletionReason)() {
    err := m.GetBackingStore().Set("completionReason", value)
    if err != nil {
        panic(err)
    }
}
type SendDtmfTonesOperationable interface {
    CommsOperationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompletionReason()(*SendDtmfCompletionReason)
    SetCompletionReason(value *SendDtmfCompletionReason)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EngagementAsyncOperation represents the status of a Viva Engage async operation that is an operation that transcends thelifetime of a single API request. These operations are long-running or too expensive to completewithin the time frame of their original request.
type EngagementAsyncOperation struct {
    LongRunningOperation
}
// NewEngagementAsyncOperation instantiates a new EngagementAsyncOperation and sets the default values.
func NewEngagementAsyncOperation()(*EngagementAsyncOperation) {
    m := &EngagementAsyncOperation{
        LongRunningOperation: *NewLongRunningOperation(),
    }
    return m
}
// CreateEngagementAsyncOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEngagementAsyncOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEngagementAsyncOperation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EngagementAsyncOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.LongRunningOperation.GetFieldDeserializers()
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEngagementAsyncOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val.(*EngagementAsyncOperationType))
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
    return res
}
// GetOperationType gets the operationType property value. The type of the long-running operation. The possible values are: createCommunity, unknownFutureValue.
// returns a *EngagementAsyncOperationType when successful
func (m *EngagementAsyncOperation) GetOperationType()(*EngagementAsyncOperationType) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EngagementAsyncOperationType)
    }
    return nil
}
// GetResourceId gets the resourceId property value. The ID of the object created or modified as a result of this async operation.
// returns a *string when successful
func (m *EngagementAsyncOperation) GetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("resourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EngagementAsyncOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.LongRunningOperation.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetOperationType() != nil {
        cast := (*m.GetOperationType()).String()
        err = writer.WriteStringValue("operationType", &cast)
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
    return nil
}
// SetOperationType sets the operationType property value. The type of the long-running operation. The possible values are: createCommunity, unknownFutureValue.
func (m *EngagementAsyncOperation) SetOperationType(value *EngagementAsyncOperationType)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceId sets the resourceId property value. The ID of the object created or modified as a result of this async operation.
func (m *EngagementAsyncOperation) SetResourceId(value *string)() {
    err := m.GetBackingStore().Set("resourceId", value)
    if err != nil {
        panic(err)
    }
}
type EngagementAsyncOperationable interface {
    LongRunningOperationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOperationType()(*EngagementAsyncOperationType)
    GetResourceId()(*string)
    SetOperationType(value *EngagementAsyncOperationType)()
    SetResourceId(value *string)()
}

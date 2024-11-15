package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DisableAndDeleteUserApplyAction struct {
    AccessReviewApplyAction
}
// NewDisableAndDeleteUserApplyAction instantiates a new DisableAndDeleteUserApplyAction and sets the default values.
func NewDisableAndDeleteUserApplyAction()(*DisableAndDeleteUserApplyAction) {
    m := &DisableAndDeleteUserApplyAction{
        AccessReviewApplyAction: *NewAccessReviewApplyAction(),
    }
    odataTypeValue := "#microsoft.graph.disableAndDeleteUserApplyAction"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDisableAndDeleteUserApplyActionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDisableAndDeleteUserApplyActionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDisableAndDeleteUserApplyAction(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DisableAndDeleteUserApplyAction) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewApplyAction.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *DisableAndDeleteUserApplyAction) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewApplyAction.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type DisableAndDeleteUserApplyActionable interface {
    AccessReviewApplyActionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

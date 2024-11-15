package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RemoveAccessApplyAction struct {
    AccessReviewApplyAction
}
// NewRemoveAccessApplyAction instantiates a new RemoveAccessApplyAction and sets the default values.
func NewRemoveAccessApplyAction()(*RemoveAccessApplyAction) {
    m := &RemoveAccessApplyAction{
        AccessReviewApplyAction: *NewAccessReviewApplyAction(),
    }
    odataTypeValue := "#microsoft.graph.removeAccessApplyAction"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRemoveAccessApplyActionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRemoveAccessApplyActionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRemoveAccessApplyAction(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RemoveAccessApplyAction) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewApplyAction.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *RemoveAccessApplyAction) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewApplyAction.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type RemoveAccessApplyActionable interface {
    AccessReviewApplyActionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

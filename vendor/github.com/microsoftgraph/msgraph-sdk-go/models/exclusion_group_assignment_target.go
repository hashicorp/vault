package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ExclusionGroupAssignmentTarget represents a group that should be excluded from an assignment.
type ExclusionGroupAssignmentTarget struct {
    GroupAssignmentTarget
}
// NewExclusionGroupAssignmentTarget instantiates a new ExclusionGroupAssignmentTarget and sets the default values.
func NewExclusionGroupAssignmentTarget()(*ExclusionGroupAssignmentTarget) {
    m := &ExclusionGroupAssignmentTarget{
        GroupAssignmentTarget: *NewGroupAssignmentTarget(),
    }
    odataTypeValue := "#microsoft.graph.exclusionGroupAssignmentTarget"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateExclusionGroupAssignmentTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateExclusionGroupAssignmentTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExclusionGroupAssignmentTarget(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ExclusionGroupAssignmentTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.GroupAssignmentTarget.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *ExclusionGroupAssignmentTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.GroupAssignmentTarget.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type ExclusionGroupAssignmentTargetable interface {
    GroupAssignmentTargetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

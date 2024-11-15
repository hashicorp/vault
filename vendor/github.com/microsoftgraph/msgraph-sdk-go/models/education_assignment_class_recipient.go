package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationAssignmentClassRecipient struct {
    EducationAssignmentRecipient
}
// NewEducationAssignmentClassRecipient instantiates a new EducationAssignmentClassRecipient and sets the default values.
func NewEducationAssignmentClassRecipient()(*EducationAssignmentClassRecipient) {
    m := &EducationAssignmentClassRecipient{
        EducationAssignmentRecipient: *NewEducationAssignmentRecipient(),
    }
    odataTypeValue := "#microsoft.graph.educationAssignmentClassRecipient"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationAssignmentClassRecipientFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationAssignmentClassRecipientFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignmentClassRecipient(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationAssignmentClassRecipient) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationAssignmentRecipient.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *EducationAssignmentClassRecipient) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationAssignmentRecipient.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type EducationAssignmentClassRecipientable interface {
    EducationAssignmentRecipientable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

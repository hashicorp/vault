package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosVppEBookAssignment contains properties used to assign an iOS VPP EBook to a group.
type IosVppEBookAssignment struct {
    ManagedEBookAssignment
}
// NewIosVppEBookAssignment instantiates a new IosVppEBookAssignment and sets the default values.
func NewIosVppEBookAssignment()(*IosVppEBookAssignment) {
    m := &IosVppEBookAssignment{
        ManagedEBookAssignment: *NewManagedEBookAssignment(),
    }
    return m
}
// CreateIosVppEBookAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosVppEBookAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosVppEBookAssignment(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IosVppEBookAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedEBookAssignment.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *IosVppEBookAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedEBookAssignment.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type IosVppEBookAssignmentable interface {
    ManagedEBookAssignmentable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

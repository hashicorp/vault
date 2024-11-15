package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SubjectRightsRequestAllMailboxLocation struct {
    SubjectRightsRequestMailboxLocation
}
// NewSubjectRightsRequestAllMailboxLocation instantiates a new SubjectRightsRequestAllMailboxLocation and sets the default values.
func NewSubjectRightsRequestAllMailboxLocation()(*SubjectRightsRequestAllMailboxLocation) {
    m := &SubjectRightsRequestAllMailboxLocation{
        SubjectRightsRequestMailboxLocation: *NewSubjectRightsRequestMailboxLocation(),
    }
    odataTypeValue := "#microsoft.graph.subjectRightsRequestAllMailboxLocation"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSubjectRightsRequestAllMailboxLocationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSubjectRightsRequestAllMailboxLocationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSubjectRightsRequestAllMailboxLocation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SubjectRightsRequestAllMailboxLocation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectRightsRequestMailboxLocation.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *SubjectRightsRequestAllMailboxLocation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectRightsRequestMailboxLocation.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type SubjectRightsRequestAllMailboxLocationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SubjectRightsRequestMailboxLocationable
}

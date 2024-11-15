package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Organizer struct {
    ParticipantBase
}
// NewOrganizer instantiates a new Organizer and sets the default values.
func NewOrganizer()(*Organizer) {
    m := &Organizer{
        ParticipantBase: *NewParticipantBase(),
    }
    odataTypeValue := "#microsoft.graph.callRecords.organizer"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOrganizerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOrganizerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOrganizer(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Organizer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ParticipantBase.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *Organizer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ParticipantBase.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type Organizerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ParticipantBaseable
}

package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Participant struct {
    ParticipantBase
}
// NewParticipant instantiates a new Participant and sets the default values.
func NewParticipant()(*Participant) {
    m := &Participant{
        ParticipantBase: *NewParticipantBase(),
    }
    odataTypeValue := "#microsoft.graph.callRecords.participant"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateParticipantFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateParticipantFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewParticipant(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Participant) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ParticipantBase.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *Participant) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ParticipantBase.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type Participantable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ParticipantBaseable
}

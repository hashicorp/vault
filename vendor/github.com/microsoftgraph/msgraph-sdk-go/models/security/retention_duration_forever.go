package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RetentionDurationForever struct {
    RetentionDuration
}
// NewRetentionDurationForever instantiates a new RetentionDurationForever and sets the default values.
func NewRetentionDurationForever()(*RetentionDurationForever) {
    m := &RetentionDurationForever{
        RetentionDuration: *NewRetentionDuration(),
    }
    odataTypeValue := "#microsoft.graph.security.retentionDurationForever"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRetentionDurationForeverFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRetentionDurationForeverFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRetentionDurationForever(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RetentionDurationForever) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RetentionDuration.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *RetentionDurationForever) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RetentionDuration.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type RetentionDurationForeverable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RetentionDurationable
}

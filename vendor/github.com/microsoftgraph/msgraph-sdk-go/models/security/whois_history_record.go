package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WhoisHistoryRecord struct {
    WhoisBaseRecord
}
// NewWhoisHistoryRecord instantiates a new WhoisHistoryRecord and sets the default values.
func NewWhoisHistoryRecord()(*WhoisHistoryRecord) {
    m := &WhoisHistoryRecord{
        WhoisBaseRecord: *NewWhoisBaseRecord(),
    }
    odataTypeValue := "#microsoft.graph.security.whoisHistoryRecord"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWhoisHistoryRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWhoisHistoryRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWhoisHistoryRecord(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WhoisHistoryRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WhoisBaseRecord.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *WhoisHistoryRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WhoisBaseRecord.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type WhoisHistoryRecordable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WhoisBaseRecordable
}

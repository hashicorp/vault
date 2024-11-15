package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WhoisRecord struct {
    WhoisBaseRecord
}
// NewWhoisRecord instantiates a new WhoisRecord and sets the default values.
func NewWhoisRecord()(*WhoisRecord) {
    m := &WhoisRecord{
        WhoisBaseRecord: *NewWhoisBaseRecord(),
    }
    odataTypeValue := "#microsoft.graph.security.whoisRecord"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWhoisRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWhoisRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWhoisRecord(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WhoisRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WhoisBaseRecord.GetFieldDeserializers()
    res["history"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWhoisHistoryRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WhoisHistoryRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WhoisHistoryRecordable)
                }
            }
            m.SetHistory(res)
        }
        return nil
    }
    return res
}
// GetHistory gets the history property value. The collection of historical records associated to this WHOIS object.
// returns a []WhoisHistoryRecordable when successful
func (m *WhoisRecord) GetHistory()([]WhoisHistoryRecordable) {
    val, err := m.GetBackingStore().Get("history")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WhoisHistoryRecordable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WhoisRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WhoisBaseRecord.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetHistory() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHistory()))
        for i, v := range m.GetHistory() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("history", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHistory sets the history property value. The collection of historical records associated to this WHOIS object.
func (m *WhoisRecord) SetHistory(value []WhoisHistoryRecordable)() {
    err := m.GetBackingStore().Set("history", value)
    if err != nil {
        panic(err)
    }
}
type WhoisRecordable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WhoisBaseRecordable
    GetHistory()([]WhoisHistoryRecordable)
    SetHistory(value []WhoisHistoryRecordable)()
}

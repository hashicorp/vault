package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DomainDnsMxRecord struct {
    DomainDnsRecord
}
// NewDomainDnsMxRecord instantiates a new DomainDnsMxRecord and sets the default values.
func NewDomainDnsMxRecord()(*DomainDnsMxRecord) {
    m := &DomainDnsMxRecord{
        DomainDnsRecord: *NewDomainDnsRecord(),
    }
    return m
}
// CreateDomainDnsMxRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDomainDnsMxRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDomainDnsMxRecord(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DomainDnsMxRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DomainDnsRecord.GetFieldDeserializers()
    res["mailExchange"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailExchange(val)
        }
        return nil
    }
    res["preference"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreference(val)
        }
        return nil
    }
    return res
}
// GetMailExchange gets the mailExchange property value. Value used when configuring the answer/destination/value of the MX record at the DNS host.
// returns a *string when successful
func (m *DomainDnsMxRecord) GetMailExchange()(*string) {
    val, err := m.GetBackingStore().Get("mailExchange")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreference gets the preference property value. Value used when configuring the Preference/Priority property of the MX record at the DNS host.
// returns a *int32 when successful
func (m *DomainDnsMxRecord) GetPreference()(*int32) {
    val, err := m.GetBackingStore().Get("preference")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DomainDnsMxRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DomainDnsRecord.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("mailExchange", m.GetMailExchange())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("preference", m.GetPreference())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMailExchange sets the mailExchange property value. Value used when configuring the answer/destination/value of the MX record at the DNS host.
func (m *DomainDnsMxRecord) SetMailExchange(value *string)() {
    err := m.GetBackingStore().Set("mailExchange", value)
    if err != nil {
        panic(err)
    }
}
// SetPreference sets the preference property value. Value used when configuring the Preference/Priority property of the MX record at the DNS host.
func (m *DomainDnsMxRecord) SetPreference(value *int32)() {
    err := m.GetBackingStore().Set("preference", value)
    if err != nil {
        panic(err)
    }
}
type DomainDnsMxRecordable interface {
    DomainDnsRecordable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMailExchange()(*string)
    GetPreference()(*int32)
    SetMailExchange(value *string)()
    SetPreference(value *int32)()
}

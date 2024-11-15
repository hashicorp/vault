package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CustomTimeZone struct {
    TimeZoneBase
}
// NewCustomTimeZone instantiates a new CustomTimeZone and sets the default values.
func NewCustomTimeZone()(*CustomTimeZone) {
    m := &CustomTimeZone{
        TimeZoneBase: *NewTimeZoneBase(),
    }
    odataTypeValue := "#microsoft.graph.customTimeZone"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCustomTimeZoneFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomTimeZoneFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomTimeZone(), nil
}
// GetBias gets the bias property value. The time offset of the time zone from Coordinated Universal Time (UTC). This value is in minutes. Time zones that are ahead of UTC have a positive offset; time zones that are behind UTC have a negative offset.
// returns a *int32 when successful
func (m *CustomTimeZone) GetBias()(*int32) {
    val, err := m.GetBackingStore().Get("bias")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDaylightOffset gets the daylightOffset property value. Specifies when the time zone switches from standard time to daylight saving time.
// returns a DaylightTimeZoneOffsetable when successful
func (m *CustomTimeZone) GetDaylightOffset()(DaylightTimeZoneOffsetable) {
    val, err := m.GetBackingStore().Get("daylightOffset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DaylightTimeZoneOffsetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomTimeZone) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TimeZoneBase.GetFieldDeserializers()
    res["bias"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBias(val)
        }
        return nil
    }
    res["daylightOffset"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDaylightTimeZoneOffsetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDaylightOffset(val.(DaylightTimeZoneOffsetable))
        }
        return nil
    }
    res["standardOffset"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateStandardTimeZoneOffsetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStandardOffset(val.(StandardTimeZoneOffsetable))
        }
        return nil
    }
    return res
}
// GetStandardOffset gets the standardOffset property value. Specifies when the time zone switches from daylight saving time to standard time.
// returns a StandardTimeZoneOffsetable when successful
func (m *CustomTimeZone) GetStandardOffset()(StandardTimeZoneOffsetable) {
    val, err := m.GetBackingStore().Get("standardOffset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(StandardTimeZoneOffsetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CustomTimeZone) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TimeZoneBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("bias", m.GetBias())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("daylightOffset", m.GetDaylightOffset())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("standardOffset", m.GetStandardOffset())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBias sets the bias property value. The time offset of the time zone from Coordinated Universal Time (UTC). This value is in minutes. Time zones that are ahead of UTC have a positive offset; time zones that are behind UTC have a negative offset.
func (m *CustomTimeZone) SetBias(value *int32)() {
    err := m.GetBackingStore().Set("bias", value)
    if err != nil {
        panic(err)
    }
}
// SetDaylightOffset sets the daylightOffset property value. Specifies when the time zone switches from standard time to daylight saving time.
func (m *CustomTimeZone) SetDaylightOffset(value DaylightTimeZoneOffsetable)() {
    err := m.GetBackingStore().Set("daylightOffset", value)
    if err != nil {
        panic(err)
    }
}
// SetStandardOffset sets the standardOffset property value. Specifies when the time zone switches from daylight saving time to standard time.
func (m *CustomTimeZone) SetStandardOffset(value StandardTimeZoneOffsetable)() {
    err := m.GetBackingStore().Set("standardOffset", value)
    if err != nil {
        panic(err)
    }
}
type CustomTimeZoneable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    TimeZoneBaseable
    GetBias()(*int32)
    GetDaylightOffset()(DaylightTimeZoneOffsetable)
    GetStandardOffset()(StandardTimeZoneOffsetable)
    SetBias(value *int32)()
    SetDaylightOffset(value DaylightTimeZoneOffsetable)()
    SetStandardOffset(value StandardTimeZoneOffsetable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AttendanceRecord struct {
    Entity
}
// NewAttendanceRecord instantiates a new AttendanceRecord and sets the default values.
func NewAttendanceRecord()(*AttendanceRecord) {
    m := &AttendanceRecord{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAttendanceRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttendanceRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttendanceRecord(), nil
}
// GetAttendanceIntervals gets the attendanceIntervals property value. List of time periods between joining and leaving a meeting.
// returns a []AttendanceIntervalable when successful
func (m *AttendanceRecord) GetAttendanceIntervals()([]AttendanceIntervalable) {
    val, err := m.GetBackingStore().Get("attendanceIntervals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AttendanceIntervalable)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. Email address of the user associated with this attendance record.
// returns a *string when successful
func (m *AttendanceRecord) GetEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("emailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AttendanceRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attendanceIntervals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttendanceIntervalFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AttendanceIntervalable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AttendanceIntervalable)
                }
            }
            m.SetAttendanceIntervals(res)
        }
        return nil
    }
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val)
        }
        return nil
    }
    res["identity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentity(val.(Identityable))
        }
        return nil
    }
    res["role"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRole(val)
        }
        return nil
    }
    res["totalAttendanceInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalAttendanceInSeconds(val)
        }
        return nil
    }
    return res
}
// GetIdentity gets the identity property value. Identity of the user associated with this attendance record.
// returns a Identityable when successful
func (m *AttendanceRecord) GetIdentity()(Identityable) {
    val, err := m.GetBackingStore().Get("identity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Identityable)
    }
    return nil
}
// GetRole gets the role property value. Role of the attendee. Possible values are: None, Attendee, Presenter, and Organizer.
// returns a *string when successful
func (m *AttendanceRecord) GetRole()(*string) {
    val, err := m.GetBackingStore().Get("role")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalAttendanceInSeconds gets the totalAttendanceInSeconds property value. Total duration of the attendances in seconds.
// returns a *int32 when successful
func (m *AttendanceRecord) GetTotalAttendanceInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("totalAttendanceInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AttendanceRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttendanceIntervals() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttendanceIntervals()))
        for i, v := range m.GetAttendanceIntervals() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("attendanceIntervals", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("identity", m.GetIdentity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("role", m.GetRole())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalAttendanceInSeconds", m.GetTotalAttendanceInSeconds())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttendanceIntervals sets the attendanceIntervals property value. List of time periods between joining and leaving a meeting.
func (m *AttendanceRecord) SetAttendanceIntervals(value []AttendanceIntervalable)() {
    err := m.GetBackingStore().Set("attendanceIntervals", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. Email address of the user associated with this attendance record.
func (m *AttendanceRecord) SetEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentity sets the identity property value. Identity of the user associated with this attendance record.
func (m *AttendanceRecord) SetIdentity(value Identityable)() {
    err := m.GetBackingStore().Set("identity", value)
    if err != nil {
        panic(err)
    }
}
// SetRole sets the role property value. Role of the attendee. Possible values are: None, Attendee, Presenter, and Organizer.
func (m *AttendanceRecord) SetRole(value *string)() {
    err := m.GetBackingStore().Set("role", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalAttendanceInSeconds sets the totalAttendanceInSeconds property value. Total duration of the attendances in seconds.
func (m *AttendanceRecord) SetTotalAttendanceInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("totalAttendanceInSeconds", value)
    if err != nil {
        panic(err)
    }
}
type AttendanceRecordable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttendanceIntervals()([]AttendanceIntervalable)
    GetEmailAddress()(*string)
    GetIdentity()(Identityable)
    GetRole()(*string)
    GetTotalAttendanceInSeconds()(*int32)
    SetAttendanceIntervals(value []AttendanceIntervalable)()
    SetEmailAddress(value *string)()
    SetIdentity(value Identityable)()
    SetRole(value *string)()
    SetTotalAttendanceInSeconds(value *int32)()
}

package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MeetingAttendanceReport struct {
    Entity
}
// NewMeetingAttendanceReport instantiates a new MeetingAttendanceReport and sets the default values.
func NewMeetingAttendanceReport()(*MeetingAttendanceReport) {
    m := &MeetingAttendanceReport{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMeetingAttendanceReportFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMeetingAttendanceReportFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingAttendanceReport(), nil
}
// GetAttendanceRecords gets the attendanceRecords property value. List of attendance records of an attendance report. Read-only.
// returns a []AttendanceRecordable when successful
func (m *MeetingAttendanceReport) GetAttendanceRecords()([]AttendanceRecordable) {
    val, err := m.GetBackingStore().Get("attendanceRecords")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AttendanceRecordable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MeetingAttendanceReport) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attendanceRecords"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttendanceRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AttendanceRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AttendanceRecordable)
                }
            }
            m.SetAttendanceRecords(res)
        }
        return nil
    }
    res["meetingEndDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingEndDateTime(val)
        }
        return nil
    }
    res["meetingStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingStartDateTime(val)
        }
        return nil
    }
    res["totalParticipantCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalParticipantCount(val)
        }
        return nil
    }
    return res
}
// GetMeetingEndDateTime gets the meetingEndDateTime property value. UTC time when the meeting ended. Read-only.
// returns a *Time when successful
func (m *MeetingAttendanceReport) GetMeetingEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("meetingEndDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMeetingStartDateTime gets the meetingStartDateTime property value. UTC time when the meeting started. Read-only.
// returns a *Time when successful
func (m *MeetingAttendanceReport) GetMeetingStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("meetingStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTotalParticipantCount gets the totalParticipantCount property value. Total number of participants. Read-only.
// returns a *int32 when successful
func (m *MeetingAttendanceReport) GetTotalParticipantCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalParticipantCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MeetingAttendanceReport) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttendanceRecords() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttendanceRecords()))
        for i, v := range m.GetAttendanceRecords() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("attendanceRecords", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("meetingEndDateTime", m.GetMeetingEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("meetingStartDateTime", m.GetMeetingStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalParticipantCount", m.GetTotalParticipantCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttendanceRecords sets the attendanceRecords property value. List of attendance records of an attendance report. Read-only.
func (m *MeetingAttendanceReport) SetAttendanceRecords(value []AttendanceRecordable)() {
    err := m.GetBackingStore().Set("attendanceRecords", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingEndDateTime sets the meetingEndDateTime property value. UTC time when the meeting ended. Read-only.
func (m *MeetingAttendanceReport) SetMeetingEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("meetingEndDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingStartDateTime sets the meetingStartDateTime property value. UTC time when the meeting started. Read-only.
func (m *MeetingAttendanceReport) SetMeetingStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("meetingStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalParticipantCount sets the totalParticipantCount property value. Total number of participants. Read-only.
func (m *MeetingAttendanceReport) SetTotalParticipantCount(value *int32)() {
    err := m.GetBackingStore().Set("totalParticipantCount", value)
    if err != nil {
        panic(err)
    }
}
type MeetingAttendanceReportable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttendanceRecords()([]AttendanceRecordable)
    GetMeetingEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMeetingStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTotalParticipantCount()(*int32)
    SetAttendanceRecords(value []AttendanceRecordable)()
    SetMeetingEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMeetingStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTotalParticipantCount(value *int32)()
}

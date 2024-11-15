package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WindowsUpdateScheduledInstall struct {
    WindowsUpdateInstallScheduleType
}
// NewWindowsUpdateScheduledInstall instantiates a new WindowsUpdateScheduledInstall and sets the default values.
func NewWindowsUpdateScheduledInstall()(*WindowsUpdateScheduledInstall) {
    m := &WindowsUpdateScheduledInstall{
        WindowsUpdateInstallScheduleType: *NewWindowsUpdateInstallScheduleType(),
    }
    odataTypeValue := "#microsoft.graph.windowsUpdateScheduledInstall"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsUpdateScheduledInstallFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsUpdateScheduledInstallFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUpdateScheduledInstall(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsUpdateScheduledInstall) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WindowsUpdateInstallScheduleType.GetFieldDeserializers()
    res["scheduledInstallDay"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWeeklySchedule)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduledInstallDay(val.(*WeeklySchedule))
        }
        return nil
    }
    res["scheduledInstallTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduledInstallTime(val)
        }
        return nil
    }
    return res
}
// GetScheduledInstallDay gets the scheduledInstallDay property value. Possible values for a weekly schedule.
// returns a *WeeklySchedule when successful
func (m *WindowsUpdateScheduledInstall) GetScheduledInstallDay()(*WeeklySchedule) {
    val, err := m.GetBackingStore().Get("scheduledInstallDay")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WeeklySchedule)
    }
    return nil
}
// GetScheduledInstallTime gets the scheduledInstallTime property value. Scheduled Install Time during day
// returns a *TimeOnly when successful
func (m *WindowsUpdateScheduledInstall) GetScheduledInstallTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    val, err := m.GetBackingStore().Get("scheduledInstallTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsUpdateScheduledInstall) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WindowsUpdateInstallScheduleType.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetScheduledInstallDay() != nil {
        cast := (*m.GetScheduledInstallDay()).String()
        err = writer.WriteStringValue("scheduledInstallDay", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeOnlyValue("scheduledInstallTime", m.GetScheduledInstallTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetScheduledInstallDay sets the scheduledInstallDay property value. Possible values for a weekly schedule.
func (m *WindowsUpdateScheduledInstall) SetScheduledInstallDay(value *WeeklySchedule)() {
    err := m.GetBackingStore().Set("scheduledInstallDay", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduledInstallTime sets the scheduledInstallTime property value. Scheduled Install Time during day
func (m *WindowsUpdateScheduledInstall) SetScheduledInstallTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    err := m.GetBackingStore().Set("scheduledInstallTime", value)
    if err != nil {
        panic(err)
    }
}
type WindowsUpdateScheduledInstallable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WindowsUpdateInstallScheduleTypeable
    GetScheduledInstallDay()(*WeeklySchedule)
    GetScheduledInstallTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    SetScheduledInstallDay(value *WeeklySchedule)()
    SetScheduledInstallTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)()
}

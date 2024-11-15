package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsMetricHistory the user experience analytics metric history.
type UserExperienceAnalyticsMetricHistory struct {
    Entity
}
// NewUserExperienceAnalyticsMetricHistory instantiates a new UserExperienceAnalyticsMetricHistory and sets the default values.
func NewUserExperienceAnalyticsMetricHistory()(*UserExperienceAnalyticsMetricHistory) {
    m := &UserExperienceAnalyticsMetricHistory{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsMetricHistoryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsMetricHistoryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsMetricHistory(), nil
}
// GetDeviceId gets the deviceId property value. The Intune device id of the device.
// returns a *string when successful
func (m *UserExperienceAnalyticsMetricHistory) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
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
func (m *UserExperienceAnalyticsMetricHistory) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceId(val)
        }
        return nil
    }
    res["metricDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMetricDateTime(val)
        }
        return nil
    }
    res["metricType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMetricType(val)
        }
        return nil
    }
    return res
}
// GetMetricDateTime gets the metricDateTime property value. The metric date time. The value cannot be modified and is automatically populated when the metric is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'. Returned by default.
// returns a *Time when successful
func (m *UserExperienceAnalyticsMetricHistory) GetMetricDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("metricDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMetricType gets the metricType property value. The user experience analytics metric type.
// returns a *string when successful
func (m *UserExperienceAnalyticsMetricHistory) GetMetricType()(*string) {
    val, err := m.GetBackingStore().Get("metricType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsMetricHistory) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("metricDateTime", m.GetMetricDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("metricType", m.GetMetricType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceId sets the deviceId property value. The Intune device id of the device.
func (m *UserExperienceAnalyticsMetricHistory) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetMetricDateTime sets the metricDateTime property value. The metric date time. The value cannot be modified and is automatically populated when the metric is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'. Returned by default.
func (m *UserExperienceAnalyticsMetricHistory) SetMetricDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("metricDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMetricType sets the metricType property value. The user experience analytics metric type.
func (m *UserExperienceAnalyticsMetricHistory) SetMetricType(value *string)() {
    err := m.GetBackingStore().Set("metricType", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsMetricHistoryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceId()(*string)
    GetMetricDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMetricType()(*string)
    SetDeviceId(value *string)()
    SetMetricDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMetricType(value *string)()
}

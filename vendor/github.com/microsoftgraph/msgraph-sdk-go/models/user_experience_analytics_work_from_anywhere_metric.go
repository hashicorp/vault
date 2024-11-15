package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsWorkFromAnywhereMetric the user experience analytics metric for work from anywhere report.
type UserExperienceAnalyticsWorkFromAnywhereMetric struct {
    Entity
}
// NewUserExperienceAnalyticsWorkFromAnywhereMetric instantiates a new UserExperienceAnalyticsWorkFromAnywhereMetric and sets the default values.
func NewUserExperienceAnalyticsWorkFromAnywhereMetric()(*UserExperienceAnalyticsWorkFromAnywhereMetric) {
    m := &UserExperienceAnalyticsWorkFromAnywhereMetric{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsWorkFromAnywhereMetricFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsWorkFromAnywhereMetricFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsWorkFromAnywhereMetric(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetric) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["metricDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserExperienceAnalyticsWorkFromAnywhereDeviceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserExperienceAnalyticsWorkFromAnywhereDeviceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserExperienceAnalyticsWorkFromAnywhereDeviceable)
                }
            }
            m.SetMetricDevices(res)
        }
        return nil
    }
    return res
}
// GetMetricDevices gets the metricDevices property value. The work from anywhere metric devices. Read-only.
// returns a []UserExperienceAnalyticsWorkFromAnywhereDeviceable when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetric) GetMetricDevices()([]UserExperienceAnalyticsWorkFromAnywhereDeviceable) {
    val, err := m.GetBackingStore().Get("metricDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserExperienceAnalyticsWorkFromAnywhereDeviceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsWorkFromAnywhereMetric) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMetricDevices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMetricDevices()))
        for i, v := range m.GetMetricDevices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("metricDevices", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMetricDevices sets the metricDevices property value. The work from anywhere metric devices. Read-only.
func (m *UserExperienceAnalyticsWorkFromAnywhereMetric) SetMetricDevices(value []UserExperienceAnalyticsWorkFromAnywhereDeviceable)() {
    err := m.GetBackingStore().Set("metricDevices", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsWorkFromAnywhereMetricable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMetricDevices()([]UserExperienceAnalyticsWorkFromAnywhereDeviceable)
    SetMetricDevices(value []UserExperienceAnalyticsWorkFromAnywhereDeviceable)()
}

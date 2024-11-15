package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsMetric the user experience analytics metric contains the score and units of a metric of a user experience anlaytics category.
type UserExperienceAnalyticsMetric struct {
    Entity
}
// NewUserExperienceAnalyticsMetric instantiates a new UserExperienceAnalyticsMetric and sets the default values.
func NewUserExperienceAnalyticsMetric()(*UserExperienceAnalyticsMetric) {
    m := &UserExperienceAnalyticsMetric{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsMetricFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsMetricFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsMetric(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsMetric) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["unit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnit(val)
        }
        return nil
    }
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValue(val)
        }
        return nil
    }
    return res
}
// GetUnit gets the unit property value. The unit of the user experience analytics metric. Examples: none, percentage, count, seconds, score.
// returns a *string when successful
func (m *UserExperienceAnalyticsMetric) GetUnit()(*string) {
    val, err := m.GetBackingStore().Get("unit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetValue gets the value property value. The value of the user experience analytics metric.
// returns a *float64 when successful
func (m *UserExperienceAnalyticsMetric) GetValue()(*float64) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsMetric) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("unit", m.GetUnit())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUnit sets the unit property value. The unit of the user experience analytics metric. Examples: none, percentage, count, seconds, score.
func (m *UserExperienceAnalyticsMetric) SetUnit(value *string)() {
    err := m.GetBackingStore().Set("unit", value)
    if err != nil {
        panic(err)
    }
}
// SetValue sets the value property value. The value of the user experience analytics metric.
func (m *UserExperienceAnalyticsMetric) SetValue(value *float64)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsMetricable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUnit()(*string)
    GetValue()(*float64)
    SetUnit(value *string)()
    SetValue(value *float64)()
}

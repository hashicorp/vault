package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsCategory the user experience analytics category entity contains the scores and insights for the various metrics of a category.
type UserExperienceAnalyticsCategory struct {
    Entity
}
// NewUserExperienceAnalyticsCategory instantiates a new UserExperienceAnalyticsCategory and sets the default values.
func NewUserExperienceAnalyticsCategory()(*UserExperienceAnalyticsCategory) {
    m := &UserExperienceAnalyticsCategory{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsCategoryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsCategoryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsCategory(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsCategory) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["insights"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserExperienceAnalyticsInsightFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserExperienceAnalyticsInsightable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserExperienceAnalyticsInsightable)
                }
            }
            m.SetInsights(res)
        }
        return nil
    }
    res["metricValues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserExperienceAnalyticsMetricFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserExperienceAnalyticsMetricable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserExperienceAnalyticsMetricable)
                }
            }
            m.SetMetricValues(res)
        }
        return nil
    }
    return res
}
// GetInsights gets the insights property value. The insights for the category. Read-only.
// returns a []UserExperienceAnalyticsInsightable when successful
func (m *UserExperienceAnalyticsCategory) GetInsights()([]UserExperienceAnalyticsInsightable) {
    val, err := m.GetBackingStore().Get("insights")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserExperienceAnalyticsInsightable)
    }
    return nil
}
// GetMetricValues gets the metricValues property value. The metric values for the user experience analytics category. Read-only.
// returns a []UserExperienceAnalyticsMetricable when successful
func (m *UserExperienceAnalyticsCategory) GetMetricValues()([]UserExperienceAnalyticsMetricable) {
    val, err := m.GetBackingStore().Get("metricValues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserExperienceAnalyticsMetricable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsCategory) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetInsights() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInsights()))
        for i, v := range m.GetInsights() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("insights", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMetricValues() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMetricValues()))
        for i, v := range m.GetMetricValues() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("metricValues", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInsights sets the insights property value. The insights for the category. Read-only.
func (m *UserExperienceAnalyticsCategory) SetInsights(value []UserExperienceAnalyticsInsightable)() {
    err := m.GetBackingStore().Set("insights", value)
    if err != nil {
        panic(err)
    }
}
// SetMetricValues sets the metricValues property value. The metric values for the user experience analytics category. Read-only.
func (m *UserExperienceAnalyticsCategory) SetMetricValues(value []UserExperienceAnalyticsMetricable)() {
    err := m.GetBackingStore().Set("metricValues", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsCategoryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInsights()([]UserExperienceAnalyticsInsightable)
    GetMetricValues()([]UserExperienceAnalyticsMetricable)
    SetInsights(value []UserExperienceAnalyticsInsightable)()
    SetMetricValues(value []UserExperienceAnalyticsMetricable)()
}

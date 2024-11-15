package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsOverview the user experience analytics overview entity contains the overall score and the scores and insights of every metric of all categories.
type UserExperienceAnalyticsOverview struct {
    Entity
}
// NewUserExperienceAnalyticsOverview instantiates a new UserExperienceAnalyticsOverview and sets the default values.
func NewUserExperienceAnalyticsOverview()(*UserExperienceAnalyticsOverview) {
    m := &UserExperienceAnalyticsOverview{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsOverview(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    return res
}
// GetInsights gets the insights property value. The user experience analytics insights. Read-only.
// returns a []UserExperienceAnalyticsInsightable when successful
func (m *UserExperienceAnalyticsOverview) GetInsights()([]UserExperienceAnalyticsInsightable) {
    val, err := m.GetBackingStore().Get("insights")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserExperienceAnalyticsInsightable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    return nil
}
// SetInsights sets the insights property value. The user experience analytics insights. Read-only.
func (m *UserExperienceAnalyticsOverview) SetInsights(value []UserExperienceAnalyticsInsightable)() {
    err := m.GetBackingStore().Set("insights", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsOverviewable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInsights()([]UserExperienceAnalyticsInsightable)
    SetInsights(value []UserExperienceAnalyticsInsightable)()
}

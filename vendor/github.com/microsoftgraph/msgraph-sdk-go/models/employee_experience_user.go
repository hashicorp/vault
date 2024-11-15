package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EmployeeExperienceUser struct {
    Entity
}
// NewEmployeeExperienceUser instantiates a new EmployeeExperienceUser and sets the default values.
func NewEmployeeExperienceUser()(*EmployeeExperienceUser) {
    m := &EmployeeExperienceUser{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEmployeeExperienceUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEmployeeExperienceUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmployeeExperienceUser(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EmployeeExperienceUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["learningCourseActivities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLearningCourseActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LearningCourseActivityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LearningCourseActivityable)
                }
            }
            m.SetLearningCourseActivities(res)
        }
        return nil
    }
    return res
}
// GetLearningCourseActivities gets the learningCourseActivities property value. The learningCourseActivities property
// returns a []LearningCourseActivityable when successful
func (m *EmployeeExperienceUser) GetLearningCourseActivities()([]LearningCourseActivityable) {
    val, err := m.GetBackingStore().Get("learningCourseActivities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LearningCourseActivityable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EmployeeExperienceUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLearningCourseActivities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLearningCourseActivities()))
        for i, v := range m.GetLearningCourseActivities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("learningCourseActivities", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLearningCourseActivities sets the learningCourseActivities property value. The learningCourseActivities property
func (m *EmployeeExperienceUser) SetLearningCourseActivities(value []LearningCourseActivityable)() {
    err := m.GetBackingStore().Set("learningCourseActivities", value)
    if err != nil {
        panic(err)
    }
}
type EmployeeExperienceUserable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLearningCourseActivities()([]LearningCourseActivityable)
    SetLearningCourseActivities(value []LearningCourseActivityable)()
}

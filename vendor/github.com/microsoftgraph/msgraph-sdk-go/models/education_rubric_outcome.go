package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationRubricOutcome struct {
    EducationOutcome
}
// NewEducationRubricOutcome instantiates a new EducationRubricOutcome and sets the default values.
func NewEducationRubricOutcome()(*EducationRubricOutcome) {
    m := &EducationRubricOutcome{
        EducationOutcome: *NewEducationOutcome(),
    }
    odataTypeValue := "#microsoft.graph.educationRubricOutcome"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationRubricOutcomeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationRubricOutcomeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationRubricOutcome(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationRubricOutcome) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationOutcome.GetFieldDeserializers()
    res["publishedRubricQualityFeedback"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRubricQualityFeedbackModelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RubricQualityFeedbackModelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RubricQualityFeedbackModelable)
                }
            }
            m.SetPublishedRubricQualityFeedback(res)
        }
        return nil
    }
    res["publishedRubricQualitySelectedLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRubricQualitySelectedColumnModelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RubricQualitySelectedColumnModelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RubricQualitySelectedColumnModelable)
                }
            }
            m.SetPublishedRubricQualitySelectedLevels(res)
        }
        return nil
    }
    res["rubricQualityFeedback"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRubricQualityFeedbackModelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RubricQualityFeedbackModelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RubricQualityFeedbackModelable)
                }
            }
            m.SetRubricQualityFeedback(res)
        }
        return nil
    }
    res["rubricQualitySelectedLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRubricQualitySelectedColumnModelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RubricQualitySelectedColumnModelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RubricQualitySelectedColumnModelable)
                }
            }
            m.SetRubricQualitySelectedLevels(res)
        }
        return nil
    }
    return res
}
// GetPublishedRubricQualityFeedback gets the publishedRubricQualityFeedback property value. A copy of the rubricQualityFeedback property that is made when the grade is released to the student.
// returns a []RubricQualityFeedbackModelable when successful
func (m *EducationRubricOutcome) GetPublishedRubricQualityFeedback()([]RubricQualityFeedbackModelable) {
    val, err := m.GetBackingStore().Get("publishedRubricQualityFeedback")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RubricQualityFeedbackModelable)
    }
    return nil
}
// GetPublishedRubricQualitySelectedLevels gets the publishedRubricQualitySelectedLevels property value. A copy of the rubricQualitySelectedLevels property that is made when the grade is released to the student.
// returns a []RubricQualitySelectedColumnModelable when successful
func (m *EducationRubricOutcome) GetPublishedRubricQualitySelectedLevels()([]RubricQualitySelectedColumnModelable) {
    val, err := m.GetBackingStore().Get("publishedRubricQualitySelectedLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RubricQualitySelectedColumnModelable)
    }
    return nil
}
// GetRubricQualityFeedback gets the rubricQualityFeedback property value. A collection of specific feedback for each quality of this rubric.
// returns a []RubricQualityFeedbackModelable when successful
func (m *EducationRubricOutcome) GetRubricQualityFeedback()([]RubricQualityFeedbackModelable) {
    val, err := m.GetBackingStore().Get("rubricQualityFeedback")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RubricQualityFeedbackModelable)
    }
    return nil
}
// GetRubricQualitySelectedLevels gets the rubricQualitySelectedLevels property value. The level that the teacher has selected for each quality while grading this assignment.
// returns a []RubricQualitySelectedColumnModelable when successful
func (m *EducationRubricOutcome) GetRubricQualitySelectedLevels()([]RubricQualitySelectedColumnModelable) {
    val, err := m.GetBackingStore().Get("rubricQualitySelectedLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RubricQualitySelectedColumnModelable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationRubricOutcome) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationOutcome.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetPublishedRubricQualityFeedback() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPublishedRubricQualityFeedback()))
        for i, v := range m.GetPublishedRubricQualityFeedback() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("publishedRubricQualityFeedback", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPublishedRubricQualitySelectedLevels() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPublishedRubricQualitySelectedLevels()))
        for i, v := range m.GetPublishedRubricQualitySelectedLevels() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("publishedRubricQualitySelectedLevels", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRubricQualityFeedback() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRubricQualityFeedback()))
        for i, v := range m.GetRubricQualityFeedback() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rubricQualityFeedback", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRubricQualitySelectedLevels() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRubricQualitySelectedLevels()))
        for i, v := range m.GetRubricQualitySelectedLevels() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rubricQualitySelectedLevels", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPublishedRubricQualityFeedback sets the publishedRubricQualityFeedback property value. A copy of the rubricQualityFeedback property that is made when the grade is released to the student.
func (m *EducationRubricOutcome) SetPublishedRubricQualityFeedback(value []RubricQualityFeedbackModelable)() {
    err := m.GetBackingStore().Set("publishedRubricQualityFeedback", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishedRubricQualitySelectedLevels sets the publishedRubricQualitySelectedLevels property value. A copy of the rubricQualitySelectedLevels property that is made when the grade is released to the student.
func (m *EducationRubricOutcome) SetPublishedRubricQualitySelectedLevels(value []RubricQualitySelectedColumnModelable)() {
    err := m.GetBackingStore().Set("publishedRubricQualitySelectedLevels", value)
    if err != nil {
        panic(err)
    }
}
// SetRubricQualityFeedback sets the rubricQualityFeedback property value. A collection of specific feedback for each quality of this rubric.
func (m *EducationRubricOutcome) SetRubricQualityFeedback(value []RubricQualityFeedbackModelable)() {
    err := m.GetBackingStore().Set("rubricQualityFeedback", value)
    if err != nil {
        panic(err)
    }
}
// SetRubricQualitySelectedLevels sets the rubricQualitySelectedLevels property value. The level that the teacher has selected for each quality while grading this assignment.
func (m *EducationRubricOutcome) SetRubricQualitySelectedLevels(value []RubricQualitySelectedColumnModelable)() {
    err := m.GetBackingStore().Set("rubricQualitySelectedLevels", value)
    if err != nil {
        panic(err)
    }
}
type EducationRubricOutcomeable interface {
    EducationOutcomeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPublishedRubricQualityFeedback()([]RubricQualityFeedbackModelable)
    GetPublishedRubricQualitySelectedLevels()([]RubricQualitySelectedColumnModelable)
    GetRubricQualityFeedback()([]RubricQualityFeedbackModelable)
    GetRubricQualitySelectedLevels()([]RubricQualitySelectedColumnModelable)
    SetPublishedRubricQualityFeedback(value []RubricQualityFeedbackModelable)()
    SetPublishedRubricQualitySelectedLevels(value []RubricQualitySelectedColumnModelable)()
    SetRubricQualityFeedback(value []RubricQualityFeedbackModelable)()
    SetRubricQualitySelectedLevels(value []RubricQualitySelectedColumnModelable)()
}

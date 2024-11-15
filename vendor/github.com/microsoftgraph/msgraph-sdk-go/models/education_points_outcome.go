package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationPointsOutcome struct {
    EducationOutcome
}
// NewEducationPointsOutcome instantiates a new EducationPointsOutcome and sets the default values.
func NewEducationPointsOutcome()(*EducationPointsOutcome) {
    m := &EducationPointsOutcome{
        EducationOutcome: *NewEducationOutcome(),
    }
    odataTypeValue := "#microsoft.graph.educationPointsOutcome"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationPointsOutcomeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationPointsOutcomeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationPointsOutcome(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationPointsOutcome) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationOutcome.GetFieldDeserializers()
    res["points"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentPointsGradeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPoints(val.(EducationAssignmentPointsGradeable))
        }
        return nil
    }
    res["publishedPoints"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentPointsGradeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublishedPoints(val.(EducationAssignmentPointsGradeable))
        }
        return nil
    }
    return res
}
// GetPoints gets the points property value. The numeric grade the teacher has given the student for this assignment.
// returns a EducationAssignmentPointsGradeable when successful
func (m *EducationPointsOutcome) GetPoints()(EducationAssignmentPointsGradeable) {
    val, err := m.GetBackingStore().Get("points")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentPointsGradeable)
    }
    return nil
}
// GetPublishedPoints gets the publishedPoints property value. A copy of the points property that is made when the grade is released to the student.
// returns a EducationAssignmentPointsGradeable when successful
func (m *EducationPointsOutcome) GetPublishedPoints()(EducationAssignmentPointsGradeable) {
    val, err := m.GetBackingStore().Get("publishedPoints")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentPointsGradeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationPointsOutcome) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationOutcome.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("points", m.GetPoints())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("publishedPoints", m.GetPublishedPoints())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPoints sets the points property value. The numeric grade the teacher has given the student for this assignment.
func (m *EducationPointsOutcome) SetPoints(value EducationAssignmentPointsGradeable)() {
    err := m.GetBackingStore().Set("points", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishedPoints sets the publishedPoints property value. A copy of the points property that is made when the grade is released to the student.
func (m *EducationPointsOutcome) SetPublishedPoints(value EducationAssignmentPointsGradeable)() {
    err := m.GetBackingStore().Set("publishedPoints", value)
    if err != nil {
        panic(err)
    }
}
type EducationPointsOutcomeable interface {
    EducationOutcomeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPoints()(EducationAssignmentPointsGradeable)
    GetPublishedPoints()(EducationAssignmentPointsGradeable)
    SetPoints(value EducationAssignmentPointsGradeable)()
    SetPublishedPoints(value EducationAssignmentPointsGradeable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationAssignmentPointsGradeType struct {
    EducationAssignmentGradeType
}
// NewEducationAssignmentPointsGradeType instantiates a new EducationAssignmentPointsGradeType and sets the default values.
func NewEducationAssignmentPointsGradeType()(*EducationAssignmentPointsGradeType) {
    m := &EducationAssignmentPointsGradeType{
        EducationAssignmentGradeType: *NewEducationAssignmentGradeType(),
    }
    odataTypeValue := "#microsoft.graph.educationAssignmentPointsGradeType"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationAssignmentPointsGradeTypeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationAssignmentPointsGradeTypeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignmentPointsGradeType(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationAssignmentPointsGradeType) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationAssignmentGradeType.GetFieldDeserializers()
    res["maxPoints"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxPoints(val)
        }
        return nil
    }
    return res
}
// GetMaxPoints gets the maxPoints property value. Max points possible for this assignment.
// returns a *float32 when successful
func (m *EducationAssignmentPointsGradeType) GetMaxPoints()(*float32) {
    val, err := m.GetBackingStore().Get("maxPoints")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationAssignmentPointsGradeType) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationAssignmentGradeType.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteFloat32Value("maxPoints", m.GetMaxPoints())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMaxPoints sets the maxPoints property value. Max points possible for this assignment.
func (m *EducationAssignmentPointsGradeType) SetMaxPoints(value *float32)() {
    err := m.GetBackingStore().Set("maxPoints", value)
    if err != nil {
        panic(err)
    }
}
type EducationAssignmentPointsGradeTypeable interface {
    EducationAssignmentGradeTypeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMaxPoints()(*float32)
    SetMaxPoints(value *float32)()
}

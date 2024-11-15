package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationAssignmentResource struct {
    Entity
}
// NewEducationAssignmentResource instantiates a new EducationAssignmentResource and sets the default values.
func NewEducationAssignmentResource()(*EducationAssignmentResource) {
    m := &EducationAssignmentResource{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationAssignmentResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationAssignmentResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignmentResource(), nil
}
// GetDistributeForStudentWork gets the distributeForStudentWork property value. Indicates whether this resource should be copied to each student submission for modification and submission. Required
// returns a *bool when successful
func (m *EducationAssignmentResource) GetDistributeForStudentWork()(*bool) {
    val, err := m.GetBackingStore().Get("distributeForStudentWork")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationAssignmentResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["distributeForStudentWork"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDistributeForStudentWork(val)
        }
        return nil
    }
    res["resource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResource(val.(EducationResourceable))
        }
        return nil
    }
    return res
}
// GetResource gets the resource property value. Resource object that has been associated with this assignment.
// returns a EducationResourceable when successful
func (m *EducationAssignmentResource) GetResource()(EducationResourceable) {
    val, err := m.GetBackingStore().Get("resource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationResourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationAssignmentResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("distributeForStudentWork", m.GetDistributeForStudentWork())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resource", m.GetResource())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDistributeForStudentWork sets the distributeForStudentWork property value. Indicates whether this resource should be copied to each student submission for modification and submission. Required
func (m *EducationAssignmentResource) SetDistributeForStudentWork(value *bool)() {
    err := m.GetBackingStore().Set("distributeForStudentWork", value)
    if err != nil {
        panic(err)
    }
}
// SetResource sets the resource property value. Resource object that has been associated with this assignment.
func (m *EducationAssignmentResource) SetResource(value EducationResourceable)() {
    err := m.GetBackingStore().Set("resource", value)
    if err != nil {
        panic(err)
    }
}
type EducationAssignmentResourceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDistributeForStudentWork()(*bool)
    GetResource()(EducationResourceable)
    SetDistributeForStudentWork(value *bool)()
    SetResource(value EducationResourceable)()
}

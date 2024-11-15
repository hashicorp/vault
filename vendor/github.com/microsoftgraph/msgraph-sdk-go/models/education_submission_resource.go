package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationSubmissionResource struct {
    Entity
}
// NewEducationSubmissionResource instantiates a new EducationSubmissionResource and sets the default values.
func NewEducationSubmissionResource()(*EducationSubmissionResource) {
    m := &EducationSubmissionResource{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationSubmissionResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationSubmissionResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationSubmissionResource(), nil
}
// GetAssignmentResourceUrl gets the assignmentResourceUrl property value. Pointer to the assignment from which the resource was copied, and if null, the student uploaded the resource.
// returns a *string when successful
func (m *EducationSubmissionResource) GetAssignmentResourceUrl()(*string) {
    val, err := m.GetBackingStore().Get("assignmentResourceUrl")
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
func (m *EducationSubmissionResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignmentResourceUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentResourceUrl(val)
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
// GetResource gets the resource property value. Resource object.
// returns a EducationResourceable when successful
func (m *EducationSubmissionResource) GetResource()(EducationResourceable) {
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
func (m *EducationSubmissionResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("assignmentResourceUrl", m.GetAssignmentResourceUrl())
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
// SetAssignmentResourceUrl sets the assignmentResourceUrl property value. Pointer to the assignment from which the resource was copied, and if null, the student uploaded the resource.
func (m *EducationSubmissionResource) SetAssignmentResourceUrl(value *string)() {
    err := m.GetBackingStore().Set("assignmentResourceUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetResource sets the resource property value. Resource object.
func (m *EducationSubmissionResource) SetResource(value EducationResourceable)() {
    err := m.GetBackingStore().Set("resource", value)
    if err != nil {
        panic(err)
    }
}
type EducationSubmissionResourceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignmentResourceUrl()(*string)
    GetResource()(EducationResourceable)
    SetAssignmentResourceUrl(value *string)()
    SetResource(value EducationResourceable)()
}

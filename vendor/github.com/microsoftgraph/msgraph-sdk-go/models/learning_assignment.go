package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type LearningAssignment struct {
    LearningCourseActivity
}
// NewLearningAssignment instantiates a new LearningAssignment and sets the default values.
func NewLearningAssignment()(*LearningAssignment) {
    m := &LearningAssignment{
        LearningCourseActivity: *NewLearningCourseActivity(),
    }
    return m
}
// CreateLearningAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLearningAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLearningAssignment(), nil
}
// GetAssignedDateTime gets the assignedDateTime property value. Assigned date for the course activity. Optional.
// returns a *Time when successful
func (m *LearningAssignment) GetAssignedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("assignedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAssignerUserId gets the assignerUserId property value. The user ID of the assigner. Optional.
// returns a *string when successful
func (m *LearningAssignment) GetAssignerUserId()(*string) {
    val, err := m.GetBackingStore().Get("assignerUserId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAssignmentType gets the assignmentType property value. The assignmentType property
// returns a *AssignmentType when successful
func (m *LearningAssignment) GetAssignmentType()(*AssignmentType) {
    val, err := m.GetBackingStore().Get("assignmentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AssignmentType)
    }
    return nil
}
// GetDueDateTime gets the dueDateTime property value. Due date for the course activity. Optional.
// returns a DateTimeTimeZoneable when successful
func (m *LearningAssignment) GetDueDateTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("dueDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LearningAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.LearningCourseActivity.GetFieldDeserializers()
    res["assignedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedDateTime(val)
        }
        return nil
    }
    res["assignerUserId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignerUserId(val)
        }
        return nil
    }
    res["assignmentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAssignmentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentType(val.(*AssignmentType))
        }
        return nil
    }
    res["dueDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDueDateTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["notes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotes(val.(ItemBodyable))
        }
        return nil
    }
    return res
}
// GetNotes gets the notes property value. Notes for the course activity. Optional.
// returns a ItemBodyable when successful
func (m *LearningAssignment) GetNotes()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("notes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LearningAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.LearningCourseActivity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("assignedDateTime", m.GetAssignedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assignerUserId", m.GetAssignerUserId())
        if err != nil {
            return err
        }
    }
    if m.GetAssignmentType() != nil {
        cast := (*m.GetAssignmentType()).String()
        err = writer.WriteStringValue("assignmentType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dueDateTime", m.GetDueDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("notes", m.GetNotes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignedDateTime sets the assignedDateTime property value. Assigned date for the course activity. Optional.
func (m *LearningAssignment) SetAssignedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("assignedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignerUserId sets the assignerUserId property value. The user ID of the assigner. Optional.
func (m *LearningAssignment) SetAssignerUserId(value *string)() {
    err := m.GetBackingStore().Set("assignerUserId", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentType sets the assignmentType property value. The assignmentType property
func (m *LearningAssignment) SetAssignmentType(value *AssignmentType)() {
    err := m.GetBackingStore().Set("assignmentType", value)
    if err != nil {
        panic(err)
    }
}
// SetDueDateTime sets the dueDateTime property value. Due date for the course activity. Optional.
func (m *LearningAssignment) SetDueDateTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("dueDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetNotes sets the notes property value. Notes for the course activity. Optional.
func (m *LearningAssignment) SetNotes(value ItemBodyable)() {
    err := m.GetBackingStore().Set("notes", value)
    if err != nil {
        panic(err)
    }
}
type LearningAssignmentable interface {
    LearningCourseActivityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAssignerUserId()(*string)
    GetAssignmentType()(*AssignmentType)
    GetDueDateTime()(DateTimeTimeZoneable)
    GetNotes()(ItemBodyable)
    SetAssignedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAssignerUserId(value *string)()
    SetAssignmentType(value *AssignmentType)()
    SetDueDateTime(value DateTimeTimeZoneable)()
    SetNotes(value ItemBodyable)()
}

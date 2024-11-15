package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type LearningCourseActivity struct {
    Entity
}
// NewLearningCourseActivity instantiates a new LearningCourseActivity and sets the default values.
func NewLearningCourseActivity()(*LearningCourseActivity) {
    m := &LearningCourseActivity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateLearningCourseActivityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLearningCourseActivityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.learningAssignment":
                        return NewLearningAssignment(), nil
                    case "#microsoft.graph.learningSelfInitiatedCourse":
                        return NewLearningSelfInitiatedCourse(), nil
                }
            }
        }
    }
    return NewLearningCourseActivity(), nil
}
// GetCompletedDateTime gets the completedDateTime property value. Date and time when the assignment was completed. Optional.
// returns a *Time when successful
func (m *LearningCourseActivity) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCompletionPercentage gets the completionPercentage property value. The percentage completion value of the course activity. Optional.
// returns a *int32 when successful
func (m *LearningCourseActivity) GetCompletionPercentage()(*int32) {
    val, err := m.GetBackingStore().Get("completionPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetExternalcourseActivityId gets the externalcourseActivityId property value. The externalcourseActivityId property
// returns a *string when successful
func (m *LearningCourseActivity) GetExternalcourseActivityId()(*string) {
    val, err := m.GetBackingStore().Get("externalcourseActivityId")
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
func (m *LearningCourseActivity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["completedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedDateTime(val)
        }
        return nil
    }
    res["completionPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletionPercentage(val)
        }
        return nil
    }
    res["externalcourseActivityId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalcourseActivityId(val)
        }
        return nil
    }
    res["learnerUserId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLearnerUserId(val)
        }
        return nil
    }
    res["learningContentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLearningContentId(val)
        }
        return nil
    }
    res["learningProviderId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLearningProviderId(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCourseStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*CourseStatus))
        }
        return nil
    }
    return res
}
// GetLearnerUserId gets the learnerUserId property value. The user ID of the learner to whom the activity is assigned. Required.
// returns a *string when successful
func (m *LearningCourseActivity) GetLearnerUserId()(*string) {
    val, err := m.GetBackingStore().Get("learnerUserId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLearningContentId gets the learningContentId property value. The ID of the learning content created in Viva Learning. Required.
// returns a *string when successful
func (m *LearningCourseActivity) GetLearningContentId()(*string) {
    val, err := m.GetBackingStore().Get("learningContentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLearningProviderId gets the learningProviderId property value. The registration ID of the provider. Required.
// returns a *string when successful
func (m *LearningCourseActivity) GetLearningProviderId()(*string) {
    val, err := m.GetBackingStore().Get("learningProviderId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status of the course activity. Possible values are: notStarted, inProgress, completed. Required.
// returns a *CourseStatus when successful
func (m *LearningCourseActivity) GetStatus()(*CourseStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CourseStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LearningCourseActivity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("completionPercentage", m.GetCompletionPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalcourseActivityId", m.GetExternalcourseActivityId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("learnerUserId", m.GetLearnerUserId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("learningContentId", m.GetLearningContentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("learningProviderId", m.GetLearningProviderId())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletedDateTime sets the completedDateTime property value. Date and time when the assignment was completed. Optional.
func (m *LearningCourseActivity) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletionPercentage sets the completionPercentage property value. The percentage completion value of the course activity. Optional.
func (m *LearningCourseActivity) SetCompletionPercentage(value *int32)() {
    err := m.GetBackingStore().Set("completionPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalcourseActivityId sets the externalcourseActivityId property value. The externalcourseActivityId property
func (m *LearningCourseActivity) SetExternalcourseActivityId(value *string)() {
    err := m.GetBackingStore().Set("externalcourseActivityId", value)
    if err != nil {
        panic(err)
    }
}
// SetLearnerUserId sets the learnerUserId property value. The user ID of the learner to whom the activity is assigned. Required.
func (m *LearningCourseActivity) SetLearnerUserId(value *string)() {
    err := m.GetBackingStore().Set("learnerUserId", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningContentId sets the learningContentId property value. The ID of the learning content created in Viva Learning. Required.
func (m *LearningCourseActivity) SetLearningContentId(value *string)() {
    err := m.GetBackingStore().Set("learningContentId", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningProviderId sets the learningProviderId property value. The registration ID of the provider. Required.
func (m *LearningCourseActivity) SetLearningProviderId(value *string)() {
    err := m.GetBackingStore().Set("learningProviderId", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of the course activity. Possible values are: notStarted, inProgress, completed. Required.
func (m *LearningCourseActivity) SetStatus(value *CourseStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type LearningCourseActivityable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCompletionPercentage()(*int32)
    GetExternalcourseActivityId()(*string)
    GetLearnerUserId()(*string)
    GetLearningContentId()(*string)
    GetLearningProviderId()(*string)
    GetStatus()(*CourseStatus)
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCompletionPercentage(value *int32)()
    SetExternalcourseActivityId(value *string)()
    SetLearnerUserId(value *string)()
    SetLearningContentId(value *string)()
    SetLearningProviderId(value *string)()
    SetStatus(value *CourseStatus)()
}

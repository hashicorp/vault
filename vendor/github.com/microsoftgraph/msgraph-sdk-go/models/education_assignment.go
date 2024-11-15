package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationAssignment struct {
    Entity
}
// NewEducationAssignment instantiates a new EducationAssignment and sets the default values.
func NewEducationAssignment()(*EducationAssignment) {
    m := &EducationAssignment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignment(), nil
}
// GetAddedStudentAction gets the addedStudentAction property value. Optional field to control the assignment behavior for students who are added after the assignment is published. If not specified, defaults to none. Supported values are: none, assignIfOpen. For example, a teacher can use assignIfOpen to indicate that an assignment should be assigned to any new student who joins the class while the assignment is still open, and none to indicate that an assignment shouldn't be assigned to new students.
// returns a *EducationAddedStudentAction when successful
func (m *EducationAssignment) GetAddedStudentAction()(*EducationAddedStudentAction) {
    val, err := m.GetBackingStore().Get("addedStudentAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationAddedStudentAction)
    }
    return nil
}
// GetAddToCalendarAction gets the addToCalendarAction property value. Optional field to control the assignment behavior  for adding assignments to students' and teachers' calendars when the assignment is published. The possible values are: none, studentsAndPublisher, studentsAndTeamOwners, unknownFutureValue, and studentsOnly. You must use the Prefer: include-unknown-enum-members request header to get the following values in this evolvable enum: studentsOnly. The default value is none.
// returns a *EducationAddToCalendarOptions when successful
func (m *EducationAssignment) GetAddToCalendarAction()(*EducationAddToCalendarOptions) {
    val, err := m.GetBackingStore().Get("addToCalendarAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationAddToCalendarOptions)
    }
    return nil
}
// GetAllowLateSubmissions gets the allowLateSubmissions property value. Identifies whether students can submit after the due date. If this property isn't specified during create, it defaults to true.
// returns a *bool when successful
func (m *EducationAssignment) GetAllowLateSubmissions()(*bool) {
    val, err := m.GetBackingStore().Get("allowLateSubmissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowStudentsToAddResourcesToSubmission gets the allowStudentsToAddResourcesToSubmission property value. Identifies whether students can add their own resources to a submission or if they can only modify resources added by the teacher.
// returns a *bool when successful
func (m *EducationAssignment) GetAllowStudentsToAddResourcesToSubmission()(*bool) {
    val, err := m.GetBackingStore().Get("allowStudentsToAddResourcesToSubmission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAssignDateTime gets the assignDateTime property value. The date when the assignment should become active. If in the future, the assignment isn't shown to the student until this date. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetAssignDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("assignDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAssignedDateTime gets the assignedDateTime property value. The moment that the assignment was published to students and the assignment shows up on the students timeline. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetAssignedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("assignedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAssignTo gets the assignTo property value. Which users, or whole class should receive a submission object once the assignment is published.
// returns a EducationAssignmentRecipientable when successful
func (m *EducationAssignment) GetAssignTo()(EducationAssignmentRecipientable) {
    val, err := m.GetBackingStore().Get("assignTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentRecipientable)
    }
    return nil
}
// GetCategories gets the categories property value. When set, enables users to easily find assignments of a given type. Read-only. Nullable.
// returns a []EducationCategoryable when successful
func (m *EducationAssignment) GetCategories()([]EducationCategoryable) {
    val, err := m.GetBackingStore().Get("categories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationCategoryable)
    }
    return nil
}
// GetClassId gets the classId property value. Class to which this assignment belongs.
// returns a *string when successful
func (m *EducationAssignment) GetClassId()(*string) {
    val, err := m.GetBackingStore().Get("classId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCloseDateTime gets the closeDateTime property value. Date when the assignment is closed for submissions. This is an optional field that can be null if the assignment doesn't allowLateSubmissions or when the closeDateTime is the same as the dueDateTime. But if specified, then the closeDateTime must be greater than or equal to the dueDateTime. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetCloseDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("closeDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Who created the assignment.
// returns a IdentitySetable when successful
func (m *EducationAssignment) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Moment when the assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the assignment.
// returns a *string when successful
func (m *EducationAssignment) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDueDateTime gets the dueDateTime property value. Date when the students assignment is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetDueDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("dueDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFeedbackResourcesFolderUrl gets the feedbackResourcesFolderUrl property value. Folder URL where all the feedback file resources for this assignment are stored.
// returns a *string when successful
func (m *EducationAssignment) GetFeedbackResourcesFolderUrl()(*string) {
    val, err := m.GetBackingStore().Get("feedbackResourcesFolderUrl")
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
func (m *EducationAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["addedStudentAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationAddedStudentAction)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddedStudentAction(val.(*EducationAddedStudentAction))
        }
        return nil
    }
    res["addToCalendarAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationAddToCalendarOptions)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddToCalendarAction(val.(*EducationAddToCalendarOptions))
        }
        return nil
    }
    res["allowLateSubmissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowLateSubmissions(val)
        }
        return nil
    }
    res["allowStudentsToAddResourcesToSubmission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowStudentsToAddResourcesToSubmission(val)
        }
        return nil
    }
    res["assignDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignDateTime(val)
        }
        return nil
    }
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
    res["assignTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignTo(val.(EducationAssignmentRecipientable))
        }
        return nil
    }
    res["categories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationCategoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationCategoryable)
                }
            }
            m.SetCategories(res)
        }
        return nil
    }
    res["classId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassId(val)
        }
        return nil
    }
    res["closeDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloseDateTime(val)
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["dueDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDueDateTime(val)
        }
        return nil
    }
    res["feedbackResourcesFolderUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeedbackResourcesFolderUrl(val)
        }
        return nil
    }
    res["grading"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentGradeTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrading(val.(EducationAssignmentGradeTypeable))
        }
        return nil
    }
    res["gradingCategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationGradingCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGradingCategory(val.(EducationGradingCategoryable))
        }
        return nil
    }
    res["instructions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstructions(val.(EducationItemBodyable))
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["moduleUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModuleUrl(val)
        }
        return nil
    }
    res["notificationChannelUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationChannelUrl(val)
        }
        return nil
    }
    res["resources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationAssignmentResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationAssignmentResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationAssignmentResourceable)
                }
            }
            m.SetResources(res)
        }
        return nil
    }
    res["resourcesFolderUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourcesFolderUrl(val)
        }
        return nil
    }
    res["rubric"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationRubricFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRubric(val.(EducationRubricable))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationAssignmentStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*EducationAssignmentStatus))
        }
        return nil
    }
    res["submissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationSubmissionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationSubmissionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationSubmissionable)
                }
            }
            m.SetSubmissions(res)
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetGrading gets the grading property value. How the assignment will be graded.
// returns a EducationAssignmentGradeTypeable when successful
func (m *EducationAssignment) GetGrading()(EducationAssignmentGradeTypeable) {
    val, err := m.GetBackingStore().Get("grading")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentGradeTypeable)
    }
    return nil
}
// GetGradingCategory gets the gradingCategory property value. When set, enables users to weight assignments differently when computing a class average grade.
// returns a EducationGradingCategoryable when successful
func (m *EducationAssignment) GetGradingCategory()(EducationGradingCategoryable) {
    val, err := m.GetBackingStore().Get("gradingCategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationGradingCategoryable)
    }
    return nil
}
// GetInstructions gets the instructions property value. Instructions for the assignment. The instructions and the display name tell the student what to do.
// returns a EducationItemBodyable when successful
func (m *EducationAssignment) GetInstructions()(EducationItemBodyable) {
    val, err := m.GetBackingStore().Get("instructions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationItemBodyable)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. Who last modified the assignment.
// returns a IdentitySetable when successful
func (m *EducationAssignment) GetLastModifiedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time on which the assignment was modified. A student submission doesn't modify the assignment; only teachers can update assignments. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EducationAssignment) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetModuleUrl gets the moduleUrl property value. The URL of the module from which to access the assignment.
// returns a *string when successful
func (m *EducationAssignment) GetModuleUrl()(*string) {
    val, err := m.GetBackingStore().Get("moduleUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationChannelUrl gets the notificationChannelUrl property value. Optional field to specify the URL of the channel to post the assignment publish notification. If not specified or null, defaults to the General channel. This field only applies to assignments where the assignTo value is educationAssignmentClassRecipient. Updating the notificationChannelUrl isn't allowed after the assignment is published.
// returns a *string when successful
func (m *EducationAssignment) GetNotificationChannelUrl()(*string) {
    val, err := m.GetBackingStore().Get("notificationChannelUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResources gets the resources property value. Learning objects that are associated with this assignment. Only teachers can modify this list. Nullable.
// returns a []EducationAssignmentResourceable when successful
func (m *EducationAssignment) GetResources()([]EducationAssignmentResourceable) {
    val, err := m.GetBackingStore().Get("resources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationAssignmentResourceable)
    }
    return nil
}
// GetResourcesFolderUrl gets the resourcesFolderUrl property value. Folder URL where all the file resources for this assignment are stored.
// returns a *string when successful
func (m *EducationAssignment) GetResourcesFolderUrl()(*string) {
    val, err := m.GetBackingStore().Get("resourcesFolderUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRubric gets the rubric property value. When set, the grading rubric attached to this assignment.
// returns a EducationRubricable when successful
func (m *EducationAssignment) GetRubric()(EducationRubricable) {
    val, err := m.GetBackingStore().Get("rubric")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationRubricable)
    }
    return nil
}
// GetStatus gets the status property value. Status of the assignment.  You can't PATCH this value. Possible values are: draft, scheduled, published, assigned, unknownFutureValue, inactive. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: inactive.
// returns a *EducationAssignmentStatus when successful
func (m *EducationAssignment) GetStatus()(*EducationAssignmentStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationAssignmentStatus)
    }
    return nil
}
// GetSubmissions gets the submissions property value. Once published, there's a submission object for each student representing their work and grade. Read-only. Nullable.
// returns a []EducationSubmissionable when successful
func (m *EducationAssignment) GetSubmissions()([]EducationSubmissionable) {
    val, err := m.GetBackingStore().Get("submissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationSubmissionable)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. The deep link URL for the given assignment.
// returns a *string when successful
func (m *EducationAssignment) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAddedStudentAction() != nil {
        cast := (*m.GetAddedStudentAction()).String()
        err = writer.WriteStringValue("addedStudentAction", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAddToCalendarAction() != nil {
        cast := (*m.GetAddToCalendarAction()).String()
        err = writer.WriteStringValue("addToCalendarAction", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowLateSubmissions", m.GetAllowLateSubmissions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowStudentsToAddResourcesToSubmission", m.GetAllowStudentsToAddResourcesToSubmission())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignTo", m.GetAssignTo())
        if err != nil {
            return err
        }
    }
    if m.GetCategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCategories()))
        for i, v := range m.GetCategories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("categories", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("classId", m.GetClassId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("closeDateTime", m.GetCloseDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("dueDateTime", m.GetDueDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("grading", m.GetGrading())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("gradingCategory", m.GetGradingCategory())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("instructions", m.GetInstructions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("moduleUrl", m.GetModuleUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationChannelUrl", m.GetNotificationChannelUrl())
        if err != nil {
            return err
        }
    }
    if m.GetResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResources()))
        for i, v := range m.GetResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resources", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("rubric", m.GetRubric())
        if err != nil {
            return err
        }
    }
    if m.GetSubmissions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSubmissions()))
        for i, v := range m.GetSubmissions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("submissions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAddedStudentAction sets the addedStudentAction property value. Optional field to control the assignment behavior for students who are added after the assignment is published. If not specified, defaults to none. Supported values are: none, assignIfOpen. For example, a teacher can use assignIfOpen to indicate that an assignment should be assigned to any new student who joins the class while the assignment is still open, and none to indicate that an assignment shouldn't be assigned to new students.
func (m *EducationAssignment) SetAddedStudentAction(value *EducationAddedStudentAction)() {
    err := m.GetBackingStore().Set("addedStudentAction", value)
    if err != nil {
        panic(err)
    }
}
// SetAddToCalendarAction sets the addToCalendarAction property value. Optional field to control the assignment behavior  for adding assignments to students' and teachers' calendars when the assignment is published. The possible values are: none, studentsAndPublisher, studentsAndTeamOwners, unknownFutureValue, and studentsOnly. You must use the Prefer: include-unknown-enum-members request header to get the following values in this evolvable enum: studentsOnly. The default value is none.
func (m *EducationAssignment) SetAddToCalendarAction(value *EducationAddToCalendarOptions)() {
    err := m.GetBackingStore().Set("addToCalendarAction", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowLateSubmissions sets the allowLateSubmissions property value. Identifies whether students can submit after the due date. If this property isn't specified during create, it defaults to true.
func (m *EducationAssignment) SetAllowLateSubmissions(value *bool)() {
    err := m.GetBackingStore().Set("allowLateSubmissions", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowStudentsToAddResourcesToSubmission sets the allowStudentsToAddResourcesToSubmission property value. Identifies whether students can add their own resources to a submission or if they can only modify resources added by the teacher.
func (m *EducationAssignment) SetAllowStudentsToAddResourcesToSubmission(value *bool)() {
    err := m.GetBackingStore().Set("allowStudentsToAddResourcesToSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignDateTime sets the assignDateTime property value. The date when the assignment should become active. If in the future, the assignment isn't shown to the student until this date. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetAssignDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("assignDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedDateTime sets the assignedDateTime property value. The moment that the assignment was published to students and the assignment shows up on the students timeline. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetAssignedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("assignedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignTo sets the assignTo property value. Which users, or whole class should receive a submission object once the assignment is published.
func (m *EducationAssignment) SetAssignTo(value EducationAssignmentRecipientable)() {
    err := m.GetBackingStore().Set("assignTo", value)
    if err != nil {
        panic(err)
    }
}
// SetCategories sets the categories property value. When set, enables users to easily find assignments of a given type. Read-only. Nullable.
func (m *EducationAssignment) SetCategories(value []EducationCategoryable)() {
    err := m.GetBackingStore().Set("categories", value)
    if err != nil {
        panic(err)
    }
}
// SetClassId sets the classId property value. Class to which this assignment belongs.
func (m *EducationAssignment) SetClassId(value *string)() {
    err := m.GetBackingStore().Set("classId", value)
    if err != nil {
        panic(err)
    }
}
// SetCloseDateTime sets the closeDateTime property value. Date when the assignment is closed for submissions. This is an optional field that can be null if the assignment doesn't allowLateSubmissions or when the closeDateTime is the same as the dueDateTime. But if specified, then the closeDateTime must be greater than or equal to the dueDateTime. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetCloseDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("closeDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Who created the assignment.
func (m *EducationAssignment) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Moment when the assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the assignment.
func (m *EducationAssignment) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDueDateTime sets the dueDateTime property value. Date when the students assignment is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetDueDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("dueDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFeedbackResourcesFolderUrl sets the feedbackResourcesFolderUrl property value. Folder URL where all the feedback file resources for this assignment are stored.
func (m *EducationAssignment) SetFeedbackResourcesFolderUrl(value *string)() {
    err := m.GetBackingStore().Set("feedbackResourcesFolderUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetGrading sets the grading property value. How the assignment will be graded.
func (m *EducationAssignment) SetGrading(value EducationAssignmentGradeTypeable)() {
    err := m.GetBackingStore().Set("grading", value)
    if err != nil {
        panic(err)
    }
}
// SetGradingCategory sets the gradingCategory property value. When set, enables users to weight assignments differently when computing a class average grade.
func (m *EducationAssignment) SetGradingCategory(value EducationGradingCategoryable)() {
    err := m.GetBackingStore().Set("gradingCategory", value)
    if err != nil {
        panic(err)
    }
}
// SetInstructions sets the instructions property value. Instructions for the assignment. The instructions and the display name tell the student what to do.
func (m *EducationAssignment) SetInstructions(value EducationItemBodyable)() {
    err := m.GetBackingStore().Set("instructions", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. Who last modified the assignment.
func (m *EducationAssignment) SetLastModifiedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time on which the assignment was modified. A student submission doesn't modify the assignment; only teachers can update assignments. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EducationAssignment) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetModuleUrl sets the moduleUrl property value. The URL of the module from which to access the assignment.
func (m *EducationAssignment) SetModuleUrl(value *string)() {
    err := m.GetBackingStore().Set("moduleUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationChannelUrl sets the notificationChannelUrl property value. Optional field to specify the URL of the channel to post the assignment publish notification. If not specified or null, defaults to the General channel. This field only applies to assignments where the assignTo value is educationAssignmentClassRecipient. Updating the notificationChannelUrl isn't allowed after the assignment is published.
func (m *EducationAssignment) SetNotificationChannelUrl(value *string)() {
    err := m.GetBackingStore().Set("notificationChannelUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetResources sets the resources property value. Learning objects that are associated with this assignment. Only teachers can modify this list. Nullable.
func (m *EducationAssignment) SetResources(value []EducationAssignmentResourceable)() {
    err := m.GetBackingStore().Set("resources", value)
    if err != nil {
        panic(err)
    }
}
// SetResourcesFolderUrl sets the resourcesFolderUrl property value. Folder URL where all the file resources for this assignment are stored.
func (m *EducationAssignment) SetResourcesFolderUrl(value *string)() {
    err := m.GetBackingStore().Set("resourcesFolderUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetRubric sets the rubric property value. When set, the grading rubric attached to this assignment.
func (m *EducationAssignment) SetRubric(value EducationRubricable)() {
    err := m.GetBackingStore().Set("rubric", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Status of the assignment.  You can't PATCH this value. Possible values are: draft, scheduled, published, assigned, unknownFutureValue, inactive. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: inactive.
func (m *EducationAssignment) SetStatus(value *EducationAssignmentStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetSubmissions sets the submissions property value. Once published, there's a submission object for each student representing their work and grade. Read-only. Nullable.
func (m *EducationAssignment) SetSubmissions(value []EducationSubmissionable)() {
    err := m.GetBackingStore().Set("submissions", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. The deep link URL for the given assignment.
func (m *EducationAssignment) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type EducationAssignmentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAddedStudentAction()(*EducationAddedStudentAction)
    GetAddToCalendarAction()(*EducationAddToCalendarOptions)
    GetAllowLateSubmissions()(*bool)
    GetAllowStudentsToAddResourcesToSubmission()(*bool)
    GetAssignDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAssignedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAssignTo()(EducationAssignmentRecipientable)
    GetCategories()([]EducationCategoryable)
    GetClassId()(*string)
    GetCloseDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedBy()(IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDisplayName()(*string)
    GetDueDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFeedbackResourcesFolderUrl()(*string)
    GetGrading()(EducationAssignmentGradeTypeable)
    GetGradingCategory()(EducationGradingCategoryable)
    GetInstructions()(EducationItemBodyable)
    GetLastModifiedBy()(IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetModuleUrl()(*string)
    GetNotificationChannelUrl()(*string)
    GetResources()([]EducationAssignmentResourceable)
    GetResourcesFolderUrl()(*string)
    GetRubric()(EducationRubricable)
    GetStatus()(*EducationAssignmentStatus)
    GetSubmissions()([]EducationSubmissionable)
    GetWebUrl()(*string)
    SetAddedStudentAction(value *EducationAddedStudentAction)()
    SetAddToCalendarAction(value *EducationAddToCalendarOptions)()
    SetAllowLateSubmissions(value *bool)()
    SetAllowStudentsToAddResourcesToSubmission(value *bool)()
    SetAssignDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAssignedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAssignTo(value EducationAssignmentRecipientable)()
    SetCategories(value []EducationCategoryable)()
    SetClassId(value *string)()
    SetCloseDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedBy(value IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDisplayName(value *string)()
    SetDueDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFeedbackResourcesFolderUrl(value *string)()
    SetGrading(value EducationAssignmentGradeTypeable)()
    SetGradingCategory(value EducationGradingCategoryable)()
    SetInstructions(value EducationItemBodyable)()
    SetLastModifiedBy(value IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetModuleUrl(value *string)()
    SetNotificationChannelUrl(value *string)()
    SetResources(value []EducationAssignmentResourceable)()
    SetResourcesFolderUrl(value *string)()
    SetRubric(value EducationRubricable)()
    SetStatus(value *EducationAssignmentStatus)()
    SetSubmissions(value []EducationSubmissionable)()
    SetWebUrl(value *string)()
}

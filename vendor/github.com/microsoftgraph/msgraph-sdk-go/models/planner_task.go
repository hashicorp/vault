package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PlannerTask struct {
    Entity
}
// NewPlannerTask instantiates a new PlannerTask and sets the default values.
func NewPlannerTask()(*PlannerTask) {
    m := &PlannerTask{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerTaskFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlannerTaskFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerTask(), nil
}
// GetActiveChecklistItemCount gets the activeChecklistItemCount property value. Number of checklist items with value set to false, representing incomplete items.
// returns a *int32 when successful
func (m *PlannerTask) GetActiveChecklistItemCount()(*int32) {
    val, err := m.GetBackingStore().Get("activeChecklistItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAppliedCategories gets the appliedCategories property value. The categories to which the task has been applied. See applied Categories for possible values.
// returns a PlannerAppliedCategoriesable when successful
func (m *PlannerTask) GetAppliedCategories()(PlannerAppliedCategoriesable) {
    val, err := m.GetBackingStore().Get("appliedCategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerAppliedCategoriesable)
    }
    return nil
}
// GetAssignedToTaskBoardFormat gets the assignedToTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by assignedTo.
// returns a PlannerAssignedToTaskBoardTaskFormatable when successful
func (m *PlannerTask) GetAssignedToTaskBoardFormat()(PlannerAssignedToTaskBoardTaskFormatable) {
    val, err := m.GetBackingStore().Get("assignedToTaskBoardFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerAssignedToTaskBoardTaskFormatable)
    }
    return nil
}
// GetAssigneePriority gets the assigneePriority property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
// returns a *string when successful
func (m *PlannerTask) GetAssigneePriority()(*string) {
    val, err := m.GetBackingStore().Get("assigneePriority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAssignments gets the assignments property value. The set of assignees the task is assigned to.
// returns a PlannerAssignmentsable when successful
func (m *PlannerTask) GetAssignments()(PlannerAssignmentsable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerAssignmentsable)
    }
    return nil
}
// GetBucketId gets the bucketId property value. Bucket ID to which the task belongs. The bucket needs to be in the plan that the task is in. It's 28 characters long and case-sensitive. Format validation is done on the service.
// returns a *string when successful
func (m *PlannerTask) GetBucketId()(*string) {
    val, err := m.GetBackingStore().Get("bucketId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBucketTaskBoardFormat gets the bucketTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by bucket.
// returns a PlannerBucketTaskBoardTaskFormatable when successful
func (m *PlannerTask) GetBucketTaskBoardFormat()(PlannerBucketTaskBoardTaskFormatable) {
    val, err := m.GetBackingStore().Get("bucketTaskBoardFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerBucketTaskBoardTaskFormatable)
    }
    return nil
}
// GetChecklistItemCount gets the checklistItemCount property value. Number of checklist items that are present on the task.
// returns a *int32 when successful
func (m *PlannerTask) GetChecklistItemCount()(*int32) {
    val, err := m.GetBackingStore().Get("checklistItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCompletedBy gets the completedBy property value. Identity of the user that completed the task.
// returns a IdentitySetable when successful
func (m *PlannerTask) GetCompletedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("completedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCompletedDateTime gets the completedDateTime property value. Read-only. Date and time at which the 'percentComplete' of the task is set to '100'. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *PlannerTask) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetConversationThreadId gets the conversationThreadId property value. Thread ID of the conversation on the task. This is the ID of the conversation thread object created in the group.
// returns a *string when successful
func (m *PlannerTask) GetConversationThreadId()(*string) {
    val, err := m.GetBackingStore().Get("conversationThreadId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Identity of the user that created the task.
// returns a IdentitySetable when successful
func (m *PlannerTask) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Read-only. Date and time at which the task is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *PlannerTask) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDetails gets the details property value. Read-only. Nullable. More details about the task.
// returns a PlannerTaskDetailsable when successful
func (m *PlannerTask) GetDetails()(PlannerTaskDetailsable) {
    val, err := m.GetBackingStore().Get("details")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerTaskDetailsable)
    }
    return nil
}
// GetDueDateTime gets the dueDateTime property value. Date and time at which the task is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *PlannerTask) GetDueDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("dueDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PlannerTask) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activeChecklistItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActiveChecklistItemCount(val)
        }
        return nil
    }
    res["appliedCategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerAppliedCategoriesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppliedCategories(val.(PlannerAppliedCategoriesable))
        }
        return nil
    }
    res["assignedToTaskBoardFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerAssignedToTaskBoardTaskFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedToTaskBoardFormat(val.(PlannerAssignedToTaskBoardTaskFormatable))
        }
        return nil
    }
    res["assigneePriority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssigneePriority(val)
        }
        return nil
    }
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerAssignmentsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignments(val.(PlannerAssignmentsable))
        }
        return nil
    }
    res["bucketId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBucketId(val)
        }
        return nil
    }
    res["bucketTaskBoardFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerBucketTaskBoardTaskFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBucketTaskBoardFormat(val.(PlannerBucketTaskBoardTaskFormatable))
        }
        return nil
    }
    res["checklistItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChecklistItemCount(val)
        }
        return nil
    }
    res["completedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedBy(val.(IdentitySetable))
        }
        return nil
    }
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
    res["conversationThreadId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConversationThreadId(val)
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
    res["details"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerTaskDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetails(val.(PlannerTaskDetailsable))
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
    res["hasDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasDescription(val)
        }
        return nil
    }
    res["orderHint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrderHint(val)
        }
        return nil
    }
    res["percentComplete"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPercentComplete(val)
        }
        return nil
    }
    res["planId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlanId(val)
        }
        return nil
    }
    res["previewType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePlannerPreviewType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviewType(val.(*PlannerPreviewType))
        }
        return nil
    }
    res["priority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPriority(val)
        }
        return nil
    }
    res["progressTaskBoardFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerProgressTaskBoardTaskFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProgressTaskBoardFormat(val.(PlannerProgressTaskBoardTaskFormatable))
        }
        return nil
    }
    res["referenceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReferenceCount(val)
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val)
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    return res
}
// GetHasDescription gets the hasDescription property value. Read-only. Value is true if the details object of the task has a nonempty description and false otherwise.
// returns a *bool when successful
func (m *PlannerTask) GetHasDescription()(*bool) {
    val, err := m.GetBackingStore().Get("hasDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOrderHint gets the orderHint property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
// returns a *string when successful
func (m *PlannerTask) GetOrderHint()(*string) {
    val, err := m.GetBackingStore().Get("orderHint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPercentComplete gets the percentComplete property value. Percentage of task completion. When set to 100, the task is considered completed.
// returns a *int32 when successful
func (m *PlannerTask) GetPercentComplete()(*int32) {
    val, err := m.GetBackingStore().Get("percentComplete")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPlanId gets the planId property value. Plan ID to which the task belongs.
// returns a *string when successful
func (m *PlannerTask) GetPlanId()(*string) {
    val, err := m.GetBackingStore().Get("planId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreviewType gets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference.
// returns a *PlannerPreviewType when successful
func (m *PlannerTask) GetPreviewType()(*PlannerPreviewType) {
    val, err := m.GetBackingStore().Get("previewType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PlannerPreviewType)
    }
    return nil
}
// GetPriority gets the priority property value. Priority of the task. The valid range of values is between 0 and 10, with the increasing value being lower priority (0 has the highest priority and 10 has the lowest priority).  Currently, Planner interprets values 0 and 1 as 'urgent', 2, 3 and 4 as 'important', 5, 6, and 7 as 'medium', and 8, 9, and 10 as 'low'.  Additionally, Planner sets the value 1 for 'urgent', 3 for 'important', 5 for 'medium', and 9 for 'low'.
// returns a *int32 when successful
func (m *PlannerTask) GetPriority()(*int32) {
    val, err := m.GetBackingStore().Get("priority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetProgressTaskBoardFormat gets the progressTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by progress.
// returns a PlannerProgressTaskBoardTaskFormatable when successful
func (m *PlannerTask) GetProgressTaskBoardFormat()(PlannerProgressTaskBoardTaskFormatable) {
    val, err := m.GetBackingStore().Get("progressTaskBoardFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerProgressTaskBoardTaskFormatable)
    }
    return nil
}
// GetReferenceCount gets the referenceCount property value. Number of external references that exist on the task.
// returns a *int32 when successful
func (m *PlannerTask) GetReferenceCount()(*int32) {
    val, err := m.GetBackingStore().Get("referenceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. Date and time at which the task starts. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *PlannerTask) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTitle gets the title property value. Title of the task.
// returns a *string when successful
func (m *PlannerTask) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PlannerTask) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("activeChecklistItemCount", m.GetActiveChecklistItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("appliedCategories", m.GetAppliedCategories())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignedToTaskBoardFormat", m.GetAssignedToTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assigneePriority", m.GetAssigneePriority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignments", m.GetAssignments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bucketId", m.GetBucketId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("bucketTaskBoardFormat", m.GetBucketTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("checklistItemCount", m.GetChecklistItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("completedBy", m.GetCompletedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("conversationThreadId", m.GetConversationThreadId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("details", m.GetDetails())
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
        err = writer.WriteBoolValue("hasDescription", m.GetHasDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("orderHint", m.GetOrderHint())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("percentComplete", m.GetPercentComplete())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("planId", m.GetPlanId())
        if err != nil {
            return err
        }
    }
    if m.GetPreviewType() != nil {
        cast := (*m.GetPreviewType()).String()
        err = writer.WriteStringValue("previewType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("priority", m.GetPriority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("progressTaskBoardFormat", m.GetProgressTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("referenceCount", m.GetReferenceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActiveChecklistItemCount sets the activeChecklistItemCount property value. Number of checklist items with value set to false, representing incomplete items.
func (m *PlannerTask) SetActiveChecklistItemCount(value *int32)() {
    err := m.GetBackingStore().Set("activeChecklistItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAppliedCategories sets the appliedCategories property value. The categories to which the task has been applied. See applied Categories for possible values.
func (m *PlannerTask) SetAppliedCategories(value PlannerAppliedCategoriesable)() {
    err := m.GetBackingStore().Set("appliedCategories", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignedToTaskBoardFormat sets the assignedToTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by assignedTo.
func (m *PlannerTask) SetAssignedToTaskBoardFormat(value PlannerAssignedToTaskBoardTaskFormatable)() {
    err := m.GetBackingStore().Set("assignedToTaskBoardFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetAssigneePriority sets the assigneePriority property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) SetAssigneePriority(value *string)() {
    err := m.GetBackingStore().Set("assigneePriority", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignments sets the assignments property value. The set of assignees the task is assigned to.
func (m *PlannerTask) SetAssignments(value PlannerAssignmentsable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetBucketId sets the bucketId property value. Bucket ID to which the task belongs. The bucket needs to be in the plan that the task is in. It's 28 characters long and case-sensitive. Format validation is done on the service.
func (m *PlannerTask) SetBucketId(value *string)() {
    err := m.GetBackingStore().Set("bucketId", value)
    if err != nil {
        panic(err)
    }
}
// SetBucketTaskBoardFormat sets the bucketTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by bucket.
func (m *PlannerTask) SetBucketTaskBoardFormat(value PlannerBucketTaskBoardTaskFormatable)() {
    err := m.GetBackingStore().Set("bucketTaskBoardFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetChecklistItemCount sets the checklistItemCount property value. Number of checklist items that are present on the task.
func (m *PlannerTask) SetChecklistItemCount(value *int32)() {
    err := m.GetBackingStore().Set("checklistItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedBy sets the completedBy property value. Identity of the user that completed the task.
func (m *PlannerTask) SetCompletedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("completedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedDateTime sets the completedDateTime property value. Read-only. Date and time at which the 'percentComplete' of the task is set to '100'. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetConversationThreadId sets the conversationThreadId property value. Thread ID of the conversation on the task. This is the ID of the conversation thread object created in the group.
func (m *PlannerTask) SetConversationThreadId(value *string)() {
    err := m.GetBackingStore().Set("conversationThreadId", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Identity of the user that created the task.
func (m *PlannerTask) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Read-only. Date and time at which the task is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDetails sets the details property value. Read-only. Nullable. More details about the task.
func (m *PlannerTask) SetDetails(value PlannerTaskDetailsable)() {
    err := m.GetBackingStore().Set("details", value)
    if err != nil {
        panic(err)
    }
}
// SetDueDateTime sets the dueDateTime property value. Date and time at which the task is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetDueDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("dueDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHasDescription sets the hasDescription property value. Read-only. Value is true if the details object of the task has a nonempty description and false otherwise.
func (m *PlannerTask) SetHasDescription(value *bool)() {
    err := m.GetBackingStore().Set("hasDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetOrderHint sets the orderHint property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) SetOrderHint(value *string)() {
    err := m.GetBackingStore().Set("orderHint", value)
    if err != nil {
        panic(err)
    }
}
// SetPercentComplete sets the percentComplete property value. Percentage of task completion. When set to 100, the task is considered completed.
func (m *PlannerTask) SetPercentComplete(value *int32)() {
    err := m.GetBackingStore().Set("percentComplete", value)
    if err != nil {
        panic(err)
    }
}
// SetPlanId sets the planId property value. Plan ID to which the task belongs.
func (m *PlannerTask) SetPlanId(value *string)() {
    err := m.GetBackingStore().Set("planId", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviewType sets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference.
func (m *PlannerTask) SetPreviewType(value *PlannerPreviewType)() {
    err := m.GetBackingStore().Set("previewType", value)
    if err != nil {
        panic(err)
    }
}
// SetPriority sets the priority property value. Priority of the task. The valid range of values is between 0 and 10, with the increasing value being lower priority (0 has the highest priority and 10 has the lowest priority).  Currently, Planner interprets values 0 and 1 as 'urgent', 2, 3 and 4 as 'important', 5, 6, and 7 as 'medium', and 8, 9, and 10 as 'low'.  Additionally, Planner sets the value 1 for 'urgent', 3 for 'important', 5 for 'medium', and 9 for 'low'.
func (m *PlannerTask) SetPriority(value *int32)() {
    err := m.GetBackingStore().Set("priority", value)
    if err != nil {
        panic(err)
    }
}
// SetProgressTaskBoardFormat sets the progressTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by progress.
func (m *PlannerTask) SetProgressTaskBoardFormat(value PlannerProgressTaskBoardTaskFormatable)() {
    err := m.GetBackingStore().Set("progressTaskBoardFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetReferenceCount sets the referenceCount property value. Number of external references that exist on the task.
func (m *PlannerTask) SetReferenceCount(value *int32)() {
    err := m.GetBackingStore().Set("referenceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. Date and time at which the task starts. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. Title of the task.
func (m *PlannerTask) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
type PlannerTaskable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActiveChecklistItemCount()(*int32)
    GetAppliedCategories()(PlannerAppliedCategoriesable)
    GetAssignedToTaskBoardFormat()(PlannerAssignedToTaskBoardTaskFormatable)
    GetAssigneePriority()(*string)
    GetAssignments()(PlannerAssignmentsable)
    GetBucketId()(*string)
    GetBucketTaskBoardFormat()(PlannerBucketTaskBoardTaskFormatable)
    GetChecklistItemCount()(*int32)
    GetCompletedBy()(IdentitySetable)
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetConversationThreadId()(*string)
    GetCreatedBy()(IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDetails()(PlannerTaskDetailsable)
    GetDueDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHasDescription()(*bool)
    GetOrderHint()(*string)
    GetPercentComplete()(*int32)
    GetPlanId()(*string)
    GetPreviewType()(*PlannerPreviewType)
    GetPriority()(*int32)
    GetProgressTaskBoardFormat()(PlannerProgressTaskBoardTaskFormatable)
    GetReferenceCount()(*int32)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTitle()(*string)
    SetActiveChecklistItemCount(value *int32)()
    SetAppliedCategories(value PlannerAppliedCategoriesable)()
    SetAssignedToTaskBoardFormat(value PlannerAssignedToTaskBoardTaskFormatable)()
    SetAssigneePriority(value *string)()
    SetAssignments(value PlannerAssignmentsable)()
    SetBucketId(value *string)()
    SetBucketTaskBoardFormat(value PlannerBucketTaskBoardTaskFormatable)()
    SetChecklistItemCount(value *int32)()
    SetCompletedBy(value IdentitySetable)()
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetConversationThreadId(value *string)()
    SetCreatedBy(value IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDetails(value PlannerTaskDetailsable)()
    SetDueDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHasDescription(value *bool)()
    SetOrderHint(value *string)()
    SetPercentComplete(value *int32)()
    SetPlanId(value *string)()
    SetPreviewType(value *PlannerPreviewType)()
    SetPriority(value *int32)()
    SetProgressTaskBoardFormat(value PlannerProgressTaskBoardTaskFormatable)()
    SetReferenceCount(value *int32)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTitle(value *string)()
}

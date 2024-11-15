package identitygovernance

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Run struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewRun instantiates a new Run and sets the default values.
func NewRun()(*Run) {
    m := &Run{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateRunFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRunFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRun(), nil
}
// GetCompletedDateTime gets the completedDateTime property value. The date time that the run completed. Value is null if the workflow hasn't completed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *Run) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFailedTasksCount gets the failedTasksCount property value. The number of tasks that failed in the run execution.
// returns a *int32 when successful
func (m *Run) GetFailedTasksCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedTasksCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedUsersCount gets the failedUsersCount property value. The number of users that failed in the run execution.
// returns a *int32 when successful
func (m *Run) GetFailedUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Run) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["failedTasksCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedTasksCount(val)
        }
        return nil
    }
    res["failedUsersCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedUsersCount(val)
        }
        return nil
    }
    res["lastUpdatedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdatedDateTime(val)
        }
        return nil
    }
    res["processingStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLifecycleWorkflowProcessingStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessingStatus(val.(*LifecycleWorkflowProcessingStatus))
        }
        return nil
    }
    res["scheduledDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduledDateTime(val)
        }
        return nil
    }
    res["startedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartedDateTime(val)
        }
        return nil
    }
    res["successfulUsersCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulUsersCount(val)
        }
        return nil
    }
    res["taskProcessingResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTaskProcessingResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TaskProcessingResultable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TaskProcessingResultable)
                }
            }
            m.SetTaskProcessingResults(res)
        }
        return nil
    }
    res["totalTasksCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalTasksCount(val)
        }
        return nil
    }
    res["totalUnprocessedTasksCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalUnprocessedTasksCount(val)
        }
        return nil
    }
    res["totalUsersCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalUsersCount(val)
        }
        return nil
    }
    res["userProcessingResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserProcessingResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserProcessingResultable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserProcessingResultable)
                }
            }
            m.SetUserProcessingResults(res)
        }
        return nil
    }
    res["workflowExecutionType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWorkflowExecutionType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflowExecutionType(val.(*WorkflowExecutionType))
        }
        return nil
    }
    return res
}
// GetLastUpdatedDateTime gets the lastUpdatedDateTime property value. The datetime that the run was last updated.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *Run) GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdatedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetProcessingStatus gets the processingStatus property value. The processingStatus property
// returns a *LifecycleWorkflowProcessingStatus when successful
func (m *Run) GetProcessingStatus()(*LifecycleWorkflowProcessingStatus) {
    val, err := m.GetBackingStore().Get("processingStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*LifecycleWorkflowProcessingStatus)
    }
    return nil
}
// GetScheduledDateTime gets the scheduledDateTime property value. The date time that the run is scheduled to be executed for a workflow.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *Run) GetScheduledDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("scheduledDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetStartedDateTime gets the startedDateTime property value. The date time that the run execution started.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *Run) GetStartedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSuccessfulUsersCount gets the successfulUsersCount property value. The number of successfully completed users in the run.
// returns a *int32 when successful
func (m *Run) GetSuccessfulUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("successfulUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTaskProcessingResults gets the taskProcessingResults property value. The related taskProcessingResults.
// returns a []TaskProcessingResultable when successful
func (m *Run) GetTaskProcessingResults()([]TaskProcessingResultable) {
    val, err := m.GetBackingStore().Get("taskProcessingResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TaskProcessingResultable)
    }
    return nil
}
// GetTotalTasksCount gets the totalTasksCount property value. The totalTasksCount property
// returns a *int32 when successful
func (m *Run) GetTotalTasksCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalTasksCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalUnprocessedTasksCount gets the totalUnprocessedTasksCount property value. The total number of unprocessed tasks in the run execution.
// returns a *int32 when successful
func (m *Run) GetTotalUnprocessedTasksCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalUnprocessedTasksCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalUsersCount gets the totalUsersCount property value. The total number of users in the workflow execution.
// returns a *int32 when successful
func (m *Run) GetTotalUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUserProcessingResults gets the userProcessingResults property value. The associated individual user execution.
// returns a []UserProcessingResultable when successful
func (m *Run) GetUserProcessingResults()([]UserProcessingResultable) {
    val, err := m.GetBackingStore().Get("userProcessingResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserProcessingResultable)
    }
    return nil
}
// GetWorkflowExecutionType gets the workflowExecutionType property value. The workflowExecutionType property
// returns a *WorkflowExecutionType when successful
func (m *Run) GetWorkflowExecutionType()(*WorkflowExecutionType) {
    val, err := m.GetBackingStore().Get("workflowExecutionType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WorkflowExecutionType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Run) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteInt32Value("failedTasksCount", m.GetFailedTasksCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("failedUsersCount", m.GetFailedUsersCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdatedDateTime", m.GetLastUpdatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetProcessingStatus() != nil {
        cast := (*m.GetProcessingStatus()).String()
        err = writer.WriteStringValue("processingStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("scheduledDateTime", m.GetScheduledDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startedDateTime", m.GetStartedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("successfulUsersCount", m.GetSuccessfulUsersCount())
        if err != nil {
            return err
        }
    }
    if m.GetTaskProcessingResults() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTaskProcessingResults()))
        for i, v := range m.GetTaskProcessingResults() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("taskProcessingResults", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalTasksCount", m.GetTotalTasksCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalUnprocessedTasksCount", m.GetTotalUnprocessedTasksCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalUsersCount", m.GetTotalUsersCount())
        if err != nil {
            return err
        }
    }
    if m.GetUserProcessingResults() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserProcessingResults()))
        for i, v := range m.GetUserProcessingResults() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userProcessingResults", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorkflowExecutionType() != nil {
        cast := (*m.GetWorkflowExecutionType()).String()
        err = writer.WriteStringValue("workflowExecutionType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletedDateTime sets the completedDateTime property value. The date time that the run completed. Value is null if the workflow hasn't completed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Run) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedTasksCount sets the failedTasksCount property value. The number of tasks that failed in the run execution.
func (m *Run) SetFailedTasksCount(value *int32)() {
    err := m.GetBackingStore().Set("failedTasksCount", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedUsersCount sets the failedUsersCount property value. The number of users that failed in the run execution.
func (m *Run) SetFailedUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("failedUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdatedDateTime sets the lastUpdatedDateTime property value. The datetime that the run was last updated.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Run) SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdatedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessingStatus sets the processingStatus property value. The processingStatus property
func (m *Run) SetProcessingStatus(value *LifecycleWorkflowProcessingStatus)() {
    err := m.GetBackingStore().Set("processingStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduledDateTime sets the scheduledDateTime property value. The date time that the run is scheduled to be executed for a workflow.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Run) SetScheduledDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("scheduledDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetStartedDateTime sets the startedDateTime property value. The date time that the run execution started.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Run) SetStartedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulUsersCount sets the successfulUsersCount property value. The number of successfully completed users in the run.
func (m *Run) SetSuccessfulUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("successfulUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskProcessingResults sets the taskProcessingResults property value. The related taskProcessingResults.
func (m *Run) SetTaskProcessingResults(value []TaskProcessingResultable)() {
    err := m.GetBackingStore().Set("taskProcessingResults", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalTasksCount sets the totalTasksCount property value. The totalTasksCount property
func (m *Run) SetTotalTasksCount(value *int32)() {
    err := m.GetBackingStore().Set("totalTasksCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUnprocessedTasksCount sets the totalUnprocessedTasksCount property value. The total number of unprocessed tasks in the run execution.
func (m *Run) SetTotalUnprocessedTasksCount(value *int32)() {
    err := m.GetBackingStore().Set("totalUnprocessedTasksCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUsersCount sets the totalUsersCount property value. The total number of users in the workflow execution.
func (m *Run) SetTotalUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("totalUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUserProcessingResults sets the userProcessingResults property value. The associated individual user execution.
func (m *Run) SetUserProcessingResults(value []UserProcessingResultable)() {
    err := m.GetBackingStore().Set("userProcessingResults", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowExecutionType sets the workflowExecutionType property value. The workflowExecutionType property
func (m *Run) SetWorkflowExecutionType(value *WorkflowExecutionType)() {
    err := m.GetBackingStore().Set("workflowExecutionType", value)
    if err != nil {
        panic(err)
    }
}
type Runable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFailedTasksCount()(*int32)
    GetFailedUsersCount()(*int32)
    GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetProcessingStatus()(*LifecycleWorkflowProcessingStatus)
    GetScheduledDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetStartedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSuccessfulUsersCount()(*int32)
    GetTaskProcessingResults()([]TaskProcessingResultable)
    GetTotalTasksCount()(*int32)
    GetTotalUnprocessedTasksCount()(*int32)
    GetTotalUsersCount()(*int32)
    GetUserProcessingResults()([]UserProcessingResultable)
    GetWorkflowExecutionType()(*WorkflowExecutionType)
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFailedTasksCount(value *int32)()
    SetFailedUsersCount(value *int32)()
    SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetProcessingStatus(value *LifecycleWorkflowProcessingStatus)()
    SetScheduledDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetStartedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSuccessfulUsersCount(value *int32)()
    SetTaskProcessingResults(value []TaskProcessingResultable)()
    SetTotalTasksCount(value *int32)()
    SetTotalUnprocessedTasksCount(value *int32)()
    SetTotalUsersCount(value *int32)()
    SetUserProcessingResults(value []UserProcessingResultable)()
    SetWorkflowExecutionType(value *WorkflowExecutionType)()
}

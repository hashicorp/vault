package identitygovernance

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type TaskReport struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewTaskReport instantiates a new TaskReport and sets the default values.
func NewTaskReport()(*TaskReport) {
    m := &TaskReport{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateTaskReportFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTaskReportFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTaskReport(), nil
}
// GetCompletedDateTime gets the completedDateTime property value. The date time that the associated run completed. Value is null if the run has not completed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *TaskReport) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFailedUsersCount gets the failedUsersCount property value. The number of users in the run execution for which the associated task failed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *TaskReport) GetFailedUsersCount()(*int32) {
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
func (m *TaskReport) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["runId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRunId(val)
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
    res["task"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTaskFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTask(val.(Taskable))
        }
        return nil
    }
    res["taskDefinition"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTaskDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTaskDefinition(val.(TaskDefinitionable))
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
    res["unprocessedUsersCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnprocessedUsersCount(val)
        }
        return nil
    }
    return res
}
// GetLastUpdatedDateTime gets the lastUpdatedDateTime property value. The date and time that the task report was last updated.
// returns a *Time when successful
func (m *TaskReport) GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *TaskReport) GetProcessingStatus()(*LifecycleWorkflowProcessingStatus) {
    val, err := m.GetBackingStore().Get("processingStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*LifecycleWorkflowProcessingStatus)
    }
    return nil
}
// GetRunId gets the runId property value. The unique identifier of the associated run.
// returns a *string when successful
func (m *TaskReport) GetRunId()(*string) {
    val, err := m.GetBackingStore().Get("runId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartedDateTime gets the startedDateTime property value. The date time that the associated run started. Value is null if the run has not started.
// returns a *Time when successful
func (m *TaskReport) GetStartedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSuccessfulUsersCount gets the successfulUsersCount property value. The number of users in the run execution for which the associated task succeeded.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *TaskReport) GetSuccessfulUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("successfulUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTask gets the task property value. The task property
// returns a Taskable when successful
func (m *TaskReport) GetTask()(Taskable) {
    val, err := m.GetBackingStore().Get("task")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Taskable)
    }
    return nil
}
// GetTaskDefinition gets the taskDefinition property value. The taskDefinition property
// returns a TaskDefinitionable when successful
func (m *TaskReport) GetTaskDefinition()(TaskDefinitionable) {
    val, err := m.GetBackingStore().Get("taskDefinition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TaskDefinitionable)
    }
    return nil
}
// GetTaskProcessingResults gets the taskProcessingResults property value. The related lifecycle workflow taskProcessingResults.
// returns a []TaskProcessingResultable when successful
func (m *TaskReport) GetTaskProcessingResults()([]TaskProcessingResultable) {
    val, err := m.GetBackingStore().Get("taskProcessingResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TaskProcessingResultable)
    }
    return nil
}
// GetTotalUsersCount gets the totalUsersCount property value. The total number of users in the run execution for which the associated task was scheduled to execute.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *TaskReport) GetTotalUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnprocessedUsersCount gets the unprocessedUsersCount property value. The number of users in the run execution for which the associated task is queued, in progress, or canceled.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *TaskReport) GetUnprocessedUsersCount()(*int32) {
    val, err := m.GetBackingStore().Get("unprocessedUsersCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TaskReport) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("runId", m.GetRunId())
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
    {
        err = writer.WriteObjectValue("task", m.GetTask())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("taskDefinition", m.GetTaskDefinition())
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
        err = writer.WriteInt32Value("totalUsersCount", m.GetTotalUsersCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unprocessedUsersCount", m.GetUnprocessedUsersCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletedDateTime sets the completedDateTime property value. The date time that the associated run completed. Value is null if the run has not completed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *TaskReport) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedUsersCount sets the failedUsersCount property value. The number of users in the run execution for which the associated task failed.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *TaskReport) SetFailedUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("failedUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdatedDateTime sets the lastUpdatedDateTime property value. The date and time that the task report was last updated.
func (m *TaskReport) SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdatedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessingStatus sets the processingStatus property value. The processingStatus property
func (m *TaskReport) SetProcessingStatus(value *LifecycleWorkflowProcessingStatus)() {
    err := m.GetBackingStore().Set("processingStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetRunId sets the runId property value. The unique identifier of the associated run.
func (m *TaskReport) SetRunId(value *string)() {
    err := m.GetBackingStore().Set("runId", value)
    if err != nil {
        panic(err)
    }
}
// SetStartedDateTime sets the startedDateTime property value. The date time that the associated run started. Value is null if the run has not started.
func (m *TaskReport) SetStartedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulUsersCount sets the successfulUsersCount property value. The number of users in the run execution for which the associated task succeeded.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *TaskReport) SetSuccessfulUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("successfulUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTask sets the task property value. The task property
func (m *TaskReport) SetTask(value Taskable)() {
    err := m.GetBackingStore().Set("task", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskDefinition sets the taskDefinition property value. The taskDefinition property
func (m *TaskReport) SetTaskDefinition(value TaskDefinitionable)() {
    err := m.GetBackingStore().Set("taskDefinition", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskProcessingResults sets the taskProcessingResults property value. The related lifecycle workflow taskProcessingResults.
func (m *TaskReport) SetTaskProcessingResults(value []TaskProcessingResultable)() {
    err := m.GetBackingStore().Set("taskProcessingResults", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUsersCount sets the totalUsersCount property value. The total number of users in the run execution for which the associated task was scheduled to execute.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *TaskReport) SetTotalUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("totalUsersCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnprocessedUsersCount sets the unprocessedUsersCount property value. The number of users in the run execution for which the associated task is queued, in progress, or canceled.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *TaskReport) SetUnprocessedUsersCount(value *int32)() {
    err := m.GetBackingStore().Set("unprocessedUsersCount", value)
    if err != nil {
        panic(err)
    }
}
type TaskReportable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFailedUsersCount()(*int32)
    GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetProcessingStatus()(*LifecycleWorkflowProcessingStatus)
    GetRunId()(*string)
    GetStartedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSuccessfulUsersCount()(*int32)
    GetTask()(Taskable)
    GetTaskDefinition()(TaskDefinitionable)
    GetTaskProcessingResults()([]TaskProcessingResultable)
    GetTotalUsersCount()(*int32)
    GetUnprocessedUsersCount()(*int32)
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFailedUsersCount(value *int32)()
    SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetProcessingStatus(value *LifecycleWorkflowProcessingStatus)()
    SetRunId(value *string)()
    SetStartedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSuccessfulUsersCount(value *int32)()
    SetTask(value Taskable)()
    SetTaskDefinition(value TaskDefinitionable)()
    SetTaskProcessingResults(value []TaskProcessingResultable)()
    SetTotalUsersCount(value *int32)()
    SetUnprocessedUsersCount(value *int32)()
}

package identitygovernance

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Workflow struct {
    WorkflowBase
}
// NewWorkflow instantiates a new Workflow and sets the default values.
func NewWorkflow()(*Workflow) {
    m := &Workflow{
        WorkflowBase: *NewWorkflowBase(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.workflow"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWorkflowFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkflowFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkflow(), nil
}
// GetDeletedDateTime gets the deletedDateTime property value. When the workflow was deleted.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *Time when successful
func (m *Workflow) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("deletedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetExecutionScope gets the executionScope property value. The unique identifier of the Microsoft Entra identity that last modified the workflow object.
// returns a []UserProcessingResultable when successful
func (m *Workflow) GetExecutionScope()([]UserProcessingResultable) {
    val, err := m.GetBackingStore().Get("executionScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserProcessingResultable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Workflow) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WorkflowBase.GetFieldDeserializers()
    res["deletedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedDateTime(val)
        }
        return nil
    }
    res["executionScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetExecutionScope(res)
        }
        return nil
    }
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["nextScheduleRunDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNextScheduleRunDateTime(val)
        }
        return nil
    }
    res["runs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRunFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Runable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Runable)
                }
            }
            m.SetRuns(res)
        }
        return nil
    }
    res["taskReports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTaskReportFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TaskReportable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TaskReportable)
                }
            }
            m.SetTaskReports(res)
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
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    res["versions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkflowVersionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkflowVersionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkflowVersionable)
                }
            }
            m.SetVersions(res)
        }
        return nil
    }
    return res
}
// GetId gets the id property value. Identifier used for individually addressing a specific workflow.Supports $filter(eq, ne) and $orderby.
// returns a *string when successful
func (m *Workflow) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNextScheduleRunDateTime gets the nextScheduleRunDateTime property value. The date time when the workflow is expected to run next based on the schedule interval, if there are any users matching the execution conditions. Supports $filter(lt,gt) and $orderby.
// returns a *Time when successful
func (m *Workflow) GetNextScheduleRunDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("nextScheduleRunDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRuns gets the runs property value. Workflow runs.
// returns a []Runable when successful
func (m *Workflow) GetRuns()([]Runable) {
    val, err := m.GetBackingStore().Get("runs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Runable)
    }
    return nil
}
// GetTaskReports gets the taskReports property value. Represents the aggregation of task execution data for tasks within a workflow object.
// returns a []TaskReportable when successful
func (m *Workflow) GetTaskReports()([]TaskReportable) {
    val, err := m.GetBackingStore().Get("taskReports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TaskReportable)
    }
    return nil
}
// GetUserProcessingResults gets the userProcessingResults property value. Per-user workflow execution results.
// returns a []UserProcessingResultable when successful
func (m *Workflow) GetUserProcessingResults()([]UserProcessingResultable) {
    val, err := m.GetBackingStore().Get("userProcessingResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserProcessingResultable)
    }
    return nil
}
// GetVersion gets the version property value. The current version number of the workflow. Value is 1 when the workflow is first created.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *Workflow) GetVersion()(*int32) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetVersions gets the versions property value. The workflow versions that are available.
// returns a []WorkflowVersionable when successful
func (m *Workflow) GetVersions()([]WorkflowVersionable) {
    val, err := m.GetBackingStore().Get("versions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkflowVersionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Workflow) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WorkflowBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetExecutionScope() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExecutionScope()))
        for i, v := range m.GetExecutionScope() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("executionScope", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("nextScheduleRunDateTime", m.GetNextScheduleRunDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetRuns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRuns()))
        for i, v := range m.GetRuns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("runs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTaskReports() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTaskReports()))
        for i, v := range m.GetTaskReports() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("taskReports", cast)
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
    {
        err = writer.WriteInt32Value("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    if m.GetVersions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetVersions()))
        for i, v := range m.GetVersions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("versions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeletedDateTime sets the deletedDateTime property value. When the workflow was deleted.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Workflow) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("deletedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetExecutionScope sets the executionScope property value. The unique identifier of the Microsoft Entra identity that last modified the workflow object.
func (m *Workflow) SetExecutionScope(value []UserProcessingResultable)() {
    err := m.GetBackingStore().Set("executionScope", value)
    if err != nil {
        panic(err)
    }
}
// SetId sets the id property value. Identifier used for individually addressing a specific workflow.Supports $filter(eq, ne) and $orderby.
func (m *Workflow) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetNextScheduleRunDateTime sets the nextScheduleRunDateTime property value. The date time when the workflow is expected to run next based on the schedule interval, if there are any users matching the execution conditions. Supports $filter(lt,gt) and $orderby.
func (m *Workflow) SetNextScheduleRunDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("nextScheduleRunDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRuns sets the runs property value. Workflow runs.
func (m *Workflow) SetRuns(value []Runable)() {
    err := m.GetBackingStore().Set("runs", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskReports sets the taskReports property value. Represents the aggregation of task execution data for tasks within a workflow object.
func (m *Workflow) SetTaskReports(value []TaskReportable)() {
    err := m.GetBackingStore().Set("taskReports", value)
    if err != nil {
        panic(err)
    }
}
// SetUserProcessingResults sets the userProcessingResults property value. Per-user workflow execution results.
func (m *Workflow) SetUserProcessingResults(value []UserProcessingResultable)() {
    err := m.GetBackingStore().Set("userProcessingResults", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The current version number of the workflow. Value is 1 when the workflow is first created.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *Workflow) SetVersion(value *int32)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
// SetVersions sets the versions property value. The workflow versions that are available.
func (m *Workflow) SetVersions(value []WorkflowVersionable)() {
    err := m.GetBackingStore().Set("versions", value)
    if err != nil {
        panic(err)
    }
}
type Workflowable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WorkflowBaseable
    GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetExecutionScope()([]UserProcessingResultable)
    GetId()(*string)
    GetNextScheduleRunDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRuns()([]Runable)
    GetTaskReports()([]TaskReportable)
    GetUserProcessingResults()([]UserProcessingResultable)
    GetVersion()(*int32)
    GetVersions()([]WorkflowVersionable)
    SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetExecutionScope(value []UserProcessingResultable)()
    SetId(value *string)()
    SetNextScheduleRunDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRuns(value []Runable)()
    SetTaskReports(value []TaskReportable)()
    SetUserProcessingResults(value []UserProcessingResultable)()
    SetVersion(value *int32)()
    SetVersions(value []WorkflowVersionable)()
}

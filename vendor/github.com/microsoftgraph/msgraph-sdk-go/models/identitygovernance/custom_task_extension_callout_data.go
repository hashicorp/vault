package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type CustomTaskExtensionCalloutData struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionData
}
// NewCustomTaskExtensionCalloutData instantiates a new CustomTaskExtensionCalloutData and sets the default values.
func NewCustomTaskExtensionCalloutData()(*CustomTaskExtensionCalloutData) {
    m := &CustomTaskExtensionCalloutData{
        CustomExtensionData: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewCustomExtensionData(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.customTaskExtensionCalloutData"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCustomTaskExtensionCalloutDataFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomTaskExtensionCalloutDataFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomTaskExtensionCalloutData(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomTaskExtensionCalloutData) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomExtensionData.GetFieldDeserializers()
    res["subject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubject(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable))
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
    res["taskProcessingresult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTaskProcessingResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTaskProcessingresult(val.(TaskProcessingResultable))
        }
        return nil
    }
    res["workflow"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkflowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflow(val.(Workflowable))
        }
        return nil
    }
    return res
}
// GetSubject gets the subject property value. The subject property
// returns a Userable when successful
func (m *CustomTaskExtensionCalloutData) GetSubject()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable)
    }
    return nil
}
// GetTask gets the task property value. The task property
// returns a Taskable when successful
func (m *CustomTaskExtensionCalloutData) GetTask()(Taskable) {
    val, err := m.GetBackingStore().Get("task")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Taskable)
    }
    return nil
}
// GetTaskProcessingresult gets the taskProcessingresult property value. The taskProcessingresult property
// returns a TaskProcessingResultable when successful
func (m *CustomTaskExtensionCalloutData) GetTaskProcessingresult()(TaskProcessingResultable) {
    val, err := m.GetBackingStore().Get("taskProcessingresult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TaskProcessingResultable)
    }
    return nil
}
// GetWorkflow gets the workflow property value. The workflow property
// returns a Workflowable when successful
func (m *CustomTaskExtensionCalloutData) GetWorkflow()(Workflowable) {
    val, err := m.GetBackingStore().Get("workflow")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Workflowable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CustomTaskExtensionCalloutData) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomExtensionData.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("subject", m.GetSubject())
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
        err = writer.WriteObjectValue("taskProcessingresult", m.GetTaskProcessingresult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("workflow", m.GetWorkflow())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSubject sets the subject property value. The subject property
func (m *CustomTaskExtensionCalloutData) SetSubject(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetTask sets the task property value. The task property
func (m *CustomTaskExtensionCalloutData) SetTask(value Taskable)() {
    err := m.GetBackingStore().Set("task", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskProcessingresult sets the taskProcessingresult property value. The taskProcessingresult property
func (m *CustomTaskExtensionCalloutData) SetTaskProcessingresult(value TaskProcessingResultable)() {
    err := m.GetBackingStore().Set("taskProcessingresult", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflow sets the workflow property value. The workflow property
func (m *CustomTaskExtensionCalloutData) SetWorkflow(value Workflowable)() {
    err := m.GetBackingStore().Set("workflow", value)
    if err != nil {
        panic(err)
    }
}
type CustomTaskExtensionCalloutDataable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionDataable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSubject()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable)
    GetTask()(Taskable)
    GetTaskProcessingresult()(TaskProcessingResultable)
    GetWorkflow()(Workflowable)
    SetSubject(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable)()
    SetTask(value Taskable)()
    SetTaskProcessingresult(value TaskProcessingResultable)()
    SetWorkflow(value Workflowable)()
}

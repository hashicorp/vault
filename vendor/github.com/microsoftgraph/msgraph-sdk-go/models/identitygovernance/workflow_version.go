package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkflowVersion struct {
    WorkflowBase
}
// NewWorkflowVersion instantiates a new WorkflowVersion and sets the default values.
func NewWorkflowVersion()(*WorkflowVersion) {
    m := &WorkflowVersion{
        WorkflowBase: *NewWorkflowBase(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.workflowVersion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWorkflowVersionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkflowVersionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkflowVersion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkflowVersion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WorkflowBase.GetFieldDeserializers()
    res["versionNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersionNumber(val)
        }
        return nil
    }
    return res
}
// GetVersionNumber gets the versionNumber property value. The version of the workflow.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
// returns a *int32 when successful
func (m *WorkflowVersion) GetVersionNumber()(*int32) {
    val, err := m.GetBackingStore().Get("versionNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkflowVersion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WorkflowBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("versionNumber", m.GetVersionNumber())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetVersionNumber sets the versionNumber property value. The version of the workflow.Supports $filter(lt, le, gt, ge, eq, ne) and $orderby.
func (m *WorkflowVersion) SetVersionNumber(value *int32)() {
    err := m.GetBackingStore().Set("versionNumber", value)
    if err != nil {
        panic(err)
    }
}
type WorkflowVersionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WorkflowBaseable
    GetVersionNumber()(*int32)
    SetVersionNumber(value *int32)()
}

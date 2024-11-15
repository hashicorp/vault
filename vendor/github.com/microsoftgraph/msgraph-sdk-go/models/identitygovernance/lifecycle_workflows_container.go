package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type LifecycleWorkflowsContainer struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewLifecycleWorkflowsContainer instantiates a new LifecycleWorkflowsContainer and sets the default values.
func NewLifecycleWorkflowsContainer()(*LifecycleWorkflowsContainer) {
    m := &LifecycleWorkflowsContainer{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateLifecycleWorkflowsContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLifecycleWorkflowsContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLifecycleWorkflowsContainer(), nil
}
// GetCustomTaskExtensions gets the customTaskExtensions property value. The customTaskExtension instance.
// returns a []CustomTaskExtensionable when successful
func (m *LifecycleWorkflowsContainer) GetCustomTaskExtensions()([]CustomTaskExtensionable) {
    val, err := m.GetBackingStore().Get("customTaskExtensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CustomTaskExtensionable)
    }
    return nil
}
// GetDeletedItems gets the deletedItems property value. Deleted workflows in your lifecycle workflows instance.
// returns a DeletedItemContainerable when successful
func (m *LifecycleWorkflowsContainer) GetDeletedItems()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable) {
    val, err := m.GetBackingStore().Get("deletedItems")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LifecycleWorkflowsContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["customTaskExtensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCustomTaskExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CustomTaskExtensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CustomTaskExtensionable)
                }
            }
            m.SetCustomTaskExtensions(res)
        }
        return nil
    }
    res["deletedItems"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeletedItemContainerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedItems(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable))
        }
        return nil
    }
    res["insights"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInsightsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInsights(val.(Insightsable))
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLifecycleManagementSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(LifecycleManagementSettingsable))
        }
        return nil
    }
    res["taskDefinitions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTaskDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TaskDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TaskDefinitionable)
                }
            }
            m.SetTaskDefinitions(res)
        }
        return nil
    }
    res["workflows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkflowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Workflowable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Workflowable)
                }
            }
            m.SetWorkflows(res)
        }
        return nil
    }
    res["workflowTemplates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkflowTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkflowTemplateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkflowTemplateable)
                }
            }
            m.SetWorkflowTemplates(res)
        }
        return nil
    }
    return res
}
// GetInsights gets the insights property value. The insight container holding workflow insight summaries for a tenant.
// returns a Insightsable when successful
func (m *LifecycleWorkflowsContainer) GetInsights()(Insightsable) {
    val, err := m.GetBackingStore().Get("insights")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Insightsable)
    }
    return nil
}
// GetSettings gets the settings property value. The settings property
// returns a LifecycleManagementSettingsable when successful
func (m *LifecycleWorkflowsContainer) GetSettings()(LifecycleManagementSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LifecycleManagementSettingsable)
    }
    return nil
}
// GetTaskDefinitions gets the taskDefinitions property value. The definition of tasks within the lifecycle workflows instance.
// returns a []TaskDefinitionable when successful
func (m *LifecycleWorkflowsContainer) GetTaskDefinitions()([]TaskDefinitionable) {
    val, err := m.GetBackingStore().Get("taskDefinitions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TaskDefinitionable)
    }
    return nil
}
// GetWorkflows gets the workflows property value. The workflows in the lifecycle workflows instance.
// returns a []Workflowable when successful
func (m *LifecycleWorkflowsContainer) GetWorkflows()([]Workflowable) {
    val, err := m.GetBackingStore().Get("workflows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Workflowable)
    }
    return nil
}
// GetWorkflowTemplates gets the workflowTemplates property value. The workflow templates in the lifecycle workflow instance.
// returns a []WorkflowTemplateable when successful
func (m *LifecycleWorkflowsContainer) GetWorkflowTemplates()([]WorkflowTemplateable) {
    val, err := m.GetBackingStore().Get("workflowTemplates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkflowTemplateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LifecycleWorkflowsContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCustomTaskExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomTaskExtensions()))
        for i, v := range m.GetCustomTaskExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customTaskExtensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deletedItems", m.GetDeletedItems())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("insights", m.GetInsights())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    if m.GetTaskDefinitions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTaskDefinitions()))
        for i, v := range m.GetTaskDefinitions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("taskDefinitions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorkflows() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWorkflows()))
        for i, v := range m.GetWorkflows() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("workflows", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorkflowTemplates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWorkflowTemplates()))
        for i, v := range m.GetWorkflowTemplates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("workflowTemplates", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCustomTaskExtensions sets the customTaskExtensions property value. The customTaskExtension instance.
func (m *LifecycleWorkflowsContainer) SetCustomTaskExtensions(value []CustomTaskExtensionable)() {
    err := m.GetBackingStore().Set("customTaskExtensions", value)
    if err != nil {
        panic(err)
    }
}
// SetDeletedItems sets the deletedItems property value. Deleted workflows in your lifecycle workflows instance.
func (m *LifecycleWorkflowsContainer) SetDeletedItems(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable)() {
    err := m.GetBackingStore().Set("deletedItems", value)
    if err != nil {
        panic(err)
    }
}
// SetInsights sets the insights property value. The insight container holding workflow insight summaries for a tenant.
func (m *LifecycleWorkflowsContainer) SetInsights(value Insightsable)() {
    err := m.GetBackingStore().Set("insights", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. The settings property
func (m *LifecycleWorkflowsContainer) SetSettings(value LifecycleManagementSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskDefinitions sets the taskDefinitions property value. The definition of tasks within the lifecycle workflows instance.
func (m *LifecycleWorkflowsContainer) SetTaskDefinitions(value []TaskDefinitionable)() {
    err := m.GetBackingStore().Set("taskDefinitions", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflows sets the workflows property value. The workflows in the lifecycle workflows instance.
func (m *LifecycleWorkflowsContainer) SetWorkflows(value []Workflowable)() {
    err := m.GetBackingStore().Set("workflows", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowTemplates sets the workflowTemplates property value. The workflow templates in the lifecycle workflow instance.
func (m *LifecycleWorkflowsContainer) SetWorkflowTemplates(value []WorkflowTemplateable)() {
    err := m.GetBackingStore().Set("workflowTemplates", value)
    if err != nil {
        panic(err)
    }
}
type LifecycleWorkflowsContainerable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCustomTaskExtensions()([]CustomTaskExtensionable)
    GetDeletedItems()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable)
    GetInsights()(Insightsable)
    GetSettings()(LifecycleManagementSettingsable)
    GetTaskDefinitions()([]TaskDefinitionable)
    GetWorkflows()([]Workflowable)
    GetWorkflowTemplates()([]WorkflowTemplateable)
    SetCustomTaskExtensions(value []CustomTaskExtensionable)()
    SetDeletedItems(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeletedItemContainerable)()
    SetInsights(value Insightsable)()
    SetSettings(value LifecycleManagementSettingsable)()
    SetTaskDefinitions(value []TaskDefinitionable)()
    SetWorkflows(value []Workflowable)()
    SetWorkflowTemplates(value []WorkflowTemplateable)()
}

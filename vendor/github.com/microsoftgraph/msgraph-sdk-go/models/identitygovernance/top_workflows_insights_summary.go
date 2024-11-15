package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TopWorkflowsInsightsSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTopWorkflowsInsightsSummary instantiates a new TopWorkflowsInsightsSummary and sets the default values.
func NewTopWorkflowsInsightsSummary()(*TopWorkflowsInsightsSummary) {
    m := &TopWorkflowsInsightsSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTopWorkflowsInsightsSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTopWorkflowsInsightsSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTopWorkflowsInsightsSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TopWorkflowsInsightsSummary) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *TopWorkflowsInsightsSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFailedRuns gets the failedRuns property value. Count of failed runs for workflow.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetFailedRuns()(*int32) {
    val, err := m.GetBackingStore().Get("failedRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedUsers gets the failedUsers property value. Count of failed users who were processed.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetFailedUsers()(*int32) {
    val, err := m.GetBackingStore().Get("failedUsers")
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
func (m *TopWorkflowsInsightsSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["failedRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedRuns(val)
        }
        return nil
    }
    res["failedUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedUsers(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["successfulRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulRuns(val)
        }
        return nil
    }
    res["successfulUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulUsers(val)
        }
        return nil
    }
    res["totalRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalRuns(val)
        }
        return nil
    }
    res["totalUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalUsers(val)
        }
        return nil
    }
    res["workflowCategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLifecycleWorkflowCategory)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflowCategory(val.(*LifecycleWorkflowCategory))
        }
        return nil
    }
    res["workflowDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflowDisplayName(val)
        }
        return nil
    }
    res["workflowId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflowId(val)
        }
        return nil
    }
    res["workflowVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkflowVersion(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TopWorkflowsInsightsSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSuccessfulRuns gets the successfulRuns property value. Count of successful runs of the workflow.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetSuccessfulRuns()(*int32) {
    val, err := m.GetBackingStore().Get("successfulRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSuccessfulUsers gets the successfulUsers property value. Count of successful users processed by the workflow.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetSuccessfulUsers()(*int32) {
    val, err := m.GetBackingStore().Get("successfulUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalRuns gets the totalRuns property value. Count of total runs of workflow.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetTotalRuns()(*int32) {
    val, err := m.GetBackingStore().Get("totalRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalUsers gets the totalUsers property value. Total number of users processed by the workflow.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetTotalUsers()(*int32) {
    val, err := m.GetBackingStore().Get("totalUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWorkflowCategory gets the workflowCategory property value. The workflowCategory property
// returns a *LifecycleWorkflowCategory when successful
func (m *TopWorkflowsInsightsSummary) GetWorkflowCategory()(*LifecycleWorkflowCategory) {
    val, err := m.GetBackingStore().Get("workflowCategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*LifecycleWorkflowCategory)
    }
    return nil
}
// GetWorkflowDisplayName gets the workflowDisplayName property value. The name of the workflow.
// returns a *string when successful
func (m *TopWorkflowsInsightsSummary) GetWorkflowDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("workflowDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWorkflowId gets the workflowId property value. The workflow ID.
// returns a *string when successful
func (m *TopWorkflowsInsightsSummary) GetWorkflowId()(*string) {
    val, err := m.GetBackingStore().Get("workflowId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWorkflowVersion gets the workflowVersion property value. The version of the workflow that was a top workflow ran.
// returns a *int32 when successful
func (m *TopWorkflowsInsightsSummary) GetWorkflowVersion()(*int32) {
    val, err := m.GetBackingStore().Get("workflowVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TopWorkflowsInsightsSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("failedRuns", m.GetFailedRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("failedUsers", m.GetFailedUsers())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("successfulRuns", m.GetSuccessfulRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("successfulUsers", m.GetSuccessfulUsers())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalRuns", m.GetTotalRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalUsers", m.GetTotalUsers())
        if err != nil {
            return err
        }
    }
    if m.GetWorkflowCategory() != nil {
        cast := (*m.GetWorkflowCategory()).String()
        err := writer.WriteStringValue("workflowCategory", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("workflowDisplayName", m.GetWorkflowDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("workflowId", m.GetWorkflowId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("workflowVersion", m.GetWorkflowVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TopWorkflowsInsightsSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TopWorkflowsInsightsSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFailedRuns sets the failedRuns property value. Count of failed runs for workflow.
func (m *TopWorkflowsInsightsSummary) SetFailedRuns(value *int32)() {
    err := m.GetBackingStore().Set("failedRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedUsers sets the failedUsers property value. Count of failed users who were processed.
func (m *TopWorkflowsInsightsSummary) SetFailedUsers(value *int32)() {
    err := m.GetBackingStore().Set("failedUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TopWorkflowsInsightsSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulRuns sets the successfulRuns property value. Count of successful runs of the workflow.
func (m *TopWorkflowsInsightsSummary) SetSuccessfulRuns(value *int32)() {
    err := m.GetBackingStore().Set("successfulRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulUsers sets the successfulUsers property value. Count of successful users processed by the workflow.
func (m *TopWorkflowsInsightsSummary) SetSuccessfulUsers(value *int32)() {
    err := m.GetBackingStore().Set("successfulUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalRuns sets the totalRuns property value. Count of total runs of workflow.
func (m *TopWorkflowsInsightsSummary) SetTotalRuns(value *int32)() {
    err := m.GetBackingStore().Set("totalRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUsers sets the totalUsers property value. Total number of users processed by the workflow.
func (m *TopWorkflowsInsightsSummary) SetTotalUsers(value *int32)() {
    err := m.GetBackingStore().Set("totalUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowCategory sets the workflowCategory property value. The workflowCategory property
func (m *TopWorkflowsInsightsSummary) SetWorkflowCategory(value *LifecycleWorkflowCategory)() {
    err := m.GetBackingStore().Set("workflowCategory", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowDisplayName sets the workflowDisplayName property value. The name of the workflow.
func (m *TopWorkflowsInsightsSummary) SetWorkflowDisplayName(value *string)() {
    err := m.GetBackingStore().Set("workflowDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowId sets the workflowId property value. The workflow ID.
func (m *TopWorkflowsInsightsSummary) SetWorkflowId(value *string)() {
    err := m.GetBackingStore().Set("workflowId", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkflowVersion sets the workflowVersion property value. The version of the workflow that was a top workflow ran.
func (m *TopWorkflowsInsightsSummary) SetWorkflowVersion(value *int32)() {
    err := m.GetBackingStore().Set("workflowVersion", value)
    if err != nil {
        panic(err)
    }
}
type TopWorkflowsInsightsSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFailedRuns()(*int32)
    GetFailedUsers()(*int32)
    GetOdataType()(*string)
    GetSuccessfulRuns()(*int32)
    GetSuccessfulUsers()(*int32)
    GetTotalRuns()(*int32)
    GetTotalUsers()(*int32)
    GetWorkflowCategory()(*LifecycleWorkflowCategory)
    GetWorkflowDisplayName()(*string)
    GetWorkflowId()(*string)
    GetWorkflowVersion()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFailedRuns(value *int32)()
    SetFailedUsers(value *int32)()
    SetOdataType(value *string)()
    SetSuccessfulRuns(value *int32)()
    SetSuccessfulUsers(value *int32)()
    SetTotalRuns(value *int32)()
    SetTotalUsers(value *int32)()
    SetWorkflowCategory(value *LifecycleWorkflowCategory)()
    SetWorkflowDisplayName(value *string)()
    SetWorkflowId(value *string)()
    SetWorkflowVersion(value *int32)()
}

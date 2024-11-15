package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TopTasksInsightsSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTopTasksInsightsSummary instantiates a new TopTasksInsightsSummary and sets the default values.
func NewTopTasksInsightsSummary()(*TopTasksInsightsSummary) {
    m := &TopTasksInsightsSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTopTasksInsightsSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTopTasksInsightsSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTopTasksInsightsSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TopTasksInsightsSummary) GetAdditionalData()(map[string]any) {
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
func (m *TopTasksInsightsSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFailedTasks gets the failedTasks property value. Count of failed runs of the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetFailedTasks()(*int32) {
    val, err := m.GetBackingStore().Get("failedTasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedUsers gets the failedUsers property value. Count of failed users who were processed by the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetFailedUsers()(*int32) {
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
func (m *TopTasksInsightsSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["failedTasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedTasks(val)
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
    res["successfulTasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulTasks(val)
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
    res["taskDefinitionDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTaskDefinitionDisplayName(val)
        }
        return nil
    }
    res["taskDefinitionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTaskDefinitionId(val)
        }
        return nil
    }
    res["totalTasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalTasks(val)
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
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TopTasksInsightsSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSuccessfulTasks gets the successfulTasks property value. Count of successful runs of the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetSuccessfulTasks()(*int32) {
    val, err := m.GetBackingStore().Get("successfulTasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSuccessfulUsers gets the successfulUsers property value. Count of successful users processed by the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetSuccessfulUsers()(*int32) {
    val, err := m.GetBackingStore().Get("successfulUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTaskDefinitionDisplayName gets the taskDefinitionDisplayName property value. The name of the task.
// returns a *string when successful
func (m *TopTasksInsightsSummary) GetTaskDefinitionDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("taskDefinitionDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTaskDefinitionId gets the taskDefinitionId property value. The task ID.
// returns a *string when successful
func (m *TopTasksInsightsSummary) GetTaskDefinitionId()(*string) {
    val, err := m.GetBackingStore().Get("taskDefinitionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalTasks gets the totalTasks property value. Count of total runs of the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetTotalTasks()(*int32) {
    val, err := m.GetBackingStore().Get("totalTasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalUsers gets the totalUsers property value. Count of total users processed by the task.
// returns a *int32 when successful
func (m *TopTasksInsightsSummary) GetTotalUsers()(*int32) {
    val, err := m.GetBackingStore().Get("totalUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TopTasksInsightsSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("failedTasks", m.GetFailedTasks())
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
        err := writer.WriteInt32Value("successfulTasks", m.GetSuccessfulTasks())
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
        err := writer.WriteStringValue("taskDefinitionDisplayName", m.GetTaskDefinitionDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("taskDefinitionId", m.GetTaskDefinitionId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalTasks", m.GetTotalTasks())
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
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TopTasksInsightsSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TopTasksInsightsSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFailedTasks sets the failedTasks property value. Count of failed runs of the task.
func (m *TopTasksInsightsSummary) SetFailedTasks(value *int32)() {
    err := m.GetBackingStore().Set("failedTasks", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedUsers sets the failedUsers property value. Count of failed users who were processed by the task.
func (m *TopTasksInsightsSummary) SetFailedUsers(value *int32)() {
    err := m.GetBackingStore().Set("failedUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TopTasksInsightsSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulTasks sets the successfulTasks property value. Count of successful runs of the task.
func (m *TopTasksInsightsSummary) SetSuccessfulTasks(value *int32)() {
    err := m.GetBackingStore().Set("successfulTasks", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulUsers sets the successfulUsers property value. Count of successful users processed by the task.
func (m *TopTasksInsightsSummary) SetSuccessfulUsers(value *int32)() {
    err := m.GetBackingStore().Set("successfulUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskDefinitionDisplayName sets the taskDefinitionDisplayName property value. The name of the task.
func (m *TopTasksInsightsSummary) SetTaskDefinitionDisplayName(value *string)() {
    err := m.GetBackingStore().Set("taskDefinitionDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetTaskDefinitionId sets the taskDefinitionId property value. The task ID.
func (m *TopTasksInsightsSummary) SetTaskDefinitionId(value *string)() {
    err := m.GetBackingStore().Set("taskDefinitionId", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalTasks sets the totalTasks property value. Count of total runs of the task.
func (m *TopTasksInsightsSummary) SetTotalTasks(value *int32)() {
    err := m.GetBackingStore().Set("totalTasks", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUsers sets the totalUsers property value. Count of total users processed by the task.
func (m *TopTasksInsightsSummary) SetTotalUsers(value *int32)() {
    err := m.GetBackingStore().Set("totalUsers", value)
    if err != nil {
        panic(err)
    }
}
type TopTasksInsightsSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFailedTasks()(*int32)
    GetFailedUsers()(*int32)
    GetOdataType()(*string)
    GetSuccessfulTasks()(*int32)
    GetSuccessfulUsers()(*int32)
    GetTaskDefinitionDisplayName()(*string)
    GetTaskDefinitionId()(*string)
    GetTotalTasks()(*int32)
    GetTotalUsers()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFailedTasks(value *int32)()
    SetFailedUsers(value *int32)()
    SetOdataType(value *string)()
    SetSuccessfulTasks(value *int32)()
    SetSuccessfulUsers(value *int32)()
    SetTaskDefinitionDisplayName(value *string)()
    SetTaskDefinitionId(value *string)()
    SetTotalTasks(value *int32)()
    SetTotalUsers(value *int32)()
}

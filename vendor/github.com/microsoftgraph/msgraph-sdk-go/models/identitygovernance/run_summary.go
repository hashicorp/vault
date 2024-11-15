package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type RunSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRunSummary instantiates a new RunSummary and sets the default values.
func NewRunSummary()(*RunSummary) {
    m := &RunSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRunSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRunSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRunSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *RunSummary) GetAdditionalData()(map[string]any) {
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
func (m *RunSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFailedRuns gets the failedRuns property value. The number of failed workflow runs.
// returns a *int32 when successful
func (m *RunSummary) GetFailedRuns()(*int32) {
    val, err := m.GetBackingStore().Get("failedRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedTasks gets the failedTasks property value. The number of failed tasks of a workflow.
// returns a *int32 when successful
func (m *RunSummary) GetFailedTasks()(*int32) {
    val, err := m.GetBackingStore().Get("failedTasks")
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
func (m *RunSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
func (m *RunSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSuccessfulRuns gets the successfulRuns property value. The number of successful workflow runs.
// returns a *int32 when successful
func (m *RunSummary) GetSuccessfulRuns()(*int32) {
    val, err := m.GetBackingStore().Get("successfulRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalRuns gets the totalRuns property value. The total number of runs for a workflow.
// returns a *int32 when successful
func (m *RunSummary) GetTotalRuns()(*int32) {
    val, err := m.GetBackingStore().Get("totalRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalTasks gets the totalTasks property value. The total number of tasks processed by a workflow.
// returns a *int32 when successful
func (m *RunSummary) GetTotalTasks()(*int32) {
    val, err := m.GetBackingStore().Get("totalTasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalUsers gets the totalUsers property value. The total number of users processed by a workflow.
// returns a *int32 when successful
func (m *RunSummary) GetTotalUsers()(*int32) {
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
func (m *RunSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("failedRuns", m.GetFailedRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("failedTasks", m.GetFailedTasks())
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
        err := writer.WriteInt32Value("totalRuns", m.GetTotalRuns())
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
func (m *RunSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *RunSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFailedRuns sets the failedRuns property value. The number of failed workflow runs.
func (m *RunSummary) SetFailedRuns(value *int32)() {
    err := m.GetBackingStore().Set("failedRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedTasks sets the failedTasks property value. The number of failed tasks of a workflow.
func (m *RunSummary) SetFailedTasks(value *int32)() {
    err := m.GetBackingStore().Set("failedTasks", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RunSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulRuns sets the successfulRuns property value. The number of successful workflow runs.
func (m *RunSummary) SetSuccessfulRuns(value *int32)() {
    err := m.GetBackingStore().Set("successfulRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalRuns sets the totalRuns property value. The total number of runs for a workflow.
func (m *RunSummary) SetTotalRuns(value *int32)() {
    err := m.GetBackingStore().Set("totalRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalTasks sets the totalTasks property value. The total number of tasks processed by a workflow.
func (m *RunSummary) SetTotalTasks(value *int32)() {
    err := m.GetBackingStore().Set("totalTasks", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUsers sets the totalUsers property value. The total number of users processed by a workflow.
func (m *RunSummary) SetTotalUsers(value *int32)() {
    err := m.GetBackingStore().Set("totalUsers", value)
    if err != nil {
        panic(err)
    }
}
type RunSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFailedRuns()(*int32)
    GetFailedTasks()(*int32)
    GetOdataType()(*string)
    GetSuccessfulRuns()(*int32)
    GetTotalRuns()(*int32)
    GetTotalTasks()(*int32)
    GetTotalUsers()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFailedRuns(value *int32)()
    SetFailedTasks(value *int32)()
    SetOdataType(value *string)()
    SetSuccessfulRuns(value *int32)()
    SetTotalRuns(value *int32)()
    SetTotalTasks(value *int32)()
    SetTotalUsers(value *int32)()
}

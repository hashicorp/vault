package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type WorkflowsInsightsByCategory struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewWorkflowsInsightsByCategory instantiates a new WorkflowsInsightsByCategory and sets the default values.
func NewWorkflowsInsightsByCategory()(*WorkflowsInsightsByCategory) {
    m := &WorkflowsInsightsByCategory{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateWorkflowsInsightsByCategoryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkflowsInsightsByCategoryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkflowsInsightsByCategory(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *WorkflowsInsightsByCategory) GetAdditionalData()(map[string]any) {
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
func (m *WorkflowsInsightsByCategory) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFailedJoinerRuns gets the failedJoinerRuns property value. Failed 'Joiner' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetFailedJoinerRuns()(*int32) {
    val, err := m.GetBackingStore().Get("failedJoinerRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedLeaverRuns gets the failedLeaverRuns property value. Failed 'Leaver' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetFailedLeaverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("failedLeaverRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedMoverRuns gets the failedMoverRuns property value. Failed 'Mover' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetFailedMoverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("failedMoverRuns")
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
func (m *WorkflowsInsightsByCategory) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["failedJoinerRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedJoinerRuns(val)
        }
        return nil
    }
    res["failedLeaverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedLeaverRuns(val)
        }
        return nil
    }
    res["failedMoverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedMoverRuns(val)
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
    res["successfulJoinerRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulJoinerRuns(val)
        }
        return nil
    }
    res["successfulLeaverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulLeaverRuns(val)
        }
        return nil
    }
    res["successfulMoverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessfulMoverRuns(val)
        }
        return nil
    }
    res["totalJoinerRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalJoinerRuns(val)
        }
        return nil
    }
    res["totalLeaverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalLeaverRuns(val)
        }
        return nil
    }
    res["totalMoverRuns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalMoverRuns(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *WorkflowsInsightsByCategory) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSuccessfulJoinerRuns gets the successfulJoinerRuns property value. Successful 'Joiner' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetSuccessfulJoinerRuns()(*int32) {
    val, err := m.GetBackingStore().Get("successfulJoinerRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSuccessfulLeaverRuns gets the successfulLeaverRuns property value. Successful 'Leaver' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetSuccessfulLeaverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("successfulLeaverRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSuccessfulMoverRuns gets the successfulMoverRuns property value. Successful 'Mover' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetSuccessfulMoverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("successfulMoverRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalJoinerRuns gets the totalJoinerRuns property value. Total 'Joiner' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetTotalJoinerRuns()(*int32) {
    val, err := m.GetBackingStore().Get("totalJoinerRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalLeaverRuns gets the totalLeaverRuns property value. Total 'Leaver' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetTotalLeaverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("totalLeaverRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalMoverRuns gets the totalMoverRuns property value. Total 'Mover' workflows processed in a tenant.
// returns a *int32 when successful
func (m *WorkflowsInsightsByCategory) GetTotalMoverRuns()(*int32) {
    val, err := m.GetBackingStore().Get("totalMoverRuns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkflowsInsightsByCategory) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("failedJoinerRuns", m.GetFailedJoinerRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("failedLeaverRuns", m.GetFailedLeaverRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("failedMoverRuns", m.GetFailedMoverRuns())
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
        err := writer.WriteInt32Value("successfulJoinerRuns", m.GetSuccessfulJoinerRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("successfulLeaverRuns", m.GetSuccessfulLeaverRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("successfulMoverRuns", m.GetSuccessfulMoverRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalJoinerRuns", m.GetTotalJoinerRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalLeaverRuns", m.GetTotalLeaverRuns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalMoverRuns", m.GetTotalMoverRuns())
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
func (m *WorkflowsInsightsByCategory) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *WorkflowsInsightsByCategory) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFailedJoinerRuns sets the failedJoinerRuns property value. Failed 'Joiner' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetFailedJoinerRuns(value *int32)() {
    err := m.GetBackingStore().Set("failedJoinerRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedLeaverRuns sets the failedLeaverRuns property value. Failed 'Leaver' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetFailedLeaverRuns(value *int32)() {
    err := m.GetBackingStore().Set("failedLeaverRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedMoverRuns sets the failedMoverRuns property value. Failed 'Mover' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetFailedMoverRuns(value *int32)() {
    err := m.GetBackingStore().Set("failedMoverRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WorkflowsInsightsByCategory) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulJoinerRuns sets the successfulJoinerRuns property value. Successful 'Joiner' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetSuccessfulJoinerRuns(value *int32)() {
    err := m.GetBackingStore().Set("successfulJoinerRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulLeaverRuns sets the successfulLeaverRuns property value. Successful 'Leaver' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetSuccessfulLeaverRuns(value *int32)() {
    err := m.GetBackingStore().Set("successfulLeaverRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessfulMoverRuns sets the successfulMoverRuns property value. Successful 'Mover' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetSuccessfulMoverRuns(value *int32)() {
    err := m.GetBackingStore().Set("successfulMoverRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalJoinerRuns sets the totalJoinerRuns property value. Total 'Joiner' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetTotalJoinerRuns(value *int32)() {
    err := m.GetBackingStore().Set("totalJoinerRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLeaverRuns sets the totalLeaverRuns property value. Total 'Leaver' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetTotalLeaverRuns(value *int32)() {
    err := m.GetBackingStore().Set("totalLeaverRuns", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalMoverRuns sets the totalMoverRuns property value. Total 'Mover' workflows processed in a tenant.
func (m *WorkflowsInsightsByCategory) SetTotalMoverRuns(value *int32)() {
    err := m.GetBackingStore().Set("totalMoverRuns", value)
    if err != nil {
        panic(err)
    }
}
type WorkflowsInsightsByCategoryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFailedJoinerRuns()(*int32)
    GetFailedLeaverRuns()(*int32)
    GetFailedMoverRuns()(*int32)
    GetOdataType()(*string)
    GetSuccessfulJoinerRuns()(*int32)
    GetSuccessfulLeaverRuns()(*int32)
    GetSuccessfulMoverRuns()(*int32)
    GetTotalJoinerRuns()(*int32)
    GetTotalLeaverRuns()(*int32)
    GetTotalMoverRuns()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFailedJoinerRuns(value *int32)()
    SetFailedLeaverRuns(value *int32)()
    SetFailedMoverRuns(value *int32)()
    SetOdataType(value *string)()
    SetSuccessfulJoinerRuns(value *int32)()
    SetSuccessfulLeaverRuns(value *int32)()
    SetSuccessfulMoverRuns(value *int32)()
    SetTotalJoinerRuns(value *int32)()
    SetTotalLeaverRuns(value *int32)()
    SetTotalMoverRuns(value *int32)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UnifiedApprovalStage struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUnifiedApprovalStage instantiates a new UnifiedApprovalStage and sets the default values.
func NewUnifiedApprovalStage()(*UnifiedApprovalStage) {
    m := &UnifiedApprovalStage{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUnifiedApprovalStageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedApprovalStageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedApprovalStage(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UnifiedApprovalStage) GetAdditionalData()(map[string]any) {
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
// GetApprovalStageTimeOutInDays gets the approvalStageTimeOutInDays property value. The number of days that a request can be pending a response before it is automatically denied.
// returns a *int32 when successful
func (m *UnifiedApprovalStage) GetApprovalStageTimeOutInDays()(*int32) {
    val, err := m.GetBackingStore().Get("approvalStageTimeOutInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *UnifiedApprovalStage) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEscalationApprovers gets the escalationApprovers property value. The escalation approvers for this stage when the primary approvers don't respond.
// returns a []SubjectSetable when successful
func (m *UnifiedApprovalStage) GetEscalationApprovers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("escalationApprovers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// GetEscalationTimeInMinutes gets the escalationTimeInMinutes property value. The time a request can be pending a response from a primary approver before it can be escalated to the escalation approvers.
// returns a *int32 when successful
func (m *UnifiedApprovalStage) GetEscalationTimeInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("escalationTimeInMinutes")
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
func (m *UnifiedApprovalStage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["approvalStageTimeOutInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApprovalStageTimeOutInDays(val)
        }
        return nil
    }
    res["escalationApprovers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SubjectSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SubjectSetable)
                }
            }
            m.SetEscalationApprovers(res)
        }
        return nil
    }
    res["escalationTimeInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEscalationTimeInMinutes(val)
        }
        return nil
    }
    res["isApproverJustificationRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsApproverJustificationRequired(val)
        }
        return nil
    }
    res["isEscalationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEscalationEnabled(val)
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
    res["primaryApprovers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SubjectSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SubjectSetable)
                }
            }
            m.SetPrimaryApprovers(res)
        }
        return nil
    }
    return res
}
// GetIsApproverJustificationRequired gets the isApproverJustificationRequired property value. Indicates whether the approver must provide justification for their reponse.
// returns a *bool when successful
func (m *UnifiedApprovalStage) GetIsApproverJustificationRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isApproverJustificationRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsEscalationEnabled gets the isEscalationEnabled property value. Indicates whether escalation if enabled.
// returns a *bool when successful
func (m *UnifiedApprovalStage) GetIsEscalationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEscalationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UnifiedApprovalStage) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrimaryApprovers gets the primaryApprovers property value. The primary approvers of this stage.
// returns a []SubjectSetable when successful
func (m *UnifiedApprovalStage) GetPrimaryApprovers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("primaryApprovers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedApprovalStage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("approvalStageTimeOutInDays", m.GetApprovalStageTimeOutInDays())
        if err != nil {
            return err
        }
    }
    if m.GetEscalationApprovers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEscalationApprovers()))
        for i, v := range m.GetEscalationApprovers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("escalationApprovers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("escalationTimeInMinutes", m.GetEscalationTimeInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApproverJustificationRequired", m.GetIsApproverJustificationRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEscalationEnabled", m.GetIsEscalationEnabled())
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
    if m.GetPrimaryApprovers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPrimaryApprovers()))
        for i, v := range m.GetPrimaryApprovers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("primaryApprovers", cast)
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
func (m *UnifiedApprovalStage) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetApprovalStageTimeOutInDays sets the approvalStageTimeOutInDays property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *UnifiedApprovalStage) SetApprovalStageTimeOutInDays(value *int32)() {
    err := m.GetBackingStore().Set("approvalStageTimeOutInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UnifiedApprovalStage) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEscalationApprovers sets the escalationApprovers property value. The escalation approvers for this stage when the primary approvers don't respond.
func (m *UnifiedApprovalStage) SetEscalationApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("escalationApprovers", value)
    if err != nil {
        panic(err)
    }
}
// SetEscalationTimeInMinutes sets the escalationTimeInMinutes property value. The time a request can be pending a response from a primary approver before it can be escalated to the escalation approvers.
func (m *UnifiedApprovalStage) SetEscalationTimeInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("escalationTimeInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetIsApproverJustificationRequired sets the isApproverJustificationRequired property value. Indicates whether the approver must provide justification for their reponse.
func (m *UnifiedApprovalStage) SetIsApproverJustificationRequired(value *bool)() {
    err := m.GetBackingStore().Set("isApproverJustificationRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEscalationEnabled sets the isEscalationEnabled property value. Indicates whether escalation if enabled.
func (m *UnifiedApprovalStage) SetIsEscalationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEscalationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UnifiedApprovalStage) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryApprovers sets the primaryApprovers property value. The primary approvers of this stage.
func (m *UnifiedApprovalStage) SetPrimaryApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("primaryApprovers", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedApprovalStageable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApprovalStageTimeOutInDays()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEscalationApprovers()([]SubjectSetable)
    GetEscalationTimeInMinutes()(*int32)
    GetIsApproverJustificationRequired()(*bool)
    GetIsEscalationEnabled()(*bool)
    GetOdataType()(*string)
    GetPrimaryApprovers()([]SubjectSetable)
    SetApprovalStageTimeOutInDays(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEscalationApprovers(value []SubjectSetable)()
    SetEscalationTimeInMinutes(value *int32)()
    SetIsApproverJustificationRequired(value *bool)()
    SetIsEscalationEnabled(value *bool)()
    SetOdataType(value *string)()
    SetPrimaryApprovers(value []SubjectSetable)()
}

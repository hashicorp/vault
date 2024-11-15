package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessPackageApprovalStage struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessPackageApprovalStage instantiates a new AccessPackageApprovalStage and sets the default values.
func NewAccessPackageApprovalStage()(*AccessPackageApprovalStage) {
    m := &AccessPackageApprovalStage{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessPackageApprovalStageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageApprovalStageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageApprovalStage(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessPackageApprovalStage) GetAdditionalData()(map[string]any) {
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
func (m *AccessPackageApprovalStage) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDurationBeforeAutomaticDenial gets the durationBeforeAutomaticDenial property value. The number of days that a request can be pending a response before it is automatically denied.
// returns a *ISODuration when successful
func (m *AccessPackageApprovalStage) GetDurationBeforeAutomaticDenial()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("durationBeforeAutomaticDenial")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetDurationBeforeEscalation gets the durationBeforeEscalation property value. If escalation is required, the time a request can be pending a response from a primary approver.
// returns a *ISODuration when successful
func (m *AccessPackageApprovalStage) GetDurationBeforeEscalation()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("durationBeforeEscalation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetEscalationApprovers gets the escalationApprovers property value. If escalation is enabled and the primary approvers do not respond before the escalation time, the escalationApprovers are the users who will be asked to approve requests.
// returns a []SubjectSetable when successful
func (m *AccessPackageApprovalStage) GetEscalationApprovers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("escalationApprovers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// GetFallbackEscalationApprovers gets the fallbackEscalationApprovers property value. The subjects, typically users, who are the fallback escalation approvers.
// returns a []SubjectSetable when successful
func (m *AccessPackageApprovalStage) GetFallbackEscalationApprovers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("fallbackEscalationApprovers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// GetFallbackPrimaryApprovers gets the fallbackPrimaryApprovers property value. The subjects, typically users, who are the fallback primary approvers.
// returns a []SubjectSetable when successful
func (m *AccessPackageApprovalStage) GetFallbackPrimaryApprovers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("fallbackPrimaryApprovers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageApprovalStage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["durationBeforeAutomaticDenial"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationBeforeAutomaticDenial(val)
        }
        return nil
    }
    res["durationBeforeEscalation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationBeforeEscalation(val)
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
    res["fallbackEscalationApprovers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFallbackEscalationApprovers(res)
        }
        return nil
    }
    res["fallbackPrimaryApprovers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFallbackPrimaryApprovers(res)
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
// GetIsApproverJustificationRequired gets the isApproverJustificationRequired property value. Indicates whether the approver is required to provide a justification for approving a request.
// returns a *bool when successful
func (m *AccessPackageApprovalStage) GetIsApproverJustificationRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isApproverJustificationRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsEscalationEnabled gets the isEscalationEnabled property value. If true, then one or more escalationApprovers are configured in this approval stage.
// returns a *bool when successful
func (m *AccessPackageApprovalStage) GetIsEscalationEnabled()(*bool) {
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
func (m *AccessPackageApprovalStage) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrimaryApprovers gets the primaryApprovers property value. The subjects, typically users, who will be asked to approve requests. A collection of singleUser, groupMembers, requestorManager, internalSponsors, externalSponsors, or targetUserSponsors.
// returns a []SubjectSetable when successful
func (m *AccessPackageApprovalStage) GetPrimaryApprovers()([]SubjectSetable) {
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
func (m *AccessPackageApprovalStage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteISODurationValue("durationBeforeAutomaticDenial", m.GetDurationBeforeAutomaticDenial())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("durationBeforeEscalation", m.GetDurationBeforeEscalation())
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
    if m.GetFallbackEscalationApprovers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFallbackEscalationApprovers()))
        for i, v := range m.GetFallbackEscalationApprovers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("fallbackEscalationApprovers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFallbackPrimaryApprovers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFallbackPrimaryApprovers()))
        for i, v := range m.GetFallbackPrimaryApprovers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("fallbackPrimaryApprovers", cast)
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
func (m *AccessPackageApprovalStage) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessPackageApprovalStage) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDurationBeforeAutomaticDenial sets the durationBeforeAutomaticDenial property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *AccessPackageApprovalStage) SetDurationBeforeAutomaticDenial(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("durationBeforeAutomaticDenial", value)
    if err != nil {
        panic(err)
    }
}
// SetDurationBeforeEscalation sets the durationBeforeEscalation property value. If escalation is required, the time a request can be pending a response from a primary approver.
func (m *AccessPackageApprovalStage) SetDurationBeforeEscalation(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("durationBeforeEscalation", value)
    if err != nil {
        panic(err)
    }
}
// SetEscalationApprovers sets the escalationApprovers property value. If escalation is enabled and the primary approvers do not respond before the escalation time, the escalationApprovers are the users who will be asked to approve requests.
func (m *AccessPackageApprovalStage) SetEscalationApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("escalationApprovers", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackEscalationApprovers sets the fallbackEscalationApprovers property value. The subjects, typically users, who are the fallback escalation approvers.
func (m *AccessPackageApprovalStage) SetFallbackEscalationApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("fallbackEscalationApprovers", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackPrimaryApprovers sets the fallbackPrimaryApprovers property value. The subjects, typically users, who are the fallback primary approvers.
func (m *AccessPackageApprovalStage) SetFallbackPrimaryApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("fallbackPrimaryApprovers", value)
    if err != nil {
        panic(err)
    }
}
// SetIsApproverJustificationRequired sets the isApproverJustificationRequired property value. Indicates whether the approver is required to provide a justification for approving a request.
func (m *AccessPackageApprovalStage) SetIsApproverJustificationRequired(value *bool)() {
    err := m.GetBackingStore().Set("isApproverJustificationRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEscalationEnabled sets the isEscalationEnabled property value. If true, then one or more escalationApprovers are configured in this approval stage.
func (m *AccessPackageApprovalStage) SetIsEscalationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEscalationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageApprovalStage) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryApprovers sets the primaryApprovers property value. The subjects, typically users, who will be asked to approve requests. A collection of singleUser, groupMembers, requestorManager, internalSponsors, externalSponsors, or targetUserSponsors.
func (m *AccessPackageApprovalStage) SetPrimaryApprovers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("primaryApprovers", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageApprovalStageable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDurationBeforeAutomaticDenial()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetDurationBeforeEscalation()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetEscalationApprovers()([]SubjectSetable)
    GetFallbackEscalationApprovers()([]SubjectSetable)
    GetFallbackPrimaryApprovers()([]SubjectSetable)
    GetIsApproverJustificationRequired()(*bool)
    GetIsEscalationEnabled()(*bool)
    GetOdataType()(*string)
    GetPrimaryApprovers()([]SubjectSetable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDurationBeforeAutomaticDenial(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetDurationBeforeEscalation(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetEscalationApprovers(value []SubjectSetable)()
    SetFallbackEscalationApprovers(value []SubjectSetable)()
    SetFallbackPrimaryApprovers(value []SubjectSetable)()
    SetIsApproverJustificationRequired(value *bool)()
    SetIsEscalationEnabled(value *bool)()
    SetOdataType(value *string)()
    SetPrimaryApprovers(value []SubjectSetable)()
}

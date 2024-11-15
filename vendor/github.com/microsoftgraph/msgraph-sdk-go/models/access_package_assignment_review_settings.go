package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessPackageAssignmentReviewSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessPackageAssignmentReviewSettings instantiates a new AccessPackageAssignmentReviewSettings and sets the default values.
func NewAccessPackageAssignmentReviewSettings()(*AccessPackageAssignmentReviewSettings) {
    m := &AccessPackageAssignmentReviewSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessPackageAssignmentReviewSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentReviewSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentReviewSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessPackageAssignmentReviewSettings) GetAdditionalData()(map[string]any) {
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
func (m *AccessPackageAssignmentReviewSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetExpirationBehavior gets the expirationBehavior property value. The default decision to apply if the access is not reviewed. The possible values are: keepAccess, removeAccess, acceptAccessRecommendation, unknownFutureValue.
// returns a *AccessReviewExpirationBehavior when successful
func (m *AccessPackageAssignmentReviewSettings) GetExpirationBehavior()(*AccessReviewExpirationBehavior) {
    val, err := m.GetBackingStore().Get("expirationBehavior")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessReviewExpirationBehavior)
    }
    return nil
}
// GetFallbackReviewers gets the fallbackReviewers property value. This collection specifies the users who will be the fallback reviewers when the primary reviewers don't respond.
// returns a []SubjectSetable when successful
func (m *AccessPackageAssignmentReviewSettings) GetFallbackReviewers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("fallbackReviewers")
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
func (m *AccessPackageAssignmentReviewSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["expirationBehavior"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessReviewExpirationBehavior)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpirationBehavior(val.(*AccessReviewExpirationBehavior))
        }
        return nil
    }
    res["fallbackReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFallbackReviewers(res)
        }
        return nil
    }
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
        }
        return nil
    }
    res["isRecommendationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRecommendationEnabled(val)
        }
        return nil
    }
    res["isReviewerJustificationRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReviewerJustificationRequired(val)
        }
        return nil
    }
    res["isSelfReview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSelfReview(val)
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
    res["primaryReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetPrimaryReviewers(res)
        }
        return nil
    }
    res["schedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntitlementManagementScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedule(val.(EntitlementManagementScheduleable))
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. If true, access reviews are required for assignments through this policy.
// returns a *bool when successful
func (m *AccessPackageAssignmentReviewSettings) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRecommendationEnabled gets the isRecommendationEnabled property value. Specifies whether to display recommendations to the reviewer. The default value is true.
// returns a *bool when successful
func (m *AccessPackageAssignmentReviewSettings) GetIsRecommendationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isRecommendationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReviewerJustificationRequired gets the isReviewerJustificationRequired property value. Specifies whether the reviewer must provide justification for the approval. The default value is true.
// returns a *bool when successful
func (m *AccessPackageAssignmentReviewSettings) GetIsReviewerJustificationRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isReviewerJustificationRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSelfReview gets the isSelfReview property value. Specifies whether the principals can review their own assignments.
// returns a *bool when successful
func (m *AccessPackageAssignmentReviewSettings) GetIsSelfReview()(*bool) {
    val, err := m.GetBackingStore().Get("isSelfReview")
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
func (m *AccessPackageAssignmentReviewSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrimaryReviewers gets the primaryReviewers property value. This collection specifies the users or group of users who will review the access package assignments.
// returns a []SubjectSetable when successful
func (m *AccessPackageAssignmentReviewSettings) GetPrimaryReviewers()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("primaryReviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// GetSchedule gets the schedule property value. When the first review should start and how often it should recur.
// returns a EntitlementManagementScheduleable when successful
func (m *AccessPackageAssignmentReviewSettings) GetSchedule()(EntitlementManagementScheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementScheduleable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentReviewSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetExpirationBehavior() != nil {
        cast := (*m.GetExpirationBehavior()).String()
        err := writer.WriteStringValue("expirationBehavior", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFallbackReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFallbackReviewers()))
        for i, v := range m.GetFallbackReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("fallbackReviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRecommendationEnabled", m.GetIsRecommendationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isReviewerJustificationRequired", m.GetIsReviewerJustificationRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isSelfReview", m.GetIsSelfReview())
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
    if m.GetPrimaryReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPrimaryReviewers()))
        for i, v := range m.GetPrimaryReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("primaryReviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("schedule", m.GetSchedule())
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
func (m *AccessPackageAssignmentReviewSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessPackageAssignmentReviewSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetExpirationBehavior sets the expirationBehavior property value. The default decision to apply if the access is not reviewed. The possible values are: keepAccess, removeAccess, acceptAccessRecommendation, unknownFutureValue.
func (m *AccessPackageAssignmentReviewSettings) SetExpirationBehavior(value *AccessReviewExpirationBehavior)() {
    err := m.GetBackingStore().Set("expirationBehavior", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackReviewers sets the fallbackReviewers property value. This collection specifies the users who will be the fallback reviewers when the primary reviewers don't respond.
func (m *AccessPackageAssignmentReviewSettings) SetFallbackReviewers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("fallbackReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. If true, access reviews are required for assignments through this policy.
func (m *AccessPackageAssignmentReviewSettings) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRecommendationEnabled sets the isRecommendationEnabled property value. Specifies whether to display recommendations to the reviewer. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) SetIsRecommendationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isRecommendationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReviewerJustificationRequired sets the isReviewerJustificationRequired property value. Specifies whether the reviewer must provide justification for the approval. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) SetIsReviewerJustificationRequired(value *bool)() {
    err := m.GetBackingStore().Set("isReviewerJustificationRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSelfReview sets the isSelfReview property value. Specifies whether the principals can review their own assignments.
func (m *AccessPackageAssignmentReviewSettings) SetIsSelfReview(value *bool)() {
    err := m.GetBackingStore().Set("isSelfReview", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageAssignmentReviewSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryReviewers sets the primaryReviewers property value. This collection specifies the users or group of users who will review the access package assignments.
func (m *AccessPackageAssignmentReviewSettings) SetPrimaryReviewers(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("primaryReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. When the first review should start and how often it should recur.
func (m *AccessPackageAssignmentReviewSettings) SetSchedule(value EntitlementManagementScheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentReviewSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetExpirationBehavior()(*AccessReviewExpirationBehavior)
    GetFallbackReviewers()([]SubjectSetable)
    GetIsEnabled()(*bool)
    GetIsRecommendationEnabled()(*bool)
    GetIsReviewerJustificationRequired()(*bool)
    GetIsSelfReview()(*bool)
    GetOdataType()(*string)
    GetPrimaryReviewers()([]SubjectSetable)
    GetSchedule()(EntitlementManagementScheduleable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetExpirationBehavior(value *AccessReviewExpirationBehavior)()
    SetFallbackReviewers(value []SubjectSetable)()
    SetIsEnabled(value *bool)()
    SetIsRecommendationEnabled(value *bool)()
    SetIsReviewerJustificationRequired(value *bool)()
    SetIsSelfReview(value *bool)()
    SetOdataType(value *string)()
    SetPrimaryReviewers(value []SubjectSetable)()
    SetSchedule(value EntitlementManagementScheduleable)()
}

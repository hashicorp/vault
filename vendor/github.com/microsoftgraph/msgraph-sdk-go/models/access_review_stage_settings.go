package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessReviewStageSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessReviewStageSettings instantiates a new AccessReviewStageSettings and sets the default values.
func NewAccessReviewStageSettings()(*AccessReviewStageSettings) {
    m := &AccessReviewStageSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessReviewStageSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewStageSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewStageSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessReviewStageSettings) GetAdditionalData()(map[string]any) {
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
func (m *AccessReviewStageSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDecisionsThatWillMoveToNextStage gets the decisionsThatWillMoveToNextStage property value. Indicate which decisions will go to the next stage. Can be a subset of Approve, Deny, Recommendation, or NotReviewed. If not provided, all decisions will go to the next stage. Optional.
// returns a []string when successful
func (m *AccessReviewStageSettings) GetDecisionsThatWillMoveToNextStage()([]string) {
    val, err := m.GetBackingStore().Get("decisionsThatWillMoveToNextStage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDependsOn gets the dependsOn property value. Defines the sequential or parallel order of the stages and depends on the stageId. Only sequential stages are currently supported. For example, if stageId is 2, then dependsOn must be 1. If stageId is 1, don't specify dependsOn. Required if stageId isn't 1.
// returns a []string when successful
func (m *AccessReviewStageSettings) GetDependsOn()([]string) {
    val, err := m.GetBackingStore().Get("dependsOn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDurationInDays gets the durationInDays property value. The duration of the stage. Required.  NOTE: The cumulative value of this property across all stages  1. Will override the instanceDurationInDays setting on the accessReviewScheduleDefinition object. 2. Can't exceed the length of one recurrence. That is, if the review recurs weekly, the cumulative durationInDays can't exceed 7.
// returns a *int32 when successful
func (m *AccessReviewStageSettings) GetDurationInDays()(*int32) {
    val, err := m.GetBackingStore().Get("durationInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFallbackReviewers gets the fallbackReviewers property value. If provided, the fallback reviewers are asked to complete a review if the primary reviewers don't exist. For example, if managers are selected as reviewers and a principal under review doesn't have a manager in Microsoft Entra ID, the fallback reviewers are asked to review that principal. NOTE: The value of this property overrides the corresponding setting on the accessReviewScheduleDefinition object.
// returns a []AccessReviewReviewerScopeable when successful
func (m *AccessReviewStageSettings) GetFallbackReviewers()([]AccessReviewReviewerScopeable) {
    val, err := m.GetBackingStore().Get("fallbackReviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewReviewerScopeable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewStageSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["decisionsThatWillMoveToNextStage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetDecisionsThatWillMoveToNextStage(res)
        }
        return nil
    }
    res["dependsOn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetDependsOn(res)
        }
        return nil
    }
    res["durationInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationInDays(val)
        }
        return nil
    }
    res["fallbackReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewReviewerScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewReviewerScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewReviewerScopeable)
                }
            }
            m.SetFallbackReviewers(res)
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
    res["recommendationInsightSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewRecommendationInsightSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewRecommendationInsightSettingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewRecommendationInsightSettingable)
                }
            }
            m.SetRecommendationInsightSettings(res)
        }
        return nil
    }
    res["recommendationsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecommendationsEnabled(val)
        }
        return nil
    }
    res["reviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewReviewerScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewReviewerScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewReviewerScopeable)
                }
            }
            m.SetReviewers(res)
        }
        return nil
    }
    res["stageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStageId(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AccessReviewStageSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecommendationInsightSettings gets the recommendationInsightSettings property value. The recommendationInsightSettings property
// returns a []AccessReviewRecommendationInsightSettingable when successful
func (m *AccessReviewStageSettings) GetRecommendationInsightSettings()([]AccessReviewRecommendationInsightSettingable) {
    val, err := m.GetBackingStore().Get("recommendationInsightSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewRecommendationInsightSettingable)
    }
    return nil
}
// GetRecommendationsEnabled gets the recommendationsEnabled property value. Indicates whether showing recommendations to reviewers is enabled. Required. NOTE: The value of this property overrides override the corresponding setting on the accessReviewScheduleDefinition object.
// returns a *bool when successful
func (m *AccessReviewStageSettings) GetRecommendationsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("recommendationsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetReviewers gets the reviewers property value. Defines who the reviewers are. If none is specified, the review is a self-review (users review their own access).  For examples of options for assigning reviewers, see Assign reviewers to your access review definition using the Microsoft Graph API. NOTE: The value of this property overrides the corresponding setting on the accessReviewScheduleDefinition.
// returns a []AccessReviewReviewerScopeable when successful
func (m *AccessReviewStageSettings) GetReviewers()([]AccessReviewReviewerScopeable) {
    val, err := m.GetBackingStore().Get("reviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewReviewerScopeable)
    }
    return nil
}
// GetStageId gets the stageId property value. Unique identifier of the accessReviewStageSettings object. The stageId is used by the dependsOn property to indicate the order of the stages. Required.
// returns a *string when successful
func (m *AccessReviewStageSettings) GetStageId()(*string) {
    val, err := m.GetBackingStore().Get("stageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewStageSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDecisionsThatWillMoveToNextStage() != nil {
        err := writer.WriteCollectionOfStringValues("decisionsThatWillMoveToNextStage", m.GetDecisionsThatWillMoveToNextStage())
        if err != nil {
            return err
        }
    }
    if m.GetDependsOn() != nil {
        err := writer.WriteCollectionOfStringValues("dependsOn", m.GetDependsOn())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("durationInDays", m.GetDurationInDays())
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
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRecommendationInsightSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRecommendationInsightSettings()))
        for i, v := range m.GetRecommendationInsightSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("recommendationInsightSettings", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("recommendationsEnabled", m.GetRecommendationsEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReviewers()))
        for i, v := range m.GetReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("reviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("stageId", m.GetStageId())
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
func (m *AccessReviewStageSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessReviewStageSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDecisionsThatWillMoveToNextStage sets the decisionsThatWillMoveToNextStage property value. Indicate which decisions will go to the next stage. Can be a subset of Approve, Deny, Recommendation, or NotReviewed. If not provided, all decisions will go to the next stage. Optional.
func (m *AccessReviewStageSettings) SetDecisionsThatWillMoveToNextStage(value []string)() {
    err := m.GetBackingStore().Set("decisionsThatWillMoveToNextStage", value)
    if err != nil {
        panic(err)
    }
}
// SetDependsOn sets the dependsOn property value. Defines the sequential or parallel order of the stages and depends on the stageId. Only sequential stages are currently supported. For example, if stageId is 2, then dependsOn must be 1. If stageId is 1, don't specify dependsOn. Required if stageId isn't 1.
func (m *AccessReviewStageSettings) SetDependsOn(value []string)() {
    err := m.GetBackingStore().Set("dependsOn", value)
    if err != nil {
        panic(err)
    }
}
// SetDurationInDays sets the durationInDays property value. The duration of the stage. Required.  NOTE: The cumulative value of this property across all stages  1. Will override the instanceDurationInDays setting on the accessReviewScheduleDefinition object. 2. Can't exceed the length of one recurrence. That is, if the review recurs weekly, the cumulative durationInDays can't exceed 7.
func (m *AccessReviewStageSettings) SetDurationInDays(value *int32)() {
    err := m.GetBackingStore().Set("durationInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackReviewers sets the fallbackReviewers property value. If provided, the fallback reviewers are asked to complete a review if the primary reviewers don't exist. For example, if managers are selected as reviewers and a principal under review doesn't have a manager in Microsoft Entra ID, the fallback reviewers are asked to review that principal. NOTE: The value of this property overrides the corresponding setting on the accessReviewScheduleDefinition object.
func (m *AccessReviewStageSettings) SetFallbackReviewers(value []AccessReviewReviewerScopeable)() {
    err := m.GetBackingStore().Set("fallbackReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessReviewStageSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendationInsightSettings sets the recommendationInsightSettings property value. The recommendationInsightSettings property
func (m *AccessReviewStageSettings) SetRecommendationInsightSettings(value []AccessReviewRecommendationInsightSettingable)() {
    err := m.GetBackingStore().Set("recommendationInsightSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendationsEnabled sets the recommendationsEnabled property value. Indicates whether showing recommendations to reviewers is enabled. Required. NOTE: The value of this property overrides override the corresponding setting on the accessReviewScheduleDefinition object.
func (m *AccessReviewStageSettings) SetRecommendationsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("recommendationsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewers sets the reviewers property value. Defines who the reviewers are. If none is specified, the review is a self-review (users review their own access).  For examples of options for assigning reviewers, see Assign reviewers to your access review definition using the Microsoft Graph API. NOTE: The value of this property overrides the corresponding setting on the accessReviewScheduleDefinition.
func (m *AccessReviewStageSettings) SetReviewers(value []AccessReviewReviewerScopeable)() {
    err := m.GetBackingStore().Set("reviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetStageId sets the stageId property value. Unique identifier of the accessReviewStageSettings object. The stageId is used by the dependsOn property to indicate the order of the stages. Required.
func (m *AccessReviewStageSettings) SetStageId(value *string)() {
    err := m.GetBackingStore().Set("stageId", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewStageSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDecisionsThatWillMoveToNextStage()([]string)
    GetDependsOn()([]string)
    GetDurationInDays()(*int32)
    GetFallbackReviewers()([]AccessReviewReviewerScopeable)
    GetOdataType()(*string)
    GetRecommendationInsightSettings()([]AccessReviewRecommendationInsightSettingable)
    GetRecommendationsEnabled()(*bool)
    GetReviewers()([]AccessReviewReviewerScopeable)
    GetStageId()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDecisionsThatWillMoveToNextStage(value []string)()
    SetDependsOn(value []string)()
    SetDurationInDays(value *int32)()
    SetFallbackReviewers(value []AccessReviewReviewerScopeable)()
    SetOdataType(value *string)()
    SetRecommendationInsightSettings(value []AccessReviewRecommendationInsightSettingable)()
    SetRecommendationsEnabled(value *bool)()
    SetReviewers(value []AccessReviewReviewerScopeable)()
    SetStageId(value *string)()
}

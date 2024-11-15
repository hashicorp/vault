package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserLastSignInRecommendationInsightSetting struct {
    AccessReviewRecommendationInsightSetting
}
// NewUserLastSignInRecommendationInsightSetting instantiates a new UserLastSignInRecommendationInsightSetting and sets the default values.
func NewUserLastSignInRecommendationInsightSetting()(*UserLastSignInRecommendationInsightSetting) {
    m := &UserLastSignInRecommendationInsightSetting{
        AccessReviewRecommendationInsightSetting: *NewAccessReviewRecommendationInsightSetting(),
    }
    odataTypeValue := "#microsoft.graph.userLastSignInRecommendationInsightSetting"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUserLastSignInRecommendationInsightSettingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserLastSignInRecommendationInsightSettingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserLastSignInRecommendationInsightSetting(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserLastSignInRecommendationInsightSetting) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewRecommendationInsightSetting.GetFieldDeserializers()
    res["recommendationLookBackDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecommendationLookBackDuration(val)
        }
        return nil
    }
    res["signInScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserSignInRecommendationScope)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInScope(val.(*UserSignInRecommendationScope))
        }
        return nil
    }
    return res
}
// GetRecommendationLookBackDuration gets the recommendationLookBackDuration property value. Optional. Indicates the time period of inactivity (with respect to the start date of the review instance) that recommendations will be configured from. The recommendation will be to deny if the user is inactive during the look-back duration. For reviews of groups and Microsoft Entra roles, any duration is accepted. For reviews of applications, 30 days is the maximum duration. If not specified, the duration is 30 days.
// returns a *ISODuration when successful
func (m *UserLastSignInRecommendationInsightSetting) GetRecommendationLookBackDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("recommendationLookBackDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetSignInScope gets the signInScope property value. Indicates whether inactivity is calculated based on the user's inactivity in the tenant or in the application. The possible values are tenant, application, unknownFutureValue. application is only relevant when the access review is a review of an assignment to an application.
// returns a *UserSignInRecommendationScope when successful
func (m *UserLastSignInRecommendationInsightSetting) GetSignInScope()(*UserSignInRecommendationScope) {
    val, err := m.GetBackingStore().Get("signInScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserSignInRecommendationScope)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserLastSignInRecommendationInsightSetting) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewRecommendationInsightSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteISODurationValue("recommendationLookBackDuration", m.GetRecommendationLookBackDuration())
        if err != nil {
            return err
        }
    }
    if m.GetSignInScope() != nil {
        cast := (*m.GetSignInScope()).String()
        err = writer.WriteStringValue("signInScope", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRecommendationLookBackDuration sets the recommendationLookBackDuration property value. Optional. Indicates the time period of inactivity (with respect to the start date of the review instance) that recommendations will be configured from. The recommendation will be to deny if the user is inactive during the look-back duration. For reviews of groups and Microsoft Entra roles, any duration is accepted. For reviews of applications, 30 days is the maximum duration. If not specified, the duration is 30 days.
func (m *UserLastSignInRecommendationInsightSetting) SetRecommendationLookBackDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("recommendationLookBackDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInScope sets the signInScope property value. Indicates whether inactivity is calculated based on the user's inactivity in the tenant or in the application. The possible values are tenant, application, unknownFutureValue. application is only relevant when the access review is a review of an assignment to an application.
func (m *UserLastSignInRecommendationInsightSetting) SetSignInScope(value *UserSignInRecommendationScope)() {
    err := m.GetBackingStore().Set("signInScope", value)
    if err != nil {
        panic(err)
    }
}
type UserLastSignInRecommendationInsightSettingable interface {
    AccessReviewRecommendationInsightSettingable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRecommendationLookBackDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetSignInScope()(*UserSignInRecommendationScope)
    SetRecommendationLookBackDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetSignInScope(value *UserSignInRecommendationScope)()
}

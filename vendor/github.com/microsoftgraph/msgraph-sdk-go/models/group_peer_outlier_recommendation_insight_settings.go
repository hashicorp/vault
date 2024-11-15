package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type GroupPeerOutlierRecommendationInsightSettings struct {
    AccessReviewRecommendationInsightSetting
}
// NewGroupPeerOutlierRecommendationInsightSettings instantiates a new GroupPeerOutlierRecommendationInsightSettings and sets the default values.
func NewGroupPeerOutlierRecommendationInsightSettings()(*GroupPeerOutlierRecommendationInsightSettings) {
    m := &GroupPeerOutlierRecommendationInsightSettings{
        AccessReviewRecommendationInsightSetting: *NewAccessReviewRecommendationInsightSetting(),
    }
    odataTypeValue := "#microsoft.graph.groupPeerOutlierRecommendationInsightSettings"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateGroupPeerOutlierRecommendationInsightSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGroupPeerOutlierRecommendationInsightSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGroupPeerOutlierRecommendationInsightSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *GroupPeerOutlierRecommendationInsightSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewRecommendationInsightSetting.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *GroupPeerOutlierRecommendationInsightSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewRecommendationInsightSetting.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type GroupPeerOutlierRecommendationInsightSettingsable interface {
    AccessReviewRecommendationInsightSettingable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FeatureRolloutPolicy struct {
    Entity
}
// NewFeatureRolloutPolicy instantiates a new FeatureRolloutPolicy and sets the default values.
func NewFeatureRolloutPolicy()(*FeatureRolloutPolicy) {
    m := &FeatureRolloutPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateFeatureRolloutPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFeatureRolloutPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFeatureRolloutPolicy(), nil
}
// GetAppliesTo gets the appliesTo property value. Nullable. Specifies a list of directoryObject resources that feature is enabled for.
// returns a []DirectoryObjectable when successful
func (m *FeatureRolloutPolicy) GetAppliesTo()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("appliesTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetDescription gets the description property value. A description for this feature rollout policy.
// returns a *string when successful
func (m *FeatureRolloutPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for this  feature rollout policy.
// returns a *string when successful
func (m *FeatureRolloutPolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFeature gets the feature property value. The feature property
// returns a *StagedFeatureName when successful
func (m *FeatureRolloutPolicy) GetFeature()(*StagedFeatureName) {
    val, err := m.GetBackingStore().Get("feature")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*StagedFeatureName)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FeatureRolloutPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appliesTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetAppliesTo(res)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["feature"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStagedFeatureName)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeature(val.(*StagedFeatureName))
        }
        return nil
    }
    res["isAppliedToOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAppliedToOrganization(val)
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
    return res
}
// GetIsAppliedToOrganization gets the isAppliedToOrganization property value. Indicates whether this feature rollout policy should be applied to the entire organization.
// returns a *bool when successful
func (m *FeatureRolloutPolicy) GetIsAppliedToOrganization()(*bool) {
    val, err := m.GetBackingStore().Get("isAppliedToOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsEnabled gets the isEnabled property value. Indicates whether the feature rollout is enabled.
// returns a *bool when successful
func (m *FeatureRolloutPolicy) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FeatureRolloutPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppliesTo() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppliesTo()))
        for i, v := range m.GetAppliesTo() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appliesTo", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetFeature() != nil {
        cast := (*m.GetFeature()).String()
        err = writer.WriteStringValue("feature", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAppliedToOrganization", m.GetIsAppliedToOrganization())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppliesTo sets the appliesTo property value. Nullable. Specifies a list of directoryObject resources that feature is enabled for.
func (m *FeatureRolloutPolicy) SetAppliesTo(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("appliesTo", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. A description for this feature rollout policy.
func (m *FeatureRolloutPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for this  feature rollout policy.
func (m *FeatureRolloutPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetFeature sets the feature property value. The feature property
func (m *FeatureRolloutPolicy) SetFeature(value *StagedFeatureName)() {
    err := m.GetBackingStore().Set("feature", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAppliedToOrganization sets the isAppliedToOrganization property value. Indicates whether this feature rollout policy should be applied to the entire organization.
func (m *FeatureRolloutPolicy) SetIsAppliedToOrganization(value *bool)() {
    err := m.GetBackingStore().Set("isAppliedToOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Indicates whether the feature rollout is enabled.
func (m *FeatureRolloutPolicy) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
type FeatureRolloutPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppliesTo()([]DirectoryObjectable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetFeature()(*StagedFeatureName)
    GetIsAppliedToOrganization()(*bool)
    GetIsEnabled()(*bool)
    SetAppliesTo(value []DirectoryObjectable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetFeature(value *StagedFeatureName)()
    SetIsAppliedToOrganization(value *bool)()
    SetIsEnabled(value *bool)()
}

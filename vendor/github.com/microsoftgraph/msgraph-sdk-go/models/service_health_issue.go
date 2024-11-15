package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServiceHealthIssue struct {
    ServiceAnnouncementBase
}
// NewServiceHealthIssue instantiates a new ServiceHealthIssue and sets the default values.
func NewServiceHealthIssue()(*ServiceHealthIssue) {
    m := &ServiceHealthIssue{
        ServiceAnnouncementBase: *NewServiceAnnouncementBase(),
    }
    odataTypeValue := "#microsoft.graph.serviceHealthIssue"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateServiceHealthIssueFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceHealthIssueFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceHealthIssue(), nil
}
// GetClassification gets the classification property value. The classification property
// returns a *ServiceHealthClassificationType when successful
func (m *ServiceHealthIssue) GetClassification()(*ServiceHealthClassificationType) {
    val, err := m.GetBackingStore().Get("classification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceHealthClassificationType)
    }
    return nil
}
// GetFeature gets the feature property value. The feature name of the service issue.
// returns a *string when successful
func (m *ServiceHealthIssue) GetFeature()(*string) {
    val, err := m.GetBackingStore().Get("feature")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFeatureGroup gets the featureGroup property value. The feature group name of the service issue.
// returns a *string when successful
func (m *ServiceHealthIssue) GetFeatureGroup()(*string) {
    val, err := m.GetBackingStore().Get("featureGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServiceHealthIssue) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ServiceAnnouncementBase.GetFieldDeserializers()
    res["classification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceHealthClassificationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassification(val.(*ServiceHealthClassificationType))
        }
        return nil
    }
    res["feature"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeature(val)
        }
        return nil
    }
    res["featureGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureGroup(val)
        }
        return nil
    }
    res["impactDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImpactDescription(val)
        }
        return nil
    }
    res["isResolved"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsResolved(val)
        }
        return nil
    }
    res["origin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceHealthOrigin)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrigin(val.(*ServiceHealthOrigin))
        }
        return nil
    }
    res["posts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceHealthIssuePostFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceHealthIssuePostable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceHealthIssuePostable)
                }
            }
            m.SetPosts(res)
        }
        return nil
    }
    res["service"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetService(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceHealthStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*ServiceHealthStatus))
        }
        return nil
    }
    return res
}
// GetImpactDescription gets the impactDescription property value. The description of the service issue impact.
// returns a *string when successful
func (m *ServiceHealthIssue) GetImpactDescription()(*string) {
    val, err := m.GetBackingStore().Get("impactDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsResolved gets the isResolved property value. Indicates whether the issue is resolved.
// returns a *bool when successful
func (m *ServiceHealthIssue) GetIsResolved()(*bool) {
    val, err := m.GetBackingStore().Get("isResolved")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOrigin gets the origin property value. The origin property
// returns a *ServiceHealthOrigin when successful
func (m *ServiceHealthIssue) GetOrigin()(*ServiceHealthOrigin) {
    val, err := m.GetBackingStore().Get("origin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceHealthOrigin)
    }
    return nil
}
// GetPosts gets the posts property value. Collection of historical posts for the service issue.
// returns a []ServiceHealthIssuePostable when successful
func (m *ServiceHealthIssue) GetPosts()([]ServiceHealthIssuePostable) {
    val, err := m.GetBackingStore().Get("posts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceHealthIssuePostable)
    }
    return nil
}
// GetService gets the service property value. Indicates the service affected by the issue.
// returns a *string when successful
func (m *ServiceHealthIssue) GetService()(*string) {
    val, err := m.GetBackingStore().Get("service")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *ServiceHealthStatus when successful
func (m *ServiceHealthIssue) GetStatus()(*ServiceHealthStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceHealthStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServiceHealthIssue) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ServiceAnnouncementBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClassification() != nil {
        cast := (*m.GetClassification()).String()
        err = writer.WriteStringValue("classification", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("feature", m.GetFeature())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("featureGroup", m.GetFeatureGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("impactDescription", m.GetImpactDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isResolved", m.GetIsResolved())
        if err != nil {
            return err
        }
    }
    if m.GetOrigin() != nil {
        cast := (*m.GetOrigin()).String()
        err = writer.WriteStringValue("origin", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPosts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPosts()))
        for i, v := range m.GetPosts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("posts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("service", m.GetService())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClassification sets the classification property value. The classification property
func (m *ServiceHealthIssue) SetClassification(value *ServiceHealthClassificationType)() {
    err := m.GetBackingStore().Set("classification", value)
    if err != nil {
        panic(err)
    }
}
// SetFeature sets the feature property value. The feature name of the service issue.
func (m *ServiceHealthIssue) SetFeature(value *string)() {
    err := m.GetBackingStore().Set("feature", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureGroup sets the featureGroup property value. The feature group name of the service issue.
func (m *ServiceHealthIssue) SetFeatureGroup(value *string)() {
    err := m.GetBackingStore().Set("featureGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetImpactDescription sets the impactDescription property value. The description of the service issue impact.
func (m *ServiceHealthIssue) SetImpactDescription(value *string)() {
    err := m.GetBackingStore().Set("impactDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetIsResolved sets the isResolved property value. Indicates whether the issue is resolved.
func (m *ServiceHealthIssue) SetIsResolved(value *bool)() {
    err := m.GetBackingStore().Set("isResolved", value)
    if err != nil {
        panic(err)
    }
}
// SetOrigin sets the origin property value. The origin property
func (m *ServiceHealthIssue) SetOrigin(value *ServiceHealthOrigin)() {
    err := m.GetBackingStore().Set("origin", value)
    if err != nil {
        panic(err)
    }
}
// SetPosts sets the posts property value. Collection of historical posts for the service issue.
func (m *ServiceHealthIssue) SetPosts(value []ServiceHealthIssuePostable)() {
    err := m.GetBackingStore().Set("posts", value)
    if err != nil {
        panic(err)
    }
}
// SetService sets the service property value. Indicates the service affected by the issue.
func (m *ServiceHealthIssue) SetService(value *string)() {
    err := m.GetBackingStore().Set("service", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *ServiceHealthIssue) SetStatus(value *ServiceHealthStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type ServiceHealthIssueable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ServiceAnnouncementBaseable
    GetClassification()(*ServiceHealthClassificationType)
    GetFeature()(*string)
    GetFeatureGroup()(*string)
    GetImpactDescription()(*string)
    GetIsResolved()(*bool)
    GetOrigin()(*ServiceHealthOrigin)
    GetPosts()([]ServiceHealthIssuePostable)
    GetService()(*string)
    GetStatus()(*ServiceHealthStatus)
    SetClassification(value *ServiceHealthClassificationType)()
    SetFeature(value *string)()
    SetFeatureGroup(value *string)()
    SetImpactDescription(value *string)()
    SetIsResolved(value *bool)()
    SetOrigin(value *ServiceHealthOrigin)()
    SetPosts(value []ServiceHealthIssuePostable)()
    SetService(value *string)()
    SetStatus(value *ServiceHealthStatus)()
}

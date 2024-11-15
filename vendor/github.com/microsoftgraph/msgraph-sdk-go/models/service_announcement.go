package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServiceAnnouncement struct {
    Entity
}
// NewServiceAnnouncement instantiates a new ServiceAnnouncement and sets the default values.
func NewServiceAnnouncement()(*ServiceAnnouncement) {
    m := &ServiceAnnouncement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateServiceAnnouncementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceAnnouncementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceAnnouncement(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServiceAnnouncement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["healthOverviews"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceHealthFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceHealthable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceHealthable)
                }
            }
            m.SetHealthOverviews(res)
        }
        return nil
    }
    res["issues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceHealthIssueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceHealthIssueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceHealthIssueable)
                }
            }
            m.SetIssues(res)
        }
        return nil
    }
    res["messages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceUpdateMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceUpdateMessageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceUpdateMessageable)
                }
            }
            m.SetMessages(res)
        }
        return nil
    }
    return res
}
// GetHealthOverviews gets the healthOverviews property value. A collection of service health information for tenant. This property is a contained navigation property, it is nullable and readonly.
// returns a []ServiceHealthable when successful
func (m *ServiceAnnouncement) GetHealthOverviews()([]ServiceHealthable) {
    val, err := m.GetBackingStore().Get("healthOverviews")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceHealthable)
    }
    return nil
}
// GetIssues gets the issues property value. A collection of service issues for tenant. This property is a contained navigation property, it is nullable and readonly.
// returns a []ServiceHealthIssueable when successful
func (m *ServiceAnnouncement) GetIssues()([]ServiceHealthIssueable) {
    val, err := m.GetBackingStore().Get("issues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceHealthIssueable)
    }
    return nil
}
// GetMessages gets the messages property value. A collection of service messages for tenant. This property is a contained navigation property, it is nullable and readonly.
// returns a []ServiceUpdateMessageable when successful
func (m *ServiceAnnouncement) GetMessages()([]ServiceUpdateMessageable) {
    val, err := m.GetBackingStore().Get("messages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceUpdateMessageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServiceAnnouncement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetHealthOverviews() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHealthOverviews()))
        for i, v := range m.GetHealthOverviews() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("healthOverviews", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIssues() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIssues()))
        for i, v := range m.GetIssues() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("issues", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessages()))
        for i, v := range m.GetMessages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messages", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHealthOverviews sets the healthOverviews property value. A collection of service health information for tenant. This property is a contained navigation property, it is nullable and readonly.
func (m *ServiceAnnouncement) SetHealthOverviews(value []ServiceHealthable)() {
    err := m.GetBackingStore().Set("healthOverviews", value)
    if err != nil {
        panic(err)
    }
}
// SetIssues sets the issues property value. A collection of service issues for tenant. This property is a contained navigation property, it is nullable and readonly.
func (m *ServiceAnnouncement) SetIssues(value []ServiceHealthIssueable)() {
    err := m.GetBackingStore().Set("issues", value)
    if err != nil {
        panic(err)
    }
}
// SetMessages sets the messages property value. A collection of service messages for tenant. This property is a contained navigation property, it is nullable and readonly.
func (m *ServiceAnnouncement) SetMessages(value []ServiceUpdateMessageable)() {
    err := m.GetBackingStore().Set("messages", value)
    if err != nil {
        panic(err)
    }
}
type ServiceAnnouncementable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHealthOverviews()([]ServiceHealthable)
    GetIssues()([]ServiceHealthIssueable)
    GetMessages()([]ServiceUpdateMessageable)
    SetHealthOverviews(value []ServiceHealthable)()
    SetIssues(value []ServiceHealthIssueable)()
    SetMessages(value []ServiceUpdateMessageable)()
}

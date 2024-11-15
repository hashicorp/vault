package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServiceHealth struct {
    Entity
}
// NewServiceHealth instantiates a new ServiceHealth and sets the default values.
func NewServiceHealth()(*ServiceHealth) {
    m := &ServiceHealth{
        Entity: *NewEntity(),
    }
    return m
}
// CreateServiceHealthFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceHealthFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceHealth(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServiceHealth) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
// GetIssues gets the issues property value. A collection of issues that happened on the service, with detailed information for each issue.
// returns a []ServiceHealthIssueable when successful
func (m *ServiceHealth) GetIssues()([]ServiceHealthIssueable) {
    val, err := m.GetBackingStore().Get("issues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceHealthIssueable)
    }
    return nil
}
// GetService gets the service property value. The service name. Use the list healthOverviews operation to get exact string names for services subscribed by the tenant.
// returns a *string when successful
func (m *ServiceHealth) GetService()(*string) {
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
func (m *ServiceHealth) GetStatus()(*ServiceHealthStatus) {
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
func (m *ServiceHealth) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
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
// SetIssues sets the issues property value. A collection of issues that happened on the service, with detailed information for each issue.
func (m *ServiceHealth) SetIssues(value []ServiceHealthIssueable)() {
    err := m.GetBackingStore().Set("issues", value)
    if err != nil {
        panic(err)
    }
}
// SetService sets the service property value. The service name. Use the list healthOverviews operation to get exact string names for services subscribed by the tenant.
func (m *ServiceHealth) SetService(value *string)() {
    err := m.GetBackingStore().Set("service", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *ServiceHealth) SetStatus(value *ServiceHealthStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type ServiceHealthable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIssues()([]ServiceHealthIssueable)
    GetService()(*string)
    GetStatus()(*ServiceHealthStatus)
    SetIssues(value []ServiceHealthIssueable)()
    SetService(value *string)()
    SetStatus(value *ServiceHealthStatus)()
}

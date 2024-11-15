package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DelegatedAdminServiceManagementDetail struct {
    Entity
}
// NewDelegatedAdminServiceManagementDetail instantiates a new DelegatedAdminServiceManagementDetail and sets the default values.
func NewDelegatedAdminServiceManagementDetail()(*DelegatedAdminServiceManagementDetail) {
    m := &DelegatedAdminServiceManagementDetail{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDelegatedAdminServiceManagementDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDelegatedAdminServiceManagementDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDelegatedAdminServiceManagementDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DelegatedAdminServiceManagementDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["serviceManagementUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceManagementUrl(val)
        }
        return nil
    }
    res["serviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceName(val)
        }
        return nil
    }
    return res
}
// GetServiceManagementUrl gets the serviceManagementUrl property value. The URL of the management portal for the managed service. Read-only.
// returns a *string when successful
func (m *DelegatedAdminServiceManagementDetail) GetServiceManagementUrl()(*string) {
    val, err := m.GetBackingStore().Get("serviceManagementUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServiceName gets the serviceName property value. The name of a managed service. Read-only.
// returns a *string when successful
func (m *DelegatedAdminServiceManagementDetail) GetServiceName()(*string) {
    val, err := m.GetBackingStore().Get("serviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DelegatedAdminServiceManagementDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("serviceManagementUrl", m.GetServiceManagementUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serviceName", m.GetServiceName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetServiceManagementUrl sets the serviceManagementUrl property value. The URL of the management portal for the managed service. Read-only.
func (m *DelegatedAdminServiceManagementDetail) SetServiceManagementUrl(value *string)() {
    err := m.GetBackingStore().Set("serviceManagementUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceName sets the serviceName property value. The name of a managed service. Read-only.
func (m *DelegatedAdminServiceManagementDetail) SetServiceName(value *string)() {
    err := m.GetBackingStore().Set("serviceName", value)
    if err != nil {
        panic(err)
    }
}
type DelegatedAdminServiceManagementDetailable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetServiceManagementUrl()(*string)
    GetServiceName()(*string)
    SetServiceManagementUrl(value *string)()
    SetServiceName(value *string)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RemoteDesktopSecurityConfiguration struct {
    Entity
}
// NewRemoteDesktopSecurityConfiguration instantiates a new RemoteDesktopSecurityConfiguration and sets the default values.
func NewRemoteDesktopSecurityConfiguration()(*RemoteDesktopSecurityConfiguration) {
    m := &RemoteDesktopSecurityConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateRemoteDesktopSecurityConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRemoteDesktopSecurityConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRemoteDesktopSecurityConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RemoteDesktopSecurityConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["isRemoteDesktopProtocolEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRemoteDesktopProtocolEnabled(val)
        }
        return nil
    }
    res["targetDeviceGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTargetDeviceGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TargetDeviceGroupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TargetDeviceGroupable)
                }
            }
            m.SetTargetDeviceGroups(res)
        }
        return nil
    }
    return res
}
// GetIsRemoteDesktopProtocolEnabled gets the isRemoteDesktopProtocolEnabled property value. Determines if Microsoft Entra ID RDS authentication protocol for RDP is enabled.
// returns a *bool when successful
func (m *RemoteDesktopSecurityConfiguration) GetIsRemoteDesktopProtocolEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isRemoteDesktopProtocolEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTargetDeviceGroups gets the targetDeviceGroups property value. The collection of target device groups that are associated with the RDS security configuration that will be enabled for SSO when a client connects to the target device over RDP using the new Microsoft Entra ID RDS authentication protocol.
// returns a []TargetDeviceGroupable when successful
func (m *RemoteDesktopSecurityConfiguration) GetTargetDeviceGroups()([]TargetDeviceGroupable) {
    val, err := m.GetBackingStore().Get("targetDeviceGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TargetDeviceGroupable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RemoteDesktopSecurityConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isRemoteDesktopProtocolEnabled", m.GetIsRemoteDesktopProtocolEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetTargetDeviceGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargetDeviceGroups()))
        for i, v := range m.GetTargetDeviceGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("targetDeviceGroups", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsRemoteDesktopProtocolEnabled sets the isRemoteDesktopProtocolEnabled property value. Determines if Microsoft Entra ID RDS authentication protocol for RDP is enabled.
func (m *RemoteDesktopSecurityConfiguration) SetIsRemoteDesktopProtocolEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isRemoteDesktopProtocolEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetDeviceGroups sets the targetDeviceGroups property value. The collection of target device groups that are associated with the RDS security configuration that will be enabled for SSO when a client connects to the target device over RDP using the new Microsoft Entra ID RDS authentication protocol.
func (m *RemoteDesktopSecurityConfiguration) SetTargetDeviceGroups(value []TargetDeviceGroupable)() {
    err := m.GetBackingStore().Set("targetDeviceGroups", value)
    if err != nil {
        panic(err)
    }
}
type RemoteDesktopSecurityConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsRemoteDesktopProtocolEnabled()(*bool)
    GetTargetDeviceGroups()([]TargetDeviceGroupable)
    SetIsRemoteDesktopProtocolEnabled(value *bool)()
    SetTargetDeviceGroups(value []TargetDeviceGroupable)()
}

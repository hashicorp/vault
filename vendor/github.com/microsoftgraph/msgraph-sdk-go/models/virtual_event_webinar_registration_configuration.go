package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventWebinarRegistrationConfiguration struct {
    VirtualEventRegistrationConfiguration
}
// NewVirtualEventWebinarRegistrationConfiguration instantiates a new VirtualEventWebinarRegistrationConfiguration and sets the default values.
func NewVirtualEventWebinarRegistrationConfiguration()(*VirtualEventWebinarRegistrationConfiguration) {
    m := &VirtualEventWebinarRegistrationConfiguration{
        VirtualEventRegistrationConfiguration: *NewVirtualEventRegistrationConfiguration(),
    }
    return m
}
// CreateVirtualEventWebinarRegistrationConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventWebinarRegistrationConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventWebinarRegistrationConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventWebinarRegistrationConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.VirtualEventRegistrationConfiguration.GetFieldDeserializers()
    res["isManualApprovalEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsManualApprovalEnabled(val)
        }
        return nil
    }
    res["isWaitlistEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsWaitlistEnabled(val)
        }
        return nil
    }
    return res
}
// GetIsManualApprovalEnabled gets the isManualApprovalEnabled property value. The isManualApprovalEnabled property
// returns a *bool when successful
func (m *VirtualEventWebinarRegistrationConfiguration) GetIsManualApprovalEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isManualApprovalEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsWaitlistEnabled gets the isWaitlistEnabled property value. The isWaitlistEnabled property
// returns a *bool when successful
func (m *VirtualEventWebinarRegistrationConfiguration) GetIsWaitlistEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isWaitlistEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventWebinarRegistrationConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.VirtualEventRegistrationConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isManualApprovalEnabled", m.GetIsManualApprovalEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isWaitlistEnabled", m.GetIsWaitlistEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsManualApprovalEnabled sets the isManualApprovalEnabled property value. The isManualApprovalEnabled property
func (m *VirtualEventWebinarRegistrationConfiguration) SetIsManualApprovalEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isManualApprovalEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsWaitlistEnabled sets the isWaitlistEnabled property value. The isWaitlistEnabled property
func (m *VirtualEventWebinarRegistrationConfiguration) SetIsWaitlistEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isWaitlistEnabled", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventWebinarRegistrationConfigurationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VirtualEventRegistrationConfigurationable
    GetIsManualApprovalEnabled()(*bool)
    GetIsWaitlistEnabled()(*bool)
    SetIsManualApprovalEnabled(value *bool)()
    SetIsWaitlistEnabled(value *bool)()
}

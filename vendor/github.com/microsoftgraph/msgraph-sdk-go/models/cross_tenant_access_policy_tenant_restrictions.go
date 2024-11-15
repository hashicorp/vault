package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CrossTenantAccessPolicyTenantRestrictions struct {
    CrossTenantAccessPolicyB2BSetting
}
// NewCrossTenantAccessPolicyTenantRestrictions instantiates a new CrossTenantAccessPolicyTenantRestrictions and sets the default values.
func NewCrossTenantAccessPolicyTenantRestrictions()(*CrossTenantAccessPolicyTenantRestrictions) {
    m := &CrossTenantAccessPolicyTenantRestrictions{
        CrossTenantAccessPolicyB2BSetting: *NewCrossTenantAccessPolicyB2BSetting(),
    }
    odataTypeValue := "#microsoft.graph.crossTenantAccessPolicyTenantRestrictions"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCrossTenantAccessPolicyTenantRestrictionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCrossTenantAccessPolicyTenantRestrictionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCrossTenantAccessPolicyTenantRestrictions(), nil
}
// GetDevices gets the devices property value. Defines the rule for filtering devices and whether devices that satisfy the rule should be allowed or blocked. This property isn't supported on the server side yet.
// returns a DevicesFilterable when successful
func (m *CrossTenantAccessPolicyTenantRestrictions) GetDevices()(DevicesFilterable) {
    val, err := m.GetBackingStore().Get("devices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DevicesFilterable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CrossTenantAccessPolicyTenantRestrictions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CrossTenantAccessPolicyB2BSetting.GetFieldDeserializers()
    res["devices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDevicesFilterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevices(val.(DevicesFilterable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CrossTenantAccessPolicyTenantRestrictions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CrossTenantAccessPolicyB2BSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("devices", m.GetDevices())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDevices sets the devices property value. Defines the rule for filtering devices and whether devices that satisfy the rule should be allowed or blocked. This property isn't supported on the server side yet.
func (m *CrossTenantAccessPolicyTenantRestrictions) SetDevices(value DevicesFilterable)() {
    err := m.GetBackingStore().Set("devices", value)
    if err != nil {
        panic(err)
    }
}
type CrossTenantAccessPolicyTenantRestrictionsable interface {
    CrossTenantAccessPolicyB2BSettingable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDevices()(DevicesFilterable)
    SetDevices(value DevicesFilterable)()
}

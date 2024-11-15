package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceRegistrationPolicy struct {
    Entity
}
// NewDeviceRegistrationPolicy instantiates a new DeviceRegistrationPolicy and sets the default values.
func NewDeviceRegistrationPolicy()(*DeviceRegistrationPolicy) {
    m := &DeviceRegistrationPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceRegistrationPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceRegistrationPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceRegistrationPolicy(), nil
}
// GetAzureADJoin gets the azureADJoin property value. The azureADJoin property
// returns a AzureADJoinPolicyable when successful
func (m *DeviceRegistrationPolicy) GetAzureADJoin()(AzureADJoinPolicyable) {
    val, err := m.GetBackingStore().Get("azureADJoin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AzureADJoinPolicyable)
    }
    return nil
}
// GetAzureADRegistration gets the azureADRegistration property value. The azureADRegistration property
// returns a AzureADRegistrationPolicyable when successful
func (m *DeviceRegistrationPolicy) GetAzureADRegistration()(AzureADRegistrationPolicyable) {
    val, err := m.GetBackingStore().Get("azureADRegistration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AzureADRegistrationPolicyable)
    }
    return nil
}
// GetDescription gets the description property value. The description property
// returns a *string when successful
func (m *DeviceRegistrationPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The displayName property
// returns a *string when successful
func (m *DeviceRegistrationPolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *DeviceRegistrationPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["azureADJoin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAzureADJoinPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureADJoin(val.(AzureADJoinPolicyable))
        }
        return nil
    }
    res["azureADRegistration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAzureADRegistrationPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureADRegistration(val.(AzureADRegistrationPolicyable))
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
    res["localAdminPassword"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocalAdminPasswordSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocalAdminPassword(val.(LocalAdminPasswordSettingsable))
        }
        return nil
    }
    res["multiFactorAuthConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMultiFactorAuthConfiguration)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultiFactorAuthConfiguration(val.(*MultiFactorAuthConfiguration))
        }
        return nil
    }
    res["userDeviceQuota"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserDeviceQuota(val)
        }
        return nil
    }
    return res
}
// GetLocalAdminPassword gets the localAdminPassword property value. The localAdminPassword property
// returns a LocalAdminPasswordSettingsable when successful
func (m *DeviceRegistrationPolicy) GetLocalAdminPassword()(LocalAdminPasswordSettingsable) {
    val, err := m.GetBackingStore().Get("localAdminPassword")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LocalAdminPasswordSettingsable)
    }
    return nil
}
// GetMultiFactorAuthConfiguration gets the multiFactorAuthConfiguration property value. The multiFactorAuthConfiguration property
// returns a *MultiFactorAuthConfiguration when successful
func (m *DeviceRegistrationPolicy) GetMultiFactorAuthConfiguration()(*MultiFactorAuthConfiguration) {
    val, err := m.GetBackingStore().Get("multiFactorAuthConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MultiFactorAuthConfiguration)
    }
    return nil
}
// GetUserDeviceQuota gets the userDeviceQuota property value. The userDeviceQuota property
// returns a *int32 when successful
func (m *DeviceRegistrationPolicy) GetUserDeviceQuota()(*int32) {
    val, err := m.GetBackingStore().Get("userDeviceQuota")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceRegistrationPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("azureADJoin", m.GetAzureADJoin())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("azureADRegistration", m.GetAzureADRegistration())
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
    {
        err = writer.WriteObjectValue("localAdminPassword", m.GetLocalAdminPassword())
        if err != nil {
            return err
        }
    }
    if m.GetMultiFactorAuthConfiguration() != nil {
        cast := (*m.GetMultiFactorAuthConfiguration()).String()
        err = writer.WriteStringValue("multiFactorAuthConfiguration", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("userDeviceQuota", m.GetUserDeviceQuota())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureADJoin sets the azureADJoin property value. The azureADJoin property
func (m *DeviceRegistrationPolicy) SetAzureADJoin(value AzureADJoinPolicyable)() {
    err := m.GetBackingStore().Set("azureADJoin", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureADRegistration sets the azureADRegistration property value. The azureADRegistration property
func (m *DeviceRegistrationPolicy) SetAzureADRegistration(value AzureADRegistrationPolicyable)() {
    err := m.GetBackingStore().Set("azureADRegistration", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description property
func (m *DeviceRegistrationPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The displayName property
func (m *DeviceRegistrationPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLocalAdminPassword sets the localAdminPassword property value. The localAdminPassword property
func (m *DeviceRegistrationPolicy) SetLocalAdminPassword(value LocalAdminPasswordSettingsable)() {
    err := m.GetBackingStore().Set("localAdminPassword", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiFactorAuthConfiguration sets the multiFactorAuthConfiguration property value. The multiFactorAuthConfiguration property
func (m *DeviceRegistrationPolicy) SetMultiFactorAuthConfiguration(value *MultiFactorAuthConfiguration)() {
    err := m.GetBackingStore().Set("multiFactorAuthConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetUserDeviceQuota sets the userDeviceQuota property value. The userDeviceQuota property
func (m *DeviceRegistrationPolicy) SetUserDeviceQuota(value *int32)() {
    err := m.GetBackingStore().Set("userDeviceQuota", value)
    if err != nil {
        panic(err)
    }
}
type DeviceRegistrationPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureADJoin()(AzureADJoinPolicyable)
    GetAzureADRegistration()(AzureADRegistrationPolicyable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetLocalAdminPassword()(LocalAdminPasswordSettingsable)
    GetMultiFactorAuthConfiguration()(*MultiFactorAuthConfiguration)
    GetUserDeviceQuota()(*int32)
    SetAzureADJoin(value AzureADJoinPolicyable)()
    SetAzureADRegistration(value AzureADRegistrationPolicyable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetLocalAdminPassword(value LocalAdminPasswordSettingsable)()
    SetMultiFactorAuthConfiguration(value *MultiFactorAuthConfiguration)()
    SetUserDeviceQuota(value *int32)()
}

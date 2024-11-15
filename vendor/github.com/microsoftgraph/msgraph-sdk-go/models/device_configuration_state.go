package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceConfigurationState device Configuration State for a given device.
type DeviceConfigurationState struct {
    Entity
}
// NewDeviceConfigurationState instantiates a new DeviceConfigurationState and sets the default values.
func NewDeviceConfigurationState()(*DeviceConfigurationState) {
    m := &DeviceConfigurationState{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceConfigurationStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceConfigurationStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceConfigurationState(), nil
}
// GetDisplayName gets the displayName property value. The name of the policy for this policyBase
// returns a *string when successful
func (m *DeviceConfigurationState) GetDisplayName()(*string) {
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
func (m *DeviceConfigurationState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["platformType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePolicyPlatformType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlatformType(val.(*PolicyPlatformType))
        }
        return nil
    }
    res["settingCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingCount(val)
        }
        return nil
    }
    res["settingStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceConfigurationSettingStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceConfigurationSettingStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceConfigurationSettingStateable)
                }
            }
            m.SetSettingStates(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseComplianceStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*ComplianceStatus))
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetPlatformType gets the platformType property value. Supported platform types for policies.
// returns a *PolicyPlatformType when successful
func (m *DeviceConfigurationState) GetPlatformType()(*PolicyPlatformType) {
    val, err := m.GetBackingStore().Get("platformType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PolicyPlatformType)
    }
    return nil
}
// GetSettingCount gets the settingCount property value. Count of how many setting a policy holds
// returns a *int32 when successful
func (m *DeviceConfigurationState) GetSettingCount()(*int32) {
    val, err := m.GetBackingStore().Get("settingCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSettingStates gets the settingStates property value. The settingStates property
// returns a []DeviceConfigurationSettingStateable when successful
func (m *DeviceConfigurationState) GetSettingStates()([]DeviceConfigurationSettingStateable) {
    val, err := m.GetBackingStore().Get("settingStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceConfigurationSettingStateable)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *ComplianceStatus when successful
func (m *DeviceConfigurationState) GetState()(*ComplianceStatus) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ComplianceStatus)
    }
    return nil
}
// GetVersion gets the version property value. The version of the policy
// returns a *int32 when successful
func (m *DeviceConfigurationState) GetVersion()(*int32) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceConfigurationState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetPlatformType() != nil {
        cast := (*m.GetPlatformType()).String()
        err = writer.WriteStringValue("platformType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("settingCount", m.GetSettingCount())
        if err != nil {
            return err
        }
    }
    if m.GetSettingStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSettingStates()))
        for i, v := range m.GetSettingStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("settingStates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The name of the policy for this policyBase
func (m *DeviceConfigurationState) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatformType sets the platformType property value. Supported platform types for policies.
func (m *DeviceConfigurationState) SetPlatformType(value *PolicyPlatformType)() {
    err := m.GetBackingStore().Set("platformType", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingCount sets the settingCount property value. Count of how many setting a policy holds
func (m *DeviceConfigurationState) SetSettingCount(value *int32)() {
    err := m.GetBackingStore().Set("settingCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingStates sets the settingStates property value. The settingStates property
func (m *DeviceConfigurationState) SetSettingStates(value []DeviceConfigurationSettingStateable)() {
    err := m.GetBackingStore().Set("settingStates", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *DeviceConfigurationState) SetState(value *ComplianceStatus)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The version of the policy
func (m *DeviceConfigurationState) SetVersion(value *int32)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type DeviceConfigurationStateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetPlatformType()(*PolicyPlatformType)
    GetSettingCount()(*int32)
    GetSettingStates()([]DeviceConfigurationSettingStateable)
    GetState()(*ComplianceStatus)
    GetVersion()(*int32)
    SetDisplayName(value *string)()
    SetPlatformType(value *PolicyPlatformType)()
    SetSettingCount(value *int32)()
    SetSettingStates(value []DeviceConfigurationSettingStateable)()
    SetState(value *ComplianceStatus)()
    SetVersion(value *int32)()
}

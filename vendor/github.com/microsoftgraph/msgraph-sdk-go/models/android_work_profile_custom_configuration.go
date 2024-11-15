package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AndroidWorkProfileCustomConfiguration android Work Profile custom configuration
type AndroidWorkProfileCustomConfiguration struct {
    DeviceConfiguration
}
// NewAndroidWorkProfileCustomConfiguration instantiates a new AndroidWorkProfileCustomConfiguration and sets the default values.
func NewAndroidWorkProfileCustomConfiguration()(*AndroidWorkProfileCustomConfiguration) {
    m := &AndroidWorkProfileCustomConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.androidWorkProfileCustomConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAndroidWorkProfileCustomConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAndroidWorkProfileCustomConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAndroidWorkProfileCustomConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AndroidWorkProfileCustomConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["omaSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOmaSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OmaSettingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OmaSettingable)
                }
            }
            m.SetOmaSettings(res)
        }
        return nil
    }
    return res
}
// GetOmaSettings gets the omaSettings property value. OMA settings. This collection can contain a maximum of 500 elements.
// returns a []OmaSettingable when successful
func (m *AndroidWorkProfileCustomConfiguration) GetOmaSettings()([]OmaSettingable) {
    val, err := m.GetBackingStore().Get("omaSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OmaSettingable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AndroidWorkProfileCustomConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetOmaSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOmaSettings()))
        for i, v := range m.GetOmaSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("omaSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOmaSettings sets the omaSettings property value. OMA settings. This collection can contain a maximum of 500 elements.
func (m *AndroidWorkProfileCustomConfiguration) SetOmaSettings(value []OmaSettingable)() {
    err := m.GetBackingStore().Set("omaSettings", value)
    if err != nil {
        panic(err)
    }
}
type AndroidWorkProfileCustomConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOmaSettings()([]OmaSettingable)
    SetOmaSettings(value []OmaSettingable)()
}

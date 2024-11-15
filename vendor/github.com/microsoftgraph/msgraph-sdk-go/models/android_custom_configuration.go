package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AndroidCustomConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the androidCustomConfiguration resource.
type AndroidCustomConfiguration struct {
    DeviceConfiguration
}
// NewAndroidCustomConfiguration instantiates a new AndroidCustomConfiguration and sets the default values.
func NewAndroidCustomConfiguration()(*AndroidCustomConfiguration) {
    m := &AndroidCustomConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.androidCustomConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAndroidCustomConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAndroidCustomConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAndroidCustomConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AndroidCustomConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
// GetOmaSettings gets the omaSettings property value. OMA settings. This collection can contain a maximum of 1000 elements.
// returns a []OmaSettingable when successful
func (m *AndroidCustomConfiguration) GetOmaSettings()([]OmaSettingable) {
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
func (m *AndroidCustomConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetOmaSettings sets the omaSettings property value. OMA settings. This collection can contain a maximum of 1000 elements.
func (m *AndroidCustomConfiguration) SetOmaSettings(value []OmaSettingable)() {
    err := m.GetBackingStore().Set("omaSettings", value)
    if err != nil {
        panic(err)
    }
}
type AndroidCustomConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOmaSettings()([]OmaSettingable)
    SetOmaSettings(value []OmaSettingable)()
}

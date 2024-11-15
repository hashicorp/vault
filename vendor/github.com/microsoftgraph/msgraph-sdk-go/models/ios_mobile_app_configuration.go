package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosMobileAppConfiguration contains properties, inherited properties and actions for iOS mobile app configurations.
type IosMobileAppConfiguration struct {
    ManagedDeviceMobileAppConfiguration
}
// NewIosMobileAppConfiguration instantiates a new IosMobileAppConfiguration and sets the default values.
func NewIosMobileAppConfiguration()(*IosMobileAppConfiguration) {
    m := &IosMobileAppConfiguration{
        ManagedDeviceMobileAppConfiguration: *NewManagedDeviceMobileAppConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.iosMobileAppConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosMobileAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosMobileAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosMobileAppConfiguration(), nil
}
// GetEncodedSettingXml gets the encodedSettingXml property value. mdm app configuration Base64 binary.
// returns a []byte when successful
func (m *IosMobileAppConfiguration) GetEncodedSettingXml()([]byte) {
    val, err := m.GetBackingStore().Get("encodedSettingXml")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IosMobileAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedDeviceMobileAppConfiguration.GetFieldDeserializers()
    res["encodedSettingXml"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEncodedSettingXml(val)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppConfigurationSettingItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppConfigurationSettingItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppConfigurationSettingItemable)
                }
            }
            m.SetSettings(res)
        }
        return nil
    }
    return res
}
// GetSettings gets the settings property value. app configuration setting items.
// returns a []AppConfigurationSettingItemable when successful
func (m *IosMobileAppConfiguration) GetSettings()([]AppConfigurationSettingItemable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppConfigurationSettingItemable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosMobileAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedDeviceMobileAppConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("encodedSettingXml", m.GetEncodedSettingXml())
        if err != nil {
            return err
        }
    }
    if m.GetSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSettings()))
        for i, v := range m.GetSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("settings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEncodedSettingXml sets the encodedSettingXml property value. mdm app configuration Base64 binary.
func (m *IosMobileAppConfiguration) SetEncodedSettingXml(value []byte)() {
    err := m.GetBackingStore().Set("encodedSettingXml", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. app configuration setting items.
func (m *IosMobileAppConfiguration) SetSettings(value []AppConfigurationSettingItemable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
type IosMobileAppConfigurationable interface {
    ManagedDeviceMobileAppConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEncodedSettingXml()([]byte)
    GetSettings()([]AppConfigurationSettingItemable)
    SetEncodedSettingXml(value []byte)()
    SetSettings(value []AppConfigurationSettingItemable)()
}

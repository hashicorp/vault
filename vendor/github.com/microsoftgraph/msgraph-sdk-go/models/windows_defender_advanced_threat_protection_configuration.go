package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsDefenderAdvancedThreatProtectionConfiguration windows Defender AdvancedThreatProtection Configuration.
type WindowsDefenderAdvancedThreatProtectionConfiguration struct {
    DeviceConfiguration
}
// NewWindowsDefenderAdvancedThreatProtectionConfiguration instantiates a new WindowsDefenderAdvancedThreatProtectionConfiguration and sets the default values.
func NewWindowsDefenderAdvancedThreatProtectionConfiguration()(*WindowsDefenderAdvancedThreatProtectionConfiguration) {
    m := &WindowsDefenderAdvancedThreatProtectionConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windowsDefenderAdvancedThreatProtectionConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsDefenderAdvancedThreatProtectionConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsDefenderAdvancedThreatProtectionConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsDefenderAdvancedThreatProtectionConfiguration(), nil
}
// GetAllowSampleSharing gets the allowSampleSharing property value. Windows Defender AdvancedThreatProtection 'Allow Sample Sharing' Rule
// returns a *bool when successful
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) GetAllowSampleSharing()(*bool) {
    val, err := m.GetBackingStore().Get("allowSampleSharing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableExpeditedTelemetryReporting gets the enableExpeditedTelemetryReporting property value. Expedite Windows Defender Advanced Threat Protection telemetry reporting frequency.
// returns a *bool when successful
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) GetEnableExpeditedTelemetryReporting()(*bool) {
    val, err := m.GetBackingStore().Get("enableExpeditedTelemetryReporting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["allowSampleSharing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowSampleSharing(val)
        }
        return nil
    }
    res["enableExpeditedTelemetryReporting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableExpeditedTelemetryReporting(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowSampleSharing", m.GetAllowSampleSharing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enableExpeditedTelemetryReporting", m.GetEnableExpeditedTelemetryReporting())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowSampleSharing sets the allowSampleSharing property value. Windows Defender AdvancedThreatProtection 'Allow Sample Sharing' Rule
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) SetAllowSampleSharing(value *bool)() {
    err := m.GetBackingStore().Set("allowSampleSharing", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableExpeditedTelemetryReporting sets the enableExpeditedTelemetryReporting property value. Expedite Windows Defender Advanced Threat Protection telemetry reporting frequency.
func (m *WindowsDefenderAdvancedThreatProtectionConfiguration) SetEnableExpeditedTelemetryReporting(value *bool)() {
    err := m.GetBackingStore().Set("enableExpeditedTelemetryReporting", value)
    if err != nil {
        panic(err)
    }
}
type WindowsDefenderAdvancedThreatProtectionConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowSampleSharing()(*bool)
    GetEnableExpeditedTelemetryReporting()(*bool)
    SetAllowSampleSharing(value *bool)()
    SetEnableExpeditedTelemetryReporting(value *bool)()
}

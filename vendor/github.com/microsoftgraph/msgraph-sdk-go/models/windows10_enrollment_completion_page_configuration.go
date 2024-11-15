package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10EnrollmentCompletionPageConfiguration windows 10 Enrollment Status Page Configuration
type Windows10EnrollmentCompletionPageConfiguration struct {
    DeviceEnrollmentConfiguration
}
// NewWindows10EnrollmentCompletionPageConfiguration instantiates a new Windows10EnrollmentCompletionPageConfiguration and sets the default values.
func NewWindows10EnrollmentCompletionPageConfiguration()(*Windows10EnrollmentCompletionPageConfiguration) {
    m := &Windows10EnrollmentCompletionPageConfiguration{
        DeviceEnrollmentConfiguration: *NewDeviceEnrollmentConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows10EnrollmentCompletionPageConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows10EnrollmentCompletionPageConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10EnrollmentCompletionPageConfiguration(), nil
}
// GetAllowNonBlockingAppInstallation gets the allowNonBlockingAppInstallation property value. When TRUE, ESP (Enrollment Status Page) installs all required apps targeted during technician phase and ignores any failures for non-blocking apps. When FALSE, ESP fails on any error during app install. The default is false.
// returns a *bool when successful
func (m *Windows10EnrollmentCompletionPageConfiguration) GetAllowNonBlockingAppInstallation()(*bool) {
    val, err := m.GetBackingStore().Get("allowNonBlockingAppInstallation")
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
func (m *Windows10EnrollmentCompletionPageConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceEnrollmentConfiguration.GetFieldDeserializers()
    res["allowNonBlockingAppInstallation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowNonBlockingAppInstallation(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *Windows10EnrollmentCompletionPageConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceEnrollmentConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowNonBlockingAppInstallation", m.GetAllowNonBlockingAppInstallation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowNonBlockingAppInstallation sets the allowNonBlockingAppInstallation property value. When TRUE, ESP (Enrollment Status Page) installs all required apps targeted during technician phase and ignores any failures for non-blocking apps. When FALSE, ESP fails on any error during app install. The default is false.
func (m *Windows10EnrollmentCompletionPageConfiguration) SetAllowNonBlockingAppInstallation(value *bool)() {
    err := m.GetBackingStore().Set("allowNonBlockingAppInstallation", value)
    if err != nil {
        panic(err)
    }
}
type Windows10EnrollmentCompletionPageConfigurationable interface {
    DeviceEnrollmentConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowNonBlockingAppInstallation()(*bool)
    SetAllowNonBlockingAppInstallation(value *bool)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceEnrollmentLimitConfiguration device Enrollment Configuration that restricts the number of devices a user can enroll
type DeviceEnrollmentLimitConfiguration struct {
    DeviceEnrollmentConfiguration
}
// NewDeviceEnrollmentLimitConfiguration instantiates a new DeviceEnrollmentLimitConfiguration and sets the default values.
func NewDeviceEnrollmentLimitConfiguration()(*DeviceEnrollmentLimitConfiguration) {
    m := &DeviceEnrollmentLimitConfiguration{
        DeviceEnrollmentConfiguration: *NewDeviceEnrollmentConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.deviceEnrollmentLimitConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDeviceEnrollmentLimitConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceEnrollmentLimitConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceEnrollmentLimitConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceEnrollmentLimitConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceEnrollmentConfiguration.GetFieldDeserializers()
    res["limit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLimit(val)
        }
        return nil
    }
    return res
}
// GetLimit gets the limit property value. The maximum number of devices that a user can enroll
// returns a *int32 when successful
func (m *DeviceEnrollmentLimitConfiguration) GetLimit()(*int32) {
    val, err := m.GetBackingStore().Get("limit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceEnrollmentLimitConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceEnrollmentConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("limit", m.GetLimit())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLimit sets the limit property value. The maximum number of devices that a user can enroll
func (m *DeviceEnrollmentLimitConfiguration) SetLimit(value *int32)() {
    err := m.GetBackingStore().Set("limit", value)
    if err != nil {
        panic(err)
    }
}
type DeviceEnrollmentLimitConfigurationable interface {
    DeviceEnrollmentConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLimit()(*int32)
    SetLimit(value *int32)()
}

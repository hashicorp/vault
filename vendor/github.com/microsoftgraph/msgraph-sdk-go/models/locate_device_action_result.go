package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LocateDeviceActionResult locate device action result
type LocateDeviceActionResult struct {
    DeviceActionResult
}
// NewLocateDeviceActionResult instantiates a new LocateDeviceActionResult and sets the default values.
func NewLocateDeviceActionResult()(*LocateDeviceActionResult) {
    m := &LocateDeviceActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateLocateDeviceActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLocateDeviceActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocateDeviceActionResult(), nil
}
// GetDeviceLocation gets the deviceLocation property value. device location
// returns a DeviceGeoLocationable when successful
func (m *LocateDeviceActionResult) GetDeviceLocation()(DeviceGeoLocationable) {
    val, err := m.GetBackingStore().Get("deviceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceGeoLocationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LocateDeviceActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["deviceLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceGeoLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceLocation(val.(DeviceGeoLocationable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *LocateDeviceActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("deviceLocation", m.GetDeviceLocation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceLocation sets the deviceLocation property value. device location
func (m *LocateDeviceActionResult) SetDeviceLocation(value DeviceGeoLocationable)() {
    err := m.GetBackingStore().Set("deviceLocation", value)
    if err != nil {
        panic(err)
    }
}
type LocateDeviceActionResultable interface {
    DeviceActionResultable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceLocation()(DeviceGeoLocationable)
    SetDeviceLocation(value DeviceGeoLocationable)()
}

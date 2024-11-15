package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsDefenderScanActionResult windows Defender last scan result
type WindowsDefenderScanActionResult struct {
    DeviceActionResult
}
// NewWindowsDefenderScanActionResult instantiates a new WindowsDefenderScanActionResult and sets the default values.
func NewWindowsDefenderScanActionResult()(*WindowsDefenderScanActionResult) {
    m := &WindowsDefenderScanActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateWindowsDefenderScanActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsDefenderScanActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsDefenderScanActionResult(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsDefenderScanActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["scanType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScanType(val)
        }
        return nil
    }
    return res
}
// GetScanType gets the scanType property value. Scan type either full scan or quick scan
// returns a *string when successful
func (m *WindowsDefenderScanActionResult) GetScanType()(*string) {
    val, err := m.GetBackingStore().Get("scanType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsDefenderScanActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("scanType", m.GetScanType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetScanType sets the scanType property value. Scan type either full scan or quick scan
func (m *WindowsDefenderScanActionResult) SetScanType(value *string)() {
    err := m.GetBackingStore().Set("scanType", value)
    if err != nil {
        panic(err)
    }
}
type WindowsDefenderScanActionResultable interface {
    DeviceActionResultable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetScanType()(*string)
    SetScanType(value *string)()
}

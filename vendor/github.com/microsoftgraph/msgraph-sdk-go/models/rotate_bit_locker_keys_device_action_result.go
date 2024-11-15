package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RotateBitLockerKeysDeviceActionResult rotateBitLockerKeys device action result
type RotateBitLockerKeysDeviceActionResult struct {
    DeviceActionResult
}
// NewRotateBitLockerKeysDeviceActionResult instantiates a new RotateBitLockerKeysDeviceActionResult and sets the default values.
func NewRotateBitLockerKeysDeviceActionResult()(*RotateBitLockerKeysDeviceActionResult) {
    m := &RotateBitLockerKeysDeviceActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateRotateBitLockerKeysDeviceActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRotateBitLockerKeysDeviceActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRotateBitLockerKeysDeviceActionResult(), nil
}
// GetErrorCode gets the errorCode property value. RotateBitLockerKeys action error code
// returns a *int32 when successful
func (m *RotateBitLockerKeysDeviceActionResult) GetErrorCode()(*int32) {
    val, err := m.GetBackingStore().Get("errorCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RotateBitLockerKeysDeviceActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["errorCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetErrorCode(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *RotateBitLockerKeysDeviceActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("errorCode", m.GetErrorCode())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetErrorCode sets the errorCode property value. RotateBitLockerKeys action error code
func (m *RotateBitLockerKeysDeviceActionResult) SetErrorCode(value *int32)() {
    err := m.GetBackingStore().Set("errorCode", value)
    if err != nil {
        panic(err)
    }
}
type RotateBitLockerKeysDeviceActionResultable interface {
    DeviceActionResultable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetErrorCode()(*int32)
    SetErrorCode(value *int32)()
}

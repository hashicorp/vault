package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ResetPasscodeActionResult reset passcode action result
type ResetPasscodeActionResult struct {
    DeviceActionResult
}
// NewResetPasscodeActionResult instantiates a new ResetPasscodeActionResult and sets the default values.
func NewResetPasscodeActionResult()(*ResetPasscodeActionResult) {
    m := &ResetPasscodeActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateResetPasscodeActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateResetPasscodeActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResetPasscodeActionResult(), nil
}
// GetErrorCode gets the errorCode property value. RotateBitLockerKeys action error code. Valid values 0 to 2147483647
// returns a *int32 when successful
func (m *ResetPasscodeActionResult) GetErrorCode()(*int32) {
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
func (m *ResetPasscodeActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["passcode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscode(val)
        }
        return nil
    }
    return res
}
// GetPasscode gets the passcode property value. Newly generated passcode for the device
// returns a *string when successful
func (m *ResetPasscodeActionResult) GetPasscode()(*string) {
    val, err := m.GetBackingStore().Get("passcode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ResetPasscodeActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteStringValue("passcode", m.GetPasscode())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetErrorCode sets the errorCode property value. RotateBitLockerKeys action error code. Valid values 0 to 2147483647
func (m *ResetPasscodeActionResult) SetErrorCode(value *int32)() {
    err := m.GetBackingStore().Set("errorCode", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscode sets the passcode property value. Newly generated passcode for the device
func (m *ResetPasscodeActionResult) SetPasscode(value *string)() {
    err := m.GetBackingStore().Set("passcode", value)
    if err != nil {
        panic(err)
    }
}
type ResetPasscodeActionResultable interface {
    DeviceActionResultable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetErrorCode()(*int32)
    GetPasscode()(*string)
    SetErrorCode(value *int32)()
    SetPasscode(value *string)()
}

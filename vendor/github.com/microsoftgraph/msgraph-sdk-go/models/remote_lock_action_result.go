package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RemoteLockActionResult lock action result with a pin to unlock
type RemoteLockActionResult struct {
    DeviceActionResult
}
// NewRemoteLockActionResult instantiates a new RemoteLockActionResult and sets the default values.
func NewRemoteLockActionResult()(*RemoteLockActionResult) {
    m := &RemoteLockActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateRemoteLockActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRemoteLockActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRemoteLockActionResult(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RemoteLockActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["unlockPin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnlockPin(val)
        }
        return nil
    }
    return res
}
// GetUnlockPin gets the unlockPin property value. Pin to unlock the client
// returns a *string when successful
func (m *RemoteLockActionResult) GetUnlockPin()(*string) {
    val, err := m.GetBackingStore().Get("unlockPin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RemoteLockActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("unlockPin", m.GetUnlockPin())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUnlockPin sets the unlockPin property value. Pin to unlock the client
func (m *RemoteLockActionResult) SetUnlockPin(value *string)() {
    err := m.GetBackingStore().Set("unlockPin", value)
    if err != nil {
        panic(err)
    }
}
type RemoteLockActionResultable interface {
    DeviceActionResultable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUnlockPin()(*string)
    SetUnlockPin(value *string)()
}

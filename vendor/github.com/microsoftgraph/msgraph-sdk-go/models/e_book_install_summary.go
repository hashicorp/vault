package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EBookInstallSummary contains properties for the installation summary of a book for a device.
type EBookInstallSummary struct {
    Entity
}
// NewEBookInstallSummary instantiates a new EBookInstallSummary and sets the default values.
func NewEBookInstallSummary()(*EBookInstallSummary) {
    m := &EBookInstallSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEBookInstallSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEBookInstallSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEBookInstallSummary(), nil
}
// GetFailedDeviceCount gets the failedDeviceCount property value. Number of Devices that have failed to install this book.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetFailedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedUserCount gets the failedUserCount property value. Number of Users that have 1 or more device that failed to install this book.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetFailedUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedUserCount")
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
func (m *EBookInstallSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["failedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedDeviceCount(val)
        }
        return nil
    }
    res["failedUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedUserCount(val)
        }
        return nil
    }
    res["installedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstalledDeviceCount(val)
        }
        return nil
    }
    res["installedUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstalledUserCount(val)
        }
        return nil
    }
    res["notInstalledDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotInstalledDeviceCount(val)
        }
        return nil
    }
    res["notInstalledUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotInstalledUserCount(val)
        }
        return nil
    }
    return res
}
// GetInstalledDeviceCount gets the installedDeviceCount property value. Number of Devices that have successfully installed this book.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetInstalledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("installedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetInstalledUserCount gets the installedUserCount property value. Number of Users whose devices have all succeeded to install this book.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetInstalledUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("installedUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotInstalledDeviceCount gets the notInstalledDeviceCount property value. Number of Devices that does not have this book installed.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetNotInstalledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("notInstalledDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotInstalledUserCount gets the notInstalledUserCount property value. Number of Users that did not install this book.
// returns a *int32 when successful
func (m *EBookInstallSummary) GetNotInstalledUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("notInstalledUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EBookInstallSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("failedDeviceCount", m.GetFailedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("failedUserCount", m.GetFailedUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("installedDeviceCount", m.GetInstalledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("installedUserCount", m.GetInstalledUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notInstalledDeviceCount", m.GetNotInstalledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notInstalledUserCount", m.GetNotInstalledUserCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFailedDeviceCount sets the failedDeviceCount property value. Number of Devices that have failed to install this book.
func (m *EBookInstallSummary) SetFailedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("failedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedUserCount sets the failedUserCount property value. Number of Users that have 1 or more device that failed to install this book.
func (m *EBookInstallSummary) SetFailedUserCount(value *int32)() {
    err := m.GetBackingStore().Set("failedUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledDeviceCount sets the installedDeviceCount property value. Number of Devices that have successfully installed this book.
func (m *EBookInstallSummary) SetInstalledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("installedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledUserCount sets the installedUserCount property value. Number of Users whose devices have all succeeded to install this book.
func (m *EBookInstallSummary) SetInstalledUserCount(value *int32)() {
    err := m.GetBackingStore().Set("installedUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotInstalledDeviceCount sets the notInstalledDeviceCount property value. Number of Devices that does not have this book installed.
func (m *EBookInstallSummary) SetNotInstalledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("notInstalledDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotInstalledUserCount sets the notInstalledUserCount property value. Number of Users that did not install this book.
func (m *EBookInstallSummary) SetNotInstalledUserCount(value *int32)() {
    err := m.GetBackingStore().Set("notInstalledUserCount", value)
    if err != nil {
        panic(err)
    }
}
type EBookInstallSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFailedDeviceCount()(*int32)
    GetFailedUserCount()(*int32)
    GetInstalledDeviceCount()(*int32)
    GetInstalledUserCount()(*int32)
    GetNotInstalledDeviceCount()(*int32)
    GetNotInstalledUserCount()(*int32)
    SetFailedDeviceCount(value *int32)()
    SetFailedUserCount(value *int32)()
    SetInstalledDeviceCount(value *int32)()
    SetInstalledUserCount(value *int32)()
    SetNotInstalledDeviceCount(value *int32)()
    SetNotInstalledUserCount(value *int32)()
}

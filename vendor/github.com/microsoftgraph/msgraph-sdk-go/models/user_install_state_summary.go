package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserInstallStateSummary contains properties for the installation state summary for a user.
type UserInstallStateSummary struct {
    Entity
}
// NewUserInstallStateSummary instantiates a new UserInstallStateSummary and sets the default values.
func NewUserInstallStateSummary()(*UserInstallStateSummary) {
    m := &UserInstallStateSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserInstallStateSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserInstallStateSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserInstallStateSummary(), nil
}
// GetDeviceStates gets the deviceStates property value. The install state of the eBook.
// returns a []DeviceInstallStateable when successful
func (m *UserInstallStateSummary) GetDeviceStates()([]DeviceInstallStateable) {
    val, err := m.GetBackingStore().Get("deviceStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceInstallStateable)
    }
    return nil
}
// GetFailedDeviceCount gets the failedDeviceCount property value. Failed Device Count.
// returns a *int32 when successful
func (m *UserInstallStateSummary) GetFailedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedDeviceCount")
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
func (m *UserInstallStateSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceInstallStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceInstallStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceInstallStateable)
                }
            }
            m.SetDeviceStates(res)
        }
        return nil
    }
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
    res["userName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserName(val)
        }
        return nil
    }
    return res
}
// GetInstalledDeviceCount gets the installedDeviceCount property value. Installed Device Count.
// returns a *int32 when successful
func (m *UserInstallStateSummary) GetInstalledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("installedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotInstalledDeviceCount gets the notInstalledDeviceCount property value. Not installed device count.
// returns a *int32 when successful
func (m *UserInstallStateSummary) GetNotInstalledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("notInstalledDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUserName gets the userName property value. User name.
// returns a *string when successful
func (m *UserInstallStateSummary) GetUserName()(*string) {
    val, err := m.GetBackingStore().Get("userName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserInstallStateSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDeviceStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceStates()))
        for i, v := range m.GetDeviceStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceStates", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("failedDeviceCount", m.GetFailedDeviceCount())
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
        err = writer.WriteInt32Value("notInstalledDeviceCount", m.GetNotInstalledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userName", m.GetUserName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceStates sets the deviceStates property value. The install state of the eBook.
func (m *UserInstallStateSummary) SetDeviceStates(value []DeviceInstallStateable)() {
    err := m.GetBackingStore().Set("deviceStates", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedDeviceCount sets the failedDeviceCount property value. Failed Device Count.
func (m *UserInstallStateSummary) SetFailedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("failedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledDeviceCount sets the installedDeviceCount property value. Installed Device Count.
func (m *UserInstallStateSummary) SetInstalledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("installedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotInstalledDeviceCount sets the notInstalledDeviceCount property value. Not installed device count.
func (m *UserInstallStateSummary) SetNotInstalledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("notInstalledDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUserName sets the userName property value. User name.
func (m *UserInstallStateSummary) SetUserName(value *string)() {
    err := m.GetBackingStore().Set("userName", value)
    if err != nil {
        panic(err)
    }
}
type UserInstallStateSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceStates()([]DeviceInstallStateable)
    GetFailedDeviceCount()(*int32)
    GetInstalledDeviceCount()(*int32)
    GetNotInstalledDeviceCount()(*int32)
    GetUserName()(*string)
    SetDeviceStates(value []DeviceInstallStateable)()
    SetFailedDeviceCount(value *int32)()
    SetInstalledDeviceCount(value *int32)()
    SetNotInstalledDeviceCount(value *int32)()
    SetUserName(value *string)()
}

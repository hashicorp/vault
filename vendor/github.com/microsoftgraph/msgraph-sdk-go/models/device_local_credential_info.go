package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceLocalCredentialInfo struct {
    Entity
}
// NewDeviceLocalCredentialInfo instantiates a new DeviceLocalCredentialInfo and sets the default values.
func NewDeviceLocalCredentialInfo()(*DeviceLocalCredentialInfo) {
    m := &DeviceLocalCredentialInfo{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceLocalCredentialInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceLocalCredentialInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceLocalCredentialInfo(), nil
}
// GetCredentials gets the credentials property value. The credentials of the device's local administrator account backed up to Azure Active Directory.
// returns a []DeviceLocalCredentialable when successful
func (m *DeviceLocalCredentialInfo) GetCredentials()([]DeviceLocalCredentialable) {
    val, err := m.GetBackingStore().Get("credentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceLocalCredentialable)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. Display name of the device that the local credentials are associated with.
// returns a *string when successful
func (m *DeviceLocalCredentialInfo) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceLocalCredentialInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["credentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceLocalCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceLocalCredentialable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceLocalCredentialable)
                }
            }
            m.SetCredentials(res)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
        }
        return nil
    }
    res["lastBackupDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastBackupDateTime(val)
        }
        return nil
    }
    res["refreshDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRefreshDateTime(val)
        }
        return nil
    }
    return res
}
// GetLastBackupDateTime gets the lastBackupDateTime property value. When the local administrator account credential was backed up to Azure Active Directory.
// returns a *Time when successful
func (m *DeviceLocalCredentialInfo) GetLastBackupDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastBackupDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRefreshDateTime gets the refreshDateTime property value. When the local administrator account credential will be refreshed and backed up to Azure Active Directory.
// returns a *Time when successful
func (m *DeviceLocalCredentialInfo) GetRefreshDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("refreshDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceLocalCredentialInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCredentials()))
        for i, v := range m.GetCredentials() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("credentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceName", m.GetDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastBackupDateTime", m.GetLastBackupDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("refreshDateTime", m.GetRefreshDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCredentials sets the credentials property value. The credentials of the device's local administrator account backed up to Azure Active Directory.
func (m *DeviceLocalCredentialInfo) SetCredentials(value []DeviceLocalCredentialable)() {
    err := m.GetBackingStore().Set("credentials", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. Display name of the device that the local credentials are associated with.
func (m *DeviceLocalCredentialInfo) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastBackupDateTime sets the lastBackupDateTime property value. When the local administrator account credential was backed up to Azure Active Directory.
func (m *DeviceLocalCredentialInfo) SetLastBackupDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastBackupDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRefreshDateTime sets the refreshDateTime property value. When the local administrator account credential will be refreshed and backed up to Azure Active Directory.
func (m *DeviceLocalCredentialInfo) SetRefreshDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("refreshDateTime", value)
    if err != nil {
        panic(err)
    }
}
type DeviceLocalCredentialInfoable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCredentials()([]DeviceLocalCredentialable)
    GetDeviceName()(*string)
    GetLastBackupDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRefreshDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetCredentials(value []DeviceLocalCredentialable)()
    SetDeviceName(value *string)()
    SetLastBackupDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRefreshDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}

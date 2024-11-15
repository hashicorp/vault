package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceLocalCredential struct {
    Entity
}
// NewDeviceLocalCredential instantiates a new DeviceLocalCredential and sets the default values.
func NewDeviceLocalCredential()(*DeviceLocalCredential) {
    m := &DeviceLocalCredential{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceLocalCredentialFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceLocalCredentialFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceLocalCredential(), nil
}
// GetAccountName gets the accountName property value. The name of the local admin account for which LAPS is enabled.
// returns a *string when successful
func (m *DeviceLocalCredential) GetAccountName()(*string) {
    val, err := m.GetBackingStore().Get("accountName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAccountSid gets the accountSid property value. The SID of the local admin account for which LAPS is enabled.
// returns a *string when successful
func (m *DeviceLocalCredential) GetAccountSid()(*string) {
    val, err := m.GetBackingStore().Get("accountSid")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackupDateTime gets the backupDateTime property value. When the local administrator account credential for the device object was backed up to Azure Active Directory.
// returns a *Time when successful
func (m *DeviceLocalCredential) GetBackupDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("backupDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceLocalCredential) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accountName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountName(val)
        }
        return nil
    }
    res["accountSid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountSid(val)
        }
        return nil
    }
    res["backupDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBackupDateTime(val)
        }
        return nil
    }
    res["passwordBase64"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBase64(val)
        }
        return nil
    }
    return res
}
// GetPasswordBase64 gets the passwordBase64 property value. The password for the local administrator account that is backed up to Azure Active Directory and returned as a Base64 encoded value.
// returns a *string when successful
func (m *DeviceLocalCredential) GetPasswordBase64()(*string) {
    val, err := m.GetBackingStore().Get("passwordBase64")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceLocalCredential) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("accountName", m.GetAccountName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("accountSid", m.GetAccountSid())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("backupDateTime", m.GetBackupDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("passwordBase64", m.GetPasswordBase64())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountName sets the accountName property value. The name of the local admin account for which LAPS is enabled.
func (m *DeviceLocalCredential) SetAccountName(value *string)() {
    err := m.GetBackingStore().Set("accountName", value)
    if err != nil {
        panic(err)
    }
}
// SetAccountSid sets the accountSid property value. The SID of the local admin account for which LAPS is enabled.
func (m *DeviceLocalCredential) SetAccountSid(value *string)() {
    err := m.GetBackingStore().Set("accountSid", value)
    if err != nil {
        panic(err)
    }
}
// SetBackupDateTime sets the backupDateTime property value. When the local administrator account credential for the device object was backed up to Azure Active Directory.
func (m *DeviceLocalCredential) SetBackupDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("backupDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBase64 sets the passwordBase64 property value. The password for the local administrator account that is backed up to Azure Active Directory and returned as a Base64 encoded value.
func (m *DeviceLocalCredential) SetPasswordBase64(value *string)() {
    err := m.GetBackingStore().Set("passwordBase64", value)
    if err != nil {
        panic(err)
    }
}
type DeviceLocalCredentialable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountName()(*string)
    GetAccountSid()(*string)
    GetBackupDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPasswordBase64()(*string)
    SetAccountName(value *string)()
    SetAccountSid(value *string)()
    SetBackupDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPasswordBase64(value *string)()
}

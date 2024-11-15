package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionAppLockerFile windows Information Protection AppLocker File
type WindowsInformationProtectionAppLockerFile struct {
    Entity
}
// NewWindowsInformationProtectionAppLockerFile instantiates a new WindowsInformationProtectionAppLockerFile and sets the default values.
func NewWindowsInformationProtectionAppLockerFile()(*WindowsInformationProtectionAppLockerFile) {
    m := &WindowsInformationProtectionAppLockerFile{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionAppLockerFile(), nil
}
// GetDisplayName gets the displayName property value. The friendly name
// returns a *string when successful
func (m *WindowsInformationProtectionAppLockerFile) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *WindowsInformationProtectionAppLockerFile) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["file"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFile(val)
        }
        return nil
    }
    res["fileHash"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileHash(val)
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetFile gets the file property value. File as a byte array
// returns a []byte when successful
func (m *WindowsInformationProtectionAppLockerFile) GetFile()([]byte) {
    val, err := m.GetBackingStore().Get("file")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetFileHash gets the fileHash property value. SHA256 hash of the file
// returns a *string when successful
func (m *WindowsInformationProtectionAppLockerFile) GetFileHash()(*string) {
    val, err := m.GetBackingStore().Get("fileHash")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVersion gets the version property value. Version of the entity.
// returns a *string when successful
func (m *WindowsInformationProtectionAppLockerFile) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionAppLockerFile) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("file", m.GetFile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("fileHash", m.GetFileHash())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The friendly name
func (m *WindowsInformationProtectionAppLockerFile) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetFile sets the file property value. File as a byte array
func (m *WindowsInformationProtectionAppLockerFile) SetFile(value []byte)() {
    err := m.GetBackingStore().Set("file", value)
    if err != nil {
        panic(err)
    }
}
// SetFileHash sets the fileHash property value. SHA256 hash of the file
func (m *WindowsInformationProtectionAppLockerFile) SetFileHash(value *string)() {
    err := m.GetBackingStore().Set("fileHash", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Version of the entity.
func (m *WindowsInformationProtectionAppLockerFile) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type WindowsInformationProtectionAppLockerFileable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetFile()([]byte)
    GetFileHash()(*string)
    GetVersion()(*string)
    SetDisplayName(value *string)()
    SetFile(value []byte)()
    SetFileHash(value *string)()
    SetVersion(value *string)()
}

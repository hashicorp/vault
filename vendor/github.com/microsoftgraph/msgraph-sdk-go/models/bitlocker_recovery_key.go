package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type BitlockerRecoveryKey struct {
    Entity
}
// NewBitlockerRecoveryKey instantiates a new BitlockerRecoveryKey and sets the default values.
func NewBitlockerRecoveryKey()(*BitlockerRecoveryKey) {
    m := &BitlockerRecoveryKey{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBitlockerRecoveryKeyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBitlockerRecoveryKeyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBitlockerRecoveryKey(), nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the key was originally backed up to Microsoft Entra ID. Not nullable.
// returns a *Time when successful
func (m *BitlockerRecoveryKey) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeviceId gets the deviceId property value. Identifier of the device the BitLocker key is originally backed up from. Supports $filter (eq).
// returns a *string when successful
func (m *BitlockerRecoveryKey) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
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
func (m *BitlockerRecoveryKey) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["deviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceId(val)
        }
        return nil
    }
    res["key"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKey(val)
        }
        return nil
    }
    res["volumeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVolumeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVolumeType(val.(*VolumeType))
        }
        return nil
    }
    return res
}
// GetKey gets the key property value. The BitLocker recovery key. Returned only on $select. Not nullable.
// returns a *string when successful
func (m *BitlockerRecoveryKey) GetKey()(*string) {
    val, err := m.GetBackingStore().Get("key")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVolumeType gets the volumeType property value. Indicates the type of volume the BitLocker key is associated with. The possible values are: 1 (for operatingSystemVolume), 2 (for fixedDataVolume), 3 (for removableDataVolume), and 4 (for unknownFutureValue).
// returns a *VolumeType when successful
func (m *BitlockerRecoveryKey) GetVolumeType()(*VolumeType) {
    val, err := m.GetBackingStore().Get("volumeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VolumeType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BitlockerRecoveryKey) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("key", m.GetKey())
        if err != nil {
            return err
        }
    }
    if m.GetVolumeType() != nil {
        cast := (*m.GetVolumeType()).String()
        err = writer.WriteStringValue("volumeType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the key was originally backed up to Microsoft Entra ID. Not nullable.
func (m *BitlockerRecoveryKey) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceId sets the deviceId property value. Identifier of the device the BitLocker key is originally backed up from. Supports $filter (eq).
func (m *BitlockerRecoveryKey) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetKey sets the key property value. The BitLocker recovery key. Returned only on $select. Not nullable.
func (m *BitlockerRecoveryKey) SetKey(value *string)() {
    err := m.GetBackingStore().Set("key", value)
    if err != nil {
        panic(err)
    }
}
// SetVolumeType sets the volumeType property value. Indicates the type of volume the BitLocker key is associated with. The possible values are: 1 (for operatingSystemVolume), 2 (for fixedDataVolume), 3 (for removableDataVolume), and 4 (for unknownFutureValue).
func (m *BitlockerRecoveryKey) SetVolumeType(value *VolumeType)() {
    err := m.GetBackingStore().Set("volumeType", value)
    if err != nil {
        panic(err)
    }
}
type BitlockerRecoveryKeyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeviceId()(*string)
    GetKey()(*string)
    GetVolumeType()(*VolumeType)
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeviceId(value *string)()
    SetKey(value *string)()
    SetVolumeType(value *VolumeType)()
}

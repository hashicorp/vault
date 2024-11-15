package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ImportedWindowsAutopilotDeviceIdentity imported windows autopilot devices.
type ImportedWindowsAutopilotDeviceIdentity struct {
    Entity
}
// NewImportedWindowsAutopilotDeviceIdentity instantiates a new ImportedWindowsAutopilotDeviceIdentity and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentity()(*ImportedWindowsAutopilotDeviceIdentity) {
    m := &ImportedWindowsAutopilotDeviceIdentity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewImportedWindowsAutopilotDeviceIdentity(), nil
}
// GetAssignedUserPrincipalName gets the assignedUserPrincipalName property value. UPN of the user the device will be assigned
// returns a *string when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetAssignedUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("assignedUserPrincipalName")
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
func (m *ImportedWindowsAutopilotDeviceIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignedUserPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedUserPrincipalName(val)
        }
        return nil
    }
    res["groupTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupTag(val)
        }
        return nil
    }
    res["hardwareIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHardwareIdentifier(val)
        }
        return nil
    }
    res["importId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImportId(val)
        }
        return nil
    }
    res["productKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductKey(val)
        }
        return nil
    }
    res["serialNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSerialNumber(val)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateImportedWindowsAutopilotDeviceIdentityStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(ImportedWindowsAutopilotDeviceIdentityStateable))
        }
        return nil
    }
    return res
}
// GetGroupTag gets the groupTag property value. Group Tag of the Windows autopilot device.
// returns a *string when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetGroupTag()(*string) {
    val, err := m.GetBackingStore().Get("groupTag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHardwareIdentifier gets the hardwareIdentifier property value. Hardware Blob of the Windows autopilot device.
// returns a []byte when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetHardwareIdentifier()([]byte) {
    val, err := m.GetBackingStore().Get("hardwareIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetImportId gets the importId property value. The Import Id of the Windows autopilot device.
// returns a *string when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetImportId()(*string) {
    val, err := m.GetBackingStore().Get("importId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductKey gets the productKey property value. Product Key of the Windows autopilot device.
// returns a *string when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetProductKey()(*string) {
    val, err := m.GetBackingStore().Get("productKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSerialNumber gets the serialNumber property value. Serial number of the Windows autopilot device.
// returns a *string when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetSerialNumber()(*string) {
    val, err := m.GetBackingStore().Get("serialNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetState gets the state property value. Current state of the imported device.
// returns a ImportedWindowsAutopilotDeviceIdentityStateable when successful
func (m *ImportedWindowsAutopilotDeviceIdentity) GetState()(ImportedWindowsAutopilotDeviceIdentityStateable) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ImportedWindowsAutopilotDeviceIdentityStateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ImportedWindowsAutopilotDeviceIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("assignedUserPrincipalName", m.GetAssignedUserPrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("groupTag", m.GetGroupTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("hardwareIdentifier", m.GetHardwareIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("importId", m.GetImportId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productKey", m.GetProductKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serialNumber", m.GetSerialNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignedUserPrincipalName sets the assignedUserPrincipalName property value. UPN of the user the device will be assigned
func (m *ImportedWindowsAutopilotDeviceIdentity) SetAssignedUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("assignedUserPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupTag sets the groupTag property value. Group Tag of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetGroupTag(value *string)() {
    err := m.GetBackingStore().Set("groupTag", value)
    if err != nil {
        panic(err)
    }
}
// SetHardwareIdentifier sets the hardwareIdentifier property value. Hardware Blob of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetHardwareIdentifier(value []byte)() {
    err := m.GetBackingStore().Set("hardwareIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetImportId sets the importId property value. The Import Id of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetImportId(value *string)() {
    err := m.GetBackingStore().Set("importId", value)
    if err != nil {
        panic(err)
    }
}
// SetProductKey sets the productKey property value. Product Key of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetProductKey(value *string)() {
    err := m.GetBackingStore().Set("productKey", value)
    if err != nil {
        panic(err)
    }
}
// SetSerialNumber sets the serialNumber property value. Serial number of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetSerialNumber(value *string)() {
    err := m.GetBackingStore().Set("serialNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. Current state of the imported device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetState(value ImportedWindowsAutopilotDeviceIdentityStateable)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
type ImportedWindowsAutopilotDeviceIdentityable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignedUserPrincipalName()(*string)
    GetGroupTag()(*string)
    GetHardwareIdentifier()([]byte)
    GetImportId()(*string)
    GetProductKey()(*string)
    GetSerialNumber()(*string)
    GetState()(ImportedWindowsAutopilotDeviceIdentityStateable)
    SetAssignedUserPrincipalName(value *string)()
    SetGroupTag(value *string)()
    SetHardwareIdentifier(value []byte)()
    SetImportId(value *string)()
    SetProductKey(value *string)()
    SetSerialNumber(value *string)()
    SetState(value ImportedWindowsAutopilotDeviceIdentityStateable)()
}

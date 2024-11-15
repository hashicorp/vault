package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type RegistryKeyState struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRegistryKeyState instantiates a new RegistryKeyState and sets the default values.
func NewRegistryKeyState()(*RegistryKeyState) {
    m := &RegistryKeyState{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRegistryKeyStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRegistryKeyStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRegistryKeyState(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *RegistryKeyState) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *RegistryKeyState) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RegistryKeyState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRegistryHive)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHive(val.(*RegistryHive))
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
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["oldKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOldKey(val)
        }
        return nil
    }
    res["oldValueData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOldValueData(val)
        }
        return nil
    }
    res["oldValueName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOldValueName(val)
        }
        return nil
    }
    res["operation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRegistryOperation)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperation(val.(*RegistryOperation))
        }
        return nil
    }
    res["processId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessId(val)
        }
        return nil
    }
    res["valueData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValueData(val)
        }
        return nil
    }
    res["valueName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValueName(val)
        }
        return nil
    }
    res["valueType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRegistryValueType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValueType(val.(*RegistryValueType))
        }
        return nil
    }
    return res
}
// GetHive gets the hive property value. A Windows registry hive : HKEYCURRENTCONFIG HKEYCURRENTUSER HKEYLOCALMACHINE/SAM HKEYLOCALMACHINE/Security HKEYLOCALMACHINE/Software HKEYLOCALMACHINE/System HKEY_USERS/.Default. Possible values are: unknown, currentConfig, currentUser, localMachineSam, localMachineSecurity, localMachineSoftware, localMachineSystem, usersDefault.
// returns a *RegistryHive when successful
func (m *RegistryKeyState) GetHive()(*RegistryHive) {
    val, err := m.GetBackingStore().Get("hive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RegistryHive)
    }
    return nil
}
// GetKey gets the key property value. Current (i.e. changed) registry key (excludes HIVE).
// returns a *string when successful
func (m *RegistryKeyState) GetKey()(*string) {
    val, err := m.GetBackingStore().Get("key")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *RegistryKeyState) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOldKey gets the oldKey property value. Previous (i.e. before changed) registry key (excludes HIVE).
// returns a *string when successful
func (m *RegistryKeyState) GetOldKey()(*string) {
    val, err := m.GetBackingStore().Get("oldKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOldValueData gets the oldValueData property value. Previous (i.e. before changed) registry key value data (contents).
// returns a *string when successful
func (m *RegistryKeyState) GetOldValueData()(*string) {
    val, err := m.GetBackingStore().Get("oldValueData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOldValueName gets the oldValueName property value. Previous (i.e. before changed) registry key value name.
// returns a *string when successful
func (m *RegistryKeyState) GetOldValueName()(*string) {
    val, err := m.GetBackingStore().Get("oldValueName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperation gets the operation property value. Operation that changed the registry key name and/or value. Possible values are: unknown, create, modify, delete.
// returns a *RegistryOperation when successful
func (m *RegistryKeyState) GetOperation()(*RegistryOperation) {
    val, err := m.GetBackingStore().Get("operation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RegistryOperation)
    }
    return nil
}
// GetProcessId gets the processId property value. Process ID (PID) of the process that modified the registry key (process details will appear in the alert 'processes' collection).
// returns a *int32 when successful
func (m *RegistryKeyState) GetProcessId()(*int32) {
    val, err := m.GetBackingStore().Get("processId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetValueData gets the valueData property value. Current (i.e. changed) registry key value data (contents).
// returns a *string when successful
func (m *RegistryKeyState) GetValueData()(*string) {
    val, err := m.GetBackingStore().Get("valueData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetValueName gets the valueName property value. Current (i.e. changed) registry key value name
// returns a *string when successful
func (m *RegistryKeyState) GetValueName()(*string) {
    val, err := m.GetBackingStore().Get("valueName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetValueType gets the valueType property value. Registry key value type REGBINARY REGDWORD REGDWORDLITTLEENDIAN REGDWORDBIGENDIANREGEXPANDSZ REGLINK REGMULTISZ REGNONE REGQWORD REGQWORDLITTLEENDIAN REG_SZ Possible values are: unknown, binary, dword, dwordLittleEndian, dwordBigEndian, expandSz, link, multiSz, none, qword, qwordlittleEndian, sz.
// returns a *RegistryValueType when successful
func (m *RegistryKeyState) GetValueType()(*RegistryValueType) {
    val, err := m.GetBackingStore().Get("valueType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RegistryValueType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RegistryKeyState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetHive() != nil {
        cast := (*m.GetHive()).String()
        err := writer.WriteStringValue("hive", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("key", m.GetKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("oldKey", m.GetOldKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("oldValueData", m.GetOldValueData())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("oldValueName", m.GetOldValueName())
        if err != nil {
            return err
        }
    }
    if m.GetOperation() != nil {
        cast := (*m.GetOperation()).String()
        err := writer.WriteStringValue("operation", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("processId", m.GetProcessId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("valueData", m.GetValueData())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("valueName", m.GetValueName())
        if err != nil {
            return err
        }
    }
    if m.GetValueType() != nil {
        cast := (*m.GetValueType()).String()
        err := writer.WriteStringValue("valueType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RegistryKeyState) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *RegistryKeyState) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetHive sets the hive property value. A Windows registry hive : HKEYCURRENTCONFIG HKEYCURRENTUSER HKEYLOCALMACHINE/SAM HKEYLOCALMACHINE/Security HKEYLOCALMACHINE/Software HKEYLOCALMACHINE/System HKEY_USERS/.Default. Possible values are: unknown, currentConfig, currentUser, localMachineSam, localMachineSecurity, localMachineSoftware, localMachineSystem, usersDefault.
func (m *RegistryKeyState) SetHive(value *RegistryHive)() {
    err := m.GetBackingStore().Set("hive", value)
    if err != nil {
        panic(err)
    }
}
// SetKey sets the key property value. Current (i.e. changed) registry key (excludes HIVE).
func (m *RegistryKeyState) SetKey(value *string)() {
    err := m.GetBackingStore().Set("key", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RegistryKeyState) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOldKey sets the oldKey property value. Previous (i.e. before changed) registry key (excludes HIVE).
func (m *RegistryKeyState) SetOldKey(value *string)() {
    err := m.GetBackingStore().Set("oldKey", value)
    if err != nil {
        panic(err)
    }
}
// SetOldValueData sets the oldValueData property value. Previous (i.e. before changed) registry key value data (contents).
func (m *RegistryKeyState) SetOldValueData(value *string)() {
    err := m.GetBackingStore().Set("oldValueData", value)
    if err != nil {
        panic(err)
    }
}
// SetOldValueName sets the oldValueName property value. Previous (i.e. before changed) registry key value name.
func (m *RegistryKeyState) SetOldValueName(value *string)() {
    err := m.GetBackingStore().Set("oldValueName", value)
    if err != nil {
        panic(err)
    }
}
// SetOperation sets the operation property value. Operation that changed the registry key name and/or value. Possible values are: unknown, create, modify, delete.
func (m *RegistryKeyState) SetOperation(value *RegistryOperation)() {
    err := m.GetBackingStore().Set("operation", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessId sets the processId property value. Process ID (PID) of the process that modified the registry key (process details will appear in the alert 'processes' collection).
func (m *RegistryKeyState) SetProcessId(value *int32)() {
    err := m.GetBackingStore().Set("processId", value)
    if err != nil {
        panic(err)
    }
}
// SetValueData sets the valueData property value. Current (i.e. changed) registry key value data (contents).
func (m *RegistryKeyState) SetValueData(value *string)() {
    err := m.GetBackingStore().Set("valueData", value)
    if err != nil {
        panic(err)
    }
}
// SetValueName sets the valueName property value. Current (i.e. changed) registry key value name
func (m *RegistryKeyState) SetValueName(value *string)() {
    err := m.GetBackingStore().Set("valueName", value)
    if err != nil {
        panic(err)
    }
}
// SetValueType sets the valueType property value. Registry key value type REGBINARY REGDWORD REGDWORDLITTLEENDIAN REGDWORDBIGENDIANREGEXPANDSZ REGLINK REGMULTISZ REGNONE REGQWORD REGQWORDLITTLEENDIAN REG_SZ Possible values are: unknown, binary, dword, dwordLittleEndian, dwordBigEndian, expandSz, link, multiSz, none, qword, qwordlittleEndian, sz.
func (m *RegistryKeyState) SetValueType(value *RegistryValueType)() {
    err := m.GetBackingStore().Set("valueType", value)
    if err != nil {
        panic(err)
    }
}
type RegistryKeyStateable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetHive()(*RegistryHive)
    GetKey()(*string)
    GetOdataType()(*string)
    GetOldKey()(*string)
    GetOldValueData()(*string)
    GetOldValueName()(*string)
    GetOperation()(*RegistryOperation)
    GetProcessId()(*int32)
    GetValueData()(*string)
    GetValueName()(*string)
    GetValueType()(*RegistryValueType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetHive(value *RegistryHive)()
    SetKey(value *string)()
    SetOdataType(value *string)()
    SetOldKey(value *string)()
    SetOldValueData(value *string)()
    SetOldValueName(value *string)()
    SetOperation(value *RegistryOperation)()
    SetProcessId(value *int32)()
    SetValueData(value *string)()
    SetValueName(value *string)()
    SetValueType(value *RegistryValueType)()
}

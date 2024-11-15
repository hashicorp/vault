package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Process struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewProcess instantiates a new Process and sets the default values.
func NewProcess()(*Process) {
    m := &Process{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateProcessFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateProcessFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProcess(), nil
}
// GetAccountName gets the accountName property value. User account identifier (user account context the process ran under) for example, AccountName, SID, and so on.
// returns a *string when successful
func (m *Process) GetAccountName()(*string) {
    val, err := m.GetBackingStore().Get("accountName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Process) GetAdditionalData()(map[string]any) {
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
func (m *Process) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCommandLine gets the commandLine property value. The full process invocation commandline including all parameters.
// returns a *string when successful
func (m *Process) GetCommandLine()(*string) {
    val, err := m.GetBackingStore().Get("commandLine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Time at which the process was started. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Process) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
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
func (m *Process) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["commandLine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCommandLine(val)
        }
        return nil
    }
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
    res["fileHash"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileHashFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileHash(val.(FileHashable))
        }
        return nil
    }
    res["integrityLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseProcessIntegrityLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIntegrityLevel(val.(*ProcessIntegrityLevel))
        }
        return nil
    }
    res["isElevated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsElevated(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
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
    res["parentProcessCreatedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessCreatedDateTime(val)
        }
        return nil
    }
    res["parentProcessId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessId(val)
        }
        return nil
    }
    res["parentProcessName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessName(val)
        }
        return nil
    }
    res["path"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPath(val)
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
    return res
}
// GetFileHash gets the fileHash property value. Complex type containing file hashes (cryptographic and location-sensitive).
// returns a FileHashable when successful
func (m *Process) GetFileHash()(FileHashable) {
    val, err := m.GetBackingStore().Get("fileHash")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileHashable)
    }
    return nil
}
// GetIntegrityLevel gets the integrityLevel property value. The integrity level of the process. Possible values are: unknown, untrusted, low, medium, high, system.
// returns a *ProcessIntegrityLevel when successful
func (m *Process) GetIntegrityLevel()(*ProcessIntegrityLevel) {
    val, err := m.GetBackingStore().Get("integrityLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ProcessIntegrityLevel)
    }
    return nil
}
// GetIsElevated gets the isElevated property value. True if the process is elevated.
// returns a *bool when successful
func (m *Process) GetIsElevated()(*bool) {
    val, err := m.GetBackingStore().Get("isElevated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetName gets the name property value. The name of the process' Image file.
// returns a *string when successful
func (m *Process) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
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
func (m *Process) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetParentProcessCreatedDateTime gets the parentProcessCreatedDateTime property value. DateTime at which the parent process was started. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Process) GetParentProcessCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("parentProcessCreatedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetParentProcessId gets the parentProcessId property value. The Process ID (PID) of the parent process.
// returns a *int32 when successful
func (m *Process) GetParentProcessId()(*int32) {
    val, err := m.GetBackingStore().Get("parentProcessId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetParentProcessName gets the parentProcessName property value. The name of the image file of the parent process.
// returns a *string when successful
func (m *Process) GetParentProcessName()(*string) {
    val, err := m.GetBackingStore().Get("parentProcessName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPath gets the path property value. Full path, including filename.
// returns a *string when successful
func (m *Process) GetPath()(*string) {
    val, err := m.GetBackingStore().Get("path")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProcessId gets the processId property value. The Process ID (PID) of the process.
// returns a *int32 when successful
func (m *Process) GetProcessId()(*int32) {
    val, err := m.GetBackingStore().Get("processId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Process) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("accountName", m.GetAccountName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("commandLine", m.GetCommandLine())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("fileHash", m.GetFileHash())
        if err != nil {
            return err
        }
    }
    if m.GetIntegrityLevel() != nil {
        cast := (*m.GetIntegrityLevel()).String()
        err := writer.WriteStringValue("integrityLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isElevated", m.GetIsElevated())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
        err := writer.WriteTimeValue("parentProcessCreatedDateTime", m.GetParentProcessCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("parentProcessId", m.GetParentProcessId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("parentProcessName", m.GetParentProcessName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("path", m.GetPath())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountName sets the accountName property value. User account identifier (user account context the process ran under) for example, AccountName, SID, and so on.
func (m *Process) SetAccountName(value *string)() {
    err := m.GetBackingStore().Set("accountName", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Process) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Process) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCommandLine sets the commandLine property value. The full process invocation commandline including all parameters.
func (m *Process) SetCommandLine(value *string)() {
    err := m.GetBackingStore().Set("commandLine", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Time at which the process was started. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Process) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFileHash sets the fileHash property value. Complex type containing file hashes (cryptographic and location-sensitive).
func (m *Process) SetFileHash(value FileHashable)() {
    err := m.GetBackingStore().Set("fileHash", value)
    if err != nil {
        panic(err)
    }
}
// SetIntegrityLevel sets the integrityLevel property value. The integrity level of the process. Possible values are: unknown, untrusted, low, medium, high, system.
func (m *Process) SetIntegrityLevel(value *ProcessIntegrityLevel)() {
    err := m.GetBackingStore().Set("integrityLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetIsElevated sets the isElevated property value. True if the process is elevated.
func (m *Process) SetIsElevated(value *bool)() {
    err := m.GetBackingStore().Set("isElevated", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name of the process' Image file.
func (m *Process) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Process) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessCreatedDateTime sets the parentProcessCreatedDateTime property value. DateTime at which the parent process was started. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Process) SetParentProcessCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("parentProcessCreatedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessId sets the parentProcessId property value. The Process ID (PID) of the parent process.
func (m *Process) SetParentProcessId(value *int32)() {
    err := m.GetBackingStore().Set("parentProcessId", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessName sets the parentProcessName property value. The name of the image file of the parent process.
func (m *Process) SetParentProcessName(value *string)() {
    err := m.GetBackingStore().Set("parentProcessName", value)
    if err != nil {
        panic(err)
    }
}
// SetPath sets the path property value. Full path, including filename.
func (m *Process) SetPath(value *string)() {
    err := m.GetBackingStore().Set("path", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessId sets the processId property value. The Process ID (PID) of the process.
func (m *Process) SetProcessId(value *int32)() {
    err := m.GetBackingStore().Set("processId", value)
    if err != nil {
        panic(err)
    }
}
type Processable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountName()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCommandLine()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFileHash()(FileHashable)
    GetIntegrityLevel()(*ProcessIntegrityLevel)
    GetIsElevated()(*bool)
    GetName()(*string)
    GetOdataType()(*string)
    GetParentProcessCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetParentProcessId()(*int32)
    GetParentProcessName()(*string)
    GetPath()(*string)
    GetProcessId()(*int32)
    SetAccountName(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCommandLine(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFileHash(value FileHashable)()
    SetIntegrityLevel(value *ProcessIntegrityLevel)()
    SetIsElevated(value *bool)()
    SetName(value *string)()
    SetOdataType(value *string)()
    SetParentProcessCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetParentProcessId(value *int32)()
    SetParentProcessName(value *string)()
    SetPath(value *string)()
    SetProcessId(value *int32)()
}

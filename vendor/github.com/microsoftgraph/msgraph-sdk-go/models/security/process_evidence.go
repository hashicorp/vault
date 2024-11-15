package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ProcessEvidence struct {
    AlertEvidence
}
// NewProcessEvidence instantiates a new ProcessEvidence and sets the default values.
func NewProcessEvidence()(*ProcessEvidence) {
    m := &ProcessEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.processEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateProcessEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateProcessEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProcessEvidence(), nil
}
// GetDetectionStatus gets the detectionStatus property value. The status of the detection.The possible values are: detected, blocked, prevented, unknownFutureValue.
// returns a *DetectionStatus when successful
func (m *ProcessEvidence) GetDetectionStatus()(*DetectionStatus) {
    val, err := m.GetBackingStore().Get("detectionStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DetectionStatus)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ProcessEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["detectionStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDetectionStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetectionStatus(val.(*DetectionStatus))
        }
        return nil
    }
    res["imageFile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImageFile(val.(FileDetailsable))
        }
        return nil
    }
    res["mdeDeviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMdeDeviceId(val)
        }
        return nil
    }
    res["parentProcessCreationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessCreationDateTime(val)
        }
        return nil
    }
    res["parentProcessId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessId(val)
        }
        return nil
    }
    res["parentProcessImageFile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentProcessImageFile(val.(FileDetailsable))
        }
        return nil
    }
    res["processCommandLine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessCommandLine(val)
        }
        return nil
    }
    res["processCreationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessCreationDateTime(val)
        }
        return nil
    }
    res["processId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessId(val)
        }
        return nil
    }
    res["userAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserAccountFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAccount(val.(UserAccountable))
        }
        return nil
    }
    return res
}
// GetImageFile gets the imageFile property value. Image file details.
// returns a FileDetailsable when successful
func (m *ProcessEvidence) GetImageFile()(FileDetailsable) {
    val, err := m.GetBackingStore().Get("imageFile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileDetailsable)
    }
    return nil
}
// GetMdeDeviceId gets the mdeDeviceId property value. A unique identifier assigned to a device by Microsoft Defender for Endpoint.
// returns a *string when successful
func (m *ProcessEvidence) GetMdeDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("mdeDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetParentProcessCreationDateTime gets the parentProcessCreationDateTime property value. Date and time when the parent of the process was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *ProcessEvidence) GetParentProcessCreationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("parentProcessCreationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetParentProcessId gets the parentProcessId property value. Process ID (PID) of the parent process that spawned the process.
// returns a *int64 when successful
func (m *ProcessEvidence) GetParentProcessId()(*int64) {
    val, err := m.GetBackingStore().Get("parentProcessId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetParentProcessImageFile gets the parentProcessImageFile property value. Parent process image file details.
// returns a FileDetailsable when successful
func (m *ProcessEvidence) GetParentProcessImageFile()(FileDetailsable) {
    val, err := m.GetBackingStore().Get("parentProcessImageFile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileDetailsable)
    }
    return nil
}
// GetProcessCommandLine gets the processCommandLine property value. Command line used to create the new process.
// returns a *string when successful
func (m *ProcessEvidence) GetProcessCommandLine()(*string) {
    val, err := m.GetBackingStore().Get("processCommandLine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProcessCreationDateTime gets the processCreationDateTime property value. Date and time when the process was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *ProcessEvidence) GetProcessCreationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("processCreationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetProcessId gets the processId property value. Process ID (PID) of the newly created process.
// returns a *int64 when successful
func (m *ProcessEvidence) GetProcessId()(*int64) {
    val, err := m.GetBackingStore().Get("processId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUserAccount gets the userAccount property value. User details of the user that ran the process.
// returns a UserAccountable when successful
func (m *ProcessEvidence) GetUserAccount()(UserAccountable) {
    val, err := m.GetBackingStore().Get("userAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserAccountable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ProcessEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDetectionStatus() != nil {
        cast := (*m.GetDetectionStatus()).String()
        err = writer.WriteStringValue("detectionStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("imageFile", m.GetImageFile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mdeDeviceId", m.GetMdeDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("parentProcessCreationDateTime", m.GetParentProcessCreationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("parentProcessId", m.GetParentProcessId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentProcessImageFile", m.GetParentProcessImageFile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("processCommandLine", m.GetProcessCommandLine())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("processCreationDateTime", m.GetProcessCreationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("processId", m.GetProcessId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("userAccount", m.GetUserAccount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDetectionStatus sets the detectionStatus property value. The status of the detection.The possible values are: detected, blocked, prevented, unknownFutureValue.
func (m *ProcessEvidence) SetDetectionStatus(value *DetectionStatus)() {
    err := m.GetBackingStore().Set("detectionStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetImageFile sets the imageFile property value. Image file details.
func (m *ProcessEvidence) SetImageFile(value FileDetailsable)() {
    err := m.GetBackingStore().Set("imageFile", value)
    if err != nil {
        panic(err)
    }
}
// SetMdeDeviceId sets the mdeDeviceId property value. A unique identifier assigned to a device by Microsoft Defender for Endpoint.
func (m *ProcessEvidence) SetMdeDeviceId(value *string)() {
    err := m.GetBackingStore().Set("mdeDeviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessCreationDateTime sets the parentProcessCreationDateTime property value. Date and time when the parent of the process was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *ProcessEvidence) SetParentProcessCreationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("parentProcessCreationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessId sets the parentProcessId property value. Process ID (PID) of the parent process that spawned the process.
func (m *ProcessEvidence) SetParentProcessId(value *int64)() {
    err := m.GetBackingStore().Set("parentProcessId", value)
    if err != nil {
        panic(err)
    }
}
// SetParentProcessImageFile sets the parentProcessImageFile property value. Parent process image file details.
func (m *ProcessEvidence) SetParentProcessImageFile(value FileDetailsable)() {
    err := m.GetBackingStore().Set("parentProcessImageFile", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessCommandLine sets the processCommandLine property value. Command line used to create the new process.
func (m *ProcessEvidence) SetProcessCommandLine(value *string)() {
    err := m.GetBackingStore().Set("processCommandLine", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessCreationDateTime sets the processCreationDateTime property value. Date and time when the process was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *ProcessEvidence) SetProcessCreationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("processCreationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessId sets the processId property value. Process ID (PID) of the newly created process.
func (m *ProcessEvidence) SetProcessId(value *int64)() {
    err := m.GetBackingStore().Set("processId", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAccount sets the userAccount property value. User details of the user that ran the process.
func (m *ProcessEvidence) SetUserAccount(value UserAccountable)() {
    err := m.GetBackingStore().Set("userAccount", value)
    if err != nil {
        panic(err)
    }
}
type ProcessEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDetectionStatus()(*DetectionStatus)
    GetImageFile()(FileDetailsable)
    GetMdeDeviceId()(*string)
    GetParentProcessCreationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetParentProcessId()(*int64)
    GetParentProcessImageFile()(FileDetailsable)
    GetProcessCommandLine()(*string)
    GetProcessCreationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetProcessId()(*int64)
    GetUserAccount()(UserAccountable)
    SetDetectionStatus(value *DetectionStatus)()
    SetImageFile(value FileDetailsable)()
    SetMdeDeviceId(value *string)()
    SetParentProcessCreationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetParentProcessId(value *int64)()
    SetParentProcessImageFile(value FileDetailsable)()
    SetProcessCommandLine(value *string)()
    SetProcessCreationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetProcessId(value *int64)()
    SetUserAccount(value UserAccountable)()
}

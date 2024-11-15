package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FileEvidence struct {
    AlertEvidence
}
// NewFileEvidence instantiates a new FileEvidence and sets the default values.
func NewFileEvidence()(*FileEvidence) {
    m := &FileEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.fileEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateFileEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileEvidence(), nil
}
// GetDetectionStatus gets the detectionStatus property value. The status of the detection.The possible values are: detected, blocked, prevented, unknownFutureValue.
// returns a *DetectionStatus when successful
func (m *FileEvidence) GetDetectionStatus()(*DetectionStatus) {
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
func (m *FileEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["fileDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileDetails(val.(FileDetailsable))
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
    return res
}
// GetFileDetails gets the fileDetails property value. The file details.
// returns a FileDetailsable when successful
func (m *FileEvidence) GetFileDetails()(FileDetailsable) {
    val, err := m.GetBackingStore().Get("fileDetails")
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
func (m *FileEvidence) GetMdeDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("mdeDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FileEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("fileDetails", m.GetFileDetails())
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
    return nil
}
// SetDetectionStatus sets the detectionStatus property value. The status of the detection.The possible values are: detected, blocked, prevented, unknownFutureValue.
func (m *FileEvidence) SetDetectionStatus(value *DetectionStatus)() {
    err := m.GetBackingStore().Set("detectionStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetFileDetails sets the fileDetails property value. The file details.
func (m *FileEvidence) SetFileDetails(value FileDetailsable)() {
    err := m.GetBackingStore().Set("fileDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetMdeDeviceId sets the mdeDeviceId property value. A unique identifier assigned to a device by Microsoft Defender for Endpoint.
func (m *FileEvidence) SetMdeDeviceId(value *string)() {
    err := m.GetBackingStore().Set("mdeDeviceId", value)
    if err != nil {
        panic(err)
    }
}
type FileEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDetectionStatus()(*DetectionStatus)
    GetFileDetails()(FileDetailsable)
    GetMdeDeviceId()(*string)
    SetDetectionStatus(value *DetectionStatus)()
    SetFileDetails(value FileDetailsable)()
    SetMdeDeviceId(value *string)()
}

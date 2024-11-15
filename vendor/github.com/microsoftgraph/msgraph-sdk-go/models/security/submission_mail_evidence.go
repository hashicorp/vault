package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SubmissionMailEvidence struct {
    AlertEvidence
}
// NewSubmissionMailEvidence instantiates a new SubmissionMailEvidence and sets the default values.
func NewSubmissionMailEvidence()(*SubmissionMailEvidence) {
    m := &SubmissionMailEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.submissionMailEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSubmissionMailEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSubmissionMailEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSubmissionMailEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SubmissionMailEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["networkMessageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkMessageId(val)
        }
        return nil
    }
    res["recipient"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipient(val)
        }
        return nil
    }
    res["reportType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReportType(val)
        }
        return nil
    }
    res["sender"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSender(val)
        }
        return nil
    }
    res["senderIp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSenderIp(val)
        }
        return nil
    }
    res["subject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubject(val)
        }
        return nil
    }
    res["submissionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubmissionDateTime(val)
        }
        return nil
    }
    res["submissionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubmissionId(val)
        }
        return nil
    }
    res["submitter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubmitter(val)
        }
        return nil
    }
    return res
}
// GetNetworkMessageId gets the networkMessageId property value. The networkMessageId property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetNetworkMessageId()(*string) {
    val, err := m.GetBackingStore().Get("networkMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecipient gets the recipient property value. The recipient property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetRecipient()(*string) {
    val, err := m.GetBackingStore().Get("recipient")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReportType gets the reportType property value. The reportType property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetReportType()(*string) {
    val, err := m.GetBackingStore().Get("reportType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSender gets the sender property value. The sender property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetSender()(*string) {
    val, err := m.GetBackingStore().Get("sender")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSenderIp gets the senderIp property value. The senderIp property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetSenderIp()(*string) {
    val, err := m.GetBackingStore().Get("senderIp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubject gets the subject property value. The subject property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubmissionDateTime gets the submissionDateTime property value. The submissionDateTime property
// returns a *Time when successful
func (m *SubmissionMailEvidence) GetSubmissionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("submissionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSubmissionId gets the submissionId property value. The submissionId property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetSubmissionId()(*string) {
    val, err := m.GetBackingStore().Get("submissionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubmitter gets the submitter property value. The submitter property
// returns a *string when successful
func (m *SubmissionMailEvidence) GetSubmitter()(*string) {
    val, err := m.GetBackingStore().Get("submitter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SubmissionMailEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("networkMessageId", m.GetNetworkMessageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recipient", m.GetRecipient())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("reportType", m.GetReportType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sender", m.GetSender())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("senderIp", m.GetSenderIp())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("submissionDateTime", m.GetSubmissionDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("submissionId", m.GetSubmissionId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("submitter", m.GetSubmitter())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetNetworkMessageId sets the networkMessageId property value. The networkMessageId property
func (m *SubmissionMailEvidence) SetNetworkMessageId(value *string)() {
    err := m.GetBackingStore().Set("networkMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipient sets the recipient property value. The recipient property
func (m *SubmissionMailEvidence) SetRecipient(value *string)() {
    err := m.GetBackingStore().Set("recipient", value)
    if err != nil {
        panic(err)
    }
}
// SetReportType sets the reportType property value. The reportType property
func (m *SubmissionMailEvidence) SetReportType(value *string)() {
    err := m.GetBackingStore().Set("reportType", value)
    if err != nil {
        panic(err)
    }
}
// SetSender sets the sender property value. The sender property
func (m *SubmissionMailEvidence) SetSender(value *string)() {
    err := m.GetBackingStore().Set("sender", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderIp sets the senderIp property value. The senderIp property
func (m *SubmissionMailEvidence) SetSenderIp(value *string)() {
    err := m.GetBackingStore().Set("senderIp", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject property
func (m *SubmissionMailEvidence) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetSubmissionDateTime sets the submissionDateTime property value. The submissionDateTime property
func (m *SubmissionMailEvidence) SetSubmissionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("submissionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSubmissionId sets the submissionId property value. The submissionId property
func (m *SubmissionMailEvidence) SetSubmissionId(value *string)() {
    err := m.GetBackingStore().Set("submissionId", value)
    if err != nil {
        panic(err)
    }
}
// SetSubmitter sets the submitter property value. The submitter property
func (m *SubmissionMailEvidence) SetSubmitter(value *string)() {
    err := m.GetBackingStore().Set("submitter", value)
    if err != nil {
        panic(err)
    }
}
type SubmissionMailEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetNetworkMessageId()(*string)
    GetRecipient()(*string)
    GetReportType()(*string)
    GetSender()(*string)
    GetSenderIp()(*string)
    GetSubject()(*string)
    GetSubmissionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSubmissionId()(*string)
    GetSubmitter()(*string)
    SetNetworkMessageId(value *string)()
    SetRecipient(value *string)()
    SetReportType(value *string)()
    SetSender(value *string)()
    SetSenderIp(value *string)()
    SetSubject(value *string)()
    SetSubmissionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSubmissionId(value *string)()
    SetSubmitter(value *string)()
}

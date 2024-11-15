package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AnalyzedMessageEvidence struct {
    AlertEvidence
}
// NewAnalyzedMessageEvidence instantiates a new AnalyzedMessageEvidence and sets the default values.
func NewAnalyzedMessageEvidence()(*AnalyzedMessageEvidence) {
    m := &AnalyzedMessageEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.analyzedMessageEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAnalyzedMessageEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAnalyzedMessageEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAnalyzedMessageEvidence(), nil
}
// GetAntiSpamDirection gets the antiSpamDirection property value. Direction of the email relative to your network. The possible values are: inbound, outbound or intraorg.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetAntiSpamDirection()(*string) {
    val, err := m.GetBackingStore().Get("antiSpamDirection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAttachmentsCount gets the attachmentsCount property value. Number of attachments in the email.
// returns a *int64 when successful
func (m *AnalyzedMessageEvidence) GetAttachmentsCount()(*int64) {
    val, err := m.GetBackingStore().Get("attachmentsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetDeliveryAction gets the deliveryAction property value. Delivery action of the email. The possible values are: delivered, deliveredAsSpam, junked, blocked, or replaced.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetDeliveryAction()(*string) {
    val, err := m.GetBackingStore().Get("deliveryAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeliveryLocation gets the deliveryLocation property value. Location where the email was delivered. The possible values are: inbox, external, junkFolder, quarantine, failed, dropped, deletedFolder or forwarded.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetDeliveryLocation()(*string) {
    val, err := m.GetBackingStore().Get("deliveryLocation")
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
func (m *AnalyzedMessageEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["antiSpamDirection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAntiSpamDirection(val)
        }
        return nil
    }
    res["attachmentsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttachmentsCount(val)
        }
        return nil
    }
    res["deliveryAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryAction(val)
        }
        return nil
    }
    res["deliveryLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryLocation(val)
        }
        return nil
    }
    res["internetMessageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternetMessageId(val)
        }
        return nil
    }
    res["language"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLanguage(val)
        }
        return nil
    }
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
    res["p1Sender"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailSenderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetP1Sender(val.(EmailSenderable))
        }
        return nil
    }
    res["p2Sender"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailSenderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetP2Sender(val.(EmailSenderable))
        }
        return nil
    }
    res["receivedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReceivedDateTime(val)
        }
        return nil
    }
    res["recipientEmailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipientEmailAddress(val)
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
    res["threatDetectionMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetThreatDetectionMethods(res)
        }
        return nil
    }
    res["threats"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetThreats(res)
        }
        return nil
    }
    res["urlCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrlCount(val)
        }
        return nil
    }
    res["urls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetUrls(res)
        }
        return nil
    }
    res["urn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrn(val)
        }
        return nil
    }
    return res
}
// GetInternetMessageId gets the internetMessageId property value. Public-facing identifier for the email that is set by the sending email system.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetInternetMessageId()(*string) {
    val, err := m.GetBackingStore().Get("internetMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLanguage gets the language property value. Detected language of the email content.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetLanguage()(*string) {
    val, err := m.GetBackingStore().Get("language")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNetworkMessageId gets the networkMessageId property value. Unique identifier for the email, generated by Microsoft 365.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetNetworkMessageId()(*string) {
    val, err := m.GetBackingStore().Get("networkMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetP1Sender gets the p1Sender property value. The P1 sender.
// returns a EmailSenderable when successful
func (m *AnalyzedMessageEvidence) GetP1Sender()(EmailSenderable) {
    val, err := m.GetBackingStore().Get("p1Sender")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailSenderable)
    }
    return nil
}
// GetP2Sender gets the p2Sender property value. The P2 sender.
// returns a EmailSenderable when successful
func (m *AnalyzedMessageEvidence) GetP2Sender()(EmailSenderable) {
    val, err := m.GetBackingStore().Get("p2Sender")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailSenderable)
    }
    return nil
}
// GetReceivedDateTime gets the receivedDateTime property value. Date and time when the email was received.
// returns a *Time when successful
func (m *AnalyzedMessageEvidence) GetReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("receivedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRecipientEmailAddress gets the recipientEmailAddress property value. Email address of the recipient, or email address of the recipient after distribution list expansion.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetRecipientEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("recipientEmailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSenderIp gets the senderIp property value. IP address of the last detected mail server that relayed the message.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetSenderIp()(*string) {
    val, err := m.GetBackingStore().Get("senderIp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubject gets the subject property value. Subject of the email.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetThreatDetectionMethods gets the threatDetectionMethods property value. Collection of methods used to detect malware, phishing, or other threats found in the email.
// returns a []string when successful
func (m *AnalyzedMessageEvidence) GetThreatDetectionMethods()([]string) {
    val, err := m.GetBackingStore().Get("threatDetectionMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetThreats gets the threats property value. Collection of detection names for malware or other threats found.
// returns a []string when successful
func (m *AnalyzedMessageEvidence) GetThreats()([]string) {
    val, err := m.GetBackingStore().Get("threats")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetUrlCount gets the urlCount property value. Number of embedded URLs in the email.
// returns a *int64 when successful
func (m *AnalyzedMessageEvidence) GetUrlCount()(*int64) {
    val, err := m.GetBackingStore().Get("urlCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUrls gets the urls property value. Collection of the URLs contained in this email.
// returns a []string when successful
func (m *AnalyzedMessageEvidence) GetUrls()([]string) {
    val, err := m.GetBackingStore().Get("urls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetUrn gets the urn property value. Uniform resource name (URN) of the automated investigation where the cluster was identified.
// returns a *string when successful
func (m *AnalyzedMessageEvidence) GetUrn()(*string) {
    val, err := m.GetBackingStore().Get("urn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AnalyzedMessageEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("antiSpamDirection", m.GetAntiSpamDirection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("attachmentsCount", m.GetAttachmentsCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deliveryAction", m.GetDeliveryAction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deliveryLocation", m.GetDeliveryLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("internetMessageId", m.GetInternetMessageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("language", m.GetLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("networkMessageId", m.GetNetworkMessageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("p1Sender", m.GetP1Sender())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("p2Sender", m.GetP2Sender())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("receivedDateTime", m.GetReceivedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recipientEmailAddress", m.GetRecipientEmailAddress())
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
    if m.GetThreatDetectionMethods() != nil {
        err = writer.WriteCollectionOfStringValues("threatDetectionMethods", m.GetThreatDetectionMethods())
        if err != nil {
            return err
        }
    }
    if m.GetThreats() != nil {
        err = writer.WriteCollectionOfStringValues("threats", m.GetThreats())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("urlCount", m.GetUrlCount())
        if err != nil {
            return err
        }
    }
    if m.GetUrls() != nil {
        err = writer.WriteCollectionOfStringValues("urls", m.GetUrls())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("urn", m.GetUrn())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAntiSpamDirection sets the antiSpamDirection property value. Direction of the email relative to your network. The possible values are: inbound, outbound or intraorg.
func (m *AnalyzedMessageEvidence) SetAntiSpamDirection(value *string)() {
    err := m.GetBackingStore().Set("antiSpamDirection", value)
    if err != nil {
        panic(err)
    }
}
// SetAttachmentsCount sets the attachmentsCount property value. Number of attachments in the email.
func (m *AnalyzedMessageEvidence) SetAttachmentsCount(value *int64)() {
    err := m.GetBackingStore().Set("attachmentsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryAction sets the deliveryAction property value. Delivery action of the email. The possible values are: delivered, deliveredAsSpam, junked, blocked, or replaced.
func (m *AnalyzedMessageEvidence) SetDeliveryAction(value *string)() {
    err := m.GetBackingStore().Set("deliveryAction", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryLocation sets the deliveryLocation property value. Location where the email was delivered. The possible values are: inbox, external, junkFolder, quarantine, failed, dropped, deletedFolder or forwarded.
func (m *AnalyzedMessageEvidence) SetDeliveryLocation(value *string)() {
    err := m.GetBackingStore().Set("deliveryLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetInternetMessageId sets the internetMessageId property value. Public-facing identifier for the email that is set by the sending email system.
func (m *AnalyzedMessageEvidence) SetInternetMessageId(value *string)() {
    err := m.GetBackingStore().Set("internetMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguage sets the language property value. Detected language of the email content.
func (m *AnalyzedMessageEvidence) SetLanguage(value *string)() {
    err := m.GetBackingStore().Set("language", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkMessageId sets the networkMessageId property value. Unique identifier for the email, generated by Microsoft 365.
func (m *AnalyzedMessageEvidence) SetNetworkMessageId(value *string)() {
    err := m.GetBackingStore().Set("networkMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetP1Sender sets the p1Sender property value. The P1 sender.
func (m *AnalyzedMessageEvidence) SetP1Sender(value EmailSenderable)() {
    err := m.GetBackingStore().Set("p1Sender", value)
    if err != nil {
        panic(err)
    }
}
// SetP2Sender sets the p2Sender property value. The P2 sender.
func (m *AnalyzedMessageEvidence) SetP2Sender(value EmailSenderable)() {
    err := m.GetBackingStore().Set("p2Sender", value)
    if err != nil {
        panic(err)
    }
}
// SetReceivedDateTime sets the receivedDateTime property value. Date and time when the email was received.
func (m *AnalyzedMessageEvidence) SetReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("receivedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientEmailAddress sets the recipientEmailAddress property value. Email address of the recipient, or email address of the recipient after distribution list expansion.
func (m *AnalyzedMessageEvidence) SetRecipientEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("recipientEmailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderIp sets the senderIp property value. IP address of the last detected mail server that relayed the message.
func (m *AnalyzedMessageEvidence) SetSenderIp(value *string)() {
    err := m.GetBackingStore().Set("senderIp", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. Subject of the email.
func (m *AnalyzedMessageEvidence) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetThreatDetectionMethods sets the threatDetectionMethods property value. Collection of methods used to detect malware, phishing, or other threats found in the email.
func (m *AnalyzedMessageEvidence) SetThreatDetectionMethods(value []string)() {
    err := m.GetBackingStore().Set("threatDetectionMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetThreats sets the threats property value. Collection of detection names for malware or other threats found.
func (m *AnalyzedMessageEvidence) SetThreats(value []string)() {
    err := m.GetBackingStore().Set("threats", value)
    if err != nil {
        panic(err)
    }
}
// SetUrlCount sets the urlCount property value. Number of embedded URLs in the email.
func (m *AnalyzedMessageEvidence) SetUrlCount(value *int64)() {
    err := m.GetBackingStore().Set("urlCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUrls sets the urls property value. Collection of the URLs contained in this email.
func (m *AnalyzedMessageEvidence) SetUrls(value []string)() {
    err := m.GetBackingStore().Set("urls", value)
    if err != nil {
        panic(err)
    }
}
// SetUrn sets the urn property value. Uniform resource name (URN) of the automated investigation where the cluster was identified.
func (m *AnalyzedMessageEvidence) SetUrn(value *string)() {
    err := m.GetBackingStore().Set("urn", value)
    if err != nil {
        panic(err)
    }
}
type AnalyzedMessageEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAntiSpamDirection()(*string)
    GetAttachmentsCount()(*int64)
    GetDeliveryAction()(*string)
    GetDeliveryLocation()(*string)
    GetInternetMessageId()(*string)
    GetLanguage()(*string)
    GetNetworkMessageId()(*string)
    GetP1Sender()(EmailSenderable)
    GetP2Sender()(EmailSenderable)
    GetReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRecipientEmailAddress()(*string)
    GetSenderIp()(*string)
    GetSubject()(*string)
    GetThreatDetectionMethods()([]string)
    GetThreats()([]string)
    GetUrlCount()(*int64)
    GetUrls()([]string)
    GetUrn()(*string)
    SetAntiSpamDirection(value *string)()
    SetAttachmentsCount(value *int64)()
    SetDeliveryAction(value *string)()
    SetDeliveryLocation(value *string)()
    SetInternetMessageId(value *string)()
    SetLanguage(value *string)()
    SetNetworkMessageId(value *string)()
    SetP1Sender(value EmailSenderable)()
    SetP2Sender(value EmailSenderable)()
    SetReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRecipientEmailAddress(value *string)()
    SetSenderIp(value *string)()
    SetSubject(value *string)()
    SetThreatDetectionMethods(value []string)()
    SetThreats(value []string)()
    SetUrlCount(value *int64)()
    SetUrls(value []string)()
    SetUrn(value *string)()
}

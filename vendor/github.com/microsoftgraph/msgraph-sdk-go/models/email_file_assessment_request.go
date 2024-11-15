package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EmailFileAssessmentRequest struct {
    ThreatAssessmentRequest
}
// NewEmailFileAssessmentRequest instantiates a new EmailFileAssessmentRequest and sets the default values.
func NewEmailFileAssessmentRequest()(*EmailFileAssessmentRequest) {
    m := &EmailFileAssessmentRequest{
        ThreatAssessmentRequest: *NewThreatAssessmentRequest(),
    }
    odataTypeValue := "#microsoft.graph.emailFileAssessmentRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEmailFileAssessmentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEmailFileAssessmentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmailFileAssessmentRequest(), nil
}
// GetContentData gets the contentData property value. Base64 encoded .eml email file content. The file content can't fetch back because it isn't stored.
// returns a *string when successful
func (m *EmailFileAssessmentRequest) GetContentData()(*string) {
    val, err := m.GetBackingStore().Get("contentData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDestinationRoutingReason gets the destinationRoutingReason property value. The reason for mail routed to its destination. Possible values are: none, mailFlowRule, safeSender, blockedSender, advancedSpamFiltering, domainAllowList, domainBlockList, notInAddressBook, firstTimeSender, autoPurgeToInbox, autoPurgeToJunk, autoPurgeToDeleted, outbound, notJunk, junk.
// returns a *MailDestinationRoutingReason when successful
func (m *EmailFileAssessmentRequest) GetDestinationRoutingReason()(*MailDestinationRoutingReason) {
    val, err := m.GetBackingStore().Get("destinationRoutingReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MailDestinationRoutingReason)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EmailFileAssessmentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ThreatAssessmentRequest.GetFieldDeserializers()
    res["contentData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentData(val)
        }
        return nil
    }
    res["destinationRoutingReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMailDestinationRoutingReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationRoutingReason(val.(*MailDestinationRoutingReason))
        }
        return nil
    }
    res["recipientEmail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipientEmail(val)
        }
        return nil
    }
    return res
}
// GetRecipientEmail gets the recipientEmail property value. The mail recipient whose policies are used to assess the mail.
// returns a *string when successful
func (m *EmailFileAssessmentRequest) GetRecipientEmail()(*string) {
    val, err := m.GetBackingStore().Get("recipientEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EmailFileAssessmentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ThreatAssessmentRequest.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("contentData", m.GetContentData())
        if err != nil {
            return err
        }
    }
    if m.GetDestinationRoutingReason() != nil {
        cast := (*m.GetDestinationRoutingReason()).String()
        err = writer.WriteStringValue("destinationRoutingReason", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recipientEmail", m.GetRecipientEmail())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContentData sets the contentData property value. Base64 encoded .eml email file content. The file content can't fetch back because it isn't stored.
func (m *EmailFileAssessmentRequest) SetContentData(value *string)() {
    err := m.GetBackingStore().Set("contentData", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationRoutingReason sets the destinationRoutingReason property value. The reason for mail routed to its destination. Possible values are: none, mailFlowRule, safeSender, blockedSender, advancedSpamFiltering, domainAllowList, domainBlockList, notInAddressBook, firstTimeSender, autoPurgeToInbox, autoPurgeToJunk, autoPurgeToDeleted, outbound, notJunk, junk.
func (m *EmailFileAssessmentRequest) SetDestinationRoutingReason(value *MailDestinationRoutingReason)() {
    err := m.GetBackingStore().Set("destinationRoutingReason", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientEmail sets the recipientEmail property value. The mail recipient whose policies are used to assess the mail.
func (m *EmailFileAssessmentRequest) SetRecipientEmail(value *string)() {
    err := m.GetBackingStore().Set("recipientEmail", value)
    if err != nil {
        panic(err)
    }
}
type EmailFileAssessmentRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ThreatAssessmentRequestable
    GetContentData()(*string)
    GetDestinationRoutingReason()(*MailDestinationRoutingReason)
    GetRecipientEmail()(*string)
    SetContentData(value *string)()
    SetDestinationRoutingReason(value *MailDestinationRoutingReason)()
    SetRecipientEmail(value *string)()
}

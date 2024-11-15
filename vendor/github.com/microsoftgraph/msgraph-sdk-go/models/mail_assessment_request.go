package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailAssessmentRequest struct {
    ThreatAssessmentRequest
}
// NewMailAssessmentRequest instantiates a new MailAssessmentRequest and sets the default values.
func NewMailAssessmentRequest()(*MailAssessmentRequest) {
    m := &MailAssessmentRequest{
        ThreatAssessmentRequest: *NewThreatAssessmentRequest(),
    }
    odataTypeValue := "#microsoft.graph.mailAssessmentRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMailAssessmentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailAssessmentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailAssessmentRequest(), nil
}
// GetDestinationRoutingReason gets the destinationRoutingReason property value. The reason for mail routed to its destination. Possible values are: none, mailFlowRule, safeSender, blockedSender, advancedSpamFiltering, domainAllowList, domainBlockList, notInAddressBook, firstTimeSender, autoPurgeToInbox, autoPurgeToJunk, autoPurgeToDeleted, outbound, notJunk, junk.
// returns a *MailDestinationRoutingReason when successful
func (m *MailAssessmentRequest) GetDestinationRoutingReason()(*MailDestinationRoutingReason) {
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
func (m *MailAssessmentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ThreatAssessmentRequest.GetFieldDeserializers()
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
    res["messageUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageUri(val)
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
// GetMessageUri gets the messageUri property value. The resource URI of the mail message for assessment.
// returns a *string when successful
func (m *MailAssessmentRequest) GetMessageUri()(*string) {
    val, err := m.GetBackingStore().Get("messageUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecipientEmail gets the recipientEmail property value. The mail recipient whose policies are used to assess the mail.
// returns a *string when successful
func (m *MailAssessmentRequest) GetRecipientEmail()(*string) {
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
func (m *MailAssessmentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ThreatAssessmentRequest.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDestinationRoutingReason() != nil {
        cast := (*m.GetDestinationRoutingReason()).String()
        err = writer.WriteStringValue("destinationRoutingReason", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("messageUri", m.GetMessageUri())
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
// SetDestinationRoutingReason sets the destinationRoutingReason property value. The reason for mail routed to its destination. Possible values are: none, mailFlowRule, safeSender, blockedSender, advancedSpamFiltering, domainAllowList, domainBlockList, notInAddressBook, firstTimeSender, autoPurgeToInbox, autoPurgeToJunk, autoPurgeToDeleted, outbound, notJunk, junk.
func (m *MailAssessmentRequest) SetDestinationRoutingReason(value *MailDestinationRoutingReason)() {
    err := m.GetBackingStore().Set("destinationRoutingReason", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageUri sets the messageUri property value. The resource URI of the mail message for assessment.
func (m *MailAssessmentRequest) SetMessageUri(value *string)() {
    err := m.GetBackingStore().Set("messageUri", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientEmail sets the recipientEmail property value. The mail recipient whose policies are used to assess the mail.
func (m *MailAssessmentRequest) SetRecipientEmail(value *string)() {
    err := m.GetBackingStore().Set("recipientEmail", value)
    if err != nil {
        panic(err)
    }
}
type MailAssessmentRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ThreatAssessmentRequestable
    GetDestinationRoutingReason()(*MailDestinationRoutingReason)
    GetMessageUri()(*string)
    GetRecipientEmail()(*string)
    SetDestinationRoutingReason(value *MailDestinationRoutingReason)()
    SetMessageUri(value *string)()
    SetRecipientEmail(value *string)()
}

package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailboxEvidence struct {
    AlertEvidence
}
// NewMailboxEvidence instantiates a new MailboxEvidence and sets the default values.
func NewMailboxEvidence()(*MailboxEvidence) {
    m := &MailboxEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.mailboxEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMailboxEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailboxEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailboxEvidence(), nil
}
// GetDisplayName gets the displayName property value. The name associated with the mailbox.
// returns a *string when successful
func (m *MailboxEvidence) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *MailboxEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["primaryAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrimaryAddress(val)
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
// GetPrimaryAddress gets the primaryAddress property value. The primary email address of the mailbox.
// returns a *string when successful
func (m *MailboxEvidence) GetPrimaryAddress()(*string) {
    val, err := m.GetBackingStore().Get("primaryAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserAccount gets the userAccount property value. The user account of the mailbox.
// returns a UserAccountable when successful
func (m *MailboxEvidence) GetUserAccount()(UserAccountable) {
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
func (m *MailboxEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("primaryAddress", m.GetPrimaryAddress())
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
// SetDisplayName sets the displayName property value. The name associated with the mailbox.
func (m *MailboxEvidence) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryAddress sets the primaryAddress property value. The primary email address of the mailbox.
func (m *MailboxEvidence) SetPrimaryAddress(value *string)() {
    err := m.GetBackingStore().Set("primaryAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAccount sets the userAccount property value. The user account of the mailbox.
func (m *MailboxEvidence) SetUserAccount(value UserAccountable)() {
    err := m.GetBackingStore().Set("userAccount", value)
    if err != nil {
        panic(err)
    }
}
type MailboxEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetPrimaryAddress()(*string)
    GetUserAccount()(UserAccountable)
    SetDisplayName(value *string)()
    SetPrimaryAddress(value *string)()
    SetUserAccount(value UserAccountable)()
}

package security

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailboxConfigurationEvidence struct {
    AlertEvidence
}
// NewMailboxConfigurationEvidence instantiates a new MailboxConfigurationEvidence and sets the default values.
func NewMailboxConfigurationEvidence()(*MailboxConfigurationEvidence) {
    m := &MailboxConfigurationEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.mailboxConfigurationEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMailboxConfigurationEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailboxConfigurationEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailboxConfigurationEvidence(), nil
}
// GetConfigurationId gets the configurationId property value. The configurationId property
// returns a *string when successful
func (m *MailboxConfigurationEvidence) GetConfigurationId()(*string) {
    val, err := m.GetBackingStore().Get("configurationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConfigurationType gets the configurationType property value. The configurationType property
// returns a *MailboxConfigurationType when successful
func (m *MailboxConfigurationEvidence) GetConfigurationType()(*MailboxConfigurationType) {
    val, err := m.GetBackingStore().Get("configurationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MailboxConfigurationType)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The displayName property
// returns a *string when successful
func (m *MailboxConfigurationEvidence) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalDirectoryObjectId gets the externalDirectoryObjectId property value. The externalDirectoryObjectId property
// returns a *UUID when successful
func (m *MailboxConfigurationEvidence) GetExternalDirectoryObjectId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("externalDirectoryObjectId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailboxConfigurationEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["configurationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationId(val)
        }
        return nil
    }
    res["configurationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMailboxConfigurationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationType(val.(*MailboxConfigurationType))
        }
        return nil
    }
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
    res["externalDirectoryObjectId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalDirectoryObjectId(val)
        }
        return nil
    }
    res["mailboxPrimaryAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailboxPrimaryAddress(val)
        }
        return nil
    }
    res["upn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpn(val)
        }
        return nil
    }
    return res
}
// GetMailboxPrimaryAddress gets the mailboxPrimaryAddress property value. The mailboxPrimaryAddress property
// returns a *string when successful
func (m *MailboxConfigurationEvidence) GetMailboxPrimaryAddress()(*string) {
    val, err := m.GetBackingStore().Get("mailboxPrimaryAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUpn gets the upn property value. The upn property
// returns a *string when successful
func (m *MailboxConfigurationEvidence) GetUpn()(*string) {
    val, err := m.GetBackingStore().Get("upn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailboxConfigurationEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("configurationId", m.GetConfigurationId())
        if err != nil {
            return err
        }
    }
    if m.GetConfigurationType() != nil {
        cast := (*m.GetConfigurationType()).String()
        err = writer.WriteStringValue("configurationType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("externalDirectoryObjectId", m.GetExternalDirectoryObjectId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mailboxPrimaryAddress", m.GetMailboxPrimaryAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("upn", m.GetUpn())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfigurationId sets the configurationId property value. The configurationId property
func (m *MailboxConfigurationEvidence) SetConfigurationId(value *string)() {
    err := m.GetBackingStore().Set("configurationId", value)
    if err != nil {
        panic(err)
    }
}
// SetConfigurationType sets the configurationType property value. The configurationType property
func (m *MailboxConfigurationEvidence) SetConfigurationType(value *MailboxConfigurationType)() {
    err := m.GetBackingStore().Set("configurationType", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The displayName property
func (m *MailboxConfigurationEvidence) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalDirectoryObjectId sets the externalDirectoryObjectId property value. The externalDirectoryObjectId property
func (m *MailboxConfigurationEvidence) SetExternalDirectoryObjectId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("externalDirectoryObjectId", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxPrimaryAddress sets the mailboxPrimaryAddress property value. The mailboxPrimaryAddress property
func (m *MailboxConfigurationEvidence) SetMailboxPrimaryAddress(value *string)() {
    err := m.GetBackingStore().Set("mailboxPrimaryAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetUpn sets the upn property value. The upn property
func (m *MailboxConfigurationEvidence) SetUpn(value *string)() {
    err := m.GetBackingStore().Set("upn", value)
    if err != nil {
        panic(err)
    }
}
type MailboxConfigurationEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConfigurationId()(*string)
    GetConfigurationType()(*MailboxConfigurationType)
    GetDisplayName()(*string)
    GetExternalDirectoryObjectId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetMailboxPrimaryAddress()(*string)
    GetUpn()(*string)
    SetConfigurationId(value *string)()
    SetConfigurationType(value *MailboxConfigurationType)()
    SetDisplayName(value *string)()
    SetExternalDirectoryObjectId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetMailboxPrimaryAddress(value *string)()
    SetUpn(value *string)()
}

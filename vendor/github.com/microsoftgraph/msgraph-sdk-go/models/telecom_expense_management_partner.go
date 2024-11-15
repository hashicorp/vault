package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TelecomExpenseManagementPartner telecomExpenseManagementPartner resources represent the metadata and status of a given TEM service. Once your organization has onboarded with a partner, the partner can be enabled or disabled to switch TEM functionality on or off.
type TelecomExpenseManagementPartner struct {
    Entity
}
// NewTelecomExpenseManagementPartner instantiates a new TelecomExpenseManagementPartner and sets the default values.
func NewTelecomExpenseManagementPartner()(*TelecomExpenseManagementPartner) {
    m := &TelecomExpenseManagementPartner{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTelecomExpenseManagementPartnerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTelecomExpenseManagementPartnerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTelecomExpenseManagementPartner(), nil
}
// GetAppAuthorized gets the appAuthorized property value. Whether the partner's AAD app has been authorized to access Intune.
// returns a *bool when successful
func (m *TelecomExpenseManagementPartner) GetAppAuthorized()(*bool) {
    val, err := m.GetBackingStore().Get("appAuthorized")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name of the TEM partner.
// returns a *string when successful
func (m *TelecomExpenseManagementPartner) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnabled gets the enabled property value. Whether Intune's connection to the TEM service is currently enabled or disabled.
// returns a *bool when successful
func (m *TelecomExpenseManagementPartner) GetEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("enabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TelecomExpenseManagementPartner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appAuthorized"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppAuthorized(val)
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
    res["enabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnabled(val)
        }
        return nil
    }
    res["lastConnectionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastConnectionDateTime(val)
        }
        return nil
    }
    res["url"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrl(val)
        }
        return nil
    }
    return res
}
// GetLastConnectionDateTime gets the lastConnectionDateTime property value. Timestamp of the last request sent to Intune by the TEM partner.
// returns a *Time when successful
func (m *TelecomExpenseManagementPartner) GetLastConnectionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastConnectionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetUrl gets the url property value. URL of the TEM partner's administrative control panel, where an administrator can configure their TEM service.
// returns a *string when successful
func (m *TelecomExpenseManagementPartner) GetUrl()(*string) {
    val, err := m.GetBackingStore().Get("url")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TelecomExpenseManagementPartner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("appAuthorized", m.GetAppAuthorized())
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
        err = writer.WriteBoolValue("enabled", m.GetEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastConnectionDateTime", m.GetLastConnectionDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("url", m.GetUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppAuthorized sets the appAuthorized property value. Whether the partner's AAD app has been authorized to access Intune.
func (m *TelecomExpenseManagementPartner) SetAppAuthorized(value *bool)() {
    err := m.GetBackingStore().Set("appAuthorized", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name of the TEM partner.
func (m *TelecomExpenseManagementPartner) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEnabled sets the enabled property value. Whether Intune's connection to the TEM service is currently enabled or disabled.
func (m *TelecomExpenseManagementPartner) SetEnabled(value *bool)() {
    err := m.GetBackingStore().Set("enabled", value)
    if err != nil {
        panic(err)
    }
}
// SetLastConnectionDateTime sets the lastConnectionDateTime property value. Timestamp of the last request sent to Intune by the TEM partner.
func (m *TelecomExpenseManagementPartner) SetLastConnectionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastConnectionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetUrl sets the url property value. URL of the TEM partner's administrative control panel, where an administrator can configure their TEM service.
func (m *TelecomExpenseManagementPartner) SetUrl(value *string)() {
    err := m.GetBackingStore().Set("url", value)
    if err != nil {
        panic(err)
    }
}
type TelecomExpenseManagementPartnerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppAuthorized()(*bool)
    GetDisplayName()(*string)
    GetEnabled()(*bool)
    GetLastConnectionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetUrl()(*string)
    SetAppAuthorized(value *bool)()
    SetDisplayName(value *string)()
    SetEnabled(value *bool)()
    SetLastConnectionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetUrl(value *string)()
}

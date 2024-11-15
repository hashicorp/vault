package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MailboxSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMailboxSettings instantiates a new MailboxSettings and sets the default values.
func NewMailboxSettings()(*MailboxSettings) {
    m := &MailboxSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMailboxSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailboxSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailboxSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MailboxSettings) GetAdditionalData()(map[string]any) {
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
// GetArchiveFolder gets the archiveFolder property value. Folder ID of an archive folder for the user.
// returns a *string when successful
func (m *MailboxSettings) GetArchiveFolder()(*string) {
    val, err := m.GetBackingStore().Get("archiveFolder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAutomaticRepliesSetting gets the automaticRepliesSetting property value. Configuration settings to automatically notify the sender of an incoming email with a message from the signed-in user.
// returns a AutomaticRepliesSettingable when successful
func (m *MailboxSettings) GetAutomaticRepliesSetting()(AutomaticRepliesSettingable) {
    val, err := m.GetBackingStore().Get("automaticRepliesSetting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AutomaticRepliesSettingable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MailboxSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDateFormat gets the dateFormat property value. The date format for the user's mailbox.
// returns a *string when successful
func (m *MailboxSettings) GetDateFormat()(*string) {
    val, err := m.GetBackingStore().Get("dateFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDelegateMeetingMessageDeliveryOptions gets the delegateMeetingMessageDeliveryOptions property value. If the user has a calendar delegate, this specifies whether the delegate, mailbox owner, or both receive meeting messages and meeting responses. Possible values are: sendToDelegateAndInformationToPrincipal, sendToDelegateAndPrincipal, sendToDelegateOnly.
// returns a *DelegateMeetingMessageDeliveryOptions when successful
func (m *MailboxSettings) GetDelegateMeetingMessageDeliveryOptions()(*DelegateMeetingMessageDeliveryOptions) {
    val, err := m.GetBackingStore().Get("delegateMeetingMessageDeliveryOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DelegateMeetingMessageDeliveryOptions)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailboxSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["archiveFolder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetArchiveFolder(val)
        }
        return nil
    }
    res["automaticRepliesSetting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAutomaticRepliesSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomaticRepliesSetting(val.(AutomaticRepliesSettingable))
        }
        return nil
    }
    res["dateFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDateFormat(val)
        }
        return nil
    }
    res["delegateMeetingMessageDeliveryOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDelegateMeetingMessageDeliveryOptions)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDelegateMeetingMessageDeliveryOptions(val.(*DelegateMeetingMessageDeliveryOptions))
        }
        return nil
    }
    res["language"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocaleInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLanguage(val.(LocaleInfoable))
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
    res["timeFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeFormat(val)
        }
        return nil
    }
    res["timeZone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeZone(val)
        }
        return nil
    }
    res["userPurpose"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserPurpose)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPurpose(val.(*UserPurpose))
        }
        return nil
    }
    res["workingHours"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkingHoursFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkingHours(val.(WorkingHoursable))
        }
        return nil
    }
    return res
}
// GetLanguage gets the language property value. The locale information for the user, including the preferred language and country/region.
// returns a LocaleInfoable when successful
func (m *MailboxSettings) GetLanguage()(LocaleInfoable) {
    val, err := m.GetBackingStore().Get("language")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LocaleInfoable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MailboxSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTimeFormat gets the timeFormat property value. The time format for the user's mailbox.
// returns a *string when successful
func (m *MailboxSettings) GetTimeFormat()(*string) {
    val, err := m.GetBackingStore().Get("timeFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTimeZone gets the timeZone property value. The default time zone for the user's mailbox.
// returns a *string when successful
func (m *MailboxSettings) GetTimeZone()(*string) {
    val, err := m.GetBackingStore().Get("timeZone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPurpose gets the userPurpose property value. The purpose of the mailbox. Differentiates a mailbox for a single user from a shared mailbox and equipment mailbox in Exchange Online. Possible values are: user, linked, shared, room, equipment, others, unknownFutureValue. Read-only.
// returns a *UserPurpose when successful
func (m *MailboxSettings) GetUserPurpose()(*UserPurpose) {
    val, err := m.GetBackingStore().Get("userPurpose")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserPurpose)
    }
    return nil
}
// GetWorkingHours gets the workingHours property value. The days of the week and hours in a specific time zone that the user works.
// returns a WorkingHoursable when successful
func (m *MailboxSettings) GetWorkingHours()(WorkingHoursable) {
    val, err := m.GetBackingStore().Get("workingHours")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkingHoursable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailboxSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("archiveFolder", m.GetArchiveFolder())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("automaticRepliesSetting", m.GetAutomaticRepliesSetting())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("dateFormat", m.GetDateFormat())
        if err != nil {
            return err
        }
    }
    if m.GetDelegateMeetingMessageDeliveryOptions() != nil {
        cast := (*m.GetDelegateMeetingMessageDeliveryOptions()).String()
        err := writer.WriteStringValue("delegateMeetingMessageDeliveryOptions", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("language", m.GetLanguage())
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
        err := writer.WriteStringValue("timeFormat", m.GetTimeFormat())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("timeZone", m.GetTimeZone())
        if err != nil {
            return err
        }
    }
    if m.GetUserPurpose() != nil {
        cast := (*m.GetUserPurpose()).String()
        err := writer.WriteStringValue("userPurpose", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("workingHours", m.GetWorkingHours())
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
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MailboxSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetArchiveFolder sets the archiveFolder property value. Folder ID of an archive folder for the user.
func (m *MailboxSettings) SetArchiveFolder(value *string)() {
    err := m.GetBackingStore().Set("archiveFolder", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomaticRepliesSetting sets the automaticRepliesSetting property value. Configuration settings to automatically notify the sender of an incoming email with a message from the signed-in user.
func (m *MailboxSettings) SetAutomaticRepliesSetting(value AutomaticRepliesSettingable)() {
    err := m.GetBackingStore().Set("automaticRepliesSetting", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MailboxSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDateFormat sets the dateFormat property value. The date format for the user's mailbox.
func (m *MailboxSettings) SetDateFormat(value *string)() {
    err := m.GetBackingStore().Set("dateFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetDelegateMeetingMessageDeliveryOptions sets the delegateMeetingMessageDeliveryOptions property value. If the user has a calendar delegate, this specifies whether the delegate, mailbox owner, or both receive meeting messages and meeting responses. Possible values are: sendToDelegateAndInformationToPrincipal, sendToDelegateAndPrincipal, sendToDelegateOnly.
func (m *MailboxSettings) SetDelegateMeetingMessageDeliveryOptions(value *DelegateMeetingMessageDeliveryOptions)() {
    err := m.GetBackingStore().Set("delegateMeetingMessageDeliveryOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguage sets the language property value. The locale information for the user, including the preferred language and country/region.
func (m *MailboxSettings) SetLanguage(value LocaleInfoable)() {
    err := m.GetBackingStore().Set("language", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MailboxSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeFormat sets the timeFormat property value. The time format for the user's mailbox.
func (m *MailboxSettings) SetTimeFormat(value *string)() {
    err := m.GetBackingStore().Set("timeFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeZone sets the timeZone property value. The default time zone for the user's mailbox.
func (m *MailboxSettings) SetTimeZone(value *string)() {
    err := m.GetBackingStore().Set("timeZone", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPurpose sets the userPurpose property value. The purpose of the mailbox. Differentiates a mailbox for a single user from a shared mailbox and equipment mailbox in Exchange Online. Possible values are: user, linked, shared, room, equipment, others, unknownFutureValue. Read-only.
func (m *MailboxSettings) SetUserPurpose(value *UserPurpose)() {
    err := m.GetBackingStore().Set("userPurpose", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkingHours sets the workingHours property value. The days of the week and hours in a specific time zone that the user works.
func (m *MailboxSettings) SetWorkingHours(value WorkingHoursable)() {
    err := m.GetBackingStore().Set("workingHours", value)
    if err != nil {
        panic(err)
    }
}
type MailboxSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetArchiveFolder()(*string)
    GetAutomaticRepliesSetting()(AutomaticRepliesSettingable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDateFormat()(*string)
    GetDelegateMeetingMessageDeliveryOptions()(*DelegateMeetingMessageDeliveryOptions)
    GetLanguage()(LocaleInfoable)
    GetOdataType()(*string)
    GetTimeFormat()(*string)
    GetTimeZone()(*string)
    GetUserPurpose()(*UserPurpose)
    GetWorkingHours()(WorkingHoursable)
    SetArchiveFolder(value *string)()
    SetAutomaticRepliesSetting(value AutomaticRepliesSettingable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDateFormat(value *string)()
    SetDelegateMeetingMessageDeliveryOptions(value *DelegateMeetingMessageDeliveryOptions)()
    SetLanguage(value LocaleInfoable)()
    SetOdataType(value *string)()
    SetTimeFormat(value *string)()
    SetTimeZone(value *string)()
    SetUserPurpose(value *UserPurpose)()
    SetWorkingHours(value WorkingHoursable)()
}

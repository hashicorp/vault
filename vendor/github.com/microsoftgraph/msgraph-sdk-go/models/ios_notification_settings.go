package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// IosNotificationSettings an item describing notification setting.
type IosNotificationSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewIosNotificationSettings instantiates a new IosNotificationSettings and sets the default values.
func NewIosNotificationSettings()(*IosNotificationSettings) {
    m := &IosNotificationSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateIosNotificationSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosNotificationSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosNotificationSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *IosNotificationSettings) GetAdditionalData()(map[string]any) {
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
// GetAlertType gets the alertType property value. Notification Settings Alert Type.
// returns a *IosNotificationAlertType when successful
func (m *IosNotificationSettings) GetAlertType()(*IosNotificationAlertType) {
    val, err := m.GetBackingStore().Get("alertType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IosNotificationAlertType)
    }
    return nil
}
// GetAppName gets the appName property value. Application name to be associated with the bundleID.
// returns a *string when successful
func (m *IosNotificationSettings) GetAppName()(*string) {
    val, err := m.GetBackingStore().Get("appName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *IosNotificationSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBadgesEnabled gets the badgesEnabled property value. Indicates whether badges are allowed for this app.
// returns a *bool when successful
func (m *IosNotificationSettings) GetBadgesEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("badgesEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBundleID gets the bundleID property value. Bundle id of app to which to apply these notification settings.
// returns a *string when successful
func (m *IosNotificationSettings) GetBundleID()(*string) {
    val, err := m.GetBackingStore().Get("bundleID")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnabled gets the enabled property value. Indicates whether notifications are allowed for this app.
// returns a *bool when successful
func (m *IosNotificationSettings) GetEnabled()(*bool) {
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
func (m *IosNotificationSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["alertType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseIosNotificationAlertType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlertType(val.(*IosNotificationAlertType))
        }
        return nil
    }
    res["appName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppName(val)
        }
        return nil
    }
    res["badgesEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBadgesEnabled(val)
        }
        return nil
    }
    res["bundleID"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBundleID(val)
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
    res["publisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisher(val)
        }
        return nil
    }
    res["showInNotificationCenter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowInNotificationCenter(val)
        }
        return nil
    }
    res["showOnLockScreen"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowOnLockScreen(val)
        }
        return nil
    }
    res["soundsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSoundsEnabled(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *IosNotificationSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublisher gets the publisher property value. Publisher to be associated with the bundleID.
// returns a *string when successful
func (m *IosNotificationSettings) GetPublisher()(*string) {
    val, err := m.GetBackingStore().Get("publisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetShowInNotificationCenter gets the showInNotificationCenter property value. Indicates whether notifications can be shown in notification center.
// returns a *bool when successful
func (m *IosNotificationSettings) GetShowInNotificationCenter()(*bool) {
    val, err := m.GetBackingStore().Get("showInNotificationCenter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowOnLockScreen gets the showOnLockScreen property value. Indicates whether notifications can be shown on the lock screen.
// returns a *bool when successful
func (m *IosNotificationSettings) GetShowOnLockScreen()(*bool) {
    val, err := m.GetBackingStore().Get("showOnLockScreen")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSoundsEnabled gets the soundsEnabled property value. Indicates whether sounds are allowed for this app.
// returns a *bool when successful
func (m *IosNotificationSettings) GetSoundsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("soundsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosNotificationSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAlertType() != nil {
        cast := (*m.GetAlertType()).String()
        err := writer.WriteStringValue("alertType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("appName", m.GetAppName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("badgesEnabled", m.GetBadgesEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bundleID", m.GetBundleID())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enabled", m.GetEnabled())
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
        err := writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("showInNotificationCenter", m.GetShowInNotificationCenter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("showOnLockScreen", m.GetShowOnLockScreen())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("soundsEnabled", m.GetSoundsEnabled())
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
func (m *IosNotificationSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAlertType sets the alertType property value. Notification Settings Alert Type.
func (m *IosNotificationSettings) SetAlertType(value *IosNotificationAlertType)() {
    err := m.GetBackingStore().Set("alertType", value)
    if err != nil {
        panic(err)
    }
}
// SetAppName sets the appName property value. Application name to be associated with the bundleID.
func (m *IosNotificationSettings) SetAppName(value *string)() {
    err := m.GetBackingStore().Set("appName", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *IosNotificationSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBadgesEnabled sets the badgesEnabled property value. Indicates whether badges are allowed for this app.
func (m *IosNotificationSettings) SetBadgesEnabled(value *bool)() {
    err := m.GetBackingStore().Set("badgesEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetBundleID sets the bundleID property value. Bundle id of app to which to apply these notification settings.
func (m *IosNotificationSettings) SetBundleID(value *string)() {
    err := m.GetBackingStore().Set("bundleID", value)
    if err != nil {
        panic(err)
    }
}
// SetEnabled sets the enabled property value. Indicates whether notifications are allowed for this app.
func (m *IosNotificationSettings) SetEnabled(value *bool)() {
    err := m.GetBackingStore().Set("enabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IosNotificationSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisher sets the publisher property value. Publisher to be associated with the bundleID.
func (m *IosNotificationSettings) SetPublisher(value *string)() {
    err := m.GetBackingStore().Set("publisher", value)
    if err != nil {
        panic(err)
    }
}
// SetShowInNotificationCenter sets the showInNotificationCenter property value. Indicates whether notifications can be shown in notification center.
func (m *IosNotificationSettings) SetShowInNotificationCenter(value *bool)() {
    err := m.GetBackingStore().Set("showInNotificationCenter", value)
    if err != nil {
        panic(err)
    }
}
// SetShowOnLockScreen sets the showOnLockScreen property value. Indicates whether notifications can be shown on the lock screen.
func (m *IosNotificationSettings) SetShowOnLockScreen(value *bool)() {
    err := m.GetBackingStore().Set("showOnLockScreen", value)
    if err != nil {
        panic(err)
    }
}
// SetSoundsEnabled sets the soundsEnabled property value. Indicates whether sounds are allowed for this app.
func (m *IosNotificationSettings) SetSoundsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("soundsEnabled", value)
    if err != nil {
        panic(err)
    }
}
type IosNotificationSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlertType()(*IosNotificationAlertType)
    GetAppName()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBadgesEnabled()(*bool)
    GetBundleID()(*string)
    GetEnabled()(*bool)
    GetOdataType()(*string)
    GetPublisher()(*string)
    GetShowInNotificationCenter()(*bool)
    GetShowOnLockScreen()(*bool)
    GetSoundsEnabled()(*bool)
    SetAlertType(value *IosNotificationAlertType)()
    SetAppName(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBadgesEnabled(value *bool)()
    SetBundleID(value *string)()
    SetEnabled(value *bool)()
    SetOdataType(value *string)()
    SetPublisher(value *string)()
    SetShowInNotificationCenter(value *bool)()
    SetShowOnLockScreen(value *bool)()
    SetSoundsEnabled(value *bool)()
}

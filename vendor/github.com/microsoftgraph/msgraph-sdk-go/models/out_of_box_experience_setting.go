package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// OutOfBoxExperienceSetting the Windows Autopilot Deployment Profile settings used by the device for the out-of-box experience. Supports: $select, $top, $skip. $Search, $orderBy and $filter are not supported.
type OutOfBoxExperienceSetting struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewOutOfBoxExperienceSetting instantiates a new OutOfBoxExperienceSetting and sets the default values.
func NewOutOfBoxExperienceSetting()(*OutOfBoxExperienceSetting) {
    m := &OutOfBoxExperienceSetting{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateOutOfBoxExperienceSettingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOutOfBoxExperienceSettingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOutOfBoxExperienceSetting(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *OutOfBoxExperienceSetting) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *OutOfBoxExperienceSetting) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDeviceUsageType gets the deviceUsageType property value. The deviceUsageType property
// returns a *WindowsDeviceUsageType when successful
func (m *OutOfBoxExperienceSetting) GetDeviceUsageType()(*WindowsDeviceUsageType) {
    val, err := m.GetBackingStore().Get("deviceUsageType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsDeviceUsageType)
    }
    return nil
}
// GetEscapeLinkHidden gets the escapeLinkHidden property value. When TRUE, the link that allows user to start over with a different account on company sign-in is hidden. When false, the link that allows user to start over with a different account on company sign-in is available. Default value is FALSE.
// returns a *bool when successful
func (m *OutOfBoxExperienceSetting) GetEscapeLinkHidden()(*bool) {
    val, err := m.GetBackingStore().Get("escapeLinkHidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEulaHidden gets the eulaHidden property value. When TRUE, EULA is hidden to the end user during OOBE. When FALSE, EULA is shown to the end user during OOBE. Default value is FALSE.
// returns a *bool when successful
func (m *OutOfBoxExperienceSetting) GetEulaHidden()(*bool) {
    val, err := m.GetBackingStore().Get("eulaHidden")
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
func (m *OutOfBoxExperienceSetting) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["deviceUsageType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsDeviceUsageType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceUsageType(val.(*WindowsDeviceUsageType))
        }
        return nil
    }
    res["escapeLinkHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEscapeLinkHidden(val)
        }
        return nil
    }
    res["eulaHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEulaHidden(val)
        }
        return nil
    }
    res["keyboardSelectionPageSkipped"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardSelectionPageSkipped(val)
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
    res["privacySettingsHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacySettingsHidden(val)
        }
        return nil
    }
    res["userType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsUserType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserType(val.(*WindowsUserType))
        }
        return nil
    }
    return res
}
// GetKeyboardSelectionPageSkipped gets the keyboardSelectionPageSkipped property value. When TRUE, the keyboard selection page is hidden to the end user during OOBE if Language and Region are set. When FALSE, the keyboard selection page is skipped during OOBE.
// returns a *bool when successful
func (m *OutOfBoxExperienceSetting) GetKeyboardSelectionPageSkipped()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardSelectionPageSkipped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *OutOfBoxExperienceSetting) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrivacySettingsHidden gets the privacySettingsHidden property value. When TRUE, privacy settings is hidden to the end user during OOBE. When FALSE, privacy settings is shown to the end user during OOBE. Default value is FALSE.
// returns a *bool when successful
func (m *OutOfBoxExperienceSetting) GetPrivacySettingsHidden()(*bool) {
    val, err := m.GetBackingStore().Get("privacySettingsHidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUserType gets the userType property value. The userType property
// returns a *WindowsUserType when successful
func (m *OutOfBoxExperienceSetting) GetUserType()(*WindowsUserType) {
    val, err := m.GetBackingStore().Get("userType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsUserType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OutOfBoxExperienceSetting) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDeviceUsageType() != nil {
        cast := (*m.GetDeviceUsageType()).String()
        err := writer.WriteStringValue("deviceUsageType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("escapeLinkHidden", m.GetEscapeLinkHidden())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("eulaHidden", m.GetEulaHidden())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("keyboardSelectionPageSkipped", m.GetKeyboardSelectionPageSkipped())
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
        err := writer.WriteBoolValue("privacySettingsHidden", m.GetPrivacySettingsHidden())
        if err != nil {
            return err
        }
    }
    if m.GetUserType() != nil {
        cast := (*m.GetUserType()).String()
        err := writer.WriteStringValue("userType", &cast)
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
func (m *OutOfBoxExperienceSetting) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *OutOfBoxExperienceSetting) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDeviceUsageType sets the deviceUsageType property value. The deviceUsageType property
func (m *OutOfBoxExperienceSetting) SetDeviceUsageType(value *WindowsDeviceUsageType)() {
    err := m.GetBackingStore().Set("deviceUsageType", value)
    if err != nil {
        panic(err)
    }
}
// SetEscapeLinkHidden sets the escapeLinkHidden property value. When TRUE, the link that allows user to start over with a different account on company sign-in is hidden. When false, the link that allows user to start over with a different account on company sign-in is available. Default value is FALSE.
func (m *OutOfBoxExperienceSetting) SetEscapeLinkHidden(value *bool)() {
    err := m.GetBackingStore().Set("escapeLinkHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetEulaHidden sets the eulaHidden property value. When TRUE, EULA is hidden to the end user during OOBE. When FALSE, EULA is shown to the end user during OOBE. Default value is FALSE.
func (m *OutOfBoxExperienceSetting) SetEulaHidden(value *bool)() {
    err := m.GetBackingStore().Set("eulaHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardSelectionPageSkipped sets the keyboardSelectionPageSkipped property value. When TRUE, the keyboard selection page is hidden to the end user during OOBE if Language and Region are set. When FALSE, the keyboard selection page is skipped during OOBE.
func (m *OutOfBoxExperienceSetting) SetKeyboardSelectionPageSkipped(value *bool)() {
    err := m.GetBackingStore().Set("keyboardSelectionPageSkipped", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OutOfBoxExperienceSetting) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacySettingsHidden sets the privacySettingsHidden property value. When TRUE, privacy settings is hidden to the end user during OOBE. When FALSE, privacy settings is shown to the end user during OOBE. Default value is FALSE.
func (m *OutOfBoxExperienceSetting) SetPrivacySettingsHidden(value *bool)() {
    err := m.GetBackingStore().Set("privacySettingsHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetUserType sets the userType property value. The userType property
func (m *OutOfBoxExperienceSetting) SetUserType(value *WindowsUserType)() {
    err := m.GetBackingStore().Set("userType", value)
    if err != nil {
        panic(err)
    }
}
type OutOfBoxExperienceSettingable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDeviceUsageType()(*WindowsDeviceUsageType)
    GetEscapeLinkHidden()(*bool)
    GetEulaHidden()(*bool)
    GetKeyboardSelectionPageSkipped()(*bool)
    GetOdataType()(*string)
    GetPrivacySettingsHidden()(*bool)
    GetUserType()(*WindowsUserType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDeviceUsageType(value *WindowsDeviceUsageType)()
    SetEscapeLinkHidden(value *bool)()
    SetEulaHidden(value *bool)()
    SetKeyboardSelectionPageSkipped(value *bool)()
    SetOdataType(value *string)()
    SetPrivacySettingsHidden(value *bool)()
    SetUserType(value *WindowsUserType)()
}

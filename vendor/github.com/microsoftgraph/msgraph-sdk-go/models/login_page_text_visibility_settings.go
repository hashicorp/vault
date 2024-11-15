package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type LoginPageTextVisibilitySettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewLoginPageTextVisibilitySettings instantiates a new LoginPageTextVisibilitySettings and sets the default values.
func NewLoginPageTextVisibilitySettings()(*LoginPageTextVisibilitySettings) {
    m := &LoginPageTextVisibilitySettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateLoginPageTextVisibilitySettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLoginPageTextVisibilitySettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLoginPageTextVisibilitySettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *LoginPageTextVisibilitySettings) GetAdditionalData()(map[string]any) {
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
func (m *LoginPageTextVisibilitySettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LoginPageTextVisibilitySettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hideAccountResetCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideAccountResetCredentials(val)
        }
        return nil
    }
    res["hideCannotAccessYourAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideCannotAccessYourAccount(val)
        }
        return nil
    }
    res["hideForgotMyPassword"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideForgotMyPassword(val)
        }
        return nil
    }
    res["hidePrivacyAndCookies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHidePrivacyAndCookies(val)
        }
        return nil
    }
    res["hideResetItNow"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideResetItNow(val)
        }
        return nil
    }
    res["hideTermsOfUse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideTermsOfUse(val)
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
    return res
}
// GetHideAccountResetCredentials gets the hideAccountResetCredentials property value. Option to hide the self-service password reset (SSPR) hyperlinks such as 'Can't access your account?', 'Forgot my password' and 'Reset it now' on the sign-in form.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHideAccountResetCredentials()(*bool) {
    val, err := m.GetBackingStore().Get("hideAccountResetCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHideCannotAccessYourAccount gets the hideCannotAccessYourAccount property value. Option to hide the self-service password reset (SSPR) 'Can't access your account?' hyperlink on the sign-in form.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHideCannotAccessYourAccount()(*bool) {
    val, err := m.GetBackingStore().Get("hideCannotAccessYourAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHideForgotMyPassword gets the hideForgotMyPassword property value. Option to hide the self-service password reset (SSPR) 'Forgot my password' hyperlink on the sign-in form.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHideForgotMyPassword()(*bool) {
    val, err := m.GetBackingStore().Get("hideForgotMyPassword")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHidePrivacyAndCookies gets the hidePrivacyAndCookies property value. Option to hide the 'Privacy & Cookies' hyperlink in the footer.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHidePrivacyAndCookies()(*bool) {
    val, err := m.GetBackingStore().Get("hidePrivacyAndCookies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHideResetItNow gets the hideResetItNow property value. Option to hide the self-service password reset (SSPR) 'reset it now' hyperlink on the sign-in form.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHideResetItNow()(*bool) {
    val, err := m.GetBackingStore().Get("hideResetItNow")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHideTermsOfUse gets the hideTermsOfUse property value. Option to hide the 'Terms of Use' hyperlink in the footer.
// returns a *bool when successful
func (m *LoginPageTextVisibilitySettings) GetHideTermsOfUse()(*bool) {
    val, err := m.GetBackingStore().Get("hideTermsOfUse")
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
func (m *LoginPageTextVisibilitySettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LoginPageTextVisibilitySettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("hideAccountResetCredentials", m.GetHideAccountResetCredentials())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hideCannotAccessYourAccount", m.GetHideCannotAccessYourAccount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hideForgotMyPassword", m.GetHideForgotMyPassword())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hidePrivacyAndCookies", m.GetHidePrivacyAndCookies())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hideResetItNow", m.GetHideResetItNow())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hideTermsOfUse", m.GetHideTermsOfUse())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *LoginPageTextVisibilitySettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *LoginPageTextVisibilitySettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetHideAccountResetCredentials sets the hideAccountResetCredentials property value. Option to hide the self-service password reset (SSPR) hyperlinks such as 'Can't access your account?', 'Forgot my password' and 'Reset it now' on the sign-in form.
func (m *LoginPageTextVisibilitySettings) SetHideAccountResetCredentials(value *bool)() {
    err := m.GetBackingStore().Set("hideAccountResetCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetHideCannotAccessYourAccount sets the hideCannotAccessYourAccount property value. Option to hide the self-service password reset (SSPR) 'Can't access your account?' hyperlink on the sign-in form.
func (m *LoginPageTextVisibilitySettings) SetHideCannotAccessYourAccount(value *bool)() {
    err := m.GetBackingStore().Set("hideCannotAccessYourAccount", value)
    if err != nil {
        panic(err)
    }
}
// SetHideForgotMyPassword sets the hideForgotMyPassword property value. Option to hide the self-service password reset (SSPR) 'Forgot my password' hyperlink on the sign-in form.
func (m *LoginPageTextVisibilitySettings) SetHideForgotMyPassword(value *bool)() {
    err := m.GetBackingStore().Set("hideForgotMyPassword", value)
    if err != nil {
        panic(err)
    }
}
// SetHidePrivacyAndCookies sets the hidePrivacyAndCookies property value. Option to hide the 'Privacy & Cookies' hyperlink in the footer.
func (m *LoginPageTextVisibilitySettings) SetHidePrivacyAndCookies(value *bool)() {
    err := m.GetBackingStore().Set("hidePrivacyAndCookies", value)
    if err != nil {
        panic(err)
    }
}
// SetHideResetItNow sets the hideResetItNow property value. Option to hide the self-service password reset (SSPR) 'reset it now' hyperlink on the sign-in form.
func (m *LoginPageTextVisibilitySettings) SetHideResetItNow(value *bool)() {
    err := m.GetBackingStore().Set("hideResetItNow", value)
    if err != nil {
        panic(err)
    }
}
// SetHideTermsOfUse sets the hideTermsOfUse property value. Option to hide the 'Terms of Use' hyperlink in the footer.
func (m *LoginPageTextVisibilitySettings) SetHideTermsOfUse(value *bool)() {
    err := m.GetBackingStore().Set("hideTermsOfUse", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *LoginPageTextVisibilitySettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type LoginPageTextVisibilitySettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetHideAccountResetCredentials()(*bool)
    GetHideCannotAccessYourAccount()(*bool)
    GetHideForgotMyPassword()(*bool)
    GetHidePrivacyAndCookies()(*bool)
    GetHideResetItNow()(*bool)
    GetHideTermsOfUse()(*bool)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetHideAccountResetCredentials(value *bool)()
    SetHideCannotAccessYourAccount(value *bool)()
    SetHideForgotMyPassword(value *bool)()
    SetHidePrivacyAndCookies(value *bool)()
    SetHideResetItNow(value *bool)()
    SetHideTermsOfUse(value *bool)()
    SetOdataType(value *string)()
}

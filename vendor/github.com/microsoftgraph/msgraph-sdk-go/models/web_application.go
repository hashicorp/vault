package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type WebApplication struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewWebApplication instantiates a new WebApplication and sets the default values.
func NewWebApplication()(*WebApplication) {
    m := &WebApplication{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateWebApplicationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWebApplicationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWebApplication(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *WebApplication) GetAdditionalData()(map[string]any) {
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
func (m *WebApplication) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WebApplication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["homePageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHomePageUrl(val)
        }
        return nil
    }
    res["implicitGrantSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateImplicitGrantSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImplicitGrantSettings(val.(ImplicitGrantSettingsable))
        }
        return nil
    }
    res["logoutUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogoutUrl(val)
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
    res["redirectUris"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRedirectUris(res)
        }
        return nil
    }
    res["redirectUriSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRedirectUriSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RedirectUriSettingsable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RedirectUriSettingsable)
                }
            }
            m.SetRedirectUriSettings(res)
        }
        return nil
    }
    return res
}
// GetHomePageUrl gets the homePageUrl property value. Home page or landing page of the application.
// returns a *string when successful
func (m *WebApplication) GetHomePageUrl()(*string) {
    val, err := m.GetBackingStore().Get("homePageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetImplicitGrantSettings gets the implicitGrantSettings property value. Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.
// returns a ImplicitGrantSettingsable when successful
func (m *WebApplication) GetImplicitGrantSettings()(ImplicitGrantSettingsable) {
    val, err := m.GetBackingStore().Get("implicitGrantSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ImplicitGrantSettingsable)
    }
    return nil
}
// GetLogoutUrl gets the logoutUrl property value. Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.
// returns a *string when successful
func (m *WebApplication) GetLogoutUrl()(*string) {
    val, err := m.GetBackingStore().Get("logoutUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *WebApplication) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRedirectUris gets the redirectUris property value. Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.
// returns a []string when successful
func (m *WebApplication) GetRedirectUris()([]string) {
    val, err := m.GetBackingStore().Get("redirectUris")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetRedirectUriSettings gets the redirectUriSettings property value. The redirectUriSettings property
// returns a []RedirectUriSettingsable when successful
func (m *WebApplication) GetRedirectUriSettings()([]RedirectUriSettingsable) {
    val, err := m.GetBackingStore().Get("redirectUriSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RedirectUriSettingsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WebApplication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("homePageUrl", m.GetHomePageUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("implicitGrantSettings", m.GetImplicitGrantSettings())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("logoutUrl", m.GetLogoutUrl())
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
    if m.GetRedirectUris() != nil {
        err := writer.WriteCollectionOfStringValues("redirectUris", m.GetRedirectUris())
        if err != nil {
            return err
        }
    }
    if m.GetRedirectUriSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRedirectUriSettings()))
        for i, v := range m.GetRedirectUriSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("redirectUriSettings", cast)
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
func (m *WebApplication) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *WebApplication) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetHomePageUrl sets the homePageUrl property value. Home page or landing page of the application.
func (m *WebApplication) SetHomePageUrl(value *string)() {
    err := m.GetBackingStore().Set("homePageUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetImplicitGrantSettings sets the implicitGrantSettings property value. Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.
func (m *WebApplication) SetImplicitGrantSettings(value ImplicitGrantSettingsable)() {
    err := m.GetBackingStore().Set("implicitGrantSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetLogoutUrl sets the logoutUrl property value. Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.
func (m *WebApplication) SetLogoutUrl(value *string)() {
    err := m.GetBackingStore().Set("logoutUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WebApplication) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRedirectUris sets the redirectUris property value. Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.
func (m *WebApplication) SetRedirectUris(value []string)() {
    err := m.GetBackingStore().Set("redirectUris", value)
    if err != nil {
        panic(err)
    }
}
// SetRedirectUriSettings sets the redirectUriSettings property value. The redirectUriSettings property
func (m *WebApplication) SetRedirectUriSettings(value []RedirectUriSettingsable)() {
    err := m.GetBackingStore().Set("redirectUriSettings", value)
    if err != nil {
        panic(err)
    }
}
type WebApplicationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetHomePageUrl()(*string)
    GetImplicitGrantSettings()(ImplicitGrantSettingsable)
    GetLogoutUrl()(*string)
    GetOdataType()(*string)
    GetRedirectUris()([]string)
    GetRedirectUriSettings()([]RedirectUriSettingsable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetHomePageUrl(value *string)()
    SetImplicitGrantSettings(value ImplicitGrantSettingsable)()
    SetLogoutUrl(value *string)()
    SetOdataType(value *string)()
    SetRedirectUris(value []string)()
    SetRedirectUriSettings(value []RedirectUriSettingsable)()
}

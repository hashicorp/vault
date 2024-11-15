package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type InformationalUrl struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewInformationalUrl instantiates a new InformationalUrl and sets the default values.
func NewInformationalUrl()(*InformationalUrl) {
    m := &InformationalUrl{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateInformationalUrlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInformationalUrlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInformationalUrl(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *InformationalUrl) GetAdditionalData()(map[string]any) {
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
func (m *InformationalUrl) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *InformationalUrl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["logoUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogoUrl(val)
        }
        return nil
    }
    res["marketingUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMarketingUrl(val)
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
    res["privacyStatementUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacyStatementUrl(val)
        }
        return nil
    }
    res["supportUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSupportUrl(val)
        }
        return nil
    }
    res["termsOfServiceUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTermsOfServiceUrl(val)
        }
        return nil
    }
    return res
}
// GetLogoUrl gets the logoUrl property value. CDN URL to the application's logo, Read-only.
// returns a *string when successful
func (m *InformationalUrl) GetLogoUrl()(*string) {
    val, err := m.GetBackingStore().Get("logoUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMarketingUrl gets the marketingUrl property value. Link to the application's marketing page. For example, https://www.contoso.com/app/marketing
// returns a *string when successful
func (m *InformationalUrl) GetMarketingUrl()(*string) {
    val, err := m.GetBackingStore().Get("marketingUrl")
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
func (m *InformationalUrl) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrivacyStatementUrl gets the privacyStatementUrl property value. Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy
// returns a *string when successful
func (m *InformationalUrl) GetPrivacyStatementUrl()(*string) {
    val, err := m.GetBackingStore().Get("privacyStatementUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSupportUrl gets the supportUrl property value. Link to the application's support page. For example, https://www.contoso.com/app/support
// returns a *string when successful
func (m *InformationalUrl) GetSupportUrl()(*string) {
    val, err := m.GetBackingStore().Get("supportUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTermsOfServiceUrl gets the termsOfServiceUrl property value. Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice
// returns a *string when successful
func (m *InformationalUrl) GetTermsOfServiceUrl()(*string) {
    val, err := m.GetBackingStore().Get("termsOfServiceUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *InformationalUrl) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("logoUrl", m.GetLogoUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("marketingUrl", m.GetMarketingUrl())
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
        err := writer.WriteStringValue("privacyStatementUrl", m.GetPrivacyStatementUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("supportUrl", m.GetSupportUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("termsOfServiceUrl", m.GetTermsOfServiceUrl())
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
func (m *InformationalUrl) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *InformationalUrl) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetLogoUrl sets the logoUrl property value. CDN URL to the application's logo, Read-only.
func (m *InformationalUrl) SetLogoUrl(value *string)() {
    err := m.GetBackingStore().Set("logoUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetMarketingUrl sets the marketingUrl property value. Link to the application's marketing page. For example, https://www.contoso.com/app/marketing
func (m *InformationalUrl) SetMarketingUrl(value *string)() {
    err := m.GetBackingStore().Set("marketingUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *InformationalUrl) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacyStatementUrl sets the privacyStatementUrl property value. Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy
func (m *InformationalUrl) SetPrivacyStatementUrl(value *string)() {
    err := m.GetBackingStore().Set("privacyStatementUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportUrl sets the supportUrl property value. Link to the application's support page. For example, https://www.contoso.com/app/support
func (m *InformationalUrl) SetSupportUrl(value *string)() {
    err := m.GetBackingStore().Set("supportUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetTermsOfServiceUrl sets the termsOfServiceUrl property value. Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice
func (m *InformationalUrl) SetTermsOfServiceUrl(value *string)() {
    err := m.GetBackingStore().Set("termsOfServiceUrl", value)
    if err != nil {
        panic(err)
    }
}
type InformationalUrlable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetLogoUrl()(*string)
    GetMarketingUrl()(*string)
    GetOdataType()(*string)
    GetPrivacyStatementUrl()(*string)
    GetSupportUrl()(*string)
    GetTermsOfServiceUrl()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetLogoUrl(value *string)()
    SetMarketingUrl(value *string)()
    SetOdataType(value *string)()
    SetPrivacyStatementUrl(value *string)()
    SetSupportUrl(value *string)()
    SetTermsOfServiceUrl(value *string)()
}

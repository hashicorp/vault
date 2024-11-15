package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MicrosoftAuthenticatorFeatureSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMicrosoftAuthenticatorFeatureSettings instantiates a new MicrosoftAuthenticatorFeatureSettings and sets the default values.
func NewMicrosoftAuthenticatorFeatureSettings()(*MicrosoftAuthenticatorFeatureSettings) {
    m := &MicrosoftAuthenticatorFeatureSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMicrosoftAuthenticatorFeatureSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMicrosoftAuthenticatorFeatureSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftAuthenticatorFeatureSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MicrosoftAuthenticatorFeatureSettings) GetAdditionalData()(map[string]any) {
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
func (m *MicrosoftAuthenticatorFeatureSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDisplayAppInformationRequiredState gets the displayAppInformationRequiredState property value. Determines whether the user's Authenticator app shows them the client app they're signing into.
// returns a AuthenticationMethodFeatureConfigurationable when successful
func (m *MicrosoftAuthenticatorFeatureSettings) GetDisplayAppInformationRequiredState()(AuthenticationMethodFeatureConfigurationable) {
    val, err := m.GetBackingStore().Get("displayAppInformationRequiredState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationMethodFeatureConfigurationable)
    }
    return nil
}
// GetDisplayLocationInformationRequiredState gets the displayLocationInformationRequiredState property value. Determines whether the user's Authenticator app shows them the geographic location of where the authentication request originated from.
// returns a AuthenticationMethodFeatureConfigurationable when successful
func (m *MicrosoftAuthenticatorFeatureSettings) GetDisplayLocationInformationRequiredState()(AuthenticationMethodFeatureConfigurationable) {
    val, err := m.GetBackingStore().Get("displayLocationInformationRequiredState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationMethodFeatureConfigurationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MicrosoftAuthenticatorFeatureSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayAppInformationRequiredState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationMethodFeatureConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayAppInformationRequiredState(val.(AuthenticationMethodFeatureConfigurationable))
        }
        return nil
    }
    res["displayLocationInformationRequiredState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationMethodFeatureConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayLocationInformationRequiredState(val.(AuthenticationMethodFeatureConfigurationable))
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
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MicrosoftAuthenticatorFeatureSettings) GetOdataType()(*string) {
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
func (m *MicrosoftAuthenticatorFeatureSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("displayAppInformationRequiredState", m.GetDisplayAppInformationRequiredState())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("displayLocationInformationRequiredState", m.GetDisplayLocationInformationRequiredState())
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
func (m *MicrosoftAuthenticatorFeatureSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MicrosoftAuthenticatorFeatureSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDisplayAppInformationRequiredState sets the displayAppInformationRequiredState property value. Determines whether the user's Authenticator app shows them the client app they're signing into.
func (m *MicrosoftAuthenticatorFeatureSettings) SetDisplayAppInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)() {
    err := m.GetBackingStore().Set("displayAppInformationRequiredState", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayLocationInformationRequiredState sets the displayLocationInformationRequiredState property value. Determines whether the user's Authenticator app shows them the geographic location of where the authentication request originated from.
func (m *MicrosoftAuthenticatorFeatureSettings) SetDisplayLocationInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)() {
    err := m.GetBackingStore().Set("displayLocationInformationRequiredState", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MicrosoftAuthenticatorFeatureSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type MicrosoftAuthenticatorFeatureSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDisplayAppInformationRequiredState()(AuthenticationMethodFeatureConfigurationable)
    GetDisplayLocationInformationRequiredState()(AuthenticationMethodFeatureConfigurationable)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDisplayAppInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)()
    SetDisplayLocationInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)()
    SetOdataType(value *string)()
}

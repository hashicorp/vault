package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessSessionControls struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessSessionControls instantiates a new ConditionalAccessSessionControls and sets the default values.
func NewConditionalAccessSessionControls()(*ConditionalAccessSessionControls) {
    m := &ConditionalAccessSessionControls{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessSessionControlsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessSessionControlsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessSessionControls(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessSessionControls) GetAdditionalData()(map[string]any) {
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
// GetApplicationEnforcedRestrictions gets the applicationEnforcedRestrictions property value. Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.
// returns a ApplicationEnforcedRestrictionsSessionControlable when successful
func (m *ConditionalAccessSessionControls) GetApplicationEnforcedRestrictions()(ApplicationEnforcedRestrictionsSessionControlable) {
    val, err := m.GetBackingStore().Get("applicationEnforcedRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ApplicationEnforcedRestrictionsSessionControlable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ConditionalAccessSessionControls) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCloudAppSecurity gets the cloudAppSecurity property value. Session control to apply cloud app security.
// returns a CloudAppSecuritySessionControlable when successful
func (m *ConditionalAccessSessionControls) GetCloudAppSecurity()(CloudAppSecuritySessionControlable) {
    val, err := m.GetBackingStore().Get("cloudAppSecurity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CloudAppSecuritySessionControlable)
    }
    return nil
}
// GetDisableResilienceDefaults gets the disableResilienceDefaults property value. Session control that determines whether it is acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.
// returns a *bool when successful
func (m *ConditionalAccessSessionControls) GetDisableResilienceDefaults()(*bool) {
    val, err := m.GetBackingStore().Get("disableResilienceDefaults")
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
func (m *ConditionalAccessSessionControls) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["applicationEnforcedRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateApplicationEnforcedRestrictionsSessionControlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationEnforcedRestrictions(val.(ApplicationEnforcedRestrictionsSessionControlable))
        }
        return nil
    }
    res["cloudAppSecurity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCloudAppSecuritySessionControlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudAppSecurity(val.(CloudAppSecuritySessionControlable))
        }
        return nil
    }
    res["disableResilienceDefaults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisableResilienceDefaults(val)
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
    res["persistentBrowser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePersistentBrowserSessionControlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersistentBrowser(val.(PersistentBrowserSessionControlable))
        }
        return nil
    }
    res["signInFrequency"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSignInFrequencySessionControlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInFrequency(val.(SignInFrequencySessionControlable))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessSessionControls) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersistentBrowser gets the persistentBrowser property value. Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.
// returns a PersistentBrowserSessionControlable when successful
func (m *ConditionalAccessSessionControls) GetPersistentBrowser()(PersistentBrowserSessionControlable) {
    val, err := m.GetBackingStore().Get("persistentBrowser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PersistentBrowserSessionControlable)
    }
    return nil
}
// GetSignInFrequency gets the signInFrequency property value. Session control to enforce signin frequency.
// returns a SignInFrequencySessionControlable when successful
func (m *ConditionalAccessSessionControls) GetSignInFrequency()(SignInFrequencySessionControlable) {
    val, err := m.GetBackingStore().Get("signInFrequency")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SignInFrequencySessionControlable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessSessionControls) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("applicationEnforcedRestrictions", m.GetApplicationEnforcedRestrictions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("cloudAppSecurity", m.GetCloudAppSecurity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("disableResilienceDefaults", m.GetDisableResilienceDefaults())
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
        err := writer.WriteObjectValue("persistentBrowser", m.GetPersistentBrowser())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("signInFrequency", m.GetSignInFrequency())
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
func (m *ConditionalAccessSessionControls) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationEnforcedRestrictions sets the applicationEnforcedRestrictions property value. Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.
func (m *ConditionalAccessSessionControls) SetApplicationEnforcedRestrictions(value ApplicationEnforcedRestrictionsSessionControlable)() {
    err := m.GetBackingStore().Set("applicationEnforcedRestrictions", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessSessionControls) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCloudAppSecurity sets the cloudAppSecurity property value. Session control to apply cloud app security.
func (m *ConditionalAccessSessionControls) SetCloudAppSecurity(value CloudAppSecuritySessionControlable)() {
    err := m.GetBackingStore().Set("cloudAppSecurity", value)
    if err != nil {
        panic(err)
    }
}
// SetDisableResilienceDefaults sets the disableResilienceDefaults property value. Session control that determines whether it is acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.
func (m *ConditionalAccessSessionControls) SetDisableResilienceDefaults(value *bool)() {
    err := m.GetBackingStore().Set("disableResilienceDefaults", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessSessionControls) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPersistentBrowser sets the persistentBrowser property value. Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.
func (m *ConditionalAccessSessionControls) SetPersistentBrowser(value PersistentBrowserSessionControlable)() {
    err := m.GetBackingStore().Set("persistentBrowser", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInFrequency sets the signInFrequency property value. Session control to enforce signin frequency.
func (m *ConditionalAccessSessionControls) SetSignInFrequency(value SignInFrequencySessionControlable)() {
    err := m.GetBackingStore().Set("signInFrequency", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessSessionControlsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationEnforcedRestrictions()(ApplicationEnforcedRestrictionsSessionControlable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCloudAppSecurity()(CloudAppSecuritySessionControlable)
    GetDisableResilienceDefaults()(*bool)
    GetOdataType()(*string)
    GetPersistentBrowser()(PersistentBrowserSessionControlable)
    GetSignInFrequency()(SignInFrequencySessionControlable)
    SetApplicationEnforcedRestrictions(value ApplicationEnforcedRestrictionsSessionControlable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCloudAppSecurity(value CloudAppSecuritySessionControlable)()
    SetDisableResilienceDefaults(value *bool)()
    SetOdataType(value *string)()
    SetPersistentBrowser(value PersistentBrowserSessionControlable)()
    SetSignInFrequency(value SignInFrequencySessionControlable)()
}

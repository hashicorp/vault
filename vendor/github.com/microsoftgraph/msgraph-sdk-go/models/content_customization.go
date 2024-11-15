package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ContentCustomization struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewContentCustomization instantiates a new ContentCustomization and sets the default values.
func NewContentCustomization()(*ContentCustomization) {
    m := &ContentCustomization{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateContentCustomizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateContentCustomizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContentCustomization(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ContentCustomization) GetAdditionalData()(map[string]any) {
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
// GetAttributeCollection gets the attributeCollection property value. Represents the content options of External Identities to be customized throughout the authentication flow for a tenant.
// returns a []KeyValueable when successful
func (m *ContentCustomization) GetAttributeCollection()([]KeyValueable) {
    val, err := m.GetBackingStore().Get("attributeCollection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValueable)
    }
    return nil
}
// GetAttributeCollectionRelativeUrl gets the attributeCollectionRelativeUrl property value. A relative URL for the content options of External Identities to be customized throughout the authentication flow for a tenant.
// returns a *string when successful
func (m *ContentCustomization) GetAttributeCollectionRelativeUrl()(*string) {
    val, err := m.GetBackingStore().Get("attributeCollectionRelativeUrl")
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
func (m *ContentCustomization) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ContentCustomization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attributeCollection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValueable)
                }
            }
            m.SetAttributeCollection(res)
        }
        return nil
    }
    res["attributeCollectionRelativeUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttributeCollectionRelativeUrl(val)
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
    res["registrationCampaign"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValueable)
                }
            }
            m.SetRegistrationCampaign(res)
        }
        return nil
    }
    res["registrationCampaignRelativeUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistrationCampaignRelativeUrl(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ContentCustomization) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistrationCampaign gets the registrationCampaign property value. Represents content options to customize during MFA proofup interruptions.
// returns a []KeyValueable when successful
func (m *ContentCustomization) GetRegistrationCampaign()([]KeyValueable) {
    val, err := m.GetBackingStore().Get("registrationCampaign")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValueable)
    }
    return nil
}
// GetRegistrationCampaignRelativeUrl gets the registrationCampaignRelativeUrl property value. The relative URL of the content options to customize during MFA proofup interruptions.
// returns a *string when successful
func (m *ContentCustomization) GetRegistrationCampaignRelativeUrl()(*string) {
    val, err := m.GetBackingStore().Get("registrationCampaignRelativeUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ContentCustomization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAttributeCollection() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttributeCollection()))
        for i, v := range m.GetAttributeCollection() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("attributeCollection", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("attributeCollectionRelativeUrl", m.GetAttributeCollectionRelativeUrl())
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
    if m.GetRegistrationCampaign() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRegistrationCampaign()))
        for i, v := range m.GetRegistrationCampaign() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("registrationCampaign", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("registrationCampaignRelativeUrl", m.GetRegistrationCampaignRelativeUrl())
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
func (m *ContentCustomization) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttributeCollection sets the attributeCollection property value. Represents the content options of External Identities to be customized throughout the authentication flow for a tenant.
func (m *ContentCustomization) SetAttributeCollection(value []KeyValueable)() {
    err := m.GetBackingStore().Set("attributeCollection", value)
    if err != nil {
        panic(err)
    }
}
// SetAttributeCollectionRelativeUrl sets the attributeCollectionRelativeUrl property value. A relative URL for the content options of External Identities to be customized throughout the authentication flow for a tenant.
func (m *ContentCustomization) SetAttributeCollectionRelativeUrl(value *string)() {
    err := m.GetBackingStore().Set("attributeCollectionRelativeUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ContentCustomization) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ContentCustomization) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistrationCampaign sets the registrationCampaign property value. Represents content options to customize during MFA proofup interruptions.
func (m *ContentCustomization) SetRegistrationCampaign(value []KeyValueable)() {
    err := m.GetBackingStore().Set("registrationCampaign", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistrationCampaignRelativeUrl sets the registrationCampaignRelativeUrl property value. The relative URL of the content options to customize during MFA proofup interruptions.
func (m *ContentCustomization) SetRegistrationCampaignRelativeUrl(value *string)() {
    err := m.GetBackingStore().Set("registrationCampaignRelativeUrl", value)
    if err != nil {
        panic(err)
    }
}
type ContentCustomizationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttributeCollection()([]KeyValueable)
    GetAttributeCollectionRelativeUrl()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetRegistrationCampaign()([]KeyValueable)
    GetRegistrationCampaignRelativeUrl()(*string)
    SetAttributeCollection(value []KeyValueable)()
    SetAttributeCollectionRelativeUrl(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetRegistrationCampaign(value []KeyValueable)()
    SetRegistrationCampaignRelativeUrl(value *string)()
}

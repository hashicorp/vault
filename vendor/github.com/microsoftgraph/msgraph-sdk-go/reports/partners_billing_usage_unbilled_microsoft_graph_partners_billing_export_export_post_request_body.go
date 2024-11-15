package reports

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
    ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4 "github.com/microsoftgraph/msgraph-sdk-go/models/partners/billing"
)

type PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody instantiates a new PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody and sets the default values.
func NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody()(*PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) {
    m := &PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetAdditionalData()(map[string]any) {
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
// GetAttributeSet gets the attributeSet property value. The attributeSet property
// returns a *AttributeSet when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetAttributeSet()(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet) {
    val, err := m.GetBackingStore().Get("attributeSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBillingPeriod gets the billingPeriod property value. The billingPeriod property
// returns a *BillingPeriod when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetBillingPeriod()(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod) {
    val, err := m.GetBackingStore().Get("billingPeriod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod)
    }
    return nil
}
// GetCurrencyCode gets the currencyCode property value. The currencyCode property
// returns a *string when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetCurrencyCode()(*string) {
    val, err := m.GetBackingStore().Get("currencyCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attributeSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.ParseAttributeSet)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttributeSet(val.(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet))
        }
        return nil
    }
    res["billingPeriod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.ParseBillingPeriod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBillingPeriod(val.(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod))
        }
        return nil
    }
    res["currencyCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCurrencyCode(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAttributeSet() != nil {
        cast := (*m.GetAttributeSet()).String()
        err := writer.WriteStringValue("attributeSet", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetBillingPeriod() != nil {
        cast := (*m.GetBillingPeriod()).String()
        err := writer.WriteStringValue("billingPeriod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("currencyCode", m.GetCurrencyCode())
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
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttributeSet sets the attributeSet property value. The attributeSet property
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) SetAttributeSet(value *ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet)() {
    err := m.GetBackingStore().Set("attributeSet", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBillingPeriod sets the billingPeriod property value. The billingPeriod property
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) SetBillingPeriod(value *ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod)() {
    err := m.GetBackingStore().Set("billingPeriod", value)
    if err != nil {
        panic(err)
    }
}
// SetCurrencyCode sets the currencyCode property value. The currencyCode property
func (m *PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBody) SetCurrencyCode(value *string)() {
    err := m.GetBackingStore().Set("currencyCode", value)
    if err != nil {
        panic(err)
    }
}
type PartnersBillingUsageUnbilledMicrosoftGraphPartnersBillingExportExportPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttributeSet()(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBillingPeriod()(*ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod)
    GetCurrencyCode()(*string)
    SetAttributeSet(value *ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.AttributeSet)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBillingPeriod(value *ieaa1d050ea8ba883c482e05cf2306cb5376cc6e2cf5966c1a6850c42c6118fa4.BillingPeriod)()
    SetCurrencyCode(value *string)()
}

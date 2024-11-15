package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type LoginPageLayoutConfiguration struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewLoginPageLayoutConfiguration instantiates a new LoginPageLayoutConfiguration and sets the default values.
func NewLoginPageLayoutConfiguration()(*LoginPageLayoutConfiguration) {
    m := &LoginPageLayoutConfiguration{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateLoginPageLayoutConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLoginPageLayoutConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLoginPageLayoutConfiguration(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *LoginPageLayoutConfiguration) GetAdditionalData()(map[string]any) {
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
func (m *LoginPageLayoutConfiguration) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LoginPageLayoutConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isFooterShown"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFooterShown(val)
        }
        return nil
    }
    res["isHeaderShown"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsHeaderShown(val)
        }
        return nil
    }
    res["layoutTemplateType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLayoutTemplateType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLayoutTemplateType(val.(*LayoutTemplateType))
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
// GetIsFooterShown gets the isFooterShown property value. Option to show the footer on the sign-in page.
// returns a *bool when successful
func (m *LoginPageLayoutConfiguration) GetIsFooterShown()(*bool) {
    val, err := m.GetBackingStore().Get("isFooterShown")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsHeaderShown gets the isHeaderShown property value. Option to show the header on the sign-in page.
// returns a *bool when successful
func (m *LoginPageLayoutConfiguration) GetIsHeaderShown()(*bool) {
    val, err := m.GetBackingStore().Get("isHeaderShown")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLayoutTemplateType gets the layoutTemplateType property value. Represents the layout template to be displayed on the login page for a tenant. The possible values are  default - Represents the default Microsoft layout with a centered lightbox.  verticalSplit - Represents a layout with a background on the left side and a full-height lightbox to the right.  unknownFutureValue - Evolvable enumeration sentinel value. Don't use.
// returns a *LayoutTemplateType when successful
func (m *LoginPageLayoutConfiguration) GetLayoutTemplateType()(*LayoutTemplateType) {
    val, err := m.GetBackingStore().Get("layoutTemplateType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*LayoutTemplateType)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *LoginPageLayoutConfiguration) GetOdataType()(*string) {
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
func (m *LoginPageLayoutConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isFooterShown", m.GetIsFooterShown())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isHeaderShown", m.GetIsHeaderShown())
        if err != nil {
            return err
        }
    }
    if m.GetLayoutTemplateType() != nil {
        cast := (*m.GetLayoutTemplateType()).String()
        err := writer.WriteStringValue("layoutTemplateType", &cast)
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
func (m *LoginPageLayoutConfiguration) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *LoginPageLayoutConfiguration) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsFooterShown sets the isFooterShown property value. Option to show the footer on the sign-in page.
func (m *LoginPageLayoutConfiguration) SetIsFooterShown(value *bool)() {
    err := m.GetBackingStore().Set("isFooterShown", value)
    if err != nil {
        panic(err)
    }
}
// SetIsHeaderShown sets the isHeaderShown property value. Option to show the header on the sign-in page.
func (m *LoginPageLayoutConfiguration) SetIsHeaderShown(value *bool)() {
    err := m.GetBackingStore().Set("isHeaderShown", value)
    if err != nil {
        panic(err)
    }
}
// SetLayoutTemplateType sets the layoutTemplateType property value. Represents the layout template to be displayed on the login page for a tenant. The possible values are  default - Represents the default Microsoft layout with a centered lightbox.  verticalSplit - Represents a layout with a background on the left side and a full-height lightbox to the right.  unknownFutureValue - Evolvable enumeration sentinel value. Don't use.
func (m *LoginPageLayoutConfiguration) SetLayoutTemplateType(value *LayoutTemplateType)() {
    err := m.GetBackingStore().Set("layoutTemplateType", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *LoginPageLayoutConfiguration) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type LoginPageLayoutConfigurationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsFooterShown()(*bool)
    GetIsHeaderShown()(*bool)
    GetLayoutTemplateType()(*LayoutTemplateType)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsFooterShown(value *bool)()
    SetIsHeaderShown(value *bool)()
    SetLayoutTemplateType(value *LayoutTemplateType)()
    SetOdataType(value *string)()
}

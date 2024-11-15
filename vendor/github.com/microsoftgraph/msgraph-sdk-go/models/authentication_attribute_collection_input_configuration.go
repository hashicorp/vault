package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AuthenticationAttributeCollectionInputConfiguration struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAuthenticationAttributeCollectionInputConfiguration instantiates a new AuthenticationAttributeCollectionInputConfiguration and sets the default values.
func NewAuthenticationAttributeCollectionInputConfiguration()(*AuthenticationAttributeCollectionInputConfiguration) {
    m := &AuthenticationAttributeCollectionInputConfiguration{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAuthenticationAttributeCollectionInputConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationAttributeCollectionInputConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationAttributeCollectionInputConfiguration(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetAdditionalData()(map[string]any) {
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
// GetAttribute gets the attribute property value. The built-in or custom attribute for which a value is being collected.
// returns a *string when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetAttribute()(*string) {
    val, err := m.GetBackingStore().Get("attribute")
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
func (m *AuthenticationAttributeCollectionInputConfiguration) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDefaultValue gets the defaultValue property value. The default value of the attribute displayed to the end user. The capability to set the default value isn't available through the Microsoft Entra admin center.
// returns a *string when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetDefaultValue()(*string) {
    val, err := m.GetBackingStore().Get("defaultValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEditable gets the editable property value. Defines whether the attribute is editable by the end user.
// returns a *bool when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetEditable()(*bool) {
    val, err := m.GetBackingStore().Get("editable")
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
func (m *AuthenticationAttributeCollectionInputConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attribute"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttribute(val)
        }
        return nil
    }
    res["defaultValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultValue(val)
        }
        return nil
    }
    res["editable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEditable(val)
        }
        return nil
    }
    res["hidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHidden(val)
        }
        return nil
    }
    res["inputType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationAttributeCollectionInputType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInputType(val.(*AuthenticationAttributeCollectionInputType))
        }
        return nil
    }
    res["label"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabel(val)
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
    res["options"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationAttributeCollectionOptionConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationAttributeCollectionOptionConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationAttributeCollectionOptionConfigurationable)
                }
            }
            m.SetOptions(res)
        }
        return nil
    }
    res["required"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequired(val)
        }
        return nil
    }
    res["validationRegEx"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValidationRegEx(val)
        }
        return nil
    }
    res["writeToDirectory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWriteToDirectory(val)
        }
        return nil
    }
    return res
}
// GetHidden gets the hidden property value. Defines whether the attribute is displayed to the end user. The capability to hide isn't available through the Microsoft Entra admin center.
// returns a *bool when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetHidden()(*bool) {
    val, err := m.GetBackingStore().Get("hidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInputType gets the inputType property value. The inputType property
// returns a *AuthenticationAttributeCollectionInputType when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetInputType()(*AuthenticationAttributeCollectionInputType) {
    val, err := m.GetBackingStore().Get("inputType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationAttributeCollectionInputType)
    }
    return nil
}
// GetLabel gets the label property value. The label of the attribute field that's displayed to end user, unless overridden.
// returns a *string when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetLabel()(*string) {
    val, err := m.GetBackingStore().Get("label")
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
func (m *AuthenticationAttributeCollectionInputConfiguration) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOptions gets the options property value. The option values for certain multiple-option input types.
// returns a []AuthenticationAttributeCollectionOptionConfigurationable when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetOptions()([]AuthenticationAttributeCollectionOptionConfigurationable) {
    val, err := m.GetBackingStore().Get("options")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationAttributeCollectionOptionConfigurationable)
    }
    return nil
}
// GetRequired gets the required property value. Defines whether the field is required.
// returns a *bool when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetRequired()(*bool) {
    val, err := m.GetBackingStore().Get("required")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetValidationRegEx gets the validationRegEx property value. The regex for the value of the field. For more information about the supported regexes, see validationRegEx values for inputType objects. To understand how to specify regexes, see the Regular expressions cheat sheet.
// returns a *string when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetValidationRegEx()(*string) {
    val, err := m.GetBackingStore().Get("validationRegEx")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWriteToDirectory gets the writeToDirectory property value. Defines whether Microsoft Entra ID stores the value that it collects.
// returns a *bool when successful
func (m *AuthenticationAttributeCollectionInputConfiguration) GetWriteToDirectory()(*bool) {
    val, err := m.GetBackingStore().Get("writeToDirectory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationAttributeCollectionInputConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("attribute", m.GetAttribute())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("editable", m.GetEditable())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    if m.GetInputType() != nil {
        cast := (*m.GetInputType()).String()
        err := writer.WriteStringValue("inputType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("label", m.GetLabel())
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
    if m.GetOptions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOptions()))
        for i, v := range m.GetOptions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("options", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("required", m.GetRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("validationRegEx", m.GetValidationRegEx())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("writeToDirectory", m.GetWriteToDirectory())
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
func (m *AuthenticationAttributeCollectionInputConfiguration) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttribute sets the attribute property value. The built-in or custom attribute for which a value is being collected.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetAttribute(value *string)() {
    err := m.GetBackingStore().Set("attribute", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDefaultValue sets the defaultValue property value. The default value of the attribute displayed to the end user. The capability to set the default value isn't available through the Microsoft Entra admin center.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetDefaultValue(value *string)() {
    err := m.GetBackingStore().Set("defaultValue", value)
    if err != nil {
        panic(err)
    }
}
// SetEditable sets the editable property value. Defines whether the attribute is editable by the end user.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetEditable(value *bool)() {
    err := m.GetBackingStore().Set("editable", value)
    if err != nil {
        panic(err)
    }
}
// SetHidden sets the hidden property value. Defines whether the attribute is displayed to the end user. The capability to hide isn't available through the Microsoft Entra admin center.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetHidden(value *bool)() {
    err := m.GetBackingStore().Set("hidden", value)
    if err != nil {
        panic(err)
    }
}
// SetInputType sets the inputType property value. The inputType property
func (m *AuthenticationAttributeCollectionInputConfiguration) SetInputType(value *AuthenticationAttributeCollectionInputType)() {
    err := m.GetBackingStore().Set("inputType", value)
    if err != nil {
        panic(err)
    }
}
// SetLabel sets the label property value. The label of the attribute field that's displayed to end user, unless overridden.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetLabel(value *string)() {
    err := m.GetBackingStore().Set("label", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AuthenticationAttributeCollectionInputConfiguration) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOptions sets the options property value. The option values for certain multiple-option input types.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetOptions(value []AuthenticationAttributeCollectionOptionConfigurationable)() {
    err := m.GetBackingStore().Set("options", value)
    if err != nil {
        panic(err)
    }
}
// SetRequired sets the required property value. Defines whether the field is required.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetRequired(value *bool)() {
    err := m.GetBackingStore().Set("required", value)
    if err != nil {
        panic(err)
    }
}
// SetValidationRegEx sets the validationRegEx property value. The regex for the value of the field. For more information about the supported regexes, see validationRegEx values for inputType objects. To understand how to specify regexes, see the Regular expressions cheat sheet.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetValidationRegEx(value *string)() {
    err := m.GetBackingStore().Set("validationRegEx", value)
    if err != nil {
        panic(err)
    }
}
// SetWriteToDirectory sets the writeToDirectory property value. Defines whether Microsoft Entra ID stores the value that it collects.
func (m *AuthenticationAttributeCollectionInputConfiguration) SetWriteToDirectory(value *bool)() {
    err := m.GetBackingStore().Set("writeToDirectory", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationAttributeCollectionInputConfigurationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttribute()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDefaultValue()(*string)
    GetEditable()(*bool)
    GetHidden()(*bool)
    GetInputType()(*AuthenticationAttributeCollectionInputType)
    GetLabel()(*string)
    GetOdataType()(*string)
    GetOptions()([]AuthenticationAttributeCollectionOptionConfigurationable)
    GetRequired()(*bool)
    GetValidationRegEx()(*string)
    GetWriteToDirectory()(*bool)
    SetAttribute(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDefaultValue(value *string)()
    SetEditable(value *bool)()
    SetHidden(value *bool)()
    SetInputType(value *AuthenticationAttributeCollectionInputType)()
    SetLabel(value *string)()
    SetOdataType(value *string)()
    SetOptions(value []AuthenticationAttributeCollectionOptionConfigurationable)()
    SetRequired(value *bool)()
    SetValidationRegEx(value *string)()
    SetWriteToDirectory(value *bool)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessGrantControls struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessGrantControls instantiates a new ConditionalAccessGrantControls and sets the default values.
func NewConditionalAccessGrantControls()(*ConditionalAccessGrantControls) {
    m := &ConditionalAccessGrantControls{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessGrantControlsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessGrantControlsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessGrantControls(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessGrantControls) GetAdditionalData()(map[string]any) {
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
// GetAuthenticationStrength gets the authenticationStrength property value. The authenticationStrength property
// returns a AuthenticationStrengthPolicyable when successful
func (m *ConditionalAccessGrantControls) GetAuthenticationStrength()(AuthenticationStrengthPolicyable) {
    val, err := m.GetBackingStore().Get("authenticationStrength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationStrengthPolicyable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ConditionalAccessGrantControls) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBuiltInControls gets the builtInControls property value. List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.
// returns a []ConditionalAccessGrantControl when successful
func (m *ConditionalAccessGrantControls) GetBuiltInControls()([]ConditionalAccessGrantControl) {
    val, err := m.GetBackingStore().Get("builtInControls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConditionalAccessGrantControl)
    }
    return nil
}
// GetCustomAuthenticationFactors gets the customAuthenticationFactors property value. List of custom controls IDs required by the policy. For more information, see Custom controls.
// returns a []string when successful
func (m *ConditionalAccessGrantControls) GetCustomAuthenticationFactors()([]string) {
    val, err := m.GetBackingStore().Get("customAuthenticationFactors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessGrantControls) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["authenticationStrength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationStrengthPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationStrength(val.(AuthenticationStrengthPolicyable))
        }
        return nil
    }
    res["builtInControls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseConditionalAccessGrantControl)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConditionalAccessGrantControl, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*ConditionalAccessGrantControl))
                }
            }
            m.SetBuiltInControls(res)
        }
        return nil
    }
    res["customAuthenticationFactors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetCustomAuthenticationFactors(res)
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
    res["operator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperator(val)
        }
        return nil
    }
    res["termsOfUse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetTermsOfUse(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessGrantControls) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperator gets the operator property value. Defines the relationship of the grant controls. Possible values: AND, OR.
// returns a *string when successful
func (m *ConditionalAccessGrantControls) GetOperator()(*string) {
    val, err := m.GetBackingStore().Get("operator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTermsOfUse gets the termsOfUse property value. List of terms of use IDs required by the policy.
// returns a []string when successful
func (m *ConditionalAccessGrantControls) GetTermsOfUse()([]string) {
    val, err := m.GetBackingStore().Get("termsOfUse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessGrantControls) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("authenticationStrength", m.GetAuthenticationStrength())
        if err != nil {
            return err
        }
    }
    if m.GetBuiltInControls() != nil {
        err := writer.WriteCollectionOfStringValues("builtInControls", SerializeConditionalAccessGrantControl(m.GetBuiltInControls()))
        if err != nil {
            return err
        }
    }
    if m.GetCustomAuthenticationFactors() != nil {
        err := writer.WriteCollectionOfStringValues("customAuthenticationFactors", m.GetCustomAuthenticationFactors())
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
        err := writer.WriteStringValue("operator", m.GetOperator())
        if err != nil {
            return err
        }
    }
    if m.GetTermsOfUse() != nil {
        err := writer.WriteCollectionOfStringValues("termsOfUse", m.GetTermsOfUse())
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
func (m *ConditionalAccessGrantControls) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthenticationStrength sets the authenticationStrength property value. The authenticationStrength property
func (m *ConditionalAccessGrantControls) SetAuthenticationStrength(value AuthenticationStrengthPolicyable)() {
    err := m.GetBackingStore().Set("authenticationStrength", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessGrantControls) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBuiltInControls sets the builtInControls property value. List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.
func (m *ConditionalAccessGrantControls) SetBuiltInControls(value []ConditionalAccessGrantControl)() {
    err := m.GetBackingStore().Set("builtInControls", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomAuthenticationFactors sets the customAuthenticationFactors property value. List of custom controls IDs required by the policy. For more information, see Custom controls.
func (m *ConditionalAccessGrantControls) SetCustomAuthenticationFactors(value []string)() {
    err := m.GetBackingStore().Set("customAuthenticationFactors", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessGrantControls) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperator sets the operator property value. Defines the relationship of the grant controls. Possible values: AND, OR.
func (m *ConditionalAccessGrantControls) SetOperator(value *string)() {
    err := m.GetBackingStore().Set("operator", value)
    if err != nil {
        panic(err)
    }
}
// SetTermsOfUse sets the termsOfUse property value. List of terms of use IDs required by the policy.
func (m *ConditionalAccessGrantControls) SetTermsOfUse(value []string)() {
    err := m.GetBackingStore().Set("termsOfUse", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessGrantControlsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationStrength()(AuthenticationStrengthPolicyable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBuiltInControls()([]ConditionalAccessGrantControl)
    GetCustomAuthenticationFactors()([]string)
    GetOdataType()(*string)
    GetOperator()(*string)
    GetTermsOfUse()([]string)
    SetAuthenticationStrength(value AuthenticationStrengthPolicyable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBuiltInControls(value []ConditionalAccessGrantControl)()
    SetCustomAuthenticationFactors(value []string)()
    SetOdataType(value *string)()
    SetOperator(value *string)()
    SetTermsOfUse(value []string)()
}

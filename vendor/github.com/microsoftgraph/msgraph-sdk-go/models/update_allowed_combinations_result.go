package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UpdateAllowedCombinationsResult struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUpdateAllowedCombinationsResult instantiates a new UpdateAllowedCombinationsResult and sets the default values.
func NewUpdateAllowedCombinationsResult()(*UpdateAllowedCombinationsResult) {
    m := &UpdateAllowedCombinationsResult{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUpdateAllowedCombinationsResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUpdateAllowedCombinationsResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUpdateAllowedCombinationsResult(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UpdateAllowedCombinationsResult) GetAdditionalData()(map[string]any) {
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
// GetAdditionalInformation gets the additionalInformation property value. Information about why the updateAllowedCombinations action was successful or failed.
// returns a *string when successful
func (m *UpdateAllowedCombinationsResult) GetAdditionalInformation()(*string) {
    val, err := m.GetBackingStore().Get("additionalInformation")
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
func (m *UpdateAllowedCombinationsResult) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConditionalAccessReferences gets the conditionalAccessReferences property value. References to existing Conditional Access policies that use this authentication strength.
// returns a []string when successful
func (m *UpdateAllowedCombinationsResult) GetConditionalAccessReferences()([]string) {
    val, err := m.GetBackingStore().Get("conditionalAccessReferences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCurrentCombinations gets the currentCombinations property value. The list of current authentication method combinations allowed by the authentication strength.
// returns a []AuthenticationMethodModes when successful
func (m *UpdateAllowedCombinationsResult) GetCurrentCombinations()([]AuthenticationMethodModes) {
    val, err := m.GetBackingStore().Get("currentCombinations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodModes)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UpdateAllowedCombinationsResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["additionalInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAdditionalInformation(val)
        }
        return nil
    }
    res["conditionalAccessReferences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetConditionalAccessReferences(res)
        }
        return nil
    }
    res["currentCombinations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseAuthenticationMethodModes)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodModes, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*AuthenticationMethodModes))
                }
            }
            m.SetCurrentCombinations(res)
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
    res["previousCombinations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseAuthenticationMethodModes)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodModes, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*AuthenticationMethodModes))
                }
            }
            m.SetPreviousCombinations(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UpdateAllowedCombinationsResult) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreviousCombinations gets the previousCombinations property value. The list of former authentication method combinations allowed by the authentication strength before they were updated through the updateAllowedCombinations action.
// returns a []AuthenticationMethodModes when successful
func (m *UpdateAllowedCombinationsResult) GetPreviousCombinations()([]AuthenticationMethodModes) {
    val, err := m.GetBackingStore().Get("previousCombinations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodModes)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UpdateAllowedCombinationsResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("additionalInformation", m.GetAdditionalInformation())
        if err != nil {
            return err
        }
    }
    if m.GetConditionalAccessReferences() != nil {
        err := writer.WriteCollectionOfStringValues("conditionalAccessReferences", m.GetConditionalAccessReferences())
        if err != nil {
            return err
        }
    }
    if m.GetCurrentCombinations() != nil {
        err := writer.WriteCollectionOfStringValues("currentCombinations", SerializeAuthenticationMethodModes(m.GetCurrentCombinations()))
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
    if m.GetPreviousCombinations() != nil {
        err := writer.WriteCollectionOfStringValues("previousCombinations", SerializeAuthenticationMethodModes(m.GetPreviousCombinations()))
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
func (m *UpdateAllowedCombinationsResult) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalInformation sets the additionalInformation property value. Information about why the updateAllowedCombinations action was successful or failed.
func (m *UpdateAllowedCombinationsResult) SetAdditionalInformation(value *string)() {
    err := m.GetBackingStore().Set("additionalInformation", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UpdateAllowedCombinationsResult) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConditionalAccessReferences sets the conditionalAccessReferences property value. References to existing Conditional Access policies that use this authentication strength.
func (m *UpdateAllowedCombinationsResult) SetConditionalAccessReferences(value []string)() {
    err := m.GetBackingStore().Set("conditionalAccessReferences", value)
    if err != nil {
        panic(err)
    }
}
// SetCurrentCombinations sets the currentCombinations property value. The list of current authentication method combinations allowed by the authentication strength.
func (m *UpdateAllowedCombinationsResult) SetCurrentCombinations(value []AuthenticationMethodModes)() {
    err := m.GetBackingStore().Set("currentCombinations", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UpdateAllowedCombinationsResult) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviousCombinations sets the previousCombinations property value. The list of former authentication method combinations allowed by the authentication strength before they were updated through the updateAllowedCombinations action.
func (m *UpdateAllowedCombinationsResult) SetPreviousCombinations(value []AuthenticationMethodModes)() {
    err := m.GetBackingStore().Set("previousCombinations", value)
    if err != nil {
        panic(err)
    }
}
type UpdateAllowedCombinationsResultable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdditionalInformation()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConditionalAccessReferences()([]string)
    GetCurrentCombinations()([]AuthenticationMethodModes)
    GetOdataType()(*string)
    GetPreviousCombinations()([]AuthenticationMethodModes)
    SetAdditionalInformation(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConditionalAccessReferences(value []string)()
    SetCurrentCombinations(value []AuthenticationMethodModes)()
    SetOdataType(value *string)()
    SetPreviousCombinations(value []AuthenticationMethodModes)()
}

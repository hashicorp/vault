package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Fido2KeyRestrictions struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewFido2KeyRestrictions instantiates a new Fido2KeyRestrictions and sets the default values.
func NewFido2KeyRestrictions()(*Fido2KeyRestrictions) {
    m := &Fido2KeyRestrictions{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateFido2KeyRestrictionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFido2KeyRestrictionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFido2KeyRestrictions(), nil
}
// GetAaGuids gets the aaGuids property value. A collection of Authenticator Attestation GUIDs. AADGUIDs define key types and manufacturers.
// returns a []string when successful
func (m *Fido2KeyRestrictions) GetAaGuids()([]string) {
    val, err := m.GetBackingStore().Get("aaGuids")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Fido2KeyRestrictions) GetAdditionalData()(map[string]any) {
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
func (m *Fido2KeyRestrictions) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEnforcementType gets the enforcementType property value. Enforcement type. Possible values are: allow, block.
// returns a *Fido2RestrictionEnforcementType when successful
func (m *Fido2KeyRestrictions) GetEnforcementType()(*Fido2RestrictionEnforcementType) {
    val, err := m.GetBackingStore().Get("enforcementType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Fido2RestrictionEnforcementType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Fido2KeyRestrictions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["aaGuids"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAaGuids(res)
        }
        return nil
    }
    res["enforcementType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFido2RestrictionEnforcementType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnforcementType(val.(*Fido2RestrictionEnforcementType))
        }
        return nil
    }
    res["isEnforced"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnforced(val)
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
// GetIsEnforced gets the isEnforced property value. Determines if the configured key enforcement is enabled.
// returns a *bool when successful
func (m *Fido2KeyRestrictions) GetIsEnforced()(*bool) {
    val, err := m.GetBackingStore().Get("isEnforced")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Fido2KeyRestrictions) GetOdataType()(*string) {
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
func (m *Fido2KeyRestrictions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAaGuids() != nil {
        err := writer.WriteCollectionOfStringValues("aaGuids", m.GetAaGuids())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcementType() != nil {
        cast := (*m.GetEnforcementType()).String()
        err := writer.WriteStringValue("enforcementType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEnforced", m.GetIsEnforced())
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
// SetAaGuids sets the aaGuids property value. A collection of Authenticator Attestation GUIDs. AADGUIDs define key types and manufacturers.
func (m *Fido2KeyRestrictions) SetAaGuids(value []string)() {
    err := m.GetBackingStore().Set("aaGuids", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Fido2KeyRestrictions) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Fido2KeyRestrictions) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEnforcementType sets the enforcementType property value. Enforcement type. Possible values are: allow, block.
func (m *Fido2KeyRestrictions) SetEnforcementType(value *Fido2RestrictionEnforcementType)() {
    err := m.GetBackingStore().Set("enforcementType", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnforced sets the isEnforced property value. Determines if the configured key enforcement is enabled.
func (m *Fido2KeyRestrictions) SetIsEnforced(value *bool)() {
    err := m.GetBackingStore().Set("isEnforced", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Fido2KeyRestrictions) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type Fido2KeyRestrictionsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAaGuids()([]string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEnforcementType()(*Fido2RestrictionEnforcementType)
    GetIsEnforced()(*bool)
    GetOdataType()(*string)
    SetAaGuids(value []string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEnforcementType(value *Fido2RestrictionEnforcementType)()
    SetIsEnforced(value *bool)()
    SetOdataType(value *string)()
}

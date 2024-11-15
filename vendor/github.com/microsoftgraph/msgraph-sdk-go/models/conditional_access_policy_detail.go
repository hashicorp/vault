package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessPolicyDetail struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessPolicyDetail instantiates a new ConditionalAccessPolicyDetail and sets the default values.
func NewConditionalAccessPolicyDetail()(*ConditionalAccessPolicyDetail) {
    m := &ConditionalAccessPolicyDetail{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessPolicyDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessPolicyDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessPolicyDetail(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessPolicyDetail) GetAdditionalData()(map[string]any) {
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
func (m *ConditionalAccessPolicyDetail) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConditions gets the conditions property value. The conditions property
// returns a ConditionalAccessConditionSetable when successful
func (m *ConditionalAccessPolicyDetail) GetConditions()(ConditionalAccessConditionSetable) {
    val, err := m.GetBackingStore().Get("conditions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessConditionSetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessPolicyDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["conditions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessConditionSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditions(val.(ConditionalAccessConditionSetable))
        }
        return nil
    }
    res["grantControls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessGrantControlsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrantControls(val.(ConditionalAccessGrantControlsable))
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
    res["sessionControls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessSessionControlsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSessionControls(val.(ConditionalAccessSessionControlsable))
        }
        return nil
    }
    return res
}
// GetGrantControls gets the grantControls property value. Represents grant controls that must be fulfilled for the policy.
// returns a ConditionalAccessGrantControlsable when successful
func (m *ConditionalAccessPolicyDetail) GetGrantControls()(ConditionalAccessGrantControlsable) {
    val, err := m.GetBackingStore().Get("grantControls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessGrantControlsable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessPolicyDetail) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSessionControls gets the sessionControls property value. Represents a complex type of session controls that is enforced after sign-in.
// returns a ConditionalAccessSessionControlsable when successful
func (m *ConditionalAccessPolicyDetail) GetSessionControls()(ConditionalAccessSessionControlsable) {
    val, err := m.GetBackingStore().Get("sessionControls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessSessionControlsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessPolicyDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("conditions", m.GetConditions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("grantControls", m.GetGrantControls())
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
        err := writer.WriteObjectValue("sessionControls", m.GetSessionControls())
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
func (m *ConditionalAccessPolicyDetail) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessPolicyDetail) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConditions sets the conditions property value. The conditions property
func (m *ConditionalAccessPolicyDetail) SetConditions(value ConditionalAccessConditionSetable)() {
    err := m.GetBackingStore().Set("conditions", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantControls sets the grantControls property value. Represents grant controls that must be fulfilled for the policy.
func (m *ConditionalAccessPolicyDetail) SetGrantControls(value ConditionalAccessGrantControlsable)() {
    err := m.GetBackingStore().Set("grantControls", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessPolicyDetail) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSessionControls sets the sessionControls property value. Represents a complex type of session controls that is enforced after sign-in.
func (m *ConditionalAccessPolicyDetail) SetSessionControls(value ConditionalAccessSessionControlsable)() {
    err := m.GetBackingStore().Set("sessionControls", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessPolicyDetailable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConditions()(ConditionalAccessConditionSetable)
    GetGrantControls()(ConditionalAccessGrantControlsable)
    GetOdataType()(*string)
    GetSessionControls()(ConditionalAccessSessionControlsable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConditions(value ConditionalAccessConditionSetable)()
    SetGrantControls(value ConditionalAccessGrantControlsable)()
    SetOdataType(value *string)()
    SetSessionControls(value ConditionalAccessSessionControlsable)()
}

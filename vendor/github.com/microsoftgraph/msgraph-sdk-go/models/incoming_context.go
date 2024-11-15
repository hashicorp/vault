package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type IncomingContext struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewIncomingContext instantiates a new IncomingContext and sets the default values.
func NewIncomingContext()(*IncomingContext) {
    m := &IncomingContext{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateIncomingContextFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIncomingContextFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIncomingContext(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *IncomingContext) GetAdditionalData()(map[string]any) {
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
func (m *IncomingContext) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IncomingContext) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["observedParticipantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetObservedParticipantId(val)
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
    res["onBehalfOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnBehalfOf(val.(IdentitySetable))
        }
        return nil
    }
    res["sourceParticipantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceParticipantId(val)
        }
        return nil
    }
    res["transferor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTransferor(val.(IdentitySetable))
        }
        return nil
    }
    return res
}
// GetObservedParticipantId gets the observedParticipantId property value. The ID of the participant that is under observation. Read-only.
// returns a *string when successful
func (m *IncomingContext) GetObservedParticipantId()(*string) {
    val, err := m.GetBackingStore().Get("observedParticipantId")
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
func (m *IncomingContext) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnBehalfOf gets the onBehalfOf property value. The identity that the call is happening on behalf of.
// returns a IdentitySetable when successful
func (m *IncomingContext) GetOnBehalfOf()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("onBehalfOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetSourceParticipantId gets the sourceParticipantId property value. The ID of the participant that triggered the incoming call. Read-only.
// returns a *string when successful
func (m *IncomingContext) GetSourceParticipantId()(*string) {
    val, err := m.GetBackingStore().Get("sourceParticipantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTransferor gets the transferor property value. The identity that transferred the call.
// returns a IdentitySetable when successful
func (m *IncomingContext) GetTransferor()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("transferor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IncomingContext) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("observedParticipantId", m.GetObservedParticipantId())
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
        err := writer.WriteObjectValue("onBehalfOf", m.GetOnBehalfOf())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sourceParticipantId", m.GetSourceParticipantId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("transferor", m.GetTransferor())
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
func (m *IncomingContext) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *IncomingContext) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetObservedParticipantId sets the observedParticipantId property value. The ID of the participant that is under observation. Read-only.
func (m *IncomingContext) SetObservedParticipantId(value *string)() {
    err := m.GetBackingStore().Set("observedParticipantId", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IncomingContext) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOnBehalfOf sets the onBehalfOf property value. The identity that the call is happening on behalf of.
func (m *IncomingContext) SetOnBehalfOf(value IdentitySetable)() {
    err := m.GetBackingStore().Set("onBehalfOf", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceParticipantId sets the sourceParticipantId property value. The ID of the participant that triggered the incoming call. Read-only.
func (m *IncomingContext) SetSourceParticipantId(value *string)() {
    err := m.GetBackingStore().Set("sourceParticipantId", value)
    if err != nil {
        panic(err)
    }
}
// SetTransferor sets the transferor property value. The identity that transferred the call.
func (m *IncomingContext) SetTransferor(value IdentitySetable)() {
    err := m.GetBackingStore().Set("transferor", value)
    if err != nil {
        panic(err)
    }
}
type IncomingContextable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetObservedParticipantId()(*string)
    GetOdataType()(*string)
    GetOnBehalfOf()(IdentitySetable)
    GetSourceParticipantId()(*string)
    GetTransferor()(IdentitySetable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetObservedParticipantId(value *string)()
    SetOdataType(value *string)()
    SetOnBehalfOf(value IdentitySetable)()
    SetSourceParticipantId(value *string)()
    SetTransferor(value IdentitySetable)()
}

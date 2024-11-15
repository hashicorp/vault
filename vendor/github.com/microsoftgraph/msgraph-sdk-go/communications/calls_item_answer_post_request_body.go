package communications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CallsItemAnswerPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCallsItemAnswerPostRequestBody instantiates a new CallsItemAnswerPostRequestBody and sets the default values.
func NewCallsItemAnswerPostRequestBody()(*CallsItemAnswerPostRequestBody) {
    m := &CallsItemAnswerPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCallsItemAnswerPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallsItemAnswerPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallsItemAnswerPostRequestBody(), nil
}
// GetAcceptedModalities gets the acceptedModalities property value. The acceptedModalities property
// returns a []Modality when successful
func (m *CallsItemAnswerPostRequestBody) GetAcceptedModalities()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality) {
    val, err := m.GetBackingStore().Get("acceptedModalities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CallsItemAnswerPostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *CallsItemAnswerPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCallbackUri gets the callbackUri property value. The callbackUri property
// returns a *string when successful
func (m *CallsItemAnswerPostRequestBody) GetCallbackUri()(*string) {
    val, err := m.GetBackingStore().Get("callbackUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallOptions gets the callOptions property value. The callOptions property
// returns a IncomingCallOptionsable when successful
func (m *CallsItemAnswerPostRequestBody) GetCallOptions()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable) {
    val, err := m.GetBackingStore().Get("callOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CallsItemAnswerPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["acceptedModalities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseModality)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality))
                }
            }
            m.SetAcceptedModalities(res)
        }
        return nil
    }
    res["callbackUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallbackUri(val)
        }
        return nil
    }
    res["callOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIncomingCallOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallOptions(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable))
        }
        return nil
    }
    res["mediaConfig"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMediaConfigFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaConfig(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable))
        }
        return nil
    }
    res["participantCapacity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParticipantCapacity(val)
        }
        return nil
    }
    return res
}
// GetMediaConfig gets the mediaConfig property value. The mediaConfig property
// returns a MediaConfigable when successful
func (m *CallsItemAnswerPostRequestBody) GetMediaConfig()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable) {
    val, err := m.GetBackingStore().Get("mediaConfig")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable)
    }
    return nil
}
// GetParticipantCapacity gets the participantCapacity property value. The participantCapacity property
// returns a *int32 when successful
func (m *CallsItemAnswerPostRequestBody) GetParticipantCapacity()(*int32) {
    val, err := m.GetBackingStore().Get("participantCapacity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CallsItemAnswerPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAcceptedModalities() != nil {
        err := writer.WriteCollectionOfStringValues("acceptedModalities", iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SerializeModality(m.GetAcceptedModalities()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("callbackUri", m.GetCallbackUri())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callOptions", m.GetCallOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("mediaConfig", m.GetMediaConfig())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("participantCapacity", m.GetParticipantCapacity())
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
// SetAcceptedModalities sets the acceptedModalities property value. The acceptedModalities property
func (m *CallsItemAnswerPostRequestBody) SetAcceptedModalities(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality)() {
    err := m.GetBackingStore().Set("acceptedModalities", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CallsItemAnswerPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CallsItemAnswerPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCallbackUri sets the callbackUri property value. The callbackUri property
func (m *CallsItemAnswerPostRequestBody) SetCallbackUri(value *string)() {
    err := m.GetBackingStore().Set("callbackUri", value)
    if err != nil {
        panic(err)
    }
}
// SetCallOptions sets the callOptions property value. The callOptions property
func (m *CallsItemAnswerPostRequestBody) SetCallOptions(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable)() {
    err := m.GetBackingStore().Set("callOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaConfig sets the mediaConfig property value. The mediaConfig property
func (m *CallsItemAnswerPostRequestBody) SetMediaConfig(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable)() {
    err := m.GetBackingStore().Set("mediaConfig", value)
    if err != nil {
        panic(err)
    }
}
// SetParticipantCapacity sets the participantCapacity property value. The participantCapacity property
func (m *CallsItemAnswerPostRequestBody) SetParticipantCapacity(value *int32)() {
    err := m.GetBackingStore().Set("participantCapacity", value)
    if err != nil {
        panic(err)
    }
}
type CallsItemAnswerPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAcceptedModalities()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCallbackUri()(*string)
    GetCallOptions()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable)
    GetMediaConfig()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable)
    GetParticipantCapacity()(*int32)
    SetAcceptedModalities(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCallbackUri(value *string)()
    SetCallOptions(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable)()
    SetMediaConfig(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable)()
    SetParticipantCapacity(value *int32)()
}

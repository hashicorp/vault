package communications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CallsItemSendDtmfTonesPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCallsItemSendDtmfTonesPostRequestBody instantiates a new CallsItemSendDtmfTonesPostRequestBody and sets the default values.
func NewCallsItemSendDtmfTonesPostRequestBody()(*CallsItemSendDtmfTonesPostRequestBody) {
    m := &CallsItemSendDtmfTonesPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCallsItemSendDtmfTonesPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallsItemSendDtmfTonesPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallsItemSendDtmfTonesPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CallsItemSendDtmfTonesPostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *CallsItemSendDtmfTonesPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetClientContext gets the clientContext property value. The clientContext property
// returns a *string when successful
func (m *CallsItemSendDtmfTonesPostRequestBody) GetClientContext()(*string) {
    val, err := m.GetBackingStore().Get("clientContext")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDelayBetweenTonesMs gets the delayBetweenTonesMs property value. The delayBetweenTonesMs property
// returns a *int32 when successful
func (m *CallsItemSendDtmfTonesPostRequestBody) GetDelayBetweenTonesMs()(*int32) {
    val, err := m.GetBackingStore().Get("delayBetweenTonesMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CallsItemSendDtmfTonesPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["clientContext"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientContext(val)
        }
        return nil
    }
    res["delayBetweenTonesMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDelayBetweenTonesMs(val)
        }
        return nil
    }
    res["tones"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseTone)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone))
                }
            }
            m.SetTones(res)
        }
        return nil
    }
    return res
}
// GetTones gets the tones property value. The tones property
// returns a []Tone when successful
func (m *CallsItemSendDtmfTonesPostRequestBody) GetTones()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone) {
    val, err := m.GetBackingStore().Get("tones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CallsItemSendDtmfTonesPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("clientContext", m.GetClientContext())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("delayBetweenTonesMs", m.GetDelayBetweenTonesMs())
        if err != nil {
            return err
        }
    }
    if m.GetTones() != nil {
        err := writer.WriteCollectionOfStringValues("tones", iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SerializeTone(m.GetTones()))
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
func (m *CallsItemSendDtmfTonesPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CallsItemSendDtmfTonesPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetClientContext sets the clientContext property value. The clientContext property
func (m *CallsItemSendDtmfTonesPostRequestBody) SetClientContext(value *string)() {
    err := m.GetBackingStore().Set("clientContext", value)
    if err != nil {
        panic(err)
    }
}
// SetDelayBetweenTonesMs sets the delayBetweenTonesMs property value. The delayBetweenTonesMs property
func (m *CallsItemSendDtmfTonesPostRequestBody) SetDelayBetweenTonesMs(value *int32)() {
    err := m.GetBackingStore().Set("delayBetweenTonesMs", value)
    if err != nil {
        panic(err)
    }
}
// SetTones sets the tones property value. The tones property
func (m *CallsItemSendDtmfTonesPostRequestBody) SetTones(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone)() {
    err := m.GetBackingStore().Set("tones", value)
    if err != nil {
        panic(err)
    }
}
type CallsItemSendDtmfTonesPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetClientContext()(*string)
    GetDelayBetweenTonesMs()(*int32)
    GetTones()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetClientContext(value *string)()
    SetDelayBetweenTonesMs(value *int32)()
    SetTones(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Tone)()
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ItemTranslateExchangeIdsPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewItemTranslateExchangeIdsPostRequestBody instantiates a new ItemTranslateExchangeIdsPostRequestBody and sets the default values.
func NewItemTranslateExchangeIdsPostRequestBody()(*ItemTranslateExchangeIdsPostRequestBody) {
    m := &ItemTranslateExchangeIdsPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateItemTranslateExchangeIdsPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemTranslateExchangeIdsPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemTranslateExchangeIdsPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ItemTranslateExchangeIdsPostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *ItemTranslateExchangeIdsPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemTranslateExchangeIdsPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["InputIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetInputIds(res)
        }
        return nil
    }
    res["SourceIdType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseExchangeIdFormat)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceIdType(val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat))
        }
        return nil
    }
    res["TargetIdType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseExchangeIdFormat)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetIdType(val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat))
        }
        return nil
    }
    return res
}
// GetInputIds gets the InputIds property value. The InputIds property
// returns a []string when successful
func (m *ItemTranslateExchangeIdsPostRequestBody) GetInputIds()([]string) {
    val, err := m.GetBackingStore().Get("inputIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSourceIdType gets the SourceIdType property value. The SourceIdType property
// returns a *ExchangeIdFormat when successful
func (m *ItemTranslateExchangeIdsPostRequestBody) GetSourceIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat) {
    val, err := m.GetBackingStore().Get("sourceIdType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)
    }
    return nil
}
// GetTargetIdType gets the TargetIdType property value. The TargetIdType property
// returns a *ExchangeIdFormat when successful
func (m *ItemTranslateExchangeIdsPostRequestBody) GetTargetIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat) {
    val, err := m.GetBackingStore().Get("targetIdType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemTranslateExchangeIdsPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetInputIds() != nil {
        err := writer.WriteCollectionOfStringValues("InputIds", m.GetInputIds())
        if err != nil {
            return err
        }
    }
    if m.GetSourceIdType() != nil {
        cast := (*m.GetSourceIdType()).String()
        err := writer.WriteStringValue("SourceIdType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetTargetIdType() != nil {
        cast := (*m.GetTargetIdType()).String()
        err := writer.WriteStringValue("TargetIdType", &cast)
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
func (m *ItemTranslateExchangeIdsPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ItemTranslateExchangeIdsPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetInputIds sets the InputIds property value. The InputIds property
func (m *ItemTranslateExchangeIdsPostRequestBody) SetInputIds(value []string)() {
    err := m.GetBackingStore().Set("inputIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceIdType sets the SourceIdType property value. The SourceIdType property
func (m *ItemTranslateExchangeIdsPostRequestBody) SetSourceIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)() {
    err := m.GetBackingStore().Set("sourceIdType", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetIdType sets the TargetIdType property value. The TargetIdType property
func (m *ItemTranslateExchangeIdsPostRequestBody) SetTargetIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)() {
    err := m.GetBackingStore().Set("targetIdType", value)
    if err != nil {
        panic(err)
    }
}
type ItemTranslateExchangeIdsPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetInputIds()([]string)
    GetSourceIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)
    GetTargetIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetInputIds(value []string)()
    SetSourceIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)()
    SetTargetIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)()
}

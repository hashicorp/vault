package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AlterationResponse struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAlterationResponse instantiates a new AlterationResponse and sets the default values.
func NewAlterationResponse()(*AlterationResponse) {
    m := &AlterationResponse{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAlterationResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAlterationResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAlterationResponse(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AlterationResponse) GetAdditionalData()(map[string]any) {
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
func (m *AlterationResponse) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AlterationResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["originalQueryString"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOriginalQueryString(val)
        }
        return nil
    }
    res["queryAlteration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchAlterationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryAlteration(val.(SearchAlterationable))
        }
        return nil
    }
    res["queryAlterationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSearchAlterationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryAlterationType(val.(*SearchAlterationType))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AlterationResponse) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOriginalQueryString gets the originalQueryString property value. Defines the original user query string.
// returns a *string when successful
func (m *AlterationResponse) GetOriginalQueryString()(*string) {
    val, err := m.GetBackingStore().Get("originalQueryString")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryAlteration gets the queryAlteration property value. Defines the details of the alteration information for the spelling correction.
// returns a SearchAlterationable when successful
func (m *AlterationResponse) GetQueryAlteration()(SearchAlterationable) {
    val, err := m.GetBackingStore().Get("queryAlteration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchAlterationable)
    }
    return nil
}
// GetQueryAlterationType gets the queryAlterationType property value. Defines the type of the spelling correction. Possible values are: suggestion, modification.
// returns a *SearchAlterationType when successful
func (m *AlterationResponse) GetQueryAlterationType()(*SearchAlterationType) {
    val, err := m.GetBackingStore().Get("queryAlterationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SearchAlterationType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AlterationResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("originalQueryString", m.GetOriginalQueryString())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("queryAlteration", m.GetQueryAlteration())
        if err != nil {
            return err
        }
    }
    if m.GetQueryAlterationType() != nil {
        cast := (*m.GetQueryAlterationType()).String()
        err := writer.WriteStringValue("queryAlterationType", &cast)
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
func (m *AlterationResponse) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AlterationResponse) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AlterationResponse) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOriginalQueryString sets the originalQueryString property value. Defines the original user query string.
func (m *AlterationResponse) SetOriginalQueryString(value *string)() {
    err := m.GetBackingStore().Set("originalQueryString", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryAlteration sets the queryAlteration property value. Defines the details of the alteration information for the spelling correction.
func (m *AlterationResponse) SetQueryAlteration(value SearchAlterationable)() {
    err := m.GetBackingStore().Set("queryAlteration", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryAlterationType sets the queryAlterationType property value. Defines the type of the spelling correction. Possible values are: suggestion, modification.
func (m *AlterationResponse) SetQueryAlterationType(value *SearchAlterationType)() {
    err := m.GetBackingStore().Set("queryAlterationType", value)
    if err != nil {
        panic(err)
    }
}
type AlterationResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetOriginalQueryString()(*string)
    GetQueryAlteration()(SearchAlterationable)
    GetQueryAlterationType()(*SearchAlterationType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetOriginalQueryString(value *string)()
    SetQueryAlteration(value SearchAlterationable)()
    SetQueryAlterationType(value *SearchAlterationType)()
}

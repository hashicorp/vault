package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type BucketAggregationDefinition struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewBucketAggregationDefinition instantiates a new BucketAggregationDefinition and sets the default values.
func NewBucketAggregationDefinition()(*BucketAggregationDefinition) {
    m := &BucketAggregationDefinition{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateBucketAggregationDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBucketAggregationDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBucketAggregationDefinition(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *BucketAggregationDefinition) GetAdditionalData()(map[string]any) {
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
func (m *BucketAggregationDefinition) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BucketAggregationDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isDescending"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDescending(val)
        }
        return nil
    }
    res["minimumCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumCount(val)
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
    res["prefixFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrefixFilter(val)
        }
        return nil
    }
    res["ranges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBucketAggregationRangeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BucketAggregationRangeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BucketAggregationRangeable)
                }
            }
            m.SetRanges(res)
        }
        return nil
    }
    res["sortBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBucketAggregationSortProperty)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSortBy(val.(*BucketAggregationSortProperty))
        }
        return nil
    }
    return res
}
// GetIsDescending gets the isDescending property value. True to specify the sort order as descending. The default is false, with the sort order as ascending. Optional.
// returns a *bool when successful
func (m *BucketAggregationDefinition) GetIsDescending()(*bool) {
    val, err := m.GetBackingStore().Get("isDescending")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMinimumCount gets the minimumCount property value. The minimum number of items that should be present in the aggregation to be returned in a bucket. Optional.
// returns a *int32 when successful
func (m *BucketAggregationDefinition) GetMinimumCount()(*int32) {
    val, err := m.GetBackingStore().Get("minimumCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *BucketAggregationDefinition) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrefixFilter gets the prefixFilter property value. A filter to define a matching criteria. The key should start with the specified prefix to be returned in the response. Optional.
// returns a *string when successful
func (m *BucketAggregationDefinition) GetPrefixFilter()(*string) {
    val, err := m.GetBackingStore().Get("prefixFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRanges gets the ranges property value. Specifies the manual ranges to compute the aggregations. This is only valid for nonstring refiners of date or numeric type. Optional.
// returns a []BucketAggregationRangeable when successful
func (m *BucketAggregationDefinition) GetRanges()([]BucketAggregationRangeable) {
    val, err := m.GetBackingStore().Get("ranges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BucketAggregationRangeable)
    }
    return nil
}
// GetSortBy gets the sortBy property value. The sortBy property
// returns a *BucketAggregationSortProperty when successful
func (m *BucketAggregationDefinition) GetSortBy()(*BucketAggregationSortProperty) {
    val, err := m.GetBackingStore().Get("sortBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BucketAggregationSortProperty)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BucketAggregationDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isDescending", m.GetIsDescending())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("minimumCount", m.GetMinimumCount())
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
        err := writer.WriteStringValue("prefixFilter", m.GetPrefixFilter())
        if err != nil {
            return err
        }
    }
    if m.GetRanges() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRanges()))
        for i, v := range m.GetRanges() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("ranges", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSortBy() != nil {
        cast := (*m.GetSortBy()).String()
        err := writer.WriteStringValue("sortBy", &cast)
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
func (m *BucketAggregationDefinition) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *BucketAggregationDefinition) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsDescending sets the isDescending property value. True to specify the sort order as descending. The default is false, with the sort order as ascending. Optional.
func (m *BucketAggregationDefinition) SetIsDescending(value *bool)() {
    err := m.GetBackingStore().Set("isDescending", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumCount sets the minimumCount property value. The minimum number of items that should be present in the aggregation to be returned in a bucket. Optional.
func (m *BucketAggregationDefinition) SetMinimumCount(value *int32)() {
    err := m.GetBackingStore().Set("minimumCount", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BucketAggregationDefinition) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrefixFilter sets the prefixFilter property value. A filter to define a matching criteria. The key should start with the specified prefix to be returned in the response. Optional.
func (m *BucketAggregationDefinition) SetPrefixFilter(value *string)() {
    err := m.GetBackingStore().Set("prefixFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetRanges sets the ranges property value. Specifies the manual ranges to compute the aggregations. This is only valid for nonstring refiners of date or numeric type. Optional.
func (m *BucketAggregationDefinition) SetRanges(value []BucketAggregationRangeable)() {
    err := m.GetBackingStore().Set("ranges", value)
    if err != nil {
        panic(err)
    }
}
// SetSortBy sets the sortBy property value. The sortBy property
func (m *BucketAggregationDefinition) SetSortBy(value *BucketAggregationSortProperty)() {
    err := m.GetBackingStore().Set("sortBy", value)
    if err != nil {
        panic(err)
    }
}
type BucketAggregationDefinitionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsDescending()(*bool)
    GetMinimumCount()(*int32)
    GetOdataType()(*string)
    GetPrefixFilter()(*string)
    GetRanges()([]BucketAggregationRangeable)
    GetSortBy()(*BucketAggregationSortProperty)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsDescending(value *bool)()
    SetMinimumCount(value *int32)()
    SetOdataType(value *string)()
    SetPrefixFilter(value *string)()
    SetRanges(value []BucketAggregationRangeable)()
    SetSortBy(value *BucketAggregationSortProperty)()
}

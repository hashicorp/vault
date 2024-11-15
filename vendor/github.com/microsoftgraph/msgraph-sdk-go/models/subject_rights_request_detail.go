package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SubjectRightsRequestDetail struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSubjectRightsRequestDetail instantiates a new SubjectRightsRequestDetail and sets the default values.
func NewSubjectRightsRequestDetail()(*SubjectRightsRequestDetail) {
    m := &SubjectRightsRequestDetail{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSubjectRightsRequestDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSubjectRightsRequestDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSubjectRightsRequestDetail(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SubjectRightsRequestDetail) GetAdditionalData()(map[string]any) {
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
func (m *SubjectRightsRequestDetail) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetExcludedItemCount gets the excludedItemCount property value. Count of items that are excluded from the request.
// returns a *int64 when successful
func (m *SubjectRightsRequestDetail) GetExcludedItemCount()(*int64) {
    val, err := m.GetBackingStore().Get("excludedItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SubjectRightsRequestDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["excludedItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExcludedItemCount(val)
        }
        return nil
    }
    res["insightCounts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValuePairable)
                }
            }
            m.SetInsightCounts(res)
        }
        return nil
    }
    res["itemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetItemCount(val)
        }
        return nil
    }
    res["itemNeedReview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetItemNeedReview(val)
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
    res["productItemCounts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValuePairable)
                }
            }
            m.SetProductItemCounts(res)
        }
        return nil
    }
    res["signedOffItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignedOffItemCount(val)
        }
        return nil
    }
    res["totalItemSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalItemSize(val)
        }
        return nil
    }
    return res
}
// GetInsightCounts gets the insightCounts property value. Count of items per insight.
// returns a []KeyValuePairable when successful
func (m *SubjectRightsRequestDetail) GetInsightCounts()([]KeyValuePairable) {
    val, err := m.GetBackingStore().Get("insightCounts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValuePairable)
    }
    return nil
}
// GetItemCount gets the itemCount property value. Count of items found.
// returns a *int64 when successful
func (m *SubjectRightsRequestDetail) GetItemCount()(*int64) {
    val, err := m.GetBackingStore().Get("itemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetItemNeedReview gets the itemNeedReview property value. Count of item that need review.
// returns a *int64 when successful
func (m *SubjectRightsRequestDetail) GetItemNeedReview()(*int64) {
    val, err := m.GetBackingStore().Get("itemNeedReview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SubjectRightsRequestDetail) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductItemCounts gets the productItemCounts property value. Count of items per product, such as Exchange, SharePoint, OneDrive, and Teams.
// returns a []KeyValuePairable when successful
func (m *SubjectRightsRequestDetail) GetProductItemCounts()([]KeyValuePairable) {
    val, err := m.GetBackingStore().Get("productItemCounts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValuePairable)
    }
    return nil
}
// GetSignedOffItemCount gets the signedOffItemCount property value. Count of items signed off by the administrator.
// returns a *int64 when successful
func (m *SubjectRightsRequestDetail) GetSignedOffItemCount()(*int64) {
    val, err := m.GetBackingStore().Get("signedOffItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetTotalItemSize gets the totalItemSize property value. Total item size in bytes.
// returns a *int64 when successful
func (m *SubjectRightsRequestDetail) GetTotalItemSize()(*int64) {
    val, err := m.GetBackingStore().Get("totalItemSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SubjectRightsRequestDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt64Value("excludedItemCount", m.GetExcludedItemCount())
        if err != nil {
            return err
        }
    }
    if m.GetInsightCounts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInsightCounts()))
        for i, v := range m.GetInsightCounts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("insightCounts", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("itemCount", m.GetItemCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("itemNeedReview", m.GetItemNeedReview())
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
    if m.GetProductItemCounts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProductItemCounts()))
        for i, v := range m.GetProductItemCounts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("productItemCounts", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("signedOffItemCount", m.GetSignedOffItemCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("totalItemSize", m.GetTotalItemSize())
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
func (m *SubjectRightsRequestDetail) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SubjectRightsRequestDetail) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetExcludedItemCount sets the excludedItemCount property value. Count of items that are excluded from the request.
func (m *SubjectRightsRequestDetail) SetExcludedItemCount(value *int64)() {
    err := m.GetBackingStore().Set("excludedItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInsightCounts sets the insightCounts property value. Count of items per insight.
func (m *SubjectRightsRequestDetail) SetInsightCounts(value []KeyValuePairable)() {
    err := m.GetBackingStore().Set("insightCounts", value)
    if err != nil {
        panic(err)
    }
}
// SetItemCount sets the itemCount property value. Count of items found.
func (m *SubjectRightsRequestDetail) SetItemCount(value *int64)() {
    err := m.GetBackingStore().Set("itemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetItemNeedReview sets the itemNeedReview property value. Count of item that need review.
func (m *SubjectRightsRequestDetail) SetItemNeedReview(value *int64)() {
    err := m.GetBackingStore().Set("itemNeedReview", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SubjectRightsRequestDetail) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetProductItemCounts sets the productItemCounts property value. Count of items per product, such as Exchange, SharePoint, OneDrive, and Teams.
func (m *SubjectRightsRequestDetail) SetProductItemCounts(value []KeyValuePairable)() {
    err := m.GetBackingStore().Set("productItemCounts", value)
    if err != nil {
        panic(err)
    }
}
// SetSignedOffItemCount sets the signedOffItemCount property value. Count of items signed off by the administrator.
func (m *SubjectRightsRequestDetail) SetSignedOffItemCount(value *int64)() {
    err := m.GetBackingStore().Set("signedOffItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalItemSize sets the totalItemSize property value. Total item size in bytes.
func (m *SubjectRightsRequestDetail) SetTotalItemSize(value *int64)() {
    err := m.GetBackingStore().Set("totalItemSize", value)
    if err != nil {
        panic(err)
    }
}
type SubjectRightsRequestDetailable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetExcludedItemCount()(*int64)
    GetInsightCounts()([]KeyValuePairable)
    GetItemCount()(*int64)
    GetItemNeedReview()(*int64)
    GetOdataType()(*string)
    GetProductItemCounts()([]KeyValuePairable)
    GetSignedOffItemCount()(*int64)
    GetTotalItemSize()(*int64)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetExcludedItemCount(value *int64)()
    SetInsightCounts(value []KeyValuePairable)()
    SetItemCount(value *int64)()
    SetItemNeedReview(value *int64)()
    SetOdataType(value *string)()
    SetProductItemCounts(value []KeyValuePairable)()
    SetSignedOffItemCount(value *int64)()
    SetTotalItemSize(value *int64)()
}

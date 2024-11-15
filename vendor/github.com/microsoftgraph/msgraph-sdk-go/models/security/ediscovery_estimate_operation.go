package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryEstimateOperation struct {
    CaseOperation
}
// NewEdiscoveryEstimateOperation instantiates a new EdiscoveryEstimateOperation and sets the default values.
func NewEdiscoveryEstimateOperation()(*EdiscoveryEstimateOperation) {
    m := &EdiscoveryEstimateOperation{
        CaseOperation: *NewCaseOperation(),
    }
    return m
}
// CreateEdiscoveryEstimateOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryEstimateOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryEstimateOperation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryEstimateOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CaseOperation.GetFieldDeserializers()
    res["indexedItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndexedItemCount(val)
        }
        return nil
    }
    res["indexedItemsSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndexedItemsSize(val)
        }
        return nil
    }
    res["mailboxCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailboxCount(val)
        }
        return nil
    }
    res["search"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoverySearchFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearch(val.(EdiscoverySearchable))
        }
        return nil
    }
    res["siteCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteCount(val)
        }
        return nil
    }
    res["unindexedItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnindexedItemCount(val)
        }
        return nil
    }
    res["unindexedItemsSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnindexedItemsSize(val)
        }
        return nil
    }
    return res
}
// GetIndexedItemCount gets the indexedItemCount property value. The estimated count of items for the search that matched the content query.
// returns a *int64 when successful
func (m *EdiscoveryEstimateOperation) GetIndexedItemCount()(*int64) {
    val, err := m.GetBackingStore().Get("indexedItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetIndexedItemsSize gets the indexedItemsSize property value. The estimated size of items for the search that matched the content query.
// returns a *int64 when successful
func (m *EdiscoveryEstimateOperation) GetIndexedItemsSize()(*int64) {
    val, err := m.GetBackingStore().Get("indexedItemsSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetMailboxCount gets the mailboxCount property value. The number of mailboxes that had search hits.
// returns a *int32 when successful
func (m *EdiscoveryEstimateOperation) GetMailboxCount()(*int32) {
    val, err := m.GetBackingStore().Get("mailboxCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSearch gets the search property value. eDiscovery search.
// returns a EdiscoverySearchable when successful
func (m *EdiscoveryEstimateOperation) GetSearch()(EdiscoverySearchable) {
    val, err := m.GetBackingStore().Get("search")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoverySearchable)
    }
    return nil
}
// GetSiteCount gets the siteCount property value. The number of mailboxes that had search hits.
// returns a *int32 when successful
func (m *EdiscoveryEstimateOperation) GetSiteCount()(*int32) {
    val, err := m.GetBackingStore().Get("siteCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnindexedItemCount gets the unindexedItemCount property value. The estimated count of unindexed items for the collection.
// returns a *int64 when successful
func (m *EdiscoveryEstimateOperation) GetUnindexedItemCount()(*int64) {
    val, err := m.GetBackingStore().Get("unindexedItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUnindexedItemsSize gets the unindexedItemsSize property value. The estimated size of unindexed items for the collection.
// returns a *int64 when successful
func (m *EdiscoveryEstimateOperation) GetUnindexedItemsSize()(*int64) {
    val, err := m.GetBackingStore().Get("unindexedItemsSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryEstimateOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CaseOperation.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt64Value("indexedItemCount", m.GetIndexedItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("indexedItemsSize", m.GetIndexedItemsSize())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("mailboxCount", m.GetMailboxCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("search", m.GetSearch())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("siteCount", m.GetSiteCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("unindexedItemCount", m.GetUnindexedItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("unindexedItemsSize", m.GetUnindexedItemsSize())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIndexedItemCount sets the indexedItemCount property value. The estimated count of items for the search that matched the content query.
func (m *EdiscoveryEstimateOperation) SetIndexedItemCount(value *int64)() {
    err := m.GetBackingStore().Set("indexedItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIndexedItemsSize sets the indexedItemsSize property value. The estimated size of items for the search that matched the content query.
func (m *EdiscoveryEstimateOperation) SetIndexedItemsSize(value *int64)() {
    err := m.GetBackingStore().Set("indexedItemsSize", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxCount sets the mailboxCount property value. The number of mailboxes that had search hits.
func (m *EdiscoveryEstimateOperation) SetMailboxCount(value *int32)() {
    err := m.GetBackingStore().Set("mailboxCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSearch sets the search property value. eDiscovery search.
func (m *EdiscoveryEstimateOperation) SetSearch(value EdiscoverySearchable)() {
    err := m.GetBackingStore().Set("search", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteCount sets the siteCount property value. The number of mailboxes that had search hits.
func (m *EdiscoveryEstimateOperation) SetSiteCount(value *int32)() {
    err := m.GetBackingStore().Set("siteCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnindexedItemCount sets the unindexedItemCount property value. The estimated count of unindexed items for the collection.
func (m *EdiscoveryEstimateOperation) SetUnindexedItemCount(value *int64)() {
    err := m.GetBackingStore().Set("unindexedItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnindexedItemsSize sets the unindexedItemsSize property value. The estimated size of unindexed items for the collection.
func (m *EdiscoveryEstimateOperation) SetUnindexedItemsSize(value *int64)() {
    err := m.GetBackingStore().Set("unindexedItemsSize", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryEstimateOperationable interface {
    CaseOperationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIndexedItemCount()(*int64)
    GetIndexedItemsSize()(*int64)
    GetMailboxCount()(*int32)
    GetSearch()(EdiscoverySearchable)
    GetSiteCount()(*int32)
    GetUnindexedItemCount()(*int64)
    GetUnindexedItemsSize()(*int64)
    SetIndexedItemCount(value *int64)()
    SetIndexedItemsSize(value *int64)()
    SetMailboxCount(value *int32)()
    SetSearch(value EdiscoverySearchable)()
    SetSiteCount(value *int32)()
    SetUnindexedItemCount(value *int64)()
    SetUnindexedItemsSize(value *int64)()
}

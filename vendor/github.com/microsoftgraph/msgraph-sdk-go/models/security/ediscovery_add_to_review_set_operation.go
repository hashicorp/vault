package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryAddToReviewSetOperation struct {
    CaseOperation
}
// NewEdiscoveryAddToReviewSetOperation instantiates a new EdiscoveryAddToReviewSetOperation and sets the default values.
func NewEdiscoveryAddToReviewSetOperation()(*EdiscoveryAddToReviewSetOperation) {
    m := &EdiscoveryAddToReviewSetOperation{
        CaseOperation: *NewCaseOperation(),
    }
    return m
}
// CreateEdiscoveryAddToReviewSetOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryAddToReviewSetOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryAddToReviewSetOperation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryAddToReviewSetOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CaseOperation.GetFieldDeserializers()
    res["reviewSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryReviewSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewSet(val.(EdiscoveryReviewSetable))
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
    return res
}
// GetReviewSet gets the reviewSet property value. eDiscovery review set to which items matching source collection query gets added.
// returns a EdiscoveryReviewSetable when successful
func (m *EdiscoveryAddToReviewSetOperation) GetReviewSet()(EdiscoveryReviewSetable) {
    val, err := m.GetBackingStore().Get("reviewSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryReviewSetable)
    }
    return nil
}
// GetSearch gets the search property value. eDiscovery search that gets added to review set.
// returns a EdiscoverySearchable when successful
func (m *EdiscoveryAddToReviewSetOperation) GetSearch()(EdiscoverySearchable) {
    val, err := m.GetBackingStore().Get("search")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoverySearchable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryAddToReviewSetOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CaseOperation.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("reviewSet", m.GetReviewSet())
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
    return nil
}
// SetReviewSet sets the reviewSet property value. eDiscovery review set to which items matching source collection query gets added.
func (m *EdiscoveryAddToReviewSetOperation) SetReviewSet(value EdiscoveryReviewSetable)() {
    err := m.GetBackingStore().Set("reviewSet", value)
    if err != nil {
        panic(err)
    }
}
// SetSearch sets the search property value. eDiscovery search that gets added to review set.
func (m *EdiscoveryAddToReviewSetOperation) SetSearch(value EdiscoverySearchable)() {
    err := m.GetBackingStore().Set("search", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryAddToReviewSetOperationable interface {
    CaseOperationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetReviewSet()(EdiscoveryReviewSetable)
    GetSearch()(EdiscoverySearchable)
    SetReviewSet(value EdiscoveryReviewSetable)()
    SetSearch(value EdiscoverySearchable)()
}

package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryReviewSetQuery struct {
    Search
}
// NewEdiscoveryReviewSetQuery instantiates a new EdiscoveryReviewSetQuery and sets the default values.
func NewEdiscoveryReviewSetQuery()(*EdiscoveryReviewSetQuery) {
    m := &EdiscoveryReviewSetQuery{
        Search: *NewSearch(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryReviewSetQuery"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoveryReviewSetQueryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryReviewSetQueryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryReviewSetQuery(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryReviewSetQuery) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Search.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *EdiscoveryReviewSetQuery) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Search.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type EdiscoveryReviewSetQueryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Searchable
}

package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryReviewSet struct {
    DataSet
}
// NewEdiscoveryReviewSet instantiates a new EdiscoveryReviewSet and sets the default values.
func NewEdiscoveryReviewSet()(*EdiscoveryReviewSet) {
    m := &EdiscoveryReviewSet{
        DataSet: *NewDataSet(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryReviewSet"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoveryReviewSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryReviewSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryReviewSet(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryReviewSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DataSet.GetFieldDeserializers()
    res["queries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryReviewSetQueryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryReviewSetQueryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryReviewSetQueryable)
                }
            }
            m.SetQueries(res)
        }
        return nil
    }
    return res
}
// GetQueries gets the queries property value. Represents queries within the review set.
// returns a []EdiscoveryReviewSetQueryable when successful
func (m *EdiscoveryReviewSet) GetQueries()([]EdiscoveryReviewSetQueryable) {
    val, err := m.GetBackingStore().Get("queries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryReviewSetQueryable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryReviewSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DataSet.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetQueries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetQueries()))
        for i, v := range m.GetQueries() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("queries", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetQueries sets the queries property value. Represents queries within the review set.
func (m *EdiscoveryReviewSet) SetQueries(value []EdiscoveryReviewSetQueryable)() {
    err := m.GetBackingStore().Set("queries", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryReviewSetable interface {
    DataSetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetQueries()([]EdiscoveryReviewSetQueryable)
    SetQueries(value []EdiscoveryReviewSetQueryable)()
}

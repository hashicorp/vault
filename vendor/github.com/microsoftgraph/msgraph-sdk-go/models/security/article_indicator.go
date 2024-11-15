package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ArticleIndicator struct {
    Indicator
}
// NewArticleIndicator instantiates a new ArticleIndicator and sets the default values.
func NewArticleIndicator()(*ArticleIndicator) {
    m := &ArticleIndicator{
        Indicator: *NewIndicator(),
    }
    odataTypeValue := "#microsoft.graph.security.articleIndicator"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateArticleIndicatorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateArticleIndicatorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewArticleIndicator(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ArticleIndicator) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Indicator.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *ArticleIndicator) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Indicator.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type ArticleIndicatorable interface {
    Indicatorable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

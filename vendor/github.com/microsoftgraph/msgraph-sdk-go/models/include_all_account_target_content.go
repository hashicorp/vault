package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IncludeAllAccountTargetContent struct {
    AccountTargetContent
}
// NewIncludeAllAccountTargetContent instantiates a new IncludeAllAccountTargetContent and sets the default values.
func NewIncludeAllAccountTargetContent()(*IncludeAllAccountTargetContent) {
    m := &IncludeAllAccountTargetContent{
        AccountTargetContent: *NewAccountTargetContent(),
    }
    odataTypeValue := "#microsoft.graph.includeAllAccountTargetContent"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIncludeAllAccountTargetContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIncludeAllAccountTargetContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIncludeAllAccountTargetContent(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IncludeAllAccountTargetContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccountTargetContent.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *IncludeAllAccountTargetContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccountTargetContent.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type IncludeAllAccountTargetContentable interface {
    AccountTargetContentable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

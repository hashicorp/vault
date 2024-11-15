package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IdentityCustomUserFlowAttribute struct {
    IdentityUserFlowAttribute
}
// NewIdentityCustomUserFlowAttribute instantiates a new IdentityCustomUserFlowAttribute and sets the default values.
func NewIdentityCustomUserFlowAttribute()(*IdentityCustomUserFlowAttribute) {
    m := &IdentityCustomUserFlowAttribute{
        IdentityUserFlowAttribute: *NewIdentityUserFlowAttribute(),
    }
    odataTypeValue := "#microsoft.graph.identityCustomUserFlowAttribute"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIdentityCustomUserFlowAttributeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityCustomUserFlowAttributeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityCustomUserFlowAttribute(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IdentityCustomUserFlowAttribute) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentityUserFlowAttribute.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *IdentityCustomUserFlowAttribute) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentityUserFlowAttribute.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type IdentityCustomUserFlowAttributeable interface {
    IdentityUserFlowAttributeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

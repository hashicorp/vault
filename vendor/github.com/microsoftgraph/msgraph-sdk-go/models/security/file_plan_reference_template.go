package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FilePlanReferenceTemplate struct {
    FilePlanDescriptorTemplate
}
// NewFilePlanReferenceTemplate instantiates a new FilePlanReferenceTemplate and sets the default values.
func NewFilePlanReferenceTemplate()(*FilePlanReferenceTemplate) {
    m := &FilePlanReferenceTemplate{
        FilePlanDescriptorTemplate: *NewFilePlanDescriptorTemplate(),
    }
    return m
}
// CreateFilePlanReferenceTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFilePlanReferenceTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFilePlanReferenceTemplate(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FilePlanReferenceTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.FilePlanDescriptorTemplate.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *FilePlanReferenceTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.FilePlanDescriptorTemplate.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type FilePlanReferenceTemplateable interface {
    FilePlanDescriptorTemplateable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

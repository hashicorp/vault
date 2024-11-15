package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FilePlanAuthority struct {
    FilePlanDescriptorBase
}
// NewFilePlanAuthority instantiates a new FilePlanAuthority and sets the default values.
func NewFilePlanAuthority()(*FilePlanAuthority) {
    m := &FilePlanAuthority{
        FilePlanDescriptorBase: *NewFilePlanDescriptorBase(),
    }
    return m
}
// CreateFilePlanAuthorityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFilePlanAuthorityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFilePlanAuthority(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FilePlanAuthority) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.FilePlanDescriptorBase.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *FilePlanAuthority) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.FilePlanDescriptorBase.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type FilePlanAuthorityable interface {
    FilePlanDescriptorBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

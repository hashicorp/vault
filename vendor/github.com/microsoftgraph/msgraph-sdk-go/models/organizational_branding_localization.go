package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OrganizationalBrandingLocalization struct {
    OrganizationalBrandingProperties
}
// NewOrganizationalBrandingLocalization instantiates a new OrganizationalBrandingLocalization and sets the default values.
func NewOrganizationalBrandingLocalization()(*OrganizationalBrandingLocalization) {
    m := &OrganizationalBrandingLocalization{
        OrganizationalBrandingProperties: *NewOrganizationalBrandingProperties(),
    }
    odataTypeValue := "#microsoft.graph.organizationalBrandingLocalization"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOrganizationalBrandingLocalizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOrganizationalBrandingLocalizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOrganizationalBrandingLocalization(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OrganizationalBrandingLocalization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OrganizationalBrandingProperties.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *OrganizationalBrandingLocalization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OrganizationalBrandingProperties.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type OrganizationalBrandingLocalizationable interface {
    OrganizationalBrandingPropertiesable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageUserDirectoryAttributeStore struct {
    AccessPackageResourceAttributeDestination
}
// NewAccessPackageUserDirectoryAttributeStore instantiates a new AccessPackageUserDirectoryAttributeStore and sets the default values.
func NewAccessPackageUserDirectoryAttributeStore()(*AccessPackageUserDirectoryAttributeStore) {
    m := &AccessPackageUserDirectoryAttributeStore{
        AccessPackageResourceAttributeDestination: *NewAccessPackageResourceAttributeDestination(),
    }
    odataTypeValue := "#microsoft.graph.accessPackageUserDirectoryAttributeStore"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessPackageUserDirectoryAttributeStoreFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageUserDirectoryAttributeStoreFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageUserDirectoryAttributeStore(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageUserDirectoryAttributeStore) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessPackageResourceAttributeDestination.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *AccessPackageUserDirectoryAttributeStore) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessPackageResourceAttributeDestination.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type AccessPackageUserDirectoryAttributeStoreable interface {
    AccessPackageResourceAttributeDestinationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

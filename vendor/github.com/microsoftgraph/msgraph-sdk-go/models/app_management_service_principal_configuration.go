package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AppManagementServicePrincipalConfiguration struct {
    AppManagementConfiguration
}
// NewAppManagementServicePrincipalConfiguration instantiates a new AppManagementServicePrincipalConfiguration and sets the default values.
func NewAppManagementServicePrincipalConfiguration()(*AppManagementServicePrincipalConfiguration) {
    m := &AppManagementServicePrincipalConfiguration{
        AppManagementConfiguration: *NewAppManagementConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.appManagementServicePrincipalConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAppManagementServicePrincipalConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppManagementServicePrincipalConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppManagementServicePrincipalConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AppManagementServicePrincipalConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AppManagementConfiguration.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *AppManagementServicePrincipalConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AppManagementConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type AppManagementServicePrincipalConfigurationable interface {
    AppManagementConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

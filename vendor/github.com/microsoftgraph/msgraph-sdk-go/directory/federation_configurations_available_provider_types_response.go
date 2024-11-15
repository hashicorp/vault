package directory

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use FederationConfigurationsAvailableProviderTypesGetResponseable instead.
type FederationConfigurationsAvailableProviderTypesResponse struct {
    FederationConfigurationsAvailableProviderTypesGetResponse
}
// NewFederationConfigurationsAvailableProviderTypesResponse instantiates a new FederationConfigurationsAvailableProviderTypesResponse and sets the default values.
func NewFederationConfigurationsAvailableProviderTypesResponse()(*FederationConfigurationsAvailableProviderTypesResponse) {
    m := &FederationConfigurationsAvailableProviderTypesResponse{
        FederationConfigurationsAvailableProviderTypesGetResponse: *NewFederationConfigurationsAvailableProviderTypesGetResponse(),
    }
    return m
}
// CreateFederationConfigurationsAvailableProviderTypesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFederationConfigurationsAvailableProviderTypesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFederationConfigurationsAvailableProviderTypesResponse(), nil
}
// Deprecated: This class is obsolete. Use FederationConfigurationsAvailableProviderTypesGetResponseable instead.
type FederationConfigurationsAvailableProviderTypesResponseable interface {
    FederationConfigurationsAvailableProviderTypesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

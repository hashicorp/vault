package identityproviders

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use AvailableProviderTypesGetResponseable instead.
type AvailableProviderTypesResponse struct {
    AvailableProviderTypesGetResponse
}
// NewAvailableProviderTypesResponse instantiates a new AvailableProviderTypesResponse and sets the default values.
func NewAvailableProviderTypesResponse()(*AvailableProviderTypesResponse) {
    m := &AvailableProviderTypesResponse{
        AvailableProviderTypesGetResponse: *NewAvailableProviderTypesGetResponse(),
    }
    return m
}
// CreateAvailableProviderTypesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAvailableProviderTypesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAvailableProviderTypesResponse(), nil
}
// Deprecated: This class is obsolete. Use AvailableProviderTypesGetResponseable instead.
type AvailableProviderTypesResponseable interface {
    AvailableProviderTypesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

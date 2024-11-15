package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemGetManagedAppPoliciesGetResponseable instead.
type ItemGetManagedAppPoliciesResponse struct {
    ItemGetManagedAppPoliciesGetResponse
}
// NewItemGetManagedAppPoliciesResponse instantiates a new ItemGetManagedAppPoliciesResponse and sets the default values.
func NewItemGetManagedAppPoliciesResponse()(*ItemGetManagedAppPoliciesResponse) {
    m := &ItemGetManagedAppPoliciesResponse{
        ItemGetManagedAppPoliciesGetResponse: *NewItemGetManagedAppPoliciesGetResponse(),
    }
    return m
}
// CreateItemGetManagedAppPoliciesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemGetManagedAppPoliciesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemGetManagedAppPoliciesResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemGetManagedAppPoliciesGetResponseable instead.
type ItemGetManagedAppPoliciesResponseable interface {
    ItemGetManagedAppPoliciesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

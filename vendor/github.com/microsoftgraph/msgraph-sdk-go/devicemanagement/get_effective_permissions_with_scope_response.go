package devicemanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use GetEffectivePermissionsWithScopeGetResponseable instead.
type GetEffectivePermissionsWithScopeResponse struct {
    GetEffectivePermissionsWithScopeGetResponse
}
// NewGetEffectivePermissionsWithScopeResponse instantiates a new GetEffectivePermissionsWithScopeResponse and sets the default values.
func NewGetEffectivePermissionsWithScopeResponse()(*GetEffectivePermissionsWithScopeResponse) {
    m := &GetEffectivePermissionsWithScopeResponse{
        GetEffectivePermissionsWithScopeGetResponse: *NewGetEffectivePermissionsWithScopeGetResponse(),
    }
    return m
}
// CreateGetEffectivePermissionsWithScopeResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGetEffectivePermissionsWithScopeResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGetEffectivePermissionsWithScopeResponse(), nil
}
// Deprecated: This class is obsolete. Use GetEffectivePermissionsWithScopeGetResponseable instead.
type GetEffectivePermissionsWithScopeResponseable interface {
    GetEffectivePermissionsWithScopeGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

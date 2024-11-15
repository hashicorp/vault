package solutions

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use VirtualEventsWebinarsGetByUserRoleWithRoleGetResponseable instead.
type VirtualEventsWebinarsGetByUserRoleWithRoleResponse struct {
    VirtualEventsWebinarsGetByUserRoleWithRoleGetResponse
}
// NewVirtualEventsWebinarsGetByUserRoleWithRoleResponse instantiates a new VirtualEventsWebinarsGetByUserRoleWithRoleResponse and sets the default values.
func NewVirtualEventsWebinarsGetByUserRoleWithRoleResponse()(*VirtualEventsWebinarsGetByUserRoleWithRoleResponse) {
    m := &VirtualEventsWebinarsGetByUserRoleWithRoleResponse{
        VirtualEventsWebinarsGetByUserRoleWithRoleGetResponse: *NewVirtualEventsWebinarsGetByUserRoleWithRoleGetResponse(),
    }
    return m
}
// CreateVirtualEventsWebinarsGetByUserRoleWithRoleResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventsWebinarsGetByUserRoleWithRoleResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventsWebinarsGetByUserRoleWithRoleResponse(), nil
}
// Deprecated: This class is obsolete. Use VirtualEventsWebinarsGetByUserRoleWithRoleGetResponseable instead.
type VirtualEventsWebinarsGetByUserRoleWithRoleResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VirtualEventsWebinarsGetByUserRoleWithRoleGetResponseable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarAllowedCalendarSharingRolesWithUserGetResponseable instead.
type ItemCalendarAllowedCalendarSharingRolesWithUserResponse struct {
    ItemCalendarAllowedCalendarSharingRolesWithUserGetResponse
}
// NewItemCalendarAllowedCalendarSharingRolesWithUserResponse instantiates a new ItemCalendarAllowedCalendarSharingRolesWithUserResponse and sets the default values.
func NewItemCalendarAllowedCalendarSharingRolesWithUserResponse()(*ItemCalendarAllowedCalendarSharingRolesWithUserResponse) {
    m := &ItemCalendarAllowedCalendarSharingRolesWithUserResponse{
        ItemCalendarAllowedCalendarSharingRolesWithUserGetResponse: *NewItemCalendarAllowedCalendarSharingRolesWithUserGetResponse(),
    }
    return m
}
// CreateItemCalendarAllowedCalendarSharingRolesWithUserResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarAllowedCalendarSharingRolesWithUserResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarAllowedCalendarSharingRolesWithUserResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarAllowedCalendarSharingRolesWithUserGetResponseable instead.
type ItemCalendarAllowedCalendarSharingRolesWithUserResponseable interface {
    ItemCalendarAllowedCalendarSharingRolesWithUserGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

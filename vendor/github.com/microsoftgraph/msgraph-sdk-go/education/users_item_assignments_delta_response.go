package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use UsersItemAssignmentsDeltaGetResponseable instead.
type UsersItemAssignmentsDeltaResponse struct {
    UsersItemAssignmentsDeltaGetResponse
}
// NewUsersItemAssignmentsDeltaResponse instantiates a new UsersItemAssignmentsDeltaResponse and sets the default values.
func NewUsersItemAssignmentsDeltaResponse()(*UsersItemAssignmentsDeltaResponse) {
    m := &UsersItemAssignmentsDeltaResponse{
        UsersItemAssignmentsDeltaGetResponse: *NewUsersItemAssignmentsDeltaGetResponse(),
    }
    return m
}
// CreateUsersItemAssignmentsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUsersItemAssignmentsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUsersItemAssignmentsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use UsersItemAssignmentsDeltaGetResponseable instead.
type UsersItemAssignmentsDeltaResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UsersItemAssignmentsDeltaGetResponseable
}

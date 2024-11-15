package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use UsersItemAssignmentsItemCategoriesDeltaGetResponseable instead.
type UsersItemAssignmentsItemCategoriesDeltaResponse struct {
    UsersItemAssignmentsItemCategoriesDeltaGetResponse
}
// NewUsersItemAssignmentsItemCategoriesDeltaResponse instantiates a new UsersItemAssignmentsItemCategoriesDeltaResponse and sets the default values.
func NewUsersItemAssignmentsItemCategoriesDeltaResponse()(*UsersItemAssignmentsItemCategoriesDeltaResponse) {
    m := &UsersItemAssignmentsItemCategoriesDeltaResponse{
        UsersItemAssignmentsItemCategoriesDeltaGetResponse: *NewUsersItemAssignmentsItemCategoriesDeltaGetResponse(),
    }
    return m
}
// CreateUsersItemAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUsersItemAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUsersItemAssignmentsItemCategoriesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use UsersItemAssignmentsItemCategoriesDeltaGetResponseable instead.
type UsersItemAssignmentsItemCategoriesDeltaResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UsersItemAssignmentsItemCategoriesDeltaGetResponseable
}

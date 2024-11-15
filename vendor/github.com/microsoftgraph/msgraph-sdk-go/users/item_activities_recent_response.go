package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemActivitiesRecentGetResponseable instead.
type ItemActivitiesRecentResponse struct {
    ItemActivitiesRecentGetResponse
}
// NewItemActivitiesRecentResponse instantiates a new ItemActivitiesRecentResponse and sets the default values.
func NewItemActivitiesRecentResponse()(*ItemActivitiesRecentResponse) {
    m := &ItemActivitiesRecentResponse{
        ItemActivitiesRecentGetResponse: *NewItemActivitiesRecentGetResponse(),
    }
    return m
}
// CreateItemActivitiesRecentResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemActivitiesRecentResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemActivitiesRecentResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemActivitiesRecentGetResponseable instead.
type ItemActivitiesRecentResponseable interface {
    ItemActivitiesRecentGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

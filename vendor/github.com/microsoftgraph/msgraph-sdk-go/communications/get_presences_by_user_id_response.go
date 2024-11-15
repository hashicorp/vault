package communications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use GetPresencesByUserIdPostResponseable instead.
type GetPresencesByUserIdResponse struct {
    GetPresencesByUserIdPostResponse
}
// NewGetPresencesByUserIdResponse instantiates a new GetPresencesByUserIdResponse and sets the default values.
func NewGetPresencesByUserIdResponse()(*GetPresencesByUserIdResponse) {
    m := &GetPresencesByUserIdResponse{
        GetPresencesByUserIdPostResponse: *NewGetPresencesByUserIdPostResponse(),
    }
    return m
}
// CreateGetPresencesByUserIdResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGetPresencesByUserIdResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGetPresencesByUserIdResponse(), nil
}
// Deprecated: This class is obsolete. Use GetPresencesByUserIdPostResponseable instead.
type GetPresencesByUserIdResponseable interface {
    GetPresencesByUserIdPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

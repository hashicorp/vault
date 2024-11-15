package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use GetAllMessagesGetResponseable instead.
type GetAllMessagesResponse struct {
    GetAllMessagesGetResponse
}
// NewGetAllMessagesResponse instantiates a new GetAllMessagesResponse and sets the default values.
func NewGetAllMessagesResponse()(*GetAllMessagesResponse) {
    m := &GetAllMessagesResponse{
        GetAllMessagesGetResponse: *NewGetAllMessagesGetResponse(),
    }
    return m
}
// CreateGetAllMessagesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGetAllMessagesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGetAllMessagesResponse(), nil
}
// Deprecated: This class is obsolete. Use GetAllMessagesGetResponseable instead.
type GetAllMessagesResponseable interface {
    GetAllMessagesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package teamwork

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedTeamsGetAllMessagesGetResponseable instead.
type DeletedTeamsGetAllMessagesResponse struct {
    DeletedTeamsGetAllMessagesGetResponse
}
// NewDeletedTeamsGetAllMessagesResponse instantiates a new DeletedTeamsGetAllMessagesResponse and sets the default values.
func NewDeletedTeamsGetAllMessagesResponse()(*DeletedTeamsGetAllMessagesResponse) {
    m := &DeletedTeamsGetAllMessagesResponse{
        DeletedTeamsGetAllMessagesGetResponse: *NewDeletedTeamsGetAllMessagesGetResponse(),
    }
    return m
}
// CreateDeletedTeamsGetAllMessagesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedTeamsGetAllMessagesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedTeamsGetAllMessagesResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedTeamsGetAllMessagesGetResponseable instead.
type DeletedTeamsGetAllMessagesResponseable interface {
    DeletedTeamsGetAllMessagesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

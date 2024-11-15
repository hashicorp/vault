package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use AddPostResponseable instead.
type AddResponse struct {
    AddPostResponse
}
// NewAddResponse instantiates a new AddResponse and sets the default values.
func NewAddResponse()(*AddResponse) {
    m := &AddResponse{
        AddPostResponse: *NewAddPostResponse(),
    }
    return m
}
// CreateAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAddResponse(), nil
}
// Deprecated: This class is obsolete. Use AddPostResponseable instead.
type AddResponseable interface {
    AddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

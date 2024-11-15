package admin

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ServiceAnnouncementMessagesMarkReadPostResponseable instead.
type ServiceAnnouncementMessagesMarkReadResponse struct {
    ServiceAnnouncementMessagesMarkReadPostResponse
}
// NewServiceAnnouncementMessagesMarkReadResponse instantiates a new ServiceAnnouncementMessagesMarkReadResponse and sets the default values.
func NewServiceAnnouncementMessagesMarkReadResponse()(*ServiceAnnouncementMessagesMarkReadResponse) {
    m := &ServiceAnnouncementMessagesMarkReadResponse{
        ServiceAnnouncementMessagesMarkReadPostResponse: *NewServiceAnnouncementMessagesMarkReadPostResponse(),
    }
    return m
}
// CreateServiceAnnouncementMessagesMarkReadResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceAnnouncementMessagesMarkReadResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceAnnouncementMessagesMarkReadResponse(), nil
}
// Deprecated: This class is obsolete. Use ServiceAnnouncementMessagesMarkReadPostResponseable instead.
type ServiceAnnouncementMessagesMarkReadResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ServiceAnnouncementMessagesMarkReadPostResponseable
}

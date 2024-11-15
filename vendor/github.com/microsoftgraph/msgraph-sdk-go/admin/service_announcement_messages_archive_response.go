package admin

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ServiceAnnouncementMessagesArchivePostResponseable instead.
type ServiceAnnouncementMessagesArchiveResponse struct {
    ServiceAnnouncementMessagesArchivePostResponse
}
// NewServiceAnnouncementMessagesArchiveResponse instantiates a new ServiceAnnouncementMessagesArchiveResponse and sets the default values.
func NewServiceAnnouncementMessagesArchiveResponse()(*ServiceAnnouncementMessagesArchiveResponse) {
    m := &ServiceAnnouncementMessagesArchiveResponse{
        ServiceAnnouncementMessagesArchivePostResponse: *NewServiceAnnouncementMessagesArchivePostResponse(),
    }
    return m
}
// CreateServiceAnnouncementMessagesArchiveResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceAnnouncementMessagesArchiveResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceAnnouncementMessagesArchiveResponse(), nil
}
// Deprecated: This class is obsolete. Use ServiceAnnouncementMessagesArchivePostResponseable instead.
type ServiceAnnouncementMessagesArchiveResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ServiceAnnouncementMessagesArchivePostResponseable
}

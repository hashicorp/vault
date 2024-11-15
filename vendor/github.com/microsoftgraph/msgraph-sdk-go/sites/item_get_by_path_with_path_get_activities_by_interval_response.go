package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemGetByPathWithPathGetActivitiesByIntervalGetResponseable instead.
type ItemGetByPathWithPathGetActivitiesByIntervalResponse struct {
    ItemGetByPathWithPathGetActivitiesByIntervalGetResponse
}
// NewItemGetByPathWithPathGetActivitiesByIntervalResponse instantiates a new ItemGetByPathWithPathGetActivitiesByIntervalResponse and sets the default values.
func NewItemGetByPathWithPathGetActivitiesByIntervalResponse()(*ItemGetByPathWithPathGetActivitiesByIntervalResponse) {
    m := &ItemGetByPathWithPathGetActivitiesByIntervalResponse{
        ItemGetByPathWithPathGetActivitiesByIntervalGetResponse: *NewItemGetByPathWithPathGetActivitiesByIntervalGetResponse(),
    }
    return m
}
// CreateItemGetByPathWithPathGetActivitiesByIntervalResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemGetByPathWithPathGetActivitiesByIntervalResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemGetByPathWithPathGetActivitiesByIntervalResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemGetByPathWithPathGetActivitiesByIntervalGetResponseable instead.
type ItemGetByPathWithPathGetActivitiesByIntervalResponseable interface {
    ItemGetByPathWithPathGetActivitiesByIntervalGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

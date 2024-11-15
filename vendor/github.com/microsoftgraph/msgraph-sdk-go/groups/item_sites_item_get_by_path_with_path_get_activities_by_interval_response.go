package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponseable instead.
type ItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse struct {
    ItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponse
}
// NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse instantiates a new ItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse and sets the default values.
func NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse()(*ItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse) {
    m := &ItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse{
        ItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponse: *NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponse(),
    }
    return m
}
// CreateItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponseable instead.
type ItemSitesItemGetByPathWithPathGetActivitiesByIntervalResponseable interface {
    ItemSitesItemGetByPathWithPathGetActivitiesByIntervalGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

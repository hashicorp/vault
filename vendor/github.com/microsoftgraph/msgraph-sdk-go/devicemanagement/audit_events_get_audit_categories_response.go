package devicemanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use AuditEventsGetAuditCategoriesGetResponseable instead.
type AuditEventsGetAuditCategoriesResponse struct {
    AuditEventsGetAuditCategoriesGetResponse
}
// NewAuditEventsGetAuditCategoriesResponse instantiates a new AuditEventsGetAuditCategoriesResponse and sets the default values.
func NewAuditEventsGetAuditCategoriesResponse()(*AuditEventsGetAuditCategoriesResponse) {
    m := &AuditEventsGetAuditCategoriesResponse{
        AuditEventsGetAuditCategoriesGetResponse: *NewAuditEventsGetAuditCategoriesGetResponse(),
    }
    return m
}
// CreateAuditEventsGetAuditCategoriesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuditEventsGetAuditCategoriesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditEventsGetAuditCategoriesResponse(), nil
}
// Deprecated: This class is obsolete. Use AuditEventsGetAuditCategoriesGetResponseable instead.
type AuditEventsGetAuditCategoriesResponseable interface {
    AuditEventsGetAuditCategoriesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

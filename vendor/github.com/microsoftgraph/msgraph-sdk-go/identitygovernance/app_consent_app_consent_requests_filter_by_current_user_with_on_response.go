package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use AppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponseable instead.
type AppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse struct {
    AppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponse
}
// NewAppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse instantiates a new AppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse and sets the default values.
func NewAppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse()(*AppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse) {
    m := &AppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse{
        AppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponse: *NewAppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponse(),
    }
    return m
}
// CreateAppConsentAppConsentRequestsFilterByCurrentUserWithOnResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppConsentAppConsentRequestsFilterByCurrentUserWithOnResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppConsentAppConsentRequestsFilterByCurrentUserWithOnResponse(), nil
}
// Deprecated: This class is obsolete. Use AppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponseable instead.
type AppConsentAppConsentRequestsFilterByCurrentUserWithOnResponseable interface {
    AppConsentAppConsentRequestsFilterByCurrentUserWithOnGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package devicemanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponseable instead.
type VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse struct {
    VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponse
}
// NewVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse instantiates a new VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse and sets the default values.
func NewVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse()(*VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse) {
    m := &VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse{
        VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponse: *NewVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponse(),
    }
    return m
}
// CreateVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponse(), nil
}
// Deprecated: This class is obsolete. Use VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponseable instead.
type VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameGetResponseable
}

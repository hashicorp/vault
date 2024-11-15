package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponseable instead.
type ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse struct {
    ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponse
}
// NewItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse instantiates a new ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse and sets the default values.
func NewItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse()(*ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse) {
    m := &ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse{
        ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponse: *NewItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponse(),
    }
    return m
}
// CreateItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponseable instead.
type ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlResponseable interface {
    ItemOnlineMeetingsItemGetVirtualAppointmentJoinWebUrlGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

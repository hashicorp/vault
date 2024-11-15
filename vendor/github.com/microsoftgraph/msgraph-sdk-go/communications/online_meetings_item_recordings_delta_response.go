package communications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use OnlineMeetingsItemRecordingsDeltaGetResponseable instead.
type OnlineMeetingsItemRecordingsDeltaResponse struct {
    OnlineMeetingsItemRecordingsDeltaGetResponse
}
// NewOnlineMeetingsItemRecordingsDeltaResponse instantiates a new OnlineMeetingsItemRecordingsDeltaResponse and sets the default values.
func NewOnlineMeetingsItemRecordingsDeltaResponse()(*OnlineMeetingsItemRecordingsDeltaResponse) {
    m := &OnlineMeetingsItemRecordingsDeltaResponse{
        OnlineMeetingsItemRecordingsDeltaGetResponse: *NewOnlineMeetingsItemRecordingsDeltaGetResponse(),
    }
    return m
}
// CreateOnlineMeetingsItemRecordingsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnlineMeetingsItemRecordingsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnlineMeetingsItemRecordingsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use OnlineMeetingsItemRecordingsDeltaGetResponseable instead.
type OnlineMeetingsItemRecordingsDeltaResponseable interface {
    OnlineMeetingsItemRecordingsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

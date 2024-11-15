package devicemanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable instead.
type ImportedWindowsAutopilotDeviceIdentitiesImportResponse struct {
    ImportedWindowsAutopilotDeviceIdentitiesImportPostResponse
}
// NewImportedWindowsAutopilotDeviceIdentitiesImportResponse instantiates a new ImportedWindowsAutopilotDeviceIdentitiesImportResponse and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentitiesImportResponse()(*ImportedWindowsAutopilotDeviceIdentitiesImportResponse) {
    m := &ImportedWindowsAutopilotDeviceIdentitiesImportResponse{
        ImportedWindowsAutopilotDeviceIdentitiesImportPostResponse: *NewImportedWindowsAutopilotDeviceIdentitiesImportPostResponse(),
    }
    return m
}
// CreateImportedWindowsAutopilotDeviceIdentitiesImportResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateImportedWindowsAutopilotDeviceIdentitiesImportResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewImportedWindowsAutopilotDeviceIdentitiesImportResponse(), nil
}
// Deprecated: This class is obsolete. Use ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable instead.
type ImportedWindowsAutopilotDeviceIdentitiesImportResponseable interface {
    ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

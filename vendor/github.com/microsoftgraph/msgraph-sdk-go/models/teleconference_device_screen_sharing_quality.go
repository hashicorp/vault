package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeleconferenceDeviceScreenSharingQuality struct {
    TeleconferenceDeviceVideoQuality
}
// NewTeleconferenceDeviceScreenSharingQuality instantiates a new TeleconferenceDeviceScreenSharingQuality and sets the default values.
func NewTeleconferenceDeviceScreenSharingQuality()(*TeleconferenceDeviceScreenSharingQuality) {
    m := &TeleconferenceDeviceScreenSharingQuality{
        TeleconferenceDeviceVideoQuality: *NewTeleconferenceDeviceVideoQuality(),
    }
    return m
}
// CreateTeleconferenceDeviceScreenSharingQualityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeleconferenceDeviceScreenSharingQualityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeleconferenceDeviceScreenSharingQuality(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeleconferenceDeviceScreenSharingQuality) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TeleconferenceDeviceVideoQuality.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *TeleconferenceDeviceScreenSharingQuality) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TeleconferenceDeviceVideoQuality.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type TeleconferenceDeviceScreenSharingQualityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    TeleconferenceDeviceVideoQualityable
}

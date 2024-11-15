package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AllDeviceRegistrationMembership struct {
    DeviceRegistrationMembership
}
// NewAllDeviceRegistrationMembership instantiates a new AllDeviceRegistrationMembership and sets the default values.
func NewAllDeviceRegistrationMembership()(*AllDeviceRegistrationMembership) {
    m := &AllDeviceRegistrationMembership{
        DeviceRegistrationMembership: *NewDeviceRegistrationMembership(),
    }
    odataTypeValue := "#microsoft.graph.allDeviceRegistrationMembership"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAllDeviceRegistrationMembershipFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAllDeviceRegistrationMembershipFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAllDeviceRegistrationMembership(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AllDeviceRegistrationMembership) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceRegistrationMembership.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *AllDeviceRegistrationMembership) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceRegistrationMembership.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type AllDeviceRegistrationMembershipable interface {
    DeviceRegistrationMembershipable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

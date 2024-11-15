package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsLicensingDetails struct {
    Entity
}
// NewTeamsLicensingDetails instantiates a new TeamsLicensingDetails and sets the default values.
func NewTeamsLicensingDetails()(*TeamsLicensingDetails) {
    m := &TeamsLicensingDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsLicensingDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsLicensingDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsLicensingDetails(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsLicensingDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["hasTeamsLicense"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasTeamsLicense(val)
        }
        return nil
    }
    return res
}
// GetHasTeamsLicense gets the hasTeamsLicense property value. Indicates whether the user has a valid license to use Microsoft Teams.
// returns a *bool when successful
func (m *TeamsLicensingDetails) GetHasTeamsLicense()(*bool) {
    val, err := m.GetBackingStore().Get("hasTeamsLicense")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamsLicensingDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("hasTeamsLicense", m.GetHasTeamsLicense())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHasTeamsLicense sets the hasTeamsLicense property value. Indicates whether the user has a valid license to use Microsoft Teams.
func (m *TeamsLicensingDetails) SetHasTeamsLicense(value *bool)() {
    err := m.GetBackingStore().Set("hasTeamsLicense", value)
    if err != nil {
        panic(err)
    }
}
type TeamsLicensingDetailsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHasTeamsLicense()(*bool)
    SetHasTeamsLicense(value *bool)()
}

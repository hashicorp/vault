package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamworkApplicationIdentity struct {
    Identity
}
// NewTeamworkApplicationIdentity instantiates a new TeamworkApplicationIdentity and sets the default values.
func NewTeamworkApplicationIdentity()(*TeamworkApplicationIdentity) {
    m := &TeamworkApplicationIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.teamworkApplicationIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTeamworkApplicationIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamworkApplicationIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkApplicationIdentity(), nil
}
// GetApplicationIdentityType gets the applicationIdentityType property value. Type of application that is referenced. Possible values are: aadApplication, bot, tenantBot, office365Connector, outgoingWebhook, and unknownFutureValue.
// returns a *TeamworkApplicationIdentityType when successful
func (m *TeamworkApplicationIdentity) GetApplicationIdentityType()(*TeamworkApplicationIdentityType) {
    val, err := m.GetBackingStore().Get("applicationIdentityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamworkApplicationIdentityType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamworkApplicationIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["applicationIdentityType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamworkApplicationIdentityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationIdentityType(val.(*TeamworkApplicationIdentityType))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TeamworkApplicationIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApplicationIdentityType() != nil {
        cast := (*m.GetApplicationIdentityType()).String()
        err = writer.WriteStringValue("applicationIdentityType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationIdentityType sets the applicationIdentityType property value. Type of application that is referenced. Possible values are: aadApplication, bot, tenantBot, office365Connector, outgoingWebhook, and unknownFutureValue.
func (m *TeamworkApplicationIdentity) SetApplicationIdentityType(value *TeamworkApplicationIdentityType)() {
    err := m.GetBackingStore().Set("applicationIdentityType", value)
    if err != nil {
        panic(err)
    }
}
type TeamworkApplicationIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationIdentityType()(*TeamworkApplicationIdentityType)
    SetApplicationIdentityType(value *TeamworkApplicationIdentityType)()
}

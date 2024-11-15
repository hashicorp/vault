package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DefaultInvitationRedemptionIdentityProviderConfiguration struct {
    InvitationRedemptionIdentityProviderConfiguration
}
// NewDefaultInvitationRedemptionIdentityProviderConfiguration instantiates a new DefaultInvitationRedemptionIdentityProviderConfiguration and sets the default values.
func NewDefaultInvitationRedemptionIdentityProviderConfiguration()(*DefaultInvitationRedemptionIdentityProviderConfiguration) {
    m := &DefaultInvitationRedemptionIdentityProviderConfiguration{
        InvitationRedemptionIdentityProviderConfiguration: *NewInvitationRedemptionIdentityProviderConfiguration(),
    }
    return m
}
// CreateDefaultInvitationRedemptionIdentityProviderConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDefaultInvitationRedemptionIdentityProviderConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDefaultInvitationRedemptionIdentityProviderConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DefaultInvitationRedemptionIdentityProviderConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.InvitationRedemptionIdentityProviderConfiguration.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *DefaultInvitationRedemptionIdentityProviderConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.InvitationRedemptionIdentityProviderConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type DefaultInvitationRedemptionIdentityProviderConfigurationable interface {
    InvitationRedemptionIdentityProviderConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

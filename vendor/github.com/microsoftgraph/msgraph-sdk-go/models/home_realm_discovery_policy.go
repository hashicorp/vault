package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type HomeRealmDiscoveryPolicy struct {
    StsPolicy
}
// NewHomeRealmDiscoveryPolicy instantiates a new HomeRealmDiscoveryPolicy and sets the default values.
func NewHomeRealmDiscoveryPolicy()(*HomeRealmDiscoveryPolicy) {
    m := &HomeRealmDiscoveryPolicy{
        StsPolicy: *NewStsPolicy(),
    }
    odataTypeValue := "#microsoft.graph.homeRealmDiscoveryPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateHomeRealmDiscoveryPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHomeRealmDiscoveryPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHomeRealmDiscoveryPolicy(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *HomeRealmDiscoveryPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.StsPolicy.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *HomeRealmDiscoveryPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.StsPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type HomeRealmDiscoveryPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    StsPolicyable
}

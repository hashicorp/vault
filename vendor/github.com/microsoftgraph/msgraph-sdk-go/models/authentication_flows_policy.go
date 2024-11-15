package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationFlowsPolicy struct {
    Entity
}
// NewAuthenticationFlowsPolicy instantiates a new AuthenticationFlowsPolicy and sets the default values.
func NewAuthenticationFlowsPolicy()(*AuthenticationFlowsPolicy) {
    m := &AuthenticationFlowsPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationFlowsPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationFlowsPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationFlowsPolicy(), nil
}
// GetDescription gets the description property value. Inherited property. A description of the policy. Optional. Read-only.
// returns a *string when successful
func (m *AuthenticationFlowsPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Inherited property. The human-readable name of the policy. Optional. Read-only.
// returns a *string when successful
func (m *AuthenticationFlowsPolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationFlowsPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["selfServiceSignUp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSelfServiceSignUpAuthenticationFlowConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSelfServiceSignUp(val.(SelfServiceSignUpAuthenticationFlowConfigurationable))
        }
        return nil
    }
    return res
}
// GetSelfServiceSignUp gets the selfServiceSignUp property value. Contains selfServiceSignUpAuthenticationFlowConfiguration settings that convey whether self-service sign-up is enabled or disabled. Optional. Read-only.
// returns a SelfServiceSignUpAuthenticationFlowConfigurationable when successful
func (m *AuthenticationFlowsPolicy) GetSelfServiceSignUp()(SelfServiceSignUpAuthenticationFlowConfigurationable) {
    val, err := m.GetBackingStore().Get("selfServiceSignUp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SelfServiceSignUpAuthenticationFlowConfigurationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationFlowsPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("selfServiceSignUp", m.GetSelfServiceSignUp())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Inherited property. A description of the policy. Optional. Read-only.
func (m *AuthenticationFlowsPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Inherited property. The human-readable name of the policy. Optional. Read-only.
func (m *AuthenticationFlowsPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetSelfServiceSignUp sets the selfServiceSignUp property value. Contains selfServiceSignUpAuthenticationFlowConfiguration settings that convey whether self-service sign-up is enabled or disabled. Optional. Read-only.
func (m *AuthenticationFlowsPolicy) SetSelfServiceSignUp(value SelfServiceSignUpAuthenticationFlowConfigurationable)() {
    err := m.GetBackingStore().Set("selfServiceSignUp", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationFlowsPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetSelfServiceSignUp()(SelfServiceSignUpAuthenticationFlowConfigurationable)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetSelfServiceSignUp(value SelfServiceSignUpAuthenticationFlowConfigurationable)()
}

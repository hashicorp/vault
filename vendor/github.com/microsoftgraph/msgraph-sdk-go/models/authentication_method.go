package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationMethod struct {
    Entity
}
// NewAuthenticationMethod instantiates a new AuthenticationMethod and sets the default values.
func NewAuthenticationMethod()(*AuthenticationMethod) {
    m := &AuthenticationMethod{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationMethodFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationMethodFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.emailAuthenticationMethod":
                        return NewEmailAuthenticationMethod(), nil
                    case "#microsoft.graph.fido2AuthenticationMethod":
                        return NewFido2AuthenticationMethod(), nil
                    case "#microsoft.graph.microsoftAuthenticatorAuthenticationMethod":
                        return NewMicrosoftAuthenticatorAuthenticationMethod(), nil
                    case "#microsoft.graph.passwordAuthenticationMethod":
                        return NewPasswordAuthenticationMethod(), nil
                    case "#microsoft.graph.phoneAuthenticationMethod":
                        return NewPhoneAuthenticationMethod(), nil
                    case "#microsoft.graph.softwareOathAuthenticationMethod":
                        return NewSoftwareOathAuthenticationMethod(), nil
                    case "#microsoft.graph.temporaryAccessPassAuthenticationMethod":
                        return NewTemporaryAccessPassAuthenticationMethod(), nil
                    case "#microsoft.graph.windowsHelloForBusinessAuthenticationMethod":
                        return NewWindowsHelloForBusinessAuthenticationMethod(), nil
                }
            }
        }
    }
    return NewAuthenticationMethod(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationMethod) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *AuthenticationMethod) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type AuthenticationMethodable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

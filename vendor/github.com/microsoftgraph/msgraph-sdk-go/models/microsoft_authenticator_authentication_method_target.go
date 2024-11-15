package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MicrosoftAuthenticatorAuthenticationMethodTarget struct {
    AuthenticationMethodTarget
}
// NewMicrosoftAuthenticatorAuthenticationMethodTarget instantiates a new MicrosoftAuthenticatorAuthenticationMethodTarget and sets the default values.
func NewMicrosoftAuthenticatorAuthenticationMethodTarget()(*MicrosoftAuthenticatorAuthenticationMethodTarget) {
    m := &MicrosoftAuthenticatorAuthenticationMethodTarget{
        AuthenticationMethodTarget: *NewAuthenticationMethodTarget(),
    }
    return m
}
// CreateMicrosoftAuthenticatorAuthenticationMethodTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMicrosoftAuthenticatorAuthenticationMethodTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftAuthenticatorAuthenticationMethodTarget(), nil
}
// GetAuthenticationMode gets the authenticationMode property value. The authenticationMode property
// returns a *MicrosoftAuthenticatorAuthenticationMode when successful
func (m *MicrosoftAuthenticatorAuthenticationMethodTarget) GetAuthenticationMode()(*MicrosoftAuthenticatorAuthenticationMode) {
    val, err := m.GetBackingStore().Get("authenticationMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MicrosoftAuthenticatorAuthenticationMode)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MicrosoftAuthenticatorAuthenticationMethodTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodTarget.GetFieldDeserializers()
    res["authenticationMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMicrosoftAuthenticatorAuthenticationMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationMode(val.(*MicrosoftAuthenticatorAuthenticationMode))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *MicrosoftAuthenticatorAuthenticationMethodTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodTarget.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationMode() != nil {
        cast := (*m.GetAuthenticationMode()).String()
        err = writer.WriteStringValue("authenticationMode", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationMode sets the authenticationMode property value. The authenticationMode property
func (m *MicrosoftAuthenticatorAuthenticationMethodTarget) SetAuthenticationMode(value *MicrosoftAuthenticatorAuthenticationMode)() {
    err := m.GetBackingStore().Set("authenticationMode", value)
    if err != nil {
        panic(err)
    }
}
type MicrosoftAuthenticatorAuthenticationMethodTargetable interface {
    AuthenticationMethodTargetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationMode()(*MicrosoftAuthenticatorAuthenticationMode)
    SetAuthenticationMode(value *MicrosoftAuthenticatorAuthenticationMode)()
}

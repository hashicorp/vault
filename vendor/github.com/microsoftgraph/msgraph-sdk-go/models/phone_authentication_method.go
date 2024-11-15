package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PhoneAuthenticationMethod struct {
    AuthenticationMethod
}
// NewPhoneAuthenticationMethod instantiates a new PhoneAuthenticationMethod and sets the default values.
func NewPhoneAuthenticationMethod()(*PhoneAuthenticationMethod) {
    m := &PhoneAuthenticationMethod{
        AuthenticationMethod: *NewAuthenticationMethod(),
    }
    odataTypeValue := "#microsoft.graph.phoneAuthenticationMethod"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePhoneAuthenticationMethodFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePhoneAuthenticationMethodFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPhoneAuthenticationMethod(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PhoneAuthenticationMethod) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethod.GetFieldDeserializers()
    res["phoneNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoneNumber(val)
        }
        return nil
    }
    res["phoneType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationPhoneType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoneType(val.(*AuthenticationPhoneType))
        }
        return nil
    }
    res["smsSignInState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationMethodSignInState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmsSignInState(val.(*AuthenticationMethodSignInState))
        }
        return nil
    }
    return res
}
// GetPhoneNumber gets the phoneNumber property value. The phone number to text or call for authentication. Phone numbers use the format +{country code} {number}x{extension}, with extension optional. For example, +1 5555551234 or +1 5555551234x123 are valid. Numbers are rejected when creating or updating if they don't match the required format.
// returns a *string when successful
func (m *PhoneAuthenticationMethod) GetPhoneNumber()(*string) {
    val, err := m.GetBackingStore().Get("phoneNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhoneType gets the phoneType property value. The type of this phone. Possible values are: mobile, alternateMobile, or office.
// returns a *AuthenticationPhoneType when successful
func (m *PhoneAuthenticationMethod) GetPhoneType()(*AuthenticationPhoneType) {
    val, err := m.GetBackingStore().Get("phoneType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationPhoneType)
    }
    return nil
}
// GetSmsSignInState gets the smsSignInState property value. Whether a phone is ready to be used for SMS sign-in or not. Possible values are: notSupported, notAllowedByPolicy, notEnabled, phoneNumberNotUnique, ready, or notConfigured, unknownFutureValue.
// returns a *AuthenticationMethodSignInState when successful
func (m *PhoneAuthenticationMethod) GetSmsSignInState()(*AuthenticationMethodSignInState) {
    val, err := m.GetBackingStore().Get("smsSignInState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationMethodSignInState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PhoneAuthenticationMethod) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethod.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("phoneNumber", m.GetPhoneNumber())
        if err != nil {
            return err
        }
    }
    if m.GetPhoneType() != nil {
        cast := (*m.GetPhoneType()).String()
        err = writer.WriteStringValue("phoneType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSmsSignInState() != nil {
        cast := (*m.GetSmsSignInState()).String()
        err = writer.WriteStringValue("smsSignInState", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPhoneNumber sets the phoneNumber property value. The phone number to text or call for authentication. Phone numbers use the format +{country code} {number}x{extension}, with extension optional. For example, +1 5555551234 or +1 5555551234x123 are valid. Numbers are rejected when creating or updating if they don't match the required format.
func (m *PhoneAuthenticationMethod) SetPhoneNumber(value *string)() {
    err := m.GetBackingStore().Set("phoneNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoneType sets the phoneType property value. The type of this phone. Possible values are: mobile, alternateMobile, or office.
func (m *PhoneAuthenticationMethod) SetPhoneType(value *AuthenticationPhoneType)() {
    err := m.GetBackingStore().Set("phoneType", value)
    if err != nil {
        panic(err)
    }
}
// SetSmsSignInState sets the smsSignInState property value. Whether a phone is ready to be used for SMS sign-in or not. Possible values are: notSupported, notAllowedByPolicy, notEnabled, phoneNumberNotUnique, ready, or notConfigured, unknownFutureValue.
func (m *PhoneAuthenticationMethod) SetSmsSignInState(value *AuthenticationMethodSignInState)() {
    err := m.GetBackingStore().Set("smsSignInState", value)
    if err != nil {
        panic(err)
    }
}
type PhoneAuthenticationMethodable interface {
    AuthenticationMethodable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPhoneNumber()(*string)
    GetPhoneType()(*AuthenticationPhoneType)
    GetSmsSignInState()(*AuthenticationMethodSignInState)
    SetPhoneNumber(value *string)()
    SetPhoneType(value *AuthenticationPhoneType)()
    SetSmsSignInState(value *AuthenticationMethodSignInState)()
}

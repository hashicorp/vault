package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Authentication struct {
    Entity
}
// NewAuthentication instantiates a new Authentication and sets the default values.
func NewAuthentication()(*Authentication) {
    m := &Authentication{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthentication(), nil
}
// GetEmailMethods gets the emailMethods property value. The email address registered to a user for authentication.
// returns a []EmailAuthenticationMethodable when successful
func (m *Authentication) GetEmailMethods()([]EmailAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("emailMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EmailAuthenticationMethodable)
    }
    return nil
}
// GetFido2Methods gets the fido2Methods property value. Represents the FIDO2 security keys registered to a user for authentication.
// returns a []Fido2AuthenticationMethodable when successful
func (m *Authentication) GetFido2Methods()([]Fido2AuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("fido2Methods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Fido2AuthenticationMethodable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Authentication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["emailMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEmailAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EmailAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EmailAuthenticationMethodable)
                }
            }
            m.SetEmailMethods(res)
        }
        return nil
    }
    res["fido2Methods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateFido2AuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Fido2AuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Fido2AuthenticationMethodable)
                }
            }
            m.SetFido2Methods(res)
        }
        return nil
    }
    res["methods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationMethodable)
                }
            }
            m.SetMethods(res)
        }
        return nil
    }
    res["microsoftAuthenticatorMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMicrosoftAuthenticatorAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MicrosoftAuthenticatorAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MicrosoftAuthenticatorAuthenticationMethodable)
                }
            }
            m.SetMicrosoftAuthenticatorMethods(res)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLongRunningOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LongRunningOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LongRunningOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["passwordMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePasswordAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PasswordAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PasswordAuthenticationMethodable)
                }
            }
            m.SetPasswordMethods(res)
        }
        return nil
    }
    res["phoneMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePhoneAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PhoneAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PhoneAuthenticationMethodable)
                }
            }
            m.SetPhoneMethods(res)
        }
        return nil
    }
    res["softwareOathMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSoftwareOathAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SoftwareOathAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SoftwareOathAuthenticationMethodable)
                }
            }
            m.SetSoftwareOathMethods(res)
        }
        return nil
    }
    res["temporaryAccessPassMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTemporaryAccessPassAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TemporaryAccessPassAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TemporaryAccessPassAuthenticationMethodable)
                }
            }
            m.SetTemporaryAccessPassMethods(res)
        }
        return nil
    }
    res["windowsHelloForBusinessMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsHelloForBusinessAuthenticationMethodFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsHelloForBusinessAuthenticationMethodable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsHelloForBusinessAuthenticationMethodable)
                }
            }
            m.SetWindowsHelloForBusinessMethods(res)
        }
        return nil
    }
    return res
}
// GetMethods gets the methods property value. Represents all authentication methods registered to a user.
// returns a []AuthenticationMethodable when successful
func (m *Authentication) GetMethods()([]AuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("methods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodable)
    }
    return nil
}
// GetMicrosoftAuthenticatorMethods gets the microsoftAuthenticatorMethods property value. The details of the Microsoft Authenticator app registered to a user for authentication.
// returns a []MicrosoftAuthenticatorAuthenticationMethodable when successful
func (m *Authentication) GetMicrosoftAuthenticatorMethods()([]MicrosoftAuthenticatorAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("microsoftAuthenticatorMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MicrosoftAuthenticatorAuthenticationMethodable)
    }
    return nil
}
// GetOperations gets the operations property value. Represents the status of a long-running operation, such as a password reset operation.
// returns a []LongRunningOperationable when successful
func (m *Authentication) GetOperations()([]LongRunningOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LongRunningOperationable)
    }
    return nil
}
// GetPasswordMethods gets the passwordMethods property value. Represents the password registered to a user for authentication. For security, the password itself is never returned in the object, but action can be taken to reset a password.
// returns a []PasswordAuthenticationMethodable when successful
func (m *Authentication) GetPasswordMethods()([]PasswordAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("passwordMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PasswordAuthenticationMethodable)
    }
    return nil
}
// GetPhoneMethods gets the phoneMethods property value. The phone numbers registered to a user for authentication.
// returns a []PhoneAuthenticationMethodable when successful
func (m *Authentication) GetPhoneMethods()([]PhoneAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("phoneMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PhoneAuthenticationMethodable)
    }
    return nil
}
// GetSoftwareOathMethods gets the softwareOathMethods property value. The software OATH time-based one-time password (TOTP) applications registered to a user for authentication.
// returns a []SoftwareOathAuthenticationMethodable when successful
func (m *Authentication) GetSoftwareOathMethods()([]SoftwareOathAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("softwareOathMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SoftwareOathAuthenticationMethodable)
    }
    return nil
}
// GetTemporaryAccessPassMethods gets the temporaryAccessPassMethods property value. Represents a Temporary Access Pass registered to a user for authentication through time-limited passcodes.
// returns a []TemporaryAccessPassAuthenticationMethodable when successful
func (m *Authentication) GetTemporaryAccessPassMethods()([]TemporaryAccessPassAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("temporaryAccessPassMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TemporaryAccessPassAuthenticationMethodable)
    }
    return nil
}
// GetWindowsHelloForBusinessMethods gets the windowsHelloForBusinessMethods property value. Represents the Windows Hello for Business authentication method registered to a user for authentication.
// returns a []WindowsHelloForBusinessAuthenticationMethodable when successful
func (m *Authentication) GetWindowsHelloForBusinessMethods()([]WindowsHelloForBusinessAuthenticationMethodable) {
    val, err := m.GetBackingStore().Get("windowsHelloForBusinessMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsHelloForBusinessAuthenticationMethodable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Authentication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetEmailMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEmailMethods()))
        for i, v := range m.GetEmailMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("emailMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFido2Methods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFido2Methods()))
        for i, v := range m.GetFido2Methods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("fido2Methods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMethods()))
        for i, v := range m.GetMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("methods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMicrosoftAuthenticatorMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMicrosoftAuthenticatorMethods()))
        for i, v := range m.GetMicrosoftAuthenticatorMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("microsoftAuthenticatorMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPasswordMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPasswordMethods()))
        for i, v := range m.GetPasswordMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("passwordMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPhoneMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPhoneMethods()))
        for i, v := range m.GetPhoneMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("phoneMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSoftwareOathMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSoftwareOathMethods()))
        for i, v := range m.GetSoftwareOathMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("softwareOathMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTemporaryAccessPassMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTemporaryAccessPassMethods()))
        for i, v := range m.GetTemporaryAccessPassMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("temporaryAccessPassMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsHelloForBusinessMethods() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWindowsHelloForBusinessMethods()))
        for i, v := range m.GetWindowsHelloForBusinessMethods() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("windowsHelloForBusinessMethods", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmailMethods sets the emailMethods property value. The email address registered to a user for authentication.
func (m *Authentication) SetEmailMethods(value []EmailAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("emailMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetFido2Methods sets the fido2Methods property value. Represents the FIDO2 security keys registered to a user for authentication.
func (m *Authentication) SetFido2Methods(value []Fido2AuthenticationMethodable)() {
    err := m.GetBackingStore().Set("fido2Methods", value)
    if err != nil {
        panic(err)
    }
}
// SetMethods sets the methods property value. Represents all authentication methods registered to a user.
func (m *Authentication) SetMethods(value []AuthenticationMethodable)() {
    err := m.GetBackingStore().Set("methods", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftAuthenticatorMethods sets the microsoftAuthenticatorMethods property value. The details of the Microsoft Authenticator app registered to a user for authentication.
func (m *Authentication) SetMicrosoftAuthenticatorMethods(value []MicrosoftAuthenticatorAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("microsoftAuthenticatorMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. Represents the status of a long-running operation, such as a password reset operation.
func (m *Authentication) SetOperations(value []LongRunningOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMethods sets the passwordMethods property value. Represents the password registered to a user for authentication. For security, the password itself is never returned in the object, but action can be taken to reset a password.
func (m *Authentication) SetPasswordMethods(value []PasswordAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("passwordMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoneMethods sets the phoneMethods property value. The phone numbers registered to a user for authentication.
func (m *Authentication) SetPhoneMethods(value []PhoneAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("phoneMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetSoftwareOathMethods sets the softwareOathMethods property value. The software OATH time-based one-time password (TOTP) applications registered to a user for authentication.
func (m *Authentication) SetSoftwareOathMethods(value []SoftwareOathAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("softwareOathMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetTemporaryAccessPassMethods sets the temporaryAccessPassMethods property value. Represents a Temporary Access Pass registered to a user for authentication through time-limited passcodes.
func (m *Authentication) SetTemporaryAccessPassMethods(value []TemporaryAccessPassAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("temporaryAccessPassMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsHelloForBusinessMethods sets the windowsHelloForBusinessMethods property value. Represents the Windows Hello for Business authentication method registered to a user for authentication.
func (m *Authentication) SetWindowsHelloForBusinessMethods(value []WindowsHelloForBusinessAuthenticationMethodable)() {
    err := m.GetBackingStore().Set("windowsHelloForBusinessMethods", value)
    if err != nil {
        panic(err)
    }
}
type Authenticationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEmailMethods()([]EmailAuthenticationMethodable)
    GetFido2Methods()([]Fido2AuthenticationMethodable)
    GetMethods()([]AuthenticationMethodable)
    GetMicrosoftAuthenticatorMethods()([]MicrosoftAuthenticatorAuthenticationMethodable)
    GetOperations()([]LongRunningOperationable)
    GetPasswordMethods()([]PasswordAuthenticationMethodable)
    GetPhoneMethods()([]PhoneAuthenticationMethodable)
    GetSoftwareOathMethods()([]SoftwareOathAuthenticationMethodable)
    GetTemporaryAccessPassMethods()([]TemporaryAccessPassAuthenticationMethodable)
    GetWindowsHelloForBusinessMethods()([]WindowsHelloForBusinessAuthenticationMethodable)
    SetEmailMethods(value []EmailAuthenticationMethodable)()
    SetFido2Methods(value []Fido2AuthenticationMethodable)()
    SetMethods(value []AuthenticationMethodable)()
    SetMicrosoftAuthenticatorMethods(value []MicrosoftAuthenticatorAuthenticationMethodable)()
    SetOperations(value []LongRunningOperationable)()
    SetPasswordMethods(value []PasswordAuthenticationMethodable)()
    SetPhoneMethods(value []PhoneAuthenticationMethodable)()
    SetSoftwareOathMethods(value []SoftwareOathAuthenticationMethodable)()
    SetTemporaryAccessPassMethods(value []TemporaryAccessPassAuthenticationMethodable)()
    SetWindowsHelloForBusinessMethods(value []WindowsHelloForBusinessAuthenticationMethodable)()
}

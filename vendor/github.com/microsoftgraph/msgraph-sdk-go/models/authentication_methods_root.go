package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationMethodsRoot struct {
    Entity
}
// NewAuthenticationMethodsRoot instantiates a new AuthenticationMethodsRoot and sets the default values.
func NewAuthenticationMethodsRoot()(*AuthenticationMethodsRoot) {
    m := &AuthenticationMethodsRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationMethodsRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationMethodsRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationMethodsRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationMethodsRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["userRegistrationDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserRegistrationDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserRegistrationDetailsable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserRegistrationDetailsable)
                }
            }
            m.SetUserRegistrationDetails(res)
        }
        return nil
    }
    return res
}
// GetUserRegistrationDetails gets the userRegistrationDetails property value. Represents the state of a user's authentication methods, including which methods are registered and which features the user is registered and capable of (such as multifactor authentication, self-service password reset, and passwordless authentication).
// returns a []UserRegistrationDetailsable when successful
func (m *AuthenticationMethodsRoot) GetUserRegistrationDetails()([]UserRegistrationDetailsable) {
    val, err := m.GetBackingStore().Get("userRegistrationDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserRegistrationDetailsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationMethodsRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetUserRegistrationDetails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserRegistrationDetails()))
        for i, v := range m.GetUserRegistrationDetails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userRegistrationDetails", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUserRegistrationDetails sets the userRegistrationDetails property value. Represents the state of a user's authentication methods, including which methods are registered and which features the user is registered and capable of (such as multifactor authentication, self-service password reset, and passwordless authentication).
func (m *AuthenticationMethodsRoot) SetUserRegistrationDetails(value []UserRegistrationDetailsable)() {
    err := m.GetBackingStore().Set("userRegistrationDetails", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationMethodsRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUserRegistrationDetails()([]UserRegistrationDetailsable)
    SetUserRegistrationDetails(value []UserRegistrationDetailsable)()
}

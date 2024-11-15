package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp struct {
    OnAuthenticationMethodLoadStartHandler
}
// NewOnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp instantiates a new OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp and sets the default values.
func NewOnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp()(*OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp) {
    m := &OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp{
        OnAuthenticationMethodLoadStartHandler: *NewOnAuthenticationMethodLoadStartHandler(),
    }
    odataTypeValue := "#microsoft.graph.onAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUpFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUpFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OnAuthenticationMethodLoadStartHandler.GetFieldDeserializers()
    res["identityProviders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityProviderBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityProviderBaseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityProviderBaseable)
                }
            }
            m.SetIdentityProviders(res)
        }
        return nil
    }
    return res
}
// GetIdentityProviders gets the identityProviders property value. The identityProviders property
// returns a []IdentityProviderBaseable when successful
func (m *OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp) GetIdentityProviders()([]IdentityProviderBaseable) {
    val, err := m.GetBackingStore().Get("identityProviders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityProviderBaseable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OnAuthenticationMethodLoadStartHandler.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetIdentityProviders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIdentityProviders()))
        for i, v := range m.GetIdentityProviders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("identityProviders", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIdentityProviders sets the identityProviders property value. The identityProviders property
func (m *OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUp) SetIdentityProviders(value []IdentityProviderBaseable)() {
    err := m.GetBackingStore().Set("identityProviders", value)
    if err != nil {
        panic(err)
    }
}
type OnAuthenticationMethodLoadStartExternalUsersSelfServiceSignUpable interface {
    OnAuthenticationMethodLoadStartHandlerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIdentityProviders()([]IdentityProviderBaseable)
    SetIdentityProviders(value []IdentityProviderBaseable)()
}

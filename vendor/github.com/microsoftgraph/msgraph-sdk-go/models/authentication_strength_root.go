package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationStrengthRoot struct {
    Entity
}
// NewAuthenticationStrengthRoot instantiates a new AuthenticationStrengthRoot and sets the default values.
func NewAuthenticationStrengthRoot()(*AuthenticationStrengthRoot) {
    m := &AuthenticationStrengthRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationStrengthRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationStrengthRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationStrengthRoot(), nil
}
// GetAuthenticationMethodModes gets the authenticationMethodModes property value. Names and descriptions of all valid authentication method modes in the system.
// returns a []AuthenticationMethodModeDetailable when successful
func (m *AuthenticationStrengthRoot) GetAuthenticationMethodModes()([]AuthenticationMethodModeDetailable) {
    val, err := m.GetBackingStore().Get("authenticationMethodModes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodModeDetailable)
    }
    return nil
}
// GetCombinations gets the combinations property value. The combinations property
// returns a []AuthenticationMethodModes when successful
func (m *AuthenticationStrengthRoot) GetCombinations()([]AuthenticationMethodModes) {
    val, err := m.GetBackingStore().Get("combinations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodModes)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthenticationStrengthRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationMethodModes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationMethodModeDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodModeDetailable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationMethodModeDetailable)
                }
            }
            m.SetAuthenticationMethodModes(res)
        }
        return nil
    }
    res["combinations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseAuthenticationMethodModes)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodModes, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*AuthenticationMethodModes))
                }
            }
            m.SetCombinations(res)
        }
        return nil
    }
    res["policies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationStrengthPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationStrengthPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationStrengthPolicyable)
                }
            }
            m.SetPolicies(res)
        }
        return nil
    }
    return res
}
// GetPolicies gets the policies property value. A collection of authentication strength policies that exist for this tenant, including both built-in and custom policies.
// returns a []AuthenticationStrengthPolicyable when successful
func (m *AuthenticationStrengthRoot) GetPolicies()([]AuthenticationStrengthPolicyable) {
    val, err := m.GetBackingStore().Get("policies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationStrengthPolicyable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationStrengthRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationMethodModes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthenticationMethodModes()))
        for i, v := range m.GetAuthenticationMethodModes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("authenticationMethodModes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCombinations() != nil {
        err = writer.WriteCollectionOfStringValues("combinations", SerializeAuthenticationMethodModes(m.GetCombinations()))
        if err != nil {
            return err
        }
    }
    if m.GetPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPolicies()))
        for i, v := range m.GetPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("policies", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationMethodModes sets the authenticationMethodModes property value. Names and descriptions of all valid authentication method modes in the system.
func (m *AuthenticationStrengthRoot) SetAuthenticationMethodModes(value []AuthenticationMethodModeDetailable)() {
    err := m.GetBackingStore().Set("authenticationMethodModes", value)
    if err != nil {
        panic(err)
    }
}
// SetCombinations sets the combinations property value. The combinations property
func (m *AuthenticationStrengthRoot) SetCombinations(value []AuthenticationMethodModes)() {
    err := m.GetBackingStore().Set("combinations", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicies sets the policies property value. A collection of authentication strength policies that exist for this tenant, including both built-in and custom policies.
func (m *AuthenticationStrengthRoot) SetPolicies(value []AuthenticationStrengthPolicyable)() {
    err := m.GetBackingStore().Set("policies", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationStrengthRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationMethodModes()([]AuthenticationMethodModeDetailable)
    GetCombinations()([]AuthenticationMethodModes)
    GetPolicies()([]AuthenticationStrengthPolicyable)
    SetAuthenticationMethodModes(value []AuthenticationMethodModeDetailable)()
    SetCombinations(value []AuthenticationMethodModes)()
    SetPolicies(value []AuthenticationStrengthPolicyable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IdentityContainer struct {
    Entity
}
// NewIdentityContainer instantiates a new IdentityContainer and sets the default values.
func NewIdentityContainer()(*IdentityContainer) {
    m := &IdentityContainer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityContainer(), nil
}
// GetApiConnectors gets the apiConnectors property value. Represents entry point for API connectors.
// returns a []IdentityApiConnectorable when successful
func (m *IdentityContainer) GetApiConnectors()([]IdentityApiConnectorable) {
    val, err := m.GetBackingStore().Get("apiConnectors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityApiConnectorable)
    }
    return nil
}
// GetAuthenticationEventListeners gets the authenticationEventListeners property value. Represents listeners for custom authentication extension events in Azure AD for workforce and customers.
// returns a []AuthenticationEventListenerable when successful
func (m *IdentityContainer) GetAuthenticationEventListeners()([]AuthenticationEventListenerable) {
    val, err := m.GetBackingStore().Get("authenticationEventListeners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationEventListenerable)
    }
    return nil
}
// GetAuthenticationEventsFlows gets the authenticationEventsFlows property value. Represents the entry point for self-service sign-up and sign-in user flows in both Microsoft Entra workforce and external tenants.
// returns a []AuthenticationEventsFlowable when successful
func (m *IdentityContainer) GetAuthenticationEventsFlows()([]AuthenticationEventsFlowable) {
    val, err := m.GetBackingStore().Get("authenticationEventsFlows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationEventsFlowable)
    }
    return nil
}
// GetB2xUserFlows gets the b2xUserFlows property value. Represents entry point for B2X/self-service sign-up identity userflows.
// returns a []B2xIdentityUserFlowable when successful
func (m *IdentityContainer) GetB2xUserFlows()([]B2xIdentityUserFlowable) {
    val, err := m.GetBackingStore().Get("b2xUserFlows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]B2xIdentityUserFlowable)
    }
    return nil
}
// GetConditionalAccess gets the conditionalAccess property value. the entry point for the Conditional Access (CA) object model.
// returns a ConditionalAccessRootable when successful
func (m *IdentityContainer) GetConditionalAccess()(ConditionalAccessRootable) {
    val, err := m.GetBackingStore().Get("conditionalAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessRootable)
    }
    return nil
}
// GetCustomAuthenticationExtensions gets the customAuthenticationExtensions property value. Represents custom extensions to authentication flows in Azure AD for workforce and customers.
// returns a []CustomAuthenticationExtensionable when successful
func (m *IdentityContainer) GetCustomAuthenticationExtensions()([]CustomAuthenticationExtensionable) {
    val, err := m.GetBackingStore().Get("customAuthenticationExtensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CustomAuthenticationExtensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IdentityContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["apiConnectors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityApiConnectorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityApiConnectorable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityApiConnectorable)
                }
            }
            m.SetApiConnectors(res)
        }
        return nil
    }
    res["authenticationEventListeners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationEventListenerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationEventListenerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationEventListenerable)
                }
            }
            m.SetAuthenticationEventListeners(res)
        }
        return nil
    }
    res["authenticationEventsFlows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationEventsFlowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationEventsFlowable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationEventsFlowable)
                }
            }
            m.SetAuthenticationEventsFlows(res)
        }
        return nil
    }
    res["b2xUserFlows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateB2xIdentityUserFlowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]B2xIdentityUserFlowable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(B2xIdentityUserFlowable)
                }
            }
            m.SetB2xUserFlows(res)
        }
        return nil
    }
    res["conditionalAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditionalAccess(val.(ConditionalAccessRootable))
        }
        return nil
    }
    res["customAuthenticationExtensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCustomAuthenticationExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CustomAuthenticationExtensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CustomAuthenticationExtensionable)
                }
            }
            m.SetCustomAuthenticationExtensions(res)
        }
        return nil
    }
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
    res["userFlowAttributes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityUserFlowAttributeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityUserFlowAttributeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityUserFlowAttributeable)
                }
            }
            m.SetUserFlowAttributes(res)
        }
        return nil
    }
    return res
}
// GetIdentityProviders gets the identityProviders property value. The identityProviders property
// returns a []IdentityProviderBaseable when successful
func (m *IdentityContainer) GetIdentityProviders()([]IdentityProviderBaseable) {
    val, err := m.GetBackingStore().Get("identityProviders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityProviderBaseable)
    }
    return nil
}
// GetUserFlowAttributes gets the userFlowAttributes property value. Represents entry point for identity userflow attributes.
// returns a []IdentityUserFlowAttributeable when successful
func (m *IdentityContainer) GetUserFlowAttributes()([]IdentityUserFlowAttributeable) {
    val, err := m.GetBackingStore().Get("userFlowAttributes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityUserFlowAttributeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApiConnectors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetApiConnectors()))
        for i, v := range m.GetApiConnectors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("apiConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAuthenticationEventListeners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthenticationEventListeners()))
        for i, v := range m.GetAuthenticationEventListeners() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("authenticationEventListeners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAuthenticationEventsFlows() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthenticationEventsFlows()))
        for i, v := range m.GetAuthenticationEventsFlows() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("authenticationEventsFlows", cast)
        if err != nil {
            return err
        }
    }
    if m.GetB2xUserFlows() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetB2xUserFlows()))
        for i, v := range m.GetB2xUserFlows() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("b2xUserFlows", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditionalAccess", m.GetConditionalAccess())
        if err != nil {
            return err
        }
    }
    if m.GetCustomAuthenticationExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomAuthenticationExtensions()))
        for i, v := range m.GetCustomAuthenticationExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customAuthenticationExtensions", cast)
        if err != nil {
            return err
        }
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
    if m.GetUserFlowAttributes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserFlowAttributes()))
        for i, v := range m.GetUserFlowAttributes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userFlowAttributes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApiConnectors sets the apiConnectors property value. Represents entry point for API connectors.
func (m *IdentityContainer) SetApiConnectors(value []IdentityApiConnectorable)() {
    err := m.GetBackingStore().Set("apiConnectors", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthenticationEventListeners sets the authenticationEventListeners property value. Represents listeners for custom authentication extension events in Azure AD for workforce and customers.
func (m *IdentityContainer) SetAuthenticationEventListeners(value []AuthenticationEventListenerable)() {
    err := m.GetBackingStore().Set("authenticationEventListeners", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthenticationEventsFlows sets the authenticationEventsFlows property value. Represents the entry point for self-service sign-up and sign-in user flows in both Microsoft Entra workforce and external tenants.
func (m *IdentityContainer) SetAuthenticationEventsFlows(value []AuthenticationEventsFlowable)() {
    err := m.GetBackingStore().Set("authenticationEventsFlows", value)
    if err != nil {
        panic(err)
    }
}
// SetB2xUserFlows sets the b2xUserFlows property value. Represents entry point for B2X/self-service sign-up identity userflows.
func (m *IdentityContainer) SetB2xUserFlows(value []B2xIdentityUserFlowable)() {
    err := m.GetBackingStore().Set("b2xUserFlows", value)
    if err != nil {
        panic(err)
    }
}
// SetConditionalAccess sets the conditionalAccess property value. the entry point for the Conditional Access (CA) object model.
func (m *IdentityContainer) SetConditionalAccess(value ConditionalAccessRootable)() {
    err := m.GetBackingStore().Set("conditionalAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomAuthenticationExtensions sets the customAuthenticationExtensions property value. Represents custom extensions to authentication flows in Azure AD for workforce and customers.
func (m *IdentityContainer) SetCustomAuthenticationExtensions(value []CustomAuthenticationExtensionable)() {
    err := m.GetBackingStore().Set("customAuthenticationExtensions", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityProviders sets the identityProviders property value. The identityProviders property
func (m *IdentityContainer) SetIdentityProviders(value []IdentityProviderBaseable)() {
    err := m.GetBackingStore().Set("identityProviders", value)
    if err != nil {
        panic(err)
    }
}
// SetUserFlowAttributes sets the userFlowAttributes property value. Represents entry point for identity userflow attributes.
func (m *IdentityContainer) SetUserFlowAttributes(value []IdentityUserFlowAttributeable)() {
    err := m.GetBackingStore().Set("userFlowAttributes", value)
    if err != nil {
        panic(err)
    }
}
type IdentityContainerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApiConnectors()([]IdentityApiConnectorable)
    GetAuthenticationEventListeners()([]AuthenticationEventListenerable)
    GetAuthenticationEventsFlows()([]AuthenticationEventsFlowable)
    GetB2xUserFlows()([]B2xIdentityUserFlowable)
    GetConditionalAccess()(ConditionalAccessRootable)
    GetCustomAuthenticationExtensions()([]CustomAuthenticationExtensionable)
    GetIdentityProviders()([]IdentityProviderBaseable)
    GetUserFlowAttributes()([]IdentityUserFlowAttributeable)
    SetApiConnectors(value []IdentityApiConnectorable)()
    SetAuthenticationEventListeners(value []AuthenticationEventListenerable)()
    SetAuthenticationEventsFlows(value []AuthenticationEventsFlowable)()
    SetB2xUserFlows(value []B2xIdentityUserFlowable)()
    SetConditionalAccess(value ConditionalAccessRootable)()
    SetCustomAuthenticationExtensions(value []CustomAuthenticationExtensionable)()
    SetIdentityProviders(value []IdentityProviderBaseable)()
    SetUserFlowAttributes(value []IdentityUserFlowAttributeable)()
}

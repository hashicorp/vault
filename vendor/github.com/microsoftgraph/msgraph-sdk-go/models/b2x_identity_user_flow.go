package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type B2xIdentityUserFlow struct {
    IdentityUserFlow
}
// NewB2xIdentityUserFlow instantiates a new B2xIdentityUserFlow and sets the default values.
func NewB2xIdentityUserFlow()(*B2xIdentityUserFlow) {
    m := &B2xIdentityUserFlow{
        IdentityUserFlow: *NewIdentityUserFlow(),
    }
    return m
}
// CreateB2xIdentityUserFlowFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateB2xIdentityUserFlowFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewB2xIdentityUserFlow(), nil
}
// GetApiConnectorConfiguration gets the apiConnectorConfiguration property value. Configuration for enabling an API connector for use as part of the self-service sign-up user flow. You can only obtain the value of this object using Get userFlowApiConnectorConfiguration.
// returns a UserFlowApiConnectorConfigurationable when successful
func (m *B2xIdentityUserFlow) GetApiConnectorConfiguration()(UserFlowApiConnectorConfigurationable) {
    val, err := m.GetBackingStore().Get("apiConnectorConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserFlowApiConnectorConfigurationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *B2xIdentityUserFlow) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentityUserFlow.GetFieldDeserializers()
    res["apiConnectorConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserFlowApiConnectorConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApiConnectorConfiguration(val.(UserFlowApiConnectorConfigurationable))
        }
        return nil
    }
    res["identityProviders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityProviderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityProviderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityProviderable)
                }
            }
            m.SetIdentityProviders(res)
        }
        return nil
    }
    res["languages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserFlowLanguageConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserFlowLanguageConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserFlowLanguageConfigurationable)
                }
            }
            m.SetLanguages(res)
        }
        return nil
    }
    res["userAttributeAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityUserFlowAttributeAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentityUserFlowAttributeAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentityUserFlowAttributeAssignmentable)
                }
            }
            m.SetUserAttributeAssignments(res)
        }
        return nil
    }
    res["userFlowIdentityProviders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetUserFlowIdentityProviders(res)
        }
        return nil
    }
    return res
}
// GetIdentityProviders gets the identityProviders property value. The identity providers included in the user flow.
// returns a []IdentityProviderable when successful
func (m *B2xIdentityUserFlow) GetIdentityProviders()([]IdentityProviderable) {
    val, err := m.GetBackingStore().Get("identityProviders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityProviderable)
    }
    return nil
}
// GetLanguages gets the languages property value. The languages supported for customization within the user flow. Language customization is enabled by default in self-service sign-up user flow. You can't create custom languages in self-service sign-up user flows.
// returns a []UserFlowLanguageConfigurationable when successful
func (m *B2xIdentityUserFlow) GetLanguages()([]UserFlowLanguageConfigurationable) {
    val, err := m.GetBackingStore().Get("languages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserFlowLanguageConfigurationable)
    }
    return nil
}
// GetUserAttributeAssignments gets the userAttributeAssignments property value. The user attribute assignments included in the user flow.
// returns a []IdentityUserFlowAttributeAssignmentable when successful
func (m *B2xIdentityUserFlow) GetUserAttributeAssignments()([]IdentityUserFlowAttributeAssignmentable) {
    val, err := m.GetBackingStore().Get("userAttributeAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityUserFlowAttributeAssignmentable)
    }
    return nil
}
// GetUserFlowIdentityProviders gets the userFlowIdentityProviders property value. The userFlowIdentityProviders property
// returns a []IdentityProviderBaseable when successful
func (m *B2xIdentityUserFlow) GetUserFlowIdentityProviders()([]IdentityProviderBaseable) {
    val, err := m.GetBackingStore().Get("userFlowIdentityProviders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentityProviderBaseable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *B2xIdentityUserFlow) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentityUserFlow.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("apiConnectorConfiguration", m.GetApiConnectorConfiguration())
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
    if m.GetLanguages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLanguages()))
        for i, v := range m.GetLanguages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("languages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserAttributeAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserAttributeAssignments()))
        for i, v := range m.GetUserAttributeAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userAttributeAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserFlowIdentityProviders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserFlowIdentityProviders()))
        for i, v := range m.GetUserFlowIdentityProviders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userFlowIdentityProviders", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApiConnectorConfiguration sets the apiConnectorConfiguration property value. Configuration for enabling an API connector for use as part of the self-service sign-up user flow. You can only obtain the value of this object using Get userFlowApiConnectorConfiguration.
func (m *B2xIdentityUserFlow) SetApiConnectorConfiguration(value UserFlowApiConnectorConfigurationable)() {
    err := m.GetBackingStore().Set("apiConnectorConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityProviders sets the identityProviders property value. The identity providers included in the user flow.
func (m *B2xIdentityUserFlow) SetIdentityProviders(value []IdentityProviderable)() {
    err := m.GetBackingStore().Set("identityProviders", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguages sets the languages property value. The languages supported for customization within the user flow. Language customization is enabled by default in self-service sign-up user flow. You can't create custom languages in self-service sign-up user flows.
func (m *B2xIdentityUserFlow) SetLanguages(value []UserFlowLanguageConfigurationable)() {
    err := m.GetBackingStore().Set("languages", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAttributeAssignments sets the userAttributeAssignments property value. The user attribute assignments included in the user flow.
func (m *B2xIdentityUserFlow) SetUserAttributeAssignments(value []IdentityUserFlowAttributeAssignmentable)() {
    err := m.GetBackingStore().Set("userAttributeAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetUserFlowIdentityProviders sets the userFlowIdentityProviders property value. The userFlowIdentityProviders property
func (m *B2xIdentityUserFlow) SetUserFlowIdentityProviders(value []IdentityProviderBaseable)() {
    err := m.GetBackingStore().Set("userFlowIdentityProviders", value)
    if err != nil {
        panic(err)
    }
}
type B2xIdentityUserFlowable interface {
    IdentityUserFlowable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApiConnectorConfiguration()(UserFlowApiConnectorConfigurationable)
    GetIdentityProviders()([]IdentityProviderable)
    GetLanguages()([]UserFlowLanguageConfigurationable)
    GetUserAttributeAssignments()([]IdentityUserFlowAttributeAssignmentable)
    GetUserFlowIdentityProviders()([]IdentityProviderBaseable)
    SetApiConnectorConfiguration(value UserFlowApiConnectorConfigurationable)()
    SetIdentityProviders(value []IdentityProviderable)()
    SetLanguages(value []UserFlowLanguageConfigurationable)()
    SetUserAttributeAssignments(value []IdentityUserFlowAttributeAssignmentable)()
    SetUserFlowIdentityProviders(value []IdentityProviderBaseable)()
}

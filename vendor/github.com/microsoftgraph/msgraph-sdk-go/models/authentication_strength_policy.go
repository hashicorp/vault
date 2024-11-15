package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationStrengthPolicy struct {
    Entity
}
// NewAuthenticationStrengthPolicy instantiates a new AuthenticationStrengthPolicy and sets the default values.
func NewAuthenticationStrengthPolicy()(*AuthenticationStrengthPolicy) {
    m := &AuthenticationStrengthPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationStrengthPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationStrengthPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationStrengthPolicy(), nil
}
// GetAllowedCombinations gets the allowedCombinations property value. A collection of authentication method modes that are required be used to satify this authentication strength.
// returns a []AuthenticationMethodModes when successful
func (m *AuthenticationStrengthPolicy) GetAllowedCombinations()([]AuthenticationMethodModes) {
    val, err := m.GetBackingStore().Get("allowedCombinations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodModes)
    }
    return nil
}
// GetCombinationConfigurations gets the combinationConfigurations property value. Settings that may be used to require specific types or instances of an authentication method to be used when authenticating with a specified combination of authentication methods.
// returns a []AuthenticationCombinationConfigurationable when successful
func (m *AuthenticationStrengthPolicy) GetCombinationConfigurations()([]AuthenticationCombinationConfigurationable) {
    val, err := m.GetBackingStore().Get("combinationConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationCombinationConfigurationable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The datetime when this policy was created.
// returns a *Time when successful
func (m *AuthenticationStrengthPolicy) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. The human-readable description of this policy.
// returns a *string when successful
func (m *AuthenticationStrengthPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The human-readable display name of this policy. Supports $filter (eq, ne, not , and in).
// returns a *string when successful
func (m *AuthenticationStrengthPolicy) GetDisplayName()(*string) {
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
func (m *AuthenticationStrengthPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowedCombinations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAllowedCombinations(res)
        }
        return nil
    }
    res["combinationConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationCombinationConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationCombinationConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationCombinationConfigurationable)
                }
            }
            m.SetCombinationConfigurations(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
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
    res["modifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModifiedDateTime(val)
        }
        return nil
    }
    res["policyType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationStrengthPolicyType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyType(val.(*AuthenticationStrengthPolicyType))
        }
        return nil
    }
    res["requirementsSatisfied"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationStrengthRequirements)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequirementsSatisfied(val.(*AuthenticationStrengthRequirements))
        }
        return nil
    }
    return res
}
// GetModifiedDateTime gets the modifiedDateTime property value. The datetime when this policy was last modified.
// returns a *Time when successful
func (m *AuthenticationStrengthPolicy) GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("modifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPolicyType gets the policyType property value. The policyType property
// returns a *AuthenticationStrengthPolicyType when successful
func (m *AuthenticationStrengthPolicy) GetPolicyType()(*AuthenticationStrengthPolicyType) {
    val, err := m.GetBackingStore().Get("policyType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationStrengthPolicyType)
    }
    return nil
}
// GetRequirementsSatisfied gets the requirementsSatisfied property value. The requirementsSatisfied property
// returns a *AuthenticationStrengthRequirements when successful
func (m *AuthenticationStrengthPolicy) GetRequirementsSatisfied()(*AuthenticationStrengthRequirements) {
    val, err := m.GetBackingStore().Get("requirementsSatisfied")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationStrengthRequirements)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationStrengthPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedCombinations() != nil {
        err = writer.WriteCollectionOfStringValues("allowedCombinations", SerializeAuthenticationMethodModes(m.GetAllowedCombinations()))
        if err != nil {
            return err
        }
    }
    if m.GetCombinationConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCombinationConfigurations()))
        for i, v := range m.GetCombinationConfigurations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("combinationConfigurations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
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
        err = writer.WriteTimeValue("modifiedDateTime", m.GetModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetPolicyType() != nil {
        cast := (*m.GetPolicyType()).String()
        err = writer.WriteStringValue("policyType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRequirementsSatisfied() != nil {
        cast := (*m.GetRequirementsSatisfied()).String()
        err = writer.WriteStringValue("requirementsSatisfied", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedCombinations sets the allowedCombinations property value. A collection of authentication method modes that are required be used to satify this authentication strength.
func (m *AuthenticationStrengthPolicy) SetAllowedCombinations(value []AuthenticationMethodModes)() {
    err := m.GetBackingStore().Set("allowedCombinations", value)
    if err != nil {
        panic(err)
    }
}
// SetCombinationConfigurations sets the combinationConfigurations property value. Settings that may be used to require specific types or instances of an authentication method to be used when authenticating with a specified combination of authentication methods.
func (m *AuthenticationStrengthPolicy) SetCombinationConfigurations(value []AuthenticationCombinationConfigurationable)() {
    err := m.GetBackingStore().Set("combinationConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The datetime when this policy was created.
func (m *AuthenticationStrengthPolicy) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The human-readable description of this policy.
func (m *AuthenticationStrengthPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The human-readable display name of this policy. Supports $filter (eq, ne, not , and in).
func (m *AuthenticationStrengthPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetModifiedDateTime sets the modifiedDateTime property value. The datetime when this policy was last modified.
func (m *AuthenticationStrengthPolicy) SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("modifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyType sets the policyType property value. The policyType property
func (m *AuthenticationStrengthPolicy) SetPolicyType(value *AuthenticationStrengthPolicyType)() {
    err := m.GetBackingStore().Set("policyType", value)
    if err != nil {
        panic(err)
    }
}
// SetRequirementsSatisfied sets the requirementsSatisfied property value. The requirementsSatisfied property
func (m *AuthenticationStrengthPolicy) SetRequirementsSatisfied(value *AuthenticationStrengthRequirements)() {
    err := m.GetBackingStore().Set("requirementsSatisfied", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationStrengthPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedCombinations()([]AuthenticationMethodModes)
    GetCombinationConfigurations()([]AuthenticationCombinationConfigurationable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPolicyType()(*AuthenticationStrengthPolicyType)
    GetRequirementsSatisfied()(*AuthenticationStrengthRequirements)
    SetAllowedCombinations(value []AuthenticationMethodModes)()
    SetCombinationConfigurations(value []AuthenticationCombinationConfigurationable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPolicyType(value *AuthenticationStrengthPolicyType)()
    SetRequirementsSatisfied(value *AuthenticationStrengthRequirements)()
}

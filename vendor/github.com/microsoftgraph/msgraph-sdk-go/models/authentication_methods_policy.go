package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthenticationMethodsPolicy struct {
    Entity
}
// NewAuthenticationMethodsPolicy instantiates a new AuthenticationMethodsPolicy and sets the default values.
func NewAuthenticationMethodsPolicy()(*AuthenticationMethodsPolicy) {
    m := &AuthenticationMethodsPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationMethodsPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthenticationMethodsPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationMethodsPolicy(), nil
}
// GetAuthenticationMethodConfigurations gets the authenticationMethodConfigurations property value. Represents the settings for each authentication method. Automatically expanded on GET /policies/authenticationMethodsPolicy.
// returns a []AuthenticationMethodConfigurationable when successful
func (m *AuthenticationMethodsPolicy) GetAuthenticationMethodConfigurations()([]AuthenticationMethodConfigurationable) {
    val, err := m.GetBackingStore().Get("authenticationMethodConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationMethodConfigurationable)
    }
    return nil
}
// GetDescription gets the description property value. A description of the policy. Read-only.
// returns a *string when successful
func (m *AuthenticationMethodsPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the policy. Read-only.
// returns a *string when successful
func (m *AuthenticationMethodsPolicy) GetDisplayName()(*string) {
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
func (m *AuthenticationMethodsPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationMethodConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationMethodConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationMethodConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AuthenticationMethodConfigurationable)
                }
            }
            m.SetAuthenticationMethodConfigurations(res)
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
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["policyMigrationState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAuthenticationMethodsPolicyMigrationState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyMigrationState(val.(*AuthenticationMethodsPolicyMigrationState))
        }
        return nil
    }
    res["policyVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyVersion(val)
        }
        return nil
    }
    res["reconfirmationInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReconfirmationInDays(val)
        }
        return nil
    }
    res["registrationEnforcement"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRegistrationEnforcementFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistrationEnforcement(val.(RegistrationEnforcementable))
        }
        return nil
    }
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time of the last update to the policy. Read-only.
// returns a *Time when successful
func (m *AuthenticationMethodsPolicy) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPolicyMigrationState gets the policyMigrationState property value. The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. The possible values are: premigration - means the authentication methods policy is used for authentication only, legacy policies are respected. migrationInProgress - means the authentication methods policy is used for both authentication and SSPR, legacy policies are respected. migrationComplete - means the authentication methods policy is used for authentication and SSPR, legacy policies are ignored. unknownFutureValue - Evolvable enumeration sentinel value. Do not use.
// returns a *AuthenticationMethodsPolicyMigrationState when successful
func (m *AuthenticationMethodsPolicy) GetPolicyMigrationState()(*AuthenticationMethodsPolicyMigrationState) {
    val, err := m.GetBackingStore().Get("policyMigrationState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AuthenticationMethodsPolicyMigrationState)
    }
    return nil
}
// GetPolicyVersion gets the policyVersion property value. The version of the policy in use. Read-only.
// returns a *string when successful
func (m *AuthenticationMethodsPolicy) GetPolicyVersion()(*string) {
    val, err := m.GetBackingStore().Get("policyVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReconfirmationInDays gets the reconfirmationInDays property value. The reconfirmationInDays property
// returns a *int32 when successful
func (m *AuthenticationMethodsPolicy) GetReconfirmationInDays()(*int32) {
    val, err := m.GetBackingStore().Get("reconfirmationInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRegistrationEnforcement gets the registrationEnforcement property value. Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.
// returns a RegistrationEnforcementable when successful
func (m *AuthenticationMethodsPolicy) GetRegistrationEnforcement()(RegistrationEnforcementable) {
    val, err := m.GetBackingStore().Get("registrationEnforcement")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RegistrationEnforcementable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthenticationMethodsPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationMethodConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthenticationMethodConfigurations()))
        for i, v := range m.GetAuthenticationMethodConfigurations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("authenticationMethodConfigurations", cast)
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
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetPolicyMigrationState() != nil {
        cast := (*m.GetPolicyMigrationState()).String()
        err = writer.WriteStringValue("policyMigrationState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("policyVersion", m.GetPolicyVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("reconfirmationInDays", m.GetReconfirmationInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("registrationEnforcement", m.GetRegistrationEnforcement())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationMethodConfigurations sets the authenticationMethodConfigurations property value. Represents the settings for each authentication method. Automatically expanded on GET /policies/authenticationMethodsPolicy.
func (m *AuthenticationMethodsPolicy) SetAuthenticationMethodConfigurations(value []AuthenticationMethodConfigurationable)() {
    err := m.GetBackingStore().Set("authenticationMethodConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. A description of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time of the last update to the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyMigrationState sets the policyMigrationState property value. The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. The possible values are: premigration - means the authentication methods policy is used for authentication only, legacy policies are respected. migrationInProgress - means the authentication methods policy is used for both authentication and SSPR, legacy policies are respected. migrationComplete - means the authentication methods policy is used for authentication and SSPR, legacy policies are ignored. unknownFutureValue - Evolvable enumeration sentinel value. Do not use.
func (m *AuthenticationMethodsPolicy) SetPolicyMigrationState(value *AuthenticationMethodsPolicyMigrationState)() {
    err := m.GetBackingStore().Set("policyMigrationState", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyVersion sets the policyVersion property value. The version of the policy in use. Read-only.
func (m *AuthenticationMethodsPolicy) SetPolicyVersion(value *string)() {
    err := m.GetBackingStore().Set("policyVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetReconfirmationInDays sets the reconfirmationInDays property value. The reconfirmationInDays property
func (m *AuthenticationMethodsPolicy) SetReconfirmationInDays(value *int32)() {
    err := m.GetBackingStore().Set("reconfirmationInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistrationEnforcement sets the registrationEnforcement property value. Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.
func (m *AuthenticationMethodsPolicy) SetRegistrationEnforcement(value RegistrationEnforcementable)() {
    err := m.GetBackingStore().Set("registrationEnforcement", value)
    if err != nil {
        panic(err)
    }
}
type AuthenticationMethodsPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationMethodConfigurations()([]AuthenticationMethodConfigurationable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPolicyMigrationState()(*AuthenticationMethodsPolicyMigrationState)
    GetPolicyVersion()(*string)
    GetReconfirmationInDays()(*int32)
    GetRegistrationEnforcement()(RegistrationEnforcementable)
    SetAuthenticationMethodConfigurations(value []AuthenticationMethodConfigurationable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPolicyMigrationState(value *AuthenticationMethodsPolicyMigrationState)()
    SetPolicyVersion(value *string)()
    SetReconfirmationInDays(value *int32)()
    SetRegistrationEnforcement(value RegistrationEnforcementable)()
}

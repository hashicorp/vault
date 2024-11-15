package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuthorizationPolicy struct {
    PolicyBase
}
// NewAuthorizationPolicy instantiates a new AuthorizationPolicy and sets the default values.
func NewAuthorizationPolicy()(*AuthorizationPolicy) {
    m := &AuthorizationPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.authorizationPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAuthorizationPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuthorizationPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthorizationPolicy(), nil
}
// GetAllowedToSignUpEmailBasedSubscriptions gets the allowedToSignUpEmailBasedSubscriptions property value. Indicates whether users can sign up for email based subscriptions.
// returns a *bool when successful
func (m *AuthorizationPolicy) GetAllowedToSignUpEmailBasedSubscriptions()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToSignUpEmailBasedSubscriptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedToUseSSPR gets the allowedToUseSSPR property value. Indicates whether administrators of the tenant can use the Self-Service Password Reset (SSPR). For more information, see Self-service password reset for administrators.
// returns a *bool when successful
func (m *AuthorizationPolicy) GetAllowedToUseSSPR()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToUseSSPR")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowEmailVerifiedUsersToJoinOrganization gets the allowEmailVerifiedUsersToJoinOrganization property value. Indicates whether a user can join the tenant by email validation.
// returns a *bool when successful
func (m *AuthorizationPolicy) GetAllowEmailVerifiedUsersToJoinOrganization()(*bool) {
    val, err := m.GetBackingStore().Get("allowEmailVerifiedUsersToJoinOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowInvitesFrom gets the allowInvitesFrom property value. Indicates who can invite guests to the organization. Possible values are: none, adminsAndGuestInviters, adminsGuestInvitersAndAllMembers, everyone.  everyone is the default setting for all cloud environments except US Government. For more information, see allowInvitesFrom values.
// returns a *AllowInvitesFrom when successful
func (m *AuthorizationPolicy) GetAllowInvitesFrom()(*AllowInvitesFrom) {
    val, err := m.GetBackingStore().Get("allowInvitesFrom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AllowInvitesFrom)
    }
    return nil
}
// GetAllowUserConsentForRiskyApps gets the allowUserConsentForRiskyApps property value. Indicates whether user consent for risky apps is allowed. We recommend keeping allowUserConsentForRiskyApps as false. Default value is false.
// returns a *bool when successful
func (m *AuthorizationPolicy) GetAllowUserConsentForRiskyApps()(*bool) {
    val, err := m.GetBackingStore().Get("allowUserConsentForRiskyApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBlockMsolPowerShell gets the blockMsolPowerShell property value. To disable the use of MSOL PowerShell, set this property to true. This also disables user-based access to the legacy service endpoint used by MSOL PowerShell. This doesn't affect Microsoft Entra Connect or Microsoft Graph.
// returns a *bool when successful
func (m *AuthorizationPolicy) GetBlockMsolPowerShell()(*bool) {
    val, err := m.GetBackingStore().Get("blockMsolPowerShell")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefaultUserRolePermissions gets the defaultUserRolePermissions property value. The defaultUserRolePermissions property
// returns a DefaultUserRolePermissionsable when successful
func (m *AuthorizationPolicy) GetDefaultUserRolePermissions()(DefaultUserRolePermissionsable) {
    val, err := m.GetBackingStore().Get("defaultUserRolePermissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DefaultUserRolePermissionsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuthorizationPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["allowedToSignUpEmailBasedSubscriptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToSignUpEmailBasedSubscriptions(val)
        }
        return nil
    }
    res["allowedToUseSSPR"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToUseSSPR(val)
        }
        return nil
    }
    res["allowEmailVerifiedUsersToJoinOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowEmailVerifiedUsersToJoinOrganization(val)
        }
        return nil
    }
    res["allowInvitesFrom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAllowInvitesFrom)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowInvitesFrom(val.(*AllowInvitesFrom))
        }
        return nil
    }
    res["allowUserConsentForRiskyApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowUserConsentForRiskyApps(val)
        }
        return nil
    }
    res["blockMsolPowerShell"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlockMsolPowerShell(val)
        }
        return nil
    }
    res["defaultUserRolePermissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDefaultUserRolePermissionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultUserRolePermissions(val.(DefaultUserRolePermissionsable))
        }
        return nil
    }
    res["guestUserRoleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGuestUserRoleId(val)
        }
        return nil
    }
    return res
}
// GetGuestUserRoleId gets the guestUserRoleId property value. Represents role templateId for the role that should be granted to guests. Currently following roles are supported:  User (a0b1b346-4d3e-4e8b-98f8-753987be4970), Guest User (10dae51f-b6af-4016-8d66-8c2a99b929b3), and Restricted Guest User (2af84b1e-32c8-42b7-82bc-daa82404023b).
// returns a *UUID when successful
func (m *AuthorizationPolicy) GetGuestUserRoleId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("guestUserRoleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuthorizationPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowedToSignUpEmailBasedSubscriptions", m.GetAllowedToSignUpEmailBasedSubscriptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowedToUseSSPR", m.GetAllowedToUseSSPR())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowEmailVerifiedUsersToJoinOrganization", m.GetAllowEmailVerifiedUsersToJoinOrganization())
        if err != nil {
            return err
        }
    }
    if m.GetAllowInvitesFrom() != nil {
        cast := (*m.GetAllowInvitesFrom()).String()
        err = writer.WriteStringValue("allowInvitesFrom", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowUserConsentForRiskyApps", m.GetAllowUserConsentForRiskyApps())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("blockMsolPowerShell", m.GetBlockMsolPowerShell())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("defaultUserRolePermissions", m.GetDefaultUserRolePermissions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("guestUserRoleId", m.GetGuestUserRoleId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedToSignUpEmailBasedSubscriptions sets the allowedToSignUpEmailBasedSubscriptions property value. Indicates whether users can sign up for email based subscriptions.
func (m *AuthorizationPolicy) SetAllowedToSignUpEmailBasedSubscriptions(value *bool)() {
    err := m.GetBackingStore().Set("allowedToSignUpEmailBasedSubscriptions", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToUseSSPR sets the allowedToUseSSPR property value. Indicates whether administrators of the tenant can use the Self-Service Password Reset (SSPR). For more information, see Self-service password reset for administrators.
func (m *AuthorizationPolicy) SetAllowedToUseSSPR(value *bool)() {
    err := m.GetBackingStore().Set("allowedToUseSSPR", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowEmailVerifiedUsersToJoinOrganization sets the allowEmailVerifiedUsersToJoinOrganization property value. Indicates whether a user can join the tenant by email validation.
func (m *AuthorizationPolicy) SetAllowEmailVerifiedUsersToJoinOrganization(value *bool)() {
    err := m.GetBackingStore().Set("allowEmailVerifiedUsersToJoinOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowInvitesFrom sets the allowInvitesFrom property value. Indicates who can invite guests to the organization. Possible values are: none, adminsAndGuestInviters, adminsGuestInvitersAndAllMembers, everyone.  everyone is the default setting for all cloud environments except US Government. For more information, see allowInvitesFrom values.
func (m *AuthorizationPolicy) SetAllowInvitesFrom(value *AllowInvitesFrom)() {
    err := m.GetBackingStore().Set("allowInvitesFrom", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowUserConsentForRiskyApps sets the allowUserConsentForRiskyApps property value. Indicates whether user consent for risky apps is allowed. We recommend keeping allowUserConsentForRiskyApps as false. Default value is false.
func (m *AuthorizationPolicy) SetAllowUserConsentForRiskyApps(value *bool)() {
    err := m.GetBackingStore().Set("allowUserConsentForRiskyApps", value)
    if err != nil {
        panic(err)
    }
}
// SetBlockMsolPowerShell sets the blockMsolPowerShell property value. To disable the use of MSOL PowerShell, set this property to true. This also disables user-based access to the legacy service endpoint used by MSOL PowerShell. This doesn't affect Microsoft Entra Connect or Microsoft Graph.
func (m *AuthorizationPolicy) SetBlockMsolPowerShell(value *bool)() {
    err := m.GetBackingStore().Set("blockMsolPowerShell", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultUserRolePermissions sets the defaultUserRolePermissions property value. The defaultUserRolePermissions property
func (m *AuthorizationPolicy) SetDefaultUserRolePermissions(value DefaultUserRolePermissionsable)() {
    err := m.GetBackingStore().Set("defaultUserRolePermissions", value)
    if err != nil {
        panic(err)
    }
}
// SetGuestUserRoleId sets the guestUserRoleId property value. Represents role templateId for the role that should be granted to guests. Currently following roles are supported:  User (a0b1b346-4d3e-4e8b-98f8-753987be4970), Guest User (10dae51f-b6af-4016-8d66-8c2a99b929b3), and Restricted Guest User (2af84b1e-32c8-42b7-82bc-daa82404023b).
func (m *AuthorizationPolicy) SetGuestUserRoleId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("guestUserRoleId", value)
    if err != nil {
        panic(err)
    }
}
type AuthorizationPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PolicyBaseable
    GetAllowedToSignUpEmailBasedSubscriptions()(*bool)
    GetAllowedToUseSSPR()(*bool)
    GetAllowEmailVerifiedUsersToJoinOrganization()(*bool)
    GetAllowInvitesFrom()(*AllowInvitesFrom)
    GetAllowUserConsentForRiskyApps()(*bool)
    GetBlockMsolPowerShell()(*bool)
    GetDefaultUserRolePermissions()(DefaultUserRolePermissionsable)
    GetGuestUserRoleId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    SetAllowedToSignUpEmailBasedSubscriptions(value *bool)()
    SetAllowedToUseSSPR(value *bool)()
    SetAllowEmailVerifiedUsersToJoinOrganization(value *bool)()
    SetAllowInvitesFrom(value *AllowInvitesFrom)()
    SetAllowUserConsentForRiskyApps(value *bool)()
    SetBlockMsolPowerShell(value *bool)()
    SetDefaultUserRolePermissions(value DefaultUserRolePermissionsable)()
    SetGuestUserRoleId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CrossTenantAccessPolicyConfigurationDefault struct {
    Entity
}
// NewCrossTenantAccessPolicyConfigurationDefault instantiates a new CrossTenantAccessPolicyConfigurationDefault and sets the default values.
func NewCrossTenantAccessPolicyConfigurationDefault()(*CrossTenantAccessPolicyConfigurationDefault) {
    m := &CrossTenantAccessPolicyConfigurationDefault{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCrossTenantAccessPolicyConfigurationDefaultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCrossTenantAccessPolicyConfigurationDefaultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCrossTenantAccessPolicyConfigurationDefault(), nil
}
// GetAutomaticUserConsentSettings gets the automaticUserConsentSettings property value. Determines the default configuration for automatic user consent settings. The inboundAllowed and outboundAllowed properties are always false and can't be updated in the default configuration. Read-only.
// returns a InboundOutboundPolicyConfigurationable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetAutomaticUserConsentSettings()(InboundOutboundPolicyConfigurationable) {
    val, err := m.GetBackingStore().Get("automaticUserConsentSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InboundOutboundPolicyConfigurationable)
    }
    return nil
}
// GetB2bCollaborationInbound gets the b2bCollaborationInbound property value. Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetB2bCollaborationInbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bCollaborationInbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bCollaborationOutbound gets the b2bCollaborationOutbound property value. Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetB2bCollaborationOutbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bCollaborationOutbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bDirectConnectInbound gets the b2bDirectConnectInbound property value. Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B direct connect.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetB2bDirectConnectInbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bDirectConnectInbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bDirectConnectOutbound gets the b2bDirectConnectOutbound property value. Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetB2bDirectConnectOutbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bDirectConnectOutbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["automaticUserConsentSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInboundOutboundPolicyConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomaticUserConsentSettings(val.(InboundOutboundPolicyConfigurationable))
        }
        return nil
    }
    res["b2bCollaborationInbound"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyB2BSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetB2bCollaborationInbound(val.(CrossTenantAccessPolicyB2BSettingable))
        }
        return nil
    }
    res["b2bCollaborationOutbound"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyB2BSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetB2bCollaborationOutbound(val.(CrossTenantAccessPolicyB2BSettingable))
        }
        return nil
    }
    res["b2bDirectConnectInbound"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyB2BSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetB2bDirectConnectInbound(val.(CrossTenantAccessPolicyB2BSettingable))
        }
        return nil
    }
    res["b2bDirectConnectOutbound"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyB2BSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetB2bDirectConnectOutbound(val.(CrossTenantAccessPolicyB2BSettingable))
        }
        return nil
    }
    res["inboundTrust"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyInboundTrustFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInboundTrust(val.(CrossTenantAccessPolicyInboundTrustable))
        }
        return nil
    }
    res["invitationRedemptionIdentityProviderConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDefaultInvitationRedemptionIdentityProviderConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitationRedemptionIdentityProviderConfiguration(val.(DefaultInvitationRedemptionIdentityProviderConfigurationable))
        }
        return nil
    }
    res["isServiceDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsServiceDefault(val)
        }
        return nil
    }
    res["tenantRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantAccessPolicyTenantRestrictionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantRestrictions(val.(CrossTenantAccessPolicyTenantRestrictionsable))
        }
        return nil
    }
    return res
}
// GetInboundTrust gets the inboundTrust property value. Determines the default configuration for trusting other Conditional Access claims from external Microsoft Entra organizations.
// returns a CrossTenantAccessPolicyInboundTrustable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetInboundTrust()(CrossTenantAccessPolicyInboundTrustable) {
    val, err := m.GetBackingStore().Get("inboundTrust")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyInboundTrustable)
    }
    return nil
}
// GetInvitationRedemptionIdentityProviderConfiguration gets the invitationRedemptionIdentityProviderConfiguration property value. Defines the priority order based on which an identity provider is selected during invitation redemption for a guest user.
// returns a DefaultInvitationRedemptionIdentityProviderConfigurationable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetInvitationRedemptionIdentityProviderConfiguration()(DefaultInvitationRedemptionIdentityProviderConfigurationable) {
    val, err := m.GetBackingStore().Get("invitationRedemptionIdentityProviderConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DefaultInvitationRedemptionIdentityProviderConfigurationable)
    }
    return nil
}
// GetIsServiceDefault gets the isServiceDefault property value. If true, the default configuration is set to the system default configuration. If false, the default settings are customized.
// returns a *bool when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetIsServiceDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isServiceDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTenantRestrictions gets the tenantRestrictions property value. Defines the default tenant restrictions configuration for users in your organization who access an external organization on your network or devices.
// returns a CrossTenantAccessPolicyTenantRestrictionsable when successful
func (m *CrossTenantAccessPolicyConfigurationDefault) GetTenantRestrictions()(CrossTenantAccessPolicyTenantRestrictionsable) {
    val, err := m.GetBackingStore().Get("tenantRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyTenantRestrictionsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CrossTenantAccessPolicyConfigurationDefault) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("automaticUserConsentSettings", m.GetAutomaticUserConsentSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("b2bCollaborationInbound", m.GetB2bCollaborationInbound())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("b2bCollaborationOutbound", m.GetB2bCollaborationOutbound())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("b2bDirectConnectInbound", m.GetB2bDirectConnectInbound())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("b2bDirectConnectOutbound", m.GetB2bDirectConnectOutbound())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("inboundTrust", m.GetInboundTrust())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("invitationRedemptionIdentityProviderConfiguration", m.GetInvitationRedemptionIdentityProviderConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isServiceDefault", m.GetIsServiceDefault())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("tenantRestrictions", m.GetTenantRestrictions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAutomaticUserConsentSettings sets the automaticUserConsentSettings property value. Determines the default configuration for automatic user consent settings. The inboundAllowed and outboundAllowed properties are always false and can't be updated in the default configuration. Read-only.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetAutomaticUserConsentSettings(value InboundOutboundPolicyConfigurationable)() {
    err := m.GetBackingStore().Set("automaticUserConsentSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bCollaborationInbound sets the b2bCollaborationInbound property value. Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetB2bCollaborationInbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bCollaborationInbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bCollaborationOutbound sets the b2bCollaborationOutbound property value. Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetB2bCollaborationOutbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bCollaborationOutbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bDirectConnectInbound sets the b2bDirectConnectInbound property value. Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B direct connect.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetB2bDirectConnectInbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bDirectConnectInbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bDirectConnectOutbound sets the b2bDirectConnectOutbound property value. Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetB2bDirectConnectOutbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bDirectConnectOutbound", value)
    if err != nil {
        panic(err)
    }
}
// SetInboundTrust sets the inboundTrust property value. Determines the default configuration for trusting other Conditional Access claims from external Microsoft Entra organizations.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetInboundTrust(value CrossTenantAccessPolicyInboundTrustable)() {
    err := m.GetBackingStore().Set("inboundTrust", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitationRedemptionIdentityProviderConfiguration sets the invitationRedemptionIdentityProviderConfiguration property value. Defines the priority order based on which an identity provider is selected during invitation redemption for a guest user.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetInvitationRedemptionIdentityProviderConfiguration(value DefaultInvitationRedemptionIdentityProviderConfigurationable)() {
    err := m.GetBackingStore().Set("invitationRedemptionIdentityProviderConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetIsServiceDefault sets the isServiceDefault property value. If true, the default configuration is set to the system default configuration. If false, the default settings are customized.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetIsServiceDefault(value *bool)() {
    err := m.GetBackingStore().Set("isServiceDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantRestrictions sets the tenantRestrictions property value. Defines the default tenant restrictions configuration for users in your organization who access an external organization on your network or devices.
func (m *CrossTenantAccessPolicyConfigurationDefault) SetTenantRestrictions(value CrossTenantAccessPolicyTenantRestrictionsable)() {
    err := m.GetBackingStore().Set("tenantRestrictions", value)
    if err != nil {
        panic(err)
    }
}
type CrossTenantAccessPolicyConfigurationDefaultable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAutomaticUserConsentSettings()(InboundOutboundPolicyConfigurationable)
    GetB2bCollaborationInbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bCollaborationOutbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bDirectConnectInbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bDirectConnectOutbound()(CrossTenantAccessPolicyB2BSettingable)
    GetInboundTrust()(CrossTenantAccessPolicyInboundTrustable)
    GetInvitationRedemptionIdentityProviderConfiguration()(DefaultInvitationRedemptionIdentityProviderConfigurationable)
    GetIsServiceDefault()(*bool)
    GetTenantRestrictions()(CrossTenantAccessPolicyTenantRestrictionsable)
    SetAutomaticUserConsentSettings(value InboundOutboundPolicyConfigurationable)()
    SetB2bCollaborationInbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bCollaborationOutbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bDirectConnectInbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bDirectConnectOutbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetInboundTrust(value CrossTenantAccessPolicyInboundTrustable)()
    SetInvitationRedemptionIdentityProviderConfiguration(value DefaultInvitationRedemptionIdentityProviderConfigurationable)()
    SetIsServiceDefault(value *bool)()
    SetTenantRestrictions(value CrossTenantAccessPolicyTenantRestrictionsable)()
}

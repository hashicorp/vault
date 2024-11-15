package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CrossTenantAccessPolicyConfigurationPartner struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCrossTenantAccessPolicyConfigurationPartner instantiates a new CrossTenantAccessPolicyConfigurationPartner and sets the default values.
func NewCrossTenantAccessPolicyConfigurationPartner()(*CrossTenantAccessPolicyConfigurationPartner) {
    m := &CrossTenantAccessPolicyConfigurationPartner{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCrossTenantAccessPolicyConfigurationPartnerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCrossTenantAccessPolicyConfigurationPartnerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCrossTenantAccessPolicyConfigurationPartner(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAutomaticUserConsentSettings gets the automaticUserConsentSettings property value. Determines the partner-specific configuration for automatic user consent settings. Unless specifically configured, the inboundAllowed and outboundAllowed properties are null and inherit from the default settings, which is always false.
// returns a InboundOutboundPolicyConfigurationable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetAutomaticUserConsentSettings()(InboundOutboundPolicyConfigurationable) {
    val, err := m.GetBackingStore().Get("automaticUserConsentSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InboundOutboundPolicyConfigurationable)
    }
    return nil
}
// GetB2bCollaborationInbound gets the b2bCollaborationInbound property value. Defines your partner-specific configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetB2bCollaborationInbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bCollaborationInbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bCollaborationOutbound gets the b2bCollaborationOutbound property value. Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetB2bCollaborationOutbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bCollaborationOutbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bDirectConnectInbound gets the b2bDirectConnectInbound property value. Defines your partner-specific configuration for users from other organizations accessing your resources via Azure B2B direct connect.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetB2bDirectConnectInbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bDirectConnectInbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetB2bDirectConnectOutbound gets the b2bDirectConnectOutbound property value. Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect.
// returns a CrossTenantAccessPolicyB2BSettingable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetB2bDirectConnectOutbound()(CrossTenantAccessPolicyB2BSettingable) {
    val, err := m.GetBackingStore().Get("b2bDirectConnectOutbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyB2BSettingable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["identitySynchronization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantIdentitySyncPolicyPartnerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentitySynchronization(val.(CrossTenantIdentitySyncPolicyPartnerable))
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
    res["isInMultiTenantOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInMultiTenantOrganization(val)
        }
        return nil
    }
    res["isServiceProvider"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsServiceProvider(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
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
// GetIdentitySynchronization gets the identitySynchronization property value. Defines the cross-tenant policy for the synchronization of users from a partner tenant. Use this user synchronization policy to streamline collaboration between users in a multitenant organization by automating the creation, update, and deletion of users from one tenant to another.
// returns a CrossTenantIdentitySyncPolicyPartnerable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetIdentitySynchronization()(CrossTenantIdentitySyncPolicyPartnerable) {
    val, err := m.GetBackingStore().Get("identitySynchronization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantIdentitySyncPolicyPartnerable)
    }
    return nil
}
// GetInboundTrust gets the inboundTrust property value. Determines the partner-specific configuration for trusting other Conditional Access claims from external Microsoft Entra organizations.
// returns a CrossTenantAccessPolicyInboundTrustable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetInboundTrust()(CrossTenantAccessPolicyInboundTrustable) {
    val, err := m.GetBackingStore().Get("inboundTrust")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantAccessPolicyInboundTrustable)
    }
    return nil
}
// GetIsInMultiTenantOrganization gets the isInMultiTenantOrganization property value. Identifies whether a tenant is a member of a multitenant organization.
// returns a *bool when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetIsInMultiTenantOrganization()(*bool) {
    val, err := m.GetBackingStore().Get("isInMultiTenantOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsServiceProvider gets the isServiceProvider property value. Identifies whether the partner-specific configuration is a Cloud Service Provider for your organization.
// returns a *bool when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetIsServiceProvider()(*bool) {
    val, err := m.GetBackingStore().Get("isServiceProvider")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The tenant identifier for the partner Microsoft Entra organization. Read-only. Key.
// returns a *string when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTenantRestrictions gets the tenantRestrictions property value. Defines the partner-specific tenant restrictions configuration for users in your organization who access a partner organization using partner supplied identities on your network or devices.
// returns a CrossTenantAccessPolicyTenantRestrictionsable when successful
func (m *CrossTenantAccessPolicyConfigurationPartner) GetTenantRestrictions()(CrossTenantAccessPolicyTenantRestrictionsable) {
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
func (m *CrossTenantAccessPolicyConfigurationPartner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("automaticUserConsentSettings", m.GetAutomaticUserConsentSettings())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("b2bCollaborationInbound", m.GetB2bCollaborationInbound())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("b2bCollaborationOutbound", m.GetB2bCollaborationOutbound())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("b2bDirectConnectInbound", m.GetB2bDirectConnectInbound())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("b2bDirectConnectOutbound", m.GetB2bDirectConnectOutbound())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("identitySynchronization", m.GetIdentitySynchronization())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("inboundTrust", m.GetInboundTrust())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isInMultiTenantOrganization", m.GetIsInMultiTenantOrganization())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isServiceProvider", m.GetIsServiceProvider())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("tenantRestrictions", m.GetTenantRestrictions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomaticUserConsentSettings sets the automaticUserConsentSettings property value. Determines the partner-specific configuration for automatic user consent settings. Unless specifically configured, the inboundAllowed and outboundAllowed properties are null and inherit from the default settings, which is always false.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetAutomaticUserConsentSettings(value InboundOutboundPolicyConfigurationable)() {
    err := m.GetBackingStore().Set("automaticUserConsentSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bCollaborationInbound sets the b2bCollaborationInbound property value. Defines your partner-specific configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetB2bCollaborationInbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bCollaborationInbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bCollaborationOutbound sets the b2bCollaborationOutbound property value. Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetB2bCollaborationOutbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bCollaborationOutbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bDirectConnectInbound sets the b2bDirectConnectInbound property value. Defines your partner-specific configuration for users from other organizations accessing your resources via Azure B2B direct connect.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetB2bDirectConnectInbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bDirectConnectInbound", value)
    if err != nil {
        panic(err)
    }
}
// SetB2bDirectConnectOutbound sets the b2bDirectConnectOutbound property value. Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetB2bDirectConnectOutbound(value CrossTenantAccessPolicyB2BSettingable)() {
    err := m.GetBackingStore().Set("b2bDirectConnectOutbound", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIdentitySynchronization sets the identitySynchronization property value. Defines the cross-tenant policy for the synchronization of users from a partner tenant. Use this user synchronization policy to streamline collaboration between users in a multitenant organization by automating the creation, update, and deletion of users from one tenant to another.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetIdentitySynchronization(value CrossTenantIdentitySyncPolicyPartnerable)() {
    err := m.GetBackingStore().Set("identitySynchronization", value)
    if err != nil {
        panic(err)
    }
}
// SetInboundTrust sets the inboundTrust property value. Determines the partner-specific configuration for trusting other Conditional Access claims from external Microsoft Entra organizations.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetInboundTrust(value CrossTenantAccessPolicyInboundTrustable)() {
    err := m.GetBackingStore().Set("inboundTrust", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInMultiTenantOrganization sets the isInMultiTenantOrganization property value. Identifies whether a tenant is a member of a multitenant organization.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetIsInMultiTenantOrganization(value *bool)() {
    err := m.GetBackingStore().Set("isInMultiTenantOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetIsServiceProvider sets the isServiceProvider property value. Identifies whether the partner-specific configuration is a Cloud Service Provider for your organization.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetIsServiceProvider(value *bool)() {
    err := m.GetBackingStore().Set("isServiceProvider", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *CrossTenantAccessPolicyConfigurationPartner) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The tenant identifier for the partner Microsoft Entra organization. Read-only. Key.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantRestrictions sets the tenantRestrictions property value. Defines the partner-specific tenant restrictions configuration for users in your organization who access a partner organization using partner supplied identities on your network or devices.
func (m *CrossTenantAccessPolicyConfigurationPartner) SetTenantRestrictions(value CrossTenantAccessPolicyTenantRestrictionsable)() {
    err := m.GetBackingStore().Set("tenantRestrictions", value)
    if err != nil {
        panic(err)
    }
}
type CrossTenantAccessPolicyConfigurationPartnerable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAutomaticUserConsentSettings()(InboundOutboundPolicyConfigurationable)
    GetB2bCollaborationInbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bCollaborationOutbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bDirectConnectInbound()(CrossTenantAccessPolicyB2BSettingable)
    GetB2bDirectConnectOutbound()(CrossTenantAccessPolicyB2BSettingable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIdentitySynchronization()(CrossTenantIdentitySyncPolicyPartnerable)
    GetInboundTrust()(CrossTenantAccessPolicyInboundTrustable)
    GetIsInMultiTenantOrganization()(*bool)
    GetIsServiceProvider()(*bool)
    GetOdataType()(*string)
    GetTenantId()(*string)
    GetTenantRestrictions()(CrossTenantAccessPolicyTenantRestrictionsable)
    SetAutomaticUserConsentSettings(value InboundOutboundPolicyConfigurationable)()
    SetB2bCollaborationInbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bCollaborationOutbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bDirectConnectInbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetB2bDirectConnectOutbound(value CrossTenantAccessPolicyB2BSettingable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIdentitySynchronization(value CrossTenantIdentitySyncPolicyPartnerable)()
    SetInboundTrust(value CrossTenantAccessPolicyInboundTrustable)()
    SetIsInMultiTenantOrganization(value *bool)()
    SetIsServiceProvider(value *bool)()
    SetOdataType(value *string)()
    SetTenantId(value *string)()
    SetTenantRestrictions(value CrossTenantAccessPolicyTenantRestrictionsable)()
}

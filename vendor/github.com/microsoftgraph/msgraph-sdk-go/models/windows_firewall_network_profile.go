package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// WindowsFirewallNetworkProfile windows Firewall Profile Policies.
type WindowsFirewallNetworkProfile struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewWindowsFirewallNetworkProfile instantiates a new WindowsFirewallNetworkProfile and sets the default values.
func NewWindowsFirewallNetworkProfile()(*WindowsFirewallNetworkProfile) {
    m := &WindowsFirewallNetworkProfile{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateWindowsFirewallNetworkProfileFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsFirewallNetworkProfileFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsFirewallNetworkProfile(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *WindowsFirewallNetworkProfile) GetAdditionalData()(map[string]any) {
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
// GetAuthorizedApplicationRulesFromGroupPolicyMerged gets the authorizedApplicationRulesFromGroupPolicyMerged property value. Configures the firewall to merge authorized application rules from group policy with those from local store instead of ignoring the local store rules. When AuthorizedApplicationRulesFromGroupPolicyNotMerged and AuthorizedApplicationRulesFromGroupPolicyMerged are both true, AuthorizedApplicationRulesFromGroupPolicyMerged takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetAuthorizedApplicationRulesFromGroupPolicyMerged()(*bool) {
    val, err := m.GetBackingStore().Get("authorizedApplicationRulesFromGroupPolicyMerged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *WindowsFirewallNetworkProfile) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConnectionSecurityRulesFromGroupPolicyMerged gets the connectionSecurityRulesFromGroupPolicyMerged property value. Configures the firewall to merge connection security rules from group policy with those from local store instead of ignoring the local store rules. When ConnectionSecurityRulesFromGroupPolicyNotMerged and ConnectionSecurityRulesFromGroupPolicyMerged are both true, ConnectionSecurityRulesFromGroupPolicyMerged takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetConnectionSecurityRulesFromGroupPolicyMerged()(*bool) {
    val, err := m.GetBackingStore().Get("connectionSecurityRulesFromGroupPolicyMerged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsFirewallNetworkProfile) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["authorizedApplicationRulesFromGroupPolicyMerged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthorizedApplicationRulesFromGroupPolicyMerged(val)
        }
        return nil
    }
    res["connectionSecurityRulesFromGroupPolicyMerged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectionSecurityRulesFromGroupPolicyMerged(val)
        }
        return nil
    }
    res["firewallEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStateManagementSetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallEnabled(val.(*StateManagementSetting))
        }
        return nil
    }
    res["globalPortRulesFromGroupPolicyMerged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGlobalPortRulesFromGroupPolicyMerged(val)
        }
        return nil
    }
    res["inboundConnectionsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInboundConnectionsBlocked(val)
        }
        return nil
    }
    res["inboundNotificationsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInboundNotificationsBlocked(val)
        }
        return nil
    }
    res["incomingTrafficBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncomingTrafficBlocked(val)
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
    res["outboundConnectionsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOutboundConnectionsBlocked(val)
        }
        return nil
    }
    res["policyRulesFromGroupPolicyMerged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyRulesFromGroupPolicyMerged(val)
        }
        return nil
    }
    res["securedPacketExemptionAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecuredPacketExemptionAllowed(val)
        }
        return nil
    }
    res["stealthModeBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStealthModeBlocked(val)
        }
        return nil
    }
    res["unicastResponsesToMulticastBroadcastsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnicastResponsesToMulticastBroadcastsBlocked(val)
        }
        return nil
    }
    return res
}
// GetFirewallEnabled gets the firewallEnabled property value. State Management Setting.
// returns a *StateManagementSetting when successful
func (m *WindowsFirewallNetworkProfile) GetFirewallEnabled()(*StateManagementSetting) {
    val, err := m.GetBackingStore().Get("firewallEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*StateManagementSetting)
    }
    return nil
}
// GetGlobalPortRulesFromGroupPolicyMerged gets the globalPortRulesFromGroupPolicyMerged property value. Configures the firewall to merge global port rules from group policy with those from local store instead of ignoring the local store rules. When GlobalPortRulesFromGroupPolicyNotMerged and GlobalPortRulesFromGroupPolicyMerged are both true, GlobalPortRulesFromGroupPolicyMerged takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetGlobalPortRulesFromGroupPolicyMerged()(*bool) {
    val, err := m.GetBackingStore().Get("globalPortRulesFromGroupPolicyMerged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInboundConnectionsBlocked gets the inboundConnectionsBlocked property value. Configures the firewall to block all incoming connections by default. When InboundConnectionsRequired and InboundConnectionsBlocked are both true, InboundConnectionsBlocked takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetInboundConnectionsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("inboundConnectionsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInboundNotificationsBlocked gets the inboundNotificationsBlocked property value. Prevents the firewall from displaying notifications when an application is blocked from listening on a port. When InboundNotificationsRequired and InboundNotificationsBlocked are both true, InboundNotificationsBlocked takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetInboundNotificationsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("inboundNotificationsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIncomingTrafficBlocked gets the incomingTrafficBlocked property value. Configures the firewall to block all incoming traffic regardless of other policy settings. When IncomingTrafficRequired and IncomingTrafficBlocked are both true, IncomingTrafficBlocked takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetIncomingTrafficBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("incomingTrafficBlocked")
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
func (m *WindowsFirewallNetworkProfile) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOutboundConnectionsBlocked gets the outboundConnectionsBlocked property value. Configures the firewall to block all outgoing connections by default. When OutboundConnectionsRequired and OutboundConnectionsBlocked are both true, OutboundConnectionsBlocked takes priority. This setting will get applied to Windows releases version 1809 and above.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetOutboundConnectionsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("outboundConnectionsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPolicyRulesFromGroupPolicyMerged gets the policyRulesFromGroupPolicyMerged property value. Configures the firewall to merge Firewall Rule policies from group policy with those from local store instead of ignoring the local store rules. When PolicyRulesFromGroupPolicyNotMerged and PolicyRulesFromGroupPolicyMerged are both true, PolicyRulesFromGroupPolicyMerged takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetPolicyRulesFromGroupPolicyMerged()(*bool) {
    val, err := m.GetBackingStore().Get("policyRulesFromGroupPolicyMerged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSecuredPacketExemptionAllowed gets the securedPacketExemptionAllowed property value. Configures the firewall to allow the host computer to respond to unsolicited network traffic of that traffic is secured by IPSec even when stealthModeBlocked is set to true. When SecuredPacketExemptionBlocked and SecuredPacketExemptionAllowed are both true, SecuredPacketExemptionAllowed takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetSecuredPacketExemptionAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("securedPacketExemptionAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStealthModeBlocked gets the stealthModeBlocked property value. Prevent the server from operating in stealth mode. When StealthModeRequired and StealthModeBlocked are both true, StealthModeBlocked takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetStealthModeBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("stealthModeBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUnicastResponsesToMulticastBroadcastsBlocked gets the unicastResponsesToMulticastBroadcastsBlocked property value. Configures the firewall to block unicast responses to multicast broadcast traffic. When UnicastResponsesToMulticastBroadcastsRequired and UnicastResponsesToMulticastBroadcastsBlocked are both true, UnicastResponsesToMulticastBroadcastsBlocked takes priority.
// returns a *bool when successful
func (m *WindowsFirewallNetworkProfile) GetUnicastResponsesToMulticastBroadcastsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("unicastResponsesToMulticastBroadcastsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsFirewallNetworkProfile) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("authorizedApplicationRulesFromGroupPolicyMerged", m.GetAuthorizedApplicationRulesFromGroupPolicyMerged())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("connectionSecurityRulesFromGroupPolicyMerged", m.GetConnectionSecurityRulesFromGroupPolicyMerged())
        if err != nil {
            return err
        }
    }
    if m.GetFirewallEnabled() != nil {
        cast := (*m.GetFirewallEnabled()).String()
        err := writer.WriteStringValue("firewallEnabled", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("globalPortRulesFromGroupPolicyMerged", m.GetGlobalPortRulesFromGroupPolicyMerged())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("inboundConnectionsBlocked", m.GetInboundConnectionsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("inboundNotificationsBlocked", m.GetInboundNotificationsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("incomingTrafficBlocked", m.GetIncomingTrafficBlocked())
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
        err := writer.WriteBoolValue("outboundConnectionsBlocked", m.GetOutboundConnectionsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("policyRulesFromGroupPolicyMerged", m.GetPolicyRulesFromGroupPolicyMerged())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("securedPacketExemptionAllowed", m.GetSecuredPacketExemptionAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("stealthModeBlocked", m.GetStealthModeBlocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("unicastResponsesToMulticastBroadcastsBlocked", m.GetUnicastResponsesToMulticastBroadcastsBlocked())
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
func (m *WindowsFirewallNetworkProfile) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthorizedApplicationRulesFromGroupPolicyMerged sets the authorizedApplicationRulesFromGroupPolicyMerged property value. Configures the firewall to merge authorized application rules from group policy with those from local store instead of ignoring the local store rules. When AuthorizedApplicationRulesFromGroupPolicyNotMerged and AuthorizedApplicationRulesFromGroupPolicyMerged are both true, AuthorizedApplicationRulesFromGroupPolicyMerged takes priority.
func (m *WindowsFirewallNetworkProfile) SetAuthorizedApplicationRulesFromGroupPolicyMerged(value *bool)() {
    err := m.GetBackingStore().Set("authorizedApplicationRulesFromGroupPolicyMerged", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *WindowsFirewallNetworkProfile) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConnectionSecurityRulesFromGroupPolicyMerged sets the connectionSecurityRulesFromGroupPolicyMerged property value. Configures the firewall to merge connection security rules from group policy with those from local store instead of ignoring the local store rules. When ConnectionSecurityRulesFromGroupPolicyNotMerged and ConnectionSecurityRulesFromGroupPolicyMerged are both true, ConnectionSecurityRulesFromGroupPolicyMerged takes priority.
func (m *WindowsFirewallNetworkProfile) SetConnectionSecurityRulesFromGroupPolicyMerged(value *bool)() {
    err := m.GetBackingStore().Set("connectionSecurityRulesFromGroupPolicyMerged", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallEnabled sets the firewallEnabled property value. State Management Setting.
func (m *WindowsFirewallNetworkProfile) SetFirewallEnabled(value *StateManagementSetting)() {
    err := m.GetBackingStore().Set("firewallEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetGlobalPortRulesFromGroupPolicyMerged sets the globalPortRulesFromGroupPolicyMerged property value. Configures the firewall to merge global port rules from group policy with those from local store instead of ignoring the local store rules. When GlobalPortRulesFromGroupPolicyNotMerged and GlobalPortRulesFromGroupPolicyMerged are both true, GlobalPortRulesFromGroupPolicyMerged takes priority.
func (m *WindowsFirewallNetworkProfile) SetGlobalPortRulesFromGroupPolicyMerged(value *bool)() {
    err := m.GetBackingStore().Set("globalPortRulesFromGroupPolicyMerged", value)
    if err != nil {
        panic(err)
    }
}
// SetInboundConnectionsBlocked sets the inboundConnectionsBlocked property value. Configures the firewall to block all incoming connections by default. When InboundConnectionsRequired and InboundConnectionsBlocked are both true, InboundConnectionsBlocked takes priority.
func (m *WindowsFirewallNetworkProfile) SetInboundConnectionsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("inboundConnectionsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetInboundNotificationsBlocked sets the inboundNotificationsBlocked property value. Prevents the firewall from displaying notifications when an application is blocked from listening on a port. When InboundNotificationsRequired and InboundNotificationsBlocked are both true, InboundNotificationsBlocked takes priority.
func (m *WindowsFirewallNetworkProfile) SetInboundNotificationsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("inboundNotificationsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetIncomingTrafficBlocked sets the incomingTrafficBlocked property value. Configures the firewall to block all incoming traffic regardless of other policy settings. When IncomingTrafficRequired and IncomingTrafficBlocked are both true, IncomingTrafficBlocked takes priority.
func (m *WindowsFirewallNetworkProfile) SetIncomingTrafficBlocked(value *bool)() {
    err := m.GetBackingStore().Set("incomingTrafficBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WindowsFirewallNetworkProfile) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOutboundConnectionsBlocked sets the outboundConnectionsBlocked property value. Configures the firewall to block all outgoing connections by default. When OutboundConnectionsRequired and OutboundConnectionsBlocked are both true, OutboundConnectionsBlocked takes priority. This setting will get applied to Windows releases version 1809 and above.
func (m *WindowsFirewallNetworkProfile) SetOutboundConnectionsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("outboundConnectionsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyRulesFromGroupPolicyMerged sets the policyRulesFromGroupPolicyMerged property value. Configures the firewall to merge Firewall Rule policies from group policy with those from local store instead of ignoring the local store rules. When PolicyRulesFromGroupPolicyNotMerged and PolicyRulesFromGroupPolicyMerged are both true, PolicyRulesFromGroupPolicyMerged takes priority.
func (m *WindowsFirewallNetworkProfile) SetPolicyRulesFromGroupPolicyMerged(value *bool)() {
    err := m.GetBackingStore().Set("policyRulesFromGroupPolicyMerged", value)
    if err != nil {
        panic(err)
    }
}
// SetSecuredPacketExemptionAllowed sets the securedPacketExemptionAllowed property value. Configures the firewall to allow the host computer to respond to unsolicited network traffic of that traffic is secured by IPSec even when stealthModeBlocked is set to true. When SecuredPacketExemptionBlocked and SecuredPacketExemptionAllowed are both true, SecuredPacketExemptionAllowed takes priority.
func (m *WindowsFirewallNetworkProfile) SetSecuredPacketExemptionAllowed(value *bool)() {
    err := m.GetBackingStore().Set("securedPacketExemptionAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetStealthModeBlocked sets the stealthModeBlocked property value. Prevent the server from operating in stealth mode. When StealthModeRequired and StealthModeBlocked are both true, StealthModeBlocked takes priority.
func (m *WindowsFirewallNetworkProfile) SetStealthModeBlocked(value *bool)() {
    err := m.GetBackingStore().Set("stealthModeBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetUnicastResponsesToMulticastBroadcastsBlocked sets the unicastResponsesToMulticastBroadcastsBlocked property value. Configures the firewall to block unicast responses to multicast broadcast traffic. When UnicastResponsesToMulticastBroadcastsRequired and UnicastResponsesToMulticastBroadcastsBlocked are both true, UnicastResponsesToMulticastBroadcastsBlocked takes priority.
func (m *WindowsFirewallNetworkProfile) SetUnicastResponsesToMulticastBroadcastsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("unicastResponsesToMulticastBroadcastsBlocked", value)
    if err != nil {
        panic(err)
    }
}
type WindowsFirewallNetworkProfileable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthorizedApplicationRulesFromGroupPolicyMerged()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConnectionSecurityRulesFromGroupPolicyMerged()(*bool)
    GetFirewallEnabled()(*StateManagementSetting)
    GetGlobalPortRulesFromGroupPolicyMerged()(*bool)
    GetInboundConnectionsBlocked()(*bool)
    GetInboundNotificationsBlocked()(*bool)
    GetIncomingTrafficBlocked()(*bool)
    GetOdataType()(*string)
    GetOutboundConnectionsBlocked()(*bool)
    GetPolicyRulesFromGroupPolicyMerged()(*bool)
    GetSecuredPacketExemptionAllowed()(*bool)
    GetStealthModeBlocked()(*bool)
    GetUnicastResponsesToMulticastBroadcastsBlocked()(*bool)
    SetAuthorizedApplicationRulesFromGroupPolicyMerged(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConnectionSecurityRulesFromGroupPolicyMerged(value *bool)()
    SetFirewallEnabled(value *StateManagementSetting)()
    SetGlobalPortRulesFromGroupPolicyMerged(value *bool)()
    SetInboundConnectionsBlocked(value *bool)()
    SetInboundNotificationsBlocked(value *bool)()
    SetIncomingTrafficBlocked(value *bool)()
    SetOdataType(value *string)()
    SetOutboundConnectionsBlocked(value *bool)()
    SetPolicyRulesFromGroupPolicyMerged(value *bool)()
    SetSecuredPacketExemptionAllowed(value *bool)()
    SetStealthModeBlocked(value *bool)()
    SetUnicastResponsesToMulticastBroadcastsBlocked(value *bool)()
}

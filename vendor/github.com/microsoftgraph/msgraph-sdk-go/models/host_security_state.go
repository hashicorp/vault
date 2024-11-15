package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type HostSecurityState struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewHostSecurityState instantiates a new HostSecurityState and sets the default values.
func NewHostSecurityState()(*HostSecurityState) {
    m := &HostSecurityState{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateHostSecurityStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHostSecurityStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHostSecurityState(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *HostSecurityState) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *HostSecurityState) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *HostSecurityState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["fqdn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFqdn(val)
        }
        return nil
    }
    res["isAzureAdJoined"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAzureAdJoined(val)
        }
        return nil
    }
    res["isAzureAdRegistered"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAzureAdRegistered(val)
        }
        return nil
    }
    res["isHybridAzureDomainJoined"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsHybridAzureDomainJoined(val)
        }
        return nil
    }
    res["netBiosName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetBiosName(val)
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
    res["os"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOs(val)
        }
        return nil
    }
    res["privateIpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivateIpAddress(val)
        }
        return nil
    }
    res["publicIpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublicIpAddress(val)
        }
        return nil
    }
    res["riskScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskScore(val)
        }
        return nil
    }
    return res
}
// GetFqdn gets the fqdn property value. Host FQDN (Fully Qualified Domain Name) (for example, machine.company.com).
// returns a *string when successful
func (m *HostSecurityState) GetFqdn()(*string) {
    val, err := m.GetBackingStore().Get("fqdn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsAzureAdJoined gets the isAzureAdJoined property value. The isAzureAdJoined property
// returns a *bool when successful
func (m *HostSecurityState) GetIsAzureAdJoined()(*bool) {
    val, err := m.GetBackingStore().Get("isAzureAdJoined")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsAzureAdRegistered gets the isAzureAdRegistered property value. The isAzureAdRegistered property
// returns a *bool when successful
func (m *HostSecurityState) GetIsAzureAdRegistered()(*bool) {
    val, err := m.GetBackingStore().Get("isAzureAdRegistered")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsHybridAzureDomainJoined gets the isHybridAzureDomainJoined property value. True if the host is domain joined to an on-premises Active Directory domain.
// returns a *bool when successful
func (m *HostSecurityState) GetIsHybridAzureDomainJoined()(*bool) {
    val, err := m.GetBackingStore().Get("isHybridAzureDomainJoined")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetBiosName gets the netBiosName property value. The local host name, without the DNS domain name.
// returns a *string when successful
func (m *HostSecurityState) GetNetBiosName()(*string) {
    val, err := m.GetBackingStore().Get("netBiosName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *HostSecurityState) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOs gets the os property value. Host Operating System. (For example, Windows 10, macOS, RHEL, etc.).
// returns a *string when successful
func (m *HostSecurityState) GetOs()(*string) {
    val, err := m.GetBackingStore().Get("os")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrivateIpAddress gets the privateIpAddress property value. Private (not routable) IPv4 or IPv6 address (see RFC 1918) at the time of the alert.
// returns a *string when successful
func (m *HostSecurityState) GetPrivateIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("privateIpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublicIpAddress gets the publicIpAddress property value. Publicly routable IPv4 or IPv6 address (see RFC 1918) at time of the alert.
// returns a *string when successful
func (m *HostSecurityState) GetPublicIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("publicIpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskScore gets the riskScore property value. Provider-generated/calculated risk score of the host.  Recommended value range of 0-1, which equates to a percentage.
// returns a *string when successful
func (m *HostSecurityState) GetRiskScore()(*string) {
    val, err := m.GetBackingStore().Get("riskScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *HostSecurityState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("fqdn", m.GetFqdn())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAzureAdJoined", m.GetIsAzureAdJoined())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAzureAdRegistered", m.GetIsAzureAdRegistered())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isHybridAzureDomainJoined", m.GetIsHybridAzureDomainJoined())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("netBiosName", m.GetNetBiosName())
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
        err := writer.WriteStringValue("os", m.GetOs())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("privateIpAddress", m.GetPrivateIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("publicIpAddress", m.GetPublicIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("riskScore", m.GetRiskScore())
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
func (m *HostSecurityState) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *HostSecurityState) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetFqdn sets the fqdn property value. Host FQDN (Fully Qualified Domain Name) (for example, machine.company.com).
func (m *HostSecurityState) SetFqdn(value *string)() {
    err := m.GetBackingStore().Set("fqdn", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAzureAdJoined sets the isAzureAdJoined property value. The isAzureAdJoined property
func (m *HostSecurityState) SetIsAzureAdJoined(value *bool)() {
    err := m.GetBackingStore().Set("isAzureAdJoined", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAzureAdRegistered sets the isAzureAdRegistered property value. The isAzureAdRegistered property
func (m *HostSecurityState) SetIsAzureAdRegistered(value *bool)() {
    err := m.GetBackingStore().Set("isAzureAdRegistered", value)
    if err != nil {
        panic(err)
    }
}
// SetIsHybridAzureDomainJoined sets the isHybridAzureDomainJoined property value. True if the host is domain joined to an on-premises Active Directory domain.
func (m *HostSecurityState) SetIsHybridAzureDomainJoined(value *bool)() {
    err := m.GetBackingStore().Set("isHybridAzureDomainJoined", value)
    if err != nil {
        panic(err)
    }
}
// SetNetBiosName sets the netBiosName property value. The local host name, without the DNS domain name.
func (m *HostSecurityState) SetNetBiosName(value *string)() {
    err := m.GetBackingStore().Set("netBiosName", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *HostSecurityState) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOs sets the os property value. Host Operating System. (For example, Windows 10, macOS, RHEL, etc.).
func (m *HostSecurityState) SetOs(value *string)() {
    err := m.GetBackingStore().Set("os", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivateIpAddress sets the privateIpAddress property value. Private (not routable) IPv4 or IPv6 address (see RFC 1918) at the time of the alert.
func (m *HostSecurityState) SetPrivateIpAddress(value *string)() {
    err := m.GetBackingStore().Set("privateIpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetPublicIpAddress sets the publicIpAddress property value. Publicly routable IPv4 or IPv6 address (see RFC 1918) at time of the alert.
func (m *HostSecurityState) SetPublicIpAddress(value *string)() {
    err := m.GetBackingStore().Set("publicIpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskScore sets the riskScore property value. Provider-generated/calculated risk score of the host.  Recommended value range of 0-1, which equates to a percentage.
func (m *HostSecurityState) SetRiskScore(value *string)() {
    err := m.GetBackingStore().Set("riskScore", value)
    if err != nil {
        panic(err)
    }
}
type HostSecurityStateable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetFqdn()(*string)
    GetIsAzureAdJoined()(*bool)
    GetIsAzureAdRegistered()(*bool)
    GetIsHybridAzureDomainJoined()(*bool)
    GetNetBiosName()(*string)
    GetOdataType()(*string)
    GetOs()(*string)
    GetPrivateIpAddress()(*string)
    GetPublicIpAddress()(*string)
    GetRiskScore()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetFqdn(value *string)()
    SetIsAzureAdJoined(value *bool)()
    SetIsAzureAdRegistered(value *bool)()
    SetIsHybridAzureDomainJoined(value *bool)()
    SetNetBiosName(value *string)()
    SetOdataType(value *string)()
    SetOs(value *string)()
    SetPrivateIpAddress(value *string)()
    SetPublicIpAddress(value *string)()
    SetRiskScore(value *string)()
}

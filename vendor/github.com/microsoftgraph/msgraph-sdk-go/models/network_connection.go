package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type NetworkConnection struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewNetworkConnection instantiates a new NetworkConnection and sets the default values.
func NewNetworkConnection()(*NetworkConnection) {
    m := &NetworkConnection{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateNetworkConnectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateNetworkConnectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNetworkConnection(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *NetworkConnection) GetAdditionalData()(map[string]any) {
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
// GetApplicationName gets the applicationName property value. Name of the application managing the network connection (for example, Facebook or SMTP).
// returns a *string when successful
func (m *NetworkConnection) GetApplicationName()(*string) {
    val, err := m.GetBackingStore().Get("applicationName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *NetworkConnection) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDestinationAddress gets the destinationAddress property value. Destination IP address (of the network connection).
// returns a *string when successful
func (m *NetworkConnection) GetDestinationAddress()(*string) {
    val, err := m.GetBackingStore().Get("destinationAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDestinationDomain gets the destinationDomain property value. Destination domain portion of the destination URL. (for example 'www.contoso.com').
// returns a *string when successful
func (m *NetworkConnection) GetDestinationDomain()(*string) {
    val, err := m.GetBackingStore().Get("destinationDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDestinationLocation gets the destinationLocation property value. Location (by IP address mapping) associated with the destination of a network connection.
// returns a *string when successful
func (m *NetworkConnection) GetDestinationLocation()(*string) {
    val, err := m.GetBackingStore().Get("destinationLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDestinationPort gets the destinationPort property value. Destination port (of the network connection).
// returns a *string when successful
func (m *NetworkConnection) GetDestinationPort()(*string) {
    val, err := m.GetBackingStore().Get("destinationPort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDestinationUrl gets the destinationUrl property value. Network connection URL/URI string - excluding parameters. (for example 'www.contoso.com/products/default.html')
// returns a *string when successful
func (m *NetworkConnection) GetDestinationUrl()(*string) {
    val, err := m.GetBackingStore().Get("destinationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDirection gets the direction property value. Network connection direction. Possible values are: unknown, inbound, outbound.
// returns a *ConnectionDirection when successful
func (m *NetworkConnection) GetDirection()(*ConnectionDirection) {
    val, err := m.GetBackingStore().Get("direction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConnectionDirection)
    }
    return nil
}
// GetDomainRegisteredDateTime gets the domainRegisteredDateTime property value. Date when the destination domain was registered. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *NetworkConnection) GetDomainRegisteredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("domainRegisteredDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *NetworkConnection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["applicationName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationName(val)
        }
        return nil
    }
    res["destinationAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationAddress(val)
        }
        return nil
    }
    res["destinationDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationDomain(val)
        }
        return nil
    }
    res["destinationLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationLocation(val)
        }
        return nil
    }
    res["destinationPort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationPort(val)
        }
        return nil
    }
    res["destinationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationUrl(val)
        }
        return nil
    }
    res["direction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConnectionDirection)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirection(val.(*ConnectionDirection))
        }
        return nil
    }
    res["domainRegisteredDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDomainRegisteredDateTime(val)
        }
        return nil
    }
    res["localDnsName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocalDnsName(val)
        }
        return nil
    }
    res["natDestinationAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNatDestinationAddress(val)
        }
        return nil
    }
    res["natDestinationPort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNatDestinationPort(val)
        }
        return nil
    }
    res["natSourceAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNatSourceAddress(val)
        }
        return nil
    }
    res["natSourcePort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNatSourcePort(val)
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
    res["protocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSecurityNetworkProtocol)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtocol(val.(*SecurityNetworkProtocol))
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
    res["sourceAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceAddress(val)
        }
        return nil
    }
    res["sourceLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceLocation(val)
        }
        return nil
    }
    res["sourcePort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourcePort(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConnectionStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*ConnectionStatus))
        }
        return nil
    }
    res["urlParameters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUrlParameters(val)
        }
        return nil
    }
    return res
}
// GetLocalDnsName gets the localDnsName property value. The local DNS name resolution as it appears in the host's local DNS cache (for example, in case the 'hosts' file was tampered with).
// returns a *string when successful
func (m *NetworkConnection) GetLocalDnsName()(*string) {
    val, err := m.GetBackingStore().Get("localDnsName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNatDestinationAddress gets the natDestinationAddress property value. Network Address Translation destination IP address.
// returns a *string when successful
func (m *NetworkConnection) GetNatDestinationAddress()(*string) {
    val, err := m.GetBackingStore().Get("natDestinationAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNatDestinationPort gets the natDestinationPort property value. Network Address Translation destination port.
// returns a *string when successful
func (m *NetworkConnection) GetNatDestinationPort()(*string) {
    val, err := m.GetBackingStore().Get("natDestinationPort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNatSourceAddress gets the natSourceAddress property value. Network Address Translation source IP address.
// returns a *string when successful
func (m *NetworkConnection) GetNatSourceAddress()(*string) {
    val, err := m.GetBackingStore().Get("natSourceAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNatSourcePort gets the natSourcePort property value. Network Address Translation source port.
// returns a *string when successful
func (m *NetworkConnection) GetNatSourcePort()(*string) {
    val, err := m.GetBackingStore().Get("natSourcePort")
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
func (m *NetworkConnection) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProtocol gets the protocol property value. Network protocol. Possible values are: unknown, ip, icmp, igmp, ggp, ipv4, tcp, pup, udp, idp, ipv6, ipv6RoutingHeader, ipv6FragmentHeader, ipSecEncapsulatingSecurityPayload, ipSecAuthenticationHeader, icmpV6, ipv6NoNextHeader, ipv6DestinationOptions, nd, raw, ipx, spx, spxII.
// returns a *SecurityNetworkProtocol when successful
func (m *NetworkConnection) GetProtocol()(*SecurityNetworkProtocol) {
    val, err := m.GetBackingStore().Get("protocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SecurityNetworkProtocol)
    }
    return nil
}
// GetRiskScore gets the riskScore property value. Provider generated/calculated risk score of the network connection. Recommended value range of 0-1, which equates to a percentage.
// returns a *string when successful
func (m *NetworkConnection) GetRiskScore()(*string) {
    val, err := m.GetBackingStore().Get("riskScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSourceAddress gets the sourceAddress property value. Source (i.e. origin) IP address (of the network connection).
// returns a *string when successful
func (m *NetworkConnection) GetSourceAddress()(*string) {
    val, err := m.GetBackingStore().Get("sourceAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSourceLocation gets the sourceLocation property value. Location (by IP address mapping) associated with the source of a network connection.
// returns a *string when successful
func (m *NetworkConnection) GetSourceLocation()(*string) {
    val, err := m.GetBackingStore().Get("sourceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSourcePort gets the sourcePort property value. Source (i.e. origin) IP port (of the network connection).
// returns a *string when successful
func (m *NetworkConnection) GetSourcePort()(*string) {
    val, err := m.GetBackingStore().Get("sourcePort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. Network connection status. Possible values are: unknown, attempted, succeeded, blocked, failed.
// returns a *ConnectionStatus when successful
func (m *NetworkConnection) GetStatus()(*ConnectionStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConnectionStatus)
    }
    return nil
}
// GetUrlParameters gets the urlParameters property value. Parameters (suffix) of the destination URL.
// returns a *string when successful
func (m *NetworkConnection) GetUrlParameters()(*string) {
    val, err := m.GetBackingStore().Get("urlParameters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *NetworkConnection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("applicationName", m.GetApplicationName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("destinationAddress", m.GetDestinationAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("destinationDomain", m.GetDestinationDomain())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("destinationLocation", m.GetDestinationLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("destinationPort", m.GetDestinationPort())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("destinationUrl", m.GetDestinationUrl())
        if err != nil {
            return err
        }
    }
    if m.GetDirection() != nil {
        cast := (*m.GetDirection()).String()
        err := writer.WriteStringValue("direction", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("domainRegisteredDateTime", m.GetDomainRegisteredDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("localDnsName", m.GetLocalDnsName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("natDestinationAddress", m.GetNatDestinationAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("natDestinationPort", m.GetNatDestinationPort())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("natSourceAddress", m.GetNatSourceAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("natSourcePort", m.GetNatSourcePort())
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
    if m.GetProtocol() != nil {
        cast := (*m.GetProtocol()).String()
        err := writer.WriteStringValue("protocol", &cast)
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
        err := writer.WriteStringValue("sourceAddress", m.GetSourceAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sourceLocation", m.GetSourceLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sourcePort", m.GetSourcePort())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err := writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("urlParameters", m.GetUrlParameters())
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
func (m *NetworkConnection) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationName sets the applicationName property value. Name of the application managing the network connection (for example, Facebook or SMTP).
func (m *NetworkConnection) SetApplicationName(value *string)() {
    err := m.GetBackingStore().Set("applicationName", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *NetworkConnection) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDestinationAddress sets the destinationAddress property value. Destination IP address (of the network connection).
func (m *NetworkConnection) SetDestinationAddress(value *string)() {
    err := m.GetBackingStore().Set("destinationAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationDomain sets the destinationDomain property value. Destination domain portion of the destination URL. (for example 'www.contoso.com').
func (m *NetworkConnection) SetDestinationDomain(value *string)() {
    err := m.GetBackingStore().Set("destinationDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationLocation sets the destinationLocation property value. Location (by IP address mapping) associated with the destination of a network connection.
func (m *NetworkConnection) SetDestinationLocation(value *string)() {
    err := m.GetBackingStore().Set("destinationLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationPort sets the destinationPort property value. Destination port (of the network connection).
func (m *NetworkConnection) SetDestinationPort(value *string)() {
    err := m.GetBackingStore().Set("destinationPort", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationUrl sets the destinationUrl property value. Network connection URL/URI string - excluding parameters. (for example 'www.contoso.com/products/default.html')
func (m *NetworkConnection) SetDestinationUrl(value *string)() {
    err := m.GetBackingStore().Set("destinationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetDirection sets the direction property value. Network connection direction. Possible values are: unknown, inbound, outbound.
func (m *NetworkConnection) SetDirection(value *ConnectionDirection)() {
    err := m.GetBackingStore().Set("direction", value)
    if err != nil {
        panic(err)
    }
}
// SetDomainRegisteredDateTime sets the domainRegisteredDateTime property value. Date when the destination domain was registered. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *NetworkConnection) SetDomainRegisteredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("domainRegisteredDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLocalDnsName sets the localDnsName property value. The local DNS name resolution as it appears in the host's local DNS cache (for example, in case the 'hosts' file was tampered with).
func (m *NetworkConnection) SetLocalDnsName(value *string)() {
    err := m.GetBackingStore().Set("localDnsName", value)
    if err != nil {
        panic(err)
    }
}
// SetNatDestinationAddress sets the natDestinationAddress property value. Network Address Translation destination IP address.
func (m *NetworkConnection) SetNatDestinationAddress(value *string)() {
    err := m.GetBackingStore().Set("natDestinationAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetNatDestinationPort sets the natDestinationPort property value. Network Address Translation destination port.
func (m *NetworkConnection) SetNatDestinationPort(value *string)() {
    err := m.GetBackingStore().Set("natDestinationPort", value)
    if err != nil {
        panic(err)
    }
}
// SetNatSourceAddress sets the natSourceAddress property value. Network Address Translation source IP address.
func (m *NetworkConnection) SetNatSourceAddress(value *string)() {
    err := m.GetBackingStore().Set("natSourceAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetNatSourcePort sets the natSourcePort property value. Network Address Translation source port.
func (m *NetworkConnection) SetNatSourcePort(value *string)() {
    err := m.GetBackingStore().Set("natSourcePort", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *NetworkConnection) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocol sets the protocol property value. Network protocol. Possible values are: unknown, ip, icmp, igmp, ggp, ipv4, tcp, pup, udp, idp, ipv6, ipv6RoutingHeader, ipv6FragmentHeader, ipSecEncapsulatingSecurityPayload, ipSecAuthenticationHeader, icmpV6, ipv6NoNextHeader, ipv6DestinationOptions, nd, raw, ipx, spx, spxII.
func (m *NetworkConnection) SetProtocol(value *SecurityNetworkProtocol)() {
    err := m.GetBackingStore().Set("protocol", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskScore sets the riskScore property value. Provider generated/calculated risk score of the network connection. Recommended value range of 0-1, which equates to a percentage.
func (m *NetworkConnection) SetRiskScore(value *string)() {
    err := m.GetBackingStore().Set("riskScore", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceAddress sets the sourceAddress property value. Source (i.e. origin) IP address (of the network connection).
func (m *NetworkConnection) SetSourceAddress(value *string)() {
    err := m.GetBackingStore().Set("sourceAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceLocation sets the sourceLocation property value. Location (by IP address mapping) associated with the source of a network connection.
func (m *NetworkConnection) SetSourceLocation(value *string)() {
    err := m.GetBackingStore().Set("sourceLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetSourcePort sets the sourcePort property value. Source (i.e. origin) IP port (of the network connection).
func (m *NetworkConnection) SetSourcePort(value *string)() {
    err := m.GetBackingStore().Set("sourcePort", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Network connection status. Possible values are: unknown, attempted, succeeded, blocked, failed.
func (m *NetworkConnection) SetStatus(value *ConnectionStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetUrlParameters sets the urlParameters property value. Parameters (suffix) of the destination URL.
func (m *NetworkConnection) SetUrlParameters(value *string)() {
    err := m.GetBackingStore().Set("urlParameters", value)
    if err != nil {
        panic(err)
    }
}
type NetworkConnectionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationName()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDestinationAddress()(*string)
    GetDestinationDomain()(*string)
    GetDestinationLocation()(*string)
    GetDestinationPort()(*string)
    GetDestinationUrl()(*string)
    GetDirection()(*ConnectionDirection)
    GetDomainRegisteredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLocalDnsName()(*string)
    GetNatDestinationAddress()(*string)
    GetNatDestinationPort()(*string)
    GetNatSourceAddress()(*string)
    GetNatSourcePort()(*string)
    GetOdataType()(*string)
    GetProtocol()(*SecurityNetworkProtocol)
    GetRiskScore()(*string)
    GetSourceAddress()(*string)
    GetSourceLocation()(*string)
    GetSourcePort()(*string)
    GetStatus()(*ConnectionStatus)
    GetUrlParameters()(*string)
    SetApplicationName(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDestinationAddress(value *string)()
    SetDestinationDomain(value *string)()
    SetDestinationLocation(value *string)()
    SetDestinationPort(value *string)()
    SetDestinationUrl(value *string)()
    SetDirection(value *ConnectionDirection)()
    SetDomainRegisteredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLocalDnsName(value *string)()
    SetNatDestinationAddress(value *string)()
    SetNatDestinationPort(value *string)()
    SetNatSourceAddress(value *string)()
    SetNatSourcePort(value *string)()
    SetOdataType(value *string)()
    SetProtocol(value *SecurityNetworkProtocol)()
    SetRiskScore(value *string)()
    SetSourceAddress(value *string)()
    SetSourceLocation(value *string)()
    SetSourcePort(value *string)()
    SetStatus(value *ConnectionStatus)()
    SetUrlParameters(value *string)()
}

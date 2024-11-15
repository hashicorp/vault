package callrecords

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type NetworkInfo struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewNetworkInfo instantiates a new NetworkInfo and sets the default values.
func NewNetworkInfo()(*NetworkInfo) {
    m := &NetworkInfo{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateNetworkInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateNetworkInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNetworkInfo(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *NetworkInfo) GetAdditionalData()(map[string]any) {
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
func (m *NetworkInfo) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBandwidthLowEventRatio gets the bandwidthLowEventRatio property value. Fraction of the call that the media endpoint detected the available bandwidth or bandwidth policy was low enough to cause poor quality of the audio sent.
// returns a *float32 when successful
func (m *NetworkInfo) GetBandwidthLowEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("bandwidthLowEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetBasicServiceSetIdentifier gets the basicServiceSetIdentifier property value. The wireless LAN basic service set identifier of the media endpoint used to connect to the network.
// returns a *string when successful
func (m *NetworkInfo) GetBasicServiceSetIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("basicServiceSetIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConnectionType gets the connectionType property value. The connectionType property
// returns a *NetworkConnectionType when successful
func (m *NetworkInfo) GetConnectionType()(*NetworkConnectionType) {
    val, err := m.GetBackingStore().Get("connectionType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*NetworkConnectionType)
    }
    return nil
}
// GetDelayEventRatio gets the delayEventRatio property value. Fraction of the call that the media endpoint detected the network delay was significant enough to impact the ability to have real-time two-way communication.
// returns a *float32 when successful
func (m *NetworkInfo) GetDelayEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("delayEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetDnsSuffix gets the dnsSuffix property value. DNS suffix associated with the network adapter of the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetDnsSuffix()(*string) {
    val, err := m.GetBackingStore().Get("dnsSuffix")
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
func (m *NetworkInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bandwidthLowEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBandwidthLowEventRatio(val)
        }
        return nil
    }
    res["basicServiceSetIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBasicServiceSetIdentifier(val)
        }
        return nil
    }
    res["connectionType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseNetworkConnectionType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectionType(val.(*NetworkConnectionType))
        }
        return nil
    }
    res["delayEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDelayEventRatio(val)
        }
        return nil
    }
    res["dnsSuffix"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDnsSuffix(val)
        }
        return nil
    }
    res["ipAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIpAddress(val)
        }
        return nil
    }
    res["linkSpeed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLinkSpeed(val)
        }
        return nil
    }
    res["macAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMacAddress(val)
        }
        return nil
    }
    res["networkTransportProtocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseNetworkTransportProtocol)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkTransportProtocol(val.(*NetworkTransportProtocol))
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
    res["port"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPort(val)
        }
        return nil
    }
    res["receivedQualityEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReceivedQualityEventRatio(val)
        }
        return nil
    }
    res["reflexiveIPAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReflexiveIPAddress(val)
        }
        return nil
    }
    res["relayIPAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRelayIPAddress(val)
        }
        return nil
    }
    res["relayPort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRelayPort(val)
        }
        return nil
    }
    res["sentQualityEventRatio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentQualityEventRatio(val)
        }
        return nil
    }
    res["subnet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubnet(val)
        }
        return nil
    }
    res["traceRouteHops"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTraceRouteHopFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TraceRouteHopable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TraceRouteHopable)
                }
            }
            m.SetTraceRouteHops(res)
        }
        return nil
    }
    res["wifiBand"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWifiBand)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiBand(val.(*WifiBand))
        }
        return nil
    }
    res["wifiBatteryCharge"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiBatteryCharge(val)
        }
        return nil
    }
    res["wifiChannel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiChannel(val)
        }
        return nil
    }
    res["wifiMicrosoftDriver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiMicrosoftDriver(val)
        }
        return nil
    }
    res["wifiMicrosoftDriverVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiMicrosoftDriverVersion(val)
        }
        return nil
    }
    res["wifiRadioType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWifiRadioType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiRadioType(val.(*WifiRadioType))
        }
        return nil
    }
    res["wifiSignalStrength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiSignalStrength(val)
        }
        return nil
    }
    res["wifiVendorDriver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiVendorDriver(val)
        }
        return nil
    }
    res["wifiVendorDriverVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiVendorDriverVersion(val)
        }
        return nil
    }
    return res
}
// GetIpAddress gets the ipAddress property value. IP address of the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("ipAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLinkSpeed gets the linkSpeed property value. Link speed in bits per second reported by the network adapter used by the media endpoint.
// returns a *int64 when successful
func (m *NetworkInfo) GetLinkSpeed()(*int64) {
    val, err := m.GetBackingStore().Get("linkSpeed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetMacAddress gets the macAddress property value. The media access control (MAC) address of the media endpoint's network device. This value may be missing or shown as 02:00:00:00:00:00 due to operating system privacy policies.
// returns a *string when successful
func (m *NetworkInfo) GetMacAddress()(*string) {
    val, err := m.GetBackingStore().Get("macAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNetworkTransportProtocol gets the networkTransportProtocol property value. The networkTransportProtocol property
// returns a *NetworkTransportProtocol when successful
func (m *NetworkInfo) GetNetworkTransportProtocol()(*NetworkTransportProtocol) {
    val, err := m.GetBackingStore().Get("networkTransportProtocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*NetworkTransportProtocol)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *NetworkInfo) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPort gets the port property value. Network port number used by media endpoint.
// returns a *int32 when successful
func (m *NetworkInfo) GetPort()(*int32) {
    val, err := m.GetBackingStore().Get("port")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetReceivedQualityEventRatio gets the receivedQualityEventRatio property value. Fraction of the call that the media endpoint detected the network was causing poor quality of the audio received.
// returns a *float32 when successful
func (m *NetworkInfo) GetReceivedQualityEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("receivedQualityEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetReflexiveIPAddress gets the reflexiveIPAddress property value. IP address of the media endpoint as seen by the media relay server. This is typically the public internet IP address associated to the endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetReflexiveIPAddress()(*string) {
    val, err := m.GetBackingStore().Get("reflexiveIPAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRelayIPAddress gets the relayIPAddress property value. IP address of the media relay server allocated by the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetRelayIPAddress()(*string) {
    val, err := m.GetBackingStore().Get("relayIPAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRelayPort gets the relayPort property value. Network port number allocated on the media relay server by the media endpoint.
// returns a *int32 when successful
func (m *NetworkInfo) GetRelayPort()(*int32) {
    val, err := m.GetBackingStore().Get("relayPort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSentQualityEventRatio gets the sentQualityEventRatio property value. Fraction of the call that the media endpoint detected the network was causing poor quality of the audio sent.
// returns a *float32 when successful
func (m *NetworkInfo) GetSentQualityEventRatio()(*float32) {
    val, err := m.GetBackingStore().Get("sentQualityEventRatio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// GetSubnet gets the subnet property value. Subnet used for media stream by the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetSubnet()(*string) {
    val, err := m.GetBackingStore().Get("subnet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTraceRouteHops gets the traceRouteHops property value. List of network trace route hops collected for this media stream.*
// returns a []TraceRouteHopable when successful
func (m *NetworkInfo) GetTraceRouteHops()([]TraceRouteHopable) {
    val, err := m.GetBackingStore().Get("traceRouteHops")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TraceRouteHopable)
    }
    return nil
}
// GetWifiBand gets the wifiBand property value. The wifiBand property
// returns a *WifiBand when successful
func (m *NetworkInfo) GetWifiBand()(*WifiBand) {
    val, err := m.GetBackingStore().Get("wifiBand")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WifiBand)
    }
    return nil
}
// GetWifiBatteryCharge gets the wifiBatteryCharge property value. Estimated remaining battery charge in percentage reported by the media endpoint.
// returns a *int32 when successful
func (m *NetworkInfo) GetWifiBatteryCharge()(*int32) {
    val, err := m.GetBackingStore().Get("wifiBatteryCharge")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWifiChannel gets the wifiChannel property value. WiFi channel used by the media endpoint.
// returns a *int32 when successful
func (m *NetworkInfo) GetWifiChannel()(*int32) {
    val, err := m.GetBackingStore().Get("wifiChannel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWifiMicrosoftDriver gets the wifiMicrosoftDriver property value. Name of the Microsoft WiFi driver used by the media endpoint. Value may be localized based on the language used by endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetWifiMicrosoftDriver()(*string) {
    val, err := m.GetBackingStore().Get("wifiMicrosoftDriver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWifiMicrosoftDriverVersion gets the wifiMicrosoftDriverVersion property value. Version of the Microsoft WiFi driver used by the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetWifiMicrosoftDriverVersion()(*string) {
    val, err := m.GetBackingStore().Get("wifiMicrosoftDriverVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWifiRadioType gets the wifiRadioType property value. The wifiRadioType property
// returns a *WifiRadioType when successful
func (m *NetworkInfo) GetWifiRadioType()(*WifiRadioType) {
    val, err := m.GetBackingStore().Get("wifiRadioType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WifiRadioType)
    }
    return nil
}
// GetWifiSignalStrength gets the wifiSignalStrength property value. WiFi signal strength in percentage reported by the media endpoint.
// returns a *int32 when successful
func (m *NetworkInfo) GetWifiSignalStrength()(*int32) {
    val, err := m.GetBackingStore().Get("wifiSignalStrength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWifiVendorDriver gets the wifiVendorDriver property value. Name of the WiFi driver used by the media endpoint. Value may be localized based on the language used by endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetWifiVendorDriver()(*string) {
    val, err := m.GetBackingStore().Get("wifiVendorDriver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWifiVendorDriverVersion gets the wifiVendorDriverVersion property value. Version of the WiFi driver used by the media endpoint.
// returns a *string when successful
func (m *NetworkInfo) GetWifiVendorDriverVersion()(*string) {
    val, err := m.GetBackingStore().Get("wifiVendorDriverVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *NetworkInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteFloat32Value("bandwidthLowEventRatio", m.GetBandwidthLowEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("basicServiceSetIdentifier", m.GetBasicServiceSetIdentifier())
        if err != nil {
            return err
        }
    }
    if m.GetConnectionType() != nil {
        cast := (*m.GetConnectionType()).String()
        err := writer.WriteStringValue("connectionType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("delayEventRatio", m.GetDelayEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("dnsSuffix", m.GetDnsSuffix())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("ipAddress", m.GetIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("linkSpeed", m.GetLinkSpeed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("macAddress", m.GetMacAddress())
        if err != nil {
            return err
        }
    }
    if m.GetNetworkTransportProtocol() != nil {
        cast := (*m.GetNetworkTransportProtocol()).String()
        err := writer.WriteStringValue("networkTransportProtocol", &cast)
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
        err := writer.WriteInt32Value("port", m.GetPort())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("receivedQualityEventRatio", m.GetReceivedQualityEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("reflexiveIPAddress", m.GetReflexiveIPAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("relayIPAddress", m.GetRelayIPAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("relayPort", m.GetRelayPort())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("sentQualityEventRatio", m.GetSentQualityEventRatio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("subnet", m.GetSubnet())
        if err != nil {
            return err
        }
    }
    if m.GetTraceRouteHops() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTraceRouteHops()))
        for i, v := range m.GetTraceRouteHops() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("traceRouteHops", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWifiBand() != nil {
        cast := (*m.GetWifiBand()).String()
        err := writer.WriteStringValue("wifiBand", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("wifiBatteryCharge", m.GetWifiBatteryCharge())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("wifiChannel", m.GetWifiChannel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("wifiMicrosoftDriver", m.GetWifiMicrosoftDriver())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("wifiMicrosoftDriverVersion", m.GetWifiMicrosoftDriverVersion())
        if err != nil {
            return err
        }
    }
    if m.GetWifiRadioType() != nil {
        cast := (*m.GetWifiRadioType()).String()
        err := writer.WriteStringValue("wifiRadioType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("wifiSignalStrength", m.GetWifiSignalStrength())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("wifiVendorDriver", m.GetWifiVendorDriver())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("wifiVendorDriverVersion", m.GetWifiVendorDriverVersion())
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
func (m *NetworkInfo) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *NetworkInfo) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBandwidthLowEventRatio sets the bandwidthLowEventRatio property value. Fraction of the call that the media endpoint detected the available bandwidth or bandwidth policy was low enough to cause poor quality of the audio sent.
func (m *NetworkInfo) SetBandwidthLowEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("bandwidthLowEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetBasicServiceSetIdentifier sets the basicServiceSetIdentifier property value. The wireless LAN basic service set identifier of the media endpoint used to connect to the network.
func (m *NetworkInfo) SetBasicServiceSetIdentifier(value *string)() {
    err := m.GetBackingStore().Set("basicServiceSetIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetConnectionType sets the connectionType property value. The connectionType property
func (m *NetworkInfo) SetConnectionType(value *NetworkConnectionType)() {
    err := m.GetBackingStore().Set("connectionType", value)
    if err != nil {
        panic(err)
    }
}
// SetDelayEventRatio sets the delayEventRatio property value. Fraction of the call that the media endpoint detected the network delay was significant enough to impact the ability to have real-time two-way communication.
func (m *NetworkInfo) SetDelayEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("delayEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetDnsSuffix sets the dnsSuffix property value. DNS suffix associated with the network adapter of the media endpoint.
func (m *NetworkInfo) SetDnsSuffix(value *string)() {
    err := m.GetBackingStore().Set("dnsSuffix", value)
    if err != nil {
        panic(err)
    }
}
// SetIpAddress sets the ipAddress property value. IP address of the media endpoint.
func (m *NetworkInfo) SetIpAddress(value *string)() {
    err := m.GetBackingStore().Set("ipAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetLinkSpeed sets the linkSpeed property value. Link speed in bits per second reported by the network adapter used by the media endpoint.
func (m *NetworkInfo) SetLinkSpeed(value *int64)() {
    err := m.GetBackingStore().Set("linkSpeed", value)
    if err != nil {
        panic(err)
    }
}
// SetMacAddress sets the macAddress property value. The media access control (MAC) address of the media endpoint's network device. This value may be missing or shown as 02:00:00:00:00:00 due to operating system privacy policies.
func (m *NetworkInfo) SetMacAddress(value *string)() {
    err := m.GetBackingStore().Set("macAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkTransportProtocol sets the networkTransportProtocol property value. The networkTransportProtocol property
func (m *NetworkInfo) SetNetworkTransportProtocol(value *NetworkTransportProtocol)() {
    err := m.GetBackingStore().Set("networkTransportProtocol", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *NetworkInfo) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPort sets the port property value. Network port number used by media endpoint.
func (m *NetworkInfo) SetPort(value *int32)() {
    err := m.GetBackingStore().Set("port", value)
    if err != nil {
        panic(err)
    }
}
// SetReceivedQualityEventRatio sets the receivedQualityEventRatio property value. Fraction of the call that the media endpoint detected the network was causing poor quality of the audio received.
func (m *NetworkInfo) SetReceivedQualityEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("receivedQualityEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetReflexiveIPAddress sets the reflexiveIPAddress property value. IP address of the media endpoint as seen by the media relay server. This is typically the public internet IP address associated to the endpoint.
func (m *NetworkInfo) SetReflexiveIPAddress(value *string)() {
    err := m.GetBackingStore().Set("reflexiveIPAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetRelayIPAddress sets the relayIPAddress property value. IP address of the media relay server allocated by the media endpoint.
func (m *NetworkInfo) SetRelayIPAddress(value *string)() {
    err := m.GetBackingStore().Set("relayIPAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetRelayPort sets the relayPort property value. Network port number allocated on the media relay server by the media endpoint.
func (m *NetworkInfo) SetRelayPort(value *int32)() {
    err := m.GetBackingStore().Set("relayPort", value)
    if err != nil {
        panic(err)
    }
}
// SetSentQualityEventRatio sets the sentQualityEventRatio property value. Fraction of the call that the media endpoint detected the network was causing poor quality of the audio sent.
func (m *NetworkInfo) SetSentQualityEventRatio(value *float32)() {
    err := m.GetBackingStore().Set("sentQualityEventRatio", value)
    if err != nil {
        panic(err)
    }
}
// SetSubnet sets the subnet property value. Subnet used for media stream by the media endpoint.
func (m *NetworkInfo) SetSubnet(value *string)() {
    err := m.GetBackingStore().Set("subnet", value)
    if err != nil {
        panic(err)
    }
}
// SetTraceRouteHops sets the traceRouteHops property value. List of network trace route hops collected for this media stream.*
func (m *NetworkInfo) SetTraceRouteHops(value []TraceRouteHopable)() {
    err := m.GetBackingStore().Set("traceRouteHops", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiBand sets the wifiBand property value. The wifiBand property
func (m *NetworkInfo) SetWifiBand(value *WifiBand)() {
    err := m.GetBackingStore().Set("wifiBand", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiBatteryCharge sets the wifiBatteryCharge property value. Estimated remaining battery charge in percentage reported by the media endpoint.
func (m *NetworkInfo) SetWifiBatteryCharge(value *int32)() {
    err := m.GetBackingStore().Set("wifiBatteryCharge", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiChannel sets the wifiChannel property value. WiFi channel used by the media endpoint.
func (m *NetworkInfo) SetWifiChannel(value *int32)() {
    err := m.GetBackingStore().Set("wifiChannel", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiMicrosoftDriver sets the wifiMicrosoftDriver property value. Name of the Microsoft WiFi driver used by the media endpoint. Value may be localized based on the language used by endpoint.
func (m *NetworkInfo) SetWifiMicrosoftDriver(value *string)() {
    err := m.GetBackingStore().Set("wifiMicrosoftDriver", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiMicrosoftDriverVersion sets the wifiMicrosoftDriverVersion property value. Version of the Microsoft WiFi driver used by the media endpoint.
func (m *NetworkInfo) SetWifiMicrosoftDriverVersion(value *string)() {
    err := m.GetBackingStore().Set("wifiMicrosoftDriverVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiRadioType sets the wifiRadioType property value. The wifiRadioType property
func (m *NetworkInfo) SetWifiRadioType(value *WifiRadioType)() {
    err := m.GetBackingStore().Set("wifiRadioType", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiSignalStrength sets the wifiSignalStrength property value. WiFi signal strength in percentage reported by the media endpoint.
func (m *NetworkInfo) SetWifiSignalStrength(value *int32)() {
    err := m.GetBackingStore().Set("wifiSignalStrength", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiVendorDriver sets the wifiVendorDriver property value. Name of the WiFi driver used by the media endpoint. Value may be localized based on the language used by endpoint.
func (m *NetworkInfo) SetWifiVendorDriver(value *string)() {
    err := m.GetBackingStore().Set("wifiVendorDriver", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiVendorDriverVersion sets the wifiVendorDriverVersion property value. Version of the WiFi driver used by the media endpoint.
func (m *NetworkInfo) SetWifiVendorDriverVersion(value *string)() {
    err := m.GetBackingStore().Set("wifiVendorDriverVersion", value)
    if err != nil {
        panic(err)
    }
}
type NetworkInfoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBandwidthLowEventRatio()(*float32)
    GetBasicServiceSetIdentifier()(*string)
    GetConnectionType()(*NetworkConnectionType)
    GetDelayEventRatio()(*float32)
    GetDnsSuffix()(*string)
    GetIpAddress()(*string)
    GetLinkSpeed()(*int64)
    GetMacAddress()(*string)
    GetNetworkTransportProtocol()(*NetworkTransportProtocol)
    GetOdataType()(*string)
    GetPort()(*int32)
    GetReceivedQualityEventRatio()(*float32)
    GetReflexiveIPAddress()(*string)
    GetRelayIPAddress()(*string)
    GetRelayPort()(*int32)
    GetSentQualityEventRatio()(*float32)
    GetSubnet()(*string)
    GetTraceRouteHops()([]TraceRouteHopable)
    GetWifiBand()(*WifiBand)
    GetWifiBatteryCharge()(*int32)
    GetWifiChannel()(*int32)
    GetWifiMicrosoftDriver()(*string)
    GetWifiMicrosoftDriverVersion()(*string)
    GetWifiRadioType()(*WifiRadioType)
    GetWifiSignalStrength()(*int32)
    GetWifiVendorDriver()(*string)
    GetWifiVendorDriverVersion()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBandwidthLowEventRatio(value *float32)()
    SetBasicServiceSetIdentifier(value *string)()
    SetConnectionType(value *NetworkConnectionType)()
    SetDelayEventRatio(value *float32)()
    SetDnsSuffix(value *string)()
    SetIpAddress(value *string)()
    SetLinkSpeed(value *int64)()
    SetMacAddress(value *string)()
    SetNetworkTransportProtocol(value *NetworkTransportProtocol)()
    SetOdataType(value *string)()
    SetPort(value *int32)()
    SetReceivedQualityEventRatio(value *float32)()
    SetReflexiveIPAddress(value *string)()
    SetRelayIPAddress(value *string)()
    SetRelayPort(value *int32)()
    SetSentQualityEventRatio(value *float32)()
    SetSubnet(value *string)()
    SetTraceRouteHops(value []TraceRouteHopable)()
    SetWifiBand(value *WifiBand)()
    SetWifiBatteryCharge(value *int32)()
    SetWifiChannel(value *int32)()
    SetWifiMicrosoftDriver(value *string)()
    SetWifiMicrosoftDriverVersion(value *string)()
    SetWifiRadioType(value *WifiRadioType)()
    SetWifiSignalStrength(value *int32)()
    SetWifiVendorDriver(value *string)()
    SetWifiVendorDriverVersion(value *string)()
}

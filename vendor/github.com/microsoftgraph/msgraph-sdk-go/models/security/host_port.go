package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type HostPort struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewHostPort instantiates a new HostPort and sets the default values.
func NewHostPort()(*HostPort) {
    m := &HostPort{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateHostPortFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHostPortFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHostPort(), nil
}
// GetBanners gets the banners property value. The hostPortBanners retrieved from scanning the port.
// returns a []HostPortBannerable when successful
func (m *HostPort) GetBanners()([]HostPortBannerable) {
    val, err := m.GetBackingStore().Get("banners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPortBannerable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *HostPort) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["banners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPortBannerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPortBannerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPortBannerable)
                }
            }
            m.SetBanners(res)
        }
        return nil
    }
    res["firstSeenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirstSeenDateTime(val)
        }
        return nil
    }
    res["host"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateHostFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHost(val.(Hostable))
        }
        return nil
    }
    res["lastScanDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastScanDateTime(val)
        }
        return nil
    }
    res["lastSeenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSeenDateTime(val)
        }
        return nil
    }
    res["mostRecentSslCertificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSslCertificateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMostRecentSslCertificate(val.(SslCertificateable))
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
    res["protocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHostPortProtocol)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtocol(val.(*HostPortProtocol))
        }
        return nil
    }
    res["services"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPortComponentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPortComponentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPortComponentable)
                }
            }
            m.SetServices(res)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseHostPortStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*HostPortStatus))
        }
        return nil
    }
    res["timesObserved"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimesObserved(val)
        }
        return nil
    }
    return res
}
// GetFirstSeenDateTime gets the firstSeenDateTime property value. The first date and time when Microsoft Defender Threat Intelligence observed the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *HostPort) GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("firstSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetHost gets the host property value. The host property
// returns a Hostable when successful
func (m *HostPort) GetHost()(Hostable) {
    val, err := m.GetBackingStore().Get("host")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Hostable)
    }
    return nil
}
// GetLastScanDateTime gets the lastScanDateTime property value. The last date and time when Microsoft Defender Threat Intelligence scanned the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *HostPort) GetLastScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastScanDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastSeenDateTime gets the lastSeenDateTime property value. The last date and time when Microsoft Defender Threat Intelligence observed the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *HostPort) GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMostRecentSslCertificate gets the mostRecentSslCertificate property value. The most recent sslCertificate used to communicate on the port.
// returns a SslCertificateable when successful
func (m *HostPort) GetMostRecentSslCertificate()(SslCertificateable) {
    val, err := m.GetBackingStore().Get("mostRecentSslCertificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SslCertificateable)
    }
    return nil
}
// GetPort gets the port property value. The numerical identifier of the port which is standardized across the internet.
// returns a *int32 when successful
func (m *HostPort) GetPort()(*int32) {
    val, err := m.GetBackingStore().Get("port")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetProtocol gets the protocol property value. The general protocol used to scan the port. The possible values are: tcp, udp, unknownFutureValue.
// returns a *HostPortProtocol when successful
func (m *HostPort) GetProtocol()(*HostPortProtocol) {
    val, err := m.GetBackingStore().Get("protocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HostPortProtocol)
    }
    return nil
}
// GetServices gets the services property value. The hostPortComponents retrieved from scanning the port.
// returns a []HostPortComponentable when successful
func (m *HostPort) GetServices()([]HostPortComponentable) {
    val, err := m.GetBackingStore().Get("services")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPortComponentable)
    }
    return nil
}
// GetStatus gets the status property value. The status of the port. The possible values are: open, filtered, closed, unknownFutureValue.
// returns a *HostPortStatus when successful
func (m *HostPort) GetStatus()(*HostPortStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*HostPortStatus)
    }
    return nil
}
// GetTimesObserved gets the timesObserved property value. The total amount of times that Microsoft Defender Threat Intelligence has observed the hostPort in all its scans.
// returns a *int32 when successful
func (m *HostPort) GetTimesObserved()(*int32) {
    val, err := m.GetBackingStore().Get("timesObserved")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *HostPort) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBanners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBanners()))
        for i, v := range m.GetBanners() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("banners", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("firstSeenDateTime", m.GetFirstSeenDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("host", m.GetHost())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastScanDateTime", m.GetLastScanDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastSeenDateTime", m.GetLastSeenDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mostRecentSslCertificate", m.GetMostRecentSslCertificate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("port", m.GetPort())
        if err != nil {
            return err
        }
    }
    if m.GetProtocol() != nil {
        cast := (*m.GetProtocol()).String()
        err = writer.WriteStringValue("protocol", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetServices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServices()))
        for i, v := range m.GetServices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("services", cast)
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("timesObserved", m.GetTimesObserved())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBanners sets the banners property value. The hostPortBanners retrieved from scanning the port.
func (m *HostPort) SetBanners(value []HostPortBannerable)() {
    err := m.GetBackingStore().Set("banners", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstSeenDateTime sets the firstSeenDateTime property value. The first date and time when Microsoft Defender Threat Intelligence observed the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
func (m *HostPort) SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("firstSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHost sets the host property value. The host property
func (m *HostPort) SetHost(value Hostable)() {
    err := m.GetBackingStore().Set("host", value)
    if err != nil {
        panic(err)
    }
}
// SetLastScanDateTime sets the lastScanDateTime property value. The last date and time when Microsoft Defender Threat Intelligence scanned the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
func (m *HostPort) SetLastScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastScanDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSeenDateTime sets the lastSeenDateTime property value. The last date and time when Microsoft Defender Threat Intelligence observed the hostPort. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z.
func (m *HostPort) SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMostRecentSslCertificate sets the mostRecentSslCertificate property value. The most recent sslCertificate used to communicate on the port.
func (m *HostPort) SetMostRecentSslCertificate(value SslCertificateable)() {
    err := m.GetBackingStore().Set("mostRecentSslCertificate", value)
    if err != nil {
        panic(err)
    }
}
// SetPort sets the port property value. The numerical identifier of the port which is standardized across the internet.
func (m *HostPort) SetPort(value *int32)() {
    err := m.GetBackingStore().Set("port", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocol sets the protocol property value. The general protocol used to scan the port. The possible values are: tcp, udp, unknownFutureValue.
func (m *HostPort) SetProtocol(value *HostPortProtocol)() {
    err := m.GetBackingStore().Set("protocol", value)
    if err != nil {
        panic(err)
    }
}
// SetServices sets the services property value. The hostPortComponents retrieved from scanning the port.
func (m *HostPort) SetServices(value []HostPortComponentable)() {
    err := m.GetBackingStore().Set("services", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of the port. The possible values are: open, filtered, closed, unknownFutureValue.
func (m *HostPort) SetStatus(value *HostPortStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTimesObserved sets the timesObserved property value. The total amount of times that Microsoft Defender Threat Intelligence has observed the hostPort in all its scans.
func (m *HostPort) SetTimesObserved(value *int32)() {
    err := m.GetBackingStore().Set("timesObserved", value)
    if err != nil {
        panic(err)
    }
}
type HostPortable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBanners()([]HostPortBannerable)
    GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHost()(Hostable)
    GetLastScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMostRecentSslCertificate()(SslCertificateable)
    GetPort()(*int32)
    GetProtocol()(*HostPortProtocol)
    GetServices()([]HostPortComponentable)
    GetStatus()(*HostPortStatus)
    GetTimesObserved()(*int32)
    SetBanners(value []HostPortBannerable)()
    SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHost(value Hostable)()
    SetLastScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMostRecentSslCertificate(value SslCertificateable)()
    SetPort(value *int32)()
    SetProtocol(value *HostPortProtocol)()
    SetServices(value []HostPortComponentable)()
    SetStatus(value *HostPortStatus)()
    SetTimesObserved(value *int32)()
}

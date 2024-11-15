package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Host struct {
    Artifact
}
// NewHost instantiates a new Host and sets the default values.
func NewHost()(*Host) {
    m := &Host{
        Artifact: *NewArtifact(),
    }
    odataTypeValue := "#microsoft.graph.security.host"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateHostFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHostFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.security.hostname":
                        return NewHostname(), nil
                    case "#microsoft.graph.security.ipAddress":
                        return NewIpAddress(), nil
                }
            }
        }
    }
    return NewHost(), nil
}
// GetChildHostPairs gets the childHostPairs property value. The hostPairs that are resources associated with a host, where that host is the parentHost and has an outgoing pairing to a childHost.
// returns a []HostPairable when successful
func (m *Host) GetChildHostPairs()([]HostPairable) {
    val, err := m.GetBackingStore().Get("childHostPairs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPairable)
    }
    return nil
}
// GetComponents gets the components property value. The hostComponents that are associated with this host.
// returns a []HostComponentable when successful
func (m *Host) GetComponents()([]HostComponentable) {
    val, err := m.GetBackingStore().Get("components")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostComponentable)
    }
    return nil
}
// GetCookies gets the cookies property value. The hostCookies that are associated with this host.
// returns a []HostCookieable when successful
func (m *Host) GetCookies()([]HostCookieable) {
    val, err := m.GetBackingStore().Get("cookies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostCookieable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Host) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Artifact.GetFieldDeserializers()
    res["childHostPairs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPairable)
                }
            }
            m.SetChildHostPairs(res)
        }
        return nil
    }
    res["components"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostComponentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostComponentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostComponentable)
                }
            }
            m.SetComponents(res)
        }
        return nil
    }
    res["cookies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostCookieFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostCookieable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostCookieable)
                }
            }
            m.SetCookies(res)
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
    res["hostPairs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPairable)
                }
            }
            m.SetHostPairs(res)
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
    res["parentHostPairs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPairable)
                }
            }
            m.SetParentHostPairs(res)
        }
        return nil
    }
    res["passiveDns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePassiveDnsRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PassiveDnsRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PassiveDnsRecordable)
                }
            }
            m.SetPassiveDns(res)
        }
        return nil
    }
    res["passiveDnsReverse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePassiveDnsRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PassiveDnsRecordable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PassiveDnsRecordable)
                }
            }
            m.SetPassiveDnsReverse(res)
        }
        return nil
    }
    res["ports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostPortFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostPortable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostPortable)
                }
            }
            m.SetPorts(res)
        }
        return nil
    }
    res["reputation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateHostReputationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReputation(val.(HostReputationable))
        }
        return nil
    }
    res["sslCertificates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostSslCertificateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostSslCertificateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostSslCertificateable)
                }
            }
            m.SetSslCertificates(res)
        }
        return nil
    }
    res["subdomains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubdomainFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Subdomainable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Subdomainable)
                }
            }
            m.SetSubdomains(res)
        }
        return nil
    }
    res["trackers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHostTrackerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HostTrackerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(HostTrackerable)
                }
            }
            m.SetTrackers(res)
        }
        return nil
    }
    res["whois"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWhoisRecordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWhois(val.(WhoisRecordable))
        }
        return nil
    }
    return res
}
// GetFirstSeenDateTime gets the firstSeenDateTime property value. The first date and time when this host was observed. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Host) GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("firstSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetHostPairs gets the hostPairs property value. The hostPairs that are associated with this host, where this host is either the parentHost or childHost.
// returns a []HostPairable when successful
func (m *Host) GetHostPairs()([]HostPairable) {
    val, err := m.GetBackingStore().Get("hostPairs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPairable)
    }
    return nil
}
// GetLastSeenDateTime gets the lastSeenDateTime property value. The most recent date and time when this host was observed. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Host) GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetParentHostPairs gets the parentHostPairs property value. The hostPairs that are associated with a host, where that host is the childHost and has an incoming pairing with a parentHost.
// returns a []HostPairable when successful
func (m *Host) GetParentHostPairs()([]HostPairable) {
    val, err := m.GetBackingStore().Get("parentHostPairs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPairable)
    }
    return nil
}
// GetPassiveDns gets the passiveDns property value. Passive DNS retrieval about this host.
// returns a []PassiveDnsRecordable when successful
func (m *Host) GetPassiveDns()([]PassiveDnsRecordable) {
    val, err := m.GetBackingStore().Get("passiveDns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PassiveDnsRecordable)
    }
    return nil
}
// GetPassiveDnsReverse gets the passiveDnsReverse property value. Reverse passive DNS retrieval about this host.
// returns a []PassiveDnsRecordable when successful
func (m *Host) GetPassiveDnsReverse()([]PassiveDnsRecordable) {
    val, err := m.GetBackingStore().Get("passiveDnsReverse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PassiveDnsRecordable)
    }
    return nil
}
// GetPorts gets the ports property value. The hostPorts associated with a host.
// returns a []HostPortable when successful
func (m *Host) GetPorts()([]HostPortable) {
    val, err := m.GetBackingStore().Get("ports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostPortable)
    }
    return nil
}
// GetReputation gets the reputation property value. Represents a calculated reputation of this host.
// returns a HostReputationable when successful
func (m *Host) GetReputation()(HostReputationable) {
    val, err := m.GetBackingStore().Get("reputation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(HostReputationable)
    }
    return nil
}
// GetSslCertificates gets the sslCertificates property value. The hostSslCertificates that are associated with this host.
// returns a []HostSslCertificateable when successful
func (m *Host) GetSslCertificates()([]HostSslCertificateable) {
    val, err := m.GetBackingStore().Get("sslCertificates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostSslCertificateable)
    }
    return nil
}
// GetSubdomains gets the subdomains property value. The subdomains that are associated with this host.
// returns a []Subdomainable when successful
func (m *Host) GetSubdomains()([]Subdomainable) {
    val, err := m.GetBackingStore().Get("subdomains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Subdomainable)
    }
    return nil
}
// GetTrackers gets the trackers property value. The hostTrackers that are associated with this host.
// returns a []HostTrackerable when successful
func (m *Host) GetTrackers()([]HostTrackerable) {
    val, err := m.GetBackingStore().Get("trackers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HostTrackerable)
    }
    return nil
}
// GetWhois gets the whois property value. The most recent whoisRecord for this host.
// returns a WhoisRecordable when successful
func (m *Host) GetWhois()(WhoisRecordable) {
    val, err := m.GetBackingStore().Get("whois")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WhoisRecordable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Host) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Artifact.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChildHostPairs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildHostPairs()))
        for i, v := range m.GetChildHostPairs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("childHostPairs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetComponents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetComponents()))
        for i, v := range m.GetComponents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("components", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCookies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCookies()))
        for i, v := range m.GetCookies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("cookies", cast)
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
    if m.GetHostPairs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHostPairs()))
        for i, v := range m.GetHostPairs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("hostPairs", cast)
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
    if m.GetParentHostPairs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetParentHostPairs()))
        for i, v := range m.GetParentHostPairs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("parentHostPairs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPassiveDns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPassiveDns()))
        for i, v := range m.GetPassiveDns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("passiveDns", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPassiveDnsReverse() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPassiveDnsReverse()))
        for i, v := range m.GetPassiveDnsReverse() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("passiveDnsReverse", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPorts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPorts()))
        for i, v := range m.GetPorts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ports", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("reputation", m.GetReputation())
        if err != nil {
            return err
        }
    }
    if m.GetSslCertificates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSslCertificates()))
        for i, v := range m.GetSslCertificates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sslCertificates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSubdomains() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSubdomains()))
        for i, v := range m.GetSubdomains() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("subdomains", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTrackers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTrackers()))
        for i, v := range m.GetTrackers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("trackers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("whois", m.GetWhois())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChildHostPairs sets the childHostPairs property value. The hostPairs that are resources associated with a host, where that host is the parentHost and has an outgoing pairing to a childHost.
func (m *Host) SetChildHostPairs(value []HostPairable)() {
    err := m.GetBackingStore().Set("childHostPairs", value)
    if err != nil {
        panic(err)
    }
}
// SetComponents sets the components property value. The hostComponents that are associated with this host.
func (m *Host) SetComponents(value []HostComponentable)() {
    err := m.GetBackingStore().Set("components", value)
    if err != nil {
        panic(err)
    }
}
// SetCookies sets the cookies property value. The hostCookies that are associated with this host.
func (m *Host) SetCookies(value []HostCookieable)() {
    err := m.GetBackingStore().Set("cookies", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstSeenDateTime sets the firstSeenDateTime property value. The first date and time when this host was observed. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Host) SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("firstSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHostPairs sets the hostPairs property value. The hostPairs that are associated with this host, where this host is either the parentHost or childHost.
func (m *Host) SetHostPairs(value []HostPairable)() {
    err := m.GetBackingStore().Set("hostPairs", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSeenDateTime sets the lastSeenDateTime property value. The most recent date and time when this host was observed. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Host) SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetParentHostPairs sets the parentHostPairs property value. The hostPairs that are associated with a host, where that host is the childHost and has an incoming pairing with a parentHost.
func (m *Host) SetParentHostPairs(value []HostPairable)() {
    err := m.GetBackingStore().Set("parentHostPairs", value)
    if err != nil {
        panic(err)
    }
}
// SetPassiveDns sets the passiveDns property value. Passive DNS retrieval about this host.
func (m *Host) SetPassiveDns(value []PassiveDnsRecordable)() {
    err := m.GetBackingStore().Set("passiveDns", value)
    if err != nil {
        panic(err)
    }
}
// SetPassiveDnsReverse sets the passiveDnsReverse property value. Reverse passive DNS retrieval about this host.
func (m *Host) SetPassiveDnsReverse(value []PassiveDnsRecordable)() {
    err := m.GetBackingStore().Set("passiveDnsReverse", value)
    if err != nil {
        panic(err)
    }
}
// SetPorts sets the ports property value. The hostPorts associated with a host.
func (m *Host) SetPorts(value []HostPortable)() {
    err := m.GetBackingStore().Set("ports", value)
    if err != nil {
        panic(err)
    }
}
// SetReputation sets the reputation property value. Represents a calculated reputation of this host.
func (m *Host) SetReputation(value HostReputationable)() {
    err := m.GetBackingStore().Set("reputation", value)
    if err != nil {
        panic(err)
    }
}
// SetSslCertificates sets the sslCertificates property value. The hostSslCertificates that are associated with this host.
func (m *Host) SetSslCertificates(value []HostSslCertificateable)() {
    err := m.GetBackingStore().Set("sslCertificates", value)
    if err != nil {
        panic(err)
    }
}
// SetSubdomains sets the subdomains property value. The subdomains that are associated with this host.
func (m *Host) SetSubdomains(value []Subdomainable)() {
    err := m.GetBackingStore().Set("subdomains", value)
    if err != nil {
        panic(err)
    }
}
// SetTrackers sets the trackers property value. The hostTrackers that are associated with this host.
func (m *Host) SetTrackers(value []HostTrackerable)() {
    err := m.GetBackingStore().Set("trackers", value)
    if err != nil {
        panic(err)
    }
}
// SetWhois sets the whois property value. The most recent whoisRecord for this host.
func (m *Host) SetWhois(value WhoisRecordable)() {
    err := m.GetBackingStore().Set("whois", value)
    if err != nil {
        panic(err)
    }
}
type Hostable interface {
    Artifactable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChildHostPairs()([]HostPairable)
    GetComponents()([]HostComponentable)
    GetCookies()([]HostCookieable)
    GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHostPairs()([]HostPairable)
    GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetParentHostPairs()([]HostPairable)
    GetPassiveDns()([]PassiveDnsRecordable)
    GetPassiveDnsReverse()([]PassiveDnsRecordable)
    GetPorts()([]HostPortable)
    GetReputation()(HostReputationable)
    GetSslCertificates()([]HostSslCertificateable)
    GetSubdomains()([]Subdomainable)
    GetTrackers()([]HostTrackerable)
    GetWhois()(WhoisRecordable)
    SetChildHostPairs(value []HostPairable)()
    SetComponents(value []HostComponentable)()
    SetCookies(value []HostCookieable)()
    SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHostPairs(value []HostPairable)()
    SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetParentHostPairs(value []HostPairable)()
    SetPassiveDns(value []PassiveDnsRecordable)()
    SetPassiveDnsReverse(value []PassiveDnsRecordable)()
    SetPorts(value []HostPortable)()
    SetReputation(value HostReputationable)()
    SetSslCertificates(value []HostSslCertificateable)()
    SetSubdomains(value []Subdomainable)()
    SetTrackers(value []HostTrackerable)()
    SetWhois(value WhoisRecordable)()
}

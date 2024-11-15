package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IpAddress struct {
    Host
}
// NewIpAddress instantiates a new IpAddress and sets the default values.
func NewIpAddress()(*IpAddress) {
    m := &IpAddress{
        Host: *NewHost(),
    }
    odataTypeValue := "#microsoft.graph.security.ipAddress"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIpAddressFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIpAddressFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIpAddress(), nil
}
// GetAutonomousSystem gets the autonomousSystem property value. The details about the autonomous system to which this IP address belongs.
// returns a AutonomousSystemable when successful
func (m *IpAddress) GetAutonomousSystem()(AutonomousSystemable) {
    val, err := m.GetBackingStore().Get("autonomousSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AutonomousSystemable)
    }
    return nil
}
// GetCountryOrRegion gets the countryOrRegion property value. The country/region for this IP address.
// returns a *string when successful
func (m *IpAddress) GetCountryOrRegion()(*string) {
    val, err := m.GetBackingStore().Get("countryOrRegion")
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
func (m *IpAddress) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Host.GetFieldDeserializers()
    res["autonomousSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAutonomousSystemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutonomousSystem(val.(AutonomousSystemable))
        }
        return nil
    }
    res["countryOrRegion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountryOrRegion(val)
        }
        return nil
    }
    res["hostingProvider"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHostingProvider(val)
        }
        return nil
    }
    res["netblock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetblock(val)
        }
        return nil
    }
    return res
}
// GetHostingProvider gets the hostingProvider property value. The hosting company listed for this host.
// returns a *string when successful
func (m *IpAddress) GetHostingProvider()(*string) {
    val, err := m.GetBackingStore().Get("hostingProvider")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNetblock gets the netblock property value. The block of IP addresses this IP address belongs to.
// returns a *string when successful
func (m *IpAddress) GetNetblock()(*string) {
    val, err := m.GetBackingStore().Get("netblock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IpAddress) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Host.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("autonomousSystem", m.GetAutonomousSystem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("countryOrRegion", m.GetCountryOrRegion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("hostingProvider", m.GetHostingProvider())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("netblock", m.GetNetblock())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAutonomousSystem sets the autonomousSystem property value. The details about the autonomous system to which this IP address belongs.
func (m *IpAddress) SetAutonomousSystem(value AutonomousSystemable)() {
    err := m.GetBackingStore().Set("autonomousSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetCountryOrRegion sets the countryOrRegion property value. The country/region for this IP address.
func (m *IpAddress) SetCountryOrRegion(value *string)() {
    err := m.GetBackingStore().Set("countryOrRegion", value)
    if err != nil {
        panic(err)
    }
}
// SetHostingProvider sets the hostingProvider property value. The hosting company listed for this host.
func (m *IpAddress) SetHostingProvider(value *string)() {
    err := m.GetBackingStore().Set("hostingProvider", value)
    if err != nil {
        panic(err)
    }
}
// SetNetblock sets the netblock property value. The block of IP addresses this IP address belongs to.
func (m *IpAddress) SetNetblock(value *string)() {
    err := m.GetBackingStore().Set("netblock", value)
    if err != nil {
        panic(err)
    }
}
type IpAddressable interface {
    Hostable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAutonomousSystem()(AutonomousSystemable)
    GetCountryOrRegion()(*string)
    GetHostingProvider()(*string)
    GetNetblock()(*string)
    SetAutonomousSystem(value AutonomousSystemable)()
    SetCountryOrRegion(value *string)()
    SetHostingProvider(value *string)()
    SetNetblock(value *string)()
}

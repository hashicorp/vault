package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IPv4CidrRange struct {
    IpRange
}
// NewIPv4CidrRange instantiates a new IPv4CidrRange and sets the default values.
func NewIPv4CidrRange()(*IPv4CidrRange) {
    m := &IPv4CidrRange{
        IpRange: *NewIpRange(),
    }
    odataTypeValue := "#microsoft.graph.iPv4CidrRange"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIPv4CidrRangeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIPv4CidrRangeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIPv4CidrRange(), nil
}
// GetCidrAddress gets the cidrAddress property value. IPv4 address in CIDR notation. Not nullable.
// returns a *string when successful
func (m *IPv4CidrRange) GetCidrAddress()(*string) {
    val, err := m.GetBackingStore().Get("cidrAddress")
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
func (m *IPv4CidrRange) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IpRange.GetFieldDeserializers()
    res["cidrAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCidrAddress(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *IPv4CidrRange) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IpRange.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("cidrAddress", m.GetCidrAddress())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCidrAddress sets the cidrAddress property value. IPv4 address in CIDR notation. Not nullable.
func (m *IPv4CidrRange) SetCidrAddress(value *string)() {
    err := m.GetBackingStore().Set("cidrAddress", value)
    if err != nil {
        panic(err)
    }
}
type IPv4CidrRangeable interface {
    IpRangeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCidrAddress()(*string)
    SetCidrAddress(value *string)()
}

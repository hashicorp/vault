package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IPv6CidrRange struct {
    IpRange
}
// NewIPv6CidrRange instantiates a new IPv6CidrRange and sets the default values.
func NewIPv6CidrRange()(*IPv6CidrRange) {
    m := &IPv6CidrRange{
        IpRange: *NewIpRange(),
    }
    odataTypeValue := "#microsoft.graph.iPv6CidrRange"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIPv6CidrRangeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIPv6CidrRangeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIPv6CidrRange(), nil
}
// GetCidrAddress gets the cidrAddress property value. IPv6 address in CIDR notation. Not nullable.
// returns a *string when successful
func (m *IPv6CidrRange) GetCidrAddress()(*string) {
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
func (m *IPv6CidrRange) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
func (m *IPv6CidrRange) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetCidrAddress sets the cidrAddress property value. IPv6 address in CIDR notation. Not nullable.
func (m *IPv6CidrRange) SetCidrAddress(value *string)() {
    err := m.GetBackingStore().Set("cidrAddress", value)
    if err != nil {
        panic(err)
    }
}
type IPv6CidrRangeable interface {
    IpRangeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCidrAddress()(*string)
    SetCidrAddress(value *string)()
}

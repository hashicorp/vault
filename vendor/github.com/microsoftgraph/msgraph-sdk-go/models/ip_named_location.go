package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IpNamedLocation struct {
    NamedLocation
}
// NewIpNamedLocation instantiates a new IpNamedLocation and sets the default values.
func NewIpNamedLocation()(*IpNamedLocation) {
    m := &IpNamedLocation{
        NamedLocation: *NewNamedLocation(),
    }
    return m
}
// CreateIpNamedLocationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIpNamedLocationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIpNamedLocation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IpNamedLocation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.NamedLocation.GetFieldDeserializers()
    res["ipRanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIpRangeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IpRangeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IpRangeable)
                }
            }
            m.SetIpRanges(res)
        }
        return nil
    }
    res["isTrusted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsTrusted(val)
        }
        return nil
    }
    return res
}
// GetIpRanges gets the ipRanges property value. List of IP address ranges in IPv4 CIDR format (for example, 1.2.3.4/32) or any allowable IPv6 format from IETF RFC5969. Required.
// returns a []IpRangeable when successful
func (m *IpNamedLocation) GetIpRanges()([]IpRangeable) {
    val, err := m.GetBackingStore().Get("ipRanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IpRangeable)
    }
    return nil
}
// GetIsTrusted gets the isTrusted property value. true if this location is explicitly trusted. Optional. Default value is false.
// returns a *bool when successful
func (m *IpNamedLocation) GetIsTrusted()(*bool) {
    val, err := m.GetBackingStore().Get("isTrusted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IpNamedLocation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.NamedLocation.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetIpRanges() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIpRanges()))
        for i, v := range m.GetIpRanges() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ipRanges", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isTrusted", m.GetIsTrusted())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIpRanges sets the ipRanges property value. List of IP address ranges in IPv4 CIDR format (for example, 1.2.3.4/32) or any allowable IPv6 format from IETF RFC5969. Required.
func (m *IpNamedLocation) SetIpRanges(value []IpRangeable)() {
    err := m.GetBackingStore().Set("ipRanges", value)
    if err != nil {
        panic(err)
    }
}
// SetIsTrusted sets the isTrusted property value. true if this location is explicitly trusted. Optional. Default value is false.
func (m *IpNamedLocation) SetIsTrusted(value *bool)() {
    err := m.GetBackingStore().Set("isTrusted", value)
    if err != nil {
        panic(err)
    }
}
type IpNamedLocationable interface {
    NamedLocationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIpRanges()([]IpRangeable)
    GetIsTrusted()(*bool)
    SetIpRanges(value []IpRangeable)()
    SetIsTrusted(value *bool)()
}

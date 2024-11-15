package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IPv6Range iPv6 Range definition.
type IPv6Range struct {
    IpRange
}
// NewIPv6Range instantiates a new IPv6Range and sets the default values.
func NewIPv6Range()(*IPv6Range) {
    m := &IPv6Range{
        IpRange: *NewIpRange(),
    }
    odataTypeValue := "#microsoft.graph.iPv6Range"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIPv6RangeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIPv6RangeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIPv6Range(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IPv6Range) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IpRange.GetFieldDeserializers()
    res["lowerAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLowerAddress(val)
        }
        return nil
    }
    res["upperAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpperAddress(val)
        }
        return nil
    }
    return res
}
// GetLowerAddress gets the lowerAddress property value. Lower address.
// returns a *string when successful
func (m *IPv6Range) GetLowerAddress()(*string) {
    val, err := m.GetBackingStore().Get("lowerAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUpperAddress gets the upperAddress property value. Upper address.
// returns a *string when successful
func (m *IPv6Range) GetUpperAddress()(*string) {
    val, err := m.GetBackingStore().Get("upperAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IPv6Range) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IpRange.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("lowerAddress", m.GetLowerAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("upperAddress", m.GetUpperAddress())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLowerAddress sets the lowerAddress property value. Lower address.
func (m *IPv6Range) SetLowerAddress(value *string)() {
    err := m.GetBackingStore().Set("lowerAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetUpperAddress sets the upperAddress property value. Upper address.
func (m *IPv6Range) SetUpperAddress(value *string)() {
    err := m.GetBackingStore().Set("upperAddress", value)
    if err != nil {
        panic(err)
    }
}
type IPv6Rangeable interface {
    IpRangeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLowerAddress()(*string)
    GetUpperAddress()(*string)
    SetLowerAddress(value *string)()
    SetUpperAddress(value *string)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Place struct {
    Entity
}
// NewPlace instantiates a new Place and sets the default values.
func NewPlace()(*Place) {
    m := &Place{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlaceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlaceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.room":
                        return NewRoom(), nil
                    case "#microsoft.graph.roomList":
                        return NewRoomList(), nil
                }
            }
        }
    }
    return NewPlace(), nil
}
// GetAddress gets the address property value. The street address of the place.
// returns a PhysicalAddressable when successful
func (m *Place) GetAddress()(PhysicalAddressable) {
    val, err := m.GetBackingStore().Get("address")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PhysicalAddressable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name associated with the place.
// returns a *string when successful
func (m *Place) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *Place) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["address"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhysicalAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddress(val.(PhysicalAddressable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["geoCoordinates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOutlookGeoCoordinatesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGeoCoordinates(val.(OutlookGeoCoordinatesable))
        }
        return nil
    }
    res["phone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhone(val)
        }
        return nil
    }
    return res
}
// GetGeoCoordinates gets the geoCoordinates property value. Specifies the place location in latitude, longitude, and (optionally) altitude coordinates.
// returns a OutlookGeoCoordinatesable when successful
func (m *Place) GetGeoCoordinates()(OutlookGeoCoordinatesable) {
    val, err := m.GetBackingStore().Get("geoCoordinates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OutlookGeoCoordinatesable)
    }
    return nil
}
// GetPhone gets the phone property value. The phone number of the place.
// returns a *string when successful
func (m *Place) GetPhone()(*string) {
    val, err := m.GetBackingStore().Get("phone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Place) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("address", m.GetAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("geoCoordinates", m.GetGeoCoordinates())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("phone", m.GetPhone())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAddress sets the address property value. The street address of the place.
func (m *Place) SetAddress(value PhysicalAddressable)() {
    err := m.GetBackingStore().Set("address", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name associated with the place.
func (m *Place) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetGeoCoordinates sets the geoCoordinates property value. Specifies the place location in latitude, longitude, and (optionally) altitude coordinates.
func (m *Place) SetGeoCoordinates(value OutlookGeoCoordinatesable)() {
    err := m.GetBackingStore().Set("geoCoordinates", value)
    if err != nil {
        panic(err)
    }
}
// SetPhone sets the phone property value. The phone number of the place.
func (m *Place) SetPhone(value *string)() {
    err := m.GetBackingStore().Set("phone", value)
    if err != nil {
        panic(err)
    }
}
type Placeable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAddress()(PhysicalAddressable)
    GetDisplayName()(*string)
    GetGeoCoordinates()(OutlookGeoCoordinatesable)
    GetPhone()(*string)
    SetAddress(value PhysicalAddressable)()
    SetDisplayName(value *string)()
    SetGeoCoordinates(value OutlookGeoCoordinatesable)()
    SetPhone(value *string)()
}

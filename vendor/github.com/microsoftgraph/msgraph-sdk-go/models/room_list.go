package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RoomList struct {
    Place
}
// NewRoomList instantiates a new RoomList and sets the default values.
func NewRoomList()(*RoomList) {
    m := &RoomList{
        Place: *NewPlace(),
    }
    odataTypeValue := "#microsoft.graph.roomList"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRoomListFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRoomListFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRoomList(), nil
}
// GetEmailAddress gets the emailAddress property value. The email address of the room list.
// returns a *string when successful
func (m *RoomList) GetEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("emailAddress")
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
func (m *RoomList) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Place.GetFieldDeserializers()
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val)
        }
        return nil
    }
    res["rooms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRoomFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Roomable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Roomable)
                }
            }
            m.SetRooms(res)
        }
        return nil
    }
    return res
}
// GetRooms gets the rooms property value. The rooms property
// returns a []Roomable when successful
func (m *RoomList) GetRooms()([]Roomable) {
    val, err := m.GetBackingStore().Get("rooms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Roomable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RoomList) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Place.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    if m.GetRooms() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRooms()))
        for i, v := range m.GetRooms() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rooms", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmailAddress sets the emailAddress property value. The email address of the room list.
func (m *RoomList) SetEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetRooms sets the rooms property value. The rooms property
func (m *RoomList) SetRooms(value []Roomable)() {
    err := m.GetBackingStore().Set("rooms", value)
    if err != nil {
        panic(err)
    }
}
type RoomListable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Placeable
    GetEmailAddress()(*string)
    GetRooms()([]Roomable)
    SetEmailAddress(value *string)()
    SetRooms(value []Roomable)()
}

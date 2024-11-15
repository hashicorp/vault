package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Room struct {
    Place
}
// NewRoom instantiates a new Room and sets the default values.
func NewRoom()(*Room) {
    m := &Room{
        Place: *NewPlace(),
    }
    odataTypeValue := "#microsoft.graph.room"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRoomFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRoomFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRoom(), nil
}
// GetAudioDeviceName gets the audioDeviceName property value. Specifies the name of the audio device in the room.
// returns a *string when successful
func (m *Room) GetAudioDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("audioDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBookingType gets the bookingType property value. Type of room. Possible values are standard, and reserved.
// returns a *BookingType when successful
func (m *Room) GetBookingType()(*BookingType) {
    val, err := m.GetBackingStore().Get("bookingType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BookingType)
    }
    return nil
}
// GetBuilding gets the building property value. Specifies the building name or building number that the room is in.
// returns a *string when successful
func (m *Room) GetBuilding()(*string) {
    val, err := m.GetBackingStore().Get("building")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCapacity gets the capacity property value. Specifies the capacity of the room.
// returns a *int32 when successful
func (m *Room) GetCapacity()(*int32) {
    val, err := m.GetBackingStore().Get("capacity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDisplayDeviceName gets the displayDeviceName property value. Specifies the name of the display device in the room.
// returns a *string when successful
func (m *Room) GetDisplayDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("displayDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. Email address of the room.
// returns a *string when successful
func (m *Room) GetEmailAddress()(*string) {
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
func (m *Room) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Place.GetFieldDeserializers()
    res["audioDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudioDeviceName(val)
        }
        return nil
    }
    res["bookingType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBookingType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBookingType(val.(*BookingType))
        }
        return nil
    }
    res["building"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBuilding(val)
        }
        return nil
    }
    res["capacity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCapacity(val)
        }
        return nil
    }
    res["displayDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayDeviceName(val)
        }
        return nil
    }
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
    res["floorLabel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFloorLabel(val)
        }
        return nil
    }
    res["floorNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFloorNumber(val)
        }
        return nil
    }
    res["isWheelChairAccessible"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsWheelChairAccessible(val)
        }
        return nil
    }
    res["label"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabel(val)
        }
        return nil
    }
    res["nickname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNickname(val)
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetTags(res)
        }
        return nil
    }
    res["videoDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVideoDeviceName(val)
        }
        return nil
    }
    return res
}
// GetFloorLabel gets the floorLabel property value. Specifies a descriptive label for the floor, for example, P.
// returns a *string when successful
func (m *Room) GetFloorLabel()(*string) {
    val, err := m.GetBackingStore().Get("floorLabel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFloorNumber gets the floorNumber property value. Specifies the floor number that the room is on.
// returns a *int32 when successful
func (m *Room) GetFloorNumber()(*int32) {
    val, err := m.GetBackingStore().Get("floorNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetIsWheelChairAccessible gets the isWheelChairAccessible property value. Specifies whether the room is wheelchair accessible.
// returns a *bool when successful
func (m *Room) GetIsWheelChairAccessible()(*bool) {
    val, err := m.GetBackingStore().Get("isWheelChairAccessible")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLabel gets the label property value. Specifies a descriptive label for the room, for example, a number or name.
// returns a *string when successful
func (m *Room) GetLabel()(*string) {
    val, err := m.GetBackingStore().Get("label")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNickname gets the nickname property value. Specifies a nickname for the room, for example, 'conf room'.
// returns a *string when successful
func (m *Room) GetNickname()(*string) {
    val, err := m.GetBackingStore().Get("nickname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTags gets the tags property value. Specifies other features of the room, for example, details like the type of view or furniture type.
// returns a []string when successful
func (m *Room) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetVideoDeviceName gets the videoDeviceName property value. Specifies the name of the video device in the room.
// returns a *string when successful
func (m *Room) GetVideoDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("videoDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Room) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Place.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("audioDeviceName", m.GetAudioDeviceName())
        if err != nil {
            return err
        }
    }
    if m.GetBookingType() != nil {
        cast := (*m.GetBookingType()).String()
        err = writer.WriteStringValue("bookingType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("building", m.GetBuilding())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("capacity", m.GetCapacity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayDeviceName", m.GetDisplayDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("floorLabel", m.GetFloorLabel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("floorNumber", m.GetFloorNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isWheelChairAccessible", m.GetIsWheelChairAccessible())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("label", m.GetLabel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("nickname", m.GetNickname())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err = writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("videoDeviceName", m.GetVideoDeviceName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAudioDeviceName sets the audioDeviceName property value. Specifies the name of the audio device in the room.
func (m *Room) SetAudioDeviceName(value *string)() {
    err := m.GetBackingStore().Set("audioDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetBookingType sets the bookingType property value. Type of room. Possible values are standard, and reserved.
func (m *Room) SetBookingType(value *BookingType)() {
    err := m.GetBackingStore().Set("bookingType", value)
    if err != nil {
        panic(err)
    }
}
// SetBuilding sets the building property value. Specifies the building name or building number that the room is in.
func (m *Room) SetBuilding(value *string)() {
    err := m.GetBackingStore().Set("building", value)
    if err != nil {
        panic(err)
    }
}
// SetCapacity sets the capacity property value. Specifies the capacity of the room.
func (m *Room) SetCapacity(value *int32)() {
    err := m.GetBackingStore().Set("capacity", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayDeviceName sets the displayDeviceName property value. Specifies the name of the display device in the room.
func (m *Room) SetDisplayDeviceName(value *string)() {
    err := m.GetBackingStore().Set("displayDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. Email address of the room.
func (m *Room) SetEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetFloorLabel sets the floorLabel property value. Specifies a descriptive label for the floor, for example, P.
func (m *Room) SetFloorLabel(value *string)() {
    err := m.GetBackingStore().Set("floorLabel", value)
    if err != nil {
        panic(err)
    }
}
// SetFloorNumber sets the floorNumber property value. Specifies the floor number that the room is on.
func (m *Room) SetFloorNumber(value *int32)() {
    err := m.GetBackingStore().Set("floorNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetIsWheelChairAccessible sets the isWheelChairAccessible property value. Specifies whether the room is wheelchair accessible.
func (m *Room) SetIsWheelChairAccessible(value *bool)() {
    err := m.GetBackingStore().Set("isWheelChairAccessible", value)
    if err != nil {
        panic(err)
    }
}
// SetLabel sets the label property value. Specifies a descriptive label for the room, for example, a number or name.
func (m *Room) SetLabel(value *string)() {
    err := m.GetBackingStore().Set("label", value)
    if err != nil {
        panic(err)
    }
}
// SetNickname sets the nickname property value. Specifies a nickname for the room, for example, 'conf room'.
func (m *Room) SetNickname(value *string)() {
    err := m.GetBackingStore().Set("nickname", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Specifies other features of the room, for example, details like the type of view or furniture type.
func (m *Room) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetVideoDeviceName sets the videoDeviceName property value. Specifies the name of the video device in the room.
func (m *Room) SetVideoDeviceName(value *string)() {
    err := m.GetBackingStore().Set("videoDeviceName", value)
    if err != nil {
        panic(err)
    }
}
type Roomable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Placeable
    GetAudioDeviceName()(*string)
    GetBookingType()(*BookingType)
    GetBuilding()(*string)
    GetCapacity()(*int32)
    GetDisplayDeviceName()(*string)
    GetEmailAddress()(*string)
    GetFloorLabel()(*string)
    GetFloorNumber()(*int32)
    GetIsWheelChairAccessible()(*bool)
    GetLabel()(*string)
    GetNickname()(*string)
    GetTags()([]string)
    GetVideoDeviceName()(*string)
    SetAudioDeviceName(value *string)()
    SetBookingType(value *BookingType)()
    SetBuilding(value *string)()
    SetCapacity(value *int32)()
    SetDisplayDeviceName(value *string)()
    SetEmailAddress(value *string)()
    SetFloorLabel(value *string)()
    SetFloorNumber(value *int32)()
    SetIsWheelChairAccessible(value *bool)()
    SetLabel(value *string)()
    SetNickname(value *string)()
    SetTags(value []string)()
    SetVideoDeviceName(value *string)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ChannelAddedEventMessageDetail struct {
    EventMessageDetail
}
// NewChannelAddedEventMessageDetail instantiates a new ChannelAddedEventMessageDetail and sets the default values.
func NewChannelAddedEventMessageDetail()(*ChannelAddedEventMessageDetail) {
    m := &ChannelAddedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.channelAddedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateChannelAddedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChannelAddedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChannelAddedEventMessageDetail(), nil
}
// GetChannelDisplayName gets the channelDisplayName property value. Display name of the channel.
// returns a *string when successful
func (m *ChannelAddedEventMessageDetail) GetChannelDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("channelDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetChannelId gets the channelId property value. Unique identifier of the channel.
// returns a *string when successful
func (m *ChannelAddedEventMessageDetail) GetChannelId()(*string) {
    val, err := m.GetBackingStore().Get("channelId")
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
func (m *ChannelAddedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["channelDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChannelDisplayName(val)
        }
        return nil
    }
    res["channelId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChannelId(val)
        }
        return nil
    }
    res["initiator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiator(val.(IdentitySetable))
        }
        return nil
    }
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *ChannelAddedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ChannelAddedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("channelDisplayName", m.GetChannelDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("channelId", m.GetChannelId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChannelDisplayName sets the channelDisplayName property value. Display name of the channel.
func (m *ChannelAddedEventMessageDetail) SetChannelDisplayName(value *string)() {
    err := m.GetBackingStore().Set("channelDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetChannelId sets the channelId property value. Unique identifier of the channel.
func (m *ChannelAddedEventMessageDetail) SetChannelId(value *string)() {
    err := m.GetBackingStore().Set("channelId", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *ChannelAddedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
type ChannelAddedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChannelDisplayName()(*string)
    GetChannelId()(*string)
    GetInitiator()(IdentitySetable)
    SetChannelDisplayName(value *string)()
    SetChannelId(value *string)()
    SetInitiator(value IdentitySetable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PinnedChatMessageInfo struct {
    Entity
}
// NewPinnedChatMessageInfo instantiates a new PinnedChatMessageInfo and sets the default values.
func NewPinnedChatMessageInfo()(*PinnedChatMessageInfo) {
    m := &PinnedChatMessageInfo{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePinnedChatMessageInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePinnedChatMessageInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPinnedChatMessageInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PinnedChatMessageInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["message"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessage(val.(ChatMessageable))
        }
        return nil
    }
    return res
}
// GetMessage gets the message property value. Represents details about the chat message that is pinned.
// returns a ChatMessageable when successful
func (m *PinnedChatMessageInfo) GetMessage()(ChatMessageable) {
    val, err := m.GetBackingStore().Get("message")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PinnedChatMessageInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("message", m.GetMessage())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMessage sets the message property value. Represents details about the chat message that is pinned.
func (m *PinnedChatMessageInfo) SetMessage(value ChatMessageable)() {
    err := m.GetBackingStore().Set("message", value)
    if err != nil {
        panic(err)
    }
}
type PinnedChatMessageInfoable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMessage()(ChatMessageable)
    SetMessage(value ChatMessageable)()
}

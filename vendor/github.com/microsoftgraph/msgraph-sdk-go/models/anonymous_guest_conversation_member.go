package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AnonymousGuestConversationMember struct {
    ConversationMember
}
// NewAnonymousGuestConversationMember instantiates a new AnonymousGuestConversationMember and sets the default values.
func NewAnonymousGuestConversationMember()(*AnonymousGuestConversationMember) {
    m := &AnonymousGuestConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.anonymousGuestConversationMember"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAnonymousGuestConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAnonymousGuestConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAnonymousGuestConversationMember(), nil
}
// GetAnonymousGuestId gets the anonymousGuestId property value. Unique ID that represents the user. Note: This ID can change if the user leaves and rejoins the meeting, or joins from a different device.
// returns a *string when successful
func (m *AnonymousGuestConversationMember) GetAnonymousGuestId()(*string) {
    val, err := m.GetBackingStore().Get("anonymousGuestId")
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
func (m *AnonymousGuestConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["anonymousGuestId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnonymousGuestId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AnonymousGuestConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("anonymousGuestId", m.GetAnonymousGuestId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnonymousGuestId sets the anonymousGuestId property value. Unique ID that represents the user. Note: This ID can change if the user leaves and rejoins the meeting, or joins from a different device.
func (m *AnonymousGuestConversationMember) SetAnonymousGuestId(value *string)() {
    err := m.GetBackingStore().Set("anonymousGuestId", value)
    if err != nil {
        panic(err)
    }
}
type AnonymousGuestConversationMemberable interface {
    ConversationMemberable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnonymousGuestId()(*string)
    SetAnonymousGuestId(value *string)()
}

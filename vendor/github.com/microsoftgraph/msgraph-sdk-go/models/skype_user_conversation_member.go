package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SkypeUserConversationMember struct {
    ConversationMember
}
// NewSkypeUserConversationMember instantiates a new SkypeUserConversationMember and sets the default values.
func NewSkypeUserConversationMember()(*SkypeUserConversationMember) {
    m := &SkypeUserConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.skypeUserConversationMember"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSkypeUserConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSkypeUserConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSkypeUserConversationMember(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SkypeUserConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["skypeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSkypeId(val)
        }
        return nil
    }
    return res
}
// GetSkypeId gets the skypeId property value. Skype ID of the user.
// returns a *string when successful
func (m *SkypeUserConversationMember) GetSkypeId()(*string) {
    val, err := m.GetBackingStore().Get("skypeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SkypeUserConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("skypeId", m.GetSkypeId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSkypeId sets the skypeId property value. Skype ID of the user.
func (m *SkypeUserConversationMember) SetSkypeId(value *string)() {
    err := m.GetBackingStore().Set("skypeId", value)
    if err != nil {
        panic(err)
    }
}
type SkypeUserConversationMemberable interface {
    ConversationMemberable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSkypeId()(*string)
    SetSkypeId(value *string)()
}

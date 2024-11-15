package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConversationMemberRoleUpdatedEventMessageDetail struct {
    EventMessageDetail
}
// NewConversationMemberRoleUpdatedEventMessageDetail instantiates a new ConversationMemberRoleUpdatedEventMessageDetail and sets the default values.
func NewConversationMemberRoleUpdatedEventMessageDetail()(*ConversationMemberRoleUpdatedEventMessageDetail) {
    m := &ConversationMemberRoleUpdatedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.conversationMemberRoleUpdatedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateConversationMemberRoleUpdatedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConversationMemberRoleUpdatedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConversationMemberRoleUpdatedEventMessageDetail(), nil
}
// GetConversationMemberRoles gets the conversationMemberRoles property value. Roles for the coversation member user.
// returns a []string when successful
func (m *ConversationMemberRoleUpdatedEventMessageDetail) GetConversationMemberRoles()([]string) {
    val, err := m.GetBackingStore().Get("conversationMemberRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetConversationMemberUser gets the conversationMemberUser property value. Identity of the conversation member user.
// returns a TeamworkUserIdentityable when successful
func (m *ConversationMemberRoleUpdatedEventMessageDetail) GetConversationMemberUser()(TeamworkUserIdentityable) {
    val, err := m.GetBackingStore().Get("conversationMemberUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamworkUserIdentityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConversationMemberRoleUpdatedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["conversationMemberRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetConversationMemberRoles(res)
        }
        return nil
    }
    res["conversationMemberUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamworkUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConversationMemberUser(val.(TeamworkUserIdentityable))
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
func (m *ConversationMemberRoleUpdatedEventMessageDetail) GetInitiator()(IdentitySetable) {
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
func (m *ConversationMemberRoleUpdatedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetConversationMemberRoles() != nil {
        err = writer.WriteCollectionOfStringValues("conversationMemberRoles", m.GetConversationMemberRoles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conversationMemberUser", m.GetConversationMemberUser())
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
// SetConversationMemberRoles sets the conversationMemberRoles property value. Roles for the coversation member user.
func (m *ConversationMemberRoleUpdatedEventMessageDetail) SetConversationMemberRoles(value []string)() {
    err := m.GetBackingStore().Set("conversationMemberRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetConversationMemberUser sets the conversationMemberUser property value. Identity of the conversation member user.
func (m *ConversationMemberRoleUpdatedEventMessageDetail) SetConversationMemberUser(value TeamworkUserIdentityable)() {
    err := m.GetBackingStore().Set("conversationMemberUser", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *ConversationMemberRoleUpdatedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
type ConversationMemberRoleUpdatedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConversationMemberRoles()([]string)
    GetConversationMemberUser()(TeamworkUserIdentityable)
    GetInitiator()(IdentitySetable)
    SetConversationMemberRoles(value []string)()
    SetConversationMemberUser(value TeamworkUserIdentityable)()
    SetInitiator(value IdentitySetable)()
}

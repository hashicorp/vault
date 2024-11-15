package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamworkConversationIdentity struct {
    Identity
}
// NewTeamworkConversationIdentity instantiates a new TeamworkConversationIdentity and sets the default values.
func NewTeamworkConversationIdentity()(*TeamworkConversationIdentity) {
    m := &TeamworkConversationIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.teamworkConversationIdentity"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTeamworkConversationIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamworkConversationIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkConversationIdentity(), nil
}
// GetConversationIdentityType gets the conversationIdentityType property value. Type of conversation. Possible values are: team, channel, chat, and unknownFutureValue.
// returns a *TeamworkConversationIdentityType when successful
func (m *TeamworkConversationIdentity) GetConversationIdentityType()(*TeamworkConversationIdentityType) {
    val, err := m.GetBackingStore().Get("conversationIdentityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamworkConversationIdentityType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamworkConversationIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["conversationIdentityType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamworkConversationIdentityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConversationIdentityType(val.(*TeamworkConversationIdentityType))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TeamworkConversationIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetConversationIdentityType() != nil {
        cast := (*m.GetConversationIdentityType()).String()
        err = writer.WriteStringValue("conversationIdentityType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConversationIdentityType sets the conversationIdentityType property value. Type of conversation. Possible values are: team, channel, chat, and unknownFutureValue.
func (m *TeamworkConversationIdentity) SetConversationIdentityType(value *TeamworkConversationIdentityType)() {
    err := m.GetBackingStore().Set("conversationIdentityType", value)
    if err != nil {
        panic(err)
    }
}
type TeamworkConversationIdentityable interface {
    Identityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConversationIdentityType()(*TeamworkConversationIdentityType)
    SetConversationIdentityType(value *TeamworkConversationIdentityType)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharedWithChannelTeamInfo struct {
    TeamInfo
}
// NewSharedWithChannelTeamInfo instantiates a new SharedWithChannelTeamInfo and sets the default values.
func NewSharedWithChannelTeamInfo()(*SharedWithChannelTeamInfo) {
    m := &SharedWithChannelTeamInfo{
        TeamInfo: *NewTeamInfo(),
    }
    return m
}
// CreateSharedWithChannelTeamInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharedWithChannelTeamInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedWithChannelTeamInfo(), nil
}
// GetAllowedMembers gets the allowedMembers property value. A collection of team members who have access to the shared channel.
// returns a []ConversationMemberable when successful
func (m *SharedWithChannelTeamInfo) GetAllowedMembers()([]ConversationMemberable) {
    val, err := m.GetBackingStore().Get("allowedMembers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConversationMemberable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharedWithChannelTeamInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TeamInfo.GetFieldDeserializers()
    res["allowedMembers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateConversationMemberFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConversationMemberable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ConversationMemberable)
                }
            }
            m.SetAllowedMembers(res)
        }
        return nil
    }
    res["isHostTeam"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsHostTeam(val)
        }
        return nil
    }
    return res
}
// GetIsHostTeam gets the isHostTeam property value. Indicates whether the team is the host of the channel.
// returns a *bool when successful
func (m *SharedWithChannelTeamInfo) GetIsHostTeam()(*bool) {
    val, err := m.GetBackingStore().Get("isHostTeam")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharedWithChannelTeamInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TeamInfo.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllowedMembers()))
        for i, v := range m.GetAllowedMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("allowedMembers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isHostTeam", m.GetIsHostTeam())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedMembers sets the allowedMembers property value. A collection of team members who have access to the shared channel.
func (m *SharedWithChannelTeamInfo) SetAllowedMembers(value []ConversationMemberable)() {
    err := m.GetBackingStore().Set("allowedMembers", value)
    if err != nil {
        panic(err)
    }
}
// SetIsHostTeam sets the isHostTeam property value. Indicates whether the team is the host of the channel.
func (m *SharedWithChannelTeamInfo) SetIsHostTeam(value *bool)() {
    err := m.GetBackingStore().Set("isHostTeam", value)
    if err != nil {
        panic(err)
    }
}
type SharedWithChannelTeamInfoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    TeamInfoable
    GetAllowedMembers()([]ConversationMemberable)
    GetIsHostTeam()(*bool)
    SetAllowedMembers(value []ConversationMemberable)()
    SetIsHostTeam(value *bool)()
}

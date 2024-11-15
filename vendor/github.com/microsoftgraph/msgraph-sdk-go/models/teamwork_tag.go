package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamworkTag struct {
    Entity
}
// NewTeamworkTag instantiates a new TeamworkTag and sets the default values.
func NewTeamworkTag()(*TeamworkTag) {
    m := &TeamworkTag{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamworkTagFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamworkTagFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkTag(), nil
}
// GetDescription gets the description property value. The description of the tag as it appears to the user in Microsoft Teams. A teamworkTag can't have more than 200 teamworkTagMembers.
// returns a *string when successful
func (m *TeamworkTag) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the tag as it appears to the user in Microsoft Teams.
// returns a *string when successful
func (m *TeamworkTag) GetDisplayName()(*string) {
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
func (m *TeamworkTag) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["memberCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMemberCount(val)
        }
        return nil
    }
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamworkTagMemberFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamworkTagMemberable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamworkTagMemberable)
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    res["tagType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamworkTagType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTagType(val.(*TeamworkTagType))
        }
        return nil
    }
    res["teamId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamId(val)
        }
        return nil
    }
    return res
}
// GetMemberCount gets the memberCount property value. The number of users assigned to the tag.
// returns a *int32 when successful
func (m *TeamworkTag) GetMemberCount()(*int32) {
    val, err := m.GetBackingStore().Get("memberCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMembers gets the members property value. Users assigned to the tag.
// returns a []TeamworkTagMemberable when successful
func (m *TeamworkTag) GetMembers()([]TeamworkTagMemberable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamworkTagMemberable)
    }
    return nil
}
// GetTagType gets the tagType property value. The type of the tag. Default is standard.
// returns a *TeamworkTagType when successful
func (m *TeamworkTag) GetTagType()(*TeamworkTagType) {
    val, err := m.GetBackingStore().Get("tagType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamworkTagType)
    }
    return nil
}
// GetTeamId gets the teamId property value. ID of the team in which the tag is defined.
// returns a *string when successful
func (m *TeamworkTag) GetTeamId()(*string) {
    val, err := m.GetBackingStore().Get("teamId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamworkTag) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteInt32Value("memberCount", m.GetMemberCount())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMembers()))
        for i, v := range m.GetMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTagType() != nil {
        cast := (*m.GetTagType()).String()
        err = writer.WriteStringValue("tagType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamId", m.GetTeamId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. The description of the tag as it appears to the user in Microsoft Teams. A teamworkTag can't have more than 200 teamworkTagMembers.
func (m *TeamworkTag) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the tag as it appears to the user in Microsoft Teams.
func (m *TeamworkTag) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberCount sets the memberCount property value. The number of users assigned to the tag.
func (m *TeamworkTag) SetMemberCount(value *int32)() {
    err := m.GetBackingStore().Set("memberCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. Users assigned to the tag.
func (m *TeamworkTag) SetMembers(value []TeamworkTagMemberable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetTagType sets the tagType property value. The type of the tag. Default is standard.
func (m *TeamworkTag) SetTagType(value *TeamworkTagType)() {
    err := m.GetBackingStore().Set("tagType", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamId sets the teamId property value. ID of the team in which the tag is defined.
func (m *TeamworkTag) SetTeamId(value *string)() {
    err := m.GetBackingStore().Set("teamId", value)
    if err != nil {
        panic(err)
    }
}
type TeamworkTagable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetMemberCount()(*int32)
    GetMembers()([]TeamworkTagMemberable)
    GetTagType()(*TeamworkTagType)
    GetTeamId()(*string)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetMemberCount(value *int32)()
    SetMembers(value []TeamworkTagMemberable)()
    SetTagType(value *TeamworkTagType)()
    SetTeamId(value *string)()
}

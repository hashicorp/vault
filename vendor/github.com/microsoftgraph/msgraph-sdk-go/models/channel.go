package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Channel struct {
    Entity
}
// NewChannel instantiates a new Channel and sets the default values.
func NewChannel()(*Channel) {
    m := &Channel{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChannelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChannelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChannel(), nil
}
// GetCreatedDateTime gets the createdDateTime property value. Read only. Timestamp at which the channel was created.
// returns a *Time when successful
func (m *Channel) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Optional textual description for the channel.
// returns a *string when successful
func (m *Channel) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Channel name as it will appear to the user in Microsoft Teams. The maximum length is 50 characters.
// returns a *string when successful
func (m *Channel) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmail gets the email property value. The email address for sending messages to the channel. Read-only.
// returns a *string when successful
func (m *Channel) GetEmail()(*string) {
    val, err := m.GetBackingStore().Get("email")
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
func (m *Channel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
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
    res["email"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmail(val)
        }
        return nil
    }
    res["filesFolder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFilesFolder(val.(DriveItemable))
        }
        return nil
    }
    res["isArchived"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsArchived(val)
        }
        return nil
    }
    res["isFavoriteByDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFavoriteByDefault(val)
        }
        return nil
    }
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetMembers(res)
        }
        return nil
    }
    res["membershipType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChannelMembershipType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMembershipType(val.(*ChannelMembershipType))
        }
        return nil
    }
    res["messages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageable)
                }
            }
            m.SetMessages(res)
        }
        return nil
    }
    res["sharedWithTeams"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSharedWithChannelTeamInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SharedWithChannelTeamInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SharedWithChannelTeamInfoable)
                }
            }
            m.SetSharedWithTeams(res)
        }
        return nil
    }
    res["summary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChannelSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSummary(val.(ChannelSummaryable))
        }
        return nil
    }
    res["tabs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamsTabFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamsTabable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamsTabable)
                }
            }
            m.SetTabs(res)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetFilesFolder gets the filesFolder property value. Metadata for the location where the channel's files are stored.
// returns a DriveItemable when successful
func (m *Channel) GetFilesFolder()(DriveItemable) {
    val, err := m.GetBackingStore().Get("filesFolder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemable)
    }
    return nil
}
// GetIsArchived gets the isArchived property value. Indicates whether the channel is archived. Read-only.
// returns a *bool when successful
func (m *Channel) GetIsArchived()(*bool) {
    val, err := m.GetBackingStore().Get("isArchived")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsFavoriteByDefault gets the isFavoriteByDefault property value. Indicates whether the channel should be marked as recommended for all members of the team to show in their channel list. Note: All recommended channels automatically show in the channels list for education and frontline worker users. The property can only be set programmatically via the Create team method. The default value is false.
// returns a *bool when successful
func (m *Channel) GetIsFavoriteByDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isFavoriteByDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMembers gets the members property value. A collection of membership records associated with the channel.
// returns a []ConversationMemberable when successful
func (m *Channel) GetMembers()([]ConversationMemberable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConversationMemberable)
    }
    return nil
}
// GetMembershipType gets the membershipType property value. The type of the channel. Can be set during creation and can't be changed. The possible values are: standard, private, unknownFutureValue, shared. The default value is standard. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value in this evolvable enum: shared.
// returns a *ChannelMembershipType when successful
func (m *Channel) GetMembershipType()(*ChannelMembershipType) {
    val, err := m.GetBackingStore().Get("membershipType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChannelMembershipType)
    }
    return nil
}
// GetMessages gets the messages property value. A collection of all the messages in the channel. A navigation property. Nullable.
// returns a []ChatMessageable when successful
func (m *Channel) GetMessages()([]ChatMessageable) {
    val, err := m.GetBackingStore().Get("messages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageable)
    }
    return nil
}
// GetSharedWithTeams gets the sharedWithTeams property value. A collection of teams with which a channel is shared.
// returns a []SharedWithChannelTeamInfoable when successful
func (m *Channel) GetSharedWithTeams()([]SharedWithChannelTeamInfoable) {
    val, err := m.GetBackingStore().Get("sharedWithTeams")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SharedWithChannelTeamInfoable)
    }
    return nil
}
// GetSummary gets the summary property value. Contains summary information about the channel, including number of owners, members, guests, and an indicator for members from other tenants. The summary property will only be returned if it is specified in the $select clause of the Get channel method.
// returns a ChannelSummaryable when successful
func (m *Channel) GetSummary()(ChannelSummaryable) {
    val, err := m.GetBackingStore().Get("summary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChannelSummaryable)
    }
    return nil
}
// GetTabs gets the tabs property value. A collection of all the tabs in the channel. A navigation property.
// returns a []TeamsTabable when successful
func (m *Channel) GetTabs()([]TeamsTabable) {
    val, err := m.GetBackingStore().Get("tabs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsTabable)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The ID of the Microsoft Entra tenant.
// returns a *string when successful
func (m *Channel) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. A hyperlink that will go to the channel in Microsoft Teams. This is the URL that you get when you right-click a channel in Microsoft Teams and select Get link to channel. This URL should be treated as an opaque blob, and not parsed. Read-only.
// returns a *string when successful
func (m *Channel) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Channel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
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
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("filesFolder", m.GetFilesFolder())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isArchived", m.GetIsArchived())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFavoriteByDefault", m.GetIsFavoriteByDefault())
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
    if m.GetMembershipType() != nil {
        cast := (*m.GetMembershipType()).String()
        err = writer.WriteStringValue("membershipType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessages()))
        for i, v := range m.GetMessages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSharedWithTeams() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharedWithTeams()))
        for i, v := range m.GetSharedWithTeams() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharedWithTeams", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("summary", m.GetSummary())
        if err != nil {
            return err
        }
    }
    if m.GetTabs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTabs()))
        for i, v := range m.GetTabs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tabs", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCreatedDateTime sets the createdDateTime property value. Read only. Timestamp at which the channel was created.
func (m *Channel) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Optional textual description for the channel.
func (m *Channel) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Channel name as it will appear to the user in Microsoft Teams. The maximum length is 50 characters.
func (m *Channel) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEmail sets the email property value. The email address for sending messages to the channel. Read-only.
func (m *Channel) SetEmail(value *string)() {
    err := m.GetBackingStore().Set("email", value)
    if err != nil {
        panic(err)
    }
}
// SetFilesFolder sets the filesFolder property value. Metadata for the location where the channel's files are stored.
func (m *Channel) SetFilesFolder(value DriveItemable)() {
    err := m.GetBackingStore().Set("filesFolder", value)
    if err != nil {
        panic(err)
    }
}
// SetIsArchived sets the isArchived property value. Indicates whether the channel is archived. Read-only.
func (m *Channel) SetIsArchived(value *bool)() {
    err := m.GetBackingStore().Set("isArchived", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFavoriteByDefault sets the isFavoriteByDefault property value. Indicates whether the channel should be marked as recommended for all members of the team to show in their channel list. Note: All recommended channels automatically show in the channels list for education and frontline worker users. The property can only be set programmatically via the Create team method. The default value is false.
func (m *Channel) SetIsFavoriteByDefault(value *bool)() {
    err := m.GetBackingStore().Set("isFavoriteByDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. A collection of membership records associated with the channel.
func (m *Channel) SetMembers(value []ConversationMemberable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetMembershipType sets the membershipType property value. The type of the channel. Can be set during creation and can't be changed. The possible values are: standard, private, unknownFutureValue, shared. The default value is standard. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value in this evolvable enum: shared.
func (m *Channel) SetMembershipType(value *ChannelMembershipType)() {
    err := m.GetBackingStore().Set("membershipType", value)
    if err != nil {
        panic(err)
    }
}
// SetMessages sets the messages property value. A collection of all the messages in the channel. A navigation property. Nullable.
func (m *Channel) SetMessages(value []ChatMessageable)() {
    err := m.GetBackingStore().Set("messages", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedWithTeams sets the sharedWithTeams property value. A collection of teams with which a channel is shared.
func (m *Channel) SetSharedWithTeams(value []SharedWithChannelTeamInfoable)() {
    err := m.GetBackingStore().Set("sharedWithTeams", value)
    if err != nil {
        panic(err)
    }
}
// SetSummary sets the summary property value. Contains summary information about the channel, including number of owners, members, guests, and an indicator for members from other tenants. The summary property will only be returned if it is specified in the $select clause of the Get channel method.
func (m *Channel) SetSummary(value ChannelSummaryable)() {
    err := m.GetBackingStore().Set("summary", value)
    if err != nil {
        panic(err)
    }
}
// SetTabs sets the tabs property value. A collection of all the tabs in the channel. A navigation property.
func (m *Channel) SetTabs(value []TeamsTabable)() {
    err := m.GetBackingStore().Set("tabs", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The ID of the Microsoft Entra tenant.
func (m *Channel) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. A hyperlink that will go to the channel in Microsoft Teams. This is the URL that you get when you right-click a channel in Microsoft Teams and select Get link to channel. This URL should be treated as an opaque blob, and not parsed. Read-only.
func (m *Channel) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type Channelable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetEmail()(*string)
    GetFilesFolder()(DriveItemable)
    GetIsArchived()(*bool)
    GetIsFavoriteByDefault()(*bool)
    GetMembers()([]ConversationMemberable)
    GetMembershipType()(*ChannelMembershipType)
    GetMessages()([]ChatMessageable)
    GetSharedWithTeams()([]SharedWithChannelTeamInfoable)
    GetSummary()(ChannelSummaryable)
    GetTabs()([]TeamsTabable)
    GetTenantId()(*string)
    GetWebUrl()(*string)
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetEmail(value *string)()
    SetFilesFolder(value DriveItemable)()
    SetIsArchived(value *bool)()
    SetIsFavoriteByDefault(value *bool)()
    SetMembers(value []ConversationMemberable)()
    SetMembershipType(value *ChannelMembershipType)()
    SetMessages(value []ChatMessageable)()
    SetSharedWithTeams(value []SharedWithChannelTeamInfoable)()
    SetSummary(value ChannelSummaryable)()
    SetTabs(value []TeamsTabable)()
    SetTenantId(value *string)()
    SetWebUrl(value *string)()
}

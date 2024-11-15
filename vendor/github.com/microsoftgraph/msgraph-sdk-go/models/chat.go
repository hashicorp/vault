package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Chat struct {
    Entity
}
// NewChat instantiates a new Chat and sets the default values.
func NewChat()(*Chat) {
    m := &Chat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChat(), nil
}
// GetChatType gets the chatType property value. The chatType property
// returns a *ChatType when successful
func (m *Chat) GetChatType()(*ChatType) {
    val, err := m.GetBackingStore().Get("chatType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatType)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time at which the chat was created. Read-only.
// returns a *Time when successful
func (m *Chat) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Chat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["chatType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChatType(val.(*ChatType))
        }
        return nil
    }
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
    res["installedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamsAppInstallationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamsAppInstallationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamsAppInstallationable)
                }
            }
            m.SetInstalledApps(res)
        }
        return nil
    }
    res["lastMessagePreview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatMessageInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastMessagePreview(val.(ChatMessageInfoable))
        }
        return nil
    }
    res["lastUpdatedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdatedDateTime(val)
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
    res["onlineMeetingInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamworkOnlineMeetingInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnlineMeetingInfo(val.(TeamworkOnlineMeetingInfoable))
        }
        return nil
    }
    res["permissionGrants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateResourceSpecificPermissionGrantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ResourceSpecificPermissionGrantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ResourceSpecificPermissionGrantable)
                }
            }
            m.SetPermissionGrants(res)
        }
        return nil
    }
    res["pinnedMessages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePinnedChatMessageInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PinnedChatMessageInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PinnedChatMessageInfoable)
                }
            }
            m.SetPinnedMessages(res)
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
    res["topic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTopic(val)
        }
        return nil
    }
    res["viewpoint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatViewpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetViewpoint(val.(ChatViewpointable))
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
// GetInstalledApps gets the installedApps property value. A collection of all the apps in the chat. Nullable.
// returns a []TeamsAppInstallationable when successful
func (m *Chat) GetInstalledApps()([]TeamsAppInstallationable) {
    val, err := m.GetBackingStore().Get("installedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsAppInstallationable)
    }
    return nil
}
// GetLastMessagePreview gets the lastMessagePreview property value. Preview of the last message sent in the chat. Null if no messages were sent in the chat. Currently, only the list chats operation supports this property.
// returns a ChatMessageInfoable when successful
func (m *Chat) GetLastMessagePreview()(ChatMessageInfoable) {
    val, err := m.GetBackingStore().Get("lastMessagePreview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessageInfoable)
    }
    return nil
}
// GetLastUpdatedDateTime gets the lastUpdatedDateTime property value. Date and time at which the chat was renamed or the list of members was last changed. Read-only.
// returns a *Time when successful
func (m *Chat) GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdatedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMembers gets the members property value. A collection of all the members in the chat. Nullable.
// returns a []ConversationMemberable when successful
func (m *Chat) GetMembers()([]ConversationMemberable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConversationMemberable)
    }
    return nil
}
// GetMessages gets the messages property value. A collection of all the messages in the chat. Nullable.
// returns a []ChatMessageable when successful
func (m *Chat) GetMessages()([]ChatMessageable) {
    val, err := m.GetBackingStore().Get("messages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageable)
    }
    return nil
}
// GetOnlineMeetingInfo gets the onlineMeetingInfo property value. Represents details about an online meeting. If the chat isn't associated with an online meeting, the property is empty. Read-only.
// returns a TeamworkOnlineMeetingInfoable when successful
func (m *Chat) GetOnlineMeetingInfo()(TeamworkOnlineMeetingInfoable) {
    val, err := m.GetBackingStore().Get("onlineMeetingInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamworkOnlineMeetingInfoable)
    }
    return nil
}
// GetPermissionGrants gets the permissionGrants property value. A collection of permissions granted to apps for the chat.
// returns a []ResourceSpecificPermissionGrantable when successful
func (m *Chat) GetPermissionGrants()([]ResourceSpecificPermissionGrantable) {
    val, err := m.GetBackingStore().Get("permissionGrants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ResourceSpecificPermissionGrantable)
    }
    return nil
}
// GetPinnedMessages gets the pinnedMessages property value. A collection of all the pinned messages in the chat. Nullable.
// returns a []PinnedChatMessageInfoable when successful
func (m *Chat) GetPinnedMessages()([]PinnedChatMessageInfoable) {
    val, err := m.GetBackingStore().Get("pinnedMessages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PinnedChatMessageInfoable)
    }
    return nil
}
// GetTabs gets the tabs property value. A collection of all the tabs in the chat. Nullable.
// returns a []TeamsTabable when successful
func (m *Chat) GetTabs()([]TeamsTabable) {
    val, err := m.GetBackingStore().Get("tabs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsTabable)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The identifier of the tenant in which the chat was created. Read-only.
// returns a *string when successful
func (m *Chat) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTopic gets the topic property value. (Optional) Subject or topic for the chat. Only available for group chats.
// returns a *string when successful
func (m *Chat) GetTopic()(*string) {
    val, err := m.GetBackingStore().Get("topic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetViewpoint gets the viewpoint property value. Represents caller-specific information about the chat, such as the last message read date and time. This property is populated only when the request is made in a delegated context.
// returns a ChatViewpointable when successful
func (m *Chat) GetViewpoint()(ChatViewpointable) {
    val, err := m.GetBackingStore().Get("viewpoint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatViewpointable)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. The URL for the chat in Microsoft Teams. The URL should be treated as an opaque blob, and not parsed. Read-only.
// returns a *string when successful
func (m *Chat) GetWebUrl()(*string) {
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
func (m *Chat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChatType() != nil {
        cast := (*m.GetChatType()).String()
        err = writer.WriteStringValue("chatType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetInstalledApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInstalledApps()))
        for i, v := range m.GetInstalledApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("installedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastMessagePreview", m.GetLastMessagePreview())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdatedDateTime", m.GetLastUpdatedDateTime())
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
    {
        err = writer.WriteObjectValue("onlineMeetingInfo", m.GetOnlineMeetingInfo())
        if err != nil {
            return err
        }
    }
    if m.GetPermissionGrants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPermissionGrants()))
        for i, v := range m.GetPermissionGrants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("permissionGrants", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPinnedMessages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPinnedMessages()))
        for i, v := range m.GetPinnedMessages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("pinnedMessages", cast)
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
        err = writer.WriteStringValue("topic", m.GetTopic())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("viewpoint", m.GetViewpoint())
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
// SetChatType sets the chatType property value. The chatType property
func (m *Chat) SetChatType(value *ChatType)() {
    err := m.GetBackingStore().Set("chatType", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time at which the chat was created. Read-only.
func (m *Chat) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledApps sets the installedApps property value. A collection of all the apps in the chat. Nullable.
func (m *Chat) SetInstalledApps(value []TeamsAppInstallationable)() {
    err := m.GetBackingStore().Set("installedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetLastMessagePreview sets the lastMessagePreview property value. Preview of the last message sent in the chat. Null if no messages were sent in the chat. Currently, only the list chats operation supports this property.
func (m *Chat) SetLastMessagePreview(value ChatMessageInfoable)() {
    err := m.GetBackingStore().Set("lastMessagePreview", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdatedDateTime sets the lastUpdatedDateTime property value. Date and time at which the chat was renamed or the list of members was last changed. Read-only.
func (m *Chat) SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdatedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. A collection of all the members in the chat. Nullable.
func (m *Chat) SetMembers(value []ConversationMemberable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetMessages sets the messages property value. A collection of all the messages in the chat. Nullable.
func (m *Chat) SetMessages(value []ChatMessageable)() {
    err := m.GetBackingStore().Set("messages", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeetingInfo sets the onlineMeetingInfo property value. Represents details about an online meeting. If the chat isn't associated with an online meeting, the property is empty. Read-only.
func (m *Chat) SetOnlineMeetingInfo(value TeamworkOnlineMeetingInfoable)() {
    err := m.GetBackingStore().Set("onlineMeetingInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionGrants sets the permissionGrants property value. A collection of permissions granted to apps for the chat.
func (m *Chat) SetPermissionGrants(value []ResourceSpecificPermissionGrantable)() {
    err := m.GetBackingStore().Set("permissionGrants", value)
    if err != nil {
        panic(err)
    }
}
// SetPinnedMessages sets the pinnedMessages property value. A collection of all the pinned messages in the chat. Nullable.
func (m *Chat) SetPinnedMessages(value []PinnedChatMessageInfoable)() {
    err := m.GetBackingStore().Set("pinnedMessages", value)
    if err != nil {
        panic(err)
    }
}
// SetTabs sets the tabs property value. A collection of all the tabs in the chat. Nullable.
func (m *Chat) SetTabs(value []TeamsTabable)() {
    err := m.GetBackingStore().Set("tabs", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The identifier of the tenant in which the chat was created. Read-only.
func (m *Chat) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetTopic sets the topic property value. (Optional) Subject or topic for the chat. Only available for group chats.
func (m *Chat) SetTopic(value *string)() {
    err := m.GetBackingStore().Set("topic", value)
    if err != nil {
        panic(err)
    }
}
// SetViewpoint sets the viewpoint property value. Represents caller-specific information about the chat, such as the last message read date and time. This property is populated only when the request is made in a delegated context.
func (m *Chat) SetViewpoint(value ChatViewpointable)() {
    err := m.GetBackingStore().Set("viewpoint", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. The URL for the chat in Microsoft Teams. The URL should be treated as an opaque blob, and not parsed. Read-only.
func (m *Chat) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type Chatable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChatType()(*ChatType)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetInstalledApps()([]TeamsAppInstallationable)
    GetLastMessagePreview()(ChatMessageInfoable)
    GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMembers()([]ConversationMemberable)
    GetMessages()([]ChatMessageable)
    GetOnlineMeetingInfo()(TeamworkOnlineMeetingInfoable)
    GetPermissionGrants()([]ResourceSpecificPermissionGrantable)
    GetPinnedMessages()([]PinnedChatMessageInfoable)
    GetTabs()([]TeamsTabable)
    GetTenantId()(*string)
    GetTopic()(*string)
    GetViewpoint()(ChatViewpointable)
    GetWebUrl()(*string)
    SetChatType(value *ChatType)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetInstalledApps(value []TeamsAppInstallationable)()
    SetLastMessagePreview(value ChatMessageInfoable)()
    SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMembers(value []ConversationMemberable)()
    SetMessages(value []ChatMessageable)()
    SetOnlineMeetingInfo(value TeamworkOnlineMeetingInfoable)()
    SetPermissionGrants(value []ResourceSpecificPermissionGrantable)()
    SetPinnedMessages(value []PinnedChatMessageInfoable)()
    SetTabs(value []TeamsTabable)()
    SetTenantId(value *string)()
    SetTopic(value *string)()
    SetViewpoint(value ChatViewpointable)()
    SetWebUrl(value *string)()
}

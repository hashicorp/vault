package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Team struct {
    Entity
}
// NewTeam instantiates a new Team and sets the default values.
func NewTeam()(*Team) {
    m := &Team{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeam(), nil
}
// GetAllChannels gets the allChannels property value. List of channels either hosted in or shared with the team (incoming channels).
// returns a []Channelable when successful
func (m *Team) GetAllChannels()([]Channelable) {
    val, err := m.GetBackingStore().Get("allChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Channelable)
    }
    return nil
}
// GetChannels gets the channels property value. The collection of channels and messages associated with the team.
// returns a []Channelable when successful
func (m *Team) GetChannels()([]Channelable) {
    val, err := m.GetBackingStore().Get("channels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Channelable)
    }
    return nil
}
// GetClassification gets the classification property value. An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
// returns a *string when successful
func (m *Team) GetClassification()(*string) {
    val, err := m.GetBackingStore().Get("classification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp at which the team was created.
// returns a *Time when successful
func (m *Team) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. An optional description for the team. Maximum length: 1024 characters.
// returns a *string when successful
func (m *Team) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the team.
// returns a *string when successful
func (m *Team) GetDisplayName()(*string) {
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
func (m *Team) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Channelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Channelable)
                }
            }
            m.SetAllChannels(res)
        }
        return nil
    }
    res["channels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Channelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Channelable)
                }
            }
            m.SetChannels(res)
        }
        return nil
    }
    res["classification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassification(val)
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
    res["funSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamFunSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFunSettings(val.(TeamFunSettingsable))
        }
        return nil
    }
    res["group"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroup(val.(Groupable))
        }
        return nil
    }
    res["guestSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamGuestSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGuestSettings(val.(TeamGuestSettingsable))
        }
        return nil
    }
    res["incomingChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Channelable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Channelable)
                }
            }
            m.SetIncomingChannels(res)
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
    res["internalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternalId(val)
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
    res["memberSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamMemberSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMemberSettings(val.(TeamMemberSettingsable))
        }
        return nil
    }
    res["messagingSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamMessagingSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessagingSettings(val.(TeamMessagingSettingsable))
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamsAsyncOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamsAsyncOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamsAsyncOperationable)
                }
            }
            m.SetOperations(res)
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
    res["photo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateProfilePhotoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoto(val.(ProfilePhotoable))
        }
        return nil
    }
    res["primaryChannel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChannelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrimaryChannel(val.(Channelable))
        }
        return nil
    }
    res["schedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedule(val.(Scheduleable))
        }
        return nil
    }
    res["specialization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamSpecialization)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpecialization(val.(*TeamSpecialization))
        }
        return nil
    }
    res["summary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSummary(val.(TeamSummaryable))
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamworkTagFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamworkTagable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamworkTagable)
                }
            }
            m.SetTags(res)
        }
        return nil
    }
    res["template"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamsTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTemplate(val.(TeamsTemplateable))
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
    res["visibility"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamVisibilityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVisibility(val.(*TeamVisibilityType))
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
// GetFunSettings gets the funSettings property value. Settings to configure use of Giphy, memes, and stickers in the team.
// returns a TeamFunSettingsable when successful
func (m *Team) GetFunSettings()(TeamFunSettingsable) {
    val, err := m.GetBackingStore().Get("funSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamFunSettingsable)
    }
    return nil
}
// GetGroup gets the group property value. The group property
// returns a Groupable when successful
func (m *Team) GetGroup()(Groupable) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Groupable)
    }
    return nil
}
// GetGuestSettings gets the guestSettings property value. Settings to configure whether guests can create, update, or delete channels in the team.
// returns a TeamGuestSettingsable when successful
func (m *Team) GetGuestSettings()(TeamGuestSettingsable) {
    val, err := m.GetBackingStore().Get("guestSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamGuestSettingsable)
    }
    return nil
}
// GetIncomingChannels gets the incomingChannels property value. List of channels shared with the team.
// returns a []Channelable when successful
func (m *Team) GetIncomingChannels()([]Channelable) {
    val, err := m.GetBackingStore().Get("incomingChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Channelable)
    }
    return nil
}
// GetInstalledApps gets the installedApps property value. The apps installed in this team.
// returns a []TeamsAppInstallationable when successful
func (m *Team) GetInstalledApps()([]TeamsAppInstallationable) {
    val, err := m.GetBackingStore().Get("installedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsAppInstallationable)
    }
    return nil
}
// GetInternalId gets the internalId property value. A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
// returns a *string when successful
func (m *Team) GetInternalId()(*string) {
    val, err := m.GetBackingStore().Get("internalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsArchived gets the isArchived property value. Whether this team is in read-only mode.
// returns a *bool when successful
func (m *Team) GetIsArchived()(*bool) {
    val, err := m.GetBackingStore().Get("isArchived")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMembers gets the members property value. Members and owners of the team.
// returns a []ConversationMemberable when successful
func (m *Team) GetMembers()([]ConversationMemberable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConversationMemberable)
    }
    return nil
}
// GetMemberSettings gets the memberSettings property value. Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.
// returns a TeamMemberSettingsable when successful
func (m *Team) GetMemberSettings()(TeamMemberSettingsable) {
    val, err := m.GetBackingStore().Get("memberSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamMemberSettingsable)
    }
    return nil
}
// GetMessagingSettings gets the messagingSettings property value. Settings to configure messaging and mentions in the team.
// returns a TeamMessagingSettingsable when successful
func (m *Team) GetMessagingSettings()(TeamMessagingSettingsable) {
    val, err := m.GetBackingStore().Get("messagingSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamMessagingSettingsable)
    }
    return nil
}
// GetOperations gets the operations property value. The async operations that ran or are running on this team.
// returns a []TeamsAsyncOperationable when successful
func (m *Team) GetOperations()([]TeamsAsyncOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsAsyncOperationable)
    }
    return nil
}
// GetPermissionGrants gets the permissionGrants property value. A collection of permissions granted to apps to access the team.
// returns a []ResourceSpecificPermissionGrantable when successful
func (m *Team) GetPermissionGrants()([]ResourceSpecificPermissionGrantable) {
    val, err := m.GetBackingStore().Get("permissionGrants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ResourceSpecificPermissionGrantable)
    }
    return nil
}
// GetPhoto gets the photo property value. The profile photo for the team.
// returns a ProfilePhotoable when successful
func (m *Team) GetPhoto()(ProfilePhotoable) {
    val, err := m.GetBackingStore().Get("photo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ProfilePhotoable)
    }
    return nil
}
// GetPrimaryChannel gets the primaryChannel property value. The general channel for the team.
// returns a Channelable when successful
func (m *Team) GetPrimaryChannel()(Channelable) {
    val, err := m.GetBackingStore().Get("primaryChannel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Channelable)
    }
    return nil
}
// GetSchedule gets the schedule property value. The schedule of shifts for this team.
// returns a Scheduleable when successful
func (m *Team) GetSchedule()(Scheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Scheduleable)
    }
    return nil
}
// GetSpecialization gets the specialization property value. Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
// returns a *TeamSpecialization when successful
func (m *Team) GetSpecialization()(*TeamSpecialization) {
    val, err := m.GetBackingStore().Get("specialization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamSpecialization)
    }
    return nil
}
// GetSummary gets the summary property value. Contains summary information about the team, including number of owners, members, and guests.
// returns a TeamSummaryable when successful
func (m *Team) GetSummary()(TeamSummaryable) {
    val, err := m.GetBackingStore().Get("summary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamSummaryable)
    }
    return nil
}
// GetTags gets the tags property value. The tags associated with the team.
// returns a []TeamworkTagable when successful
func (m *Team) GetTags()([]TeamworkTagable) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamworkTagable)
    }
    return nil
}
// GetTemplate gets the template property value. The template this team was created from. See available templates.
// returns a TeamsTemplateable when successful
func (m *Team) GetTemplate()(TeamsTemplateable) {
    val, err := m.GetBackingStore().Get("template")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamsTemplateable)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The ID of the Microsoft Entra tenant.
// returns a *string when successful
func (m *Team) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVisibility gets the visibility property value. The visibility of the group and team. Defaults to Public.
// returns a *TeamVisibilityType when successful
func (m *Team) GetVisibility()(*TeamVisibilityType) {
    val, err := m.GetBackingStore().Get("visibility")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamVisibilityType)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.
// returns a *string when successful
func (m *Team) GetWebUrl()(*string) {
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
func (m *Team) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllChannels() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllChannels()))
        for i, v := range m.GetAllChannels() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("allChannels", cast)
        if err != nil {
            return err
        }
    }
    if m.GetChannels() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChannels()))
        for i, v := range m.GetChannels() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("channels", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("classification", m.GetClassification())
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
        err = writer.WriteObjectValue("funSettings", m.GetFunSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("guestSettings", m.GetGuestSettings())
        if err != nil {
            return err
        }
    }
    if m.GetIncomingChannels() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIncomingChannels()))
        for i, v := range m.GetIncomingChannels() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("incomingChannels", cast)
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
        err = writer.WriteStringValue("internalId", m.GetInternalId())
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
    {
        err = writer.WriteObjectValue("memberSettings", m.GetMemberSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("messagingSettings", m.GetMessagingSettings())
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
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
    {
        err = writer.WriteObjectValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("primaryChannel", m.GetPrimaryChannel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("schedule", m.GetSchedule())
        if err != nil {
            return err
        }
    }
    if m.GetSpecialization() != nil {
        cast := (*m.GetSpecialization()).String()
        err = writer.WriteStringValue("specialization", &cast)
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
    if m.GetTags() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTags()))
        for i, v := range m.GetTags() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tags", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("template", m.GetTemplate())
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
    if m.GetVisibility() != nil {
        cast := (*m.GetVisibility()).String()
        err = writer.WriteStringValue("visibility", &cast)
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
// SetAllChannels sets the allChannels property value. List of channels either hosted in or shared with the team (incoming channels).
func (m *Team) SetAllChannels(value []Channelable)() {
    err := m.GetBackingStore().Set("allChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetChannels sets the channels property value. The collection of channels and messages associated with the team.
func (m *Team) SetChannels(value []Channelable)() {
    err := m.GetBackingStore().Set("channels", value)
    if err != nil {
        panic(err)
    }
}
// SetClassification sets the classification property value. An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
func (m *Team) SetClassification(value *string)() {
    err := m.GetBackingStore().Set("classification", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp at which the team was created.
func (m *Team) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. An optional description for the team. Maximum length: 1024 characters.
func (m *Team) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the team.
func (m *Team) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetFunSettings sets the funSettings property value. Settings to configure use of Giphy, memes, and stickers in the team.
func (m *Team) SetFunSettings(value TeamFunSettingsable)() {
    err := m.GetBackingStore().Set("funSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetGroup sets the group property value. The group property
func (m *Team) SetGroup(value Groupable)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
// SetGuestSettings sets the guestSettings property value. Settings to configure whether guests can create, update, or delete channels in the team.
func (m *Team) SetGuestSettings(value TeamGuestSettingsable)() {
    err := m.GetBackingStore().Set("guestSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetIncomingChannels sets the incomingChannels property value. List of channels shared with the team.
func (m *Team) SetIncomingChannels(value []Channelable)() {
    err := m.GetBackingStore().Set("incomingChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledApps sets the installedApps property value. The apps installed in this team.
func (m *Team) SetInstalledApps(value []TeamsAppInstallationable)() {
    err := m.GetBackingStore().Set("installedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetInternalId sets the internalId property value. A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
func (m *Team) SetInternalId(value *string)() {
    err := m.GetBackingStore().Set("internalId", value)
    if err != nil {
        panic(err)
    }
}
// SetIsArchived sets the isArchived property value. Whether this team is in read-only mode.
func (m *Team) SetIsArchived(value *bool)() {
    err := m.GetBackingStore().Set("isArchived", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. Members and owners of the team.
func (m *Team) SetMembers(value []ConversationMemberable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberSettings sets the memberSettings property value. Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.
func (m *Team) SetMemberSettings(value TeamMemberSettingsable)() {
    err := m.GetBackingStore().Set("memberSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetMessagingSettings sets the messagingSettings property value. Settings to configure messaging and mentions in the team.
func (m *Team) SetMessagingSettings(value TeamMessagingSettingsable)() {
    err := m.GetBackingStore().Set("messagingSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The async operations that ran or are running on this team.
func (m *Team) SetOperations(value []TeamsAsyncOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionGrants sets the permissionGrants property value. A collection of permissions granted to apps to access the team.
func (m *Team) SetPermissionGrants(value []ResourceSpecificPermissionGrantable)() {
    err := m.GetBackingStore().Set("permissionGrants", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoto sets the photo property value. The profile photo for the team.
func (m *Team) SetPhoto(value ProfilePhotoable)() {
    err := m.GetBackingStore().Set("photo", value)
    if err != nil {
        panic(err)
    }
}
// SetPrimaryChannel sets the primaryChannel property value. The general channel for the team.
func (m *Team) SetPrimaryChannel(value Channelable)() {
    err := m.GetBackingStore().Set("primaryChannel", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. The schedule of shifts for this team.
func (m *Team) SetSchedule(value Scheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
// SetSpecialization sets the specialization property value. Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
func (m *Team) SetSpecialization(value *TeamSpecialization)() {
    err := m.GetBackingStore().Set("specialization", value)
    if err != nil {
        panic(err)
    }
}
// SetSummary sets the summary property value. Contains summary information about the team, including number of owners, members, and guests.
func (m *Team) SetSummary(value TeamSummaryable)() {
    err := m.GetBackingStore().Set("summary", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. The tags associated with the team.
func (m *Team) SetTags(value []TeamworkTagable)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetTemplate sets the template property value. The template this team was created from. See available templates.
func (m *Team) SetTemplate(value TeamsTemplateable)() {
    err := m.GetBackingStore().Set("template", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The ID of the Microsoft Entra tenant.
func (m *Team) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetVisibility sets the visibility property value. The visibility of the group and team. Defaults to Public.
func (m *Team) SetVisibility(value *TeamVisibilityType)() {
    err := m.GetBackingStore().Set("visibility", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.
func (m *Team) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type Teamable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllChannels()([]Channelable)
    GetChannels()([]Channelable)
    GetClassification()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetFunSettings()(TeamFunSettingsable)
    GetGroup()(Groupable)
    GetGuestSettings()(TeamGuestSettingsable)
    GetIncomingChannels()([]Channelable)
    GetInstalledApps()([]TeamsAppInstallationable)
    GetInternalId()(*string)
    GetIsArchived()(*bool)
    GetMembers()([]ConversationMemberable)
    GetMemberSettings()(TeamMemberSettingsable)
    GetMessagingSettings()(TeamMessagingSettingsable)
    GetOperations()([]TeamsAsyncOperationable)
    GetPermissionGrants()([]ResourceSpecificPermissionGrantable)
    GetPhoto()(ProfilePhotoable)
    GetPrimaryChannel()(Channelable)
    GetSchedule()(Scheduleable)
    GetSpecialization()(*TeamSpecialization)
    GetSummary()(TeamSummaryable)
    GetTags()([]TeamworkTagable)
    GetTemplate()(TeamsTemplateable)
    GetTenantId()(*string)
    GetVisibility()(*TeamVisibilityType)
    GetWebUrl()(*string)
    SetAllChannels(value []Channelable)()
    SetChannels(value []Channelable)()
    SetClassification(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetFunSettings(value TeamFunSettingsable)()
    SetGroup(value Groupable)()
    SetGuestSettings(value TeamGuestSettingsable)()
    SetIncomingChannels(value []Channelable)()
    SetInstalledApps(value []TeamsAppInstallationable)()
    SetInternalId(value *string)()
    SetIsArchived(value *bool)()
    SetMembers(value []ConversationMemberable)()
    SetMemberSettings(value TeamMemberSettingsable)()
    SetMessagingSettings(value TeamMessagingSettingsable)()
    SetOperations(value []TeamsAsyncOperationable)()
    SetPermissionGrants(value []ResourceSpecificPermissionGrantable)()
    SetPhoto(value ProfilePhotoable)()
    SetPrimaryChannel(value Channelable)()
    SetSchedule(value Scheduleable)()
    SetSpecialization(value *TeamSpecialization)()
    SetSummary(value TeamSummaryable)()
    SetTags(value []TeamworkTagable)()
    SetTemplate(value TeamsTemplateable)()
    SetTenantId(value *string)()
    SetVisibility(value *TeamVisibilityType)()
    SetWebUrl(value *string)()
}

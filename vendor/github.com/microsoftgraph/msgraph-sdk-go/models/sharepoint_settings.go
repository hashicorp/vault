package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharepointSettings struct {
    Entity
}
// NewSharepointSettings instantiates a new SharepointSettings and sets the default values.
func NewSharepointSettings()(*SharepointSettings) {
    m := &SharepointSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSharepointSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharepointSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharepointSettings(), nil
}
// GetAllowedDomainGuidsForSyncApp gets the allowedDomainGuidsForSyncApp property value. Collection of trusted domain GUIDs for the OneDrive sync app.
// returns a []UUID when successful
func (m *SharepointSettings) GetAllowedDomainGuidsForSyncApp()([]i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("allowedDomainGuidsForSyncApp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetAvailableManagedPathsForSiteCreation gets the availableManagedPathsForSiteCreation property value. Collection of managed paths available for site creation. Read-only.
// returns a []string when successful
func (m *SharepointSettings) GetAvailableManagedPathsForSiteCreation()([]string) {
    val, err := m.GetBackingStore().Get("availableManagedPathsForSiteCreation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDeletedUserPersonalSiteRetentionPeriodInDays gets the deletedUserPersonalSiteRetentionPeriodInDays property value. The number of days for preserving a deleted user's OneDrive.
// returns a *int32 when successful
func (m *SharepointSettings) GetDeletedUserPersonalSiteRetentionPeriodInDays()(*int32) {
    val, err := m.GetBackingStore().Get("deletedUserPersonalSiteRetentionPeriodInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetExcludedFileExtensionsForSyncApp gets the excludedFileExtensionsForSyncApp property value. Collection of file extensions not uploaded by the OneDrive sync app.
// returns a []string when successful
func (m *SharepointSettings) GetExcludedFileExtensionsForSyncApp()([]string) {
    val, err := m.GetBackingStore().Get("excludedFileExtensionsForSyncApp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharepointSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowedDomainGuidsForSyncApp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("uuid")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID))
                }
            }
            m.SetAllowedDomainGuidsForSyncApp(res)
        }
        return nil
    }
    res["availableManagedPathsForSiteCreation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAvailableManagedPathsForSiteCreation(res)
        }
        return nil
    }
    res["deletedUserPersonalSiteRetentionPeriodInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedUserPersonalSiteRetentionPeriodInDays(val)
        }
        return nil
    }
    res["excludedFileExtensionsForSyncApp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetExcludedFileExtensionsForSyncApp(res)
        }
        return nil
    }
    res["idleSessionSignOut"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdleSessionSignOutFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdleSessionSignOut(val.(IdleSessionSignOutable))
        }
        return nil
    }
    res["imageTaggingOption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseImageTaggingChoice)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImageTaggingOption(val.(*ImageTaggingChoice))
        }
        return nil
    }
    res["isCommentingOnSitePagesEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCommentingOnSitePagesEnabled(val)
        }
        return nil
    }
    res["isFileActivityNotificationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFileActivityNotificationEnabled(val)
        }
        return nil
    }
    res["isLegacyAuthProtocolsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLegacyAuthProtocolsEnabled(val)
        }
        return nil
    }
    res["isLoopEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLoopEnabled(val)
        }
        return nil
    }
    res["isMacSyncAppEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMacSyncAppEnabled(val)
        }
        return nil
    }
    res["isRequireAcceptingUserToMatchInvitedUserEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRequireAcceptingUserToMatchInvitedUserEnabled(val)
        }
        return nil
    }
    res["isResharingByExternalUsersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsResharingByExternalUsersEnabled(val)
        }
        return nil
    }
    res["isSharePointMobileNotificationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSharePointMobileNotificationEnabled(val)
        }
        return nil
    }
    res["isSharePointNewsfeedEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSharePointNewsfeedEnabled(val)
        }
        return nil
    }
    res["isSiteCreationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSiteCreationEnabled(val)
        }
        return nil
    }
    res["isSiteCreationUIEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSiteCreationUIEnabled(val)
        }
        return nil
    }
    res["isSitePagesCreationEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSitePagesCreationEnabled(val)
        }
        return nil
    }
    res["isSitesStorageLimitAutomatic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSitesStorageLimitAutomatic(val)
        }
        return nil
    }
    res["isSyncButtonHiddenOnPersonalSite"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSyncButtonHiddenOnPersonalSite(val)
        }
        return nil
    }
    res["isUnmanagedSyncAppForTenantRestricted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsUnmanagedSyncAppForTenantRestricted(val)
        }
        return nil
    }
    res["personalSiteDefaultStorageLimitInMB"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonalSiteDefaultStorageLimitInMB(val)
        }
        return nil
    }
    res["sharingAllowedDomainList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSharingAllowedDomainList(res)
        }
        return nil
    }
    res["sharingBlockedDomainList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSharingBlockedDomainList(res)
        }
        return nil
    }
    res["sharingCapability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSharingCapabilities)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharingCapability(val.(*SharingCapabilities))
        }
        return nil
    }
    res["sharingDomainRestrictionMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSharingDomainRestrictionMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharingDomainRestrictionMode(val.(*SharingDomainRestrictionMode))
        }
        return nil
    }
    res["siteCreationDefaultManagedPath"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteCreationDefaultManagedPath(val)
        }
        return nil
    }
    res["siteCreationDefaultStorageLimitInMB"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteCreationDefaultStorageLimitInMB(val)
        }
        return nil
    }
    res["tenantDefaultTimezone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantDefaultTimezone(val)
        }
        return nil
    }
    return res
}
// GetIdleSessionSignOut gets the idleSessionSignOut property value. Specifies the idle session sign-out policies for the tenant.
// returns a IdleSessionSignOutable when successful
func (m *SharepointSettings) GetIdleSessionSignOut()(IdleSessionSignOutable) {
    val, err := m.GetBackingStore().Get("idleSessionSignOut")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdleSessionSignOutable)
    }
    return nil
}
// GetImageTaggingOption gets the imageTaggingOption property value. Specifies the image tagging option for the tenant. Possible values are: disabled, basic, enhanced.
// returns a *ImageTaggingChoice when successful
func (m *SharepointSettings) GetImageTaggingOption()(*ImageTaggingChoice) {
    val, err := m.GetBackingStore().Get("imageTaggingOption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ImageTaggingChoice)
    }
    return nil
}
// GetIsCommentingOnSitePagesEnabled gets the isCommentingOnSitePagesEnabled property value. Indicates whether comments are allowed on modern site pages in SharePoint.
// returns a *bool when successful
func (m *SharepointSettings) GetIsCommentingOnSitePagesEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isCommentingOnSitePagesEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsFileActivityNotificationEnabled gets the isFileActivityNotificationEnabled property value. Indicates whether push notifications are enabled for OneDrive events.
// returns a *bool when successful
func (m *SharepointSettings) GetIsFileActivityNotificationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isFileActivityNotificationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsLegacyAuthProtocolsEnabled gets the isLegacyAuthProtocolsEnabled property value. Indicates whether legacy authentication protocols are enabled for the tenant.
// returns a *bool when successful
func (m *SharepointSettings) GetIsLegacyAuthProtocolsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isLegacyAuthProtocolsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsLoopEnabled gets the isLoopEnabled property value. Indicates whether if Fluid Framework is allowed on SharePoint sites.
// returns a *bool when successful
func (m *SharepointSettings) GetIsLoopEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isLoopEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMacSyncAppEnabled gets the isMacSyncAppEnabled property value. Indicates whether files can be synced using the OneDrive sync app for Mac.
// returns a *bool when successful
func (m *SharepointSettings) GetIsMacSyncAppEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isMacSyncAppEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRequireAcceptingUserToMatchInvitedUserEnabled gets the isRequireAcceptingUserToMatchInvitedUserEnabled property value. Indicates whether guests must sign in using the same account to which sharing invitations are sent.
// returns a *bool when successful
func (m *SharepointSettings) GetIsRequireAcceptingUserToMatchInvitedUserEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isRequireAcceptingUserToMatchInvitedUserEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsResharingByExternalUsersEnabled gets the isResharingByExternalUsersEnabled property value. Indicates whether guests are allowed to reshare files, folders, and sites they don't own.
// returns a *bool when successful
func (m *SharepointSettings) GetIsResharingByExternalUsersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isResharingByExternalUsersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSharePointMobileNotificationEnabled gets the isSharePointMobileNotificationEnabled property value. Indicates whether mobile push notifications are enabled for SharePoint.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSharePointMobileNotificationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isSharePointMobileNotificationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSharePointNewsfeedEnabled gets the isSharePointNewsfeedEnabled property value. Indicates whether the newsfeed is allowed on the modern site pages in SharePoint.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSharePointNewsfeedEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isSharePointNewsfeedEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSiteCreationEnabled gets the isSiteCreationEnabled property value. Indicates whether users are allowed to create sites.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSiteCreationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isSiteCreationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSiteCreationUIEnabled gets the isSiteCreationUIEnabled property value. Indicates whether the UI commands for creating sites are shown.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSiteCreationUIEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isSiteCreationUIEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSitePagesCreationEnabled gets the isSitePagesCreationEnabled property value. Indicates whether creating new modern pages is allowed on SharePoint sites.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSitePagesCreationEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isSitePagesCreationEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSitesStorageLimitAutomatic gets the isSitesStorageLimitAutomatic property value. Indicates whether site storage space is automatically managed or if specific storage limits are set per site.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSitesStorageLimitAutomatic()(*bool) {
    val, err := m.GetBackingStore().Get("isSitesStorageLimitAutomatic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSyncButtonHiddenOnPersonalSite gets the isSyncButtonHiddenOnPersonalSite property value. Indicates whether the sync button in OneDrive is hidden.
// returns a *bool when successful
func (m *SharepointSettings) GetIsSyncButtonHiddenOnPersonalSite()(*bool) {
    val, err := m.GetBackingStore().Get("isSyncButtonHiddenOnPersonalSite")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsUnmanagedSyncAppForTenantRestricted gets the isUnmanagedSyncAppForTenantRestricted property value. Indicates whether users are allowed to sync files only on PCs joined to specific domains.
// returns a *bool when successful
func (m *SharepointSettings) GetIsUnmanagedSyncAppForTenantRestricted()(*bool) {
    val, err := m.GetBackingStore().Get("isUnmanagedSyncAppForTenantRestricted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPersonalSiteDefaultStorageLimitInMB gets the personalSiteDefaultStorageLimitInMB property value. The default OneDrive storage limit for all new and existing users who are assigned a qualifying license. Measured in megabytes (MB).
// returns a *int64 when successful
func (m *SharepointSettings) GetPersonalSiteDefaultStorageLimitInMB()(*int64) {
    val, err := m.GetBackingStore().Get("personalSiteDefaultStorageLimitInMB")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetSharingAllowedDomainList gets the sharingAllowedDomainList property value. Collection of email domains that are allowed for sharing outside the organization.
// returns a []string when successful
func (m *SharepointSettings) GetSharingAllowedDomainList()([]string) {
    val, err := m.GetBackingStore().Get("sharingAllowedDomainList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSharingBlockedDomainList gets the sharingBlockedDomainList property value. Collection of email domains that are blocked for sharing outside the organization.
// returns a []string when successful
func (m *SharepointSettings) GetSharingBlockedDomainList()([]string) {
    val, err := m.GetBackingStore().Get("sharingBlockedDomainList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSharingCapability gets the sharingCapability property value. Sharing capability for the tenant. Possible values are: disabled, externalUserSharingOnly, externalUserAndGuestSharing, existingExternalUserSharingOnly.
// returns a *SharingCapabilities when successful
func (m *SharepointSettings) GetSharingCapability()(*SharingCapabilities) {
    val, err := m.GetBackingStore().Get("sharingCapability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SharingCapabilities)
    }
    return nil
}
// GetSharingDomainRestrictionMode gets the sharingDomainRestrictionMode property value. Specifies the external sharing mode for domains. Possible values are: none, allowList, blockList.
// returns a *SharingDomainRestrictionMode when successful
func (m *SharepointSettings) GetSharingDomainRestrictionMode()(*SharingDomainRestrictionMode) {
    val, err := m.GetBackingStore().Get("sharingDomainRestrictionMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SharingDomainRestrictionMode)
    }
    return nil
}
// GetSiteCreationDefaultManagedPath gets the siteCreationDefaultManagedPath property value. The value of the team site managed path. This is the path under which new team sites will be created.
// returns a *string when successful
func (m *SharepointSettings) GetSiteCreationDefaultManagedPath()(*string) {
    val, err := m.GetBackingStore().Get("siteCreationDefaultManagedPath")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSiteCreationDefaultStorageLimitInMB gets the siteCreationDefaultStorageLimitInMB property value. The default storage quota for a new site upon creation. Measured in megabytes (MB).
// returns a *int32 when successful
func (m *SharepointSettings) GetSiteCreationDefaultStorageLimitInMB()(*int32) {
    val, err := m.GetBackingStore().Get("siteCreationDefaultStorageLimitInMB")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTenantDefaultTimezone gets the tenantDefaultTimezone property value. The default timezone of a tenant for newly created sites. For a list of possible values, see SPRegionalSettings.TimeZones property.
// returns a *string when successful
func (m *SharepointSettings) GetTenantDefaultTimezone()(*string) {
    val, err := m.GetBackingStore().Get("tenantDefaultTimezone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharepointSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedDomainGuidsForSyncApp() != nil {
        err = writer.WriteCollectionOfUUIDValues("allowedDomainGuidsForSyncApp", m.GetAllowedDomainGuidsForSyncApp())
        if err != nil {
            return err
        }
    }
    if m.GetAvailableManagedPathsForSiteCreation() != nil {
        err = writer.WriteCollectionOfStringValues("availableManagedPathsForSiteCreation", m.GetAvailableManagedPathsForSiteCreation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deletedUserPersonalSiteRetentionPeriodInDays", m.GetDeletedUserPersonalSiteRetentionPeriodInDays())
        if err != nil {
            return err
        }
    }
    if m.GetExcludedFileExtensionsForSyncApp() != nil {
        err = writer.WriteCollectionOfStringValues("excludedFileExtensionsForSyncApp", m.GetExcludedFileExtensionsForSyncApp())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("idleSessionSignOut", m.GetIdleSessionSignOut())
        if err != nil {
            return err
        }
    }
    if m.GetImageTaggingOption() != nil {
        cast := (*m.GetImageTaggingOption()).String()
        err = writer.WriteStringValue("imageTaggingOption", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isCommentingOnSitePagesEnabled", m.GetIsCommentingOnSitePagesEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFileActivityNotificationEnabled", m.GetIsFileActivityNotificationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isLegacyAuthProtocolsEnabled", m.GetIsLegacyAuthProtocolsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isLoopEnabled", m.GetIsLoopEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isMacSyncAppEnabled", m.GetIsMacSyncAppEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isRequireAcceptingUserToMatchInvitedUserEnabled", m.GetIsRequireAcceptingUserToMatchInvitedUserEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isResharingByExternalUsersEnabled", m.GetIsResharingByExternalUsersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSharePointMobileNotificationEnabled", m.GetIsSharePointMobileNotificationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSharePointNewsfeedEnabled", m.GetIsSharePointNewsfeedEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSiteCreationEnabled", m.GetIsSiteCreationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSiteCreationUIEnabled", m.GetIsSiteCreationUIEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSitePagesCreationEnabled", m.GetIsSitePagesCreationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSitesStorageLimitAutomatic", m.GetIsSitesStorageLimitAutomatic())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSyncButtonHiddenOnPersonalSite", m.GetIsSyncButtonHiddenOnPersonalSite())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isUnmanagedSyncAppForTenantRestricted", m.GetIsUnmanagedSyncAppForTenantRestricted())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("personalSiteDefaultStorageLimitInMB", m.GetPersonalSiteDefaultStorageLimitInMB())
        if err != nil {
            return err
        }
    }
    if m.GetSharingAllowedDomainList() != nil {
        err = writer.WriteCollectionOfStringValues("sharingAllowedDomainList", m.GetSharingAllowedDomainList())
        if err != nil {
            return err
        }
    }
    if m.GetSharingBlockedDomainList() != nil {
        err = writer.WriteCollectionOfStringValues("sharingBlockedDomainList", m.GetSharingBlockedDomainList())
        if err != nil {
            return err
        }
    }
    if m.GetSharingCapability() != nil {
        cast := (*m.GetSharingCapability()).String()
        err = writer.WriteStringValue("sharingCapability", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSharingDomainRestrictionMode() != nil {
        cast := (*m.GetSharingDomainRestrictionMode()).String()
        err = writer.WriteStringValue("sharingDomainRestrictionMode", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("siteCreationDefaultManagedPath", m.GetSiteCreationDefaultManagedPath())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("siteCreationDefaultStorageLimitInMB", m.GetSiteCreationDefaultStorageLimitInMB())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantDefaultTimezone", m.GetTenantDefaultTimezone())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedDomainGuidsForSyncApp sets the allowedDomainGuidsForSyncApp property value. Collection of trusted domain GUIDs for the OneDrive sync app.
func (m *SharepointSettings) SetAllowedDomainGuidsForSyncApp(value []i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("allowedDomainGuidsForSyncApp", value)
    if err != nil {
        panic(err)
    }
}
// SetAvailableManagedPathsForSiteCreation sets the availableManagedPathsForSiteCreation property value. Collection of managed paths available for site creation. Read-only.
func (m *SharepointSettings) SetAvailableManagedPathsForSiteCreation(value []string)() {
    err := m.GetBackingStore().Set("availableManagedPathsForSiteCreation", value)
    if err != nil {
        panic(err)
    }
}
// SetDeletedUserPersonalSiteRetentionPeriodInDays sets the deletedUserPersonalSiteRetentionPeriodInDays property value. The number of days for preserving a deleted user's OneDrive.
func (m *SharepointSettings) SetDeletedUserPersonalSiteRetentionPeriodInDays(value *int32)() {
    err := m.GetBackingStore().Set("deletedUserPersonalSiteRetentionPeriodInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetExcludedFileExtensionsForSyncApp sets the excludedFileExtensionsForSyncApp property value. Collection of file extensions not uploaded by the OneDrive sync app.
func (m *SharepointSettings) SetExcludedFileExtensionsForSyncApp(value []string)() {
    err := m.GetBackingStore().Set("excludedFileExtensionsForSyncApp", value)
    if err != nil {
        panic(err)
    }
}
// SetIdleSessionSignOut sets the idleSessionSignOut property value. Specifies the idle session sign-out policies for the tenant.
func (m *SharepointSettings) SetIdleSessionSignOut(value IdleSessionSignOutable)() {
    err := m.GetBackingStore().Set("idleSessionSignOut", value)
    if err != nil {
        panic(err)
    }
}
// SetImageTaggingOption sets the imageTaggingOption property value. Specifies the image tagging option for the tenant. Possible values are: disabled, basic, enhanced.
func (m *SharepointSettings) SetImageTaggingOption(value *ImageTaggingChoice)() {
    err := m.GetBackingStore().Set("imageTaggingOption", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCommentingOnSitePagesEnabled sets the isCommentingOnSitePagesEnabled property value. Indicates whether comments are allowed on modern site pages in SharePoint.
func (m *SharepointSettings) SetIsCommentingOnSitePagesEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isCommentingOnSitePagesEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFileActivityNotificationEnabled sets the isFileActivityNotificationEnabled property value. Indicates whether push notifications are enabled for OneDrive events.
func (m *SharepointSettings) SetIsFileActivityNotificationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isFileActivityNotificationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsLegacyAuthProtocolsEnabled sets the isLegacyAuthProtocolsEnabled property value. Indicates whether legacy authentication protocols are enabled for the tenant.
func (m *SharepointSettings) SetIsLegacyAuthProtocolsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isLegacyAuthProtocolsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsLoopEnabled sets the isLoopEnabled property value. Indicates whether if Fluid Framework is allowed on SharePoint sites.
func (m *SharepointSettings) SetIsLoopEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isLoopEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMacSyncAppEnabled sets the isMacSyncAppEnabled property value. Indicates whether files can be synced using the OneDrive sync app for Mac.
func (m *SharepointSettings) SetIsMacSyncAppEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isMacSyncAppEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRequireAcceptingUserToMatchInvitedUserEnabled sets the isRequireAcceptingUserToMatchInvitedUserEnabled property value. Indicates whether guests must sign in using the same account to which sharing invitations are sent.
func (m *SharepointSettings) SetIsRequireAcceptingUserToMatchInvitedUserEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isRequireAcceptingUserToMatchInvitedUserEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsResharingByExternalUsersEnabled sets the isResharingByExternalUsersEnabled property value. Indicates whether guests are allowed to reshare files, folders, and sites they don't own.
func (m *SharepointSettings) SetIsResharingByExternalUsersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isResharingByExternalUsersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSharePointMobileNotificationEnabled sets the isSharePointMobileNotificationEnabled property value. Indicates whether mobile push notifications are enabled for SharePoint.
func (m *SharepointSettings) SetIsSharePointMobileNotificationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isSharePointMobileNotificationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSharePointNewsfeedEnabled sets the isSharePointNewsfeedEnabled property value. Indicates whether the newsfeed is allowed on the modern site pages in SharePoint.
func (m *SharepointSettings) SetIsSharePointNewsfeedEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isSharePointNewsfeedEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSiteCreationEnabled sets the isSiteCreationEnabled property value. Indicates whether users are allowed to create sites.
func (m *SharepointSettings) SetIsSiteCreationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isSiteCreationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSiteCreationUIEnabled sets the isSiteCreationUIEnabled property value. Indicates whether the UI commands for creating sites are shown.
func (m *SharepointSettings) SetIsSiteCreationUIEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isSiteCreationUIEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSitePagesCreationEnabled sets the isSitePagesCreationEnabled property value. Indicates whether creating new modern pages is allowed on SharePoint sites.
func (m *SharepointSettings) SetIsSitePagesCreationEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isSitePagesCreationEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSitesStorageLimitAutomatic sets the isSitesStorageLimitAutomatic property value. Indicates whether site storage space is automatically managed or if specific storage limits are set per site.
func (m *SharepointSettings) SetIsSitesStorageLimitAutomatic(value *bool)() {
    err := m.GetBackingStore().Set("isSitesStorageLimitAutomatic", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSyncButtonHiddenOnPersonalSite sets the isSyncButtonHiddenOnPersonalSite property value. Indicates whether the sync button in OneDrive is hidden.
func (m *SharepointSettings) SetIsSyncButtonHiddenOnPersonalSite(value *bool)() {
    err := m.GetBackingStore().Set("isSyncButtonHiddenOnPersonalSite", value)
    if err != nil {
        panic(err)
    }
}
// SetIsUnmanagedSyncAppForTenantRestricted sets the isUnmanagedSyncAppForTenantRestricted property value. Indicates whether users are allowed to sync files only on PCs joined to specific domains.
func (m *SharepointSettings) SetIsUnmanagedSyncAppForTenantRestricted(value *bool)() {
    err := m.GetBackingStore().Set("isUnmanagedSyncAppForTenantRestricted", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonalSiteDefaultStorageLimitInMB sets the personalSiteDefaultStorageLimitInMB property value. The default OneDrive storage limit for all new and existing users who are assigned a qualifying license. Measured in megabytes (MB).
func (m *SharepointSettings) SetPersonalSiteDefaultStorageLimitInMB(value *int64)() {
    err := m.GetBackingStore().Set("personalSiteDefaultStorageLimitInMB", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingAllowedDomainList sets the sharingAllowedDomainList property value. Collection of email domains that are allowed for sharing outside the organization.
func (m *SharepointSettings) SetSharingAllowedDomainList(value []string)() {
    err := m.GetBackingStore().Set("sharingAllowedDomainList", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingBlockedDomainList sets the sharingBlockedDomainList property value. Collection of email domains that are blocked for sharing outside the organization.
func (m *SharepointSettings) SetSharingBlockedDomainList(value []string)() {
    err := m.GetBackingStore().Set("sharingBlockedDomainList", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingCapability sets the sharingCapability property value. Sharing capability for the tenant. Possible values are: disabled, externalUserSharingOnly, externalUserAndGuestSharing, existingExternalUserSharingOnly.
func (m *SharepointSettings) SetSharingCapability(value *SharingCapabilities)() {
    err := m.GetBackingStore().Set("sharingCapability", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingDomainRestrictionMode sets the sharingDomainRestrictionMode property value. Specifies the external sharing mode for domains. Possible values are: none, allowList, blockList.
func (m *SharepointSettings) SetSharingDomainRestrictionMode(value *SharingDomainRestrictionMode)() {
    err := m.GetBackingStore().Set("sharingDomainRestrictionMode", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteCreationDefaultManagedPath sets the siteCreationDefaultManagedPath property value. The value of the team site managed path. This is the path under which new team sites will be created.
func (m *SharepointSettings) SetSiteCreationDefaultManagedPath(value *string)() {
    err := m.GetBackingStore().Set("siteCreationDefaultManagedPath", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteCreationDefaultStorageLimitInMB sets the siteCreationDefaultStorageLimitInMB property value. The default storage quota for a new site upon creation. Measured in megabytes (MB).
func (m *SharepointSettings) SetSiteCreationDefaultStorageLimitInMB(value *int32)() {
    err := m.GetBackingStore().Set("siteCreationDefaultStorageLimitInMB", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantDefaultTimezone sets the tenantDefaultTimezone property value. The default timezone of a tenant for newly created sites. For a list of possible values, see SPRegionalSettings.TimeZones property.
func (m *SharepointSettings) SetTenantDefaultTimezone(value *string)() {
    err := m.GetBackingStore().Set("tenantDefaultTimezone", value)
    if err != nil {
        panic(err)
    }
}
type SharepointSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedDomainGuidsForSyncApp()([]i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetAvailableManagedPathsForSiteCreation()([]string)
    GetDeletedUserPersonalSiteRetentionPeriodInDays()(*int32)
    GetExcludedFileExtensionsForSyncApp()([]string)
    GetIdleSessionSignOut()(IdleSessionSignOutable)
    GetImageTaggingOption()(*ImageTaggingChoice)
    GetIsCommentingOnSitePagesEnabled()(*bool)
    GetIsFileActivityNotificationEnabled()(*bool)
    GetIsLegacyAuthProtocolsEnabled()(*bool)
    GetIsLoopEnabled()(*bool)
    GetIsMacSyncAppEnabled()(*bool)
    GetIsRequireAcceptingUserToMatchInvitedUserEnabled()(*bool)
    GetIsResharingByExternalUsersEnabled()(*bool)
    GetIsSharePointMobileNotificationEnabled()(*bool)
    GetIsSharePointNewsfeedEnabled()(*bool)
    GetIsSiteCreationEnabled()(*bool)
    GetIsSiteCreationUIEnabled()(*bool)
    GetIsSitePagesCreationEnabled()(*bool)
    GetIsSitesStorageLimitAutomatic()(*bool)
    GetIsSyncButtonHiddenOnPersonalSite()(*bool)
    GetIsUnmanagedSyncAppForTenantRestricted()(*bool)
    GetPersonalSiteDefaultStorageLimitInMB()(*int64)
    GetSharingAllowedDomainList()([]string)
    GetSharingBlockedDomainList()([]string)
    GetSharingCapability()(*SharingCapabilities)
    GetSharingDomainRestrictionMode()(*SharingDomainRestrictionMode)
    GetSiteCreationDefaultManagedPath()(*string)
    GetSiteCreationDefaultStorageLimitInMB()(*int32)
    GetTenantDefaultTimezone()(*string)
    SetAllowedDomainGuidsForSyncApp(value []i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetAvailableManagedPathsForSiteCreation(value []string)()
    SetDeletedUserPersonalSiteRetentionPeriodInDays(value *int32)()
    SetExcludedFileExtensionsForSyncApp(value []string)()
    SetIdleSessionSignOut(value IdleSessionSignOutable)()
    SetImageTaggingOption(value *ImageTaggingChoice)()
    SetIsCommentingOnSitePagesEnabled(value *bool)()
    SetIsFileActivityNotificationEnabled(value *bool)()
    SetIsLegacyAuthProtocolsEnabled(value *bool)()
    SetIsLoopEnabled(value *bool)()
    SetIsMacSyncAppEnabled(value *bool)()
    SetIsRequireAcceptingUserToMatchInvitedUserEnabled(value *bool)()
    SetIsResharingByExternalUsersEnabled(value *bool)()
    SetIsSharePointMobileNotificationEnabled(value *bool)()
    SetIsSharePointNewsfeedEnabled(value *bool)()
    SetIsSiteCreationEnabled(value *bool)()
    SetIsSiteCreationUIEnabled(value *bool)()
    SetIsSitePagesCreationEnabled(value *bool)()
    SetIsSitesStorageLimitAutomatic(value *bool)()
    SetIsSyncButtonHiddenOnPersonalSite(value *bool)()
    SetIsUnmanagedSyncAppForTenantRestricted(value *bool)()
    SetPersonalSiteDefaultStorageLimitInMB(value *int64)()
    SetSharingAllowedDomainList(value []string)()
    SetSharingBlockedDomainList(value []string)()
    SetSharingCapability(value *SharingCapabilities)()
    SetSharingDomainRestrictionMode(value *SharingDomainRestrictionMode)()
    SetSiteCreationDefaultManagedPath(value *string)()
    SetSiteCreationDefaultStorageLimitInMB(value *int32)()
    SetTenantDefaultTimezone(value *string)()
}

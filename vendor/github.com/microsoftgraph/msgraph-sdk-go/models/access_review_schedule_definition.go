package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewScheduleDefinition struct {
    Entity
}
// NewAccessReviewScheduleDefinition instantiates a new AccessReviewScheduleDefinition and sets the default values.
func NewAccessReviewScheduleDefinition()(*AccessReviewScheduleDefinition) {
    m := &AccessReviewScheduleDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessReviewScheduleDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewScheduleDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewScheduleDefinition(), nil
}
// GetAdditionalNotificationRecipients gets the additionalNotificationRecipients property value. Defines the list of additional users or group members to be notified of the access review progress.
// returns a []AccessReviewNotificationRecipientItemable when successful
func (m *AccessReviewScheduleDefinition) GetAdditionalNotificationRecipients()([]AccessReviewNotificationRecipientItemable) {
    val, err := m.GetBackingStore().Get("additionalNotificationRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewNotificationRecipientItemable)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. User who created this review. Read-only.
// returns a UserIdentityable when successful
func (m *AccessReviewScheduleDefinition) GetCreatedBy()(UserIdentityable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserIdentityable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp when the access review series was created. Supports $select. Read-only.
// returns a *Time when successful
func (m *AccessReviewScheduleDefinition) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescriptionForAdmins gets the descriptionForAdmins property value. Description provided by review creators to provide more context of the review to admins. Supports $select.
// returns a *string when successful
func (m *AccessReviewScheduleDefinition) GetDescriptionForAdmins()(*string) {
    val, err := m.GetBackingStore().Get("descriptionForAdmins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescriptionForReviewers gets the descriptionForReviewers property value. Description provided  by review creators to provide more context of the review to reviewers. Reviewers see this description in the email sent to them requesting their review. Email notifications support up to 256 characters. Supports $select.
// returns a *string when successful
func (m *AccessReviewScheduleDefinition) GetDescriptionForReviewers()(*string) {
    val, err := m.GetBackingStore().Get("descriptionForReviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the access review series. Supports $select and $orderby. Required on create.
// returns a *string when successful
func (m *AccessReviewScheduleDefinition) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFallbackReviewers gets the fallbackReviewers property value. This collection of reviewer scopes is used to define the list of fallback reviewers. These fallback reviewers are notified to take action if no users are found from the list of reviewers specified. This could occur when either the group owner is specified as the reviewer but the group owner doesn't exist, or manager is specified as reviewer but a user's manager doesn't exist. See accessReviewReviewerScope. Replaces backupReviewers. Supports $select. NOTE: The value of this property will be ignored if fallback reviewers are assigned through the stageSettings property.
// returns a []AccessReviewReviewerScopeable when successful
func (m *AccessReviewScheduleDefinition) GetFallbackReviewers()([]AccessReviewReviewerScopeable) {
    val, err := m.GetBackingStore().Get("fallbackReviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewReviewerScopeable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewScheduleDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["additionalNotificationRecipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewNotificationRecipientItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewNotificationRecipientItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewNotificationRecipientItemable)
                }
            }
            m.SetAdditionalNotificationRecipients(res)
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(UserIdentityable))
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
    res["descriptionForAdmins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescriptionForAdmins(val)
        }
        return nil
    }
    res["descriptionForReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescriptionForReviewers(val)
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
    res["fallbackReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewReviewerScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewReviewerScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewReviewerScopeable)
                }
            }
            m.SetFallbackReviewers(res)
        }
        return nil
    }
    res["instanceEnumerationScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstanceEnumerationScope(val.(AccessReviewScopeable))
        }
        return nil
    }
    res["instances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewInstanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewInstanceable)
                }
            }
            m.SetInstances(res)
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["reviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewReviewerScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewReviewerScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewReviewerScopeable)
                }
            }
            m.SetReviewers(res)
        }
        return nil
    }
    res["scope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScope(val.(AccessReviewScopeable))
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewScheduleSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(AccessReviewScheduleSettingsable))
        }
        return nil
    }
    res["stageSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewStageSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewStageSettingsable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewStageSettingsable)
                }
            }
            m.SetStageSettings(res)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val)
        }
        return nil
    }
    return res
}
// GetInstanceEnumerationScope gets the instanceEnumerationScope property value. This property is required when scoping a review to guest users' access across all Microsoft 365 groups and determines which Microsoft 365 groups are reviewed. Each group becomes a unique accessReviewInstance of the access review series.  For supported scopes, see accessReviewScope. Supports $select. For examples of options for configuring instanceEnumerationScope, see Configure the scope of your access review definition using the Microsoft Graph API.
// returns a AccessReviewScopeable when successful
func (m *AccessReviewScheduleDefinition) GetInstanceEnumerationScope()(AccessReviewScopeable) {
    val, err := m.GetBackingStore().Get("instanceEnumerationScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewScopeable)
    }
    return nil
}
// GetInstances gets the instances property value. If the accessReviewScheduleDefinition is a recurring access review, instances represent each recurrence. A review that doesn't recur will have exactly one instance. Instances also represent each unique resource under review in the accessReviewScheduleDefinition. If a review has multiple resources and multiple instances, each resource has a unique instance for each recurrence.
// returns a []AccessReviewInstanceable when successful
func (m *AccessReviewScheduleDefinition) GetInstances()([]AccessReviewInstanceable) {
    val, err := m.GetBackingStore().Get("instances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewInstanceable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Timestamp when the access review series was last modified. Supports $select. Read-only.
// returns a *Time when successful
func (m *AccessReviewScheduleDefinition) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetReviewers gets the reviewers property value. This collection of access review scopes is used to define who are the reviewers. The reviewers property is only updatable if individual users are assigned as reviewers. Required on create. Supports $select. For examples of options for assigning reviewers, see Assign reviewers to your access review definition using the Microsoft Graph API. NOTE: The value of this property will be ignored if reviewers are assigned through the stageSettings property.
// returns a []AccessReviewReviewerScopeable when successful
func (m *AccessReviewScheduleDefinition) GetReviewers()([]AccessReviewReviewerScopeable) {
    val, err := m.GetBackingStore().Get("reviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewReviewerScopeable)
    }
    return nil
}
// GetScope gets the scope property value. Defines the entities whose access is reviewed. For supported scopes, see accessReviewScope. Required on create. Supports $select and $filter (contains only). For examples of options for configuring scope, see Configure the scope of your access review definition using the Microsoft Graph API.
// returns a AccessReviewScopeable when successful
func (m *AccessReviewScheduleDefinition) GetScope()(AccessReviewScopeable) {
    val, err := m.GetBackingStore().Get("scope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewScopeable)
    }
    return nil
}
// GetSettings gets the settings property value. The settings for an access review series, see type definition below. Supports $select. Required on create.
// returns a AccessReviewScheduleSettingsable when successful
func (m *AccessReviewScheduleDefinition) GetSettings()(AccessReviewScheduleSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewScheduleSettingsable)
    }
    return nil
}
// GetStageSettings gets the stageSettings property value. Required only for a multi-stage access review to define the stages and their settings. You can break down each review instance into up to three sequential stages, where each stage can have a different set of reviewers, fallback reviewers, and settings. Stages are created sequentially based on the dependsOn property. Optional.  When this property is defined, its settings are used instead of the corresponding settings in the accessReviewScheduleDefinition object and its settings, reviewers, and fallbackReviewers properties.
// returns a []AccessReviewStageSettingsable when successful
func (m *AccessReviewScheduleDefinition) GetStageSettings()([]AccessReviewStageSettingsable) {
    val, err := m.GetBackingStore().Get("stageSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewStageSettingsable)
    }
    return nil
}
// GetStatus gets the status property value. This read-only field specifies the status of an access review. The typical states include Initializing, NotStarted, Starting, InProgress, Completing, Completed, AutoReviewing, and AutoReviewed.  Supports $select, $orderby, and $filter (eq only). Read-only.
// returns a *string when successful
func (m *AccessReviewScheduleDefinition) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewScheduleDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAdditionalNotificationRecipients() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAdditionalNotificationRecipients()))
        for i, v := range m.GetAdditionalNotificationRecipients() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("additionalNotificationRecipients", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
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
        err = writer.WriteStringValue("descriptionForAdmins", m.GetDescriptionForAdmins())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("descriptionForReviewers", m.GetDescriptionForReviewers())
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
    if m.GetFallbackReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFallbackReviewers()))
        for i, v := range m.GetFallbackReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("fallbackReviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("instanceEnumerationScope", m.GetInstanceEnumerationScope())
        if err != nil {
            return err
        }
    }
    if m.GetInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInstances()))
        for i, v := range m.GetInstances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("instances", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReviewers()))
        for i, v := range m.GetReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("reviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("scope", m.GetScope())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    if m.GetStageSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetStageSettings()))
        for i, v := range m.GetStageSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("stageSettings", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalNotificationRecipients sets the additionalNotificationRecipients property value. Defines the list of additional users or group members to be notified of the access review progress.
func (m *AccessReviewScheduleDefinition) SetAdditionalNotificationRecipients(value []AccessReviewNotificationRecipientItemable)() {
    err := m.GetBackingStore().Set("additionalNotificationRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. User who created this review. Read-only.
func (m *AccessReviewScheduleDefinition) SetCreatedBy(value UserIdentityable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp when the access review series was created. Supports $select. Read-only.
func (m *AccessReviewScheduleDefinition) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescriptionForAdmins sets the descriptionForAdmins property value. Description provided by review creators to provide more context of the review to admins. Supports $select.
func (m *AccessReviewScheduleDefinition) SetDescriptionForAdmins(value *string)() {
    err := m.GetBackingStore().Set("descriptionForAdmins", value)
    if err != nil {
        panic(err)
    }
}
// SetDescriptionForReviewers sets the descriptionForReviewers property value. Description provided  by review creators to provide more context of the review to reviewers. Reviewers see this description in the email sent to them requesting their review. Email notifications support up to 256 characters. Supports $select.
func (m *AccessReviewScheduleDefinition) SetDescriptionForReviewers(value *string)() {
    err := m.GetBackingStore().Set("descriptionForReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the access review series. Supports $select and $orderby. Required on create.
func (m *AccessReviewScheduleDefinition) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetFallbackReviewers sets the fallbackReviewers property value. This collection of reviewer scopes is used to define the list of fallback reviewers. These fallback reviewers are notified to take action if no users are found from the list of reviewers specified. This could occur when either the group owner is specified as the reviewer but the group owner doesn't exist, or manager is specified as reviewer but a user's manager doesn't exist. See accessReviewReviewerScope. Replaces backupReviewers. Supports $select. NOTE: The value of this property will be ignored if fallback reviewers are assigned through the stageSettings property.
func (m *AccessReviewScheduleDefinition) SetFallbackReviewers(value []AccessReviewReviewerScopeable)() {
    err := m.GetBackingStore().Set("fallbackReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetInstanceEnumerationScope sets the instanceEnumerationScope property value. This property is required when scoping a review to guest users' access across all Microsoft 365 groups and determines which Microsoft 365 groups are reviewed. Each group becomes a unique accessReviewInstance of the access review series.  For supported scopes, see accessReviewScope. Supports $select. For examples of options for configuring instanceEnumerationScope, see Configure the scope of your access review definition using the Microsoft Graph API.
func (m *AccessReviewScheduleDefinition) SetInstanceEnumerationScope(value AccessReviewScopeable)() {
    err := m.GetBackingStore().Set("instanceEnumerationScope", value)
    if err != nil {
        panic(err)
    }
}
// SetInstances sets the instances property value. If the accessReviewScheduleDefinition is a recurring access review, instances represent each recurrence. A review that doesn't recur will have exactly one instance. Instances also represent each unique resource under review in the accessReviewScheduleDefinition. If a review has multiple resources and multiple instances, each resource has a unique instance for each recurrence.
func (m *AccessReviewScheduleDefinition) SetInstances(value []AccessReviewInstanceable)() {
    err := m.GetBackingStore().Set("instances", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Timestamp when the access review series was last modified. Supports $select. Read-only.
func (m *AccessReviewScheduleDefinition) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewers sets the reviewers property value. This collection of access review scopes is used to define who are the reviewers. The reviewers property is only updatable if individual users are assigned as reviewers. Required on create. Supports $select. For examples of options for assigning reviewers, see Assign reviewers to your access review definition using the Microsoft Graph API. NOTE: The value of this property will be ignored if reviewers are assigned through the stageSettings property.
func (m *AccessReviewScheduleDefinition) SetReviewers(value []AccessReviewReviewerScopeable)() {
    err := m.GetBackingStore().Set("reviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetScope sets the scope property value. Defines the entities whose access is reviewed. For supported scopes, see accessReviewScope. Required on create. Supports $select and $filter (contains only). For examples of options for configuring scope, see Configure the scope of your access review definition using the Microsoft Graph API.
func (m *AccessReviewScheduleDefinition) SetScope(value AccessReviewScopeable)() {
    err := m.GetBackingStore().Set("scope", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. The settings for an access review series, see type definition below. Supports $select. Required on create.
func (m *AccessReviewScheduleDefinition) SetSettings(value AccessReviewScheduleSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetStageSettings sets the stageSettings property value. Required only for a multi-stage access review to define the stages and their settings. You can break down each review instance into up to three sequential stages, where each stage can have a different set of reviewers, fallback reviewers, and settings. Stages are created sequentially based on the dependsOn property. Optional.  When this property is defined, its settings are used instead of the corresponding settings in the accessReviewScheduleDefinition object and its settings, reviewers, and fallbackReviewers properties.
func (m *AccessReviewScheduleDefinition) SetStageSettings(value []AccessReviewStageSettingsable)() {
    err := m.GetBackingStore().Set("stageSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. This read-only field specifies the status of an access review. The typical states include Initializing, NotStarted, Starting, InProgress, Completing, Completed, AutoReviewing, and AutoReviewed.  Supports $select, $orderby, and $filter (eq only). Read-only.
func (m *AccessReviewScheduleDefinition) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewScheduleDefinitionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdditionalNotificationRecipients()([]AccessReviewNotificationRecipientItemable)
    GetCreatedBy()(UserIdentityable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescriptionForAdmins()(*string)
    GetDescriptionForReviewers()(*string)
    GetDisplayName()(*string)
    GetFallbackReviewers()([]AccessReviewReviewerScopeable)
    GetInstanceEnumerationScope()(AccessReviewScopeable)
    GetInstances()([]AccessReviewInstanceable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetReviewers()([]AccessReviewReviewerScopeable)
    GetScope()(AccessReviewScopeable)
    GetSettings()(AccessReviewScheduleSettingsable)
    GetStageSettings()([]AccessReviewStageSettingsable)
    GetStatus()(*string)
    SetAdditionalNotificationRecipients(value []AccessReviewNotificationRecipientItemable)()
    SetCreatedBy(value UserIdentityable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescriptionForAdmins(value *string)()
    SetDescriptionForReviewers(value *string)()
    SetDisplayName(value *string)()
    SetFallbackReviewers(value []AccessReviewReviewerScopeable)()
    SetInstanceEnumerationScope(value AccessReviewScopeable)()
    SetInstances(value []AccessReviewInstanceable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetReviewers(value []AccessReviewReviewerScopeable)()
    SetScope(value AccessReviewScopeable)()
    SetSettings(value AccessReviewScheduleSettingsable)()
    SetStageSettings(value []AccessReviewStageSettingsable)()
    SetStatus(value *string)()
}

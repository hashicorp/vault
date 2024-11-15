package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUpdateForBusinessConfiguration windows Update for business configuration, allows you to specify how and when Windows as a Service updates your Windows 10/11 devices with feature and quality updates. Supports ODATA clauses that DeviceConfiguration entity supports: $filter by types of DeviceConfiguration, $top, $select only DeviceConfiguration base properties, $orderby only DeviceConfiguration base properties, and $skip. The query parameter '$search' is not supported.
type WindowsUpdateForBusinessConfiguration struct {
    DeviceConfiguration
}
// NewWindowsUpdateForBusinessConfiguration instantiates a new WindowsUpdateForBusinessConfiguration and sets the default values.
func NewWindowsUpdateForBusinessConfiguration()(*WindowsUpdateForBusinessConfiguration) {
    m := &WindowsUpdateForBusinessConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windowsUpdateForBusinessConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsUpdateForBusinessConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsUpdateForBusinessConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUpdateForBusinessConfiguration(), nil
}
// GetAllowWindows11Upgrade gets the allowWindows11Upgrade property value. When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetAllowWindows11Upgrade()(*bool) {
    val, err := m.GetBackingStore().Get("allowWindows11Upgrade")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAutomaticUpdateMode gets the automaticUpdateMode property value. Possible values for automatic update mode.
// returns a *AutomaticUpdateMode when successful
func (m *WindowsUpdateForBusinessConfiguration) GetAutomaticUpdateMode()(*AutomaticUpdateMode) {
    val, err := m.GetBackingStore().Get("automaticUpdateMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AutomaticUpdateMode)
    }
    return nil
}
// GetAutoRestartNotificationDismissal gets the autoRestartNotificationDismissal property value. Auto restart required notification dismissal method
// returns a *AutoRestartNotificationDismissalMethod when successful
func (m *WindowsUpdateForBusinessConfiguration) GetAutoRestartNotificationDismissal()(*AutoRestartNotificationDismissalMethod) {
    val, err := m.GetBackingStore().Get("autoRestartNotificationDismissal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AutoRestartNotificationDismissalMethod)
    }
    return nil
}
// GetBusinessReadyUpdatesOnly gets the businessReadyUpdatesOnly property value. Which branch devices will receive their updates from
// returns a *WindowsUpdateType when successful
func (m *WindowsUpdateForBusinessConfiguration) GetBusinessReadyUpdatesOnly()(*WindowsUpdateType) {
    val, err := m.GetBackingStore().Get("businessReadyUpdatesOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsUpdateType)
    }
    return nil
}
// GetDeadlineForFeatureUpdatesInDays gets the deadlineForFeatureUpdatesInDays property value. Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetDeadlineForFeatureUpdatesInDays()(*int32) {
    val, err := m.GetBackingStore().Get("deadlineForFeatureUpdatesInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeadlineForQualityUpdatesInDays gets the deadlineForQualityUpdatesInDays property value. Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetDeadlineForQualityUpdatesInDays()(*int32) {
    val, err := m.GetBackingStore().Get("deadlineForQualityUpdatesInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeadlineGracePeriodInDays gets the deadlineGracePeriodInDays property value. Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetDeadlineGracePeriodInDays()(*int32) {
    val, err := m.GetBackingStore().Get("deadlineGracePeriodInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeliveryOptimizationMode gets the deliveryOptimizationMode property value. Delivery optimization mode for peer distribution
// returns a *WindowsDeliveryOptimizationMode when successful
func (m *WindowsUpdateForBusinessConfiguration) GetDeliveryOptimizationMode()(*WindowsDeliveryOptimizationMode) {
    val, err := m.GetBackingStore().Get("deliveryOptimizationMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsDeliveryOptimizationMode)
    }
    return nil
}
// GetDriversExcluded gets the driversExcluded property value. When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetDriversExcluded()(*bool) {
    val, err := m.GetBackingStore().Get("driversExcluded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEngagedRestartDeadlineInDays gets the engagedRestartDeadlineInDays property value. Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetEngagedRestartDeadlineInDays()(*int32) {
    val, err := m.GetBackingStore().Get("engagedRestartDeadlineInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEngagedRestartSnoozeScheduleInDays gets the engagedRestartSnoozeScheduleInDays property value. Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetEngagedRestartSnoozeScheduleInDays()(*int32) {
    val, err := m.GetBackingStore().Get("engagedRestartSnoozeScheduleInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEngagedRestartTransitionScheduleInDays gets the engagedRestartTransitionScheduleInDays property value. Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetEngagedRestartTransitionScheduleInDays()(*int32) {
    val, err := m.GetBackingStore().Get("engagedRestartTransitionScheduleInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFeatureUpdatesDeferralPeriodInDays gets the featureUpdatesDeferralPeriodInDays property value. Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesDeferralPeriodInDays()(*int32) {
    val, err := m.GetBackingStore().Get("featureUpdatesDeferralPeriodInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFeatureUpdatesPaused gets the featureUpdatesPaused property value. When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesPaused()(*bool) {
    val, err := m.GetBackingStore().Get("featureUpdatesPaused")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFeatureUpdatesPauseExpiryDateTime gets the featureUpdatesPauseExpiryDateTime property value. The Feature Updates Pause Expiry datetime. This value is 35 days from the time admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported.
// returns a *Time when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesPauseExpiryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("featureUpdatesPauseExpiryDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFeatureUpdatesPauseStartDate gets the featureUpdatesPauseStartDate property value. The Feature Updates Pause start date. This value is the time when the admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported. This property is read-only.
// returns a *DateOnly when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesPauseStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("featureUpdatesPauseStartDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// GetFeatureUpdatesRollbackStartDateTime gets the featureUpdatesRollbackStartDateTime property value. The Feature Updates Rollback Start datetime.This value is the time when the admin rolled back the Feature update for the ring.Returned by default.Query parameters are not supported.
// returns a *Time when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesRollbackStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("featureUpdatesRollbackStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFeatureUpdatesRollbackWindowInDays gets the featureUpdatesRollbackWindowInDays property value. The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesRollbackWindowInDays()(*int32) {
    val, err := m.GetBackingStore().Get("featureUpdatesRollbackWindowInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFeatureUpdatesWillBeRolledBack gets the featureUpdatesWillBeRolledBack property value. When TRUE, rollback Feature Updates on the next device check in. When FALSE, do not rollback Feature Updates on the next device check in. Returned by default.Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFeatureUpdatesWillBeRolledBack()(*bool) {
    val, err := m.GetBackingStore().Get("featureUpdatesWillBeRolledBack")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsUpdateForBusinessConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["allowWindows11Upgrade"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowWindows11Upgrade(val)
        }
        return nil
    }
    res["automaticUpdateMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAutomaticUpdateMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomaticUpdateMode(val.(*AutomaticUpdateMode))
        }
        return nil
    }
    res["autoRestartNotificationDismissal"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAutoRestartNotificationDismissalMethod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutoRestartNotificationDismissal(val.(*AutoRestartNotificationDismissalMethod))
        }
        return nil
    }
    res["businessReadyUpdatesOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsUpdateType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBusinessReadyUpdatesOnly(val.(*WindowsUpdateType))
        }
        return nil
    }
    res["deadlineForFeatureUpdatesInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeadlineForFeatureUpdatesInDays(val)
        }
        return nil
    }
    res["deadlineForQualityUpdatesInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeadlineForQualityUpdatesInDays(val)
        }
        return nil
    }
    res["deadlineGracePeriodInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeadlineGracePeriodInDays(val)
        }
        return nil
    }
    res["deliveryOptimizationMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsDeliveryOptimizationMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryOptimizationMode(val.(*WindowsDeliveryOptimizationMode))
        }
        return nil
    }
    res["driversExcluded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriversExcluded(val)
        }
        return nil
    }
    res["engagedRestartDeadlineInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEngagedRestartDeadlineInDays(val)
        }
        return nil
    }
    res["engagedRestartSnoozeScheduleInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEngagedRestartSnoozeScheduleInDays(val)
        }
        return nil
    }
    res["engagedRestartTransitionScheduleInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEngagedRestartTransitionScheduleInDays(val)
        }
        return nil
    }
    res["featureUpdatesDeferralPeriodInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesDeferralPeriodInDays(val)
        }
        return nil
    }
    res["featureUpdatesPaused"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesPaused(val)
        }
        return nil
    }
    res["featureUpdatesPauseExpiryDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesPauseExpiryDateTime(val)
        }
        return nil
    }
    res["featureUpdatesPauseStartDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesPauseStartDate(val)
        }
        return nil
    }
    res["featureUpdatesRollbackStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesRollbackStartDateTime(val)
        }
        return nil
    }
    res["featureUpdatesRollbackWindowInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesRollbackWindowInDays(val)
        }
        return nil
    }
    res["featureUpdatesWillBeRolledBack"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdatesWillBeRolledBack(val)
        }
        return nil
    }
    res["installationSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsUpdateInstallScheduleTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallationSchedule(val.(WindowsUpdateInstallScheduleTypeable))
        }
        return nil
    }
    res["microsoftUpdateServiceAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftUpdateServiceAllowed(val)
        }
        return nil
    }
    res["postponeRebootUntilAfterDeadline"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPostponeRebootUntilAfterDeadline(val)
        }
        return nil
    }
    res["prereleaseFeatures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrereleaseFeatures)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrereleaseFeatures(val.(*PrereleaseFeatures))
        }
        return nil
    }
    res["qualityUpdatesDeferralPeriodInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesDeferralPeriodInDays(val)
        }
        return nil
    }
    res["qualityUpdatesPaused"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesPaused(val)
        }
        return nil
    }
    res["qualityUpdatesPauseExpiryDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesPauseExpiryDateTime(val)
        }
        return nil
    }
    res["qualityUpdatesPauseStartDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesPauseStartDate(val)
        }
        return nil
    }
    res["qualityUpdatesRollbackStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesRollbackStartDateTime(val)
        }
        return nil
    }
    res["qualityUpdatesWillBeRolledBack"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQualityUpdatesWillBeRolledBack(val)
        }
        return nil
    }
    res["scheduleImminentRestartWarningInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduleImminentRestartWarningInMinutes(val)
        }
        return nil
    }
    res["scheduleRestartWarningInHours"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduleRestartWarningInHours(val)
        }
        return nil
    }
    res["skipChecksBeforeRestart"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSkipChecksBeforeRestart(val)
        }
        return nil
    }
    res["updateNotificationLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsUpdateNotificationDisplayOption)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpdateNotificationLevel(val.(*WindowsUpdateNotificationDisplayOption))
        }
        return nil
    }
    res["updateWeeks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsUpdateForBusinessUpdateWeeks)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpdateWeeks(val.(*WindowsUpdateForBusinessUpdateWeeks))
        }
        return nil
    }
    res["userPauseAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEnablement)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPauseAccess(val.(*Enablement))
        }
        return nil
    }
    res["userWindowsUpdateScanAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEnablement)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserWindowsUpdateScanAccess(val.(*Enablement))
        }
        return nil
    }
    return res
}
// GetInstallationSchedule gets the installationSchedule property value. The Installation Schedule. Possible values are: ActiveHoursStart, ActiveHoursEnd, ScheduledInstallDay, ScheduledInstallTime. Returned by default. Query parameters are not supported.
// returns a WindowsUpdateInstallScheduleTypeable when successful
func (m *WindowsUpdateForBusinessConfiguration) GetInstallationSchedule()(WindowsUpdateInstallScheduleTypeable) {
    val, err := m.GetBackingStore().Get("installationSchedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsUpdateInstallScheduleTypeable)
    }
    return nil
}
// GetMicrosoftUpdateServiceAllowed gets the microsoftUpdateServiceAllowed property value. When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetMicrosoftUpdateServiceAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("microsoftUpdateServiceAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPostponeRebootUntilAfterDeadline gets the postponeRebootUntilAfterDeadline property value. When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetPostponeRebootUntilAfterDeadline()(*bool) {
    val, err := m.GetBackingStore().Get("postponeRebootUntilAfterDeadline")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPrereleaseFeatures gets the prereleaseFeatures property value. Possible values for pre-release features.
// returns a *PrereleaseFeatures when successful
func (m *WindowsUpdateForBusinessConfiguration) GetPrereleaseFeatures()(*PrereleaseFeatures) {
    val, err := m.GetBackingStore().Get("prereleaseFeatures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrereleaseFeatures)
    }
    return nil
}
// GetQualityUpdatesDeferralPeriodInDays gets the qualityUpdatesDeferralPeriodInDays property value. Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesDeferralPeriodInDays()(*int32) {
    val, err := m.GetBackingStore().Get("qualityUpdatesDeferralPeriodInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetQualityUpdatesPaused gets the qualityUpdatesPaused property value. When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesPaused()(*bool) {
    val, err := m.GetBackingStore().Get("qualityUpdatesPaused")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetQualityUpdatesPauseExpiryDateTime gets the qualityUpdatesPauseExpiryDateTime property value. The Quality Updates Pause Expiry datetime. This value is 35 days from the time admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported.
// returns a *Time when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesPauseExpiryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("qualityUpdatesPauseExpiryDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetQualityUpdatesPauseStartDate gets the qualityUpdatesPauseStartDate property value. The Quality Updates Pause start date. This value is the time when the admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported. This property is read-only.
// returns a *DateOnly when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesPauseStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("qualityUpdatesPauseStartDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// GetQualityUpdatesRollbackStartDateTime gets the qualityUpdatesRollbackStartDateTime property value. The Quality Updates Rollback Start datetime. This value is the time when the admin rolled back the Quality update for the ring. Returned by default. Query parameters are not supported.
// returns a *Time when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesRollbackStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("qualityUpdatesRollbackStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetQualityUpdatesWillBeRolledBack gets the qualityUpdatesWillBeRolledBack property value. When TRUE, rollback Quality Updates on the next device check in. When FALSE, do not rollback Quality Updates on the next device check in. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetQualityUpdatesWillBeRolledBack()(*bool) {
    val, err := m.GetBackingStore().Get("qualityUpdatesWillBeRolledBack")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetScheduleImminentRestartWarningInMinutes gets the scheduleImminentRestartWarningInMinutes property value. Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetScheduleImminentRestartWarningInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("scheduleImminentRestartWarningInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetScheduleRestartWarningInHours gets the scheduleRestartWarningInHours property value. Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.
// returns a *int32 when successful
func (m *WindowsUpdateForBusinessConfiguration) GetScheduleRestartWarningInHours()(*int32) {
    val, err := m.GetBackingStore().Get("scheduleRestartWarningInHours")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSkipChecksBeforeRestart gets the skipChecksBeforeRestart property value. When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.
// returns a *bool when successful
func (m *WindowsUpdateForBusinessConfiguration) GetSkipChecksBeforeRestart()(*bool) {
    val, err := m.GetBackingStore().Get("skipChecksBeforeRestart")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUpdateNotificationLevel gets the updateNotificationLevel property value. Windows Update Notification Display Options
// returns a *WindowsUpdateNotificationDisplayOption when successful
func (m *WindowsUpdateForBusinessConfiguration) GetUpdateNotificationLevel()(*WindowsUpdateNotificationDisplayOption) {
    val, err := m.GetBackingStore().Get("updateNotificationLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsUpdateNotificationDisplayOption)
    }
    return nil
}
// GetUpdateWeeks gets the updateWeeks property value. Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. Possible values are: userDefined, firstWeek, secondWeek, thirdWeek, fourthWeek, everyWeek, unknownFutureValue.
// returns a *WindowsUpdateForBusinessUpdateWeeks when successful
func (m *WindowsUpdateForBusinessConfiguration) GetUpdateWeeks()(*WindowsUpdateForBusinessUpdateWeeks) {
    val, err := m.GetBackingStore().Get("updateWeeks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsUpdateForBusinessUpdateWeeks)
    }
    return nil
}
// GetUserPauseAccess gets the userPauseAccess property value. Possible values of a property
// returns a *Enablement when successful
func (m *WindowsUpdateForBusinessConfiguration) GetUserPauseAccess()(*Enablement) {
    val, err := m.GetBackingStore().Get("userPauseAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Enablement)
    }
    return nil
}
// GetUserWindowsUpdateScanAccess gets the userWindowsUpdateScanAccess property value. Possible values of a property
// returns a *Enablement when successful
func (m *WindowsUpdateForBusinessConfiguration) GetUserWindowsUpdateScanAccess()(*Enablement) {
    val, err := m.GetBackingStore().Get("userWindowsUpdateScanAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Enablement)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsUpdateForBusinessConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowWindows11Upgrade", m.GetAllowWindows11Upgrade())
        if err != nil {
            return err
        }
    }
    if m.GetAutomaticUpdateMode() != nil {
        cast := (*m.GetAutomaticUpdateMode()).String()
        err = writer.WriteStringValue("automaticUpdateMode", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAutoRestartNotificationDismissal() != nil {
        cast := (*m.GetAutoRestartNotificationDismissal()).String()
        err = writer.WriteStringValue("autoRestartNotificationDismissal", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetBusinessReadyUpdatesOnly() != nil {
        cast := (*m.GetBusinessReadyUpdatesOnly()).String()
        err = writer.WriteStringValue("businessReadyUpdatesOnly", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deadlineForFeatureUpdatesInDays", m.GetDeadlineForFeatureUpdatesInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deadlineForQualityUpdatesInDays", m.GetDeadlineForQualityUpdatesInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deadlineGracePeriodInDays", m.GetDeadlineGracePeriodInDays())
        if err != nil {
            return err
        }
    }
    if m.GetDeliveryOptimizationMode() != nil {
        cast := (*m.GetDeliveryOptimizationMode()).String()
        err = writer.WriteStringValue("deliveryOptimizationMode", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("driversExcluded", m.GetDriversExcluded())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("engagedRestartDeadlineInDays", m.GetEngagedRestartDeadlineInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("engagedRestartSnoozeScheduleInDays", m.GetEngagedRestartSnoozeScheduleInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("engagedRestartTransitionScheduleInDays", m.GetEngagedRestartTransitionScheduleInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("featureUpdatesDeferralPeriodInDays", m.GetFeatureUpdatesDeferralPeriodInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("featureUpdatesPaused", m.GetFeatureUpdatesPaused())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("featureUpdatesPauseExpiryDateTime", m.GetFeatureUpdatesPauseExpiryDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("featureUpdatesRollbackStartDateTime", m.GetFeatureUpdatesRollbackStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("featureUpdatesRollbackWindowInDays", m.GetFeatureUpdatesRollbackWindowInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("featureUpdatesWillBeRolledBack", m.GetFeatureUpdatesWillBeRolledBack())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("installationSchedule", m.GetInstallationSchedule())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("microsoftUpdateServiceAllowed", m.GetMicrosoftUpdateServiceAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("postponeRebootUntilAfterDeadline", m.GetPostponeRebootUntilAfterDeadline())
        if err != nil {
            return err
        }
    }
    if m.GetPrereleaseFeatures() != nil {
        cast := (*m.GetPrereleaseFeatures()).String()
        err = writer.WriteStringValue("prereleaseFeatures", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("qualityUpdatesDeferralPeriodInDays", m.GetQualityUpdatesDeferralPeriodInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("qualityUpdatesPaused", m.GetQualityUpdatesPaused())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("qualityUpdatesPauseExpiryDateTime", m.GetQualityUpdatesPauseExpiryDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("qualityUpdatesRollbackStartDateTime", m.GetQualityUpdatesRollbackStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("qualityUpdatesWillBeRolledBack", m.GetQualityUpdatesWillBeRolledBack())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("scheduleImminentRestartWarningInMinutes", m.GetScheduleImminentRestartWarningInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("scheduleRestartWarningInHours", m.GetScheduleRestartWarningInHours())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("skipChecksBeforeRestart", m.GetSkipChecksBeforeRestart())
        if err != nil {
            return err
        }
    }
    if m.GetUpdateNotificationLevel() != nil {
        cast := (*m.GetUpdateNotificationLevel()).String()
        err = writer.WriteStringValue("updateNotificationLevel", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetUpdateWeeks() != nil {
        cast := (*m.GetUpdateWeeks()).String()
        err = writer.WriteStringValue("updateWeeks", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserPauseAccess() != nil {
        cast := (*m.GetUserPauseAccess()).String()
        err = writer.WriteStringValue("userPauseAccess", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserWindowsUpdateScanAccess() != nil {
        cast := (*m.GetUserWindowsUpdateScanAccess()).String()
        err = writer.WriteStringValue("userWindowsUpdateScanAccess", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowWindows11Upgrade sets the allowWindows11Upgrade property value. When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetAllowWindows11Upgrade(value *bool)() {
    err := m.GetBackingStore().Set("allowWindows11Upgrade", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomaticUpdateMode sets the automaticUpdateMode property value. Possible values for automatic update mode.
func (m *WindowsUpdateForBusinessConfiguration) SetAutomaticUpdateMode(value *AutomaticUpdateMode)() {
    err := m.GetBackingStore().Set("automaticUpdateMode", value)
    if err != nil {
        panic(err)
    }
}
// SetAutoRestartNotificationDismissal sets the autoRestartNotificationDismissal property value. Auto restart required notification dismissal method
func (m *WindowsUpdateForBusinessConfiguration) SetAutoRestartNotificationDismissal(value *AutoRestartNotificationDismissalMethod)() {
    err := m.GetBackingStore().Set("autoRestartNotificationDismissal", value)
    if err != nil {
        panic(err)
    }
}
// SetBusinessReadyUpdatesOnly sets the businessReadyUpdatesOnly property value. Which branch devices will receive their updates from
func (m *WindowsUpdateForBusinessConfiguration) SetBusinessReadyUpdatesOnly(value *WindowsUpdateType)() {
    err := m.GetBackingStore().Set("businessReadyUpdatesOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetDeadlineForFeatureUpdatesInDays sets the deadlineForFeatureUpdatesInDays property value. Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetDeadlineForFeatureUpdatesInDays(value *int32)() {
    err := m.GetBackingStore().Set("deadlineForFeatureUpdatesInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetDeadlineForQualityUpdatesInDays sets the deadlineForQualityUpdatesInDays property value. Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetDeadlineForQualityUpdatesInDays(value *int32)() {
    err := m.GetBackingStore().Set("deadlineForQualityUpdatesInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetDeadlineGracePeriodInDays sets the deadlineGracePeriodInDays property value. Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetDeadlineGracePeriodInDays(value *int32)() {
    err := m.GetBackingStore().Set("deadlineGracePeriodInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryOptimizationMode sets the deliveryOptimizationMode property value. Delivery optimization mode for peer distribution
func (m *WindowsUpdateForBusinessConfiguration) SetDeliveryOptimizationMode(value *WindowsDeliveryOptimizationMode)() {
    err := m.GetBackingStore().Set("deliveryOptimizationMode", value)
    if err != nil {
        panic(err)
    }
}
// SetDriversExcluded sets the driversExcluded property value. When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetDriversExcluded(value *bool)() {
    err := m.GetBackingStore().Set("driversExcluded", value)
    if err != nil {
        panic(err)
    }
}
// SetEngagedRestartDeadlineInDays sets the engagedRestartDeadlineInDays property value. Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetEngagedRestartDeadlineInDays(value *int32)() {
    err := m.GetBackingStore().Set("engagedRestartDeadlineInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetEngagedRestartSnoozeScheduleInDays sets the engagedRestartSnoozeScheduleInDays property value. Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetEngagedRestartSnoozeScheduleInDays(value *int32)() {
    err := m.GetBackingStore().Set("engagedRestartSnoozeScheduleInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetEngagedRestartTransitionScheduleInDays sets the engagedRestartTransitionScheduleInDays property value. Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetEngagedRestartTransitionScheduleInDays(value *int32)() {
    err := m.GetBackingStore().Set("engagedRestartTransitionScheduleInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesDeferralPeriodInDays sets the featureUpdatesDeferralPeriodInDays property value. Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesDeferralPeriodInDays(value *int32)() {
    err := m.GetBackingStore().Set("featureUpdatesDeferralPeriodInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesPaused sets the featureUpdatesPaused property value. When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesPaused(value *bool)() {
    err := m.GetBackingStore().Set("featureUpdatesPaused", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesPauseExpiryDateTime sets the featureUpdatesPauseExpiryDateTime property value. The Feature Updates Pause Expiry datetime. This value is 35 days from the time admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesPauseExpiryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("featureUpdatesPauseExpiryDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesPauseStartDate sets the featureUpdatesPauseStartDate property value. The Feature Updates Pause start date. This value is the time when the admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported. This property is read-only.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesPauseStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("featureUpdatesPauseStartDate", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesRollbackStartDateTime sets the featureUpdatesRollbackStartDateTime property value. The Feature Updates Rollback Start datetime.This value is the time when the admin rolled back the Feature update for the ring.Returned by default.Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesRollbackStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("featureUpdatesRollbackStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesRollbackWindowInDays sets the featureUpdatesRollbackWindowInDays property value. The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesRollbackWindowInDays(value *int32)() {
    err := m.GetBackingStore().Set("featureUpdatesRollbackWindowInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdatesWillBeRolledBack sets the featureUpdatesWillBeRolledBack property value. When TRUE, rollback Feature Updates on the next device check in. When FALSE, do not rollback Feature Updates on the next device check in. Returned by default.Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetFeatureUpdatesWillBeRolledBack(value *bool)() {
    err := m.GetBackingStore().Set("featureUpdatesWillBeRolledBack", value)
    if err != nil {
        panic(err)
    }
}
// SetInstallationSchedule sets the installationSchedule property value. The Installation Schedule. Possible values are: ActiveHoursStart, ActiveHoursEnd, ScheduledInstallDay, ScheduledInstallTime. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetInstallationSchedule(value WindowsUpdateInstallScheduleTypeable)() {
    err := m.GetBackingStore().Set("installationSchedule", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftUpdateServiceAllowed sets the microsoftUpdateServiceAllowed property value. When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetMicrosoftUpdateServiceAllowed(value *bool)() {
    err := m.GetBackingStore().Set("microsoftUpdateServiceAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetPostponeRebootUntilAfterDeadline sets the postponeRebootUntilAfterDeadline property value. When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetPostponeRebootUntilAfterDeadline(value *bool)() {
    err := m.GetBackingStore().Set("postponeRebootUntilAfterDeadline", value)
    if err != nil {
        panic(err)
    }
}
// SetPrereleaseFeatures sets the prereleaseFeatures property value. Possible values for pre-release features.
func (m *WindowsUpdateForBusinessConfiguration) SetPrereleaseFeatures(value *PrereleaseFeatures)() {
    err := m.GetBackingStore().Set("prereleaseFeatures", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesDeferralPeriodInDays sets the qualityUpdatesDeferralPeriodInDays property value. Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesDeferralPeriodInDays(value *int32)() {
    err := m.GetBackingStore().Set("qualityUpdatesDeferralPeriodInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesPaused sets the qualityUpdatesPaused property value. When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesPaused(value *bool)() {
    err := m.GetBackingStore().Set("qualityUpdatesPaused", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesPauseExpiryDateTime sets the qualityUpdatesPauseExpiryDateTime property value. The Quality Updates Pause Expiry datetime. This value is 35 days from the time admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesPauseExpiryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("qualityUpdatesPauseExpiryDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesPauseStartDate sets the qualityUpdatesPauseStartDate property value. The Quality Updates Pause start date. This value is the time when the admin paused or extended the pause for the ring. Returned by default. Query parameters are not supported. This property is read-only.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesPauseStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("qualityUpdatesPauseStartDate", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesRollbackStartDateTime sets the qualityUpdatesRollbackStartDateTime property value. The Quality Updates Rollback Start datetime. This value is the time when the admin rolled back the Quality update for the ring. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesRollbackStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("qualityUpdatesRollbackStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetQualityUpdatesWillBeRolledBack sets the qualityUpdatesWillBeRolledBack property value. When TRUE, rollback Quality Updates on the next device check in. When FALSE, do not rollback Quality Updates on the next device check in. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetQualityUpdatesWillBeRolledBack(value *bool)() {
    err := m.GetBackingStore().Set("qualityUpdatesWillBeRolledBack", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleImminentRestartWarningInMinutes sets the scheduleImminentRestartWarningInMinutes property value. Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetScheduleImminentRestartWarningInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("scheduleImminentRestartWarningInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleRestartWarningInHours sets the scheduleRestartWarningInHours property value. Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetScheduleRestartWarningInHours(value *int32)() {
    err := m.GetBackingStore().Set("scheduleRestartWarningInHours", value)
    if err != nil {
        panic(err)
    }
}
// SetSkipChecksBeforeRestart sets the skipChecksBeforeRestart property value. When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.
func (m *WindowsUpdateForBusinessConfiguration) SetSkipChecksBeforeRestart(value *bool)() {
    err := m.GetBackingStore().Set("skipChecksBeforeRestart", value)
    if err != nil {
        panic(err)
    }
}
// SetUpdateNotificationLevel sets the updateNotificationLevel property value. Windows Update Notification Display Options
func (m *WindowsUpdateForBusinessConfiguration) SetUpdateNotificationLevel(value *WindowsUpdateNotificationDisplayOption)() {
    err := m.GetBackingStore().Set("updateNotificationLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetUpdateWeeks sets the updateWeeks property value. Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. Possible values are: userDefined, firstWeek, secondWeek, thirdWeek, fourthWeek, everyWeek, unknownFutureValue.
func (m *WindowsUpdateForBusinessConfiguration) SetUpdateWeeks(value *WindowsUpdateForBusinessUpdateWeeks)() {
    err := m.GetBackingStore().Set("updateWeeks", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPauseAccess sets the userPauseAccess property value. Possible values of a property
func (m *WindowsUpdateForBusinessConfiguration) SetUserPauseAccess(value *Enablement)() {
    err := m.GetBackingStore().Set("userPauseAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetUserWindowsUpdateScanAccess sets the userWindowsUpdateScanAccess property value. Possible values of a property
func (m *WindowsUpdateForBusinessConfiguration) SetUserWindowsUpdateScanAccess(value *Enablement)() {
    err := m.GetBackingStore().Set("userWindowsUpdateScanAccess", value)
    if err != nil {
        panic(err)
    }
}
type WindowsUpdateForBusinessConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowWindows11Upgrade()(*bool)
    GetAutomaticUpdateMode()(*AutomaticUpdateMode)
    GetAutoRestartNotificationDismissal()(*AutoRestartNotificationDismissalMethod)
    GetBusinessReadyUpdatesOnly()(*WindowsUpdateType)
    GetDeadlineForFeatureUpdatesInDays()(*int32)
    GetDeadlineForQualityUpdatesInDays()(*int32)
    GetDeadlineGracePeriodInDays()(*int32)
    GetDeliveryOptimizationMode()(*WindowsDeliveryOptimizationMode)
    GetDriversExcluded()(*bool)
    GetEngagedRestartDeadlineInDays()(*int32)
    GetEngagedRestartSnoozeScheduleInDays()(*int32)
    GetEngagedRestartTransitionScheduleInDays()(*int32)
    GetFeatureUpdatesDeferralPeriodInDays()(*int32)
    GetFeatureUpdatesPaused()(*bool)
    GetFeatureUpdatesPauseExpiryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFeatureUpdatesPauseStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    GetFeatureUpdatesRollbackStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFeatureUpdatesRollbackWindowInDays()(*int32)
    GetFeatureUpdatesWillBeRolledBack()(*bool)
    GetInstallationSchedule()(WindowsUpdateInstallScheduleTypeable)
    GetMicrosoftUpdateServiceAllowed()(*bool)
    GetPostponeRebootUntilAfterDeadline()(*bool)
    GetPrereleaseFeatures()(*PrereleaseFeatures)
    GetQualityUpdatesDeferralPeriodInDays()(*int32)
    GetQualityUpdatesPaused()(*bool)
    GetQualityUpdatesPauseExpiryDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetQualityUpdatesPauseStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    GetQualityUpdatesRollbackStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetQualityUpdatesWillBeRolledBack()(*bool)
    GetScheduleImminentRestartWarningInMinutes()(*int32)
    GetScheduleRestartWarningInHours()(*int32)
    GetSkipChecksBeforeRestart()(*bool)
    GetUpdateNotificationLevel()(*WindowsUpdateNotificationDisplayOption)
    GetUpdateWeeks()(*WindowsUpdateForBusinessUpdateWeeks)
    GetUserPauseAccess()(*Enablement)
    GetUserWindowsUpdateScanAccess()(*Enablement)
    SetAllowWindows11Upgrade(value *bool)()
    SetAutomaticUpdateMode(value *AutomaticUpdateMode)()
    SetAutoRestartNotificationDismissal(value *AutoRestartNotificationDismissalMethod)()
    SetBusinessReadyUpdatesOnly(value *WindowsUpdateType)()
    SetDeadlineForFeatureUpdatesInDays(value *int32)()
    SetDeadlineForQualityUpdatesInDays(value *int32)()
    SetDeadlineGracePeriodInDays(value *int32)()
    SetDeliveryOptimizationMode(value *WindowsDeliveryOptimizationMode)()
    SetDriversExcluded(value *bool)()
    SetEngagedRestartDeadlineInDays(value *int32)()
    SetEngagedRestartSnoozeScheduleInDays(value *int32)()
    SetEngagedRestartTransitionScheduleInDays(value *int32)()
    SetFeatureUpdatesDeferralPeriodInDays(value *int32)()
    SetFeatureUpdatesPaused(value *bool)()
    SetFeatureUpdatesPauseExpiryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFeatureUpdatesPauseStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
    SetFeatureUpdatesRollbackStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFeatureUpdatesRollbackWindowInDays(value *int32)()
    SetFeatureUpdatesWillBeRolledBack(value *bool)()
    SetInstallationSchedule(value WindowsUpdateInstallScheduleTypeable)()
    SetMicrosoftUpdateServiceAllowed(value *bool)()
    SetPostponeRebootUntilAfterDeadline(value *bool)()
    SetPrereleaseFeatures(value *PrereleaseFeatures)()
    SetQualityUpdatesDeferralPeriodInDays(value *int32)()
    SetQualityUpdatesPaused(value *bool)()
    SetQualityUpdatesPauseExpiryDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetQualityUpdatesPauseStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
    SetQualityUpdatesRollbackStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetQualityUpdatesWillBeRolledBack(value *bool)()
    SetScheduleImminentRestartWarningInMinutes(value *int32)()
    SetScheduleRestartWarningInHours(value *int32)()
    SetSkipChecksBeforeRestart(value *bool)()
    SetUpdateNotificationLevel(value *WindowsUpdateNotificationDisplayOption)()
    SetUpdateWeeks(value *WindowsUpdateForBusinessUpdateWeeks)()
    SetUserPauseAccess(value *Enablement)()
    SetUserWindowsUpdateScanAccess(value *Enablement)()
}

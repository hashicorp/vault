package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsDeviceStartupHistory the user experience analytics device startup history entity contains device boot performance history details.
type UserExperienceAnalyticsDeviceStartupHistory struct {
    Entity
}
// NewUserExperienceAnalyticsDeviceStartupHistory instantiates a new UserExperienceAnalyticsDeviceStartupHistory and sets the default values.
func NewUserExperienceAnalyticsDeviceStartupHistory()(*UserExperienceAnalyticsDeviceStartupHistory) {
    m := &UserExperienceAnalyticsDeviceStartupHistory{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsDeviceStartupHistoryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsDeviceStartupHistoryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsDeviceStartupHistory(), nil
}
// GetCoreBootTimeInMs gets the coreBootTimeInMs property value. The device core boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetCoreBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("coreBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCoreLoginTimeInMs gets the coreLoginTimeInMs property value. The device core login time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetCoreLoginTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("coreLoginTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeviceId gets the deviceId property value. The Intune device id of the device. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("deviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFeatureUpdateBootTimeInMs gets the featureUpdateBootTimeInMs property value. The impact of device feature updates on boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetFeatureUpdateBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("featureUpdateBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["coreBootTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCoreBootTimeInMs(val)
        }
        return nil
    }
    res["coreLoginTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCoreLoginTimeInMs(val)
        }
        return nil
    }
    res["deviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceId(val)
        }
        return nil
    }
    res["featureUpdateBootTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFeatureUpdateBootTimeInMs(val)
        }
        return nil
    }
    res["groupPolicyBootTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupPolicyBootTimeInMs(val)
        }
        return nil
    }
    res["groupPolicyLoginTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupPolicyLoginTimeInMs(val)
        }
        return nil
    }
    res["isFeatureUpdate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFeatureUpdate(val)
        }
        return nil
    }
    res["isFirstLogin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFirstLogin(val)
        }
        return nil
    }
    res["operatingSystemVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystemVersion(val)
        }
        return nil
    }
    res["responsiveDesktopTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResponsiveDesktopTimeInMs(val)
        }
        return nil
    }
    res["restartCategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserExperienceAnalyticsOperatingSystemRestartCategory)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestartCategory(val.(*UserExperienceAnalyticsOperatingSystemRestartCategory))
        }
        return nil
    }
    res["restartFaultBucket"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestartFaultBucket(val)
        }
        return nil
    }
    res["restartStopCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestartStopCode(val)
        }
        return nil
    }
    res["startTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartTime(val)
        }
        return nil
    }
    res["totalBootTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalBootTimeInMs(val)
        }
        return nil
    }
    res["totalLoginTimeInMs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalLoginTimeInMs(val)
        }
        return nil
    }
    return res
}
// GetGroupPolicyBootTimeInMs gets the groupPolicyBootTimeInMs property value. The impact of device group policy client on boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetGroupPolicyBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("groupPolicyBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetGroupPolicyLoginTimeInMs gets the groupPolicyLoginTimeInMs property value. The impact of device group policy client on login time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetGroupPolicyLoginTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("groupPolicyLoginTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetIsFeatureUpdate gets the isFeatureUpdate property value. When TRUE, indicates the device boot record is associated with feature updates. When FALSE, indicates the device boot record is not associated with feature updates. Supports: $select, $OrderBy. Read-only.
// returns a *bool when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetIsFeatureUpdate()(*bool) {
    val, err := m.GetBackingStore().Get("isFeatureUpdate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsFirstLogin gets the isFirstLogin property value. When TRUE, indicates the device login is the first login after a reboot. When FALSE, indicates the device login is not the first login after a reboot. Supports: $select, $OrderBy. Read-only.
// returns a *bool when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetIsFirstLogin()(*bool) {
    val, err := m.GetBackingStore().Get("isFirstLogin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOperatingSystemVersion gets the operatingSystemVersion property value. The user experience analytics device boot record's operating system version. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetOperatingSystemVersion()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystemVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResponsiveDesktopTimeInMs gets the responsiveDesktopTimeInMs property value. The time for desktop to become responsive during login process in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetResponsiveDesktopTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("responsiveDesktopTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRestartCategory gets the restartCategory property value. Operating System restart category.
// returns a *UserExperienceAnalyticsOperatingSystemRestartCategory when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetRestartCategory()(*UserExperienceAnalyticsOperatingSystemRestartCategory) {
    val, err := m.GetBackingStore().Get("restartCategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserExperienceAnalyticsOperatingSystemRestartCategory)
    }
    return nil
}
// GetRestartFaultBucket gets the restartFaultBucket property value. OS restart fault bucket. The fault bucket is used to find additional information about a system crash. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetRestartFaultBucket()(*string) {
    val, err := m.GetBackingStore().Get("restartFaultBucket")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRestartStopCode gets the restartStopCode property value. OS restart stop code. This shows the bug check code which can be used to look up the blue screen reason. Supports: $select, $OrderBy. Read-only.
// returns a *string when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetRestartStopCode()(*string) {
    val, err := m.GetBackingStore().Get("restartStopCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartTime gets the startTime property value. The device boot start time. The value cannot be modified and is automatically populated when the device performs a reboot. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2022 would look like this: '2022-01-01T00:00:00Z'. Returned by default. Read-only.
// returns a *Time when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetStartTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTotalBootTimeInMs gets the totalBootTimeInMs property value. The device total boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetTotalBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("totalBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalLoginTimeInMs gets the totalLoginTimeInMs property value. The device total login time in milliseconds. Supports: $select, $OrderBy. Read-only.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDeviceStartupHistory) GetTotalLoginTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("totalLoginTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsDeviceStartupHistory) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("coreBootTimeInMs", m.GetCoreBootTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("coreLoginTimeInMs", m.GetCoreLoginTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("featureUpdateBootTimeInMs", m.GetFeatureUpdateBootTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("groupPolicyBootTimeInMs", m.GetGroupPolicyBootTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("groupPolicyLoginTimeInMs", m.GetGroupPolicyLoginTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFeatureUpdate", m.GetIsFeatureUpdate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFirstLogin", m.GetIsFirstLogin())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("operatingSystemVersion", m.GetOperatingSystemVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("responsiveDesktopTimeInMs", m.GetResponsiveDesktopTimeInMs())
        if err != nil {
            return err
        }
    }
    if m.GetRestartCategory() != nil {
        cast := (*m.GetRestartCategory()).String()
        err = writer.WriteStringValue("restartCategory", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("restartFaultBucket", m.GetRestartFaultBucket())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("restartStopCode", m.GetRestartStopCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startTime", m.GetStartTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalBootTimeInMs", m.GetTotalBootTimeInMs())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalLoginTimeInMs", m.GetTotalLoginTimeInMs())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCoreBootTimeInMs sets the coreBootTimeInMs property value. The device core boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetCoreBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("coreBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetCoreLoginTimeInMs sets the coreLoginTimeInMs property value. The device core login time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetCoreLoginTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("coreLoginTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceId sets the deviceId property value. The Intune device id of the device. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetDeviceId(value *string)() {
    err := m.GetBackingStore().Set("deviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetFeatureUpdateBootTimeInMs sets the featureUpdateBootTimeInMs property value. The impact of device feature updates on boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetFeatureUpdateBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("featureUpdateBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupPolicyBootTimeInMs sets the groupPolicyBootTimeInMs property value. The impact of device group policy client on boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetGroupPolicyBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("groupPolicyBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupPolicyLoginTimeInMs sets the groupPolicyLoginTimeInMs property value. The impact of device group policy client on login time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetGroupPolicyLoginTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("groupPolicyLoginTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFeatureUpdate sets the isFeatureUpdate property value. When TRUE, indicates the device boot record is associated with feature updates. When FALSE, indicates the device boot record is not associated with feature updates. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetIsFeatureUpdate(value *bool)() {
    err := m.GetBackingStore().Set("isFeatureUpdate", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFirstLogin sets the isFirstLogin property value. When TRUE, indicates the device login is the first login after a reboot. When FALSE, indicates the device login is not the first login after a reboot. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetIsFirstLogin(value *bool)() {
    err := m.GetBackingStore().Set("isFirstLogin", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystemVersion sets the operatingSystemVersion property value. The user experience analytics device boot record's operating system version. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetOperatingSystemVersion(value *string)() {
    err := m.GetBackingStore().Set("operatingSystemVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetResponsiveDesktopTimeInMs sets the responsiveDesktopTimeInMs property value. The time for desktop to become responsive during login process in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetResponsiveDesktopTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("responsiveDesktopTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetRestartCategory sets the restartCategory property value. Operating System restart category.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetRestartCategory(value *UserExperienceAnalyticsOperatingSystemRestartCategory)() {
    err := m.GetBackingStore().Set("restartCategory", value)
    if err != nil {
        panic(err)
    }
}
// SetRestartFaultBucket sets the restartFaultBucket property value. OS restart fault bucket. The fault bucket is used to find additional information about a system crash. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetRestartFaultBucket(value *string)() {
    err := m.GetBackingStore().Set("restartFaultBucket", value)
    if err != nil {
        panic(err)
    }
}
// SetRestartStopCode sets the restartStopCode property value. OS restart stop code. This shows the bug check code which can be used to look up the blue screen reason. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetRestartStopCode(value *string)() {
    err := m.GetBackingStore().Set("restartStopCode", value)
    if err != nil {
        panic(err)
    }
}
// SetStartTime sets the startTime property value. The device boot start time. The value cannot be modified and is automatically populated when the device performs a reboot. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2022 would look like this: '2022-01-01T00:00:00Z'. Returned by default. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetStartTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalBootTimeInMs sets the totalBootTimeInMs property value. The device total boot time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetTotalBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("totalBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLoginTimeInMs sets the totalLoginTimeInMs property value. The device total login time in milliseconds. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsDeviceStartupHistory) SetTotalLoginTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("totalLoginTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsDeviceStartupHistoryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCoreBootTimeInMs()(*int32)
    GetCoreLoginTimeInMs()(*int32)
    GetDeviceId()(*string)
    GetFeatureUpdateBootTimeInMs()(*int32)
    GetGroupPolicyBootTimeInMs()(*int32)
    GetGroupPolicyLoginTimeInMs()(*int32)
    GetIsFeatureUpdate()(*bool)
    GetIsFirstLogin()(*bool)
    GetOperatingSystemVersion()(*string)
    GetResponsiveDesktopTimeInMs()(*int32)
    GetRestartCategory()(*UserExperienceAnalyticsOperatingSystemRestartCategory)
    GetRestartFaultBucket()(*string)
    GetRestartStopCode()(*string)
    GetStartTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTotalBootTimeInMs()(*int32)
    GetTotalLoginTimeInMs()(*int32)
    SetCoreBootTimeInMs(value *int32)()
    SetCoreLoginTimeInMs(value *int32)()
    SetDeviceId(value *string)()
    SetFeatureUpdateBootTimeInMs(value *int32)()
    SetGroupPolicyBootTimeInMs(value *int32)()
    SetGroupPolicyLoginTimeInMs(value *int32)()
    SetIsFeatureUpdate(value *bool)()
    SetIsFirstLogin(value *bool)()
    SetOperatingSystemVersion(value *string)()
    SetResponsiveDesktopTimeInMs(value *int32)()
    SetRestartCategory(value *UserExperienceAnalyticsOperatingSystemRestartCategory)()
    SetRestartFaultBucket(value *string)()
    SetRestartStopCode(value *string)()
    SetStartTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTotalBootTimeInMs(value *int32)()
    SetTotalLoginTimeInMs(value *int32)()
}

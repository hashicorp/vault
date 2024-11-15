package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails the user experience analytics application performance entity contains application performance by application version details.
type UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails struct {
    Entity
}
// NewUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails instantiates a new UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails and sets the default values.
func NewUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails()(*UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) {
    m := &UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails(), nil
}
// GetAppCrashCount gets the appCrashCount property value. The number of crashes for the app. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetAppCrashCount()(*int32) {
    val, err := m.GetBackingStore().Get("appCrashCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAppDisplayName gets the appDisplayName property value. The friendly name of the application.
// returns a *string when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetAppDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("appDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppName gets the appName property value. The name of the application.
// returns a *string when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetAppName()(*string) {
    val, err := m.GetBackingStore().Get("appName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppPublisher gets the appPublisher property value. The publisher of the application.
// returns a *string when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetAppPublisher()(*string) {
    val, err := m.GetBackingStore().Get("appPublisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppVersion gets the appVersion property value. The version of the application.
// returns a *string when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetAppVersion()(*string) {
    val, err := m.GetBackingStore().Get("appVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceCountWithCrashes gets the deviceCountWithCrashes property value. The total number of devices that have reported one or more application crashes for this application and version. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetDeviceCountWithCrashes()(*int32) {
    val, err := m.GetBackingStore().Get("deviceCountWithCrashes")
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
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appCrashCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppCrashCount(val)
        }
        return nil
    }
    res["appDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppDisplayName(val)
        }
        return nil
    }
    res["appName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppName(val)
        }
        return nil
    }
    res["appPublisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppPublisher(val)
        }
        return nil
    }
    res["appVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppVersion(val)
        }
        return nil
    }
    res["deviceCountWithCrashes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCountWithCrashes(val)
        }
        return nil
    }
    res["isLatestUsedVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLatestUsedVersion(val)
        }
        return nil
    }
    res["isMostUsedVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMostUsedVersion(val)
        }
        return nil
    }
    return res
}
// GetIsLatestUsedVersion gets the isLatestUsedVersion property value. When TRUE, indicates the version of application is the latest version for that application that is in use. When FALSE, indicates the version is not the latest version. FALSE by default. Supports: $select, $OrderBy.
// returns a *bool when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetIsLatestUsedVersion()(*bool) {
    val, err := m.GetBackingStore().Get("isLatestUsedVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMostUsedVersion gets the isMostUsedVersion property value. When TRUE, indicates the version of application is the most used version for that application. When FALSE, indicates the version is not the most used version. FALSE by default. Supports: $select, $OrderBy. Read-only.
// returns a *bool when successful
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) GetIsMostUsedVersion()(*bool) {
    val, err := m.GetBackingStore().Get("isMostUsedVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("appCrashCount", m.GetAppCrashCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appDisplayName", m.GetAppDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appName", m.GetAppName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appPublisher", m.GetAppPublisher())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appVersion", m.GetAppVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("deviceCountWithCrashes", m.GetDeviceCountWithCrashes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isLatestUsedVersion", m.GetIsLatestUsedVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isMostUsedVersion", m.GetIsMostUsedVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppCrashCount sets the appCrashCount property value. The number of crashes for the app. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetAppCrashCount(value *int32)() {
    err := m.GetBackingStore().Set("appCrashCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAppDisplayName sets the appDisplayName property value. The friendly name of the application.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetAppDisplayName(value *string)() {
    err := m.GetBackingStore().Set("appDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAppName sets the appName property value. The name of the application.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetAppName(value *string)() {
    err := m.GetBackingStore().Set("appName", value)
    if err != nil {
        panic(err)
    }
}
// SetAppPublisher sets the appPublisher property value. The publisher of the application.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetAppPublisher(value *string)() {
    err := m.GetBackingStore().Set("appPublisher", value)
    if err != nil {
        panic(err)
    }
}
// SetAppVersion sets the appVersion property value. The version of the application.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetAppVersion(value *string)() {
    err := m.GetBackingStore().Set("appVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCountWithCrashes sets the deviceCountWithCrashes property value. The total number of devices that have reported one or more application crashes for this application and version. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetDeviceCountWithCrashes(value *int32)() {
    err := m.GetBackingStore().Set("deviceCountWithCrashes", value)
    if err != nil {
        panic(err)
    }
}
// SetIsLatestUsedVersion sets the isLatestUsedVersion property value. When TRUE, indicates the version of application is the latest version for that application that is in use. When FALSE, indicates the version is not the latest version. FALSE by default. Supports: $select, $OrderBy.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetIsLatestUsedVersion(value *bool)() {
    err := m.GetBackingStore().Set("isLatestUsedVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMostUsedVersion sets the isMostUsedVersion property value. When TRUE, indicates the version of application is the most used version for that application. When FALSE, indicates the version is not the most used version. FALSE by default. Supports: $select, $OrderBy. Read-only.
func (m *UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails) SetIsMostUsedVersion(value *bool)() {
    err := m.GetBackingStore().Set("isMostUsedVersion", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetailsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppCrashCount()(*int32)
    GetAppDisplayName()(*string)
    GetAppName()(*string)
    GetAppPublisher()(*string)
    GetAppVersion()(*string)
    GetDeviceCountWithCrashes()(*int32)
    GetIsLatestUsedVersion()(*bool)
    GetIsMostUsedVersion()(*bool)
    SetAppCrashCount(value *int32)()
    SetAppDisplayName(value *string)()
    SetAppName(value *string)()
    SetAppPublisher(value *string)()
    SetAppVersion(value *string)()
    SetDeviceCountWithCrashes(value *int32)()
    SetIsLatestUsedVersion(value *bool)()
    SetIsMostUsedVersion(value *bool)()
}

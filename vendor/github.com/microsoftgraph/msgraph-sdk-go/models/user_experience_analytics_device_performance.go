package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsDevicePerformance the user experience analytics device performance entity contains device boot performance details.
type UserExperienceAnalyticsDevicePerformance struct {
    Entity
}
// NewUserExperienceAnalyticsDevicePerformance instantiates a new UserExperienceAnalyticsDevicePerformance and sets the default values.
func NewUserExperienceAnalyticsDevicePerformance()(*UserExperienceAnalyticsDevicePerformance) {
    m := &UserExperienceAnalyticsDevicePerformance{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsDevicePerformanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsDevicePerformanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsDevicePerformance(), nil
}
// GetAverageBlueScreens gets the averageBlueScreens property value. Average (mean) number of Blue Screens per device in the last 30 days. Valid values 0 to 9999999
// returns a *float64 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetAverageBlueScreens()(*float64) {
    val, err := m.GetBackingStore().Get("averageBlueScreens")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetAverageRestarts gets the averageRestarts property value. Average (mean) number of Restarts per device in the last 30 days. Valid values 0 to 9999999
// returns a *float64 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetAverageRestarts()(*float64) {
    val, err := m.GetBackingStore().Get("averageRestarts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetBlueScreenCount gets the blueScreenCount property value. Number of Blue Screens in the last 30 days. Valid values 0 to 9999999
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetBlueScreenCount()(*int32) {
    val, err := m.GetBackingStore().Get("blueScreenCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBootScore gets the bootScore property value. The user experience analytics device boot score.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetBootScore()(*int32) {
    val, err := m.GetBackingStore().Get("bootScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCoreBootTimeInMs gets the coreBootTimeInMs property value. The user experience analytics device core boot time in milliseconds.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetCoreBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("coreBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCoreLoginTimeInMs gets the coreLoginTimeInMs property value. The user experience analytics device core login time in milliseconds.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetCoreLoginTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("coreLoginTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeviceCount gets the deviceCount property value. User experience analytics summarized device count.
// returns a *int64 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetDeviceCount()(*int64) {
    val, err := m.GetBackingStore().Get("deviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. The user experience analytics device name.
// returns a *string when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDiskType gets the diskType property value. The diskType property
// returns a *DiskType when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetDiskType()(*DiskType) {
    val, err := m.GetBackingStore().Get("diskType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DiskType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["averageBlueScreens"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageBlueScreens(val)
        }
        return nil
    }
    res["averageRestarts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAverageRestarts(val)
        }
        return nil
    }
    res["blueScreenCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlueScreenCount(val)
        }
        return nil
    }
    res["bootScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootScore(val)
        }
        return nil
    }
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
    res["deviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCount(val)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
        }
        return nil
    }
    res["diskType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDiskType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiskType(val.(*DiskType))
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
    res["healthStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserExperienceAnalyticsHealthState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHealthStatus(val.(*UserExperienceAnalyticsHealthState))
        }
        return nil
    }
    res["loginScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoginScore(val)
        }
        return nil
    }
    res["manufacturer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManufacturer(val)
        }
        return nil
    }
    res["model"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModel(val)
        }
        return nil
    }
    res["modelStartupPerformanceScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModelStartupPerformanceScore(val)
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
    res["restartCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestartCount(val)
        }
        return nil
    }
    res["startupPerformanceScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartupPerformanceScore(val)
        }
        return nil
    }
    return res
}
// GetGroupPolicyBootTimeInMs gets the groupPolicyBootTimeInMs property value. The user experience analytics device group policy boot time in milliseconds.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetGroupPolicyBootTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("groupPolicyBootTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetGroupPolicyLoginTimeInMs gets the groupPolicyLoginTimeInMs property value. The user experience analytics device group policy login time in milliseconds.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetGroupPolicyLoginTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("groupPolicyLoginTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetHealthStatus gets the healthStatus property value. The healthStatus property
// returns a *UserExperienceAnalyticsHealthState when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetHealthStatus()(*UserExperienceAnalyticsHealthState) {
    val, err := m.GetBackingStore().Get("healthStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserExperienceAnalyticsHealthState)
    }
    return nil
}
// GetLoginScore gets the loginScore property value. The user experience analytics device login score.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetLoginScore()(*int32) {
    val, err := m.GetBackingStore().Get("loginScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetManufacturer gets the manufacturer property value. The user experience analytics device manufacturer.
// returns a *string when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetManufacturer()(*string) {
    val, err := m.GetBackingStore().Get("manufacturer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModel gets the model property value. The user experience analytics device model.
// returns a *string when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModelStartupPerformanceScore gets the modelStartupPerformanceScore property value. The user experience analytics model level startup performance score. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetModelStartupPerformanceScore()(*float64) {
    val, err := m.GetBackingStore().Get("modelStartupPerformanceScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetOperatingSystemVersion gets the operatingSystemVersion property value. The user experience analytics device Operating System version.
// returns a *string when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetOperatingSystemVersion()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystemVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResponsiveDesktopTimeInMs gets the responsiveDesktopTimeInMs property value. The user experience analytics responsive desktop time in milliseconds.
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetResponsiveDesktopTimeInMs()(*int32) {
    val, err := m.GetBackingStore().Get("responsiveDesktopTimeInMs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRestartCount gets the restartCount property value. Number of Restarts in the last 30 days. Valid values 0 to 9999999
// returns a *int32 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetRestartCount()(*int32) {
    val, err := m.GetBackingStore().Get("restartCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetStartupPerformanceScore gets the startupPerformanceScore property value. The user experience analytics device startup performance score. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsDevicePerformance) GetStartupPerformanceScore()(*float64) {
    val, err := m.GetBackingStore().Get("startupPerformanceScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsDevicePerformance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteFloat64Value("averageBlueScreens", m.GetAverageBlueScreens())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("averageRestarts", m.GetAverageRestarts())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("blueScreenCount", m.GetBlueScreenCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("bootScore", m.GetBootScore())
        if err != nil {
            return err
        }
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
        err = writer.WriteInt64Value("deviceCount", m.GetDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceName", m.GetDeviceName())
        if err != nil {
            return err
        }
    }
    if m.GetDiskType() != nil {
        cast := (*m.GetDiskType()).String()
        err = writer.WriteStringValue("diskType", &cast)
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
    if m.GetHealthStatus() != nil {
        cast := (*m.GetHealthStatus()).String()
        err = writer.WriteStringValue("healthStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("loginScore", m.GetLoginScore())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("modelStartupPerformanceScore", m.GetModelStartupPerformanceScore())
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
    {
        err = writer.WriteInt32Value("restartCount", m.GetRestartCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("startupPerformanceScore", m.GetStartupPerformanceScore())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAverageBlueScreens sets the averageBlueScreens property value. Average (mean) number of Blue Screens per device in the last 30 days. Valid values 0 to 9999999
func (m *UserExperienceAnalyticsDevicePerformance) SetAverageBlueScreens(value *float64)() {
    err := m.GetBackingStore().Set("averageBlueScreens", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageRestarts sets the averageRestarts property value. Average (mean) number of Restarts per device in the last 30 days. Valid values 0 to 9999999
func (m *UserExperienceAnalyticsDevicePerformance) SetAverageRestarts(value *float64)() {
    err := m.GetBackingStore().Set("averageRestarts", value)
    if err != nil {
        panic(err)
    }
}
// SetBlueScreenCount sets the blueScreenCount property value. Number of Blue Screens in the last 30 days. Valid values 0 to 9999999
func (m *UserExperienceAnalyticsDevicePerformance) SetBlueScreenCount(value *int32)() {
    err := m.GetBackingStore().Set("blueScreenCount", value)
    if err != nil {
        panic(err)
    }
}
// SetBootScore sets the bootScore property value. The user experience analytics device boot score.
func (m *UserExperienceAnalyticsDevicePerformance) SetBootScore(value *int32)() {
    err := m.GetBackingStore().Set("bootScore", value)
    if err != nil {
        panic(err)
    }
}
// SetCoreBootTimeInMs sets the coreBootTimeInMs property value. The user experience analytics device core boot time in milliseconds.
func (m *UserExperienceAnalyticsDevicePerformance) SetCoreBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("coreBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetCoreLoginTimeInMs sets the coreLoginTimeInMs property value. The user experience analytics device core login time in milliseconds.
func (m *UserExperienceAnalyticsDevicePerformance) SetCoreLoginTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("coreLoginTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCount sets the deviceCount property value. User experience analytics summarized device count.
func (m *UserExperienceAnalyticsDevicePerformance) SetDeviceCount(value *int64)() {
    err := m.GetBackingStore().Set("deviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. The user experience analytics device name.
func (m *UserExperienceAnalyticsDevicePerformance) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetDiskType sets the diskType property value. The diskType property
func (m *UserExperienceAnalyticsDevicePerformance) SetDiskType(value *DiskType)() {
    err := m.GetBackingStore().Set("diskType", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupPolicyBootTimeInMs sets the groupPolicyBootTimeInMs property value. The user experience analytics device group policy boot time in milliseconds.
func (m *UserExperienceAnalyticsDevicePerformance) SetGroupPolicyBootTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("groupPolicyBootTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupPolicyLoginTimeInMs sets the groupPolicyLoginTimeInMs property value. The user experience analytics device group policy login time in milliseconds.
func (m *UserExperienceAnalyticsDevicePerformance) SetGroupPolicyLoginTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("groupPolicyLoginTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetHealthStatus sets the healthStatus property value. The healthStatus property
func (m *UserExperienceAnalyticsDevicePerformance) SetHealthStatus(value *UserExperienceAnalyticsHealthState)() {
    err := m.GetBackingStore().Set("healthStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetLoginScore sets the loginScore property value. The user experience analytics device login score.
func (m *UserExperienceAnalyticsDevicePerformance) SetLoginScore(value *int32)() {
    err := m.GetBackingStore().Set("loginScore", value)
    if err != nil {
        panic(err)
    }
}
// SetManufacturer sets the manufacturer property value. The user experience analytics device manufacturer.
func (m *UserExperienceAnalyticsDevicePerformance) SetManufacturer(value *string)() {
    err := m.GetBackingStore().Set("manufacturer", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. The user experience analytics device model.
func (m *UserExperienceAnalyticsDevicePerformance) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
// SetModelStartupPerformanceScore sets the modelStartupPerformanceScore property value. The user experience analytics model level startup performance score. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsDevicePerformance) SetModelStartupPerformanceScore(value *float64)() {
    err := m.GetBackingStore().Set("modelStartupPerformanceScore", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystemVersion sets the operatingSystemVersion property value. The user experience analytics device Operating System version.
func (m *UserExperienceAnalyticsDevicePerformance) SetOperatingSystemVersion(value *string)() {
    err := m.GetBackingStore().Set("operatingSystemVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetResponsiveDesktopTimeInMs sets the responsiveDesktopTimeInMs property value. The user experience analytics responsive desktop time in milliseconds.
func (m *UserExperienceAnalyticsDevicePerformance) SetResponsiveDesktopTimeInMs(value *int32)() {
    err := m.GetBackingStore().Set("responsiveDesktopTimeInMs", value)
    if err != nil {
        panic(err)
    }
}
// SetRestartCount sets the restartCount property value. Number of Restarts in the last 30 days. Valid values 0 to 9999999
func (m *UserExperienceAnalyticsDevicePerformance) SetRestartCount(value *int32)() {
    err := m.GetBackingStore().Set("restartCount", value)
    if err != nil {
        panic(err)
    }
}
// SetStartupPerformanceScore sets the startupPerformanceScore property value. The user experience analytics device startup performance score. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsDevicePerformance) SetStartupPerformanceScore(value *float64)() {
    err := m.GetBackingStore().Set("startupPerformanceScore", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsDevicePerformanceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAverageBlueScreens()(*float64)
    GetAverageRestarts()(*float64)
    GetBlueScreenCount()(*int32)
    GetBootScore()(*int32)
    GetCoreBootTimeInMs()(*int32)
    GetCoreLoginTimeInMs()(*int32)
    GetDeviceCount()(*int64)
    GetDeviceName()(*string)
    GetDiskType()(*DiskType)
    GetGroupPolicyBootTimeInMs()(*int32)
    GetGroupPolicyLoginTimeInMs()(*int32)
    GetHealthStatus()(*UserExperienceAnalyticsHealthState)
    GetLoginScore()(*int32)
    GetManufacturer()(*string)
    GetModel()(*string)
    GetModelStartupPerformanceScore()(*float64)
    GetOperatingSystemVersion()(*string)
    GetResponsiveDesktopTimeInMs()(*int32)
    GetRestartCount()(*int32)
    GetStartupPerformanceScore()(*float64)
    SetAverageBlueScreens(value *float64)()
    SetAverageRestarts(value *float64)()
    SetBlueScreenCount(value *int32)()
    SetBootScore(value *int32)()
    SetCoreBootTimeInMs(value *int32)()
    SetCoreLoginTimeInMs(value *int32)()
    SetDeviceCount(value *int64)()
    SetDeviceName(value *string)()
    SetDiskType(value *DiskType)()
    SetGroupPolicyBootTimeInMs(value *int32)()
    SetGroupPolicyLoginTimeInMs(value *int32)()
    SetHealthStatus(value *UserExperienceAnalyticsHealthState)()
    SetLoginScore(value *int32)()
    SetManufacturer(value *string)()
    SetModel(value *string)()
    SetModelStartupPerformanceScore(value *float64)()
    SetOperatingSystemVersion(value *string)()
    SetResponsiveDesktopTimeInMs(value *int32)()
    SetRestartCount(value *int32)()
    SetStartupPerformanceScore(value *float64)()
}

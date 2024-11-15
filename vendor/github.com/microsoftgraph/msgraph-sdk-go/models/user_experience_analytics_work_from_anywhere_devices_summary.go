package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// UserExperienceAnalyticsWorkFromAnywhereDevicesSummary the user experience analytics Work From Anywhere metrics devices summary.
type UserExperienceAnalyticsWorkFromAnywhereDevicesSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserExperienceAnalyticsWorkFromAnywhereDevicesSummary instantiates a new UserExperienceAnalyticsWorkFromAnywhereDevicesSummary and sets the default values.
func NewUserExperienceAnalyticsWorkFromAnywhereDevicesSummary()(*UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) {
    m := &UserExperienceAnalyticsWorkFromAnywhereDevicesSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserExperienceAnalyticsWorkFromAnywhereDevicesSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsWorkFromAnywhereDevicesSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsWorkFromAnywhereDevicesSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAutopilotDevicesSummary gets the autopilotDevicesSummary property value. The user experience analytics work from anywhere Autopilot devices summary. Read-only.
// returns a UserExperienceAnalyticsAutopilotDevicesSummaryable when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetAutopilotDevicesSummary()(UserExperienceAnalyticsAutopilotDevicesSummaryable) {
    val, err := m.GetBackingStore().Get("autopilotDevicesSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserExperienceAnalyticsAutopilotDevicesSummaryable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCloudIdentityDevicesSummary gets the cloudIdentityDevicesSummary property value. The user experience analytics work from anywhere Cloud Identity devices summary. Read-only.
// returns a UserExperienceAnalyticsCloudIdentityDevicesSummaryable when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetCloudIdentityDevicesSummary()(UserExperienceAnalyticsCloudIdentityDevicesSummaryable) {
    val, err := m.GetBackingStore().Get("cloudIdentityDevicesSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserExperienceAnalyticsCloudIdentityDevicesSummaryable)
    }
    return nil
}
// GetCloudManagementDevicesSummary gets the cloudManagementDevicesSummary property value. The user experience analytics work from anywhere Cloud management devices summary. Read-only.
// returns a UserExperienceAnalyticsCloudManagementDevicesSummaryable when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetCloudManagementDevicesSummary()(UserExperienceAnalyticsCloudManagementDevicesSummaryable) {
    val, err := m.GetBackingStore().Get("cloudManagementDevicesSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserExperienceAnalyticsCloudManagementDevicesSummaryable)
    }
    return nil
}
// GetCoManagedDevices gets the coManagedDevices property value. Total number of co-managed devices. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetCoManagedDevices()(*int32) {
    val, err := m.GetBackingStore().Get("coManagedDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDevicesNotAutopilotRegistered gets the devicesNotAutopilotRegistered property value. The count of intune devices that are not autopilot registerd. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetDevicesNotAutopilotRegistered()(*int32) {
    val, err := m.GetBackingStore().Get("devicesNotAutopilotRegistered")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDevicesWithoutAutopilotProfileAssigned gets the devicesWithoutAutopilotProfileAssigned property value. The count of intune devices not autopilot profile assigned. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetDevicesWithoutAutopilotProfileAssigned()(*int32) {
    val, err := m.GetBackingStore().Get("devicesWithoutAutopilotProfileAssigned")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDevicesWithoutCloudIdentity gets the devicesWithoutCloudIdentity property value. The count of devices that are not cloud identity. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetDevicesWithoutCloudIdentity()(*int32) {
    val, err := m.GetBackingStore().Get("devicesWithoutCloudIdentity")
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
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["autopilotDevicesSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserExperienceAnalyticsAutopilotDevicesSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutopilotDevicesSummary(val.(UserExperienceAnalyticsAutopilotDevicesSummaryable))
        }
        return nil
    }
    res["cloudIdentityDevicesSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserExperienceAnalyticsCloudIdentityDevicesSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudIdentityDevicesSummary(val.(UserExperienceAnalyticsCloudIdentityDevicesSummaryable))
        }
        return nil
    }
    res["cloudManagementDevicesSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserExperienceAnalyticsCloudManagementDevicesSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudManagementDevicesSummary(val.(UserExperienceAnalyticsCloudManagementDevicesSummaryable))
        }
        return nil
    }
    res["coManagedDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCoManagedDevices(val)
        }
        return nil
    }
    res["devicesNotAutopilotRegistered"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicesNotAutopilotRegistered(val)
        }
        return nil
    }
    res["devicesWithoutAutopilotProfileAssigned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicesWithoutAutopilotProfileAssigned(val)
        }
        return nil
    }
    res["devicesWithoutCloudIdentity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevicesWithoutCloudIdentity(val)
        }
        return nil
    }
    res["intuneDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIntuneDevices(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["tenantAttachDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantAttachDevices(val)
        }
        return nil
    }
    res["totalDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalDevices(val)
        }
        return nil
    }
    res["unsupportedOSversionDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnsupportedOSversionDevices(val)
        }
        return nil
    }
    res["windows10Devices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindows10Devices(val)
        }
        return nil
    }
    res["windows10DevicesSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserExperienceAnalyticsWindows10DevicesSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindows10DevicesSummary(val.(UserExperienceAnalyticsWindows10DevicesSummaryable))
        }
        return nil
    }
    res["windows10DevicesWithoutTenantAttach"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindows10DevicesWithoutTenantAttach(val)
        }
        return nil
    }
    return res
}
// GetIntuneDevices gets the intuneDevices property value. The count of intune devices that are not autopilot registerd. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetIntuneDevices()(*int32) {
    val, err := m.GetBackingStore().Get("intuneDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTenantAttachDevices gets the tenantAttachDevices property value. Total count of tenant attach devices. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetTenantAttachDevices()(*int32) {
    val, err := m.GetBackingStore().Get("tenantAttachDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalDevices gets the totalDevices property value. The total count of devices. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetTotalDevices()(*int32) {
    val, err := m.GetBackingStore().Get("totalDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnsupportedOSversionDevices gets the unsupportedOSversionDevices property value. The count of Windows 10 devices that have unsupported OS versions. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetUnsupportedOSversionDevices()(*int32) {
    val, err := m.GetBackingStore().Get("unsupportedOSversionDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWindows10Devices gets the windows10Devices property value. The count of windows 10 devices. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetWindows10Devices()(*int32) {
    val, err := m.GetBackingStore().Get("windows10Devices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWindows10DevicesSummary gets the windows10DevicesSummary property value. The user experience analytics work from anywhere Windows 10 devices summary. Read-only.
// returns a UserExperienceAnalyticsWindows10DevicesSummaryable when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetWindows10DevicesSummary()(UserExperienceAnalyticsWindows10DevicesSummaryable) {
    val, err := m.GetBackingStore().Get("windows10DevicesSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserExperienceAnalyticsWindows10DevicesSummaryable)
    }
    return nil
}
// GetWindows10DevicesWithoutTenantAttach gets the windows10DevicesWithoutTenantAttach property value. The count of windows 10 devices that are Intune and co-managed. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) GetWindows10DevicesWithoutTenantAttach()(*int32) {
    val, err := m.GetBackingStore().Get("windows10DevicesWithoutTenantAttach")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("autopilotDevicesSummary", m.GetAutopilotDevicesSummary())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("cloudIdentityDevicesSummary", m.GetCloudIdentityDevicesSummary())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("cloudManagementDevicesSummary", m.GetCloudManagementDevicesSummary())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("coManagedDevices", m.GetCoManagedDevices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("devicesNotAutopilotRegistered", m.GetDevicesNotAutopilotRegistered())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("devicesWithoutAutopilotProfileAssigned", m.GetDevicesWithoutAutopilotProfileAssigned())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("devicesWithoutCloudIdentity", m.GetDevicesWithoutCloudIdentity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("intuneDevices", m.GetIntuneDevices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("tenantAttachDevices", m.GetTenantAttachDevices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalDevices", m.GetTotalDevices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("unsupportedOSversionDevices", m.GetUnsupportedOSversionDevices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("windows10Devices", m.GetWindows10Devices())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("windows10DevicesSummary", m.GetWindows10DevicesSummary())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("windows10DevicesWithoutTenantAttach", m.GetWindows10DevicesWithoutTenantAttach())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAutopilotDevicesSummary sets the autopilotDevicesSummary property value. The user experience analytics work from anywhere Autopilot devices summary. Read-only.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetAutopilotDevicesSummary(value UserExperienceAnalyticsAutopilotDevicesSummaryable)() {
    err := m.GetBackingStore().Set("autopilotDevicesSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCloudIdentityDevicesSummary sets the cloudIdentityDevicesSummary property value. The user experience analytics work from anywhere Cloud Identity devices summary. Read-only.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetCloudIdentityDevicesSummary(value UserExperienceAnalyticsCloudIdentityDevicesSummaryable)() {
    err := m.GetBackingStore().Set("cloudIdentityDevicesSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudManagementDevicesSummary sets the cloudManagementDevicesSummary property value. The user experience analytics work from anywhere Cloud management devices summary. Read-only.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetCloudManagementDevicesSummary(value UserExperienceAnalyticsCloudManagementDevicesSummaryable)() {
    err := m.GetBackingStore().Set("cloudManagementDevicesSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetCoManagedDevices sets the coManagedDevices property value. Total number of co-managed devices. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetCoManagedDevices(value *int32)() {
    err := m.GetBackingStore().Set("coManagedDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetDevicesNotAutopilotRegistered sets the devicesNotAutopilotRegistered property value. The count of intune devices that are not autopilot registerd. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetDevicesNotAutopilotRegistered(value *int32)() {
    err := m.GetBackingStore().Set("devicesNotAutopilotRegistered", value)
    if err != nil {
        panic(err)
    }
}
// SetDevicesWithoutAutopilotProfileAssigned sets the devicesWithoutAutopilotProfileAssigned property value. The count of intune devices not autopilot profile assigned. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetDevicesWithoutAutopilotProfileAssigned(value *int32)() {
    err := m.GetBackingStore().Set("devicesWithoutAutopilotProfileAssigned", value)
    if err != nil {
        panic(err)
    }
}
// SetDevicesWithoutCloudIdentity sets the devicesWithoutCloudIdentity property value. The count of devices that are not cloud identity. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetDevicesWithoutCloudIdentity(value *int32)() {
    err := m.GetBackingStore().Set("devicesWithoutCloudIdentity", value)
    if err != nil {
        panic(err)
    }
}
// SetIntuneDevices sets the intuneDevices property value. The count of intune devices that are not autopilot registerd. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetIntuneDevices(value *int32)() {
    err := m.GetBackingStore().Set("intuneDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantAttachDevices sets the tenantAttachDevices property value. Total count of tenant attach devices. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetTenantAttachDevices(value *int32)() {
    err := m.GetBackingStore().Set("tenantAttachDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalDevices sets the totalDevices property value. The total count of devices. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetTotalDevices(value *int32)() {
    err := m.GetBackingStore().Set("totalDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetUnsupportedOSversionDevices sets the unsupportedOSversionDevices property value. The count of Windows 10 devices that have unsupported OS versions. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetUnsupportedOSversionDevices(value *int32)() {
    err := m.GetBackingStore().Set("unsupportedOSversionDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetWindows10Devices sets the windows10Devices property value. The count of windows 10 devices. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetWindows10Devices(value *int32)() {
    err := m.GetBackingStore().Set("windows10Devices", value)
    if err != nil {
        panic(err)
    }
}
// SetWindows10DevicesSummary sets the windows10DevicesSummary property value. The user experience analytics work from anywhere Windows 10 devices summary. Read-only.
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetWindows10DevicesSummary(value UserExperienceAnalyticsWindows10DevicesSummaryable)() {
    err := m.GetBackingStore().Set("windows10DevicesSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetWindows10DevicesWithoutTenantAttach sets the windows10DevicesWithoutTenantAttach property value. The count of windows 10 devices that are Intune and co-managed. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereDevicesSummary) SetWindows10DevicesWithoutTenantAttach(value *int32)() {
    err := m.GetBackingStore().Set("windows10DevicesWithoutTenantAttach", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsWorkFromAnywhereDevicesSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAutopilotDevicesSummary()(UserExperienceAnalyticsAutopilotDevicesSummaryable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCloudIdentityDevicesSummary()(UserExperienceAnalyticsCloudIdentityDevicesSummaryable)
    GetCloudManagementDevicesSummary()(UserExperienceAnalyticsCloudManagementDevicesSummaryable)
    GetCoManagedDevices()(*int32)
    GetDevicesNotAutopilotRegistered()(*int32)
    GetDevicesWithoutAutopilotProfileAssigned()(*int32)
    GetDevicesWithoutCloudIdentity()(*int32)
    GetIntuneDevices()(*int32)
    GetOdataType()(*string)
    GetTenantAttachDevices()(*int32)
    GetTotalDevices()(*int32)
    GetUnsupportedOSversionDevices()(*int32)
    GetWindows10Devices()(*int32)
    GetWindows10DevicesSummary()(UserExperienceAnalyticsWindows10DevicesSummaryable)
    GetWindows10DevicesWithoutTenantAttach()(*int32)
    SetAutopilotDevicesSummary(value UserExperienceAnalyticsAutopilotDevicesSummaryable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCloudIdentityDevicesSummary(value UserExperienceAnalyticsCloudIdentityDevicesSummaryable)()
    SetCloudManagementDevicesSummary(value UserExperienceAnalyticsCloudManagementDevicesSummaryable)()
    SetCoManagedDevices(value *int32)()
    SetDevicesNotAutopilotRegistered(value *int32)()
    SetDevicesWithoutAutopilotProfileAssigned(value *int32)()
    SetDevicesWithoutCloudIdentity(value *int32)()
    SetIntuneDevices(value *int32)()
    SetOdataType(value *string)()
    SetTenantAttachDevices(value *int32)()
    SetTotalDevices(value *int32)()
    SetUnsupportedOSversionDevices(value *int32)()
    SetWindows10Devices(value *int32)()
    SetWindows10DevicesSummary(value UserExperienceAnalyticsWindows10DevicesSummaryable)()
    SetWindows10DevicesWithoutTenantAttach(value *int32)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric the user experience analytics hardware readiness entity contains account level information about hardware blockers for windows upgrade.
type UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric struct {
    Entity
}
// NewUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric instantiates a new UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric and sets the default values.
func NewUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric()(*UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) {
    m := &UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["osCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsCheckFailedPercentage(val)
        }
        return nil
    }
    res["processor64BitCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessor64BitCheckFailedPercentage(val)
        }
        return nil
    }
    res["processorCoreCountCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessorCoreCountCheckFailedPercentage(val)
        }
        return nil
    }
    res["processorFamilyCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessorFamilyCheckFailedPercentage(val)
        }
        return nil
    }
    res["processorSpeedCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProcessorSpeedCheckFailedPercentage(val)
        }
        return nil
    }
    res["ramCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRamCheckFailedPercentage(val)
        }
        return nil
    }
    res["secureBootCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecureBootCheckFailedPercentage(val)
        }
        return nil
    }
    res["storageCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageCheckFailedPercentage(val)
        }
        return nil
    }
    res["totalDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalDeviceCount(val)
        }
        return nil
    }
    res["tpmCheckFailedPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTpmCheckFailedPercentage(val)
        }
        return nil
    }
    res["upgradeEligibleDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpgradeEligibleDeviceCount(val)
        }
        return nil
    }
    return res
}
// GetOsCheckFailedPercentage gets the osCheckFailedPercentage property value. The percentage of devices for which OS check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetOsCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("osCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetProcessor64BitCheckFailedPercentage gets the processor64BitCheckFailedPercentage property value. The percentage of devices for which processor hardware 64-bit architecture check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetProcessor64BitCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("processor64BitCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetProcessorCoreCountCheckFailedPercentage gets the processorCoreCountCheckFailedPercentage property value. The percentage of devices for which processor hardware core count check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetProcessorCoreCountCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("processorCoreCountCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetProcessorFamilyCheckFailedPercentage gets the processorFamilyCheckFailedPercentage property value. The percentage of devices for which processor hardware family check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetProcessorFamilyCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("processorFamilyCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetProcessorSpeedCheckFailedPercentage gets the processorSpeedCheckFailedPercentage property value. The percentage of devices for which processor hardware speed check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetProcessorSpeedCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("processorSpeedCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetRamCheckFailedPercentage gets the ramCheckFailedPercentage property value. The percentage of devices for which RAM hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetRamCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("ramCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetSecureBootCheckFailedPercentage gets the secureBootCheckFailedPercentage property value. The percentage of devices for which secure boot hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetSecureBootCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("secureBootCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetStorageCheckFailedPercentage gets the storageCheckFailedPercentage property value. The percentage of devices for which storage hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetStorageCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("storageCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetTotalDeviceCount gets the totalDeviceCount property value. The count of total devices in an organization. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetTotalDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTpmCheckFailedPercentage gets the tpmCheckFailedPercentage property value. The percentage of devices for which Trusted Platform Module (TPM) hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
// returns a *float64 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetTpmCheckFailedPercentage()(*float64) {
    val, err := m.GetBackingStore().Get("tpmCheckFailedPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetUpgradeEligibleDeviceCount gets the upgradeEligibleDeviceCount property value. The count of devices in an organization eligible for windows upgrade. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) GetUpgradeEligibleDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("upgradeEligibleDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteFloat64Value("osCheckFailedPercentage", m.GetOsCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("processor64BitCheckFailedPercentage", m.GetProcessor64BitCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("processorCoreCountCheckFailedPercentage", m.GetProcessorCoreCountCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("processorFamilyCheckFailedPercentage", m.GetProcessorFamilyCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("processorSpeedCheckFailedPercentage", m.GetProcessorSpeedCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("ramCheckFailedPercentage", m.GetRamCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("secureBootCheckFailedPercentage", m.GetSecureBootCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("storageCheckFailedPercentage", m.GetStorageCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalDeviceCount", m.GetTotalDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("tpmCheckFailedPercentage", m.GetTpmCheckFailedPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("upgradeEligibleDeviceCount", m.GetUpgradeEligibleDeviceCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOsCheckFailedPercentage sets the osCheckFailedPercentage property value. The percentage of devices for which OS check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetOsCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("osCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessor64BitCheckFailedPercentage sets the processor64BitCheckFailedPercentage property value. The percentage of devices for which processor hardware 64-bit architecture check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetProcessor64BitCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("processor64BitCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessorCoreCountCheckFailedPercentage sets the processorCoreCountCheckFailedPercentage property value. The percentage of devices for which processor hardware core count check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetProcessorCoreCountCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("processorCoreCountCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessorFamilyCheckFailedPercentage sets the processorFamilyCheckFailedPercentage property value. The percentage of devices for which processor hardware family check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetProcessorFamilyCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("processorFamilyCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetProcessorSpeedCheckFailedPercentage sets the processorSpeedCheckFailedPercentage property value. The percentage of devices for which processor hardware speed check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetProcessorSpeedCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("processorSpeedCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetRamCheckFailedPercentage sets the ramCheckFailedPercentage property value. The percentage of devices for which RAM hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetRamCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("ramCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetSecureBootCheckFailedPercentage sets the secureBootCheckFailedPercentage property value. The percentage of devices for which secure boot hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetSecureBootCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("secureBootCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageCheckFailedPercentage sets the storageCheckFailedPercentage property value. The percentage of devices for which storage hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetStorageCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("storageCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalDeviceCount sets the totalDeviceCount property value. The count of total devices in an organization. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetTotalDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("totalDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTpmCheckFailedPercentage sets the tpmCheckFailedPercentage property value. The percentage of devices for which Trusted Platform Module (TPM) hardware check has failed. Valid values 0 to 100. Supports: $select, $OrderBy. Read-only. Valid values -1.79769313486232E+308 to 1.79769313486232E+308
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetTpmCheckFailedPercentage(value *float64)() {
    err := m.GetBackingStore().Set("tpmCheckFailedPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetUpgradeEligibleDeviceCount sets the upgradeEligibleDeviceCount property value. The count of devices in an organization eligible for windows upgrade. Valid values 0 to 2147483647. Supports: $select, $OrderBy. Read-only. Valid values -2147483648 to 2147483647
func (m *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric) SetUpgradeEligibleDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("upgradeEligibleDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOsCheckFailedPercentage()(*float64)
    GetProcessor64BitCheckFailedPercentage()(*float64)
    GetProcessorCoreCountCheckFailedPercentage()(*float64)
    GetProcessorFamilyCheckFailedPercentage()(*float64)
    GetProcessorSpeedCheckFailedPercentage()(*float64)
    GetRamCheckFailedPercentage()(*float64)
    GetSecureBootCheckFailedPercentage()(*float64)
    GetStorageCheckFailedPercentage()(*float64)
    GetTotalDeviceCount()(*int32)
    GetTpmCheckFailedPercentage()(*float64)
    GetUpgradeEligibleDeviceCount()(*int32)
    SetOsCheckFailedPercentage(value *float64)()
    SetProcessor64BitCheckFailedPercentage(value *float64)()
    SetProcessorCoreCountCheckFailedPercentage(value *float64)()
    SetProcessorFamilyCheckFailedPercentage(value *float64)()
    SetProcessorSpeedCheckFailedPercentage(value *float64)()
    SetRamCheckFailedPercentage(value *float64)()
    SetSecureBootCheckFailedPercentage(value *float64)()
    SetStorageCheckFailedPercentage(value *float64)()
    SetTotalDeviceCount(value *int32)()
    SetTpmCheckFailedPercentage(value *float64)()
    SetUpgradeEligibleDeviceCount(value *int32)()
}

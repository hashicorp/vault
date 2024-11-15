package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// DeviceOperatingSystemSummary device operating system summary.
type DeviceOperatingSystemSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDeviceOperatingSystemSummary instantiates a new DeviceOperatingSystemSummary and sets the default values.
func NewDeviceOperatingSystemSummary()(*DeviceOperatingSystemSummary) {
    m := &DeviceOperatingSystemSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDeviceOperatingSystemSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceOperatingSystemSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceOperatingSystemSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DeviceOperatingSystemSummary) GetAdditionalData()(map[string]any) {
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
// GetAndroidCorporateWorkProfileCount gets the androidCorporateWorkProfileCount property value. The count of Corporate work profile Android devices. Also known as Corporate Owned Personally Enabled (COPE). Valid values -1 to 2147483647
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidCorporateWorkProfileCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidCorporateWorkProfileCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAndroidCount gets the androidCount property value. Number of android device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAndroidDedicatedCount gets the androidDedicatedCount property value. Number of dedicated Android devices.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidDedicatedCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidDedicatedCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAndroidDeviceAdminCount gets the androidDeviceAdminCount property value. Number of device admin Android devices.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidDeviceAdminCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidDeviceAdminCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAndroidFullyManagedCount gets the androidFullyManagedCount property value. Number of fully managed Android devices.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidFullyManagedCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidFullyManagedCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAndroidWorkProfileCount gets the androidWorkProfileCount property value. Number of work profile Android devices.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetAndroidWorkProfileCount()(*int32) {
    val, err := m.GetBackingStore().Get("androidWorkProfileCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DeviceOperatingSystemSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConfigMgrDeviceCount gets the configMgrDeviceCount property value. Number of ConfigMgr managed devices.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetConfigMgrDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("configMgrDeviceCount")
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
func (m *DeviceOperatingSystemSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["androidCorporateWorkProfileCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidCorporateWorkProfileCount(val)
        }
        return nil
    }
    res["androidCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidCount(val)
        }
        return nil
    }
    res["androidDedicatedCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidDedicatedCount(val)
        }
        return nil
    }
    res["androidDeviceAdminCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidDeviceAdminCount(val)
        }
        return nil
    }
    res["androidFullyManagedCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidFullyManagedCount(val)
        }
        return nil
    }
    res["androidWorkProfileCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidWorkProfileCount(val)
        }
        return nil
    }
    res["configMgrDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigMgrDeviceCount(val)
        }
        return nil
    }
    res["iosCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIosCount(val)
        }
        return nil
    }
    res["macOSCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMacOSCount(val)
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
    res["unknownCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnknownCount(val)
        }
        return nil
    }
    res["windowsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsCount(val)
        }
        return nil
    }
    res["windowsMobileCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsMobileCount(val)
        }
        return nil
    }
    return res
}
// GetIosCount gets the iosCount property value. Number of iOS device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetIosCount()(*int32) {
    val, err := m.GetBackingStore().Get("iosCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMacOSCount gets the macOSCount property value. Number of Mac OS X device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetMacOSCount()(*int32) {
    val, err := m.GetBackingStore().Get("macOSCount")
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
func (m *DeviceOperatingSystemSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUnknownCount gets the unknownCount property value. Number of unknown device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetUnknownCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWindowsCount gets the windowsCount property value. Number of Windows device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetWindowsCount()(*int32) {
    val, err := m.GetBackingStore().Get("windowsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWindowsMobileCount gets the windowsMobileCount property value. Number of Windows mobile device count.
// returns a *int32 when successful
func (m *DeviceOperatingSystemSummary) GetWindowsMobileCount()(*int32) {
    val, err := m.GetBackingStore().Get("windowsMobileCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceOperatingSystemSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("androidCorporateWorkProfileCount", m.GetAndroidCorporateWorkProfileCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("androidCount", m.GetAndroidCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("androidDedicatedCount", m.GetAndroidDedicatedCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("androidDeviceAdminCount", m.GetAndroidDeviceAdminCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("androidFullyManagedCount", m.GetAndroidFullyManagedCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("androidWorkProfileCount", m.GetAndroidWorkProfileCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("configMgrDeviceCount", m.GetConfigMgrDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("iosCount", m.GetIosCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("macOSCount", m.GetMacOSCount())
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
        err := writer.WriteInt32Value("unknownCount", m.GetUnknownCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("windowsCount", m.GetWindowsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("windowsMobileCount", m.GetWindowsMobileCount())
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
func (m *DeviceOperatingSystemSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidCorporateWorkProfileCount sets the androidCorporateWorkProfileCount property value. The count of Corporate work profile Android devices. Also known as Corporate Owned Personally Enabled (COPE). Valid values -1 to 2147483647
func (m *DeviceOperatingSystemSummary) SetAndroidCorporateWorkProfileCount(value *int32)() {
    err := m.GetBackingStore().Set("androidCorporateWorkProfileCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidCount sets the androidCount property value. Number of android device count.
func (m *DeviceOperatingSystemSummary) SetAndroidCount(value *int32)() {
    err := m.GetBackingStore().Set("androidCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidDedicatedCount sets the androidDedicatedCount property value. Number of dedicated Android devices.
func (m *DeviceOperatingSystemSummary) SetAndroidDedicatedCount(value *int32)() {
    err := m.GetBackingStore().Set("androidDedicatedCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidDeviceAdminCount sets the androidDeviceAdminCount property value. Number of device admin Android devices.
func (m *DeviceOperatingSystemSummary) SetAndroidDeviceAdminCount(value *int32)() {
    err := m.GetBackingStore().Set("androidDeviceAdminCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidFullyManagedCount sets the androidFullyManagedCount property value. Number of fully managed Android devices.
func (m *DeviceOperatingSystemSummary) SetAndroidFullyManagedCount(value *int32)() {
    err := m.GetBackingStore().Set("androidFullyManagedCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidWorkProfileCount sets the androidWorkProfileCount property value. Number of work profile Android devices.
func (m *DeviceOperatingSystemSummary) SetAndroidWorkProfileCount(value *int32)() {
    err := m.GetBackingStore().Set("androidWorkProfileCount", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DeviceOperatingSystemSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConfigMgrDeviceCount sets the configMgrDeviceCount property value. Number of ConfigMgr managed devices.
func (m *DeviceOperatingSystemSummary) SetConfigMgrDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("configMgrDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIosCount sets the iosCount property value. Number of iOS device count.
func (m *DeviceOperatingSystemSummary) SetIosCount(value *int32)() {
    err := m.GetBackingStore().Set("iosCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMacOSCount sets the macOSCount property value. Number of Mac OS X device count.
func (m *DeviceOperatingSystemSummary) SetMacOSCount(value *int32)() {
    err := m.GetBackingStore().Set("macOSCount", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DeviceOperatingSystemSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownCount sets the unknownCount property value. Number of unknown device count.
func (m *DeviceOperatingSystemSummary) SetUnknownCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownCount", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsCount sets the windowsCount property value. Number of Windows device count.
func (m *DeviceOperatingSystemSummary) SetWindowsCount(value *int32)() {
    err := m.GetBackingStore().Set("windowsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsMobileCount sets the windowsMobileCount property value. Number of Windows mobile device count.
func (m *DeviceOperatingSystemSummary) SetWindowsMobileCount(value *int32)() {
    err := m.GetBackingStore().Set("windowsMobileCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceOperatingSystemSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAndroidCorporateWorkProfileCount()(*int32)
    GetAndroidCount()(*int32)
    GetAndroidDedicatedCount()(*int32)
    GetAndroidDeviceAdminCount()(*int32)
    GetAndroidFullyManagedCount()(*int32)
    GetAndroidWorkProfileCount()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConfigMgrDeviceCount()(*int32)
    GetIosCount()(*int32)
    GetMacOSCount()(*int32)
    GetOdataType()(*string)
    GetUnknownCount()(*int32)
    GetWindowsCount()(*int32)
    GetWindowsMobileCount()(*int32)
    SetAndroidCorporateWorkProfileCount(value *int32)()
    SetAndroidCount(value *int32)()
    SetAndroidDedicatedCount(value *int32)()
    SetAndroidDeviceAdminCount(value *int32)()
    SetAndroidFullyManagedCount(value *int32)()
    SetAndroidWorkProfileCount(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConfigMgrDeviceCount(value *int32)()
    SetIosCount(value *int32)()
    SetMacOSCount(value *int32)()
    SetOdataType(value *string)()
    SetUnknownCount(value *int32)()
    SetWindowsCount(value *int32)()
    SetWindowsMobileCount(value *int32)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// DeviceProtectionOverview hardware information of a given device.
type DeviceProtectionOverview struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDeviceProtectionOverview instantiates a new DeviceProtectionOverview and sets the default values.
func NewDeviceProtectionOverview()(*DeviceProtectionOverview) {
    m := &DeviceProtectionOverview{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDeviceProtectionOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceProtectionOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceProtectionOverview(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DeviceProtectionOverview) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DeviceProtectionOverview) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCleanDeviceCount gets the cleanDeviceCount property value. Indicates number of devices reporting as clean
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetCleanDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("cleanDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCriticalFailuresDeviceCount gets the criticalFailuresDeviceCount property value. Indicates number of devices with critical failures
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetCriticalFailuresDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("criticalFailuresDeviceCount")
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
func (m *DeviceProtectionOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["cleanDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCleanDeviceCount(val)
        }
        return nil
    }
    res["criticalFailuresDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCriticalFailuresDeviceCount(val)
        }
        return nil
    }
    res["inactiveThreatAgentDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInactiveThreatAgentDeviceCount(val)
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
    res["pendingFullScanDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingFullScanDeviceCount(val)
        }
        return nil
    }
    res["pendingManualStepsDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingManualStepsDeviceCount(val)
        }
        return nil
    }
    res["pendingOfflineScanDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingOfflineScanDeviceCount(val)
        }
        return nil
    }
    res["pendingQuickScanDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingQuickScanDeviceCount(val)
        }
        return nil
    }
    res["pendingRestartDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingRestartDeviceCount(val)
        }
        return nil
    }
    res["pendingSignatureUpdateDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingSignatureUpdateDeviceCount(val)
        }
        return nil
    }
    res["totalReportedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalReportedDeviceCount(val)
        }
        return nil
    }
    res["unknownStateThreatAgentDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnknownStateThreatAgentDeviceCount(val)
        }
        return nil
    }
    return res
}
// GetInactiveThreatAgentDeviceCount gets the inactiveThreatAgentDeviceCount property value. Indicates number of devices with inactive threat agent
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetInactiveThreatAgentDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("inactiveThreatAgentDeviceCount")
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
func (m *DeviceProtectionOverview) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPendingFullScanDeviceCount gets the pendingFullScanDeviceCount property value. Indicates number of devices pending full scan
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingFullScanDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingFullScanDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingManualStepsDeviceCount gets the pendingManualStepsDeviceCount property value. Indicates number of devices with pending manual steps
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingManualStepsDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingManualStepsDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingOfflineScanDeviceCount gets the pendingOfflineScanDeviceCount property value. Indicates number of pending offline scan devices
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingOfflineScanDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingOfflineScanDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingQuickScanDeviceCount gets the pendingQuickScanDeviceCount property value. Indicates the number of devices that have a pending full scan. Valid values -2147483648 to 2147483647
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingQuickScanDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingQuickScanDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingRestartDeviceCount gets the pendingRestartDeviceCount property value. Indicates number of devices pending restart
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingRestartDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingRestartDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingSignatureUpdateDeviceCount gets the pendingSignatureUpdateDeviceCount property value. Indicates number of devices with an old signature
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetPendingSignatureUpdateDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingSignatureUpdateDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTotalReportedDeviceCount gets the totalReportedDeviceCount property value. Total device count.
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetTotalReportedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalReportedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnknownStateThreatAgentDeviceCount gets the unknownStateThreatAgentDeviceCount property value. Indicates number of devices with threat agent state as unknown
// returns a *int32 when successful
func (m *DeviceProtectionOverview) GetUnknownStateThreatAgentDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownStateThreatAgentDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceProtectionOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("cleanDeviceCount", m.GetCleanDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("criticalFailuresDeviceCount", m.GetCriticalFailuresDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("inactiveThreatAgentDeviceCount", m.GetInactiveThreatAgentDeviceCount())
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
        err := writer.WriteInt32Value("pendingFullScanDeviceCount", m.GetPendingFullScanDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pendingManualStepsDeviceCount", m.GetPendingManualStepsDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pendingOfflineScanDeviceCount", m.GetPendingOfflineScanDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pendingQuickScanDeviceCount", m.GetPendingQuickScanDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pendingRestartDeviceCount", m.GetPendingRestartDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pendingSignatureUpdateDeviceCount", m.GetPendingSignatureUpdateDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalReportedDeviceCount", m.GetTotalReportedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("unknownStateThreatAgentDeviceCount", m.GetUnknownStateThreatAgentDeviceCount())
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
func (m *DeviceProtectionOverview) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DeviceProtectionOverview) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCleanDeviceCount sets the cleanDeviceCount property value. Indicates number of devices reporting as clean
func (m *DeviceProtectionOverview) SetCleanDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("cleanDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCriticalFailuresDeviceCount sets the criticalFailuresDeviceCount property value. Indicates number of devices with critical failures
func (m *DeviceProtectionOverview) SetCriticalFailuresDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("criticalFailuresDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInactiveThreatAgentDeviceCount sets the inactiveThreatAgentDeviceCount property value. Indicates number of devices with inactive threat agent
func (m *DeviceProtectionOverview) SetInactiveThreatAgentDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("inactiveThreatAgentDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DeviceProtectionOverview) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingFullScanDeviceCount sets the pendingFullScanDeviceCount property value. Indicates number of devices pending full scan
func (m *DeviceProtectionOverview) SetPendingFullScanDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingFullScanDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingManualStepsDeviceCount sets the pendingManualStepsDeviceCount property value. Indicates number of devices with pending manual steps
func (m *DeviceProtectionOverview) SetPendingManualStepsDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingManualStepsDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingOfflineScanDeviceCount sets the pendingOfflineScanDeviceCount property value. Indicates number of pending offline scan devices
func (m *DeviceProtectionOverview) SetPendingOfflineScanDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingOfflineScanDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingQuickScanDeviceCount sets the pendingQuickScanDeviceCount property value. Indicates the number of devices that have a pending full scan. Valid values -2147483648 to 2147483647
func (m *DeviceProtectionOverview) SetPendingQuickScanDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingQuickScanDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingRestartDeviceCount sets the pendingRestartDeviceCount property value. Indicates number of devices pending restart
func (m *DeviceProtectionOverview) SetPendingRestartDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingRestartDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingSignatureUpdateDeviceCount sets the pendingSignatureUpdateDeviceCount property value. Indicates number of devices with an old signature
func (m *DeviceProtectionOverview) SetPendingSignatureUpdateDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingSignatureUpdateDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalReportedDeviceCount sets the totalReportedDeviceCount property value. Total device count.
func (m *DeviceProtectionOverview) SetTotalReportedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("totalReportedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownStateThreatAgentDeviceCount sets the unknownStateThreatAgentDeviceCount property value. Indicates number of devices with threat agent state as unknown
func (m *DeviceProtectionOverview) SetUnknownStateThreatAgentDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownStateThreatAgentDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceProtectionOverviewable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCleanDeviceCount()(*int32)
    GetCriticalFailuresDeviceCount()(*int32)
    GetInactiveThreatAgentDeviceCount()(*int32)
    GetOdataType()(*string)
    GetPendingFullScanDeviceCount()(*int32)
    GetPendingManualStepsDeviceCount()(*int32)
    GetPendingOfflineScanDeviceCount()(*int32)
    GetPendingQuickScanDeviceCount()(*int32)
    GetPendingRestartDeviceCount()(*int32)
    GetPendingSignatureUpdateDeviceCount()(*int32)
    GetTotalReportedDeviceCount()(*int32)
    GetUnknownStateThreatAgentDeviceCount()(*int32)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCleanDeviceCount(value *int32)()
    SetCriticalFailuresDeviceCount(value *int32)()
    SetInactiveThreatAgentDeviceCount(value *int32)()
    SetOdataType(value *string)()
    SetPendingFullScanDeviceCount(value *int32)()
    SetPendingManualStepsDeviceCount(value *int32)()
    SetPendingOfflineScanDeviceCount(value *int32)()
    SetPendingQuickScanDeviceCount(value *int32)()
    SetPendingRestartDeviceCount(value *int32)()
    SetPendingSignatureUpdateDeviceCount(value *int32)()
    SetTotalReportedDeviceCount(value *int32)()
    SetUnknownStateThreatAgentDeviceCount(value *int32)()
}

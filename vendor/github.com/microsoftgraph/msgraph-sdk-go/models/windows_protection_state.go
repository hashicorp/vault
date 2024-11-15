package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsProtectionState device protection status entity.
type WindowsProtectionState struct {
    Entity
}
// NewWindowsProtectionState instantiates a new WindowsProtectionState and sets the default values.
func NewWindowsProtectionState()(*WindowsProtectionState) {
    m := &WindowsProtectionState{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWindowsProtectionStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsProtectionStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsProtectionState(), nil
}
// GetAntiMalwareVersion gets the antiMalwareVersion property value. Current anti malware version
// returns a *string when successful
func (m *WindowsProtectionState) GetAntiMalwareVersion()(*string) {
    val, err := m.GetBackingStore().Get("antiMalwareVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDetectedMalwareState gets the detectedMalwareState property value. Device malware list
// returns a []WindowsDeviceMalwareStateable when successful
func (m *WindowsProtectionState) GetDetectedMalwareState()([]WindowsDeviceMalwareStateable) {
    val, err := m.GetBackingStore().Get("detectedMalwareState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsDeviceMalwareStateable)
    }
    return nil
}
// GetDeviceState gets the deviceState property value. Indicates device's health state. Possible values are: clean, fullScanPending, rebootPending, manualStepsPending, offlineScanPending, critical. Possible values are: clean, fullScanPending, rebootPending, manualStepsPending, offlineScanPending, critical.
// returns a *WindowsDeviceHealthState when successful
func (m *WindowsProtectionState) GetDeviceState()(*WindowsDeviceHealthState) {
    val, err := m.GetBackingStore().Get("deviceState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsDeviceHealthState)
    }
    return nil
}
// GetEngineVersion gets the engineVersion property value. Current endpoint protection engine's version
// returns a *string when successful
func (m *WindowsProtectionState) GetEngineVersion()(*string) {
    val, err := m.GetBackingStore().Get("engineVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsProtectionState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["antiMalwareVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAntiMalwareVersion(val)
        }
        return nil
    }
    res["detectedMalwareState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsDeviceMalwareStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsDeviceMalwareStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsDeviceMalwareStateable)
                }
            }
            m.SetDetectedMalwareState(res)
        }
        return nil
    }
    res["deviceState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsDeviceHealthState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceState(val.(*WindowsDeviceHealthState))
        }
        return nil
    }
    res["engineVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEngineVersion(val)
        }
        return nil
    }
    res["fullScanOverdue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFullScanOverdue(val)
        }
        return nil
    }
    res["fullScanRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFullScanRequired(val)
        }
        return nil
    }
    res["isVirtualMachine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVirtualMachine(val)
        }
        return nil
    }
    res["lastFullScanDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastFullScanDateTime(val)
        }
        return nil
    }
    res["lastFullScanSignatureVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastFullScanSignatureVersion(val)
        }
        return nil
    }
    res["lastQuickScanDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastQuickScanDateTime(val)
        }
        return nil
    }
    res["lastQuickScanSignatureVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastQuickScanSignatureVersion(val)
        }
        return nil
    }
    res["lastReportedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastReportedDateTime(val)
        }
        return nil
    }
    res["malwareProtectionEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMalwareProtectionEnabled(val)
        }
        return nil
    }
    res["networkInspectionSystemEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkInspectionSystemEnabled(val)
        }
        return nil
    }
    res["productStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsDefenderProductStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductStatus(val.(*WindowsDefenderProductStatus))
        }
        return nil
    }
    res["quickScanOverdue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuickScanOverdue(val)
        }
        return nil
    }
    res["realTimeProtectionEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRealTimeProtectionEnabled(val)
        }
        return nil
    }
    res["rebootRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRebootRequired(val)
        }
        return nil
    }
    res["signatureUpdateOverdue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignatureUpdateOverdue(val)
        }
        return nil
    }
    res["signatureVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignatureVersion(val)
        }
        return nil
    }
    res["tamperProtectionEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTamperProtectionEnabled(val)
        }
        return nil
    }
    return res
}
// GetFullScanOverdue gets the fullScanOverdue property value. When TRUE indicates full scan is overdue, when FALSE indicates full scan is not overdue. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetFullScanOverdue()(*bool) {
    val, err := m.GetBackingStore().Get("fullScanOverdue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFullScanRequired gets the fullScanRequired property value. When TRUE indicates full scan is required, when FALSE indicates full scan is not required. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetFullScanRequired()(*bool) {
    val, err := m.GetBackingStore().Get("fullScanRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsVirtualMachine gets the isVirtualMachine property value. When TRUE indicates the device is a virtual machine, when FALSE indicates the device is not a virtual machine. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetIsVirtualMachine()(*bool) {
    val, err := m.GetBackingStore().Get("isVirtualMachine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastFullScanDateTime gets the lastFullScanDateTime property value. Last quick scan datetime
// returns a *Time when successful
func (m *WindowsProtectionState) GetLastFullScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastFullScanDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastFullScanSignatureVersion gets the lastFullScanSignatureVersion property value. Last full scan signature version
// returns a *string when successful
func (m *WindowsProtectionState) GetLastFullScanSignatureVersion()(*string) {
    val, err := m.GetBackingStore().Get("lastFullScanSignatureVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastQuickScanDateTime gets the lastQuickScanDateTime property value. Last quick scan datetime
// returns a *Time when successful
func (m *WindowsProtectionState) GetLastQuickScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastQuickScanDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastQuickScanSignatureVersion gets the lastQuickScanSignatureVersion property value. Last quick scan signature version
// returns a *string when successful
func (m *WindowsProtectionState) GetLastQuickScanSignatureVersion()(*string) {
    val, err := m.GetBackingStore().Get("lastQuickScanSignatureVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastReportedDateTime gets the lastReportedDateTime property value. Last device health status reported time
// returns a *Time when successful
func (m *WindowsProtectionState) GetLastReportedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastReportedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMalwareProtectionEnabled gets the malwareProtectionEnabled property value. When TRUE indicates anti malware is enabled when FALSE indicates anti malware is not enabled.
// returns a *bool when successful
func (m *WindowsProtectionState) GetMalwareProtectionEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("malwareProtectionEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetworkInspectionSystemEnabled gets the networkInspectionSystemEnabled property value. When TRUE indicates network inspection system enabled, when FALSE indicates network inspection system is not enabled. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetNetworkInspectionSystemEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("networkInspectionSystemEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetProductStatus gets the productStatus property value. Product Status of Windows Defender Antivirus. Possible values are: noStatus, serviceNotRunning, serviceStartedWithoutMalwareProtection, pendingFullScanDueToThreatAction, pendingRebootDueToThreatAction, pendingManualStepsDueToThreatAction, avSignaturesOutOfDate, asSignaturesOutOfDate, noQuickScanHappenedForSpecifiedPeriod, noFullScanHappenedForSpecifiedPeriod, systemInitiatedScanInProgress, systemInitiatedCleanInProgress, samplesPendingSubmission, productRunningInEvaluationMode, productRunningInNonGenuineMode, productExpired, offlineScanRequired, serviceShutdownAsPartOfSystemShutdown, threatRemediationFailedCritically, threatRemediationFailedNonCritically, noStatusFlagsSet, platformOutOfDate, platformUpdateInProgress, platformAboutToBeOutdated, signatureOrPlatformEndOfLifeIsPastOrIsImpending, windowsSModeSignaturesInUseOnNonWin10SInstall. Possible values are: noStatus, serviceNotRunning, serviceStartedWithoutMalwareProtection, pendingFullScanDueToThreatAction, pendingRebootDueToThreatAction, pendingManualStepsDueToThreatAction, avSignaturesOutOfDate, asSignaturesOutOfDate, noQuickScanHappenedForSpecifiedPeriod, noFullScanHappenedForSpecifiedPeriod, systemInitiatedScanInProgress, systemInitiatedCleanInProgress, samplesPendingSubmission, productRunningInEvaluationMode, productRunningInNonGenuineMode, productExpired, offlineScanRequired, serviceShutdownAsPartOfSystemShutdown, threatRemediationFailedCritically, threatRemediationFailedNonCritically, noStatusFlagsSet, platformOutOfDate, platformUpdateInProgress, platformAboutToBeOutdated, signatureOrPlatformEndOfLifeIsPastOrIsImpending, windowsSModeSignaturesInUseOnNonWin10SInstall.
// returns a *WindowsDefenderProductStatus when successful
func (m *WindowsProtectionState) GetProductStatus()(*WindowsDefenderProductStatus) {
    val, err := m.GetBackingStore().Get("productStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsDefenderProductStatus)
    }
    return nil
}
// GetQuickScanOverdue gets the quickScanOverdue property value. When TRUE indicates quick scan is overdue, when FALSE indicates quick scan is not overdue. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetQuickScanOverdue()(*bool) {
    val, err := m.GetBackingStore().Get("quickScanOverdue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRealTimeProtectionEnabled gets the realTimeProtectionEnabled property value. When TRUE indicates real time protection is enabled, when FALSE indicates real time protection is not enabled. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetRealTimeProtectionEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("realTimeProtectionEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRebootRequired gets the rebootRequired property value. When TRUE indicates reboot is required, when FALSE indicates when TRUE indicates reboot is not required. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetRebootRequired()(*bool) {
    val, err := m.GetBackingStore().Get("rebootRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSignatureUpdateOverdue gets the signatureUpdateOverdue property value. When TRUE indicates signature is out of date, when FALSE indicates signature is not out of date. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetSignatureUpdateOverdue()(*bool) {
    val, err := m.GetBackingStore().Get("signatureUpdateOverdue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSignatureVersion gets the signatureVersion property value. Current malware definitions version
// returns a *string when successful
func (m *WindowsProtectionState) GetSignatureVersion()(*string) {
    val, err := m.GetBackingStore().Get("signatureVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTamperProtectionEnabled gets the tamperProtectionEnabled property value. When TRUE indicates the Windows Defender tamper protection feature is enabled, when FALSE indicates the Windows Defender tamper protection feature is not enabled. Defaults to setting on client device.
// returns a *bool when successful
func (m *WindowsProtectionState) GetTamperProtectionEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("tamperProtectionEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsProtectionState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("antiMalwareVersion", m.GetAntiMalwareVersion())
        if err != nil {
            return err
        }
    }
    if m.GetDetectedMalwareState() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDetectedMalwareState()))
        for i, v := range m.GetDetectedMalwareState() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("detectedMalwareState", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceState() != nil {
        cast := (*m.GetDeviceState()).String()
        err = writer.WriteStringValue("deviceState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("engineVersion", m.GetEngineVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("fullScanOverdue", m.GetFullScanOverdue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("fullScanRequired", m.GetFullScanRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isVirtualMachine", m.GetIsVirtualMachine())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastFullScanDateTime", m.GetLastFullScanDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lastFullScanSignatureVersion", m.GetLastFullScanSignatureVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastQuickScanDateTime", m.GetLastQuickScanDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lastQuickScanSignatureVersion", m.GetLastQuickScanSignatureVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastReportedDateTime", m.GetLastReportedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("malwareProtectionEnabled", m.GetMalwareProtectionEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("networkInspectionSystemEnabled", m.GetNetworkInspectionSystemEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetProductStatus() != nil {
        cast := (*m.GetProductStatus()).String()
        err = writer.WriteStringValue("productStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("quickScanOverdue", m.GetQuickScanOverdue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("realTimeProtectionEnabled", m.GetRealTimeProtectionEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("rebootRequired", m.GetRebootRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("signatureUpdateOverdue", m.GetSignatureUpdateOverdue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("signatureVersion", m.GetSignatureVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("tamperProtectionEnabled", m.GetTamperProtectionEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAntiMalwareVersion sets the antiMalwareVersion property value. Current anti malware version
func (m *WindowsProtectionState) SetAntiMalwareVersion(value *string)() {
    err := m.GetBackingStore().Set("antiMalwareVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectedMalwareState sets the detectedMalwareState property value. Device malware list
func (m *WindowsProtectionState) SetDetectedMalwareState(value []WindowsDeviceMalwareStateable)() {
    err := m.GetBackingStore().Set("detectedMalwareState", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceState sets the deviceState property value. Indicates device's health state. Possible values are: clean, fullScanPending, rebootPending, manualStepsPending, offlineScanPending, critical. Possible values are: clean, fullScanPending, rebootPending, manualStepsPending, offlineScanPending, critical.
func (m *WindowsProtectionState) SetDeviceState(value *WindowsDeviceHealthState)() {
    err := m.GetBackingStore().Set("deviceState", value)
    if err != nil {
        panic(err)
    }
}
// SetEngineVersion sets the engineVersion property value. Current endpoint protection engine's version
func (m *WindowsProtectionState) SetEngineVersion(value *string)() {
    err := m.GetBackingStore().Set("engineVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetFullScanOverdue sets the fullScanOverdue property value. When TRUE indicates full scan is overdue, when FALSE indicates full scan is not overdue. Defaults to setting on client device.
func (m *WindowsProtectionState) SetFullScanOverdue(value *bool)() {
    err := m.GetBackingStore().Set("fullScanOverdue", value)
    if err != nil {
        panic(err)
    }
}
// SetFullScanRequired sets the fullScanRequired property value. When TRUE indicates full scan is required, when FALSE indicates full scan is not required. Defaults to setting on client device.
func (m *WindowsProtectionState) SetFullScanRequired(value *bool)() {
    err := m.GetBackingStore().Set("fullScanRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVirtualMachine sets the isVirtualMachine property value. When TRUE indicates the device is a virtual machine, when FALSE indicates the device is not a virtual machine. Defaults to setting on client device.
func (m *WindowsProtectionState) SetIsVirtualMachine(value *bool)() {
    err := m.GetBackingStore().Set("isVirtualMachine", value)
    if err != nil {
        panic(err)
    }
}
// SetLastFullScanDateTime sets the lastFullScanDateTime property value. Last quick scan datetime
func (m *WindowsProtectionState) SetLastFullScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastFullScanDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastFullScanSignatureVersion sets the lastFullScanSignatureVersion property value. Last full scan signature version
func (m *WindowsProtectionState) SetLastFullScanSignatureVersion(value *string)() {
    err := m.GetBackingStore().Set("lastFullScanSignatureVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetLastQuickScanDateTime sets the lastQuickScanDateTime property value. Last quick scan datetime
func (m *WindowsProtectionState) SetLastQuickScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastQuickScanDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastQuickScanSignatureVersion sets the lastQuickScanSignatureVersion property value. Last quick scan signature version
func (m *WindowsProtectionState) SetLastQuickScanSignatureVersion(value *string)() {
    err := m.GetBackingStore().Set("lastQuickScanSignatureVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetLastReportedDateTime sets the lastReportedDateTime property value. Last device health status reported time
func (m *WindowsProtectionState) SetLastReportedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastReportedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMalwareProtectionEnabled sets the malwareProtectionEnabled property value. When TRUE indicates anti malware is enabled when FALSE indicates anti malware is not enabled.
func (m *WindowsProtectionState) SetMalwareProtectionEnabled(value *bool)() {
    err := m.GetBackingStore().Set("malwareProtectionEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkInspectionSystemEnabled sets the networkInspectionSystemEnabled property value. When TRUE indicates network inspection system enabled, when FALSE indicates network inspection system is not enabled. Defaults to setting on client device.
func (m *WindowsProtectionState) SetNetworkInspectionSystemEnabled(value *bool)() {
    err := m.GetBackingStore().Set("networkInspectionSystemEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetProductStatus sets the productStatus property value. Product Status of Windows Defender Antivirus. Possible values are: noStatus, serviceNotRunning, serviceStartedWithoutMalwareProtection, pendingFullScanDueToThreatAction, pendingRebootDueToThreatAction, pendingManualStepsDueToThreatAction, avSignaturesOutOfDate, asSignaturesOutOfDate, noQuickScanHappenedForSpecifiedPeriod, noFullScanHappenedForSpecifiedPeriod, systemInitiatedScanInProgress, systemInitiatedCleanInProgress, samplesPendingSubmission, productRunningInEvaluationMode, productRunningInNonGenuineMode, productExpired, offlineScanRequired, serviceShutdownAsPartOfSystemShutdown, threatRemediationFailedCritically, threatRemediationFailedNonCritically, noStatusFlagsSet, platformOutOfDate, platformUpdateInProgress, platformAboutToBeOutdated, signatureOrPlatformEndOfLifeIsPastOrIsImpending, windowsSModeSignaturesInUseOnNonWin10SInstall. Possible values are: noStatus, serviceNotRunning, serviceStartedWithoutMalwareProtection, pendingFullScanDueToThreatAction, pendingRebootDueToThreatAction, pendingManualStepsDueToThreatAction, avSignaturesOutOfDate, asSignaturesOutOfDate, noQuickScanHappenedForSpecifiedPeriod, noFullScanHappenedForSpecifiedPeriod, systemInitiatedScanInProgress, systemInitiatedCleanInProgress, samplesPendingSubmission, productRunningInEvaluationMode, productRunningInNonGenuineMode, productExpired, offlineScanRequired, serviceShutdownAsPartOfSystemShutdown, threatRemediationFailedCritically, threatRemediationFailedNonCritically, noStatusFlagsSet, platformOutOfDate, platformUpdateInProgress, platformAboutToBeOutdated, signatureOrPlatformEndOfLifeIsPastOrIsImpending, windowsSModeSignaturesInUseOnNonWin10SInstall.
func (m *WindowsProtectionState) SetProductStatus(value *WindowsDefenderProductStatus)() {
    err := m.GetBackingStore().Set("productStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetQuickScanOverdue sets the quickScanOverdue property value. When TRUE indicates quick scan is overdue, when FALSE indicates quick scan is not overdue. Defaults to setting on client device.
func (m *WindowsProtectionState) SetQuickScanOverdue(value *bool)() {
    err := m.GetBackingStore().Set("quickScanOverdue", value)
    if err != nil {
        panic(err)
    }
}
// SetRealTimeProtectionEnabled sets the realTimeProtectionEnabled property value. When TRUE indicates real time protection is enabled, when FALSE indicates real time protection is not enabled. Defaults to setting on client device.
func (m *WindowsProtectionState) SetRealTimeProtectionEnabled(value *bool)() {
    err := m.GetBackingStore().Set("realTimeProtectionEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetRebootRequired sets the rebootRequired property value. When TRUE indicates reboot is required, when FALSE indicates when TRUE indicates reboot is not required. Defaults to setting on client device.
func (m *WindowsProtectionState) SetRebootRequired(value *bool)() {
    err := m.GetBackingStore().Set("rebootRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetSignatureUpdateOverdue sets the signatureUpdateOverdue property value. When TRUE indicates signature is out of date, when FALSE indicates signature is not out of date. Defaults to setting on client device.
func (m *WindowsProtectionState) SetSignatureUpdateOverdue(value *bool)() {
    err := m.GetBackingStore().Set("signatureUpdateOverdue", value)
    if err != nil {
        panic(err)
    }
}
// SetSignatureVersion sets the signatureVersion property value. Current malware definitions version
func (m *WindowsProtectionState) SetSignatureVersion(value *string)() {
    err := m.GetBackingStore().Set("signatureVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetTamperProtectionEnabled sets the tamperProtectionEnabled property value. When TRUE indicates the Windows Defender tamper protection feature is enabled, when FALSE indicates the Windows Defender tamper protection feature is not enabled. Defaults to setting on client device.
func (m *WindowsProtectionState) SetTamperProtectionEnabled(value *bool)() {
    err := m.GetBackingStore().Set("tamperProtectionEnabled", value)
    if err != nil {
        panic(err)
    }
}
type WindowsProtectionStateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAntiMalwareVersion()(*string)
    GetDetectedMalwareState()([]WindowsDeviceMalwareStateable)
    GetDeviceState()(*WindowsDeviceHealthState)
    GetEngineVersion()(*string)
    GetFullScanOverdue()(*bool)
    GetFullScanRequired()(*bool)
    GetIsVirtualMachine()(*bool)
    GetLastFullScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastFullScanSignatureVersion()(*string)
    GetLastQuickScanDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastQuickScanSignatureVersion()(*string)
    GetLastReportedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMalwareProtectionEnabled()(*bool)
    GetNetworkInspectionSystemEnabled()(*bool)
    GetProductStatus()(*WindowsDefenderProductStatus)
    GetQuickScanOverdue()(*bool)
    GetRealTimeProtectionEnabled()(*bool)
    GetRebootRequired()(*bool)
    GetSignatureUpdateOverdue()(*bool)
    GetSignatureVersion()(*string)
    GetTamperProtectionEnabled()(*bool)
    SetAntiMalwareVersion(value *string)()
    SetDetectedMalwareState(value []WindowsDeviceMalwareStateable)()
    SetDeviceState(value *WindowsDeviceHealthState)()
    SetEngineVersion(value *string)()
    SetFullScanOverdue(value *bool)()
    SetFullScanRequired(value *bool)()
    SetIsVirtualMachine(value *bool)()
    SetLastFullScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastFullScanSignatureVersion(value *string)()
    SetLastQuickScanDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastQuickScanSignatureVersion(value *string)()
    SetLastReportedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMalwareProtectionEnabled(value *bool)()
    SetNetworkInspectionSystemEnabled(value *bool)()
    SetProductStatus(value *WindowsDefenderProductStatus)()
    SetQuickScanOverdue(value *bool)()
    SetRealTimeProtectionEnabled(value *bool)()
    SetRebootRequired(value *bool)()
    SetSignatureUpdateOverdue(value *bool)()
    SetSignatureVersion(value *string)()
    SetTamperProtectionEnabled(value *bool)()
}

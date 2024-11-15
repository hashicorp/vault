package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DeviceHealthAttestationState struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDeviceHealthAttestationState instantiates a new DeviceHealthAttestationState and sets the default values.
func NewDeviceHealthAttestationState()(*DeviceHealthAttestationState) {
    m := &DeviceHealthAttestationState{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDeviceHealthAttestationStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceHealthAttestationStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceHealthAttestationState(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DeviceHealthAttestationState) GetAdditionalData()(map[string]any) {
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
// GetAttestationIdentityKey gets the attestationIdentityKey property value. TWhen an Attestation Identity Key (AIK) is present on a device, it indicates that the device has an endorsement key (EK) certificate.
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetAttestationIdentityKey()(*string) {
    val, err := m.GetBackingStore().Get("attestationIdentityKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DeviceHealthAttestationState) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBitLockerStatus gets the bitLockerStatus property value. On or Off of BitLocker Drive Encryption
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBitLockerStatus()(*string) {
    val, err := m.GetBackingStore().Get("bitLockerStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBootAppSecurityVersion gets the bootAppSecurityVersion property value. The security version number of the Boot Application
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBootAppSecurityVersion()(*string) {
    val, err := m.GetBackingStore().Get("bootAppSecurityVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBootDebugging gets the bootDebugging property value. When bootDebugging is enabled, the device is used in development and testing
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBootDebugging()(*string) {
    val, err := m.GetBackingStore().Get("bootDebugging")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBootManagerSecurityVersion gets the bootManagerSecurityVersion property value. The security version number of the Boot Application
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBootManagerSecurityVersion()(*string) {
    val, err := m.GetBackingStore().Get("bootManagerSecurityVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBootManagerVersion gets the bootManagerVersion property value. The version of the Boot Manager
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBootManagerVersion()(*string) {
    val, err := m.GetBackingStore().Get("bootManagerVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBootRevisionListInfo gets the bootRevisionListInfo property value. The Boot Revision List that was loaded during initial boot on the attested device
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetBootRevisionListInfo()(*string) {
    val, err := m.GetBackingStore().Get("bootRevisionListInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCodeIntegrity gets the codeIntegrity property value. When code integrity is enabled, code execution is restricted to integrity verified code
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetCodeIntegrity()(*string) {
    val, err := m.GetBackingStore().Get("codeIntegrity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCodeIntegrityCheckVersion gets the codeIntegrityCheckVersion property value. The version of the Boot Manager
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetCodeIntegrityCheckVersion()(*string) {
    val, err := m.GetBackingStore().Get("codeIntegrityCheckVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCodeIntegrityPolicy gets the codeIntegrityPolicy property value. The Code Integrity policy that is controlling the security of the boot environment
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetCodeIntegrityPolicy()(*string) {
    val, err := m.GetBackingStore().Get("codeIntegrityPolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentNamespaceUrl gets the contentNamespaceUrl property value. The DHA report version. (Namespace version)
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetContentNamespaceUrl()(*string) {
    val, err := m.GetBackingStore().Get("contentNamespaceUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentVersion gets the contentVersion property value. The HealthAttestation state schema version
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetContentVersion()(*string) {
    val, err := m.GetBackingStore().Get("contentVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDataExcutionPolicy gets the dataExcutionPolicy property value. DEP Policy defines a set of hardware and software technologies that perform additional checks on memory
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetDataExcutionPolicy()(*string) {
    val, err := m.GetBackingStore().Get("dataExcutionPolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceHealthAttestationStatus gets the deviceHealthAttestationStatus property value. The DHA report version. (Namespace version)
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetDeviceHealthAttestationStatus()(*string) {
    val, err := m.GetBackingStore().Get("deviceHealthAttestationStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEarlyLaunchAntiMalwareDriverProtection gets the earlyLaunchAntiMalwareDriverProtection property value. ELAM provides protection for the computers in your network when they start up
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetEarlyLaunchAntiMalwareDriverProtection()(*string) {
    val, err := m.GetBackingStore().Get("earlyLaunchAntiMalwareDriverProtection")
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
func (m *DeviceHealthAttestationState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attestationIdentityKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttestationIdentityKey(val)
        }
        return nil
    }
    res["bitLockerStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitLockerStatus(val)
        }
        return nil
    }
    res["bootAppSecurityVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootAppSecurityVersion(val)
        }
        return nil
    }
    res["bootDebugging"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootDebugging(val)
        }
        return nil
    }
    res["bootManagerSecurityVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootManagerSecurityVersion(val)
        }
        return nil
    }
    res["bootManagerVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootManagerVersion(val)
        }
        return nil
    }
    res["bootRevisionListInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBootRevisionListInfo(val)
        }
        return nil
    }
    res["codeIntegrity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCodeIntegrity(val)
        }
        return nil
    }
    res["codeIntegrityCheckVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCodeIntegrityCheckVersion(val)
        }
        return nil
    }
    res["codeIntegrityPolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCodeIntegrityPolicy(val)
        }
        return nil
    }
    res["contentNamespaceUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentNamespaceUrl(val)
        }
        return nil
    }
    res["contentVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentVersion(val)
        }
        return nil
    }
    res["dataExcutionPolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataExcutionPolicy(val)
        }
        return nil
    }
    res["deviceHealthAttestationStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceHealthAttestationStatus(val)
        }
        return nil
    }
    res["earlyLaunchAntiMalwareDriverProtection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEarlyLaunchAntiMalwareDriverProtection(val)
        }
        return nil
    }
    res["healthAttestationSupportedStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHealthAttestationSupportedStatus(val)
        }
        return nil
    }
    res["healthStatusMismatchInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHealthStatusMismatchInfo(val)
        }
        return nil
    }
    res["issuedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIssuedDateTime(val)
        }
        return nil
    }
    res["lastUpdateDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdateDateTime(val)
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
    res["operatingSystemKernelDebugging"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystemKernelDebugging(val)
        }
        return nil
    }
    res["operatingSystemRevListInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystemRevListInfo(val)
        }
        return nil
    }
    res["pcr0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPcr0(val)
        }
        return nil
    }
    res["pcrHashAlgorithm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPcrHashAlgorithm(val)
        }
        return nil
    }
    res["resetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResetCount(val)
        }
        return nil
    }
    res["restartCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestartCount(val)
        }
        return nil
    }
    res["safeMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafeMode(val)
        }
        return nil
    }
    res["secureBoot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecureBoot(val)
        }
        return nil
    }
    res["secureBootConfigurationPolicyFingerPrint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecureBootConfigurationPolicyFingerPrint(val)
        }
        return nil
    }
    res["testSigning"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTestSigning(val)
        }
        return nil
    }
    res["tpmVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTpmVersion(val)
        }
        return nil
    }
    res["virtualSecureMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVirtualSecureMode(val)
        }
        return nil
    }
    res["windowsPE"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsPE(val)
        }
        return nil
    }
    return res
}
// GetHealthAttestationSupportedStatus gets the healthAttestationSupportedStatus property value. This attribute indicates if DHA is supported for the device
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetHealthAttestationSupportedStatus()(*string) {
    val, err := m.GetBackingStore().Get("healthAttestationSupportedStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHealthStatusMismatchInfo gets the healthStatusMismatchInfo property value. This attribute appears if DHA-Service detects an integrity issue
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetHealthStatusMismatchInfo()(*string) {
    val, err := m.GetBackingStore().Get("healthStatusMismatchInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIssuedDateTime gets the issuedDateTime property value. The DateTime when device was evaluated or issued to MDM
// returns a *Time when successful
func (m *DeviceHealthAttestationState) GetIssuedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("issuedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastUpdateDateTime gets the lastUpdateDateTime property value. The Timestamp of the last update.
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetLastUpdateDateTime()(*string) {
    val, err := m.GetBackingStore().Get("lastUpdateDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperatingSystemKernelDebugging gets the operatingSystemKernelDebugging property value. When operatingSystemKernelDebugging is enabled, the device is used in development and testing
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetOperatingSystemKernelDebugging()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystemKernelDebugging")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperatingSystemRevListInfo gets the operatingSystemRevListInfo property value. The Operating System Revision List that was loaded during initial boot on the attested device
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetOperatingSystemRevListInfo()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystemRevListInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPcr0 gets the pcr0 property value. The measurement that is captured in PCR[0]
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetPcr0()(*string) {
    val, err := m.GetBackingStore().Get("pcr0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPcrHashAlgorithm gets the pcrHashAlgorithm property value. Informational attribute that identifies the HASH algorithm that was used by TPM
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetPcrHashAlgorithm()(*string) {
    val, err := m.GetBackingStore().Get("pcrHashAlgorithm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResetCount gets the resetCount property value. The number of times a PC device has hibernated or resumed
// returns a *int64 when successful
func (m *DeviceHealthAttestationState) GetResetCount()(*int64) {
    val, err := m.GetBackingStore().Get("resetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetRestartCount gets the restartCount property value. The number of times a PC device has rebooted
// returns a *int64 when successful
func (m *DeviceHealthAttestationState) GetRestartCount()(*int64) {
    val, err := m.GetBackingStore().Get("restartCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetSafeMode gets the safeMode property value. Safe mode is a troubleshooting option for Windows that starts your computer in a limited state
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetSafeMode()(*string) {
    val, err := m.GetBackingStore().Get("safeMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSecureBoot gets the secureBoot property value. When Secure Boot is enabled, the core components must have the correct cryptographic signatures
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetSecureBoot()(*string) {
    val, err := m.GetBackingStore().Get("secureBoot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSecureBootConfigurationPolicyFingerPrint gets the secureBootConfigurationPolicyFingerPrint property value. Fingerprint of the Custom Secure Boot Configuration Policy
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetSecureBootConfigurationPolicyFingerPrint()(*string) {
    val, err := m.GetBackingStore().Get("secureBootConfigurationPolicyFingerPrint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTestSigning gets the testSigning property value. When test signing is allowed, the device does not enforce signature validation during boot
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetTestSigning()(*string) {
    val, err := m.GetBackingStore().Get("testSigning")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTpmVersion gets the tpmVersion property value. The security version number of the Boot Application
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetTpmVersion()(*string) {
    val, err := m.GetBackingStore().Get("tpmVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVirtualSecureMode gets the virtualSecureMode property value. VSM is a container that protects high value assets from a compromised kernel
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetVirtualSecureMode()(*string) {
    val, err := m.GetBackingStore().Get("virtualSecureMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWindowsPE gets the windowsPE property value. Operating system running with limited services that is used to prepare a computer for Windows
// returns a *string when successful
func (m *DeviceHealthAttestationState) GetWindowsPE()(*string) {
    val, err := m.GetBackingStore().Get("windowsPE")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceHealthAttestationState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("attestationIdentityKey", m.GetAttestationIdentityKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bitLockerStatus", m.GetBitLockerStatus())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bootAppSecurityVersion", m.GetBootAppSecurityVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bootDebugging", m.GetBootDebugging())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bootManagerSecurityVersion", m.GetBootManagerSecurityVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bootManagerVersion", m.GetBootManagerVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bootRevisionListInfo", m.GetBootRevisionListInfo())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("codeIntegrity", m.GetCodeIntegrity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("codeIntegrityCheckVersion", m.GetCodeIntegrityCheckVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("codeIntegrityPolicy", m.GetCodeIntegrityPolicy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("contentNamespaceUrl", m.GetContentNamespaceUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("contentVersion", m.GetContentVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("dataExcutionPolicy", m.GetDataExcutionPolicy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deviceHealthAttestationStatus", m.GetDeviceHealthAttestationStatus())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("earlyLaunchAntiMalwareDriverProtection", m.GetEarlyLaunchAntiMalwareDriverProtection())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("healthAttestationSupportedStatus", m.GetHealthAttestationSupportedStatus())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("healthStatusMismatchInfo", m.GetHealthStatusMismatchInfo())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("issuedDateTime", m.GetIssuedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("lastUpdateDateTime", m.GetLastUpdateDateTime())
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
        err := writer.WriteStringValue("operatingSystemKernelDebugging", m.GetOperatingSystemKernelDebugging())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("operatingSystemRevListInfo", m.GetOperatingSystemRevListInfo())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("pcr0", m.GetPcr0())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("pcrHashAlgorithm", m.GetPcrHashAlgorithm())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("resetCount", m.GetResetCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("restartCount", m.GetRestartCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("safeMode", m.GetSafeMode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("secureBoot", m.GetSecureBoot())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("secureBootConfigurationPolicyFingerPrint", m.GetSecureBootConfigurationPolicyFingerPrint())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("testSigning", m.GetTestSigning())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("tpmVersion", m.GetTpmVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("virtualSecureMode", m.GetVirtualSecureMode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("windowsPE", m.GetWindowsPE())
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
func (m *DeviceHealthAttestationState) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttestationIdentityKey sets the attestationIdentityKey property value. TWhen an Attestation Identity Key (AIK) is present on a device, it indicates that the device has an endorsement key (EK) certificate.
func (m *DeviceHealthAttestationState) SetAttestationIdentityKey(value *string)() {
    err := m.GetBackingStore().Set("attestationIdentityKey", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DeviceHealthAttestationState) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBitLockerStatus sets the bitLockerStatus property value. On or Off of BitLocker Drive Encryption
func (m *DeviceHealthAttestationState) SetBitLockerStatus(value *string)() {
    err := m.GetBackingStore().Set("bitLockerStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetBootAppSecurityVersion sets the bootAppSecurityVersion property value. The security version number of the Boot Application
func (m *DeviceHealthAttestationState) SetBootAppSecurityVersion(value *string)() {
    err := m.GetBackingStore().Set("bootAppSecurityVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetBootDebugging sets the bootDebugging property value. When bootDebugging is enabled, the device is used in development and testing
func (m *DeviceHealthAttestationState) SetBootDebugging(value *string)() {
    err := m.GetBackingStore().Set("bootDebugging", value)
    if err != nil {
        panic(err)
    }
}
// SetBootManagerSecurityVersion sets the bootManagerSecurityVersion property value. The security version number of the Boot Application
func (m *DeviceHealthAttestationState) SetBootManagerSecurityVersion(value *string)() {
    err := m.GetBackingStore().Set("bootManagerSecurityVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetBootManagerVersion sets the bootManagerVersion property value. The version of the Boot Manager
func (m *DeviceHealthAttestationState) SetBootManagerVersion(value *string)() {
    err := m.GetBackingStore().Set("bootManagerVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetBootRevisionListInfo sets the bootRevisionListInfo property value. The Boot Revision List that was loaded during initial boot on the attested device
func (m *DeviceHealthAttestationState) SetBootRevisionListInfo(value *string)() {
    err := m.GetBackingStore().Set("bootRevisionListInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetCodeIntegrity sets the codeIntegrity property value. When code integrity is enabled, code execution is restricted to integrity verified code
func (m *DeviceHealthAttestationState) SetCodeIntegrity(value *string)() {
    err := m.GetBackingStore().Set("codeIntegrity", value)
    if err != nil {
        panic(err)
    }
}
// SetCodeIntegrityCheckVersion sets the codeIntegrityCheckVersion property value. The version of the Boot Manager
func (m *DeviceHealthAttestationState) SetCodeIntegrityCheckVersion(value *string)() {
    err := m.GetBackingStore().Set("codeIntegrityCheckVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetCodeIntegrityPolicy sets the codeIntegrityPolicy property value. The Code Integrity policy that is controlling the security of the boot environment
func (m *DeviceHealthAttestationState) SetCodeIntegrityPolicy(value *string)() {
    err := m.GetBackingStore().Set("codeIntegrityPolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetContentNamespaceUrl sets the contentNamespaceUrl property value. The DHA report version. (Namespace version)
func (m *DeviceHealthAttestationState) SetContentNamespaceUrl(value *string)() {
    err := m.GetBackingStore().Set("contentNamespaceUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetContentVersion sets the contentVersion property value. The HealthAttestation state schema version
func (m *DeviceHealthAttestationState) SetContentVersion(value *string)() {
    err := m.GetBackingStore().Set("contentVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetDataExcutionPolicy sets the dataExcutionPolicy property value. DEP Policy defines a set of hardware and software technologies that perform additional checks on memory
func (m *DeviceHealthAttestationState) SetDataExcutionPolicy(value *string)() {
    err := m.GetBackingStore().Set("dataExcutionPolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceHealthAttestationStatus sets the deviceHealthAttestationStatus property value. The DHA report version. (Namespace version)
func (m *DeviceHealthAttestationState) SetDeviceHealthAttestationStatus(value *string)() {
    err := m.GetBackingStore().Set("deviceHealthAttestationStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetEarlyLaunchAntiMalwareDriverProtection sets the earlyLaunchAntiMalwareDriverProtection property value. ELAM provides protection for the computers in your network when they start up
func (m *DeviceHealthAttestationState) SetEarlyLaunchAntiMalwareDriverProtection(value *string)() {
    err := m.GetBackingStore().Set("earlyLaunchAntiMalwareDriverProtection", value)
    if err != nil {
        panic(err)
    }
}
// SetHealthAttestationSupportedStatus sets the healthAttestationSupportedStatus property value. This attribute indicates if DHA is supported for the device
func (m *DeviceHealthAttestationState) SetHealthAttestationSupportedStatus(value *string)() {
    err := m.GetBackingStore().Set("healthAttestationSupportedStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetHealthStatusMismatchInfo sets the healthStatusMismatchInfo property value. This attribute appears if DHA-Service detects an integrity issue
func (m *DeviceHealthAttestationState) SetHealthStatusMismatchInfo(value *string)() {
    err := m.GetBackingStore().Set("healthStatusMismatchInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetIssuedDateTime sets the issuedDateTime property value. The DateTime when device was evaluated or issued to MDM
func (m *DeviceHealthAttestationState) SetIssuedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("issuedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdateDateTime sets the lastUpdateDateTime property value. The Timestamp of the last update.
func (m *DeviceHealthAttestationState) SetLastUpdateDateTime(value *string)() {
    err := m.GetBackingStore().Set("lastUpdateDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DeviceHealthAttestationState) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystemKernelDebugging sets the operatingSystemKernelDebugging property value. When operatingSystemKernelDebugging is enabled, the device is used in development and testing
func (m *DeviceHealthAttestationState) SetOperatingSystemKernelDebugging(value *string)() {
    err := m.GetBackingStore().Set("operatingSystemKernelDebugging", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystemRevListInfo sets the operatingSystemRevListInfo property value. The Operating System Revision List that was loaded during initial boot on the attested device
func (m *DeviceHealthAttestationState) SetOperatingSystemRevListInfo(value *string)() {
    err := m.GetBackingStore().Set("operatingSystemRevListInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetPcr0 sets the pcr0 property value. The measurement that is captured in PCR[0]
func (m *DeviceHealthAttestationState) SetPcr0(value *string)() {
    err := m.GetBackingStore().Set("pcr0", value)
    if err != nil {
        panic(err)
    }
}
// SetPcrHashAlgorithm sets the pcrHashAlgorithm property value. Informational attribute that identifies the HASH algorithm that was used by TPM
func (m *DeviceHealthAttestationState) SetPcrHashAlgorithm(value *string)() {
    err := m.GetBackingStore().Set("pcrHashAlgorithm", value)
    if err != nil {
        panic(err)
    }
}
// SetResetCount sets the resetCount property value. The number of times a PC device has hibernated or resumed
func (m *DeviceHealthAttestationState) SetResetCount(value *int64)() {
    err := m.GetBackingStore().Set("resetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetRestartCount sets the restartCount property value. The number of times a PC device has rebooted
func (m *DeviceHealthAttestationState) SetRestartCount(value *int64)() {
    err := m.GetBackingStore().Set("restartCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSafeMode sets the safeMode property value. Safe mode is a troubleshooting option for Windows that starts your computer in a limited state
func (m *DeviceHealthAttestationState) SetSafeMode(value *string)() {
    err := m.GetBackingStore().Set("safeMode", value)
    if err != nil {
        panic(err)
    }
}
// SetSecureBoot sets the secureBoot property value. When Secure Boot is enabled, the core components must have the correct cryptographic signatures
func (m *DeviceHealthAttestationState) SetSecureBoot(value *string)() {
    err := m.GetBackingStore().Set("secureBoot", value)
    if err != nil {
        panic(err)
    }
}
// SetSecureBootConfigurationPolicyFingerPrint sets the secureBootConfigurationPolicyFingerPrint property value. Fingerprint of the Custom Secure Boot Configuration Policy
func (m *DeviceHealthAttestationState) SetSecureBootConfigurationPolicyFingerPrint(value *string)() {
    err := m.GetBackingStore().Set("secureBootConfigurationPolicyFingerPrint", value)
    if err != nil {
        panic(err)
    }
}
// SetTestSigning sets the testSigning property value. When test signing is allowed, the device does not enforce signature validation during boot
func (m *DeviceHealthAttestationState) SetTestSigning(value *string)() {
    err := m.GetBackingStore().Set("testSigning", value)
    if err != nil {
        panic(err)
    }
}
// SetTpmVersion sets the tpmVersion property value. The security version number of the Boot Application
func (m *DeviceHealthAttestationState) SetTpmVersion(value *string)() {
    err := m.GetBackingStore().Set("tpmVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetVirtualSecureMode sets the virtualSecureMode property value. VSM is a container that protects high value assets from a compromised kernel
func (m *DeviceHealthAttestationState) SetVirtualSecureMode(value *string)() {
    err := m.GetBackingStore().Set("virtualSecureMode", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsPE sets the windowsPE property value. Operating system running with limited services that is used to prepare a computer for Windows
func (m *DeviceHealthAttestationState) SetWindowsPE(value *string)() {
    err := m.GetBackingStore().Set("windowsPE", value)
    if err != nil {
        panic(err)
    }
}
type DeviceHealthAttestationStateable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttestationIdentityKey()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBitLockerStatus()(*string)
    GetBootAppSecurityVersion()(*string)
    GetBootDebugging()(*string)
    GetBootManagerSecurityVersion()(*string)
    GetBootManagerVersion()(*string)
    GetBootRevisionListInfo()(*string)
    GetCodeIntegrity()(*string)
    GetCodeIntegrityCheckVersion()(*string)
    GetCodeIntegrityPolicy()(*string)
    GetContentNamespaceUrl()(*string)
    GetContentVersion()(*string)
    GetDataExcutionPolicy()(*string)
    GetDeviceHealthAttestationStatus()(*string)
    GetEarlyLaunchAntiMalwareDriverProtection()(*string)
    GetHealthAttestationSupportedStatus()(*string)
    GetHealthStatusMismatchInfo()(*string)
    GetIssuedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastUpdateDateTime()(*string)
    GetOdataType()(*string)
    GetOperatingSystemKernelDebugging()(*string)
    GetOperatingSystemRevListInfo()(*string)
    GetPcr0()(*string)
    GetPcrHashAlgorithm()(*string)
    GetResetCount()(*int64)
    GetRestartCount()(*int64)
    GetSafeMode()(*string)
    GetSecureBoot()(*string)
    GetSecureBootConfigurationPolicyFingerPrint()(*string)
    GetTestSigning()(*string)
    GetTpmVersion()(*string)
    GetVirtualSecureMode()(*string)
    GetWindowsPE()(*string)
    SetAttestationIdentityKey(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBitLockerStatus(value *string)()
    SetBootAppSecurityVersion(value *string)()
    SetBootDebugging(value *string)()
    SetBootManagerSecurityVersion(value *string)()
    SetBootManagerVersion(value *string)()
    SetBootRevisionListInfo(value *string)()
    SetCodeIntegrity(value *string)()
    SetCodeIntegrityCheckVersion(value *string)()
    SetCodeIntegrityPolicy(value *string)()
    SetContentNamespaceUrl(value *string)()
    SetContentVersion(value *string)()
    SetDataExcutionPolicy(value *string)()
    SetDeviceHealthAttestationStatus(value *string)()
    SetEarlyLaunchAntiMalwareDriverProtection(value *string)()
    SetHealthAttestationSupportedStatus(value *string)()
    SetHealthStatusMismatchInfo(value *string)()
    SetIssuedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastUpdateDateTime(value *string)()
    SetOdataType(value *string)()
    SetOperatingSystemKernelDebugging(value *string)()
    SetOperatingSystemRevListInfo(value *string)()
    SetPcr0(value *string)()
    SetPcrHashAlgorithm(value *string)()
    SetResetCount(value *int64)()
    SetRestartCount(value *int64)()
    SetSafeMode(value *string)()
    SetSecureBoot(value *string)()
    SetSecureBootConfigurationPolicyFingerPrint(value *string)()
    SetTestSigning(value *string)()
    SetTpmVersion(value *string)()
    SetVirtualSecureMode(value *string)()
    SetWindowsPE(value *string)()
}

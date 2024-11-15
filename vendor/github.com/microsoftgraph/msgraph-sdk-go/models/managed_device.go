package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedDevice devices that are managed or pre-enrolled through Intune
type ManagedDevice struct {
    Entity
}
// NewManagedDevice instantiates a new ManagedDevice and sets the default values.
func NewManagedDevice()(*ManagedDevice) {
    m := &ManagedDevice{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedDeviceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedDeviceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedDevice(), nil
}
// GetActivationLockBypassCode gets the activationLockBypassCode property value. The code that allows the Activation Lock on managed device to be bypassed. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity in LIST call. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetActivationLockBypassCode()(*string) {
    val, err := m.GetBackingStore().Get("activationLockBypassCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAndroidSecurityPatchLevel gets the androidSecurityPatchLevel property value. Android security patch level. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetAndroidSecurityPatchLevel()(*string) {
    val, err := m.GetBackingStore().Get("androidSecurityPatchLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAzureADDeviceId gets the azureADDeviceId property value. The unique identifier for the Azure Active Directory device. Read only. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetAzureADDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("azureADDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAzureADRegistered gets the azureADRegistered property value. Whether the device is Azure Active Directory registered. This property is read-only.
// returns a *bool when successful
func (m *ManagedDevice) GetAzureADRegistered()(*bool) {
    val, err := m.GetBackingStore().Get("azureADRegistered")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetComplianceGracePeriodExpirationDateTime gets the complianceGracePeriodExpirationDateTime property value. The DateTime when device compliance grace period expires. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetComplianceGracePeriodExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("complianceGracePeriodExpirationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetComplianceState gets the complianceState property value. Compliance state.
// returns a *ComplianceState when successful
func (m *ManagedDevice) GetComplianceState()(*ComplianceState) {
    val, err := m.GetBackingStore().Get("complianceState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ComplianceState)
    }
    return nil
}
// GetConfigurationManagerClientEnabledFeatures gets the configurationManagerClientEnabledFeatures property value. ConfigrMgr client enabled features. This property is read-only.
// returns a ConfigurationManagerClientEnabledFeaturesable when successful
func (m *ManagedDevice) GetConfigurationManagerClientEnabledFeatures()(ConfigurationManagerClientEnabledFeaturesable) {
    val, err := m.GetBackingStore().Get("configurationManagerClientEnabledFeatures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConfigurationManagerClientEnabledFeaturesable)
    }
    return nil
}
// GetDeviceActionResults gets the deviceActionResults property value. List of ComplexType deviceActionResult objects. This property is read-only.
// returns a []DeviceActionResultable when successful
func (m *ManagedDevice) GetDeviceActionResults()([]DeviceActionResultable) {
    val, err := m.GetBackingStore().Get("deviceActionResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceActionResultable)
    }
    return nil
}
// GetDeviceCategory gets the deviceCategory property value. Device category
// returns a DeviceCategoryable when successful
func (m *ManagedDevice) GetDeviceCategory()(DeviceCategoryable) {
    val, err := m.GetBackingStore().Get("deviceCategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceCategoryable)
    }
    return nil
}
// GetDeviceCategoryDisplayName gets the deviceCategoryDisplayName property value. Device category display name. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetDeviceCategoryDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("deviceCategoryDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceCompliancePolicyStates gets the deviceCompliancePolicyStates property value. Device compliance policy states for this device.
// returns a []DeviceCompliancePolicyStateable when successful
func (m *ManagedDevice) GetDeviceCompliancePolicyStates()([]DeviceCompliancePolicyStateable) {
    val, err := m.GetBackingStore().Get("deviceCompliancePolicyStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceCompliancePolicyStateable)
    }
    return nil
}
// GetDeviceConfigurationStates gets the deviceConfigurationStates property value. Device configuration states for this device.
// returns a []DeviceConfigurationStateable when successful
func (m *ManagedDevice) GetDeviceConfigurationStates()([]DeviceConfigurationStateable) {
    val, err := m.GetBackingStore().Get("deviceConfigurationStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceConfigurationStateable)
    }
    return nil
}
// GetDeviceEnrollmentType gets the deviceEnrollmentType property value. Possible ways of adding a mobile device to management.
// returns a *DeviceEnrollmentType when successful
func (m *ManagedDevice) GetDeviceEnrollmentType()(*DeviceEnrollmentType) {
    val, err := m.GetBackingStore().Get("deviceEnrollmentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceEnrollmentType)
    }
    return nil
}
// GetDeviceHealthAttestationState gets the deviceHealthAttestationState property value. The device health attestation state. This property is read-only.
// returns a DeviceHealthAttestationStateable when successful
func (m *ManagedDevice) GetDeviceHealthAttestationState()(DeviceHealthAttestationStateable) {
    val, err := m.GetBackingStore().Get("deviceHealthAttestationState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceHealthAttestationStateable)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. Name of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceRegistrationState gets the deviceRegistrationState property value. Device registration status.
// returns a *DeviceRegistrationState when successful
func (m *ManagedDevice) GetDeviceRegistrationState()(*DeviceRegistrationState) {
    val, err := m.GetBackingStore().Get("deviceRegistrationState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceRegistrationState)
    }
    return nil
}
// GetEasActivated gets the easActivated property value. Whether the device is Exchange ActiveSync activated. This property is read-only.
// returns a *bool when successful
func (m *ManagedDevice) GetEasActivated()(*bool) {
    val, err := m.GetBackingStore().Get("easActivated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEasActivationDateTime gets the easActivationDateTime property value. Exchange ActivationSync activation time of the device. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetEasActivationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("easActivationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEasDeviceId gets the easDeviceId property value. Exchange ActiveSync Id of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetEasDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("easDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. Email(s) for the user associated with the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("emailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnrolledDateTime gets the enrolledDateTime property value. Enrollment time of the device. Supports $filter operator 'lt' and 'gt'. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetEnrolledDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("enrolledDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEnrollmentProfileName gets the enrollmentProfileName property value. Name of the enrollment profile assigned to the device. Default value is empty string, indicating no enrollment profile was assgined. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetEnrollmentProfileName()(*string) {
    val, err := m.GetBackingStore().Get("enrollmentProfileName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEthernetMacAddress gets the ethernetMacAddress property value. Indicates Ethernet MAC Address of the device. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity. Individual get call with select query options is needed to retrieve actual values. Example: deviceManagement/managedDevices({managedDeviceId})?$select=ethernetMacAddress Supports: $select. $Search is not supported. Read-only. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetEthernetMacAddress()(*string) {
    val, err := m.GetBackingStore().Get("ethernetMacAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExchangeAccessState gets the exchangeAccessState property value. Device Exchange Access State.
// returns a *DeviceManagementExchangeAccessState when successful
func (m *ManagedDevice) GetExchangeAccessState()(*DeviceManagementExchangeAccessState) {
    val, err := m.GetBackingStore().Get("exchangeAccessState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementExchangeAccessState)
    }
    return nil
}
// GetExchangeAccessStateReason gets the exchangeAccessStateReason property value. Device Exchange Access State Reason.
// returns a *DeviceManagementExchangeAccessStateReason when successful
func (m *ManagedDevice) GetExchangeAccessStateReason()(*DeviceManagementExchangeAccessStateReason) {
    val, err := m.GetBackingStore().Get("exchangeAccessStateReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementExchangeAccessStateReason)
    }
    return nil
}
// GetExchangeLastSuccessfulSyncDateTime gets the exchangeLastSuccessfulSyncDateTime property value. Last time the device contacted Exchange. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetExchangeLastSuccessfulSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("exchangeLastSuccessfulSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ManagedDevice) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activationLockBypassCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivationLockBypassCode(val)
        }
        return nil
    }
    res["androidSecurityPatchLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidSecurityPatchLevel(val)
        }
        return nil
    }
    res["azureADDeviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureADDeviceId(val)
        }
        return nil
    }
    res["azureADRegistered"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureADRegistered(val)
        }
        return nil
    }
    res["complianceGracePeriodExpirationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComplianceGracePeriodExpirationDateTime(val)
        }
        return nil
    }
    res["complianceState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseComplianceState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComplianceState(val.(*ComplianceState))
        }
        return nil
    }
    res["configurationManagerClientEnabledFeatures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConfigurationManagerClientEnabledFeaturesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationManagerClientEnabledFeatures(val.(ConfigurationManagerClientEnabledFeaturesable))
        }
        return nil
    }
    res["deviceActionResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceActionResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceActionResultable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceActionResultable)
                }
            }
            m.SetDeviceActionResults(res)
        }
        return nil
    }
    res["deviceCategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCategory(val.(DeviceCategoryable))
        }
        return nil
    }
    res["deviceCategoryDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCategoryDisplayName(val)
        }
        return nil
    }
    res["deviceCompliancePolicyStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceCompliancePolicyStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceCompliancePolicyStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceCompliancePolicyStateable)
                }
            }
            m.SetDeviceCompliancePolicyStates(res)
        }
        return nil
    }
    res["deviceConfigurationStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceConfigurationStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceConfigurationStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceConfigurationStateable)
                }
            }
            m.SetDeviceConfigurationStates(res)
        }
        return nil
    }
    res["deviceEnrollmentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceEnrollmentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceEnrollmentType(val.(*DeviceEnrollmentType))
        }
        return nil
    }
    res["deviceHealthAttestationState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceHealthAttestationStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceHealthAttestationState(val.(DeviceHealthAttestationStateable))
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
    res["deviceRegistrationState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceRegistrationState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceRegistrationState(val.(*DeviceRegistrationState))
        }
        return nil
    }
    res["easActivated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEasActivated(val)
        }
        return nil
    }
    res["easActivationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEasActivationDateTime(val)
        }
        return nil
    }
    res["easDeviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEasDeviceId(val)
        }
        return nil
    }
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val)
        }
        return nil
    }
    res["enrolledDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnrolledDateTime(val)
        }
        return nil
    }
    res["enrollmentProfileName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnrollmentProfileName(val)
        }
        return nil
    }
    res["ethernetMacAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEthernetMacAddress(val)
        }
        return nil
    }
    res["exchangeAccessState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementExchangeAccessState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeAccessState(val.(*DeviceManagementExchangeAccessState))
        }
        return nil
    }
    res["exchangeAccessStateReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementExchangeAccessStateReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeAccessStateReason(val.(*DeviceManagementExchangeAccessStateReason))
        }
        return nil
    }
    res["exchangeLastSuccessfulSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExchangeLastSuccessfulSyncDateTime(val)
        }
        return nil
    }
    res["freeStorageSpaceInBytes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFreeStorageSpaceInBytes(val)
        }
        return nil
    }
    res["iccid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIccid(val)
        }
        return nil
    }
    res["imei"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImei(val)
        }
        return nil
    }
    res["isEncrypted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEncrypted(val)
        }
        return nil
    }
    res["isSupervised"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSupervised(val)
        }
        return nil
    }
    res["jailBroken"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJailBroken(val)
        }
        return nil
    }
    res["lastSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSyncDateTime(val)
        }
        return nil
    }
    res["logCollectionRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceLogCollectionResponseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceLogCollectionResponseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceLogCollectionResponseable)
                }
            }
            m.SetLogCollectionRequests(res)
        }
        return nil
    }
    res["managedDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedDeviceName(val)
        }
        return nil
    }
    res["managedDeviceOwnerType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedDeviceOwnerType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedDeviceOwnerType(val.(*ManagedDeviceOwnerType))
        }
        return nil
    }
    res["managementAgent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagementAgentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagementAgent(val.(*ManagementAgentType))
        }
        return nil
    }
    res["managementCertificateExpirationDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagementCertificateExpirationDate(val)
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
    res["meid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeid(val)
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
    res["notes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotes(val)
        }
        return nil
    }
    res["operatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperatingSystem(val)
        }
        return nil
    }
    res["osVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsVersion(val)
        }
        return nil
    }
    res["partnerReportedThreatState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedDevicePartnerReportedHealthState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerReportedThreatState(val.(*ManagedDevicePartnerReportedHealthState))
        }
        return nil
    }
    res["phoneNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoneNumber(val)
        }
        return nil
    }
    res["physicalMemoryInBytes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhysicalMemoryInBytes(val)
        }
        return nil
    }
    res["remoteAssistanceSessionErrorDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoteAssistanceSessionErrorDetails(val)
        }
        return nil
    }
    res["remoteAssistanceSessionUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoteAssistanceSessionUrl(val)
        }
        return nil
    }
    res["requireUserEnrollmentApproval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequireUserEnrollmentApproval(val)
        }
        return nil
    }
    res["serialNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSerialNumber(val)
        }
        return nil
    }
    res["subscriberCarrier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubscriberCarrier(val)
        }
        return nil
    }
    res["totalStorageSpaceInBytes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalStorageSpaceInBytes(val)
        }
        return nil
    }
    res["udid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUdid(val)
        }
        return nil
    }
    res["userDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserDisplayName(val)
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
        }
        return nil
    }
    res["userPrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserPrincipalName(val)
        }
        return nil
    }
    res["users"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Userable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Userable)
                }
            }
            m.SetUsers(res)
        }
        return nil
    }
    res["wiFiMacAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiMacAddress(val)
        }
        return nil
    }
    res["windowsProtectionState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsProtectionStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsProtectionState(val.(WindowsProtectionStateable))
        }
        return nil
    }
    return res
}
// GetFreeStorageSpaceInBytes gets the freeStorageSpaceInBytes property value. Free Storage in Bytes. Default value is 0. Read-only. This property is read-only.
// returns a *int64 when successful
func (m *ManagedDevice) GetFreeStorageSpaceInBytes()(*int64) {
    val, err := m.GetBackingStore().Get("freeStorageSpaceInBytes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetIccid gets the iccid property value. Integrated Circuit Card Identifier, it is A SIM card's unique identification number. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetIccid()(*string) {
    val, err := m.GetBackingStore().Get("iccid")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetImei gets the imei property value. IMEI. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetImei()(*string) {
    val, err := m.GetBackingStore().Get("imei")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsEncrypted gets the isEncrypted property value. Device encryption status. This property is read-only.
// returns a *bool when successful
func (m *ManagedDevice) GetIsEncrypted()(*bool) {
    val, err := m.GetBackingStore().Get("isEncrypted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSupervised gets the isSupervised property value. Device supervised status. This property is read-only.
// returns a *bool when successful
func (m *ManagedDevice) GetIsSupervised()(*bool) {
    val, err := m.GetBackingStore().Get("isSupervised")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJailBroken gets the jailBroken property value. Whether the device is jail broken or rooted. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetJailBroken()(*string) {
    val, err := m.GetBackingStore().Get("jailBroken")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastSyncDateTime gets the lastSyncDateTime property value. The date and time that the device last completed a successful sync with Intune. Supports $filter operator 'lt' and 'gt'. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLogCollectionRequests gets the logCollectionRequests property value. List of log collection requests
// returns a []DeviceLogCollectionResponseable when successful
func (m *ManagedDevice) GetLogCollectionRequests()([]DeviceLogCollectionResponseable) {
    val, err := m.GetBackingStore().Get("logCollectionRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceLogCollectionResponseable)
    }
    return nil
}
// GetManagedDeviceName gets the managedDeviceName property value. Automatically generated name to identify a device. Can be overwritten to a user friendly name.
// returns a *string when successful
func (m *ManagedDevice) GetManagedDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("managedDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetManagedDeviceOwnerType gets the managedDeviceOwnerType property value. Owner type of device.
// returns a *ManagedDeviceOwnerType when successful
func (m *ManagedDevice) GetManagedDeviceOwnerType()(*ManagedDeviceOwnerType) {
    val, err := m.GetBackingStore().Get("managedDeviceOwnerType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedDeviceOwnerType)
    }
    return nil
}
// GetManagementAgent gets the managementAgent property value. The managementAgent property
// returns a *ManagementAgentType when successful
func (m *ManagedDevice) GetManagementAgent()(*ManagementAgentType) {
    val, err := m.GetBackingStore().Get("managementAgent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagementAgentType)
    }
    return nil
}
// GetManagementCertificateExpirationDate gets the managementCertificateExpirationDate property value. Reports device management certificate expiration date. This property is read-only.
// returns a *Time when successful
func (m *ManagedDevice) GetManagementCertificateExpirationDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("managementCertificateExpirationDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetManufacturer gets the manufacturer property value. Manufacturer of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetManufacturer()(*string) {
    val, err := m.GetBackingStore().Get("manufacturer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMeid gets the meid property value. MEID. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetMeid()(*string) {
    val, err := m.GetBackingStore().Get("meid")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetModel gets the model property value. Model of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotes gets the notes property value. Notes on the device created by IT Admin. Default is null. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported.
// returns a *string when successful
func (m *ManagedDevice) GetNotes()(*string) {
    val, err := m.GetBackingStore().Get("notes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperatingSystem gets the operatingSystem property value. Operating system of the device. Windows, iOS, etc. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetOperatingSystem()(*string) {
    val, err := m.GetBackingStore().Get("operatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOsVersion gets the osVersion property value. Operating system version of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetOsVersion()(*string) {
    val, err := m.GetBackingStore().Get("osVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPartnerReportedThreatState gets the partnerReportedThreatState property value. Available health states for the Device Health API
// returns a *ManagedDevicePartnerReportedHealthState when successful
func (m *ManagedDevice) GetPartnerReportedThreatState()(*ManagedDevicePartnerReportedHealthState) {
    val, err := m.GetBackingStore().Get("partnerReportedThreatState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedDevicePartnerReportedHealthState)
    }
    return nil
}
// GetPhoneNumber gets the phoneNumber property value. Phone number of the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetPhoneNumber()(*string) {
    val, err := m.GetBackingStore().Get("phoneNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhysicalMemoryInBytes gets the physicalMemoryInBytes property value. Total Memory in Bytes. Default is 0. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. Read-only. This property is read-only.
// returns a *int64 when successful
func (m *ManagedDevice) GetPhysicalMemoryInBytes()(*int64) {
    val, err := m.GetBackingStore().Get("physicalMemoryInBytes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetRemoteAssistanceSessionErrorDetails gets the remoteAssistanceSessionErrorDetails property value. An error string that identifies issues when creating Remote Assistance session objects. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetRemoteAssistanceSessionErrorDetails()(*string) {
    val, err := m.GetBackingStore().Get("remoteAssistanceSessionErrorDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemoteAssistanceSessionUrl gets the remoteAssistanceSessionUrl property value. Url that allows a Remote Assistance session to be established with the device. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetRemoteAssistanceSessionUrl()(*string) {
    val, err := m.GetBackingStore().Get("remoteAssistanceSessionUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRequireUserEnrollmentApproval gets the requireUserEnrollmentApproval property value. Reports if the managed iOS device is user approval enrollment. This property is read-only.
// returns a *bool when successful
func (m *ManagedDevice) GetRequireUserEnrollmentApproval()(*bool) {
    val, err := m.GetBackingStore().Get("requireUserEnrollmentApproval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSerialNumber gets the serialNumber property value. SerialNumber. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetSerialNumber()(*string) {
    val, err := m.GetBackingStore().Get("serialNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubscriberCarrier gets the subscriberCarrier property value. Subscriber Carrier. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetSubscriberCarrier()(*string) {
    val, err := m.GetBackingStore().Get("subscriberCarrier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalStorageSpaceInBytes gets the totalStorageSpaceInBytes property value. Total Storage in Bytes. This property is read-only.
// returns a *int64 when successful
func (m *ManagedDevice) GetTotalStorageSpaceInBytes()(*int64) {
    val, err := m.GetBackingStore().Get("totalStorageSpaceInBytes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUdid gets the udid property value. Unique Device Identifier for iOS and macOS devices. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetUdid()(*string) {
    val, err := m.GetBackingStore().Get("udid")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserDisplayName gets the userDisplayName property value. User display name. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetUserDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("userDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. Unique Identifier for the user associated with the device. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserPrincipalName gets the userPrincipalName property value. Device user principal name. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetUserPrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("userPrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUsers gets the users property value. The primary users associated with the managed device.
// returns a []Userable when successful
func (m *ManagedDevice) GetUsers()([]Userable) {
    val, err := m.GetBackingStore().Get("users")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Userable)
    }
    return nil
}
// GetWiFiMacAddress gets the wiFiMacAddress property value. Wi-Fi MAC. This property is read-only.
// returns a *string when successful
func (m *ManagedDevice) GetWiFiMacAddress()(*string) {
    val, err := m.GetBackingStore().Get("wiFiMacAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWindowsProtectionState gets the windowsProtectionState property value. The device protection status. This property is read-only.
// returns a WindowsProtectionStateable when successful
func (m *ManagedDevice) GetWindowsProtectionState()(WindowsProtectionStateable) {
    val, err := m.GetBackingStore().Get("windowsProtectionState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsProtectionStateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedDevice) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetComplianceState() != nil {
        cast := (*m.GetComplianceState()).String()
        err = writer.WriteStringValue("complianceState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceCategory", m.GetDeviceCategory())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCompliancePolicyStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceCompliancePolicyStates()))
        for i, v := range m.GetDeviceCompliancePolicyStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceCompliancePolicyStates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceConfigurationStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceConfigurationStates()))
        for i, v := range m.GetDeviceConfigurationStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceConfigurationStates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceEnrollmentType() != nil {
        cast := (*m.GetDeviceEnrollmentType()).String()
        err = writer.WriteStringValue("deviceEnrollmentType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceRegistrationState() != nil {
        cast := (*m.GetDeviceRegistrationState()).String()
        err = writer.WriteStringValue("deviceRegistrationState", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeAccessState() != nil {
        cast := (*m.GetExchangeAccessState()).String()
        err = writer.WriteStringValue("exchangeAccessState", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeAccessStateReason() != nil {
        cast := (*m.GetExchangeAccessStateReason()).String()
        err = writer.WriteStringValue("exchangeAccessStateReason", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetLogCollectionRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLogCollectionRequests()))
        for i, v := range m.GetLogCollectionRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("logCollectionRequests", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managedDeviceName", m.GetManagedDeviceName())
        if err != nil {
            return err
        }
    }
    if m.GetManagedDeviceOwnerType() != nil {
        cast := (*m.GetManagedDeviceOwnerType()).String()
        err = writer.WriteStringValue("managedDeviceOwnerType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetManagementAgent() != nil {
        cast := (*m.GetManagementAgent()).String()
        err = writer.WriteStringValue("managementAgent", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notes", m.GetNotes())
        if err != nil {
            return err
        }
    }
    if m.GetPartnerReportedThreatState() != nil {
        cast := (*m.GetPartnerReportedThreatState()).String()
        err = writer.WriteStringValue("partnerReportedThreatState", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetUsers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUsers()))
        for i, v := range m.GetUsers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("users", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("windowsProtectionState", m.GetWindowsProtectionState())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivationLockBypassCode sets the activationLockBypassCode property value. The code that allows the Activation Lock on managed device to be bypassed. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity in LIST call. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
func (m *ManagedDevice) SetActivationLockBypassCode(value *string)() {
    err := m.GetBackingStore().Set("activationLockBypassCode", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidSecurityPatchLevel sets the androidSecurityPatchLevel property value. Android security patch level. This property is read-only.
func (m *ManagedDevice) SetAndroidSecurityPatchLevel(value *string)() {
    err := m.GetBackingStore().Set("androidSecurityPatchLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureADDeviceId sets the azureADDeviceId property value. The unique identifier for the Azure Active Directory device. Read only. This property is read-only.
func (m *ManagedDevice) SetAzureADDeviceId(value *string)() {
    err := m.GetBackingStore().Set("azureADDeviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureADRegistered sets the azureADRegistered property value. Whether the device is Azure Active Directory registered. This property is read-only.
func (m *ManagedDevice) SetAzureADRegistered(value *bool)() {
    err := m.GetBackingStore().Set("azureADRegistered", value)
    if err != nil {
        panic(err)
    }
}
// SetComplianceGracePeriodExpirationDateTime sets the complianceGracePeriodExpirationDateTime property value. The DateTime when device compliance grace period expires. This property is read-only.
func (m *ManagedDevice) SetComplianceGracePeriodExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("complianceGracePeriodExpirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetComplianceState sets the complianceState property value. Compliance state.
func (m *ManagedDevice) SetComplianceState(value *ComplianceState)() {
    err := m.GetBackingStore().Set("complianceState", value)
    if err != nil {
        panic(err)
    }
}
// SetConfigurationManagerClientEnabledFeatures sets the configurationManagerClientEnabledFeatures property value. ConfigrMgr client enabled features. This property is read-only.
func (m *ManagedDevice) SetConfigurationManagerClientEnabledFeatures(value ConfigurationManagerClientEnabledFeaturesable)() {
    err := m.GetBackingStore().Set("configurationManagerClientEnabledFeatures", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceActionResults sets the deviceActionResults property value. List of ComplexType deviceActionResult objects. This property is read-only.
func (m *ManagedDevice) SetDeviceActionResults(value []DeviceActionResultable)() {
    err := m.GetBackingStore().Set("deviceActionResults", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCategory sets the deviceCategory property value. Device category
func (m *ManagedDevice) SetDeviceCategory(value DeviceCategoryable)() {
    err := m.GetBackingStore().Set("deviceCategory", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCategoryDisplayName sets the deviceCategoryDisplayName property value. Device category display name. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
func (m *ManagedDevice) SetDeviceCategoryDisplayName(value *string)() {
    err := m.GetBackingStore().Set("deviceCategoryDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCompliancePolicyStates sets the deviceCompliancePolicyStates property value. Device compliance policy states for this device.
func (m *ManagedDevice) SetDeviceCompliancePolicyStates(value []DeviceCompliancePolicyStateable)() {
    err := m.GetBackingStore().Set("deviceCompliancePolicyStates", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceConfigurationStates sets the deviceConfigurationStates property value. Device configuration states for this device.
func (m *ManagedDevice) SetDeviceConfigurationStates(value []DeviceConfigurationStateable)() {
    err := m.GetBackingStore().Set("deviceConfigurationStates", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceEnrollmentType sets the deviceEnrollmentType property value. Possible ways of adding a mobile device to management.
func (m *ManagedDevice) SetDeviceEnrollmentType(value *DeviceEnrollmentType)() {
    err := m.GetBackingStore().Set("deviceEnrollmentType", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceHealthAttestationState sets the deviceHealthAttestationState property value. The device health attestation state. This property is read-only.
func (m *ManagedDevice) SetDeviceHealthAttestationState(value DeviceHealthAttestationStateable)() {
    err := m.GetBackingStore().Set("deviceHealthAttestationState", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. Name of the device. This property is read-only.
func (m *ManagedDevice) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceRegistrationState sets the deviceRegistrationState property value. Device registration status.
func (m *ManagedDevice) SetDeviceRegistrationState(value *DeviceRegistrationState)() {
    err := m.GetBackingStore().Set("deviceRegistrationState", value)
    if err != nil {
        panic(err)
    }
}
// SetEasActivated sets the easActivated property value. Whether the device is Exchange ActiveSync activated. This property is read-only.
func (m *ManagedDevice) SetEasActivated(value *bool)() {
    err := m.GetBackingStore().Set("easActivated", value)
    if err != nil {
        panic(err)
    }
}
// SetEasActivationDateTime sets the easActivationDateTime property value. Exchange ActivationSync activation time of the device. This property is read-only.
func (m *ManagedDevice) SetEasActivationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("easActivationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEasDeviceId sets the easDeviceId property value. Exchange ActiveSync Id of the device. This property is read-only.
func (m *ManagedDevice) SetEasDeviceId(value *string)() {
    err := m.GetBackingStore().Set("easDeviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. Email(s) for the user associated with the device. This property is read-only.
func (m *ManagedDevice) SetEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetEnrolledDateTime sets the enrolledDateTime property value. Enrollment time of the device. Supports $filter operator 'lt' and 'gt'. This property is read-only.
func (m *ManagedDevice) SetEnrolledDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("enrolledDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEnrollmentProfileName sets the enrollmentProfileName property value. Name of the enrollment profile assigned to the device. Default value is empty string, indicating no enrollment profile was assgined. This property is read-only.
func (m *ManagedDevice) SetEnrollmentProfileName(value *string)() {
    err := m.GetBackingStore().Set("enrollmentProfileName", value)
    if err != nil {
        panic(err)
    }
}
// SetEthernetMacAddress sets the ethernetMacAddress property value. Indicates Ethernet MAC Address of the device. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity. Individual get call with select query options is needed to retrieve actual values. Example: deviceManagement/managedDevices({managedDeviceId})?$select=ethernetMacAddress Supports: $select. $Search is not supported. Read-only. This property is read-only.
func (m *ManagedDevice) SetEthernetMacAddress(value *string)() {
    err := m.GetBackingStore().Set("ethernetMacAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeAccessState sets the exchangeAccessState property value. Device Exchange Access State.
func (m *ManagedDevice) SetExchangeAccessState(value *DeviceManagementExchangeAccessState)() {
    err := m.GetBackingStore().Set("exchangeAccessState", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeAccessStateReason sets the exchangeAccessStateReason property value. Device Exchange Access State Reason.
func (m *ManagedDevice) SetExchangeAccessStateReason(value *DeviceManagementExchangeAccessStateReason)() {
    err := m.GetBackingStore().Set("exchangeAccessStateReason", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeLastSuccessfulSyncDateTime sets the exchangeLastSuccessfulSyncDateTime property value. Last time the device contacted Exchange. This property is read-only.
func (m *ManagedDevice) SetExchangeLastSuccessfulSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("exchangeLastSuccessfulSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFreeStorageSpaceInBytes sets the freeStorageSpaceInBytes property value. Free Storage in Bytes. Default value is 0. Read-only. This property is read-only.
func (m *ManagedDevice) SetFreeStorageSpaceInBytes(value *int64)() {
    err := m.GetBackingStore().Set("freeStorageSpaceInBytes", value)
    if err != nil {
        panic(err)
    }
}
// SetIccid sets the iccid property value. Integrated Circuit Card Identifier, it is A SIM card's unique identification number. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
func (m *ManagedDevice) SetIccid(value *string)() {
    err := m.GetBackingStore().Set("iccid", value)
    if err != nil {
        panic(err)
    }
}
// SetImei sets the imei property value. IMEI. This property is read-only.
func (m *ManagedDevice) SetImei(value *string)() {
    err := m.GetBackingStore().Set("imei", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEncrypted sets the isEncrypted property value. Device encryption status. This property is read-only.
func (m *ManagedDevice) SetIsEncrypted(value *bool)() {
    err := m.GetBackingStore().Set("isEncrypted", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSupervised sets the isSupervised property value. Device supervised status. This property is read-only.
func (m *ManagedDevice) SetIsSupervised(value *bool)() {
    err := m.GetBackingStore().Set("isSupervised", value)
    if err != nil {
        panic(err)
    }
}
// SetJailBroken sets the jailBroken property value. Whether the device is jail broken or rooted. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
func (m *ManagedDevice) SetJailBroken(value *string)() {
    err := m.GetBackingStore().Set("jailBroken", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSyncDateTime sets the lastSyncDateTime property value. The date and time that the device last completed a successful sync with Intune. Supports $filter operator 'lt' and 'gt'. This property is read-only.
func (m *ManagedDevice) SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLogCollectionRequests sets the logCollectionRequests property value. List of log collection requests
func (m *ManagedDevice) SetLogCollectionRequests(value []DeviceLogCollectionResponseable)() {
    err := m.GetBackingStore().Set("logCollectionRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDeviceName sets the managedDeviceName property value. Automatically generated name to identify a device. Can be overwritten to a user friendly name.
func (m *ManagedDevice) SetManagedDeviceName(value *string)() {
    err := m.GetBackingStore().Set("managedDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDeviceOwnerType sets the managedDeviceOwnerType property value. Owner type of device.
func (m *ManagedDevice) SetManagedDeviceOwnerType(value *ManagedDeviceOwnerType)() {
    err := m.GetBackingStore().Set("managedDeviceOwnerType", value)
    if err != nil {
        panic(err)
    }
}
// SetManagementAgent sets the managementAgent property value. The managementAgent property
func (m *ManagedDevice) SetManagementAgent(value *ManagementAgentType)() {
    err := m.GetBackingStore().Set("managementAgent", value)
    if err != nil {
        panic(err)
    }
}
// SetManagementCertificateExpirationDate sets the managementCertificateExpirationDate property value. Reports device management certificate expiration date. This property is read-only.
func (m *ManagedDevice) SetManagementCertificateExpirationDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("managementCertificateExpirationDate", value)
    if err != nil {
        panic(err)
    }
}
// SetManufacturer sets the manufacturer property value. Manufacturer of the device. This property is read-only.
func (m *ManagedDevice) SetManufacturer(value *string)() {
    err := m.GetBackingStore().Set("manufacturer", value)
    if err != nil {
        panic(err)
    }
}
// SetMeid sets the meid property value. MEID. This property is read-only.
func (m *ManagedDevice) SetMeid(value *string)() {
    err := m.GetBackingStore().Set("meid", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. Model of the device. This property is read-only.
func (m *ManagedDevice) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
// SetNotes sets the notes property value. Notes on the device created by IT Admin. Default is null. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported.
func (m *ManagedDevice) SetNotes(value *string)() {
    err := m.GetBackingStore().Set("notes", value)
    if err != nil {
        panic(err)
    }
}
// SetOperatingSystem sets the operatingSystem property value. Operating system of the device. Windows, iOS, etc. This property is read-only.
func (m *ManagedDevice) SetOperatingSystem(value *string)() {
    err := m.GetBackingStore().Set("operatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetOsVersion sets the osVersion property value. Operating system version of the device. This property is read-only.
func (m *ManagedDevice) SetOsVersion(value *string)() {
    err := m.GetBackingStore().Set("osVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerReportedThreatState sets the partnerReportedThreatState property value. Available health states for the Device Health API
func (m *ManagedDevice) SetPartnerReportedThreatState(value *ManagedDevicePartnerReportedHealthState)() {
    err := m.GetBackingStore().Set("partnerReportedThreatState", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoneNumber sets the phoneNumber property value. Phone number of the device. This property is read-only.
func (m *ManagedDevice) SetPhoneNumber(value *string)() {
    err := m.GetBackingStore().Set("phoneNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetPhysicalMemoryInBytes sets the physicalMemoryInBytes property value. Total Memory in Bytes. Default is 0. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. Read-only. This property is read-only.
func (m *ManagedDevice) SetPhysicalMemoryInBytes(value *int64)() {
    err := m.GetBackingStore().Set("physicalMemoryInBytes", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoteAssistanceSessionErrorDetails sets the remoteAssistanceSessionErrorDetails property value. An error string that identifies issues when creating Remote Assistance session objects. This property is read-only.
func (m *ManagedDevice) SetRemoteAssistanceSessionErrorDetails(value *string)() {
    err := m.GetBackingStore().Set("remoteAssistanceSessionErrorDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoteAssistanceSessionUrl sets the remoteAssistanceSessionUrl property value. Url that allows a Remote Assistance session to be established with the device. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. This property is read-only.
func (m *ManagedDevice) SetRemoteAssistanceSessionUrl(value *string)() {
    err := m.GetBackingStore().Set("remoteAssistanceSessionUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetRequireUserEnrollmentApproval sets the requireUserEnrollmentApproval property value. Reports if the managed iOS device is user approval enrollment. This property is read-only.
func (m *ManagedDevice) SetRequireUserEnrollmentApproval(value *bool)() {
    err := m.GetBackingStore().Set("requireUserEnrollmentApproval", value)
    if err != nil {
        panic(err)
    }
}
// SetSerialNumber sets the serialNumber property value. SerialNumber. This property is read-only.
func (m *ManagedDevice) SetSerialNumber(value *string)() {
    err := m.GetBackingStore().Set("serialNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetSubscriberCarrier sets the subscriberCarrier property value. Subscriber Carrier. This property is read-only.
func (m *ManagedDevice) SetSubscriberCarrier(value *string)() {
    err := m.GetBackingStore().Set("subscriberCarrier", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalStorageSpaceInBytes sets the totalStorageSpaceInBytes property value. Total Storage in Bytes. This property is read-only.
func (m *ManagedDevice) SetTotalStorageSpaceInBytes(value *int64)() {
    err := m.GetBackingStore().Set("totalStorageSpaceInBytes", value)
    if err != nil {
        panic(err)
    }
}
// SetUdid sets the udid property value. Unique Device Identifier for iOS and macOS devices. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
func (m *ManagedDevice) SetUdid(value *string)() {
    err := m.GetBackingStore().Set("udid", value)
    if err != nil {
        panic(err)
    }
}
// SetUserDisplayName sets the userDisplayName property value. User display name. This property is read-only.
func (m *ManagedDevice) SetUserDisplayName(value *string)() {
    err := m.GetBackingStore().Set("userDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. Unique Identifier for the user associated with the device. This property is read-only.
func (m *ManagedDevice) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
// SetUserPrincipalName sets the userPrincipalName property value. Device user principal name. This property is read-only.
func (m *ManagedDevice) SetUserPrincipalName(value *string)() {
    err := m.GetBackingStore().Set("userPrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetUsers sets the users property value. The primary users associated with the managed device.
func (m *ManagedDevice) SetUsers(value []Userable)() {
    err := m.GetBackingStore().Set("users", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiMacAddress sets the wiFiMacAddress property value. Wi-Fi MAC. This property is read-only.
func (m *ManagedDevice) SetWiFiMacAddress(value *string)() {
    err := m.GetBackingStore().Set("wiFiMacAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsProtectionState sets the windowsProtectionState property value. The device protection status. This property is read-only.
func (m *ManagedDevice) SetWindowsProtectionState(value WindowsProtectionStateable)() {
    err := m.GetBackingStore().Set("windowsProtectionState", value)
    if err != nil {
        panic(err)
    }
}
type ManagedDeviceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivationLockBypassCode()(*string)
    GetAndroidSecurityPatchLevel()(*string)
    GetAzureADDeviceId()(*string)
    GetAzureADRegistered()(*bool)
    GetComplianceGracePeriodExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetComplianceState()(*ComplianceState)
    GetConfigurationManagerClientEnabledFeatures()(ConfigurationManagerClientEnabledFeaturesable)
    GetDeviceActionResults()([]DeviceActionResultable)
    GetDeviceCategory()(DeviceCategoryable)
    GetDeviceCategoryDisplayName()(*string)
    GetDeviceCompliancePolicyStates()([]DeviceCompliancePolicyStateable)
    GetDeviceConfigurationStates()([]DeviceConfigurationStateable)
    GetDeviceEnrollmentType()(*DeviceEnrollmentType)
    GetDeviceHealthAttestationState()(DeviceHealthAttestationStateable)
    GetDeviceName()(*string)
    GetDeviceRegistrationState()(*DeviceRegistrationState)
    GetEasActivated()(*bool)
    GetEasActivationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEasDeviceId()(*string)
    GetEmailAddress()(*string)
    GetEnrolledDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEnrollmentProfileName()(*string)
    GetEthernetMacAddress()(*string)
    GetExchangeAccessState()(*DeviceManagementExchangeAccessState)
    GetExchangeAccessStateReason()(*DeviceManagementExchangeAccessStateReason)
    GetExchangeLastSuccessfulSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFreeStorageSpaceInBytes()(*int64)
    GetIccid()(*string)
    GetImei()(*string)
    GetIsEncrypted()(*bool)
    GetIsSupervised()(*bool)
    GetJailBroken()(*string)
    GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLogCollectionRequests()([]DeviceLogCollectionResponseable)
    GetManagedDeviceName()(*string)
    GetManagedDeviceOwnerType()(*ManagedDeviceOwnerType)
    GetManagementAgent()(*ManagementAgentType)
    GetManagementCertificateExpirationDate()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetManufacturer()(*string)
    GetMeid()(*string)
    GetModel()(*string)
    GetNotes()(*string)
    GetOperatingSystem()(*string)
    GetOsVersion()(*string)
    GetPartnerReportedThreatState()(*ManagedDevicePartnerReportedHealthState)
    GetPhoneNumber()(*string)
    GetPhysicalMemoryInBytes()(*int64)
    GetRemoteAssistanceSessionErrorDetails()(*string)
    GetRemoteAssistanceSessionUrl()(*string)
    GetRequireUserEnrollmentApproval()(*bool)
    GetSerialNumber()(*string)
    GetSubscriberCarrier()(*string)
    GetTotalStorageSpaceInBytes()(*int64)
    GetUdid()(*string)
    GetUserDisplayName()(*string)
    GetUserId()(*string)
    GetUserPrincipalName()(*string)
    GetUsers()([]Userable)
    GetWiFiMacAddress()(*string)
    GetWindowsProtectionState()(WindowsProtectionStateable)
    SetActivationLockBypassCode(value *string)()
    SetAndroidSecurityPatchLevel(value *string)()
    SetAzureADDeviceId(value *string)()
    SetAzureADRegistered(value *bool)()
    SetComplianceGracePeriodExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetComplianceState(value *ComplianceState)()
    SetConfigurationManagerClientEnabledFeatures(value ConfigurationManagerClientEnabledFeaturesable)()
    SetDeviceActionResults(value []DeviceActionResultable)()
    SetDeviceCategory(value DeviceCategoryable)()
    SetDeviceCategoryDisplayName(value *string)()
    SetDeviceCompliancePolicyStates(value []DeviceCompliancePolicyStateable)()
    SetDeviceConfigurationStates(value []DeviceConfigurationStateable)()
    SetDeviceEnrollmentType(value *DeviceEnrollmentType)()
    SetDeviceHealthAttestationState(value DeviceHealthAttestationStateable)()
    SetDeviceName(value *string)()
    SetDeviceRegistrationState(value *DeviceRegistrationState)()
    SetEasActivated(value *bool)()
    SetEasActivationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEasDeviceId(value *string)()
    SetEmailAddress(value *string)()
    SetEnrolledDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEnrollmentProfileName(value *string)()
    SetEthernetMacAddress(value *string)()
    SetExchangeAccessState(value *DeviceManagementExchangeAccessState)()
    SetExchangeAccessStateReason(value *DeviceManagementExchangeAccessStateReason)()
    SetExchangeLastSuccessfulSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFreeStorageSpaceInBytes(value *int64)()
    SetIccid(value *string)()
    SetImei(value *string)()
    SetIsEncrypted(value *bool)()
    SetIsSupervised(value *bool)()
    SetJailBroken(value *string)()
    SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLogCollectionRequests(value []DeviceLogCollectionResponseable)()
    SetManagedDeviceName(value *string)()
    SetManagedDeviceOwnerType(value *ManagedDeviceOwnerType)()
    SetManagementAgent(value *ManagementAgentType)()
    SetManagementCertificateExpirationDate(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetManufacturer(value *string)()
    SetMeid(value *string)()
    SetModel(value *string)()
    SetNotes(value *string)()
    SetOperatingSystem(value *string)()
    SetOsVersion(value *string)()
    SetPartnerReportedThreatState(value *ManagedDevicePartnerReportedHealthState)()
    SetPhoneNumber(value *string)()
    SetPhysicalMemoryInBytes(value *int64)()
    SetRemoteAssistanceSessionErrorDetails(value *string)()
    SetRemoteAssistanceSessionUrl(value *string)()
    SetRequireUserEnrollmentApproval(value *bool)()
    SetSerialNumber(value *string)()
    SetSubscriberCarrier(value *string)()
    SetTotalStorageSpaceInBytes(value *int64)()
    SetUdid(value *string)()
    SetUserDisplayName(value *string)()
    SetUserId(value *string)()
    SetUserPrincipalName(value *string)()
    SetUsers(value []Userable)()
    SetWiFiMacAddress(value *string)()
    SetWindowsProtectionState(value WindowsProtectionStateable)()
}

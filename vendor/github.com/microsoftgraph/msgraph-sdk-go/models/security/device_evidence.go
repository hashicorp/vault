package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceEvidence struct {
    AlertEvidence
}
// NewDeviceEvidence instantiates a new DeviceEvidence and sets the default values.
func NewDeviceEvidence()(*DeviceEvidence) {
    m := &DeviceEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.deviceEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDeviceEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceEvidence(), nil
}
// GetAzureAdDeviceId gets the azureAdDeviceId property value. A unique identifier assigned to a device by Microsoft Entra ID when device is Microsoft Entra joined.
// returns a *string when successful
func (m *DeviceEvidence) GetAzureAdDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("azureAdDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDefenderAvStatus gets the defenderAvStatus property value. State of the Defender AntiMalware engine. The possible values are: notReporting, disabled, notUpdated, updated, unknown, notSupported, unknownFutureValue.
// returns a *DefenderAvStatus when successful
func (m *DeviceEvidence) GetDefenderAvStatus()(*DefenderAvStatus) {
    val, err := m.GetBackingStore().Get("defenderAvStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefenderAvStatus)
    }
    return nil
}
// GetDeviceDnsName gets the deviceDnsName property value. The fully qualified domain name (FQDN) for the device.
// returns a *string when successful
func (m *DeviceEvidence) GetDeviceDnsName()(*string) {
    val, err := m.GetBackingStore().Get("deviceDnsName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDnsDomain gets the dnsDomain property value. The DNS domain that this computer belongs to. A sequence of labels separated by dots.
// returns a *string when successful
func (m *DeviceEvidence) GetDnsDomain()(*string) {
    val, err := m.GetBackingStore().Get("dnsDomain")
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
func (m *DeviceEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["azureAdDeviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureAdDeviceId(val)
        }
        return nil
    }
    res["defenderAvStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefenderAvStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderAvStatus(val.(*DefenderAvStatus))
        }
        return nil
    }
    res["deviceDnsName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceDnsName(val)
        }
        return nil
    }
    res["dnsDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDnsDomain(val)
        }
        return nil
    }
    res["firstSeenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirstSeenDateTime(val)
        }
        return nil
    }
    res["healthStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceHealthStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHealthStatus(val.(*DeviceHealthStatus))
        }
        return nil
    }
    res["hostName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHostName(val)
        }
        return nil
    }
    res["ipInterfaces"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetIpInterfaces(res)
        }
        return nil
    }
    res["lastExternalIpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastExternalIpAddress(val)
        }
        return nil
    }
    res["lastIpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastIpAddress(val)
        }
        return nil
    }
    res["loggedOnUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLoggedOnUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LoggedOnUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LoggedOnUserable)
                }
            }
            m.SetLoggedOnUsers(res)
        }
        return nil
    }
    res["mdeDeviceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMdeDeviceId(val)
        }
        return nil
    }
    res["ntDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNtDomain(val)
        }
        return nil
    }
    res["onboardingStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnboardingStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnboardingStatus(val.(*OnboardingStatus))
        }
        return nil
    }
    res["osBuild"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsBuild(val)
        }
        return nil
    }
    res["osPlatform"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOsPlatform(val)
        }
        return nil
    }
    res["rbacGroupId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRbacGroupId(val)
        }
        return nil
    }
    res["rbacGroupName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRbacGroupName(val)
        }
        return nil
    }
    res["riskScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceRiskScore)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRiskScore(val.(*DeviceRiskScore))
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    res["vmMetadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVmMetadataFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVmMetadata(val.(VmMetadataable))
        }
        return nil
    }
    return res
}
// GetFirstSeenDateTime gets the firstSeenDateTime property value. The date and time when the device was first seen.
// returns a *Time when successful
func (m *DeviceEvidence) GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("firstSeenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetHealthStatus gets the healthStatus property value. The health state of the device. The possible values are: active, inactive, impairedCommunication, noSensorData, noSensorDataImpairedCommunication, unknown, unknownFutureValue.
// returns a *DeviceHealthStatus when successful
func (m *DeviceEvidence) GetHealthStatus()(*DeviceHealthStatus) {
    val, err := m.GetBackingStore().Get("healthStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceHealthStatus)
    }
    return nil
}
// GetHostName gets the hostName property value. The hostname without the domain suffix.
// returns a *string when successful
func (m *DeviceEvidence) GetHostName()(*string) {
    val, err := m.GetBackingStore().Get("hostName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIpInterfaces gets the ipInterfaces property value. Ip interfaces of the device during the time of the alert.
// returns a []string when successful
func (m *DeviceEvidence) GetIpInterfaces()([]string) {
    val, err := m.GetBackingStore().Get("ipInterfaces")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetLastExternalIpAddress gets the lastExternalIpAddress property value. The lastExternalIpAddress property
// returns a *string when successful
func (m *DeviceEvidence) GetLastExternalIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("lastExternalIpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastIpAddress gets the lastIpAddress property value. The lastIpAddress property
// returns a *string when successful
func (m *DeviceEvidence) GetLastIpAddress()(*string) {
    val, err := m.GetBackingStore().Get("lastIpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLoggedOnUsers gets the loggedOnUsers property value. Users that were logged on the machine during the time of the alert.
// returns a []LoggedOnUserable when successful
func (m *DeviceEvidence) GetLoggedOnUsers()([]LoggedOnUserable) {
    val, err := m.GetBackingStore().Get("loggedOnUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LoggedOnUserable)
    }
    return nil
}
// GetMdeDeviceId gets the mdeDeviceId property value. A unique identifier assigned to a device by Microsoft Defender for Endpoint.
// returns a *string when successful
func (m *DeviceEvidence) GetMdeDeviceId()(*string) {
    val, err := m.GetBackingStore().Get("mdeDeviceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNtDomain gets the ntDomain property value. A logical grouping of computers within a Microsoft Windows network.
// returns a *string when successful
func (m *DeviceEvidence) GetNtDomain()(*string) {
    val, err := m.GetBackingStore().Get("ntDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnboardingStatus gets the onboardingStatus property value. The status of the machine onboarding to Microsoft Defender for Endpoint. The possible values are: insufficientInfo, onboarded, canBeOnboarded, unsupported, unknownFutureValue.
// returns a *OnboardingStatus when successful
func (m *DeviceEvidence) GetOnboardingStatus()(*OnboardingStatus) {
    val, err := m.GetBackingStore().Get("onboardingStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnboardingStatus)
    }
    return nil
}
// GetOsBuild gets the osBuild property value. The build version for the operating system the device is running.
// returns a *int64 when successful
func (m *DeviceEvidence) GetOsBuild()(*int64) {
    val, err := m.GetBackingStore().Get("osBuild")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetOsPlatform gets the osPlatform property value. The operating system platform the device is running.
// returns a *string when successful
func (m *DeviceEvidence) GetOsPlatform()(*string) {
    val, err := m.GetBackingStore().Get("osPlatform")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRbacGroupId gets the rbacGroupId property value. The ID of the role-based access control (RBAC) device group.
// returns a *int32 when successful
func (m *DeviceEvidence) GetRbacGroupId()(*int32) {
    val, err := m.GetBackingStore().Get("rbacGroupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRbacGroupName gets the rbacGroupName property value. The name of the RBAC device group.
// returns a *string when successful
func (m *DeviceEvidence) GetRbacGroupName()(*string) {
    val, err := m.GetBackingStore().Get("rbacGroupName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskScore gets the riskScore property value. Risk score as evaluated by Microsoft Defender for Endpoint. The possible values are: none, informational, low, medium, high, unknownFutureValue.
// returns a *DeviceRiskScore when successful
func (m *DeviceEvidence) GetRiskScore()(*DeviceRiskScore) {
    val, err := m.GetBackingStore().Get("riskScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceRiskScore)
    }
    return nil
}
// GetVersion gets the version property value. The version of the operating system platform.
// returns a *string when successful
func (m *DeviceEvidence) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVmMetadata gets the vmMetadata property value. Metadata of the virtual machine (VM) on which Microsoft Defender for Endpoint is running.
// returns a VmMetadataable when successful
func (m *DeviceEvidence) GetVmMetadata()(VmMetadataable) {
    val, err := m.GetBackingStore().Get("vmMetadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VmMetadataable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("azureAdDeviceId", m.GetAzureAdDeviceId())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderAvStatus() != nil {
        cast := (*m.GetDefenderAvStatus()).String()
        err = writer.WriteStringValue("defenderAvStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceDnsName", m.GetDeviceDnsName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("dnsDomain", m.GetDnsDomain())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("firstSeenDateTime", m.GetFirstSeenDateTime())
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
        err = writer.WriteStringValue("hostName", m.GetHostName())
        if err != nil {
            return err
        }
    }
    if m.GetIpInterfaces() != nil {
        err = writer.WriteCollectionOfStringValues("ipInterfaces", m.GetIpInterfaces())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lastExternalIpAddress", m.GetLastExternalIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lastIpAddress", m.GetLastIpAddress())
        if err != nil {
            return err
        }
    }
    if m.GetLoggedOnUsers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLoggedOnUsers()))
        for i, v := range m.GetLoggedOnUsers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("loggedOnUsers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mdeDeviceId", m.GetMdeDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ntDomain", m.GetNtDomain())
        if err != nil {
            return err
        }
    }
    if m.GetOnboardingStatus() != nil {
        cast := (*m.GetOnboardingStatus()).String()
        err = writer.WriteStringValue("onboardingStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("osBuild", m.GetOsBuild())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osPlatform", m.GetOsPlatform())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("rbacGroupId", m.GetRbacGroupId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("rbacGroupName", m.GetRbacGroupName())
        if err != nil {
            return err
        }
    }
    if m.GetRiskScore() != nil {
        cast := (*m.GetRiskScore()).String()
        err = writer.WriteStringValue("riskScore", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("vmMetadata", m.GetVmMetadata())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureAdDeviceId sets the azureAdDeviceId property value. A unique identifier assigned to a device by Microsoft Entra ID when device is Microsoft Entra joined.
func (m *DeviceEvidence) SetAzureAdDeviceId(value *string)() {
    err := m.GetBackingStore().Set("azureAdDeviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderAvStatus sets the defenderAvStatus property value. State of the Defender AntiMalware engine. The possible values are: notReporting, disabled, notUpdated, updated, unknown, notSupported, unknownFutureValue.
func (m *DeviceEvidence) SetDefenderAvStatus(value *DefenderAvStatus)() {
    err := m.GetBackingStore().Set("defenderAvStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceDnsName sets the deviceDnsName property value. The fully qualified domain name (FQDN) for the device.
func (m *DeviceEvidence) SetDeviceDnsName(value *string)() {
    err := m.GetBackingStore().Set("deviceDnsName", value)
    if err != nil {
        panic(err)
    }
}
// SetDnsDomain sets the dnsDomain property value. The DNS domain that this computer belongs to. A sequence of labels separated by dots.
func (m *DeviceEvidence) SetDnsDomain(value *string)() {
    err := m.GetBackingStore().Set("dnsDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetFirstSeenDateTime sets the firstSeenDateTime property value. The date and time when the device was first seen.
func (m *DeviceEvidence) SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("firstSeenDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetHealthStatus sets the healthStatus property value. The health state of the device. The possible values are: active, inactive, impairedCommunication, noSensorData, noSensorDataImpairedCommunication, unknown, unknownFutureValue.
func (m *DeviceEvidence) SetHealthStatus(value *DeviceHealthStatus)() {
    err := m.GetBackingStore().Set("healthStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetHostName sets the hostName property value. The hostname without the domain suffix.
func (m *DeviceEvidence) SetHostName(value *string)() {
    err := m.GetBackingStore().Set("hostName", value)
    if err != nil {
        panic(err)
    }
}
// SetIpInterfaces sets the ipInterfaces property value. Ip interfaces of the device during the time of the alert.
func (m *DeviceEvidence) SetIpInterfaces(value []string)() {
    err := m.GetBackingStore().Set("ipInterfaces", value)
    if err != nil {
        panic(err)
    }
}
// SetLastExternalIpAddress sets the lastExternalIpAddress property value. The lastExternalIpAddress property
func (m *DeviceEvidence) SetLastExternalIpAddress(value *string)() {
    err := m.GetBackingStore().Set("lastExternalIpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetLastIpAddress sets the lastIpAddress property value. The lastIpAddress property
func (m *DeviceEvidence) SetLastIpAddress(value *string)() {
    err := m.GetBackingStore().Set("lastIpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetLoggedOnUsers sets the loggedOnUsers property value. Users that were logged on the machine during the time of the alert.
func (m *DeviceEvidence) SetLoggedOnUsers(value []LoggedOnUserable)() {
    err := m.GetBackingStore().Set("loggedOnUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetMdeDeviceId sets the mdeDeviceId property value. A unique identifier assigned to a device by Microsoft Defender for Endpoint.
func (m *DeviceEvidence) SetMdeDeviceId(value *string)() {
    err := m.GetBackingStore().Set("mdeDeviceId", value)
    if err != nil {
        panic(err)
    }
}
// SetNtDomain sets the ntDomain property value. A logical grouping of computers within a Microsoft Windows network.
func (m *DeviceEvidence) SetNtDomain(value *string)() {
    err := m.GetBackingStore().Set("ntDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetOnboardingStatus sets the onboardingStatus property value. The status of the machine onboarding to Microsoft Defender for Endpoint. The possible values are: insufficientInfo, onboarded, canBeOnboarded, unsupported, unknownFutureValue.
func (m *DeviceEvidence) SetOnboardingStatus(value *OnboardingStatus)() {
    err := m.GetBackingStore().Set("onboardingStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetOsBuild sets the osBuild property value. The build version for the operating system the device is running.
func (m *DeviceEvidence) SetOsBuild(value *int64)() {
    err := m.GetBackingStore().Set("osBuild", value)
    if err != nil {
        panic(err)
    }
}
// SetOsPlatform sets the osPlatform property value. The operating system platform the device is running.
func (m *DeviceEvidence) SetOsPlatform(value *string)() {
    err := m.GetBackingStore().Set("osPlatform", value)
    if err != nil {
        panic(err)
    }
}
// SetRbacGroupId sets the rbacGroupId property value. The ID of the role-based access control (RBAC) device group.
func (m *DeviceEvidence) SetRbacGroupId(value *int32)() {
    err := m.GetBackingStore().Set("rbacGroupId", value)
    if err != nil {
        panic(err)
    }
}
// SetRbacGroupName sets the rbacGroupName property value. The name of the RBAC device group.
func (m *DeviceEvidence) SetRbacGroupName(value *string)() {
    err := m.GetBackingStore().Set("rbacGroupName", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskScore sets the riskScore property value. Risk score as evaluated by Microsoft Defender for Endpoint. The possible values are: none, informational, low, medium, high, unknownFutureValue.
func (m *DeviceEvidence) SetRiskScore(value *DeviceRiskScore)() {
    err := m.GetBackingStore().Set("riskScore", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The version of the operating system platform.
func (m *DeviceEvidence) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
// SetVmMetadata sets the vmMetadata property value. Metadata of the virtual machine (VM) on which Microsoft Defender for Endpoint is running.
func (m *DeviceEvidence) SetVmMetadata(value VmMetadataable)() {
    err := m.GetBackingStore().Set("vmMetadata", value)
    if err != nil {
        panic(err)
    }
}
type DeviceEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureAdDeviceId()(*string)
    GetDefenderAvStatus()(*DefenderAvStatus)
    GetDeviceDnsName()(*string)
    GetDnsDomain()(*string)
    GetFirstSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetHealthStatus()(*DeviceHealthStatus)
    GetHostName()(*string)
    GetIpInterfaces()([]string)
    GetLastExternalIpAddress()(*string)
    GetLastIpAddress()(*string)
    GetLoggedOnUsers()([]LoggedOnUserable)
    GetMdeDeviceId()(*string)
    GetNtDomain()(*string)
    GetOnboardingStatus()(*OnboardingStatus)
    GetOsBuild()(*int64)
    GetOsPlatform()(*string)
    GetRbacGroupId()(*int32)
    GetRbacGroupName()(*string)
    GetRiskScore()(*DeviceRiskScore)
    GetVersion()(*string)
    GetVmMetadata()(VmMetadataable)
    SetAzureAdDeviceId(value *string)()
    SetDefenderAvStatus(value *DefenderAvStatus)()
    SetDeviceDnsName(value *string)()
    SetDnsDomain(value *string)()
    SetFirstSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetHealthStatus(value *DeviceHealthStatus)()
    SetHostName(value *string)()
    SetIpInterfaces(value []string)()
    SetLastExternalIpAddress(value *string)()
    SetLastIpAddress(value *string)()
    SetLoggedOnUsers(value []LoggedOnUserable)()
    SetMdeDeviceId(value *string)()
    SetNtDomain(value *string)()
    SetOnboardingStatus(value *OnboardingStatus)()
    SetOsBuild(value *int64)()
    SetOsPlatform(value *string)()
    SetRbacGroupId(value *int32)()
    SetRbacGroupName(value *string)()
    SetRiskScore(value *DeviceRiskScore)()
    SetVersion(value *string)()
    SetVmMetadata(value VmMetadataable)()
}

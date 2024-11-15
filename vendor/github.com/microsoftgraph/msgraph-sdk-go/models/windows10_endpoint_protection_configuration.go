package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10EndpointProtectionConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the Windows10EndpointProtectionConfiguration resource.
type Windows10EndpointProtectionConfiguration struct {
    DeviceConfiguration
}
// NewWindows10EndpointProtectionConfiguration instantiates a new Windows10EndpointProtectionConfiguration and sets the default values.
func NewWindows10EndpointProtectionConfiguration()(*Windows10EndpointProtectionConfiguration) {
    m := &Windows10EndpointProtectionConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10EndpointProtectionConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows10EndpointProtectionConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows10EndpointProtectionConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10EndpointProtectionConfiguration(), nil
}
// GetApplicationGuardAllowPersistence gets the applicationGuardAllowPersistence property value. Allow persisting user generated data inside the App Guard Containter (favorites, cookies, web passwords, etc.)
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardAllowPersistence()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardAllowPersistence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardAllowPrintToLocalPrinters gets the applicationGuardAllowPrintToLocalPrinters property value. Allow printing to Local Printers from Container
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardAllowPrintToLocalPrinters()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardAllowPrintToLocalPrinters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardAllowPrintToNetworkPrinters gets the applicationGuardAllowPrintToNetworkPrinters property value. Allow printing to Network Printers from Container
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardAllowPrintToNetworkPrinters()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardAllowPrintToNetworkPrinters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardAllowPrintToPDF gets the applicationGuardAllowPrintToPDF property value. Allow printing to PDF from Container
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardAllowPrintToPDF()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardAllowPrintToPDF")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardAllowPrintToXPS gets the applicationGuardAllowPrintToXPS property value. Allow printing to XPS from Container
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardAllowPrintToXPS()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardAllowPrintToXPS")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardBlockClipboardSharing gets the applicationGuardBlockClipboardSharing property value. Possible values for applicationGuardBlockClipboardSharingType
// returns a *ApplicationGuardBlockClipboardSharingType when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardBlockClipboardSharing()(*ApplicationGuardBlockClipboardSharingType) {
    val, err := m.GetBackingStore().Get("applicationGuardBlockClipboardSharing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ApplicationGuardBlockClipboardSharingType)
    }
    return nil
}
// GetApplicationGuardBlockFileTransfer gets the applicationGuardBlockFileTransfer property value. Possible values for applicationGuardBlockFileTransfer
// returns a *ApplicationGuardBlockFileTransferType when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardBlockFileTransfer()(*ApplicationGuardBlockFileTransferType) {
    val, err := m.GetBackingStore().Get("applicationGuardBlockFileTransfer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ApplicationGuardBlockFileTransferType)
    }
    return nil
}
// GetApplicationGuardBlockNonEnterpriseContent gets the applicationGuardBlockNonEnterpriseContent property value. Block enterprise sites to load non-enterprise content, such as third party plug-ins
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardBlockNonEnterpriseContent()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardBlockNonEnterpriseContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardEnabled gets the applicationGuardEnabled property value. Enable Windows Defender Application Guard
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplicationGuardForceAuditing gets the applicationGuardForceAuditing property value. Force auditing will persist Windows logs and events to meet security/compliance criteria (sample events are user login-logoff, use of privilege rights, software installation, system changes, etc.)
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetApplicationGuardForceAuditing()(*bool) {
    val, err := m.GetBackingStore().Get("applicationGuardForceAuditing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppLockerApplicationControl gets the appLockerApplicationControl property value. Possible values of AppLocker Application Control Types
// returns a *AppLockerApplicationControlType when successful
func (m *Windows10EndpointProtectionConfiguration) GetAppLockerApplicationControl()(*AppLockerApplicationControlType) {
    val, err := m.GetBackingStore().Get("appLockerApplicationControl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AppLockerApplicationControlType)
    }
    return nil
}
// GetBitLockerDisableWarningForOtherDiskEncryption gets the bitLockerDisableWarningForOtherDiskEncryption property value. Allows the Admin to disable the warning prompt for other disk encryption on the user machines.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetBitLockerDisableWarningForOtherDiskEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("bitLockerDisableWarningForOtherDiskEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBitLockerEnableStorageCardEncryptionOnMobile gets the bitLockerEnableStorageCardEncryptionOnMobile property value. Allows the admin to require encryption to be turned on using BitLocker. This policy is valid only for a mobile SKU.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetBitLockerEnableStorageCardEncryptionOnMobile()(*bool) {
    val, err := m.GetBackingStore().Get("bitLockerEnableStorageCardEncryptionOnMobile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBitLockerEncryptDevice gets the bitLockerEncryptDevice property value. Allows the admin to require encryption to be turned on using BitLocker.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetBitLockerEncryptDevice()(*bool) {
    val, err := m.GetBackingStore().Get("bitLockerEncryptDevice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBitLockerRemovableDrivePolicy gets the bitLockerRemovableDrivePolicy property value. BitLocker Removable Drive Policy.
// returns a BitLockerRemovableDrivePolicyable when successful
func (m *Windows10EndpointProtectionConfiguration) GetBitLockerRemovableDrivePolicy()(BitLockerRemovableDrivePolicyable) {
    val, err := m.GetBackingStore().Get("bitLockerRemovableDrivePolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BitLockerRemovableDrivePolicyable)
    }
    return nil
}
// GetDefenderAdditionalGuardedFolders gets the defenderAdditionalGuardedFolders property value. List of folder paths to be added to the list of protected folders
// returns a []string when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderAdditionalGuardedFolders()([]string) {
    val, err := m.GetBackingStore().Get("defenderAdditionalGuardedFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderAttackSurfaceReductionExcludedPaths gets the defenderAttackSurfaceReductionExcludedPaths property value. List of exe files and folders to be excluded from attack surface reduction rules
// returns a []string when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderAttackSurfaceReductionExcludedPaths()([]string) {
    val, err := m.GetBackingStore().Get("defenderAttackSurfaceReductionExcludedPaths")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderExploitProtectionXml gets the defenderExploitProtectionXml property value. Xml content containing information regarding exploit protection details.
// returns a []byte when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderExploitProtectionXml()([]byte) {
    val, err := m.GetBackingStore().Get("defenderExploitProtectionXml")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetDefenderExploitProtectionXmlFileName gets the defenderExploitProtectionXmlFileName property value. Name of the file from which DefenderExploitProtectionXml was obtained.
// returns a *string when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderExploitProtectionXmlFileName()(*string) {
    val, err := m.GetBackingStore().Get("defenderExploitProtectionXmlFileName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDefenderGuardedFoldersAllowedAppPaths gets the defenderGuardedFoldersAllowedAppPaths property value. List of paths to exe that are allowed to access protected folders
// returns a []string when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderGuardedFoldersAllowedAppPaths()([]string) {
    val, err := m.GetBackingStore().Get("defenderGuardedFoldersAllowedAppPaths")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderSecurityCenterBlockExploitProtectionOverride gets the defenderSecurityCenterBlockExploitProtectionOverride property value. Indicates whether or not to block user from overriding Exploit Protection settings.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetDefenderSecurityCenterBlockExploitProtectionOverride()(*bool) {
    val, err := m.GetBackingStore().Get("defenderSecurityCenterBlockExploitProtectionOverride")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Windows10EndpointProtectionConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["applicationGuardAllowPersistence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardAllowPersistence(val)
        }
        return nil
    }
    res["applicationGuardAllowPrintToLocalPrinters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardAllowPrintToLocalPrinters(val)
        }
        return nil
    }
    res["applicationGuardAllowPrintToNetworkPrinters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardAllowPrintToNetworkPrinters(val)
        }
        return nil
    }
    res["applicationGuardAllowPrintToPDF"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardAllowPrintToPDF(val)
        }
        return nil
    }
    res["applicationGuardAllowPrintToXPS"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardAllowPrintToXPS(val)
        }
        return nil
    }
    res["applicationGuardBlockClipboardSharing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseApplicationGuardBlockClipboardSharingType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardBlockClipboardSharing(val.(*ApplicationGuardBlockClipboardSharingType))
        }
        return nil
    }
    res["applicationGuardBlockFileTransfer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseApplicationGuardBlockFileTransferType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardBlockFileTransfer(val.(*ApplicationGuardBlockFileTransferType))
        }
        return nil
    }
    res["applicationGuardBlockNonEnterpriseContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardBlockNonEnterpriseContent(val)
        }
        return nil
    }
    res["applicationGuardEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardEnabled(val)
        }
        return nil
    }
    res["applicationGuardForceAuditing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationGuardForceAuditing(val)
        }
        return nil
    }
    res["appLockerApplicationControl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAppLockerApplicationControlType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppLockerApplicationControl(val.(*AppLockerApplicationControlType))
        }
        return nil
    }
    res["bitLockerDisableWarningForOtherDiskEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitLockerDisableWarningForOtherDiskEncryption(val)
        }
        return nil
    }
    res["bitLockerEnableStorageCardEncryptionOnMobile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitLockerEnableStorageCardEncryptionOnMobile(val)
        }
        return nil
    }
    res["bitLockerEncryptDevice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitLockerEncryptDevice(val)
        }
        return nil
    }
    res["bitLockerRemovableDrivePolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBitLockerRemovableDrivePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitLockerRemovableDrivePolicy(val.(BitLockerRemovableDrivePolicyable))
        }
        return nil
    }
    res["defenderAdditionalGuardedFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDefenderAdditionalGuardedFolders(res)
        }
        return nil
    }
    res["defenderAttackSurfaceReductionExcludedPaths"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDefenderAttackSurfaceReductionExcludedPaths(res)
        }
        return nil
    }
    res["defenderExploitProtectionXml"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderExploitProtectionXml(val)
        }
        return nil
    }
    res["defenderExploitProtectionXmlFileName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderExploitProtectionXmlFileName(val)
        }
        return nil
    }
    res["defenderGuardedFoldersAllowedAppPaths"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetDefenderGuardedFoldersAllowedAppPaths(res)
        }
        return nil
    }
    res["defenderSecurityCenterBlockExploitProtectionOverride"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderSecurityCenterBlockExploitProtectionOverride(val)
        }
        return nil
    }
    res["firewallBlockStatefulFTP"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallBlockStatefulFTP(val)
        }
        return nil
    }
    res["firewallCertificateRevocationListCheckMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFirewallCertificateRevocationListCheckMethodType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallCertificateRevocationListCheckMethod(val.(*FirewallCertificateRevocationListCheckMethodType))
        }
        return nil
    }
    res["firewallIdleTimeoutForSecurityAssociationInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallIdleTimeoutForSecurityAssociationInSeconds(val)
        }
        return nil
    }
    res["firewallIPSecExemptionsAllowDHCP"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallIPSecExemptionsAllowDHCP(val)
        }
        return nil
    }
    res["firewallIPSecExemptionsAllowICMP"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallIPSecExemptionsAllowICMP(val)
        }
        return nil
    }
    res["firewallIPSecExemptionsAllowNeighborDiscovery"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallIPSecExemptionsAllowNeighborDiscovery(val)
        }
        return nil
    }
    res["firewallIPSecExemptionsAllowRouterDiscovery"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallIPSecExemptionsAllowRouterDiscovery(val)
        }
        return nil
    }
    res["firewallMergeKeyingModuleSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallMergeKeyingModuleSettings(val)
        }
        return nil
    }
    res["firewallPacketQueueingMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFirewallPacketQueueingMethodType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallPacketQueueingMethod(val.(*FirewallPacketQueueingMethodType))
        }
        return nil
    }
    res["firewallPreSharedKeyEncodingMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFirewallPreSharedKeyEncodingMethodType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallPreSharedKeyEncodingMethod(val.(*FirewallPreSharedKeyEncodingMethodType))
        }
        return nil
    }
    res["firewallProfileDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsFirewallNetworkProfileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallProfileDomain(val.(WindowsFirewallNetworkProfileable))
        }
        return nil
    }
    res["firewallProfilePrivate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsFirewallNetworkProfileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallProfilePrivate(val.(WindowsFirewallNetworkProfileable))
        }
        return nil
    }
    res["firewallProfilePublic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsFirewallNetworkProfileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFirewallProfilePublic(val.(WindowsFirewallNetworkProfileable))
        }
        return nil
    }
    res["smartScreenBlockOverrideForFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmartScreenBlockOverrideForFiles(val)
        }
        return nil
    }
    res["smartScreenEnableInShell"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmartScreenEnableInShell(val)
        }
        return nil
    }
    return res
}
// GetFirewallBlockStatefulFTP gets the firewallBlockStatefulFTP property value. Blocks stateful FTP connections to the device
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallBlockStatefulFTP()(*bool) {
    val, err := m.GetBackingStore().Get("firewallBlockStatefulFTP")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallCertificateRevocationListCheckMethod gets the firewallCertificateRevocationListCheckMethod property value. Possible values for firewallCertificateRevocationListCheckMethod
// returns a *FirewallCertificateRevocationListCheckMethodType when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallCertificateRevocationListCheckMethod()(*FirewallCertificateRevocationListCheckMethodType) {
    val, err := m.GetBackingStore().Get("firewallCertificateRevocationListCheckMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FirewallCertificateRevocationListCheckMethodType)
    }
    return nil
}
// GetFirewallIdleTimeoutForSecurityAssociationInSeconds gets the firewallIdleTimeoutForSecurityAssociationInSeconds property value. Configures the idle timeout for security associations, in seconds, from 300 to 3600 inclusive. This is the period after which security associations will expire and be deleted. Valid values 300 to 3600
// returns a *int32 when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallIdleTimeoutForSecurityAssociationInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("firewallIdleTimeoutForSecurityAssociationInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFirewallIPSecExemptionsAllowDHCP gets the firewallIPSecExemptionsAllowDHCP property value. Configures IPSec exemptions to allow both IPv4 and IPv6 DHCP traffic
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallIPSecExemptionsAllowDHCP()(*bool) {
    val, err := m.GetBackingStore().Get("firewallIPSecExemptionsAllowDHCP")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallIPSecExemptionsAllowICMP gets the firewallIPSecExemptionsAllowICMP property value. Configures IPSec exemptions to allow ICMP
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallIPSecExemptionsAllowICMP()(*bool) {
    val, err := m.GetBackingStore().Get("firewallIPSecExemptionsAllowICMP")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallIPSecExemptionsAllowNeighborDiscovery gets the firewallIPSecExemptionsAllowNeighborDiscovery property value. Configures IPSec exemptions to allow neighbor discovery IPv6 ICMP type-codes
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallIPSecExemptionsAllowNeighborDiscovery()(*bool) {
    val, err := m.GetBackingStore().Get("firewallIPSecExemptionsAllowNeighborDiscovery")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallIPSecExemptionsAllowRouterDiscovery gets the firewallIPSecExemptionsAllowRouterDiscovery property value. Configures IPSec exemptions to allow router discovery IPv6 ICMP type-codes
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallIPSecExemptionsAllowRouterDiscovery()(*bool) {
    val, err := m.GetBackingStore().Get("firewallIPSecExemptionsAllowRouterDiscovery")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallMergeKeyingModuleSettings gets the firewallMergeKeyingModuleSettings property value. If an authentication set is not fully supported by a keying module, direct the module to ignore only unsupported authentication suites rather than the entire set
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallMergeKeyingModuleSettings()(*bool) {
    val, err := m.GetBackingStore().Get("firewallMergeKeyingModuleSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFirewallPacketQueueingMethod gets the firewallPacketQueueingMethod property value. Possible values for firewallPacketQueueingMethod
// returns a *FirewallPacketQueueingMethodType when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallPacketQueueingMethod()(*FirewallPacketQueueingMethodType) {
    val, err := m.GetBackingStore().Get("firewallPacketQueueingMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FirewallPacketQueueingMethodType)
    }
    return nil
}
// GetFirewallPreSharedKeyEncodingMethod gets the firewallPreSharedKeyEncodingMethod property value. Possible values for firewallPreSharedKeyEncodingMethod
// returns a *FirewallPreSharedKeyEncodingMethodType when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallPreSharedKeyEncodingMethod()(*FirewallPreSharedKeyEncodingMethodType) {
    val, err := m.GetBackingStore().Get("firewallPreSharedKeyEncodingMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FirewallPreSharedKeyEncodingMethodType)
    }
    return nil
}
// GetFirewallProfileDomain gets the firewallProfileDomain property value. Configures the firewall profile settings for domain networks
// returns a WindowsFirewallNetworkProfileable when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallProfileDomain()(WindowsFirewallNetworkProfileable) {
    val, err := m.GetBackingStore().Get("firewallProfileDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsFirewallNetworkProfileable)
    }
    return nil
}
// GetFirewallProfilePrivate gets the firewallProfilePrivate property value. Configures the firewall profile settings for private networks
// returns a WindowsFirewallNetworkProfileable when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallProfilePrivate()(WindowsFirewallNetworkProfileable) {
    val, err := m.GetBackingStore().Get("firewallProfilePrivate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsFirewallNetworkProfileable)
    }
    return nil
}
// GetFirewallProfilePublic gets the firewallProfilePublic property value. Configures the firewall profile settings for public networks
// returns a WindowsFirewallNetworkProfileable when successful
func (m *Windows10EndpointProtectionConfiguration) GetFirewallProfilePublic()(WindowsFirewallNetworkProfileable) {
    val, err := m.GetBackingStore().Get("firewallProfilePublic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsFirewallNetworkProfileable)
    }
    return nil
}
// GetSmartScreenBlockOverrideForFiles gets the smartScreenBlockOverrideForFiles property value. Allows IT Admins to control whether users can can ignore SmartScreen warnings and run malicious files.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetSmartScreenBlockOverrideForFiles()(*bool) {
    val, err := m.GetBackingStore().Get("smartScreenBlockOverrideForFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSmartScreenEnableInShell gets the smartScreenEnableInShell property value. Allows IT Admins to configure SmartScreen for Windows.
// returns a *bool when successful
func (m *Windows10EndpointProtectionConfiguration) GetSmartScreenEnableInShell()(*bool) {
    val, err := m.GetBackingStore().Get("smartScreenEnableInShell")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Windows10EndpointProtectionConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("applicationGuardAllowPersistence", m.GetApplicationGuardAllowPersistence())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardAllowPrintToLocalPrinters", m.GetApplicationGuardAllowPrintToLocalPrinters())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardAllowPrintToNetworkPrinters", m.GetApplicationGuardAllowPrintToNetworkPrinters())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardAllowPrintToPDF", m.GetApplicationGuardAllowPrintToPDF())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardAllowPrintToXPS", m.GetApplicationGuardAllowPrintToXPS())
        if err != nil {
            return err
        }
    }
    if m.GetApplicationGuardBlockClipboardSharing() != nil {
        cast := (*m.GetApplicationGuardBlockClipboardSharing()).String()
        err = writer.WriteStringValue("applicationGuardBlockClipboardSharing", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetApplicationGuardBlockFileTransfer() != nil {
        cast := (*m.GetApplicationGuardBlockFileTransfer()).String()
        err = writer.WriteStringValue("applicationGuardBlockFileTransfer", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardBlockNonEnterpriseContent", m.GetApplicationGuardBlockNonEnterpriseContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardEnabled", m.GetApplicationGuardEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("applicationGuardForceAuditing", m.GetApplicationGuardForceAuditing())
        if err != nil {
            return err
        }
    }
    if m.GetAppLockerApplicationControl() != nil {
        cast := (*m.GetAppLockerApplicationControl()).String()
        err = writer.WriteStringValue("appLockerApplicationControl", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bitLockerDisableWarningForOtherDiskEncryption", m.GetBitLockerDisableWarningForOtherDiskEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bitLockerEnableStorageCardEncryptionOnMobile", m.GetBitLockerEnableStorageCardEncryptionOnMobile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bitLockerEncryptDevice", m.GetBitLockerEncryptDevice())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("bitLockerRemovableDrivePolicy", m.GetBitLockerRemovableDrivePolicy())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderAdditionalGuardedFolders() != nil {
        err = writer.WriteCollectionOfStringValues("defenderAdditionalGuardedFolders", m.GetDefenderAdditionalGuardedFolders())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderAttackSurfaceReductionExcludedPaths() != nil {
        err = writer.WriteCollectionOfStringValues("defenderAttackSurfaceReductionExcludedPaths", m.GetDefenderAttackSurfaceReductionExcludedPaths())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("defenderExploitProtectionXml", m.GetDefenderExploitProtectionXml())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("defenderExploitProtectionXmlFileName", m.GetDefenderExploitProtectionXmlFileName())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderGuardedFoldersAllowedAppPaths() != nil {
        err = writer.WriteCollectionOfStringValues("defenderGuardedFoldersAllowedAppPaths", m.GetDefenderGuardedFoldersAllowedAppPaths())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderSecurityCenterBlockExploitProtectionOverride", m.GetDefenderSecurityCenterBlockExploitProtectionOverride())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallBlockStatefulFTP", m.GetFirewallBlockStatefulFTP())
        if err != nil {
            return err
        }
    }
    if m.GetFirewallCertificateRevocationListCheckMethod() != nil {
        cast := (*m.GetFirewallCertificateRevocationListCheckMethod()).String()
        err = writer.WriteStringValue("firewallCertificateRevocationListCheckMethod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("firewallIdleTimeoutForSecurityAssociationInSeconds", m.GetFirewallIdleTimeoutForSecurityAssociationInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallIPSecExemptionsAllowDHCP", m.GetFirewallIPSecExemptionsAllowDHCP())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallIPSecExemptionsAllowICMP", m.GetFirewallIPSecExemptionsAllowICMP())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallIPSecExemptionsAllowNeighborDiscovery", m.GetFirewallIPSecExemptionsAllowNeighborDiscovery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallIPSecExemptionsAllowRouterDiscovery", m.GetFirewallIPSecExemptionsAllowRouterDiscovery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallMergeKeyingModuleSettings", m.GetFirewallMergeKeyingModuleSettings())
        if err != nil {
            return err
        }
    }
    if m.GetFirewallPacketQueueingMethod() != nil {
        cast := (*m.GetFirewallPacketQueueingMethod()).String()
        err = writer.WriteStringValue("firewallPacketQueueingMethod", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFirewallPreSharedKeyEncodingMethod() != nil {
        cast := (*m.GetFirewallPreSharedKeyEncodingMethod()).String()
        err = writer.WriteStringValue("firewallPreSharedKeyEncodingMethod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("firewallProfileDomain", m.GetFirewallProfileDomain())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("firewallProfilePrivate", m.GetFirewallProfilePrivate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("firewallProfilePublic", m.GetFirewallProfilePublic())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("smartScreenBlockOverrideForFiles", m.GetSmartScreenBlockOverrideForFiles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("smartScreenEnableInShell", m.GetSmartScreenEnableInShell())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationGuardAllowPersistence sets the applicationGuardAllowPersistence property value. Allow persisting user generated data inside the App Guard Containter (favorites, cookies, web passwords, etc.)
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardAllowPersistence(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardAllowPersistence", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardAllowPrintToLocalPrinters sets the applicationGuardAllowPrintToLocalPrinters property value. Allow printing to Local Printers from Container
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardAllowPrintToLocalPrinters(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardAllowPrintToLocalPrinters", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardAllowPrintToNetworkPrinters sets the applicationGuardAllowPrintToNetworkPrinters property value. Allow printing to Network Printers from Container
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardAllowPrintToNetworkPrinters(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardAllowPrintToNetworkPrinters", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardAllowPrintToPDF sets the applicationGuardAllowPrintToPDF property value. Allow printing to PDF from Container
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardAllowPrintToPDF(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardAllowPrintToPDF", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardAllowPrintToXPS sets the applicationGuardAllowPrintToXPS property value. Allow printing to XPS from Container
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardAllowPrintToXPS(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardAllowPrintToXPS", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardBlockClipboardSharing sets the applicationGuardBlockClipboardSharing property value. Possible values for applicationGuardBlockClipboardSharingType
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardBlockClipboardSharing(value *ApplicationGuardBlockClipboardSharingType)() {
    err := m.GetBackingStore().Set("applicationGuardBlockClipboardSharing", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardBlockFileTransfer sets the applicationGuardBlockFileTransfer property value. Possible values for applicationGuardBlockFileTransfer
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardBlockFileTransfer(value *ApplicationGuardBlockFileTransferType)() {
    err := m.GetBackingStore().Set("applicationGuardBlockFileTransfer", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardBlockNonEnterpriseContent sets the applicationGuardBlockNonEnterpriseContent property value. Block enterprise sites to load non-enterprise content, such as third party plug-ins
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardBlockNonEnterpriseContent(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardBlockNonEnterpriseContent", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardEnabled sets the applicationGuardEnabled property value. Enable Windows Defender Application Guard
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardEnabled(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationGuardForceAuditing sets the applicationGuardForceAuditing property value. Force auditing will persist Windows logs and events to meet security/compliance criteria (sample events are user login-logoff, use of privilege rights, software installation, system changes, etc.)
func (m *Windows10EndpointProtectionConfiguration) SetApplicationGuardForceAuditing(value *bool)() {
    err := m.GetBackingStore().Set("applicationGuardForceAuditing", value)
    if err != nil {
        panic(err)
    }
}
// SetAppLockerApplicationControl sets the appLockerApplicationControl property value. Possible values of AppLocker Application Control Types
func (m *Windows10EndpointProtectionConfiguration) SetAppLockerApplicationControl(value *AppLockerApplicationControlType)() {
    err := m.GetBackingStore().Set("appLockerApplicationControl", value)
    if err != nil {
        panic(err)
    }
}
// SetBitLockerDisableWarningForOtherDiskEncryption sets the bitLockerDisableWarningForOtherDiskEncryption property value. Allows the Admin to disable the warning prompt for other disk encryption on the user machines.
func (m *Windows10EndpointProtectionConfiguration) SetBitLockerDisableWarningForOtherDiskEncryption(value *bool)() {
    err := m.GetBackingStore().Set("bitLockerDisableWarningForOtherDiskEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetBitLockerEnableStorageCardEncryptionOnMobile sets the bitLockerEnableStorageCardEncryptionOnMobile property value. Allows the admin to require encryption to be turned on using BitLocker. This policy is valid only for a mobile SKU.
func (m *Windows10EndpointProtectionConfiguration) SetBitLockerEnableStorageCardEncryptionOnMobile(value *bool)() {
    err := m.GetBackingStore().Set("bitLockerEnableStorageCardEncryptionOnMobile", value)
    if err != nil {
        panic(err)
    }
}
// SetBitLockerEncryptDevice sets the bitLockerEncryptDevice property value. Allows the admin to require encryption to be turned on using BitLocker.
func (m *Windows10EndpointProtectionConfiguration) SetBitLockerEncryptDevice(value *bool)() {
    err := m.GetBackingStore().Set("bitLockerEncryptDevice", value)
    if err != nil {
        panic(err)
    }
}
// SetBitLockerRemovableDrivePolicy sets the bitLockerRemovableDrivePolicy property value. BitLocker Removable Drive Policy.
func (m *Windows10EndpointProtectionConfiguration) SetBitLockerRemovableDrivePolicy(value BitLockerRemovableDrivePolicyable)() {
    err := m.GetBackingStore().Set("bitLockerRemovableDrivePolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderAdditionalGuardedFolders sets the defenderAdditionalGuardedFolders property value. List of folder paths to be added to the list of protected folders
func (m *Windows10EndpointProtectionConfiguration) SetDefenderAdditionalGuardedFolders(value []string)() {
    err := m.GetBackingStore().Set("defenderAdditionalGuardedFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderAttackSurfaceReductionExcludedPaths sets the defenderAttackSurfaceReductionExcludedPaths property value. List of exe files and folders to be excluded from attack surface reduction rules
func (m *Windows10EndpointProtectionConfiguration) SetDefenderAttackSurfaceReductionExcludedPaths(value []string)() {
    err := m.GetBackingStore().Set("defenderAttackSurfaceReductionExcludedPaths", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderExploitProtectionXml sets the defenderExploitProtectionXml property value. Xml content containing information regarding exploit protection details.
func (m *Windows10EndpointProtectionConfiguration) SetDefenderExploitProtectionXml(value []byte)() {
    err := m.GetBackingStore().Set("defenderExploitProtectionXml", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderExploitProtectionXmlFileName sets the defenderExploitProtectionXmlFileName property value. Name of the file from which DefenderExploitProtectionXml was obtained.
func (m *Windows10EndpointProtectionConfiguration) SetDefenderExploitProtectionXmlFileName(value *string)() {
    err := m.GetBackingStore().Set("defenderExploitProtectionXmlFileName", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderGuardedFoldersAllowedAppPaths sets the defenderGuardedFoldersAllowedAppPaths property value. List of paths to exe that are allowed to access protected folders
func (m *Windows10EndpointProtectionConfiguration) SetDefenderGuardedFoldersAllowedAppPaths(value []string)() {
    err := m.GetBackingStore().Set("defenderGuardedFoldersAllowedAppPaths", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderSecurityCenterBlockExploitProtectionOverride sets the defenderSecurityCenterBlockExploitProtectionOverride property value. Indicates whether or not to block user from overriding Exploit Protection settings.
func (m *Windows10EndpointProtectionConfiguration) SetDefenderSecurityCenterBlockExploitProtectionOverride(value *bool)() {
    err := m.GetBackingStore().Set("defenderSecurityCenterBlockExploitProtectionOverride", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallBlockStatefulFTP sets the firewallBlockStatefulFTP property value. Blocks stateful FTP connections to the device
func (m *Windows10EndpointProtectionConfiguration) SetFirewallBlockStatefulFTP(value *bool)() {
    err := m.GetBackingStore().Set("firewallBlockStatefulFTP", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallCertificateRevocationListCheckMethod sets the firewallCertificateRevocationListCheckMethod property value. Possible values for firewallCertificateRevocationListCheckMethod
func (m *Windows10EndpointProtectionConfiguration) SetFirewallCertificateRevocationListCheckMethod(value *FirewallCertificateRevocationListCheckMethodType)() {
    err := m.GetBackingStore().Set("firewallCertificateRevocationListCheckMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallIdleTimeoutForSecurityAssociationInSeconds sets the firewallIdleTimeoutForSecurityAssociationInSeconds property value. Configures the idle timeout for security associations, in seconds, from 300 to 3600 inclusive. This is the period after which security associations will expire and be deleted. Valid values 300 to 3600
func (m *Windows10EndpointProtectionConfiguration) SetFirewallIdleTimeoutForSecurityAssociationInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("firewallIdleTimeoutForSecurityAssociationInSeconds", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallIPSecExemptionsAllowDHCP sets the firewallIPSecExemptionsAllowDHCP property value. Configures IPSec exemptions to allow both IPv4 and IPv6 DHCP traffic
func (m *Windows10EndpointProtectionConfiguration) SetFirewallIPSecExemptionsAllowDHCP(value *bool)() {
    err := m.GetBackingStore().Set("firewallIPSecExemptionsAllowDHCP", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallIPSecExemptionsAllowICMP sets the firewallIPSecExemptionsAllowICMP property value. Configures IPSec exemptions to allow ICMP
func (m *Windows10EndpointProtectionConfiguration) SetFirewallIPSecExemptionsAllowICMP(value *bool)() {
    err := m.GetBackingStore().Set("firewallIPSecExemptionsAllowICMP", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallIPSecExemptionsAllowNeighborDiscovery sets the firewallIPSecExemptionsAllowNeighborDiscovery property value. Configures IPSec exemptions to allow neighbor discovery IPv6 ICMP type-codes
func (m *Windows10EndpointProtectionConfiguration) SetFirewallIPSecExemptionsAllowNeighborDiscovery(value *bool)() {
    err := m.GetBackingStore().Set("firewallIPSecExemptionsAllowNeighborDiscovery", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallIPSecExemptionsAllowRouterDiscovery sets the firewallIPSecExemptionsAllowRouterDiscovery property value. Configures IPSec exemptions to allow router discovery IPv6 ICMP type-codes
func (m *Windows10EndpointProtectionConfiguration) SetFirewallIPSecExemptionsAllowRouterDiscovery(value *bool)() {
    err := m.GetBackingStore().Set("firewallIPSecExemptionsAllowRouterDiscovery", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallMergeKeyingModuleSettings sets the firewallMergeKeyingModuleSettings property value. If an authentication set is not fully supported by a keying module, direct the module to ignore only unsupported authentication suites rather than the entire set
func (m *Windows10EndpointProtectionConfiguration) SetFirewallMergeKeyingModuleSettings(value *bool)() {
    err := m.GetBackingStore().Set("firewallMergeKeyingModuleSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallPacketQueueingMethod sets the firewallPacketQueueingMethod property value. Possible values for firewallPacketQueueingMethod
func (m *Windows10EndpointProtectionConfiguration) SetFirewallPacketQueueingMethod(value *FirewallPacketQueueingMethodType)() {
    err := m.GetBackingStore().Set("firewallPacketQueueingMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallPreSharedKeyEncodingMethod sets the firewallPreSharedKeyEncodingMethod property value. Possible values for firewallPreSharedKeyEncodingMethod
func (m *Windows10EndpointProtectionConfiguration) SetFirewallPreSharedKeyEncodingMethod(value *FirewallPreSharedKeyEncodingMethodType)() {
    err := m.GetBackingStore().Set("firewallPreSharedKeyEncodingMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallProfileDomain sets the firewallProfileDomain property value. Configures the firewall profile settings for domain networks
func (m *Windows10EndpointProtectionConfiguration) SetFirewallProfileDomain(value WindowsFirewallNetworkProfileable)() {
    err := m.GetBackingStore().Set("firewallProfileDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallProfilePrivate sets the firewallProfilePrivate property value. Configures the firewall profile settings for private networks
func (m *Windows10EndpointProtectionConfiguration) SetFirewallProfilePrivate(value WindowsFirewallNetworkProfileable)() {
    err := m.GetBackingStore().Set("firewallProfilePrivate", value)
    if err != nil {
        panic(err)
    }
}
// SetFirewallProfilePublic sets the firewallProfilePublic property value. Configures the firewall profile settings for public networks
func (m *Windows10EndpointProtectionConfiguration) SetFirewallProfilePublic(value WindowsFirewallNetworkProfileable)() {
    err := m.GetBackingStore().Set("firewallProfilePublic", value)
    if err != nil {
        panic(err)
    }
}
// SetSmartScreenBlockOverrideForFiles sets the smartScreenBlockOverrideForFiles property value. Allows IT Admins to control whether users can can ignore SmartScreen warnings and run malicious files.
func (m *Windows10EndpointProtectionConfiguration) SetSmartScreenBlockOverrideForFiles(value *bool)() {
    err := m.GetBackingStore().Set("smartScreenBlockOverrideForFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetSmartScreenEnableInShell sets the smartScreenEnableInShell property value. Allows IT Admins to configure SmartScreen for Windows.
func (m *Windows10EndpointProtectionConfiguration) SetSmartScreenEnableInShell(value *bool)() {
    err := m.GetBackingStore().Set("smartScreenEnableInShell", value)
    if err != nil {
        panic(err)
    }
}
type Windows10EndpointProtectionConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicationGuardAllowPersistence()(*bool)
    GetApplicationGuardAllowPrintToLocalPrinters()(*bool)
    GetApplicationGuardAllowPrintToNetworkPrinters()(*bool)
    GetApplicationGuardAllowPrintToPDF()(*bool)
    GetApplicationGuardAllowPrintToXPS()(*bool)
    GetApplicationGuardBlockClipboardSharing()(*ApplicationGuardBlockClipboardSharingType)
    GetApplicationGuardBlockFileTransfer()(*ApplicationGuardBlockFileTransferType)
    GetApplicationGuardBlockNonEnterpriseContent()(*bool)
    GetApplicationGuardEnabled()(*bool)
    GetApplicationGuardForceAuditing()(*bool)
    GetAppLockerApplicationControl()(*AppLockerApplicationControlType)
    GetBitLockerDisableWarningForOtherDiskEncryption()(*bool)
    GetBitLockerEnableStorageCardEncryptionOnMobile()(*bool)
    GetBitLockerEncryptDevice()(*bool)
    GetBitLockerRemovableDrivePolicy()(BitLockerRemovableDrivePolicyable)
    GetDefenderAdditionalGuardedFolders()([]string)
    GetDefenderAttackSurfaceReductionExcludedPaths()([]string)
    GetDefenderExploitProtectionXml()([]byte)
    GetDefenderExploitProtectionXmlFileName()(*string)
    GetDefenderGuardedFoldersAllowedAppPaths()([]string)
    GetDefenderSecurityCenterBlockExploitProtectionOverride()(*bool)
    GetFirewallBlockStatefulFTP()(*bool)
    GetFirewallCertificateRevocationListCheckMethod()(*FirewallCertificateRevocationListCheckMethodType)
    GetFirewallIdleTimeoutForSecurityAssociationInSeconds()(*int32)
    GetFirewallIPSecExemptionsAllowDHCP()(*bool)
    GetFirewallIPSecExemptionsAllowICMP()(*bool)
    GetFirewallIPSecExemptionsAllowNeighborDiscovery()(*bool)
    GetFirewallIPSecExemptionsAllowRouterDiscovery()(*bool)
    GetFirewallMergeKeyingModuleSettings()(*bool)
    GetFirewallPacketQueueingMethod()(*FirewallPacketQueueingMethodType)
    GetFirewallPreSharedKeyEncodingMethod()(*FirewallPreSharedKeyEncodingMethodType)
    GetFirewallProfileDomain()(WindowsFirewallNetworkProfileable)
    GetFirewallProfilePrivate()(WindowsFirewallNetworkProfileable)
    GetFirewallProfilePublic()(WindowsFirewallNetworkProfileable)
    GetSmartScreenBlockOverrideForFiles()(*bool)
    GetSmartScreenEnableInShell()(*bool)
    SetApplicationGuardAllowPersistence(value *bool)()
    SetApplicationGuardAllowPrintToLocalPrinters(value *bool)()
    SetApplicationGuardAllowPrintToNetworkPrinters(value *bool)()
    SetApplicationGuardAllowPrintToPDF(value *bool)()
    SetApplicationGuardAllowPrintToXPS(value *bool)()
    SetApplicationGuardBlockClipboardSharing(value *ApplicationGuardBlockClipboardSharingType)()
    SetApplicationGuardBlockFileTransfer(value *ApplicationGuardBlockFileTransferType)()
    SetApplicationGuardBlockNonEnterpriseContent(value *bool)()
    SetApplicationGuardEnabled(value *bool)()
    SetApplicationGuardForceAuditing(value *bool)()
    SetAppLockerApplicationControl(value *AppLockerApplicationControlType)()
    SetBitLockerDisableWarningForOtherDiskEncryption(value *bool)()
    SetBitLockerEnableStorageCardEncryptionOnMobile(value *bool)()
    SetBitLockerEncryptDevice(value *bool)()
    SetBitLockerRemovableDrivePolicy(value BitLockerRemovableDrivePolicyable)()
    SetDefenderAdditionalGuardedFolders(value []string)()
    SetDefenderAttackSurfaceReductionExcludedPaths(value []string)()
    SetDefenderExploitProtectionXml(value []byte)()
    SetDefenderExploitProtectionXmlFileName(value *string)()
    SetDefenderGuardedFoldersAllowedAppPaths(value []string)()
    SetDefenderSecurityCenterBlockExploitProtectionOverride(value *bool)()
    SetFirewallBlockStatefulFTP(value *bool)()
    SetFirewallCertificateRevocationListCheckMethod(value *FirewallCertificateRevocationListCheckMethodType)()
    SetFirewallIdleTimeoutForSecurityAssociationInSeconds(value *int32)()
    SetFirewallIPSecExemptionsAllowDHCP(value *bool)()
    SetFirewallIPSecExemptionsAllowICMP(value *bool)()
    SetFirewallIPSecExemptionsAllowNeighborDiscovery(value *bool)()
    SetFirewallIPSecExemptionsAllowRouterDiscovery(value *bool)()
    SetFirewallMergeKeyingModuleSettings(value *bool)()
    SetFirewallPacketQueueingMethod(value *FirewallPacketQueueingMethodType)()
    SetFirewallPreSharedKeyEncodingMethod(value *FirewallPreSharedKeyEncodingMethodType)()
    SetFirewallProfileDomain(value WindowsFirewallNetworkProfileable)()
    SetFirewallProfilePrivate(value WindowsFirewallNetworkProfileable)()
    SetFirewallProfilePublic(value WindowsFirewallNetworkProfileable)()
    SetSmartScreenBlockOverrideForFiles(value *bool)()
    SetSmartScreenEnableInShell(value *bool)()
}

package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtection policy for Windows information protection to configure detailed management settings
type WindowsInformationProtection struct {
    ManagedAppPolicy
}
// NewWindowsInformationProtection instantiates a new WindowsInformationProtection and sets the default values.
func NewWindowsInformationProtection()(*WindowsInformationProtection) {
    m := &WindowsInformationProtection{
        ManagedAppPolicy: *NewManagedAppPolicy(),
    }
    odataTypeValue := "#microsoft.graph.windowsInformationProtection"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsInformationProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsInformationProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.mdmWindowsInformationProtectionPolicy":
                        return NewMdmWindowsInformationProtectionPolicy(), nil
                    case "#microsoft.graph.windowsInformationProtectionPolicy":
                        return NewWindowsInformationProtectionPolicy(), nil
                }
            }
        }
    }
    return NewWindowsInformationProtection(), nil
}
// GetAssignments gets the assignments property value. Navigation property to list of security groups targeted for policy.
// returns a []TargetedManagedAppPolicyAssignmentable when successful
func (m *WindowsInformationProtection) GetAssignments()([]TargetedManagedAppPolicyAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TargetedManagedAppPolicyAssignmentable)
    }
    return nil
}
// GetAzureRightsManagementServicesAllowed gets the azureRightsManagementServicesAllowed property value. Specifies whether to allow Azure RMS encryption for WIP
// returns a *bool when successful
func (m *WindowsInformationProtection) GetAzureRightsManagementServicesAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("azureRightsManagementServicesAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDataRecoveryCertificate gets the dataRecoveryCertificate property value. Specifies a recovery certificate that can be used for data recovery of encrypted files. This is the same as the data recovery agent(DRA) certificate for encrypting file system(EFS)
// returns a WindowsInformationProtectionDataRecoveryCertificateable when successful
func (m *WindowsInformationProtection) GetDataRecoveryCertificate()(WindowsInformationProtectionDataRecoveryCertificateable) {
    val, err := m.GetBackingStore().Get("dataRecoveryCertificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsInformationProtectionDataRecoveryCertificateable)
    }
    return nil
}
// GetEnforcementLevel gets the enforcementLevel property value. Possible values for WIP Protection enforcement levels
// returns a *WindowsInformationProtectionEnforcementLevel when successful
func (m *WindowsInformationProtection) GetEnforcementLevel()(*WindowsInformationProtectionEnforcementLevel) {
    val, err := m.GetBackingStore().Get("enforcementLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsInformationProtectionEnforcementLevel)
    }
    return nil
}
// GetEnterpriseDomain gets the enterpriseDomain property value. Primary enterprise domain
// returns a *string when successful
func (m *WindowsInformationProtection) GetEnterpriseDomain()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseDomain")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnterpriseInternalProxyServers gets the enterpriseInternalProxyServers property value. This is the comma-separated list of internal proxy servers. For example, '157.54.14.28, 157.54.11.118, 10.202.14.167, 157.53.14.163, 157.69.210.59'. These proxies have been configured by the admin to connect to specific resources on the Internet. They are considered to be enterprise network locations. The proxies are only leveraged in configuring the EnterpriseProxiedDomains policy to force traffic to the matched domains through these proxies
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseInternalProxyServers()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseInternalProxyServers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// GetEnterpriseIPRanges gets the enterpriseIPRanges property value. Sets the enterprise IP ranges that define the computers in the enterprise network. Data that comes from those computers will be considered part of the enterprise and protected. These locations will be considered a safe destination for enterprise data to be shared to
// returns a []WindowsInformationProtectionIPRangeCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseIPRanges()([]WindowsInformationProtectionIPRangeCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseIPRanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionIPRangeCollectionable)
    }
    return nil
}
// GetEnterpriseIPRangesAreAuthoritative gets the enterpriseIPRangesAreAuthoritative property value. Boolean value that tells the client to accept the configured list and not to use heuristics to attempt to find other subnets. Default is false
// returns a *bool when successful
func (m *WindowsInformationProtection) GetEnterpriseIPRangesAreAuthoritative()(*bool) {
    val, err := m.GetBackingStore().Get("enterpriseIPRangesAreAuthoritative")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnterpriseNetworkDomainNames gets the enterpriseNetworkDomainNames property value. This is the list of domains that comprise the boundaries of the enterprise. Data from one of these domains that is sent to a device will be considered enterprise data and protected These locations will be considered a safe destination for enterprise data to be shared to
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseNetworkDomainNames()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseNetworkDomainNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// GetEnterpriseProtectedDomainNames gets the enterpriseProtectedDomainNames property value. List of enterprise domains to be protected
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseProtectedDomainNames()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseProtectedDomainNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// GetEnterpriseProxiedDomains gets the enterpriseProxiedDomains property value. Contains a list of Enterprise resource domains hosted in the cloud that need to be protected. Connections to these resources are considered enterprise data. If a proxy is paired with a cloud resource, traffic to the cloud resource will be routed through the enterprise network via the denoted proxy server (on Port 80). A proxy server used for this purpose must also be configured using the EnterpriseInternalProxyServers policy
// returns a []WindowsInformationProtectionProxiedDomainCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseProxiedDomains()([]WindowsInformationProtectionProxiedDomainCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseProxiedDomains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionProxiedDomainCollectionable)
    }
    return nil
}
// GetEnterpriseProxyServers gets the enterpriseProxyServers property value. This is a list of proxy servers. Any server not on this list is considered non-enterprise
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetEnterpriseProxyServers()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("enterpriseProxyServers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// GetEnterpriseProxyServersAreAuthoritative gets the enterpriseProxyServersAreAuthoritative property value. Boolean value that tells the client to accept the configured list of proxies and not try to detect other work proxies. Default is false
// returns a *bool when successful
func (m *WindowsInformationProtection) GetEnterpriseProxyServersAreAuthoritative()(*bool) {
    val, err := m.GetBackingStore().Get("enterpriseProxyServersAreAuthoritative")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetExemptAppLockerFiles gets the exemptAppLockerFiles property value. Another way to input exempt apps through xml files
// returns a []WindowsInformationProtectionAppLockerFileable when successful
func (m *WindowsInformationProtection) GetExemptAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable) {
    val, err := m.GetBackingStore().Get("exemptAppLockerFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionAppLockerFileable)
    }
    return nil
}
// GetExemptApps gets the exemptApps property value. Exempt applications can also access enterprise data, but the data handled by those applications are not protected. This is because some critical enterprise applications may have compatibility problems with encrypted data.
// returns a []WindowsInformationProtectionAppable when successful
func (m *WindowsInformationProtection) GetExemptApps()([]WindowsInformationProtectionAppable) {
    val, err := m.GetBackingStore().Get("exemptApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionAppable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsInformationProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppPolicy.GetFieldDeserializers()
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTargetedManagedAppPolicyAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TargetedManagedAppPolicyAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TargetedManagedAppPolicyAssignmentable)
                }
            }
            m.SetAssignments(res)
        }
        return nil
    }
    res["azureRightsManagementServicesAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureRightsManagementServicesAllowed(val)
        }
        return nil
    }
    res["dataRecoveryCertificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsInformationProtectionDataRecoveryCertificateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataRecoveryCertificate(val.(WindowsInformationProtectionDataRecoveryCertificateable))
        }
        return nil
    }
    res["enforcementLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsInformationProtectionEnforcementLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnforcementLevel(val.(*WindowsInformationProtectionEnforcementLevel))
        }
        return nil
    }
    res["enterpriseDomain"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseDomain(val)
        }
        return nil
    }
    res["enterpriseInternalProxyServers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetEnterpriseInternalProxyServers(res)
        }
        return nil
    }
    res["enterpriseIPRanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionIPRangeCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionIPRangeCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionIPRangeCollectionable)
                }
            }
            m.SetEnterpriseIPRanges(res)
        }
        return nil
    }
    res["enterpriseIPRangesAreAuthoritative"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseIPRangesAreAuthoritative(val)
        }
        return nil
    }
    res["enterpriseNetworkDomainNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetEnterpriseNetworkDomainNames(res)
        }
        return nil
    }
    res["enterpriseProtectedDomainNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetEnterpriseProtectedDomainNames(res)
        }
        return nil
    }
    res["enterpriseProxiedDomains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionProxiedDomainCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionProxiedDomainCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionProxiedDomainCollectionable)
                }
            }
            m.SetEnterpriseProxiedDomains(res)
        }
        return nil
    }
    res["enterpriseProxyServers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetEnterpriseProxyServers(res)
        }
        return nil
    }
    res["enterpriseProxyServersAreAuthoritative"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseProxyServersAreAuthoritative(val)
        }
        return nil
    }
    res["exemptAppLockerFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionAppLockerFileable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionAppLockerFileable)
                }
            }
            m.SetExemptAppLockerFiles(res)
        }
        return nil
    }
    res["exemptApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionAppable)
                }
            }
            m.SetExemptApps(res)
        }
        return nil
    }
    res["iconsVisible"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIconsVisible(val)
        }
        return nil
    }
    res["indexingEncryptedStoresOrItemsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndexingEncryptedStoresOrItemsBlocked(val)
        }
        return nil
    }
    res["isAssigned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAssigned(val)
        }
        return nil
    }
    res["neutralDomainResources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetNeutralDomainResources(res)
        }
        return nil
    }
    res["protectedAppLockerFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionAppLockerFileable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionAppLockerFileable)
                }
            }
            m.SetProtectedAppLockerFiles(res)
        }
        return nil
    }
    res["protectedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionAppable)
                }
            }
            m.SetProtectedApps(res)
        }
        return nil
    }
    res["protectionUnderLockConfigRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtectionUnderLockConfigRequired(val)
        }
        return nil
    }
    res["revokeOnUnenrollDisabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRevokeOnUnenrollDisabled(val)
        }
        return nil
    }
    res["rightsManagementServicesTemplateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRightsManagementServicesTemplateId(val)
        }
        return nil
    }
    res["smbAutoEncryptedFileExtensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionResourceCollectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionResourceCollectionable)
                }
            }
            m.SetSmbAutoEncryptedFileExtensions(res)
        }
        return nil
    }
    return res
}
// GetIconsVisible gets the iconsVisible property value. Determines whether overlays are added to icons for WIP protected files in Explorer and enterprise only app tiles in the Start menu. Starting in Windows 10, version 1703 this setting also configures the visibility of the WIP icon in the title bar of a WIP-protected app
// returns a *bool when successful
func (m *WindowsInformationProtection) GetIconsVisible()(*bool) {
    val, err := m.GetBackingStore().Get("iconsVisible")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIndexingEncryptedStoresOrItemsBlocked gets the indexingEncryptedStoresOrItemsBlocked property value. This switch is for the Windows Search Indexer, to allow or disallow indexing of items
// returns a *bool when successful
func (m *WindowsInformationProtection) GetIndexingEncryptedStoresOrItemsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("indexingEncryptedStoresOrItemsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsAssigned gets the isAssigned property value. Indicates if the policy is deployed to any inclusion groups or not.
// returns a *bool when successful
func (m *WindowsInformationProtection) GetIsAssigned()(*bool) {
    val, err := m.GetBackingStore().Get("isAssigned")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNeutralDomainResources gets the neutralDomainResources property value. List of domain names that can used for work or personal resource
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetNeutralDomainResources()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("neutralDomainResources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// GetProtectedAppLockerFiles gets the protectedAppLockerFiles property value. Another way to input protected apps through xml files
// returns a []WindowsInformationProtectionAppLockerFileable when successful
func (m *WindowsInformationProtection) GetProtectedAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable) {
    val, err := m.GetBackingStore().Get("protectedAppLockerFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionAppLockerFileable)
    }
    return nil
}
// GetProtectedApps gets the protectedApps property value. Protected applications can access enterprise data and the data handled by those applications are protected with encryption
// returns a []WindowsInformationProtectionAppable when successful
func (m *WindowsInformationProtection) GetProtectedApps()([]WindowsInformationProtectionAppable) {
    val, err := m.GetBackingStore().Get("protectedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionAppable)
    }
    return nil
}
// GetProtectionUnderLockConfigRequired gets the protectionUnderLockConfigRequired property value. Specifies whether the protection under lock feature (also known as encrypt under pin) should be configured
// returns a *bool when successful
func (m *WindowsInformationProtection) GetProtectionUnderLockConfigRequired()(*bool) {
    val, err := m.GetBackingStore().Get("protectionUnderLockConfigRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRevokeOnUnenrollDisabled gets the revokeOnUnenrollDisabled property value. This policy controls whether to revoke the WIP keys when a device unenrolls from the management service. If set to 1 (Don't revoke keys), the keys will not be revoked and the user will continue to have access to protected files after unenrollment. If the keys are not revoked, there will be no revoked file cleanup subsequently.
// returns a *bool when successful
func (m *WindowsInformationProtection) GetRevokeOnUnenrollDisabled()(*bool) {
    val, err := m.GetBackingStore().Get("revokeOnUnenrollDisabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRightsManagementServicesTemplateId gets the rightsManagementServicesTemplateId property value. TemplateID GUID to use for RMS encryption. The RMS template allows the IT admin to configure the details about who has access to RMS-protected file and how long they have access
// returns a *UUID when successful
func (m *WindowsInformationProtection) GetRightsManagementServicesTemplateId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("rightsManagementServicesTemplateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetSmbAutoEncryptedFileExtensions gets the smbAutoEncryptedFileExtensions property value. Specifies a list of file extensions, so that files with these extensions are encrypted when copying from an SMB share within the corporate boundary
// returns a []WindowsInformationProtectionResourceCollectionable when successful
func (m *WindowsInformationProtection) GetSmbAutoEncryptedFileExtensions()([]WindowsInformationProtectionResourceCollectionable) {
    val, err := m.GetBackingStore().Get("smbAutoEncryptedFileExtensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionResourceCollectionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsInformationProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("azureRightsManagementServicesAllowed", m.GetAzureRightsManagementServicesAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dataRecoveryCertificate", m.GetDataRecoveryCertificate())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcementLevel() != nil {
        cast := (*m.GetEnforcementLevel()).String()
        err = writer.WriteStringValue("enforcementLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseDomain", m.GetEnterpriseDomain())
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseInternalProxyServers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseInternalProxyServers()))
        for i, v := range m.GetEnterpriseInternalProxyServers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseInternalProxyServers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseIPRanges() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseIPRanges()))
        for i, v := range m.GetEnterpriseIPRanges() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseIPRanges", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseIPRangesAreAuthoritative", m.GetEnterpriseIPRangesAreAuthoritative())
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseNetworkDomainNames() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseNetworkDomainNames()))
        for i, v := range m.GetEnterpriseNetworkDomainNames() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseNetworkDomainNames", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProtectedDomainNames() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseProtectedDomainNames()))
        for i, v := range m.GetEnterpriseProtectedDomainNames() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseProtectedDomainNames", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProxiedDomains() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseProxiedDomains()))
        for i, v := range m.GetEnterpriseProxiedDomains() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseProxiedDomains", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProxyServers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEnterpriseProxyServers()))
        for i, v := range m.GetEnterpriseProxyServers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("enterpriseProxyServers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseProxyServersAreAuthoritative", m.GetEnterpriseProxyServersAreAuthoritative())
        if err != nil {
            return err
        }
    }
    if m.GetExemptAppLockerFiles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExemptAppLockerFiles()))
        for i, v := range m.GetExemptAppLockerFiles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("exemptAppLockerFiles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExemptApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExemptApps()))
        for i, v := range m.GetExemptApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("exemptApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iconsVisible", m.GetIconsVisible())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("indexingEncryptedStoresOrItemsBlocked", m.GetIndexingEncryptedStoresOrItemsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAssigned", m.GetIsAssigned())
        if err != nil {
            return err
        }
    }
    if m.GetNeutralDomainResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNeutralDomainResources()))
        for i, v := range m.GetNeutralDomainResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("neutralDomainResources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectedAppLockerFiles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProtectedAppLockerFiles()))
        for i, v := range m.GetProtectedAppLockerFiles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("protectedAppLockerFiles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectedApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProtectedApps()))
        for i, v := range m.GetProtectedApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("protectedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("protectionUnderLockConfigRequired", m.GetProtectionUnderLockConfigRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("revokeOnUnenrollDisabled", m.GetRevokeOnUnenrollDisabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("rightsManagementServicesTemplateId", m.GetRightsManagementServicesTemplateId())
        if err != nil {
            return err
        }
    }
    if m.GetSmbAutoEncryptedFileExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSmbAutoEncryptedFileExtensions()))
        for i, v := range m.GetSmbAutoEncryptedFileExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("smbAutoEncryptedFileExtensions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignments sets the assignments property value. Navigation property to list of security groups targeted for policy.
func (m *WindowsInformationProtection) SetAssignments(value []TargetedManagedAppPolicyAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureRightsManagementServicesAllowed sets the azureRightsManagementServicesAllowed property value. Specifies whether to allow Azure RMS encryption for WIP
func (m *WindowsInformationProtection) SetAzureRightsManagementServicesAllowed(value *bool)() {
    err := m.GetBackingStore().Set("azureRightsManagementServicesAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetDataRecoveryCertificate sets the dataRecoveryCertificate property value. Specifies a recovery certificate that can be used for data recovery of encrypted files. This is the same as the data recovery agent(DRA) certificate for encrypting file system(EFS)
func (m *WindowsInformationProtection) SetDataRecoveryCertificate(value WindowsInformationProtectionDataRecoveryCertificateable)() {
    err := m.GetBackingStore().Set("dataRecoveryCertificate", value)
    if err != nil {
        panic(err)
    }
}
// SetEnforcementLevel sets the enforcementLevel property value. Possible values for WIP Protection enforcement levels
func (m *WindowsInformationProtection) SetEnforcementLevel(value *WindowsInformationProtectionEnforcementLevel)() {
    err := m.GetBackingStore().Set("enforcementLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseDomain sets the enterpriseDomain property value. Primary enterprise domain
func (m *WindowsInformationProtection) SetEnterpriseDomain(value *string)() {
    err := m.GetBackingStore().Set("enterpriseDomain", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseInternalProxyServers sets the enterpriseInternalProxyServers property value. This is the comma-separated list of internal proxy servers. For example, '157.54.14.28, 157.54.11.118, 10.202.14.167, 157.53.14.163, 157.69.210.59'. These proxies have been configured by the admin to connect to specific resources on the Internet. They are considered to be enterprise network locations. The proxies are only leveraged in configuring the EnterpriseProxiedDomains policy to force traffic to the matched domains through these proxies
func (m *WindowsInformationProtection) SetEnterpriseInternalProxyServers(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseInternalProxyServers", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseIPRanges sets the enterpriseIPRanges property value. Sets the enterprise IP ranges that define the computers in the enterprise network. Data that comes from those computers will be considered part of the enterprise and protected. These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) SetEnterpriseIPRanges(value []WindowsInformationProtectionIPRangeCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseIPRanges", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseIPRangesAreAuthoritative sets the enterpriseIPRangesAreAuthoritative property value. Boolean value that tells the client to accept the configured list and not to use heuristics to attempt to find other subnets. Default is false
func (m *WindowsInformationProtection) SetEnterpriseIPRangesAreAuthoritative(value *bool)() {
    err := m.GetBackingStore().Set("enterpriseIPRangesAreAuthoritative", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseNetworkDomainNames sets the enterpriseNetworkDomainNames property value. This is the list of domains that comprise the boundaries of the enterprise. Data from one of these domains that is sent to a device will be considered enterprise data and protected These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) SetEnterpriseNetworkDomainNames(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseNetworkDomainNames", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseProtectedDomainNames sets the enterpriseProtectedDomainNames property value. List of enterprise domains to be protected
func (m *WindowsInformationProtection) SetEnterpriseProtectedDomainNames(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseProtectedDomainNames", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseProxiedDomains sets the enterpriseProxiedDomains property value. Contains a list of Enterprise resource domains hosted in the cloud that need to be protected. Connections to these resources are considered enterprise data. If a proxy is paired with a cloud resource, traffic to the cloud resource will be routed through the enterprise network via the denoted proxy server (on Port 80). A proxy server used for this purpose must also be configured using the EnterpriseInternalProxyServers policy
func (m *WindowsInformationProtection) SetEnterpriseProxiedDomains(value []WindowsInformationProtectionProxiedDomainCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseProxiedDomains", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseProxyServers sets the enterpriseProxyServers property value. This is a list of proxy servers. Any server not on this list is considered non-enterprise
func (m *WindowsInformationProtection) SetEnterpriseProxyServers(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("enterpriseProxyServers", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseProxyServersAreAuthoritative sets the enterpriseProxyServersAreAuthoritative property value. Boolean value that tells the client to accept the configured list of proxies and not try to detect other work proxies. Default is false
func (m *WindowsInformationProtection) SetEnterpriseProxyServersAreAuthoritative(value *bool)() {
    err := m.GetBackingStore().Set("enterpriseProxyServersAreAuthoritative", value)
    if err != nil {
        panic(err)
    }
}
// SetExemptAppLockerFiles sets the exemptAppLockerFiles property value. Another way to input exempt apps through xml files
func (m *WindowsInformationProtection) SetExemptAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)() {
    err := m.GetBackingStore().Set("exemptAppLockerFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetExemptApps sets the exemptApps property value. Exempt applications can also access enterprise data, but the data handled by those applications are not protected. This is because some critical enterprise applications may have compatibility problems with encrypted data.
func (m *WindowsInformationProtection) SetExemptApps(value []WindowsInformationProtectionAppable)() {
    err := m.GetBackingStore().Set("exemptApps", value)
    if err != nil {
        panic(err)
    }
}
// SetIconsVisible sets the iconsVisible property value. Determines whether overlays are added to icons for WIP protected files in Explorer and enterprise only app tiles in the Start menu. Starting in Windows 10, version 1703 this setting also configures the visibility of the WIP icon in the title bar of a WIP-protected app
func (m *WindowsInformationProtection) SetIconsVisible(value *bool)() {
    err := m.GetBackingStore().Set("iconsVisible", value)
    if err != nil {
        panic(err)
    }
}
// SetIndexingEncryptedStoresOrItemsBlocked sets the indexingEncryptedStoresOrItemsBlocked property value. This switch is for the Windows Search Indexer, to allow or disallow indexing of items
func (m *WindowsInformationProtection) SetIndexingEncryptedStoresOrItemsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("indexingEncryptedStoresOrItemsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAssigned sets the isAssigned property value. Indicates if the policy is deployed to any inclusion groups or not.
func (m *WindowsInformationProtection) SetIsAssigned(value *bool)() {
    err := m.GetBackingStore().Set("isAssigned", value)
    if err != nil {
        panic(err)
    }
}
// SetNeutralDomainResources sets the neutralDomainResources property value. List of domain names that can used for work or personal resource
func (m *WindowsInformationProtection) SetNeutralDomainResources(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("neutralDomainResources", value)
    if err != nil {
        panic(err)
    }
}
// SetProtectedAppLockerFiles sets the protectedAppLockerFiles property value. Another way to input protected apps through xml files
func (m *WindowsInformationProtection) SetProtectedAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)() {
    err := m.GetBackingStore().Set("protectedAppLockerFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetProtectedApps sets the protectedApps property value. Protected applications can access enterprise data and the data handled by those applications are protected with encryption
func (m *WindowsInformationProtection) SetProtectedApps(value []WindowsInformationProtectionAppable)() {
    err := m.GetBackingStore().Set("protectedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetProtectionUnderLockConfigRequired sets the protectionUnderLockConfigRequired property value. Specifies whether the protection under lock feature (also known as encrypt under pin) should be configured
func (m *WindowsInformationProtection) SetProtectionUnderLockConfigRequired(value *bool)() {
    err := m.GetBackingStore().Set("protectionUnderLockConfigRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetRevokeOnUnenrollDisabled sets the revokeOnUnenrollDisabled property value. This policy controls whether to revoke the WIP keys when a device unenrolls from the management service. If set to 1 (Don't revoke keys), the keys will not be revoked and the user will continue to have access to protected files after unenrollment. If the keys are not revoked, there will be no revoked file cleanup subsequently.
func (m *WindowsInformationProtection) SetRevokeOnUnenrollDisabled(value *bool)() {
    err := m.GetBackingStore().Set("revokeOnUnenrollDisabled", value)
    if err != nil {
        panic(err)
    }
}
// SetRightsManagementServicesTemplateId sets the rightsManagementServicesTemplateId property value. TemplateID GUID to use for RMS encryption. The RMS template allows the IT admin to configure the details about who has access to RMS-protected file and how long they have access
func (m *WindowsInformationProtection) SetRightsManagementServicesTemplateId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("rightsManagementServicesTemplateId", value)
    if err != nil {
        panic(err)
    }
}
// SetSmbAutoEncryptedFileExtensions sets the smbAutoEncryptedFileExtensions property value. Specifies a list of file extensions, so that files with these extensions are encrypted when copying from an SMB share within the corporate boundary
func (m *WindowsInformationProtection) SetSmbAutoEncryptedFileExtensions(value []WindowsInformationProtectionResourceCollectionable)() {
    err := m.GetBackingStore().Set("smbAutoEncryptedFileExtensions", value)
    if err != nil {
        panic(err)
    }
}
type WindowsInformationProtectionable interface {
    ManagedAppPolicyable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignments()([]TargetedManagedAppPolicyAssignmentable)
    GetAzureRightsManagementServicesAllowed()(*bool)
    GetDataRecoveryCertificate()(WindowsInformationProtectionDataRecoveryCertificateable)
    GetEnforcementLevel()(*WindowsInformationProtectionEnforcementLevel)
    GetEnterpriseDomain()(*string)
    GetEnterpriseInternalProxyServers()([]WindowsInformationProtectionResourceCollectionable)
    GetEnterpriseIPRanges()([]WindowsInformationProtectionIPRangeCollectionable)
    GetEnterpriseIPRangesAreAuthoritative()(*bool)
    GetEnterpriseNetworkDomainNames()([]WindowsInformationProtectionResourceCollectionable)
    GetEnterpriseProtectedDomainNames()([]WindowsInformationProtectionResourceCollectionable)
    GetEnterpriseProxiedDomains()([]WindowsInformationProtectionProxiedDomainCollectionable)
    GetEnterpriseProxyServers()([]WindowsInformationProtectionResourceCollectionable)
    GetEnterpriseProxyServersAreAuthoritative()(*bool)
    GetExemptAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable)
    GetExemptApps()([]WindowsInformationProtectionAppable)
    GetIconsVisible()(*bool)
    GetIndexingEncryptedStoresOrItemsBlocked()(*bool)
    GetIsAssigned()(*bool)
    GetNeutralDomainResources()([]WindowsInformationProtectionResourceCollectionable)
    GetProtectedAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable)
    GetProtectedApps()([]WindowsInformationProtectionAppable)
    GetProtectionUnderLockConfigRequired()(*bool)
    GetRevokeOnUnenrollDisabled()(*bool)
    GetRightsManagementServicesTemplateId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetSmbAutoEncryptedFileExtensions()([]WindowsInformationProtectionResourceCollectionable)
    SetAssignments(value []TargetedManagedAppPolicyAssignmentable)()
    SetAzureRightsManagementServicesAllowed(value *bool)()
    SetDataRecoveryCertificate(value WindowsInformationProtectionDataRecoveryCertificateable)()
    SetEnforcementLevel(value *WindowsInformationProtectionEnforcementLevel)()
    SetEnterpriseDomain(value *string)()
    SetEnterpriseInternalProxyServers(value []WindowsInformationProtectionResourceCollectionable)()
    SetEnterpriseIPRanges(value []WindowsInformationProtectionIPRangeCollectionable)()
    SetEnterpriseIPRangesAreAuthoritative(value *bool)()
    SetEnterpriseNetworkDomainNames(value []WindowsInformationProtectionResourceCollectionable)()
    SetEnterpriseProtectedDomainNames(value []WindowsInformationProtectionResourceCollectionable)()
    SetEnterpriseProxiedDomains(value []WindowsInformationProtectionProxiedDomainCollectionable)()
    SetEnterpriseProxyServers(value []WindowsInformationProtectionResourceCollectionable)()
    SetEnterpriseProxyServersAreAuthoritative(value *bool)()
    SetExemptAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)()
    SetExemptApps(value []WindowsInformationProtectionAppable)()
    SetIconsVisible(value *bool)()
    SetIndexingEncryptedStoresOrItemsBlocked(value *bool)()
    SetIsAssigned(value *bool)()
    SetNeutralDomainResources(value []WindowsInformationProtectionResourceCollectionable)()
    SetProtectedAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)()
    SetProtectedApps(value []WindowsInformationProtectionAppable)()
    SetProtectionUnderLockConfigRequired(value *bool)()
    SetRevokeOnUnenrollDisabled(value *bool)()
    SetRightsManagementServicesTemplateId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetSmbAutoEncryptedFileExtensions(value []WindowsInformationProtectionResourceCollectionable)()
}

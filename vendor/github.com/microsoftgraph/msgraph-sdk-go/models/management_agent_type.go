package models
type ManagementAgentType int

const (
    // The device is managed by Exchange server.
    EAS_MANAGEMENTAGENTTYPE ManagementAgentType = iota
    // The device is managed by Intune MDM.
    MDM_MANAGEMENTAGENTTYPE
    // The device is managed by both Exchange server and Intune MDM.
    EASMDM_MANAGEMENTAGENTTYPE
    // Intune client managed.
    INTUNECLIENT_MANAGEMENTAGENTTYPE
    // The device is EAS and Intune client dual managed.
    EASINTUNECLIENT_MANAGEMENTAGENTTYPE
    // The device is managed by Configuration Manager.
    CONFIGURATIONMANAGERCLIENT_MANAGEMENTAGENTTYPE
    // The device is managed by Configuration Manager and MDM.
    CONFIGURATIONMANAGERCLIENTMDM_MANAGEMENTAGENTTYPE
    // The device is managed by Configuration Manager, MDM and Eas.
    CONFIGURATIONMANAGERCLIENTMDMEAS_MANAGEMENTAGENTTYPE
    // Unknown management agent type.
    UNKNOWN_MANAGEMENTAGENTTYPE
    // The device attributes are fetched from Jamf.
    JAMF_MANAGEMENTAGENTTYPE
    // The device is managed by Google's CloudDPC.
    GOOGLECLOUDDEVICEPOLICYCONTROLLER_MANAGEMENTAGENTTYPE
    // This device is managed by Microsoft 365 through Intune.
    MICROSOFT365MANAGEDMDM_MANAGEMENTAGENTTYPE
    MSSENSE_MANAGEMENTAGENTTYPE
)

func (i ManagementAgentType) String() string {
    return []string{"eas", "mdm", "easMdm", "intuneClient", "easIntuneClient", "configurationManagerClient", "configurationManagerClientMdm", "configurationManagerClientMdmEas", "unknown", "jamf", "googleCloudDevicePolicyController", "microsoft365ManagedMdm", "msSense"}[i]
}
func ParseManagementAgentType(v string) (any, error) {
    result := EAS_MANAGEMENTAGENTTYPE
    switch v {
        case "eas":
            result = EAS_MANAGEMENTAGENTTYPE
        case "mdm":
            result = MDM_MANAGEMENTAGENTTYPE
        case "easMdm":
            result = EASMDM_MANAGEMENTAGENTTYPE
        case "intuneClient":
            result = INTUNECLIENT_MANAGEMENTAGENTTYPE
        case "easIntuneClient":
            result = EASINTUNECLIENT_MANAGEMENTAGENTTYPE
        case "configurationManagerClient":
            result = CONFIGURATIONMANAGERCLIENT_MANAGEMENTAGENTTYPE
        case "configurationManagerClientMdm":
            result = CONFIGURATIONMANAGERCLIENTMDM_MANAGEMENTAGENTTYPE
        case "configurationManagerClientMdmEas":
            result = CONFIGURATIONMANAGERCLIENTMDMEAS_MANAGEMENTAGENTTYPE
        case "unknown":
            result = UNKNOWN_MANAGEMENTAGENTTYPE
        case "jamf":
            result = JAMF_MANAGEMENTAGENTTYPE
        case "googleCloudDevicePolicyController":
            result = GOOGLECLOUDDEVICEPOLICYCONTROLLER_MANAGEMENTAGENTTYPE
        case "microsoft365ManagedMdm":
            result = MICROSOFT365MANAGEDMDM_MANAGEMENTAGENTTYPE
        case "msSense":
            result = MSSENSE_MANAGEMENTAGENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagementAgentType(values []ManagementAgentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagementAgentType) isMultiValue() bool {
    return false
}

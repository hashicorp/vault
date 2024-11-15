package models
type SynchronizationMetadata int

const (
    GALLERYAPPLICATIONIDENTIFIER_SYNCHRONIZATIONMETADATA SynchronizationMetadata = iota
    GALLERYAPPLICATIONKEY_SYNCHRONIZATIONMETADATA
    ISOAUTHENABLED_SYNCHRONIZATIONMETADATA
    ISSYNCHRONIZATIONAGENTASSIGNMENTREQUIRED_SYNCHRONIZATIONMETADATA
    ISSYNCHRONIZATIONAGENTREQUIRED_SYNCHRONIZATIONMETADATA
    ISSYNCHRONIZATIONINPREVIEW_SYNCHRONIZATIONMETADATA
    OAUTHSETTINGS_SYNCHRONIZATIONMETADATA
    SYNCHRONIZATIONLEARNMOREIBIZAFWLINK_SYNCHRONIZATIONMETADATA
    CONFIGURATIONFIELDS_SYNCHRONIZATIONMETADATA
)

func (i SynchronizationMetadata) String() string {
    return []string{"GalleryApplicationIdentifier", "GalleryApplicationKey", "IsOAuthEnabled", "IsSynchronizationAgentAssignmentRequired", "IsSynchronizationAgentRequired", "IsSynchronizationInPreview", "OAuthSettings", "SynchronizationLearnMoreIbizaFwLink", "ConfigurationFields"}[i]
}
func ParseSynchronizationMetadata(v string) (any, error) {
    result := GALLERYAPPLICATIONIDENTIFIER_SYNCHRONIZATIONMETADATA
    switch v {
        case "GalleryApplicationIdentifier":
            result = GALLERYAPPLICATIONIDENTIFIER_SYNCHRONIZATIONMETADATA
        case "GalleryApplicationKey":
            result = GALLERYAPPLICATIONKEY_SYNCHRONIZATIONMETADATA
        case "IsOAuthEnabled":
            result = ISOAUTHENABLED_SYNCHRONIZATIONMETADATA
        case "IsSynchronizationAgentAssignmentRequired":
            result = ISSYNCHRONIZATIONAGENTASSIGNMENTREQUIRED_SYNCHRONIZATIONMETADATA
        case "IsSynchronizationAgentRequired":
            result = ISSYNCHRONIZATIONAGENTREQUIRED_SYNCHRONIZATIONMETADATA
        case "IsSynchronizationInPreview":
            result = ISSYNCHRONIZATIONINPREVIEW_SYNCHRONIZATIONMETADATA
        case "OAuthSettings":
            result = OAUTHSETTINGS_SYNCHRONIZATIONMETADATA
        case "SynchronizationLearnMoreIbizaFwLink":
            result = SYNCHRONIZATIONLEARNMOREIBIZAFWLINK_SYNCHRONIZATIONMETADATA
        case "ConfigurationFields":
            result = CONFIGURATIONFIELDS_SYNCHRONIZATIONMETADATA
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSynchronizationMetadata(values []SynchronizationMetadata) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SynchronizationMetadata) isMultiValue() bool {
    return false
}

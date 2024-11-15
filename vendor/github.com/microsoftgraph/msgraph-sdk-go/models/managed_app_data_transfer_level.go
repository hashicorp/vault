package models
// Data can be transferred from/to these classes of apps
type ManagedAppDataTransferLevel int

const (
    // All apps.
    ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL ManagedAppDataTransferLevel = iota
    // Managed apps.
    MANAGEDAPPS_MANAGEDAPPDATATRANSFERLEVEL
    // No apps.
    NONE_MANAGEDAPPDATATRANSFERLEVEL
)

func (i ManagedAppDataTransferLevel) String() string {
    return []string{"allApps", "managedApps", "none"}[i]
}
func ParseManagedAppDataTransferLevel(v string) (any, error) {
    result := ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL
    switch v {
        case "allApps":
            result = ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL
        case "managedApps":
            result = MANAGEDAPPS_MANAGEDAPPDATATRANSFERLEVEL
        case "none":
            result = NONE_MANAGEDAPPDATATRANSFERLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedAppDataTransferLevel(values []ManagedAppDataTransferLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedAppDataTransferLevel) isMultiValue() bool {
    return false
}

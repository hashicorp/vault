package models
// Represents the level to which the device's clipboard may be shared between apps
type ManagedAppClipboardSharingLevel int

const (
    // Sharing is allowed between all apps, managed or not
    ALLAPPS_MANAGEDAPPCLIPBOARDSHARINGLEVEL ManagedAppClipboardSharingLevel = iota
    // Sharing is allowed between all managed apps with paste in enabled
    MANAGEDAPPSWITHPASTEIN_MANAGEDAPPCLIPBOARDSHARINGLEVEL
    // Sharing is allowed between all managed apps
    MANAGEDAPPS_MANAGEDAPPCLIPBOARDSHARINGLEVEL
    // Sharing between apps is disabled
    BLOCKED_MANAGEDAPPCLIPBOARDSHARINGLEVEL
)

func (i ManagedAppClipboardSharingLevel) String() string {
    return []string{"allApps", "managedAppsWithPasteIn", "managedApps", "blocked"}[i]
}
func ParseManagedAppClipboardSharingLevel(v string) (any, error) {
    result := ALLAPPS_MANAGEDAPPCLIPBOARDSHARINGLEVEL
    switch v {
        case "allApps":
            result = ALLAPPS_MANAGEDAPPCLIPBOARDSHARINGLEVEL
        case "managedAppsWithPasteIn":
            result = MANAGEDAPPSWITHPASTEIN_MANAGEDAPPCLIPBOARDSHARINGLEVEL
        case "managedApps":
            result = MANAGEDAPPS_MANAGEDAPPCLIPBOARDSHARINGLEVEL
        case "blocked":
            result = BLOCKED_MANAGEDAPPCLIPBOARDSHARINGLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedAppClipboardSharingLevel(values []ManagedAppClipboardSharingLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedAppClipboardSharingLevel) isMultiValue() bool {
    return false
}

package models
// The reason for which a user has been flagged
type ManagedAppFlaggedReason int

const (
    // No issue.
    NONE_MANAGEDAPPFLAGGEDREASON ManagedAppFlaggedReason = iota
    // The app registration is running on a rooted/unlocked device.
    ROOTEDDEVICE_MANAGEDAPPFLAGGEDREASON
)

func (i ManagedAppFlaggedReason) String() string {
    return []string{"none", "rootedDevice"}[i]
}
func ParseManagedAppFlaggedReason(v string) (any, error) {
    result := NONE_MANAGEDAPPFLAGGEDREASON
    switch v {
        case "none":
            result = NONE_MANAGEDAPPFLAGGEDREASON
        case "rootedDevice":
            result = ROOTEDDEVICE_MANAGEDAPPFLAGGEDREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedAppFlaggedReason(values []ManagedAppFlaggedReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedAppFlaggedReason) isMultiValue() bool {
    return false
}

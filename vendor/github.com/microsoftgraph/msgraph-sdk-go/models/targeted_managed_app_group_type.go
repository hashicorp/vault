package models
// Indicates a collection of apps to target which can be one of several pre-defined lists of apps or a manually selected list of apps
type TargetedManagedAppGroupType int

const (
    // Target the collection of apps manually selected by the admin.
    SELECTEDPUBLICAPPS_TARGETEDMANAGEDAPPGROUPTYPE TargetedManagedAppGroupType = iota
    // Target the core set of Microsoft apps (Office, Edge, etc).
    ALLCOREMICROSOFTAPPS_TARGETEDMANAGEDAPPGROUPTYPE
    // Target all apps with Microsoft as publisher.
    ALLMICROSOFTAPPS_TARGETEDMANAGEDAPPGROUPTYPE
    // Target all apps with an available assignment.
    ALLAPPS_TARGETEDMANAGEDAPPGROUPTYPE
)

func (i TargetedManagedAppGroupType) String() string {
    return []string{"selectedPublicApps", "allCoreMicrosoftApps", "allMicrosoftApps", "allApps"}[i]
}
func ParseTargetedManagedAppGroupType(v string) (any, error) {
    result := SELECTEDPUBLICAPPS_TARGETEDMANAGEDAPPGROUPTYPE
    switch v {
        case "selectedPublicApps":
            result = SELECTEDPUBLICAPPS_TARGETEDMANAGEDAPPGROUPTYPE
        case "allCoreMicrosoftApps":
            result = ALLCOREMICROSOFTAPPS_TARGETEDMANAGEDAPPGROUPTYPE
        case "allMicrosoftApps":
            result = ALLMICROSOFTAPPS_TARGETEDMANAGEDAPPGROUPTYPE
        case "allApps":
            result = ALLAPPS_TARGETEDMANAGEDAPPGROUPTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTargetedManagedAppGroupType(values []TargetedManagedAppGroupType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TargetedManagedAppGroupType) isMultiValue() bool {
    return false
}

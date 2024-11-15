package models
// Android Work Profile cross profile data sharing type.
type AndroidWorkProfileCrossProfileDataSharingType int

const (
    // Device default value, no intent.
    DEVICEDEFAULT_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE AndroidWorkProfileCrossProfileDataSharingType = iota
    // Prevent any sharing.
    PREVENTANY_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
    // Allow data sharing request from personal profile to work profile.
    ALLOWPERSONALTOWORK_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
    // No restrictions on sharing.
    NORESTRICTIONS_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
)

func (i AndroidWorkProfileCrossProfileDataSharingType) String() string {
    return []string{"deviceDefault", "preventAny", "allowPersonalToWork", "noRestrictions"}[i]
}
func ParseAndroidWorkProfileCrossProfileDataSharingType(v string) (any, error) {
    result := DEVICEDEFAULT_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
    switch v {
        case "deviceDefault":
            result = DEVICEDEFAULT_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
        case "preventAny":
            result = PREVENTANY_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
        case "allowPersonalToWork":
            result = ALLOWPERSONALTOWORK_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
        case "noRestrictions":
            result = NORESTRICTIONS_ANDROIDWORKPROFILECROSSPROFILEDATASHARINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAndroidWorkProfileCrossProfileDataSharingType(values []AndroidWorkProfileCrossProfileDataSharingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AndroidWorkProfileCrossProfileDataSharingType) isMultiValue() bool {
    return false
}

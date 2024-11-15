package models
// Possible values of the compliance app list.
type AppListType int

const (
    // Default value, no intent.
    NONE_APPLISTTYPE AppListType = iota
    // The list represents the apps that will be considered compliant (only apps on the list are compliant).
    APPSINLISTCOMPLIANT_APPLISTTYPE
    // The list represents the apps that will be considered non compliant (all apps are compliant except apps on the list).
    APPSNOTINLISTCOMPLIANT_APPLISTTYPE
)

func (i AppListType) String() string {
    return []string{"none", "appsInListCompliant", "appsNotInListCompliant"}[i]
}
func ParseAppListType(v string) (any, error) {
    result := NONE_APPLISTTYPE
    switch v {
        case "none":
            result = NONE_APPLISTTYPE
        case "appsInListCompliant":
            result = APPSINLISTCOMPLIANT_APPLISTTYPE
        case "appsNotInListCompliant":
            result = APPSNOTINLISTCOMPLIANT_APPLISTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppListType(values []AppListType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppListType) isMultiValue() bool {
    return false
}

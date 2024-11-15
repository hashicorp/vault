package models
type PermissionType int

const (
    DELEGATEDUSERCONSENTABLE_PERMISSIONTYPE PermissionType = iota
    DELEGATED_PERMISSIONTYPE
    APPLICATION_PERMISSIONTYPE
)

func (i PermissionType) String() string {
    return []string{"delegatedUserConsentable", "delegated", "application"}[i]
}
func ParsePermissionType(v string) (any, error) {
    result := DELEGATEDUSERCONSENTABLE_PERMISSIONTYPE
    switch v {
        case "delegatedUserConsentable":
            result = DELEGATEDUSERCONSENTABLE_PERMISSIONTYPE
        case "delegated":
            result = DELEGATED_PERMISSIONTYPE
        case "application":
            result = APPLICATION_PERMISSIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePermissionType(values []PermissionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PermissionType) isMultiValue() bool {
    return false
}

package models
type IncludedUserRoles int

const (
    ALL_INCLUDEDUSERROLES IncludedUserRoles = iota
    PRIVILEGEDADMIN_INCLUDEDUSERROLES
    ADMIN_INCLUDEDUSERROLES
    USER_INCLUDEDUSERROLES
    UNKNOWNFUTUREVALUE_INCLUDEDUSERROLES
)

func (i IncludedUserRoles) String() string {
    return []string{"all", "privilegedAdmin", "admin", "user", "unknownFutureValue"}[i]
}
func ParseIncludedUserRoles(v string) (any, error) {
    result := ALL_INCLUDEDUSERROLES
    switch v {
        case "all":
            result = ALL_INCLUDEDUSERROLES
        case "privilegedAdmin":
            result = PRIVILEGEDADMIN_INCLUDEDUSERROLES
        case "admin":
            result = ADMIN_INCLUDEDUSERROLES
        case "user":
            result = USER_INCLUDEDUSERROLES
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INCLUDEDUSERROLES
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIncludedUserRoles(values []IncludedUserRoles) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IncludedUserRoles) isMultiValue() bool {
    return false
}

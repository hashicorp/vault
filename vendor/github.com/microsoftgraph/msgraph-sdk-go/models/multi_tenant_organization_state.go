package models
type MultiTenantOrganizationState int

const (
    ACTIVE_MULTITENANTORGANIZATIONSTATE MultiTenantOrganizationState = iota
    INACTIVE_MULTITENANTORGANIZATIONSTATE
    UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONSTATE
)

func (i MultiTenantOrganizationState) String() string {
    return []string{"active", "inactive", "unknownFutureValue"}[i]
}
func ParseMultiTenantOrganizationState(v string) (any, error) {
    result := ACTIVE_MULTITENANTORGANIZATIONSTATE
    switch v {
        case "active":
            result = ACTIVE_MULTITENANTORGANIZATIONSTATE
        case "inactive":
            result = INACTIVE_MULTITENANTORGANIZATIONSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMultiTenantOrganizationState(values []MultiTenantOrganizationState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MultiTenantOrganizationState) isMultiValue() bool {
    return false
}

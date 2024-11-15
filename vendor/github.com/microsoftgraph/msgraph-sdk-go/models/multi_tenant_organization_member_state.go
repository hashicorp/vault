package models
type MultiTenantOrganizationMemberState int

const (
    PENDING_MULTITENANTORGANIZATIONMEMBERSTATE MultiTenantOrganizationMemberState = iota
    ACTIVE_MULTITENANTORGANIZATIONMEMBERSTATE
    REMOVED_MULTITENANTORGANIZATIONMEMBERSTATE
    UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONMEMBERSTATE
)

func (i MultiTenantOrganizationMemberState) String() string {
    return []string{"pending", "active", "removed", "unknownFutureValue"}[i]
}
func ParseMultiTenantOrganizationMemberState(v string) (any, error) {
    result := PENDING_MULTITENANTORGANIZATIONMEMBERSTATE
    switch v {
        case "pending":
            result = PENDING_MULTITENANTORGANIZATIONMEMBERSTATE
        case "active":
            result = ACTIVE_MULTITENANTORGANIZATIONMEMBERSTATE
        case "removed":
            result = REMOVED_MULTITENANTORGANIZATIONMEMBERSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONMEMBERSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMultiTenantOrganizationMemberState(values []MultiTenantOrganizationMemberState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MultiTenantOrganizationMemberState) isMultiValue() bool {
    return false
}

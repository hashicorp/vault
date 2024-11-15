package models
type MultiTenantOrganizationMemberProcessingStatus int

const (
    NOTSTARTED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS MultiTenantOrganizationMemberProcessingStatus = iota
    RUNNING_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
    SUCCEEDED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
    FAILED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
    UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
)

func (i MultiTenantOrganizationMemberProcessingStatus) String() string {
    return []string{"notStarted", "running", "succeeded", "failed", "unknownFutureValue"}[i]
}
func ParseMultiTenantOrganizationMemberProcessingStatus(v string) (any, error) {
    result := NOTSTARTED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
        case "running":
            result = RUNNING_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
        case "succeeded":
            result = SUCCEEDED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
        case "failed":
            result = FAILED_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MULTITENANTORGANIZATIONMEMBERPROCESSINGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMultiTenantOrganizationMemberProcessingStatus(values []MultiTenantOrganizationMemberProcessingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MultiTenantOrganizationMemberProcessingStatus) isMultiValue() bool {
    return false
}

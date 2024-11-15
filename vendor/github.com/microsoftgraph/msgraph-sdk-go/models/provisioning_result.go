package models
type ProvisioningResult int

const (
    SUCCESS_PROVISIONINGRESULT ProvisioningResult = iota
    FAILURE_PROVISIONINGRESULT
    SKIPPED_PROVISIONINGRESULT
    WARNING_PROVISIONINGRESULT
    UNKNOWNFUTUREVALUE_PROVISIONINGRESULT
)

func (i ProvisioningResult) String() string {
    return []string{"success", "failure", "skipped", "warning", "unknownFutureValue"}[i]
}
func ParseProvisioningResult(v string) (any, error) {
    result := SUCCESS_PROVISIONINGRESULT
    switch v {
        case "success":
            result = SUCCESS_PROVISIONINGRESULT
        case "failure":
            result = FAILURE_PROVISIONINGRESULT
        case "skipped":
            result = SKIPPED_PROVISIONINGRESULT
        case "warning":
            result = WARNING_PROVISIONINGRESULT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROVISIONINGRESULT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProvisioningResult(values []ProvisioningResult) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProvisioningResult) isMultiValue() bool {
    return false
}

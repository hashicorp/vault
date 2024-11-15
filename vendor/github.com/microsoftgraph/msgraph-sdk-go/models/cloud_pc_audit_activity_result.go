package models
type CloudPcAuditActivityResult int

const (
    SUCCESS_CLOUDPCAUDITACTIVITYRESULT CloudPcAuditActivityResult = iota
    CLIENTERROR_CLOUDPCAUDITACTIVITYRESULT
    FAILURE_CLOUDPCAUDITACTIVITYRESULT
    TIMEOUT_CLOUDPCAUDITACTIVITYRESULT
    UNKNOWNFUTUREVALUE_CLOUDPCAUDITACTIVITYRESULT
)

func (i CloudPcAuditActivityResult) String() string {
    return []string{"success", "clientError", "failure", "timeout", "unknownFutureValue"}[i]
}
func ParseCloudPcAuditActivityResult(v string) (any, error) {
    result := SUCCESS_CLOUDPCAUDITACTIVITYRESULT
    switch v {
        case "success":
            result = SUCCESS_CLOUDPCAUDITACTIVITYRESULT
        case "clientError":
            result = CLIENTERROR_CLOUDPCAUDITACTIVITYRESULT
        case "failure":
            result = FAILURE_CLOUDPCAUDITACTIVITYRESULT
        case "timeout":
            result = TIMEOUT_CLOUDPCAUDITACTIVITYRESULT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCAUDITACTIVITYRESULT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcAuditActivityResult(values []CloudPcAuditActivityResult) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcAuditActivityResult) isMultiValue() bool {
    return false
}

package models
type OperationResult int

const (
    SUCCESS_OPERATIONRESULT OperationResult = iota
    FAILURE_OPERATIONRESULT
    TIMEOUT_OPERATIONRESULT
    UNKNOWNFUTUREVALUE_OPERATIONRESULT
)

func (i OperationResult) String() string {
    return []string{"success", "failure", "timeout", "unknownFutureValue"}[i]
}
func ParseOperationResult(v string) (any, error) {
    result := SUCCESS_OPERATIONRESULT
    switch v {
        case "success":
            result = SUCCESS_OPERATIONRESULT
        case "failure":
            result = FAILURE_OPERATIONRESULT
        case "timeout":
            result = TIMEOUT_OPERATIONRESULT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_OPERATIONRESULT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOperationResult(values []OperationResult) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OperationResult) isMultiValue() bool {
    return false
}

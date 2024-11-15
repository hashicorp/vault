package security
type CaseOperationStatus int

const (
    NOTSTARTED_CASEOPERATIONSTATUS CaseOperationStatus = iota
    SUBMISSIONFAILED_CASEOPERATIONSTATUS
    RUNNING_CASEOPERATIONSTATUS
    SUCCEEDED_CASEOPERATIONSTATUS
    PARTIALLYSUCCEEDED_CASEOPERATIONSTATUS
    FAILED_CASEOPERATIONSTATUS
    UNKNOWNFUTUREVALUE_CASEOPERATIONSTATUS
)

func (i CaseOperationStatus) String() string {
    return []string{"notStarted", "submissionFailed", "running", "succeeded", "partiallySucceeded", "failed", "unknownFutureValue"}[i]
}
func ParseCaseOperationStatus(v string) (any, error) {
    result := NOTSTARTED_CASEOPERATIONSTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_CASEOPERATIONSTATUS
        case "submissionFailed":
            result = SUBMISSIONFAILED_CASEOPERATIONSTATUS
        case "running":
            result = RUNNING_CASEOPERATIONSTATUS
        case "succeeded":
            result = SUCCEEDED_CASEOPERATIONSTATUS
        case "partiallySucceeded":
            result = PARTIALLYSUCCEEDED_CASEOPERATIONSTATUS
        case "failed":
            result = FAILED_CASEOPERATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CASEOPERATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCaseOperationStatus(values []CaseOperationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CaseOperationStatus) isMultiValue() bool {
    return false
}

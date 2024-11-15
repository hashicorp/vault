package models
type LongRunningOperationStatus int

const (
    NOTSTARTED_LONGRUNNINGOPERATIONSTATUS LongRunningOperationStatus = iota
    RUNNING_LONGRUNNINGOPERATIONSTATUS
    SUCCEEDED_LONGRUNNINGOPERATIONSTATUS
    FAILED_LONGRUNNINGOPERATIONSTATUS
    UNKNOWNFUTUREVALUE_LONGRUNNINGOPERATIONSTATUS
)

func (i LongRunningOperationStatus) String() string {
    return []string{"notStarted", "running", "succeeded", "failed", "unknownFutureValue"}[i]
}
func ParseLongRunningOperationStatus(v string) (any, error) {
    result := NOTSTARTED_LONGRUNNINGOPERATIONSTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_LONGRUNNINGOPERATIONSTATUS
        case "running":
            result = RUNNING_LONGRUNNINGOPERATIONSTATUS
        case "succeeded":
            result = SUCCEEDED_LONGRUNNINGOPERATIONSTATUS
        case "failed":
            result = FAILED_LONGRUNNINGOPERATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LONGRUNNINGOPERATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLongRunningOperationStatus(values []LongRunningOperationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LongRunningOperationStatus) isMultiValue() bool {
    return false
}

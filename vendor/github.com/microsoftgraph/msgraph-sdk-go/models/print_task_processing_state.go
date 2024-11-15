package models
type PrintTaskProcessingState int

const (
    PENDING_PRINTTASKPROCESSINGSTATE PrintTaskProcessingState = iota
    PROCESSING_PRINTTASKPROCESSINGSTATE
    COMPLETED_PRINTTASKPROCESSINGSTATE
    ABORTED_PRINTTASKPROCESSINGSTATE
    UNKNOWNFUTUREVALUE_PRINTTASKPROCESSINGSTATE
)

func (i PrintTaskProcessingState) String() string {
    return []string{"pending", "processing", "completed", "aborted", "unknownFutureValue"}[i]
}
func ParsePrintTaskProcessingState(v string) (any, error) {
    result := PENDING_PRINTTASKPROCESSINGSTATE
    switch v {
        case "pending":
            result = PENDING_PRINTTASKPROCESSINGSTATE
        case "processing":
            result = PROCESSING_PRINTTASKPROCESSINGSTATE
        case "completed":
            result = COMPLETED_PRINTTASKPROCESSINGSTATE
        case "aborted":
            result = ABORTED_PRINTTASKPROCESSINGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTTASKPROCESSINGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintTaskProcessingState(values []PrintTaskProcessingState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintTaskProcessingState) isMultiValue() bool {
    return false
}

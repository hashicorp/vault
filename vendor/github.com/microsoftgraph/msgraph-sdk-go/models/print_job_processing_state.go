package models
type PrintJobProcessingState int

const (
    UNKNOWN_PRINTJOBPROCESSINGSTATE PrintJobProcessingState = iota
    PENDING_PRINTJOBPROCESSINGSTATE
    PROCESSING_PRINTJOBPROCESSINGSTATE
    PAUSED_PRINTJOBPROCESSINGSTATE
    STOPPED_PRINTJOBPROCESSINGSTATE
    COMPLETED_PRINTJOBPROCESSINGSTATE
    CANCELED_PRINTJOBPROCESSINGSTATE
    ABORTED_PRINTJOBPROCESSINGSTATE
    UNKNOWNFUTUREVALUE_PRINTJOBPROCESSINGSTATE
)

func (i PrintJobProcessingState) String() string {
    return []string{"unknown", "pending", "processing", "paused", "stopped", "completed", "canceled", "aborted", "unknownFutureValue"}[i]
}
func ParsePrintJobProcessingState(v string) (any, error) {
    result := UNKNOWN_PRINTJOBPROCESSINGSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_PRINTJOBPROCESSINGSTATE
        case "pending":
            result = PENDING_PRINTJOBPROCESSINGSTATE
        case "processing":
            result = PROCESSING_PRINTJOBPROCESSINGSTATE
        case "paused":
            result = PAUSED_PRINTJOBPROCESSINGSTATE
        case "stopped":
            result = STOPPED_PRINTJOBPROCESSINGSTATE
        case "completed":
            result = COMPLETED_PRINTJOBPROCESSINGSTATE
        case "canceled":
            result = CANCELED_PRINTJOBPROCESSINGSTATE
        case "aborted":
            result = ABORTED_PRINTJOBPROCESSINGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTJOBPROCESSINGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintJobProcessingState(values []PrintJobProcessingState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintJobProcessingState) isMultiValue() bool {
    return false
}

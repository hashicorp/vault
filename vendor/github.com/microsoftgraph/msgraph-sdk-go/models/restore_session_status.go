package models
type RestoreSessionStatus int

const (
    DRAFT_RESTORESESSIONSTATUS RestoreSessionStatus = iota
    ACTIVATING_RESTORESESSIONSTATUS
    ACTIVE_RESTORESESSIONSTATUS
    COMPLETEDWITHERROR_RESTORESESSIONSTATUS
    COMPLETED_RESTORESESSIONSTATUS
    UNKNOWNFUTUREVALUE_RESTORESESSIONSTATUS
    FAILED_RESTORESESSIONSTATUS
)

func (i RestoreSessionStatus) String() string {
    return []string{"draft", "activating", "active", "completedWithError", "completed", "unknownFutureValue", "failed"}[i]
}
func ParseRestoreSessionStatus(v string) (any, error) {
    result := DRAFT_RESTORESESSIONSTATUS
    switch v {
        case "draft":
            result = DRAFT_RESTORESESSIONSTATUS
        case "activating":
            result = ACTIVATING_RESTORESESSIONSTATUS
        case "active":
            result = ACTIVE_RESTORESESSIONSTATUS
        case "completedWithError":
            result = COMPLETEDWITHERROR_RESTORESESSIONSTATUS
        case "completed":
            result = COMPLETED_RESTORESESSIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RESTORESESSIONSTATUS
        case "failed":
            result = FAILED_RESTORESESSIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRestoreSessionStatus(values []RestoreSessionStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RestoreSessionStatus) isMultiValue() bool {
    return false
}

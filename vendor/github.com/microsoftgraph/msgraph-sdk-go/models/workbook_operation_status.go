package models
type WorkbookOperationStatus int

const (
    NOTSTARTED_WORKBOOKOPERATIONSTATUS WorkbookOperationStatus = iota
    RUNNING_WORKBOOKOPERATIONSTATUS
    SUCCEEDED_WORKBOOKOPERATIONSTATUS
    FAILED_WORKBOOKOPERATIONSTATUS
)

func (i WorkbookOperationStatus) String() string {
    return []string{"notStarted", "running", "succeeded", "failed"}[i]
}
func ParseWorkbookOperationStatus(v string) (any, error) {
    result := NOTSTARTED_WORKBOOKOPERATIONSTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_WORKBOOKOPERATIONSTATUS
        case "running":
            result = RUNNING_WORKBOOKOPERATIONSTATUS
        case "succeeded":
            result = SUCCEEDED_WORKBOOKOPERATIONSTATUS
        case "failed":
            result = FAILED_WORKBOOKOPERATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWorkbookOperationStatus(values []WorkbookOperationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WorkbookOperationStatus) isMultiValue() bool {
    return false
}

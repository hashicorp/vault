package identitygovernance
type LifecycleWorkflowProcessingStatus int

const (
    QUEUED_LIFECYCLEWORKFLOWPROCESSINGSTATUS LifecycleWorkflowProcessingStatus = iota
    INPROGRESS_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    COMPLETED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    COMPLETEDWITHERRORS_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    CANCELED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    FAILED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    UNKNOWNFUTUREVALUE_LIFECYCLEWORKFLOWPROCESSINGSTATUS
)

func (i LifecycleWorkflowProcessingStatus) String() string {
    return []string{"queued", "inProgress", "completed", "completedWithErrors", "canceled", "failed", "unknownFutureValue"}[i]
}
func ParseLifecycleWorkflowProcessingStatus(v string) (any, error) {
    result := QUEUED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
    switch v {
        case "queued":
            result = QUEUED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "inProgress":
            result = INPROGRESS_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "completed":
            result = COMPLETED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "completedWithErrors":
            result = COMPLETEDWITHERRORS_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "canceled":
            result = CANCELED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "failed":
            result = FAILED_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LIFECYCLEWORKFLOWPROCESSINGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLifecycleWorkflowProcessingStatus(values []LifecycleWorkflowProcessingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LifecycleWorkflowProcessingStatus) isMultiValue() bool {
    return false
}

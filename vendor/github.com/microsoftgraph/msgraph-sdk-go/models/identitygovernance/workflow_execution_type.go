package identitygovernance
type WorkflowExecutionType int

const (
    SCHEDULED_WORKFLOWEXECUTIONTYPE WorkflowExecutionType = iota
    ONDEMAND_WORKFLOWEXECUTIONTYPE
    UNKNOWNFUTUREVALUE_WORKFLOWEXECUTIONTYPE
)

func (i WorkflowExecutionType) String() string {
    return []string{"scheduled", "onDemand", "unknownFutureValue"}[i]
}
func ParseWorkflowExecutionType(v string) (any, error) {
    result := SCHEDULED_WORKFLOWEXECUTIONTYPE
    switch v {
        case "scheduled":
            result = SCHEDULED_WORKFLOWEXECUTIONTYPE
        case "onDemand":
            result = ONDEMAND_WORKFLOWEXECUTIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WORKFLOWEXECUTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWorkflowExecutionType(values []WorkflowExecutionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WorkflowExecutionType) isMultiValue() bool {
    return false
}

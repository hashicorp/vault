package identitygovernance
type LifecycleWorkflowCategory int

const (
    JOINER_LIFECYCLEWORKFLOWCATEGORY LifecycleWorkflowCategory = iota
    LEAVER_LIFECYCLEWORKFLOWCATEGORY
    UNKNOWNFUTUREVALUE_LIFECYCLEWORKFLOWCATEGORY
    MOVER_LIFECYCLEWORKFLOWCATEGORY
)

func (i LifecycleWorkflowCategory) String() string {
    return []string{"joiner", "leaver", "unknownFutureValue", "mover"}[i]
}
func ParseLifecycleWorkflowCategory(v string) (any, error) {
    result := JOINER_LIFECYCLEWORKFLOWCATEGORY
    switch v {
        case "joiner":
            result = JOINER_LIFECYCLEWORKFLOWCATEGORY
        case "leaver":
            result = LEAVER_LIFECYCLEWORKFLOWCATEGORY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LIFECYCLEWORKFLOWCATEGORY
        case "mover":
            result = MOVER_LIFECYCLEWORKFLOWCATEGORY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLifecycleWorkflowCategory(values []LifecycleWorkflowCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LifecycleWorkflowCategory) isMultiValue() bool {
    return false
}

package models
type TrainingStatus int

const (
    UNKNOWN_TRAININGSTATUS TrainingStatus = iota
    ASSIGNED_TRAININGSTATUS
    INPROGRESS_TRAININGSTATUS
    COMPLETED_TRAININGSTATUS
    OVERDUE_TRAININGSTATUS
    UNKNOWNFUTUREVALUE_TRAININGSTATUS
)

func (i TrainingStatus) String() string {
    return []string{"unknown", "assigned", "inProgress", "completed", "overdue", "unknownFutureValue"}[i]
}
func ParseTrainingStatus(v string) (any, error) {
    result := UNKNOWN_TRAININGSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_TRAININGSTATUS
        case "assigned":
            result = ASSIGNED_TRAININGSTATUS
        case "inProgress":
            result = INPROGRESS_TRAININGSTATUS
        case "completed":
            result = COMPLETED_TRAININGSTATUS
        case "overdue":
            result = OVERDUE_TRAININGSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingStatus(values []TrainingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingStatus) isMultiValue() bool {
    return false
}

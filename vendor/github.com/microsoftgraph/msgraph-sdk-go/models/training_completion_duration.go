package models
type TrainingCompletionDuration int

const (
    WEEK_TRAININGCOMPLETIONDURATION TrainingCompletionDuration = iota
    FORTNITE_TRAININGCOMPLETIONDURATION
    MONTH_TRAININGCOMPLETIONDURATION
    UNKNOWNFUTUREVALUE_TRAININGCOMPLETIONDURATION
)

func (i TrainingCompletionDuration) String() string {
    return []string{"week", "fortnite", "month", "unknownFutureValue"}[i]
}
func ParseTrainingCompletionDuration(v string) (any, error) {
    result := WEEK_TRAININGCOMPLETIONDURATION
    switch v {
        case "week":
            result = WEEK_TRAININGCOMPLETIONDURATION
        case "fortnite":
            result = FORTNITE_TRAININGCOMPLETIONDURATION
        case "month":
            result = MONTH_TRAININGCOMPLETIONDURATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGCOMPLETIONDURATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingCompletionDuration(values []TrainingCompletionDuration) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingCompletionDuration) isMultiValue() bool {
    return false
}

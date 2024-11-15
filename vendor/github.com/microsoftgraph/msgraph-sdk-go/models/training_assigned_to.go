package models
type TrainingAssignedTo int

const (
    NONE_TRAININGASSIGNEDTO TrainingAssignedTo = iota
    ALLUSERS_TRAININGASSIGNEDTO
    CLICKEDPAYLOAD_TRAININGASSIGNEDTO
    COMPROMISED_TRAININGASSIGNEDTO
    REPORTEDPHISH_TRAININGASSIGNEDTO
    READBUTNOTCLICKED_TRAININGASSIGNEDTO
    DIDNOTHING_TRAININGASSIGNEDTO
    UNKNOWNFUTUREVALUE_TRAININGASSIGNEDTO
)

func (i TrainingAssignedTo) String() string {
    return []string{"none", "allUsers", "clickedPayload", "compromised", "reportedPhish", "readButNotClicked", "didNothing", "unknownFutureValue"}[i]
}
func ParseTrainingAssignedTo(v string) (any, error) {
    result := NONE_TRAININGASSIGNEDTO
    switch v {
        case "none":
            result = NONE_TRAININGASSIGNEDTO
        case "allUsers":
            result = ALLUSERS_TRAININGASSIGNEDTO
        case "clickedPayload":
            result = CLICKEDPAYLOAD_TRAININGASSIGNEDTO
        case "compromised":
            result = COMPROMISED_TRAININGASSIGNEDTO
        case "reportedPhish":
            result = REPORTEDPHISH_TRAININGASSIGNEDTO
        case "readButNotClicked":
            result = READBUTNOTCLICKED_TRAININGASSIGNEDTO
        case "didNothing":
            result = DIDNOTHING_TRAININGASSIGNEDTO
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGASSIGNEDTO
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingAssignedTo(values []TrainingAssignedTo) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingAssignedTo) isMultiValue() bool {
    return false
}

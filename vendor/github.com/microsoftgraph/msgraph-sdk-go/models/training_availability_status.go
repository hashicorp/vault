package models
type TrainingAvailabilityStatus int

const (
    UNKNOWN_TRAININGAVAILABILITYSTATUS TrainingAvailabilityStatus = iota
    NOTAVAILABLE_TRAININGAVAILABILITYSTATUS
    AVAILABLE_TRAININGAVAILABILITYSTATUS
    ARCHIVE_TRAININGAVAILABILITYSTATUS
    DELETE_TRAININGAVAILABILITYSTATUS
    UNKNOWNFUTUREVALUE_TRAININGAVAILABILITYSTATUS
)

func (i TrainingAvailabilityStatus) String() string {
    return []string{"unknown", "notAvailable", "available", "archive", "delete", "unknownFutureValue"}[i]
}
func ParseTrainingAvailabilityStatus(v string) (any, error) {
    result := UNKNOWN_TRAININGAVAILABILITYSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_TRAININGAVAILABILITYSTATUS
        case "notAvailable":
            result = NOTAVAILABLE_TRAININGAVAILABILITYSTATUS
        case "available":
            result = AVAILABLE_TRAININGAVAILABILITYSTATUS
        case "archive":
            result = ARCHIVE_TRAININGAVAILABILITYSTATUS
        case "delete":
            result = DELETE_TRAININGAVAILABILITYSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGAVAILABILITYSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingAvailabilityStatus(values []TrainingAvailabilityStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingAvailabilityStatus) isMultiValue() bool {
    return false
}

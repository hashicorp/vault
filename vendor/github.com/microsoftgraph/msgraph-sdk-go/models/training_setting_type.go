package models
type TrainingSettingType int

const (
    MICROSOFTCUSTOM_TRAININGSETTINGTYPE TrainingSettingType = iota
    MICROSOFTMANAGED_TRAININGSETTINGTYPE
    NOTRAINING_TRAININGSETTINGTYPE
    CUSTOM_TRAININGSETTINGTYPE
    UNKNOWNFUTUREVALUE_TRAININGSETTINGTYPE
)

func (i TrainingSettingType) String() string {
    return []string{"microsoftCustom", "microsoftManaged", "noTraining", "custom", "unknownFutureValue"}[i]
}
func ParseTrainingSettingType(v string) (any, error) {
    result := MICROSOFTCUSTOM_TRAININGSETTINGTYPE
    switch v {
        case "microsoftCustom":
            result = MICROSOFTCUSTOM_TRAININGSETTINGTYPE
        case "microsoftManaged":
            result = MICROSOFTMANAGED_TRAININGSETTINGTYPE
        case "noTraining":
            result = NOTRAINING_TRAININGSETTINGTYPE
        case "custom":
            result = CUSTOM_TRAININGSETTINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGSETTINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingSettingType(values []TrainingSettingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingSettingType) isMultiValue() bool {
    return false
}

package models
// Indicates severity of insights. Possible values are: None, Informational, Warning, Error.
type UserExperienceAnalyticsInsightSeverity int

const (
    // Indicates that the insight severity is none.
    NONE_USEREXPERIENCEANALYTICSINSIGHTSEVERITY UserExperienceAnalyticsInsightSeverity = iota
    // Indicates that the insight severity is informational.
    INFORMATIONAL_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
    // Indicates that the insight severity is warning.
    WARNING_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
    // Indicates that the insight severity is error.
    ERROR_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
)

func (i UserExperienceAnalyticsInsightSeverity) String() string {
    return []string{"none", "informational", "warning", "error", "unknownFutureValue"}[i]
}
func ParseUserExperienceAnalyticsInsightSeverity(v string) (any, error) {
    result := NONE_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
    switch v {
        case "none":
            result = NONE_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
        case "informational":
            result = INFORMATIONAL_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
        case "warning":
            result = WARNING_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
        case "error":
            result = ERROR_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSINSIGHTSEVERITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserExperienceAnalyticsInsightSeverity(values []UserExperienceAnalyticsInsightSeverity) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserExperienceAnalyticsInsightSeverity) isMultiValue() bool {
    return false
}

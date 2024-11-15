package models
type UserExperienceAnalyticsHealthState int

const (
    // Indicates that the health state is unknown.
    UNKNOWN_USEREXPERIENCEANALYTICSHEALTHSTATE UserExperienceAnalyticsHealthState = iota
    // Indicates that the health state is insufficient data.
    INSUFFICIENTDATA_USEREXPERIENCEANALYTICSHEALTHSTATE
    // Indicates that the health state needs attention.
    NEEDSATTENTION_USEREXPERIENCEANALYTICSHEALTHSTATE
    // Indicates that the health state is meeting goals.
    MEETINGGOALS_USEREXPERIENCEANALYTICSHEALTHSTATE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSHEALTHSTATE
)

func (i UserExperienceAnalyticsHealthState) String() string {
    return []string{"unknown", "insufficientData", "needsAttention", "meetingGoals", "unknownFutureValue"}[i]
}
func ParseUserExperienceAnalyticsHealthState(v string) (any, error) {
    result := UNKNOWN_USEREXPERIENCEANALYTICSHEALTHSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_USEREXPERIENCEANALYTICSHEALTHSTATE
        case "insufficientData":
            result = INSUFFICIENTDATA_USEREXPERIENCEANALYTICSHEALTHSTATE
        case "needsAttention":
            result = NEEDSATTENTION_USEREXPERIENCEANALYTICSHEALTHSTATE
        case "meetingGoals":
            result = MEETINGGOALS_USEREXPERIENCEANALYTICSHEALTHSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSHEALTHSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserExperienceAnalyticsHealthState(values []UserExperienceAnalyticsHealthState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserExperienceAnalyticsHealthState) isMultiValue() bool {
    return false
}

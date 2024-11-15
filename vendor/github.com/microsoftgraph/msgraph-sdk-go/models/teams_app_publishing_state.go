package models
type TeamsAppPublishingState int

const (
    SUBMITTED_TEAMSAPPPUBLISHINGSTATE TeamsAppPublishingState = iota
    REJECTED_TEAMSAPPPUBLISHINGSTATE
    PUBLISHED_TEAMSAPPPUBLISHINGSTATE
    UNKNOWNFUTUREVALUE_TEAMSAPPPUBLISHINGSTATE
)

func (i TeamsAppPublishingState) String() string {
    return []string{"submitted", "rejected", "published", "unknownFutureValue"}[i]
}
func ParseTeamsAppPublishingState(v string) (any, error) {
    result := SUBMITTED_TEAMSAPPPUBLISHINGSTATE
    switch v {
        case "submitted":
            result = SUBMITTED_TEAMSAPPPUBLISHINGSTATE
        case "rejected":
            result = REJECTED_TEAMSAPPPUBLISHINGSTATE
        case "published":
            result = PUBLISHED_TEAMSAPPPUBLISHINGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TEAMSAPPPUBLISHINGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamsAppPublishingState(values []TeamsAppPublishingState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamsAppPublishingState) isMultiValue() bool {
    return false
}

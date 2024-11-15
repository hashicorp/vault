package models
type MeetingAudience int

const (
    EVERYONE_MEETINGAUDIENCE MeetingAudience = iota
    ORGANIZATION_MEETINGAUDIENCE
    UNKNOWNFUTUREVALUE_MEETINGAUDIENCE
)

func (i MeetingAudience) String() string {
    return []string{"everyone", "organization", "unknownFutureValue"}[i]
}
func ParseMeetingAudience(v string) (any, error) {
    result := EVERYONE_MEETINGAUDIENCE
    switch v {
        case "everyone":
            result = EVERYONE_MEETINGAUDIENCE
        case "organization":
            result = ORGANIZATION_MEETINGAUDIENCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MEETINGAUDIENCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMeetingAudience(values []MeetingAudience) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MeetingAudience) isMultiValue() bool {
    return false
}

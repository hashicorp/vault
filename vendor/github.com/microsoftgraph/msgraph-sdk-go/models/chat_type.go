package models
type ChatType int

const (
    ONEONONE_CHATTYPE ChatType = iota
    GROUP_CHATTYPE
    MEETING_CHATTYPE
    UNKNOWNFUTUREVALUE_CHATTYPE
)

func (i ChatType) String() string {
    return []string{"oneOnOne", "group", "meeting", "unknownFutureValue"}[i]
}
func ParseChatType(v string) (any, error) {
    result := ONEONONE_CHATTYPE
    switch v {
        case "oneOnOne":
            result = ONEONONE_CHATTYPE
        case "group":
            result = GROUP_CHATTYPE
        case "meeting":
            result = MEETING_CHATTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CHATTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeChatType(values []ChatType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatType) isMultiValue() bool {
    return false
}

package models
type TeamworkConversationIdentityType int

const (
    TEAM_TEAMWORKCONVERSATIONIDENTITYTYPE TeamworkConversationIdentityType = iota
    CHANNEL_TEAMWORKCONVERSATIONIDENTITYTYPE
    CHAT_TEAMWORKCONVERSATIONIDENTITYTYPE
    UNKNOWNFUTUREVALUE_TEAMWORKCONVERSATIONIDENTITYTYPE
)

func (i TeamworkConversationIdentityType) String() string {
    return []string{"team", "channel", "chat", "unknownFutureValue"}[i]
}
func ParseTeamworkConversationIdentityType(v string) (any, error) {
    result := TEAM_TEAMWORKCONVERSATIONIDENTITYTYPE
    switch v {
        case "team":
            result = TEAM_TEAMWORKCONVERSATIONIDENTITYTYPE
        case "channel":
            result = CHANNEL_TEAMWORKCONVERSATIONIDENTITYTYPE
        case "chat":
            result = CHAT_TEAMWORKCONVERSATIONIDENTITYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TEAMWORKCONVERSATIONIDENTITYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamworkConversationIdentityType(values []TeamworkConversationIdentityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamworkConversationIdentityType) isMultiValue() bool {
    return false
}

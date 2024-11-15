package models
type ChatMessageImportance int

const (
    NORMAL_CHATMESSAGEIMPORTANCE ChatMessageImportance = iota
    HIGH_CHATMESSAGEIMPORTANCE
    URGENT_CHATMESSAGEIMPORTANCE
    UNKNOWNFUTUREVALUE_CHATMESSAGEIMPORTANCE
)

func (i ChatMessageImportance) String() string {
    return []string{"normal", "high", "urgent", "unknownFutureValue"}[i]
}
func ParseChatMessageImportance(v string) (any, error) {
    result := NORMAL_CHATMESSAGEIMPORTANCE
    switch v {
        case "normal":
            result = NORMAL_CHATMESSAGEIMPORTANCE
        case "high":
            result = HIGH_CHATMESSAGEIMPORTANCE
        case "urgent":
            result = URGENT_CHATMESSAGEIMPORTANCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CHATMESSAGEIMPORTANCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeChatMessageImportance(values []ChatMessageImportance) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessageImportance) isMultiValue() bool {
    return false
}

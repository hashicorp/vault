package models
type ChatMessageType int

const (
    MESSAGE_CHATMESSAGETYPE ChatMessageType = iota
    CHATEVENT_CHATMESSAGETYPE
    TYPING_CHATMESSAGETYPE
    UNKNOWNFUTUREVALUE_CHATMESSAGETYPE
    SYSTEMEVENTMESSAGE_CHATMESSAGETYPE
)

func (i ChatMessageType) String() string {
    return []string{"message", "chatEvent", "typing", "unknownFutureValue", "systemEventMessage"}[i]
}
func ParseChatMessageType(v string) (any, error) {
    result := MESSAGE_CHATMESSAGETYPE
    switch v {
        case "message":
            result = MESSAGE_CHATMESSAGETYPE
        case "chatEvent":
            result = CHATEVENT_CHATMESSAGETYPE
        case "typing":
            result = TYPING_CHATMESSAGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CHATMESSAGETYPE
        case "systemEventMessage":
            result = SYSTEMEVENTMESSAGE_CHATMESSAGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeChatMessageType(values []ChatMessageType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessageType) isMultiValue() bool {
    return false
}

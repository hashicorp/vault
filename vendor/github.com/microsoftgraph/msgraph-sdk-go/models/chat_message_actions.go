package models
import (
    "math"
    "strings"
)
type ChatMessageActions int

const (
    REACTIONADDED_CHATMESSAGEACTIONS = 1
    REACTIONREMOVED_CHATMESSAGEACTIONS = 2
    ACTIONUNDEFINED_CHATMESSAGEACTIONS = 4
    UNKNOWNFUTUREVALUE_CHATMESSAGEACTIONS = 8
)

func (i ChatMessageActions) String() string {
    var values []string
    options := []string{"reactionAdded", "reactionRemoved", "actionUndefined", "unknownFutureValue"}
    for p := 0; p < 4; p++ {
        mantis := ChatMessageActions(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseChatMessageActions(v string) (any, error) {
    var result ChatMessageActions
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "reactionAdded":
                result |= REACTIONADDED_CHATMESSAGEACTIONS
            case "reactionRemoved":
                result |= REACTIONREMOVED_CHATMESSAGEACTIONS
            case "actionUndefined":
                result |= ACTIONUNDEFINED_CHATMESSAGEACTIONS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_CHATMESSAGEACTIONS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeChatMessageActions(values []ChatMessageActions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessageActions) isMultiValue() bool {
    return true
}

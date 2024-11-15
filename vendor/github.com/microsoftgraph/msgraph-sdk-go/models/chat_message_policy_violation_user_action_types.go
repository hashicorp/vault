package models
import (
    "math"
    "strings"
)
type ChatMessagePolicyViolationUserActionTypes int

const (
    NONE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES = 1
    OVERRIDE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES = 2
    REPORTFALSEPOSITIVE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES = 4
)

func (i ChatMessagePolicyViolationUserActionTypes) String() string {
    var values []string
    options := []string{"none", "override", "reportFalsePositive"}
    for p := 0; p < 3; p++ {
        mantis := ChatMessagePolicyViolationUserActionTypes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseChatMessagePolicyViolationUserActionTypes(v string) (any, error) {
    var result ChatMessagePolicyViolationUserActionTypes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES
            case "override":
                result |= OVERRIDE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES
            case "reportFalsePositive":
                result |= REPORTFALSEPOSITIVE_CHATMESSAGEPOLICYVIOLATIONUSERACTIONTYPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeChatMessagePolicyViolationUserActionTypes(values []ChatMessagePolicyViolationUserActionTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessagePolicyViolationUserActionTypes) isMultiValue() bool {
    return true
}

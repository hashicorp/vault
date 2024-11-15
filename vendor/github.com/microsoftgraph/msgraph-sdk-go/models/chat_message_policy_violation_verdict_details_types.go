package models
import (
    "math"
    "strings"
)
type ChatMessagePolicyViolationVerdictDetailsTypes int

const (
    NONE_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES = 1
    ALLOWFALSEPOSITIVEOVERRIDE_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES = 2
    ALLOWOVERRIDEWITHOUTJUSTIFICATION_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES = 4
    ALLOWOVERRIDEWITHJUSTIFICATION_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES = 8
)

func (i ChatMessagePolicyViolationVerdictDetailsTypes) String() string {
    var values []string
    options := []string{"none", "allowFalsePositiveOverride", "allowOverrideWithoutJustification", "allowOverrideWithJustification"}
    for p := 0; p < 4; p++ {
        mantis := ChatMessagePolicyViolationVerdictDetailsTypes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseChatMessagePolicyViolationVerdictDetailsTypes(v string) (any, error) {
    var result ChatMessagePolicyViolationVerdictDetailsTypes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES
            case "allowFalsePositiveOverride":
                result |= ALLOWFALSEPOSITIVEOVERRIDE_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES
            case "allowOverrideWithoutJustification":
                result |= ALLOWOVERRIDEWITHOUTJUSTIFICATION_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES
            case "allowOverrideWithJustification":
                result |= ALLOWOVERRIDEWITHJUSTIFICATION_CHATMESSAGEPOLICYVIOLATIONVERDICTDETAILSTYPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeChatMessagePolicyViolationVerdictDetailsTypes(values []ChatMessagePolicyViolationVerdictDetailsTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessagePolicyViolationVerdictDetailsTypes) isMultiValue() bool {
    return true
}

package models
import (
    "math"
    "strings"
)
type ChatMessagePolicyViolationDlpActionTypes int

const (
    NONE_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES = 1
    NOTIFYSENDER_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES = 2
    BLOCKACCESS_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES = 4
    BLOCKACCESSEXTERNAL_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES = 8
)

func (i ChatMessagePolicyViolationDlpActionTypes) String() string {
    var values []string
    options := []string{"none", "notifySender", "blockAccess", "blockAccessExternal"}
    for p := 0; p < 4; p++ {
        mantis := ChatMessagePolicyViolationDlpActionTypes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseChatMessagePolicyViolationDlpActionTypes(v string) (any, error) {
    var result ChatMessagePolicyViolationDlpActionTypes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES
            case "notifySender":
                result |= NOTIFYSENDER_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES
            case "blockAccess":
                result |= BLOCKACCESS_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES
            case "blockAccessExternal":
                result |= BLOCKACCESSEXTERNAL_CHATMESSAGEPOLICYVIOLATIONDLPACTIONTYPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeChatMessagePolicyViolationDlpActionTypes(values []ChatMessagePolicyViolationDlpActionTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ChatMessagePolicyViolationDlpActionTypes) isMultiValue() bool {
    return true
}

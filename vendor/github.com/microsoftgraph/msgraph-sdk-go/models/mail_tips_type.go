package models
import (
    "math"
    "strings"
)
type MailTipsType int

const (
    AUTOMATICREPLIES_MAILTIPSTYPE = 1
    MAILBOXFULLSTATUS_MAILTIPSTYPE = 2
    CUSTOMMAILTIP_MAILTIPSTYPE = 4
    EXTERNALMEMBERCOUNT_MAILTIPSTYPE = 8
    TOTALMEMBERCOUNT_MAILTIPSTYPE = 16
    MAXMESSAGESIZE_MAILTIPSTYPE = 32
    DELIVERYRESTRICTION_MAILTIPSTYPE = 64
    MODERATIONSTATUS_MAILTIPSTYPE = 128
    RECIPIENTSCOPE_MAILTIPSTYPE = 256
    RECIPIENTSUGGESTIONS_MAILTIPSTYPE = 512
)

func (i MailTipsType) String() string {
    var values []string
    options := []string{"automaticReplies", "mailboxFullStatus", "customMailTip", "externalMemberCount", "totalMemberCount", "maxMessageSize", "deliveryRestriction", "moderationStatus", "recipientScope", "recipientSuggestions"}
    for p := 0; p < 10; p++ {
        mantis := MailTipsType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseMailTipsType(v string) (any, error) {
    var result MailTipsType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "automaticReplies":
                result |= AUTOMATICREPLIES_MAILTIPSTYPE
            case "mailboxFullStatus":
                result |= MAILBOXFULLSTATUS_MAILTIPSTYPE
            case "customMailTip":
                result |= CUSTOMMAILTIP_MAILTIPSTYPE
            case "externalMemberCount":
                result |= EXTERNALMEMBERCOUNT_MAILTIPSTYPE
            case "totalMemberCount":
                result |= TOTALMEMBERCOUNT_MAILTIPSTYPE
            case "maxMessageSize":
                result |= MAXMESSAGESIZE_MAILTIPSTYPE
            case "deliveryRestriction":
                result |= DELIVERYRESTRICTION_MAILTIPSTYPE
            case "moderationStatus":
                result |= MODERATIONSTATUS_MAILTIPSTYPE
            case "recipientScope":
                result |= RECIPIENTSCOPE_MAILTIPSTYPE
            case "recipientSuggestions":
                result |= RECIPIENTSUGGESTIONS_MAILTIPSTYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeMailTipsType(values []MailTipsType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MailTipsType) isMultiValue() bool {
    return true
}

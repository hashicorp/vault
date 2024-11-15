package security
import (
    "math"
    "strings"
)
type SourceType int

const (
    MAILBOX_SOURCETYPE = 1
    SITE_SOURCETYPE = 2
    UNKNOWNFUTUREVALUE_SOURCETYPE = 4
)

func (i SourceType) String() string {
    var values []string
    options := []string{"mailbox", "site", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := SourceType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseSourceType(v string) (any, error) {
    var result SourceType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "mailbox":
                result |= MAILBOX_SOURCETYPE
            case "site":
                result |= SITE_SOURCETYPE
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_SOURCETYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeSourceType(values []SourceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SourceType) isMultiValue() bool {
    return true
}

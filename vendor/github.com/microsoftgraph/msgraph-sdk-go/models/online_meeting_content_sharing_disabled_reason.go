package models
import (
    "math"
    "strings"
)
type OnlineMeetingContentSharingDisabledReason int

const (
    WATERMARKPROTECTION_ONLINEMEETINGCONTENTSHARINGDISABLEDREASON = 1
    UNKNOWNFUTUREVALUE_ONLINEMEETINGCONTENTSHARINGDISABLEDREASON = 2
)

func (i OnlineMeetingContentSharingDisabledReason) String() string {
    var values []string
    options := []string{"watermarkProtection", "unknownFutureValue"}
    for p := 0; p < 2; p++ {
        mantis := OnlineMeetingContentSharingDisabledReason(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseOnlineMeetingContentSharingDisabledReason(v string) (any, error) {
    var result OnlineMeetingContentSharingDisabledReason
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "watermarkProtection":
                result |= WATERMARKPROTECTION_ONLINEMEETINGCONTENTSHARINGDISABLEDREASON
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_ONLINEMEETINGCONTENTSHARINGDISABLEDREASON
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeOnlineMeetingContentSharingDisabledReason(values []OnlineMeetingContentSharingDisabledReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnlineMeetingContentSharingDisabledReason) isMultiValue() bool {
    return true
}

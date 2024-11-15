package models
import (
    "math"
    "strings"
)
type OnlineMeetingVideoDisabledReason int

const (
    WATERMARKPROTECTION_ONLINEMEETINGVIDEODISABLEDREASON = 1
    UNKNOWNFUTUREVALUE_ONLINEMEETINGVIDEODISABLEDREASON = 2
)

func (i OnlineMeetingVideoDisabledReason) String() string {
    var values []string
    options := []string{"watermarkProtection", "unknownFutureValue"}
    for p := 0; p < 2; p++ {
        mantis := OnlineMeetingVideoDisabledReason(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseOnlineMeetingVideoDisabledReason(v string) (any, error) {
    var result OnlineMeetingVideoDisabledReason
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "watermarkProtection":
                result |= WATERMARKPROTECTION_ONLINEMEETINGVIDEODISABLEDREASON
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_ONLINEMEETINGVIDEODISABLEDREASON
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeOnlineMeetingVideoDisabledReason(values []OnlineMeetingVideoDisabledReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnlineMeetingVideoDisabledReason) isMultiValue() bool {
    return true
}

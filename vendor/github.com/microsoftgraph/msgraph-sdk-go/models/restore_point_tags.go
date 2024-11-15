package models
import (
    "math"
    "strings"
)
type RestorePointTags int

const (
    NONE_RESTOREPOINTTAGS = 1
    FASTRESTORE_RESTOREPOINTTAGS = 2
    UNKNOWNFUTUREVALUE_RESTOREPOINTTAGS = 4
)

func (i RestorePointTags) String() string {
    var values []string
    options := []string{"none", "fastRestore", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := RestorePointTags(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseRestorePointTags(v string) (any, error) {
    var result RestorePointTags
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_RESTOREPOINTTAGS
            case "fastRestore":
                result |= FASTRESTORE_RESTOREPOINTTAGS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_RESTOREPOINTTAGS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeRestorePointTags(values []RestorePointTags) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RestorePointTags) isMultiValue() bool {
    return true
}

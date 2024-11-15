package models
import (
    "math"
    "strings"
)
type WeakAlgorithms int

const (
    RSASHA1_WEAKALGORITHMS = 1
    UNKNOWNFUTUREVALUE_WEAKALGORITHMS = 2
)

func (i WeakAlgorithms) String() string {
    var values []string
    options := []string{"rsaSha1", "unknownFutureValue"}
    for p := 0; p < 2; p++ {
        mantis := WeakAlgorithms(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWeakAlgorithms(v string) (any, error) {
    var result WeakAlgorithms
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "rsaSha1":
                result |= RSASHA1_WEAKALGORITHMS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_WEAKALGORITHMS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWeakAlgorithms(values []WeakAlgorithms) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WeakAlgorithms) isMultiValue() bool {
    return true
}

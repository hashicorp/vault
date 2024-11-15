package models
import (
    "math"
    "strings"
)
type ConditionalAccessInsiderRiskLevels int

const (
    MINOR_CONDITIONALACCESSINSIDERRISKLEVELS = 1
    MODERATE_CONDITIONALACCESSINSIDERRISKLEVELS = 2
    ELEVATED_CONDITIONALACCESSINSIDERRISKLEVELS = 4
    UNKNOWNFUTUREVALUE_CONDITIONALACCESSINSIDERRISKLEVELS = 8
)

func (i ConditionalAccessInsiderRiskLevels) String() string {
    var values []string
    options := []string{"minor", "moderate", "elevated", "unknownFutureValue"}
    for p := 0; p < 4; p++ {
        mantis := ConditionalAccessInsiderRiskLevels(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseConditionalAccessInsiderRiskLevels(v string) (any, error) {
    var result ConditionalAccessInsiderRiskLevels
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "minor":
                result |= MINOR_CONDITIONALACCESSINSIDERRISKLEVELS
            case "moderate":
                result |= MODERATE_CONDITIONALACCESSINSIDERRISKLEVELS
            case "elevated":
                result |= ELEVATED_CONDITIONALACCESSINSIDERRISKLEVELS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_CONDITIONALACCESSINSIDERRISKLEVELS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeConditionalAccessInsiderRiskLevels(values []ConditionalAccessInsiderRiskLevels) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConditionalAccessInsiderRiskLevels) isMultiValue() bool {
    return true
}

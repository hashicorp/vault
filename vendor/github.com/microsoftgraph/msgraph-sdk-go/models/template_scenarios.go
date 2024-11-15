package models
import (
    "math"
    "strings"
)
type TemplateScenarios int

const (
    NEW_TEMPLATESCENARIOS = 1
    SECUREFOUNDATION_TEMPLATESCENARIOS = 2
    ZEROTRUST_TEMPLATESCENARIOS = 4
    REMOTEWORK_TEMPLATESCENARIOS = 8
    PROTECTADMINS_TEMPLATESCENARIOS = 16
    EMERGINGTHREATS_TEMPLATESCENARIOS = 32
    UNKNOWNFUTUREVALUE_TEMPLATESCENARIOS = 64
)

func (i TemplateScenarios) String() string {
    var values []string
    options := []string{"new", "secureFoundation", "zeroTrust", "remoteWork", "protectAdmins", "emergingThreats", "unknownFutureValue"}
    for p := 0; p < 7; p++ {
        mantis := TemplateScenarios(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseTemplateScenarios(v string) (any, error) {
    var result TemplateScenarios
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "new":
                result |= NEW_TEMPLATESCENARIOS
            case "secureFoundation":
                result |= SECUREFOUNDATION_TEMPLATESCENARIOS
            case "zeroTrust":
                result |= ZEROTRUST_TEMPLATESCENARIOS
            case "remoteWork":
                result |= REMOTEWORK_TEMPLATESCENARIOS
            case "protectAdmins":
                result |= PROTECTADMINS_TEMPLATESCENARIOS
            case "emergingThreats":
                result |= EMERGINGTHREATS_TEMPLATESCENARIOS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_TEMPLATESCENARIOS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeTemplateScenarios(values []TemplateScenarios) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TemplateScenarios) isMultiValue() bool {
    return true
}

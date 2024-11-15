package models
import (
    "math"
    "strings"
)
type TemplateApplicationLevel int

const (
    NONE_TEMPLATEAPPLICATIONLEVEL = 1
    NEWPARTNERS_TEMPLATEAPPLICATIONLEVEL = 2
    EXISTINGPARTNERS_TEMPLATEAPPLICATIONLEVEL = 4
    UNKNOWNFUTUREVALUE_TEMPLATEAPPLICATIONLEVEL = 8
)

func (i TemplateApplicationLevel) String() string {
    var values []string
    options := []string{"none", "newPartners", "existingPartners", "unknownFutureValue"}
    for p := 0; p < 4; p++ {
        mantis := TemplateApplicationLevel(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseTemplateApplicationLevel(v string) (any, error) {
    var result TemplateApplicationLevel
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_TEMPLATEAPPLICATIONLEVEL
            case "newPartners":
                result |= NEWPARTNERS_TEMPLATEAPPLICATIONLEVEL
            case "existingPartners":
                result |= EXISTINGPARTNERS_TEMPLATEAPPLICATIONLEVEL
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_TEMPLATEAPPLICATIONLEVEL
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeTemplateApplicationLevel(values []TemplateApplicationLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TemplateApplicationLevel) isMultiValue() bool {
    return true
}

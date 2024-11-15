package security
import (
    "math"
    "strings"
)
type AdditionalDataOptions int

const (
    ALLVERSIONS_ADDITIONALDATAOPTIONS = 1
    LINKEDFILES_ADDITIONALDATAOPTIONS = 2
    UNKNOWNFUTUREVALUE_ADDITIONALDATAOPTIONS = 4
)

func (i AdditionalDataOptions) String() string {
    var values []string
    options := []string{"allVersions", "linkedFiles", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := AdditionalDataOptions(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseAdditionalDataOptions(v string) (any, error) {
    var result AdditionalDataOptions
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "allVersions":
                result |= ALLVERSIONS_ADDITIONALDATAOPTIONS
            case "linkedFiles":
                result |= LINKEDFILES_ADDITIONALDATAOPTIONS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_ADDITIONALDATAOPTIONS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeAdditionalDataOptions(values []AdditionalDataOptions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AdditionalDataOptions) isMultiValue() bool {
    return true
}

package security
import (
    "math"
    "strings"
)
type ExportCriteria int

const (
    SEARCHHITS_EXPORTCRITERIA = 1
    PARTIALLYINDEXED_EXPORTCRITERIA = 2
    UNKNOWNFUTUREVALUE_EXPORTCRITERIA = 4
)

func (i ExportCriteria) String() string {
    var values []string
    options := []string{"searchHits", "partiallyIndexed", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := ExportCriteria(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseExportCriteria(v string) (any, error) {
    var result ExportCriteria
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "searchHits":
                result |= SEARCHHITS_EXPORTCRITERIA
            case "partiallyIndexed":
                result |= PARTIALLYINDEXED_EXPORTCRITERIA
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_EXPORTCRITERIA
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeExportCriteria(values []ExportCriteria) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExportCriteria) isMultiValue() bool {
    return true
}

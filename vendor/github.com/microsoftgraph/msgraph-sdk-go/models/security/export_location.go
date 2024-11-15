package security
import (
    "math"
    "strings"
)
type ExportLocation int

const (
    RESPONSIVELOCATIONS_EXPORTLOCATION = 1
    NONRESPONSIVELOCATIONS_EXPORTLOCATION = 2
    UNKNOWNFUTUREVALUE_EXPORTLOCATION = 4
)

func (i ExportLocation) String() string {
    var values []string
    options := []string{"responsiveLocations", "nonresponsiveLocations", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := ExportLocation(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseExportLocation(v string) (any, error) {
    var result ExportLocation
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "responsiveLocations":
                result |= RESPONSIVELOCATIONS_EXPORTLOCATION
            case "nonresponsiveLocations":
                result |= NONRESPONSIVELOCATIONS_EXPORTLOCATION
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_EXPORTLOCATION
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeExportLocation(values []ExportLocation) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExportLocation) isMultiValue() bool {
    return true
}

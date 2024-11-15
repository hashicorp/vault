package security
import (
    "math"
    "strings"
)
type ExportOptions int

const (
    ORIGINALFILES_EXPORTOPTIONS = 1
    TEXT_EXPORTOPTIONS = 2
    PDFREPLACEMENT_EXPORTOPTIONS = 4
    TAGS_EXPORTOPTIONS = 8
    UNKNOWNFUTUREVALUE_EXPORTOPTIONS = 16
)

func (i ExportOptions) String() string {
    var values []string
    options := []string{"originalFiles", "text", "pdfReplacement", "tags", "unknownFutureValue"}
    for p := 0; p < 5; p++ {
        mantis := ExportOptions(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseExportOptions(v string) (any, error) {
    var result ExportOptions
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "originalFiles":
                result |= ORIGINALFILES_EXPORTOPTIONS
            case "text":
                result |= TEXT_EXPORTOPTIONS
            case "pdfReplacement":
                result |= PDFREPLACEMENT_EXPORTOPTIONS
            case "tags":
                result |= TAGS_EXPORTOPTIONS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_EXPORTOPTIONS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeExportOptions(values []ExportOptions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExportOptions) isMultiValue() bool {
    return true
}

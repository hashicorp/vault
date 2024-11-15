package models
import (
    "math"
    "strings"
)
type SearchContent int

const (
    SHAREDCONTENT_SEARCHCONTENT = 1
    PRIVATECONTENT_SEARCHCONTENT = 2
    UNKNOWNFUTUREVALUE_SEARCHCONTENT = 4
)

func (i SearchContent) String() string {
    var values []string
    options := []string{"sharedContent", "privateContent", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := SearchContent(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseSearchContent(v string) (any, error) {
    var result SearchContent
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "sharedContent":
                result |= SHAREDCONTENT_SEARCHCONTENT
            case "privateContent":
                result |= PRIVATECONTENT_SEARCHCONTENT
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_SEARCHCONTENT
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeSearchContent(values []SearchContent) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SearchContent) isMultiValue() bool {
    return true
}

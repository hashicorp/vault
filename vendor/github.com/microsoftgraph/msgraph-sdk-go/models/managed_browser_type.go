package models
import (
    "math"
    "strings"
)
// Type of managed browser
type ManagedBrowserType int

const (
    // Not configured
    NOTCONFIGURED_MANAGEDBROWSERTYPE = 1
    // Microsoft Edge
    MICROSOFTEDGE_MANAGEDBROWSERTYPE = 2
)

func (i ManagedBrowserType) String() string {
    var values []string
    options := []string{"notConfigured", "microsoftEdge"}
    for p := 0; p < 2; p++ {
        mantis := ManagedBrowserType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseManagedBrowserType(v string) (any, error) {
    var result ManagedBrowserType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "notConfigured":
                result |= NOTCONFIGURED_MANAGEDBROWSERTYPE
            case "microsoftEdge":
                result |= MICROSOFTEDGE_MANAGEDBROWSERTYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeManagedBrowserType(values []ManagedBrowserType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedBrowserType) isMultiValue() bool {
    return true
}

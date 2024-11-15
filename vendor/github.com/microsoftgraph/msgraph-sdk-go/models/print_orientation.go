package models
type PrintOrientation int

const (
    PORTRAIT_PRINTORIENTATION PrintOrientation = iota
    LANDSCAPE_PRINTORIENTATION
    REVERSELANDSCAPE_PRINTORIENTATION
    REVERSEPORTRAIT_PRINTORIENTATION
    UNKNOWNFUTUREVALUE_PRINTORIENTATION
)

func (i PrintOrientation) String() string {
    return []string{"portrait", "landscape", "reverseLandscape", "reversePortrait", "unknownFutureValue"}[i]
}
func ParsePrintOrientation(v string) (any, error) {
    result := PORTRAIT_PRINTORIENTATION
    switch v {
        case "portrait":
            result = PORTRAIT_PRINTORIENTATION
        case "landscape":
            result = LANDSCAPE_PRINTORIENTATION
        case "reverseLandscape":
            result = REVERSELANDSCAPE_PRINTORIENTATION
        case "reversePortrait":
            result = REVERSEPORTRAIT_PRINTORIENTATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTORIENTATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintOrientation(values []PrintOrientation) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintOrientation) isMultiValue() bool {
    return false
}

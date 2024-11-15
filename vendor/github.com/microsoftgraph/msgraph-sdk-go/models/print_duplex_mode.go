package models
type PrintDuplexMode int

const (
    FLIPONLONGEDGE_PRINTDUPLEXMODE PrintDuplexMode = iota
    FLIPONSHORTEDGE_PRINTDUPLEXMODE
    ONESIDED_PRINTDUPLEXMODE
    UNKNOWNFUTUREVALUE_PRINTDUPLEXMODE
)

func (i PrintDuplexMode) String() string {
    return []string{"flipOnLongEdge", "flipOnShortEdge", "oneSided", "unknownFutureValue"}[i]
}
func ParsePrintDuplexMode(v string) (any, error) {
    result := FLIPONLONGEDGE_PRINTDUPLEXMODE
    switch v {
        case "flipOnLongEdge":
            result = FLIPONLONGEDGE_PRINTDUPLEXMODE
        case "flipOnShortEdge":
            result = FLIPONSHORTEDGE_PRINTDUPLEXMODE
        case "oneSided":
            result = ONESIDED_PRINTDUPLEXMODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTDUPLEXMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintDuplexMode(values []PrintDuplexMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintDuplexMode) isMultiValue() bool {
    return false
}

package models
type PrintScaling int

const (
    AUTO_PRINTSCALING PrintScaling = iota
    SHRINKTOFIT_PRINTSCALING
    FILL_PRINTSCALING
    FIT_PRINTSCALING
    NONE_PRINTSCALING
    UNKNOWNFUTUREVALUE_PRINTSCALING
)

func (i PrintScaling) String() string {
    return []string{"auto", "shrinkToFit", "fill", "fit", "none", "unknownFutureValue"}[i]
}
func ParsePrintScaling(v string) (any, error) {
    result := AUTO_PRINTSCALING
    switch v {
        case "auto":
            result = AUTO_PRINTSCALING
        case "shrinkToFit":
            result = SHRINKTOFIT_PRINTSCALING
        case "fill":
            result = FILL_PRINTSCALING
        case "fit":
            result = FIT_PRINTSCALING
        case "none":
            result = NONE_PRINTSCALING
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTSCALING
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintScaling(values []PrintScaling) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintScaling) isMultiValue() bool {
    return false
}

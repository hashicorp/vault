package models
type PrintMultipageLayout int

const (
    CLOCKWISEFROMTOPLEFT_PRINTMULTIPAGELAYOUT PrintMultipageLayout = iota
    COUNTERCLOCKWISEFROMTOPLEFT_PRINTMULTIPAGELAYOUT
    COUNTERCLOCKWISEFROMTOPRIGHT_PRINTMULTIPAGELAYOUT
    CLOCKWISEFROMTOPRIGHT_PRINTMULTIPAGELAYOUT
    COUNTERCLOCKWISEFROMBOTTOMLEFT_PRINTMULTIPAGELAYOUT
    CLOCKWISEFROMBOTTOMLEFT_PRINTMULTIPAGELAYOUT
    COUNTERCLOCKWISEFROMBOTTOMRIGHT_PRINTMULTIPAGELAYOUT
    CLOCKWISEFROMBOTTOMRIGHT_PRINTMULTIPAGELAYOUT
    UNKNOWNFUTUREVALUE_PRINTMULTIPAGELAYOUT
)

func (i PrintMultipageLayout) String() string {
    return []string{"clockwiseFromTopLeft", "counterclockwiseFromTopLeft", "counterclockwiseFromTopRight", "clockwiseFromTopRight", "counterclockwiseFromBottomLeft", "clockwiseFromBottomLeft", "counterclockwiseFromBottomRight", "clockwiseFromBottomRight", "unknownFutureValue"}[i]
}
func ParsePrintMultipageLayout(v string) (any, error) {
    result := CLOCKWISEFROMTOPLEFT_PRINTMULTIPAGELAYOUT
    switch v {
        case "clockwiseFromTopLeft":
            result = CLOCKWISEFROMTOPLEFT_PRINTMULTIPAGELAYOUT
        case "counterclockwiseFromTopLeft":
            result = COUNTERCLOCKWISEFROMTOPLEFT_PRINTMULTIPAGELAYOUT
        case "counterclockwiseFromTopRight":
            result = COUNTERCLOCKWISEFROMTOPRIGHT_PRINTMULTIPAGELAYOUT
        case "clockwiseFromTopRight":
            result = CLOCKWISEFROMTOPRIGHT_PRINTMULTIPAGELAYOUT
        case "counterclockwiseFromBottomLeft":
            result = COUNTERCLOCKWISEFROMBOTTOMLEFT_PRINTMULTIPAGELAYOUT
        case "clockwiseFromBottomLeft":
            result = CLOCKWISEFROMBOTTOMLEFT_PRINTMULTIPAGELAYOUT
        case "counterclockwiseFromBottomRight":
            result = COUNTERCLOCKWISEFROMBOTTOMRIGHT_PRINTMULTIPAGELAYOUT
        case "clockwiseFromBottomRight":
            result = CLOCKWISEFROMBOTTOMRIGHT_PRINTMULTIPAGELAYOUT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTMULTIPAGELAYOUT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintMultipageLayout(values []PrintMultipageLayout) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintMultipageLayout) isMultiValue() bool {
    return false
}

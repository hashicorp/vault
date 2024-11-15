package models
type BinaryOperator int

const (
    OR_BINARYOPERATOR BinaryOperator = iota
    AND_BINARYOPERATOR
)

func (i BinaryOperator) String() string {
    return []string{"or", "and"}[i]
}
func ParseBinaryOperator(v string) (any, error) {
    result := OR_BINARYOPERATOR
    switch v {
        case "or":
            result = OR_BINARYOPERATOR
        case "and":
            result = AND_BINARYOPERATOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBinaryOperator(values []BinaryOperator) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BinaryOperator) isMultiValue() bool {
    return false
}

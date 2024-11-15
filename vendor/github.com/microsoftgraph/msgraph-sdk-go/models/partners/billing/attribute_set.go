package billing
type AttributeSet int

const (
    FULL_ATTRIBUTESET AttributeSet = iota
    BASIC_ATTRIBUTESET
    UNKNOWNFUTUREVALUE_ATTRIBUTESET
)

func (i AttributeSet) String() string {
    return []string{"full", "basic", "unknownFutureValue"}[i]
}
func ParseAttributeSet(v string) (any, error) {
    result := FULL_ATTRIBUTESET
    switch v {
        case "full":
            result = FULL_ATTRIBUTESET
        case "basic":
            result = BASIC_ATTRIBUTESET
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ATTRIBUTESET
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAttributeSet(values []AttributeSet) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AttributeSet) isMultiValue() bool {
    return false
}

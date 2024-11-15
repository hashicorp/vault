package termstore
type RelationType int

const (
    PIN_RELATIONTYPE RelationType = iota
    REUSE_RELATIONTYPE
    UNKNOWNFUTUREVALUE_RELATIONTYPE
)

func (i RelationType) String() string {
    return []string{"pin", "reuse", "unknownFutureValue"}[i]
}
func ParseRelationType(v string) (any, error) {
    result := PIN_RELATIONTYPE
    switch v {
        case "pin":
            result = PIN_RELATIONTYPE
        case "reuse":
            result = REUSE_RELATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RELATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRelationType(values []RelationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RelationType) isMultiValue() bool {
    return false
}

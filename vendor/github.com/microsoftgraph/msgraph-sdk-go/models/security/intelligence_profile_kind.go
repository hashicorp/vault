package security
type IntelligenceProfileKind int

const (
    ACTOR_INTELLIGENCEPROFILEKIND IntelligenceProfileKind = iota
    TOOL_INTELLIGENCEPROFILEKIND
    UNKNOWNFUTUREVALUE_INTELLIGENCEPROFILEKIND
)

func (i IntelligenceProfileKind) String() string {
    return []string{"actor", "tool", "unknownFutureValue"}[i]
}
func ParseIntelligenceProfileKind(v string) (any, error) {
    result := ACTOR_INTELLIGENCEPROFILEKIND
    switch v {
        case "actor":
            result = ACTOR_INTELLIGENCEPROFILEKIND
        case "tool":
            result = TOOL_INTELLIGENCEPROFILEKIND
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INTELLIGENCEPROFILEKIND
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIntelligenceProfileKind(values []IntelligenceProfileKind) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IntelligenceProfileKind) isMultiValue() bool {
    return false
}

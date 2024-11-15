package externalconnectors
type RuleOperation int

const (
    EQUALS_RULEOPERATION RuleOperation = iota
    NOTEQUALS_RULEOPERATION
    CONTAINS_RULEOPERATION
    NOTCONTAINS_RULEOPERATION
    LESSTHAN_RULEOPERATION
    GREATERTHAN_RULEOPERATION
    STARTSWITH_RULEOPERATION
    UNKNOWNFUTUREVALUE_RULEOPERATION
)

func (i RuleOperation) String() string {
    return []string{"equals", "notEquals", "contains", "notContains", "lessThan", "greaterThan", "startsWith", "unknownFutureValue"}[i]
}
func ParseRuleOperation(v string) (any, error) {
    result := EQUALS_RULEOPERATION
    switch v {
        case "equals":
            result = EQUALS_RULEOPERATION
        case "notEquals":
            result = NOTEQUALS_RULEOPERATION
        case "contains":
            result = CONTAINS_RULEOPERATION
        case "notContains":
            result = NOTCONTAINS_RULEOPERATION
        case "lessThan":
            result = LESSTHAN_RULEOPERATION
        case "greaterThan":
            result = GREATERTHAN_RULEOPERATION
        case "startsWith":
            result = STARTSWITH_RULEOPERATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RULEOPERATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRuleOperation(values []RuleOperation) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RuleOperation) isMultiValue() bool {
    return false
}

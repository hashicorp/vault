package termstore
type TermGroupScope int

const (
    GLOBAL_TERMGROUPSCOPE TermGroupScope = iota
    SYSTEM_TERMGROUPSCOPE
    SITECOLLECTION_TERMGROUPSCOPE
    UNKNOWNFUTUREVALUE_TERMGROUPSCOPE
)

func (i TermGroupScope) String() string {
    return []string{"global", "system", "siteCollection", "unknownFutureValue"}[i]
}
func ParseTermGroupScope(v string) (any, error) {
    result := GLOBAL_TERMGROUPSCOPE
    switch v {
        case "global":
            result = GLOBAL_TERMGROUPSCOPE
        case "system":
            result = SYSTEM_TERMGROUPSCOPE
        case "siteCollection":
            result = SITECOLLECTION_TERMGROUPSCOPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TERMGROUPSCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTermGroupScope(values []TermGroupScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TermGroupScope) isMultiValue() bool {
    return false
}

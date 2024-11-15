package security
type QueryType int

const (
    FILES_QUERYTYPE QueryType = iota
    MESSAGES_QUERYTYPE
    UNKNOWNFUTUREVALUE_QUERYTYPE
)

func (i QueryType) String() string {
    return []string{"files", "messages", "unknownFutureValue"}[i]
}
func ParseQueryType(v string) (any, error) {
    result := FILES_QUERYTYPE
    switch v {
        case "files":
            result = FILES_QUERYTYPE
        case "messages":
            result = MESSAGES_QUERYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_QUERYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeQueryType(values []QueryType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i QueryType) isMultiValue() bool {
    return false
}

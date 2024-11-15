package search
type AnswerState int

const (
    PUBLISHED_ANSWERSTATE AnswerState = iota
    DRAFT_ANSWERSTATE
    EXCLUDED_ANSWERSTATE
    UNKNOWNFUTUREVALUE_ANSWERSTATE
)

func (i AnswerState) String() string {
    return []string{"published", "draft", "excluded", "unknownFutureValue"}[i]
}
func ParseAnswerState(v string) (any, error) {
    result := PUBLISHED_ANSWERSTATE
    switch v {
        case "published":
            result = PUBLISHED_ANSWERSTATE
        case "draft":
            result = DRAFT_ANSWERSTATE
        case "excluded":
            result = EXCLUDED_ANSWERSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ANSWERSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAnswerState(values []AnswerState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AnswerState) isMultiValue() bool {
    return false
}

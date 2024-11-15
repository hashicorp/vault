package models
type VirtualEventRegistrationQuestionAnswerInputType int

const (
    TEXT_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE VirtualEventRegistrationQuestionAnswerInputType = iota
    MULTILINETEXT_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
    SINGLECHOICE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
    MULTICHOICE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
    BOOLEAN_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
    UNKNOWNFUTUREVALUE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
)

func (i VirtualEventRegistrationQuestionAnswerInputType) String() string {
    return []string{"text", "multilineText", "singleChoice", "multiChoice", "boolean", "unknownFutureValue"}[i]
}
func ParseVirtualEventRegistrationQuestionAnswerInputType(v string) (any, error) {
    result := TEXT_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
    switch v {
        case "text":
            result = TEXT_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        case "multilineText":
            result = MULTILINETEXT_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        case "singleChoice":
            result = SINGLECHOICE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        case "multiChoice":
            result = MULTICHOICE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        case "boolean":
            result = BOOLEAN_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIRTUALEVENTREGISTRATIONQUESTIONANSWERINPUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVirtualEventRegistrationQuestionAnswerInputType(values []VirtualEventRegistrationQuestionAnswerInputType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VirtualEventRegistrationQuestionAnswerInputType) isMultiValue() bool {
    return false
}

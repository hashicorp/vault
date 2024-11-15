package models
type ImageTaggingChoice int

const (
    DISABLED_IMAGETAGGINGCHOICE ImageTaggingChoice = iota
    BASIC_IMAGETAGGINGCHOICE
    ENHANCED_IMAGETAGGINGCHOICE
    UNKNOWNFUTUREVALUE_IMAGETAGGINGCHOICE
)

func (i ImageTaggingChoice) String() string {
    return []string{"disabled", "basic", "enhanced", "unknownFutureValue"}[i]
}
func ParseImageTaggingChoice(v string) (any, error) {
    result := DISABLED_IMAGETAGGINGCHOICE
    switch v {
        case "disabled":
            result = DISABLED_IMAGETAGGINGCHOICE
        case "basic":
            result = BASIC_IMAGETAGGINGCHOICE
        case "enhanced":
            result = ENHANCED_IMAGETAGGINGCHOICE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_IMAGETAGGINGCHOICE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeImageTaggingChoice(values []ImageTaggingChoice) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ImageTaggingChoice) isMultiValue() bool {
    return false
}

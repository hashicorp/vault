package models
type SectionEmphasisType int

const (
    NONE_SECTIONEMPHASISTYPE SectionEmphasisType = iota
    NEUTRAL_SECTIONEMPHASISTYPE
    SOFT_SECTIONEMPHASISTYPE
    STRONG_SECTIONEMPHASISTYPE
    UNKNOWNFUTUREVALUE_SECTIONEMPHASISTYPE
)

func (i SectionEmphasisType) String() string {
    return []string{"none", "neutral", "soft", "strong", "unknownFutureValue"}[i]
}
func ParseSectionEmphasisType(v string) (any, error) {
    result := NONE_SECTIONEMPHASISTYPE
    switch v {
        case "none":
            result = NONE_SECTIONEMPHASISTYPE
        case "neutral":
            result = NEUTRAL_SECTIONEMPHASISTYPE
        case "soft":
            result = SOFT_SECTIONEMPHASISTYPE
        case "strong":
            result = STRONG_SECTIONEMPHASISTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SECTIONEMPHASISTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSectionEmphasisType(values []SectionEmphasisType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SectionEmphasisType) isMultiValue() bool {
    return false
}

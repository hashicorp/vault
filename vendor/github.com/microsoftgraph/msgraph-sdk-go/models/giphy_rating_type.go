package models
type GiphyRatingType int

const (
    STRICT_GIPHYRATINGTYPE GiphyRatingType = iota
    MODERATE_GIPHYRATINGTYPE
    UNKNOWNFUTUREVALUE_GIPHYRATINGTYPE
)

func (i GiphyRatingType) String() string {
    return []string{"strict", "moderate", "unknownFutureValue"}[i]
}
func ParseGiphyRatingType(v string) (any, error) {
    result := STRICT_GIPHYRATINGTYPE
    switch v {
        case "strict":
            result = STRICT_GIPHYRATINGTYPE
        case "moderate":
            result = MODERATE_GIPHYRATINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_GIPHYRATINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeGiphyRatingType(values []GiphyRatingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i GiphyRatingType) isMultiValue() bool {
    return false
}

package models
// Movies rating labels in Japan
type RatingJapanMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGJAPANMOVIESTYPE RatingJapanMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGJAPANMOVIESTYPE
    // Suitable for all ages
    GENERAL_RATINGJAPANMOVIESTYPE
    // The PG-12 classification requests parental guidance for young people under 12
    PARENTALGUIDANCE_RATINGJAPANMOVIESTYPE
    // The R15+ classification is suitable for viewers of 15 or older
    AGESABOVE15_RATINGJAPANMOVIESTYPE
    // The R18+ classification is suitable for viewers of 18 or older
    AGESABOVE18_RATINGJAPANMOVIESTYPE
)

func (i RatingJapanMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove15", "agesAbove18"}[i]
}
func ParseRatingJapanMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGJAPANMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGJAPANMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGJAPANMOVIESTYPE
        case "general":
            result = GENERAL_RATINGJAPANMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGJAPANMOVIESTYPE
        case "agesAbove15":
            result = AGESABOVE15_RATINGJAPANMOVIESTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGJAPANMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingJapanMoviesType(values []RatingJapanMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingJapanMoviesType) isMultiValue() bool {
    return false
}

package models
// Movies rating labels in Ireland
type RatingIrelandMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGIRELANDMOVIESTYPE RatingIrelandMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGIRELANDMOVIESTYPE
    // Suitable for children of school going age
    GENERAL_RATINGIRELANDMOVIESTYPE
    // The PG classification advises parental guidance
    PARENTALGUIDANCE_RATINGIRELANDMOVIESTYPE
    // The 12A classification is suitable for viewers of 12 or older
    AGESABOVE12_RATINGIRELANDMOVIESTYPE
    // The 15A classification is suitable for viewers of 15 or older
    AGESABOVE15_RATINGIRELANDMOVIESTYPE
    // The 16 classification is suitable for viewers of 16 or older
    AGESABOVE16_RATINGIRELANDMOVIESTYPE
    // The 18 classification, suitable only for adults
    ADULTS_RATINGIRELANDMOVIESTYPE
)

func (i RatingIrelandMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove12", "agesAbove15", "agesAbove16", "adults"}[i]
}
func ParseRatingIrelandMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGIRELANDMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGIRELANDMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGIRELANDMOVIESTYPE
        case "general":
            result = GENERAL_RATINGIRELANDMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGIRELANDMOVIESTYPE
        case "agesAbove12":
            result = AGESABOVE12_RATINGIRELANDMOVIESTYPE
        case "agesAbove15":
            result = AGESABOVE15_RATINGIRELANDMOVIESTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGIRELANDMOVIESTYPE
        case "adults":
            result = ADULTS_RATINGIRELANDMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingIrelandMoviesType(values []RatingIrelandMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingIrelandMoviesType) isMultiValue() bool {
    return false
}

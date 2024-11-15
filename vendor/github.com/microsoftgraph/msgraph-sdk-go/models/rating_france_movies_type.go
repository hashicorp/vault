package models
// Movies rating labels in France
type RatingFranceMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGFRANCEMOVIESTYPE RatingFranceMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGFRANCEMOVIESTYPE
    // The 10 classification prohibits the screening of the film to minors under 10
    AGESABOVE10_RATINGFRANCEMOVIESTYPE
    // The 12 classification prohibits the screening of the film to minors under 12
    AGESABOVE12_RATINGFRANCEMOVIESTYPE
    // The 16 classification prohibits the screening of the film to minors under 16
    AGESABOVE16_RATINGFRANCEMOVIESTYPE
    // The 18 classification prohibits the screening to minors under 18
    AGESABOVE18_RATINGFRANCEMOVIESTYPE
)

func (i RatingFranceMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"}[i]
}
func ParseRatingFranceMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGFRANCEMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGFRANCEMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGFRANCEMOVIESTYPE
        case "agesAbove10":
            result = AGESABOVE10_RATINGFRANCEMOVIESTYPE
        case "agesAbove12":
            result = AGESABOVE12_RATINGFRANCEMOVIESTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGFRANCEMOVIESTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGFRANCEMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingFranceMoviesType(values []RatingFranceMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingFranceMoviesType) isMultiValue() bool {
    return false
}

package models
// Movies rating labels in Canada
type RatingCanadaMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGCANADAMOVIESTYPE RatingCanadaMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGCANADAMOVIESTYPE
    // The G classification is suitable for all ages
    GENERAL_RATINGCANADAMOVIESTYPE
    // The PG classification advises parental guidance
    PARENTALGUIDANCE_RATINGCANADAMOVIESTYPE
    // The 14A classification is suitable for viewers above 14 or older
    AGESABOVE14_RATINGCANADAMOVIESTYPE
    // The 18A classification is suitable for viewers above 18 or older
    AGESABOVE18_RATINGCANADAMOVIESTYPE
    // The R classification is restricted to 18 years and older
    RESTRICTED_RATINGCANADAMOVIESTYPE
)

func (i RatingCanadaMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove14", "agesAbove18", "restricted"}[i]
}
func ParseRatingCanadaMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGCANADAMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGCANADAMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGCANADAMOVIESTYPE
        case "general":
            result = GENERAL_RATINGCANADAMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGCANADAMOVIESTYPE
        case "agesAbove14":
            result = AGESABOVE14_RATINGCANADAMOVIESTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGCANADAMOVIESTYPE
        case "restricted":
            result = RESTRICTED_RATINGCANADAMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingCanadaMoviesType(values []RatingCanadaMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingCanadaMoviesType) isMultiValue() bool {
    return false
}

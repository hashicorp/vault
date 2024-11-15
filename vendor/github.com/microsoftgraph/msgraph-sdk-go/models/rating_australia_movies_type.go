package models
// Movies rating labels in Australia
type RatingAustraliaMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGAUSTRALIAMOVIESTYPE RatingAustraliaMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGAUSTRALIAMOVIESTYPE
    // The G classification is suitable for everyone
    GENERAL_RATINGAUSTRALIAMOVIESTYPE
    // The PG recommends viewers under 15 with guidance from parents or guardians
    PARENTALGUIDANCE_RATINGAUSTRALIAMOVIESTYPE
    // The M classification is not recommended for viewers under 15
    MATURE_RATINGAUSTRALIAMOVIESTYPE
    // The MA15+ classification is not suitable for viewers under 15
    AGESABOVE15_RATINGAUSTRALIAMOVIESTYPE
    // The R18+ classification is not suitable for viewers under 18
    AGESABOVE18_RATINGAUSTRALIAMOVIESTYPE
)

func (i RatingAustraliaMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove15", "agesAbove18"}[i]
}
func ParseRatingAustraliaMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGAUSTRALIAMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGAUSTRALIAMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGAUSTRALIAMOVIESTYPE
        case "general":
            result = GENERAL_RATINGAUSTRALIAMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGAUSTRALIAMOVIESTYPE
        case "mature":
            result = MATURE_RATINGAUSTRALIAMOVIESTYPE
        case "agesAbove15":
            result = AGESABOVE15_RATINGAUSTRALIAMOVIESTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGAUSTRALIAMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingAustraliaMoviesType(values []RatingAustraliaMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingAustraliaMoviesType) isMultiValue() bool {
    return false
}

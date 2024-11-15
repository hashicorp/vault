package models
// Movies rating labels in United States
type RatingUnitedStatesMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGUNITEDSTATESMOVIESTYPE RatingUnitedStatesMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGUNITEDSTATESMOVIESTYPE
    // G, all ages admitted
    GENERAL_RATINGUNITEDSTATESMOVIESTYPE
    // PG, some material may not be suitable for children
    PARENTALGUIDANCE_RATINGUNITEDSTATESMOVIESTYPE
    // PG13, some material may be inappropriate for children under 13
    PARENTALGUIDANCE13_RATINGUNITEDSTATESMOVIESTYPE
    // R, viewers under 17 require accompanying parent or adult guardian
    RESTRICTED_RATINGUNITEDSTATESMOVIESTYPE
    // NC17, adults only
    ADULTS_RATINGUNITEDSTATESMOVIESTYPE
)

func (i RatingUnitedStatesMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "parentalGuidance13", "restricted", "adults"}[i]
}
func ParseRatingUnitedStatesMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGUNITEDSTATESMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGUNITEDSTATESMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGUNITEDSTATESMOVIESTYPE
        case "general":
            result = GENERAL_RATINGUNITEDSTATESMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGUNITEDSTATESMOVIESTYPE
        case "parentalGuidance13":
            result = PARENTALGUIDANCE13_RATINGUNITEDSTATESMOVIESTYPE
        case "restricted":
            result = RESTRICTED_RATINGUNITEDSTATESMOVIESTYPE
        case "adults":
            result = ADULTS_RATINGUNITEDSTATESMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingUnitedStatesMoviesType(values []RatingUnitedStatesMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingUnitedStatesMoviesType) isMultiValue() bool {
    return false
}

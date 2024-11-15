package models
// Movies rating labels in New Zealand
type RatingNewZealandMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGNEWZEALANDMOVIESTYPE RatingNewZealandMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGNEWZEALANDMOVIESTYPE
    // Suitable for general audience
    GENERAL_RATINGNEWZEALANDMOVIESTYPE
    // The PG classification recommends parental guidance
    PARENTALGUIDANCE_RATINGNEWZEALANDMOVIESTYPE
    // The M classification is suitable for mature audience
    MATURE_RATINGNEWZEALANDMOVIESTYPE
    // The R13 classification is restricted to persons 13 years and over
    AGESABOVE13_RATINGNEWZEALANDMOVIESTYPE
    // The R15 classification is restricted to persons 15 years and over
    AGESABOVE15_RATINGNEWZEALANDMOVIESTYPE
    // The R16 classification is restricted to persons 16 years and over
    AGESABOVE16_RATINGNEWZEALANDMOVIESTYPE
    // The R18 classification is restricted to persons 18 years and over
    AGESABOVE18_RATINGNEWZEALANDMOVIESTYPE
    // The R classification is restricted to a certain audience
    RESTRICTED_RATINGNEWZEALANDMOVIESTYPE
    // The RP16 classification requires viewers under 16 accompanied by a parent or an adult
    AGESABOVE16RESTRICTED_RATINGNEWZEALANDMOVIESTYPE
)

func (i RatingNewZealandMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove13", "agesAbove15", "agesAbove16", "agesAbove18", "restricted", "agesAbove16Restricted"}[i]
}
func ParseRatingNewZealandMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGNEWZEALANDMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGNEWZEALANDMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGNEWZEALANDMOVIESTYPE
        case "general":
            result = GENERAL_RATINGNEWZEALANDMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGNEWZEALANDMOVIESTYPE
        case "mature":
            result = MATURE_RATINGNEWZEALANDMOVIESTYPE
        case "agesAbove13":
            result = AGESABOVE13_RATINGNEWZEALANDMOVIESTYPE
        case "agesAbove15":
            result = AGESABOVE15_RATINGNEWZEALANDMOVIESTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGNEWZEALANDMOVIESTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGNEWZEALANDMOVIESTYPE
        case "restricted":
            result = RESTRICTED_RATINGNEWZEALANDMOVIESTYPE
        case "agesAbove16Restricted":
            result = AGESABOVE16RESTRICTED_RATINGNEWZEALANDMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingNewZealandMoviesType(values []RatingNewZealandMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingNewZealandMoviesType) isMultiValue() bool {
    return false
}

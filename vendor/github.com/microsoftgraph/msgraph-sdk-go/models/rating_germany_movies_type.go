package models
// Movies rating labels in Germany
type RatingGermanyMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGGERMANYMOVIESTYPE RatingGermanyMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGGERMANYMOVIESTYPE
    // Ab 0 Jahren, no age restrictions
    GENERAL_RATINGGERMANYMOVIESTYPE
    // Ab 6 Jahren, ages 6 and older
    AGESABOVE6_RATINGGERMANYMOVIESTYPE
    // Ab 12 Jahren, ages 12 and older
    AGESABOVE12_RATINGGERMANYMOVIESTYPE
    // Ab 16 Jahren, ages 16 and older
    AGESABOVE16_RATINGGERMANYMOVIESTYPE
    // Ab 18 Jahren, adults only
    ADULTS_RATINGGERMANYMOVIESTYPE
)

func (i RatingGermanyMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"}[i]
}
func ParseRatingGermanyMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGGERMANYMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGGERMANYMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGGERMANYMOVIESTYPE
        case "general":
            result = GENERAL_RATINGGERMANYMOVIESTYPE
        case "agesAbove6":
            result = AGESABOVE6_RATINGGERMANYMOVIESTYPE
        case "agesAbove12":
            result = AGESABOVE12_RATINGGERMANYMOVIESTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGGERMANYMOVIESTYPE
        case "adults":
            result = ADULTS_RATINGGERMANYMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingGermanyMoviesType(values []RatingGermanyMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingGermanyMoviesType) isMultiValue() bool {
    return false
}

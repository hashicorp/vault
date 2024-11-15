package models
// Movies rating labels in United Kingdom
type RatingUnitedKingdomMoviesType int

const (
    // Default value, allow all movies content
    ALLALLOWED_RATINGUNITEDKINGDOMMOVIESTYPE RatingUnitedKingdomMoviesType = iota
    // Do not allow any movies content
    ALLBLOCKED_RATINGUNITEDKINGDOMMOVIESTYPE
    // The U classification is suitable for all ages
    GENERAL_RATINGUNITEDKINGDOMMOVIESTYPE
    // The UC classification is suitable for pre-school children, an old rating label
    UNIVERSALCHILDREN_RATINGUNITEDKINGDOMMOVIESTYPE
    // The PG classification is suitable for mature
    PARENTALGUIDANCE_RATINGUNITEDKINGDOMMOVIESTYPE
    // 12, video release suitable for 12 years and over
    AGESABOVE12VIDEO_RATINGUNITEDKINGDOMMOVIESTYPE
    // 12A, cinema release suitable for 12 years and over
    AGESABOVE12CINEMA_RATINGUNITEDKINGDOMMOVIESTYPE
    // 15, suitable only for 15 years and older
    AGESABOVE15_RATINGUNITEDKINGDOMMOVIESTYPE
    // Suitable only for adults
    ADULTS_RATINGUNITEDKINGDOMMOVIESTYPE
)

func (i RatingUnitedKingdomMoviesType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "universalChildren", "parentalGuidance", "agesAbove12Video", "agesAbove12Cinema", "agesAbove15", "adults"}[i]
}
func ParseRatingUnitedKingdomMoviesType(v string) (any, error) {
    result := ALLALLOWED_RATINGUNITEDKINGDOMMOVIESTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGUNITEDKINGDOMMOVIESTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGUNITEDKINGDOMMOVIESTYPE
        case "general":
            result = GENERAL_RATINGUNITEDKINGDOMMOVIESTYPE
        case "universalChildren":
            result = UNIVERSALCHILDREN_RATINGUNITEDKINGDOMMOVIESTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGUNITEDKINGDOMMOVIESTYPE
        case "agesAbove12Video":
            result = AGESABOVE12VIDEO_RATINGUNITEDKINGDOMMOVIESTYPE
        case "agesAbove12Cinema":
            result = AGESABOVE12CINEMA_RATINGUNITEDKINGDOMMOVIESTYPE
        case "agesAbove15":
            result = AGESABOVE15_RATINGUNITEDKINGDOMMOVIESTYPE
        case "adults":
            result = ADULTS_RATINGUNITEDKINGDOMMOVIESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingUnitedKingdomMoviesType(values []RatingUnitedKingdomMoviesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingUnitedKingdomMoviesType) isMultiValue() bool {
    return false
}

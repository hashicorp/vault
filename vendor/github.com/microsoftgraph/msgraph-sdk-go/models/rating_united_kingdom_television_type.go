package models
// TV content rating labels in United Kingdom
type RatingUnitedKingdomTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGUNITEDKINGDOMTELEVISIONTYPE RatingUnitedKingdomTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGUNITEDKINGDOMTELEVISIONTYPE
    // Allowing TV contents with a warning message
    CAUTION_RATINGUNITEDKINGDOMTELEVISIONTYPE
)

func (i RatingUnitedKingdomTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "caution"}[i]
}
func ParseRatingUnitedKingdomTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGUNITEDKINGDOMTELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGUNITEDKINGDOMTELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGUNITEDKINGDOMTELEVISIONTYPE
        case "caution":
            result = CAUTION_RATINGUNITEDKINGDOMTELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingUnitedKingdomTelevisionType(values []RatingUnitedKingdomTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingUnitedKingdomTelevisionType) isMultiValue() bool {
    return false
}

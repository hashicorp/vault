package models
// TV content rating labels in France
type RatingFranceTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGFRANCETELEVISIONTYPE RatingFranceTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGFRANCETELEVISIONTYPE
    // The -10 classification is not recommended for children under 10
    AGESABOVE10_RATINGFRANCETELEVISIONTYPE
    // The -12 classification is not recommended for children under 12
    AGESABOVE12_RATINGFRANCETELEVISIONTYPE
    // The -16 classification is not recommended for children under 16
    AGESABOVE16_RATINGFRANCETELEVISIONTYPE
    // The -18 classification is not recommended for persons under 18
    AGESABOVE18_RATINGFRANCETELEVISIONTYPE
)

func (i RatingFranceTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"}[i]
}
func ParseRatingFranceTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGFRANCETELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGFRANCETELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGFRANCETELEVISIONTYPE
        case "agesAbove10":
            result = AGESABOVE10_RATINGFRANCETELEVISIONTYPE
        case "agesAbove12":
            result = AGESABOVE12_RATINGFRANCETELEVISIONTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGFRANCETELEVISIONTYPE
        case "agesAbove18":
            result = AGESABOVE18_RATINGFRANCETELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingFranceTelevisionType(values []RatingFranceTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingFranceTelevisionType) isMultiValue() bool {
    return false
}

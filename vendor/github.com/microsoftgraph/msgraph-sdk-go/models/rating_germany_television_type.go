package models
// TV content rating labels in Germany
type RatingGermanyTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGGERMANYTELEVISIONTYPE RatingGermanyTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGGERMANYTELEVISIONTYPE
    // Ab 0 Jahren, no age restrictions
    GENERAL_RATINGGERMANYTELEVISIONTYPE
    // Ab 6 Jahren, ages 6 and older
    AGESABOVE6_RATINGGERMANYTELEVISIONTYPE
    // Ab 12 Jahren, ages 12 and older
    AGESABOVE12_RATINGGERMANYTELEVISIONTYPE
    // Ab 16 Jahren, ages 16 and older
    AGESABOVE16_RATINGGERMANYTELEVISIONTYPE
    // Ab 18 Jahren, adults only
    ADULTS_RATINGGERMANYTELEVISIONTYPE
)

func (i RatingGermanyTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"}[i]
}
func ParseRatingGermanyTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGGERMANYTELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGGERMANYTELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGGERMANYTELEVISIONTYPE
        case "general":
            result = GENERAL_RATINGGERMANYTELEVISIONTYPE
        case "agesAbove6":
            result = AGESABOVE6_RATINGGERMANYTELEVISIONTYPE
        case "agesAbove12":
            result = AGESABOVE12_RATINGGERMANYTELEVISIONTYPE
        case "agesAbove16":
            result = AGESABOVE16_RATINGGERMANYTELEVISIONTYPE
        case "adults":
            result = ADULTS_RATINGGERMANYTELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingGermanyTelevisionType(values []RatingGermanyTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingGermanyTelevisionType) isMultiValue() bool {
    return false
}

package models
// TV content rating labels in New Zealand
type RatingNewZealandTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGNEWZEALANDTELEVISIONTYPE RatingNewZealandTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGNEWZEALANDTELEVISIONTYPE
    // The G classification excludes materials likely to harm children under 14
    GENERAL_RATINGNEWZEALANDTELEVISIONTYPE
    // The PGR classification encourages parents and guardians to supervise younger viewers
    PARENTALGUIDANCE_RATINGNEWZEALANDTELEVISIONTYPE
    // The AO classification is not suitable for children
    ADULTS_RATINGNEWZEALANDTELEVISIONTYPE
)

func (i RatingNewZealandTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "parentalGuidance", "adults"}[i]
}
func ParseRatingNewZealandTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGNEWZEALANDTELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGNEWZEALANDTELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGNEWZEALANDTELEVISIONTYPE
        case "general":
            result = GENERAL_RATINGNEWZEALANDTELEVISIONTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGNEWZEALANDTELEVISIONTYPE
        case "adults":
            result = ADULTS_RATINGNEWZEALANDTELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingNewZealandTelevisionType(values []RatingNewZealandTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingNewZealandTelevisionType) isMultiValue() bool {
    return false
}

package models
// TV content rating labels in Ireland
type RatingIrelandTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGIRELANDTELEVISIONTYPE RatingIrelandTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGIRELANDTELEVISIONTYPE
    // The GA classification is suitable for all audiences
    GENERAL_RATINGIRELANDTELEVISIONTYPE
    // The CH classification is suitable for children
    CHILDREN_RATINGIRELANDTELEVISIONTYPE
    // The YA classification is suitable for teenage audience
    YOUNGADULTS_RATINGIRELANDTELEVISIONTYPE
    // The PS classification invites parents and guardians to consider restriction children’s access
    PARENTALSUPERVISION_RATINGIRELANDTELEVISIONTYPE
    // The MA classification is suitable for adults
    MATURE_RATINGIRELANDTELEVISIONTYPE
)

func (i RatingIrelandTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "general", "children", "youngAdults", "parentalSupervision", "mature"}[i]
}
func ParseRatingIrelandTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGIRELANDTELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGIRELANDTELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGIRELANDTELEVISIONTYPE
        case "general":
            result = GENERAL_RATINGIRELANDTELEVISIONTYPE
        case "children":
            result = CHILDREN_RATINGIRELANDTELEVISIONTYPE
        case "youngAdults":
            result = YOUNGADULTS_RATINGIRELANDTELEVISIONTYPE
        case "parentalSupervision":
            result = PARENTALSUPERVISION_RATINGIRELANDTELEVISIONTYPE
        case "mature":
            result = MATURE_RATINGIRELANDTELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingIrelandTelevisionType(values []RatingIrelandTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingIrelandTelevisionType) isMultiValue() bool {
    return false
}

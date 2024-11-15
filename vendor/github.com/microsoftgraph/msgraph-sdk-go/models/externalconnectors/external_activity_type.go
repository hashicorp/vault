package externalconnectors
type ExternalActivityType int

const (
    VIEWED_EXTERNALACTIVITYTYPE ExternalActivityType = iota
    MODIFIED_EXTERNALACTIVITYTYPE
    CREATED_EXTERNALACTIVITYTYPE
    COMMENTED_EXTERNALACTIVITYTYPE
    UNKNOWNFUTUREVALUE_EXTERNALACTIVITYTYPE
)

func (i ExternalActivityType) String() string {
    return []string{"viewed", "modified", "created", "commented", "unknownFutureValue"}[i]
}
func ParseExternalActivityType(v string) (any, error) {
    result := VIEWED_EXTERNALACTIVITYTYPE
    switch v {
        case "viewed":
            result = VIEWED_EXTERNALACTIVITYTYPE
        case "modified":
            result = MODIFIED_EXTERNALACTIVITYTYPE
        case "created":
            result = CREATED_EXTERNALACTIVITYTYPE
        case "commented":
            result = COMMENTED_EXTERNALACTIVITYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EXTERNALACTIVITYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExternalActivityType(values []ExternalActivityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExternalActivityType) isMultiValue() bool {
    return false
}

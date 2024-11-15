package models
type ContactRelationship int

const (
    PARENT_CONTACTRELATIONSHIP ContactRelationship = iota
    RELATIVE_CONTACTRELATIONSHIP
    AIDE_CONTACTRELATIONSHIP
    DOCTOR_CONTACTRELATIONSHIP
    GUARDIAN_CONTACTRELATIONSHIP
    CHILD_CONTACTRELATIONSHIP
    OTHER_CONTACTRELATIONSHIP
    UNKNOWNFUTUREVALUE_CONTACTRELATIONSHIP
)

func (i ContactRelationship) String() string {
    return []string{"parent", "relative", "aide", "doctor", "guardian", "child", "other", "unknownFutureValue"}[i]
}
func ParseContactRelationship(v string) (any, error) {
    result := PARENT_CONTACTRELATIONSHIP
    switch v {
        case "parent":
            result = PARENT_CONTACTRELATIONSHIP
        case "relative":
            result = RELATIVE_CONTACTRELATIONSHIP
        case "aide":
            result = AIDE_CONTACTRELATIONSHIP
        case "doctor":
            result = DOCTOR_CONTACTRELATIONSHIP
        case "guardian":
            result = GUARDIAN_CONTACTRELATIONSHIP
        case "child":
            result = CHILD_CONTACTRELATIONSHIP
        case "other":
            result = OTHER_CONTACTRELATIONSHIP
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONTACTRELATIONSHIP
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeContactRelationship(values []ContactRelationship) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ContactRelationship) isMultiValue() bool {
    return false
}

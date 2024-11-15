package models
type SensitivityLabelAssignmentMethod int

const (
    STANDARD_SENSITIVITYLABELASSIGNMENTMETHOD SensitivityLabelAssignmentMethod = iota
    PRIVILEGED_SENSITIVITYLABELASSIGNMENTMETHOD
    AUTO_SENSITIVITYLABELASSIGNMENTMETHOD
    UNKNOWNFUTUREVALUE_SENSITIVITYLABELASSIGNMENTMETHOD
)

func (i SensitivityLabelAssignmentMethod) String() string {
    return []string{"standard", "privileged", "auto", "unknownFutureValue"}[i]
}
func ParseSensitivityLabelAssignmentMethod(v string) (any, error) {
    result := STANDARD_SENSITIVITYLABELASSIGNMENTMETHOD
    switch v {
        case "standard":
            result = STANDARD_SENSITIVITYLABELASSIGNMENTMETHOD
        case "privileged":
            result = PRIVILEGED_SENSITIVITYLABELASSIGNMENTMETHOD
        case "auto":
            result = AUTO_SENSITIVITYLABELASSIGNMENTMETHOD
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SENSITIVITYLABELASSIGNMENTMETHOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSensitivityLabelAssignmentMethod(values []SensitivityLabelAssignmentMethod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SensitivityLabelAssignmentMethod) isMultiValue() bool {
    return false
}

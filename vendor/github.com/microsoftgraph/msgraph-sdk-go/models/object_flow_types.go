package models
import (
    "math"
    "strings"
)
type ObjectFlowTypes int

const (
    NONE_OBJECTFLOWTYPES = 1
    ADD_OBJECTFLOWTYPES = 2
    UPDATE_OBJECTFLOWTYPES = 4
    DELETE_OBJECTFLOWTYPES = 8
)

func (i ObjectFlowTypes) String() string {
    var values []string
    options := []string{"None", "Add", "Update", "Delete"}
    for p := 0; p < 4; p++ {
        mantis := ObjectFlowTypes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseObjectFlowTypes(v string) (any, error) {
    var result ObjectFlowTypes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "None":
                result |= NONE_OBJECTFLOWTYPES
            case "Add":
                result |= ADD_OBJECTFLOWTYPES
            case "Update":
                result |= UPDATE_OBJECTFLOWTYPES
            case "Delete":
                result |= DELETE_OBJECTFLOWTYPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeObjectFlowTypes(values []ObjectFlowTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ObjectFlowTypes) isMultiValue() bool {
    return true
}

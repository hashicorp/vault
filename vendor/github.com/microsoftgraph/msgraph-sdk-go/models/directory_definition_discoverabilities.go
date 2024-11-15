package models
import (
    "math"
    "strings"
)
type DirectoryDefinitionDiscoverabilities int

const (
    NONE_DIRECTORYDEFINITIONDISCOVERABILITIES = 1
    ATTRIBUTENAMES_DIRECTORYDEFINITIONDISCOVERABILITIES = 2
    ATTRIBUTEDATATYPES_DIRECTORYDEFINITIONDISCOVERABILITIES = 4
    ATTRIBUTEREADONLY_DIRECTORYDEFINITIONDISCOVERABILITIES = 8
    REFERENCEATTRIBUTES_DIRECTORYDEFINITIONDISCOVERABILITIES = 16
    UNKNOWNFUTUREVALUE_DIRECTORYDEFINITIONDISCOVERABILITIES = 32
)

func (i DirectoryDefinitionDiscoverabilities) String() string {
    var values []string
    options := []string{"None", "AttributeNames", "AttributeDataTypes", "AttributeReadOnly", "ReferenceAttributes", "UnknownFutureValue"}
    for p := 0; p < 6; p++ {
        mantis := DirectoryDefinitionDiscoverabilities(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseDirectoryDefinitionDiscoverabilities(v string) (any, error) {
    var result DirectoryDefinitionDiscoverabilities
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "None":
                result |= NONE_DIRECTORYDEFINITIONDISCOVERABILITIES
            case "AttributeNames":
                result |= ATTRIBUTENAMES_DIRECTORYDEFINITIONDISCOVERABILITIES
            case "AttributeDataTypes":
                result |= ATTRIBUTEDATATYPES_DIRECTORYDEFINITIONDISCOVERABILITIES
            case "AttributeReadOnly":
                result |= ATTRIBUTEREADONLY_DIRECTORYDEFINITIONDISCOVERABILITIES
            case "ReferenceAttributes":
                result |= REFERENCEATTRIBUTES_DIRECTORYDEFINITIONDISCOVERABILITIES
            case "UnknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_DIRECTORYDEFINITIONDISCOVERABILITIES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeDirectoryDefinitionDiscoverabilities(values []DirectoryDefinitionDiscoverabilities) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DirectoryDefinitionDiscoverabilities) isMultiValue() bool {
    return true
}

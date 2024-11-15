package models
// Possible types of an Apple Volume Purchase Program token.
type VppTokenAccountType int

const (
    // Apple Volume Purchase Program token associated with an business program.
    BUSINESS_VPPTOKENACCOUNTTYPE VppTokenAccountType = iota
    // Apple Volume Purchase Program token associated with an education program.
    EDUCATION_VPPTOKENACCOUNTTYPE
)

func (i VppTokenAccountType) String() string {
    return []string{"business", "education"}[i]
}
func ParseVppTokenAccountType(v string) (any, error) {
    result := BUSINESS_VPPTOKENACCOUNTTYPE
    switch v {
        case "business":
            result = BUSINESS_VPPTOKENACCOUNTTYPE
        case "education":
            result = EDUCATION_VPPTOKENACCOUNTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVppTokenAccountType(values []VppTokenAccountType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VppTokenAccountType) isMultiValue() bool {
    return false
}

package models
type ThreatCategory int

const (
    UNDEFINED_THREATCATEGORY ThreatCategory = iota
    SPAM_THREATCATEGORY
    PHISHING_THREATCATEGORY
    MALWARE_THREATCATEGORY
    UNKNOWNFUTUREVALUE_THREATCATEGORY
)

func (i ThreatCategory) String() string {
    return []string{"undefined", "spam", "phishing", "malware", "unknownFutureValue"}[i]
}
func ParseThreatCategory(v string) (any, error) {
    result := UNDEFINED_THREATCATEGORY
    switch v {
        case "undefined":
            result = UNDEFINED_THREATCATEGORY
        case "spam":
            result = SPAM_THREATCATEGORY
        case "phishing":
            result = PHISHING_THREATCATEGORY
        case "malware":
            result = MALWARE_THREATCATEGORY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_THREATCATEGORY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeThreatCategory(values []ThreatCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ThreatCategory) isMultiValue() bool {
    return false
}

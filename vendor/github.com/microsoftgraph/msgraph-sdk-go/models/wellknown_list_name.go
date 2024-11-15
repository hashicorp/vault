package models
type WellknownListName int

const (
    NONE_WELLKNOWNLISTNAME WellknownListName = iota
    DEFAULTLIST_WELLKNOWNLISTNAME
    FLAGGEDEMAILS_WELLKNOWNLISTNAME
    UNKNOWNFUTUREVALUE_WELLKNOWNLISTNAME
)

func (i WellknownListName) String() string {
    return []string{"none", "defaultList", "flaggedEmails", "unknownFutureValue"}[i]
}
func ParseWellknownListName(v string) (any, error) {
    result := NONE_WELLKNOWNLISTNAME
    switch v {
        case "none":
            result = NONE_WELLKNOWNLISTNAME
        case "defaultList":
            result = DEFAULTLIST_WELLKNOWNLISTNAME
        case "flaggedEmails":
            result = FLAGGEDEMAILS_WELLKNOWNLISTNAME
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WELLKNOWNLISTNAME
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWellknownListName(values []WellknownListName) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WellknownListName) isMultiValue() bool {
    return false
}

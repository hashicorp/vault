package models
type OnenoteUserRole int

const (
    NONE_ONENOTEUSERROLE OnenoteUserRole = iota
    OWNER_ONENOTEUSERROLE
    CONTRIBUTOR_ONENOTEUSERROLE
    READER_ONENOTEUSERROLE
)

func (i OnenoteUserRole) String() string {
    return []string{"None", "Owner", "Contributor", "Reader"}[i]
}
func ParseOnenoteUserRole(v string) (any, error) {
    result := NONE_ONENOTEUSERROLE
    switch v {
        case "None":
            result = NONE_ONENOTEUSERROLE
        case "Owner":
            result = OWNER_ONENOTEUSERROLE
        case "Contributor":
            result = CONTRIBUTOR_ONENOTEUSERROLE
        case "Reader":
            result = READER_ONENOTEUSERROLE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOnenoteUserRole(values []OnenoteUserRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnenoteUserRole) isMultiValue() bool {
    return false
}

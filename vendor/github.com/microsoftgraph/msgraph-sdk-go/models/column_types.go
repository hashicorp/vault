package models
type ColumnTypes int

const (
    NOTE_COLUMNTYPES ColumnTypes = iota
    TEXT_COLUMNTYPES
    CHOICE_COLUMNTYPES
    MULTICHOICE_COLUMNTYPES
    NUMBER_COLUMNTYPES
    CURRENCY_COLUMNTYPES
    DATETIME_COLUMNTYPES
    LOOKUP_COLUMNTYPES
    BOOLEAN_COLUMNTYPES
    USER_COLUMNTYPES
    URL_COLUMNTYPES
    CALCULATED_COLUMNTYPES
    LOCATION_COLUMNTYPES
    GEOLOCATION_COLUMNTYPES
    TERM_COLUMNTYPES
    MULTITERM_COLUMNTYPES
    THUMBNAIL_COLUMNTYPES
    APPROVALSTATUS_COLUMNTYPES
    UNKNOWNFUTUREVALUE_COLUMNTYPES
)

func (i ColumnTypes) String() string {
    return []string{"note", "text", "choice", "multichoice", "number", "currency", "dateTime", "lookup", "boolean", "user", "url", "calculated", "location", "geolocation", "term", "multiterm", "thumbnail", "approvalStatus", "unknownFutureValue"}[i]
}
func ParseColumnTypes(v string) (any, error) {
    result := NOTE_COLUMNTYPES
    switch v {
        case "note":
            result = NOTE_COLUMNTYPES
        case "text":
            result = TEXT_COLUMNTYPES
        case "choice":
            result = CHOICE_COLUMNTYPES
        case "multichoice":
            result = MULTICHOICE_COLUMNTYPES
        case "number":
            result = NUMBER_COLUMNTYPES
        case "currency":
            result = CURRENCY_COLUMNTYPES
        case "dateTime":
            result = DATETIME_COLUMNTYPES
        case "lookup":
            result = LOOKUP_COLUMNTYPES
        case "boolean":
            result = BOOLEAN_COLUMNTYPES
        case "user":
            result = USER_COLUMNTYPES
        case "url":
            result = URL_COLUMNTYPES
        case "calculated":
            result = CALCULATED_COLUMNTYPES
        case "location":
            result = LOCATION_COLUMNTYPES
        case "geolocation":
            result = GEOLOCATION_COLUMNTYPES
        case "term":
            result = TERM_COLUMNTYPES
        case "multiterm":
            result = MULTITERM_COLUMNTYPES
        case "thumbnail":
            result = THUMBNAIL_COLUMNTYPES
        case "approvalStatus":
            result = APPROVALSTATUS_COLUMNTYPES
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_COLUMNTYPES
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeColumnTypes(values []ColumnTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ColumnTypes) isMultiValue() bool {
    return false
}

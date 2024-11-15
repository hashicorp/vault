package models
type ExchangeIdFormat int

const (
    ENTRYID_EXCHANGEIDFORMAT ExchangeIdFormat = iota
    EWSID_EXCHANGEIDFORMAT
    IMMUTABLEENTRYID_EXCHANGEIDFORMAT
    RESTID_EXCHANGEIDFORMAT
    RESTIMMUTABLEENTRYID_EXCHANGEIDFORMAT
)

func (i ExchangeIdFormat) String() string {
    return []string{"entryId", "ewsId", "immutableEntryId", "restId", "restImmutableEntryId"}[i]
}
func ParseExchangeIdFormat(v string) (any, error) {
    result := ENTRYID_EXCHANGEIDFORMAT
    switch v {
        case "entryId":
            result = ENTRYID_EXCHANGEIDFORMAT
        case "ewsId":
            result = EWSID_EXCHANGEIDFORMAT
        case "immutableEntryId":
            result = IMMUTABLEENTRYID_EXCHANGEIDFORMAT
        case "restId":
            result = RESTID_EXCHANGEIDFORMAT
        case "restImmutableEntryId":
            result = RESTIMMUTABLEENTRYID_EXCHANGEIDFORMAT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExchangeIdFormat(values []ExchangeIdFormat) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExchangeIdFormat) isMultiValue() bool {
    return false
}

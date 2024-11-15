package models
type FreeBusyStatus int

const (
    UNKNOWN_FREEBUSYSTATUS FreeBusyStatus = iota
    FREE_FREEBUSYSTATUS
    TENTATIVE_FREEBUSYSTATUS
    BUSY_FREEBUSYSTATUS
    OOF_FREEBUSYSTATUS
    WORKINGELSEWHERE_FREEBUSYSTATUS
)

func (i FreeBusyStatus) String() string {
    return []string{"unknown", "free", "tentative", "busy", "oof", "workingElsewhere"}[i]
}
func ParseFreeBusyStatus(v string) (any, error) {
    result := UNKNOWN_FREEBUSYSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_FREEBUSYSTATUS
        case "free":
            result = FREE_FREEBUSYSTATUS
        case "tentative":
            result = TENTATIVE_FREEBUSYSTATUS
        case "busy":
            result = BUSY_FREEBUSYSTATUS
        case "oof":
            result = OOF_FREEBUSYSTATUS
        case "workingElsewhere":
            result = WORKINGELSEWHERE_FREEBUSYSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFreeBusyStatus(values []FreeBusyStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FreeBusyStatus) isMultiValue() bool {
    return false
}

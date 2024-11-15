package models
type EventType int

const (
    SINGLEINSTANCE_EVENTTYPE EventType = iota
    OCCURRENCE_EVENTTYPE
    EXCEPTION_EVENTTYPE
    SERIESMASTER_EVENTTYPE
)

func (i EventType) String() string {
    return []string{"singleInstance", "occurrence", "exception", "seriesMaster"}[i]
}
func ParseEventType(v string) (any, error) {
    result := SINGLEINSTANCE_EVENTTYPE
    switch v {
        case "singleInstance":
            result = SINGLEINSTANCE_EVENTTYPE
        case "occurrence":
            result = OCCURRENCE_EVENTTYPE
        case "exception":
            result = EXCEPTION_EVENTTYPE
        case "seriesMaster":
            result = SERIESMASTER_EVENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEventType(values []EventType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EventType) isMultiValue() bool {
    return false
}

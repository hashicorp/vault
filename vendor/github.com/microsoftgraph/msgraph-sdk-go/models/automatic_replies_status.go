package models
type AutomaticRepliesStatus int

const (
    DISABLED_AUTOMATICREPLIESSTATUS AutomaticRepliesStatus = iota
    ALWAYSENABLED_AUTOMATICREPLIESSTATUS
    SCHEDULED_AUTOMATICREPLIESSTATUS
)

func (i AutomaticRepliesStatus) String() string {
    return []string{"disabled", "alwaysEnabled", "scheduled"}[i]
}
func ParseAutomaticRepliesStatus(v string) (any, error) {
    result := DISABLED_AUTOMATICREPLIESSTATUS
    switch v {
        case "disabled":
            result = DISABLED_AUTOMATICREPLIESSTATUS
        case "alwaysEnabled":
            result = ALWAYSENABLED_AUTOMATICREPLIESSTATUS
        case "scheduled":
            result = SCHEDULED_AUTOMATICREPLIESSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAutomaticRepliesStatus(values []AutomaticRepliesStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AutomaticRepliesStatus) isMultiValue() bool {
    return false
}

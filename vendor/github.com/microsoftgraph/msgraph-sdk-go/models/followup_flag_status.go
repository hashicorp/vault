package models
type FollowupFlagStatus int

const (
    NOTFLAGGED_FOLLOWUPFLAGSTATUS FollowupFlagStatus = iota
    COMPLETE_FOLLOWUPFLAGSTATUS
    FLAGGED_FOLLOWUPFLAGSTATUS
)

func (i FollowupFlagStatus) String() string {
    return []string{"notFlagged", "complete", "flagged"}[i]
}
func ParseFollowupFlagStatus(v string) (any, error) {
    result := NOTFLAGGED_FOLLOWUPFLAGSTATUS
    switch v {
        case "notFlagged":
            result = NOTFLAGGED_FOLLOWUPFLAGSTATUS
        case "complete":
            result = COMPLETE_FOLLOWUPFLAGSTATUS
        case "flagged":
            result = FLAGGED_FOLLOWUPFLAGSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFollowupFlagStatus(values []FollowupFlagStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FollowupFlagStatus) isMultiValue() bool {
    return false
}

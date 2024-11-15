package models
// State of the action on the device
type ActionState int

const (
    // Not a valid action state
    NONE_ACTIONSTATE ActionState = iota
    // Action is pending
    PENDING_ACTIONSTATE
    // Action has been cancelled.
    CANCELED_ACTIONSTATE
    // Action is active.
    ACTIVE_ACTIONSTATE
    // Action completed without errors.
    DONE_ACTIONSTATE
    // Action failed
    FAILED_ACTIONSTATE
    // Action is not supported.
    NOTSUPPORTED_ACTIONSTATE
)

func (i ActionState) String() string {
    return []string{"none", "pending", "canceled", "active", "done", "failed", "notSupported"}[i]
}
func ParseActionState(v string) (any, error) {
    result := NONE_ACTIONSTATE
    switch v {
        case "none":
            result = NONE_ACTIONSTATE
        case "pending":
            result = PENDING_ACTIONSTATE
        case "canceled":
            result = CANCELED_ACTIONSTATE
        case "active":
            result = ACTIVE_ACTIONSTATE
        case "done":
            result = DONE_ACTIONSTATE
        case "failed":
            result = FAILED_ACTIONSTATE
        case "notSupported":
            result = NOTSUPPORTED_ACTIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeActionState(values []ActionState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ActionState) isMultiValue() bool {
    return false
}

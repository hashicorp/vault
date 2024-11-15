package models
// Possible states associated with an Apple Volume Purchase Program token.
type VppTokenState int

const (
    // Default state.
    UNKNOWN_VPPTOKENSTATE VppTokenState = iota
    // Token is valid.
    VALID_VPPTOKENSTATE
    // Token is expired.
    EXPIRED_VPPTOKENSTATE
    // Token is invalid.
    INVALID_VPPTOKENSTATE
    // Token is managed by another MDM Service.
    ASSIGNEDTOEXTERNALMDM_VPPTOKENSTATE
)

func (i VppTokenState) String() string {
    return []string{"unknown", "valid", "expired", "invalid", "assignedToExternalMDM"}[i]
}
func ParseVppTokenState(v string) (any, error) {
    result := UNKNOWN_VPPTOKENSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_VPPTOKENSTATE
        case "valid":
            result = VALID_VPPTOKENSTATE
        case "expired":
            result = EXPIRED_VPPTOKENSTATE
        case "invalid":
            result = INVALID_VPPTOKENSTATE
        case "assignedToExternalMDM":
            result = ASSIGNEDTOEXTERNALMDM_VPPTOKENSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVppTokenState(values []VppTokenState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VppTokenState) isMultiValue() bool {
    return false
}

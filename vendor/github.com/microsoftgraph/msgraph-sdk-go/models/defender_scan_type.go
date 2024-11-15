package models
// Possible values for system scan type.
type DefenderScanType int

const (
    // User Defined, default value, no intent.
    USERDEFINED_DEFENDERSCANTYPE DefenderScanType = iota
    // System scan disabled.
    DISABLED_DEFENDERSCANTYPE
    // Quick system scan.
    QUICK_DEFENDERSCANTYPE
    // Full system scan.
    FULL_DEFENDERSCANTYPE
)

func (i DefenderScanType) String() string {
    return []string{"userDefined", "disabled", "quick", "full"}[i]
}
func ParseDefenderScanType(v string) (any, error) {
    result := USERDEFINED_DEFENDERSCANTYPE
    switch v {
        case "userDefined":
            result = USERDEFINED_DEFENDERSCANTYPE
        case "disabled":
            result = DISABLED_DEFENDERSCANTYPE
        case "quick":
            result = QUICK_DEFENDERSCANTYPE
        case "full":
            result = FULL_DEFENDERSCANTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDefenderScanType(values []DefenderScanType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DefenderScanType) isMultiValue() bool {
    return false
}

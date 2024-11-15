package models
// Possible values of Cloud Block Level
type DefenderCloudBlockLevelType int

const (
    // Default value, uses the default Windows Defender Antivirus blocking level and provides strong detection without increasing the risk of detecting legitimate files
    NOTCONFIGURED_DEFENDERCLOUDBLOCKLEVELTYPE DefenderCloudBlockLevelType = iota
    // High applies a strong level of detection.
    HIGH_DEFENDERCLOUDBLOCKLEVELTYPE
    // High + uses the High level and applies addition protection measures
    HIGHPLUS_DEFENDERCLOUDBLOCKLEVELTYPE
    // Zero tolerance blocks all unknown executables
    ZEROTOLERANCE_DEFENDERCLOUDBLOCKLEVELTYPE
)

func (i DefenderCloudBlockLevelType) String() string {
    return []string{"notConfigured", "high", "highPlus", "zeroTolerance"}[i]
}
func ParseDefenderCloudBlockLevelType(v string) (any, error) {
    result := NOTCONFIGURED_DEFENDERCLOUDBLOCKLEVELTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_DEFENDERCLOUDBLOCKLEVELTYPE
        case "high":
            result = HIGH_DEFENDERCLOUDBLOCKLEVELTYPE
        case "highPlus":
            result = HIGHPLUS_DEFENDERCLOUDBLOCKLEVELTYPE
        case "zeroTolerance":
            result = ZEROTOLERANCE_DEFENDERCLOUDBLOCKLEVELTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDefenderCloudBlockLevelType(values []DefenderCloudBlockLevelType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DefenderCloudBlockLevelType) isMultiValue() bool {
    return false
}

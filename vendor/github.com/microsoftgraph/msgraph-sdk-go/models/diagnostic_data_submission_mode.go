package models
// Allow the device to send diagnostic and usage telemetry data, such as Watson.
type DiagnosticDataSubmissionMode int

const (
    // Allow the user to set.
    USERDEFINED_DIAGNOSTICDATASUBMISSIONMODE DiagnosticDataSubmissionMode = iota
    // No telemetry data is sent from OS components. Note: This value is only applicable to enterprise and server devices. Using this setting on other devices is equivalent to setting the value of 1.
    NONE_DIAGNOSTICDATASUBMISSIONMODE
    // Sends basic telemetry data.
    BASIC_DIAGNOSTICDATASUBMISSIONMODE
    // Sends enhanced telemetry data including usage and insights data.
    ENHANCED_DIAGNOSTICDATASUBMISSIONMODE
    // Sends full telemetry data including diagnostic data, such as system state.
    FULL_DIAGNOSTICDATASUBMISSIONMODE
)

func (i DiagnosticDataSubmissionMode) String() string {
    return []string{"userDefined", "none", "basic", "enhanced", "full"}[i]
}
func ParseDiagnosticDataSubmissionMode(v string) (any, error) {
    result := USERDEFINED_DIAGNOSTICDATASUBMISSIONMODE
    switch v {
        case "userDefined":
            result = USERDEFINED_DIAGNOSTICDATASUBMISSIONMODE
        case "none":
            result = NONE_DIAGNOSTICDATASUBMISSIONMODE
        case "basic":
            result = BASIC_DIAGNOSTICDATASUBMISSIONMODE
        case "enhanced":
            result = ENHANCED_DIAGNOSTICDATASUBMISSIONMODE
        case "full":
            result = FULL_DIAGNOSTICDATASUBMISSIONMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDiagnosticDataSubmissionMode(values []DiagnosticDataSubmissionMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DiagnosticDataSubmissionMode) isMultiValue() bool {
    return false
}

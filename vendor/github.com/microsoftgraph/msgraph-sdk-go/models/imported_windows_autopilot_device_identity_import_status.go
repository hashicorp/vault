package models
type ImportedWindowsAutopilotDeviceIdentityImportStatus int

const (
    // Unknown status.
    UNKNOWN_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS ImportedWindowsAutopilotDeviceIdentityImportStatus = iota
    // Pending status.
    PENDING_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
    // Partial status.
    PARTIAL_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
    // Complete status.
    COMPLETE_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
    // Error status.
    ERROR_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
)

func (i ImportedWindowsAutopilotDeviceIdentityImportStatus) String() string {
    return []string{"unknown", "pending", "partial", "complete", "error"}[i]
}
func ParseImportedWindowsAutopilotDeviceIdentityImportStatus(v string) (any, error) {
    result := UNKNOWN_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
        case "pending":
            result = PENDING_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
        case "partial":
            result = PARTIAL_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
        case "complete":
            result = COMPLETE_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
        case "error":
            result = ERROR_IMPORTEDWINDOWSAUTOPILOTDEVICEIDENTITYIMPORTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeImportedWindowsAutopilotDeviceIdentityImportStatus(values []ImportedWindowsAutopilotDeviceIdentityImportStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ImportedWindowsAutopilotDeviceIdentityImportStatus) isMultiValue() bool {
    return false
}

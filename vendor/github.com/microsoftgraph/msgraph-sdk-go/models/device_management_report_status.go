package models
// Possible statuses associated with a generated report.
type DeviceManagementReportStatus int

const (
    // Report generation status is unknown.
    UNKNOWN_DEVICEMANAGEMENTREPORTSTATUS DeviceManagementReportStatus = iota
    // Report generation has not started.
    NOTSTARTED_DEVICEMANAGEMENTREPORTSTATUS
    // Report generation is in progress.
    INPROGRESS_DEVICEMANAGEMENTREPORTSTATUS
    // Report generation is completed.
    COMPLETED_DEVICEMANAGEMENTREPORTSTATUS
    // Report generation has failed.
    FAILED_DEVICEMANAGEMENTREPORTSTATUS
)

func (i DeviceManagementReportStatus) String() string {
    return []string{"unknown", "notStarted", "inProgress", "completed", "failed"}[i]
}
func ParseDeviceManagementReportStatus(v string) (any, error) {
    result := UNKNOWN_DEVICEMANAGEMENTREPORTSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_DEVICEMANAGEMENTREPORTSTATUS
        case "notStarted":
            result = NOTSTARTED_DEVICEMANAGEMENTREPORTSTATUS
        case "inProgress":
            result = INPROGRESS_DEVICEMANAGEMENTREPORTSTATUS
        case "completed":
            result = COMPLETED_DEVICEMANAGEMENTREPORTSTATUS
        case "failed":
            result = FAILED_DEVICEMANAGEMENTREPORTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementReportStatus(values []DeviceManagementReportStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementReportStatus) isMultiValue() bool {
    return false
}

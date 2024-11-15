package models
// Possible values for the file format of a report.
type DeviceManagementReportFileFormat int

const (
    // CSV Format.
    CSV_DEVICEMANAGEMENTREPORTFILEFORMAT DeviceManagementReportFileFormat = iota
    // PDF Format (Deprecate later).
    PDF_DEVICEMANAGEMENTREPORTFILEFORMAT
    // JSON Format.
    JSON_DEVICEMANAGEMENTREPORTFILEFORMAT
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTREPORTFILEFORMAT
)

func (i DeviceManagementReportFileFormat) String() string {
    return []string{"csv", "pdf", "json", "unknownFutureValue"}[i]
}
func ParseDeviceManagementReportFileFormat(v string) (any, error) {
    result := CSV_DEVICEMANAGEMENTREPORTFILEFORMAT
    switch v {
        case "csv":
            result = CSV_DEVICEMANAGEMENTREPORTFILEFORMAT
        case "pdf":
            result = PDF_DEVICEMANAGEMENTREPORTFILEFORMAT
        case "json":
            result = JSON_DEVICEMANAGEMENTREPORTFILEFORMAT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTREPORTFILEFORMAT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementReportFileFormat(values []DeviceManagementReportFileFormat) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementReportFileFormat) isMultiValue() bool {
    return false
}

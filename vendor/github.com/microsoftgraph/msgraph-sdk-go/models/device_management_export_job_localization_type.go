package models
// Configures how the requested export job is localized.
type DeviceManagementExportJobLocalizationType int

const (
    // Configures the export job to expose localized values as an additional column.
    LOCALIZEDVALUESASADDITIONALCOLUMN_DEVICEMANAGEMENTEXPORTJOBLOCALIZATIONTYPE DeviceManagementExportJobLocalizationType = iota
    // Configures the export job to replace enumerable values with their localized values.
    REPLACELOCALIZABLEVALUES_DEVICEMANAGEMENTEXPORTJOBLOCALIZATIONTYPE
)

func (i DeviceManagementExportJobLocalizationType) String() string {
    return []string{"localizedValuesAsAdditionalColumn", "replaceLocalizableValues"}[i]
}
func ParseDeviceManagementExportJobLocalizationType(v string) (any, error) {
    result := LOCALIZEDVALUESASADDITIONALCOLUMN_DEVICEMANAGEMENTEXPORTJOBLOCALIZATIONTYPE
    switch v {
        case "localizedValuesAsAdditionalColumn":
            result = LOCALIZEDVALUESASADDITIONALCOLUMN_DEVICEMANAGEMENTEXPORTJOBLOCALIZATIONTYPE
        case "replaceLocalizableValues":
            result = REPLACELOCALIZABLEVALUES_DEVICEMANAGEMENTEXPORTJOBLOCALIZATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementExportJobLocalizationType(values []DeviceManagementExportJobLocalizationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementExportJobLocalizationType) isMultiValue() bool {
    return false
}

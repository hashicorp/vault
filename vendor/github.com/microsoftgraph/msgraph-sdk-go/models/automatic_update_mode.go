package models
// Possible values for automatic update mode.
type AutomaticUpdateMode int

const (
    // User Defined, default value, no intent.
    USERDEFINED_AUTOMATICUPDATEMODE AutomaticUpdateMode = iota
    // Notify on download.
    NOTIFYDOWNLOAD_AUTOMATICUPDATEMODE
    // Auto-install at maintenance time.
    AUTOINSTALLATMAINTENANCETIME_AUTOMATICUPDATEMODE
    // Auto-install and reboot at maintenance time.
    AUTOINSTALLANDREBOOTATMAINTENANCETIME_AUTOMATICUPDATEMODE
    // Auto-install and reboot at scheduled time.
    AUTOINSTALLANDREBOOTATSCHEDULEDTIME_AUTOMATICUPDATEMODE
    // Auto-install and restart without end-user control
    AUTOINSTALLANDREBOOTWITHOUTENDUSERCONTROL_AUTOMATICUPDATEMODE
)

func (i AutomaticUpdateMode) String() string {
    return []string{"userDefined", "notifyDownload", "autoInstallAtMaintenanceTime", "autoInstallAndRebootAtMaintenanceTime", "autoInstallAndRebootAtScheduledTime", "autoInstallAndRebootWithoutEndUserControl"}[i]
}
func ParseAutomaticUpdateMode(v string) (any, error) {
    result := USERDEFINED_AUTOMATICUPDATEMODE
    switch v {
        case "userDefined":
            result = USERDEFINED_AUTOMATICUPDATEMODE
        case "notifyDownload":
            result = NOTIFYDOWNLOAD_AUTOMATICUPDATEMODE
        case "autoInstallAtMaintenanceTime":
            result = AUTOINSTALLATMAINTENANCETIME_AUTOMATICUPDATEMODE
        case "autoInstallAndRebootAtMaintenanceTime":
            result = AUTOINSTALLANDREBOOTATMAINTENANCETIME_AUTOMATICUPDATEMODE
        case "autoInstallAndRebootAtScheduledTime":
            result = AUTOINSTALLANDREBOOTATSCHEDULEDTIME_AUTOMATICUPDATEMODE
        case "autoInstallAndRebootWithoutEndUserControl":
            result = AUTOINSTALLANDREBOOTWITHOUTENDUSERCONTROL_AUTOMATICUPDATEMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAutomaticUpdateMode(values []AutomaticUpdateMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AutomaticUpdateMode) isMultiValue() bool {
    return false
}

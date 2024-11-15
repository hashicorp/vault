package models
// Possible values for monitoring file activity.
type DefenderMonitorFileActivity int

const (
    // User Defined, default value, no intent.
    USERDEFINED_DEFENDERMONITORFILEACTIVITY DefenderMonitorFileActivity = iota
    // Disable monitoring file activity.
    DISABLE_DEFENDERMONITORFILEACTIVITY
    // Monitor all files.
    MONITORALLFILES_DEFENDERMONITORFILEACTIVITY
    //  Monitor incoming files only.
    MONITORINCOMINGFILESONLY_DEFENDERMONITORFILEACTIVITY
    // Monitor outgoing files only.
    MONITOROUTGOINGFILESONLY_DEFENDERMONITORFILEACTIVITY
)

func (i DefenderMonitorFileActivity) String() string {
    return []string{"userDefined", "disable", "monitorAllFiles", "monitorIncomingFilesOnly", "monitorOutgoingFilesOnly"}[i]
}
func ParseDefenderMonitorFileActivity(v string) (any, error) {
    result := USERDEFINED_DEFENDERMONITORFILEACTIVITY
    switch v {
        case "userDefined":
            result = USERDEFINED_DEFENDERMONITORFILEACTIVITY
        case "disable":
            result = DISABLE_DEFENDERMONITORFILEACTIVITY
        case "monitorAllFiles":
            result = MONITORALLFILES_DEFENDERMONITORFILEACTIVITY
        case "monitorIncomingFilesOnly":
            result = MONITORINCOMINGFILESONLY_DEFENDERMONITORFILEACTIVITY
        case "monitorOutgoingFilesOnly":
            result = MONITOROUTGOINGFILESONLY_DEFENDERMONITORFILEACTIVITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDefenderMonitorFileActivity(values []DefenderMonitorFileActivity) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DefenderMonitorFileActivity) isMultiValue() bool {
    return false
}

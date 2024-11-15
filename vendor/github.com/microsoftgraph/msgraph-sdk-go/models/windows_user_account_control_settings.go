package models
// Possible values for Windows user account control settings.
type WindowsUserAccountControlSettings int

const (
    // User Defined, default value, no intent.
    USERDEFINED_WINDOWSUSERACCOUNTCONTROLSETTINGS WindowsUserAccountControlSettings = iota
    // Always notify.
    ALWAYSNOTIFY_WINDOWSUSERACCOUNTCONTROLSETTINGS
    // Notify on app changes.
    NOTIFYONAPPCHANGES_WINDOWSUSERACCOUNTCONTROLSETTINGS
    // Notify on app changes without dimming desktop.
    NOTIFYONAPPCHANGESWITHOUTDIMMING_WINDOWSUSERACCOUNTCONTROLSETTINGS
    // Never notify.
    NEVERNOTIFY_WINDOWSUSERACCOUNTCONTROLSETTINGS
)

func (i WindowsUserAccountControlSettings) String() string {
    return []string{"userDefined", "alwaysNotify", "notifyOnAppChanges", "notifyOnAppChangesWithoutDimming", "neverNotify"}[i]
}
func ParseWindowsUserAccountControlSettings(v string) (any, error) {
    result := USERDEFINED_WINDOWSUSERACCOUNTCONTROLSETTINGS
    switch v {
        case "userDefined":
            result = USERDEFINED_WINDOWSUSERACCOUNTCONTROLSETTINGS
        case "alwaysNotify":
            result = ALWAYSNOTIFY_WINDOWSUSERACCOUNTCONTROLSETTINGS
        case "notifyOnAppChanges":
            result = NOTIFYONAPPCHANGES_WINDOWSUSERACCOUNTCONTROLSETTINGS
        case "notifyOnAppChangesWithoutDimming":
            result = NOTIFYONAPPCHANGESWITHOUTDIMMING_WINDOWSUSERACCOUNTCONTROLSETTINGS
        case "neverNotify":
            result = NEVERNOTIFY_WINDOWSUSERACCOUNTCONTROLSETTINGS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsUserAccountControlSettings(values []WindowsUserAccountControlSettings) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsUserAccountControlSettings) isMultiValue() bool {
    return false
}

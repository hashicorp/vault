package models
// Contains value for notification status.
type Win32LobAppNotification int

const (
    // Show all notifications.
    SHOWALL_WIN32LOBAPPNOTIFICATION Win32LobAppNotification = iota
    // Only show restart notification and suppress other notifications.
    SHOWREBOOT_WIN32LOBAPPNOTIFICATION
    // Hide all notifications.
    HIDEALL_WIN32LOBAPPNOTIFICATION
)

func (i Win32LobAppNotification) String() string {
    return []string{"showAll", "showReboot", "hideAll"}[i]
}
func ParseWin32LobAppNotification(v string) (any, error) {
    result := SHOWALL_WIN32LOBAPPNOTIFICATION
    switch v {
        case "showAll":
            result = SHOWALL_WIN32LOBAPPNOTIFICATION
        case "showReboot":
            result = SHOWREBOOT_WIN32LOBAPPNOTIFICATION
        case "hideAll":
            result = HIDEALL_WIN32LOBAPPNOTIFICATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppNotification(values []Win32LobAppNotification) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppNotification) isMultiValue() bool {
    return false
}

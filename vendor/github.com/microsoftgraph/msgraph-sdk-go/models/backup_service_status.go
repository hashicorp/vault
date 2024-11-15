package models
type BackupServiceStatus int

const (
    DISABLED_BACKUPSERVICESTATUS BackupServiceStatus = iota
    ENABLED_BACKUPSERVICESTATUS
    PROTECTIONCHANGELOCKED_BACKUPSERVICESTATUS
    RESTORELOCKED_BACKUPSERVICESTATUS
    UNKNOWNFUTUREVALUE_BACKUPSERVICESTATUS
)

func (i BackupServiceStatus) String() string {
    return []string{"disabled", "enabled", "protectionChangeLocked", "restoreLocked", "unknownFutureValue"}[i]
}
func ParseBackupServiceStatus(v string) (any, error) {
    result := DISABLED_BACKUPSERVICESTATUS
    switch v {
        case "disabled":
            result = DISABLED_BACKUPSERVICESTATUS
        case "enabled":
            result = ENABLED_BACKUPSERVICESTATUS
        case "protectionChangeLocked":
            result = PROTECTIONCHANGELOCKED_BACKUPSERVICESTATUS
        case "restoreLocked":
            result = RESTORELOCKED_BACKUPSERVICESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BACKUPSERVICESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBackupServiceStatus(values []BackupServiceStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BackupServiceStatus) isMultiValue() bool {
    return false
}

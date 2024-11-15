package models
type BackupServiceConsumer int

const (
    UNKNOWN_BACKUPSERVICECONSUMER BackupServiceConsumer = iota
    FIRSTPARTY_BACKUPSERVICECONSUMER
    THIRDPARTY_BACKUPSERVICECONSUMER
    UNKNOWNFUTUREVALUE_BACKUPSERVICECONSUMER
)

func (i BackupServiceConsumer) String() string {
    return []string{"unknown", "firstparty", "thirdparty", "unknownFutureValue"}[i]
}
func ParseBackupServiceConsumer(v string) (any, error) {
    result := UNKNOWN_BACKUPSERVICECONSUMER
    switch v {
        case "unknown":
            result = UNKNOWN_BACKUPSERVICECONSUMER
        case "firstparty":
            result = FIRSTPARTY_BACKUPSERVICECONSUMER
        case "thirdparty":
            result = THIRDPARTY_BACKUPSERVICECONSUMER
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BACKUPSERVICECONSUMER
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBackupServiceConsumer(values []BackupServiceConsumer) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BackupServiceConsumer) isMultiValue() bool {
    return false
}

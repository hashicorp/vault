package models
type RestorePointPreference int

const (
    LATEST_RESTOREPOINTPREFERENCE RestorePointPreference = iota
    OLDEST_RESTOREPOINTPREFERENCE
    UNKNOWNFUTUREVALUE_RESTOREPOINTPREFERENCE
)

func (i RestorePointPreference) String() string {
    return []string{"latest", "oldest", "unknownFutureValue"}[i]
}
func ParseRestorePointPreference(v string) (any, error) {
    result := LATEST_RESTOREPOINTPREFERENCE
    switch v {
        case "latest":
            result = LATEST_RESTOREPOINTPREFERENCE
        case "oldest":
            result = OLDEST_RESTOREPOINTPREFERENCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RESTOREPOINTPREFERENCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRestorePointPreference(values []RestorePointPreference) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RestorePointPreference) isMultiValue() bool {
    return false
}

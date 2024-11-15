package models
// A list of possible operations for rules used to make determinations about an application based on files or folders. Unless noted, can be used with either detection or requirement rules.
type Win32LobAppFileSystemOperationType int

const (
    // Default. Indicates that the rule does not have the operation type configured.
    NOTCONFIGURED_WIN32LOBAPPFILESYSTEMOPERATIONTYPE Win32LobAppFileSystemOperationType = iota
    // Indicates that the rule evaluates whether the specified file or folder exists.
    EXISTS_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
    // Indicates that the rule compares the modified date of the specified file against a provided comparison value by DateTime comparison.
    MODIFIEDDATE_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
    // Indicates that the rule compares the created date of the specified file against a provided comparison value by DateTime comparison.
    CREATEDDATE_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
    // Indicates that the rule compares the detected version of the specified file against a provided comparison value via version semantics (both operand values will be parsed as versions and directly compared). If the value read at the given registry value is not discovered to be in version-compatible format, a string comparison will be used instead.
    VERSION_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
    // Indicates that the rule compares the size of the file in MiB (rounded down) against a provided comparison value by integer comparison.
    SIZEINMB_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
)

func (i Win32LobAppFileSystemOperationType) String() string {
    return []string{"notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB"}[i]
}
func ParseWin32LobAppFileSystemOperationType(v string) (any, error) {
    result := NOTCONFIGURED_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        case "exists":
            result = EXISTS_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        case "modifiedDate":
            result = MODIFIEDDATE_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        case "createdDate":
            result = CREATEDDATE_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        case "version":
            result = VERSION_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        case "sizeInMB":
            result = SIZEINMB_WIN32LOBAPPFILESYSTEMOPERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppFileSystemOperationType(values []Win32LobAppFileSystemOperationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppFileSystemOperationType) isMultiValue() bool {
    return false
}

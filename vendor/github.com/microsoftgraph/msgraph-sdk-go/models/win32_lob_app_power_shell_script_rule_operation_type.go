package models
// Contains all supported Powershell Script output detection type.
type Win32LobAppPowerShellScriptRuleOperationType int

const (
    // Not configured.
    NOTCONFIGURED_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE Win32LobAppPowerShellScriptRuleOperationType = iota
    // Output data type is string.
    STRING_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    // Output data type is date time.
    DATETIME_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    // Output data type is integer.
    INTEGER_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    // Output data type is float.
    FLOAT_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    // Output data type is version.
    VERSION_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    // Output data type is boolean.
    BOOLEAN_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
)

func (i Win32LobAppPowerShellScriptRuleOperationType) String() string {
    return []string{"notConfigured", "string", "dateTime", "integer", "float", "version", "boolean"}[i]
}
func ParseWin32LobAppPowerShellScriptRuleOperationType(v string) (any, error) {
    result := NOTCONFIGURED_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "string":
            result = STRING_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "dateTime":
            result = DATETIME_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "integer":
            result = INTEGER_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "float":
            result = FLOAT_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "version":
            result = VERSION_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        case "boolean":
            result = BOOLEAN_WIN32LOBAPPPOWERSHELLSCRIPTRULEOPERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppPowerShellScriptRuleOperationType(values []Win32LobAppPowerShellScriptRuleOperationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppPowerShellScriptRuleOperationType) isMultiValue() bool {
    return false
}

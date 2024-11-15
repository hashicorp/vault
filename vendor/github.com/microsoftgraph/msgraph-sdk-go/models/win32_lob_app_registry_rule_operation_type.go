package models
// A list of possible operations for rules used to make determinations about an application based on registry keys or values. Unless noted, the values can be used with either detection or requirement rules.
type Win32LobAppRegistryRuleOperationType int

const (
    // Default. Indicates that the rule does not have the operation type configured.
    NOTCONFIGURED_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE Win32LobAppRegistryRuleOperationType = iota
    // Indicates that the rule evaluates whether the specified registry key or value exists.
    EXISTS_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
    // Indicates that the rule evaluates whether the specified registry key or value does not exist. It is the functional inverse of an equivalent rule that uses operation type `exists`.
    DOESNOTEXIST_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
    // Indicates that the rule compares the value read at the given registry value against a provided comparison value by string comparison.
    STRING_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
    // Indicates that the rule compares the value read at the given registry value against a provided comparison value by integer comparison.
    INTEGER_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
    // Indicates that the rule compares the value read at the given registry value against a provided comparison value via version semantics (both operand values will be parsed as versions and directly compared). If the value read at the given registry value is not discovered to be in version-compatible format, a string comparison will be used instead.
    VERSION_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
)

func (i Win32LobAppRegistryRuleOperationType) String() string {
    return []string{"notConfigured", "exists", "doesNotExist", "string", "integer", "version"}[i]
}
func ParseWin32LobAppRegistryRuleOperationType(v string) (any, error) {
    result := NOTCONFIGURED_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        case "exists":
            result = EXISTS_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        case "doesNotExist":
            result = DOESNOTEXIST_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        case "string":
            result = STRING_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        case "integer":
            result = INTEGER_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        case "version":
            result = VERSION_WIN32LOBAPPREGISTRYRULEOPERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppRegistryRuleOperationType(values []Win32LobAppRegistryRuleOperationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppRegistryRuleOperationType) isMultiValue() bool {
    return false
}

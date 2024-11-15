package models
type AuthenticationMethodsPolicyMigrationState int

const (
    PREMIGRATION_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE AuthenticationMethodsPolicyMigrationState = iota
    MIGRATIONINPROGRESS_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
    MIGRATIONCOMPLETE_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
)

func (i AuthenticationMethodsPolicyMigrationState) String() string {
    return []string{"preMigration", "migrationInProgress", "migrationComplete", "unknownFutureValue"}[i]
}
func ParseAuthenticationMethodsPolicyMigrationState(v string) (any, error) {
    result := PREMIGRATION_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
    switch v {
        case "preMigration":
            result = PREMIGRATION_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
        case "migrationInProgress":
            result = MIGRATIONINPROGRESS_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
        case "migrationComplete":
            result = MIGRATIONCOMPLETE_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationMethodsPolicyMigrationState(values []AuthenticationMethodsPolicyMigrationState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodsPolicyMigrationState) isMultiValue() bool {
    return false
}

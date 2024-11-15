package models
// Possible values of AppLocker Application Control Types
type AppLockerApplicationControlType int

const (
    // Device default value, no Application Control type selected.
    NOTCONFIGURED_APPLOCKERAPPLICATIONCONTROLTYPE AppLockerApplicationControlType = iota
    // Enforce Windows component and store apps.
    ENFORCECOMPONENTSANDSTOREAPPS_APPLOCKERAPPLICATIONCONTROLTYPE
    // Audit Windows component and store apps.
    AUDITCOMPONENTSANDSTOREAPPS_APPLOCKERAPPLICATIONCONTROLTYPE
    // Enforce Windows components, store apps and smart locker.
    ENFORCECOMPONENTSSTOREAPPSANDSMARTLOCKER_APPLOCKERAPPLICATIONCONTROLTYPE
    // Audit Windows components, store apps and smart locker​.
    AUDITCOMPONENTSSTOREAPPSANDSMARTLOCKER_APPLOCKERAPPLICATIONCONTROLTYPE
)

func (i AppLockerApplicationControlType) String() string {
    return []string{"notConfigured", "enforceComponentsAndStoreApps", "auditComponentsAndStoreApps", "enforceComponentsStoreAppsAndSmartlocker", "auditComponentsStoreAppsAndSmartlocker"}[i]
}
func ParseAppLockerApplicationControlType(v string) (any, error) {
    result := NOTCONFIGURED_APPLOCKERAPPLICATIONCONTROLTYPE
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_APPLOCKERAPPLICATIONCONTROLTYPE
        case "enforceComponentsAndStoreApps":
            result = ENFORCECOMPONENTSANDSTOREAPPS_APPLOCKERAPPLICATIONCONTROLTYPE
        case "auditComponentsAndStoreApps":
            result = AUDITCOMPONENTSANDSTOREAPPS_APPLOCKERAPPLICATIONCONTROLTYPE
        case "enforceComponentsStoreAppsAndSmartlocker":
            result = ENFORCECOMPONENTSSTOREAPPSANDSMARTLOCKER_APPLOCKERAPPLICATIONCONTROLTYPE
        case "auditComponentsStoreAppsAndSmartlocker":
            result = AUDITCOMPONENTSSTOREAPPSANDSMARTLOCKER_APPLOCKERAPPLICATIONCONTROLTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppLockerApplicationControlType(values []AppLockerApplicationControlType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppLockerApplicationControlType) isMultiValue() bool {
    return false
}

package models
type RegistryHive int

const (
    UNKNOWN_REGISTRYHIVE RegistryHive = iota
    CURRENTCONFIG_REGISTRYHIVE
    CURRENTUSER_REGISTRYHIVE
    LOCALMACHINESAM_REGISTRYHIVE
    LOCALMACHINESECURITY_REGISTRYHIVE
    LOCALMACHINESOFTWARE_REGISTRYHIVE
    LOCALMACHINESYSTEM_REGISTRYHIVE
    USERSDEFAULT_REGISTRYHIVE
    UNKNOWNFUTUREVALUE_REGISTRYHIVE
)

func (i RegistryHive) String() string {
    return []string{"unknown", "currentConfig", "currentUser", "localMachineSam", "localMachineSecurity", "localMachineSoftware", "localMachineSystem", "usersDefault", "unknownFutureValue"}[i]
}
func ParseRegistryHive(v string) (any, error) {
    result := UNKNOWN_REGISTRYHIVE
    switch v {
        case "unknown":
            result = UNKNOWN_REGISTRYHIVE
        case "currentConfig":
            result = CURRENTCONFIG_REGISTRYHIVE
        case "currentUser":
            result = CURRENTUSER_REGISTRYHIVE
        case "localMachineSam":
            result = LOCALMACHINESAM_REGISTRYHIVE
        case "localMachineSecurity":
            result = LOCALMACHINESECURITY_REGISTRYHIVE
        case "localMachineSoftware":
            result = LOCALMACHINESOFTWARE_REGISTRYHIVE
        case "localMachineSystem":
            result = LOCALMACHINESYSTEM_REGISTRYHIVE
        case "usersDefault":
            result = USERSDEFAULT_REGISTRYHIVE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_REGISTRYHIVE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRegistryHive(values []RegistryHive) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RegistryHive) isMultiValue() bool {
    return false
}

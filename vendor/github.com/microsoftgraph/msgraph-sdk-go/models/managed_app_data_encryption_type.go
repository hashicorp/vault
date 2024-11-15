package models
// Represents the level to which app data is encrypted for managed apps
type ManagedAppDataEncryptionType int

const (
    // App data is encrypted based on the default settings on the device.
    USEDEVICESETTINGS_MANAGEDAPPDATAENCRYPTIONTYPE ManagedAppDataEncryptionType = iota
    // App data is encrypted when the device is restarted.
    AFTERDEVICERESTART_MANAGEDAPPDATAENCRYPTIONTYPE
    // App data associated with this policy is encrypted when the device is locked, except data in files that are open
    WHENDEVICELOCKEDEXCEPTOPENFILES_MANAGEDAPPDATAENCRYPTIONTYPE
    // App data associated with this policy is encrypted when the device is locked
    WHENDEVICELOCKED_MANAGEDAPPDATAENCRYPTIONTYPE
)

func (i ManagedAppDataEncryptionType) String() string {
    return []string{"useDeviceSettings", "afterDeviceRestart", "whenDeviceLockedExceptOpenFiles", "whenDeviceLocked"}[i]
}
func ParseManagedAppDataEncryptionType(v string) (any, error) {
    result := USEDEVICESETTINGS_MANAGEDAPPDATAENCRYPTIONTYPE
    switch v {
        case "useDeviceSettings":
            result = USEDEVICESETTINGS_MANAGEDAPPDATAENCRYPTIONTYPE
        case "afterDeviceRestart":
            result = AFTERDEVICERESTART_MANAGEDAPPDATAENCRYPTIONTYPE
        case "whenDeviceLockedExceptOpenFiles":
            result = WHENDEVICELOCKEDEXCEPTOPENFILES_MANAGEDAPPDATAENCRYPTIONTYPE
        case "whenDeviceLocked":
            result = WHENDEVICELOCKED_MANAGEDAPPDATAENCRYPTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedAppDataEncryptionType(values []ManagedAppDataEncryptionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedAppDataEncryptionType) isMultiValue() bool {
    return false
}

package models
type IosUpdatesInstallStatus int

const (
    DEVICEOSHIGHERTHANDESIREDOSVERSION_IOSUPDATESINSTALLSTATUS IosUpdatesInstallStatus = iota
    SHAREDDEVICEUSERLOGGEDINERROR_IOSUPDATESINSTALLSTATUS
    NOTSUPPORTEDOPERATION_IOSUPDATESINSTALLSTATUS
    INSTALLFAILED_IOSUPDATESINSTALLSTATUS
    INSTALLPHONECALLINPROGRESS_IOSUPDATESINSTALLSTATUS
    INSTALLINSUFFICIENTPOWER_IOSUPDATESINSTALLSTATUS
    INSTALLINSUFFICIENTSPACE_IOSUPDATESINSTALLSTATUS
    INSTALLING_IOSUPDATESINSTALLSTATUS
    DOWNLOADINSUFFICIENTNETWORK_IOSUPDATESINSTALLSTATUS
    DOWNLOADINSUFFICIENTPOWER_IOSUPDATESINSTALLSTATUS
    DOWNLOADINSUFFICIENTSPACE_IOSUPDATESINSTALLSTATUS
    DOWNLOADREQUIRESCOMPUTER_IOSUPDATESINSTALLSTATUS
    DOWNLOADFAILED_IOSUPDATESINSTALLSTATUS
    DOWNLOADING_IOSUPDATESINSTALLSTATUS
    SUCCESS_IOSUPDATESINSTALLSTATUS
    AVAILABLE_IOSUPDATESINSTALLSTATUS
    IDLE_IOSUPDATESINSTALLSTATUS
    UNKNOWN_IOSUPDATESINSTALLSTATUS
)

func (i IosUpdatesInstallStatus) String() string {
    return []string{"deviceOsHigherThanDesiredOsVersion", "sharedDeviceUserLoggedInError", "notSupportedOperation", "installFailed", "installPhoneCallInProgress", "installInsufficientPower", "installInsufficientSpace", "installing", "downloadInsufficientNetwork", "downloadInsufficientPower", "downloadInsufficientSpace", "downloadRequiresComputer", "downloadFailed", "downloading", "success", "available", "idle", "unknown"}[i]
}
func ParseIosUpdatesInstallStatus(v string) (any, error) {
    result := DEVICEOSHIGHERTHANDESIREDOSVERSION_IOSUPDATESINSTALLSTATUS
    switch v {
        case "deviceOsHigherThanDesiredOsVersion":
            result = DEVICEOSHIGHERTHANDESIREDOSVERSION_IOSUPDATESINSTALLSTATUS
        case "sharedDeviceUserLoggedInError":
            result = SHAREDDEVICEUSERLOGGEDINERROR_IOSUPDATESINSTALLSTATUS
        case "notSupportedOperation":
            result = NOTSUPPORTEDOPERATION_IOSUPDATESINSTALLSTATUS
        case "installFailed":
            result = INSTALLFAILED_IOSUPDATESINSTALLSTATUS
        case "installPhoneCallInProgress":
            result = INSTALLPHONECALLINPROGRESS_IOSUPDATESINSTALLSTATUS
        case "installInsufficientPower":
            result = INSTALLINSUFFICIENTPOWER_IOSUPDATESINSTALLSTATUS
        case "installInsufficientSpace":
            result = INSTALLINSUFFICIENTSPACE_IOSUPDATESINSTALLSTATUS
        case "installing":
            result = INSTALLING_IOSUPDATESINSTALLSTATUS
        case "downloadInsufficientNetwork":
            result = DOWNLOADINSUFFICIENTNETWORK_IOSUPDATESINSTALLSTATUS
        case "downloadInsufficientPower":
            result = DOWNLOADINSUFFICIENTPOWER_IOSUPDATESINSTALLSTATUS
        case "downloadInsufficientSpace":
            result = DOWNLOADINSUFFICIENTSPACE_IOSUPDATESINSTALLSTATUS
        case "downloadRequiresComputer":
            result = DOWNLOADREQUIRESCOMPUTER_IOSUPDATESINSTALLSTATUS
        case "downloadFailed":
            result = DOWNLOADFAILED_IOSUPDATESINSTALLSTATUS
        case "downloading":
            result = DOWNLOADING_IOSUPDATESINSTALLSTATUS
        case "success":
            result = SUCCESS_IOSUPDATESINSTALLSTATUS
        case "available":
            result = AVAILABLE_IOSUPDATESINSTALLSTATUS
        case "idle":
            result = IDLE_IOSUPDATESINSTALLSTATUS
        case "unknown":
            result = UNKNOWN_IOSUPDATESINSTALLSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIosUpdatesInstallStatus(values []IosUpdatesInstallStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IosUpdatesInstallStatus) isMultiValue() bool {
    return false
}

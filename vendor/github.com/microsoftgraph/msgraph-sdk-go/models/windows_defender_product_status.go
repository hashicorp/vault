package models
import (
    "math"
    "strings"
)
// Product Status of Windows Defender
type WindowsDefenderProductStatus int

const (
    // No status
    NOSTATUS_WINDOWSDEFENDERPRODUCTSTATUS = 1
    // Service not running
    SERVICENOTRUNNING_WINDOWSDEFENDERPRODUCTSTATUS = 2
    // Service started without any malware protection engine
    SERVICESTARTEDWITHOUTMALWAREPROTECTION_WINDOWSDEFENDERPRODUCTSTATUS = 4
    // Pending full scan due to threat action
    PENDINGFULLSCANDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS = 8
    // Pending reboot due to threat action
    PENDINGREBOOTDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS = 16
    // Pending manual steps due to threat action 
    PENDINGMANUALSTEPSDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS = 32
    // AV signatures out of date
    AVSIGNATURESOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS = 64
    // AS signatures out of date
    ASSIGNATURESOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS = 128
    // No quick scan has happened for a specified period
    NOQUICKSCANHAPPENEDFORSPECIFIEDPERIOD_WINDOWSDEFENDERPRODUCTSTATUS = 256
    // No full scan has happened for a specified period
    NOFULLSCANHAPPENEDFORSPECIFIEDPERIOD_WINDOWSDEFENDERPRODUCTSTATUS = 512
    // System initiated scan in progress
    SYSTEMINITIATEDSCANINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS = 1024
    // System initiated clean in progress
    SYSTEMINITIATEDCLEANINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS = 2048
    // There are samples pending submission
    SAMPLESPENDINGSUBMISSION_WINDOWSDEFENDERPRODUCTSTATUS = 4096
    // Product running in evaluation mode
    PRODUCTRUNNINGINEVALUATIONMODE_WINDOWSDEFENDERPRODUCTSTATUS = 8192
    // Product running in non-genuine Windows mode
    PRODUCTRUNNINGINNONGENUINEMODE_WINDOWSDEFENDERPRODUCTSTATUS = 16384
    // Product expired
    PRODUCTEXPIRED_WINDOWSDEFENDERPRODUCTSTATUS = 32768
    // Off-line scan required
    OFFLINESCANREQUIRED_WINDOWSDEFENDERPRODUCTSTATUS = 65536
    // Service is shutting down as part of system shutdown
    SERVICESHUTDOWNASPARTOFSYSTEMSHUTDOWN_WINDOWSDEFENDERPRODUCTSTATUS = 131072
    // Threat remediation failed critically
    THREATREMEDIATIONFAILEDCRITICALLY_WINDOWSDEFENDERPRODUCTSTATUS = 262144
    // Threat remediation failed non-critically
    THREATREMEDIATIONFAILEDNONCRITICALLY_WINDOWSDEFENDERPRODUCTSTATUS = 524288
    // No status flags set (well initialized state)
    NOSTATUSFLAGSSET_WINDOWSDEFENDERPRODUCTSTATUS = 1048576
    // Platform is out of date
    PLATFORMOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS = 2097152
    // Platform update is in progress
    PLATFORMUPDATEINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS = 4194304
    // Platform is about to be outdated
    PLATFORMABOUTTOBEOUTDATED_WINDOWSDEFENDERPRODUCTSTATUS = 8388608
    // Signature or platform end of life is past or is impending
    SIGNATUREORPLATFORMENDOFLIFEISPASTORISIMPENDING_WINDOWSDEFENDERPRODUCTSTATUS = 16777216
    // Windows SMode signatures still in use on non-Win10S install
    WINDOWSSMODESIGNATURESINUSEONNONWIN10SINSTALL_WINDOWSDEFENDERPRODUCTSTATUS = 33554432
)

func (i WindowsDefenderProductStatus) String() string {
    var values []string
    options := []string{"noStatus", "serviceNotRunning", "serviceStartedWithoutMalwareProtection", "pendingFullScanDueToThreatAction", "pendingRebootDueToThreatAction", "pendingManualStepsDueToThreatAction", "avSignaturesOutOfDate", "asSignaturesOutOfDate", "noQuickScanHappenedForSpecifiedPeriod", "noFullScanHappenedForSpecifiedPeriod", "systemInitiatedScanInProgress", "systemInitiatedCleanInProgress", "samplesPendingSubmission", "productRunningInEvaluationMode", "productRunningInNonGenuineMode", "productExpired", "offlineScanRequired", "serviceShutdownAsPartOfSystemShutdown", "threatRemediationFailedCritically", "threatRemediationFailedNonCritically", "noStatusFlagsSet", "platformOutOfDate", "platformUpdateInProgress", "platformAboutToBeOutdated", "signatureOrPlatformEndOfLifeIsPastOrIsImpending", "windowsSModeSignaturesInUseOnNonWin10SInstall"}
    for p := 0; p < 26; p++ {
        mantis := WindowsDefenderProductStatus(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsDefenderProductStatus(v string) (any, error) {
    var result WindowsDefenderProductStatus
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "noStatus":
                result |= NOSTATUS_WINDOWSDEFENDERPRODUCTSTATUS
            case "serviceNotRunning":
                result |= SERVICENOTRUNNING_WINDOWSDEFENDERPRODUCTSTATUS
            case "serviceStartedWithoutMalwareProtection":
                result |= SERVICESTARTEDWITHOUTMALWAREPROTECTION_WINDOWSDEFENDERPRODUCTSTATUS
            case "pendingFullScanDueToThreatAction":
                result |= PENDINGFULLSCANDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS
            case "pendingRebootDueToThreatAction":
                result |= PENDINGREBOOTDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS
            case "pendingManualStepsDueToThreatAction":
                result |= PENDINGMANUALSTEPSDUETOTHREATACTION_WINDOWSDEFENDERPRODUCTSTATUS
            case "avSignaturesOutOfDate":
                result |= AVSIGNATURESOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS
            case "asSignaturesOutOfDate":
                result |= ASSIGNATURESOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS
            case "noQuickScanHappenedForSpecifiedPeriod":
                result |= NOQUICKSCANHAPPENEDFORSPECIFIEDPERIOD_WINDOWSDEFENDERPRODUCTSTATUS
            case "noFullScanHappenedForSpecifiedPeriod":
                result |= NOFULLSCANHAPPENEDFORSPECIFIEDPERIOD_WINDOWSDEFENDERPRODUCTSTATUS
            case "systemInitiatedScanInProgress":
                result |= SYSTEMINITIATEDSCANINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS
            case "systemInitiatedCleanInProgress":
                result |= SYSTEMINITIATEDCLEANINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS
            case "samplesPendingSubmission":
                result |= SAMPLESPENDINGSUBMISSION_WINDOWSDEFENDERPRODUCTSTATUS
            case "productRunningInEvaluationMode":
                result |= PRODUCTRUNNINGINEVALUATIONMODE_WINDOWSDEFENDERPRODUCTSTATUS
            case "productRunningInNonGenuineMode":
                result |= PRODUCTRUNNINGINNONGENUINEMODE_WINDOWSDEFENDERPRODUCTSTATUS
            case "productExpired":
                result |= PRODUCTEXPIRED_WINDOWSDEFENDERPRODUCTSTATUS
            case "offlineScanRequired":
                result |= OFFLINESCANREQUIRED_WINDOWSDEFENDERPRODUCTSTATUS
            case "serviceShutdownAsPartOfSystemShutdown":
                result |= SERVICESHUTDOWNASPARTOFSYSTEMSHUTDOWN_WINDOWSDEFENDERPRODUCTSTATUS
            case "threatRemediationFailedCritically":
                result |= THREATREMEDIATIONFAILEDCRITICALLY_WINDOWSDEFENDERPRODUCTSTATUS
            case "threatRemediationFailedNonCritically":
                result |= THREATREMEDIATIONFAILEDNONCRITICALLY_WINDOWSDEFENDERPRODUCTSTATUS
            case "noStatusFlagsSet":
                result |= NOSTATUSFLAGSSET_WINDOWSDEFENDERPRODUCTSTATUS
            case "platformOutOfDate":
                result |= PLATFORMOUTOFDATE_WINDOWSDEFENDERPRODUCTSTATUS
            case "platformUpdateInProgress":
                result |= PLATFORMUPDATEINPROGRESS_WINDOWSDEFENDERPRODUCTSTATUS
            case "platformAboutToBeOutdated":
                result |= PLATFORMABOUTTOBEOUTDATED_WINDOWSDEFENDERPRODUCTSTATUS
            case "signatureOrPlatformEndOfLifeIsPastOrIsImpending":
                result |= SIGNATUREORPLATFORMENDOFLIFEISPASTORISIMPENDING_WINDOWSDEFENDERPRODUCTSTATUS
            case "windowsSModeSignaturesInUseOnNonWin10SInstall":
                result |= WINDOWSSMODESIGNATURESINUSEONNONWIN10SINSTALL_WINDOWSDEFENDERPRODUCTSTATUS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsDefenderProductStatus(values []WindowsDefenderProductStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsDefenderProductStatus) isMultiValue() bool {
    return true
}
